CREATE TABLE IF NOT EXISTS mata_kuliah (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	mata_kuliah TEXT NOT NULL,
	is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_note TEXT,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now()
	);