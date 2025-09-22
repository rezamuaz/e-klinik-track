CREATE TABLE IF NOT EXISTS kontrak (
	id UUID PRIMARY KEY  DEFAULT uuidv7(),
	fasilitas_id UUID NOT NULL,
	nama TEXT NOT NULL, 
	periode_mulai TIMESTAMPTZ,
    periode_selesai TIMESTAMPTZ,
	durasi INTERVAL,
	deskripsi TEXT,
	is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_note TEXT,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now()
	);