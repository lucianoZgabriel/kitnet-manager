package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository/sqlc"
	"github.com/shopspring/decimal"
)

// Compile-time check to ensure LeaseRentAdjustmentRepository implements repository.LeaseRentAdjustmentRepository
var _ repository.LeaseRentAdjustmentRepository = (*LeaseRentAdjustmentRepository)(nil)

// LeaseRentAdjustmentRepository implementa repository.LeaseRentAdjustmentRepository usando PostgreSQL via SQLC
type LeaseRentAdjustmentRepository struct {
	queries *sqlc.Queries
	db      *sql.DB
}

// NewLeaseRentAdjustmentRepository cria uma nova instância do repositório de reajustes
func NewLeaseRentAdjustmentRepository(db *sql.DB) *LeaseRentAdjustmentRepository {
	return &LeaseRentAdjustmentRepository{
		queries: sqlc.New(db),
		db:      db,
	}
}

// Create insere um novo reajuste no banco de dados
func (r *LeaseRentAdjustmentRepository) Create(ctx context.Context, adjustment *domain.LeaseRentAdjustment) error {
	params := sqlc.CreateLeaseRentAdjustmentParams{
		ID:                   adjustment.ID,
		LeaseID:              adjustment.LeaseID,
		PreviousRentValue:    adjustment.PreviousRentValue.String(),
		NewRentValue:         adjustment.NewRentValue.String(),
		AdjustmentPercentage: adjustment.AdjustmentPercentage.String(),
		AppliedAt:            adjustment.AppliedAt,
		Reason:               toNullStringPtr(adjustment.Reason),
		AppliedBy:            toNullUUIDPtr(adjustment.AppliedBy),
		CreatedAt:            adjustment.CreatedAt,
	}

	_, err := r.queries.CreateLeaseRentAdjustment(ctx, params)
	return err
}

// GetByID busca um reajuste pelo ID
func (r *LeaseRentAdjustmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.LeaseRentAdjustment, error) {
	dbAdj, err := r.queries.GetLeaseRentAdjustmentByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(dbAdj), nil
}

// ListByLeaseID retorna todos os reajustes de um contrato
func (r *LeaseRentAdjustmentRepository) ListByLeaseID(ctx context.Context, leaseID uuid.UUID) ([]*domain.LeaseRentAdjustment, error) {
	dbAdjs, err := r.queries.ListLeaseRentAdjustmentsByLeaseID(ctx, leaseID)
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(dbAdjs), nil
}

// GetLatestByLeaseID retorna o reajuste mais recente de um contrato
func (r *LeaseRentAdjustmentRepository) GetLatestByLeaseID(ctx context.Context, leaseID uuid.UUID) (*domain.LeaseRentAdjustment, error) {
	dbAdj, err := r.queries.GetLatestAdjustmentByLeaseID(ctx, leaseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(dbAdj), nil
}

// CountByLeaseID retorna o total de reajustes de um contrato
func (r *LeaseRentAdjustmentRepository) CountByLeaseID(ctx context.Context, leaseID uuid.UUID) (int64, error) {
	return r.queries.CountAdjustmentsByLeaseID(ctx, leaseID)
}

// Delete remove um reajuste do banco de dados
func (r *LeaseRentAdjustmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteLeaseRentAdjustment(ctx, id)
}

// toDomain converte sqlc.LeaseRentAdjustment para domain.LeaseRentAdjustment
func (r *LeaseRentAdjustmentRepository) toDomain(dbAdj sqlc.LeaseRentAdjustment) *domain.LeaseRentAdjustment {
	previousValue, _ := decimal.NewFromString(dbAdj.PreviousRentValue)
	newValue, _ := decimal.NewFromString(dbAdj.NewRentValue)
	percentage, _ := decimal.NewFromString(dbAdj.AdjustmentPercentage)

	return &domain.LeaseRentAdjustment{
		ID:                   dbAdj.ID,
		LeaseID:              dbAdj.LeaseID,
		PreviousRentValue:    previousValue,
		NewRentValue:         newValue,
		AdjustmentPercentage: percentage,
		AppliedAt:            dbAdj.AppliedAt,
		Reason:               fromNullStringPtr(dbAdj.Reason),
		AppliedBy:            fromNullUUIDPtr(dbAdj.AppliedBy),
		CreatedAt:            dbAdj.CreatedAt,
	}
}

// toDomainSlice converte []sqlc.LeaseRentAdjustment para []*domain.LeaseRentAdjustment
func (r *LeaseRentAdjustmentRepository) toDomainSlice(dbAdjs []sqlc.LeaseRentAdjustment) []*domain.LeaseRentAdjustment {
	adjustments := make([]*domain.LeaseRentAdjustment, len(dbAdjs))
	for i, dbAdj := range dbAdjs {
		adjustments[i] = r.toDomain(dbAdj)
	}
	return adjustments
}
