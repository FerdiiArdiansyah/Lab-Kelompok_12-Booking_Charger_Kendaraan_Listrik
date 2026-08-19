-- station-db Seed Data

INSERT INTO stations (id, name, location, latitude, longitude, total_power_kw, status) VALUES
('stn-001', 'SPKLU Grand Indonesia',         'Jl. M.H. Thamrin No.1, Jakarta Pusat',             -6.19510000, 106.82000000, 200.00, 'ACTIVE'),
('stn-002', 'SPKLU Rest Area KM 19',         'Tol Jakarta-Cikampek KM 19, Bekasi',                -6.26250000, 107.02100000, 150.00, 'ACTIVE'),
('stn-003', 'SPKLU Summarecon Mall Serpong', 'Jl. Boulevard Gading Serpong, Tangerang',           -6.23900000, 106.63000000, 120.00, 'ACTIVE'),
('stn-004', 'SPKLU Bandara Soetta T3',       'Kawasan Bandara Soekarno-Hatta, Tangerang',          -6.12750000, 106.65370000, 200.00, 'ACTIVE'),
('stn-005', 'SPKLU AEON Mall BSD City',      'Jl. BSD Raya Utama, BSD City, Tangerang Selatan',   -6.29890000, 106.64880000, 100.00, 'MAINTENANCE')
ON CONFLICT (id) DO NOTHING;

INSERT INTO charger_slots (id, station_id, slot_number, connector_type, max_power_kw, status) VALUES
('slot-001', 'stn-001', 1, 'CCS2',     100.00, 'AVAILABLE'),
('slot-002', 'stn-001', 2, 'CCS2',     100.00, 'OCCUPIED'),
('slot-003', 'stn-001', 3, 'TYPE2',     50.00, 'AVAILABLE'),
('slot-004', 'stn-001', 4, 'CHADEMO',   50.00, 'AVAILABLE'),
('slot-005', 'stn-002', 1, 'TYPE2',     50.00, 'AVAILABLE'),
('slot-006', 'stn-002', 2, 'CHADEMO',  100.00, 'AVAILABLE'),
('slot-007', 'stn-002', 3, 'CCS2',     100.00, 'MAINTENANCE'),
('slot-008', 'stn-003', 1, 'CCS2',      60.00, 'AVAILABLE'),
('slot-009', 'stn-003', 2, 'TYPE2',     22.00, 'AVAILABLE'),
('slot-010', 'stn-003', 3, 'CCS2',      60.00, 'AVAILABLE'),
('slot-011', 'stn-004', 1, 'CCS2',     100.00, 'AVAILABLE'),
('slot-012', 'stn-004', 2, 'CCS2',     100.00, 'AVAILABLE')
ON CONFLICT (id) DO NOTHING;

INSERT INTO tariffs (id, station_id, price_per_kwh, currency, valid_from, is_active) VALUES
('trf-001', 'stn-001', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-002', 'stn-002', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-003', 'stn-003', 2350.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-004', 'stn-004', 2600.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-005', 'stn-005', 2350.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE)
ON CONFLICT (id) DO NOTHING;
