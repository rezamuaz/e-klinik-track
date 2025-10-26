package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func Summary(group *gin.RouterGroup, h *handler.SummaryHandlerImpl) {

	//Rekap
	group.GET("/chart/kehadiran/sekarang/global", h.GetRekapKehadiranGlobalHarian)
	group.GET("/chart/kehadiran/sekarang/fasilitas", h.GetRekapKehadiranPerFasilitasHarian)
	group.GET("/chart/skp/sekarang/global", h.GetRekapSKPGlobalHarian)
	group.GET("/chart/skp/seminggu", h.ChartGetHarianSKPPersentase)
	group.GET("/chart/skp/hariini", h.ChartGetHariIniSKPPersentase)
	group.GET("/block/skp/date", h.RekapSkpTercapaiMahasiswaByDate)

}
