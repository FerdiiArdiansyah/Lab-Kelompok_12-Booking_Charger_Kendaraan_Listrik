-- station-db Seed Data
-- Lokasi SPKLU PLN yang nyata beroperasi di Jabodetabek per 2025

INSERT INTO stations (id, name, location, latitude, longitude, total_power_kw, status) VALUES
('stn-001', 'SPKLU PLN Unit Grand Indonesia',        'Jl. M.H. Thamrin No.1, Menteng, Jakarta Pusat 10310',                          -6.19510000, 106.82000000, 200.00, 'ACTIVE'),
('stn-002', 'SPKLU PLN Rest Area KM 19+400 A',       'Tol Jakarta-Cikampek KM 19+400 Arah Cikampek, Bekasi Barat, Jawa Barat',       -6.26180000, 107.01850000, 300.00, 'ACTIVE'),
('stn-003', 'SPKLU PLN Mall Taman Anggrek',          'Jl. Letjen S. Parman Kav 21, Tanjung Duren Utara, Jakarta Barat 11470',        -6.17840000, 106.79030000, 120.00, 'ACTIVE'),
('stn-004', 'SPKLU PLN Summarecon Mall Serpong',     'Jl. Boulevard Gading Serpong, Pakulonan, Serpong Utara, Tangerang 15810',     -6.23900000, 106.63000000, 100.00, 'ACTIVE'),
('stn-005', 'SPKLU PLN Bandara Soekarno-Hatta T3',  'Jl. Raya Bandara Soekarno-Hatta, Selapajang, Neglasari, Tangerang 15126',    -6.12750000, 106.65370000, 200.00, 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

INSERT INTO charger_slots (id, station_id, slot_number, connector_type, max_power_kw, status) VALUES
-- stn-001: Grand Indonesia – 2x CCS2, 1x TYPE2, 1x CHAdeMO
('slot-001', 'stn-001', 1, 'CCS2',     100.00, 'AVAILABLE'),
('slot-002', 'stn-001', 2, 'CCS2',     100.00, 'AVAILABLE'),
('slot-003', 'stn-001', 3, 'TYPE2',     22.00, 'AVAILABLE'),
('slot-004', 'stn-001', 4, 'CHADEMO',   50.00, 'AVAILABLE'),
-- stn-002: Rest Area KM 19+400 – 2x CCS2 150kW, 1x TYPE2 (sedang maintenance)
('slot-005', 'stn-002', 1, 'CCS2',     150.00, 'AVAILABLE'),
('slot-006', 'stn-002', 2, 'CCS2',     150.00, 'AVAILABLE'),
('slot-007', 'stn-002', 3, 'TYPE2',     22.00, 'MAINTENANCE'),
-- stn-003: Mall Taman Anggrek – 1x CCS2, 2x TYPE2 (slot-009 sedang digunakan)
('slot-008', 'stn-003', 1, 'CCS2',     100.00, 'AVAILABLE'),
('slot-009', 'stn-003', 2, 'TYPE2',     22.00, 'OCCUPIED'),
('slot-010', 'stn-003', 3, 'TYPE2',     22.00, 'AVAILABLE'),
-- stn-004: Summarecon Mall Serpong – 1x CCS2, 2x TYPE2
('slot-011', 'stn-004', 1, 'CCS2',      60.00, 'AVAILABLE'),
('slot-012', 'stn-004', 2, 'TYPE2',     22.00, 'AVAILABLE'),
('slot-013', 'stn-004', 3, 'TYPE2',     22.00, 'AVAILABLE'),
-- stn-005: Bandara Soetta T3 – 2x CCS2, 1x CHAdeMO
('slot-014', 'stn-005', 1, 'CCS2',     100.00, 'AVAILABLE'),
('slot-015', 'stn-005', 2, 'CCS2',     100.00, 'AVAILABLE'),
('slot-016', 'stn-005', 3, 'CHADEMO',   50.00, 'AVAILABLE')
ON CONFLICT (id) DO NOTHING;

-- Tarif PLN SPKLU berlaku: Rp 2.466,77/kWh (dibulatkan Rp 2.467), bandara tarif premium
INSERT INTO tariffs (id, station_id, price_per_kwh, currency, valid_from, is_active) VALUES
('trf-001', 'stn-001', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-002', 'stn-002', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-003', 'stn-003', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-004', 'stn-004', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-005', 'stn-005', 2600.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE)
ON CONFLICT (id) DO NOTHING;
