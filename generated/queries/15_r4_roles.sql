-- name: CreateR4Role :one
INSERT INTO r4_roles (
  id, tag, nama, created_by,created_at
)
VALUES (
  sqlc.arg('id'),
  sqlc.arg('tag'),
  sqlc.arg('nama'),
  sqlc.narg('created_by'),
  now()
)
RETURNING *;

-- name: GetR4RoleByID :one
SELECT *
FROM r4_roles
WHERE id = sqlc.arg('id')
  AND deleted_at IS NULL;

-- name: ListR4Roles :many
SELECT *
FROM r4_roles
WHERE deleted_at IS NULL
ORDER BY id;

-- name: UpdateR4Role :one
UPDATE r4_roles
SET
  tag = COALESCE(sqlc.narg('tag'), tag),
  nama = COALESCE(sqlc.narg('nama'), nama),
  is_active = COALESCE(sqlc.narg('is_active'), is_active),
  updated_by = sqlc.narg('updated_by'),
  updated_at = now()
WHERE id = sqlc.arg('id')
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteR4Role :exec
UPDATE r4_roles
SET
  deleted_by = sqlc.narg('deleted_by'),
  deleted_at = now()
WHERE id = sqlc.arg('id')
  AND deleted_at IS NULL;