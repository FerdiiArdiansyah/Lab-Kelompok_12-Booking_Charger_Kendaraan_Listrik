-- session-db DDL Schema
-- Service: session-service
-- Database: session_db

CREATE TABLE IF NOT EXISTS charging_sessions (
    id VARCHAR(36) PRIMARY KEY,
    booking_id VARCHAR(36) NOT NULL UNIQUE,
    slot_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMPTZ,
    consumed_kwh DECIMAL(10, 3) NOT NULL DEFAULT 0.000,
    status VARCHAR(50) NOT NULL DEFAULT 'IN_PROGRESS', -- IN_PROGRESS, COMPLETED, INTERRUPTED, CANCELLED
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS meter_readings (
    id BIGSERIAL PRIMARY KEY,
    session_id VARCHAR(36) NOT NULL REFERENCES charging_sessions(id) ON DELETE CASCADE,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    current_kwh DECIMAL(10, 3) NOT NULL,
    power_kw DECIMAL(8, 2) NOT NULL,
    voltage DECIMAL(6, 2),
    current_ampere DECIMAL(6, 2)
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

CREATE OR REPLACE FUNCTION update_charging_sessions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_charging_sessions_updated_at BEFORE UPDATE ON charging_sessions
FOR EACH ROW EXECUTE PROCEDURE update_charging_sessions_updated_at();
