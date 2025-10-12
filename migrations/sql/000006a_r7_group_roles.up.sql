CREATE TABLE IF NOT EXISTS r7_group_roles (
    id UUID NOT NULL DEFAULT uuidv7(),
    group_id INT NOT NULL,
    role_id INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT group_roles_pkey PRIMARY KEY (id)
);