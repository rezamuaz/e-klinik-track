CREATE TABLE IF NOT EXISTS r1_permissions (
    id UUID NOT NULL DEFAULT uuidv7(),
    label VARCHAR(60) NOT NULL,
    level SMALLINT,
    parent_id VARCHAR(20),
    route VARCHAR(255),
    method VARCHAR(8),
    type VARCHAR(10) DEFAULT 'data',
    is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT r1_permissions_pkey PRIMARY KEY (id)
);