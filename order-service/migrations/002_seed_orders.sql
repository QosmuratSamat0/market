INSERT INTO orders (id, user_id, status, total, created_at, updated_at)
VALUES
    (
        'cccccccc-cccc-cccc-cccc-ccccccccc001',
        '11111111-1111-1111-1111-111111111111',
        'paid',
        598.00,
        NOW() - INTERVAL '5 days',
        NOW() - INTERVAL '5 days'
    ),
    (
        'cccccccc-cccc-cccc-cccc-ccccccccc002',
        '11111111-1111-1111-1111-111111111111',
        'shipped',
        297.00,
        NOW() - INTERVAL '2 days',
        NOW() - INTERVAL '1 day'
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO order_items (id, order_id, product_id, seller_id, quantity, price)
VALUES
    (
        'dddddddd-dddd-dddd-dddd-ddddddddd001',
        'cccccccc-cccc-cccc-cccc-ccccccccc001',
        'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb001',
        '22222222-2222-2222-2222-222222222222',
        1,
        349.00
    ),
    (
        'dddddddd-dddd-dddd-dddd-ddddddddd002',
        'cccccccc-cccc-cccc-cccc-ccccccccc001',
        'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb002',
        '22222222-2222-2222-2222-222222222222',
        1,
        249.00
    ),
    (
        'dddddddd-dddd-dddd-dddd-ddddddddd003',
        'cccccccc-cccc-cccc-cccc-ccccccccc002',
        'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb007',
        '22222222-2222-2222-2222-222222222222',
        1,
        69.00
    ),
    (
        'dddddddd-dddd-dddd-dddd-ddddddddd004',
        'cccccccc-cccc-cccc-cccc-ccccccccc002',
        'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb006',
        '22222222-2222-2222-2222-222222222222',
        1,
        129.00
    ),
    (
        'dddddddd-dddd-dddd-dddd-ddddddddd005',
        'cccccccc-cccc-cccc-cccc-ccccccccc002',
        'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb008',
        '22222222-2222-2222-2222-222222222222',
        1,
        99.00
    )
ON CONFLICT (id) DO NOTHING;
