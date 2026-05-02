INSERT INTO users (id, name, email, password_hash, role)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'Aruzhan Buyer', 'buyer@example.com', '$2a$10$vxUwF308AkkSc4cS7PR46O2fPaOOrikuk.2gIRxZeUH05sZgTL4n2', 'user'),
    ('22222222-2222-2222-2222-222222222222', 'Daniyar Seller', 'seller@example.com', '$2a$10$vxUwF308AkkSc4cS7PR46O2fPaOOrikuk.2gIRxZeUH05sZgTL4n2', 'seller'),
    ('33333333-3333-3333-3333-333333333333', 'Admin User', 'admin@example.com', '$2a$10$vxUwF308AkkSc4cS7PR46O2fPaOOrikuk.2gIRxZeUH05sZgTL4n2', 'admin')
ON CONFLICT (email) DO NOTHING;
