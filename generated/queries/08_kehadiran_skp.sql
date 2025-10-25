

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
UPDATE public.kehadiran_skp ks
SET
  status = CASE
    WHEN ks.id = ANY(sqlc.arg('skp_kehadiran_id')::uuid[])
      THEN 'disetujui'
    ELSE 'ditolak'
  END,
  locked = TRUE,
  updated_at = now(),
  updated_by = sqlc.arg('updated_by')
WHERE
  ks.kehadiran_id = sqlc.arg('kehadiran_id')::uuid
  AND ks.deleted_at IS NULL
  AND ks.is_active = TRUE;

-- name: GetRekapSKPHarian :one
WITH DataSKP AS (
  SELECT
    COUNT(*) AS total_kompetensi,
    COUNT(*) FILTER (
      WHERE (ks.status = 'verified' OR ks.locked = TRUE)
    ) AS diverifikasi
  FROM kehadiran_skp ks
  JOIN kehadiran k ON k.id = ks.kehadiran_id
  WHERE k.tgl_kehadiran = sqlc.arg('tgl')
    AND ks.is_active = TRUE
    AND k.is_active = TRUE
)
SELECT
  COALESCE(ds.total_kompetensi, 0) AS total_kompetensi_dicatat,
  COALESCE(ds.diverifikasi, 0) AS sudah_diverifikasi_pembimbing,
  ROUND(
    (COALESCE(ds.diverifikasi, 0)::numeric / NULLIF(ds.total_kompetensi, 0)) * 100,
    2
  ) AS persentase_selesai
FROM DataSKP ds;



-- name: GetGlobalSKPPersentase :many
WITH Total AS (
  SELECT COUNT(*) AS total
  FROM kehadiran_skp ks
  JOIN kehadiran k ON k.id = ks.kehadiran_id
  WHERE ks.is_active = TRUE
    AND k.is_active = TRUE
    AND k.tgl_kehadiran BETWEEN (sqlc.arg('tgl')::date - INTERVAL '6 day') AND sqlc.arg('tgl')::date
)
SELECT
  CASE
    WHEN (ks.status = 'verified' OR ks.locked = TRUE) THEN 'Diverifikasi'
    WHEN ks.status = 'rejected' THEN 'Ditolak'
    ELSE 'Belum Diverifikasi'
  END AS kategori,
  COUNT(*) AS jumlah,
  ROUND(
    (COUNT(*)::numeric / NULLIF(t.total, 0)) * 100,
    2
  ) AS persentase
FROM kehadiran_skp ks
JOIN kehadiran k ON k.id = ks.kehadiran_id
CROSS JOIN Total t
WHERE ks.is_active = TRUE
  AND k.is_active = TRUE
  AND k.tgl_kehadiran BETWEEN (sqlc.arg('tgl')::date - INTERVAL '6 day') AND sqlc.arg('tgl')::date
GROUP BY kategori, t.total
ORDER BY persentase DESC;


-- name: GetCapaianSKP7HariTerakhir :many
WITH date_series AS (
  SELECT generate_series(
    (sqlc.arg('tgl')::date - INTERVAL '6 day'),
    sqlc.arg('tgl')::date,
    INTERVAL '1 day'
  )::date AS tanggal
),
rekap AS (
  SELECT
    k.tgl_kehadiran::date AS tanggal,
    COUNT(ks.id) AS total_skp,
    COUNT(*) FILTER (WHERE (ks.status = 'verified' OR ks.locked = TRUE)) AS diverifikasi,
    COUNT(*) FILTER (WHERE ks.status = 'rejected') AS ditolak,
    COUNT(*) FILTER (
      WHERE ks.status IS NULL OR ks.status NOT IN ('verified', 'rejected')
    ) AS belum_diverifikasi
  FROM kehadiran_skp ks
  JOIN kehadiran k ON k.id = ks.kehadiran_id
  WHERE ks.is_active = TRUE
    AND k.is_active = TRUE
    AND k.tgl_kehadiran BETWEEN (sqlc.arg('tgl')::date - INTERVAL '6 day') AND sqlc.arg('tgl')::date
  GROUP BY k.tgl_kehadiran
)
SELECT
  ds.tanggal,
  COALESCE(r.total_skp, 0) AS total_skp,
  COALESCE(r.diverifikasi, 0) AS diverifikasi,
  COALESCE(r.ditolak, 0) AS ditolak,
  COALESCE(r.belum_diverifikasi, 0) AS belum_diverifikasi,
  ROUND(
    (COALESCE(r.diverifikasi, 0)::numeric / NULLIF(r.total_skp, 0)) * 100,
    2
  ) AS persentase_diverifikasi
FROM date_series ds
LEFT JOIN rekap r ON r.tanggal = ds.tanggal
ORDER BY ds.tanggal ASC;


-- name: GetCapaianSKPPerHari :many
SELECT
  k.tgl_kehadiran::date AS tanggal,
  COUNT(ks.id) AS total_skp,
  COUNT(*) FILTER (WHERE (ks.status = 'disetujui' OR ks.locked = TRUE)) AS diverifikasi,
  COUNT(*) FILTER (WHERE ks.status = 'ditolak') AS ditolak,
  COUNT(*) FILTER (
    WHERE ks.status IS NULL OR ks.status NOT IN ('disetujui', 'ditolak')
  ) AS belum_diverifikasi,
  ROUND(
    (COUNT(*) FILTER (WHERE (ks.status = 'disetujui' OR ks.locked = TRUE))::numeric
      / NULLIF(COUNT(ks.id), 0)) * 100,
    2
  ) AS persentase_diverifikasi
FROM kehadiran_skp ks
JOIN kehadiran k ON k.id = ks.kehadiran_id
WHERE ks.is_active = TRUE
  AND k.is_active = TRUE
  AND k.tgl_kehadiran BETWEEN sqlc.arg('start_date')::date AND sqlc.arg('end_date')::date
GROUP BY k.tgl_kehadiran
ORDER BY k.tgl_kehadiran ASC;


-- name: GetRekapSkpTercapaiByUser :many
SELECT
  si.nama AS nama_skp,
  COUNT(ks.id) AS jumlah_tercapai,
  STRING_AGG(DISTINCT to_char(k.tgl_kehadiran, 'DD-MM-YYYY'), ',') AS tanggal_tercapai
FROM kehadiran_skp ks
JOIN kehadiran k ON k.id = ks.kehadiran_id
JOIN skp_intervensi si ON si.id = ks.skp_intervensi_id
WHERE
  ks.user_id = sqlc.arg('user_id')
  AND ks.status = 'disetujui'
  AND ks.is_active = TRUE
  AND k.is_active = TRUE
  AND k.tgl_kehadiran BETWEEN sqlc.arg('tgl_awal') AND sqlc.arg('tgl_akhir')
GROUP BY si.nama
ORDER BY si.nama ASC;




