package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PaymentType representa os tipos de pagamento
type PaymentType string

const (
	PaymentTypeRent        PaymentType = "rent"
	PaymentTypePaintingFee PaymentType = "painting_fee"
	PaymentTypeAdjustment  PaymentType = "adjustment"
)

// PaymentStatus representa os possíveis status de um pagamento
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusOverdue   PaymentStatus = "overdue"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// PaymentMethod representa os métodos de pagamento disponíveis
type PaymentMethod string

const (
	PaymentMethodPix          PaymentMethod = "pix"
	PaymentMethodCash         PaymentMethod = "cash"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
)

// ValidPaymentTypes contém todos os tipos válidos de pagamento
var ValidPaymentTypes = []PaymentType{
	PaymentTypeRent,
	PaymentTypePaintingFee,
	PaymentTypeAdjustment,
}

// ValidPaymentStatuses contém todos os status válidos de pagamento
var ValidPaymentStatuses = []PaymentStatus{
	PaymentStatusPending,
	PaymentStatusPaid,
	PaymentStatusOverdue,
	PaymentStatusCancelled,
}

// ValidPaymentMethods contém todos os métodos válidos de pagamento
var ValidPaymentMethods = []PaymentMethod{
	PaymentMethodPix,
	PaymentMethodCash,
	PaymentMethodBankTransfer,
	PaymentMethodCreditCard,
}

// Payment representa um pagamento relacionado a um contrato de locação
type Payment struct {
	ID             uuid.UUID       `json:"id"`
	LeaseID        uuid.UUID       `json:"lease_id"`
	PaymentType    PaymentType     `json:"payment_type"`
	ReferenceMonth time.Time       `json:"reference_month"`
	Amount         decimal.Decimal `json:"amount"`
	Status         PaymentStatus   `json:"status"`
	DueDate        time.Time       `json:"due_date"`
	PaymentDate    *time.Time      `json:"payment_date,omitempty"`
	PaymentMethod  *PaymentMethod  `json:"payment_method,omitempty"`
	ProofURL       *string         `json:"proof_url,omitempty"`
	Notes          *string         `json:"notes,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// Domain errors específicos de Payment
var (
	ErrInvalidPaymentType   = errors.New("invalid payment type")
	ErrInvalidPaymentStatus = errors.New("invalid payment status")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrInvalidAmount        = errors.New("amount must be greater than zero")
	ErrInvalidDueDate       = errors.New("due date cannot be in the past")
	ErrPaymentAlreadyPaid   = errors.New("payment already paid")
	ErrPaymentNotPending    = errors.New("payment must be pending or overdue to be paid")
)

// NewPayment cria um novo pagamento
func NewPayment(leaseID uuid.UUID, paymentType PaymentType, referenceMonth time.Time, amount decimal.Decimal, dueDate time.Time) (*Payment, error) {
	payment := &Payment{
		ID:             uuid.New(),
		LeaseID:        leaseID,
		PaymentType:    paymentType,
		ReferenceMonth: referenceMonth,
		Amount:         amount,
		Status:         PaymentStatusPending,
		DueDate:        dueDate,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Valida o pagamento
	if err := payment.Validate(); err != nil {
		return nil, err
	}

	return payment, nil
}

// Validate verifica se o pagamento possui dados válidos
func (p *Payment) Validate() error {
	// Validar tipo de pagamento
	if !p.IsValidType() {
		return ErrInvalidPaymentType
	}

	// Validar status
	if !p.IsValidStatus() {
		return ErrInvalidPaymentStatus
	}

	// Validar métodos de pagamento (se especificado)
	if p.PaymentMethod != nil && !p.IsValidMethod() {
		return ErrInvalidPaymentMethod
	}

	// Validar valor
	if p.Amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}

	return nil
}

// IsValidType verifica se o tipo de pagamento é válido
func (p *Payment) IsValidType() bool {
	for _, validType := range ValidPaymentTypes {
		if p.PaymentType == validType {
			return true
		}
	}
	return false
}

// IsValidStatus verifica se o status do pagamento é válido
func (p *Payment) IsValidStatus() bool {
	for _, validStatus := range ValidPaymentStatuses {
		if p.Status == validStatus {
			return true
		}
	}
	return false
}

// IsValidMethod verifica se o método de pagamento é válido
func (p *Payment) IsValidMethod() bool {
	if p.PaymentMethod == nil {
		return true
	}
	for _, validMethod := range ValidPaymentMethods {
		if *p.PaymentMethod == validMethod {
			return true
		}
	}
	return false
}

// IsOverdue verifica se o pagamento está atrasado
func (p *Payment) IsOverdue() bool {
	if p.Status == PaymentStatusPaid || p.Status == PaymentStatusCancelled {
		return false
	}
	return time.Now().After(p.DueDate)
}

// IsPaid verifica se o pagamento foi realizado
func (p *Payment) IsPaid() bool {
	return p.Status == PaymentStatusPaid
}

// IsPending verifica se o pagamento está pendente
func (p *Payment) IsPending() bool {
	return p.Status == PaymentStatusPending
}

// IsCancelled verifica se o pagamento foi cancelado
func (p *Payment) IsCancelled() bool {
	return p.Status == PaymentStatusCancelled
}

// CanBePaid verifica se o pagamento pode ser marcado como pago
func (p *Payment) CanBePaid() bool {
	return p.Status == PaymentStatusPending || p.Status == PaymentStatusOverdue
}

// MarkAsPaid marca o pagamento como pago
func (p *Payment) MarkAsPaid(paymentMethod PaymentMethod) error {
	if !p.CanBePaid() {
		return ErrPaymentNotPending
	}

	// Validar método de pagamento
	tempPayment := &Payment{PaymentMethod: &paymentMethod}
	if !tempPayment.IsValidMethod() {
		return ErrInvalidPaymentMethod
	}

	now := time.Now()
	p.Status = PaymentStatusPaid
	p.PaymentDate = &now
	p.PaymentMethod = &paymentMethod
	p.UpdatedAt = now

	return nil
}

// MarkAsOverdue marca o pagamento como atrasado
func (p *Payment) MarkAsOverdue() {
	if p.Status == PaymentStatusPending && p.IsOverdue() {
		p.Status = PaymentStatusOverdue
		p.UpdatedAt = time.Now()
	}
}

// MarkAsCancelled marca o pagamento como cancelado
func (p *Payment) MarkAsCancelled() {
	if p.Status != PaymentStatusPaid {
		p.Status = PaymentStatusCancelled
		p.UpdatedAt = time.Now()
	}
}

// DaysUntilDue retorna quantos dias faltam até o vencimento
func (p *Payment) DaysUntilDue() int {
	duration := time.Until(p.DueDate)
	days := int(duration.Hours() / 24)
	return days
}

// DaysOverdue retorna quantos dias de atraso
func (p *Payment) DaysOverdue() int {
	if !p.IsOverdue() {
		return 0
	}
	duration := time.Since(p.DueDate)
	days := int(duration.Hours() / 24)
	return days
}

// AddNote adiciona uma observação ao pagamento
func (p *Payment) AddNote(note string) {
	p.Notes = &note
	p.UpdatedAt = time.Now()
}

// AddProof adiciona um comprovante ao pagamento
func (p *Payment) AddProof(proofURL string) {
	p.ProofURL = &proofURL
	p.UpdatedAt = time.Now()
}

// String retorna uma representação em string do pagamento
func (p *Payment) String() string {
	return "Payment " + p.ID.String() + " (" + string(p.PaymentType) + " - " + string(p.Status) + ")"
}
