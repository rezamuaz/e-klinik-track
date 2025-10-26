package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func User(group *gin.RouterGroup, h *handler.UserHandlerImpl) {

	group.POST("/register", h.Register)
	group.DELETE("/logout", h.Logout)

	group.POST("/roles", h.CreateRole)
	group.GET("/roles", h.ListRole)
	group.GET("/roles/:id", h.RoleById)
	group.PUT("/roles/:id", h.UpdateRole)
	group.DELETE("/roles/:id", h.DelRole)
	group.PUT("/roles/policies/:id", h.UpdateRolePolicyByRoleId)
	group.POST("/group", h.CreateGroup)
	group.GET("/group", h.ListGroup)
	group.GET("/group/:id", h.GroupById)
	group.PUT("/group/:id", h.UpdateGroup)
	group.DELETE("/group/:id", h.DelGroup)
	group.GET("", h.ListUsers)
	group.GET("/:id", h.UserById)
	group.PUT("/:id", h.UpdateUser)
	group.POST("", h.CreateNewUser)
	group.POST("/user-roles", h.AddRoleUser)
	group.GET("/user-roles/:id", h.UserRoleByUserId)

}
