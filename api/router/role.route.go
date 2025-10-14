package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func Role(group *gin.RouterGroup, h *handler.AuthHandlerImpl) {
	group.POST("/user_role", h.AddRoleUser)
	group.POST("/menu", h.AddMenu)
	group.GET("/menu", h.ListMenu)
	group.GET("/menu/:id", h.MenuDetail)
	group.PUT("/menu/:id", h.UpdateMenu)
	group.DELETE("/menu/:id", h.DelMenu)
	group.GET("/access", h.ListAccess)
	group.POST("/role", h.CreateRole)
	group.GET("/role", h.ListRole)
	group.GET("/role/:id", h.RoleById)
	group.PUT("/role/:id", h.UpdateRole)
	group.DELETE("/role/:id", h.DelRole)
	group.POST("/group", h.CreateGroup)
	group.GET("/group", h.ListGroup)
	group.GET("/group/:id", h.GroupById)
	group.PUT("/group/:id", h.UpdateGroup)
	group.DELETE("/group/:id", h.DelGroup)
	group.GET("/users", h.ListUsers)
	group.GET("/users/:id", h.UserById)
	group.PUT("/users/:id", h.UpdateUser)
	group.POST("/users", h.CreateNewUser)
	group.GET("/user_role/:id", h.UserRoleByUserId)
	group.GET("/view_role/role/:id", h.ViewRoleId)
	group.PUT("/policy/role/:id", h.AddRolePolicyByRoleId)
	group.GET("/user_view", h.GetViewUser)

}
