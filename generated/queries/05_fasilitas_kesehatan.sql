-- name: CreateFasilitasKesehatan :one
INSERT INTO fasilitas_kesehatan (
  nama, propinsi, kab, alamat, thumbnail, telepon, pemilik, kelas,
  longitude, latitude, gmap_name, gmap_address, gmap_thumbnail,
  similarity, tipe, is_active, created_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8,
  $9, $10, $11, $12, $13,
  $14, $15, COALESCE($16, true), $17
)
RETURNING *;

-- name: GetFasilitasKesehatan :one
SELECT * FROM fasilitas_kesehatan
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListFasilitasKesehatan :many
SELECT
  id,
  nama,
  propinsi,
  kab,
  alamat,
  thumbnail,
  telepon,
  pemilik,
  kelas,
  longitude,
  latitude,
  gmap_name,
  gmap_address,
  gmap_thumbnail,
  similarity,
  tipe,
  is_active,
  created_by,
  created_at
FROM fasilitas_kesehatan
WHERE deleted_at IS NULL
  AND (sqlc.narg('nama')::text IS NULL OR nama ILIKE '%' || sqlc.narg('nama') || '%')
  AND (sqlc.narg('propinsi')::text IS NULL OR propinsi ILIKE '%' || sqlc.narg('propinsi') || '%')
  AND (sqlc.narg('kab_id')::uuid IS NULL OR kab_id = sqlc.narg('kab_id')::uuid)
  AND (sqlc.narg('kab')::text IS NULL OR kab ILIKE '%' || sqlc.narg('kab') || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('tipe')::text IS NULL OR tipe = sqlc.narg('tipe')::text)
ORDER BY
  CASE WHEN sqlc.narg('order_by')::text = 'nama' AND sqlc.narg('sort')::text = 'asc'  THEN nama END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'nama' AND sqlc.narg('sort')::text = 'desc' THEN nama END DESC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'asc'  THEN created_at END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'desc' THEN created_at END DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountFasilitasKesehatan :one
SELECT COUNT(*)::bigint
FROM fasilitas_kesehatan
WHERE deleted_at IS NULL
  AND (sqlc.narg('nama')::text IS NULL OR nama ILIKE '%' || sqlc.narg('nama') || '%')
  AND (sqlc.narg('propinsi')::text IS NULL OR propinsi ILIKE '%' || sqlc.narg('propinsi') || '%')
  AND (sqlc.narg('kab')::text IS NULL OR kab ILIKE '%' || sqlc.narg('kab') || '%')
  AND (sqlc.narg('kab_id')::uuid IS NULL OR kab_id = sqlc.narg('kab_id')::uuid)
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('tipe')::text IS NULL OR tipe = sqlc.narg('tipe')::text);

-- name: ListDistinctKabupaten :many
SELECT id,nama
FROM kabupaten
WHERE (sqlc.narg('nama')::text IS NULL OR nama ILIKE '%' || sqlc.narg('nama')::text || '%')
  AND (sqlc.narg('propinsi_id')::uuid IS NULL OR propinsi_id = sqlc.narg('propinsi_id')::uuid)
ORDER BY nama ASC
LIMIT $1 OFFSET $2;

-- name: ListDistinctPropinsi :many
SELECT id,nama
FROM propinsi
WHERE (sqlc.narg('propinsi')::text IS NULL OR nama ILIKE '%' || sqlc.narg('propinsi')::text || '%')
ORDER BY nama ASC
LIMIT $1 OFFSET $2;


-- name: UpdateFasilitasKesehatanPartial :one
UPDATE fasilitas_kesehatan
SET
  nama           = COALESCE(sqlc.narg('nama'), nama),
  propinsi       = COALESCE(sqlc.narg('propinsi'), propinsi),
  kab            = COALESCE(sqlc.narg('kab'), kab),
  alamat         = COALESCE(sqlc.narg('alamat'), alamat),
  thumbnail      = COALESCE(sqlc.narg('thumbnail'), thumbnail),
  telepon        = COALESCE(sqlc.narg('telepon'), telepon),
  pemilik        = COALESCE(sqlc.narg('pemilik'), pemilik),
  kelas          = COALESCE(sqlc.narg('kelas'), kelas),
  longitude      = COALESCE(sqlc.narg('longitude'), longitude),
  latitude       = COALESCE(sqlc.narg('latitude'), latitude),
  gmap_name      = COALESCE(sqlc.narg('gmap_name'), gmap_name),
  gmap_address   = COALESCE(sqlc.narg('gmap_address'), gmap_address),
  gmap_thumbnail = COALESCE(sqlc.narg('gmap_thumbnail'), gmap_thumbnail),
  similarity     = COALESCE(sqlc.narg('similarity'), similarity),
  tipe           = COALESCE(sqlc.narg('tipe'), tipe),
  is_active      = COALESCE(sqlc.narg('is_active'), is_active),
  updated_by     = COALESCE(sqlc.narg('updated_by'), updated_by),
  updated_note   = COALESCE(sqlc.narg('updated_note'), updated_note),
  updated_at     = now()
WHERE id = $1
RETURNING *;

-- name: DeleteFasilitasKesehatan :exec
UPDATE fasilitas_kesehatan
SET deleted_at = now(), deleted_by = $2
WHERE id = $1 AND deleted_at IS NULL;
