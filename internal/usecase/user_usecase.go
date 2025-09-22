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
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/matthewhartstonge/argon2"
)

type UserUsecase interface {
	LoginWithPassword(c context.Context, username string, password string) (resp.LoginResponse, error)
	RegisterWithPassword(c context.Context, u request.Register) (any, error)
	Refresh(c context.Context, refresh string) (any, error)
}

type UserUsecaseImpl struct {
	db    *pg.Queries
	pg    *pkg.Postgres
	cfg   *config.Config
	cache *pkg.RistrettoCache
}

func NewUserUsecase(postgre *pkg.Postgres, cfg *config.Config, cache *pkg.RistrettoCache) *UserUsecaseImpl {
	return &UserUsecaseImpl{
		db:    pg.New(postgre.Pool),
		pg:    postgre,
		cfg:   cfg,
		cache: cache,
	}
}

func (uu *UserUsecaseImpl) LoginWithPassword(c context.Context, username string, password string) (resp.LoginResponse, error) {

	var err error

	///TODO - Check Email & Password
	res, err := uu.db.UsersFindByUsername(c, username)
	if err != nil {
		return resp.LoginResponse{}, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed login")
	}
	if res.Password == "" {
		return resp.LoginResponse{}, pkg.WrapErrorf(err, pkg.ErrorCodeNotFound, "password empty")
	}
	storeHash := []byte(res.Password)
	bytePassword := []byte(password)
	ok, err := argon2.VerifyEncoded(bytePassword, storeHash)
	if err != nil {
		return resp.LoginResponse{}, pkg.WrapErrorf(err, pkg.ErrorCodeNotFound, "verify password error")
	}
	if !ok {
		return resp.LoginResponse{}, pkg.WrapErrorf(err, pkg.ErrorCodeNotFound, "password invalid")
	}

	session_id := pkg.NewUlid()

	user := entity.User{
		ID:       res.ID.String(),
		Username: res.Username,
		Nama:     res.Nama,
		Role:     res.Role,
		Session:  session_id,
	}
	//Generated Access Token
	accessToken, accessExp, err := pkg.CreateAccessToken(
		user,
		uu.cfg.JWT.AccessTokenSecret,
		uu.cfg.JWT.AccessTokenExpireHour)
	if err != nil {
		// uc.Log.Error(logging.JWT, logging.GenerateToken, err.Error(), nil)
		return resp.LoginResponse{}, err
	}
	//remove session from refresh token
	user.Session = ""
	//Generated Refresh Token
	refreshToken, refreshExp, err := pkg.CreateRefreshToken(
		user,
		uu.cfg.JWT.RefreshTokenSecret,
		uu.cfg.JWT.RefreshTokenExpireHour)
	if err != nil {
		// uc.Log.Error(logging.JWT, logging.GenerateToken, err.Error(), nil)
		return resp.LoginResponse{}, err
	}

	expire := time.Duration(uu.cfg.JWT.AccessTokenExpireHour)*time.Minute - 1*time.Minute

	uu.cache.SetWithTTL(session_id, session_id, expire)

	arg := pg.UpdateUserPartialParams{
		ID:      res.ID,
		Refresh: &refreshToken,
	}
	uu.db.UpdateUserPartial(c, arg)

	return resp.LoginResponse{
		User: resp.User{
			ID:          res.ID.String(),
			Username:    res.Username,
			Nama:        res.Nama,
			Role:        res.Role,
			AccessToken: accessToken,
			Exp:         accessExp,
		},
		RefreshToken: resp.RefreshToken{Token: refreshToken, Exp: refreshExp},
	}, nil

}

func (uu *UserUsecaseImpl) RegisterWithPassword(c context.Context, u request.Register) (any, error) {
	return utils.WithTransactionResult(c, uu.pg.Pool, func(qtx *pg.Queries, tx pgx.Tx) (any, error) {
		var err error
		///TODO -  (maps recipe user ID to primary user ID)
		// userid := pkg.NewUlid()
		user := pg.CreateOrUpdateUserParams{
			Username: u.Username,
			Nama:     u.Nama,
			Password: u.Password,
			Role:     uuid.Must(uuid.FromString(u.Role)),
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

		return res, nil
	})
}

func (uu *UserUsecaseImpl) Refresh(c context.Context, refresh string) (any, error) {

	isAuthorize, err := pkg.IsAuthorized(refresh, uu.cfg.JWT.RefreshTokenSecret)
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

	session_id := pkg.NewUlid()

	user := entity.User{
		ID:       res.ID.String(),
		Username: res.Username,
		Nama:     res.Nama,
		Role:     res.Role,
		Session:  session_id,
	}

	access, exp, err := pkg.CreateAccessToken(user,
		uu.cfg.JWT.AccessTokenSecret,
		uu.cfg.JWT.AccessTokenExpireHour)
	if err != nil {
		return nil, err
	}

	return resp.User{
		ID:          res.ID.String(),
		Username:    res.Username,
		Nama:        res.Nama,
		Role:        res.Role,
		AccessToken: access,
		Exp:         exp,
	}, nil
}
