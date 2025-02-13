CREATE TABLE shop (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    addresses jsonb NOT NULL,
    shop_logo VARCHAR DEFAULT NULL,
    work_schedule jsonb NOT NULL,
    pro_user_id INT REFERENCES pro_users(id)
)
