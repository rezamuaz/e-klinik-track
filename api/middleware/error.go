package middleware

import (
	"e-klinik/internal/domain/resp"
	"e-klinik/pkg"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err
			appErr := resp.ToAppError(err, pkg.ErrorCodeInternal, "internal server error")
			resp.HandleErrorResponse(c, appErr.Message, appErr)
			c.Abort()
		}
	}
}
