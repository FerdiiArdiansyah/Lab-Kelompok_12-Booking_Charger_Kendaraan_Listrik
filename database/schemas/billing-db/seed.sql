-- billing-db Seed Data

INSERT INTO invoices (id, session_id, user_id, tariff_id, consumed_kwh, price_per_kwh, subtotal, tax, total, status) VALUES
('inv-001', 'ses-001', 'usr-001', 'trf-001', 28.550, 2467.00, 70432.85, 7747.61, 78180.46, 'PAID')
ON CONFLICT (id) DO NOTHING;

INSERT INTO payments (id, invoice_id, payment_method, amount, status, transaction_ref, paid_at) VALUES
('pay-001', 'inv-001', 'QRIS', 78180.46, 'SUCCESS', 'TX-QRIS-9920192', CURRENT_TIMESTAMP - INTERVAL '25 minutes')
ON CONFLICT (id) DO NOTHING;

INSERT INTO audit_logs (entity_name, entity_id, action, old_value, new_value, performed_by) VALUES
('Invoice', 'inv-001', 'UPDATE', '{"status": "UNPAID"}', '{"status": "PAID"}', 'billing-service-worker');
