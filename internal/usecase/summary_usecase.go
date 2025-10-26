package usecase

import (
	"context"
	"e-klinik/infra/pg"
	"e-klinik/infra/worker"
	"e-klinik/internal/domain/request"
	"e-klinik/internal/domain/resp"
	"e-klinik/pkg"
	"e-klinik/utils"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jinzhu/copier"
)

type SummaryUsecase interface {
	RekapKehadiranMahasiswa(c context.Context, arg request.SearchRekapKehadiranMahasiswa) (any, error)
	RekapKehadiranMahasiswaDetail(c context.Context, arg request.SearchRekapKehadiranMahasiswa) (any, error)
	GetKehadiranByMahasiswaStatus(c context.Context, arg pg.GetKehadiranByPembimbingUserIdParams) (any, error)
	GetRekapKehadiranGlobalHarian(c context.Context) (any, error)
	GetRekapSKPGlobalHarian(c context.Context) (any, error)
	GetRekapKehadiranPerFasilitasHarian(c context.Context) (any, error)
	ChartGetHarianSKPPersentase(c context.Context) (any, error)
	ChartGetHariIniSKPPersentase(c context.Context) (any, error)
	RekapSkpTercapaiMahasiswaByDate(c context.Context, arg request.SearchSkpTercapai) (any, error)
}

type SummaryUsecaseImpl struct {
	worker *worker.ProducerService
	db     *pg.Queries
	pg     *pkg.Postgres
	cache  *pkg.RedisCache
}

func NewSummaryUsecase(postgre *pkg.Postgres, worker *worker.ProducerService, cache *pkg.RedisCache) *SummaryUsecaseImpl {
	return &SummaryUsecaseImpl{
		db:     pg.New(postgre.Pool),
		pg:     postgre,
		worker: worker,
		cache:  cache,
	}
}

func (mu *SummaryUsecaseImpl) RekapKehadiranMahasiswa(c context.Context, arg request.SearchRekapKehadiranMahasiswa) (any, error) {
	var params pg.RekapKehadiranMahasiswaParams
	if err := copier.Copy(&params, &arg); err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to copy params")
	}

	if arg.UserID != "" {
		params.UserID = uuid.FromStringOrNil(arg.UserID)
	}

	var tglAwal, tglAkhir pgtype.Date
	_ = tglAwal.Scan(arg.TglAwal)
	_ = tglAkhir.Scan(arg.TglAkhir)
	params.TglAwal = tglAwal
	params.TglAkhir = tglAkhir

	res, err := mu.db.RekapKehadiranMahasiswa(c, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return resp.WithPaginate([]any{}, nil), nil
		}
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get rekap kehadiran mahasiswa")
	}

	return resp.WithPaginate(res, nil), nil
}

func (mu *SummaryUsecaseImpl) RekapKehadiranMahasiswaDetail(c context.Context, arg request.SearchRekapKehadiranMahasiswa) (any, error) {
	var params pg.RekapKehadiranMahasiswaDetailParams
	if err := copier.Copy(&params, &arg); err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to copy params")
	}

	if arg.UserID != "" {
		params.UserID = uuid.FromStringOrNil(arg.UserID)
	}

	var tglAwal, tglAkhir pgtype.Date
	_ = tglAwal.Scan(arg.TglAwal)
	_ = tglAkhir.Scan(arg.TglAkhir)
	params.TglAwal = tglAwal
	params.TglAkhir = tglAkhir

	res, err := mu.db.RekapKehadiranMahasiswaDetail(c, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return resp.WithPaginate([]any{}, nil), nil
		}
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get rekap kehadiran detail")
	}

	return resp.WithPaginate(res, nil), nil
}

func (mu *SummaryUsecaseImpl) GetKehadiranByMahasiswaStatus(c context.Context, arg pg.GetKehadiranByPembimbingUserIdParams) (any, error) {
	res, err := mu.db.GetKehadiranByPembimbingUserId(c, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return resp.WithPaginate([]any{}, nil), nil
		}
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get kehadiran by mahasiswa")
	}
	return resp.WithPaginate(res, nil), nil
}

func (mu *SummaryUsecaseImpl) GetRekapKehadiranGlobalHarian(c context.Context) (any, error) {
	tgl, err := utils.GetJakartaDateObject()
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get jakarta time")
	}
	jktDate := pgtype.Date{Valid: true, Time: tgl}

	// 2️⃣ Buat cache key unik berdasarkan tanggal
	cacheKey := fmt.Sprintf("rekap:kehadiran:global:harian:%s", tgl.Format("2006-01-02"))

	// 3️⃣ Coba ambil dari Redis
	if mu.cache != nil {
		if cached, err := mu.cache.GetRaw(c, cacheKey); err == nil && cached != "" {
			var data any
			if err := json.Unmarshal([]byte(cached), &data); err == nil {
				return resp.WithPaginate(data, nil), nil
			}
		}
	}
	res, err := mu.db.GetRekapGlobalHarian(c, jktDate)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get kehadiran")
	}
	// 5️⃣ Simpan ke Redis selama 5 menit
	if mu.cache != nil {
		go func() {
			// marshal json supaya ringan

			_ = mu.cache.SetWithTTL(context.Background(), cacheKey, res, 5*time.Minute)

		}()
	}

	return resp.WithPaginate(res, nil), nil
}

func (mu *SummaryUsecaseImpl) GetRekapSKPGlobalHarian(c context.Context) (any, error) {
	tgl, err := utils.GetJakartaDateObject()
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get jakarta time")
	}
	jktDate := pgtype.Date{Valid: true, Time: tgl}
	// 2️⃣ Buat cache key unik berdasarkan tanggal
	cacheKey := fmt.Sprintf("rekap:skp:global:harian:%s", tgl.Format("2006-01-02"))

	// 3️⃣ Coba ambil dari Redis
	if mu.cache != nil {
		if cached, err := mu.cache.GetRaw(c, cacheKey); err == nil && cached != "" {
			var data any
			if err := json.Unmarshal([]byte(cached), &data); err == nil {
				return resp.WithPaginate(data, nil), nil
			}
		}
	}
	res, err := mu.db.GetRekapSKPHarian(c, jktDate)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get kehadiran")
	}
	// 5️⃣ Simpan ke Redis selama 5 menit
	if mu.cache != nil {
		go func() {
			// marshal json supaya ringan

			_ = mu.cache.SetWithTTL(context.Background(), cacheKey, res, 5*time.Minute)

		}()
	}

	return resp.WithPaginate(res, nil), nil
}

func (mu *SummaryUsecaseImpl) GetRekapKehadiranPerFasilitasHarian(c context.Context) (any, error) {
	tgl, err := utils.GetJakartaDateObject()
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get jakarta time")
	}
	jktDate := pgtype.Date{Valid: true, Time: tgl}
	cacheKey := fmt.Sprintf("rekap:kehadiran:fasilitas:harian:%s", tgl.Format("2006-01-02"))

	// 3️⃣ Coba ambil dari Redis
	if mu.cache != nil {
		if cached, err := mu.cache.GetRaw(c, cacheKey); err == nil && cached != "" {
			var data any
			if err := json.Unmarshal([]byte(cached), &data); err == nil {
				return resp.WithPaginate(data, nil), nil
			}
		}
	}
	res, err := mu.db.GetRekapKehadiranPerFasilitasHarian(c, jktDate)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get kehadiran")
	}
	// 5️⃣ Simpan ke Redis selama 5 menit
	if mu.cache != nil {
		go func() {
			// marshal json supaya ringan

			_ = mu.cache.SetWithTTL(context.Background(), cacheKey, res, 5*time.Minute)

		}()
	}

	return resp.WithPaginate(res, nil), nil
}

func (mu *SummaryUsecaseImpl) ChartGetHarianSKPPersentase(c context.Context) (any, error) {
	tgl, err := utils.GetJakartaDateObject()
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get jakarta time")
	}
	jktDate := pgtype.Date{Valid: true, Time: tgl}
	cacheKey := fmt.Sprintf("rekap:skp:7:harian:%s", tgl.Format("2006-01-02"))

	// 3️⃣ Coba ambil dari Redis
	if mu.cache != nil {
		if cached, err := mu.cache.GetRaw(c, cacheKey); err == nil && cached != "" {
			var data any
			if err := json.Unmarshal([]byte(cached), &data); err == nil {
				return resp.WithPaginate(data, nil), nil
			}
		}
	}
	res, err := mu.db.GetCapaianSKP7HariTerakhir(c, jktDate)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get kehadiran")
	}
	// 5️⃣ Simpan ke Redis selama 5 menit
	if mu.cache != nil {
		go func() {
			// marshal json supaya ringan

			_ = mu.cache.SetWithTTL(context.Background(), cacheKey, res, 5*time.Minute)

		}()
	}

	return resp.WithPaginate(res, nil), nil
}

func (mu *SummaryUsecaseImpl) ChartGetHariIniSKPPersentase(c context.Context) (any, error) {
	tgl, err := utils.GetJakartaDateObject()
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get jakarta time")
	}

	jktDate := pgtype.Date{Valid: true, Time: tgl}
	cacheKey := fmt.Sprintf("rekap:global:harian:%s", tgl.Format("2006-01-02"))

	// 3️⃣ Coba ambil dari Redis
	if mu.cache != nil {
		if cached, err := mu.cache.GetRaw(c, cacheKey); err == nil && cached != "" {
			var data any
			if err := json.Unmarshal([]byte(cached), &data); err == nil {
				return resp.WithPaginate(data, nil), nil
			}
		}
	}
	res, err := mu.db.GetCapaianSKPPerHari(c, pg.GetCapaianSKPPerHariParams{StartDate: jktDate, EndDate: jktDate})
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get kehadiran")
	}
	// 5️⃣ Simpan ke Redis selama 5 menit
	if mu.cache != nil {
		go func() {
			// marshal json supaya ringan

			_ = mu.cache.SetWithTTL(context.Background(), cacheKey, res, 5*time.Minute)

		}()
	}

	return resp.WithPaginate(res, nil), nil
}

func (mu *SummaryUsecaseImpl) RekapSkpTercapaiMahasiswaByDate(c context.Context, arg request.SearchSkpTercapai) (any, error) {

	var params pg.GetRekapSkpTercapaiByUserParams
	err := copier.Copy(&params, &arg)
	if err != nil {
		return nil, err
	}

	if arg.UserID != "" {
		params.UserID = uuid.FromStringOrNil(arg.UserID)
	}

	var tglAwal pgtype.Date
	tglAwal.Scan(arg.TglAwal)
	var tglAkhir pgtype.Date
	tglAkhir.Scan(arg.TglAkhir)
	params.TglAwal = tglAwal
	params.TglAkhir = tglAkhir

	res, err := mu.db.GetRekapSkpTercapaiByUser(c, params)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get kehadiran")
	}

	return resp.WithPaginate(res, nil), nil
}
