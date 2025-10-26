package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func Kontrak(group *gin.RouterGroup, h *handler.KontrakHandlerImpl) {

	//Kontak/Kerjasama
	group.POST("", h.CreateKontrak)
	group.GET("", h.ListKontrak)
	group.GET("/aktif", h.ListAktifKontrak)
	group.GET("/:id", h.KontrakDetail)
	group.PUT("/:id", h.UpdateKontrak)
	group.DELETE("/:id", h.DelKontrak)

}
