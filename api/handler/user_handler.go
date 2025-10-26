package handler

import (
	"context"
	"e-klinik/config"
	"e-klinik/infra/pg"
	"e-klinik/pkg"
	"e-klinik/utils"

	"e-klinik/internal/domain/request"
	"e-klinik/internal/domain/resp"
	"e-klinik/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofrs/uuid/v5"
)

type UserHandler interface {
	Register(c *gin.Context)
	Logout(c *gin.Context)
	AddRoleUser(c *gin.Context)
	AddRolePolicy(c *gin.Context)

	CreateRole(c *gin.Context)
	RoleById(c *gin.Context)
	ListRole(c *gin.Context)
	UpdateRole(c *gin.Context)
	DelRole(c *gin.Context)
	CreateGroup(c *gin.Context)
	GroupById(c *gin.Context)
	ListGroup(c *gin.Context)
	UpdateGroup(c *gin.Context)
	DelGroup(c *gin.Context)
	ListUsers(c *gin.Context)
	UpdateUser(c *gin.Context)
	UserById(c *gin.Context)
	UserRoleByUserId(c *gin.Context)

	UpdateRolePolicyByRoleId(c *gin.Context)
	CreateNewUser(c *gin.Context)
	UserViewPermission(c *gin.Context)
}

type UserHandlerImpl struct {
	Uu *usecase.UserUsecaseImpl

	Cfg *config.Config
}

func NewUserHandler(Uu *usecase.UserUsecaseImpl, cfg *config.Config) *UserHandlerImpl {
	return &UserHandlerImpl{
		Uu:  Uu,
		Cfg: cfg,
	}
}

func (lc *UserHandlerImpl) Register(c *gin.Context) {
	var req request.Register

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.HandleErrorResponse(
			c,
			"failed to bind JSON",
			pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"),
		)
		return
	}

	// ✅ Validasi dengan ozzo-validation
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Nama, validation.Required.Error("nama tidak boleh kosong")),
		validation.Field(&req.Username, validation.Required.Error("username tidak boleh kosong")),
	)

	if err != nil {
		// cek jika termasuk validation.Errors
		if verrs, ok := err.(validation.Errors); ok {
			resp.HandleErrorResponse(c, "validation failed", pkg.WrapValidationError(verrs, "invalid input"))
			return
		}
		// jika bukan validation error biasa
		resp.HandleErrorResponse(c, "validation failed", err)
		return
	}

	user, err := lc.Uu.RegisterWithPassword(c, req)
	if err != nil {
		resp.HandleErrorResponse(
			c,
			"failed to register user",
			err,
		)
		return
	}

	resp.HandleSuccessResponse(c, "register berhasil", user)
}

func (lc *UserHandlerImpl) Logout(c *gin.Context) {
	value, exists := c.Get("Id")
	if !exists {
		resp.HandleErrorResponse(
			c,
			"user context not found",
			pkg.ExposeError(pkg.ErrorCodeUnauthorized, "token tidak valid"),
		)
		return
	}

	userID := value.(string)
	user, err := lc.Uu.Logout(c, userID)
	if err != nil {
		resp.HandleErrorResponse(
			c,
			"failed to logout user",
			pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"),
		)
		return
	}

	// Hapus cookie
	c.SetCookie(
		"auth_token",
		"",
		-1,
		"/",
		"localhost",
		false, // ubah ke true di production (HTTPS)
		true,
	)

	resp.HandleSuccessResponse(c, "logout berhasil", user)
}

// func (lc *UserHandler) SearchMU(c *gin.Context) {
// 	targetURL := "https://api.mangaupdates.com/v1/series/search"
// 	search := c.Query("search")
// 	// Make the GET request to the target URL
// 	resp, err := fetch.Post(targetURL, &fetch.Config{Body: map[string]interface{}{"search": search,
// 		"exclude_filtered_genres": true,
// 		"stype":                   "description",
// 		"licensed":                "no",
// 		"page":                    1,
// 		"perpage":                 6,
// 		"include_rank_metadata":   false,
// 		"orderby":                 "title"}})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching the target URL"})
// 		return
// 	}

// 	// Copy the status code from the target response
// 	c.Status(resp.StatusCode())

// 	// Copy headers from the target response
// 	for key, values := range resp.Headers {
// 		for _, value := range values {
// 			c.Writer.Header().Add(key, value)
// 		}
// 	}
// 	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
// 	// Copy the response body from the target to the client
// 	c.Writer.Write(resp.Body)
// }

// func (lc *UserHandler) SearchMAL(c *gin.Context) {
// 	search := c.Query("search")
// 	targetURL := "https://api.jikan.moe/v4/manga"

// 	// Make the GET request to the target URL
// 	resp, err := fetch.Get(targetURL, &fetch.Config{Query: fetch.Query{"q": search, "page": "1", "limit": "6"}})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching the target URL"})
// 		return
// 	}

// 	// Copy the status code from the target response
// 	c.Status(resp.StatusCode())

// 	// Copy headers from the target response
// 	for key, values := range resp.Headers {

// 		for _, value := range values {
// 			c.Writer.Header().Add(key, value)
// 		}
// 	}

// 	// Copy the response body from the target to the client
// 	c.Writer.Write(resp.Body)
// }

// func (lc *UserHandler) MangaDetailMAL(c *gin.Context) {
// 	search := c.Param("id")
// 	targetURL := fmt.Sprintf("https://api.jikan.moe/v4/manga/%s", search)

// 	// Make the GET request to the target URL
// 	resp, err := fetch.Get(targetURL)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching the target URL"})
// 		return
// 	}

// 	// Copy the status code from the target response
// 	c.Status(resp.StatusCode())

// 	// Copy headers from the target response
// 	for key, values := range resp.Headers {

// 		for _, value := range values {
// 			c.Writer.Header().Add(key, value)
// 		}
// 	}

// 	// Copy the response body from the target to the client
// 	c.Writer.Write(resp.Body)
// }

// func (lc *UserHandler) MangaDetailMU(c *gin.Context) {
// 	id := c.Param("id")
// 	targetURL := fmt.Sprintf("https://api.mangaupdates.com/v1/series/%s", id)

// 	// Make the GET request to the target URL
// 	resp, err := fetch.Get(targetURL)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching the target URL"})
// 		return
// 	}

// 	// Copy the status code from the target response
// 	c.Status(resp.StatusCode())

// 	// Copy headers from the target response
// 	for key, values := range resp.Headers {
// 		for _, value := range values {
// 			c.Writer.Header().Add(key, value)
// 		}
// 	}
// 	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
// 	// Copy the response body from the target to the client
// 	c.Writer.Write(resp.Body)
// }

func (lc *UserHandlerImpl) AddRoleUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var p pg.CreateUserRoleParams
	if err := c.ShouldBind(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, exists := c.Get("nama")
	if !exists {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak ditemukan"))
		return
	}
	if v, ok := value.(string); ok {
		p.CreatedBy = utils.StringPtr(v)
	} else {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak valid"))
		return
	}

	res, err := lc.Uu.AddRoleForUser(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to add role user", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success add role", res)
}

func (lc *UserHandlerImpl) AddRolePolicy(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var p pg.CreateUserRoleParams
	if err := c.ShouldBind(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, exists := c.Get("nama")
	if !exists {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak ditemukan"))
		return
	}
	if v, ok := value.(string); ok {
		p.CreatedBy = utils.StringPtr(v)
	} else {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak valid"))
		return
	}

	res, err := lc.Uu.AddRoleForUser(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to add role policy", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success add role policy", res)
}

func (lc *UserHandlerImpl) CreateRole(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var p pg.CreateR4RoleParams
	if err := c.ShouldBind(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak ditemukan"))
		return
	}
	if v, ok := value.(string); ok {
		p.CreatedBy = utils.StringPtr(v)
	} else {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak valid"))
		return
	}

	res, err := lc.Uu.AddRole(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed create role", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success create role", res)
}

func (lc *UserHandlerImpl) RoleById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		resp.HandleErrorResponse(c, "detail failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "parameter id tidak boleh kosong"))
		return
	}

	p := utils.StrToInt32(id)
	if p == 0 {
		resp.HandleErrorResponse(c, "detail failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "id tidak valid"))
		return
	}

	res, err := lc.Uu.GetRoleById(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get role", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}
	resp.HandleSuccessResponse(c, "success get role", res)
}

func (lc *UserHandlerImpl) ListRole(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.SearchRuangan
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	res, err := lc.Uu.ListRole(ctx)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get list role", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success get list role", res)
}

func (lc *UserHandlerImpl) UpdateRole(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	var p pg.UpdateR4RoleParams
	p.ID = utils.StrToInt32(id)

	if err := c.ShouldBind(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak ditemukan"))
		return
	}
	if v, ok := value.(string); ok {
		p.UpdatedBy = utils.StringPtr(v)
	} else {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak valid"))
		return
	}

	res, err := lc.Uu.UpdateRole(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed update role", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}
	resp.HandleSuccessResponse(c, "success update", res)
}

func (lc *UserHandlerImpl) DelRole(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Query("id")
	if id == "" {
		resp.HandleErrorResponse(c, "delete failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid id"))
		return
	}
	p := pg.DeleteR4RoleParams{
		ID:        utils.StrToInt32(id),
		DeletedBy: nil,
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak ditemukan"))
		return
	}
	if v, ok := value.(string); ok {
		p.DeletedBy = utils.StringPtr(v)
	} else {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak valid"))
		return
	}

	if err := lc.Uu.DeleteRoleById(ctx, p); err != nil {
		resp.HandleErrorResponse(c, "failed delete role", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}
	resp.HandleSuccessResponse(c, "success delete role", nil)
}

func (lc *UserHandlerImpl) CreateGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var p pg.CreateR2GroupParams
	if err := c.ShouldBind(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}
	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak ditemukan"))
		return
	}
	if v, ok := value.(string); ok {
		p.CreatedBy = utils.StringPtr(v)
	} else {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak valid"))
		return
	}

	res, err := lc.Uu.AddGroup(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed create group", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success create group", res)
}

func (lc *UserHandlerImpl) GroupById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		resp.HandleErrorResponse(c, "detail failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid id"))
		return
	}
	p := utils.StrToInt32(id)

	res, err := lc.Uu.GetGroupById(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "detail failed", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}
	resp.HandleSuccessResponse(c, "success get group", res)
}

func (lc *UserHandlerImpl) ListGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.SearchRuangan
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	res, err := lc.Uu.ListGroup(ctx)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get data", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success list group", res)
}

func (lc *UserHandlerImpl) UpdateGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	var p pg.UpdateR2GroupParams
	p.ID = utils.StrToInt32(id)

	if err := c.ShouldBind(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak ditemukan"))
		return
	}
	if v, ok := value.(string); ok {
		p.UpdatedBy = utils.StringPtr(v)
	} else {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak valid"))
		return
	}

	res, err := lc.Uu.UpdateGroup(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed update", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}
	resp.HandleSuccessResponse(c, "success update group", res)
}

func (lc *UserHandlerImpl) DelGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Query("id")
	if id == "" {
		resp.HandleErrorResponse(c, "delete failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid id"))
		return
	}
	p := pg.DeleteR2GroupParams{
		ID:        utils.StrToInt32(id),
		DeletedBy: nil,
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak ditemukan"))
		return
	}
	if v, ok := value.(string); ok {
		p.DeletedBy = utils.StringPtr(v)
	} else {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak valid"))
		return
	}

	if err := lc.Uu.DeleteGroupById(ctx, p); err != nil {
		resp.HandleErrorResponse(c, "delete failed", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}
	resp.HandleSuccessResponse(c, "success create role", nil)
}

func (lc *UserHandlerImpl) ListUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.SearchUser
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	res, err := lc.Uu.ListUser(ctx, req)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get data", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success get list users", res)
}

func (lc *UserHandlerImpl) UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	uid, errParse := uuid.FromString(id)
	if errParse != nil {
		resp.HandleErrorResponse(c, "failed to parse id", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid uuid"))
		return
	}

	var p request.UpdateUser
	p.ID = uid

	if err := c.ShouldBind(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak ditemukan"))
		return
	}
	if v, ok := value.(string); ok {
		p.UpdatedBy = utils.StringPtr(v)
	} else {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak valid"))
		return
	}

	res, err := lc.Uu.UpdateUserPartial(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed update", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}
	resp.HandleSuccessResponse(c, "success update user", res)
}

func (lc *UserHandlerImpl) UserById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		resp.HandleErrorResponse(c, "detail failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid id"))
		return
	}
	uid, errParse := uuid.FromString(id)
	if errParse != nil {
		resp.HandleErrorResponse(c, "detail failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid uuid"))
		return
	}

	res, err := lc.Uu.GetUserId(ctx, uid)
	if err != nil {
		resp.HandleErrorResponse(c, "detail failed", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}
	resp.HandleSuccessResponse(c, "success get user detail", res)
}

func (lc *UserHandlerImpl) UserRoleByUserId(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		resp.HandleErrorResponse(c, "detail failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid id"))
		return
	}
	uid, errParse := uuid.FromString(id)
	if errParse != nil {
		resp.HandleErrorResponse(c, "detail failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid uuid"))
		return
	}

	res, err := lc.Uu.GetUserRoleByUserId(ctx, uid)
	if err != nil {
		resp.HandleErrorResponse(c, "detail failed", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}
	resp.HandleSuccessResponse(c, "success get user role", res)
}

func (lc *UserHandlerImpl) UpdateRolePolicyByRoleId(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	var p request.UpdateRolePolicy
	p.RoleID = utils.StrToInt32(id)

	if err := c.ShouldBind(&p); err != nil {
		resp.HandleErrorResponse(c, "failed bind json", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}
	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak ditemukan"))
		return
	}
	if v, ok := value.(string); ok {
		p.CreatedBy = utils.StringPtr(v)
	} else {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak valid"))
		return
	}

	res, err := lc.Uu.UpdateRolePolicy(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to update policy", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success policy to role", res)
}

func (lc *UserHandlerImpl) CreateNewUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var p request.Register
	if err := c.ShouldBind(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if ok {
		if v, ok2 := value.(string); ok2 {
			p.CreatedBy = utils.StringPtr(v)
		}
	}

	res, err := lc.Uu.RegisterWithPassword(ctx, p)
	if err != nil {
		// registration failure likely internal or business error; wrap
		resp.HandleErrorResponse(c, "failed create user", pkg.WrapError(err, pkg.ErrorCodeInternal, "failed create user"))
		return
	}
	resp.HandleSuccessResponse(c, "success create new user", res)
}
