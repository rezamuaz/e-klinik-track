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

type RuanganUsecase interface {
	AddRuangan(c context.Context, arg pg.CreateRuanganParams) (any, error)
	ListRuangan(c context.Context, arg request.SearchRuangan) (any, error)
	RuanganById(c context.Context, arg uuid.UUID) (any, error)
	ListRuanganByKontrak(c context.Context, arg request.SearchRuanganByKontrak) (any, error)
	UpdateRuangan(c context.Context, arg pg.UpdateRuanganPartialParams) (any, error)
	DeleteRuangan(c context.Context, arg pg.DeleteRuanganParams) error
}

type RuanganUsecaseImpl struct {
	worker *worker.ProducerService
	db     *pg.Queries
	pg     *pkg.Postgres
	cache  *pkg.RedisCache
}

func NewRuanganUsecase(postgre *pkg.Postgres, worker *worker.ProducerService, cache *pkg.RedisCache) *RuanganUsecaseImpl {
	return &RuanganUsecaseImpl{
		db:     pg.New(postgre.Pool),
		pg:     postgre,
		worker: worker,
		cache:  cache,
	}
}

func (mu *RuanganUsecaseImpl) RuanganById(c context.Context, arg uuid.UUID) (any, error) {
	res, err := mu.db.GetRuanganById(c, arg)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get kontrak")
	}
	return resp.WithPaginate(res, nil), nil
}

func (mu *RuanganUsecaseImpl) AddRuangan(c context.Context, arg pg.CreateRuanganParams) (any, error) {
	res, err := mu.db.CreateRuangan(c, arg)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed create ruangan")
	}
	return res, nil
}

func (mu *RuanganUsecaseImpl) ListRuangan(c context.Context, arg request.SearchRuangan) (any, error) {
	if arg.Limit <= 0 {
		arg.Limit = 10
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	arg.Offset = utils.GetOffset(arg.Page, arg.Limit)

	var params pg.ListRuanganParams
	if err := copier.Copy(&params, &arg); err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to copy list params")
	}

	res, err := mu.db.ListRuangan(c, params)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get ruangan")
	}

	if len(res) == 0 {
		return resp.WithPaginate([]any{}, resp.CalculatePagination(arg.Page, arg.Limit, 0)), nil
	}

	var cparams pg.CountRuanganParams
	if err := copier.Copy(&cparams, &arg); err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to copy count params")
	}

	count, err := mu.db.CountRuangan(c, cparams)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to count ruangan")
	}

	return resp.WithPaginate(res, resp.CalculatePagination(arg.Page, arg.Limit, count)), nil
}

func (mu *RuanganUsecaseImpl) ListRuanganByKontrak(c context.Context, arg request.SearchRuanganByKontrak) (any, error) {
	var params pg.GetRuanganBYKontrakParams
	err := copier.Copy(&params, &arg)
	if err != nil {
		return nil, err
	}

	fid := uuid.FromStringOrNil(arg.FasilitasID)
	params.FasilitasID = fid

	kid := uuid.FromStringOrNil(arg.KontrakID)
	params.KontrakID = kid

	res, err := mu.db.GetRuanganBYKontrak(c, params)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get ruangan")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}

	return resp.WithPaginate(res, nil), nil
}

func (mu *RuanganUsecaseImpl) UpdateRuangan(c context.Context, arg pg.UpdateRuanganPartialParams) (any, error) {
	res, err := mu.db.UpdateRuanganPartial(c, arg)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed update ruangan")
	}
	return res, nil
}

func (mu *RuanganUsecaseImpl) DeleteRuangan(c context.Context, arg pg.DeleteRuanganParams) error {
	err := mu.db.DeleteRuangan(c, arg)
	if err != nil {
		return pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed delete ruangan")
	}
	return nil
}
