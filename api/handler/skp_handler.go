package handler

import (
	"e-klinik/config"
	"e-klinik/utils"

	"e-klinik/internal/domain/resp"
	"e-klinik/internal/usecase"

	"time"

	"github.com/gin-gonic/gin"
)

type SkpHandler interface {
	ListIntervensi(c *gin.Context)
}

type SkpHandlerImpl struct {
	cfg *config.Config
	su  usecase.SkpUsecase
}

func NewSkpHandler(su usecase.SkpUsecase, cfg *config.Config) *SkpHandlerImpl {
	return &SkpHandlerImpl{
		cfg: cfg,
		su:  su,
	}
}

func (h *SkpHandlerImpl) ListIntervensi(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	result, err := h.su.ListIntervensi(ctx)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get intervensi", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get intervensi", result)
}
