

-- name: GetKehadiranSkp :one
SELECT *
FROM kehadiran_skp
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ListKehadiranSkp :many
SELECT
  id,
  kehadiran_id,
  skp_intervensi_id,
  status,
  is_active,
  created_by,
  created_at
FROM kehadiran_skp
WHERE deleted_at IS NULL
  AND (sqlc.narg('status')::text IS NULL OR status ILIKE '%' || sqlc.narg('status') || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('kehadiran_id')::uuid IS NULL OR kehadiran_id = sqlc.narg('kehadiran_id')::uuid)
  AND (sqlc.narg('skp_intervensi_id')::uuid IS NULL OR skp_intervensi_id = sqlc.narg('skp_intervensi_id')::uuid)
ORDER BY
  CASE WHEN sqlc.narg('order_by')::text = 'status' AND sqlc.narg('sort')::text = 'asc'  THEN status END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'status' AND sqlc.narg('sort')::text = 'desc' THEN status END DESC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'asc'  THEN created_at END ASC,
  CASE WHEN sqlc.narg('order_by')::text = 'created_at' AND sqlc.narg('sort')::text = 'desc' THEN created_at END DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountKehadiranSkp :one
SELECT COUNT(*)::bigint
FROM kehadiran_skp
WHERE deleted_at IS NULL
  AND (sqlc.narg('status')::text IS NULL OR status ILIKE '%' || sqlc.narg('status') || '%')
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('kehadiran_id')::uuid IS NULL OR kehadiran_id = sqlc.narg('kehadiran_id')::uuid)
  AND (sqlc.narg('skp_intervensi_id')::uuid IS NULL OR skp_intervensi_id = sqlc.narg('skp_intervensi_id')::uuid);

-- name: UpdateKehadiranSkp :one
UPDATE kehadiran_skp
SET
  status       = COALESCE(sqlc.narg('status'), status),
  is_active    = COALESCE(sqlc.narg('is_active'), is_active),
  updated_by   = COALESCE(sqlc.narg('updated_by'), updated_by),
  updated_note = COALESCE(sqlc.narg('updated_note'), updated_note),
  updated_at   = now()
WHERE id = sqlc.arg('id')
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteKehadiranSkp :exec
UPDATE kehadiran_skp
SET deleted_at = now(),
    deleted_by = $2
WHERE id = $1
  AND deleted_at IS NULL;




-- name: SyncKehadiranSkp :many
WITH
-- 1️⃣ Data input user
input_data AS (
    SELECT 
        unnest(@skp_intervensi_ids::uuid[]) AS skp_intervensi_id,
        @kehadiran_id::uuid AS kehadiran_id,
        @user_id::uuid AS user_id,
        @actor::varchar AS actor
),

-- 2️⃣ Insert baru hanya jika:
--   - baris belum ada, atau
--   - belum locked
inserted AS (
    INSERT INTO kehadiran_skp (
        kehadiran_id,
        skp_intervensi_id,
        user_id,
        created_by
    )
    SELECT 
        i.kehadiran_id,
        i.skp_intervensi_id,
        i.user_id,
        i.actor
    FROM input_data i
    LEFT JOIN kehadiran_skp k
      ON k.kehadiran_id = i.kehadiran_id 
     AND k.skp_intervensi_id = i.skp_intervensi_id
    WHERE k.id IS NULL                 -- belum ada baris
       OR k.locked = false             -- atau baris belum terkunci
    ON CONFLICT (kehadiran_id, skp_intervensi_id)
    DO NOTHING
    RETURNING kehadiran_skp.*
),

-- 3️⃣ Hapus baris yang:
--   - tidak ada di input user
--   - tidak locked
deleted AS (
    DELETE FROM kehadiran_skp k
    WHERE k.kehadiran_id = @kehadiran_id
      AND k.locked = false
      AND k.skp_intervensi_id NOT IN (SELECT skp_intervensi_id FROM input_data)
    RETURNING k.*
)

-- 4️⃣ Gabungkan hasil operasi insert + delete
SELECT * FROM inserted
UNION ALL
SELECT * FROM deleted;


-- name: SkpKehadiranID :many
SELECT skp_intervensi_id
FROM kehadiran_skp
WHERE kehadiran_id = $1
  AND is_active = true
  AND deleted_at IS NULL
ORDER BY created_at ASC;


-- name: IntervensiKehadiranID :many
SELECT
ks.id,
  si.nama, 
  ks.skp_intervensi_id,
  ks.locked
FROM public.kehadiran_skp ks
LEFT JOIN public.skp_intervensi si
  ON ks.skp_intervensi_id = si.id
WHERE 
  ks.kehadiran_id = $1
  AND ks.is_active = TRUE
  AND ks.deleted_at IS NULL
ORDER BY ks.created_at DESC;

-- name: ApproveKehadiranSkpByIds :exec
UPDATE public.kehadiran_skp
SET
  status = 'disetujui',
  locked = TRUE,
  updated_at = now(),
  updated_by = $1 -- Parameter opsional untuk menyimpan siapa yang melakukan update (misalnya user_id atau username)
WHERE
  id = ANY (sqlc.arg('skp_kehadiran_id')::uuid[]) -- $1 adalah list UUID yang Anda masukkan
  AND deleted_at IS NULL;

