-- name: CreateKontrak :one
INSERT INTO kontrak (
  fasilitas_id,
  nama,
  periode_mulai,
  periode_selesai,
  durasi,
  deskripsi,
  created_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetKontrak :one
SELECT * FROM kontrak
WHERE id = $1;


-- name: ListKontrak :many
SELECT
  k.id,
  k.fasilitas_id,
  k.nama,
  k.periode_mulai,
  k.periode_selesai,
  k.durasi,
  k.deskripsi,
  k.is_active,
  k.created_by,
  k.created_at,
  f.nama AS fasilitas_nama,
  f.kab AS fasilitas_kab,
  f.propinsi AS fasilitas_propinsi
FROM kontrak k
LEFT JOIN fasilitas_kesehatan f
  ON k.fasilitas_id = f.id
WHERE k.deleted_at IS NULL
  AND (sqlc.narg('nama')::text IS NULL OR k.nama ILIKE '%' || sqlc.narg('nama') || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR k.is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('fasilitas_nama')::text IS NULL OR f.nama ILIKE '%' || sqlc.narg('fasilitas_nama') || '%')
  AND (sqlc.narg('fasilitas_kab')::text IS NULL OR f.kab ILIKE '%' || sqlc.narg('fasilitas_kab') || '%')
  AND (sqlc.narg('fasilitas_propinsi')::text IS NULL OR f.propinsi ILIKE '%' || sqlc.narg('fasilitas_propinsi') || '%')
  AND (sqlc.narg('periode_mulai')::timestamptz IS NULL OR k.periode_mulai >= sqlc.narg('periode_mulai')::timestamptz)
  AND (sqlc.narg('periode_selesai')::timestamptz IS NULL OR k.periode_selesai <= sqlc.narg('periode_selesai')::timestamptz)
ORDER BY
  CASE WHEN sqlc.narg('order_by')::text = 'nama' AND sqlc.narg('sort')::text = 'asc'  THEN k.nama END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'nama' AND sqlc.narg('sort')::text = 'desc' THEN k.nama END DESC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'asc'  THEN k.created_at END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'desc' THEN k.created_at END DESC,
  CASE WHEN sqlc.narg('order_by')::text = 'periode_mulai' AND sqlc.narg('sort')::text = 'asc'  THEN k.periode_mulai END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'periode_mulai' AND sqlc.narg('sort')::text = 'desc' THEN k.periode_mulai END DESC,
  CASE WHEN sqlc.narg('order_by')::text = 'periode_selesai' AND sqlc.narg('sort')::text = 'asc'  THEN k.periode_selesai END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'periode_selesai' AND sqlc.narg('sort')::text = 'desc' THEN k.periode_selesai END DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountKontrak :one
SELECT COUNT(*)::bigint
FROM kontrak k
LEFT JOIN fasilitas_kesehatan f
  ON k.fasilitas_id = f.id
WHERE k.deleted_at IS NULL
  AND (sqlc.narg('nama')::text IS NULL OR k.nama ILIKE '%' || sqlc.narg('nama') || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR k.is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('fasilitas_nama')::text IS NULL OR f.nama ILIKE '%' || sqlc.narg('fasilitas_nama') || '%')
  AND (sqlc.narg('fasilitas_kab')::text IS NULL OR f.kab ILIKE '%' || sqlc.narg('fasilitas_kab') || '%')
  AND (sqlc.narg('fasilitas_propinsi')::text IS NULL OR f.propinsi ILIKE '%' || sqlc.narg('fasilitas_propinsi') || '%')
  AND (sqlc.narg('periode_mulai')::timestamptz IS NULL OR k.periode_mulai >= sqlc.narg('periode_mulai')::timestamptz)
  AND (sqlc.narg('periode_selesai')::timestamptz IS NULL OR k.periode_selesai <= sqlc.narg('periode_selesai')::timestamptz);



-- name: UpdateKontrakPartial :one
UPDATE kontrak
SET
  nama            = COALESCE(sqlc.narg('nama'), nama),
  periode_mulai   = COALESCE(sqlc.narg('periode_mulai'), periode_mulai),
  periode_selesai = COALESCE(sqlc.narg('periode_selesai'), periode_selesai),
  durasi          = COALESCE(sqlc.narg('durasi'), durasi),
  deskripsi       = COALESCE(sqlc.narg('deskripsi'), deskripsi),
  is_active       = COALESCE(sqlc.narg('is_active'), is_active),
  updated_by      = COALESCE(sqlc.narg('updated_by'), updated_by),
  updated_note    = COALESCE(sqlc.narg('updated_note'), updated_note),
  updated_at      = now()
WHERE id = $1
RETURNING *;

-- name: DeleteKontrak :exec
UPDATE kontrak
SET
  deleted_by = $2,
  deleted_at = now()
WHERE id = $1;
