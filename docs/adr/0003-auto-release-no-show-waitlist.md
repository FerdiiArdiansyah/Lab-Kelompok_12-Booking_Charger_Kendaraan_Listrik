# ADR 0003: Auto-release No-show dan Waitlist Promotion

- Status: Accepted
- Date: 2026-08-15

## Context
Pada jam sibuk, banyak pengguna memesan tapi tidak datang. Slot yang terkunci oleh no-show menurunkan utilisasi dan meningkatkan waktu tunggu.

## Decision
Booking CONFIRMED wajib check-in sebelum grace period berakhir.
Jika tidak check-in:
- status booking menjadi EXPIRED_NO_SHOW,
- event BookingExpiredNoShow dipublish,
- proses waitlist mempromosikan antrean berikutnya.

Mekanisme dilakukan oleh scheduler periodik di booking-service.

## Alternatives Considered
1. Slot tetap ditahan sampai endTime.
   Tidak efisien, utilisasi rendah.
2. Pelepasan slot manual oleh operator.
   Tidak scalable dan rawan keterlambatan.

## Consequences
Positif:
- utilisasi slot meningkat,
- fairness antrean lebih baik,
- pendapatan stasiun lebih optimal.

Negatif:
- perlu sinkronisasi notifikasi real-time ke user,
- ada risiko user terlambat check-in karena faktor eksternal.

Mitigasi:
- grace period yang masuk akal (mis. 10-15 menit),
- notifikasi pengingat sebelum waktu mulai,
- SLA untuk proses promotion waitlist.
