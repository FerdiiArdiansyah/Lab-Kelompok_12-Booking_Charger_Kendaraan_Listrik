-- Pattern 2: Idempotency Key Store
-- Digunakan untuk menjamin request berulang (misal retry karena timeout network) tidak memicu double booking atau double payment.

CREATE TABLE IF NOT EXISTS idempotency_keys (
    key VARCHAR(255) PRIMARY KEY,
    request_hash VARCHAR(64) NOT NULL, -- SHA256 dari request body/params
    response_payload JSONB,
    status_code INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_idempotency_expiry ON idempotency_keys(expires_at);

-- Cleanup query untuk menghapus key yang sudah kadaluarsa
/*
  DELETE FROM idempotency_keys WHERE expires_at < CURRENT_TIMESTAMP;
*/
