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

// PaymentRepo implementa o repository de Payment usando SQLC
type PaymentRepo struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewPaymentRepo cria uma nova instância do repository de Payment
func NewPaymentRepo(db *sql.DB) repository.PaymentRepository {
	return &PaymentRepo{
		db:      db,
		queries: sqlc.New(db),
	}
}

// Create insere um novo pagamento no banco
func (r *PaymentRepo) Create(ctx context.Context, payment *domain.Payment) error {
	params := sqlc.CreatePaymentParams{
		ID:             payment.ID,
		LeaseID:        payment.LeaseID,
		PaymentType:    string(payment.PaymentType),
		ReferenceMonth: payment.ReferenceMonth,
		Amount:         payment.Amount.String(),
		Status:         string(payment.Status),
		DueDate:        payment.DueDate,
		PaymentDate:    toNullTimePtr(payment.PaymentDate),
		PaymentMethod:  toNullStringPtr(paymentMethodToStringPtr(payment.PaymentMethod)),
		ProofUrl:       toNullStringPtr(payment.ProofURL),
		Notes:          toNullStringPtr(payment.Notes),
		CreatedAt:      payment.CreatedAt,
		UpdatedAt:      payment.UpdatedAt,
	}

	_, err := r.queries.CreatePayment(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

// GetByID busca um pagamento pelo ID
func (r *PaymentRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	row, err := r.queries.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return r.toDomain(row), nil
}

// List retorna todos os pagamentos
func (r *PaymentRepo) List(ctx context.Context) ([]*domain.Payment, error) {
	rows, err := r.queries.ListPayments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}

	return r.toDomainList(rows), nil
}

// ListByLeaseID retorna todos os pagamentos de um contrato
func (r *PaymentRepo) ListByLeaseID(ctx context.Context, leaseID uuid.UUID) ([]*domain.Payment, error) {
	rows, err := r.queries.ListPaymentsByLeaseID(ctx, leaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments by lease: %w", err)
	}

	return r.toDomainList(rows), nil
}

// ListByStatus retorna pagamentos filtrados por status
func (r *PaymentRepo) ListByStatus(ctx context.Context, status domain.PaymentStatus) ([]*domain.Payment, error) {
	rows, err := r.queries.ListPaymentsByStatus(ctx, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to list payments by status: %w", err)
	}

	return r.toDomainList(rows), nil
}

// GetOverdue retorna pagamentos atrasados
func (r *PaymentRepo) GetOverdue(ctx context.Context) ([]*domain.Payment, error) {
	rows, err := r.queries.GetOverduePayments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue payments: %w", err)
	}

	return r.toDomainList(rows), nil
}

// GetUpcoming retorna pagamentos com vencimento nos próximos X dias
func (r *PaymentRepo) GetUpcoming(ctx context.Context, days int) ([]*domain.Payment, error) {
	rows, err := r.queries.GetUpcomingPayments(ctx, int32(days))
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming payments: %w", err)
	}

	return r.toDomainList(rows), nil
}

// Update atualiza um pagamento existente
func (r *PaymentRepo) Update(ctx context.Context, payment *domain.Payment) error {
	params := sqlc.UpdatePaymentParams{
		ID:             payment.ID,
		LeaseID:        payment.LeaseID,
		PaymentType:    string(payment.PaymentType),
		ReferenceMonth: payment.ReferenceMonth,
		Amount:         payment.Amount.String(),
		Status:         string(payment.Status),
		DueDate:        payment.DueDate,
		PaymentDate:    toNullTimePtr(payment.PaymentDate),
		PaymentMethod:  toNullStringPtr(paymentMethodToStringPtr(payment.PaymentMethod)),
		ProofUrl:       toNullStringPtr(payment.ProofURL),
		Notes:          toNullStringPtr(payment.Notes),
		UpdatedAt:      payment.UpdatedAt,
	}

	_, err := r.queries.UpdatePayment(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}

// UpdateStatus atualiza apenas o status do pagamento
func (r *PaymentRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus) error {
	params := sqlc.UpdatePaymentStatusParams{
		ID:        id,
		Status:    string(status),
		UpdatedAt: time.Now(),
	}

	_, err := r.queries.UpdatePaymentStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}

// MarkAsPaid marca um pagamento como pago
func (r *PaymentRepo) MarkAsPaid(ctx context.Context, id uuid.UUID, paymentDate time.Time, method domain.PaymentMethod) error {
	params := sqlc.MarkPaymentAsPaidParams{
		ID:            id,
		PaymentDate:   sql.NullTime{Time: paymentDate, Valid: true},
		PaymentMethod: sql.NullString{String: string(method), Valid: true},
		UpdatedAt:     time.Now(),
	}

	_, err := r.queries.MarkPaymentAsPaid(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to mark payment as paid: %w", err)
	}

	return nil
}

// MarkOverduePayments marca todos os pagamentos pendentes vencidos como atrasados
func (r *PaymentRepo) MarkOverduePayments(ctx context.Context) error {
	err := r.queries.MarkPaymentsAsOverdue(ctx)
	if err != nil {
		return fmt.Errorf("failed to mark payments as overdue: %w", err)
	}

	return nil
}

// Cancel cancela um pagamento
func (r *PaymentRepo) Cancel(ctx context.Context, id uuid.UUID) error {
	params := sqlc.CancelPaymentParams{
		ID:        id,
		UpdatedAt: time.Now(),
	}

	_, err := r.queries.CancelPayment(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to cancel payment: %w", err)
	}

	return nil
}

// Delete remove um pagamento
func (r *PaymentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeletePayment(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}

	return nil
}

// Count retorna o total de pagamentos
func (r *PaymentRepo) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountPayments(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count payments: %w", err)
	}

	return count, nil
}

// CountByStatus retorna o total de pagamentos por status
func (r *PaymentRepo) CountByStatus(ctx context.Context, status domain.PaymentStatus) (int64, error) {
	count, err := r.queries.CountPaymentsByStatus(ctx, string(status))
	if err != nil {
		return 0, fmt.Errorf("failed to count payments by status: %w", err)
	}

	return count, nil
}

// CountByLeaseID retorna o total de pagamentos de um contrato
func (r *PaymentRepo) CountByLeaseID(ctx context.Context, leaseID uuid.UUID) (int64, error) {
	count, err := r.queries.CountPaymentsByLeaseID(ctx, leaseID)
	if err != nil {
		return 0, fmt.Errorf("failed to count payments by lease: %w", err)
	}

	return count, nil
}

// CountByLeaseIDAndStatus retorna o total de pagamentos de um contrato por status
func (r *PaymentRepo) CountByLeaseIDAndStatus(ctx context.Context, leaseID uuid.UUID, status domain.PaymentStatus) (int64, error) {
	query := `SELECT COUNT(*) FROM payments WHERE lease_id = $1 AND status = $2`

	var count int64
	err := r.db.QueryRowContext(ctx, query, leaseID, status).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetTotalPaidByLease retorna o valor total pago de um contrato
func (r *PaymentRepo) GetTotalPaidByLease(ctx context.Context, leaseID uuid.UUID) (decimal.Decimal, error) {
	total, err := r.queries.GetTotalPaidByLease(ctx, leaseID)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get total paid by lease: %w", err)
	}

	// Converter string para decimal
	result, err := decimal.NewFromString(total)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to parse total: %w", err)
	}
	return result, nil
}

// GetPendingAmountByLease retorna o valor total pendente de um contrato
func (r *PaymentRepo) GetPendingAmountByLease(ctx context.Context, leaseID uuid.UUID) (decimal.Decimal, error) {
	total, err := r.queries.GetPendingAmountByLease(ctx, leaseID)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get pending amount by lease: %w", err)
	}

	// Converter string para decimal
	result, err := decimal.NewFromString(total)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to parse total: %w", err)
	}
	return result, nil
}

// toDomain converte um registro do banco para o domain model
func (r *PaymentRepo) toDomain(row sqlc.Payment) *domain.Payment {
	amount, _ := decimal.NewFromString(row.Amount)

	return &domain.Payment{
		ID:             row.ID,
		LeaseID:        row.LeaseID,
		PaymentType:    domain.PaymentType(row.PaymentType),
		ReferenceMonth: row.ReferenceMonth,
		Amount:         amount,
		Status:         domain.PaymentStatus(row.Status),
		DueDate:        row.DueDate,
		PaymentDate:    fromNullTimePtr(row.PaymentDate),
		PaymentMethod:  stringToPaymentMethodPtr(fromNullStringPtr(row.PaymentMethod)),
		ProofURL:       fromNullStringPtr(row.ProofUrl),
		Notes:          fromNullStringPtr(row.Notes),
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
	}
}

// toDomainList converte múltiplos registros para domain models
func (r *PaymentRepo) toDomainList(rows []sqlc.Payment) []*domain.Payment {
	payments := make([]*domain.Payment, len(rows))
	for i, row := range rows {
		payments[i] = r.toDomain(row)
	}
	return payments
}
