package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func Main(group *gin.RouterGroup, h *handler.MainHandlerImpl) {

	//Fasilitas
	group.POST("/fasilitas", h.CreateFasilitasKesehatan)
	group.GET("/fasilitas", h.ListFaslitasKesehatan)
	group.GET("/propinsi", h.ListPropinsi)
	group.GET("/kabupaten", h.ListKabupaten)
	group.PUT("/fasilitas", h.UpdateFasilitasKesehatan)
	group.DELETE("/fasilitas", h.DelFasilitasKesehatan)
	//Matakliah
	group.POST("/mata_kuliah", h.CreateMatakuliah)
	group.GET("/mata_kuliah", h.ListMataKuliah)
	group.PUT("/mata_kuliah", h.UpdateMataKuliah)
	group.DELETE("/mata_kuliah", h.DelMataKuliah)
	//Kontak/Kerjasama
	group.POST("/kontrak", h.CreateKontrak)
	group.GET("/kontrak", h.ListKontrak)
	group.GET("/aktif_kontrak", h.ListAktifKontrak)
	group.GET("/kontrak/:id", h.KontrakDetail)
	group.PUT("/kontrak/:id", h.UpdateKontrak)
	group.DELETE("/kontrak/:id", h.DelKontrak)
	//Ruangan
	group.POST("/ruangan", h.CreateRuangan)
	group.GET("/ruangan", h.ListRuangan)
	group.GET("/ruangan/kontrak", h.GetRuanganKontrak)
	group.GET("/ruangan/:id", h.RuanganDetail)
	group.PUT("/ruangan/:id", h.UpdateRuangan)
	group.DELETE("/ruangan/:id", h.DelRuangan)
	//Kehadiran
	group.POST("/kehadiran", h.CreateKehadiran)
	group.GET("/kehadiran", h.ListKehadiran)
	group.PUT("/kehadiran", h.UpdateKehadiran)
	group.DELETE("/kehadiran", h.DeleteKehadiran)
	group.GET("/kehadiran/status", h.CheckKehadiran)
	group.GET("/kehadiran/rekap_total", h.RekapKehadiranMahasiswa)
	group.GET("/kehadiran/rekap", h.RekapKehadiranMahasiswaDetail)
	group.GET("/kehadiran/pembimbing_klinik", h.GetKehadiranByPembimbingStatus)
	group.GET("/kehadiran/mahasiswa", h.GetKehadiranByMahasiswaStatus)

	//SKP Kehadiran
	group.POST("/kehadiran_skp", h.CreateSyncKehadiranSkp)
	group.GET("/kehadiran_skp", h.ListKehadiranSkp)
	group.GET("/kehadiran_skp/active", h.SkpByKehadiranId)
	group.GET("/kehadiran_skp/active_nama", h.IntervensiKehadiranId)
	group.PUT("/kehadiran_skp", h.UpdateKehadiranSkp)
	group.DELETE("/kehadiran_skp", h.DeleteKehadiranSkp)
	group.POST("/kehadiran_skp/approve", h.ApproveKehadiranSkp)
	//SKP
	group.GET("/intervensi", h.ListIntervensi)
	group.GET("/users/roles", h.GetUsersByRoles)
	//Pembimbing Klinik
	group.POST("/pembimbing_klinik", h.CreatePembimbingKlinik)
	group.GET("/pembimbing_klinik/kontrak/:id", h.ListPembimbingKlinikByKontrak)

}
