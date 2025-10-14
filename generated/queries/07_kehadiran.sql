-- name: CreateKehadiran :one
INSERT INTO kehadiran (
  fasilitas_id, kontrak_id, ruangan_id, pembimbing_id,user_id,mata_kuliah_id,
  jadwal_dinas, created_by, tgl_kehadiran
) VALUES (
  $1, $2, $3, $4,
  $5,$6, $7,$8,$9
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