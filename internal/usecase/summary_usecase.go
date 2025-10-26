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
	GetGlobalSKPPersentaseTahunanOtomatis(c context.Context) (any, error)
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

func (mu *SummaryUsecaseImpl) RekapSkpTercapaiMahasiswaByDate(
	c context.Context,
	arg request.SearchSkpTercapai,
) (any, error) {

	var params pg.GetRekapSkpTercapaiByUserParams
	if err := copier.Copy(&params, &arg); err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to copy params")
	}

	// 🔹 Validasi & parsing UserID (hindari panic)
	if arg.UserID != "" {
		uid, err := uuid.FromString(arg.UserID)
		if err != nil {
			return nil, pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "user_id invalid")
		}
		params.UserID = uid
	}

	// 🔹 Parsing tanggal aman (pastikan nil check)
	if arg.TglAwal == "" {
		return nil, pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "tanggal awal tidak boleh kosong")
	}
	if arg.TglAkhir == "" {
		return nil, pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "tanggal akhir tidak boleh kosong")
	}

	var tglAwal, tglAkhir time.Time
	var err error

	tglAwal, err = time.Parse("2006-01-02", arg.TglAwal)
	if err != nil {
		return nil, pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "format tanggal awal tidak valid (gunakan YYYY-MM-DD)")
	}

	tglAkhir, err = time.Parse("2006-01-02", arg.TglAkhir)
	if err != nil {
		return nil, pkg.ExposeError(pkg.ErrorCodeInvalidArgument, "format tanggal akhir tidak valid (gunakan YYYY-MM-DD)")
	}

	params.TglAwal = pgtype.Date{Time: tglAwal, Valid: true}
	params.TglAkhir = pgtype.Date{Time: tglAkhir, Valid: true}

	// 🔹 Eksekusi query sqlc
	res, err := mu.db.GetRekapSkpTercapaiByUser(c, params)
	if err != nil {
		// misalnya kalau tidak ada data, kembalikan slice kosong bukan error
		if errors.Is(err, pgx.ErrNoRows) {
			return resp.WithPaginate([]any{}, nil), nil
		}
		return nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed get rekap skp mahasiswa")
	}

	return resp.WithPaginate(res, nil), nil
}

func (mu *SummaryUsecaseImpl) GetGlobalSKPPersentaseTahunanOtomatis(c context.Context) (any, error) {
	// 1. Ambil tanggal dan tahun saat ini
	tgl, err := utils.GetJakartaDateObject()
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get jakarta time")
	}

	cacheKey := fmt.Sprintf("rekap:kehadiran:skp:tahunan:%d", tgl.Year())

	// 2. Coba ambil dari Redis
	if mu.cache != nil {
		if cached, err := mu.cache.GetRaw(c, cacheKey); err == nil && cached != "" {
			var data []pg.GetGlobalSKPPersentaseTahunanOtomatisRow // ⚠️ PERBAIKAN: Gunakan tipe data hasil yang benar
			if err := json.Unmarshal([]byte(cached), &data); err == nil {
				// Return dengan tipe yang sesuai
				return resp.WithPaginate(data, nil), nil
			}
			// Log error Unmarshal jika terjadi, lalu lanjutkan ke DB
		}
	}
	// 3. Eksekusi query utama ke Database
	res, err := mu.db.GetGlobalSKPPersentaseTahunanOtomatis(c)
	if err != nil {
		return nil, pkg.WrapError(err, pkg.ErrorCodeUnknown, "failed get kehadiran")
	}
	// 4. Simpan ke Redis selama 5 menit
	if mu.cache != nil {
		go func() {
			// ⚠️ PERBAIKAN: Marshal data ke JSON sebelum disimpan.
			// Menggunakan tipe yang benar (res) lebih aman daripada 'any'.
			jsonBytes, err := json.Marshal(res)
			if err == nil {
				// Simpan sebagai string JSON. TTL 5 menit cocok untuk data yang sering berubah.
				_ = mu.cache.SetWithTTL(context.Background(), cacheKey, string(jsonBytes), 5*time.Minute)
			}
		}()
	}

	// 5. Kembalikan hasil dari Database
	return resp.WithPaginate(res, nil), nil
}
