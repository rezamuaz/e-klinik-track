package middleware

import (
	"e-klinik/internal/domain/resp"
	"e-klinik/pkg"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			resp.HandleErrorResponse(c, "missing authorization header", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "unauthorized"))
			c.Abort()
			return
		}

		// Format header: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			resp.HandleErrorResponse(c, "invalid authorization format", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "invalid token format"))
			c.Abort()
			return
		}

		token := parts[1]
		authorized, _, err := pkg.IsAuthorized(token, secret)
		if err != nil {
			resp.HandleErrorResponse(c, "token validation failed", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "token validation error"))
			c.Abort()
			return
		}

		if !authorized {
			resp.HandleErrorResponse(c, "unauthorized access", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "unauthorized"))
			c.Abort()
			return
		}

		user, err := pkg.ExtractClaimsFromToken(token, secret)
		if err != nil {
			if err.Error() == "expire" {
				resp.HandleErrorResponse(c, "token expired", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "token expired"))
				c.Abort()
				return
			}
			resp.HandleErrorResponse(c, "invalid token", pkg.WrapError(err, pkg.ErrorCodeUnauthorized, "invalid token"))
			c.Abort()
			return
		}

		// Set user data ke context untuk digunakan downstream
		c.Set("username", user.Username)
		c.Set("Id", user.Subject)
		c.Set("nama", user.Nama)

		c.Next()
	}
}
