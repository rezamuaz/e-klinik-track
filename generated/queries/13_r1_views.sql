-- name: CreateR1View :one
INSERT INTO r1_views (
    label, level, parent_id, path, method,
    resource_key, action, view, data,
    is_active, created_by
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9,
    COALESCE($10, true), $11
)
RETURNING *;

-- name: GetR1View :one
SELECT *
FROM r1_views
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListR1Views :many
SELECT *
FROM r1_views
WHERE deleted_at IS NULL 
AND (sqlc.narg('label')::text IS NULL OR label ILIKE '%' || sqlc.narg('label') || '%')
ORDER BY id;

-- name: UpdateR1View :one
UPDATE r1_views
SET 
    label = COALESCE($1, label),
    level = COALESCE($2, level),
    parent_id = COALESCE($3, parent_id),
    path = COALESCE($4, path),
    method = COALESCE($5, method),
    resource_key = COALESCE($6, resource_key),
    action = COALESCE($7, action),
    view = COALESCE($8, view),
    data = COALESCE($9, data),
    is_active = COALESCE($10, is_active),
    updated_by = $11,
    updated_at = now()
WHERE id = $12
RETURNING *;

-- name: DeleteR1View :exec
UPDATE r1_views
SET deleted_at = now(),
    deleted_by = $2
WHERE id = $1 AND deleted_at IS NULL;

-- name: HardDeleteR1View :exec
DELETE FROM r1_views
WHERE id = $1;


-- name: GetR1ViewRecursive :many
SELECT 
    id,
    label,
    parent_id,
    resource_key,
    action,
    view,
    data,
    level,
    path
FROM r1_views
WHERE view = 'view'
ORDER BY level, id;


-- name: GetViewsByIdsWithChildren :many
SELECT v.id,
       v.resource_key,
       v.action,
       v.view
FROM r1_views v
WHERE (v.resource_key, v.action) IN (
    SELECT rv.resource_key, rv.action
    FROM r1_views rv
    WHERE rv.id = ANY(@ids::int[])
)
AND v.is_active = true
AND v.deleted_at IS NULL;
