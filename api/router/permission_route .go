package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func Permission(group *gin.RouterGroup, h *handler.PermissionHandlerImpl) {

	group.GET("/tree", h.ListPermission)
	group.GET("/role/:id", h.PermissionByRoleId)
	group.GET("/users", h.UserViewPermission)
	group.POST("", h.AddMenu)
	group.GET("", h.ListMenu)
	group.GET("/:id", h.MenuDetail)
	group.PUT("/:id", h.UpdateMenu)
	group.DELETE("/:id", h.DelMenu)

}
