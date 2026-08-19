-- station-db Indexes
-- Optimasi pencarian stasiun & slot

CREATE INDEX IF NOT EXISTS idx_stations_status ON stations(status);
CREATE INDEX IF NOT EXISTS idx_charger_slots_station_status ON charger_slots(station_id, status);
CREATE INDEX IF NOT EXISTS idx_tariffs_station_valid ON tariffs(station_id, valid_from, is_active);
