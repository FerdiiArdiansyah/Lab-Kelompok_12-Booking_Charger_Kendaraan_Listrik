-- user-db Seed Data
-- password_hash adalah bcrypt dari "Password@123" — hanya untuk development

INSERT INTO users (id, name, email, phone, password_hash, role, status) VALUES
('usr-001', 'Budi Santoso',    'budi.santoso@gmail.com',    '081234567001', '$2a$10$kXBhIEuIVcJQqBbHUq6oPuLpxNJCv4j8PalGwUh4E5e8pnVMEWNKe', 'USER',  'ACTIVE'),
('usr-002', 'Siti Rahayu',     'siti.rahayu@gmail.com',     '081234567002', '$2a$10$kXBhIEuIVcJQqBbHUq6oPuLpxNJCv4j8PalGwUh4E5e8pnVMEWNKe', 'USER',  'ACTIVE'),
('usr-003', 'Andi Wijaya',     'andi.wijaya@gmail.com',     '081234567003', '$2a$10$kXBhIEuIVcJQqBbHUq6oPuLpxNJCv4j8PalGwUh4E5e8pnVMEWNKe', 'USER',  'ACTIVE'),
('usr-004', 'Dewi Lestari',    'dewi.lestari@yahoo.com',    '081234567004', '$2a$10$kXBhIEuIVcJQqBbHUq6oPuLpxNJCv4j8PalGwUh4E5e8pnVMEWNKe', 'USER',  'ACTIVE'),
('usr-005', 'Rizky Pratama',   'rizky.pratama@gmail.com',   '081234567005', '$2a$10$kXBhIEuIVcJQqBbHUq6oPuLpxNJCv4j8PalGwUh4E5e8pnVMEWNKe', 'USER',  'ACTIVE'),
('usr-006', 'Nurul Hidayah',   'nurul.hidayah@gmail.com',   '081234567006', '$2a$10$kXBhIEuIVcJQqBbHUq6oPuLpxNJCv4j8PalGwUh4E5e8pnVMEWNKe', 'USER',  'ACTIVE'),
('usr-007', 'Fajar Setiawan',  'fajar.setiawan@gmail.com',  '081234567007', '$2a$10$kXBhIEuIVcJQqBbHUq6oPuLpxNJCv4j8PalGwUh4E5e8pnVMEWNKe', 'USER',  'ACTIVE'),
('usr-008', 'Maya Putri',      'maya.putri@gmail.com',      '081234567008', '$2a$10$kXBhIEuIVcJQqBbHUq6oPuLpxNJCv4j8PalGwUh4E5e8pnVMEWNKe', 'USER',  'ACTIVE'),
('usr-009', 'Hendra Gunawan',  'hendra.gunawan@gmail.com',  '081234567009', '$2a$10$kXBhIEuIVcJQqBbHUq6oPuLpxNJCv4j8PalGwUh4E5e8pnVMEWNKe', 'USER',  'ACTIVE'),
('usr-010', 'Linda Susanti',   'linda.susanti@gmail.com',   '081234567010', '$2a$10$kXBhIEuIVcJQqBbHUq6oPuLpxNJCv4j8PalGwUh4E5e8pnVMEWNKe', 'USER',  'ACTIVE'),
('usr-admin', 'Admin EV Charger', 'admin@evcharger.id',     '081200000000', '$2a$10$kXBhIEuIVcJQqBbHUq6oPuLpxNJCv4j8PalGwUh4E5e8pnVMEWNKe', 'ADMIN', 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

INSERT INTO user_vehicles (id, user_id, brand, model, license_plate, connector_type, battery_capacity_kwh, is_default) VALUES
('veh-001', 'usr-001', 'Hyundai', 'Ioniq 5',       'B 1234 EV', 'CCS2',    72.60, TRUE),
('veh-002', 'usr-002', 'Wuling',  'Air ev',        'B 5678 EV', 'TYPE2',   17.30, TRUE),
('veh-003', 'usr-003', 'Tesla',   'Model 3',       'B 9012 EV', 'CCS2',    75.00, TRUE),
('veh-004', 'usr-004', 'BYD',     'Atto 3',        'B 3456 EV', 'CCS2',    60.48, TRUE),
('veh-005', 'usr-005', 'Toyota',  'bZ4X',          'B 7890 EV', 'TYPE2',   71.40, TRUE),
('veh-006', 'usr-006', 'Nissan',  'Leaf',          'B 2345 EV', 'CHADEMO', 40.00, TRUE),
('veh-007', 'usr-007', 'BMW',     'iX3',           'B 6789 EV', 'CCS2',    74.00, TRUE),
('veh-008', 'usr-008', 'MG',      'ZS EV',         'B 0123 EV', 'CCS2',    50.30, TRUE),
('veh-009', 'usr-009', 'Hyundai', 'Kona Electric', 'B 4567 EV', 'CCS2',    64.00, TRUE),
('veh-010', 'usr-010', 'Wuling',  'Almaz RS EV',   'B 8901 EV', 'TYPE2',   84.00, TRUE)
ON CONFLICT (id) DO NOTHING;
