package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository"
	"github.com/shopspring/decimal"
)

// Service layer errors específicos de Lease
var (
	ErrLeaseNotFound               = errors.New("lease not found")
	ErrUnitAlreadyHasActiveLease   = errors.New("unit already has an active lease")
	ErrTenantAlreadyHasActiveLease = errors.New("tenant already has an active lease")
	ErrUnitNotAvailable            = errors.New("unit is not available for rent")
	ErrCannotCancelLease           = errors.New("cannot cancel lease")
	ErrCannotRenewLease            = errors.New("cannot renew lease")
	ErrLeaseAlreadyExpired         = errors.New("lease already expired")
)

// LeaseService contém a lógica de negócio para gestão de contratos
type LeaseService struct {
	leaseRepo  repository.LeaseRepository
	unitRepo   repository.UnitRepository
	tenantRepo repository.TenantRepository
}

// NewLeaseService cria uma nova instância do serviço de contratos
func NewLeaseService(
	leaseRepo repository.LeaseRepository,
	unitRepo repository.UnitRepository,
	tenantRepo repository.TenantRepository,
) *LeaseService {
	return &LeaseService{
		leaseRepo:  leaseRepo,
		unitRepo:   unitRepo,
		tenantRepo: tenantRepo,
	}
}

// CreateLeaseRequest representa os dados necessários para criar um contrato
type CreateLeaseRequest struct {
	UnitID                  uuid.UUID       `json:"unit_id" validate:"required"`
	TenantID                uuid.UUID       `json:"tenant_id" validate:"required"`
	ContractSignedDate      time.Time       `json:"contract_signed_date" validate:"required"`
	StartDate               time.Time       `json:"start_date" validate:"required"`
	PaymentDueDay           int             `json:"payment_due_day" validate:"required,min=1,max=31"`
	MonthlyRentValue        decimal.Decimal `json:"monthly_rent_value" validate:"required"`
	PaintingFeeTotal        decimal.Decimal `json:"painting_fee_total" validate:"required"`
	PaintingFeeInstallments int             `json:"painting_fee_installments" validate:"required,min=1,max=4"`
}

// CreateLease cria um novo contrato de locação com todas as validações de negócio
func (s *LeaseService) CreateLease(ctx context.Context, req CreateLeaseRequest) (*domain.Lease, error) {
	// 1. Validar que a unidade existe
	unit, err := s.unitRepo.GetByID(ctx, req.UnitID)
	if err != nil {
		return nil, fmt.Errorf("error getting unit: %w", err)
	}
	if unit == nil {
		return nil, ErrUnitNotFound
	}

	// 2. Validar que a unidade está disponível
	if !unit.IsAvailable() {
		return nil, ErrUnitNotAvailable
	}

	// 3. Validar que não há um contrato ativo para essa unidade
	existingUnitLease, err := s.leaseRepo.GetActiveByUnitID(ctx, req.UnitID)
	if err != nil {
		return nil, fmt.Errorf("error getting active lease by unit: %w", err)
	}
	if existingUnitLease != nil {
		return nil, ErrUnitAlreadyHasActiveLease
	}

	// 4. Validar que o morador existe
	tenant, err := s.tenantRepo.GetByID(ctx, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("error getting tenant: %w", err)
	}
	if tenant == nil {
		return nil, ErrTenantNotFound
	}

	// 5. Validar que não há um contrato ativo para esse morador
	existingTenantLease, err := s.leaseRepo.GetActiveByTenantID(ctx, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("error getting active lease by tenant: %w", err)
	}
	if existingTenantLease != nil {
		return nil, ErrTenantAlreadyHasActiveLease
	}

	// 6. Criar o contrato usando o domain model
	lease, err := domain.NewLease(req.UnitID, req.TenantID, req.ContractSignedDate, req.StartDate, req.PaymentDueDay, req.MonthlyRentValue, req.PaintingFeeTotal, req.PaintingFeeInstallments)
	if err != nil {
		return nil, fmt.Errorf("error creating lease: %w", err)
	}

	// 7. Persistir o contrato no banco
	if err := s.leaseRepo.Create(ctx, lease); err != nil {
		return nil, fmt.Errorf("error saving lease: %w", err)
	}

	// 8. Atualizar o status da unidade para ocupada
	if err := s.unitRepo.UpdateStatus(ctx, req.UnitID, domain.UnitStatusOccupied); err != nil {
		// TODO: Rollback do lease criado (em um cenário ideal, seria uma transação)
		return nil, fmt.Errorf("error updating unit status: %w", err)
	}

	return lease, nil
}

// GetLeaseByID busca um contrato pelo ID
func (s *LeaseService) GetLeaseByID(ctx context.Context, id uuid.UUID) (*domain.Lease, error) {
	lease, err := s.leaseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting lease: %w", err)
	}
	if lease == nil {
		return nil, ErrLeaseNotFound
	}

	return lease, nil
}

// ListLeases retorna todos os contratos
func (s *LeaseService) ListLeases(ctx context.Context) ([]*domain.Lease, error) {
	leases, err := s.leaseRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing leases: %w", err)
	}
	return leases, nil
}

// ListLeasesByStatus retorna contratos filtrados por status
func (s *LeaseService) ListLeasesByStatus(ctx context.Context, status domain.LeaseStatus) ([]*domain.Lease, error) {
	leases, err := s.leaseRepo.ListByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("error listing leases by status: %w", err)
	}
	return leases, nil
}

// ListLeasesByUnitID retorna todos os contratos de uma unidade
func (s *LeaseService) ListLeasesByUnitID(ctx context.Context, unitID uuid.UUID) ([]*domain.Lease, error) {
	leases, err := s.leaseRepo.ListByUnitID(ctx, unitID)
	if err != nil {
		return nil, fmt.Errorf("error listing leases by unit: %w", err)
	}
	return leases, nil
}

// ListLeasesByTenantID retorna todos os contratos de um morador
func (s *LeaseService) ListLeasesByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*domain.Lease, error) {
	leases, err := s.leaseRepo.ListByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error listing leases by tenant: %w", err)
	}
	return leases, nil
}

// GetExpiringSoonLeases retorna contratos que expiram nos próximos 45 dias
func (s *LeaseService) GetExpiringSoonLeases(ctx context.Context) ([]*domain.Lease, error) {
	leases, err := s.leaseRepo.GetExpiringSoon(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting expiring soon leases: %w", err)
	}
	return leases, nil
}
