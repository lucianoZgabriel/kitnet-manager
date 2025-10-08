package service

import (
	"context"
	"testing"

	"github.com/lucianoZgabriel/kitnet-manager/internal/repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDashboardRepo - Mock do DashboardRepository
type MockDashboardRepo struct {
	mock.Mock
}

func (m *MockDashboardRepo) GetOccupancyMetrics(ctx context.Context) (*repository.OccupancyMetrics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.OccupancyMetrics), args.Error(1)
}

func (m *MockDashboardRepo) GetMonthlyProjectedRevenue(ctx context.Context) (decimal.Decimal, error) {
	args := m.Called(ctx)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func (m *MockDashboardRepo) GetMonthlyRealizedRevenue(ctx context.Context) (decimal.Decimal, error) {
	args := m.Called(ctx)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func (m *MockDashboardRepo) GetOverdueAmount(ctx context.Context) (decimal.Decimal, error) {
	args := m.Called(ctx)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func (m *MockDashboardRepo) GetTotalPendingAmount(ctx context.Context) (decimal.Decimal, error) {
	args := m.Called(ctx)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

// Test GetOccupancyMetrics - Success with data
func TestGetOccupancyMetrics_Success(t *testing.T) {
	// Arrange
	mockDashboardRepo := new(MockDashboardRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockPaymentRepo := new(MockPaymentRepo)
	service := NewDashboardService(mockDashboardRepo, mockLeaseRepo, mockPaymentRepo)

	ctx := context.Background()

	repoMetrics := &repository.OccupancyMetrics{
		TotalUnits:       31,
		OccupiedUnits:    25,
		AvailableUnits:   4,
		MaintenanceUnits: 1,
		RenovationUnits:  1,
	}

	mockDashboardRepo.On("GetOccupancyMetrics", ctx).Return(repoMetrics, nil)

	// Act
	metrics, err := service.GetOccupancyMetrics(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.Equal(t, int64(31), metrics.TotalUnits)
	assert.Equal(t, int64(25), metrics.OccupiedUnits)
	assert.Equal(t, int64(4), metrics.AvailableUnits)
	assert.Equal(t, int64(1), metrics.MaintenanceUnits)
	assert.Equal(t, int64(1), metrics.RenovationUnits)

	// Verificar cálculos de taxas
	expectedOccupancyRate := (float64(25) / float64(31)) * 100   // ~80.65%
	expectedAvailabilityRate := (float64(4) / float64(31)) * 100 // ~12.90%

	assert.InDelta(t, expectedOccupancyRate, metrics.OccupancyRate, 0.01)
	assert.InDelta(t, expectedAvailabilityRate, metrics.AvailabilityRate, 0.01)

	mockDashboardRepo.AssertExpectations(t)
}

// Test GetOccupancyMetrics - Empty database (zero units)
func TestGetOccupancyMetrics_ZeroUnits(t *testing.T) {
	// Arrange
	mockDashboardRepo := new(MockDashboardRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockPaymentRepo := new(MockPaymentRepo)
	service := NewDashboardService(mockDashboardRepo, mockLeaseRepo, mockPaymentRepo)

	ctx := context.Background()

	repoMetrics := &repository.OccupancyMetrics{
		TotalUnits:       0,
		OccupiedUnits:    0,
		AvailableUnits:   0,
		MaintenanceUnits: 0,
		RenovationUnits:  0,
	}

	mockDashboardRepo.On("GetOccupancyMetrics", ctx).Return(repoMetrics, nil)

	// Act
	metrics, err := service.GetOccupancyMetrics(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.Equal(t, int64(0), metrics.TotalUnits)
	assert.Equal(t, 0.0, metrics.OccupancyRate)    // Deve ser 0 quando não há unidades
	assert.Equal(t, 0.0, metrics.AvailabilityRate) // Deve ser 0 quando não há unidades

	mockDashboardRepo.AssertExpectations(t)
}

// Test GetOccupancyMetrics - Repository error
func TestGetOccupancyMetrics_RepositoryError(t *testing.T) {
	// Arrange
	mockDashboardRepo := new(MockDashboardRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockPaymentRepo := new(MockPaymentRepo)
	service := NewDashboardService(mockDashboardRepo, mockLeaseRepo, mockPaymentRepo)

	ctx := context.Background()

	mockDashboardRepo.On("GetOccupancyMetrics", ctx).Return(nil, assert.AnError)

	// Act
	metrics, err := service.GetOccupancyMetrics(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, metrics)

	mockDashboardRepo.AssertExpectations(t)
}

// Test GetFinancialMetrics - Success
func TestGetFinancialMetrics_Success(t *testing.T) {
	// Arrange
	mockDashboardRepo := new(MockDashboardRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockPaymentRepo := new(MockPaymentRepo)
	service := NewDashboardService(mockDashboardRepo, mockLeaseRepo, mockPaymentRepo)

	ctx := context.Background()

	// Valores de exemplo: 25 unidades ocupadas x R$800 = R$20.000 projetado
	projectedRevenue := decimal.NewFromInt(20000)
	realizedRevenue := decimal.NewFromInt(18000) // R$18.000 recebido (90%)
	overdueAmount := decimal.NewFromInt(1500)    // R$1.500 em atraso (7.5%)
	pendingAmount := decimal.NewFromInt(3500)    // R$3.500 total pendente

	mockDashboardRepo.On("GetMonthlyProjectedRevenue", ctx).Return(projectedRevenue, nil)
	mockDashboardRepo.On("GetMonthlyRealizedRevenue", ctx).Return(realizedRevenue, nil)
	mockDashboardRepo.On("GetOverdueAmount", ctx).Return(overdueAmount, nil)
	mockDashboardRepo.On("GetTotalPendingAmount", ctx).Return(pendingAmount, nil)

	// Act
	metrics, err := service.GetFinancialMetrics(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.Equal(t, projectedRevenue, metrics.MonthlyProjectedRevenue)
	assert.Equal(t, realizedRevenue, metrics.MonthlyRealizedRevenue)
	assert.Equal(t, overdueAmount, metrics.OverdueAmount)
	assert.Equal(t, pendingAmount, metrics.TotalPendingAmount)

	// Verificar cálculos de taxas
	expectedDefaultRate := (1500.0 / 20000.0) * 100     // 7.5%
	expectedCollectionRate := (18000.0 / 20000.0) * 100 // 90%

	assert.InDelta(t, expectedDefaultRate, metrics.DefaultRate, 0.01)
	assert.InDelta(t, expectedCollectionRate, metrics.CollectionRate, 0.01)

	mockDashboardRepo.AssertExpectations(t)
}

// Test GetFinancialMetrics - Zero projected revenue (edge case)
func TestGetFinancialMetrics_ZeroProjectedRevenue(t *testing.T) {
	// Arrange
	mockDashboardRepo := new(MockDashboardRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockPaymentRepo := new(MockPaymentRepo)
	service := NewDashboardService(mockDashboardRepo, mockLeaseRepo, mockPaymentRepo)

	ctx := context.Background()

	// Nenhuma receita projetada (sem contratos ativos)
	projectedRevenue := decimal.Zero
	realizedRevenue := decimal.Zero
	overdueAmount := decimal.Zero
	pendingAmount := decimal.Zero

	mockDashboardRepo.On("GetMonthlyProjectedRevenue", ctx).Return(projectedRevenue, nil)
	mockDashboardRepo.On("GetMonthlyRealizedRevenue", ctx).Return(realizedRevenue, nil)
	mockDashboardRepo.On("GetOverdueAmount", ctx).Return(overdueAmount, nil)
	mockDashboardRepo.On("GetTotalPendingAmount", ctx).Return(pendingAmount, nil)

	// Act
	metrics, err := service.GetFinancialMetrics(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.True(t, metrics.MonthlyProjectedRevenue.IsZero())
	assert.Equal(t, 0.0, metrics.DefaultRate)    // Deve ser 0 quando não há receita projetada
	assert.Equal(t, 0.0, metrics.CollectionRate) // Deve ser 0 quando não há receita projetada

	mockDashboardRepo.AssertExpectations(t)
}

// Test GetFinancialMetrics - Error on GetMonthlyProjectedRevenue
func TestGetFinancialMetrics_ErrorOnProjectedRevenue(t *testing.T) {
	// Arrange
	mockDashboardRepo := new(MockDashboardRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockPaymentRepo := new(MockPaymentRepo)
	service := NewDashboardService(mockDashboardRepo, mockLeaseRepo, mockPaymentRepo)

	ctx := context.Background()

	mockDashboardRepo.On("GetMonthlyProjectedRevenue", ctx).Return(decimal.Zero, assert.AnError)

	// Act
	metrics, err := service.GetFinancialMetrics(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, metrics)

	mockDashboardRepo.AssertExpectations(t)
}

// Test GetFinancialMetrics - Error on GetMonthlyRealizedRevenue
func TestGetFinancialMetrics_ErrorOnRealizedRevenue(t *testing.T) {
	// Arrange
	mockDashboardRepo := new(MockDashboardRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockPaymentRepo := new(MockPaymentRepo)
	service := NewDashboardService(mockDashboardRepo, mockLeaseRepo, mockPaymentRepo)

	ctx := context.Background()

	mockDashboardRepo.On("GetMonthlyProjectedRevenue", ctx).Return(decimal.NewFromInt(20000), nil)
	mockDashboardRepo.On("GetMonthlyRealizedRevenue", ctx).Return(decimal.Zero, assert.AnError)

	// Act
	metrics, err := service.GetFinancialMetrics(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, metrics)

	mockDashboardRepo.AssertExpectations(t)
}

// Test GetFinancialMetrics - Error on GetOverdueAmount
func TestGetFinancialMetrics_ErrorOnOverdueAmount(t *testing.T) {
	// Arrange
	mockDashboardRepo := new(MockDashboardRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockPaymentRepo := new(MockPaymentRepo)
	service := NewDashboardService(mockDashboardRepo, mockLeaseRepo, mockPaymentRepo)

	ctx := context.Background()

	mockDashboardRepo.On("GetMonthlyProjectedRevenue", ctx).Return(decimal.NewFromInt(20000), nil)
	mockDashboardRepo.On("GetMonthlyRealizedRevenue", ctx).Return(decimal.NewFromInt(18000), nil)
	mockDashboardRepo.On("GetOverdueAmount", ctx).Return(decimal.Zero, assert.AnError)

	// Act
	metrics, err := service.GetFinancialMetrics(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, metrics)

	mockDashboardRepo.AssertExpectations(t)
}

// Test GetFinancialMetrics - Error on GetTotalPendingAmount
func TestGetFinancialMetrics_ErrorOnPendingAmount(t *testing.T) {
	// Arrange
	mockDashboardRepo := new(MockDashboardRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockPaymentRepo := new(MockPaymentRepo)
	service := NewDashboardService(mockDashboardRepo, mockLeaseRepo, mockPaymentRepo)

	ctx := context.Background()

	mockDashboardRepo.On("GetMonthlyProjectedRevenue", ctx).Return(decimal.NewFromInt(20000), nil)
	mockDashboardRepo.On("GetMonthlyRealizedRevenue", ctx).Return(decimal.NewFromInt(18000), nil)
	mockDashboardRepo.On("GetOverdueAmount", ctx).Return(decimal.NewFromInt(1500), nil)
	mockDashboardRepo.On("GetTotalPendingAmount", ctx).Return(decimal.Zero, assert.AnError)

	// Act
	metrics, err := service.GetFinancialMetrics(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, metrics)

	mockDashboardRepo.AssertExpectations(t)
}
