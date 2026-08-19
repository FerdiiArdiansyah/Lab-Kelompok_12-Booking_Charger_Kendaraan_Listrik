# 💾 Kebijakan Backup, Recovery, & Retensi Data

## 1. Strategi Backup (RPO & RTO Targets)

- **Recovery Point Objective (RPO)**: < 5 menit (Maksimum kehilangan data saat insiden fatal).
- **Recovery Time Objective (RTO)**: < 30 menit (Maksimum waktu downtime pemulihan sistem).

### Metode Backup:
1. **Full Backup Daily**: `pg_dumpall` atau snapshot volume setiap jam 02:00 UTC.
2. **Continuous WAL Archiving (Point-In-Time Recovery / PITR)**:
   - Mengaktifkan `wal_level = replica` & `archive_mode = on`.
   - WAL log di-stream secara kontinu ke S3/GCS Object Storage terenskripsi (AES-256).

---

## 2. Kebijakan Retensi Data (Data Archiving Schedule)

| Database | Nama Tabel | Masa Simpan Live (Hot Data) | Kebijakan Archiving (Warm/Cold Data) |
| :--- | :--- | :--- | :--- |
| `station_db` | `stations`, `charger_slots`, `tariffs` | Selamanya | Audit history tarif > 2 tahun dipindah ke Cold Storage. |
| `booking_db` | `bookings`, `waitlists` | 90 Hari | Booking status `COMPLETED`/`CANCELLED` > 90 hari di-archive ke Data Warehouse (BigQuery). |
| `booking_db` | `outbox_events` | 7 Hari | Hapus event status `PUBLISHED` > 7 hari. |
| `session_db` | `meter_readings` | 30 Hari | Telemetry interval 5 detik diagregasi per jam setelah 30 hari. |
| `billing_db` | `invoices`, `payments` | 7 Tahun | Wajib disimpan 7 tahun sesuai Regulasi Kepatuhan Pajak & Keuangan. |
