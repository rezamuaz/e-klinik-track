-- name: CreateR3ViewRole :one
INSERT INTO r3_view_roles (
  view_id, role_id, action, created_by
)
VALUES (
  sqlc.arg('view_id'),
  sqlc.arg('role_id'),
  sqlc.narg('action'),
  sqlc.narg('created_by')
)
RETURNING *;

-- name: GetR3ViewRoleByID :one
SELECT *
FROM r3_view_roles
WHERE id = sqlc.arg('id')
  AND deleted_at IS NULL;

-- name: ListR3ViewRoles :many
SELECT *
FROM r3_view_roles
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateR3ViewRole :one
UPDATE r3_view_roles
SET
  view_id = COALESCE(sqlc.narg('view_id'), view_id),
  role_id = COALESCE(sqlc.narg('role_id'), role_id),
  action = COALESCE(sqlc.narg('action'), action),
  updated_by = sqlc.narg('updated_by'),
  updated_at = now()
WHERE id = sqlc.arg('id')
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteR3ViewRole :exec
UPDATE r3_view_roles
SET
  deleted_by = sqlc.narg('deleted_by'),
  deleted_at = now()
WHERE id = sqlc.arg('id')
  AND deleted_at IS NULL;

-- name: ViewRolesSyncDeleteHard :exec
WITH new_views AS (
    -- sqlc.arg('view_ids') adalah array of INTEGER yang dikirim dari input admin
    SELECT unnest(sqlc.arg('view_ids')::int[]) AS view_id  
)
DELETE FROM r3_view_roles r3
WHERE r3.role_id = sqlc.arg('role_id') -- $1 = role_id target
  -- Hapus semua policy yang dimiliki role, tetapi TIDAK ada di daftar view_id baru
  AND r3.view_id NOT IN (SELECT view_id FROM new_views);

-- name: ViewRolesSyncInsertHard :exec
INSERT INTO r3_view_roles (
    view_id, 
    role_id, 
    created_by
)
SELECT 
    view_id, 
    sqlc.arg('role_id'), 
    sqlc.arg('current_user_id')
FROM 
    unnest(sqlc.arg('view_ids')::int[]) AS new_view(view_id) -- Mengambil array view_ids baru
ON CONFLICT (view_id, role_id) DO NOTHING;


-- name: GetR3ViewRoleByRoleID :many
SELECT id,view_id,role_id
FROM r3_view_roles
WHERE role_id = sqlc.arg('role_id')
  AND deleted_at IS NULL;