-- billing-db Seed Data
-- subtotal = consumed_kwh * price_per_kwh | tax = subtotal * 11% PPN | total = subtotal + tax

INSERT INTO invoices (id, session_id, user_id, tariff_id, consumed_kwh, price_per_kwh, subtotal, tax, total, status) VALUES
('inv-001', 'ses-001', 'usr-001', 'trf-001', 24.500, 2467.00, 60441.50, 6648.57, 67090.07, 'PAID'),
('inv-002', 'ses-003', 'usr-006', 'trf-001', 18.750, 2467.00, 46256.25, 5088.19, 51344.44, 'PAID')
ON CONFLICT (id) DO NOTHING;

INSERT INTO payments (id, invoice_id, payment_method, amount, status, transaction_ref, paid_at) VALUES
('pay-001', 'inv-001', 'QRIS',   67090.07, 'SUCCESS', 'TX-QRIS-20260819-001',  CURRENT_TIMESTAMP - INTERVAL '90 minutes'),
('pay-002', 'inv-002', 'VA_BCA', 51344.44, 'SUCCESS', 'TX-VABCA-20260819-002', CURRENT_TIMESTAMP - INTERVAL '210 minutes')
ON CONFLICT (id) DO NOTHING;

INSERT INTO audit_logs (entity_name, entity_id, action, old_value, new_value, performed_by) VALUES
('Invoice', 'inv-001', 'UPDATE', '{"status":"UNPAID"}', '{"status":"PAID"}', 'billing-worker'),
('Invoice', 'inv-002', 'UPDATE', '{"status":"UNPAID"}', '{"status":"PAID"}', 'billing-worker');

INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload, status) VALUES
('evt-inv-001', 'Invoice', 'inv-001', 'InvoiceCreated',   '{"invoiceId":"inv-001","userId":"usr-001","total":67090.07}',  'PUBLISHED'),
('evt-pay-001', 'Payment', 'pay-001', 'PaymentCompleted', '{"paymentId":"pay-001","invoiceId":"inv-001","amount":67090.07}', 'PUBLISHED'),
('evt-inv-002', 'Invoice', 'inv-002', 'InvoiceCreated',   '{"invoiceId":"inv-002","userId":"usr-006","total":51344.44}',  'PUBLISHED'),
('evt-pay-002', 'Payment', 'pay-002', 'PaymentCompleted', '{"paymentId":"pay-002","invoiceId":"inv-002","amount":51344.44}', 'PUBLISHED')
ON CONFLICT (id) DO NOTHING;
