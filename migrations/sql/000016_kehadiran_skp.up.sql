CREATE TABLE IF NOT EXISTS kehadiran_skp (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    kehadiran_id UUID NOT NULL REFERENCES kehadiran(id) ON DELETE CASCADE,
    skp_intervensi_id UUID NOT NULL REFERENCES skp_intervensi(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    status VARCHAR,
    is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    updated_note TEXT,
    updated_by VARCHAR,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now()
);


ALTER TABLE public.kehadiran_skp
ADD CONSTRAINT kehadiran_skp_unique UNIQUE (kehadiran_id, skp_intervensi_id);