CREATE TABLE IF NOT EXISTS kehadiran (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    fasilitas_id UUID NOT NULL,
    kontrak_id UUID NOT NULL, 
    ruangan_id UUID NOT NULL,
    mata_kuliah_id UUID NOT NULL,
    pembimbing_id UUID NOT NULL,
    jadwal_dinas VARCHAR,
    user_id UUID NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_note TEXT,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    tgl_kehadiran DATE NOT NULL,
    presensi VARCHAR NOT NULL DEFAULT 'hadir',
    created_by VARCHAR,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE public.kehadiran
ADD CONSTRAINT uq_kehadiran_per_hari
UNIQUE (user_id, tanggal_kehadiran);