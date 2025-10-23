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
	leaseRepo      repository.LeaseRepository
	unitRepo       repository.UnitRepository
	tenantRepo     repository.TenantRepository
	paymentService *PaymentService
}

// NewLeaseService cria uma nova instância do serviço de contratos
func NewLeaseService(
	leaseRepo repository.LeaseRepository,
	unitRepo repository.UnitRepository,
	tenantRepo repository.TenantRepository,
	paymentService *PaymentService,
) *LeaseService {
	return &LeaseService{
		leaseRepo:      leaseRepo,
		unitRepo:       unitRepo,
		tenantRepo:     tenantRepo,
		paymentService: paymentService,
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

// CreateLeaseResponse representa o resultado da criação de um contrato com pagamentos
type CreateLeaseResponse struct {
	Lease    *domain.Lease     `json:"lease"`
	Payments []*domain.Payment `json:"payments"`
}

// CreateLease cria um novo contrato de locação com todas as validações de negócio
func (s *LeaseService) CreateLease(ctx context.Context, req CreateLeaseRequest) (*CreateLeaseResponse, error) {
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

	// 9. Gerar pagamentos automaticamente se paymentService estiver disponível
	var payments []*domain.Payment
	if s.paymentService != nil {
		// Gerar TODOS os 6 pagamentos de aluguel mensal (contrato de 6 meses)
		for month := 0; month < 6; month++ {
			referenceMonth := lease.StartDate.AddDate(0, month, 0)
			referenceMonth = time.Date(referenceMonth.Year(), referenceMonth.Month(), 1, 0, 0, 0, 0, time.UTC)

			rentPayment, err := s.paymentService.GenerateMonthlyRentPayment(ctx, GenerateMonthlyRentPaymentRequest{
				LeaseID:        lease.ID,
				ReferenceMonth: referenceMonth,
			})
			if err != nil {
				// Erro ao gerar não deve impedir criação de contrato
				fmt.Printf("Warning: failed to generate rent payment for month %d: %v\n", month+1, err)
			} else {
				payments = append(payments, rentPayment)
			}
		}

		// Gerar pagamentos de taxa de pintura
		paintingFeePayments, err := s.paymentService.GeneratePaintingFeePayments(ctx, GeneratePaintingFeePaymentsRequest{
			LeaseID:      lease.ID,
			Installments: req.PaintingFeeInstallments,
		})
		if err != nil {
			fmt.Printf("Warning: failed to generate painting fee payments: %v\n", err)
		} else {
			payments = append(payments, paintingFeePayments...)
		}
	}

	return &CreateLeaseResponse{
		Lease:    lease,
		Payments: payments,
	}, nil
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

// CancelLease cancela um contrato de locação e libera a unidade
func (s *LeaseService) CancelLease(ctx context.Context, id uuid.UUID) error {
	// 1. Buscar o contrato
	lease, err := s.GetLeaseByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting lease: %w", err)
	}
	if lease == nil {
		return ErrLeaseNotFound
	}

	// 2. Validar que o contrato pode ser cancelado (não expirou e não está cancelado)
	if lease.Status == domain.LeaseStatusExpired {
		return ErrLeaseAlreadyExpired
	}
	if lease.Status == domain.LeaseStatusCancelled {
		return ErrCannotCancelLease
	}

	// 3. Marcar o contrato como cancelado
	lease.MarkAsCancelled()

	// 4. Persistir o contrato no banco
	if err := s.leaseRepo.Update(ctx, lease); err != nil {
		return fmt.Errorf("error updating lease: %w", err)
	}

	// 5. Atualizar o status da unidade para disponível
	if err := s.unitRepo.UpdateStatus(ctx, lease.UnitID, domain.UnitStatusAvailable); err != nil {
		return fmt.Errorf("error updating unit status: %w", err)
	}

	return nil
}

// CancelLeaseWithPaymentsRequest representa os dados para cancelar um contrato com seleção de pagamentos
type CancelLeaseWithPaymentsRequest struct {
	PaymentIDs []uuid.UUID `json:"payment_ids" validate:"required"`
}

// CancelLeaseWithPayments cancela um contrato e os pagamentos selecionados
func (s *LeaseService) CancelLeaseWithPayments(ctx context.Context, id uuid.UUID, paymentIDs []uuid.UUID) error {
	// 1. Buscar o contrato
	lease, err := s.GetLeaseByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting lease: %w", err)
	}
	if lease == nil {
		return ErrLeaseNotFound
	}

	// 2. Validar que o contrato pode ser cancelado (não expirou e não está cancelado)
	if lease.Status == domain.LeaseStatusExpired {
		return ErrLeaseAlreadyExpired
	}
	if lease.Status == domain.LeaseStatusCancelled {
		return ErrCannotCancelLease
	}

	// 3. Cancelar os pagamentos selecionados usando o payment service
	if s.paymentService != nil && len(paymentIDs) > 0 {
		if err := s.paymentService.CancelPayments(ctx, id, paymentIDs); err != nil {
			return fmt.Errorf("error cancelling payments: %w", err)
		}
	}

	// 4. Marcar o contrato como cancelado
	lease.MarkAsCancelled()

	// 5. Persistir o contrato no banco
	if err := s.leaseRepo.Update(ctx, lease); err != nil {
		return fmt.Errorf("error updating lease: %w", err)
	}

	// 6. Atualizar o status da unidade para disponível
	if err := s.unitRepo.UpdateStatus(ctx, lease.UnitID, domain.UnitStatusAvailable); err != nil {
		return fmt.Errorf("error updating unit status: %w", err)
	}

	return nil
}

// UpdatePaintingFeePaid atualiza o valor pago da taxa de pintura
func (s *LeaseService) UpdatePaintingFeePaid(ctx context.Context, leaseID uuid.UUID, amountPaid decimal.Decimal) error {
	// 1. Buscar o contrato
	lease, err := s.GetLeaseByID(ctx, leaseID)
	if err != nil {
		return fmt.Errorf("error getting lease: %w", err)
	}
	if lease == nil {
		return ErrLeaseNotFound
	}

	// 2. Atualizar usando método do domain (valida se não excede o total)
	if err := lease.UpdatePaintingFeePaid(amountPaid); err != nil {
		return fmt.Errorf("error updating painting fee: %w", err)
	}

	// 3. Persistir no banco usando método específico do repository
	if err := s.leaseRepo.UpdatePaintingFeePaid(ctx, leaseID, lease.PaintingFeePaid); err != nil {
		return fmt.Errorf("error saving painting fee update: %w", err)
	}

	return nil
}

// CheckExpiringSoonLeases verifica contratos próximos de expirar e atualiza status
// Este método será usado por um cronjob diário no futuro
func (s *LeaseService) CheckExpiringSoonLeases(ctx context.Context) (int, error) {
	// 1. Buscar contratos que expiram nos próximos 45 dias
	leases, err := s.leaseRepo.GetExpiringSoon(ctx)
	if err != nil {
		return 0, fmt.Errorf("error getting expiring soon leases: %w", err)
	}

	updatedCount := 0

	// 2. Para cada contrato, atualizar status se necessário
	for _, lease := range leases {
		// Só atualiza se ainda estiver como 'active'
		if lease.Status == domain.LeaseStatusActive && lease.IsExpiringSoon() {
			lease.MarkAsExpiringSoon()

			if err := s.leaseRepo.UpdateStatus(ctx, lease.ID, domain.LeaseStatusExpiringSoon); err != nil {
				// Log do erro mas continua processando os outros
				fmt.Printf("error updating lease %s status: %v\n", lease.ID, err)
				continue
			}

			updatedCount++
		}
	}

	return updatedCount, nil
}

// MarkLeaseAsExpired marca um contrato como expirado (usado quando a data passa)
func (s *LeaseService) MarkLeaseAsExpired(ctx context.Context, id uuid.UUID) error {
	// 1. Buscar o contrato
	lease, err := s.GetLeaseByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting lease: %w", err)
	}
	if lease == nil {
		return ErrLeaseNotFound
	}

	// 2. Validar se já expirou a data
	if !lease.IsExpired() {
		return errors.New("lease has not expired yet")
	}

	// 3. Marcar como expirado
	lease.MarkAsExpired()

	// 4. Atualizar no banco
	if err := s.leaseRepo.UpdateStatus(ctx, id, domain.LeaseStatusExpired); err != nil {
		return fmt.Errorf("error marking lease as expired: %w", err)
	}

	// 5. Liberar a unidade (marcar como disponível)
	if err := s.unitRepo.UpdateStatus(ctx, lease.UnitID, domain.UnitStatusAvailable); err != nil {
		return fmt.Errorf("error updating unit status: %w", err)
	}

	return nil
}

// GetLeaseStats retorna estatísticas de contratos
func (s *LeaseService) GetLeaseStats(ctx context.Context) (*LeaseStats, error) {
	total, err := s.leaseRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("error counting leases: %w", err)
	}

	active, err := s.leaseRepo.CountByStatus(ctx, domain.LeaseStatusActive)
	if err != nil {
		return nil, fmt.Errorf("error counting active leases: %w", err)
	}

	expiringSoon, err := s.leaseRepo.CountByStatus(ctx, domain.LeaseStatusExpiringSoon)
	if err != nil {
		return nil, fmt.Errorf("error counting expiring soon leases: %w", err)
	}

	expired, err := s.leaseRepo.CountByStatus(ctx, domain.LeaseStatusExpired)
	if err != nil {
		return nil, fmt.Errorf("error counting expired leases: %w", err)
	}

	cancelled, err := s.leaseRepo.CountByStatus(ctx, domain.LeaseStatusCancelled)
	if err != nil {
		return nil, fmt.Errorf("error counting cancelled leases: %w", err)
	}

	return &LeaseStats{
		Total:        total,
		Active:       active,
		ExpiringSoon: expiringSoon,
		Expired:      expired,
		Cancelled:    cancelled,
	}, nil
}

// LeaseStats representa estatísticas de contratos
type LeaseStats struct {
	Total        int64 `json:"total"`
	Active       int64 `json:"active"`
	ExpiringSoon int64 `json:"expiring_soon"`
	Expired      int64 `json:"expired"`
	Cancelled    int64 `json:"cancelled"`
}

// RenewLease renova um contrato existente criando um novo contrato
func (s *LeaseService) RenewLease(ctx context.Context, oldLeaseID uuid.UUID, paintingFeeTotal decimal.Decimal, paintingFeeInstallments int) (*CreateLeaseResponse, error) {
	// 1. Buscar o contrato antigo
	oldLease, err := s.GetLeaseByID(ctx, oldLeaseID)
	if err != nil {
		return nil, fmt.Errorf("error getting old lease: %w", err)
	}
	if oldLease == nil {
		return nil, ErrLeaseNotFound
	}

	// 2. Validar que o contrato pode ser renovado
	if !oldLease.CanBeRenewed() {
		return nil, ErrCannotRenewLease
	}

	// 3. Buscar dados atualizados da unidade
	unit, err := s.unitRepo.GetByID(ctx, oldLease.UnitID)
	if err != nil {
		return nil, fmt.Errorf("error getting unit: %w", err)
	}
	if unit == nil {
		return nil, ErrUnitNotFound
	}

	// 4. Criar novo contrato
	// Start date = 1 dia após a data de término do contrato antigo
	newStartDate := oldLease.EndDate.AddDate(0, 0, 1)
	newLease, err := domain.NewLease(
		oldLease.UnitID,
		oldLease.TenantID,
		time.Now(),
		newStartDate,
		oldLease.PaymentDueDay,
		unit.CurrentRentValue,
		paintingFeeTotal,
		paintingFeeInstallments,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating new lease: %w", err)
	}

	// 5. Marcar contrato antigo como expirado
	oldLease.MarkAsExpired()

	// 6. Persistir as mudanças (idealmente em uma transação)
	if err := s.leaseRepo.Update(ctx, oldLease); err != nil {
		return nil, fmt.Errorf("error updating old lease: %w", err)
	}

	if err := s.leaseRepo.Create(ctx, newLease); err != nil {
		// TODO: Rollback do update do oldLease
		return nil, fmt.Errorf("erro creating new lease: %w", err)
	}

	// Aqui a unidade já está como occupied, não precisa atualizar

	// 7. Gerar pagamentos para o contrato renovado
	var payments []*domain.Payment
	if s.paymentService != nil {
		// Gerar TODOS os 6 pagamentos de aluguel mensal (contrato de 6 meses)
		for month := 0; month < 6; month++ {
			referenceMonth := newLease.StartDate.AddDate(0, month, 0)
			referenceMonth = time.Date(referenceMonth.Year(), referenceMonth.Month(), 1, 0, 0, 0, 0, time.UTC)

			rentPayment, err := s.paymentService.GenerateMonthlyRentPayment(ctx, GenerateMonthlyRentPaymentRequest{
				LeaseID:        newLease.ID,
				ReferenceMonth: referenceMonth,
			})
			if err != nil {
				fmt.Printf("Warning: failed to generate rent payment for month %d: %v\n", month+1, err)
			} else {
				payments = append(payments, rentPayment)
			}
		}

		// Gerar pagamentos de taxa de pintura
		paintingFeePayments, err := s.paymentService.GeneratePaintingFeePayments(ctx, GeneratePaintingFeePaymentsRequest{
			LeaseID:      newLease.ID,
			Installments: paintingFeeInstallments,
		})
		if err != nil {
			fmt.Printf("Warning: failed to generate painting fee payments: %v\n", err)
		} else {
			payments = append(payments, paintingFeePayments...)
		}
	}

	return &CreateLeaseResponse{
		Lease:    newLease,
		Payments: payments,
	}, nil
}

// RenewLeaseRequest representa os dados para renovação de contrato
type RenewLeaseRequest struct {
	PaintingFeeTotal        decimal.Decimal `json:"painting_fee_total" validate:"required"`
	PaintingFeeInstallments int             `json:"painting_fee_installments" validate:"required,min=1,max=4"`
}
