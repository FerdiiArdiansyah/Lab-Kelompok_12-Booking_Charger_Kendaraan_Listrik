-- station-db DDL Schema
-- Service: station-service
-- Database: station_db

CREATE TABLE IF NOT EXISTS stations (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    total_power_kw DECIMAL(8, 2) NOT NULL DEFAULT 0.00,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, INACTIVE, MAINTENANCE
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS charger_slots (
    id VARCHAR(36) PRIMARY KEY,
    station_id VARCHAR(36) NOT NULL REFERENCES stations(id) ON DELETE CASCADE,
    slot_number INT NOT NULL,
    connector_type VARCHAR(50) NOT NULL, -- CCS2, TYPE2, CHADEMO
    max_power_kw DECIMAL(8, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'AVAILABLE', -- AVAILABLE, OCCUPIED, MAINTENANCE, OUT_OF_SERVICE
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_station_slot_number UNIQUE (station_id, slot_number)
);

CREATE TABLE IF NOT EXISTS tariffs (
    id VARCHAR(36) PRIMARY KEY,
    station_id VARCHAR(36) NOT NULL REFERENCES stations(id) ON DELETE CASCADE,
    price_per_kwh DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
    valid_from TIMESTAMPTZ NOT NULL,
    valid_to TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Trigger untuk update updated_at secara otomatis
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_stations_updated_at BEFORE UPDATE ON stations
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TRIGGER update_charger_slots_updated_at BEFORE UPDATE ON charger_slots
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
