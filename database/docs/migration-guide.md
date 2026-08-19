# 🔄 Panduan Migrasi Database & Zero-Downtime SOP

## 1. Aturan Migrasi DDL (Expand-Contract Pattern)

Sebagai Data & Persistence Engineer, seluruh perubahan skema database harus mematuhi prinsip **Zero-Downtime Deployment**. Jangan pernah menghapus atau mengubah kolom secara mendadak (breaking changes) yang membuat instance service lama crash saat rolling deployment.

### 3-Fase Expand-Contract:

1. **Fase 1 (Expand)**:
   - Tambahkan kolom baru / tabel baru (izinkan `NULL` atau sediakan `DEFAULT`).
   - Deploy versi service baru yang menulis ke kolom lama dan kolom baru.
2. **Fase 2 (Backfill & Dual Write)**:
   - Jalankan script data migration untuk mengisi data lama ke kolom baru.
   - Switch service baru agar membaca data dari kolom baru.
3. **Fase 3 (Contract)**:
   - Setelah verified 100% stabil, lepas kolom lama (drop column) melalui script migrasi terpisah.

---

## 2. Naming Convention File Migrasi

Format penamaan skrip DDL migrasi:
`V{YYYYMMDDHHMMSS}__{deskripsi_singkat}.sql`

Contoh:
- `V20260819100000__create_stations_and_charger_slots.sql`
- `V20260819101500__add_anti_overlap_gist_constraint.sql`

---

## 3. Aturan DDL Aman

- Selalu gunakan `CREATE TABLE IF NOT EXISTS`.
- Saat membuat indeks baru pada database production besar, gunakan `CREATE INDEX CONCURRENTLY` (agar tidak mengunci tabel dari operasi READ/WRITE).
- Jangan gunakan `ALTER TABLE DROP COLUMN` tanpa persetujuan Lead Architecture & fase grace period minimum 7 hari.
