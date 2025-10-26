package usecase

import (
	"context"
	"e-klinik/infra/pg"
	"e-klinik/infra/worker"
	"e-klinik/internal/domain/resp"
	"e-klinik/pkg"
)

type SkpUsecase interface {
	ListIntervensi(c context.Context) (any, error)
}

type SkpUsecaseImpl struct {
	worker *worker.ProducerService
	db     *pg.Queries
	pg     *pkg.Postgres
	cache  *pkg.RedisCache
}

func NewSkpUsecase(postgre *pkg.Postgres, worker *worker.ProducerService, cache *pkg.RedisCache) *SkpUsecaseImpl {
	return &SkpUsecaseImpl{
		db:     pg.New(postgre.Pool),
		pg:     postgre,
		worker: worker,
		cache:  cache,
	}
}

func (mu *SkpUsecaseImpl) ListIntervensi(c context.Context) (any, error) {

	res, err := mu.db.ListKategoriSubkategoriIntervensi(c)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get intervensi")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}

	return resp.WithPaginate(res, nil), nil
}
