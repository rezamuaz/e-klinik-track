CREATE TABLE IF NOT EXISTS user_logs (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    entity_name VARCHAR(50) NOT NULL,   -- type of object: post, comment, profile, etc.
    entity_id UUID,       
    action VARCHAR(100) NOT NULL,
    description TEXT,
    ip_address INET,
    user_agent TEXT,
    device VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    meta JSONB,
    username VARCHAR,
    CONSTRAINT user_logs_users_fkey FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_user_logs_user_id
ON user_logs (user_id ASC);

CREATE INDEX IF NOT EXISTS idx_user_logs_created_at
ON user_logs (created_at ASC);

CREATE INDEX IF NOT EXISTS idx_user_logs_action
ON user_logs (action ASC);