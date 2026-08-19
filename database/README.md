# 🗄️ Data & Persistence Engineering Repository

Selamat datang di repositori utama **Data & Persistence Engineer** untuk **Sistem Booking Charger Kendaraan Listrik**.

Repositori ini memuat seluruh aset, skema database, constraint anti-bentrok, pola persistence, dokumentasi ERD, serta skrip operasional database yang menjadi tanggung jawab **Data & Persistence Engineer**.

---

## 🎯 Tanggung Jawab Utama Data & Persistence Engineer

1. **Database-per-Service Architecture**: Mengelola & memelihara skema database independen untuk setiap microservice (`station_db`, `booking_db`, `session_db`, `billing_db`).
2. **Strict Anti-Overlap Engine**: Menerapkan PostgreSQL *Exclusion Constraints* (`btree_gist` + `tstzrange`) untuk menjamin 0% booking ganda pada slot charger di level transaksi database.
3. **Eventual Consistency & Outbox Pattern**: Menyediakan struktur tabel `outbox_events` pada setiap database untuk mendukung integrasi *Event-Driven Architecture* yang reliable.
4. **Idempotency & Auditability**: Menyusun skema tabel idempotency key untuk API kritikal serta audit logging untuk melacak setiap perubahan state data sensitif (booking status, payment, invoice).
5. **High Performance & Indexing**: Optimasi indeks untuk pencarian lokasi stasiun, query telemetry meteran kWh, serta transaksi pembayaran.
6. **Data Lifecycle & Compliance**: Menyusun kebijakan backup, Point-in-Time Recovery (PITR), serta aturan retensi data.

---

## 📁 Struktur Folder Tanggung Jawab (`/database`)

```text
database/
├── README.md                          # Dokumentasi utama Data & Persistence Engineer
├── docker-compose.db.yml              # Environment local multi-database (PostgreSQL + PgBouncer)
├── init-scripts/                      # Skrip inisialisasi database & user postgres
│   └── 01-init-databases.sql
├── schemas/                           # DDL, Indeks, & Seed Data per Service Domain
│   ├── station-db/                    # Master data stasiun, slot charger, & tarif
│   │   ├── 01_schema.sql
│   │   ├── 02_indexes.sql
│   │   └── seed.sql
│   ├── booking-db/                    # Core booking, exclusion constraint, waitlist
│   │   ├── 01_schema.sql
│   │   ├── 02_anti_overlap_constraint.sql
│   │   └── seed.sql
│   ├── session-db/                    # Telemetry charging session & kWh meter readings
│   │   ├── 01_schema.sql
│   │   ├── 02_indexes.sql
│   │   └── seed.sql
│   └── billing-db/                    # Invoices, payments, & financial audit trail
│       ├── 01_schema.sql
│       ├── 02_indexes.sql
│       └── seed.sql
├── patterns/                          # Standardized Data Access & Storage Patterns
│   ├── transactional-outbox.sql       # Skema & query pattern outbox event publishing
│   ├── idempotency-store.sql          # Tabel pencatatan idempotency token API
│   └── audit-logging.sql              # Triggers & CDC audit trail
├── docs/                              # Dokumentasi Arsitektur Data & Panduan Teknik
│   ├── ERD.md                         # Entity Relationship Diagrams (Mermaid format)
│   ├── anti-overlap-strategy.md       # Spesifikasi teknis btree_gist & range exclusion
│   ├── migration-guide.md             # Standard operating procedure (SOP) migrasi zero-downtime
│   └── backup-retention-policy.md     # Kebijakan backup, recovery, & retensi data
└── scripts/                           # Utility scripts untuk development & DevOps
    ├── run-migrations.sh              # Skrip otomasi eksekusi migrasi DDL
    └── seed-all.sh                    # Skrip otomasi seeding data awal
```

---

## ⚡ Ringkasan Peta Database Per Microservice

| Service | Database Name | Primary Engine | Core Entities & Constraints |
| :--- | :--- | :--- | :--- |
| **station-service** | `station_db` | PostgreSQL 16 | `stations`, `charger_slots`, `tariffs` |
| **booking-service** | `booking_db` | PostgreSQL 16 (`btree_gist`) | `bookings` (Anti-Overlap Range Exclusion), `waitlists`, `outbox_events`, `idempotency_keys` |
| **session-service** | `session_db` | PostgreSQL 16 | `charging_sessions`, `meter_readings` (Time-series telemetry), `outbox_events` |
| **billing-service** | `billing_db` | PostgreSQL 16 | `invoices`, `payments`, `audit_logs`, `outbox_events` |

---

## 🚀 Panduan Cepat Menjalankan Environment Database Local

### 1. Jalankan PostgreSQL Services
```bash
docker-compose -f database/docker-compose.db.yml up -d
```

### 2. Jalankan Migrasi DDL & Seed Data
```bash
bash database/scripts/run-migrations.sh
bash database/scripts/seed-all.sh
```

---

## 📑 Referensi & Dokumentasi Lengkap
- 📐 [Entity Relationship Diagram (ERD)](docs/ERD.md)
- 🔒 [Mekanisme Anti-Overlap Booking](docs/anti-overlap-strategy.md)
- 🔄 [Panduan Zero-Downtime Migration](docs/migration-guide.md)
- 💾 [Kebijakan Backup & Retensi Data](docs/backup-retention-policy.md)
