package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// LeaseRentAdjustment representa um reajuste de valor de aluguel aplicado a um contrato
type LeaseRentAdjustment struct {
	ID                   uuid.UUID       `json:"id"`
	LeaseID              uuid.UUID       `json:"lease_id"`
	PreviousRentValue    decimal.Decimal `json:"previous_rent_value"`
	NewRentValue         decimal.Decimal `json:"new_rent_value"`
	AdjustmentPercentage decimal.Decimal `json:"adjustment_percentage"`
	AppliedAt            time.Time       `json:"applied_at"`
	Reason               *string         `json:"reason,omitempty"`
	AppliedBy            *uuid.UUID      `json:"applied_by,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
}

// NewLeaseRentAdjustment cria um novo registro de reajuste de aluguel
func NewLeaseRentAdjustment(
	leaseID uuid.UUID,
	previousValue, newValue decimal.Decimal,
	reason *string,
	appliedBy *uuid.UUID,
) *LeaseRentAdjustment {
	// Calcular percentual de reajuste: ((novo - antigo) / antigo) * 100
	diff := newValue.Sub(previousValue)
	percentage := diff.Div(previousValue).Mul(decimal.NewFromInt(100))

	return &LeaseRentAdjustment{
		ID:                   uuid.New(),
		LeaseID:              leaseID,
		PreviousRentValue:    previousValue,
		NewRentValue:         newValue,
		AdjustmentPercentage: percentage,
		AppliedAt:            time.Now(),
		Reason:               reason,
		AppliedBy:            appliedBy,
		CreatedAt:            time.Now(),
	}
}

// IsIncrease verifica se o reajuste foi um aumento
func (a *LeaseRentAdjustment) IsIncrease() bool {
	return a.AdjustmentPercentage.GreaterThan(decimal.Zero)
}

// IsDecrease verifica se o reajuste foi uma redução
func (a *LeaseRentAdjustment) IsDecrease() bool {
	return a.AdjustmentPercentage.LessThan(decimal.Zero)
}

// GetAbsoluteDifference retorna o valor absoluto da diferença entre novo e antigo
func (a *LeaseRentAdjustment) GetAbsoluteDifference() decimal.Decimal {
	return a.NewRentValue.Sub(a.PreviousRentValue).Abs()
}

// String retorna uma representação em string do reajuste
func (a *LeaseRentAdjustment) String() string {
	return "LeaseRentAdjustment " + a.ID.String() + " for Lease: " + a.LeaseID.String()
}
