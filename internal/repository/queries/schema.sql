-- Schema reference for SQLC code generation
-- This file must match the actual database schema from migrations

CREATE TYPE unit_status AS ENUM (
    'available',
    'occupied',
    'maintenance',
    'renovation'
);

CREATE TABLE units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number VARCHAR(10) NOT NULL UNIQUE,
    floor INTEGER NOT NULL CHECK (floor >= 1),
    status unit_status NOT NULL DEFAULT 'available',
    is_renovated BOOLEAN NOT NULL DEFAULT FALSE,
    base_rent_value DECIMAL(10,2) NOT NULL CHECK (base_rent_value >= 0),
    renovated_rent_value DECIMAL(10,2) NOT NULL CHECK (renovated_rent_value >= 0),
    current_rent_value DECIMAL(10,2) NOT NULL CHECK (current_rent_value >= 0),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_units_status ON units(status);
CREATE INDEX idx_units_floor ON units(floor);
CREATE INDEX idx_units_is_renovated ON units(is_renovated);