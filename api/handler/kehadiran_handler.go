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

type KehadiranHandler interface {
	CreateKehadiran(c *gin.Context)
	ListKehadiran(c *gin.Context)
	UpdateKehadiran(c *gin.Context)
	DeleteKehadiran(c *gin.Context)
	CheckKehadiran(c *gin.Context)
	GetKehadiranByPembimbingStatus(c *gin.Context)
	GetKehadiranByMahasiswaStatus(c *gin.Context)
	ListDistinctUserKehadiran(c *gin.Context)
}

type KehadiranHandlerImpl struct {
	cfg *config.Config
	ku  usecase.KehadiranUsecase
}

func NewKehadiranHandler(ku usecase.KehadiranUsecase, cfg *config.Config) *KehadiranHandlerImpl {
	return &KehadiranHandlerImpl{
		cfg: cfg,
		ku:  ku,
	}
}

func (h *KehadiranHandlerImpl) CreateKehadiran(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreateKehadiranParams
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to create kehadiran", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, vok := c.Get("nama")
	idVal, iok := c.Get("Id")
	if !vok || !iok {
		resp.HandleErrorResponse(c, "failed to create kehadiran", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}

	p.CreatedBy = utils.StringPtr(value.(string))
	p.UserID = uuid.Must(uuid.FromString(idVal.(string)))

	result, err := h.ku.AddKehadiran(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to create kehadiran", pkg.WrapError(err, pkg.ErrorCodeInternal, "database error"))
		return
	}

	resp.HandleSuccessResponse(c, "success create kehadiran", result)
}

func (h *KehadiranHandlerImpl) ListKehadiran(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchKehadiran
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	result, err := h.ku.ListKehadiran(ctx, req)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to get kehadiran list", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success get kehadiran list", result)
}

func (h *KehadiranHandlerImpl) UpdateKehadiran(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	idStr := c.Query("id")
	if idStr == "" {
		resp.HandleErrorResponse(c, "invalid request", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "missing kehadiran id"))
		return
	}

	var p pg.UpdateKehadiranPartialParams
	p.ID = uuid.Must(uuid.FromString(idStr))

	if err := c.ShouldBindJSON(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "failed to update kehadiran", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := h.ku.UpdateKehadiran(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to update kehadiran", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success update kehadiran", res)
}

func (h *KehadiranHandlerImpl) DeleteKehadiran(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	idStr := c.Query("id")
	if idStr == "" {
		resp.HandleErrorResponse(c, "invalid request", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "missing kehadiran id"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "failed to delete kehadiran", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context not found"))
		return
	}

	p := pg.DeleteKehadiranParams{
		ID:        uuid.Must(uuid.FromString(idStr)),
		DeletedBy: utils.StringPtr(value.(string)),
	}

	if err := h.ku.DeleteKehadiran(ctx, p); err != nil {
		resp.HandleErrorResponse(c, "failed to delete kehadiran", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success delete kehadiran", gin.H{"id": idStr})
}
func (h *KehadiranHandlerImpl) CheckKehadiran(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	value, _ := c.Get("Id")
	id := uuid.Must(uuid.FromString(value.(string)))

	result, err := h.ku.CheckKehadiran(ctx, id)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to delete kehadiran", err)
		return
	}

	resp.HandleSuccessResponse(c, "success check kehadiran", result)

}
func (h *KehadiranHandlerImpl) GetKehadiranByPembimbingStatus(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var arg pg.GetKehadiranByPembimbingUserIdParams
	value, _ := c.Get("Id")
	id := uuid.Must(uuid.FromString(value.(string)))
	arg.PembimbingKlinik = &id

	result, err := h.ku.GetKehadiranByPembimbingStatus(ctx, arg)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get kehadiran by pembimbing", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get kehadiran by pembimbing", result)
}

func (h *KehadiranHandlerImpl) GetKehadiranByMahasiswaStatus(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var arg pg.GetKehadiranByPembimbingUserIdParams
	value, _ := c.Get("Id")
	id := uuid.Must(uuid.FromString(value.(string)))
	arg.UserID = &id

	result, err := h.ku.GetKehadiranByPembimbingStatus(ctx, arg)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get kehadiran by pembimbing", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get kehadiran by pembimbing", result)
}

func (h *KehadiranHandlerImpl) ListDistinctUserKehadiran(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchUserKehadiran
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	res, err := h.ku.ListDistinctUserKehadiran(ctx, req)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get user kehadiran", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get user kehadiran", res)
}
