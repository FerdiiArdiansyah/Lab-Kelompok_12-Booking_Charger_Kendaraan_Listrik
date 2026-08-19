# 📐 Entity Relationship Diagram (ERD) & Database Schemas

Dokumen ini menjelaskan struktur data dan hubungan entitas di setiap database microservice yang dikelola oleh **Data & Persistence Engineer**.

---

## 1. `station_db` (station-service)

Mengelola master data stasiun pengisian listrik, slot konektor, dan struktur tarif.

```mermaid
erDiagram
    STATIONS ||--|{ CHARGER_SLOTS : "has slots"
    STATIONS ||--|{ TARIFFS : "has pricing"

    STATIONS {
        string id PK
        string name
        string location
        decimal latitude
        decimal longitude
        decimal total_power_kw
        string status
        timestamp created_at
        timestamp updated_at
    }

    CHARGER_SLOTS {
        string id PK
        string station_id FK
        int slot_number
        string connector_type
        decimal max_power_kw
        string status
        timestamp created_at
        timestamp updated_at
    }

    TARIFFS {
        string id PK
        string station_id FK
        decimal price_per_kwh
        string currency
        timestamp valid_from
        timestamp valid_to
        boolean is_active
        timestamp created_at
    }
```

---

## 2. `booking_db` (booking-service)

Mengelola booking slot, validasi jadwal bebas bentrok, antrean waitlist, dan event outbox.

```mermaid
erDiagram
    BOOKINGS ||--o| WAITLISTS : "promotes to"
    BOOKINGS ||--o| OUTBOX_EVENTS : "emits"

    BOOKINGS {
        string id PK
        string user_id
        string station_id
        string slot_id
        timestamp start_time
        timestamp end_time
        tstzrange booking_period "STORED GENERATED"
        string status
        string idempotency_key UK
        timestamp created_at
        timestamp updated_at
    }

    WAITLISTS {
        string id PK
        string station_id
        string user_id
        timestamp requested_start
        timestamp requested_end
        int queue_number
        string status
        timestamp created_at
        timestamp updated_at
    }

    IDEMPOTENCY_KEYS {
        string key PK
        string request_hash
        jsonb response_payload
        int status_code
        timestamp created_at
        timestamp expires_at
    }

    OUTBOX_EVENTS {
        string id PK
        string aggregate_type
        string aggregate_id
        string event_type
        jsonb payload
        string status
        timestamp created_at
        timestamp processed_at
    }
```

---

## 3. `session_db` (session-service)

Mengelola data telemetry sesi charging real-time dan log konsumsi energi kWh.

```mermaid
erDiagram
    CHARGING_SESSIONS ||--|{ METER_READINGS : "records"

    CHARGING_SESSIONS {
        string id PK
        string booking_id UK
        string slot_id
        string user_id
        timestamp started_at
        timestamp ended_at
        decimal consumed_kwh
        string status
        timestamp created_at
        timestamp updated_at
    }

    METER_READINGS {
        bigint id PK
        string session_id FK
        timestamp recorded_at
        decimal current_kwh
        decimal power_kw
        decimal voltage
        decimal current_ampere
    }
```

---

## 4. `billing_db` (billing-service)

Mengelola invoice, transaksi pembayaran, log audit keuangan, dan outbox event.

```mermaid
erDiagram
    INVOICES ||--|{ PAYMENTS : "paid by"

    INVOICES {
        string id PK
        string session_id UK
        string user_id
        string tariff_id
        decimal consumed_kwh
        decimal price_per_kwh
        decimal subtotal
        decimal tax
        decimal total
        string status
        timestamp created_at
        timestamp updated_at
    }

    PAYMENTS {
        string id PK
        string invoice_id FK
        string payment_method
        decimal amount
        string status
        string transaction_ref
        timestamp paid_at
        timestamp created_at
    }

    AUDIT_LOGS {
        bigint id PK
        string entity_name
        string entity_id
        string action
        jsonb old_value
        jsonb new_value
        string performed_by
        timestamp created_at
    }
```
