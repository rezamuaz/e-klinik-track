package usecase

import (
	"context"
	"e-klinik/infra/pg"
	"e-klinik/infra/worker"
	"e-klinik/internal/domain/resp"
	"e-klinik/pkg"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type ActorUsecase interface {
	GetUsersByRoles(c context.Context, arg []int32) (any, error)
	AddPembimbingKlinik(c context.Context, arg pg.CreatePembimbingKlinikParams) (any, error)
	ListPembimbingKlinikByKontrak(c context.Context, arg uuid.UUID) (any, error)
}

type ActorUsecaseImpl struct {
	worker *worker.ProducerService
	db     *pg.Queries
	pg     *pkg.Postgres
	cache  *pkg.RedisCache
}

func NewActorUsecase(postgre *pkg.Postgres, worker *worker.ProducerService, cache *pkg.RedisCache) *ActorUsecaseImpl {
	return &ActorUsecaseImpl{
		db:     pg.New(postgre.Pool),
		pg:     postgre,
		worker: worker,
		cache:  cache,
	}
}

func (mu *ActorUsecaseImpl) GetUsersByRoles(c context.Context, arg []int32) (any, error) {

	// var err error
	res, err := mu.db.GetUsersByRoles(c, arg)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get data")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]string{}, nil), err
	}

	return resp.WithPaginate(res, nil), nil

}

func (mu *ActorUsecaseImpl) AddPembimbingKlinik(c context.Context, arg pg.CreatePembimbingKlinikParams) (any, error) {
	kontrak, err := mu.db.GetKontrakByID(c, arg.KontrakID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if pgErr.ConstraintName == "pembimbing_klinik_kontrak_user_ukey" {
					return nil, pkg.ExposeError(pkg.ErrorCodeConflict, "kombinasi kontrak dan user sudah terdaftar")
				}
			case pgerrcode.NotNullViolation:
				return nil, pkg.ExposeError(pkg.ErrorCodeBadRequest, "data wajib diisi")
			}
		}
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed create pembimbing klinik")
	}

	arg.FasilitasID = kontrak.FasilitasID
	res, err := mu.db.CreatePembimbingKlinik(c, arg)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed create pembimbing klinik")
	}

	return res, nil
}

func (mu *ActorUsecaseImpl) ListPembimbingKlinikByKontrak(c context.Context, arg uuid.UUID) (any, error) {
	res, err := mu.db.ListPembimbingKlinikByKontrakID(c, arg)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get pembimbing klinik")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]any{}, nil), nil
	}
	return resp.WithPaginate(res, nil), nil
}
