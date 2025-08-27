-- SQL migration for Postgres
CREATE TABLE IF NOT EXISTS otps (
id UUID PRIMARY KEY,
tenant_id TEXT NOT NULL,
channel TEXT NOT NULL CHECK (channel IN ('sms','email')),
destination TEXT NOT NULL,
code TEXT NOT NULL,
attempts INT NOT NULL DEFAULT 0,
max_attempts INT NOT NULL DEFAULT 3,
expires_at TIMESTAMPTZ NOT NULL,
used_at TIMESTAMPTZ NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_otps_tenant ON otps(tenant_id);
CREATE INDEX IF NOT EXISTS idx_otps_expires ON otps(expires_at);