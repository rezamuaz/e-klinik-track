-- name: CreateKehadiran :one
INSERT INTO kehadiran (
  fasilitas_id, kontrak_id, ruangan_id, pembimbing_id,user_id,pembimbing_klinik,mata_kuliah_id,
  jadwal_dinas, created_by, tgl_kehadiran, presensi
) VALUES (
  $1, $2, $3, $4,
  $5,$6, $7,$8,$9,$10,$11
)
ON CONFLICT (user_id, tgl_kehadiran) DO NOTHING
RETURNING *;

-- name: GetKehadiran :one
SELECT * FROM kehadiran
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListKehadiran :many
SELECT
  id,
  fasilitas_id,
  kontrak_id,
  ruangan_id,
  pembimbing_id,
  jadwal_dinas,
  is_active,
  created_by,
  created_at
FROM kehadiran
WHERE deleted_at IS NULL
  AND (sqlc.narg('jadwal_dinas')::text IS NULL OR jadwal_dinas ILIKE '%' || sqlc.narg('jadwal_dinas') || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('fasilitas_id')::uuid IS NULL OR fasilitas_id = sqlc.narg('fasilitas_id')::uuid)
  AND (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id')::uuid)
  AND (sqlc.narg('kontrak_id')::uuid IS NULL OR kontrak_id = sqlc.narg('kontrak_id')::uuid)
  AND (sqlc.narg('pembimbing_id')::uuid IS NULL OR pembimbing_id = sqlc.narg('pembimbing_id')::uuid)
ORDER BY
  CASE WHEN sqlc.narg('order_by')::text = 'jadwal_dinas' AND sqlc.narg('sort')::text = 'asc'  THEN jadwal_dinas END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'jadwal_dinas' AND sqlc.narg('sort')::text = 'desc' THEN jadwal_dinas END DESC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'asc'  THEN created_at END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'desc' THEN created_at END DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountKehadiran :one
SELECT COUNT(*)::bigint
FROM kehadiran
WHERE deleted_at IS NULL
  AND (sqlc.narg('jadwal_dinas')::text IS NULL OR jadwal_dinas ILIKE '%' || sqlc.narg('jadwal_dinas') || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('fasilitas_id')::uuid IS NULL OR fasilitas_id = sqlc.narg('fasilitas_id')::uuid)
  AND (sqlc.narg('kontrak_id')::uuid IS NULL OR kontrak_id = sqlc.narg('kontrak_id')::uuid)
  AND (sqlc.narg('pembimbing_id')::uuid IS NULL OR pembimbing_id = sqlc.narg('pembimbing_id')::uuid);

-- name: UpdateKehadiranPartial :one
UPDATE kehadiran
SET
  fasilitas_id  = COALESCE(sqlc.narg('fasilitas_id'), fasilitas_id),
  kontrak_id    = COALESCE(sqlc.narg('kontrak_id'), kontrak_id),
  ruangan_id    = COALESCE(sqlc.narg('ruangan_id'), ruangan_id),
  pembimbing_id = COALESCE(sqlc.narg('pembimbing_id'), pembimbing_id),
  jadwal_dinas  = COALESCE(sqlc.narg('jadwal_dinas'), jadwal_dinas),
  is_active     = COALESCE(sqlc.narg('is_active'), is_active),
  status        = COALESCE(sqlc.narg('status'), status),
  updated_by    = COALESCE(sqlc.narg('updated_by'), updated_by),
  updated_note  = COALESCE(sqlc.narg('updated_note'), updated_note),
  updated_at    = now()
WHERE id = $1
RETURNING *;

-- name: DeleteKehadiran :exec
UPDATE kehadiran
SET deleted_at = now(), deleted_by = $2
WHERE id = $1 AND deleted_at IS NULL;

-- name: CheckKehadiran :one
SELECT id,created_at,presensi
FROM kehadiran
WHERE tgl_kehadiran = (
    CURRENT_DATE AT TIME ZONE 'Asia/Jakarta'
)
AND user_id = sqlc.arg('user_id')
AND is_active = TRUE;

-- name: RekapKehadiranMahasiswa :one
SELECT
    user_id,
    COUNT(*) FILTER (WHERE presensi = 'hadir') AS total_hadir,
    COUNT(*) FILTER (WHERE presensi = 'izin') AS total_izin,
    COUNT(*) FILTER (WHERE presensi = 'sakit') AS total_sakit,
    COUNT(*) AS total_semua
FROM kehadiran
WHERE is_active = true
  AND user_id = sqlc.arg('user_id')
  AND tgl_kehadiran BETWEEN sqlc.arg('tgl_awal') AND sqlc.arg('tgl_akhir')
GROUP BY user_id;


-- name: RekapKehadiranMahasiswaDetail :many
SELECT
    tgl_kehadiran,
    COUNT(*) FILTER (WHERE presensi = 'hadir') AS total_hadir,
    COUNT(*) FILTER (WHERE presensi = 'izin') AS total_izin,
    COUNT(*) FILTER (WHERE presensi = 'sakit') AS total_sakit
FROM kehadiran
WHERE is_active = true
  AND user_id = sqlc.arg('user_id')
  AND tgl_kehadiran BETWEEN sqlc.arg('tgl_awal') AND sqlc.arg('tgl_akhir')
GROUP BY tgl_kehadiran
ORDER BY tgl_kehadiran;


-- name: GetKehadiranByPembimbingUserId :many
SELECT
    k.id,
    u.id AS user_id,
    u.nama,
    k.tgl_kehadiran
FROM kehadiran k
JOIN users u ON u.id = k.user_id
WHERE k.is_active = true
  AND (sqlc.narg('pembimbing_klinik')::uuid IS NULL OR k.pembimbing_klinik = sqlc.narg('pembimbing_klinik')::uuid)
  AND (sqlc.narg('user_id')::uuid IS NULL OR k.user_id = sqlc.narg('user_id')::uuid)
  AND k.status IS NULL
  AND k.presensi = 'hadir'
ORDER BY k.tgl_kehadiran DESC;


--////////////////////////////////////////////////////////////////

-- name: GetRekapGlobalHarian :one
WITH
-- Data hari ini
TotalData AS (
  SELECT
    COUNT(DISTINCT k.user_id) AS total_mahasiswa_unik,
    COUNT(*) FILTER (WHERE k.presensi = 'hadir') AS hadir,
    COUNT(*) FILTER (WHERE k.presensi = 'izin') AS izin,
    COUNT(*) FILTER (WHERE k.presensi = 'sakit') AS sakit
  FROM kehadiran k
  WHERE k.tgl_kehadiran = sqlc.arg('tgl')
    AND k.is_active = TRUE
),

-- Data kemarin
YesterdayData AS (
  SELECT
    COUNT(DISTINCT k.user_id) AS total_mahasiswa_unik,
    COUNT(*) FILTER (WHERE k.presensi = 'hadir') AS hadir,
    COUNT(*) FILTER (WHERE k.presensi = 'izin') AS izin,
    COUNT(*) FILTER (WHERE k.presensi = 'sakit') AS sakit
  FROM kehadiran k
  WHERE k.tgl_kehadiran = (sqlc.arg('tgl') - INTERVAL '1 day')
    AND k.is_active = TRUE
)

SELECT
  sqlc.arg('tgl')::date AS tanggal,
  COALESCE(td.total_mahasiswa_unik, 0) AS total_mahasiswa,
  COALESCE(td.hadir, 0) AS hadir,
  COALESCE(td.izin, 0) AS izin,
  COALESCE(td.sakit, 0) AS sakit,

  -- Persentase hari ini
  ROUND((COALESCE(td.hadir, 0)::numeric / NULLIF(td.total_mahasiswa_unik, 0)) * 100, 2) AS persentase_hadir,
  ROUND((COALESCE(td.izin, 0)::numeric / NULLIF(td.total_mahasiswa_unik, 0)) * 100, 2) AS persentase_izin,
  ROUND((COALESCE(td.sakit, 0)::numeric / NULLIF(td.total_mahasiswa_unik, 0)) * 100, 2) AS persentase_sakit,

  -- Growth aman
  CASE
    WHEN COALESCE(yd.hadir, 0) = 0 AND COALESCE(td.hadir, 0) > 0
      THEN 100
    WHEN COALESCE(yd.hadir, 0) = 0 AND COALESCE(td.hadir, 0) = 0
      THEN 0
    ELSE ROUND(((td.hadir - yd.hadir)::numeric / yd.hadir) * 100, 2)
  END AS growth_hadir,

  CASE
    WHEN COALESCE(yd.izin, 0) = 0 AND COALESCE(td.izin, 0) > 0
      THEN 100
    WHEN COALESCE(yd.izin, 0) = 0 AND COALESCE(td.izin, 0) = 0
      THEN 0
    ELSE ROUND(((td.izin - yd.izin)::numeric / yd.izin) * 100, 2)
  END AS growth_izin,

  CASE
    WHEN COALESCE(yd.sakit, 0) = 0 AND COALESCE(td.sakit, 0) > 0
      THEN 100
    WHEN COALESCE(yd.sakit, 0) = 0 AND COALESCE(td.sakit, 0) = 0
      THEN 0
    ELSE ROUND(((td.sakit - yd.sakit)::numeric / yd.sakit) * 100, 2)
  END AS growth_sakit,

  CASE
    WHEN COALESCE(yd.total_mahasiswa_unik, 0) = 0 AND COALESCE(td.total_mahasiswa_unik, 0) > 0
      THEN 100
    WHEN COALESCE(yd.total_mahasiswa_unik, 0) = 0 AND COALESCE(td.total_mahasiswa_unik, 0) = 0
      THEN 0
    ELSE ROUND(((td.total_mahasiswa_unik - yd.total_mahasiswa_unik)::numeric / yd.total_mahasiswa_unik) * 100, 2)
  END AS growth_total_mahasiswa

FROM TotalData td
LEFT JOIN YesterdayData yd ON TRUE;



-- name: GetRekapKehadiranPerFasilitasHarian :many
SELECT
  f.id AS fasilitas_id,
  f.nama AS nama_fasilitas,
  COUNT(DISTINCT k.user_id) AS total_mahasiswa,
  COUNT(*) FILTER (WHERE k.presensi = 'hadir') AS hadir,
  COUNT(*) FILTER (WHERE k.presensi = 'izin') AS izin,
  COUNT(*) FILTER (WHERE k.presensi = 'sakit') AS sakit,
  -- COUNT(*) FILTER (WHERE k.presensi = 'alpa') AS alpa,
  ROUND(
    (COUNT(*) FILTER (WHERE k.presensi = 'hadir')::numeric / NULLIF(COUNT(*), 0)) * 100,
    2
  ) AS persentase_hadir
FROM kehadiran k
JOIN fasilitas_kesehatan f ON f.id = k.fasilitas_id
WHERE k.tgl_kehadiran = $1
  AND k.is_active = TRUE
GROUP BY f.id, f.nama
ORDER BY f.nama;

-- name: GetRekapDetailFasilitasHarian :many
SELECT
  r.id AS ruangan_id,
  r.nama_ruangan,
  COUNT(DISTINCT k.user_id) AS total_mahasiswa,
  COUNT(*) FILTER (WHERE k.presensi = 'hadir') AS hadir,
  COUNT(*) FILTER (WHERE k.presensi = 'izin') AS izin,
  COUNT(*) FILTER (WHERE k.presensi = 'sakit') AS sakit,
  COUNT(*) FILTER (WHERE k.presensi = 'alpa') AS alpa,
  ROUND(
    (COUNT(*) FILTER (WHERE k.presensi = 'hadir')::numeric / NULLIF(COUNT(*), 0)) * 100,
    2
  ) AS persentase_hadir
FROM kehadiran k
JOIN ruangan r ON r.id = k.ruangan_id
WHERE k.tgl_kehadiran = $1
  AND k.fasilitas_id = $2
  AND k.is_active = TRUE
GROUP BY r.id, r.nama_ruangan
ORDER BY r.nama_ruangan;


-- name: GetRekapPembimbingFasilitasHarian :many
SELECT
  u.id AS pembimbing_id,
  u.nama AS nama_pembimbing,
  COUNT(DISTINCT k.user_id) AS total_mahasiswa,
  COUNT(*) FILTER (WHERE k.presensi = 'hadir') AS hadir,
  ROUND(
    (COUNT(*) FILTER (WHERE k.presensi = 'hadir')::numeric / NULLIF(COUNT(*), 0)) * 100,
    2
  ) AS persentase_hadir
FROM kehadiran k
JOIN users u ON u.id = k.pembimbing_id
WHERE k.tgl_kehadiran = $1
  AND k.fasilitas_id = $2
  AND k.is_active = TRUE
GROUP BY u.id, u.nama
ORDER BY u.nama;

-- name: GetMahasiswaTidakHadir :many
SELECT
  k.user_id,
  u.nama AS nama_mahasiswa,
  r.nama_ruangan,
  p.nama AS nama_pembimbing,
  k.presensi,
  k.status
FROM kehadiran k
JOIN users u ON u.id = k.user_id
JOIN ruangan r ON r.id = k.ruangan_id
LEFT JOIN users p ON p.id = k.pembimbing_id
WHERE k.tgl_kehadiran = $1
  AND k.fasilitas_id = $2
  AND k.presensi IN ('izin', 'sakit')
  AND k.is_active = TRUE
ORDER BY r.nama_ruangan, u.nama;




-- name: GetTrenKehadiran7Hari :many
SELECT
  f.id AS fasilitas_id,
  f.nama AS nama_fasilitas,
  k.tgl_kehadiran AS tanggal,
  COUNT(DISTINCT k.user_id) AS total_mahasiswa,
  COUNT(*) FILTER (WHERE k.presensi = 'hadir') AS hadir,
  COUNT(*) FILTER (WHERE k.presensi = 'izin') AS izin,
  COUNT(*) FILTER (WHERE k.presensi = 'sakit') AS sakit,
  COUNT(*) FILTER (WHERE k.presensi = 'alpa') AS alpa,
  ROUND(
    (COUNT(*) FILTER (WHERE k.presensi = 'hadir')::numeric / NULLIF(COUNT(*), 0)) * 100,
    2
  ) AS persentase_hadir
FROM kehadiran k
JOIN fasilitas_kesehatan f ON f.id = k.fasilitas_id
WHERE k.tgl_kehadiran BETWEEN (CURRENT_DATE - INTERVAL '6 days') AND CURRENT_DATE
  AND k.is_active = TRUE
GROUP BY f.id, f.nama, k.tgl_kehadiran
ORDER BY f.nama, k.tgl_kehadiran;
