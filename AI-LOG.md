# AI-LOG — Catatan Jejak Penggunaan Copilot / AI

Dokumen ini mencatat riwayat interaksi kelompok dengan GitHub Copilot / AI selama perancangan, eksekusi, dan optimasi sistem **VoltHub (Booking Charger Kendaraan Listrik)** dari Lapisan 1 hingga Lapisan 3.

---

### [Arsitek] · Lapisan 1 · Entri 1
- **Konteks:** Merancang struktur arsitektur microservices untuk sistem booking SPKLU.
- **Prompt:** `Buatkan struktur folder dan docker-compose.yml untuk 6 microservices Go berbasis database-per-service untuk SPKLU.`
- **Diterima:** Pemisahan folder layanan (`auth`, `station`, `booking`, `payment`, `notification`, `telemetry`) serta alokasi port 8081-8086.
- **Ditolak:** Opsi menggunakan 1 database PostgreSQL terpusat untuk seluruh service.
- **Alasan Penolakan:** Melanggar prinsip *Database-per-Service* pada arsitektur microservices yang dapat memicu *tight coupling* dan ketergantungan antar-modul.
- **Verifikasi:** Memastikan tiap layanan memiliki kontainer database PostgreSQL independen di file `docker-compose.yml`.

---

### [Backend] · Lapisan 1 · Entri 2
- **Konteks:** Pembuatan handler otentikasi JWT pada Auth Service (Port 8086).
- **Prompt:** `Buatkan fungsi handler login di Go yang mengecek email password dan mengembalikan JWT token.`
- **Diterima:** Struktur skema claims JWT dan hashing password bcrypt.
- **Ditolak:** Menyimpan *secret key* JWT secara *hardcoded* di dalam berkas kode.
- **Alasan Penolakan:** Menimbulkan celah keamanan serius (*vulnerability*) jika kode di-push ke repository publik.
- **Verifikasi:** Mengubah variabel *secret key* agar dibaca dari *Environment Variable* (`process.env.JWT_SECRET`).

---

### [QA] · Lapisan 1 · Entri 3
- **Konteks:** Pembuatan skrip *Smoke Test* otomatis untuk verifikasi kesehatan endpoint kritis.
- **Prompt:** `Buatkan skrip smoke test menggunakan node:test native untuk mengecek kesehatan endpoint login dan stasiun.`
- **Diterima:** Penggunaan *native runner* `node:test` dan `node:assert` tanpa dependensi eksternal.
- **Ditolak:** Penggunaan framework `Jest` dan pembuatan *mock response*.
- **Alasan Penolakan:** Smoke test harus menguji server sungguhan yang sedang berjalan (*live integration test*) secara ringan, bukan menguji data tiruan (*mock*).
- **Verifikasi:** Menjalankan `node --test tests/smoke.test.js` dengan hasil `pass 3, fail 0`.

---

### [Backend] · Lapisan 2 · Entri 4
- **Konteks:** Mengatasi kendala *race condition* / *double booking* slot charger saat lonjakan trafik.
- **Prompt:** `Bagaimana cara mencegah race condition booking slot charger yang sama di PostgreSQL saat diakses bersamaan?`
- **Diterima:** Penerapan transaksi SQL dengan *Pessimistic Locking* (`SELECT ... FOR UPDATE`) pada tabel slot charger.
- **Ditolak:** Pengecekan ketersediaan slot di level aplikasi (*in-memory check*).
- **Alasan Penolakan:** Pengecekan di memori aplikasi tidak *thread-safe* saat backend berjalan dalam *multiple instances/replicas*, sehingga *overbooking* tetap terjadi.
- **Verifikasi:** Memastikan konsistensi ketersediaan slot setelah beberapa request berjalan serentak.

---

### [QA] · Lapisan 2 · Entri 5
- **Konteks:** Menguji ketahanan sistem terhadap penyerbuan 300 request simultan pada sumber daya slot charger.
- **Prompt:** `Buatkan skrip penyerbuan serentak 300 request bersamaan di Node.js menggunakan Promise.all untuk menguji endpoint booking.`
- **Diterima:** Struktur `Promise.all` untuk mengeksekusi request login dan booking secara paralel.
- **Ditolak:** Penambahan jeda waktu (*delay/sleep*) antar-request.
- **Alasan Penolakan:** Penambahan delay membatalkan skenario *concurrency test* yang mewajibkan lonjakan eksekusi tiba pada milidetik yang sama.
- **Verifikasi:** Menjalankan `node tests/rebutan.test.js` dengan hasil 300/300 request sukses (`0 error`).

---

### [DevOps] · Lapisan 2 · Entri 6
- **Konteks:** Mengoptimalkan throughput server saat pengujian beban tinggi.
- **Prompt:** `Bagaimana konfigurasi rate limiting dan connection pooling di Go untuk menangani 2000 req/sec?`
- **Diterima:** Pengaturan batas koneksi `SetMaxOpenConns(50)` dan `SetMaxIdleConns(25)` pada driver PostgreSQL.
- **Ditolak:** Menaikkan `SetMaxOpenConns` hingga 1000 koneksi.
- **Alasan Penolakan:** Menyebabkan database mengalami *Out-Of-Memory* (OOM) akibat keterbatasan RAM pada mesin uji.
- **Verifikasi:** Pengujian `autocannon` menunjukkan latensi p50 di 4.299 ms tanpa koneksi terputus.

---

### [QA] · Lapisan 3 · Entri 7
- **Konteks:** Konfigurasi pemetaan port pada lingkungan GitHub Codespaces untuk pengujian antarmuka.
- **Prompt:** `Buatkan command gh CLI untuk mengubah port visibility microservices menjadi public sekaligus.`
- **Diterima:** Perintah CLI `gh codespace ports visibility <port>:public`.
- **Ditolak:** Pengubahan akses port satu per satu secara manual via antarmuka grafik VS Code.
- **Alasan Penolakan:** Proses manual lambat dan berisiko melewatkan salah satu port dari 6 microservices yang ada.
- **Verifikasi:** Seluruh port (8081-8086 & 5173) terkonfirmasi berstatus *Public* dan dapat dipanggil oleh *frontend*.