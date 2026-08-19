-- Pattern 3: Audit Logging Trigger Pattern
-- Otomatis mencatat histori mutasi (UPDATE/DELETE) pada data sensitif seperti pembayaran & invoice.

CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    entity_name VARCHAR(100) NOT NULL,
    entity_id VARCHAR(36) NOT NULL,
    action VARCHAR(50) NOT NULL,
    old_value JSONB,
    new_value JSONB,
    performed_by VARCHAR(100) DEFAULT CURRENT_USER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION audit_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO audit_logs (entity_name, entity_id, action, old_value, new_value)
        VALUES (TG_TABLE_NAME, OLD.id::text, 'DELETE', to_jsonb(OLD), NULL);
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO audit_logs (entity_name, entity_id, action, old_value, new_value)
        VALUES (TG_TABLE_NAME, NEW.id::text, 'UPDATE', to_jsonb(OLD), to_jsonb(NEW));
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO audit_logs (entity_name, entity_id, action, old_value, new_value)
        VALUES (TG_TABLE_NAME, NEW.id::text, 'INSERT', NULL, to_jsonb(NEW));
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
