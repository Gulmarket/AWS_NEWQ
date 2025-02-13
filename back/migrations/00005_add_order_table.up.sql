CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    pro_user_id INT REFERENCES pro_users(id),
    plantation_id INT REFERENCES plantations(id),
    delivery_id INT REFERENCES delivery(id),
    delivery_address VARCHAR NOT NULL,
    quantity INT NOT NULL,
    price_for_one FLOAT NOT NULL DEFAULT 0,
    total_price FLOAT NOT NULL DEFAULT 0,
    status VARCHAR NOT NULL DEFAULT 'pending'
)
