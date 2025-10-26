package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func MataKuliah(group *gin.RouterGroup, h *handler.MataKuliahHandlerImpl) {
	//Matakliah
	group.POST("", h.CreateMatakuliah)
	group.GET("", h.ListMataKuliah)
	group.PUT("", h.UpdateMataKuliah)
	group.DELETE("", h.DelMataKuliah)

}
