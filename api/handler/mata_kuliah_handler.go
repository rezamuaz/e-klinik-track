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

type MataKuliahHandler interface {
	CreateMatakuliah(c *gin.Context)
	ListMataKuliah(c *gin.Context)
	UpdateMataKuliah(c *gin.Context)
	DelMataKuliah(c *gin.Context)
}

type MataKuliahHandlerImpl struct {
	cfg *config.Config
	mu  usecase.MataKuliahUsecase
}

func NewMataKuliahHandler(mu usecase.MataKuliahUsecase, cfg *config.Config) *MataKuliahHandlerImpl {
	return &MataKuliahHandlerImpl{
		cfg: cfg,
		mu:  mu,
	}
}

func (h *MataKuliahHandlerImpl) CreateMatakuliah(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreateMataKuliahParams
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to create mata kuliah", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "failed to create mata kuliah", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}
	p.CreatedBy = utils.StringPtr(value.(string))

	result, err := h.mu.AddMataKuliah(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to create mata kuliah", pkg.WrapError(err, pkg.ErrorCodeInternal, "database error"))
		return
	}

	resp.HandleSuccessResponse(c, "success create mata kuliah", result)
}

func (h *MataKuliahHandlerImpl) ListMataKuliah(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchMataKuliah
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	result, err := h.mu.ListMataKuliah(ctx, req)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to get list mata kuliah", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success get list mata kuliah", result)
}

func (h *MataKuliahHandlerImpl) UpdateMataKuliah(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	idStr := c.Query("id")
	if idStr == "" {
		resp.HandleErrorResponse(c, "invalid request", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "missing id parameter"))
		return
	}

	var p pg.UpdateMataKuliahParams
	p.ID = uuid.Must(uuid.FromString(idStr))

	if err := c.ShouldBindJSON(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "failed to update mata kuliah", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := h.mu.UpdateMataKuliah(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to update mata kuliah", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success update mata kuliah", res)
}

func (h *MataKuliahHandlerImpl) DelMataKuliah(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	idStr := c.Query("id")
	if idStr == "" {
		resp.HandleErrorResponse(c, "invalid request", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "missing id parameter"))
		return
	}

	p := pg.DeleteMataKuliahParams{
		ID: uuid.Must(uuid.FromString(idStr)),
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "failed to delete mata kuliah", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}
	p.DeletedBy = utils.StringPtr(value.(string))

	if err := h.mu.DeleteMataKuliah(ctx, p); err != nil {
		resp.HandleErrorResponse(c, "failed to delete mata kuliah", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success delete mata kuliah", gin.H{"id": idStr})
}
