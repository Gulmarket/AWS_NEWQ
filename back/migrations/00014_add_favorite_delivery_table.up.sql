CREATE TABLE favorite_delivery (
    id SERIAL PRIMARY KEY,
    pro_user_id INT REFERENCES pro_users(id),
    delivery_id INT REFERENCES delivery(id)
);
