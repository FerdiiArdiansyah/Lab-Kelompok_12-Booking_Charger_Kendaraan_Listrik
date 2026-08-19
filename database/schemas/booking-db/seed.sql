-- booking-db Seed Data

INSERT INTO bookings (id, user_id, station_id, slot_id, start_time, end_time, status, idempotency_key) VALUES
('bkg-001', 'usr-001', 'stn-001', 'slot-001', CURRENT_TIMESTAMP + INTERVAL '1 hour', CURRENT_TIMESTAMP + INTERVAL '3 hours', 'CONFIRMED', 'key-bkg-001'),
('bkg-002', 'usr-002', 'stn-001', 'slot-002', CURRENT_TIMESTAMP + INTERVAL '2 hours', CURRENT_TIMESTAMP + INTERVAL '4 hours', 'CONFIRMED', 'key-bkg-002')
ON CONFLICT (id) DO NOTHING;

INSERT INTO waitlists (id, station_id, user_id, requested_start, requested_end, queue_number, status) VALUES
('wt-001', 'stn-001', 'usr-003', CURRENT_TIMESTAMP + INTERVAL '1 hour', CURRENT_TIMESTAMP + INTERVAL '3 hours', 1, 'WAITING')
ON CONFLICT (id) DO NOTHING;

INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload, status) VALUES
('evt-001', 'Booking', 'bkg-001', 'BookingConfirmed', '{"bookingId": "bkg-001", "userId": "usr-001", "slotId": "slot-001"}', 'PUBLISHED')
ON CONFLICT (id) DO NOTHING;
