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

type PermissionHandler interface {
	ListPermission(c *gin.Context)
	PermissionByRoleId(c *gin.Context)
	UserViewPermission(c *gin.Context)
	AddMenu(c *gin.Context)
	ListMenu(c *gin.Context)
	UpdateMenu(c *gin.Context)
	DelMenu(c *gin.Context)
	MenuDetail(c *gin.Context)
}

type PermissionHandlerImpl struct {
	Uu *usecase.UserUsecaseImpl

	Cfg *config.Config
}

func NewPermissionHandler(Uu *usecase.UserUsecaseImpl, cfg *config.Config) *PermissionHandlerImpl {
	return &PermissionHandlerImpl{
		Uu:  Uu,
		Cfg: cfg,
	}
}

func (lc *PermissionHandlerImpl) ListPermission(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.SearchMenu
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	res, err := lc.Uu.ListPermission(ctx, req)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get list permission", err)
		return
	}

	resp.HandleSuccessResponse(c, "success get list permission", res)
}
func (lc *PermissionHandlerImpl) PermissionByRoleId(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		resp.HandleErrorResponse(c, "detail failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid id"))
		return
	}
	p := utils.StrToInt32(id)
	if p == 0 {
		resp.HandleErrorResponse(c, "detail failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "id tidak valid"))
		return
	}

	res, err := lc.Uu.GetViewByRoleId(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "detail failed", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}
	resp.HandleSuccessResponse(c, "success get view", res)
}

func (lc *PermissionHandlerImpl) UserViewPermission(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	value, ok := c.Get("Id")
	if !ok {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user id context not found"))
		return
	}
	idStr, ok := value.(string)
	if !ok {
		resp.HandleErrorResponse(c, "context invalid", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user id context invalid"))
		return
	}
	id, err := uuid.FromString(idStr)
	if err != nil {
		resp.HandleErrorResponse(c, "context invalid", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid uuid in context"))
		return
	}

	res, err := lc.Uu.UserViewPermission(ctx, id)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get view user", pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get view user"))
		return
	}
	resp.HandleSuccessResponse(c, "success get view", res)
}
func (lc *PermissionHandlerImpl) AddMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var p pg.CreateR1ViewParams
	if err := c.ShouldBind(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	res, err := lc.Uu.AddMenu(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed to add menu", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success add menu", res)
}

func (lc *PermissionHandlerImpl) ListMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.SearchMenu
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.HandleErrorResponse(c, "invalid query parameters", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid query parameters"))
		return
	}

	res, err := lc.Uu.ListMenu(ctx, req)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get list menu", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}

	resp.HandleSuccessResponse(c, "success get list menu", res)
}

func (lc *PermissionHandlerImpl) UpdateMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	var p pg.UpdateR1ViewParams
	p.ID = utils.StrToInt32(id)

	if err := c.ShouldBind(&p); err != nil {
		resp.HandleErrorResponse(c, "failed to bind JSON", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid JSON payload"))
		return
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak ditemukan"))
		return
	}
	if v, ok := value.(string); ok {
		p.UpdatedBy = utils.StringPtr(v)
	} else {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak valid"))
		return
	}

	res, err := lc.Uu.EditMenu(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed edit menu", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}
	resp.HandleSuccessResponse(c, "success update menu", res)
}

func (lc *PermissionHandlerImpl) DelMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Query("id")
	if id == "" {
		resp.HandleErrorResponse(c, "delete menu failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid id"))
		return
	}
	p := pg.DeleteR1ViewParams{
		ID:        utils.StrToInt32(id),
		DeletedBy: nil,
	}

	value, ok := c.Get("nama")
	if !ok {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak ditemukan"))
		return
	}
	if v, ok := value.(string); ok {
		p.DeletedBy = utils.StringPtr(v)
	} else {
		resp.HandleErrorResponse(c, "context missing", pkg.ExposeError(pkg.ErrorCodeUnauthorized, "user context tidak valid"))
		return
	}

	if err := lc.Uu.DeleteMenu(ctx, p); err != nil {
		resp.HandleErrorResponse(c, "failed delete menu", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}
	resp.HandleSuccessResponse(c, "success delete menu", nil)
}

func (lc *PermissionHandlerImpl) MenuDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		resp.HandleErrorResponse(c, "detail failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "invalid id"))
		return
	}
	p := utils.StrToInt32(id)
	if p == 0 {
		resp.HandleErrorResponse(c, "detail failed", pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "id tidak valid"))
		return
	}

	res, err := lc.Uu.MenuById(ctx, p)
	if err != nil {
		resp.HandleErrorResponse(c, "failed get menu detail", pkg.WrapError(err, pkg.ErrorCodeInternal, "internal server error"))
		return
	}
	resp.HandleSuccessResponse(c, "success get menu detail", res)
}
