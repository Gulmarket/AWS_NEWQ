CREATE TYPE order_status AS ENUM ('active', 'expired', 'unavailable');

CREATE TABLE cart (
    id SERIAL PRIMARY KEY,
    pro_user_id INT REFERENCES pro_users(id),
    delivery_id INT REFERENCES delivery(id),
    delivery_address VARCHAR NOT NULL,
    quantity INT NOT NULL,
    total_price FLOAT NOT NULL DEFAULT 0,
    status order_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    time_to_overdue TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '5 minutes')
)
