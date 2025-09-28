CREATE TABLE IF NOT EXISTS kehadiran (
	id UUID PRIMARY KEY  DEFAULT uuidv7(),
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
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now()
);