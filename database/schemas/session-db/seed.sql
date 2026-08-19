-- session-db Seed Data
-- kWh dan profil daya disesuaikan dengan kapasitas baterai & kemampuan charger tiap EV

INSERT INTO charging_sessions (id, booking_id, slot_id, user_id, started_at, ended_at, consumed_kwh, status) VALUES
-- Ioniq 5 (72.6 kWh): CCS2 100kW, ~15%→80% SoC = 47.2 kWh, efektif 38 menit
('ses-001', 'bkg-001', 'slot-001', 'usr-001', CURRENT_TIMESTAMP - INTERVAL '3 hours 28 minutes', CURRENT_TIMESTAMP - INTERVAL '2 hours 50 minutes', 50.600, 'COMPLETED'),
-- Wuling Air ev (17.3 kWh): TYPE2 22kW, onboard charger max 6.6 kW, ~20%→93% SoC, sedang berlangsung
('ses-002', 'bkg-002', 'slot-009', 'usr-002', CURRENT_TIMESTAMP - INTERVAL '55 minutes',         NULL,                                              6.050,  'IN_PROGRESS'),
-- Nissan Leaf (40 kWh): CHAdeMO 50kW, ~10%→76% SoC = 26.4 kWh, daya taper di 60% SoC, 48 menit
('ses-003', 'bkg-006', 'slot-004', 'usr-006', CURRENT_TIMESTAMP - INTERVAL '5 hours 58 minutes', CURRENT_TIMESTAMP - INTERVAL '5 hours 10 minutes', 26.400, 'COMPLETED')
ON CONFLICT (id) DO NOTHING;

INSERT INTO meter_readings (session_id, recorded_at, current_kwh, power_kw, voltage, current_ampere) VALUES
-- ses-001: Ioniq 5, CCS2 100kW – daya penuh di awal, taper menjelang 80% SoC
('ses-001', CURRENT_TIMESTAMP - INTERVAL '3 hours 28 minutes',  0.000, 100.00, 520.00, 192.30),
('ses-001', CURRENT_TIMESTAMP - INTERVAL '3 hours 8 minutes',  25.300,  97.50, 523.00, 186.40),
('ses-001', CURRENT_TIMESTAMP - INTERVAL '2 hours 50 minutes', 50.600,  32.00, 530.00,  60.40),
-- ses-002: Wuling Air ev, TYPE2 – onboard charger max 6.6 kW (konstan, tidak taper)
('ses-002', CURRENT_TIMESTAMP - INTERVAL '55 minutes',  0.000, 6.60, 229.00, 28.80),
('ses-002', CURRENT_TIMESTAMP - INTERVAL '30 minutes',  2.750, 6.60, 229.00, 28.80),
('ses-002', CURRENT_TIMESTAMP,                          6.050, 6.60, 229.00, 28.80),
-- ses-003: Nissan Leaf, CHAdeMO 50kW – taper mulai di 60% SoC karena limit BMS Leaf
('ses-003', CURRENT_TIMESTAMP - INTERVAL '5 hours 58 minutes',  0.000, 50.00, 380.00, 131.60),
('ses-003', CURRENT_TIMESTAMP - INTERVAL '5 hours 32 minutes', 13.300, 44.00, 388.00, 113.40),
('ses-003', CURRENT_TIMESTAMP - INTERVAL '5 hours 10 minutes', 26.400, 18.00, 395.00,  45.60);
