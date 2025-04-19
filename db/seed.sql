-- seed.sql
-- Insert sample data for Protein Shop database

-- Insert sample products
INSERT INTO products (id, name, description, price, stock, created_at) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'Whey Protein', 'Chocolate, 2lb', 29.99, 100, CURRENT_TIMESTAMP),
    ('550e8400-e29b-41d4-a716-446655440002', 'BCAA Supplement', 'Lemon-lime, 300g', 19.99, 50, CURRENT_TIMESTAMP),
    ('550e8400-e29b-41d4-a716-446655440003', 'Protein Bar', 'Peanut butter', 2.99, 200, CURRENT_TIMESTAMP);

-- Insert sample user
INSERT INTO users (id, email, password_hash, created_at, updated_at) VALUES
    ('550e8400-e29b-41d4-a716-446655440004', 'test@example.com', '$argon2id$v=19$m=65536,t=100,p=4$Q2hhbmdlVGhpc1NhbHQ$...', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
    -- Note: Replace password_hash with a real Argon2 hash for 'StrongPass123!' (generated via backend).

-- Insert sample order
INSERT INTO orders (id, user_id, status, total, created_at) VALUES
    ('550e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440004', 'delivered', 59.98, CURRENT_TIMESTAMP - INTERVAL '2 days');

-- Insert sample order items
INSERT INTO order_items (id, order_id, product_id, quantity, price_at_purchase) VALUES
    ('550e8400-e29b-41d4-a716-446655440006', '550e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440001', 2, 29.99);

-- Insert sample loyalty points
INSERT INTO loyalty_points (id, user_id, points, source, created_at) VALUES
    ('550e8400-e29b-41d4-a716-446655440007', '550e8400-e29b-41d4-a716-446655440004', 5, 'purchase', CURRENT_TIMESTAMP - INTERVAL '2 days');