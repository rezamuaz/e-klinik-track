-- name: CreatePembimbingKlinik :one
INSERT INTO pembimbing_klinik (
    fasilitas_id, user_id, is_active, created_by
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPembimbingKlinikByID :one
SELECT * FROM pembimbing_klinik
WHERE id = $1;

-- name: GetAllPembimbingKlinik :many
SELECT * FROM pembimbing_klinik
ORDER BY created_at DESC;

-- name: UpdatePembimbingKlinik :one
UPDATE pembimbing_klinik
SET fasilitas_id = $2,
    user_id = $3,
    is_active = $4,
    updated_note = $5,
    updated_by = $6,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePembimbingKlinik :exec
DELETE FROM pembimbing_klinik
WHERE id = $1;