package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// UnitStatus representa os possíveis status de uma unidada
type UnitStatus string

const (
	UnitStatusAvailable   UnitStatus = "available"
	UnitStatusOccupied    UnitStatus = "occupied"
	UnitStatusMaintenance UnitStatus = "maintenance"
	UnitStatusRenovation  UnitStatus = "renovation"
)

// ValidStatuses contém todos os status válidos
var ValidStatuses = []UnitStatus{
	UnitStatusAvailable,
	UnitStatusOccupied,
	UnitStatusMaintenance,
	UnitStatusRenovation,
}

// Unit representa uma unidade/kitnet do prédio
type Unit struct {
	ID                 uuid.UUID       `json:"id"`
	Number             string          `json:"number"`
	Floor              int             `json:"floor"`
	Status             UnitStatus      `json:"status"`
	IsRenovated        bool            `json:"is_renovated"`
	BaseRentValue      decimal.Decimal `json:"base_rent_value"`
	RenovatedRentValue decimal.Decimal `json:"renovated_rent_value"`
	CurrentRentValue   decimal.Decimal `json:"current_rent_value"`
	Notes              string          `json:"notes,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// Domain errors
var (
	ErrInvalidUnitNumber   = errors.New("unit number cannot be empty")
	ErrInvalidFloor        = errors.New("floor must be greater than or equal to 1")
	ErrInvalidRentValue    = errors.New("rent value must be greater than zero")
	ErrInvalidStatus       = errors.New("invalid unit status")
	ErrRenovatedValueLower = errors.New("renovated rent value must be greater than or equal to base rent value")
)

// NewUnit cria uma nova unidade com valores padrão
func NewUnit(number string, floor int, baseRentValue, renovatedRentValue decimal.Decimal) (*Unit, error) {
	unit := &Unit{
		ID:                 uuid.New(),
		Number:             number,
		Floor:              floor,
		Status:             UnitStatusAvailable,
		IsRenovated:        false,
		BaseRentValue:      baseRentValue,
		RenovatedRentValue: renovatedRentValue,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Calcula o valor atual baseado no status de renovação
	unit.CalculateCurrentRentValue()

	// Valida a unidade
	if err := unit.Validate(); err != nil {
		return nil, err
	}

	return unit, nil
}

// Validate verifica se a unidade possui dados válidos
func (u *Unit) Validate() error {
	if u.Number == "" {
		return ErrInvalidUnitNumber
	}

	if u.Floor < 1 {
		return ErrInvalidFloor
	}

	if u.BaseRentValue.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidRentValue
	}

	if u.RenovatedRentValue.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidRentValue
	}

	if u.RenovatedRentValue.LessThan(u.BaseRentValue) {
		return ErrRenovatedValueLower
	}

	if !u.IsValidStatus() {
		return ErrInvalidStatus
	}

	return nil
}

// IsValidStatus verifica se o status da unidade é válido
func (u *Unit) IsValidStatus() bool {
	for _, validStatus := range ValidStatuses {
		if u.Status == validStatus {
			return true
		}
	}
	return false
}

// CalculateCurrentRentValue calcula o valor atual do aluguel
func (u *Unit) CalculateCurrentRentValue() {
	if u.IsRenovated {
		u.CurrentRentValue = u.RenovatedRentValue
	} else {
		u.CurrentRentValue = u.BaseRentValue
	}
}

// MarkAsRenovated marca a unidade como reformada e atualiza o valor do aluguel
func (u *Unit) MarkAsRenovated() {
	u.IsRenovated = true
	u.CalculateCurrentRentValue()
	u.UpdatedAt = time.Now()
}

// ChangeStatus alter o status da unidade com validação
func (u *Unit) ChangeStatus(newStatus UnitStatus) error {
	// Valida se o novo status é válido
	tempUnit := &Unit{Status: newStatus}
	if !tempUnit.IsValidStatus() {
		return ErrInvalidStatus
	}

	u.Status = newStatus
	u.UpdatedAt = time.Now()
	return nil
}

// IsAvailable verifica se a unidade está disponível para locação
func (u *Unit) IsAvailable() bool {
	return u.Status == UnitStatusAvailable
}

// IsOccupied verifica se a unidade está ocupada
func (u *Unit) IsOccupied() bool {
	return u.Status == UnitStatusOccupied
}

// CanBeRented verifica se a unidade pode ser alugada
func (u *Unit) CanBeRented() bool {
	return u.IsAvailable()
}

// Occupy marca a unidade como ocupada (usada ao criar contrato)
func (u *Unit) Occupy() error {
	if !u.CanBeRented() {
		return errors.New("unit is not available for rent")
	}

	u.Status = UnitStatusOccupied
	u.UpdatedAt = time.Now()
	return nil
}

// MakeAvailable marca a unidade como disponível (usada ao encerrar contrato)
func (u *Unit) MakeAvailable() {
	u.Status = UnitStatusAvailable
	u.UpdatedAt = time.Now()
}

// String retorna uma representação em string da unidade
func (u *Unit) String() string {
	return "Unit " + u.Number + " (Floor " + string(rune(u.Floor)) + ")"
}
