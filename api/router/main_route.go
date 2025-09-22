package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func Main(group *gin.RouterGroup, h *handler.MainHandlerImpl) {

	group.POST("/fasilitas", h.CreateFasilitasKesehatan)
	group.GET("/fasilitas", h.ListFaslitasKesehatan)
	group.PUT("/fasilita", h.UpdateFasilitasKesehatan)
	group.DELETE("/fasilitas", h.DelFasilitasKesehatan)
	group.POST("/mata_kuliah", h.CreateMatakuliah)
	group.GET("/mata_kuliah", h.ListMataKuliah)
	group.PUT("/mata_kuliah", h.UpdateMataKuliah)
	group.DELETE("/mata_kuliah", h.DelMataKuliah)
	group.POST("/kontrak", h.CreateKontrak)
	group.GET("/kontrak", h.ListKontrak)
	group.PUT("/kontrak", h.UpdateKontrak)
	group.DELETE("/kontrak", h.DelKontrak)
	group.POST("/ruangan", h.CreateRuangan)
	group.GET("/ruangan", h.ListRuangan)
	group.PUT("/ruangan", h.UpdateRuangan)
	group.DELETE("/ruangan", h.DelRuangan)
	group.POST("/kehadiran", h.CreateKehadiran)
	group.GET("/kehadiran", h.ListKehadiran)
	group.PUT("/kehadiran", h.UpdateKehadiran)
	group.DELETE("/kehadiran", h.DeleteKehadiran)

}
