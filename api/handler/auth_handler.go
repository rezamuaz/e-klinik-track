package handler

import (
	"context"
	"e-klinik/api/helper"
	"e-klinik/config"
	"e-klinik/infra/pg"
	"e-klinik/pkg"
	"e-klinik/utils"

	"e-klinik/internal/domain/request"
	"e-klinik/internal/domain/resp"
	"e-klinik/internal/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	Refresh(c *gin.Context)
}

type AuthHandlerImpl struct {
	Uu  *usecase.UserUsecaseImpl
	Mu  *usecase.MainUsecaseImpl
	Cfg *config.Config
}

func NewAuthHandler(Uu *usecase.UserUsecaseImpl, Mu *usecase.MainUsecaseImpl, cfg *config.Config) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		Uu:  Uu,
		Cfg: cfg,
		Mu:  Mu,
	}
}

func (lc *AuthHandlerImpl) Login(c *gin.Context) {
	var request request.Login

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	user, err := lc.Uu.LoginWithPassword(c, request.Username, request.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, helper.GenerateBaseResponseWithAnyError(nil, false, helper.NotFoundError, err.Error()))
		return
	}

	// expUnix := int64(user.RefreshToken.Exp) // example UNIX timestamp
	// expTime := time.Unix(expUnix, 0)

	// c.SetCookie("auth_token", user.RefreshToken.Token, int(expUnix-time.Now().Unix()), "/", "localhost", true, // true in production with HTTPS
	// 	true)

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(user, true, helper.Success))

}

func (lc *AuthHandlerImpl) Register(c *gin.Context) {
	var request request.Register

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	user, err := lc.Uu.RegisterWithPassword(c, request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.GenerateBaseResponseWithAnyError(nil, false, helper.InternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(user, true, helper.Success))

}

func (lc *AuthHandlerImpl) Refresh(c *gin.Context) {
	var request request.Refresh
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	user, err := lc.Uu.Refresh(c, request.RefreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.GenerateBaseResponseWithAnyError(nil, false, helper.InternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(user, true, helper.Success))

}

func (lc *AuthHandlerImpl) Logout(c *gin.Context) {

	cookie, err := c.Cookie("auth_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponse(nil, false, helper.NotFoundError))
		return
	}

	user, err := lc.Uu.Logout(c, cookie)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.GenerateBaseResponseWithAnyError(nil, false, helper.InternalError, err.Error()))
		return
	}

	c.SetCookie(
		"auth_token", // nama cookie kamu
		"",           // kosongkan
		-1,           // expired segera
		"/",          // seluruh path
		"localhost",  // domain default
		false,        // secure: true jika pakai https
		true,         // httpOnly
	)

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(user, true, helper.Success))

}

// func (lc *AuthHandler) SearchMU(c *gin.Context) {
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

// func (lc *AuthHandler) SearchMAL(c *gin.Context) {
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

// func (lc *AuthHandler) MangaDetailMAL(c *gin.Context) {
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

// func (lc *AuthHandler) MangaDetailMU(c *gin.Context) {
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

func (lc *AuthHandlerImpl) AddRoleUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreateUserRoleParams
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created role", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.CreatedBy = utils.StringPtr(value.(string))

	user, err := lc.Uu.AddRoleForUser(ctx, p)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.GenerateBaseResponseWithAnyError(nil, false, helper.InternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(user, true, helper.Success))

}

func (lc *AuthHandlerImpl) AddRolePolicy(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreateUserRoleParams
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created role", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.CreatedBy = utils.StringPtr(value.(string))

	user, err := lc.Uu.AddRoleForUser(ctx, p)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.GenerateBaseResponseWithAnyError(nil, false, helper.InternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(user, true, helper.Success))

}

func (lc *AuthHandlerImpl) AddMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreateR1ViewParams
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created role", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}

	user, err := lc.Uu.AddMenu(ctx, p)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.GenerateBaseResponseWithAnyError(nil, false, helper.InternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(user, true, helper.Success))
}
func (lc *AuthHandlerImpl) ListMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchMenu

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	result, err := lc.Uu.ListMenu(ctx, req)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed get menu", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, 0))

}

func (lc *AuthHandlerImpl) UpdateMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Param("id")
	var p pg.UpdateR1ViewParams
	p.ID = utils.StrToInt32(id)
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update menu", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := lc.Uu.EditMenu(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update menu", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (lc *AuthHandlerImpl) DelMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Query("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "delete menu failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}
	p := pg.DeleteR1ViewParams{
		ID:        utils.StrToInt32(id),
		DeletedBy: nil}
	value, _ := c.Get("nama")
	p.DeletedBy = utils.StringPtr(value.(string))

	err := lc.Uu.DeleteMenu(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "delete menu failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(id, true, 0))
}
func (lc *AuthHandlerImpl) MenuDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Param("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "detail  failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}
	p := utils.StrToInt32(id)

	res, err := lc.Uu.MenuById(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "detail failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}
func (lc *AuthHandlerImpl) ListAccess(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchMenu

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	result, err := lc.Uu.AccessList(ctx, req)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed get menu", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, 0))

}

func (lc *AuthHandlerImpl) CreateRole(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreateR4RoleParams
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.CreatedBy = utils.StringPtr(value.(string))

	result, err := lc.Uu.AddRole(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))

}
func (lc *AuthHandlerImpl) RoleById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Param("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "detail failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}
	p := utils.StrToInt32(id)

	res, err := lc.Uu.GetRoleById(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "detail failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}
func (lc *AuthHandlerImpl) ListRole(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchRuangan

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	result, err := lc.Uu.ListRole(ctx)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed get data", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, 0))

}

func (lc *AuthHandlerImpl) UpdateRole(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Param("id")
	var p pg.UpdateR4RoleParams
	p.ID = utils.StrToInt32(id)
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := lc.Uu.UpdateRole(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update ", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (lc *AuthHandlerImpl) DelRole(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Query("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "delete failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}
	p := pg.DeleteR4RoleParams{
		ID:        utils.StrToInt32(id),
		DeletedBy: nil}
	value, _ := c.Get("nama")
	p.DeletedBy = utils.StringPtr(value.(string))

	err := lc.Uu.DeleteRoleById(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "delete failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(id, true, 0))
}

func (lc *AuthHandlerImpl) CreateGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreateR2GroupParams
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.CreatedBy = utils.StringPtr(value.(string))

	result, err := lc.Uu.AddGroup(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))

}
func (lc *AuthHandlerImpl) GroupById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Param("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "detail failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}
	p := utils.StrToInt32(id)

	res, err := lc.Uu.GetGroupById(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "detail failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}
func (lc *AuthHandlerImpl) ListGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchRuangan

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	result, err := lc.Uu.ListGroup(ctx)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed get data", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, 0))

}

func (lc *AuthHandlerImpl) UpdateGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Param("id")
	var p pg.UpdateR2GroupParams
	p.ID = utils.StrToInt32(id)
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := lc.Uu.UpdateGroup(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update ", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (lc *AuthHandlerImpl) DelGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Query("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "delete failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}
	p := pg.DeleteR2GroupParams{
		ID:        utils.StrToInt32(id),
		DeletedBy: nil}
	value, _ := c.Get("nama")
	p.DeletedBy = utils.StringPtr(value.(string))

	err := lc.Uu.DeleteGroupById(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "delete failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(id, true, 0))
}

func (lc *AuthHandlerImpl) ListUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchUser

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	result, err := lc.Uu.ListUser(ctx, req)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed get data", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, 0))

}

func (lc *AuthHandlerImpl) UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Param("id")
	var p request.UpdateUser
	p.ID = uuid.Must(uuid.FromString(id))
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := lc.Uu.UpdateUserPartial(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update ", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (lc *AuthHandlerImpl) UserById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Param("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "detail failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}

	p := uuid.Must(uuid.FromString(id))

	res, err := lc.Uu.GetUserId(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "detail failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (lc *AuthHandlerImpl) UserRoleByUserId(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Param("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "detail failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}

	p := uuid.Must(uuid.FromString(id))

	res, err := lc.Uu.GetUserRoleByUserId(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "detail failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (lc *AuthHandlerImpl) ViewRoleId(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Param("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "detail failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}

	p := utils.StrToInt32(id)

	res, err := lc.Uu.GetViewByRoleId(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "detail failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (lc *AuthHandlerImpl) AddRolePolicyByRoleId(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Param("id")
	var p request.UpdateRolePolicy
	p.RoleID = utils.StrToInt32(id)
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update policy", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.CreatedBy = utils.StringPtr(value.(string))

	res, err := lc.Uu.AddRolePolicy(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update policy", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (lc *AuthHandlerImpl) CreateNewUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	var p request.Register

	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed create user", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.CreatedBy = utils.StringPtr(value.(string))

	res, err := lc.Uu.RegisterWithPassword(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed create user", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (lc *AuthHandlerImpl) GetViewUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	value, _ := c.Get("Id")
	id := uuid.Must(uuid.FromString(value.(string)))

	res, err := lc.Uu.GetViewUser(ctx, id)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed get menu", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}
