-- name: CreateGroupRole :one
INSERT INTO r7_group_roles (
  group_id,
  role_id,
  created_by
) VALUES (
  $1, $2, $3
)
RETURNING *;