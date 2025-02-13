CREATE TABLE plantations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL DEFAULT '',
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    country VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    logo_url VARCHAR DEFAULT NULL,
    longitude FLOAT NOT NULL,
    latitude FLOAT NOT NULL,
    work_schedule jsonb NOT NULL,
    phone VARCHAR(18) NOT NULL
)
