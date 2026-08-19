-- session-db Seed Data

INSERT INTO charging_sessions (id, booking_id, slot_id, user_id, started_at, ended_at, consumed_kwh, status) VALUES
('ses-001', 'bkg-001', 'slot-001', 'usr-001', CURRENT_TIMESTAMP - INTERVAL '4 hours',    CURRENT_TIMESTAMP - INTERVAL '2 hours', 24.500, 'COMPLETED'),
('ses-002', 'bkg-002', 'slot-002', 'usr-002', CURRENT_TIMESTAMP - INTERVAL '30 minutes', NULL,                                   8.200,  'IN_PROGRESS'),
('ses-003', 'bkg-006', 'slot-003', 'usr-006', CURRENT_TIMESTAMP - INTERVAL '6 hours',    CURRENT_TIMESTAMP - INTERVAL '4 hours', 18.750, 'COMPLETED')
ON CONFLICT (id) DO NOTHING;

INSERT INTO meter_readings (session_id, recorded_at, current_kwh, power_kw, voltage, current_ampere) VALUES
-- ses-001
('ses-001', CURRENT_TIMESTAMP - INTERVAL '4 hours',              0.000, 50.00, 400.00, 125.00),
('ses-001', CURRENT_TIMESTAMP - INTERVAL '3 hours',             12.250, 50.00, 399.50, 125.20),
('ses-001', CURRENT_TIMESTAMP - INTERVAL '2 hours 10 minutes',  24.500,  8.50, 396.00,  21.40),
-- ses-002 (ongoing)
('ses-002', CURRENT_TIMESTAMP - INTERVAL '30 minutes',           0.000, 48.00, 398.00, 120.60),
('ses-002', CURRENT_TIMESTAMP - INTERVAL '15 minutes',           4.100, 48.50, 399.20, 121.50),
('ses-002', CURRENT_TIMESTAMP,                                   8.200, 49.00, 400.10, 122.40),
-- ses-003
('ses-003', CURRENT_TIMESTAMP - INTERVAL '6 hours',              0.000, 37.00, 380.00,  97.40),
('ses-003', CURRENT_TIMESTAMP - INTERVAL '5 hours',              9.375, 37.50, 381.00,  98.50),
('ses-003', CURRENT_TIMESTAMP - INTERVAL '4 hours 10 minutes',  18.750,  5.20, 376.00,  13.80);
