-- session-db Indexes
-- Optimasi telemetry & timeseries query meter_readings

CREATE INDEX IF NOT EXISTS idx_sessions_user_status ON charging_sessions(user_id, status);
CREATE INDEX IF NOT EXISTS idx_sessions_booking ON charging_sessions(booking_id);
CREATE INDEX IF NOT EXISTS idx_meter_readings_session_time ON meter_readings(session_id, recorded_at DESC);
CREATE INDEX IF NOT EXISTS idx_outbox_pending_session ON outbox_events(status, created_at) WHERE status = 'PENDING';
