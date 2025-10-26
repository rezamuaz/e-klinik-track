package usecase

import (
	"context"
	"e-klinik/infra/pg"
	"e-klinik/infra/worker"
	"e-klinik/internal/domain/request"
	"e-klinik/internal/domain/resp"
	"e-klinik/pkg"
	"e-klinik/utils"

	"github.com/gofrs/uuid/v5"

	"github.com/jinzhu/copier"
)

type FasilitasUsecase interface {
	AddFasilitasKesehatan(c context.Context, arg pg.CreateFasilitasKesehatanParams) (any, error)
	ListFasilitasKesehatan(c context.Context, arg request.SearchFasilitasKesehatan) (any, error)
	UpdateFasilitasKesehatan(c context.Context, arg pg.UpdateFasilitasKesehatanPartialParams) (any, error)
	DeleteFasilitasKesehatan(c context.Context, arg pg.DeleteFasilitasKesehatanParams) error
	ListPropinsi(c context.Context, arg request.SearchPropinsi) (any, error)
	ListKabupaten(c context.Context, arg request.SearchKabupaten) (any, error)
}

type FasilitasUsecaseImpl struct {
	worker *worker.ProducerService
	db     *pg.Queries
	pg     *pkg.Postgres
	cache  *pkg.RedisCache
}

func NewFasilitasUseCase(postgre *pkg.Postgres, worker *worker.ProducerService, cache *pkg.RedisCache) *FasilitasUsecaseImpl {
	return &FasilitasUsecaseImpl{
		db:     pg.New(postgre.Pool),
		pg:     postgre,
		worker: worker,
		cache:  cache,
	}
}

func (mu *FasilitasUsecaseImpl) AddFasilitasKesehatan(c context.Context, arg pg.CreateFasilitasKesehatanParams) (any, error) {
	res, err := mu.db.CreateFasilitasKesehatan(c, arg)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to create fasilitas kesehatan")
	}
	return res, nil
}

func (mu *FasilitasUsecaseImpl) ListFasilitasKesehatan(c context.Context, arg request.SearchFasilitasKesehatan) (any, error) {
	if arg.Limit <= 0 {
		arg.Limit = 10
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	arg.Offset = utils.GetOffset(arg.Page, arg.Limit)

	var params pg.ListFasilitasKesehatanParams
	if err := copier.Copy(&params, &arg); err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to copy list params")
	}

	if arg.KabID != nil && *arg.KabID != "" {
		kid := uuid.FromStringOrNil(*arg.KabID)
		params.KabID = &kid
	}

	res, err := mu.db.ListFasilitasKesehatan(c, params)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to get fasilitas kesehatan list")
	}

	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), nil
	}

	var cparams pg.CountFasilitasKesehatanParams
	if err := copier.Copy(&cparams, &arg); err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to copy count params")
	}

	count, err := mu.db.CountFasilitasKesehatan(c, cparams)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to count fasilitas kesehatan")
	}

	return resp.WithPaginate(res, resp.CalculatePagination(arg.Page, arg.Limit, count)), nil
}

func (mu *FasilitasUsecaseImpl) UpdateFasilitasKesehatan(c context.Context, arg pg.UpdateFasilitasKesehatanPartialParams) (any, error) {
	res, err := mu.db.UpdateFasilitasKesehatanPartial(c, arg)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to update fasilitas kesehatan")
	}
	return res, nil
}

func (mu *FasilitasUsecaseImpl) DeleteFasilitasKesehatan(c context.Context, arg pg.DeleteFasilitasKesehatanParams) error {
	if err := mu.db.DeleteFasilitasKesehatan(c, arg); err != nil {
		return pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to delete fasilitas kesehatan")
	}
	return nil
}
func (mu *FasilitasUsecaseImpl) ListPropinsi(c context.Context, arg request.SearchPropinsi) (any, error) {

	var params pg.ListDistinctPropinsiParams
	err := copier.Copy(&params, &arg)
	if err != nil {
		return nil, err
	}

	params.Limit = 100
	params.Offset = 0

	res, err := mu.db.ListDistinctPropinsi(c, params)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get propinsi")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}

	return resp.WithPaginate(res, nil), nil
}

func (mu *FasilitasUsecaseImpl) ListKabupaten(c context.Context, arg request.SearchKabupaten) (any, error) {

	var params pg.ListDistinctKabupatenParams
	err := copier.Copy(&params, &arg)
	if err != nil {
		return nil, err
	}

	if arg.PropinsiID != nil && *arg.PropinsiID != "" {
		kid := uuid.FromStringOrNil(*arg.PropinsiID)
		params.PropinsiID = &kid
	}

	params.Limit = 100
	params.Offset = 0

	res, err := mu.db.ListDistinctKabupaten(c, params)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get kabupaten")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}

	return resp.WithPaginate(res, nil), nil
}
