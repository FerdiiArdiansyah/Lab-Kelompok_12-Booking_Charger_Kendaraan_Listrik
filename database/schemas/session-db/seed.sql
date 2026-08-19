-- session-db Seed Data

INSERT INTO charging_sessions (id, booking_id, slot_id, user_id, started_at, ended_at, consumed_kwh, status) VALUES
('ses-001', 'bkg-000-prev', 'slot-001', 'usr-001', CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '30 minutes', 28.550, 'COMPLETED'),
('ses-002', 'bkg-001', 'slot-001', 'usr-001', CURRENT_TIMESTAMP - INTERVAL '10 minutes', NULL, 3.200, 'IN_PROGRESS')
ON CONFLICT (id) DO NOTHING;

INSERT INTO meter_readings (session_id, recorded_at, current_kwh, power_kw, voltage, current_ampere) VALUES
('ses-002', CURRENT_TIMESTAMP - INTERVAL '10 minutes', 0.000, 48.50, 398.00, 121.80),
('ses-002', CURRENT_TIMESTAMP - INTERVAL '5 minutes', 1.600, 49.10, 399.20, 123.00),
('ses-002', CURRENT_TIMESTAMP, 3.200, 48.80, 398.50, 122.40);
