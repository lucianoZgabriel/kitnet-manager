package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
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
