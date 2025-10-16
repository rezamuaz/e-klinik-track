package middleware

import (
	"e-klinik/internal/domain/resp"
	"e-klinik/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, _, err := pkg.IsAuthorized(authToken, secret)

			if authorized {
				user, err := pkg.ExtractClaimsFromToken(authToken, secret)
				if err != nil {
					if err.Error() == "expire" {
						c.AbortWithStatusJSON(http.StatusForbidden, resp.GenerateBaseResponseWithAnyError(nil, false, resp.ForbiddenError, err.Error()))
						return
					}
					c.AbortWithStatusJSON(http.StatusUnauthorized, resp.GenerateBaseResponseWithAnyError(nil, false, resp.InternalError, err.Error()))
					return
				}
				c.Set("username", user.Username)
				c.Set("Id", user.Subject)
				c.Set("nama", user.Nama)
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp.GenerateBaseResponseWithAnyError(nil, false, resp.AuthError, err.Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp.GenerateBaseResponse(nil, false, resp.AuthError))
	}
}
