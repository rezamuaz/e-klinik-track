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


-- name: GetRuanganById :one
SELECT
  r.id,
  r.fasilitas_id,
  r.kontrak_id,
  r.nama_ruangan,
  r.is_active,
  f.nama as rumah_sakit,
  k.no_utama as kontrak,
  f.propinsi,
  f.kab,
  r.created_by,
  r.created_at
FROM ruangan r
LEFT JOIN fasilitas_kesehatan f ON r.fasilitas_id = f.id
LEFT JOIN kontrak k ON r.kontrak_id = k.id
WHERE r.deleted_at IS NULL AND r.id = $1;

-- name: ListRuangan :many
SELECT
  r.id,
  r.fasilitas_id,
  r.kontrak_id,
  r.nama_ruangan,
  r.is_active,
  f.nama as rumah_sakit,
  k.no_utama as kontrak,
  f.propinsi,
  f.kab,
  r.created_by,
  r.created_at
FROM ruangan r
LEFT JOIN fasilitas_kesehatan f ON r.fasilitas_id = f.id
LEFT JOIN kontrak k ON r.kontrak_id = k.id
WHERE r.deleted_at IS NULL
  AND (sqlc.narg('nama_ruangan')::text IS NULL OR r.nama_ruangan ILIKE '%' || sqlc.narg('nama_ruangan')::text || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR r.is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('fasilitas')::text IS NULL OR f.nama ILIKE '%' || sqlc.narg('fasilitas')::text || '%')
  AND (sqlc.narg('kontrak')::text IS NULL OR k.no_utama ILIKE '%' || sqlc.narg('kontrak')::text || '%')
ORDER BY
  CASE WHEN sqlc.narg('order_by')::text = 'nama_ruangan' AND sqlc.narg('sort')::text = 'asc'  THEN r.nama_ruangan END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'nama_ruangan' AND sqlc.narg('sort')::text = 'desc' THEN r.nama_ruangan END DESC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'asc'  THEN r.created_at END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'desc' THEN r.created_at END DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountRuangan :one
SELECT COUNT(*)::bigint
FROM ruangan r
LEFT JOIN fasilitas_kesehatan f ON r.fasilitas_id = f.id
LEFT JOIN kontrak k ON r.kontrak_id = k.id
WHERE r.deleted_at IS NULL
  AND (sqlc.narg('nama_ruangan')::text IS NULL OR r.nama_ruangan ILIKE '%' || sqlc.narg('nama_ruangan')::text || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR r.is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('fasilitas')::text IS NULL OR f.nama ILIKE '%' || sqlc.narg('fasilitas')::text || '%')
  AND (sqlc.narg('kontrak')::text IS NULL OR k.no_utama ILIKE '%' || sqlc.narg('kontrak')::text || '%');
;



-- name: UpdateRuanganPartial :one
UPDATE ruangan
SET
  -- fasilitas_id = COALESCE(sqlc.narg('fasilitas_id'), fasilitas_id),
  -- kontrak_id   = COALESCE(sqlc.narg('kontrak_id'), kontrak_id),
  nama_ruangan = COALESCE(sqlc.narg('nama_ruangan'), nama_ruangan),
  is_active    = COALESCE(sqlc.narg('is_active'), is_active),
  updated_by   = COALESCE(sqlc.narg('updated_by'), updated_by),
  -- updated_note = COALESCE(sqlc.narg('updated_note'), updated_note),
  updated_at   = now()
WHERE id = $1
RETURNING *;

-- name: DeleteRuangan :exec
UPDATE ruangan
SET deleted_at = now(), deleted_by = $2
WHERE id = $1 AND deleted_at IS NULL;
