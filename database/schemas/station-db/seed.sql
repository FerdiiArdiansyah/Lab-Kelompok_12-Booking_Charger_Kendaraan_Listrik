-- station-db Seed Data
-- Lokasi SPKLU PLN yang nyata beroperasi di Indonesia (termasuk Sulawesi Selatan) per 2026

INSERT INTO stations (id, name, location, latitude, longitude, total_power_kw, status) VALUES
('stn-001', 'SPKLU PLN UID Jakarta Gambir',                    'Jl. M.I. Ridwan Rais No.1, Gambir, Jakarta Pusat 10110',                                       -6.18220000, 106.83440000, 200.00, 'ACTIVE'),
('stn-002', 'SPKLU PLN Senayan City Mall',                     'Jl. Asia Afrika No.19, Gelora, Tanah Abang, Jakarta Selatan 10270',                            -6.22720000, 106.79740000, 150.00, 'ACTIVE'),
('stn-003', 'SPKLU Rest Area KM 57A Tol Jakarta-Cikampek',      'Tol Jakarta-Cikampek KM 57A, Klari, Karawang, Jawa Barat 41371',                               -6.36830000, 107.35120000, 200.00, 'ACTIVE'),
('stn-004', 'SPKLU PLN Rest Area KM 207A Tol Palimanan',        'Tol Palikanci KM 207A, Mundu, Cirebon, Jawa Barat 45173',                                      -6.77250000, 108.53610000, 100.00, 'ACTIVE'),
('stn-005', 'SPKLU PLN UID Jabar - Gedung Sate',                'Jl. Asia Afrika No.63, Braga, Sumur Bandung, Bandung, Jawa Barat 40111',                       -6.90250000, 107.61860000, 120.00, 'ACTIVE'),
('stn-006', 'SPKLU PLN UID Jawa Tengah Semarang',               'Jl. Pemuda No.93, Sekyu, Semarang Tengah, Semarang, Jawa Tengah 50132',                        -6.98220000, 110.42030000, 150.00, 'ACTIVE'),
('stn-007', 'SPKLU PLN UID Jawa Timur Surabaya',                'Jl. Yos Sudarso No.11, Embong Kaliasin, Genteng, Surabaya, Jawa Timur 60271',                  -7.26560000, 112.74830000, 180.00, 'ACTIVE'),
('stn-008', 'SPKLU PLN UID Sulselrabar Hertasning Makassar',    'Jl. Letjen Hertasning No.99, Kassi-Kassi, Rappocini, Kota Makassar, Sulawesi Selatan 90222',   -5.16780000, 119.44850000, 200.00, 'ACTIVE'),
('stn-009', 'SPKLU PLN Mattoanging Makassar',                   'Jl. Andi Mappanyukki No.14, Kunjung Mae, Mariso, Kota Makassar, Sulawesi Selatan 90125',       -5.15530000, 119.41420000,  50.00, 'ACTIVE'),
('stn-010', 'SPKLU PLN UP3 Parepare',                           'Jl. Ahmad Yani No.51, Ujung Baru, Soreang, Kota Parepare, Sulawesi Selatan 91112',             -4.01520000, 119.62890000,  50.00, 'ACTIVE'),
('stn-011', 'SPKLU PLN UP3 Palopo',                             'Jl. Kelapa No.1, Dangerakko, Wara, Kota Palopo, Sulawesi Selatan 91911',                       -2.99280000, 120.19830000,  50.00, 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

INSERT INTO charger_slots (id, station_id, slot_number, connector_type, max_power_kw, status) VALUES
-- stn-001: Gambir – 2x CCS2
('slot-001', 'stn-001', 1, 'CCS2',     100.00, 'AVAILABLE'),
('slot-002', 'stn-001', 2, 'CCS2',     100.00, 'AVAILABLE'),
-- stn-002: Senayan City – 1x CCS2
('slot-003', 'stn-002', 1, 'CCS2',     150.00, 'AVAILABLE'),
-- stn-003: KM 57A – 2x CCS2
('slot-004', 'stn-003', 1, 'CCS2',     100.00, 'AVAILABLE'),
('slot-005', 'stn-003', 2, 'CCS2',     100.00, 'AVAILABLE'),
-- stn-004: KM 207A – 1x CCS2
('slot-006', 'stn-004', 1, 'CCS2',     100.00, 'AVAILABLE'),
-- stn-005: Gedung Sate – 1x CCS2
('slot-007', 'stn-005', 1, 'CCS2',     120.00, 'AVAILABLE'),
-- stn-006: Semarang – 1x CCS2
('slot-008', 'stn-006', 1, 'CCS2',     150.00, 'AVAILABLE'),
-- stn-007: Surabaya – 1x CCS2
('slot-009', 'stn-007', 1, 'CCS2',     180.00, 'AVAILABLE'),
-- stn-008: Sulselrabar Hertasning Makassar – 2x CCS2
('slot-010', 'stn-008', 1, 'CCS2',     100.00, 'AVAILABLE'),
('slot-011', 'stn-008', 2, 'CCS2',     100.00, 'AVAILABLE'),
-- stn-009: Mattoanging Makassar – 1x CCS2
('slot-012', 'stn-009', 1, 'CCS2',      50.00, 'AVAILABLE'),
-- stn-010: Parepare – 1x CCS2
('slot-013', 'stn-010', 1, 'CCS2',      50.00, 'AVAILABLE'),
-- stn-011: Palopo – 1x CCS2
('slot-014', 'stn-011', 1, 'CCS2',      50.00, 'AVAILABLE')
ON CONFLICT (id) DO NOTHING;

-- Tarif PLN SPKLU berlaku: Rp 2.467/kWh
INSERT INTO tariffs (id, station_id, price_per_kwh, currency, valid_from, is_active) VALUES
('trf-001', 'stn-001', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-002', 'stn-002', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-003', 'stn-003', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-004', 'stn-004', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-005', 'stn-005', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-006', 'stn-006', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-007', 'stn-007', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-008', 'stn-008', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-009', 'stn-009', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-010', 'stn-010', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE),
('trf-011', 'stn-011', 2467.00, 'IDR', CURRENT_TIMESTAMP - INTERVAL '30 days', TRUE)
ON CONFLICT (id) DO NOTHING;
