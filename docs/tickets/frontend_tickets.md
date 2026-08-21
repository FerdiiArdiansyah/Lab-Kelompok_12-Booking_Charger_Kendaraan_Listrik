# Daftar Tiket Implementasi Frontend (Web Dashboard & Consumer App)
**Sistem Booking & Management SPKLU Kendaraan Listrik**

Dokumen ini berisi breakdown tiket pengembangan Frontend berbasis **React / Next.js / Vanilla Modern UI** yang terintegrasi langsung dengan 5 Backend Microservices (`user`, `station`, `booking`, `session`, `billing`).

---

## 📊 Summary Matriks Tiket Frontend

| Ticket ID | Target Module | Title | Priority | Target Backend Service | Est. Complexity |
|---|---|---|---|---|---|
| **`FRONTEND-01`** | Design System | Futuristic Dark Mode UI Theme & Component Library | `P0 (Critical)` | - | Low |
| **`FRONTEND-02`** | Station Explorer | Interactive SPKLU Map & Real-Time Charger Search | `P0 (Critical)` | `station-service` | Medium |
| **`FRONTEND-03`** | Booking Flow | Slot Time Picker & Anti-Overlap Booking Engine UI | `P0 (Critical)` | `booking-service` | High |
| **`FRONTEND-04`** | Waitlist UI | FIFO Waitlist Banner & Live Queue Position Indicator | `P1 (High)` | `booking-service` | Medium |
| **`FRONTEND-05`** | Live Telemetry | Active Charging Telemetry Gauge & Real-Time kWh Graph | `P1 (High)` | `session-service` | High |
| **`FRONTEND-06`** | Billing & Checkout | ESDM Tariff Breakdown & QRIS Payment Checkout Modal | `P0 (Critical)` | `billing-service` | Medium |
| **`FRONTEND-07`** | EV Profile | User Vehicle Garage & Connector Type Compatibility Manager | `P2 (Medium)` | `user-service` | Low |
| **`FRONTEND-08`** | Admin Portal | Operator Station Control & No-Show Monitoring Dashboard | `P2 (Medium)` | All Microservices | Medium |

---

## 🎯 Detail Spesifikasi Tiket Frontend

### **`FRONTEND-01`**: Futuristic Dark Mode UI Theme & Design System
- **Objective**: Membangun fondasi CSS Design System bernuansa EV futuristik (*Dark Mode, Glassmorphism, Neon Cyan & Emerald accent, Inter typography*).
- **Key Features**:
  - CSS Variables untuk tema HSL (`--bg-primary`, `--accent-cyan`, `--accent-emerald`, `--card-glass`).
  - Reusable components: Button, Input, Modal, Badge Status (`AVAILABLE`, `IN_USE`, `OUT_OF_SERVICE`).
  - Micro-animations untuk hover state dan status pulsing.

---

### **`FRONTEND-02`**: Interactive SPKLU Map & Real-Time Station Explorer
- **Objective**: Halaman penjelajah stasiun charger SPKLU berbasis peta interaktif dan filter multi-kriteria.
- **Key Features**:
  - Peta lokasi stasiun (Leaflet / Google Maps / Mapbox) dengan marker ketersediaan slot.
  - Filter tipe konektor (`CCS2`, `CHAdeMO`, `Type 2`, `AC Type 2`).
  - Search bar berdasarkan nama stasiun / kota (contoh: SPKLU Gambir, Makassar, Surabaya).
  - Tampilan *Active Tariff badge* (Rp 2.467/kWh + Biaya Layanan ESDM).

---

### **`FRONTEND-03`**: Slot Time Picker & Anti-Overlap Booking Engine UI
- **Objective**: Interaksi jadwal booking slot charger yang interaktif dengan preventif bentrok jadwal di UI.
- **Key Features**:
  - Timeline visual slot charger (slot 1, slot 2, slot 3).
  - Selection range waktu booking (Waktu Mulai & Waktu Selesai).
  - Validasi *Anti-Overlap* langsung di sisi UI sebelum kirim HTTP POST `X-Idempotency-Key`.
  - Notifikasi visual jika slot pilihan berbenturan dengan booking pengguna lain.

---

### **`FRONTEND-04`**: FIFO Waitlist Banner & Live Queue Position Indicator
- **Objective**: Antarmuka antrean *Waitlist* saat slot pilihan pengguna penuh.
- **Key Features**:
  - Tombol **"Gabung Antrean Waitlist"** saat slot *FULL*.
  - Live Queue Card: Menampilkan nomor antrean pengguna (misal: "Anda Nomor Antrean #1").
  - Auto-update status ketika antrean dipromosikan (*WaitlistPromoted*) menjadi `CONFIRMED`.

---

### **`FRONTEND-05`**: Active Charging Telemetry Gauge & Live kWh Dashboard
- **Objective**: Dashboard real-time saat mobil listrik sedang diisi daya (*IN_PROGRESS*).
- **Key Features**:
  - Circular Meter Gauge animasi untuk mengukur persentase daya & energi terkonsumsi (kWh).
  - Stat Card real-time: Daya (`kW`), Tegangan (`V`), Arus (`A`).
  - Tombol **"Selesaikan Pengisian (Finish Charging)"** yang memicu kalkulasi final billing.

---

### **`FRONTEND-06`**: ESDM Tariff Breakdown & QRIS Payment Checkout Modal
- **Objective**: Modal rincian tagihan resmi ESDM dan sistem pembayaran idempoten.
- **Key Features**:
  - Rincian biaya: $\text{Energi Base} + \text{Biaya Layanan ESDM (Fast / Ultra-Fast)} + \text{PPN 11\%}$.
  - QR Code QRIS dinamis untuk pembayaran mock.
  - Penanganan status pembayaran idempoten (`UNPAID` $\rightarrow$ `PAID`) secara real-time via WebSocket / polling.

---

### **`FRONTEND-07`**: User Vehicle Garage & Connector Compatibility Manager
- **Objective**: Manajemen kendaraan listrik pengguna di profil akun.
- **Key Features**:
  - Form pendaftaran kendaraan (Merk, Model, Plat Nomor, Kapasitas Baterai kWh).
  - Selector tipe konektor bawaan mobil (`CCS2`, `CHAdeMO`, `Type 2`).
  - Badge kesesuaian (*Compatibility Match*) saat memilih stasiun pengisian.

---

### **`FRONTEND-08`**: Operator Station Control & No-Show Monitoring Dashboard
- **Objective**: Portal pemantauan stasiun untuk operator SPKLU.
- **Key Features**:
  - Overview statistik: Total Stasiun, Total Booking Hari Ini, KWh Terdistribusi, Pendapatan.
  - Manual trigger & monitoring worker *Auto-Release No-Show*.
  - Tabel outbox event streaming & audit log transaksi.
