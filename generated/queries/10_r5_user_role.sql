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
