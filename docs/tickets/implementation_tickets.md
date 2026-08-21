# Implementation Tickets & TDD Test Plan
## System Booking Charger Kendaraan Listrik (SPKLU)

- **Source Document**: [`docs/prd_booking_spklu_system.md`](file:///home/aseppp/Documents/Lab-Kelompok_12-Booking_Charger_Kendaraan_Listrik/docs/prd_booking_spklu_system.md)
- **Target Architecture**: Go Clean Architecture (Domain -> Repository -> Usecase -> HTTP Handler)
- **Methodology**: Test-Driven Development (TDD)

---

## 📋 Summary of Tickets

| Ticket ID | Target Service | Title | Priority | Est. Complexity |
|---|---|---|---|---|
| **`TICKET-01`** | `booking-service` | Anti-Overlap Booking Validation Engine & Unit Tests | `P0 (Critical)` | ✅ DONE |
| **`TICKET-02`** | `booking-service` | Auto-Release No-Show Worker & Event Trigger | `P0 (Critical)` | ✅ DONE |
| **`TICKET-03`** | `booking-service` | FIFO Waitlist Queue Allocation Engine | `P1 (High)` | ✅ DONE |
| **`TICKET-04`** | `session-service` | Real-Time Telemetry & Session Finish Event Publisher | `P1 (High)` | ✅ DONE |
| **`TICKET-05`** | `billing-service` | ESDM Tariff & PPN 11% Calculation Engine (TDD Core) | `P0 (Critical)` | ✅ DONE |
| **`TICKET-06`** | `billing-service` | Idempotent Payment Processing & Webhook Handler | `P1 (High)` | ✅ DONE |
| **`TICKET-07`** | `user-service` / `station-service` | Vehicle Connector Compatibility & Master Tariff Sync | `P2 (Medium)` | ✅ DONE |

---

## 🎟️ Detailed Ticket Specifications

### `TICKET-01`: Anti-Overlap Booking Validation Engine & Unit Tests
- **Service**: `services/booking-service`
- **File Focus**: `internal/usecase/booking_usecase.go`, `internal/repository/postgres/booking_repository.go`
- **Objective**: Mencegah pemesanan slot charger yang sama pada rentang waktu `start_time` dan `end_time` yang bertabrakan.

#### 🧪 TDD Specification:
1. **Red Stage (Test Creation)**:
   - Buat unit test `TestBookingUsecase_CreateBooking_OverlapConflict(t *testing.T)`.
   - Setup mock repository dengan data booking eksisting: Slot `S1`, Rentang `10:00 - 11:00`.
   - Jalankan test request pembuatan booking baru: Slot `S1`, Rentang `10:30 - 11:30`.
   - Assert: Ekspektasi mengembalikan error `ErrSlotUnavailable` / status code `409 Conflict`.
2. **Green Stage (Implementation)**:
   - Implementasikan query database `WHERE slot_id = $1 AND status NOT IN ('CANCELLED', 'EXPIRED_NO_SHOW') AND start_time < $3 AND end_time > $2`.
3. **Refactor Stage**:
   - Ekstrak logic validasi waktu ke pure function `ValidateBookingTimeSlot(start, end time.Time, existing []Booking) bool`.

---

### `TICKET-02`: Auto-Release No-Show Worker & Event Trigger
- **Service**: `services/booking-service`
- **File Focus**: `internal/usecase/booking_usecase.go`, `main.go`
- **Objective**: Otomatis mengubah status booking `CONFIRMED` menjadi `EXPIRED_NO_SHOW` jika pengguna tidak melakukan check-in setelah *grace period* (misal: 15 menit).

#### 🧪 TDD Specification:
1. **Red Stage (Test Creation)**:
   - Buat unit test `TestBookingUsecase_AutoReleaseNoShow(t *testing.T)`.
   - Setup mock data: Booking `B1` status `CONFIRMED`, `start_time` = 20 menit lalu, `check_in_at` = nil.
   - Assert: Sesi booking berubah menjadi `EXPIRED_NO_SHOW` dan melepaskan ketersediaan slot.
2. **Green Stage (Implementation)**:
   - Tambahkan method `AutoReleaseNoShow(ctx context.Context, gracePeriodMinutes int) ([]Booking, error)`.
   - Tambahkan outbox event `booking.expired` untuk memberi sinyal ke `station-service`.

---

### `TICKET-03`: FIFO Waitlist Queue Allocation Engine
- **Service**: `services/booking-service`
- **File Focus**: `internal/usecase/waitlist_usecase.go`
- **Objective**: Memasukkan pengguna ke antrean FIFO jika slot penuh dan mengalokasikan slot otomatis saat ada rilis/pembatalan booking.

#### 🧪 TDD Specification:
1. **Red Stage (Test Creation)**:
   - Buat unit test `TestWaitlist_FIFO_Allocation(t *testing.T)`.
   - Setup mock waitlist queue: User A (`queue_number: 1`), User B (`queue_number: 2`).
   - Trigger event slot released.
   - Assert: User A dipromosikan statusnya menjadi `CONFIRMED`, User B tetap `WAITLISTED`.
2. **Green Stage (Implementation)**:
   - Tulis logic `PromoteNextWaitlist(ctx context.Context, stationID, slotID uint)` mengurutkan berdasarkan `created_at ASC`.

---

### `TICKET-04`: Real-Time Telemetry & Session Finish Event Publisher
- **Service**: `services/session-service`
- **File Focus**: `internal/usecase/session_usecase.go`
- **Objective**: Menerima telemetry kWh meter dan mempublikasikan event penutupan sesi `session.finished` beserta `total_consumed_kwh`.

#### 🧪 TDD Specification:
1. **Red Stage (Test Creation)**:
   - Buat unit test `TestSessionUsecase_FinishSession_PublishesEvent(t *testing.T)`.
   - Jalankan `FinishSession(sessionID)`.
   - Assert: `consumed_kwh` dihitung akurat dari `last_meter - initial_meter`, status berubah `COMPLETED`.
2. **Green Stage (Implementation)**:
   - Implementasikan method `FinishSession` dan penyimpanan telemetry di `MeterReading` table.

---

### `TICKET-05`: ESDM Tariff & PPN 11% Calculation Engine (TDD Core)
- **Service**: `services/billing-service`
- **File Focus**: `internal/usecase/invoice_usecase.go`
- **Objective**: Menghitung invoice tagihan berdasarkan regulasi biaya layanan ESDM (Fast / Ultra-Fast Charging) dan PPN 11%.

#### 🧪 TDD Specification:
1. **Red Stage (Test Creation)**:
   - Buat unit test `TestCalculateInvoiceAmount(t *testing.T)`.
   - Input Test Case 1: `consumed_kwh = 20`, `price_per_kwh = 2467`, `service_fee = 25000` (Fast Charging).
     - Expected Subtotal: $20 \times 2467 + 25000 = 74.340$.
     - Expected PPN 11%: $74.340 \times 0.11 = 8.177,4$.
     - Expected Total: $82.517,4$.
   - Assert: Hasil kalkulasi persis sesuai ekspektasi math.
2. **Green Stage (Implementation)**:
   - Tulis fungsi murni `CalculateInvoice(consumedKWh float64, pricePerKWh float64, serviceFee float64) (subtotal, tax, total float64)`.

---

### `TICKET-06`: Idempotent Payment Processing & Webhook Handler
- **Service**: `services/billing-service`
- **File Focus**: `internal/usecase/payment_usecase.go`, `internal/delivery/http/payment_handler.go`
- **Objective**: Menjamin inisiasi pembayaran aman dari duplikasi (*double payment*) menggunakan header `X-Idempotency-Key`.

#### 🧪 TDD Specification:
1. **Red Stage (Test Creation)**:
   - Buat unit test `TestPayment_IdempotencyKey_DuplicateRequest(t *testing.T)`.
   - Kirim request pembayaran pertama dengan `Idempotency-Key: KEY-123` -> Success.
   - Kirim request pembayaran kedua dengan `Idempotency-Key: KEY-123` -> Return cached response / no double insert.
2. **Green Stage (Implementation)**:
   - Buat middleware/check idempotency table `idempotency_keys(key, response_payload, created_at)`.

---

## 🎨 Frontend Implementation Tickets

Daftar tiket pengembangan komponen Frontend (UI/UX, Station Explorer, Booking Calendar, Telemetry Gauge & Payment Modal) didokumentasikan di:
- 👉 **[`docs/tickets/frontend_tickets.md`](file:///home/aseppp/Documents/Lab-Kelompok_12-Booking_Charger_Kendaraan_Listrik/docs/tickets/frontend_tickets.md)** (`FRONTEND-01` s/d `FRONTEND-08`)

---

## 🚀 Workflow Eksekusi TDD (`/implement`)

Untuk memicu pengerjaan tiket berikutnya, jalankan perintah:
> `Jalankan /implement untuk mengerjakan FRONTEND-01`

