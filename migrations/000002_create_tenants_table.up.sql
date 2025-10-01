-- Create tenants table
CREATE TABLE tenants (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name            VARCHAR(255) NOT NULL,
    cpf                  VARCHAR(14) NOT NULL UNIQUE,
    phone                VARCHAR(20) NOT NULL,
    email                VARCHAR(255),
    id_document_type     VARCHAR(10),
    id_document_number   VARCHAR(50),
    created_at           TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT check_full_name_not_empty CHECK (length(trim(full_name)) > 0),
    CONSTRAINT check_cpf_format CHECK (cpf ~ '^\d{3}\.\d{3}\.\d{3}-\d{2}$'),
    CONSTRAINT check_phone_not_empty CHECK (length(trim(phone)) > 0)
);

-- Create indexes for common queries
CREATE INDEX idx_tenants_cpf ON tenants(cpf);
CREATE INDEX idx_tenants_full_name ON tenants(full_name);
CREATE INDEX idx_tenants_created_at ON tenants(created_at);

-- Create trigger to auto-update updated_at (reutiliza a função já criada na migration de units)
CREATE TRIGGER update_tenants_updated_at 
    BEFORE UPDATE ON tenants
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE tenants IS 'Tenant/resident information';
COMMENT ON COLUMN tenants.full_name IS 'Complete legal name of the tenant';
COMMENT ON COLUMN tenants.cpf IS 'Brazilian CPF document (formatted: XXX.XXX.XXX-XX)';
COMMENT ON COLUMN tenants.phone IS 'Primary contact phone number';
COMMENT ON COLUMN tenants.email IS 'Email address (optional)';
COMMENT ON COLUMN tenants.id_document_type IS 'Type of ID document (RG, CNH, etc.)';
COMMENT ON COLUMN tenants.id_document_number IS 'ID document number';
