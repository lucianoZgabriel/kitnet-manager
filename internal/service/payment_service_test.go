package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPaymentRepo - Mock do PaymentRepository
type MockPaymentRepo struct {
	mock.Mock
}

func (m *MockPaymentRepo) Create(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) List(ctx context.Context) ([]*domain.Payment, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) ListByLeaseID(ctx context.Context, leaseID uuid.UUID) ([]*domain.Payment, error) {
	args := m.Called(ctx, leaseID)
	return args.Get(0).([]*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) ListByStatus(ctx context.Context, status domain.PaymentStatus) ([]*domain.Payment, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) GetOverdue(ctx context.Context) ([]*domain.Payment, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) GetUpcoming(ctx context.Context, days int) ([]*domain.Payment, error) {
	args := m.Called(ctx, days)
	return args.Get(0).([]*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) Update(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockPaymentRepo) MarkAsPaid(ctx context.Context, id uuid.UUID, paymentDate time.Time, method domain.PaymentMethod) error {
	args := m.Called(ctx, id, paymentDate, method)
	return args.Error(0)
}

func (m *MockPaymentRepo) MarkOverduePayments(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPaymentRepo) Cancel(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentRepo) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPaymentRepo) CountByStatus(ctx context.Context, status domain.PaymentStatus) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPaymentRepo) CountByLeaseID(ctx context.Context, leaseID uuid.UUID) (int64, error) {
	args := m.Called(ctx, leaseID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPaymentRepo) GetTotalPaidByLease(ctx context.Context, leaseID uuid.UUID) (decimal.Decimal, error) {
	args := m.Called(ctx, leaseID)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func (m *MockPaymentRepo) GetPendingAmountByLease(ctx context.Context, leaseID uuid.UUID) (decimal.Decimal, error) {
	args := m.Called(ctx, leaseID)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

// Helper function para criar um lease de teste
func createTestLease() *domain.Lease {
	return &domain.Lease{
		ID:                      uuid.New(),
		UnitID:                  uuid.New(),
		TenantID:                uuid.New(),
		ContractSignedDate:      time.Now(),
		StartDate:               time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		EndDate:                 time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
		PaymentDueDay:           10,
		MonthlyRentValue:        decimal.NewFromInt(800),
		PaintingFeeTotal:        decimal.NewFromInt(250),
		PaintingFeeInstallments: 3,
		PaintingFeePaid:         decimal.Zero,
		Status:                  domain.LeaseStatusActive,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}
}

// Test GenerateMonthlyRentPayment - Success
func TestGenerateMonthlyRentPayment_Success(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	service := NewPaymentService(mockPaymentRepo, mockLeaseRepo)

	lease := createTestLease()
	ctx := context.Background()

	req := GenerateMonthlyRentPaymentRequest{
		LeaseID:        lease.ID,
		ReferenceMonth: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
	}

	mockLeaseRepo.On("GetByID", ctx, lease.ID).Return(lease, nil)
	mockPaymentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil)

	// Act
	payment, err := service.GenerateMonthlyRentPayment(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, lease.ID, payment.LeaseID)
	assert.Equal(t, domain.PaymentTypeRent, payment.PaymentType)
	assert.Equal(t, lease.MonthlyRentValue, payment.Amount)
	assert.Equal(t, 10, payment.DueDate.Day()) // payment_due_day = 10
	assert.Equal(t, time.March, payment.DueDate.Month())
	assert.Equal(t, domain.PaymentStatusPending, payment.Status)

	mockLeaseRepo.AssertExpectations(t)
	mockPaymentRepo.AssertExpectations(t)
}

// Test GenerateMonthlyRentPayment - Lease Not Found
func TestGenerateMonthlyRentPayment_LeaseNotFound(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	service := NewPaymentService(mockPaymentRepo, mockLeaseRepo)

	ctx := context.Background()
	leaseID := uuid.New()

	req := GenerateMonthlyRentPaymentRequest{
		LeaseID:        leaseID,
		ReferenceMonth: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
	}

	mockLeaseRepo.On("GetByID", ctx, leaseID).Return(nil, nil)

	// Act
	payment, err := service.GenerateMonthlyRentPayment(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, payment)
	assert.Equal(t, ErrLeaseNotFoundForPayment, err)

	mockLeaseRepo.AssertExpectations(t)
}

// Test GeneratePaintingFeePayments - Success with 3 installments
func TestGeneratePaintingFeePayments_Success_3Installments(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	service := NewPaymentService(mockPaymentRepo, mockLeaseRepo)

	lease := createTestLease()
	lease.PaintingFeeInstallments = 3
	lease.PaintingFeeTotal = decimal.NewFromInt(300) // R$ 300

	ctx := context.Background()

	req := GeneratePaintingFeePaymentsRequest{
		LeaseID:      lease.ID,
		Installments: 3,
	}

	mockLeaseRepo.On("GetByID", ctx, lease.ID).Return(lease, nil)
	mockPaymentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil).Times(3)

	// Act
	payments, err := service.GeneratePaintingFeePayments(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, payments)
	assert.Len(t, payments, 3)

	expectedValue := decimal.NewFromInt(100) // 300 / 3 = 100

	for i, payment := range payments {
		assert.Equal(t, lease.ID, payment.LeaseID)
		assert.Equal(t, domain.PaymentTypePaintingFee, payment.PaymentType)
		assert.True(t, expectedValue.Equal(payment.Amount), "Expected amount %s, got %s", expectedValue, payment.Amount)
		assert.Equal(t, domain.PaymentStatusPending, payment.Status)

		// Verificar que as datas de vencimento são escalonadas (mês a mês)
		expectedMonth := lease.StartDate.AddDate(0, i, 0)
		assert.Equal(t, expectedMonth.Month(), payment.DueDate.Month())
		assert.Equal(t, 10, payment.DueDate.Day()) // payment_due_day = 10
	}

	mockLeaseRepo.AssertExpectations(t)
	mockPaymentRepo.AssertExpectations(t)
}

// Test GeneratePaintingFeePayments - Success with 1 installment
func TestGeneratePaintingFeePayments_Success_1Installment(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	service := NewPaymentService(mockPaymentRepo, mockLeaseRepo)

	lease := createTestLease()
	lease.PaintingFeeTotal = decimal.NewFromInt(250)

	ctx := context.Background()

	req := GeneratePaintingFeePaymentsRequest{
		LeaseID:      lease.ID,
		Installments: 1,
	}

	mockLeaseRepo.On("GetByID", ctx, lease.ID).Return(lease, nil)
	mockPaymentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil)

	// Act
	payments, err := service.GeneratePaintingFeePayments(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, payments)
	assert.Len(t, payments, 1)
	assert.True(t, decimal.NewFromInt(250).Equal(payments[0].Amount), "Expected amount 250, got %s", payments[0].Amount)

	mockLeaseRepo.AssertExpectations(t)
	mockPaymentRepo.AssertExpectations(t)
}

// Test GeneratePaintingFeePayments - Invalid Installments
func TestGeneratePaintingFeePayments_InvalidInstallments(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	service := NewPaymentService(mockPaymentRepo, mockLeaseRepo)

	lease := createTestLease()
	ctx := context.Background()

	req := GeneratePaintingFeePaymentsRequest{
		LeaseID:      lease.ID,
		Installments: 5, // Inválido: máximo é 4
	}

	mockLeaseRepo.On("GetByID", ctx, lease.ID).Return(lease, nil)

	// Act
	payments, err := service.GeneratePaintingFeePayments(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, payments)
	assert.Equal(t, ErrInvalidInstallments, err)

	mockLeaseRepo.AssertExpectations(t)
}

// Test GenerateAdjustmentPayment - Success
func TestGenerateAdjustmentPayment_Success(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	service := NewPaymentService(mockPaymentRepo, mockLeaseRepo)

	lease := createTestLease()
	ctx := context.Background()

	req := GenerateAdjustmentPaymentRequest{
		LeaseID:        lease.ID,
		Amount:         decimal.NewFromInt(50), // R$ 50 de ajuste
		ReferenceMonth: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		DueDate:        time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
		Notes:          "Ajuste proporcional",
	}

	mockLeaseRepo.On("GetByID", ctx, lease.ID).Return(lease, nil)
	mockPaymentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil)

	// Act
	payment, err := service.GenerateAdjustmentPayment(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, lease.ID, payment.LeaseID)
	assert.Equal(t, domain.PaymentTypeAdjustment, payment.PaymentType)
	assert.True(t, decimal.NewFromInt(50).Equal(payment.Amount), "Expected amount 50, got %s", payment.Amount)
	assert.NotNil(t, payment.Notes)
	assert.Equal(t, "Ajuste proporcional", *payment.Notes)

	mockLeaseRepo.AssertExpectations(t)
	mockPaymentRepo.AssertExpectations(t)
}

// Test GenerateAdjustmentPayment - Invalid Amount
func TestGenerateAdjustmentPayment_InvalidAmount(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	service := NewPaymentService(mockPaymentRepo, mockLeaseRepo)

	lease := createTestLease()
	ctx := context.Background()

	req := GenerateAdjustmentPaymentRequest{
		LeaseID:        lease.ID,
		Amount:         decimal.Zero, // Valor inválido
		ReferenceMonth: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		DueDate:        time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
	}

	mockLeaseRepo.On("GetByID", ctx, lease.ID).Return(lease, nil)

	// Act
	payment, err := service.GenerateAdjustmentPayment(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, payment)
	assert.Equal(t, ErrInvalidPaymentAmount, err)

	mockLeaseRepo.AssertExpectations(t)
}

// Test GetPaymentByID - Success
func TestGetPaymentByID_Success(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	service := NewPaymentService(mockPaymentRepo, mockLeaseRepo)

	ctx := context.Background()
	paymentID := uuid.New()

	expectedPayment := &domain.Payment{
		ID:             paymentID,
		LeaseID:        uuid.New(),
		PaymentType:    domain.PaymentTypeRent,
		ReferenceMonth: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		Amount:         decimal.NewFromInt(800),
		Status:         domain.PaymentStatusPending,
		DueDate:        time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mockPaymentRepo.On("GetByID", ctx, paymentID).Return(expectedPayment, nil)

	// Act
	payment, err := service.GetPaymentByID(ctx, paymentID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, paymentID, payment.ID)

	mockPaymentRepo.AssertExpectations(t)
}

// Test GetPaymentByID - Not Found
func TestGetPaymentByID_NotFound(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	service := NewPaymentService(mockPaymentRepo, mockLeaseRepo)

	ctx := context.Background()
	paymentID := uuid.New()

	mockPaymentRepo.On("GetByID", ctx, paymentID).Return(nil, nil)

	// Act
	payment, err := service.GetPaymentByID(ctx, paymentID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, payment)
	assert.Equal(t, ErrPaymentNotFound, err)

	mockPaymentRepo.AssertExpectations(t)
}

// Test GetPaymentsByLease - Success
func TestGetPaymentsByLease_Success(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	service := NewPaymentService(mockPaymentRepo, mockLeaseRepo)

	lease := createTestLease()
	ctx := context.Background()

	expectedPayments := []*domain.Payment{
		{
			ID:          uuid.New(),
			LeaseID:     lease.ID,
			PaymentType: domain.PaymentTypeRent,
			Amount:      decimal.NewFromInt(800),
			Status:      domain.PaymentStatusPending,
		},
		{
			ID:          uuid.New(),
			LeaseID:     lease.ID,
			PaymentType: domain.PaymentTypePaintingFee,
			Amount:      decimal.NewFromInt(100),
			Status:      domain.PaymentStatusPaid,
		},
	}

	mockLeaseRepo.On("GetByID", ctx, lease.ID).Return(lease, nil)
	mockPaymentRepo.On("ListByLeaseID", ctx, lease.ID).Return(expectedPayments, nil)

	// Act
	payments, err := service.GetPaymentsByLease(ctx, lease.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, payments)
	assert.Len(t, payments, 2)

	mockLeaseRepo.AssertExpectations(t)
	mockPaymentRepo.AssertExpectations(t)
}

// Test GetOverduePayments - Success
func TestGetOverduePayments_Success(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	service := NewPaymentService(mockPaymentRepo, mockLeaseRepo)

	ctx := context.Background()

	expectedPayments := []*domain.Payment{
		{
			ID:          uuid.New(),
			LeaseID:     uuid.New(),
			PaymentType: domain.PaymentTypeRent,
			Amount:      decimal.NewFromInt(800),
			Status:      domain.PaymentStatusOverdue,
			DueDate:     time.Now().AddDate(0, 0, -5), // 5 dias atrás
		},
	}

	mockPaymentRepo.On("GetOverdue", ctx).Return(expectedPayments, nil)

	// Act
	payments, err := service.GetOverduePayments(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, payments)
	assert.Len(t, payments, 1)
	assert.Equal(t, domain.PaymentStatusOverdue, payments[0].Status)

	mockPaymentRepo.AssertExpectations(t)
}

// Test GetUpcomingPayments - Success
func TestGetUpcomingPayments_Success(t *testing.T) {
	// Arrange
	mockPaymentRepo := new(MockPaymentRepo)
	mockLeaseRepo := new(MockLeaseRepo)
	service := NewPaymentService(mockPaymentRepo, mockLeaseRepo)

	ctx := context.Background()
	days := 7

	expectedPayments := []*domain.Payment{
		{
			ID:          uuid.New(),
			LeaseID:     uuid.New(),
			PaymentType: domain.PaymentTypeRent,
			Amount:      decimal.NewFromInt(800),
			Status:      domain.PaymentStatusPending,
			DueDate:     time.Now().AddDate(0, 0, 3), // Daqui a 3 dias
		},
	}

	mockPaymentRepo.On("GetUpcoming", ctx, days).Return(expectedPayments, nil)

	// Act
	payments, err := service.GetUpcomingPayments(ctx, days)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, payments)
	assert.Len(t, payments, 1)

	mockPaymentRepo.AssertExpectations(t)
}
