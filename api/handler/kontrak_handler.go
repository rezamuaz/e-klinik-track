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

type KontrakHandler interface {
	CreateKontrak(c *gin.Context)
	ListKontrak(c *gin.Context)
	ListAktifKontrak(c *gin.Context)
	UpdateKontrak(c *gin.Context)
	KontrakDetail(c *gin.Context)
	DelKontrak(c *gin.Context)
}

type KontrakHandlerImpl struct {
	cfg *config.Config
	ku  usecase.KontrakUsecase
}

func NewKontrakHandler(ku usecase.KontrakUsecase, cfg *config.Config) *KontrakHandlerImpl {
	return &KontrakHandlerImpl{
		cfg: cfg,
		ku:  ku,
	}
}

func (h *KontrakHandlerImpl) CreateKontrak(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var p request.CreateKontrak
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to create kontrak", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "failed to create kontrak", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}
	p.CreatedBy = utils.StringPtr(value.(string))

	result, err := h.ku.AddKontrak(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to create kontrak", err)
		return
	}

	resp.HandleSuccessResponse(c, "success create kontrak", result)
}

func (h *KontrakHandlerImpl) ListKontrak(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchKontrak
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	result, err := h.ku.ListKontrak(ctx, req)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to get kontrak", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get kontrak list", result)
}

func (h *KontrakHandlerImpl) ListAktifKontrak(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchAktifKontrak
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	result, err := h.ku.ListAktifKontrak(ctx, req.FasilitasNama)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to get active kontrak", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success get active kontrak", result)
}

func (h *KontrakHandlerImpl) UpdateKontrak(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	idStr := c.Param("id")
	if idStr == "" {
		resp.HandleErrorResponse(c, "invalid request", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "missing kontrak id"))
		return
	}

	var p pg.UpdateKontrakPartialParams
	p.ID = uuid.Must(uuid.FromString(idStr))

	if err := c.ShouldBindJSON(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "failed to update kontrak", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := h.ku.UpdateKontrak(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to update kontrak", err)
		return
	}

	resp.HandleSuccessResponse(c, "success update kontrak", res)
}

func (h *KontrakHandlerImpl) KontrakDetail(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	idStr := c.Param("id")
	if idStr == "" {
		resp.HandleErrorResponse(c, "invalid request", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "missing kontrak id"))
		return
	}

	id := uuid.Must(uuid.FromString(idStr))
	res, err := h.ku.KontrakById(ctx, id)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to get kontrak detail", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success get kontrak detail", res)
}

func (h *KontrakHandlerImpl) DelKontrak(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	idStr := c.Param("id")
	if idStr == "" {
		resp.HandleErrorResponse(c, "invalid request", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "missing kontrak id"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "failed to delete kontrak", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}

	p := pg.DeleteKontrakParams{
		ID:        uuid.Must(uuid.FromString(idStr)),
		DeletedBy: utils.StringPtr(value.(string)),
	}

	if err := h.ku.DeleteKontrak(ctx, p); err != nil {
		resp.HandleErrorResponse(c, "failed to delete kontrak", err)
		return
	}

	resp.HandleSuccessResponse(c, "success delete kontrak", gin.H{"id": idStr})
}
