package handler

import (
	"e-klinik/config"
	"e-klinik/infra/pg"
	"e-klinik/pkg"
	"e-klinik/utils"
	"strconv"
	"strings"

	"e-klinik/internal/domain/resp"
	"e-klinik/internal/usecase"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type ActorHandler interface {
	GetUsersByRoles(c *gin.Context)
	CreatePembimbingKlinik(c *gin.Context)
	ListPembimbingKlinikByKontrak(c *gin.Context)
}

type ActorHandlerImpl struct {
	cfg *config.Config
	au  usecase.ActorUsecase
}

func NewActorHandler(au usecase.ActorUsecase, cfg *config.Config) *ActorHandlerImpl {
	return &ActorHandlerImpl{
		cfg: cfg,
		au:  au,
	}
}

func (h *ActorHandlerImpl) GetUsersByRoles(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)

	defer cancel()

	raw := c.Query("role_ids")
	var roleIDs []int32
	for _, p := range strings.Split(raw, ",") {
		if n, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			roleIDs = append(roleIDs, int32(n))
		}
	}

	result, err := h.au.GetUsersByRoles(ctx, roleIDs)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get users by roles", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get users by roles", result)
}

func (h *ActorHandlerImpl) CreatePembimbingKlinik(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreatePembimbingKlinikParams
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, _ := c.Get("nama")
	p.CreatedBy = utils.StringPtr(value.(string))

	result, err := h.au.AddPembimbingKlinik(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed create pembimbing klinik", err)
		return
	}

	resp.HandleSuccessResponse(c, "success create pembimbing klinik", result)
}

func (h *ActorHandlerImpl) ListPembimbingKlinikByKontrak(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout(c, 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		resp.HandleErrorResponse(c, "detail failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid id"))
		return
	}

	uuid := uuid.Must(uuid.FromString(id))
	result, err := h.au.ListPembimbingKlinikByKontrak(ctx, uuid)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get pembimbing klinik", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get pembimbing klinik", result)
}
