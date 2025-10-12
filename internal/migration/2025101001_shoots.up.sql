CREATE TABLE shoots
(
    id SERIAL PRIMARY KEY,
    client_id INTEGER NOT NULL REFERENCES public.clients(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    start_time TIME,
    end_time TIME,
    shoot_price DECIMAL(10,0),
    location VARCHAR(255),
    client_first_name VARCHAR(255),
    client_last_name VARCHAR(255),
    shoot_type VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);