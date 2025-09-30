package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUnitRepository é um mock do repository para testes
type MockUnitRepository struct {
	mock.Mock
}

func (m *MockUnitRepository) Create(ctx context.Context, unit *domain.Unit) error {
	args := m.Called(ctx, unit)
	return args.Error(0)
}

func (m *MockUnitRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Unit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Unit), args.Error(1)
}

func (m *MockUnitRepository) GetByNumber(ctx context.Context, number string) (*domain.Unit, error) {
	args := m.Called(ctx, number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Unit), args.Error(1)
}

func (m *MockUnitRepository) List(ctx context.Context) ([]*domain.Unit, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Unit), args.Error(1)
}

func (m *MockUnitRepository) ListByStatus(ctx context.Context, status domain.UnitStatus) ([]*domain.Unit, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Unit), args.Error(1)
}

func (m *MockUnitRepository) ListByFloor(ctx context.Context, floor int) ([]*domain.Unit, error) {
	args := m.Called(ctx, floor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Unit), args.Error(1)
}

func (m *MockUnitRepository) ListAvailable(ctx context.Context) ([]*domain.Unit, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Unit), args.Error(1)
}

func (m *MockUnitRepository) Update(ctx context.Context, unit *domain.Unit) error {
	args := m.Called(ctx, unit)
	return args.Error(0)
}

func (m *MockUnitRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.UnitStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockUnitRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUnitRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUnitRepository) CountByStatus(ctx context.Context, status domain.UnitStatus) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

// Tests

func TestUnitService_CreateUnit(t *testing.T) {
	ctx := context.Background()
	baseRent := decimal.NewFromInt(800)
	renovatedRent := decimal.NewFromInt(900)

	t.Run("should create unit successfully", func(t *testing.T) {
		mockRepo := new(MockUnitRepository)
		service := NewUnitService(mockRepo)

		// Mock: número não existe
		mockRepo.On("GetByNumber", ctx, "101").Return(nil, nil)
		// Mock: criação bem-sucedida
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Unit")).Return(nil)

		unit, err := service.CreateUnit(ctx, "101", 1, baseRent, renovatedRent)

		require.NoError(t, err)
		assert.NotNil(t, unit)
		assert.Equal(t, "101", unit.Number)
		assert.Equal(t, 1, unit.Floor)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when number already exists", func(t *testing.T) {
		mockRepo := new(MockUnitRepository)
		service := NewUnitService(mockRepo)

		existingUnit, _ := domain.NewUnit("101", 1, baseRent, renovatedRent)
		mockRepo.On("GetByNumber", ctx, "101").Return(existingUnit, nil)

		unit, err := service.CreateUnit(ctx, "101", 1, baseRent, renovatedRent)

		assert.Error(t, err)
		assert.Nil(t, unit)
		assert.Equal(t, ErrUnitNumberAlreadyExists, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail with invalid data", func(t *testing.T) {
		mockRepo := new(MockUnitRepository)
		service := NewUnitService(mockRepo)

		mockRepo.On("GetByNumber", ctx, "").Return(nil, nil)

		unit, err := service.CreateUnit(ctx, "", 1, baseRent, renovatedRent)

		assert.Error(t, err)
		assert.Nil(t, unit)
		mockRepo.AssertExpectations(t)
	})
}

func TestUnitService_GetUnitByID(t *testing.T) {
	ctx := context.Background()
	unitID := uuid.New()

	t.Run("should get unit by id", func(t *testing.T) {
		mockRepo := new(MockUnitRepository)
		service := NewUnitService(mockRepo)

		expectedUnit, _ := domain.NewUnit("101", 1, decimal.NewFromInt(800), decimal.NewFromInt(900))
		expectedUnit.ID = unitID

		mockRepo.On("GetByID", ctx, unitID).Return(expectedUnit, nil)

		unit, err := service.GetUnitByID(ctx, unitID)

		require.NoError(t, err)
		assert.Equal(t, expectedUnit, unit)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when unit not found", func(t *testing.T) {
		mockRepo := new(MockUnitRepository)
		service := NewUnitService(mockRepo)

		mockRepo.On("GetByID", ctx, unitID).Return(nil, nil)

		unit, err := service.GetUnitByID(ctx, unitID)

		assert.Error(t, err)
		assert.Nil(t, unit)
		assert.Equal(t, ErrUnitNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUnitService_DeleteUnit(t *testing.T) {
	ctx := context.Background()
	unitID := uuid.New()

	t.Run("should delete available unit", func(t *testing.T) {
		mockRepo := new(MockUnitRepository)
		service := NewUnitService(mockRepo)

		unit, _ := domain.NewUnit("101", 1, decimal.NewFromInt(800), decimal.NewFromInt(900))
		unit.ID = unitID
		unit.Status = domain.UnitStatusAvailable

		mockRepo.On("GetByID", ctx, unitID).Return(unit, nil)
		mockRepo.On("Delete", ctx, unitID).Return(nil)

		err := service.DeleteUnit(ctx, unitID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should not delete occupied unit", func(t *testing.T) {
		mockRepo := new(MockUnitRepository)
		service := NewUnitService(mockRepo)

		unit, _ := domain.NewUnit("101", 1, decimal.NewFromInt(800), decimal.NewFromInt(900))
		unit.ID = unitID
		unit.Status = domain.UnitStatusOccupied

		mockRepo.On("GetByID", ctx, unitID).Return(unit, nil)

		err := service.DeleteUnit(ctx, unitID)

		assert.Error(t, err)
		assert.Equal(t, ErrCannotDeleteOccupiedUnit, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUnitService_GetOccupancyStats(t *testing.T) {
	ctx := context.Background()

	t.Run("should calculate occupancy stats correctly", func(t *testing.T) {
		mockRepo := new(MockUnitRepository)
		service := NewUnitService(mockRepo)

		mockRepo.On("Count", ctx).Return(int64(31), nil)
		mockRepo.On("CountByStatus", ctx, domain.UnitStatusOccupied).Return(int64(25), nil)
		mockRepo.On("CountByStatus", ctx, domain.UnitStatusAvailable).Return(int64(4), nil)
		mockRepo.On("CountByStatus", ctx, domain.UnitStatusMaintenance).Return(int64(1), nil)
		mockRepo.On("CountByStatus", ctx, domain.UnitStatusRenovation).Return(int64(1), nil)

		stats, err := service.GetOccupancyStats(ctx)

		require.NoError(t, err)
		assert.Equal(t, int64(31), stats.Total)
		assert.Equal(t, int64(25), stats.Occupied)
		assert.Equal(t, int64(4), stats.Available)
		assert.InDelta(t, 80.64, stats.OccupancyRate, 0.01) // 25/31 * 100
		mockRepo.AssertExpectations(t)
	})
}
