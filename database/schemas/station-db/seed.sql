-- station-db Seed Data

INSERT INTO stations (id, name, location, latitude, longitude, total_power_kw, status) VALUES
('stn-001', 'SPKLU Grand Indonesia', 'Jl. M.H. Thamrin No.1, Jakarta Pusat', -6.19510000, 106.82000000, 200.00, 'ACTIVE'),
('stn-002', 'SPKLU Rest Area KM 19', 'Tol Jakarta-Cikampek KM 19', -6.26250000, 107.02100000, 150.00, 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

INSERT INTO charger_slots (id, station_id, slot_number, connector_type, max_power_kw, status) VALUES
('slot-001', 'stn-001', 1, 'CCS2', 100.00, 'AVAILABLE'),
('slot-002', 'stn-001', 2, 'CCS2', 100.00, 'AVAILABLE'),
('slot-003', 'stn-002', 1, 'TYPE2', 50.00, 'AVAILABLE'),
('slot-004', 'stn-002', 2, 'CHADEMO', 100.00, 'AVAILABLE')
ON CONFLICT (id) DO NOTHING;

INSERT INTO tariffs (id, station_id, price_per_kwh, currency, valid_from, is_active) VALUES
('trf-001', 'stn-001', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-002', 'stn-002', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE)
ON CONFLICT (id) DO NOTHING;
