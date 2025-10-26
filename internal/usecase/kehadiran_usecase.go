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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jinzhu/copier"
)

type KehadiranUsecase interface {
	AddKehadiran(c context.Context, arg pg.CreateKehadiranParams) (any, error)
	ListKehadiran(c context.Context, arg request.SearchKehadiran) (any, error)
	UpdateKehadiran(c context.Context, arg pg.UpdateKehadiranPartialParams) (any, error)
	DeleteKehadiran(c context.Context, arg pg.DeleteKehadiranParams) error
	CheckKehadiran(c context.Context, arg uuid.UUID) (any, error)
	GetKehadiranByPembimbingStatus(c context.Context, arg pg.GetKehadiranByPembimbingUserIdParams) (any, error)
	GetKehadiranByMahasiswaStatus(c context.Context, arg pg.GetKehadiranByPembimbingUserIdParams) (any, error)
}

type KehadiranUsecaseImpl struct {
	worker *worker.ProducerService
	db     *pg.Queries
	pg     *pkg.Postgres
	cache  *pkg.RedisCache
}

func NewKehadiranUsecase(postgre *pkg.Postgres, worker *worker.ProducerService, cache *pkg.RedisCache) *KehadiranUsecaseImpl {
	return &KehadiranUsecaseImpl{
		db:     pg.New(postgre.Pool),
		pg:     postgre,
		worker: worker,
		cache:  cache,
	}
}

func (mu *KehadiranUsecaseImpl) AddKehadiran(c context.Context, arg pg.CreateKehadiranParams) (any, error) {
	tgl, err := utils.GetJakartaDateObject()
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get jakarta time")
	}

	arg.TglKehadiran = pgtype.Date{Valid: true, Time: tgl}

	res, err := mu.db.CreateKehadiran(c, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23505":
			return nil, pkg.ExposeError(pkg.ErrorCodeConflict, "Anda telah absen hari ini. Silakan coba lagi besok.")
		case errors.Is(err, pgx.ErrNoRows):
			return nil, pkg.ExposeError(pkg.ErrorCodeConflict, "Anda telah absen hari ini. Silakan coba lagi besok.")
		default:
			return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed create kehadiran")
		}
	}

	return res, nil
}

func (mu *KehadiranUsecaseImpl) ListKehadiran(c context.Context, arg request.SearchKehadiran) (any, error) {
	if arg.Limit <= 0 {
		arg.Limit = 10
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	arg.Offset = utils.GetOffset(arg.Page, arg.Limit)

	var params pg.ListKehadiranParams
	if err := copier.Copy(&params, &arg); err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to copy search params")
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
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get kehadiran")
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

func (mu *KehadiranUsecaseImpl) UpdateKehadiran(c context.Context, arg pg.UpdateKehadiranPartialParams) (any, error) {
	res, err := mu.db.UpdateKehadiranPartial(c, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkg.ExposeError(pkg.ErrorCodeNotFound, "kehadiran not found")
		}
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed update kehadiran")
	}
	return res, nil
}

func (mu *KehadiranUsecaseImpl) DeleteKehadiran(c context.Context, arg pg.DeleteKehadiranParams) error {
	err := mu.db.DeleteKehadiran(c, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pkg.ExposeError(pkg.ErrorCodeNotFound, "kehadiran not found")
		}
		return pkg.WrapError(err, pkg.ErrorCodeInternal, "failed delete kehadiran")
	}
	return nil
}

func (mu *KehadiranUsecaseImpl) CheckKehadiran(c context.Context, arg uuid.UUID) (any, error) {
	res, err := mu.db.CheckKehadiran(c, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return resp.WithPaginate(map[string]any{}, nil), nil
		}
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed check kehadiran")
	}
	return resp.WithPaginate(res, nil), nil
}
func (mu *KehadiranUsecaseImpl) GetKehadiranByPembimbingStatus(c context.Context, arg pg.GetKehadiranByPembimbingUserIdParams) (any, error) {
	res, err := mu.db.GetKehadiranByPembimbingUserId(c, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return resp.WithPaginate([]any{}, nil), nil
		}
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get kehadiran by pembimbing")
	}
	return resp.WithPaginate(res, nil), nil
}

func (mu *KehadiranUsecaseImpl) GetKehadiranByMahasiswaStatus(c context.Context, arg pg.GetKehadiranByPembimbingUserIdParams) (any, error) {
	res, err := mu.db.GetKehadiranByPembimbingUserId(c, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return resp.WithPaginate([]any{}, nil), nil
		}
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get kehadiran by mahasiswa")
	}
	return resp.WithPaginate(res, nil), nil
}
