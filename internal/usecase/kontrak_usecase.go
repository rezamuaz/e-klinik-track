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

	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jinzhu/copier"
)

type KontrakUsecase interface {
	AddKontrak(c context.Context, arg request.CreateKontrak) (any, error)
	KontrakById(c context.Context, arg uuid.UUID) (any, error)
	ListKontrak(c context.Context, arg request.SearchKontrak) (any, error)
	ListAktifKontrak(c context.Context, arg *string) (any, error)
	UpdateKontrak(c context.Context, arg pg.UpdateKontrakPartialParams) (any, error)
	DeleteKontrak(c context.Context, arg pg.DeleteKontrakParams) error
}

type KontrakUsecaseImpl struct {
	worker *worker.ProducerService
	db     *pg.Queries
	pg     *pkg.Postgres
	cache  *pkg.RedisCache
}

func NewKontrakUsecase(postgre *pkg.Postgres, worker *worker.ProducerService, cache *pkg.RedisCache) *KontrakUsecaseImpl {
	return &KontrakUsecaseImpl{
		db:     pg.New(postgre.Pool),
		pg:     postgre,
		worker: worker,
		cache:  cache,
	}
}

func (mu *KontrakUsecaseImpl) AddKontrak(c context.Context, arg request.CreateKontrak) (any, error) {
	return utils.WithTransactionResult(c, mu.pg.Pool, func(qtx *pg.Queries, tx pgx.Tx) (any, error) {
		var params pg.CreateKontrakParams

		if err := copier.Copy(&params, &arg); err != nil {
			return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to copy kontrak params")
		}

		var overlap pg.CheckKontrakOverlapParams

		if err := copier.Copy(&overlap, &arg); err != nil {
			return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to copy kontrak params")
		}

		if arg.PeriodeMulai != nil {
			mulai, err := time.Parse(time.RFC3339, *arg.PeriodeMulai)
			if err != nil {
				return nil, pkg.ExposeError(pkg.ErrorCodeBadRequest, "periode mulai invalid")
			}
			params.PeriodeMulai = pgtype.Timestamptz{Time: mulai, Valid: true}
			overlap.PeriodeMulai = pgtype.Timestamptz{Time: mulai, Valid: true}
		}

		if arg.PeriodeSelesai != nil {
			selesai, err := time.Parse(time.RFC3339, *arg.PeriodeSelesai)
			if err != nil {
				return nil, pkg.ExposeError(pkg.ErrorCodeBadRequest, "periode selesai invalid")
			}
			params.PeriodeSelesai = pgtype.Timestamptz{Time: selesai, Valid: true}
			overlap.PeriodeSelesai = pgtype.Timestamptz{Time: selesai, Valid: true}
		}

		isOverlap, err := qtx.CheckKontrakOverlap(c, overlap)
		if err != nil {
			return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed check overlap kontrak")
		}
		if isOverlap {
			return nil, pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "periode kontrak bertabrakan dengan data lain")
		}

		res, err := qtx.CreateKontrak(c, params)
		if err != nil {
			return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed create kontrak")
		}
		return res, nil
	})
}

func (mu *KontrakUsecaseImpl) ListKontrak(c context.Context, arg request.SearchKontrak) (any, error) {
	if arg.Limit <= 0 {
		arg.Limit = 10
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	arg.Offset = utils.GetOffset(arg.Page, arg.Limit)

	var params pg.ListKontrakParams
	if err := copier.Copy(&params, &arg); err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to copy list params")
	}

	params.PeriodeMulai = utils.StringToTimestamptz(arg.PeriodeMulai)
	params.PeriodeSelesai = utils.StringToTimestamptz(arg.PeriodeSelesai)

	res, err := mu.db.ListKontrak(c, params)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get kontrak")
	}

	if len(res) == 0 {
		return resp.WithPaginate([]any{}, resp.CalculatePagination(arg.Page, arg.Limit, 0)), nil
	}

	var cparams pg.CountKontrakParams
	if err := copier.Copy(&cparams, &arg); err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to copy count params")
	}

	count, err := mu.db.CountKontrak(c, cparams)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to count kontrak")
	}

	return resp.WithPaginate(res, resp.CalculatePagination(arg.Page, arg.Limit, count)), nil
}

func (mu *KontrakUsecaseImpl) ListAktifKontrak(c context.Context, arg *string) (any, error) {
	res, err := mu.db.ListAktifKontrak(c, arg)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get aktif kontrak")
	}

	if len(res) == 0 {
		return resp.WithPaginate([]any{}, nil), nil
	}

	return resp.WithPaginate(res, nil), nil
}

func (mu *KontrakUsecaseImpl) UpdateKontrak(c context.Context, arg pg.UpdateKontrakPartialParams) (any, error) {
	res, err := mu.db.UpdateKontrakPartial(c, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkg.ExposeError(pkg.ErrorCodeNotFound, "kontrak not found")
		}
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed update kontrak")
	}
	return res, nil
}

func (mu *KontrakUsecaseImpl) KontrakById(c context.Context, arg uuid.UUID) (any, error) {
	res, err := mu.db.GetKontrakByID(c, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkg.ExposeError(pkg.ErrorCodeNotFound, "kontrak not found")
		}
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get kontrak")
	}
	return resp.WithPaginate(res, nil), nil
}

func (mu *KontrakUsecaseImpl) DeleteKontrak(c context.Context, arg pg.DeleteKontrakParams) error {
	err := mu.db.DeleteKontrak(c, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pkg.ExposeError(pkg.ErrorCodeNotFound, "kontrak not found")
		}
		return pkg.WrapError(err, pkg.ErrorCodeInternal, "failed delete kontrak")
	}
	return nil
}
