CREATE TABLE IF NOT EXISTS r5_user_roles (
    id UUID NOT NULL DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    role_id INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT r5_user_role_pkey PRIMARY KEY (id),
    CONSTRAINT r5_user_role_ukey UNIQUE(user_id, role_id) 
);