CREATE TABLE IF NOT EXISTS user_groups (
    id UUID NOT NULL DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    grup_id VARCHAR,
    tipe VARCHAR(10),
    is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT akses_grup_user_pkey PRIMARY KEY (id)
);