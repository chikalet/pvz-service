CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('client', 'moderator', 'employee')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE pickup_points (
    id SERIAL PRIMARY KEY,
    city VARCHAR(100) NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE intakes (
    id SERIAL PRIMARY KEY,
    pvz_id INTEGER NOT NULL REFERENCES pickup_points(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('in_progress', 'closed')) DEFAULT 'in_progress',
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    closed_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT valid_closing CHECK (
        (status = 'in_progress' AND closed_at IS NULL) OR
        (status = 'closed' AND closed_at IS NOT NULL)
    )
);
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    intake_id INTEGER NOT NULL REFERENCES intakes(id),
    type VARCHAR(50) NOT NULL CHECK (type IN ('electronics', 'clothing', 'shoes')),
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE intake_items (
    id SERIAL PRIMARY KEY,
    intake_id INTEGER NOT NULL REFERENCES intakes(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    received_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_intakes_pvz_id ON intakes(pvz_id);
CREATE INDEX idx_intakes_status ON intakes(status);
CREATE INDEX idx_products_intake_id ON products(intake_id);