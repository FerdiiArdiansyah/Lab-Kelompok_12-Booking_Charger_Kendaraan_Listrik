-- billing-db Seed Data
-- subtotal = consumed_kwh * price_per_kwh | tax = subtotal * 11% PPN | total = subtotal + tax
-- inv-001: Ioniq 5, 50.6 kWh × Rp2.467 = Rp124.830,20 + PPN Rp13.731,32 = Rp138.561,52
-- inv-002: Nissan Leaf, 26.4 kWh × Rp2.467 = Rp65.128,80 + PPN Rp7.164,17 = Rp72.292,97

INSERT INTO invoices (id, session_id, user_id, tariff_id, consumed_kwh, price_per_kwh, subtotal, tax, total, status) VALUES
('inv-001', 'ses-001', 'usr-001', 'trf-001', 50.600, 2467.00, 124830.20, 13731.32, 138561.52, 'PAID'),
('inv-002', 'ses-003', 'usr-006', 'trf-001', 26.400, 2467.00,  65128.80,  7164.17,  72292.97, 'PAID')
ON CONFLICT (id) DO NOTHING;

INSERT INTO payments (id, invoice_id, payment_method, amount, status, transaction_ref, paid_at) VALUES
('pay-001', 'inv-001', 'QRIS',    138561.52, 'SUCCESS', 'TX-QRIS-20260819-001',  CURRENT_TIMESTAMP - INTERVAL '2 hours 20 minutes'),
('pay-002', 'inv-002', 'VA_BCA',   72292.97, 'SUCCESS', 'TX-VABCA-20260819-002', CURRENT_TIMESTAMP - INTERVAL '4 hours 40 minutes')
ON CONFLICT (id) DO NOTHING;

INSERT INTO audit_logs (entity_name, entity_id, action, old_value, new_value, performed_by) VALUES
('Invoice', 'inv-001', 'UPDATE', '{"status":"UNPAID"}', '{"status":"PAID"}', 'billing-worker'),
('Invoice', 'inv-002', 'UPDATE', '{"status":"UNPAID"}', '{"status":"PAID"}', 'billing-worker');

INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload, status) VALUES
('evt-inv-001', 'Invoice', 'inv-001', 'InvoiceCreated',   '{"invoiceId":"inv-001","userId":"usr-001","consumedKwh":50.6,"total":138561.52}',  'PUBLISHED'),
('evt-pay-001', 'Payment', 'pay-001', 'PaymentCompleted', '{"paymentId":"pay-001","invoiceId":"inv-001","method":"QRIS","amount":138561.52}', 'PUBLISHED'),
('evt-inv-002', 'Invoice', 'inv-002', 'InvoiceCreated',   '{"invoiceId":"inv-002","userId":"usr-006","consumedKwh":26.4,"total":72292.97}',   'PUBLISHED'),
('evt-pay-002', 'Payment', 'pay-002', 'PaymentCompleted', '{"paymentId":"pay-002","invoiceId":"inv-002","method":"VA_BCA","amount":72292.97}', 'PUBLISHED')
ON CONFLICT (id) DO NOTHING;
