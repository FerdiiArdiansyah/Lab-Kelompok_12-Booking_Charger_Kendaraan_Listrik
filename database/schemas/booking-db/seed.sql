-- booking-db Seed Data
-- Skenario booking realistis: connector EV harus cocok dengan tipe slot yang dipilih

INSERT INTO bookings (id, user_id, station_id, slot_id, start_time, end_time, status, idempotency_key) VALUES
-- Ioniq 5 (CCS2) – DC fast charge pagi hari di Grand Indonesia slot-001 (100kW)
('bkg-001', 'usr-001', 'stn-001', 'slot-001', CURRENT_TIMESTAMP - INTERVAL '3 hours 30 minutes', CURRENT_TIMESTAMP - INTERVAL '2 hours 30 minutes', 'COMPLETED',       'key-bkg-001'),
-- Wuling Air ev (TYPE2) – AC charge sambil belanja di Taman Anggrek slot-009 (22kW)
('bkg-002', 'usr-002', 'stn-003', 'slot-009', CURRENT_TIMESTAMP - INTERVAL '1 hour',              CURRENT_TIMESTAMP + INTERVAL '1 hour 30 minutes',  'IN_SESSION',      'key-bkg-002'),
-- Tesla Model 3 (CCS2) – DC fast charge di Rest Area KM 19 slot-005 (150kW) perjalanan Bandung
('bkg-003', 'usr-003', 'stn-002', 'slot-005', CURRENT_TIMESTAMP + INTERVAL '1 hour',              CURRENT_TIMESTAMP + INTERVAL '2 hours',             'CONFIRMED',       'key-bkg-003'),
-- BYD Atto 3 (CCS2) – DC fast charge siang di Grand Indonesia slot-002 (100kW)
('bkg-004', 'usr-004', 'stn-001', 'slot-002', CURRENT_TIMESTAMP + INTERVAL '3 hours',             CURRENT_TIMESTAMP + INTERVAL '4 hours',             'CONFIRMED',       'key-bkg-004'),
-- Toyota bZ4X (TYPE2) – AC charge di Summarecon Mall Serpong slot-012 (22kW)
('bkg-005', 'usr-005', 'stn-004', 'slot-012', CURRENT_TIMESTAMP + INTERVAL '30 minutes',          CURRENT_TIMESTAMP + INTERVAL '3 hours',             'CONFIRMED',       'key-bkg-005'),
-- Nissan Leaf (CHAdeMO) – DC fast charge dini hari di Grand Indonesia slot-004 (50kW)
('bkg-006', 'usr-006', 'stn-001', 'slot-004', CURRENT_TIMESTAMP - INTERVAL '6 hours',             CURRENT_TIMESTAMP - INTERVAL '4 hours 30 minutes',  'COMPLETED',       'key-bkg-006'),
-- BMW iX3 (CCS2) – DC fast charge sebelum terbang di Bandara Soetta T3 slot-014 (100kW)
('bkg-007', 'usr-007', 'stn-005', 'slot-014', CURRENT_TIMESTAMP + INTERVAL '4 hours',             CURRENT_TIMESTAMP + INTERVAL '5 hours 30 minutes',  'CONFIRMED',       'key-bkg-007'),
-- MG ZS EV (CCS2) – DC fast charge di Rest Area KM 19 slot-006 (150kW)
('bkg-008', 'usr-008', 'stn-002', 'slot-006', CURRENT_TIMESTAMP + INTERVAL '1 hour 30 minutes',   CURRENT_TIMESTAMP + INTERVAL '2 hours 30 minutes',  'CONFIRMED',       'key-bkg-008'),
-- Kona Electric (CCS2) – booking tengah malam, tidak check-in (no-show) slot-001
('bkg-009', 'usr-009', 'stn-001', 'slot-001', CURRENT_TIMESTAMP - INTERVAL '8 hours',             CURRENT_TIMESTAMP - INTERVAL '6 hours 30 minutes',  'EXPIRED_NO_SHOW', 'key-bkg-009'),
-- Wuling Almaz RS EV (TYPE2) – AC charge sore di Summarecon slot-013 (22kW)
('bkg-010', 'usr-010', 'stn-004', 'slot-013', CURRENT_TIMESTAMP + INTERVAL '5 hours',             CURRENT_TIMESTAMP + INTERVAL '8 hours',             'REQUESTED',       'key-bkg-010')
ON CONFLICT (id) DO NOTHING;

INSERT INTO waitlists (id, station_id, user_id, requested_start, requested_end, queue_number, status) VALUES
('wt-001', 'stn-001', 'usr-010', CURRENT_TIMESTAMP + INTERVAL '1 hour',  CURRENT_TIMESTAMP + INTERVAL '3 hours',  1, 'WAITING'),
('wt-002', 'stn-003', 'usr-007', CURRENT_TIMESTAMP,                      CURRENT_TIMESTAMP + INTERVAL '2 hours',  1, 'WAITING'),
('wt-003', 'stn-002', 'usr-009', CURRENT_TIMESTAMP + INTERVAL '2 hours', CURRENT_TIMESTAMP + INTERVAL '4 hours',  1, 'WAITING')
ON CONFLICT (id) DO NOTHING;

INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload, status) VALUES
('evt-bkg-001', 'Booking', 'bkg-001', 'BookingCompleted',     '{"bookingId":"bkg-001","userId":"usr-001","slotId":"slot-001","stationId":"stn-001"}', 'PUBLISHED'),
('evt-bkg-002', 'Booking', 'bkg-002', 'BookingInSession',     '{"bookingId":"bkg-002","userId":"usr-002","slotId":"slot-009","stationId":"stn-003"}', 'PUBLISHED'),
('evt-bkg-003', 'Booking', 'bkg-003', 'BookingConfirmed',     '{"bookingId":"bkg-003","userId":"usr-003","slotId":"slot-005","stationId":"stn-002"}', 'PUBLISHED'),
('evt-bkg-006', 'Booking', 'bkg-006', 'BookingCompleted',     '{"bookingId":"bkg-006","userId":"usr-006","slotId":"slot-004","stationId":"stn-001"}', 'PUBLISHED'),
('evt-bkg-009', 'Booking', 'bkg-009', 'BookingExpiredNoShow', '{"bookingId":"bkg-009","userId":"usr-009","slotId":"slot-001","stationId":"stn-001"}', 'PUBLISHED')
ON CONFLICT (id) DO NOTHING;
