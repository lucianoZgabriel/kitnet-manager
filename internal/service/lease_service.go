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
	adjustmentRepo repository.LeaseRentAdjustmentRepository
}

// NewLeaseService cria uma nova instância do serviço de contratos
func NewLeaseService(
	leaseRepo repository.LeaseRepository,
	unitRepo repository.UnitRepository,
	tenantRepo repository.TenantRepository,
	paymentService *PaymentService,
	adjustmentRepo repository.LeaseRentAdjustmentRepository,
) *LeaseService {
	return &LeaseService{
		leaseRepo:      leaseRepo,
		unitRepo:       unitRepo,
		tenantRepo:     tenantRepo,
		paymentService: paymentService,
		adjustmentRepo: adjustmentRepo,
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
func (s *LeaseService) RenewLease(ctx context.Context, oldLeaseID uuid.UUID, req RenewLeaseRequest, userID *uuid.UUID) (*CreateLeaseResponse, error) {
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

	// 4. Determinar valor do aluguel (usar reajuste se fornecido, senão usar valor da unidade)
	rentValue := unit.CurrentRentValue
	var rentAdjustment *domain.LeaseRentAdjustment

	if req.NewRentValue != nil {
		// Aplicar reajuste
		rentAdjustment = domain.NewLeaseRentAdjustment(
			oldLeaseID, // Registra o reajuste no contrato antigo
			oldLease.MonthlyRentValue,
			*req.NewRentValue,
			req.AdjustmentReason,
			userID,
		)
		rentValue = *req.NewRentValue
	}

	// 5. Criar novo contrato
	// Start date = 1 dia após a data de término do contrato antigo
	newStartDate := oldLease.EndDate.AddDate(0, 0, 1)
	newGeneration := oldLease.Generation + 1

	newLease, err := domain.NewLease(
		oldLease.UnitID,
		oldLease.TenantID,
		time.Now(),
		newStartDate,
		oldLease.PaymentDueDay,
		rentValue,
		req.PaintingFeeTotal,
		req.PaintingFeeInstallments,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating new lease: %w", err)
	}

	// Definir parent e generation
	newLease.ParentLeaseID = &oldLeaseID
	newLease.Generation = newGeneration

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

	// Salvar registro de reajuste (se aplicável)
	if rentAdjustment != nil && s.adjustmentRepo != nil {
		if err := s.adjustmentRepo.Create(ctx, rentAdjustment); err != nil {
			return nil, fmt.Errorf("error saving rent adjustment: %w", err)
		}
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

		// NOTA: Taxa de pintura NÃO é gerada em renovações
		// Taxa de pintura é paga apenas no primeiro contrato (contrato original)
		// O inquilino paga adiantado para que quando sair não precise pagar novamente
	}

	return &CreateLeaseResponse{
		Lease:    newLease,
		Payments: payments,
	}, nil
}

// RenewLeaseRequest representa os dados para renovação de contrato
type RenewLeaseRequest struct {
	PaintingFeeTotal        decimal.Decimal  `json:"painting_fee_total" validate:"required"`
	PaintingFeeInstallments int              `json:"painting_fee_installments" validate:"required,min=1,max=4"`
	NewRentValue            *decimal.Decimal `json:"new_rent_value,omitempty"`     // Valor reajustado (opcional)
	AdjustmentReason        *string          `json:"adjustment_reason,omitempty"`  // Motivo do reajuste (opcional)
}

// ChangePaymentDueDayRequest representa a requisição para alterar dia de vencimento
type ChangePaymentDueDayRequest struct {
	LeaseID          uuid.UUID `json:"lease_id" validate:"required"`
	NewPaymentDueDay int       `json:"new_payment_due_day" validate:"required,min=1,max=31"`
	EffectiveDate    time.Time `json:"effective_date" validate:"required"`
	Reason           string    `json:"reason"`
}

// ProportionalPaymentInfo contém informações do pagamento proporcional gerado
type ProportionalPaymentInfo struct {
	ID              uuid.UUID       `json:"id"`
	ReferencePeriod string          `json:"reference_period"`
	Days            int             `json:"days"`
	Amount          decimal.Decimal `json:"amount"`
	DueDate         time.Time       `json:"due_date"`
	Status          string          `json:"status"`
}

// UpdatedPaymentInfo contém informações sobre pagamentos que tiveram data alterada
type UpdatedPaymentInfo struct {
	ID             uuid.UUID `json:"id"`
	ReferenceMonth time.Time `json:"reference_month"`
	OldDueDate     time.Time `json:"old_due_date"`
	NewDueDate     time.Time `json:"new_due_date"`
}

// ChangePaymentDueDayResponse representa a resposta da mudança
type ChangePaymentDueDayResponse struct {
	LeaseID              uuid.UUID                `json:"lease_id"`
	OldPaymentDueDay     int                      `json:"old_payment_due_day"`
	NewPaymentDueDay     int                      `json:"new_payment_due_day"`
	EffectiveDate        time.Time                `json:"effective_date"`
	ProportionalPayment  *ProportionalPaymentInfo `json:"proportional_payment,omitempty"`
	UpdatedPaymentsCount int                      `json:"updated_payments_count"`
	UpdatedPayments      []UpdatedPaymentInfo     `json:"updated_payments"`
}

// ChangePaymentDueDay altera o dia de vencimento de um contrato e recalcula pagamentos futuros
func (s *LeaseService) ChangePaymentDueDay(ctx context.Context, req ChangePaymentDueDayRequest) (*ChangePaymentDueDayResponse, error) {
	// ==================================================
	// ETAPA 1: VALIDAÇÕES
	// ==================================================

	// 1.1. Buscar o contrato
	lease, err := s.leaseRepo.GetByID(ctx, req.LeaseID)
	if err != nil {
		return nil, fmt.Errorf("error getting lease: %w", err)
	}
	if lease == nil {
		return nil, ErrLeaseNotFound
	}

	// 1.2. Validar que contrato está ativo
	if lease.Status != domain.LeaseStatusActive && lease.Status != domain.LeaseStatusExpiringSoon {
		return nil, errors.New("lease must be active to change payment due day")
	}

	// 1.3. Validar que o novo dia é diferente do atual
	if req.NewPaymentDueDay == lease.PaymentDueDay {
		return nil, errors.New("new payment due day must be different from current")
	}

	// 1.4. Validar que o novo dia está no range válido (1-31)
	if req.NewPaymentDueDay < 1 || req.NewPaymentDueDay > 31 {
		return nil, errors.New("payment due day must be between 1 and 31")
	}

	// 1.5. Validar que a data efetiva não está no passado
	// Comparar apenas a data (ignorando hora/minuto/segundo)
	// Usar UTC para garantir comparação consistente
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	effectiveDateOnly := time.Date(
		req.EffectiveDate.Year(),
		req.EffectiveDate.Month(),
		req.EffectiveDate.Day(),
		0, 0, 0, 0,
		time.UTC,
	)
	if effectiveDateOnly.Before(today) {
		return nil, errors.New("effective date cannot be in the past")
	}

	// 1.6. Validar que a data efetiva está dentro da vigência do contrato
	if req.EffectiveDate.Before(lease.StartDate) || req.EffectiveDate.After(lease.EndDate) {
		return nil, errors.New("effective date must be within lease period")
	}

	// 1.7. Buscar todos os pagamentos para validar a data efetiva
	allPayments, err := s.paymentService.paymentRepo.ListByLeaseID(ctx, lease.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting lease payments for validation: %w", err)
	}

	// 1.8. Validar que a data efetiva faz sentido em relação aos pagamentos existentes
	// A data efetiva deve ser após o último pagamento pago/cancelado ou pelo menos
	// na mesma data ou depois do primeiro pagamento pendente
	var lastPaidOrCancelledDate time.Time
	var firstPendingPayment *domain.Payment
	for _, payment := range allPayments {
		if payment.Status == domain.PaymentStatusPaid || payment.Status == domain.PaymentStatusCancelled {
			if payment.DueDate.After(lastPaidOrCancelledDate) {
				lastPaidOrCancelledDate = payment.DueDate
			}
		} else if (payment.Status == domain.PaymentStatusPending || payment.Status == domain.PaymentStatusOverdue) &&
			(firstPendingPayment == nil || payment.DueDate.Before(firstPendingPayment.DueDate)) {
			firstPendingPayment = payment
		}
	}

	// Se existe um pagamento já pago/cancelado, a data efetiva deve ser após ele
	if !lastPaidOrCancelledDate.IsZero() && req.EffectiveDate.Before(lastPaidOrCancelledDate) {
		return nil, fmt.Errorf("effective date (%s) cannot be before the last paid/cancelled payment date (%s)",
			req.EffectiveDate.Format("2006-01-02"), lastPaidOrCancelledDate.Format("2006-01-02"))
	}

	// ==================================================
	// ETAPA 2: CALCULAR PAGAMENTO PROPORCIONAL
	// ==================================================

	oldDueDay := lease.PaymentDueDay
	newDueDay := req.NewPaymentDueDay

	// Determinar a data do último vencimento no dia antigo
	lastOldDueDate := time.Date(
		req.EffectiveDate.Year(),
		req.EffectiveDate.Month(),
		oldDueDay,
		0, 0, 0, 0,
		time.UTC,
	)

	// Se a data efetiva é antes do dia antigo no mês atual,
	// o último vencimento foi no mês anterior
	if req.EffectiveDate.Day() < oldDueDay {
		lastOldDueDate = lastOldDueDate.AddDate(0, -1, 0)
	}

	// A nova data de vencimento (data efetiva)
	firstNewDueDate := req.EffectiveDate

	// Calcular quantos dias entre o último vencimento antigo e o primeiro novo
	proportionalDays := int(firstNewDueDate.Sub(lastOldDueDate).Hours() / 24)

	// Calcular valor proporcional
	// Valor proporcional = (valor_mensal / 30) * dias_proporcionais
	dailyRate := lease.MonthlyRentValue.Div(decimal.NewFromInt(30))
	proportionalAmount := dailyRate.Mul(decimal.NewFromInt(int64(proportionalDays)))

	// Criar pagamento proporcional
	var proportionalPayment *domain.Payment
	if proportionalDays > 0 && proportionalAmount.GreaterThan(decimal.Zero) {
		// Usar o mês de referência da data efetiva
		referenceMonth := time.Date(
			firstNewDueDate.Year(),
			firstNewDueDate.Month(),
			1, 0, 0, 0, 0,
			time.UTC,
		)

		proportionalPayment, err = domain.NewPayment(
			lease.ID,
			domain.PaymentTypeAdjustment,
			referenceMonth,
			proportionalAmount,
			firstNewDueDate,
		)
		if err != nil {
			return nil, fmt.Errorf("error creating proportional payment: %w", err)
		}

		// Adicionar nota explicativa
		note := fmt.Sprintf(
			"Pagamento proporcional devido à mudança de vencimento do dia %d para dia %d. Período: %s a %s (%d dias)",
			oldDueDay,
			newDueDay,
			lastOldDueDate.Format("02/01/2006"),
			firstNewDueDate.Format("02/01/2006"),
			proportionalDays,
		)
		proportionalPayment.AddNote(note)

		// Salvar no banco
		if err := s.paymentService.paymentRepo.Create(ctx, proportionalPayment); err != nil {
			return nil, fmt.Errorf("error saving proportional payment: %w", err)
		}
	}

	// ==================================================
	// ETAPA 3: RECALCULAR PAGAMENTOS FUTUROS
	// ==================================================

	// Nota: já temos allPayments da validação anterior (linha 625)

	// Filtrar pagamentos futuros que ainda não foram pagos
	var paymentsToUpdate []*domain.Payment
	var paymentToCancel *domain.Payment // Pagamento que será substituído pelo proporcional

	for _, payment := range allPayments {
		// Só considera pagamentos pending ou overdue
		if payment.Status != domain.PaymentStatusPending && payment.Status != domain.PaymentStatusOverdue {
			continue
		}

		// Se a data de vencimento é após a data efetiva, será recalculado
		if payment.DueDate.After(req.EffectiveDate) {
			// Se criamos um pagamento proporcional, o primeiro pagamento futuro
			// será cancelado (pois está sendo substituído pelo proporcional)
			if proportionalPayment != nil && paymentToCancel == nil {
				paymentToCancel = payment
			} else {
				paymentsToUpdate = append(paymentsToUpdate, payment)
			}
		}
	}

	// Cancelar o pagamento que foi substituído pelo proporcional
	var cancelledPaymentInfo *UpdatedPaymentInfo
	if paymentToCancel != nil {
		cancelledPaymentInfo = &UpdatedPaymentInfo{
			ID:             paymentToCancel.ID,
			ReferenceMonth: paymentToCancel.ReferenceMonth,
			OldDueDate:     paymentToCancel.DueDate,
			NewDueDate:     time.Time{}, // Zero value indica cancelamento
		}

		paymentToCancel.MarkAsCancelled()
		if err := s.paymentService.paymentRepo.Update(ctx, paymentToCancel); err != nil {
			return nil, fmt.Errorf("error cancelling replaced payment %s: %w", paymentToCancel.ID, err)
		}
	}

	// Atualizar a due_date de cada pagamento futuro
	updatedPaymentsInfo := make([]UpdatedPaymentInfo, 0, len(paymentsToUpdate))

	for _, payment := range paymentsToUpdate {
		oldDueDate := payment.DueDate

		// Calcular nova due_date mantendo o ano/mês, mas mudando o dia
		newDueDate := time.Date(
			payment.ReferenceMonth.Year(),
			payment.ReferenceMonth.Month(),
			req.NewPaymentDueDay,
			0, 0, 0, 0,
			time.UTC,
		)

		// Atualizar o pagamento
		payment.DueDate = newDueDate
		payment.UpdatedAt = time.Now()

		// Salvar no banco
		if err := s.paymentService.paymentRepo.Update(ctx, payment); err != nil {
			return nil, fmt.Errorf("error updating payment %s: %w", payment.ID, err)
		}

		// Registrar a mudança
		updatedPaymentsInfo = append(updatedPaymentsInfo, UpdatedPaymentInfo{
			ID:             payment.ID,
			ReferenceMonth: payment.ReferenceMonth,
			OldDueDate:     oldDueDate,
			NewDueDate:     newDueDate,
		})
	}

	// Incluir o pagamento cancelado na lista de atualizações (se existir)
	if cancelledPaymentInfo != nil {
		updatedPaymentsInfo = append(updatedPaymentsInfo, *cancelledPaymentInfo)
	}

	// ==================================================
	// ETAPA 4: ATUALIZAR O CONTRATO
	// ==================================================

	oldPaymentDueDay := lease.PaymentDueDay
	lease.PaymentDueDay = req.NewPaymentDueDay
	lease.UpdatedAt = time.Now()

	if err := s.leaseRepo.Update(ctx, lease); err != nil {
		return nil, fmt.Errorf("error updating lease: %w", err)
	}

	// ==================================================
	// ETAPA 5: MONTAR RESPOSTA
	// ==================================================

	response := &ChangePaymentDueDayResponse{
		LeaseID:              lease.ID,
		OldPaymentDueDay:     oldPaymentDueDay,
		NewPaymentDueDay:     req.NewPaymentDueDay,
		EffectiveDate:        req.EffectiveDate,
		UpdatedPaymentsCount: len(updatedPaymentsInfo),
		UpdatedPayments:      updatedPaymentsInfo,
	}

	// Incluir informações do pagamento proporcional se foi criado
	if proportionalPayment != nil {
		response.ProportionalPayment = &ProportionalPaymentInfo{
			ID:              proportionalPayment.ID,
			ReferencePeriod: fmt.Sprintf("%s - %s", lastOldDueDate.Format("02/01/2006"), firstNewDueDate.Format("02/01/2006")),
			Days:            proportionalDays,
			Amount:          proportionalAmount,
			DueDate:         firstNewDueDate,
			Status:          string(proportionalPayment.Status),
		}
	}

	return response, nil
}

// GetLeaseRentAdjustments retorna o histórico de reajustes de aluguel de um contrato
func (s *LeaseService) GetLeaseRentAdjustments(ctx context.Context, leaseID uuid.UUID) ([]*domain.LeaseRentAdjustment, error) {
	// Verificar se o contrato existe
	lease, err := s.leaseRepo.GetByID(ctx, leaseID)
	if err != nil {
		return nil, fmt.Errorf("error getting lease: %w", err)
	}
	if lease == nil {
		return nil, ErrLeaseNotFound
	}

	// Buscar ajustes
	adjustments, err := s.adjustmentRepo.ListByLeaseID(ctx, leaseID)
	if err != nil {
		return nil, fmt.Errorf("error getting lease rent adjustments: %w", err)
	}

	return adjustments, nil
}

// AutoRenewLeases renova automaticamente contratos expirando que não precisam de reajuste
// Contratos que devem aplicar reajuste (gerações pares) não são renovados automaticamente
func (s *LeaseService) AutoRenewLeases(ctx context.Context) (int, error) {
	// Buscar contratos expirando em breve
	expiringLeases, err := s.leaseRepo.GetExpiringSoon(ctx)
	if err != nil {
		return 0, fmt.Errorf("error getting expiring leases: %w", err)
	}

	renewedCount := 0

	for _, lease := range expiringLeases {
		// Pular contratos que devem aplicar reajuste (renovação manual)
		if lease.ShouldApplyAnnualAdjustment() {
			continue
		}

		// Pular se não está em status apropriado
		if lease.Status != domain.LeaseStatusActive && lease.Status != domain.LeaseStatusExpiringSoon {
			continue
		}

		// Verificar se unidade ainda está ocupada
		unit, err := s.unitRepo.GetByID(ctx, lease.UnitID)
		if err != nil {
			fmt.Printf("Warning: failed to get unit for auto-renewal: %v\n", err)
			continue
		}

		if unit.Status != domain.UnitStatusOccupied {
			continue
		}

		// Renovação automática sem taxa de pintura
		// Taxa de pintura é paga apenas no primeiro contrato
		// Usa o valor atual do aluguel (sem reajuste)
		req := RenewLeaseRequest{
			PaintingFeeTotal:        decimal.Zero,
			PaintingFeeInstallments: 0,
		}

		// Renovar contrato
		_, err = s.RenewLease(ctx, lease.ID, req, nil)
		if err != nil {
			fmt.Printf("Warning: failed to auto-renew lease %s: %v\n", lease.ID, err)
			continue
		}

		renewedCount++
		fmt.Printf("✅ Contrato %s renovado automaticamente (geração %d → %d)\n",
			lease.ID, lease.Generation, lease.Generation+1)
	}

	return renewedCount, nil
}
