-- Migration: Add lease tracking and rent adjustments
-- Description: Adiciona rastreamento de renovações (parent_lease_id, generation) e histórico de reajustes

-- Adicionar campos de rastreamento de renovações na tabela leases
ALTER TABLE leases
  ADD COLUMN parent_lease_id UUID REFERENCES leases(id) ON DELETE SET NULL,
  ADD COLUMN generation INTEGER NOT NULL DEFAULT 1;

-- Criar índices para otimizar queries de rastreamento
CREATE INDEX idx_leases_parent_lease_id ON leases(parent_lease_id);
CREATE INDEX idx_leases_generation ON leases(generation);

-- Criar tabela de histórico de reajustes de aluguel
CREATE TABLE IF NOT EXISTS lease_rent_adjustments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Relacionamento
    lease_id UUID NOT NULL REFERENCES leases(id) ON DELETE CASCADE,

    -- Valores do reajuste
    previous_rent_value DECIMAL(10,2) NOT NULL CHECK (previous_rent_value > 0),
    new_rent_value DECIMAL(10,2) NOT NULL CHECK (new_rent_value > 0),
    adjustment_percentage DECIMAL(5,2) NOT NULL, -- pode ser negativo (redução)

    -- Contexto e auditoria
    applied_at TIMESTAMP NOT NULL DEFAULT NOW(),
    reason TEXT,
    applied_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Auditoria
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Índices para otimizar queries de histórico
CREATE INDEX idx_lease_rent_adjustments_lease_id ON lease_rent_adjustments(lease_id);
CREATE INDEX idx_lease_rent_adjustments_applied_at ON lease_rent_adjustments(applied_at);

-- Comentários explicativos
COMMENT ON COLUMN leases.parent_lease_id IS 'ID do contrato anterior (null = contrato original)';
COMMENT ON COLUMN leases.generation IS 'Geração do contrato: 1=original, 2=1ª renovação, 3=2ª renovação, etc.';
COMMENT ON TABLE lease_rent_adjustments IS 'Histórico de reajustes de valor de aluguel aplicados';
COMMENT ON COLUMN lease_rent_adjustments.previous_rent_value IS 'Valor do aluguel antes do reajuste';
COMMENT ON COLUMN lease_rent_adjustments.new_rent_value IS 'Valor do aluguel após o reajuste';
COMMENT ON COLUMN lease_rent_adjustments.adjustment_percentage IS 'Percentual de reajuste calculado ((novo - antigo) / antigo * 100)';
COMMENT ON COLUMN lease_rent_adjustments.reason IS 'Motivo do reajuste (ex: "Reajuste anual IGPM")';
COMMENT ON COLUMN lease_rent_adjustments.applied_by IS 'Usuário que aplicou o reajuste';
