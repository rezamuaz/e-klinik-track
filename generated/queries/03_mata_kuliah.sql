-- name: CreateMataKuliah :one
INSERT INTO mata_kuliah (
  mata_kuliah,
  created_by
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetMataKuliah :one
SELECT * FROM mata_kuliah
WHERE id = $1;

-- name: ListMataKuliah :many
SELECT
  id,
  mata_kuliah,
  is_active,
  created_by,
  created_at
FROM mata_kuliah
WHERE deleted_at IS NULL
  AND (sqlc.narg('mata_kuliah')::text IS NULL OR mata_kuliah ILIKE '%' || sqlc.narg('mata_kuliah') || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
ORDER BY
  CASE WHEN sqlc.narg('order_by')::text = 'mata_kuliah' AND sqlc.narg('sort')::text = 'asc'  THEN mata_kuliah END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'mata_kuliah' AND sqlc.narg('sort')::text = 'desc' THEN mata_kuliah END DESC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'asc'  THEN created_at END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'desc' THEN created_at END DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountMataKuliah :one
SELECT COUNT(*)::bigint
FROM mata_kuliah
WHERE deleted_at IS NULL
  AND (sqlc.narg('mata_kuliah')::text IS NULL OR mata_kuliah ILIKE '%' || sqlc.narg('mata_kuliah') || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean);

-- name: UpdateMataKuliah :one
UPDATE mata_kuliah
SET
  mata_kuliah = COALESCE(sqlc.narg('mata_kuliah'), mata_kuliah),
  is_active   = COALESCE(sqlc.narg('is_active'), is_active),
  updated_by  = COALESCE(sqlc.narg('updated_by'), updated_by),
  updated_note = COALESCE(sqlc.narg('updated_note'), updated_note),
  updated_at  = now()
WHERE id = $1
RETURNING *;

-- name: DeleteMataKuliah :exec
UPDATE mata_kuliah
SET deleted_at = now(), deleted_by = $2
WHERE id = $1 AND deleted_at IS NULL;

