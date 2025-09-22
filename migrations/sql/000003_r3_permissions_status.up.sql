CREATE TABLE IF NOT EXISTS r3_permission_status (
    id UUID NOT NULL DEFAULT uuidv7(),
    permission_id UUID NOT NULL,
    group_id UUID NOT NULL,
    value BOOLEAN,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT r3_permission_status_pkey PRIMARY KEY (id),
    CONSTRAINT r3_group_permissions_ukey UNIQUE (permission_id, group_id)
);