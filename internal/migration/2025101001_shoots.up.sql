CREATE TABLE IF NOT EXISTS shoots
(
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    shoot_price DECIMAL(10,0),
    location VARCHAR(255),
    shoot_type VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);