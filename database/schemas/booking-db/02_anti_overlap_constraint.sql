-- 02_anti_overlap_constraint.sql
-- Penjelasan & Demonstrasi Pengujian Exclusion Constraint Anti-Overlap pada booking_db

/*
  DOKUMEN TEKNIS EXCLUSION CONSTRAINT:
  
  Setiap slot charger hanya boleh di-booking oleh 1 user pada rentang waktu `[start_time, end_time)`.
  Jika ada 2 transaksi bersamaan (race condition) yang mencoba membuat booking pada slot yang sama
  dengan interval waktu yang saling beririsan (overlap), database PostgreSQL akan langsung MENOLAK
  transaksi kedua dengan error:
  `ERROR: conflicting key value violates exclusion constraint "no_overlapping_bookings"`

  Fitur ini memanfaatkan extension PostgreSQL `btree_gist` dan tipe data `tstzrange` (Timestamp with Time Zone Range).
*/

-- CONTOH SKENARIO PENGUJIANKU INTEGRITAS DATA (SQL TEST PROOF):

-- 1. Insert Booking 1 (Berhasil)
INSERT INTO bookings (id, user_id, station_id, slot_id, start_time, end_time, status, idempotency_key)
VALUES (
    'bkg-101', 
    'usr-001', 
    'stn-001', 
    'slot-001', 
    '2026-08-20 10:00:00+07', 
    '2026-08-20 12:00:00+07', 
    'CONFIRMED', 
    'idem-001'
) ON CONFLICT (id) DO NOTHING;

-- 2. Insert Booking 2 pada Slot Berbeda (Berhasil)
INSERT INTO bookings (id, user_id, station_id, slot_id, start_time, end_time, status, idempotency_key)
VALUES (
    'bkg-102', 
    'usr-002', 
    'stn-001', 
    'slot-002', 
    '2026-08-20 10:30:00+07', 
    '2026-08-20 11:30:00+07', 
    'CONFIRMED', 
    'idem-002'
) ON CONFLICT (id) DO NOTHING;

-- 3. Attempt Insert Booking 3 yang OVERLAP dengan Booking 1 di slot-001 (Akan Ditolak DB Engine)
-- Jalankan ini di SQL client untuk memverifikasi proteksi DB:
/*
INSERT INTO bookings (id, user_id, station_id, slot_id, start_time, end_time, status, idempotency_key)
VALUES (
    'bkg-103', 
    'usr-003', 
    'stn-001', 
    'slot-001', 
    '2026-08-20 11:00:00+07', -- Overlap jam 11:00 - 12:00
    '2026-08-20 13:00:00+07', 
    'CONFIRMED', 
    'idem-003'
);
-- EXPECTED RESULT: ERROR (exclusion constraint violation)
*/
