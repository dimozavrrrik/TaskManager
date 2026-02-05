-- Refresh tokens table for JWT refresh token management
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP WITH TIME ZONE,
    user_agent VARCHAR(500),
    ip_address VARCHAR(45)
);

-- Indexes for performance
CREATE INDEX idx_refresh_tokens_employee ON refresh_tokens(employee_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at);
CREATE INDEX idx_refresh_tokens_revoked ON refresh_tokens(revoked_at);

-- Composite index for looking up active tokens
CREATE INDEX idx_refresh_tokens_active ON refresh_tokens(employee_id, expires_at, revoked_at)
    WHERE revoked_at IS NULL;
