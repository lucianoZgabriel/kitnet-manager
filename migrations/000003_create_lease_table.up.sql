-- Migration: Create leases table
-- Description: Tabela de contratos de locação com relacionamento entre unidades e moradores

CREATE TABLE IF NOT EXISTS leases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relacionamentos
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE RESTRICT,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    
    -- Datas do contrato
    contract_signed_date DATE NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    
    -- Configurações de pagamento
    payment_due_day INTEGER NOT NULL CHECK (payment_due_day BETWEEN 1 AND 31),
    
    -- Valores financeiros
    monthly_rent_value DECIMAL(10,2) NOT NULL CHECK (monthly_rent_value > 0),
    painting_fee_total DECIMAL(10,2) NOT NULL DEFAULT 250.00 CHECK (painting_fee_total >= 0),
    painting_fee_installments INTEGER NOT NULL CHECK (painting_fee_installments IN (1, 2, 3, 4)),
    painting_fee_paid DECIMAL(10,2) NOT NULL DEFAULT 0.00 CHECK (painting_fee_paid >= 0),
    
    -- Status do contrato
    status VARCHAR(20) NOT NULL,
    
    -- Auditoria
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints adicionais
    CONSTRAINT chk_dates CHECK (start_date < end_date),
    CONSTRAINT chk_painting_fee_paid CHECK (painting_fee_paid <= painting_fee_total)
);

-- Índices para otimizar queries mais comuns
CREATE INDEX idx_leases_unit_id ON leases(unit_id);
CREATE INDEX idx_leases_tenant_id ON leases(tenant_id);
CREATE INDEX idx_leases_status ON leases(status);
CREATE INDEX idx_leases_end_date ON leases(end_date);

-- Índice composto para buscar contratos ativos por unidade
CREATE INDEX idx_leases_unit_status ON leases(unit_id, status);

-- Índice composto para buscar contratos ativos por morador
CREATE INDEX idx_leases_tenant_status ON leases(tenant_id, status);

-- Comentários explicativos
COMMENT ON TABLE leases IS 'Contratos de locação das kitnets com período de 6 meses';
COMMENT ON COLUMN leases.contract_signed_date IS 'Data em que o contrato foi assinado';
COMMENT ON COLUMN leases.start_date IS 'Data de início da vigência do contrato';
COMMENT ON COLUMN leases.end_date IS 'Data de término do contrato (calculada: start_date + 6 meses)';
COMMENT ON COLUMN leases.payment_due_day IS 'Dia do mês para vencimento do aluguel (1-31)';
COMMENT ON COLUMN leases.monthly_rent_value IS 'Valor mensal do aluguel em reais';
COMMENT ON COLUMN leases.painting_fee_total IS 'Valor total da taxa de pintura (padrão R$ 250,00)';
COMMENT ON COLUMN leases.painting_fee_installments IS 'Número de parcelas da taxa de pintura (1, 2 ou 3)';
COMMENT ON COLUMN leases.painting_fee_paid IS 'Valor já pago da taxa de pintura';
COMMENT ON COLUMN leases.status IS 'Status do contrato: active, expiring_soon, expired, cancelled';