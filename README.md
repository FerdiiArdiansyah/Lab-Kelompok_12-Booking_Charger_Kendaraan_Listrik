# Lab-Kelompok_12-Booking_Charger_Kendaraan_Listrik

![Arsitektur Sistem Booking Charger Kendaraan Listrik](arsitektur-v2.png)
Sistem ini dirancang untuk proses booking charger kendaraan listrik (SPKLU) dengan kapasitas slot terbatas, integrasi tarif resmi ESDM, dan arsitektur Event-Driven Microservices.

Tantangan utama:
- permintaan tinggi saat jam sibuk,
- konflik jadwal pada slot yang sama (anti-overlap),
- no-show yang membuat slot terblokir,
- tagihan harus berdasarkan konsumsi kWh aktual dan regulasi biaya layanan ESDM.

Tujuan sistem:
- mencegah booking overlap,
- menjaga fairness antrean dengan waitlist FIFO,
- melepas slot otomatis saat no-show,
- menghasilkan billing yang akurat dari sesi charging.

---

## 🛠️ Microservices Architecture

| Service Name | Port | Tanggung Jawab Utama |
| :--- | :---: | :--- |
| **`user-service`** | `:8084` | Registrasi & Autentikasi User (JWT), Profil, & Manajemen Kendaraan Listrik (EV) |
| **`station-service`** | `:8081` | Master Data Stasiun SPKLU, Slot Charger, & Tarif Listrik per kWh |
| **`booking-service`** | `:8082` | Reservasi Slot Rentang Waktu, Anti-Overlap, Antrean Waitlist FIFO, & Auto-Release No-Show |
| **`session-service`** | `:8083` | Inisiasi Sesi Charging, Real-time Meter kWh Telemetry, & Finalisasi Sesi |
| **`billing-service`** | `:8085` | Perhitungan Tagihan Invoice (kWh + Biaya Layanan ESDM + PPN), & Transaksi Pembayaran |

---

## 🚀 Daftar Lengkap Endpoint per Service

### 1. `user-service` (`:8084`)
* `POST /auth/register` — Registrasi pengguna baru
* `POST /auth/login` — Autentikasi & penerbitan token JWT
* `GET /users/me` — Ambil profil pengguna aktif
* `PUT /users/me` — Perbarui profil pengguna
* `GET /users` — Daftar seluruh pengguna (*Admin*)
* `GET /users/:id` — Detail pengguna berdasarkan ID
* `GET /users/me/vehicles` — Daftar kendaraan listrik milik pengguna
* `POST /users/me/vehicles` — Tambah kendaraan listrik baru
* `DELETE /users/me/vehicles/:vehicle_id` — Hapus kendaraan listrik

### 2. `station-service` (`:8081`)
* `GET /stations` — Daftar seluruh stasiun SPKLU
* `POST /stations` — Tambah stasiun pengisian baru
* `GET /stations/:id` — Detail stasiun beserta slot & tarif aktif
* `PUT /stations/:id` — Update informasi stasiun
* `DELETE /stations/:id` — Deaktivasi / hapus stasiun
* `GET /stations/:id/slots` — Daftar slot charger pada stasiun
* `POST /stations/:id/slots` — Tambah slot charger ke stasiun
* `PUT /stations/:id/slots/:slot_id` — Update status / spesifikasi slot
* `GET /stations/:id/tariff` — Ambil tarif aktif stasiun
* `POST /stations/:id/tariffs` — Tambah / perbarui tarif stasiun
* `GET /tariffs` — Daftar seluruh tarif yang terdaftar

### 3. `booking-service` (`:8082`)
* `POST /bookings` — Buat booking slot charger (dengan validasi anti-overlap & fallback waitlist)
* `GET /bookings/:id` — Detail booking berdasarkan ID
* `GET /bookings/user/:user_id` — Riwayat booking milik pengguna
* `GET /stations/:id/availability` — Cek ketersediaan slot stasiun pada rentang waktu (`?start=&end=`)
* `GET /stations/:id/waitlist` — Daftar antrean waitlist stasiun
* `POST /bookings/:id/check-in` — Check-in lokasi & aktifkan status sesi
* `POST /bookings/:id/cancel` — Pembatalan booking oleh pengguna
* `POST /bookings/auto-release` — Trigger otomatis pelepas booking no-show (`?grace_period_minutes=15`)

### 4. `session-service` (`:8083`)
* `POST /sessions/start` — Mulai sesi pengisian daya
* `GET /sessions/:id` — Detail sesi pengisian & telemetry meteran
* `GET /sessions/booking/:booking_id` — Ambil sesi pengisian berdasarkan ID booking
* `GET /sessions/user/:user_id` — Riwayat sesi pengisian pengguna
* `POST /sessions/:id/meter` — Kirim pembacaan meter kWh real-time
* `POST /sessions/:id/finish` — Akhiri sesi pengisian & finalisasi total kWh

### 5. `billing-service` (`:8085`)
* `POST /invoices` — Generasi invoice pengisian dari kWh & tarif
* `POST /invoices/generate` — Endpoint alternatif pembuat invoice
* `GET /invoices/:id` — Detail rincian invoice
* `GET /invoices/session/:session_id` — Ambil invoice berdasarkan ID sesi
* `GET /invoices/user/:user_id` — Riwayat invoice milik pengguna
* `POST /payments` — Inisiasi transaksi pembayaran invoice
* `GET /payments/:id` — Detail status transaksi pembayaran
* `POST /payments/:id/confirm` — Konfirmasi/Webhook penyelesaian pembayaran

---

## 📊 Dokumen Data & Excel SPKLU Indonesia
- Dokumentasi Data Aktual SPKLU: [`docs/data_aktual_spklu.md`](docs/data_aktual_spklu.md)
- Dataset CSV Populasi EV: [`docs/populasi_ev_indonesia.csv`](docs/populasi_ev_indonesia.csv)
- File Spreadsheet Excel Resmi: [`docs/Data_Aktual_SPKLU_Indonesia.xlsx`](docs/Data_Aktual_SPKLU_Indonesia.xlsx)