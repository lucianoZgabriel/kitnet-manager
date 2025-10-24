package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// LeaseStatus representa os possíveis status de um contrato
type LeaseStatus string

const (
	LeaseStatusActive       LeaseStatus = "active"
	LeaseStatusExpiringSoon LeaseStatus = "expiring_soon"
	LeaseStatusExpired      LeaseStatus = "expired"
	LeaseStatusCancelled    LeaseStatus = "cancelled"
)

// ValidLeaseStatuses contém todos os status válidos de contrato
var ValidLeaseStatuses = []LeaseStatus{
	LeaseStatusActive,
	LeaseStatusExpiringSoon,
	LeaseStatusExpired,
	LeaseStatusCancelled,
}

// Lease representa um contrato de locação entre uma unidade e um morador
type Lease struct {
	ID                      uuid.UUID       `json:"id"`
	UnitID                  uuid.UUID       `json:"unit_id"`
	TenantID                uuid.UUID       `json:"tenant_id"`
	ContractSignedDate      time.Time       `json:"contract_signed_date"`
	StartDate               time.Time       `json:"start_date"`
	EndDate                 time.Time       `json:"end_date"`
	PaymentDueDay           int             `json:"payment_due_day"`
	MonthlyRentValue        decimal.Decimal `json:"monthly_rent_value"`
	PaintingFeeTotal        decimal.Decimal `json:"painting_fee_total"`
	PaintingFeeInstallments int             `json:"painting_fee_installments"`
	PaintingFeePaid         decimal.Decimal `json:"painting_fee_paid"`
	Status                  LeaseStatus     `json:"status"`
	ParentLeaseID           *uuid.UUID      `json:"parent_lease_id,omitempty"` // ID do contrato anterior (renovação)
	Generation              int             `json:"generation"`                // Geração: 1=original, 2=1ª renovação, etc.
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
}

// Domain errors específicos de Lease
var (
	ErrInvalidLeaseStatus             = errors.New("invalid lease status")
	ErrInvalidPaymentDueDay           = errors.New("payment due day must be between 1 and 31")
	ErrInvalidPaintingFeeInstallments = errors.New("painting fee installments must be 1, 2, 3, or 4")
	ErrInvalidMonthlyRentValue        = errors.New("monthly rent value must be greater than zero")
	ErrInvalidPaintingFeeTotal        = errors.New("painting fee total must be greater than or equal to zero")
	ErrInvalidDates                   = errors.New("start date must be before end date")
	ErrPaintingFeePaidExceedsTotal    = errors.New("painting fee paid cannot exceed total")
	ErrInvalidContractDuration        = errors.New("contract duration must be 6 months")
)

// NewLease cria um novo contrato de locação com valores padrão
func NewLease(unitID, tenantID uuid.UUID, contractSignedDate, startDate time.Time, paymentDueDay int, monthlyRentValue, paintingFeeTotal decimal.Decimal, paintingFeeInstallments int) (*Lease, error) {
	// Calcula automaticamente o end_date (6 meses após start_date)
	endDate := startDate.AddDate(0, 6, 0)

	lease := &Lease{
		ID:                      uuid.New(),
		UnitID:                  unitID,
		TenantID:                tenantID,
		ContractSignedDate:      contractSignedDate,
		StartDate:               startDate,
		EndDate:                 endDate,
		PaymentDueDay:           paymentDueDay,
		MonthlyRentValue:        monthlyRentValue,
		PaintingFeeTotal:        paintingFeeTotal,
		PaintingFeeInstallments: paintingFeeInstallments,
		PaintingFeePaid:         decimal.Zero,
		Status:                  LeaseStatusActive,
		ParentLeaseID:           nil, // Contrato original não tem pai
		Generation:              1,   // Contrato original é geração 1
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}

	// Valida o contrato
	if err := lease.Validate(); err != nil {
		return nil, err
	}

	return lease, nil
}

// Validate verifica se o contrato possui dados válidos
func (l *Lease) Validate() error {
	// Validar payment due day
	if l.PaymentDueDay < 1 || l.PaymentDueDay > 31 {
		return ErrInvalidPaymentDueDay
	}

	// Validar número de parcelas da taxa de pintura
	if l.PaintingFeeInstallments < 1 || l.PaintingFeeInstallments > 4 {
		return ErrInvalidPaintingFeeInstallments
	}

	// Validar valor do aluguel
	if l.MonthlyRentValue.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidMonthlyRentValue
	}

	// Validar taxa de pintura
	if l.PaintingFeeTotal.LessThan(decimal.Zero) {
		return ErrInvalidPaintingFeeTotal
	}

	// Validar datas
	if !l.StartDate.Before(l.EndDate) {
		return ErrInvalidDates
	}

	// Validar que painting_fee_paid não excede o total
	if l.PaintingFeePaid.GreaterThan(l.PaintingFeeTotal) {
		return ErrPaintingFeePaidExceedsTotal
	}

	// Validar status
	if !l.IsValidStatus() {
		return ErrInvalidLeaseStatus
	}

	return nil
}

// IsValidStatus verifica se o status do contrato é válido
func (l *Lease) IsValidStatus() bool {
	for _, validStatus := range ValidLeaseStatuses {
		if l.Status == validStatus {
			return true
		}
	}
	return false
}

// CalculateEndDate calcula a data de término do contrato (6 meses após o início)
func (l *Lease) CalculateEndDate() time.Time {
	return l.StartDate.AddDate(0, 6, 0)
}

// IsExpiringSoon verifica se o contrato está próximo de expirar (menos de 45 dias)
func (l *Lease) IsExpiringSoon() bool {
	now := time.Now()
	daysUntilExpiry := int(l.EndDate.Sub(now).Hours() / 24)
	return daysUntilExpiry <= 45 && daysUntilExpiry > 0
}

// IsExpired verifica se o contrato já expirou
func (l *Lease) IsExpired() bool {
	return time.Now().After(l.EndDate)
}

// IsActive verifica se o contrato está ativo
func (l *Lease) IsActive() bool {
	return l.Status == LeaseStatusActive
}

// IsCancelled verifica se o contrato foi cancelado
func (l *Lease) IsCancelled() bool {
	return l.Status == LeaseStatusCancelled
}

// CanBeRenewed verifica se o contrato pode ser renovado
func (l *Lease) CanBeRenewed() bool {
	// Pode renovar se estiver ativo ou expirando em breve
	return l.Status == LeaseStatusActive || l.Status == LeaseStatusExpiringSoon
}

// RemainingPaintingFee retorna o valor restante da taxa de pintura a ser pago
func (l *Lease) RemainingPaintingFee() decimal.Decimal {
	remaining := l.PaintingFeeTotal.Sub(l.PaintingFeePaid)
	if remaining.LessThan(decimal.Zero) {
		return decimal.Zero
	}
	return remaining
}

// IsPaintingFeeFullyPaid verifica se a taxa de pintura foi totalmente paga
func (l *Lease) IsPaintingFeeFullyPaid() bool {
	return l.PaintingFeePaid.GreaterThanOrEqual(l.PaintingFeeTotal)
}

// CalculatePaintingFeeInstallmentValue calcula o valor de cada parcela da taxa de pintura
func (l *Lease) CalculatePaintingFeeInstallmentValue() decimal.Decimal {
	if l.PaintingFeeInstallments <= 0 {
		return decimal.Zero
	}
	return l.PaintingFeeTotal.Div(decimal.NewFromInt(int64(l.PaintingFeeInstallments)))
}

// MarkAsCancelled marca o contrato como cancelado
func (l *Lease) MarkAsCancelled() {
	l.Status = LeaseStatusCancelled
	l.UpdatedAt = time.Now()
}

// MarkAsExpired marca o contrato como expirado
func (l *Lease) MarkAsExpired() {
	l.Status = LeaseStatusExpired
	l.UpdatedAt = time.Now()
}

// MarkAsExpiringSoon marca o contrato como expirando em breve
func (l *Lease) MarkAsExpiringSoon() {
	if l.Status == LeaseStatusActive && l.IsExpiringSoon() {
		l.Status = LeaseStatusExpiringSoon
		l.UpdatedAt = time.Now()
	}
}

// UpdatePaintingFeePaid atualiza o valor pago da taxa de pintura
func (l *Lease) UpdatePaintingFeePaid(amountPaid decimal.Decimal) error {
	newTotal := l.PaintingFeePaid.Add(amountPaid)

	if newTotal.GreaterThan(l.PaintingFeeTotal) {
		return ErrPaintingFeePaidExceedsTotal
	}

	l.PaintingFeePaid = newTotal
	l.UpdatedAt = time.Now()
	return nil
}

// DaysUntilExpiry retorna quantos dias faltam até o contrato expirar
func (l *Lease) DaysUntilExpiry() int {
	duration := time.Until(l.EndDate)
	days := int(duration.Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

// DurationInMonths retorna a duração do contrato em meses
func (l *Lease) DurationInMonths() int {
	years := l.EndDate.Year() - l.StartDate.Year()
	months := int(l.EndDate.Month() - l.StartDate.Month())
	return years*12 + months
}

// ChangeStatus altera o status do contrato com validação
func (l *Lease) ChangeStatus(newStatus LeaseStatus) error {
	// Validar se o novo status é válido
	tempLease := &Lease{Status: newStatus}
	if !tempLease.IsValidStatus() {
		return ErrInvalidLeaseStatus
	}

	l.Status = newStatus
	l.UpdatedAt = time.Now()
	return nil
}

// String retorna uma representação em string do contrato
func (l *Lease) String() string {
	return "Lease " + l.ID.String() + " (Unit: " + l.UnitID.String() + ", Tenant: " + l.TenantID.String() + ")"
}

// ShouldApplyAnnualAdjustment verifica se o contrato está na geração de reajuste anual
// Reajuste anual ocorre a cada 12 meses = a cada 2 gerações de contrato (6 meses cada)
func (l *Lease) ShouldApplyAnnualAdjustment() bool {
	return l.Generation > 1 && l.Generation%2 == 0
}

// GetTotalMonths retorna o total de meses desde o contrato original
func (l *Lease) GetTotalMonths() int {
	return l.Generation * 6
}

// IsOriginalContract verifica se este é o contrato original (não renovado)
func (l *Lease) IsOriginalContract() bool {
	return l.ParentLeaseID == nil && l.Generation == 1
}

// IsRenewal verifica se este contrato é uma renovação
func (l *Lease) IsRenewal() bool {
	return l.ParentLeaseID != nil && l.Generation > 1
}
