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

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type MainHandler interface {
	CreateFasilitasKesehatan(c *gin.Context)
	ListFaslitasKesehatan(c *gin.Context)
	UpdateFasilitasKesehatan(c *gin.Context)
	DelFasilitasKesehatan(c *gin.Context)
	CreateMatakuliah(c *gin.Context)
	ListMataKuliah(c *gin.Context)
	UpdateMataKuliah(c *gin.Context)
	DelMataKuliah(c *gin.Context)
	CreateKontrak(c *gin.Context)
	ListKontrak(c *gin.Context)
	UpdateKontrak(c *gin.Context)
	DelKontrak(c *gin.Context)
	CreateRuangan(c *gin.Context)
	ListRuangan(c *gin.Context)
	UpdateRuangan(c *gin.Context)
	DelRuangan(c *gin.Context)
	CreateKehadiran(c *gin.Context)
	ListKehadiran(c *gin.Context)
	UpdateKehadiran(c *gin.Context)
	DeleteKehadiran(c *gin.Context)
}

type MainHandlerImpl struct {
	cfg         *config.Config
	mainUsecase usecase.MainUsecase
}

func NewMainHandler(mainUsecase usecase.MainUsecase, cfg *config.Config) *MainHandlerImpl {
	return &MainHandlerImpl{
		cfg:         cfg,
		mainUsecase: mainUsecase,
	}
}

func (h *MainHandlerImpl) CreateFasilitasKesehatan(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreateFasilitasKesehatanParams
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created fasilitas kesehatan", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.CreatedBy = utils.StringPtr(value.(string))

	result, err := h.mainUsecase.AddFasilitasKesehatan(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created fasilitas kesehatan", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))

}

func (h *MainHandlerImpl) ListFaslitasKesehatan(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchFasilitasKesehatan

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	result, err := h.mainUsecase.ListFasilitasKesehatan(ctx, req)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed get fasilitas kesehatan", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, 0))

}

func (h *MainHandlerImpl) UpdateFasilitasKesehatan(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.UpdateFasilitasKesehatanPartialParams
	p.ID = uuid.Must(uuid.FromString(c.Query("id")))
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update fasilitas kesehatan", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := h.mainUsecase.UpdateFasilitasKesehatan(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update fasilitas kesehatan", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (h *MainHandlerImpl) DelFasilitasKesehatan(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Query("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "delete fasilitas kesehatan failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}
	p := pg.DeleteFasilitasKesehatanParams{
		ID:        uuid.Must(uuid.FromString(id)),
		DeletedBy: nil}
	value, _ := c.Get("nama")
	p.DeletedBy = utils.StringPtr(value.(string))

	err := h.mainUsecase.DeleteFasilitasKesehatan(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "delete fasilitas kesehatan failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(id, true, 0))
}

func (h *MainHandlerImpl) CreateMatakuliah(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreateMataKuliahParams
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created mata kuliah", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.CreatedBy = utils.StringPtr(value.(string))

	result, err := h.mainUsecase.AddMataKuliah(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created matakuliah", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))

}
func (h *MainHandlerImpl) ListMataKuliah(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchMataKuliah

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	result, err := h.mainUsecase.ListMataKuliah(ctx, req)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed get mata kulah", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, 0))

}

func (h *MainHandlerImpl) UpdateMataKuliah(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.UpdateMataKuliahParams
	p.ID = uuid.Must(uuid.FromString(c.Query("id")))
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update mata kuliah", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := h.mainUsecase.UpdateMataKuliah(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update mata kuliah", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (h *MainHandlerImpl) DelMataKuliah(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Query("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "delete mata kuliah failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}
	p := pg.DeleteMataKuliahParams{
		ID:        uuid.Must(uuid.FromString(id)),
		DeletedBy: nil}
	value, _ := c.Get("nama")
	p.DeletedBy = utils.StringPtr(value.(string))

	err := h.mainUsecase.DeleteMataKuliah(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "delete mata kuliah failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(id, true, 0))
}

func (h *MainHandlerImpl) CreateKontrak(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p request.CreateKontrak
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created kontrak", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.CreatedBy = utils.StringPtr(value.(string))

	result, err := h.mainUsecase.AddKontrak(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created kontrak", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))

}
func (h *MainHandlerImpl) ListKontrak(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchKontrak

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	result, err := h.mainUsecase.ListKontrak(ctx, req)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed get kontrak", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, 0))

}

func (h *MainHandlerImpl) UpdateKontrak(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.UpdateKontrakPartialParams
	p.ID = uuid.Must(uuid.FromString(c.Query("id")))
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update kontrak", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := h.mainUsecase.UpdateKontrak(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update kontrak", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (h *MainHandlerImpl) DelKontrak(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Query("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "delete kontrak failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}
	p := pg.DeleteKontrakParams{
		ID:        uuid.Must(uuid.FromString(id)),
		DeletedBy: nil}
	value, _ := c.Get("nama")
	p.DeletedBy = utils.StringPtr(value.(string))

	err := h.mainUsecase.DeleteKontrak(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "delete kontrak failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(id, true, 0))
}

func (h *MainHandlerImpl) CreateRuangan(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreateRuanganParams
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created ruangan", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.CreatedBy = utils.StringPtr(value.(string))

	result, err := h.mainUsecase.AddRuangan(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created ruangan", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))

}
func (h *MainHandlerImpl) ListRuangan(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchRuangan

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	result, err := h.mainUsecase.ListRuangan(ctx, req)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed get ruangan", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, 0))

}

func (h *MainHandlerImpl) UpdateRuangan(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.UpdateRuanganPartialParams
	p.ID = uuid.Must(uuid.FromString(c.Query("id")))
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update ruangan", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := h.mainUsecase.UpdateRuangan(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update ruangan", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (h *MainHandlerImpl) DelRuangan(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Query("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "delete ruangan failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}
	p := pg.DeleteRuanganParams{
		ID:        uuid.Must(uuid.FromString(id)),
		DeletedBy: nil}
	value, _ := c.Get("nama")
	p.DeletedBy = utils.StringPtr(value.(string))

	err := h.mainUsecase.DeleteRuangan(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "delete ruangan failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(id, true, 0))
}

func (h *MainHandlerImpl) CreateKehadiran(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreateKehadiranParams
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created kehadiran", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.CreatedBy = utils.StringPtr(value.(string))

	result, err := h.mainUsecase.AddKehadiran(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created kehadiran", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))

}
func (h *MainHandlerImpl) ListKehadiran(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchKehadiran

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	result, err := h.mainUsecase.ListKehadiran(ctx, req)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed get kehadiran", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, 0))

}

func (h *MainHandlerImpl) UpdateKehadiran(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.UpdateKehadiranPartialParams
	p.ID = uuid.Must(uuid.FromString(c.Query("id")))
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update kehadiran", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := h.mainUsecase.UpdateKehadiran(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update kehadiran", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (h *MainHandlerImpl) DeleteKehadiran(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Query("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "delete kehadiran failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}
	p := pg.DeleteKehadiranParams{
		ID:        uuid.Must(uuid.FromString(id)),
		DeletedBy: nil}
	value, _ := c.Get("nama")
	p.DeletedBy = utils.StringPtr(value.(string))

	err := h.mainUsecase.DeleteKehadiran(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "delete kehadiran failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(id, true, 0))
}

func (h *MainHandlerImpl) CreateKehadiranSkp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.CreateKehadiranSkpParams
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created kehadiran", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.CreatedBy = utils.StringPtr(value.(string))

	result, err := h.mainUsecase.AddSkpKehadiran(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed created kehadiran", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))

}
func (h *MainHandlerImpl) ListKehadiranSkp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var req request.SearchKehadiranSkp

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	result, err := h.mainUsecase.ListKehadiranSkp(ctx, req)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed get kehadiran", err)
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, 0))

}

func (h *MainHandlerImpl) UpdateKehadiranSkp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var p pg.UpdateKehadiranSkpParams
	p.ID = uuid.Must(uuid.FromString(c.Query("id")))
	err := c.ShouldBindJSON(&p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update kehadiran", pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "invalid json"))
		return
	}
	value, _ := c.Get("nama")
	p.UpdatedBy = utils.StringPtr(value.(string))

	res, err := h.mainUsecase.UpdateKehadiranSkp(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "failed update kehadiran", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(res, true, 0))
}

func (h *MainHandlerImpl) DeleteKehadiranSkp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	id := c.Query("id")
	if id == "" {
		resp.GenerateBaseResponseWithError(c, "delete kehadiran failed", pkg.NewErrorf(pkg.ErrorCodeInvalidArgument, "invalid id"))

	}
	p := pg.DeleteKehadiranSkpParams{
		ID:        uuid.Must(uuid.FromString(id)),
		DeletedBy: nil}
	value, _ := c.Get("nama")
	p.DeletedBy = utils.StringPtr(value.(string))

	err := h.mainUsecase.DeleteKehadiranSkp(ctx, p)
	if err != nil {
		resp.GenerateBaseResponseWithError(c, "delete kehadiran failed", err)
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(id, true, 0))
}
