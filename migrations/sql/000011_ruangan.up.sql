CREATE TABLE IF NOT EXISTS ruangan (
	id UUID PRIMARY KEY  DEFAULT uuidv7(),
	fasilitas_id UUID NOT NULL,
	kontrak_id UUID NOT NULL, 
	nama_ruangan TEXT NOT NULL,
	is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_note TEXT,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now()
	);