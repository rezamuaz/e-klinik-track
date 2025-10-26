package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func Fasilitas(group *gin.RouterGroup, h *handler.FasilitasHandlerImpl) {
	//Fasilitas
	group.POST("", h.CreateFasilitasKesehatan)
	group.GET("", h.ListFaslitasKesehatan)
	group.GET("/propinsi", h.ListPropinsi)
	group.GET("/kabupaten", h.ListKabupaten)
	group.PUT("/:id", h.UpdateFasilitasKesehatan)
	group.DELETE("/:id", h.DelFasilitasKesehatan)
}
