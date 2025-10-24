package handler

import (
	"time"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
	"github.com/shopspring/decimal"
)

// CreateLeaseRequestDTO representa os dados para criar um contrato
type CreateLeaseRequestDTO struct {
	UnitID                  uuid.UUID       `json:"unit_id" validate:"required"`
	TenantID                uuid.UUID       `json:"tenant_id" validate:"required"`
	ContractSignedDate      time.Time       `json:"contract_signed_date" validate:"required"`
	StartDate               time.Time       `json:"start_date" validate:"required"`
	PaymentDueDay           int             `json:"payment_due_day" validate:"required,min=1,max=31"`
	MonthlyRentValue        decimal.Decimal `json:"monthly_rent_value" validate:"required"`
	PaintingFeeTotal        decimal.Decimal `json:"painting_fee_total" validate:"required"`
	PaintingFeeInstallments int             `json:"painting_fee_installments" validate:"required,min=1,max=4"`
}

// LeaseResponse representa a resposta de um contrato
type LeaseResponse struct {
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
	Status                  string          `json:"status"`
	ParentLeaseID           *uuid.UUID      `json:"parent_lease_id,omitempty"`
	Generation              int             `json:"generation"`
	TotalMonths             int             `json:"total_months"` // Total de meses desde contrato original
	ShouldApplyAdjustment   bool            `json:"should_apply_adjustment"` // Indica se está na geração de reajuste
	DaysUntilExpiry         int             `json:"days_until_expiry"`
	IsExpiringSoon          bool            `json:"is_expiring_soon"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
}

// ToLeaseResponse converte domain.Lease para LeaseResponse
func ToLeaseResponse(lease *domain.Lease) *LeaseResponse {
	return &LeaseResponse{
		ID:                      lease.ID,
		UnitID:                  lease.UnitID,
		TenantID:                lease.TenantID,
		ContractSignedDate:      lease.ContractSignedDate,
		StartDate:               lease.StartDate,
		EndDate:                 lease.EndDate,
		PaymentDueDay:           lease.PaymentDueDay,
		MonthlyRentValue:        lease.MonthlyRentValue,
		PaintingFeeTotal:        lease.PaintingFeeTotal,
		PaintingFeeInstallments: lease.PaintingFeeInstallments,
		PaintingFeePaid:         lease.PaintingFeePaid,
		Status:                  string(lease.Status),
		ParentLeaseID:           lease.ParentLeaseID,
		Generation:              lease.Generation,
		TotalMonths:             lease.GetTotalMonths(),
		ShouldApplyAdjustment:   lease.ShouldApplyAnnualAdjustment(),
		DaysUntilExpiry:         lease.DaysUntilExpiry(),
		IsExpiringSoon:          lease.IsExpiringSoon(),
		CreatedAt:               lease.CreatedAt,
		UpdatedAt:               lease.UpdatedAt,
	}
}

// ToLeaseResponseList converte uma lista de domain.Lease para lista de LeaseResponse
func ToLeaseResponseList(leases []*domain.Lease) []*LeaseResponse {
	responses := make([]*LeaseResponse, len(leases))
	for i, lease := range leases {
		responses[i] = ToLeaseResponse(lease)
	}
	return responses
}

// CreateLeaseResponseDTO representa a resposta ao criar um contrato com pagamentos
type CreateLeaseResponseDTO struct {
	Lease    *LeaseResponse    `json:"lease"`
	Payments []*PaymentResponse `json:"payments"`
}

// ToCreateLeaseResponse converte service.CreateLeaseResponse para CreateLeaseResponseDTO
func ToCreateLeaseResponse(response *service.CreateLeaseResponse) *CreateLeaseResponseDTO {
	return &CreateLeaseResponseDTO{
		Lease:    ToLeaseResponse(response.Lease),
		Payments: ToPaymentResponseList(response.Payments),
	}
}

// RenewLeaseRequestDTO representa os dados para renovação de contrato
type RenewLeaseRequestDTO struct {
	PaintingFeeTotal        decimal.Decimal  `json:"painting_fee_total" validate:"required"`
	PaintingFeeInstallments int              `json:"painting_fee_installments" validate:"required,min=1,max=4"`
	NewRentValue            *decimal.Decimal `json:"new_rent_value,omitempty"`     // Opcional: valor reajustado
	AdjustmentReason        *string          `json:"adjustment_reason,omitempty"`  // Opcional: motivo do reajuste
}

// UpdatePaintingFeePaidRequestDTO representa o valor pago da taxa de pintura
type UpdatePaintingFeePaidRequestDTO struct {
	AmountPaid decimal.Decimal `json:"amount_paid" validate:"required"`
}

// CancelLeaseWithPaymentsRequestDTO representa os dados para cancelar um contrato com seleção de pagamentos
type CancelLeaseWithPaymentsRequestDTO struct {
	PaymentIDs []uuid.UUID `json:"payment_ids" validate:"required,min=1"`
}
