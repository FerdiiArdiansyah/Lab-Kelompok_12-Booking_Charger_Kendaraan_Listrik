# ADR 0001: Arsitektur Microservices Berbasis Domain

- Status: Accepted
- Date: 2026-08-15

## Context
Sistem memiliki domain yang jelas terpisah:
- manajemen stasiun/slot/tarif,
- booking dan konflik jadwal,
- sesi charging real-time,
- billing dan pembayaran.

Beban saat jam sibuk tinggi, dan kebutuhan scaling tiap domain berbeda.

## Decision
Mengadopsi arsitektur microservices dengan service:
- station-service,
- booking-service,
- session-service,
- billing-service.

Setiap service memiliki database sendiri (database per service) dan berkomunikasi:
- sinkron untuk operasi query cepat,
- asinkron via event bus untuk integrasi lintas domain.

## Consequences
Positif:
- isolasi domain dan ownership jelas,
- scaling per service lebih fleksibel,
- kegagalan lebih terlokalisasi.

Negatif:
- kompleksitas operasional meningkat,
- butuh observability dan tracing yang matang,
- konsistensi lintas service menjadi eventual consistency.

Mitigasi:
- outbox pattern,
- idempotent consumer,
- distributed tracing dengan correlationId.
