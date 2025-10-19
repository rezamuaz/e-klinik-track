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
SELECT id,created_at
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
ORDER BY k.tgl_kehadiran DESC;
