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

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jinzhu/copier"
)

type SkpKehadiranUsecase interface {
	SyncSkpKehadiran(c context.Context, arg pg.SyncKehadiranSkpParams) (any, error)
	ListKehadiranSkp(c context.Context, arg request.SearchKehadiranSkp) (any, error)
	UpdateKehadiranSkp(c context.Context, arg pg.UpdateKehadiranSkpParams) (any, error)
	DeleteKehadiranSkp(c context.Context, arg pg.DeleteKehadiranSkpParams) error
	SkpByKehadiranId(c context.Context, arg uuid.UUID) (any, error)
	IntervensiByKehadiranId(c context.Context, arg uuid.UUID) (any, error)
	ApproveSkpKehadiran(c context.Context, arg request.ApproveKehadiranSkp) (any, error)
}

type SkpKehadiranUsecaseImpl struct {
	worker *worker.ProducerService
	db     *pg.Queries
	pg     *pkg.Postgres
	cache  *pkg.RedisCache
}

func NewSkpKehadiranUsecase(postgre *pkg.Postgres, worker *worker.ProducerService, cache *pkg.RedisCache) *SkpKehadiranUsecaseImpl {
	return &SkpKehadiranUsecaseImpl{
		db:     pg.New(postgre.Pool),
		pg:     postgre,
		worker: worker,
		cache:  cache,
	}
}

func (mu *SkpKehadiranUsecaseImpl) SyncSkpKehadiran(c context.Context, arg pg.SyncKehadiranSkpParams) (any, error) {
	return utils.WithTransactionResult(c, mu.pg.Pool, func(qtx *pg.Queries, tx pgx.Tx) (any, error) {
		result, err := qtx.GetKehadiran(c, arg.KehadiranID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, pkg.ExposeError(pkg.ErrorCodeNotFound, "kehadiran not found")
			}
			return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get kehadiran")
		}

		if result.Status == nil {
			if _, err := qtx.SyncKehadiranSkp(c, arg); err != nil {
				return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed sync kehadiran skp")
			}
		}
		return result, nil
	})
}

func (mu *SkpKehadiranUsecaseImpl) ListKehadiranSkp(c context.Context, arg request.SearchKehadiranSkp) (any, error) {
	if arg.Limit <= 0 {
		arg.Limit = 10
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	arg.Offset = utils.GetOffset(arg.Page, arg.Limit)

	var params pg.ListKehadiranSkpParams
	if err := copier.Copy(&params, &arg); err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed copy params")
	}

	if arg.KehadiranID != nil && *arg.KehadiranID != "" {
		fid := uuid.FromStringOrNil(*arg.KehadiranID)
		params.KehadiranID = &fid
	}
	if arg.SkpIntervensiID != nil && *arg.SkpIntervensiID != "" {
		kid := uuid.FromStringOrNil(*arg.SkpIntervensiID)
		params.SkpIntervensiID = &kid
	}

	res, err := mu.db.ListKehadiranSkp(c, params)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get kehadiran skp")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]any{}, resp.CalculatePagination(arg.Page, arg.Limit, 0)), nil
	}

	var cparams pg.CountKehadiranParams
	if err := copier.Copy(&cparams, &arg); err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed copy count params")
	}

	count, err := mu.db.CountKehadiran(c, cparams)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed count kehadiran")
	}

	return resp.WithPaginate(res, resp.CalculatePagination(arg.Page, arg.Limit, count)), nil
}

func (mu *SkpKehadiranUsecaseImpl) UpdateKehadiranSkp(c context.Context, arg pg.UpdateKehadiranSkpParams) (any, error) {
	res, err := mu.db.UpdateKehadiranSkp(c, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkg.ExposeError(pkg.ErrorCodeNotFound, "kehadiran skp not found")
		}
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed update kehadiran skp")
	}
	return res, nil
}

func (mu *SkpKehadiranUsecaseImpl) DeleteKehadiranSkp(c context.Context, arg pg.DeleteKehadiranSkpParams) error {
	err := mu.db.DeleteKehadiranSkp(c, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pkg.ExposeError(pkg.ErrorCodeNotFound, "kehadiran skp not found")
		}
		return pkg.WrapError(err, pkg.ErrorCodeInternal, "failed delete kehadiran skp")
	}
	return nil
}

func (mu *SkpKehadiranUsecaseImpl) SkpByKehadiranId(c context.Context, arg uuid.UUID) (any, error) {
	res, err := mu.db.SkpKehadiranID(c, arg)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get data")
	}
	if len(res) == 0 {
		return resp.WithPaginate([]any{}, nil), nil
	}
	return resp.WithPaginate(res, nil), nil
}

func (mu *SkpKehadiranUsecaseImpl) IntervensiByKehadiranId(c context.Context, arg uuid.UUID) (any, error) {
	res, err := mu.db.IntervensiKehadiranID(c, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return resp.WithPaginate([]any{}, nil), nil
		}
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get intervensi")
	}
	return resp.WithPaginate(res, nil), nil
}

func (mu *SkpKehadiranUsecaseImpl) ApproveSkpKehadiran(c context.Context, arg request.ApproveKehadiranSkp) (any, error) {
	return utils.WithTransactionResult(c, mu.pg.Pool, func(qtx *pg.Queries, tx pgx.Tx) (any, error) {
		updateParams := pg.UpdateKehadiranPartialParams{
			ID:     arg.KehadiranID,
			Status: utils.StringPtr("disetujui"),
		}
		res, err := qtx.UpdateKehadiranPartial(c, updateParams)
		if err != nil {
			return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed update kehadiran")
		}

		approveParams := pg.ApproveKehadiranSkpByIdsParams{
			UpdatedBy:      arg.UpdatedBy,
			SkpKehadiranID: arg.SkpKehadiranID,
			KehadiranID:    arg.KehadiranID,
		}
		if err := qtx.ApproveKehadiranSkpByIds(c, approveParams); err != nil {
			return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed approve kehadiran skp")
		}

		return res, nil
	})
}
