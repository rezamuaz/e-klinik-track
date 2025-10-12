-- name: ListResourceMappings :many
SELECT 
    path, 
    method, 
    resource_key, 
    action
FROM 
    r1_views WHERE view = 'data'
ORDER BY 
    path, method;