package usecase

import (
	"context"
	"e-klinik/infra/pg"
	"e-klinik/infra/worker"
	"e-klinik/internal/domain/request"
	"e-klinik/internal/domain/resp"
	"e-klinik/pkg"
	"e-klinik/utils"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jinzhu/copier"
)

type MainUsecase interface {
	AddFasilitasKesehatan(c context.Context, arg pg.CreateFasilitasKesehatanParams) (any, error)
	ListFasilitasKesehatan(c context.Context, arg request.SearchFasilitasKesehatan) (any, error)
	UpdateFasilitasKesehatan(c context.Context, arg pg.UpdateFasilitasKesehatanPartialParams) (any, error)
	DeleteFasilitasKesehatan(c context.Context, arg pg.DeleteFasilitasKesehatanParams) error
	AddMataKuliah(c context.Context, arg pg.CreateMataKuliahParams) (any, error)
	ListMataKuliah(c context.Context, arg request.SearchMataKuliah) (any, error)
	UpdateMataKuliah(c context.Context, arg pg.UpdateMataKuliahParams) (any, error)
	DeleteMataKuliah(c context.Context, arg pg.DeleteMataKuliahParams) error
	AddKontrak(c context.Context, arg request.CreateKontrak) (any, error)
	ListKontrak(c context.Context, arg request.SearchKontrak) (any, error)
	UpdateKontrak(c context.Context, arg pg.UpdateKontrakPartialParams) (any, error)
	DeleteKontrak(c context.Context, arg pg.DeleteKontrakParams) error
	AddRuangan(c context.Context, arg pg.CreateRuanganParams) (any, error)
	ListRuangan(c context.Context, arg request.SearchRuangan) (any, error)
	UpdateRuangan(c context.Context, arg pg.UpdateRuanganPartialParams) (any, error)
	DeleteRuangan(c context.Context, arg pg.DeleteRuanganParams) error
	AddKehadiran(c context.Context, arg pg.CreateKehadiranParams) (any, error)
	ListKehadiran(c context.Context, arg request.SearchKehadiran) (any, error)
	UpdateKehadiran(c context.Context, arg pg.UpdateKehadiranPartialParams) (any, error)
	DeleteKehadiran(c context.Context, arg pg.DeleteKehadiranParams) error
}

type MainUsecaseImpl struct {
	worker *worker.ProducerService
	db     *pg.Queries
	pg     *pkg.Postgres
}

func NewMainUsecase(postgre *pkg.Postgres, worker *worker.ProducerService) *MainUsecaseImpl {
	return &MainUsecaseImpl{
		db:     pg.New(postgre.Pool),
		pg:     postgre,
		worker: worker,
	}
}

func (mu *MainUsecaseImpl) AddFasilitasKesehatan(c context.Context, arg pg.CreateFasilitasKesehatanParams) (any, error) {
	res, err := mu.db.CreateFasilitasKesehatan(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed create fasilitas kesehatan")
	}
	return res, nil
}

func (mu *MainUsecaseImpl) ListFasilitasKesehatan(c context.Context, arg request.SearchFasilitasKesehatan) (any, error) {
	if arg.Limit <= 0 {
		arg.Limit = 10
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	arg.Offset = utils.GetOffset(arg.Page, arg.Limit)

	var params pg.ListFasilitasKesehatanParams
	err := copier.Copy(&params, &arg)
	if err != nil {
		return nil, err
	}

	res, err := mu.db.ListFasilitasKesehatan(c, params)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get falitas kesehatan")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}
	var cparams pg.CountFasilitasKesehatanParams
	err = copier.Copy(&cparams, &arg)
	if err != nil {
		return nil, err
	}
	count, err := mu.db.CountFasilitasKesehatan(c, cparams)
	if err != nil {
		return nil, err
	}
	return resp.WithPaginate(res, resp.CalculatePagination(arg.Page, arg.Limit, count)), nil
}

func (mu *MainUsecaseImpl) UpdateFasilitasKesehatan(c context.Context, arg pg.UpdateFasilitasKesehatanPartialParams) (any, error) {
	res, err := mu.db.UpdateFasilitasKesehatanPartial(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed update fasilitas kesehatan")
	}
	return res, nil
}

func (mu *MainUsecaseImpl) DeleteFasilitasKesehatan(c context.Context, arg pg.DeleteFasilitasKesehatanParams) error {
	err := mu.db.DeleteFasilitasKesehatan(c, arg)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed delete fasilitas kesehatan")
	}
	return nil
}

func (mu *MainUsecaseImpl) AddMataKuliah(c context.Context, arg pg.CreateMataKuliahParams) (any, error) {
	res, err := mu.db.CreateMataKuliah(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed create matakuliah")
	}
	return res, nil
}

func (mu *MainUsecaseImpl) ListMataKuliah(c context.Context, arg request.SearchMataKuliah) (any, error) {
	if arg.Limit <= 0 {
		arg.Limit = 10
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	arg.Offset = utils.GetOffset(arg.Page, arg.Limit)

	var params pg.ListMataKuliahParams
	err := copier.Copy(&params, &arg)
	if err != nil {
		return nil, err
	}

	res, err := mu.db.ListMataKuliah(c, params)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get matakuliah")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}
	var cparams pg.CountMataKuliahParams
	err = copier.Copy(&cparams, &arg)
	if err != nil {
		return nil, err
	}
	count, err := mu.db.CountMataKuliah(c, cparams)
	if err != nil {
		return nil, err
	}
	return resp.WithPaginate(res, resp.CalculatePagination(arg.Page, arg.Limit, count)), nil
}

func (mu *MainUsecaseImpl) UpdateMataKuliah(c context.Context, arg pg.UpdateMataKuliahParams) (any, error) {
	res, err := mu.db.UpdateMataKuliah(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed update matakuliah")
	}
	return res, nil
}

func (mu *MainUsecaseImpl) DeleteMataKuliah(c context.Context, arg pg.DeleteMataKuliahParams) error {
	err := mu.db.DeleteMataKuliah(c, arg)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed delete matakuliah")
	}
	return nil
}

func (mu *MainUsecaseImpl) AddKontrak(c context.Context, arg request.CreateKontrak) (any, error) {
	var params pg.CreateKontrakParams
	err := copier.Copy(&params, &arg)
	if err != nil {
		return nil, err
	}
	if arg.PeriodeMulai != nil {
		mulai, err := time.Parse(time.RFC3339, *arg.PeriodeMulai)
		if err != nil {
			return nil, fmt.Errorf("periode mulai invalid")
		}
		params.PeriodeMulai = pgtype.Timestamptz{Time: mulai, Valid: true}
	}

	if arg.PeriodeSelesai != nil {
		selesai, err := time.Parse(time.RFC3339, *arg.PeriodeSelesai)
		if err != nil {
			return nil, fmt.Errorf("periode selesai invalid")
		}
		params.PeriodeSelesai = pgtype.Timestamptz{Time: selesai, Valid: true}
	}

	res, err := mu.db.CreateKontrak(c, params)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed create kontrak")
	}
	return res, nil
}

func (mu *MainUsecaseImpl) ListKontrak(c context.Context, arg request.SearchKontrak) (any, error) {
	if arg.Limit <= 0 {
		arg.Limit = 10
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	arg.Offset = utils.GetOffset(arg.Page, arg.Limit)

	var params pg.ListKontrakParams
	err := copier.Copy(&params, &arg)
	if err != nil {
		return nil, err
	}
	// ✅ Handle manual conversion for periode fields
	params.PeriodeMulai = utils.StringToTimestamptz(arg.PeriodeMulai)
	params.PeriodeSelesai = utils.StringToTimestamptz(arg.PeriodeSelesai)

	res, err := mu.db.ListKontrak(c, params)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get kontrak")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}
	var cparams pg.CountKontrakParams
	err = copier.Copy(&cparams, &arg)
	if err != nil {
		return nil, err
	}
	count, err := mu.db.CountKontrak(c, cparams)
	if err != nil {
		return nil, err
	}
	return resp.WithPaginate(res, resp.CalculatePagination(arg.Page, arg.Limit, count)), nil
}

func (mu *MainUsecaseImpl) UpdateKontrak(c context.Context, arg pg.UpdateKontrakPartialParams) (any, error) {
	res, err := mu.db.UpdateKontrakPartial(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed update kontrak")
	}
	return res, nil
}

func (mu *MainUsecaseImpl) DeleteKontrak(c context.Context, arg pg.DeleteKontrakParams) error {
	err := mu.db.DeleteKontrak(c, arg)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed delete kontrak")
	}
	return nil
}

func (mu *MainUsecaseImpl) AddRuangan(c context.Context, arg pg.CreateRuanganParams) (any, error) {
	res, err := mu.db.CreateRuangan(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed create ruangan")
	}
	return res, nil
}

func (mu *MainUsecaseImpl) ListRuangan(c context.Context, arg request.SearchRuangan) (any, error) {
	if arg.Limit <= 0 {
		arg.Limit = 10
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	arg.Offset = utils.GetOffset(arg.Page, arg.Limit)

	var params pg.ListRuanganParams
	err := copier.Copy(&params, &arg)
	if err != nil {
		return nil, err
	}

	if arg.FasilitasID != nil && *arg.FasilitasID != "" {
		fid := uuid.FromStringOrNil(*arg.FasilitasID)
		params.FasilitasID = &fid
	}
	if arg.KontrakID != nil && *arg.KontrakID != "" {
		kid := uuid.FromStringOrNil(*arg.KontrakID)
		params.KontrakID = &kid
	}

	res, err := mu.db.ListRuangan(c, params)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get ruangan")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}
	var cparams pg.CountRuanganParams
	err = copier.Copy(&cparams, &arg)
	if err != nil {
		return nil, err
	}
	count, err := mu.db.CountRuangan(c, cparams)
	if err != nil {
		return nil, err
	}
	return resp.WithPaginate(res, resp.CalculatePagination(arg.Page, arg.Limit, count)), nil
}

func (mu *MainUsecaseImpl) UpdateRuangan(c context.Context, arg pg.UpdateRuanganPartialParams) (any, error) {
	res, err := mu.db.UpdateRuanganPartial(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed update ruangan")
	}
	return res, nil
}

func (mu *MainUsecaseImpl) DeleteRuangan(c context.Context, arg pg.DeleteRuanganParams) error {
	err := mu.db.DeleteRuangan(c, arg)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed delete ruangan")
	}
	return nil
}

func (mu *MainUsecaseImpl) AddKehadiran(c context.Context, arg pg.CreateKehadiranParams) (any, error) {
	res, err := mu.db.CreateKehadiran(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed create kehadiran")
	}
	return res, nil
}

func (mu *MainUsecaseImpl) ListKehadiran(c context.Context, arg request.SearchKehadiran) (any, error) {
	if arg.Limit <= 0 {
		arg.Limit = 10
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	arg.Offset = utils.GetOffset(arg.Page, arg.Limit)

	var params pg.ListKehadiranParams
	err := copier.Copy(&params, &arg)
	if err != nil {
		return nil, err
	}

	if arg.FasilitasID != nil && *arg.FasilitasID != "" {
		fid := uuid.FromStringOrNil(*arg.FasilitasID)
		params.FasilitasID = &fid
	}
	if arg.KontrakID != nil && *arg.KontrakID != "" {
		kid := uuid.FromStringOrNil(*arg.KontrakID)
		params.KontrakID = &kid
	}

	res, err := mu.db.ListKehadiran(c, params)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed get kehadiran")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}
	var cparams pg.CountKehadiranParams
	err = copier.Copy(&cparams, &arg)
	if err != nil {
		return nil, err
	}
	count, err := mu.db.CountKehadiran(c, cparams)
	if err != nil {
		return nil, err
	}
	return resp.WithPaginate(res, resp.CalculatePagination(arg.Page, arg.Limit, count)), nil
}

func (mu *MainUsecaseImpl) UpdateKehadiran(c context.Context, arg pg.UpdateKehadiranPartialParams) (any, error) {
	res, err := mu.db.UpdateKehadiranPartial(c, arg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed update kehadiran")
	}
	return res, nil
}

func (mu *MainUsecaseImpl) DeleteKehadiran(c context.Context, arg pg.DeleteKehadiranParams) error {
	err := mu.db.DeleteKehadiran(c, arg)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "failed delete kehadiran")
	}
	return nil
}
