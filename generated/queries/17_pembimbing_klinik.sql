-- name: CreatePembimbingKlinik :one
INSERT INTO pembimbing_klinik (
    fasilitas_id, kontrak_id,user_id, created_by
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPembimbingKlinikByID :one
SELECT * FROM pembimbing_klinik
WHERE id = $1;

-- name: GetAllPembimbingKlinik :many
SELECT * FROM pembimbing_klinik
ORDER BY created_at DESC;

-- name: UpdatePembimbingKlinikPartial :one
UPDATE pembimbing_klinik
SET
  fasilitas_id = COALESCE(sqlc.narg('fasilitas_id'), fasilitas_id),
  kontrak_id   = COALESCE(sqlc.narg('kontrak_id'), kontrak_id),
  user_id      = COALESCE(sqlc.narg('user_id'), user_id),
  is_active    = COALESCE(sqlc.narg('is_active'), is_active),
  updated_by   = COALESCE(sqlc.narg('updated_by'), updated_by),
  updated_note = COALESCE(sqlc.narg('updated_note'), updated_note),
  updated_at   = now()
WHERE id = $1
RETURNING *;

-- name: DeletePembimbingKlinik :exec
DELETE FROM pembimbing_klinik
WHERE id = $1;


-- name: ListPembimbingKlinikByKontrakID :many
SELECT
    t2.id,
    t2.nama,
    t1.is_active
FROM
    pembimbing_klinik t1
JOIN
    users t2 ON t1.user_id = t2.id
WHERE
    t1.kontrak_id = $1
    AND t1.deleted_at IS NULL;