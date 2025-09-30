package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository"
	"github.com/shopspring/decimal"
)

// Service layer errors
var (
	ErrUnitNotFound             = errors.New("unit not found")
	ErrUnitNumberAlreadyExists  = errors.New("unit number already exists")
	ErrCannotDeleteOccupiedUnit = errors.New("cannot delete occupied unit")
	ErrInvalidStatusTransition  = errors.New("invalid status transition")
)

// UnitService contém a lógica de negócio para gestão de unidades
type UnitService struct {
	unitRepo repository.UnitRepository
}

// NewUnitService cria uma nova instância do serviço de unidades
func NewUnitService(unitRepo repository.UnitRepository) *UnitService {
	return &UnitService{
		unitRepo: unitRepo,
	}
}

// CreateUnit cria uma nova unidade com validações de negócio
func (s *UnitService) CreateUnit(ctx context.Context, number string, floor int, baseRentValue, renovatedRentValue decimal.Decimal) (*domain.Unit, error) {
	// Verifica se número já existe
	existing, err := s.unitRepo.GetByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("error checking unit number: %w", err)
	}
	if existing != nil {
		return nil, ErrUnitNumberAlreadyExists
	}

	// Cria nova unidade usando o domain model
	unit, err := domain.NewUnit(number, floor, baseRentValue, renovatedRentValue)
	if err != nil {
		return nil, fmt.Errorf("error creating unit: %w", err)
	}

	// Persistir no banco
	if err := s.unitRepo.Create(ctx, unit); err != nil {
		return nil, fmt.Errorf("error saving unit: %w", err)
	}

	return unit, nil
}

// GetUnitByID busca uma unidade pelo ID
func (s *UnitService) GetUnitByID(ctx context.Context, id uuid.UUID) (*domain.Unit, error) {
	unit, err := s.unitRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting unit: %w", err)
	}
	if unit == nil {
		return nil, ErrUnitNotFound
	}

	return unit, nil
}

// GetUnitByNumber busca uma unidade pelo número
func (s *UnitService) GetUnitByNumber(ctx context.Context, number string) (*domain.Unit, error) {
	unit, err := s.unitRepo.GetByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("error getting unit: %w", err)
	}
	if unit == nil {
		return nil, ErrUnitNotFound
	}
	return unit, nil
}

// ListUnits retorna todas as unidades
func (s *UnitService) ListUnits(ctx context.Context) ([]*domain.Unit, error) {
	units, err := s.unitRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing units: %w", err)
	}
	return units, nil
}

// ListUnitsByStatus retorna unidades filtradas por status
func (s *UnitService) ListUnitsByStatus(ctx context.Context, status domain.UnitStatus) ([]*domain.Unit, error) {
	units, err := s.unitRepo.ListByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("error listing units by status: %w", err)
	}
	return units, nil
}

// ListUnitsByFloor retorna unidades de um andar específico
func (s *UnitService) ListUnitsByFloor(ctx context.Context, floor int) ([]*domain.Unit, error) {
	units, err := s.unitRepo.ListByFloor(ctx, floor)
	if err != nil {
		return nil, fmt.Errorf("error listing units by floor: %w", err)
	}
	return units, nil
}

// ListAvailableUnits retorna apenas unidades disponíveis
func (s *UnitService) ListAvailableUnits(ctx context.Context) ([]*domain.Unit, error) {
	units, err := s.unitRepo.ListAvailable(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing available units: %w", err)
	}
	return units, nil
}

// UpdateUnit atualiza uma unidade existente
func (s *UnitService) UpdateUnit(ctx context.Context, id uuid.UUID, number string, floor int, isRenovated bool, baseRentValue, renovatedRentValue decimal.Decimal, notes string) (*domain.Unit, error) {
	// Busca unidade
	unit, err := s.GetUnitByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validar se o novo número já existe (se mudou)
	if unit.Number != number {
		existing, err := s.unitRepo.GetByNumber(ctx, number)
		if err != nil {
			return nil, fmt.Errorf("error checking unit number: %w", err)
		}
		if existing != nil && existing.ID != id {
			return nil, ErrUnitNumberAlreadyExists
		}
	}

	// Atualizar campos
	unit.Number = number
	unit.Floor = floor
	unit.IsRenovated = isRenovated
	unit.BaseRentValue = baseRentValue
	unit.RenovatedRentValue = renovatedRentValue
	unit.Notes = notes

	// Recalcular valor atual baseado no status de renovação
	unit.CalculateCurrentRentValue()

	// Validar unidade atualizada
	if err := unit.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Persistir mudanças
	if err := s.unitRepo.Update(ctx, unit); err != nil {
		return nil, fmt.Errorf("error updating unit: %w", err)
	}

	return unit, nil
}

// UpdateUnitStatus atualiza apenas o status de uma unidade
func (s *UnitService) UpdateUnitStatus(ctx context.Context, id uuid.UUID, newStatus domain.UnitStatus) error {
	// Buscar unidade existente
	unit, err := s.GetUnitByID(ctx, id)
	if err != nil {
		return err
	}

	// Validar transição de status usando método do domain
	if err := unit.ChangeStatus(newStatus); err != nil {
		return fmt.Errorf("invalid status change: %w", err)
	}

	// Persistir mudança
	if err := s.unitRepo.UpdateStatus(ctx, id, newStatus); err != nil {
		return fmt.Errorf("error updating unit status: %w", err)
	}

	return nil
}

// MarkUnitAsRenovated marca uma unidade como reformada
func (s *UnitService) MarkUnitAsRenovated(ctx context.Context, id uuid.UUID) error {
	// Buscar unidade
	unit, err := s.GetUnitByID(ctx, id)
	if err != nil {
		return err
	}

	// Usar método do domain
	unit.MarkAsRenovated()

	// Persistir
	if err := s.unitRepo.Update(ctx, unit); err != nil {
		return fmt.Errorf("error marking unit as renovated: %w", err)
	}

	return nil
}

// DeleteUnit remove uma unidade
func (s *UnitService) DeleteUnit(ctx context.Context, id uuid.UUID) error {
	// Buscar unidade
	unit, err := s.GetUnitByID(ctx, id)
	if err != nil {
		return err
	}

	// Regra de negócio: não pode deletar unidade ocupada
	if unit.IsOccupied() {
		return ErrCannotDeleteOccupiedUnit
	}

	// Deletar
	if err := s.unitRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting unit: %w", err)
	}

	return nil
}

// GetOccupancyStats retorna estatísticas de ocupação
func (s *UnitService) GetOccupancyStats(ctx context.Context) (*OccupancyStats, error) {
	total, err := s.unitRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("error counting units: %w", err)
	}

	occupied, err := s.unitRepo.CountByStatus(ctx, domain.UnitStatusOccupied)
	if err != nil {
		return nil, fmt.Errorf("error counting occupied units: %w", err)
	}

	available, err := s.unitRepo.CountByStatus(ctx, domain.UnitStatusAvailable)
	if err != nil {
		return nil, fmt.Errorf("error counting available units: %w", err)
	}

	maintenance, err := s.unitRepo.CountByStatus(ctx, domain.UnitStatusMaintenance)
	if err != nil {
		return nil, fmt.Errorf("error counting maintenance units: %w", err)
	}

	renovation, err := s.unitRepo.CountByStatus(ctx, domain.UnitStatusRenovation)
	if err != nil {
		return nil, fmt.Errorf("error counting renovation units: %w", err)
	}

	occupancyRate := 0.0
	if total > 0 {
		occupancyRate = (float64(occupied) / float64(total)) * 100
	}

	return &OccupancyStats{
		Total:         total,
		Occupied:      occupied,
		Available:     available,
		Maintenance:   maintenance,
		Renovation:    renovation,
		OccupancyRate: occupancyRate,
	}, nil
}

// OccupancyStats representa estatísticas de ocupação
type OccupancyStats struct {
	Total         int64   `json:"total"`
	Occupied      int64   `json:"occupied"`
	Available     int64   `json:"available"`
	Maintenance   int64   `json:"maintenance"`
	Renovation    int64   `json:"renovation"`
	OccupancyRate float64 `json:"occupancy_rate"` // Percentual
}
