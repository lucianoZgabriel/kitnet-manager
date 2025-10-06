package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/shopspring/decimal"
)

// UnitRepository define o contrato para operações de persistência de Units
type UnitRepository interface {
	// Create cria uma nova unidade no banco de dados
	Create(ctx context.Context, unit *domain.Unit) error

	// GetByID busca uma unidade pelo ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Unit, error)

	// GetByNumber busca uma unidade pelo número
	GetByNumber(ctx context.Context, number string) (*domain.Unit, error)

	// List retorna todas as unidades ordenadas por andar e número
	List(ctx context.Context) ([]*domain.Unit, error)

	// ListByStatus retorna unidades filtradas por status
	ListByStatus(ctx context.Context, status domain.UnitStatus) ([]*domain.Unit, error)

	// ListByFloor retorna unidades de um andar específico
	ListByFloor(ctx context.Context, floor int) ([]*domain.Unit, error)

	// ListAvailable retorna apenas unidades disponíveis para locação
	ListAvailable(ctx context.Context) ([]*domain.Unit, error)

	// Update atualiza uma unidade existente
	Update(ctx context.Context, unit *domain.Unit) error

	// UpdateStatus atualiza apenas o status de uma unidade
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.UnitStatus) error

	// Delete remove uma unidade do banco de dados
	Delete(ctx context.Context, id uuid.UUID) error

	// Count retorna o total de unidades
	Count(ctx context.Context) (int64, error)

	// CountByStatus retorna o total de unidades por status
	CountByStatus(ctx context.Context, status domain.UnitStatus) (int64, error)
}

// TenantRepository define o contrato para operações de persistência de Tenants
type TenantRepository interface {
	// Create cria um novo morador no banco de dados
	Create(ctx context.Context, tenant *domain.Tenant) error

	// GetByID busca um morador pelo ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)

	// GetByCPF busca um morador pelo CPF
	GetByCPF(ctx context.Context, cpf string) (*domain.Tenant, error)

	// List retorna todos os moradores ordenados por nome
	List(ctx context.Context) ([]*domain.Tenant, error)

	// SearchByName busca moradores por nome (case-insensitive)
	SearchByName(ctx context.Context, name string) ([]*domain.Tenant, error)

	// Update atualiza um morador existente
	Update(ctx context.Context, tenant *domain.Tenant) error

	// Delete remove um morador do banco de dados
	Delete(ctx context.Context, id uuid.UUID) error

	// Count retorna o total de moradores
	Count(ctx context.Context) (int64, error)

	// ExistsByCPF verifica se já existe um morador com o CPF
	ExistsByCPF(ctx context.Context, cpf string) (bool, error)
}

// LeaseRepository define as operações de persistência para Lease
type LeaseRepository interface {
	Create(ctx context.Context, lease *domain.Lease) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Lease, error)
	List(ctx context.Context) ([]*domain.Lease, error)
	ListByStatus(ctx context.Context, status domain.LeaseStatus) ([]*domain.Lease, error)
	ListByUnitID(ctx context.Context, unitID uuid.UUID) ([]*domain.Lease, error)
	ListByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*domain.Lease, error)
	GetActiveByUnitID(ctx context.Context, unitID uuid.UUID) (*domain.Lease, error)
	GetActiveByTenantID(ctx context.Context, tenantID uuid.UUID) (*domain.Lease, error)
	GetExpiringSoon(ctx context.Context) ([]*domain.Lease, error)
	Update(ctx context.Context, lease *domain.Lease) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.LeaseStatus) error
	UpdatePaintingFeePaid(ctx context.Context, id uuid.UUID, paintingFeePaid decimal.Decimal) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status domain.LeaseStatus) (int64, error)
}
