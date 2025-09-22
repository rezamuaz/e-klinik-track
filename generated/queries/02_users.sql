-- name: UsersFindByUsername :one
SELECT u.id,
       u.username,
       u.password,
	     u.nama,
       ur.id  AS role_id,
       ur.tag AS role 
       FROM users u
LEFT JOIN r4_user_roles ur ON ur.id = u.role
WHERE u.username = $1;

-- name: UsersFindById :one
SELECT u.id,
       u.username,
       u.password,
	     u.nama,
       u.refresh,
       ur.id  AS role_id,
       ur.tag AS role 
       FROM users u
LEFT JOIN r4_user_roles ur ON ur.id = u.role
WHERE u.id = $1;

-- name: ListUsers :many
SELECT
  id,
  nama,
  username,
  role,
  last_active,
  is_active,
  locked_until,
  failed_attempts,
  last_failed_at,
  refresh,
  created_by,
  created_at
FROM public.users
WHERE deleted_at IS NULL
  AND (sqlc.narg('nama')::text IS NULL OR nama ILIKE '%' || sqlc.narg('nama') || '%')
  AND (sqlc.narg('username')::text IS NULL OR username ILIKE '%' || sqlc.narg('username') || '%')
  AND (sqlc.narg('role')::uuid IS NULL OR role = sqlc.narg('role')::uuid)
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
ORDER BY
  CASE WHEN sqlc.narg('order_by')::text = 'nama' AND sqlc.narg('sort')::text = 'asc'  THEN nama END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'nama' AND sqlc.narg('sort')::text = 'desc' THEN nama END DESC,
  CASE WHEN sqlc.narg('order_by')::text = 'username' AND sqlc.narg('sort')::text = 'asc'  THEN username END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'username' AND sqlc.narg('sort')::text = 'desc' THEN username END DESC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'asc'  THEN created_at END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'desc' THEN created_at END DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountUsers :one
SELECT COUNT(*)::bigint
FROM public.users
WHERE deleted_at IS NULL
  AND (sqlc.narg('nama')::text IS NULL OR nama ILIKE '%' || sqlc.narg('nama') || '%')
  AND (sqlc.narg('username')::text IS NULL OR username ILIKE '%' || sqlc.narg('username') || '%')
  AND (sqlc.narg('role')::uuid IS NULL OR role = sqlc.narg('role')::uuid)
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean);


-- name: CreateOrUpdateUser :one
INSERT INTO users (
  nama,
  username,
  password,
  role
) VALUES (
  @nama,@username,@password,@role
) ON CONFLICT (username) DO UPDATE SET 
nama = @nama,
password = @password
RETURNING *, CASE WHEN xmax = 0 THEN 'inserted' ELSE 'updated' END as operation;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUserPartial :exec
UPDATE users
SET 
    nama        = COALESCE(sqlc.narg('nama'), nama),
    is_active   = COALESCE(sqlc.narg('is_active'), is_active),
    password    = COALESCE(sqlc.narg('password'), password),
    refresh     = COALESCE(sqlc.narg('refresh'), refresh),
    updated_note= COALESCE(sqlc.narg('updated_note'), updated_note),
    updated_by  = COALESCE(sqlc.narg('updated_by'), updated_by),
    updated_at  = now()
WHERE id = @id;

-- name: GetUserActiveStatus :one
SELECT is_active FROM users
WHERE username = $1 LIMIT 1;

-- name: GetByUsername :one
SELECT username,password FROM users
WHERE username = $1;

-- name: UpdateUserActive :exec
UPDATE users SET is_active = $1 WHERE username = $2;

-- name: SoftDelUser :exec
UPDATE users SET deleted_at = NOW() WHERE username = $1;

-- name: DelUser :exec
DELETE FROM users
WHERE id = $1;

