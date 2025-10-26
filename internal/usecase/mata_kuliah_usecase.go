package usecase

import (
	"context"
	"e-klinik/infra/pg"
	"e-klinik/infra/worker"
	"e-klinik/internal/domain/request"
	"e-klinik/internal/domain/resp"
	"e-klinik/pkg"
	"e-klinik/utils"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jinzhu/copier"
)

type MataKuliahUsecase interface {
	AddMataKuliah(c context.Context, arg pg.CreateMataKuliahParams) (any, error)
	ListMataKuliah(c context.Context, arg request.SearchMataKuliah) (any, error)
	UpdateMataKuliah(c context.Context, arg pg.UpdateMataKuliahParams) (any, error)
	DeleteMataKuliah(c context.Context, arg pg.DeleteMataKuliahParams) error
}

type MataKuliahUsecaseImpl struct {
	worker *worker.ProducerService
	db     *pg.Queries
	pg     *pkg.Postgres
	cache  *pkg.RedisCache
}

func NewMataKuliahUsecase(postgre *pkg.Postgres, worker *worker.ProducerService, cache *pkg.RedisCache) *MataKuliahUsecaseImpl {
	return &MataKuliahUsecaseImpl{
		db:     pg.New(postgre.Pool),
		pg:     postgre,
		worker: worker,
		cache:  cache,
	}
}

func (mu *MataKuliahUsecaseImpl) AddMataKuliah(c context.Context, arg pg.CreateMataKuliahParams) (any, error) {
	res, err := mu.db.CreateMataKuliah(c, arg)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to create mata kuliah")
	}
	return res, nil
}

func (mu *MataKuliahUsecaseImpl) ListMataKuliah(c context.Context, arg request.SearchMataKuliah) (any, error) {
	if arg.Limit <= 0 {
		arg.Limit = 10
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	arg.Offset = utils.GetOffset(arg.Page, arg.Limit)

	var params pg.ListMataKuliahParams
	if err := copier.Copy(&params, &arg); err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to copy list params")
	}

	res, err := mu.db.ListMataKuliah(c, params)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to get mata kuliah list")
	}

	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), nil
	}

	var cparams pg.CountMataKuliahParams
	if err := copier.Copy(&cparams, &arg); err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to copy count params")
	}

	count, err := mu.db.CountMataKuliah(c, cparams)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to count mata kuliah")
	}

	return resp.WithPaginate(res, resp.CalculatePagination(arg.Page, arg.Limit, count)), nil
}

func (mu *MataKuliahUsecaseImpl) UpdateMataKuliah(c context.Context, arg pg.UpdateMataKuliahParams) (any, error) {
	res, err := mu.db.UpdateMataKuliah(c, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkg.ExposeError(pkg.ErrorCodeNotFound, "mata kuliah not found")
		}
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed update mata kuliah")
	}
	return res, nil
}

func (mu *MataKuliahUsecaseImpl) DeleteMataKuliah(c context.Context, arg pg.DeleteMataKuliahParams) error {
	err := mu.db.DeleteMataKuliah(c, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pkg.ExposeError(pkg.ErrorCodeNotFound, "mata kuliah not found")
		}
		return pkg.WrapError(err, pkg.ErrorCodeInternal, "failed delete mata kuliah")
	}
	return nil
}
