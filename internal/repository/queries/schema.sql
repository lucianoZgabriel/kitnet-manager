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

-- Tenants table
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255),
    id_document_type VARCHAR(10),
    id_document_number VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tenants_cpf ON tenants(cpf);
CREATE INDEX idx_tenants_full_name ON tenants(full_name);
CREATE INDEX idx_tenants_created_at ON tenants(created_at);

-- Tabela leases 
CREATE TABLE leases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE RESTRICT,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    contract_signed_date DATE NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    payment_due_day INTEGER NOT NULL CHECK (payment_due_day BETWEEN 1 AND 31),
    monthly_rent_value DECIMAL(10,2) NOT NULL CHECK (monthly_rent_value > 0),
    painting_fee_total DECIMAL(10,2) NOT NULL DEFAULT 250.00 CHECK (painting_fee_total >= 0),
    painting_fee_installments INTEGER NOT NULL CHECK (painting_fee_installments IN (1, 2, 3, 4)),
    painting_fee_paid DECIMAL(10,2) NOT NULL DEFAULT 0.00 CHECK (painting_fee_paid >= 0),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_dates CHECK (start_date < end_date),
    CONSTRAINT chk_painting_fee_paid CHECK (painting_fee_paid <= painting_fee_total)
);

CREATE INDEX idx_leases_unit_id ON leases(unit_id);
CREATE INDEX idx_leases_tenant_id ON leases(tenant_id);
CREATE INDEX idx_leases_status ON leases(status);
CREATE INDEX idx_leases_end_date ON leases(end_date);
CREATE INDEX idx_leases_unit_status ON leases(unit_id, status);
CREATE INDEX idx_leases_tenant_status ON leases(tenant_id, status);

-- Payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lease_id UUID NOT NULL REFERENCES leases(id) ON DELETE RESTRICT,
    payment_type VARCHAR(20) NOT NULL CHECK (payment_type IN ('rent', 'painting_fee', 'adjustment')),
    reference_month DATE NOT NULL,
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'paid', 'overdue', 'cancelled')),
    due_date DATE NOT NULL,
    payment_date DATE,
    payment_method VARCHAR(20),
    proof_url TEXT,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_lease_id ON payments(lease_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_due_date ON payments(due_date);
CREATE INDEX idx_payments_payment_type ON payments(payment_type);
CREATE INDEX idx_payments_status_due_date ON payments(status, due_date);
CREATE INDEX idx_payments_lease_status ON payments(lease_id, status);

-- User roles enum
CREATE TYPE user_role AS ENUM (
    'admin',
    'manager',
    'viewer'
);

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'viewer',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);
