package handler

import (
	"e-klinik/config"
	"e-klinik/infra/pg"
	"e-klinik/pkg"
	"e-klinik/utils"

	"e-klinik/internal/domain/request"
	"e-klinik/internal/domain/resp"
	"e-klinik/internal/usecase"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type FasilitasHandler interface {
	CreateFasilitasKesehatan(c *gin.Context)
	ListFaslitasKesehatan(c *gin.Context)
	UpdateFasilitasKesehatan(c *gin.Context)
	DelFasilitasKesehatan(c *gin.Context)
	ListPropinsi(c *gin.Context)
	ListKabupaten(c *gin.Context)
}

type FasilitasHandlerImpl struct {
	cfg *config.Config
	fu  usecase.FasilitasUsecase
}

func NewFasilitasHandler(fu usecase.FasilitasUsecase, cfg *config.Config) *FasilitasHandlerImpl {
	return &FasilitasHandlerImpl{
		cfg: cfg,
		fu:  fu,
	}
}

func (h *FasilitasHandlerImpl) CreateFasilitasKesehatan(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreateFasilitasKesehatanParams
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to create fasilitas kesehatan", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "failed to create fasilitas kesehatan", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}
	p.CreatedBy = utils.StringPtr(value.(string))

	result, err := h.fu.AddFasilitasKesehatan(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to create fasilitas kesehatan", pkg.WrapError(err, pkg.ErrorCodeInternal, "database error"))
		return
	}

	resp.HandleSuccessResponse(c, "success create fasilitas kesehatan", result)
}

func (h *FasilitasHandlerImpl) ListFaslitasKesehatan(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchFasilitasKesehatan
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	result, err := h.fu.ListFasilitasKesehatan(ctx, req)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to get list fasilitas kesehatan", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success get list fasilitas kesehatan", result)
}

func (h *FasilitasHandlerImpl) UpdateFasilitasKesehatan(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	idStr := c.Query("id")
	if idStr == "" {
		resp.HandleErrorResponse(c, "invalid request", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "missing id parameter"))
		return
	}

	var p pg.UpdateFasilitasKesehatanPartialParams
	p.ID = uuid.Must(uuid.FromString(idStr))

	if err := c.ShouldBindJSON(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "failed to update fasilitas kesehatan", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := h.fu.UpdateFasilitasKesehatan(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to update fasilitas kesehatan", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success update fasilitas kesehatan", res)
}

func (h *FasilitasHandlerImpl) DelFasilitasKesehatan(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	idStr := c.Query("id")
	if idStr == "" {
		resp.HandleErrorResponse(c, "invalid request", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "missing id parameter"))
		return
	}

	p := pg.DeleteFasilitasKesehatanParams{
		ID: uuid.Must(uuid.FromString(idStr)),
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "failed to delete fasilitas kesehatan", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}
	p.DeletedBy = utils.StringPtr(value.(string))

	if err := h.fu.DeleteFasilitasKesehatan(ctx, p); err != nil {
		resp.HandleErrorResponse(c, "failed to delete fasilitas kesehatan", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success delete fasilitas kesehatan", gin.H{"id": idStr})
}
func (h *FasilitasHandlerImpl) ListPropinsi(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchPropinsi
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.WrapError(err, pkg.ErrorCodeInvalidArgument, "invalid query"))
		return
	}

	result, err := h.fu.ListPropinsi(ctx, req)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get propinsi", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get propinsi", result)
}

func (h *FasilitasHandlerImpl) ListKabupaten(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchKabupaten
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	result, err := h.fu.ListKabupaten(ctx, req)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get kabupaten", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get kabupaten", result)
}
