CREATE TABLE IF NOT EXISTS r2_group_permissions (
    id UUID NOT NULL DEFAULT uuidv7(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT r2_group_permissions_pkey PRIMARY KEY (id)
);