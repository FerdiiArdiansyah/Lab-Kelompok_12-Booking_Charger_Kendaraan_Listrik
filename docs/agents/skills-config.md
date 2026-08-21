# Matt Pocock Skills Repository Setup & Configuration

## 📌 Project Overview
- **Repository**: `Lab-Kelompok_12-Booking_Charger_Kendaraan_Listrik`
- **Architecture**: Event-Driven Microservices (Golang)
- **Domain**: System Booking SPKLU (Electric Vehicle Charging Station)

---

## 🏗️ Services & Ports Mapping
| Service | Port | Clean Architecture Structure |
|---|---|---|
| `user-service` | `:8084` | `internal/domain`, `internal/repository`, `internal/usecase`, `internal/delivery/http` |
| `station-service` | `:8081` | `internal/domain`, `internal/repository`, `internal/usecase`, `internal/delivery/http` |
| `booking-service` | `:8082` | `internal/domain`, `internal/repository`, `internal/usecase`, `internal/delivery/http` |
| `session-service` | `:8083` | `internal/domain`, `internal/repository`, `internal/usecase`, `internal/delivery/http` |
| `billing-service` | `:8085` | `internal/domain`, `internal/repository`, `internal/usecase`, `internal/delivery/http` |

---

## 📜 Key Architectural Rules & Constraints
1. **Architecture Pattern**: Clean Architecture (Domain Layer -> Repository Layer -> Usecase Layer -> HTTP Delivery Handler).
2. **Database & Data Layer**: PostgreSQL (gorm / std sql), Outbox pattern for event publishing.
3. **Idempotency**: All mutation endpoints require/support `idempotency_key`.
4. **Anti-Overlap Logic**: Reservation slot validation in `booking-service`.
5. **Billing Rules**: Calculation based on actual consumed kWh + ESDM Service Fee (Fast/Ultra-Fast) + PPN (11%).

---

## 🛠️ Configured Skills Workflows
- **`/grill-me` / `/grill-with-docs`**: Reviews feature requests against `docs/architecture.md` & `README.md`.
- **`/to-spec`**: Generates RFC/Technical Specs in `docs/` adhering to microservices boundaries.
- **`/to-tickets`**: Breaks down technical specs into step-by-step Go task tickets.
- **`/implement` / `/tdd`**: Implements code following Go Clean Architecture & TDD standards.
