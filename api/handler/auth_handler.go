package handler

import (
	"e-klinik/api/helper"
	"e-klinik/config"

	"e-klinik/internal/domain/request"
	"e-klinik/internal/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	Refresh(c *gin.Context)
}

type AuthHandlerImpl struct {
	Uu  *usecase.UserUsecaseImpl
	Cfg *config.Config
}

func NewAuthHandler(Uu *usecase.UserUsecaseImpl, cfg *config.Config) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		Uu:  Uu,
		Cfg: cfg,
	}
}

func (lc *AuthHandlerImpl) Login(c *gin.Context) {
	var request request.Login

	err := c.ShouldBind(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	user, err := lc.Uu.LoginWithPassword(c, request.Username, request.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.GenerateBaseResponseWithAnyError(nil, false, helper.InternalError, err.Error()))
		return
	}

	expUnix := int64(user.RefreshToken.Exp) // example UNIX timestamp
	// expTime := time.Unix(expUnix, 0)

	c.SetCookie("auth_token", user.RefreshToken.Token, int(expUnix-time.Now().Unix()), "/", "localhost", true, // true in production with HTTPS
		true)

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(user.User, true, helper.Success))

}

func (lc *AuthHandlerImpl) Register(c *gin.Context) {
	var request request.Register

	err := c.ShouldBind(&request)
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

	cookie, err := c.Cookie("auth_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponse(nil, false, helper.NotFoundError))
		return
	}

	user, err := lc.Uu.Refresh(c, cookie)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.GenerateBaseResponseWithAnyError(nil, false, helper.InternalError, err.Error()))
		return
	}

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
