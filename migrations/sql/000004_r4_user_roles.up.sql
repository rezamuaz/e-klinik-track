CREATE TABLE IF NOT EXISTS r4_user_roles (
    id UUID NOT NULL DEFAULT uuidv7(),
    tag VARCHAR NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT r4_user_roles_pkey PRIMARY KEY (id)
);