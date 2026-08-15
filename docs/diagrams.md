# Diagram Arsitektur dan Alur

## 1. Container Diagram (High Level)

```mermaid
flowchart LR
    U[User App] --> APIGW[API Gateway]

    APIGW --> ST[station-service]
    APIGW --> BK[booking-service]
    APIGW --> SS[session-service]
    APIGW --> BL[billing-service]

    ST --> STDB[(Station DB)]
    BK --> BKDB[(Booking DB)]
    SS --> SSDB[(Session DB)]
    BL --> BLDB[(Billing DB)]

    BK --> EB[(Event Bus)]
    SS --> EB
    BL --> EB

    EB --> BK
    EB --> SS
    EB --> BL

    BK --> REDIS[(Availability Cache)]
```

## 2. Sequence Diagram Booking Saat Jam Sibuk

```mermaid
sequenceDiagram
    participant User
    participant Booking as booking-service
    participant DB as Booking DB
    participant WL as Waitlist
    participant Bus as Event Bus

    User->>Booking: POST /bookings (station, start, end)
    Booking->>DB: cek overlap slot dalam transaksi

    alt Slot tersedia
        Booking->>DB: simpan booking CONFIRMED
        Booking->>Bus: publish BookingConfirmed
        Booking-->>User: 201 CONFIRMED
    else Slot penuh
        Booking->>WL: enqueue waitlist FIFO
        Booking-->>User: 202 WAITLISTED
    end
```

## 3. Sequence Diagram No-show Auto Release

```mermaid
sequenceDiagram
    participant Scheduler as booking-scheduler
    participant Booking as booking-service
    participant DB as Booking DB
    participant Bus as Event Bus
    participant Waitlist as waitlist-processor

    Scheduler->>Booking: trigger periodic scan
    Booking->>DB: cari booking CONFIRMED melewati grace period
    Booking->>DB: update status EXPIRED_NO_SHOW
    Booking->>Bus: publish BookingExpiredNoShow
    Bus->>Waitlist: consume BookingExpiredNoShow
    Waitlist->>DB: promote antrean berikutnya ke CONFIRMED
```

## 4. Sequence Diagram Session dan Billing

```mermaid
sequenceDiagram
    participant User
    participant Session as session-service
    participant Booking as booking-service
    participant SDB as Session DB
    participant Bus as Event Bus
    participant Billing as billing-service
    participant Station as station-service
    participant BDB as Billing DB

    User->>Session: start charging
    Session->>Booking: validasi booking check-in
    Session->>SDB: create session STARTED
    Session->>Bus: publish SessionStarted

    User->>Session: finish charging
    Session->>SDB: update consumedKwh + ENDED
    Session->>Bus: publish SessionFinished

    Bus->>Billing: consume SessionFinished
    Billing->>Station: get applicable tariff
    Billing->>BDB: create invoice
    Billing->>Bus: publish InvoiceCreated
```

## 5. State Diagram Booking

```mermaid
stateDiagram-v2
    [*] --> REQUESTED
    REQUESTED --> CONFIRMED: slot tersedia
    REQUESTED --> WAITLISTED: slot penuh

    WAITLISTED --> CONFIRMED: dipromosikan
    CONFIRMED --> IN_SESSION: check-in + start session
    CONFIRMED --> EXPIRED_NO_SHOW: lewat grace period
    CONFIRMED --> CANCELLED: dibatalkan user

    IN_SESSION --> COMPLETED: session selesai
    CANCELLED --> [*]
    EXPIRED_NO_SHOW --> [*]
    COMPLETED --> [*]
```
