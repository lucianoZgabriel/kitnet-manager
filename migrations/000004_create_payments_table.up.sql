-- Migration: Create payments table
-- Description: Tabela de pagamentos relacionados aos contratos de locação

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relacionamento
    lease_id UUID NOT NULL REFERENCES leases(id) ON DELETE RESTRICT,
    
    -- Informações do pagamento
    payment_type VARCHAR(20) NOT NULL CHECK (payment_type IN ('rent', 'painting_fee', 'adjustment')),
    reference_month DATE NOT NULL, -- mês/ano de referência
    
    -- Valores e status
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'paid', 'overdue', 'cancelled')),
    
    -- Datas
    due_date DATE NOT NULL,
    payment_date DATE,
    
    -- Método e comprovante
    payment_method VARCHAR(20),
    proof_url TEXT,
    
    -- Observações
    notes TEXT,
    
    -- Auditoria
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Índices para otimizar queries mais comuns
CREATE INDEX idx_payments_lease_id ON payments(lease_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_due_date ON payments(due_date);
CREATE INDEX idx_payments_payment_type ON payments(payment_type);

-- Índice composto para buscar pagamentos atrasados
CREATE INDEX idx_payments_status_due_date ON payments(status, due_date);

-- Índice composto para buscar pagamentos de um contrato por status
CREATE INDEX idx_payments_lease_status ON payments(lease_id, status);

-- Comentários explicativos
COMMENT ON TABLE payments IS 'Pagamentos de aluguéis e taxas relacionados aos contratos';
COMMENT ON COLUMN payments.payment_type IS 'Tipo: rent (aluguel), painting_fee (taxa pintura), adjustment (ajuste)';
COMMENT ON COLUMN payments.reference_month IS 'Mês/ano de referência do pagamento (ex: 2024-03-01 para março/2024)';
COMMENT ON COLUMN payments.amount IS 'Valor do pagamento em reais';
COMMENT ON COLUMN payments.status IS 'Status: pending, paid, overdue, cancelled';
COMMENT ON COLUMN payments.due_date IS 'Data de vencimento do pagamento';
COMMENT ON COLUMN payments.payment_date IS 'Data efetiva do pagamento (NULL se não pago)';
COMMENT ON COLUMN payments.payment_method IS 'Método usado: pix, cash, bank_transfer, credit_card';
COMMENT ON COLUMN payments.proof_url IS 'URL do comprovante de pagamento (futuro)';