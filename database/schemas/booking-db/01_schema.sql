-- booking-db DDL Schema
-- Service: booking-service
-- Database: booking_db

CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS bookings (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    station_id VARCHAR(36) NOT NULL,
    slot_id VARCHAR(36) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    booking_period TSTZRANGE GENERATED ALWAYS AS (tstzrange(start_time, end_time, '[)')) STORED,
    status VARCHAR(50) NOT NULL DEFAULT 'REQUESTED', 
    -- Status: REQUESTED, CONFIRMED, WAITLISTED, IN_SESSION, COMPLETED, EXPIRED_NO_SHOW, CANCELLED
    idempotency_key VARCHAR(255) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_booking_time CHECK (end_time > start_time)
);

-- Constraint Anti-Overlap: Mencegah 2 booking aktif beririsan rentang waktunya pada slot_id yang sama
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS no_overlapping_bookings;
ALTER TABLE bookings ADD CONSTRAINT no_overlapping_bookings
EXCLUDE USING GIST (
    slot_id WITH =,
    booking_period WITH &&
) WHERE (status IN ('REQUESTED', 'CONFIRMED', 'IN_SESSION'));

CREATE TABLE IF NOT EXISTS waitlists (
    id VARCHAR(36) PRIMARY KEY,
    station_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    requested_start TIMESTAMPTZ NOT NULL,
    requested_end TIMESTAMPTZ NOT NULL,
    queue_number INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'WAITING', -- WAITING, PROMOTED, EXPIRED, CANCELLED
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS idempotency_keys (
    key VARCHAR(255) PRIMARY KEY,
    request_hash VARCHAR(64) NOT NULL,
    response_payload JSONB,
    status_code INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS outbox_events (
    id VARCHAR(36) PRIMARY KEY,
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id VARCHAR(36) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- PENDING, PUBLISHED, FAILED
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMPTZ
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_bookings_user ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_slot_status ON bookings(slot_id, status);
CREATE INDEX IF NOT EXISTS idx_waitlists_queue ON waitlists(station_id, queue_number, status);
CREATE INDEX IF NOT EXISTS idx_outbox_pending ON outbox_events(status, created_at) WHERE status = 'PENDING';
