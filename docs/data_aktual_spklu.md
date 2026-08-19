# Data Aktual & Sumber Valid Sistem Booking SPKLU (EV Charger) Indonesia

Dokumen ini berisi data aktual, regulasi resmi, skema tarif, biaya layanan, spesifikasi teknologi charger, serta populasi Kendaraan Bermotor Listrik Berbasis Baterai (KBLBB) di Indonesia. Data ini disusun untuk keperluan riset, perancangan database (`station-service`, `billing-service`), dan validasi sistem pada website booking charger.

---

## 1. Regulasi Resmi Pemerintah Indonesia

Pengisian daya listrik pada Stasiun Pengisian Kendaraan Listrik Umum (SPKLU) diatur secara resmi oleh Kementerian Energi dan Sumber Daya Mineral (ESDM) Republik Indonesia:

1. **Peraturan Menteri ESDM Nomor 1 Tahun 2023**  
   *Tentang Penyediaan Infrastruktur Pengisian Listrik untuk Kendaraan Bermotor Listrik Berbasis Baterai.*  
   Mengatur skema operasional, perizinan (NIB), standar keselamatan, dan klasifikasi SPKLU.
2. **Keputusan Menteri ESDM Nomor 182.K/TL.04/MEM.S/2023**  
   *Tentang Petunjuk Teknis Biaya Layanan Pengisian Listrik Pada Stasiun Pengisian Kendaraan Listrik Umum.*  
   Mengatur batas atas biaya layanan (*service fee*) per transaksi untuk pengisian cepat (*Fast*) dan sangat cepat (*Ultra Fast*).

---

## 2. Biaya Layanan per Transaksi (Kepmen ESDM No. 182.K/2023)

Biaya yang dikenakan kepada pengguna kendaraan listrik terdiri dari **2 Komponen Utama**:
- **Biaya Pemakaian Listrik**: kWh Terpakai x Tarif Listrik per kWh
- **Biaya Layanan SPKLU**: Dikenakan per 1x transaksi pengisian daya

| Kategori Pengisian | Spesifikasi Daya | Jenis Arus | Biaya Layanan Maksimal (Sebelum PPN) | Biaya Layanan Maksimal (Plus PPN 11%) |
| :--- | :--- | :--- | :--- | :--- |
| **Slow Charging** | <= 7 kW | AC | Rp 0 (Bebas Biaya Layanan) | Rp 0 |
| **Medium / Normal** | > 7 kW s.d. <= 22 kW | AC | Rp 0 (Bebas Biaya Layanan) | Rp 0 |
| **Fast Charging** | 25 kW s.d. 50 kW | DC | Rp 25.000 / transaksi | Rp 27.750 / transaksi |
| **Ultra Fast Charging** | > 50 kW | DC | Rp 57.000 / transaksi | Rp 63.270 / transaksi |

> **Catatan Pajak:**
> - **PPN (Pajak Pertambahan Nilai)**: 11%.
> - **PPJ / PBJT (Pajak Penerangan Jalan)**: 3% - 10% (tergantung Perda Pemerintah Daerah setempat).
> - **Tarif Dasar Listrik SPKLU (PLN)**: Maksimal Rp 2.467 / kWh (Faktor Pengali 1,5).

---

## 3. Komparasi Tarif Operator SPKLU (PLN & Swasta)

| Operator | Jenis Pengisian | Skema Tarif | Biaya per kWh / Sesi |
| :--- | :--- | :--- | :--- |
| **PLN Charge.IN** | AC Slow / DC Fast / Ultra Fast | Per kWh + Biaya Layanan ESDM | Rp 2.467 / kWh (+ Service Fee Rp 25rb / Rp 57rb) |
| **Voltron** | AC / DC Fast (< 50 kW) | Per kWh + Admin Fee | Rp 2.466 / kWh (+ Admin Rp 4.000 + PPJ 5% + PPN 11%) |
| **Voltron** | DC Ultra Fast (>= 50 kW) | Per kWh + Admin Fee | Rp 3.600 / kWh (+ Admin Rp 4.000 + PPJ 5% + PPN 11%) |
| **Shell Recharge** | DC Fast Charging | Per kWh + Biaya Layanan | Rp 2.166 / kWh (+ Service Fee Rp 25.000 / transaksi) |
| **Casion** | AC / DC Charging | Berbasis Waktu (Time-based) | Rp 1.000 / menit (Rp 15.000 / 15 menit) |

---

## 4. Tipe Charger & Konektor Standar di Indonesia

| Tipe Konektor | Jenis Arus | Penggunaan Umum | Kendaraan Kompatibel |
| :--- | :--- | :--- | :--- |
| **Type 2 (Mennekes)** | AC | Home / Public Slow-Medium | Hyundai Ioniq 5/6, BYD Seal/Atto 3, Chery, MG, BMW |
| **CCS2** | DC | Public Fast & Ultra Fast | Standar Utama SPKLU Indonesia (Hyundai, BYD, MG, Chery, Wuling Binguo DC) |
| **CHAdeMO** | DC | Public Fast Charging | Nissan Leaf, Mitsubishi Outlander PHEV |
| **GB/T AC/DC** | AC / DC | Standar China khusus | Wuling Air EV (Standard/Long Range AC), Seres E1 |

---

## 5. Data Populasi Kendaraan Listrik & SPKLU Indonesia

| Tahun | Populasi KBLBB | Jumlah SPKLU (Unit) | Jumlah Lokasi SPKLU | Target SPKLU ESDM |
| :---: | :---: | :---: | :---: | :---: |
| **2021** | 14.400 | 267 | 170 | 1.000 |
| **2022** | 35.200 | 570 | 411 | 3.000 |
| **2023** | 116.400 | 1.081 | 748 | 7.500 |
| **2024** | 207.478 | 3.202 | 2.180 | 20.000 |
| **2025** | 333.561 | 4.655 | 3.007 | 35.000 |
| **2030** | 2.000.000 | 62.918 | 40.000 | 62.918 |
| **2034** | 7.000.000 | 192.253 | 120.000 | 192.253 |

---

## 6. Daftar Sumber Valid & Referensi Resmi

1. **Kementerian ESDM RI (Siaran Pers & JDIH ESDM)**
   - [Peraturan Menteri ESDM No. 1 Tahun 2023](https://www.esdm.go.id)
   - [Keputusan Menteri ESDM No. 182.K/TL.04/MEM.S/2023](https://www.esdm.go.id)
2. **Portal Resmi Informasi Publik Pemerintah RI**
   - [InfoPublik - Biaya Layanan Pengisian Listrik SPKLU](https://infopublik.id)
   - [Indonesia Baik - Aturan Tarif SPKLU](https://indonesiabaik.id)
3. **PT PLN (Persero)**
   - [PLN Charge.IN & Tarif Listrik Layanan Khusus](https://www.pln.co.id)
4. **Media Berita Nasional Terverifikasi**
   - [ANTARA News - Rincian Tarif & Biaya Layanan SPKLU](https://www.antaranews.com)
   - [Kompas.com - Regulasi Biaya Layanan Fast Charging](https://www.kompas.com)
5. **Gaikindo & Kementerian Perhubungan RI**
   - Data Statistik Kendaraan Listrik Indonesia: [Gaikindo.or.id](https://www.gaikindo.or.id)
