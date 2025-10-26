package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func Ruangan(group *gin.RouterGroup, h *handler.RuanganHandlerImpl) {

	group.POST("", h.CreateRuangan)
	group.GET("", h.ListRuangan)
	group.GET("/kontrak", h.GetRuanganKontrak)
	group.GET("/:id", h.RuanganDetail)
	group.PUT("/:id", h.UpdateRuangan)
	group.DELETE("/:id", h.DelRuangan)

}
