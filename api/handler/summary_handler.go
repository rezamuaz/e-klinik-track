package handler

import (
	"e-klinik/config"
	"e-klinik/pkg"
	"e-klinik/utils"
	"time"

	"e-klinik/internal/domain/request"
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
	GetGlobalSKPPersentaseTahunanOtomatis(c *gin.Context)
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

func (h *SummaryHandlerImpl) RekapSkpTercapaiMahasiswaByDate(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchSkpTercapai
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	res, err := h.su.RekapSkpTercapaiMahasiswaByDate(ctx, req)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get chart hari ini SKP persentase", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get chart hari ini SKP persentase", res)
}

func (h *SummaryHandlerImpl) GetGlobalSKPPersentaseTahunanOtomatis(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	res, err := h.su.GetGlobalSKPPersentaseTahunanOtomatis(ctx)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get chart hari ini SKP persentase", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get chart hari ini SKP persentase", res)
}
