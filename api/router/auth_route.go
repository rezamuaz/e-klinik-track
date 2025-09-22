package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func Auth(group *gin.RouterGroup, h *handler.AuthHandlerImpl) {

	group.POST("/register", h.Register)
	group.POST("/login", h.Login)
	group.GET("/refresh", h.Refresh)

}
