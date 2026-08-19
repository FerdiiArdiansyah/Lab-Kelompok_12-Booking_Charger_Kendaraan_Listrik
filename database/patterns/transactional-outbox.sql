-- Pattern 1: Transactional Outbox Pattern
-- Digunakan di setiap microservice database untuk menjamin dual-write problem teratasi (DB commit & Event Publish atomic).

-- 1. Skema Tabel Outbox Standard
CREATE TABLE IF NOT EXISTS outbox_events (
    id VARCHAR(36) PRIMARY KEY,
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id VARCHAR(36) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- PENDING, PUBLISHED, FAILED
    retry_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMPTZ
);

-- Partial Index untuk Polling Worker agar ultra-fast (O(1) lookups)
CREATE INDEX IF NOT EXISTS idx_outbox_pending_events 
ON outbox_events (created_at ASC) 
WHERE status = 'PENDING';

-- 2. Contoh Query Polling oleh Outbox Relay / Debezium / Worker Service dengan ROW LOCK (FOR UPDATE SKIP LOCKED)
/*
  SELECT id, aggregate_type, aggregate_id, event_type, payload
  FROM outbox_events
  WHERE status = 'PENDING'
  ORDER BY created_at ASC
  LIMIT 50
  FOR UPDATE SKIP LOCKED;
*/

-- 3. Query Cleanup Outbox Lawas (Retention Strategy: hapus event PUBLISHED > 7 hari)
/*
  DELETE FROM outbox_events
  WHERE status = 'PUBLISHED' AND processed_at < NOW() - INTERVAL '7 days';
*/
