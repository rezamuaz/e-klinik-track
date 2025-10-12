CREATE TABLE IF NOT EXISTS r3_view_roles (
    id UUID NOT NULL DEFAULT uuidv7(),
    view_id INT NOT NULL,
    role_id INT NOT NULL,
    action VARCHAR,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT r3_view_roles_pkey PRIMARY KEY (id));