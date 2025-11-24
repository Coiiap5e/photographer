CREATE TABLE shoot_clients
(
    shoot_id INTEGER REFERENCES shoots(id) ON DELETE CASCADE,
    client_id INTEGER REFERENCES clients(id) ON DELETE CASCADE,
    is_main_client BOOLEAN DEFAULT false,
    relationship_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,

    PRIMARY KEY (shoot_id, client_id)
);