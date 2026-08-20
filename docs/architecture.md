# Arsitektur Sistem Booking Charger Kendaraan Listrik

## 1. Tujuan Sistem
Sistem Booking Charger Kendaraan Listrik dirancang untuk:
- Mencegah bentrok slot pada rentang waktu yang bertabrakan (anti-overlap),
- Menangani lonjakan permintaan pada jam sibuk dengan antrean FIFO (waitlist),
- Melepaskan slot otomatis saat pengguna mengalami no-show,
- Memastikan tagihan berdasarkan konsumsi kWh aktual dan regulasi tarif ESDM (Fast & Ultra Fast charging service fee),
- Menjamin konsistensi data antar domain berbasis Event-Driven Microservices.

-----

## 2. Bounded Context dan Service

### 2.1 user-service (Port: 8084)
Tanggung jawab:
- Registrasi dan autentikasi user (JWT Auth),
- Manajemen profil pengguna dan hak akses (Role: USER, ADMIN, OPERATOR),
- Manajemen data kendaraan listrik (EV) terdaftar milik pengguna.

Data utama:
- `User(id, name, email, password_hash, phone, role, status, created_at, updated_at)`
- `UserVehicle(id, user_id, brand, model, license_plate, connector_type, battery_capacity_kwh, created_at)`

---

### 2.2 station-service (Port: 8081)
Tanggung jawab:
- Master data stasiun SPKLU (lokasi, koordinat GPS, kapasitas total daya kW),
- Manajemen slot charger (nomor slot, tipe konektor, daya kW maksimal, status operasional),
- Manajemen tarif per kWh dan biaya layanan.

Data utama:
- `Station(id, name, location, latitude, longitude, total_power_kw, status)`
- `ChargerSlot(id, station_id, slot_number, connector_type, max_power_kw, status)`
- `Tariff(id, station_id, price_per_kwh, currency, valid_from, is_active)`

---

### 2.3 booking-service (Port: 8082)
Tanggung jawab:
- Pembuatan dan reservasi slot rentang waktu (anti-overlap),
- Manajemen antrean waitlist (FIFO) saat kapasitas penuh,
- Penanganan check-in pengguna dan pembatalan booking,
- Auto-release slot booking saat lewat grace period (no-show).

Data utama:
- `Booking(id, user_id, station_id, slot_id, start_time, end_time, status, idempotency_key)`
- `Waitlist(id, station_id, user_id, requested_start, requested_end, queue_number, status)`

Status Booking:
- `REQUESTED`, `CONFIRMED`, `WAITLISTED`, `IN_SESSION`, `COMPLETED`, `EXPIRED_NO_SHOW`, `CANCELLED`

---

### 2.4 session-service (Port: 8083)
Tanggung jawab:
- Memulai dan menghentikan sesi pengisian daya (*charging session*),
- Pelacakan telemetry / meteran listrik (*real-time kWh meter reading*),
- Publish event konsumsi kWh final saat sesi pengisian selesai.

Data utama:
- `ChargingSession(id, booking_id, slot_id, user_id, started_at, ended_at, consumed_kwh, status)`
- `MeterReading(id, session_id, recorded_at, current_kwh, power_kw, voltage, current_ampere)`

---

### 2.5 billing-service (Port: 8085)
Tanggung jawab:
- Perhitungan invoice berdasarkan pemakaian kWh aktual & komponen biaya layanan ESDM,
- Penerapan PPN (11%) & Pajak Penerangan Jalan (PPJ),
- Eksekusi dan konfirmasi transaksi pembayaran (QRIS, VA, E-Wallet, Credit Card),
- Penyimpanan audit log perubahan status invoice.

Data utama:
- `Invoice(id, session_id, user_id, tariff_id, consumed_kwh, price_per_kwh, service_fee, subtotal, tax, total, status)`
- `Payment(id, invoice_id, payment_method, amount, status, transaction_ref, paid_at)`

---

## 3. Daftar Lengkap API Endpoint per Service

### 🔑 1. user-service (`:8084`)
* **`POST /auth/register`** — Registrasi pengguna baru
* **`POST /auth/login`** — Autentikasi & penerbitan token JWT
* **`GET /users/me`** — Ambil profil pengguna aktif
* **`PUT /users/me`** — Perbarui profil pengguna
* **`GET /users`** — Daftar seluruh pengguna (*Admin*)
* **`GET /users/:id`** — Detail pengguna berdasarkan ID
* **`GET /users/me/vehicles`** — Daftar kendaraan listrik milik pengguna
* **`POST /users/me/vehicles`** — Tambah kendaraan listrik baru
* **`DELETE /users/me/vehicles/:vehicle_id`** — Hapus kendaraan listrik

### ⚡ 2. station-service (`:8081`)
* **`GET /stations`** — Daftar seluruh stasiun SPKLU
* **`POST /stations`** — Tambah stasiun pengisian baru
* **`GET /stations/:id`** — Detail stasiun beserta slot & tarif aktif
* **`PUT /stations/:id`** — Update informasi stasiun
* **`DELETE /stations/:id`** — Deaktivasi / hapus stasiun
* **`GET /stations/:id/slots`** — Daftar slot charger pada stasiun
* **`POST /stations/:id/slots`** — Tambah slot charger ke stasiun
* **`PUT /stations/:id/slots/:slot_id`** — Update status / spesifikasi slot
* **`GET /stations/:id/tariff`** — Ambil tarif aktif stasiun
* **`POST /stations/:id/tariffs`** — Tambah / perbarui tarif stasiun
* **`GET /tariffs`** — Daftar seluruh tarif yang terdaftar

### 📅 3. booking-service (`:8082`)
* **`POST /bookings`** — Buat booking slot charger (dengan validasi anti-overlap & fallback waitlist)
* **`GET /bookings/:id`** — Detail booking berdasarkan ID
* **`GET /bookings/user/:user_id`** — Riwayat booking milik pengguna
* **`GET /stations/:id/availability`** — Cek ketersediaan slot stasiun pada rentang waktu (`?start=&end=`)
* **`GET /stations/:id/waitlist`** — Daftar antrean waitlist stasiun
* **`POST /bookings/:id/check-in`** — Check-in lokasi & aktifkan status sesi
* **`POST /bookings/:id/cancel`** — Pembatalan booking oleh pengguna
* **`POST /bookings/auto-release`** — Trigger otomatis pelepas booking no-show (`?grace_period_minutes=15`)

### 🔌 4. session-service (`:8083`)
* **`POST /sessions/start`** — Mulai sesi pengisian daya
* **`GET /sessions/:id`** — Detail sesi pengisian & telemetry meteran
* **`GET /sessions/booking/:booking_id`** — Ambil sesi pengisian berdasarkan ID booking
* **`GET /sessions/user/:user_id`** — Riwayat sesi pengisian pengguna
* **`POST /sessions/:id/meter`** — Kirim pembacaan meter kWh real-time
* **`POST /sessions/:id/finish`** — Akhiri sesi pengisian & finalisasi total kWh

### 💳 5. billing-service (`:8085`)
* **`POST /invoices`** — Generasi invoice pengisian dari kWh & tarif
* **`POST /invoices/generate`** — Endpoint alternatif pembuat invoice
* **`GET /invoices/:id`** — Detail rincian invoice
* **`GET /invoices/session/:session_id`** — Ambil invoice berdasarkan ID sesi
* **`GET /invoices/user/:user_id`** — Riwayat invoice milik pengguna
* **`POST /payments`** — Inisiasi transaksi pembayaran invoice
* **`GET /payments/:id`** — Detail status transaksi pembayaran
* **`POST /payments/:id/confirm`** — Konfirmasi/Webhook penyelesaian pembayaran

---

## 4. Pola Integrasi dan Konsistensi
1. **Outbox Pattern**: Setiap mutasi status utama (Booking, Session, Payment) menyimpan event ke tabel `outbox_events` dalam transaksi basis data yang sama.
2. **Event Bus Asinkron**: Event dipublish ke broker (RabbitMQ/Kafka) untuk konsistensi antar-service (eventual consistency).
3. **Idempotensi**: Seluruh API mutasi (Create Booking, Start Session, Payment Process) mendukung `idempotency_key`.
