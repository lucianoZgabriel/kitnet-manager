package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository/sqlc"
)

// Compile-time check to ensure UnitRepository implements repository.UnitRepository
var _ repository.UnitRepository = (*UnitRepository)(nil)

// UnitRepository implementa repository.UnitRepository usando PostgreSQL via SQLC
type UnitRepository struct {
	queries *sqlc.Queries
	db      *sql.DB
}

// NewUnitRepository cria uma nova instância do repositório de unidades
func NewUnitRepository(db *sql.DB) *UnitRepository {
	return &UnitRepository{
		queries: sqlc.New(db),
		db:      db,
	}
}

// Create insere uma nova unidade no banco de dados
func (r *UnitRepository) Create(ctx context.Context, unit *domain.Unit) error {
	params := sqlc.CreateUnitParams{
		ID:                 unit.ID,
		Number:             unit.Number,
		Floor:              int32(unit.Floor),
		Status:             sqlc.UnitStatus(unit.Status),
		IsRenovated:        unit.IsRenovated,
		BaseRentValue:      unit.BaseRentValue,
		RenovatedRentValue: unit.RenovatedRentValue,
		CurrentRentValue:   unit.CurrentRentValue,
		Notes:              toNullString(unit.Notes),
		CreatedAt:          unit.CreatedAt,
		UpdatedAt:          unit.UpdatedAt,
	}

	created, err := r.queries.CreateUnit(ctx, params)
	if err != nil {
		return err
	}

	// Atualiza o objeto unit com dados do banco (caso o DB tenha modificado algo)
	unit.CreatedAt = created.CreatedAt
	unit.UpdatedAt = created.UpdatedAt

	return nil
}

// GetById busca uma unidade pelo ID
func (r *UnitRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Unit, error) {
	dbUnit, err := r.queries.GetUnitByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Não encontrou, retorna nil sem erro
		}
		return nil, err
	}

	return r.toDomain(dbUnit), nil
}

// GetByNumber busca uma unidade pelo número
func (r *UnitRepository) GetByNumber(ctx context.Context, number string) (*domain.Unit, error) {
	dbUnit, err := r.queries.GetUnitByNumber(ctx, number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(dbUnit), nil
}

// List retorna todas as unidades
func (r *UnitRepository) List(ctx context.Context) ([]*domain.Unit, error) {
	dbUnits, err := r.queries.ListUnits(ctx)
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(dbUnits), nil
}

// ListByStatus retorna unidades filtradas por status
func (r *UnitRepository) ListByStatus(ctx context.Context, status domain.UnitStatus) ([]*domain.Unit, error) {
	dbUnits, err := r.queries.ListUnitsByStatus(ctx, sqlc.UnitStatus(status))
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(dbUnits), nil
}

// ListByFloor retorna unidades de um andar específico
func (r *UnitRepository) ListByFloor(ctx context.Context, floor int) ([]*domain.Unit, error) {
	dbUnits, err := r.queries.ListUnitsByFloor(ctx, int32(floor))
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(dbUnits), nil
}

// ListAvailable retorna apenas unidades disponíveis
func (r *UnitRepository) ListAvailable(ctx context.Context) ([]*domain.Unit, error) {
	dbUnits, err := r.queries.ListAvailableUnits(ctx)
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(dbUnits), nil
}

// Update atualiza uma unidade existente
func (r *UnitRepository) Update(ctx context.Context, unit *domain.Unit) error {
	unit.UpdatedAt = time.Now()

	params := sqlc.UpdateUnitParams{
		ID:                 unit.ID,
		Number:             unit.Number,
		Floor:              int32(unit.Floor),
		Status:             sqlc.UnitStatus(unit.Status),
		IsRenovated:        unit.IsRenovated,
		BaseRentValue:      unit.BaseRentValue,
		RenovatedRentValue: unit.RenovatedRentValue,
		CurrentRentValue:   unit.CurrentRentValue,
		Notes:              toNullString(unit.Notes),
		UpdatedAt:          unit.UpdatedAt,
	}

	updated, err := r.queries.UpdateUnit(ctx, params)
	if err != nil {
		return err
	}

	unit.UpdatedAt = updated.UpdatedAt
	return nil
}

// UpdateStatus atualiza apenas o status de uma unidade
func (r *UnitRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.UnitStatus) error {
	params := sqlc.UpdateUnitStatusParams{
		ID:        id,
		Status:    sqlc.UnitStatus(status),
		UpdatedAt: time.Now(),
	}

	_, err := r.queries.UpdateUnitStatus(ctx, params)
	return err
}

// Delete remove uma unidade do banco de dados
func (r *UnitRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUnit(ctx, id)
}

// Count retorna o total de unidades
func (r *UnitRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUnits(ctx)
}

// CountByStatus retorna o total de unidades por status
func (r *UnitRepository) CountByStatus(ctx context.Context, status domain.UnitStatus) (int64, error) {
	return r.queries.CountUnitsByStatus(ctx, sqlc.UnitStatus(status))
}

// toDomain converte sqlc.Unit para domain.Unit
func (r *UnitRepository) toDomain(dbUnit sqlc.Unit) *domain.Unit {
	return &domain.Unit{
		ID:                 dbUnit.ID,
		Number:             dbUnit.Number,
		Floor:              int(dbUnit.Floor),
		Status:             domain.UnitStatus(dbUnit.Status),
		IsRenovated:        dbUnit.IsRenovated,
		BaseRentValue:      dbUnit.BaseRentValue,
		RenovatedRentValue: dbUnit.RenovatedRentValue,
		CurrentRentValue:   dbUnit.CurrentRentValue,
		Notes:              dbUnit.Notes.String,
		CreatedAt:          dbUnit.CreatedAt,
		UpdatedAt:          dbUnit.UpdatedAt,
	}
}

// toDomainSlice converte []sqlc.Unit para []*domain.Unit
func (r *UnitRepository) toDomainSlice(dbUnits []sqlc.Unit) []*domain.Unit {
	units := make([]*domain.Unit, len(dbUnits))
	for i, dbUnit := range dbUnits {
		units[i] = r.toDomain(dbUnit)
	}
	return units
}

// toNullString converte string para sql.NullString
func toNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}
