package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func Actor(group *gin.RouterGroup, h *handler.ActorHandlerImpl) {

	group.GET("/users/roles", h.GetUsersByRoles)
	//Pembimbing Klinik
	group.POST("/pembimbing-klinik", h.CreatePembimbingKlinik)
	group.GET("/pembimbing-klinik/kontrak/:id", h.ListPembimbingKlinikByKontrak)

}
