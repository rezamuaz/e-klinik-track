package handler

import (
	"e-klinik/config"

	"e-klinik/internal/domain/resp"
	"e-klinik/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SummaryHandler interface {
	GetRekapKehadiranGlobalHarian(c *gin.Context)
	GetRekapSKPGlobalHarian(c *gin.Context)
	GetRekapKehadiranPerFasilitasHarian(c *gin.Context)
	ChartGetHarianSKPPersentase(c *gin.Context)
	ChartGetHariIniSKPPersentase(c *gin.Context)
}

type SummaryHandlerImpl struct {
	cfg *config.Config
	su  usecase.SummaryUsecase
}

func NewSummaryHandler(su usecase.SummaryUsecase, cfg *config.Config) *SummaryHandlerImpl {
	return &SummaryHandlerImpl{
		cfg: cfg,
		su:  su,
	}
}

func (h *SummaryHandlerImpl) GetRekapKehadiranGlobalHarian(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := h.su.GetRekapKehadiranGlobalHarian(ctx)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get rekap kehadiran global harian", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get rekap kehadiran global harian", res)
}

func (h *SummaryHandlerImpl) GetRekapSKPGlobalHarian(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := h.su.GetRekapSKPGlobalHarian(ctx)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get rekap SKP global harian", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get rekap SKP global harian", res)
}

func (h *SummaryHandlerImpl) GetRekapKehadiranPerFasilitasHarian(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := h.su.GetRekapKehadiranPerFasilitasHarian(ctx)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get rekap kehadiran per fasilitas harian", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get rekap kehadiran per fasilitas harian", res)
}

func (h *SummaryHandlerImpl) ChartGetHarianSKPPersentase(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := h.su.ChartGetHarianSKPPersentase(ctx)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get chart harian SKP persentase", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get chart harian SKP persentase", res)
}

func (h *SummaryHandlerImpl) ChartGetHariIniSKPPersentase(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := h.su.ChartGetHariIniSKPPersentase(ctx)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get chart hari ini SKP persentase", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get chart hari ini SKP persentase", res)
}

// func (h *SummaryHandlerImpl) GetKehadiranByMahasiswaStatus(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	var arg pg.GetKehadiranByPembimbingUserIdParams
// 	value, _ := c.Get("Id")
// 	id := uuid.Must(uuid.FromString(value.(string)))
// 	arg.UserID = &id

// 	result, err := h.su.GetKehadiranByPembimbingStatus(ctx, arg)
// 	if err != nil {
// 		resp.HandleErrorResponse(c, "failed get kehadiran by mahasiswa", err)
// 		return
// 	}

// 	resp.HandleSuccessResponse(c, "success get kehadiran by mahasiswa", result)
// }
