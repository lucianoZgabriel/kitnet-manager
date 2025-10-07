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

// Service layer errors específicos de Payment
var (
	ErrPaymentNotFound         = errors.New("payment not found")
	ErrLeaseNotFoundForPayment = errors.New("lease not found for payment")
	ErrInvalidPaymentAmount    = errors.New("invalid payment amount")
	ErrInvalidInstallments     = errors.New("invalid number of installments")
	ErrPaymentCannotBePaid     = errors.New("payment cannot be paid")
	ErrPaymentAlreadyPaid      = errors.New("payment already paid")
)

// PaymentService contém a lógica de negócio para gestão de pagamentos
type PaymentService struct {
	paymentRepo repository.PaymentRepository
	leaseRepo   repository.LeaseRepository
}

// NewPaymentService cria uma nova instância do serviço de pagamentos
func NewPaymentService(paymentRepo repository.PaymentRepository, leaseRepo repository.LeaseRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		leaseRepo:   leaseRepo,
	}
}

// GenerateMonthlyRentPaymentRequest representa os dados para gerar um pagamento de aluguel
type GenerateMonthlyRentPaymentRequest struct {
	LeaseID        uuid.UUID `json:"lease_id" validaate:"required"`
	ReferenceMonth time.Time `json:"reference_month" validate:"required"`
}

// GenerateMonthlyRentPayment gera um pagamento de aluguel mensal
func (s *PaymentService) GenerateMonthlyRentPayment(ctx context.Context, req GenerateMonthlyRentPaymentRequest) (*domain.Payment, error) {
	// 1. Buscar o contrato
	lease, err := s.leaseRepo.GetByID(ctx, req.LeaseID)
	if err != nil {
		return nil, fmt.Errorf("erro getting lease: %w", err)
	}
	if lease == nil {
		return nil, ErrLeaseNotFoundForPayment
	}

	// 2. Calcular a data de vencimento baseada no payment_due_day
	// Se mês de ref é março/2025 e payment_day_due é 10, due_date será 10/03/2025
	dueDate := time.Date(
		req.ReferenceMonth.Year(),
		req.ReferenceMonth.Month(),
		lease.PaymentDueDay,
		0, 0, 0, 0,
		time.UTC,
	)

	// 3. Criar o pagamento
	payment, err := domain.NewPayment(
		lease.ID,
		domain.PaymentTypeRent,
		req.ReferenceMonth,
		lease.MonthlyRentValue,
		dueDate,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating payment: %w", err)
	}

	// 4. Salvar no banco
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("error saving payment: %w", err)
	}

	return payment, nil
}

// GeneratePaintingFeePaymentsRequest representa os dados para gerar pagamentos de taxa de pintura
type GeneratePaintingFeePaymentsRequest struct {
	LeaseID      uuid.UUID `json:"lease_id" validate:"required"`
	Installments int       `json:"installments" validate:"required,min=1,max=4"`
}

// GeneratePaintingFeePayments gera os pagamentos da taxa de pintura (1x, 2x, 3x ou 4x)
func (s *PaymentService) GeneratePaintingFeePayments(ctx context.Context, req GeneratePaintingFeePaymentsRequest) ([]*domain.Payment, error) {
	// 1. Buscar o contrato
	lease, err := s.leaseRepo.GetByID(ctx, req.LeaseID)
	if err != nil {
		return nil, fmt.Errorf("error getting lease: %w", err)
	}
	if lease == nil {
		return nil, ErrLeaseNotFoundForPayment
	}

	// 2. Validar número de parcelas
	if req.Installments < 1 || req.Installments > 4 {
		return nil, ErrInvalidInstallments
	}

	// 3. Calcular valor de cada parcela
	installmentValue := lease.PaintingFeeTotal.Div(decimal.NewFromInt(int64(req.Installments)))

	// 4. Criar os pagamentos
	payments := make([]*domain.Payment, req.Installments)

	for i := 0; i < req.Installments; i++ {
		// Caluclar o mês de referência (a partir da data do contrato)
		referenceMonth := lease.StartDate.AddDate(0, i, 0)
		referenceMonth = time.Date(referenceMonth.Year(), referenceMonth.Month(), 1, 0, 0, 0, 0, time.UTC)

		// Calcular a data de vencimento (mesmo dia do payment_due_day)
		dueDate := time.Date(
			referenceMonth.Year(),
			referenceMonth.Month(),
			lease.PaymentDueDay,
			0, 0, 0, 0,
			time.UTC,
		)

		// Criar o pagamento
		payment, err := domain.NewPayment(
			lease.ID,
			domain.PaymentTypePaintingFee,
			referenceMonth,
			installmentValue,
			dueDate,
		)
		if err != nil {
			return nil, fmt.Errorf("error creating painting fee payment %d: %w", i+1, err)
		}

		// Salvar no banco
		if err := s.paymentRepo.Create(ctx, payment); err != nil {
			return nil, fmt.Errorf("error saving painting fee payment %d: %w", i+1, err)
		}

		payments[i] = payment
	}

	return payments, nil
}

// GenerateAdjustmentPaymentRequest representa os dados para gerar um pagamento de ajuste
type GenerateAdjustmentPaymentRequest struct {
	LeaseID        uuid.UUID       `json:"lease_id" validate:"required"`
	Amount         decimal.Decimal `json:"amount" validate:"required"`
	ReferenceMonth time.Time       `json:"reference_month" validate:"required"`
	DueDate        time.Time       `json:"due_date" validate:"required"`
	Notes          string          `json:"notes"`
}

// GenerateAdjustmentPayment gera um pagamento de ajuste (proporcional ou outro motivo)
func (s *PaymentService) GenerateAdjustmentPayment(ctx context.Context, req GenerateAdjustmentPaymentRequest) (*domain.Payment, error) {
	// 1. Buscar o contrato
	lease, err := s.leaseRepo.GetByID(ctx, req.LeaseID)
	if err != nil {
		return nil, fmt.Errorf("error getting lease: %w", err)
	}
	if lease == nil {
		return nil, ErrLeaseNotFoundForPayment
	}

	// 2. Validar o valor
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidPaymentAmount
	}

	// 3. Criar o pagamento de ajuste
	payment, err := domain.NewPayment(
		lease.ID,
		domain.PaymentTypeAdjustment,
		req.ReferenceMonth,
		req.Amount,
		req.DueDate,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating adjustment payment: %w", err)
	}

	// 4. Adicionar nota se fornecida
	if req.Notes != "" {
		payment.AddNote(req.Notes)
	}

	// 5. Salvar no banco
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("error saving adjustment payment: %w", err)
	}

	return payment, nil
}

// GetPaymentByID busca um pagamento pelo ID
func (s *PaymentService) GetPaymentByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting payment: %w", err)
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}

// GetPaymentsByLease retorna todos os pagamentos de um contrato
func (s *PaymentService) GetPaymentsByLease(ctx context.Context, leaseID uuid.UUID) ([]*domain.Payment, error) {
	// Validar que o contrato existe
	lease, err := s.leaseRepo.GetByID(ctx, leaseID)
	if err != nil {
		return nil, fmt.Errorf("error getting lease: %w", err)
	}
	if lease == nil {
		return nil, ErrLeaseNotFoundForPayment
	}

	payments, err := s.paymentRepo.ListByLeaseID(ctx, leaseID)
	if err != nil {
		return nil, fmt.Errorf("error listing payments by lease: %w", err)
	}

	return payments, nil
}

// GetOverduePayments retorna todos os pagamentos atrasados
func (s *PaymentService) GetOverduePayments(ctx context.Context) ([]*domain.Payment, error) {
	payments, err := s.paymentRepo.GetOverdue(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting overdue payments: %w", err)
	}
	return payments, nil
}

// GetUpcomingPayments retorna pagamentos com vencimento nos próximos X dias
func (s *PaymentService) GetUpcomingPayments(ctx context.Context, days int) ([]*domain.Payment, error) {
	if days <= 0 {
		days = 7 // Default: próximos 7 dias
	}

	payments, err := s.paymentRepo.GetUpcoming(ctx, days)
	if err != nil {
		return nil, fmt.Errorf("error getting upcoming payments: %w", err)
	}
	return payments, nil
}
