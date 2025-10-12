-- name: CreateR2Group :one
INSERT INTO r2_groups (
  name, description, created_by
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetR2GroupByID :one
SELECT * FROM r2_groups
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListR2Groups :many
SELECT * FROM r2_groups
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateR2Group :one
UPDATE r2_groups
SET
  name        = COALESCE(sqlc.narg('name'), name),
  description = COALESCE(sqlc.narg('description'), description),
  is_active   = COALESCE(sqlc.narg('is_active'), is_active),
  updated_by  = sqlc.narg('updated_by'),
  updated_at  = now()
WHERE id = sqlc.arg('id')
RETURNING *;


-- name: DeleteR2Group :exec
UPDATE r2_groups
SET 
  deleted_by = sqlc.arg('deleted_by'),
  deleted_at = now()
WHERE id = sqlc.arg('id');