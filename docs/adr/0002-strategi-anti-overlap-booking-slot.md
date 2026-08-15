# ADR 0002: Strategi Anti-overlap Booking Slot

- Status: Accepted
- Date: 2026-08-15

## Context
Sumber konflik utama adalah slot yang sama pada rentang waktu bertabrakan.
Aturan bisnis kritis: satu slot tidak boleh memiliki dua booking aktif dengan waktu overlap.

## Decision
Penegakan anti-overlap dilakukan di booking-service pada level database dan transaksi:
- transaksi create booking menggunakan isolation yang kuat,
- constraint non-overlap interval per slot,
- retry terbatas pada konflik transaksi,
- idempotency key untuk mencegah duplikasi request.

Contoh implementasi PostgreSQL:
- exclusion constraint berbasis tsrange + gist index.

## Alternatives Considered
1. Cek overlap hanya di level aplikasi.
   Risiko race condition tinggi saat trafik padat.
2. Distributed lock eksternal untuk seluruh booking.
   Menambah latensi dan single hot lock saat jam sibuk.

## Consequences
Positif:
- integritas data booking terjaga kuat,
- race condition jauh berkurang,
- perilaku deterministik saat konflik.

Negatif:
- desain schema lebih kompleks,
- throughput write dapat turun pada contention ekstrem.

Mitigasi:
- partitioning data booking,
- retry policy adaptif,
- antrean waitlist untuk request gagal alokasi.
