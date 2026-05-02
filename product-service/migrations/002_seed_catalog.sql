INSERT INTO categories (id, name)
VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'Audio'),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2', 'Wearables'),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3', 'Smart Home'),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa4', 'Accessories')
ON CONFLICT (name) DO NOTHING;

INSERT INTO products (id, name, description, price, category_id, seller_id, image_url, stock)
VALUES
    (
        'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb001',
        'Sony WH-1000XM5 Headphones',
        'Wireless noise cancelling headphones with up to 30 hours of battery life.',
        349.00,
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1',
        '22222222-2222-2222-2222-222222222222',
        'https://images.unsplash.com/photo-1618366712010-f4ae9c647dcb?auto=format&fit=crop&w=900&q=80',
        15
    ),
    (
        'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb002',
        'Apple AirPods Pro 2',
        'In-ear wireless earbuds with active noise cancellation and transparency mode.',
        249.00,
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1',
        '22222222-2222-2222-2222-222222222222',
        'https://images.unsplash.com/photo-1600294037681-c80b4cb5b434?auto=format&fit=crop&w=900&q=80',
        22
    ),
    (
        'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb003',
        'Apple Watch Series 9',
        'Smart watch with fitness tracking, notifications, and all-day battery life.',
        399.00,
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2',
        '22222222-2222-2222-2222-222222222222',
        'https://images.unsplash.com/photo-1434493789847-2f02dc6ca35d?auto=format&fit=crop&w=900&q=80',
        10
    ),
    (
        'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb004',
        'Samsung Galaxy Watch 6',
        'Health-focused smartwatch with sleep insights and advanced activity tracking.',
        299.00,
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2',
        '22222222-2222-2222-2222-222222222222',
        'https://images.unsplash.com/photo-1508685096489-7aacd43bd3b1?auto=format&fit=crop&w=900&q=80',
        12
    ),
    (
        'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb005',
        'Philips Hue Starter Kit',
        'Smart lighting kit with bridge and color bulbs for voice-controlled rooms.',
        179.00,
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3',
        '22222222-2222-2222-2222-222222222222',
        'https://images.unsplash.com/photo-1558002038-1055907df827?auto=format&fit=crop&w=900&q=80',
        18
    ),
    (
        'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb006',
        'Google Nest Thermostat',
        'Energy-saving smart thermostat with scheduling and remote app control.',
        129.00,
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3',
        '22222222-2222-2222-2222-222222222222',
        'https://images.unsplash.com/photo-1558089687-f282ffcbc0d4?auto=format&fit=crop&w=900&q=80',
        8
    ),
    (
        'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb007',
        'Anker USB-C Hub',
        'Compact 7-in-1 USB-C hub with HDMI, card reader, and pass-through charging.',
        69.00,
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa4',
        '22222222-2222-2222-2222-222222222222',
        'https://images.unsplash.com/photo-1625842268584-8f3296236761?auto=format&fit=crop&w=900&q=80',
        30
    ),
    (
        'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb008',
        'Logitech MX Master 3S',
        'Ergonomic wireless mouse with quiet clicks and fast electromagnetic scrolling.',
        99.00,
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa4',
        '22222222-2222-2222-2222-222222222222',
        'https://images.unsplash.com/photo-1615663245857-ac93bb7c39e7?auto=format&fit=crop&w=900&q=80',
        25
    )
ON CONFLICT (id) DO NOTHING;
