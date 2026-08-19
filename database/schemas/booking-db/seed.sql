-- booking-db Seed Data

INSERT INTO bookings (id, user_id, station_id, slot_id, start_time, end_time, status, idempotency_key) VALUES
('bkg-001', 'usr-001', 'stn-001', 'slot-001', CURRENT_TIMESTAMP - INTERVAL '4 hours',    CURRENT_TIMESTAMP - INTERVAL '2 hours',    'COMPLETED',       'key-bkg-001'),
('bkg-002', 'usr-002', 'stn-001', 'slot-002', CURRENT_TIMESTAMP - INTERVAL '30 minutes', CURRENT_TIMESTAMP + INTERVAL '90 minutes', 'IN_SESSION',      'key-bkg-002'),
('bkg-003', 'usr-003', 'stn-002', 'slot-005', CURRENT_TIMESTAMP + INTERVAL '1 hour',     CURRENT_TIMESTAMP + INTERVAL '3 hours',    'CONFIRMED',       'key-bkg-003'),
('bkg-004', 'usr-004', 'stn-002', 'slot-006', CURRENT_TIMESTAMP + INTERVAL '2 hours',    CURRENT_TIMESTAMP + INTERVAL '4 hours',    'CONFIRMED',       'key-bkg-004'),
('bkg-005', 'usr-005', 'stn-003', 'slot-008', CURRENT_TIMESTAMP + INTERVAL '30 minutes', CURRENT_TIMESTAMP + INTERVAL '150 minutes','CONFIRMED',       'key-bkg-005'),
('bkg-006', 'usr-006', 'stn-001', 'slot-003', CURRENT_TIMESTAMP - INTERVAL '6 hours',    CURRENT_TIMESTAMP - INTERVAL '4 hours',    'COMPLETED',       'key-bkg-006'),
('bkg-007', 'usr-007', 'stn-001', 'slot-001', CURRENT_TIMESTAMP + INTERVAL '3 hours',    CURRENT_TIMESTAMP + INTERVAL '5 hours',    'CONFIRMED',       'key-bkg-007'),
('bkg-008', 'usr-008', 'stn-004', 'slot-011', CURRENT_TIMESTAMP + INTERVAL '1 hour',     CURRENT_TIMESTAMP + INTERVAL '3 hours',    'CONFIRMED',       'key-bkg-008'),
('bkg-009', 'usr-009', 'stn-001', 'slot-004', CURRENT_TIMESTAMP - INTERVAL '8 hours',    CURRENT_TIMESTAMP - INTERVAL '6 hours',    'EXPIRED_NO_SHOW', 'key-bkg-009'),
('bkg-010', 'usr-010', 'stn-003', 'slot-010', CURRENT_TIMESTAMP + INTERVAL '4 hours',    CURRENT_TIMESTAMP + INTERVAL '6 hours',    'REQUESTED',       'key-bkg-010')
ON CONFLICT (id) DO NOTHING;

INSERT INTO waitlists (id, station_id, user_id, requested_start, requested_end, queue_number, status) VALUES
('wt-001', 'stn-001', 'usr-004', CURRENT_TIMESTAMP + INTERVAL '30 minutes', CURRENT_TIMESTAMP + INTERVAL '150 minutes', 1, 'WAITING'),
('wt-002', 'stn-001', 'usr-010', CURRENT_TIMESTAMP + INTERVAL '1 hour',     CURRENT_TIMESTAMP + INTERVAL '3 hours',     2, 'WAITING'),
('wt-003', 'stn-004', 'usr-009', CURRENT_TIMESTAMP + INTERVAL '2 hours',    CURRENT_TIMESTAMP + INTERVAL '4 hours',     1, 'WAITING')
ON CONFLICT (id) DO NOTHING;

INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload, status) VALUES
('evt-bkg-001', 'Booking', 'bkg-001', 'BookingCompleted',     '{"bookingId":"bkg-001","userId":"usr-001","slotId":"slot-001"}', 'PUBLISHED'),
('evt-bkg-002', 'Booking', 'bkg-002', 'BookingInSession',     '{"bookingId":"bkg-002","userId":"usr-002","slotId":"slot-002"}', 'PUBLISHED'),
('evt-bkg-003', 'Booking', 'bkg-003', 'BookingConfirmed',     '{"bookingId":"bkg-003","userId":"usr-003","slotId":"slot-005"}', 'PUBLISHED'),
('evt-bkg-006', 'Booking', 'bkg-006', 'BookingCompleted',     '{"bookingId":"bkg-006","userId":"usr-006","slotId":"slot-003"}', 'PUBLISHED'),
('evt-bkg-009', 'Booking', 'bkg-009', 'BookingExpiredNoShow', '{"bookingId":"bkg-009","userId":"usr-009","slotId":"slot-004"}', 'PUBLISHED')
ON CONFLICT (id) DO NOTHING;
