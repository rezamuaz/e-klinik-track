package handler

import (
	"context"
	"e-klinik/config"
	"e-klinik/pkg"
	"time"

	"e-klinik/internal/domain/request"
	"e-klinik/internal/domain/resp"
	"e-klinik/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Refresh(c *gin.Context)
}

type AuthHandlerImpl struct {
	Uu *usecase.UserUsecaseImpl

	Cfg *config.Config
}

func NewAuthHandler(Uu *usecase.UserUsecaseImpl, cfg *config.Config) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		Uu:  Uu,
		Cfg: cfg,
	}
}

func (lc *AuthHandlerImpl) Login(c *gin.Context) {
	var req request.Login

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.HandleErrorResponse(
			c,
			"failed to bind JSON",
			pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"),
		)
		return
	}

	user, err := lc.Uu.LoginWithPassword(c, req.Username, req.Password)
	if err != nil {
		resp.HandleErrorResponse(
			c,
			"invalid username or password",
			pkg.ExposeError(pkg.ErrorCodeUnauthorized, "username atau password salah"),
		)
		return
	}

	// Set cookie (optional, uncomment jika mau aktif)
	// expUnix := int64(user.RefreshToken.Exp)
	// expTime := time.Unix(expUnix, 0)
	// c.SetCookie("auth_token", user.RefreshToken.Token, int(expUnix-time.Now().Unix()), "/", "localhost", true, true)

	resp.HandleSuccessResponse(c, "login berhasil", user)
}

func (lc *AuthHandlerImpl) Refresh(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var req request.Refresh

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.HandleErrorResponse(
			c,
			"failed to bind JSON",
			pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"),
		)
		return
	}

	user, err := lc.Uu.Refresh(ctx, req.RefreshToken)
	if err != nil {
		resp.HandleErrorResponse(
			c,
			"failed to refresh access token",
			pkg.WrapError(err, pkg.ErrorCodeInternal, "refresh token invalid or expired"),
		)
		return
	}

	resp.HandleSuccessResponse(c, "refresh token berhasil", user)
}
