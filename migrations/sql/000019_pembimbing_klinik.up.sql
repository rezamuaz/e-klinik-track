CREATE TABLE IF NOT EXISTS pembimbing_klinik (
	id UUID PRIMARY KEY  DEFAULT uuidv7(),
    fasilitas_id UUID,
    user_id UUID,
	is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_note TEXT,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now()
);