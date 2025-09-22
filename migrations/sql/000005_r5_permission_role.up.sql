CREATE TABLE IF NOT EXISTS r5_permission_role (
    id UUID NOT NULL DEFAULT uuidv7(),
    permission_id UUID NOT NULL,
    user_roles_id UUID NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT r5_permission_role_pkey PRIMARY KEY (id)
);