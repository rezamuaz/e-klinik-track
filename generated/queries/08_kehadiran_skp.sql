-- name: CreateKehadiranSkp :one
INSERT INTO kehadiran_skp (
    kehadiran_id,
    skp_intervensi_id,
    status,
    is_active,
    created_by
) VALUES (
    $1, $2, $3, COALESCE($4, true), $5
)
RETURNING *;

-- name: GetKehadiranSkp :one
SELECT *
FROM kehadiran_skp
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ListKehadiranSkp :many
SELECT
  id,
  kehadiran_id,
  skp_intervensi_id,
  status,
  is_active,
  created_by,
  created_at
FROM kehadiran_skp
WHERE deleted_at IS NULL
  AND (sqlc.narg('status')::text IS NULL OR status ILIKE '%' || sqlc.narg('status') || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('kehadiran_id')::uuid IS NULL OR kehadiran_id = sqlc.narg('kehadiran_id')::uuid)
  AND (sqlc.narg('skp_intervensi_id')::uuid IS NULL OR skp_intervensi_id = sqlc.narg('skp_intervensi_id')::uuid)
ORDER BY
  CASE WHEN sqlc.narg('order_by')::text = 'status' AND sqlc.narg('sort')::text = 'asc'  THEN status END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'status' AND sqlc.narg('sort')::text = 'desc' THEN status END DESC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'asc'  THEN created_at END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'desc' THEN created_at END DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountKehadiranSkp :one
SELECT COUNT(*)::bigint
FROM kehadiran_skp
WHERE deleted_at IS NULL
  AND (sqlc.narg('status')::text IS NULL OR status ILIKE '%' || sqlc.narg('status') || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('kehadiran_id')::uuid IS NULL OR kehadiran_id = sqlc.narg('kehadiran_id')::uuid)
  AND (sqlc.narg('skp_intervensi_id')::uuid IS NULL OR skp_intervensi_id = sqlc.narg('skp_intervensi_id')::uuid);

-- name: UpdateKehadiranSkp :one
UPDATE kehadiran_skp
SET
  status       = COALESCE(sqlc.narg('status'), status),
  is_active    = COALESCE(sqlc.narg('is_active'), is_active),
  updated_by   = COALESCE(sqlc.narg('updated_by'), updated_by),
  updated_note = COALESCE(sqlc.narg('updated_note'), updated_note),
  updated_at   = now()
WHERE id = sqlc.arg('id')
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteKehadiranSkp :exec
UPDATE kehadiran_skp
SET deleted_at = now(),
    deleted_by = $2
WHERE id = $1
  AND deleted_at IS NULL;
