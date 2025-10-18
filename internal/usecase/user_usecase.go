package usecase

import (
	"context"
	"e-klinik/config"
	"e-klinik/infra/pg"
	"e-klinik/internal/domain/entity"
	"e-klinik/internal/domain/request"
	"e-klinik/internal/domain/resp"
	"e-klinik/pkg"
	"e-klinik/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jinzhu/copier"
	"github.com/matthewhartstonge/argon2"
	"github.com/redis/go-redis/v9"
)

type UserUsecase interface {
	LoginWithPassword(c context.Context, username string, password string) (resp.User, error)
	Logout(c context.Context, id string) (any, error)
	RegisterWithPassword(c context.Context, arg request.Register) (any, error)
	Refresh(c context.Context, refresh string) (any, error)
	AddRoleForUser(c context.Context, u pg.CreateUserRoleParams) (any, error)
	AddMenu(c context.Context, arg pg.CreateR1ViewParams) (any, error)
	ListMenu(c context.Context, arg request.SearchMenu) (any, error)
	EditMenu(c context.Context, arg pg.UpdateR1ViewParams) (any, error)
	DeleteMenu(c context.Context, arg pg.DeleteR1ViewParams) error
	MenuById(c context.Context, arg int32) (any, error)
	AccessList(c context.Context, arg request.SearchMenu) (any, error)
	AddRole(c context.Context, arg pg.CreateR4RoleParams) (any, error)
	ListRole(c context.Context) (any, error)
	UpdateRole(c context.Context, arg pg.UpdateR4RoleParams) (any, error)
	GetRoleById(c context.Context, arg int32) (any, error)
	DeleteRoleById(c context.Context, arg pg.DeleteR4RoleParams) error
	AddGroup(c context.Context, arg pg.CreateR2GroupParams) (any, error)
	ListGroup(c context.Context) (any, error)
	UpdateGroup(c context.Context, arg pg.UpdateR2GroupParams) (any, error)
	GetGroupById(c context.Context, arg int32) (any, error)
	DeleteGroupById(c context.Context, arg pg.DeleteR2GroupParams) error
	ListUser(c context.Context, arg request.SearchUser) (any, error)
	UpdateUserPartial(c context.Context, arg request.UpdateUser) (any, error)
	GetUserId(c context.Context, arg uuid.UUID) (any, error)
	GetUserRoleByUserId(c context.Context, arg uuid.UUID) (any, error)
	GetViewByRoleId(c context.Context, arg int32) (any, error)
	AddRolePolicy(c context.Context, arg request.UpdateRolePolicy) (any, error)
	GetViewUser(c context.Context, arg uuid.UUID) (any, error)
}

type UserUsecaseImpl struct {
	db    *pg.Queries
	pg    *pkg.Postgres
	cfg   *config.Config
	cache *pkg.RedisCache
	cbn   *casbin.Enforcer
}

func NewUserUsecase(postgre *pkg.Postgres, cfg *config.Config, cache *pkg.RedisCache, cbn *casbin.Enforcer) *UserUsecaseImpl {
	return &UserUsecaseImpl{
		db:    pg.New(postgre.Pool),
		pg:    postgre,
		cfg:   cfg,
		cache: cache,
		cbn:   cbn,
	}
}

func (uu *UserUsecaseImpl) LoginWithPassword(c context.Context, username string, password string) (resp.User, error) {

	var err error

	///TODO - Check Email & Password
	res, err := uu.db.UsersFindByUsername(c, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Tidak ada user → return empty response, error nil
			return resp.User{}, pkg.WrapErrorf(err, pkg.ErrorCodeNotFound, "user not found")
		}
		return resp.User{}, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed login")
	}

	if res.Password == "" {
		return resp.User{}, pkg.WrapErrorf(err, pkg.ErrorCodeNotFound, "password empty")
	}
	storeHash := []byte(res.Password)
	bytePassword := []byte(password)
	ok, err := argon2.VerifyEncoded(bytePassword, storeHash)
	if err != nil {
		return resp.User{}, pkg.WrapErrorf(err, pkg.ErrorCodeNotFound, "verify password error")
	}
	if !ok {
		return resp.User{}, pkg.WrapErrorf(err, pkg.ErrorCodeNotFound, "password invalid")
	}

	sessionId := pkg.NewUlid()

	user := entity.User{
		ID:       res.ID.String(),
		Username: res.Username,
		Nama:     res.Nama,
		Role:     res.Role,
		Session:  sessionId,
	}
	//Generated Access Token
	accessToken, accessExp, err := pkg.CreateAccessToken(
		user,
		uu.cfg.JWT.AccessTokenSecret,
		uu.cfg.JWT.AccessTokenExpireHour)
	if err != nil {
		// uc.Log.Error(logging.JWT, logging.GenerateToken, err.Error(), nil)
		return resp.User{}, err
	}
	//remove session from refresh token
	user.Session = ""
	//Generated Refresh Token
	refreshToken, _, err := pkg.CreateRefreshToken(
		user,
		uu.cfg.JWT.RefreshTokenSecret,
		uu.cfg.JWT.RefreshTokenExpireHour)
	if err != nil {
		// uc.Log.Error(logging.JWT, logging.GenerateToken, err.Error(), nil)
		return resp.User{}, err
	}

	view, err := uu.db.GetUserMenuViews(c, res.ID)
	if err != nil {
		return resp.User{}, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get menu")
	}

	expire := time.Duration(uu.cfg.JWT.AccessTokenExpireHour)*time.Minute - 1*time.Minute

	redisKey := fmt.Sprintf("view:%d", res.ID)

	uu.cache.SetWithTTL(c, redisKey, view, time.Duration(uu.cfg.JWT.AccessTokenExpireHour)*time.Minute)

	uu.cache.SetWithTTL(c, sessionId, sessionId, expire)

	arg := pg.UpdateUserPartialParams{
		ID:      res.ID,
		Refresh: &refreshToken,
	}
	err = uu.db.UpdateUserPartial(c, arg)

	return resp.User{
			ID:                 res.ID.String(),
			Username:           res.Username,
			Nama:               res.Nama,
			Role:               res.Role,
			AccessToken:        accessToken,
			RefreshToken:       refreshToken,
			AccessTokenExpires: accessExp,
		},
		nil

}

func (uu *UserUsecaseImpl) RegisterWithPassword(c context.Context, arg request.Register) (any, error) {
	return utils.WithTransactionResult(c, uu.pg.Pool, func(qtx *pg.Queries, tx pgx.Tx) (any, error) {
		var err error
		///TODO -  (maps recipe user ID to primary user ID)
		// userid := pkg.NewUlid()
		user := pg.CreateOrUpdateUserParams{
			Username: arg.Username,
			Nama:     arg.Nama,
			Password: arg.Password,
			// 	Role:     uuid.Must(uuid.FromString(u.Role)
			// ),
		}
		config := argon2.DefaultConfig()
		encoded, err := config.HashEncoded([]byte("123456678"))
		if err != nil {
			panic(err) // 💥
		}
		user.Password = string(encoded)
		res, err := qtx.CreateOrUpdateUser(c, user)
		if err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed register")
		}

		// 5️⃣ Tambahkan role baru di SQL + Casbin
		for _, roleID := range arg.Role {
			carg := pg.CreateUserRoleParams{
				UserID:    res.ID,
				RoleID:    utils.StrToInt32(roleID.Value),
				CreatedBy: arg.CreatedBy,
			}

			if err := qtx.CreateUserRole(c, carg); err != nil {
				return nil, fmt.Errorf("failed add role %s for user %s: %w", roleID.Value, res.ID, err)
			}

			if _, err := uu.cbn.AddRoleForUser(res.ID.String(), roleID.Value, "", "", "", ""); err != nil {
				return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed add casbin role")
			}
		}

		return res, nil
	})
}

func (uu *UserUsecaseImpl) Refresh(c context.Context, refresh string) (any, error) {

	isAuthorize, _, err := pkg.IsAuthorized(refresh, uu.cfg.JWT.RefreshTokenSecret)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "checking auth failed")
	}

	if !isAuthorize {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "UnAuthorize")
	}
	id, err := pkg.ExtractIDFromToken(refresh, uu.cfg.JWT.RefreshTokenSecret)
	// 3. Convert string -> UUID safely
	// userUUID, err := uuid.FromString(id)
	// if err != nil {
	// 	return nil, pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid UUID in token")
	// }
	res, err := uu.db.UsersFindById(c, uuid.Must(uuid.FromString(id)))
	if err != nil {
		return nil, err
	}
	if res.Refresh == nil {
		return nil, pkg.WrapErrorf(nil, pkg.ErrorCodeInvalidArgument, "user has no refresh token stored")
	}
	//Return false when refresh token doesn't match
	if *res.Refresh != refresh {
		return nil, errors.New("token doesn't match")
	}

	sessionId := pkg.NewUlid()

	user := entity.User{
		ID:       res.ID.String(),
		Username: res.Username,
		Nama:     res.Nama,
		Role:     res.Role,
		Session:  sessionId,
	}

	access, exp, err := pkg.CreateAccessToken(user,
		uu.cfg.JWT.AccessTokenSecret,
		uu.cfg.JWT.AccessTokenExpireHour)
	if err != nil {
		return nil, err
	}

	view, err := uu.db.GetUserMenuViews(c, res.ID)
	if err != nil {
		return resp.User{}, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get menu")
	}

	redisKey := fmt.Sprintf("view:%d", res.ID)

	uu.cache.SetWithTTL(c, redisKey, view, time.Duration(uu.cfg.JWT.AccessTokenExpireHour)*time.Minute)

	return resp.User{
		ID:                 res.ID.String(),
		Username:           res.Username,
		Nama:               res.Nama,
		Role:               res.Role,
		AccessToken:        access,
		RefreshToken:       utils.DerefString(res.Refresh),
		AccessTokenExpires: exp,
	}, nil
}

func (uu *UserUsecaseImpl) Logout(c context.Context, id string) (any, error) {

	var err error
	err = uu.db.UpdateUserPartial(c, pg.UpdateUserPartialParams{ID: uuid.Must(uuid.FromString(id)), Refresh: utils.StringPtr(""), UpdatedNote: utils.StringPtr("logout")})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (uu *UserUsecaseImpl) AddRoleForUser(c context.Context, u pg.CreateUserRoleParams) (any, error) {
	return utils.WithTransactionResult(c, uu.pg.Pool, func(qtx *pg.Queries, tx pgx.Tx) (any, error) {
		var err error

		_, err = uu.cbn.AddRoleForUser(u.UserID.String(), utils.Int32ToStr(u.RoleID), "", "", "", "")
		if err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed add role")
		}

		// res, err := qtx.CreateUserRole(c, u)
		// if err != nil {
		// 	return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed register")
		// }

		return nil, nil
	})
}

func (uu *UserUsecaseImpl) AddGroupRole(c context.Context, u pg.CreateGroupRoleParams) (any, error) {
	return utils.WithTransactionResult(c, uu.pg.Pool, func(qtx *pg.Queries, tx pgx.Tx) (any, error) {
		var err error

		_, err = uu.cbn.AddNamedGroupingPolicy("g2", utils.Int32ToStr(u.GroupID), utils.Int32ToStr(u.RoleID), "", "", "")
		if err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed add group role")
		}

		res, err := qtx.CreateGroupRole(c, u)
		if err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed group role")
		}

		return res, nil
	})
}

func (uu *UserUsecaseImpl) AddMenu(c context.Context, arg pg.CreateR1ViewParams) (any, error) {
	return utils.WithTransactionResult(c, uu.pg.Pool, func(qtx *pg.Queries, tx pgx.Tx) (any, error) {
		// var err error

		res, err := qtx.CreateR1View(c, arg)
		if err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed create menu")
		}

		return res, nil
	})
}

func (uu *UserUsecaseImpl) ListMenu(c context.Context, arg request.SearchMenu) (any, error) {

	res, err := uu.db.ListR1Views(c, arg.Label)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed list menu")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}
	return resp.WithPaginate(res, nil), nil
}

func (uu *UserUsecaseImpl) EditMenu(c context.Context, arg pg.UpdateR1ViewParams) (any, error) {

	res, err := uu.db.UpdateR1View(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed edit menu")
	}
	return res, nil
}

func (uu *UserUsecaseImpl) DeleteMenu(c context.Context, arg pg.DeleteR1ViewParams) error {

	err := uu.db.DeleteR1View(c, arg)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed delete menu")
	}
	return nil
}

func (uu *UserUsecaseImpl) MenuById(c context.Context, arg int32) (any, error) {
	res, err := uu.db.GetR1View(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get")
	}
	return resp.WithPaginate(res, nil), nil
}

func (uu *UserUsecaseImpl) AccessList(c context.Context, arg request.SearchMenu) (any, error) {

	res, err := uu.db.GetR1ViewRecursive(c)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed list")
	}
	tree := resp.BuildMenuTree(res, nil)
	return resp.WithPaginate(tree, nil), nil
}

func (uu *UserUsecaseImpl) AddRole(c context.Context, arg pg.CreateR4RoleParams) (any, error) {
	return utils.WithTransactionResult(c, uu.pg.Pool, func(qtx *pg.Queries, tx pgx.Tx) (any, error) {
		// var err error
		res, err := qtx.CreateR4Role(c, arg)
		if err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed create")
		}

		return res, nil
	})
}

func (uu *UserUsecaseImpl) ListRole(c context.Context) (any, error) {

	// var err error
	res, err := uu.db.ListR4Roles(c)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed list")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}

	return resp.WithPaginate(res, nil), nil

}

func (uu *UserUsecaseImpl) UpdateRole(c context.Context, arg pg.UpdateR4RoleParams) (any, error) {

	// var err error
	res, err := uu.db.UpdateR4Role(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed update")
	}

	return res, nil

}

func (uu *UserUsecaseImpl) GetRoleById(c context.Context, arg int32) (any, error) {

	// var err error
	res, err := uu.db.GetR4RoleByID(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get data")
	}

	return resp.WithPaginate(res, nil), nil

}

func (uu *UserUsecaseImpl) DeleteRoleById(c context.Context, arg pg.DeleteR4RoleParams) error {
	// var err error
	err := uu.db.DeleteR4Role(c, arg)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get data")
	}

	return nil

}

func (uu *UserUsecaseImpl) AddGroup(c context.Context, arg pg.CreateR2GroupParams) (any, error) {
	return utils.WithTransactionResult(c, uu.pg.Pool, func(qtx *pg.Queries, tx pgx.Tx) (any, error) {
		// var err error
		res, err := qtx.CreateR2Group(c, arg)
		if err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed create")
		}

		return res, nil
	})
}

func (uu *UserUsecaseImpl) ListGroup(c context.Context) (any, error) {

	// var err error
	res, err := uu.db.ListR2Groups(c)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed list")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}

	return resp.WithPaginate(res, nil), nil

}

func (uu *UserUsecaseImpl) UpdateGroup(c context.Context, arg pg.UpdateR2GroupParams) (any, error) {

	// var err error
	res, err := uu.db.UpdateR2Group(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed update")
	}

	return res, nil

}

func (uu *UserUsecaseImpl) GetGroupById(c context.Context, arg int32) (any, error) {

	// var err error
	res, err := uu.db.GetR2GroupByID(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get data")
	}

	return resp.WithPaginate(res, nil), nil

}

func (uu *UserUsecaseImpl) DeleteGroupById(c context.Context, arg pg.DeleteR2GroupParams) error {
	// var err error
	err := uu.db.DeleteR2Group(c, arg)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get data")
	}

	return nil

}

func (uu *UserUsecaseImpl) ListUser(c context.Context, arg request.SearchUser) (any, error) {

	if arg.Limit <= 0 {
		arg.Limit = 10
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	arg.Offset = utils.GetOffset(arg.Page, arg.Limit)

	var params pg.ListUsersParams

	err := copier.Copy(&params, &arg)
	if err != nil {
		return nil, err
	}
	res, err := uu.db.ListUsers(c, params)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed list")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}
	var cparams pg.CountUsersParams
	err = copier.Copy(&cparams, &arg)
	if err != nil {
		return nil, err
	}
	count, err := uu.db.CountUsers(c, cparams)
	if err != nil {
		return nil, err
	}

	return resp.WithPaginate(res, resp.CalculatePagination(arg.Page, arg.Limit, count)), nil

}

func (uu *UserUsecaseImpl) UpdateUserPartial(c context.Context, arg request.UpdateUser) (any, error) {
	return utils.WithTransactionResult(c, uu.pg.Pool, func(qtx *pg.Queries, tx pgx.Tx) (any, error) {
		// 1️⃣ Update user partial fields
		var params pg.UpdateUserPartialParams
		if err := copier.Copy(&params, &arg); err != nil {
			return nil, err
		}

		if err := qtx.UpdateUserPartial(c, params); err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed update user")
		}

		// 2️⃣ Persiapkan list role dari input
		valuesStr := make([]string, len(arg.Role))
		valuesInt := make([]int32, len(arg.Role))
		for i, opt := range arg.Role {
			valuesStr[i] = opt.Value
			valuesInt[i] = utils.StrToInt32(opt.Value)
		}

		// 3️⃣ Hapus role lama di SQL yang tidak ada di input
		darg := pg.DeleteUnRegisterRoleParams{
			UserID:  arg.ID,
			RoleIds: valuesInt, // jika kosong → hapus semua
		}
		if err := qtx.DeleteUnRegisterRole(c, darg); err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed delete role in SQL")
		}

		// 4️⃣ Hapus role lama di Casbin yang tidak ada di input
		currentRoles, err := uu.cbn.GetRolesForUser(arg.ID.String())
		if err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get casbin roles")
		}

		newRolesSet := make(map[string]struct{})
		for _, r := range valuesStr {
			newRolesSet[r] = struct{}{}
		}

		for _, r := range currentRoles {
			if _, ok := newRolesSet[r]; !ok {
				if _, err := uu.cbn.DeleteRoleForUser(arg.ID.String(), r, "", "", "", ""); err != nil {
					return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed delete casbin role")
				}
			}
		}
		// ======================
		// ⚙️ Sync CASBIN POLICY
		// ======================

		// 5️⃣ Tambahkan role baru di SQL + Casbin
		for _, roleID := range arg.Role {
			carg := pg.CreateUserRoleParams{
				UserID:    arg.ID,
				RoleID:    utils.StrToInt32(roleID.Value),
				CreatedBy: arg.UpdatedBy,
			}

			if err := qtx.CreateUserRole(c, carg); err != nil {
				return nil, fmt.Errorf("failed add role %s for user %s: %w", roleID.Value, arg.ID, err)
			}

			if _, err := uu.cbn.AddRoleForUser(arg.ID.String(), roleID.Value, "", "", "", ""); err != nil {
				return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed add casbin role")
			}
		}

		return nil, nil
	})
}

func (uu *UserUsecaseImpl) GetUserId(c context.Context, arg uuid.UUID) (any, error) {

	// var err error
	res, err := uu.db.GetUserDetail(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get data")
	}

	return resp.WithPaginate(res, nil), nil

}

func (uu *UserUsecaseImpl) GetUserRoleByUserId(c context.Context, arg uuid.UUID) (any, error) {

	// var err error
	res, err := uu.db.GetUserRolesByUserID(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get data")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}

	return resp.WithPaginate(res, nil), nil

}
func (uu *UserUsecaseImpl) GetViewByRoleId(c context.Context, arg int32) (any, error) {

	// var err error
	res, err := uu.db.GetR3ViewRoleByRoleID(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get data")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}

	return resp.WithPaginate(res, nil), nil

}

func (uu *UserUsecaseImpl) AddRolePolicy(c context.Context, arg request.UpdateRolePolicy) (any, error) {
	return utils.WithTransactionResult(c, uu.pg.Pool, func(qtx *pg.Queries, tx pgx.Tx) (any, error) {

		var err error
		rows, err := qtx.GetViewsByIdsWithChildren(c, arg.Policy)
		if err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get data")
		}

		var (
			allIDs   []int32
			dataRows []pg.GetViewsByIdsWithChildrenRow
		)

		for _, r := range rows {
			// Simpan semua ID
			allIDs = append(allIDs, r.ID)

			// Ambil baris penuh yang punya view = 'data'
			if r.View != nil && *r.View == "data" {
				dataRows = append(dataRows, r)
			}
		}
		darg := pg.ViewRolesSyncDeleteHardParams{
			RoleID:  arg.RoleID,
			ViewIds: allIDs,
		}
		err = qtx.ViewRolesSyncDeleteHard(c, darg)
		if err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed delete data")
		}

		iarg := pg.ViewRolesSyncInsertHardParams{
			RoleID:        arg.RoleID,
			ViewIds:       allIDs,
			CurrentUserID: arg.CreatedBy,
		}
		err = qtx.ViewRolesSyncInsertHard(c, iarg)
		if err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed inser data")
		}
		// ======================
		// ⚙️ 3️⃣ Sync CASBIN POLICY
		// ======================

		roleID := fmt.Sprintf("%d", arg.RoleID)

		// Ambil semua policy lama untuk role ini
		oldPolicies, err := uu.cbn.GetFilteredPolicy(0, roleID)
		if err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed filter policy")
		}

		// Buat policy baru dari hasil query
		var newPolicies [][]string
		for _, r := range dataRows {
			p := []string{
				roleID,        // v0 = role
				r.ResourceKey, // v1 = resource
				r.Action,      // v2 = action
			}
			newPolicies = append(newPolicies, padPolicy(p))
		}

		// Buat map cepat untuk cek perubahan
		newMap := make(map[string]bool)
		for _, p := range newPolicies {
			key := p[1] + "|" + p[2]
			newMap[key] = true
		}

		// 1️⃣ Hapus policy lama yang tidak ada di input baru
		var removeList [][]string
		for _, old := range oldPolicies {
			if len(old) < 3 {
				continue
			}
			key := old[1] + "|" + old[2]
			if !newMap[key] {
				removeList = append(removeList, padPolicy(old))
			}
		}

		if len(removeList) > 0 {
			_, err := uu.cbn.RemovePolicies(removeList)
			if err != nil {
				return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed remove old casbin policies")
			}
		}

		// 2️⃣ Tambahkan policy baru yang belum ada
		var addList [][]string
		for _, np := range newPolicies {
			// hanya kirim jika belum ada
			exist, _ := uu.cbn.HasPolicy(np)
			if !exist {
				addList = append(addList, np)
			}
		}

		if len(addList) > 0 {
			_, err := uu.cbn.AddPolicies(addList)
			if err != nil {
				return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed add new casbin policies")
			}
		}

		// 3️⃣ Simpan ke storage Casbin
		if err := uu.cbn.SavePolicy(); err != nil {
			return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed save casbin policy")
		}
		return resp.WithPaginate(map[string]any{
			"role_id":          arg.RoleID,
			"total_ids":        len(allIDs),
			"data_rows":        len(dataRows),
			"added_policies":   len(addList),
			"removed_policies": len(removeList),
		}, nil), nil
	})
}

func (uu *UserUsecaseImpl) GetViewUser(c context.Context, arg uuid.UUID) (any, error) {

	redisKey := fmt.Sprintf("view:%d", arg)

	// 1️⃣ Cek Redis
	val, err := uu.cache.Client.Get(c, redisKey).Result()
	if err == nil {
		// Redis ada → unmarshal JSON
		var views []pg.GetUserMenuViewsRow
		if err := json.Unmarshal([]byte(val), &views); err != nil {
			log.Println("failed unmarshal cached JSON:", err)
		} else {
			return views, nil
		}
	} else if err != redis.Nil {
		// Error Redis lain
		log.Println("failed get from redis:", err)
	}

	// 2️⃣ Ambil dari DB
	views, err := uu.db.GetUserMenuViews(c, arg)
	if err != nil {
		return nil, err
	}

	expire := time.Duration(uu.cfg.JWT.AccessTokenExpireHour)*time.Hour - 1*time.Minute
	ok := uu.cache.SetWithTTL(c, redisKey, views, expire)
	if !ok {
		log.Println("failed set cache")
	}

	return views, nil

}

// ✅ Helper untuk pastikan policy punya 6 field (v0–v5)
func padPolicy(policy []string) []string {
	for len(policy) < 6 {
		policy = append(policy, "")
	}
	return policy
}
