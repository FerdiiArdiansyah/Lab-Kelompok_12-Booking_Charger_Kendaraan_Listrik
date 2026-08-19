-- user-db Indexes

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_user_vehicles_user ON user_vehicles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_vehicles_license_plate ON user_vehicles(license_plate);
