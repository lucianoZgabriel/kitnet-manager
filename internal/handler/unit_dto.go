package handler

import (
	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/shopspring/decimal"
)

// CreateUnitRequest representa o payload para criar uma unidade
type CreateUnitRequest struct {
	Number             string          `json:"number" validate:"required,min=1,max=10"`
	Floor              int             `json:"floor" validate:"required,min=1"`
	BaseRentValue      decimal.Decimal `json:"base_rent_value" validate:"required"`
	RenovatedRentValue decimal.Decimal `json:"renovated_rent_value" validate:"required"`
}

// UpdateUnitRequest representa o payload para atualizar uma unidade
type UpdateUnitRequest struct {
	Number             string          `json:"number" validate:"required,min=1,max=10"`
	Floor              int             `json:"floor" validate:"required,min=1"`
	IsRenovated        bool            `json:"is_renovated"`
	BaseRentValue      decimal.Decimal `json:"base_rent_value" validate:"required"`
	RenovatedRentValue decimal.Decimal `json:"renovated_rent_value" validate:"required"`
	Notes              string          `json:"notes" validate:"max=500"`
}

// UpdateUnitStatusRequest representa o payload para atualizar status
type UpdateUnitStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=available occupied maintenance renovation"`
}

// UnitResponse representa a resposta com dados de uma unidade
type UnitResponse struct {
	ID                 uuid.UUID       `json:"id"`
	Number             string          `json:"number"`
	Floor              int             `json:"floor"`
	Status             string          `json:"status"`
	IsRenovated        bool            `json:"is_renovated"`
	BaseRentValue      decimal.Decimal `json:"base_rent_value"`
	RenovatedRentValue decimal.Decimal `json:"renovated_rent_value"`
	CurrentRentValue   decimal.Decimal `json:"current_rent_value"`
	Notes              string          `json:"notes,omitempty"`
	CreatedAt          string          `json:"created_at"`
	UpdatedAt          string          `json:"updated_at"`
}

// ToUnitResponse converte domain.Unit para UnitResponse
func ToUnitResponse(unit *domain.Unit) *UnitResponse {
	return &UnitResponse{
		ID:                 unit.ID,
		Number:             unit.Number,
		Floor:              unit.Floor,
		Status:             string(unit.Status),
		IsRenovated:        unit.IsRenovated,
		BaseRentValue:      unit.BaseRentValue,
		RenovatedRentValue: unit.RenovatedRentValue,
		CurrentRentValue:   unit.CurrentRentValue,
		Notes:              unit.Notes,
		CreatedAt:          unit.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          unit.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToUnitResponseList converte slice de units para slice de responses
func ToUnitResponseList(units []*domain.Unit) []*UnitResponse {
	responses := make([]*UnitResponse, len(units))
	for i, unit := range units {
		responses[i] = ToUnitResponse(unit)
	}
	return responses
}
