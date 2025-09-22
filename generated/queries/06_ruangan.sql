-- name: CreateRuangan :one
INSERT INTO ruangan (
  fasilitas_id, kontrak_id, nama_ruangan, created_by
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetRuangan :one
SELECT * FROM ruangan
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListRuangan :many
SELECT
  id,
  fasilitas_id,
  kontrak_id,
  nama_ruangan,
  is_active,
  created_by,
  created_at
FROM ruangan
WHERE deleted_at IS NULL
  AND (sqlc.narg('nama_ruangan')::text IS NULL OR nama_ruangan ILIKE '%' || sqlc.narg('nama_ruangan') || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('fasilitas_id')::uuid IS NULL OR fasilitas_id = sqlc.narg('fasilitas_id')::uuid)
  AND (sqlc.narg('kontrak_id')::uuid IS NULL OR kontrak_id = sqlc.narg('kontrak_id')::uuid)
ORDER BY
  CASE WHEN sqlc.narg('order_by')::text = 'nama_ruangan' AND sqlc.narg('sort')::text = 'asc'  THEN nama_ruangan END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'nama_ruangan' AND sqlc.narg('sort')::text = 'desc' THEN nama_ruangan END DESC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'asc'  THEN created_at END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'desc' THEN created_at END DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountRuangan :one
SELECT COUNT(*)::bigint
FROM ruangan
WHERE deleted_at IS NULL
  AND (sqlc.narg('nama_ruangan')::text IS NULL OR nama_ruangan ILIKE '%' || sqlc.narg('nama_ruangan') || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('fasilitas_id')::uuid IS NULL OR fasilitas_id = sqlc.narg('fasilitas_id')::uuid)
  AND (sqlc.narg('kontrak_id')::uuid IS NULL OR kontrak_id = sqlc.narg('kontrak_id')::uuid);


-- name: UpdateRuanganPartial :one
UPDATE ruangan
SET
  fasilitas_id = COALESCE(sqlc.narg('fasilitas_id'), fasilitas_id),
  kontrak_id   = COALESCE(sqlc.narg('kontrak_id'), kontrak_id),
  nama_ruangan = COALESCE(sqlc.narg('nama_ruangan'), nama_ruangan),
  is_active    = COALESCE(sqlc.narg('is_active'), is_active),
  updated_by   = COALESCE(sqlc.narg('updated_by'), updated_by),
  updated_note = COALESCE(sqlc.narg('updated_note'), updated_note),
  updated_at   = now()
WHERE id = $1
RETURNING *;

-- name: DeleteRuangan :exec
UPDATE ruangan
SET deleted_at = now(), deleted_by = $2
WHERE id = $1 AND deleted_at IS NULL;
