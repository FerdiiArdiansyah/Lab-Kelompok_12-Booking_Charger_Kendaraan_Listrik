-- 01-init-databases.sql
-- Script inisialisasi database per service, roles, & ekstensi PostgreSQL

CREATE DATABASE station_db;
CREATE DATABASE booking_db;
CREATE DATABASE session_db;
CREATE DATABASE billing_db;

-- Koneksi ke booking_db untuk install btree_gist extension (Wajib untuk Range Exclusion Constraint)
\c booking_db;
CREATE EXTENSION IF NOT EXISTS btree_gist;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Koneksi ke station_db
\c station_db;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Koneksi ke session_db
\c session_db;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Koneksi ke billing_db
\c billing_db;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
