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

type RuanganHandler interface {
	CreateRuangan(c *gin.Context)
	RuanganDetail(c *gin.Context)
	ListRuangan(c *gin.Context)
	UpdateRuangan(c *gin.Context)
	DelRuangan(c *gin.Context)
	GetRuanganKontrak(c *gin.Context)
}

type RuanganHandlerImpl struct {
	cfg *config.Config
	ru  usecase.RuanganUsecase
}

func NewRuanganHandler(ru usecase.RuanganUsecase, cfg *config.Config) *RuanganHandlerImpl {
	return &RuanganHandlerImpl{
		cfg: cfg,
		ru:  ru,
	}
}

func (h *RuanganHandlerImpl) CreateRuangan(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreateRuanganParams
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to create ruangan", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "failed to create ruangan", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}
	p.CreatedBy = utils.StringPtr(value.(string))

	result, err := h.ru.AddRuangan(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to create ruangan", pkg.WrapError(err, pkg.ErrorCodeInternal, "database error"))
		return
	}

	resp.HandleSuccessResponse(c, "success create ruangan", result)
}

func (h *RuanganHandlerImpl) RuanganDetail(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	idStr := c.Param("id")
	if idStr == "" {
		resp.HandleErrorResponse(c, "invalid request", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "missing ruangan id"))
		return
	}

	id := uuid.Must(uuid.FromString(idStr))
	res, err := h.ru.RuanganById(ctx, id)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to get ruangan detail", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success get ruangan detail", res)
}

func (h *RuanganHandlerImpl) ListRuangan(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchRuangan
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	result, err := h.ru.ListRuangan(ctx, req)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to get ruangan list", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success get ruangan list", result)
}

func (h *RuanganHandlerImpl) UpdateRuangan(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	idStr := c.Param("id")
	if idStr == "" {
		resp.HandleErrorResponse(c, "invalid request", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "missing ruangan id"))
		return
	}

	var p pg.UpdateRuanganPartialParams
	p.ID = uuid.Must(uuid.FromString(idStr))

	if err := c.ShouldBindJSON(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "failed to update ruangan", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := h.ru.UpdateRuangan(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to update ruangan", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success update ruangan", res)
}

func (h *RuanganHandlerImpl) DelRuangan(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	idStr := c.Param("id")
	if idStr == "" {
		resp.HandleErrorResponse(c, "invalid request", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "missing ruangan id"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "failed to delete ruangan", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}

	p := pg.DeleteRuanganParams{
		ID:        uuid.Must(uuid.FromString(idStr)),
		DeletedBy: utils.StringPtr(value.(string)),
	}

	if err := h.ru.DeleteRuangan(ctx, p); err != nil {
		resp.HandleErrorResponse(c, "failed to delete ruangan", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success delete ruangan", gin.H{"id": idStr})
}

func (h *RuanganHandlerImpl) GetRuanganKontrak(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchRuanganByKontrak
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	result, err := h.ru.ListRuanganByKontrak(ctx, req)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get ruangan", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get ruangan", result)
}
