BEGIN;

TRUNCATE TABLE
    payment_refunds,
    payment_webhooks,
    payments
RESTART IDENTITY CASCADE;

INSERT INTO payments (
    id, order_id, payment_method, transaction_id, gateway_id, amount, currency, status,
    gateway_status, gateway_response, failure_reason, processed_at, metadata, created_at, updated_at
) VALUES
    (
        1, 1, 'card', 'pi_demo_20260328_0001', 'stripe', 143.00, 'USD', 'completed',
        'succeeded', '{"charge_id":"ch_demo_0001","receipt_url":"https://example.com/receipts/0001"}', '',
        '2026-03-28T10:05:00Z', '{"brand":"visa","last4":"4242","order_number":"GBX-20260328-0001"}',
        '2026-03-28T10:04:00Z', '2026-03-28T10:05:00Z'
    ),
    (
        2, 2, 'card', 'pi_demo_20260329_0002', 'stripe', 138.99, 'USD', 'completed',
        'succeeded', '{"charge_id":"ch_demo_0002","receipt_url":"https://example.com/receipts/0002"}', '',
        '2026-03-29T11:20:00Z', '{"brand":"mastercard","last4":"4444","order_number":"GBX-20260329-0002"}',
        '2026-03-29T11:18:00Z', '2026-03-29T11:20:00Z'
    ),
    (
        3, 3, 'card', 'pi_demo_20260330_0003', 'stripe', 129.00, 'USD', 'failed',
        'requires_payment_method', '{"code":"card_declined","message":"Your card was declined."}', 'Card was declined by issuer.',
        NULL, '{"brand":"visa","last4":"0341","order_number":"GBX-20260330-0003"}',
        '2026-03-30T13:41:00Z', '2026-03-30T15:10:00Z'
    ),
    (
        4, 4, 'card', 'pi_demo_20260331_0004', 'stripe', 83.99, 'USD', 'refunded',
        'refunded', '{"charge_id":"ch_demo_0004","refund_id":"re_demo_0001"}', '',
        '2026-03-31T08:15:00Z', '{"brand":"amex","last4":"8431","order_number":"GBX-20260331-0004"}',
        '2026-03-31T08:14:00Z', '2026-03-31T11:00:00Z'
    );

INSERT INTO payment_refunds (
    id, payment_id, refund_id, amount, reason, status, gateway_response, processed_at, created_by, created_at
) VALUES
    (
        1, 4, 're_demo_0001', 83.99, 'Duplicate order reported by customer.', 'processed',
        '{"status":"succeeded","balance_transaction":"txn_demo_refund_01"}',
        '2026-03-31T11:00:00Z', 1, '2026-03-31T10:45:00Z'
    );

INSERT INTO payment_webhooks (
    id, gateway_id, event_type, event_id, payload, signature, is_processed, processed_at, created_at
) VALUES
    (
        1, 'stripe', 'payment_intent.succeeded', 'evt_demo_0001',
        '{"id":"evt_demo_0001","type":"payment_intent.succeeded","data":{"object":{"id":"pi_demo_20260328_0001"}}}',
        'whsec_demo_signature_0001', TRUE, '2026-03-28T10:05:30Z', '2026-03-28T10:05:20Z'
    ),
    (
        2, 'stripe', 'payment_intent.payment_failed', 'evt_demo_0002',
        '{"id":"evt_demo_0002","type":"payment_intent.payment_failed","data":{"object":{"id":"pi_demo_20260330_0003"}}}',
        'whsec_demo_signature_0002', TRUE, '2026-03-30T15:10:15Z', '2026-03-30T15:10:10Z'
    ),
    (
        3, 'stripe', 'charge.refunded', 'evt_demo_0003',
        '{"id":"evt_demo_0003","type":"charge.refunded","data":{"object":{"payment_intent":"pi_demo_20260331_0004"}}}',
        'whsec_demo_signature_0003', TRUE, '2026-03-31T11:00:20Z', '2026-03-31T11:00:10Z'
    );

SELECT setval('payments_id_seq', COALESCE((SELECT MAX(id) FROM payments), 1), TRUE);
SELECT setval('payment_refunds_id_seq', COALESCE((SELECT MAX(id) FROM payment_refunds), 1), TRUE);
SELECT setval('payment_webhooks_id_seq', COALESCE((SELECT MAX(id) FROM payment_webhooks), 1), TRUE);

COMMIT;
