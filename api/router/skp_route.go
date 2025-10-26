package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func Skp(group *gin.RouterGroup, h *handler.SkpHandlerImpl) {

	//SKP
	group.GET("/intervensi", h.ListIntervensi)

}
