package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUnit(t *testing.T) {
	baseRent := decimal.NewFromInt(800)
	renovatedRent := decimal.NewFromInt(900)

	t.Run("should create valid unit", func(t *testing.T) {
		unit, err := NewUnit("101", 1, baseRent, renovatedRent)

		require.NoError(t, err)
		assert.NotNil(t, unit)
		assert.Equal(t, "101", unit.Number)
		assert.Equal(t, 1, unit.Floor)
		assert.Equal(t, UnitStatusAvailable, unit.Status)
		assert.False(t, unit.IsRenovated)
		assert.Equal(t, baseRent, unit.BaseRentValue)
		assert.Equal(t, renovatedRent, unit.RenovatedRentValue)
		assert.Equal(t, baseRent, unit.CurrentRentValue) // Deve ser base pois não está renovada
	})

	t.Run("should fail with empty number", func(t *testing.T) {
		unit, err := NewUnit("", 1, baseRent, renovatedRent)

		assert.Error(t, err)
		assert.Nil(t, unit)
		assert.Equal(t, ErrInvalidUnitNumber, err)
	})

	t.Run("should fail with invalid floor", func(t *testing.T) {
		unit, err := NewUnit("101", 0, baseRent, renovatedRent)

		assert.Error(t, err)
		assert.Nil(t, unit)
		assert.Equal(t, ErrInvalidFloor, err)
	})

	t.Run("should fail with zero rent value", func(t *testing.T) {
		unit, err := NewUnit("101", 1, decimal.Zero, renovatedRent)

		assert.Error(t, err)
		assert.Nil(t, unit)
		assert.Equal(t, ErrInvalidRentValue, err)
	})

	t.Run("should fail when renovated value is lower than base", func(t *testing.T) {
		unit, err := NewUnit("101", 1, baseRent, decimal.NewFromInt(700))

		assert.Error(t, err)
		assert.Nil(t, unit)
		assert.Equal(t, ErrRenovatedValueLower, err)
	})
}

func TestUnit_CalculateCurrentRentValue(t *testing.T) {
	baseRent := decimal.NewFromInt(800)
	renovatedRent := decimal.NewFromInt(900)

	t.Run("should use base rent when not renovated", func(t *testing.T) {
		unit, _ := NewUnit("101", 1, baseRent, renovatedRent)
		unit.IsRenovated = false

		unit.CalculateCurrentRentValue()

		assert.Equal(t, baseRent, unit.CurrentRentValue)
	})

	t.Run("should use renovated rent when renovated", func(t *testing.T) {
		unit, _ := NewUnit("101", 1, baseRent, renovatedRent)
		unit.IsRenovated = true

		unit.CalculateCurrentRentValue()

		assert.Equal(t, renovatedRent, unit.CurrentRentValue)
	})
}

func TestUnit_MarkAsRenovated(t *testing.T) {
	baseRent := decimal.NewFromInt(800)
	renovatedRent := decimal.NewFromInt(900)

	unit, _ := NewUnit("101", 1, baseRent, renovatedRent)
	assert.False(t, unit.IsRenovated)
	assert.Equal(t, baseRent, unit.CurrentRentValue)

	unit.MarkAsRenovated()

	assert.True(t, unit.IsRenovated)
	assert.Equal(t, renovatedRent, unit.CurrentRentValue)
}

func TestUnit_ChangeStatus(t *testing.T) {
	baseRent := decimal.NewFromInt(800)
	renovatedRent := decimal.NewFromInt(900)

	t.Run("should change to valid status", func(t *testing.T) {
		unit, _ := NewUnit("101", 1, baseRent, renovatedRent)

		err := unit.ChangeStatus(UnitStatusMaintenance)

		assert.NoError(t, err)
		assert.Equal(t, UnitStatusMaintenance, unit.Status)
	})

	t.Run("should fail with invalid status", func(t *testing.T) {
		unit, _ := NewUnit("101", 1, baseRent, renovatedRent)

		err := unit.ChangeStatus("invalid_status")

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidStatus, err)
	})
}

func TestUnit_Occupy(t *testing.T) {
	baseRent := decimal.NewFromInt(800)
	renovatedRent := decimal.NewFromInt(900)

	t.Run("should occupy available unit", func(t *testing.T) {
		unit, _ := NewUnit("101", 1, baseRent, renovatedRent)

		err := unit.Occupy()

		assert.NoError(t, err)
		assert.Equal(t, UnitStatusOccupied, unit.Status)
		assert.True(t, unit.IsOccupied())
	})

	t.Run("should fail to occupy non-available unit", func(t *testing.T) {
		unit, _ := NewUnit("101", 1, baseRent, renovatedRent)
		unit.Status = UnitStatusMaintenance

		err := unit.Occupy()

		assert.Error(t, err)
		assert.Equal(t, UnitStatusMaintenance, unit.Status)
	})
}

func TestUnit_MakeAvailable(t *testing.T) {
	baseRent := decimal.NewFromInt(800)
	renovatedRent := decimal.NewFromInt(900)

	unit, _ := NewUnit("101", 1, baseRent, renovatedRent)
	unit.Status = UnitStatusOccupied

	unit.MakeAvailable()

	assert.Equal(t, UnitStatusAvailable, unit.Status)
	assert.True(t, unit.IsAvailable())
}

func TestUnit_CanBeRented(t *testing.T) {
	baseRent := decimal.NewFromInt(800)
	renovatedRent := decimal.NewFromInt(900)

	tests := []struct {
		name     string
		status   UnitStatus
		expected bool
	}{
		{"available can be rented", UnitStatusAvailable, true},
		{"occupied cannot be rented", UnitStatusOccupied, false},
		{"maintenance cannot be rented", UnitStatusMaintenance, false},
		{"renovation cannot be rented", UnitStatusRenovation, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unit, _ := NewUnit("101", 1, baseRent, renovatedRent)
			unit.Status = tt.status

			result := unit.CanBeRented()

			assert.Equal(t, tt.expected, result)
		})
	}
}
