package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// Helper function para criar pagamentos de teste
func createTestPayments() []*domain.Payment {
	leaseID := uuid.New()

	// Datas de pagamento dentro do período do teste (2024-03-01 a 2024-04-30)
	paymentDate1 := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	paymentDate3 := time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)

	return []*domain.Payment{
		{
			ID:             uuid.New(),
			LeaseID:        leaseID,
			PaymentType:    domain.PaymentTypeRent,
			ReferenceMonth: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			Amount:         decimal.NewFromInt(800),
			Status:         domain.PaymentStatusPaid,
			DueDate:        time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC),
			PaymentDate:    &paymentDate1, // Data de pagamento em março de 2024
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			LeaseID:        leaseID,
			PaymentType:    domain.PaymentTypeRent,
			ReferenceMonth: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			Amount:         decimal.NewFromInt(800),
			Status:         domain.PaymentStatusPending,
			DueDate:        time.Date(2024, 4, 10, 0, 0, 0, 0, time.UTC),
			PaymentDate:    nil, // Pending, sem data de pagamento - usa DueDate
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			LeaseID:        leaseID,
			PaymentType:    domain.PaymentTypePaintingFee,
			ReferenceMonth: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			Amount:         decimal.NewFromInt(250),
			Status:         domain.PaymentStatusPaid,
			DueDate:        time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC),
			PaymentDate:    &paymentDate3, // Data de pagamento em março de 2024
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}
}

// Test GetFinancialReport - Success
func TestGetFinancialReport_Success(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)
	service := NewReportService(mockPaymentRepo, mockLeaseRepo, mockUnitRepo, mockTenantRepo)

	ctx := context.Background()
	payments := createTestPayments()

	req := FinancialReportRequest{
		StartDate: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
	}

	// Mock repositories
	mockPaymentRepo.On("List", ctx).Return(payments, nil)

	// Mock lease and unit for groupByUnit
	lease := &domain.Lease{
		ID:     payments[0].LeaseID,
		UnitID: uuid.New(),
	}
	unit := &domain.Unit{
		ID:     lease.UnitID,
		Number: "101",
	}

	mockLeaseRepo.On("GetByID", ctx, payments[0].LeaseID).Return(lease, nil)
	mockUnitRepo.On("GetByID", ctx, lease.UnitID).Return(unit, nil)

	// Act
	report, err := service.GetFinancialReport(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 3, report.TotalPayments)

	// Verificar período
	assert.Equal(t, req.StartDate, report.Period.StartDate)
	assert.Equal(t, req.EndDate, report.Period.EndDate)

	// Verificar summary
	expectedTotal := decimal.NewFromInt(1850) // 800 + 800 + 250
	assert.Equal(t, expectedTotal, report.Summary.TotalRevenue)
	assert.Equal(t, decimal.NewFromInt(1050), report.Summary.PaidAmount)   // 800 + 250
	assert.Equal(t, decimal.NewFromInt(800), report.Summary.PendingAmount) // 800

	// Verificar agrupamento por tipo
	assert.Len(t, report.ByType, 2)
	assert.Contains(t, report.ByType, "rent")
	assert.Contains(t, report.ByType, "painting_fee")
	assert.Equal(t, decimal.NewFromInt(1600), report.ByType["rent"].Amount)
	assert.Equal(t, 2, report.ByType["rent"].Count)

	// Verificar agrupamento por mês
	assert.True(t, len(report.ByMonth) >= 2)

	// Verificar agrupamento por unidade
	assert.Len(t, report.ByUnit, 1)
	assert.Equal(t, "101", report.ByUnit[0].UnitNumber)

	mockPaymentRepo.AssertExpectations(t)
	mockLeaseRepo.AssertExpectations(t)
	mockUnitRepo.AssertExpectations(t)
}

// Test GetFinancialReport - Invalid date range
func TestGetFinancialReport_InvalidDateRange(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)
	service := NewReportService(mockPaymentRepo, mockLeaseRepo, mockUnitRepo, mockTenantRepo)

	ctx := context.Background()

	req := FinancialReportRequest{
		StartDate: time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC), // End before start
	}

	// Act
	report, err := service.GetFinancialReport(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidDateRange, err)
	assert.Nil(t, report)
}

// Test GetFinancialReport - Filter by payment type
func TestGetFinancialReport_FilterByPaymentType(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)
	service := NewReportService(mockPaymentRepo, mockLeaseRepo, mockUnitRepo, mockTenantRepo)

	ctx := context.Background()
	payments := createTestPayments()

	rentType := domain.PaymentTypeRent
	req := FinancialReportRequest{
		StartDate:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
		PaymentType: &rentType,
	}

	mockPaymentRepo.On("List", ctx).Return(payments, nil)

	// Mock lease and unit
	lease := &domain.Lease{
		ID:     payments[0].LeaseID,
		UnitID: uuid.New(),
	}
	unit := &domain.Unit{
		ID:     lease.UnitID,
		Number: "101",
	}

	mockLeaseRepo.On("GetByID", ctx, payments[0].LeaseID).Return(lease, nil)
	mockUnitRepo.On("GetByID", ctx, lease.UnitID).Return(unit, nil)

	// Act
	report, err := service.GetFinancialReport(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 2, report.TotalPayments) // Apenas rent, não painting_fee

	// Verificar que só tem rent
	assert.Len(t, report.ByType, 1)
	assert.Contains(t, report.ByType, "rent")
	assert.NotContains(t, report.ByType, "painting_fee")

	mockPaymentRepo.AssertExpectations(t)
}

// Test GetFinancialReport - Empty payments
func TestGetFinancialReport_EmptyPayments(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)
	service := NewReportService(mockPaymentRepo, mockLeaseRepo, mockUnitRepo, mockTenantRepo)

	ctx := context.Background()

	req := FinancialReportRequest{
		StartDate: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
	}

	mockPaymentRepo.On("List", ctx).Return([]*domain.Payment{}, nil)

	// Act
	report, err := service.GetFinancialReport(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 0, report.TotalPayments)
	assert.True(t, report.Summary.TotalRevenue.IsZero())
	assert.Len(t, report.ByType, 0)
	assert.Len(t, report.ByUnit, 0)

	mockPaymentRepo.AssertExpectations(t)
}

// Test GetPaymentHistoryReport - Success
func TestGetPaymentHistoryReport_Success(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)
	service := NewReportService(mockPaymentRepo, mockLeaseRepo, mockUnitRepo, mockTenantRepo)

	ctx := context.Background()
	payments := createTestPayments()

	req := PaymentHistoryRequest{}

	// Mock repositories
	mockPaymentRepo.On("List", ctx).Return(payments, nil)

	// Mock lease, unit, tenant para cada payment
	for i, p := range payments {
		lease := &domain.Lease{
			ID:       p.LeaseID,
			UnitID:   uuid.New(),
			TenantID: uuid.New(),
		}
		unit := &domain.Unit{
			ID:     lease.UnitID,
			Number: "10" + string(rune('1'+i)),
		}
		tenant := &domain.Tenant{
			ID:       lease.TenantID,
			FullName: "Test Tenant " + string(rune('A'+i)),
		}

		mockLeaseRepo.On("GetByID", ctx, p.LeaseID).Return(lease, nil).Once()
		mockUnitRepo.On("GetByID", ctx, lease.UnitID).Return(unit, nil).Once()
		mockTenantRepo.On("GetByID", ctx, lease.TenantID).Return(tenant, nil).Once()
	}

	// Act
	report, err := service.GetPaymentHistoryReport(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 3, report.TotalCount)
	assert.Len(t, report.Payments, 3)

	// Verificar total amount
	expectedTotal := decimal.NewFromInt(1850)
	assert.Equal(t, expectedTotal, report.TotalAmount)

	// Verificar que cada item tem os dados necessários
	for _, item := range report.Payments {
		assert.NotEqual(t, uuid.Nil, item.PaymentID)
		assert.NotEqual(t, uuid.Nil, item.LeaseID)
		assert.NotEmpty(t, item.UnitNumber)
		assert.NotEmpty(t, item.TenantName)
	}

	mockPaymentRepo.AssertExpectations(t)
	mockLeaseRepo.AssertExpectations(t)
	mockUnitRepo.AssertExpectations(t)
	mockTenantRepo.AssertExpectations(t)
}

// Test GetPaymentHistoryReport - Filter by LeaseID
func TestGetPaymentHistoryReport_FilterByLeaseID(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)
	service := NewReportService(mockPaymentRepo, mockLeaseRepo, mockUnitRepo, mockTenantRepo)

	ctx := context.Background()
	payments := createTestPayments()
	leaseID := payments[0].LeaseID

	req := PaymentHistoryRequest{
		LeaseID: &leaseID,
	}

	mockPaymentRepo.On("ListByLeaseID", ctx, leaseID).Return(payments, nil)

	// Mock lease, unit, tenant
	lease := &domain.Lease{
		ID:       leaseID,
		UnitID:   uuid.New(),
		TenantID: uuid.New(),
	}
	unit := &domain.Unit{
		ID:     lease.UnitID,
		Number: "101",
	}
	tenant := &domain.Tenant{
		ID:       lease.TenantID,
		FullName: "Test Tenant",
	}

	mockLeaseRepo.On("GetByID", ctx, leaseID).Return(lease, nil)
	mockUnitRepo.On("GetByID", ctx, lease.UnitID).Return(unit, nil)
	mockTenantRepo.On("GetByID", ctx, lease.TenantID).Return(tenant, nil)

	// Act
	report, err := service.GetPaymentHistoryReport(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 3, report.TotalCount)

	mockPaymentRepo.AssertExpectations(t)
}

// Test GetPaymentHistoryReport - Filter by Status
func TestGetPaymentHistoryReport_FilterByStatus(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)
	service := NewReportService(mockPaymentRepo, mockLeaseRepo, mockUnitRepo, mockTenantRepo)

	ctx := context.Background()

	// Criar apenas pagamentos pagos
	paidPayments := []*domain.Payment{createTestPayments()[0], createTestPayments()[2]}

	status := domain.PaymentStatusPaid
	req := PaymentHistoryRequest{
		Status: &status,
	}

	mockPaymentRepo.On("ListByStatus", ctx, status).Return(paidPayments, nil)

	// Mock lease, unit, tenant
	for _, p := range paidPayments {
		lease := &domain.Lease{
			ID:       p.LeaseID,
			UnitID:   uuid.New(),
			TenantID: uuid.New(),
		}
		unit := &domain.Unit{
			ID:     lease.UnitID,
			Number: "101",
		}
		tenant := &domain.Tenant{
			ID:       lease.TenantID,
			FullName: "Test Tenant",
		}

		mockLeaseRepo.On("GetByID", ctx, p.LeaseID).Return(lease, nil).Once()
		mockUnitRepo.On("GetByID", ctx, lease.UnitID).Return(unit, nil).Once()
		mockTenantRepo.On("GetByID", ctx, lease.TenantID).Return(tenant, nil).Once()
	}

	// Act
	report, err := service.GetPaymentHistoryReport(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 2, report.TotalCount) // Apenas pagamentos paid

	mockPaymentRepo.AssertExpectations(t)
}

// Test GetPaymentHistoryReport - Filter by date range
func TestGetPaymentHistoryReport_FilterByDateRange(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)
	service := NewReportService(mockPaymentRepo, mockLeaseRepo, mockUnitRepo, mockTenantRepo)

	ctx := context.Background()
	payments := createTestPayments()

	startDate := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)

	req := PaymentHistoryRequest{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	mockPaymentRepo.On("List", ctx).Return(payments, nil)

	// Mock lease, unit, tenant apenas para pagamentos de março
	for i, p := range payments {
		if p.DueDate.After(endDate) {
			continue // Pula pagamentos fora do range
		}

		lease := &domain.Lease{
			ID:       p.LeaseID,
			UnitID:   uuid.New(),
			TenantID: uuid.New(),
		}
		unit := &domain.Unit{
			ID:     lease.UnitID,
			Number: "10" + string(rune('1'+i)),
		}
		tenant := &domain.Tenant{
			ID:       lease.TenantID,
			FullName: "Test Tenant",
		}

		mockLeaseRepo.On("GetByID", ctx, p.LeaseID).Return(lease, nil).Once()
		mockUnitRepo.On("GetByID", ctx, lease.UnitID).Return(unit, nil).Once()
		mockTenantRepo.On("GetByID", ctx, lease.TenantID).Return(tenant, nil).Once()
	}

	// Act
	report, err := service.GetPaymentHistoryReport(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 2, report.TotalCount) // Apenas pagamentos de março

	mockPaymentRepo.AssertExpectations(t)
}
