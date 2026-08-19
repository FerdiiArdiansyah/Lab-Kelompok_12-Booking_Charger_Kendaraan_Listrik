# 🔒 Strategi Anti-Overlap Booking Slot Charger

## 1. Problem Statement

Pada sistem booking charger EV saat jam sibuk (peak hours), masalah race condition sangat rawan terjadi. Dua atau lebih user dapat mengirimkan permintaan booking untuk **slot charger yang sama** pada **rentang waktu yang saling beririsan (overlap)** secara bersamaan.

Jika validasi hanya dilakukan di layer aplikasi (seperti `SELECT COUNT(*) WHERE slot_id = X AND time OVERLAPS...`), race condition dalam transaksi concurrent tetap dapat meloloskan booking ganda (*double booking*).

---

## 2. Solusi Data & Persistence Engineer: PostgreSQL Range Exclusion Constraint

Data & Persistence Engineer memindahkan enforcement integritas anti-bentrok dari layer aplikasi ke **level terendah Engine Database PostgreSQL** menggunakan:
1. Extension **`btree_gist`**.
2. Tipe Data **`tstzrange`** (Timestamp with Time Zone Range).
3. **`EXCLUSION CONSTRAINT`** berbasis GIST Index.

### DDL Implementation Snippet

```sql
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE bookings (
    id VARCHAR(36) PRIMARY KEY,
    slot_id VARCHAR(36) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    booking_period TSTZRANGE GENERATED ALWAYS AS (tstzrange(start_time, end_time, '[)')) STORED,
    status VARCHAR(50) NOT NULL DEFAULT 'REQUESTED',
    CONSTRAINT chk_booking_time CHECK (end_time > start_time)
);

ALTER TABLE bookings ADD CONSTRAINT no_overlapping_bookings
EXCLUSION USING GIST (
    slot_id WITH =,
    booking_period WITH &&
) WHERE (status IN ('REQUESTED', 'CONFIRMED', 'IN_SESSION'));
```

---

## 3. Cara Kerja Komponen

1. **`tstzrange(start_time, end_time, '[)')`**:
   - Boundary `[` berarti *inclusive* (termasuk start time).
   - Boundary `)` berarti *exclusive* (tidak termasuk end time).
   - Contoh: User A booking 10:00 - 11:00. User B booking 11:00 - 12:00 -> **TIDAK BENTROK** karena jam 11:00 tepat sudah dilepas untuk User B.

2. **Operator `&&` (Overlap Operator)**:
   - Memeriksa apakah dua range interval waktu memiliki irisian set sekecil apapun.

3. **`WHERE (status IN ('REQUESTED', 'CONFIRMED', 'IN_SESSION'))`**:
   - Booking yang berstatus `CANCELLED` atau `EXPIRED_NO_SHOW` diabaikan oleh index constraint ini, sehingga slot otomatis bebas tanpa perlu hapus data fisik (hard delete).

---

## 4. Keunggulan Dibandingkan Locking Aplikasi

| Parameter | Application Lock (Redis/Distributed) | PostgreSQL Exclusion Constraint |
| :--- | :--- | :--- |
| **Penyebab Fail** | Redis downtime / Lock Timeout / Crash | Terjamin 100% oleh ACID Engine PostgreSQL |
| **Performa** | Membutuhkan network roundtrip tambahan ke Redis | O(log N) GIST Index Lookup |
| **Race Condition** | Masih ada celah jika TTL expired saat query lambat | Impossible (Di-enforce pada level row serialization DB) |
| **Kompleksitas** | Membutuhkan Redlock / Distributed Lock Manager | Cukup 1 Constraint DDL SQL |
