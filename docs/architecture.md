# Arsitektur Sistem

## 1. Tujuan
Sistem Booking Charger Kendaraan Listrik dirancang untuk:
- mencegah bentrok slot pada rentang waktu yang bertabrakan,
- menangani lonjakan permintaan pada jam sibuk,
- melepaskan slot otomatis saat no-show,
- memastikan tagihan berdasarkan kWh aktual,
- tetap konsisten tanpa mengorbankan skalabilitas.

## 2. Bounded Context dan Service

### station-service
Tanggung jawab:
- data stasiun (lokasi, kapasitas daya),
- data slot charger,
- tarif dasar per kWh,
- status slot operasional (aktif/nonaktif).

Data utama:
- Station(id, name, location, totalPowerKw)
- ChargerSlot(id, stationId, connectorType, maxPowerKw, status)
- Tariff(id, stationId, pricePerKwh, validFrom)

### booking-service
Tanggung jawab:
- membuat booking rentang waktu,
- mencegah overlap booking slot,
- antrean saat kapasitas penuh,
- auto-release jika no-show.

Data utama:
- Booking(id, userId, stationId, slotId, startTime, endTime, status)
- Waitlist(id, stationId, requestedStart, requestedEnd, queueNumber, status)

Status booking:
- REQUESTED
- CONFIRMED
- EXPIRED_NO_SHOW
- CANCELLED
- IN_SESSION
- COMPLETED

### session-service
Tanggung jawab:
- mulai sesi charging,
- tracking pemakaian kWh,
- selesai sesi,
- publish pemakaian akhir.

Data utama:
- ChargingSession(id, bookingId, startedAt, endedAt, consumedKwh, status)

### billing-service
Tanggung jawab:
- hitung tagihan dari consumedKwh,
- terapkan tarif yang berlaku,
- proses pembayaran,
- simpan invoice.

Data utama:
- Invoice(id, sessionId, tariffId, consumedKwh, subtotal, tax, total, status)
- Payment(id, invoiceId, method, amount, status, paidAt)

## 3. Gaya Arsitektur
- Microservices per domain utama.
- Database per service (no shared database).
- Komunikasi sinkron untuk query cepat (HTTP/gRPC).
- Komunikasi asinkron berbasis event untuk proses lintas domain.

Event kunci:
- BookingConfirmed
- BookingExpiredNoShow
- SessionStarted
- SessionFinished
- InvoiceCreated
- PaymentCompleted

## 4. Aturan Konsistensi Kritis

### 4.1 Anti-overlap booking
Aturan:
- satu slot tidak boleh punya dua booking berstatus aktif pada waktu bertabrakan.

Implementasi di booking-service:
- transaksi database dengan isolation minimal REPEATABLE READ atau SERIALIZABLE untuk operasi booking,
- exclusion constraint (jika PostgreSQL) pada interval waktu per slot,
- idempotency key untuk create booking API,
- unique booking reference.

Contoh konsep constraint PostgreSQL:
- EXCLUDE USING gist (slot_id WITH =, tsrange(start_time, end_time, '[)') WITH &&)

### 4.2 No-show auto-release
Aturan:
- booking CONFIRMED harus check-in sebelum grace period berakhir.

Implementasi:
- scheduler di booking-service memeriksa booking yang melewati startTime + gracePeriod,
- ubah status menjadi EXPIRED_NO_SHOW,
- publish BookingExpiredNoShow,
- proses waitlist untuk promosi antrian berikutnya.

### 4.3 Konsistensi antar service
- Gunakan eventual consistency antar service via event bus.
- Gunakan outbox pattern per service untuk reliabilitas publish event.
- Consumer wajib idempotent memakai eventId/version.

## 5. Alur Utama

### Alur booking
1. User request booking (station, waktu, preferensi slot).
2. booking-service cek slot tersedia dan cek overlap.
3. Jika ada slot valid, booking CONFIRMED.
4. Jika tidak ada, user masuk waitlist.
5. Event BookingConfirmed dipublish.

### Alur check-in dan sesi
1. User datang dan check-in sebelum grace period habis.
2. session-service menerima perintah start session berdasarkan booking valid.
3. Status booking berubah menjadi IN_SESSION.
4. SessionFinished dipublish saat charging selesai.

### Alur billing
1. billing-service consume SessionFinished.
2. Ambil tarif dari station-service (atau cache tarif terverifikasi).
3. Buat invoice dari consumedKwh x pricePerKwh.
4. Pembayaran diproses, status invoice diperbarui.

## 6. Ketahanan Saat Jam Sibuk
- Rate limiting per user untuk endpoint booking.
- Prioritas antrean FIFO per stasiun dan window waktu.
- Optimistic retry terbatas pada konflik booking.
- Read model/caching untuk ketersediaan slot agar query cepat.
- Partitioning tabel booking berdasarkan station_id atau waktu jika volume tinggi.

## 7. API Ringkas yang Disarankan

station-service:
- GET /stations
- GET /stations/{id}/slots
- GET /stations/{id}/tariff

booking-service:
- POST /bookings
- POST /bookings/{id}/check-in
- POST /bookings/{id}/cancel
- GET /stations/{id}/availability?start=&end=

session-service:
- POST /sessions/start
- POST /sessions/{id}/meter
- POST /sessions/{id}/finish

billing-service:
- GET /invoices/{id}
- POST /payments

## 8. Keamanan dan Audit
- JWT/OAuth2 antar client dan API gateway.
- mTLS antar service internal.
- Audit log immutable untuk perubahan status booking, session, invoice, dan payment.
- Korelasi tracing memakai correlationId.

## 9. Teknologi yang Cocok
- Runtime: Java/Spring Boot, Node.js, atau Go (pilih seragam tim).
- Database transaksional: PostgreSQL (direkomendasikan untuk range constraint).
- Message broker: Kafka atau RabbitMQ.
- Cache: Redis untuk availability read model.
- Observability: OpenTelemetry + Prometheus + Grafana.
