package handler

import (
	"context"
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

type SkpKehadiranHandler interface {
	CreateSyncKehadiranSkp(c *gin.Context)
	ListKehadiranSkp(c *gin.Context)
	UpdateKehadiranSkp(c *gin.Context)
	DeleteKehadiranSkp(c *gin.Context)
	SkpByKehadiranId(c *gin.Context)
	IntervensiKehadiranId(c *gin.Context)
	ApproveKehadiranSkp(c *gin.Context)
}

type SkpKehadiranHandlerImpl struct {
	cfg *config.Config
	sk  usecase.SkpKehadiranUsecase
}

func NewSkpKehadiranHandler(sk usecase.SkpKehadiranUsecase, cfg *config.Config) *SkpKehadiranHandlerImpl {
	return &SkpKehadiranHandlerImpl{
		cfg: cfg,
		sk:  sk,
	}
}

func (h *SkpKehadiranHandlerImpl) CreateSyncKehadiranSkp(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.SyncKehadiranSkpParams
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.HandleErrorResponse(c, "failed create kehadiran", pkg.WrapError(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}

	value, _ := c.Get("nama")
	p.Actor = value.(string)

	id, _ := c.Get("Id")
	p.UserID = uuid.Must(uuid.FromString(id.(string)))

	result, err := h.sk.SyncSkpKehadiran(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed create kehadiran", err)
		return
	}

	resp.HandleSuccessResponse(c, "success create kehadiran", result)
}

func (h *SkpKehadiranHandlerImpl) ListKehadiranSkp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req request.SearchKehadiranSkp
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.WrapError(err, pkg.ErrorCodeInvalidArgument, "invalid query"))
		return
	}

	result, err := h.sk.ListKehadiranSkp(ctx, req)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get kehadiran", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get kehadiran", result)
}

func (h *SkpKehadiranHandlerImpl) UpdateKehadiranSkp(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.UpdateKehadiranSkpParams
	p.ID = uuid.Must(uuid.FromString(c.Param("id")))
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.HandleErrorResponse(c, "failed update kehadiran", pkg.WrapError(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}

	value, _ := c.Get("nama")
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := h.sk.UpdateKehadiranSkp(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed update kehadiran", err)
		return
	}

	resp.HandleSuccessResponse(c, "success update kehadiran", res)
}

func (h *SkpKehadiranHandlerImpl) DeleteKehadiranSkp(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		resp.HandleErrorResponse(c, "delete kehadiran failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid id"))
		return
	}

	p := pg.DeleteKehadiranSkpParams{
		ID: uuid.Must(uuid.FromString(id)),
	}

	value, _ := c.Get("nama")
	p.DeletedBy = utils.StringPtr(value.(string))

	if err := h.sk.DeleteKehadiranSkp(ctx, p); err != nil {
		resp.HandleErrorResponse(c, "failed delete kehadiran", err)
		return
	}

	resp.HandleSuccessResponse(c, "success delete kehadiran", id)
}

func (h *SkpKehadiranHandlerImpl) SkpByKehadiranId(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	raw := c.Query("kehadiran_id")
	id := uuid.Must(uuid.FromString(raw))

	result, err := h.sk.SkpByKehadiranId(ctx, id)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get skp by kehadiran", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get skp by kehadiran", result)
}

func (h *SkpKehadiranHandlerImpl) IntervensiKehadiranId(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	raw := c.Query("kehadiran_id")
	id := uuid.Must(uuid.FromString(raw))

	result, err := h.sk.IntervensiByKehadiranId(ctx, id)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get intervensi by kehadiran", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get intervensi by kehadiran", result)
}

func (h *SkpKehadiranHandlerImpl) ApproveKehadiranSkp(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var p request.ApproveKehadiranSkp
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, _ := c.Get("nama")
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := h.sk.ApproveSkpKehadiran(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed approve kehadiran", err)
		return
	}

	resp.HandleSuccessResponse(c, "success approve kehadiran", res)
}
