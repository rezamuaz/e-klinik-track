CREATE TABLE IF NOT EXISTS r6_user_groups (
    id UUID NOT NULL DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    grup_id INT NOT NULL,
    tipe VARCHAR(10),
    is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT user_groups_pkey PRIMARY KEY (id)
);