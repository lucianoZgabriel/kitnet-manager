-- Create ENUM type for unit status
CREATE TYPE unit_status AS ENUM (
    'available',
    'occupied',
    'maintenance',
    'renovation'
);

-- Create units table
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

-- Create indexes for common queries
CREATE INDEX idx_units_status ON units(status);
CREATE INDEX idx_units_floor ON units(floor);
CREATE INDEX idx_units_is_renovated ON units(is_renovated);

-- Create trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER update_units_updated_at 
    BEFORE UPDATE ON units
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE units IS 'Physical apartment units in the building';
COMMENT ON COLUMN units.number IS 'Unit number (e.g., 101, 201)';
COMMENT ON COLUMN units.status IS 'Current availability status';
COMMENT ON COLUMN units.current_rent_value IS 'Active rent value based on renovation status';