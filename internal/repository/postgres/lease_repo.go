package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository/sqlc"
	"github.com/shopspring/decimal"
)

// LeaseRepo implementa o repository de Lease usando SQLC
type LeaseRepo struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewLeaseRepo cria uma nova instância do repository de Lease
func NewLeaseRepo(db *sql.DB) repository.LeaseRepository {
	return &LeaseRepo{
		db:      db,
		queries: sqlc.New(db),
	}
}

// Create insere um novo contrato no banco
func (r *LeaseRepo) Create(ctx context.Context, lease *domain.Lease) error {
	params := sqlc.CreateLeaseParams{
		ID:                      lease.ID,
		UnitID:                  lease.UnitID,
		TenantID:                lease.TenantID,
		ContractSignedDate:      lease.ContractSignedDate,
		StartDate:               lease.StartDate,
		EndDate:                 lease.EndDate,
		PaymentDueDay:           int32(lease.PaymentDueDay),
		MonthlyRentValue:        lease.MonthlyRentValue.String(),
		PaintingFeeTotal:        lease.PaintingFeeTotal.String(),
		PaintingFeeInstallments: int32(lease.PaintingFeeInstallments),
		PaintingFeePaid:         lease.PaintingFeePaid.String(),
		Status:                  string(lease.Status),
		ParentLeaseID:           toNullUUIDPtr(lease.ParentLeaseID),
		Generation:              int32(lease.Generation),
		CreatedAt:               lease.CreatedAt,
		UpdatedAt:               lease.UpdatedAt,
	}

	_, err := r.queries.CreateLease(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create lease: %w", err)
	}

	return nil
}

// GetByID busca um contrato pelo ID
func (r *LeaseRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Lease, error) {
	row, err := r.queries.GetLeaseByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get lease: %w", err)
	}

	return r.toDomain(row), nil
}

// List retorna todos os contratos
func (r *LeaseRepo) List(ctx context.Context) ([]*domain.Lease, error) {
	rows, err := r.queries.ListLeases(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list leases: %w", err)
	}

	return r.toDomainList(rows), nil
}

// ListByStatus retorna contratos filtrados por status
func (r *LeaseRepo) ListByStatus(ctx context.Context, status domain.LeaseStatus) ([]*domain.Lease, error) {
	rows, err := r.queries.ListLeasesByStatus(ctx, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to list leases by status: %w", err)
	}

	return r.toDomainList(rows), nil
}

// ListByUnitID retorna todos os contratos de uma unidade
func (r *LeaseRepo) ListByUnitID(ctx context.Context, unitID uuid.UUID) ([]*domain.Lease, error) {
	rows, err := r.queries.ListLeasesByUnitID(ctx, unitID)
	if err != nil {
		return nil, fmt.Errorf("failed to list leases by unit: %w", err)
	}

	return r.toDomainList(rows), nil
}

// ListByTenantID retorna todos os contratos de um morador
func (r *LeaseRepo) ListByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*domain.Lease, error) {
	rows, err := r.queries.ListLeasesByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list leases by tenant: %w", err)
	}

	return r.toDomainList(rows), nil
}

// GetActiveByUnitID busca contrato ativo de uma unidade
func (r *LeaseRepo) GetActiveByUnitID(ctx context.Context, unitID uuid.UUID) (*domain.Lease, error) {
	row, err := r.queries.GetActiveLeaseByUnitID(ctx, unitID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get active lease by unit: %w", err)
	}

	return r.toDomain(row), nil
}

// GetActiveByTenantID busca contrato ativo de um morador
func (r *LeaseRepo) GetActiveByTenantID(ctx context.Context, tenantID uuid.UUID) (*domain.Lease, error) {
	row, err := r.queries.GetActiveLeaseByTenantID(ctx, tenantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get active lease by tenant: %w", err)
	}

	return r.toDomain(row), nil
}

// GetExpiringSoon retorna contratos que expiram nos próximos 45 dias
func (r *LeaseRepo) GetExpiringSoon(ctx context.Context) ([]*domain.Lease, error) {
	rows, err := r.queries.GetExpiringSoonLeases(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring soon leases: %w", err)
	}

	return r.toDomainList(rows), nil
}

// Update atualiza um contrato existente
func (r *LeaseRepo) Update(ctx context.Context, lease *domain.Lease) error {
	params := sqlc.UpdateLeaseParams{
		ID:                      lease.ID,
		UnitID:                  lease.UnitID,
		TenantID:                lease.TenantID,
		ContractSignedDate:      lease.ContractSignedDate,
		StartDate:               lease.StartDate,
		EndDate:                 lease.EndDate,
		PaymentDueDay:           int32(lease.PaymentDueDay),
		MonthlyRentValue:        lease.MonthlyRentValue.String(),
		PaintingFeeTotal:        lease.PaintingFeeTotal.String(),
		PaintingFeeInstallments: int32(lease.PaintingFeeInstallments),
		PaintingFeePaid:         lease.PaintingFeePaid.String(),
		Status:                  string(lease.Status),
		ParentLeaseID:           toNullUUIDPtr(lease.ParentLeaseID),
		Generation:              int32(lease.Generation),
		UpdatedAt:               lease.UpdatedAt,
	}

	_, err := r.queries.UpdateLease(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update lease: %w", err)
	}

	return nil
}

// UpdateStatus atualiza apenas o status do contrato
func (r *LeaseRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.LeaseStatus) error {
	params := sqlc.UpdateLeaseStatusParams{
		ID:        id,
		Status:    string(status),
		UpdatedAt: time.Now(),
	}

	_, err := r.queries.UpdateLeaseStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update lease status: %w", err)
	}

	return nil
}

// UpdatePaintingFeePaid atualiza o valor pago da taxa de pintura
func (r *LeaseRepo) UpdatePaintingFeePaid(ctx context.Context, id uuid.UUID, paintingFeePaid decimal.Decimal) error {
	params := sqlc.UpdatePaintingFeePaidParams{
		ID:              id,
		PaintingFeePaid: paintingFeePaid.String(),
		UpdatedAt:       time.Now(),
	}

	_, err := r.queries.UpdatePaintingFeePaid(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update painting fee paid: %w", err)
	}

	return nil
}

// Delete remove um contrato
func (r *LeaseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteLease(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete lease: %w", err)
	}

	return nil
}

// Count retorna o total de contratos
func (r *LeaseRepo) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountLeases(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count leases: %w", err)
	}

	return count, nil
}

// CountByStatus retorna o total de contratos por status
func (r *LeaseRepo) CountByStatus(ctx context.Context, status domain.LeaseStatus) (int64, error) {
	count, err := r.queries.CountLeasesByStatus(ctx, string(status))
	if err != nil {
		return 0, fmt.Errorf("failed to count leases by status: %w", err)
	}

	return count, nil
}

// toDomain converte um registro do banco para o domain model
func (r *LeaseRepo) toDomain(row sqlc.Lease) *domain.Lease {
	monthlyRentValue, _ := decimal.NewFromString(row.MonthlyRentValue)
	paintingFeeTotal, _ := decimal.NewFromString(row.PaintingFeeTotal)
	paintingFeePaid, _ := decimal.NewFromString(row.PaintingFeePaid)

	return &domain.Lease{
		ID:                      row.ID,
		UnitID:                  row.UnitID,
		TenantID:                row.TenantID,
		ContractSignedDate:      row.ContractSignedDate,
		StartDate:               row.StartDate,
		EndDate:                 row.EndDate,
		PaymentDueDay:           int(row.PaymentDueDay),
		MonthlyRentValue:        monthlyRentValue,
		PaintingFeeTotal:        paintingFeeTotal,
		PaintingFeeInstallments: int(row.PaintingFeeInstallments),
		PaintingFeePaid:         paintingFeePaid,
		Status:                  domain.LeaseStatus(row.Status),
		ParentLeaseID:           fromNullUUIDPtr(row.ParentLeaseID),
		Generation:              int(row.Generation),
		CreatedAt:               row.CreatedAt,
		UpdatedAt:               row.UpdatedAt,
	}
}

// toDomainList converte múltiplos registros para domain models
func (r *LeaseRepo) toDomainList(rows []sqlc.Lease) []*domain.Lease {
	leases := make([]*domain.Lease, len(rows))
	for i, row := range rows {
		leases[i] = r.toDomain(row)
	}
	return leases
}
