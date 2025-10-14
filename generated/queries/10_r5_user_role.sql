-- name: CreateUserRolesBulk :many
INSERT INTO r5_user_roles (
  user_id,
  role_id,
  is_active,
  created_by
)
SELECT
  unnest(sqlc.arg('user_ids')::uuid[]),
  unnest(sqlc.arg('role_ids')::uuid[]),
  unnest(sqlc.arg('is_actives')::bool[]),
  unnest(sqlc.arg('created_bys')::text[])
RETURNING *;


-- name: CreateUserRole :exec
INSERT INTO r5_user_roles (user_id, role_id, is_active, created_by)
VALUES ($1, $2, true, $3)
ON CONFLICT (user_id, role_id)
DO NOTHING;

-- name: DeleteUnRegisterRole :exec
WITH new_roles AS (
    SELECT unnest(sqlc.arg('role_ids')::int[]) AS role_id  -- $2 = array of role_ids dari input user
)
DELETE FROM r5_user_roles ur
WHERE ur.user_id = $1
  AND ur.role_id NOT IN (SELECT role_id FROM new_roles);



-- name: GetUserRolesByUserID :many
SELECT
    ur.id,
    ur.nama
FROM r5_user_roles uur
JOIN r4_roles ur
    ON ur.id = uur.role_id
WHERE uur.user_id = $1
  AND uur.deleted_at IS NULL
  AND ur.deleted_at IS NULL;

-- name: GetUserMenuViews :many
SELECT DISTINCT
    r1.id AS view_id,
    r1.label AS view_label,
    r1.view,
    r1.resource_key,
    r1.action
FROM
    r5_user_roles ur 
JOIN
    r4_roles r4 ON ur.role_id = r4.id
JOIN
    r3_view_roles r3 ON ur.role_id = r3.role_id
JOIN
    r1_views r1 ON r3.view_id = r1.id
WHERE
    ur.user_id = sqlc.arg('user_id')
    
    -- Filter View Murni
    AND r1.view = 'view'
    
    -- Filter Status Aktif
    AND ur.deleted_at IS NULL
    AND r3.deleted_at IS NULL
    AND r1.is_active = TRUE
    AND r4.is_active = TRUE;
