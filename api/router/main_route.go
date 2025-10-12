package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func Main(group *gin.RouterGroup, h *handler.MainHandlerImpl) {

	group.POST("/fasilitas", h.CreateFasilitasKesehatan)
	group.GET("/fasilitas", h.ListFaslitasKesehatan)
	group.GET("/propinsi", h.ListPropinsi)
	group.GET("/kabupaten", h.ListKabupaten)
	group.PUT("/fasilitas", h.UpdateFasilitasKesehatan)
	group.DELETE("/fasilitas", h.DelFasilitasKesehatan)
	group.POST("/mata_kuliah", h.CreateMatakuliah)
	group.GET("/mata_kuliah", h.ListMataKuliah)
	group.PUT("/mata_kuliah", h.UpdateMataKuliah)
	group.DELETE("/mata_kuliah", h.DelMataKuliah)
	group.POST("/kontrak", h.CreateKontrak)
	group.GET("/kontrak", h.ListKontrak)
	group.GET("/aktif_kontrak", h.ListAktifKontrak)
	group.GET("/kontrak/:id", h.KontrakDetail)
	group.PUT("/kontrak", h.UpdateKontrak)
	group.DELETE("/kontrak/:id", h.DelKontrak)
	group.POST("/ruangan", h.CreateRuangan)
	group.GET("/ruangan", h.ListRuangan)
	group.GET("/ruangan/:id", h.RuanganDetail)
	group.PUT("/ruangan/:id", h.UpdateRuangan)
	group.DELETE("/ruangan/:id", h.DelRuangan)
	group.POST("/kehadiran", h.CreateKehadiran)
	group.GET("/kehadiran", h.ListKehadiran)
	group.PUT("/kehadiran", h.UpdateKehadiran)
	group.DELETE("/kehadiran", h.DeleteKehadiran)
	group.POST("/kehadiran_skp", h.CreateKehadiranSkp)
	group.GET("/kehadiran_skp", h.ListKehadiranSkp)
	group.PUT("/kehadiran_skp", h.UpdateKehadiranSkp)
	group.DELETE("/kehadiran_skp", h.DeleteKehadiranSkp)
	group.GET("/intervensi", h.ListIntervensi)

}
