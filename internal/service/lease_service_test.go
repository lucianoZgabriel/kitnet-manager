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

// Mock repositories
type MockLeaseRepo struct {
	mock.Mock
}

func (m *MockLeaseRepo) Create(ctx context.Context, lease *domain.Lease) error {
	args := m.Called(ctx, lease)
	return args.Error(0)
}

func (m *MockLeaseRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Lease, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Lease), args.Error(1)
}

func (m *MockLeaseRepo) List(ctx context.Context) ([]*domain.Lease, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Lease), args.Error(1)
}

func (m *MockLeaseRepo) ListByStatus(ctx context.Context, status domain.LeaseStatus) ([]*domain.Lease, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]*domain.Lease), args.Error(1)
}

func (m *MockLeaseRepo) ListByUnitID(ctx context.Context, unitID uuid.UUID) ([]*domain.Lease, error) {
	args := m.Called(ctx, unitID)
	return args.Get(0).([]*domain.Lease), args.Error(1)
}

func (m *MockLeaseRepo) ListByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*domain.Lease, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]*domain.Lease), args.Error(1)
}

func (m *MockLeaseRepo) GetActiveByUnitID(ctx context.Context, unitID uuid.UUID) (*domain.Lease, error) {
	args := m.Called(ctx, unitID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Lease), args.Error(1)
}

func (m *MockLeaseRepo) GetActiveByTenantID(ctx context.Context, tenantID uuid.UUID) (*domain.Lease, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Lease), args.Error(1)
}

func (m *MockLeaseRepo) GetExpiringSoon(ctx context.Context) ([]*domain.Lease, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Lease), args.Error(1)
}

func (m *MockLeaseRepo) Update(ctx context.Context, lease *domain.Lease) error {
	args := m.Called(ctx, lease)
	return args.Error(0)
}

func (m *MockLeaseRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.LeaseStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockLeaseRepo) UpdatePaintingFeePaid(ctx context.Context, id uuid.UUID, paintingFeePaid decimal.Decimal) error {
	args := m.Called(ctx, id, paintingFeePaid)
	return args.Error(0)
}

func (m *MockLeaseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLeaseRepo) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLeaseRepo) CountByStatus(ctx context.Context, status domain.LeaseStatus) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

type MockUnitRepo struct {
	mock.Mock
}

func (m *MockUnitRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Unit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Unit), args.Error(1)
}

func (m *MockUnitRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.UnitStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

// Outros métodos necessários para implementar a interface
func (m *MockUnitRepo) Create(ctx context.Context, unit *domain.Unit) error {
	return nil
}
func (m *MockUnitRepo) GetByNumber(ctx context.Context, number string) (*domain.Unit, error) {
	return nil, nil
}
func (m *MockUnitRepo) List(ctx context.Context) ([]*domain.Unit, error) {
	return nil, nil
}
func (m *MockUnitRepo) ListByStatus(ctx context.Context, status domain.UnitStatus) ([]*domain.Unit, error) {
	return nil, nil
}
func (m *MockUnitRepo) ListByFloor(ctx context.Context, floor int) ([]*domain.Unit, error) {
	return nil, nil
}
func (m *MockUnitRepo) ListAvailable(ctx context.Context) ([]*domain.Unit, error) {
	return nil, nil
}
func (m *MockUnitRepo) Update(ctx context.Context, unit *domain.Unit) error {
	return nil
}
func (m *MockUnitRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *MockUnitRepo) Count(ctx context.Context) (int64, error) {
	return 0, nil
}
func (m *MockUnitRepo) CountByStatus(ctx context.Context, status domain.UnitStatus) (int64, error) {
	return 0, nil
}

type MockTenantRepo struct {
	mock.Mock
}

func (m *MockTenantRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tenant), args.Error(1)
}

// Outros métodos necessários para implementar a interface
func (m *MockTenantRepo) ExistsByCPF(ctx context.Context, cpf string) (bool, error) {
	args := m.Called(ctx, cpf)
	return args.Bool(0), args.Error(1)
}

func (m *MockTenantRepo) SearchByName(ctx context.Context, name string) ([]*domain.Tenant, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Tenant), args.Error(1)
}

func (m *MockTenantRepo) Create(ctx context.Context, tenant *domain.Tenant) error {
	return nil
}
func (m *MockTenantRepo) GetByCPF(ctx context.Context, cpf string) (*domain.Tenant, error) {
	return nil, nil
}
func (m *MockTenantRepo) List(ctx context.Context) ([]*domain.Tenant, error) {
	return nil, nil
}
func (m *MockTenantRepo) Update(ctx context.Context, tenant *domain.Tenant) error {
	return nil
}
func (m *MockTenantRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *MockTenantRepo) Count(ctx context.Context) (int64, error) {
	return 0, nil
}

// Helper functions para criar objetos de teste
func createTestUnit(id uuid.UUID, status domain.UnitStatus) *domain.Unit {
	return &domain.Unit{
		ID:                 id,
		Number:             "101",
		Floor:              1,
		Status:             status,
		IsRenovated:        false,
		BaseRentValue:      decimal.NewFromFloat(800),
		RenovatedRentValue: decimal.NewFromFloat(850),
		CurrentRentValue:   decimal.NewFromFloat(800),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

func createTestTenant(id uuid.UUID) *domain.Tenant {
	return &domain.Tenant{
		ID:        id,
		FullName:  "João Silva",
		CPF:       "123.456.789-00",
		Phone:     "(11) 98765-4321",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// TESTES

func TestCreateLease_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	unitID := uuid.New()
	tenantID := uuid.New()

	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil)

	unit := createTestUnit(unitID, domain.UnitStatusAvailable)
	tenant := createTestTenant(tenantID)

	req := CreateLeaseRequest{
		UnitID:                  unitID,
		TenantID:                tenantID,
		ContractSignedDate:      time.Now(),
		StartDate:               time.Now().AddDate(0, 0, 1),
		PaymentDueDay:           5,
		MonthlyRentValue:        decimal.NewFromFloat(800),
		PaintingFeeTotal:        decimal.NewFromFloat(250),
		PaintingFeeInstallments: 3,
	}

	// Setup mocks
	mockUnitRepo.On("GetByID", ctx, unitID).Return(unit, nil)
	mockLeaseRepo.On("GetActiveByUnitID", ctx, unitID).Return(nil, nil)
	mockTenantRepo.On("GetByID", ctx, tenantID).Return(tenant, nil)
	mockLeaseRepo.On("GetActiveByTenantID", ctx, tenantID).Return(nil, nil)
	mockLeaseRepo.On("Create", ctx, mock.AnythingOfType("*domain.Lease")).Return(nil)
	mockUnitRepo.On("UpdateStatus", ctx, unitID, domain.UnitStatusOccupied).Return(nil)

	// Act
	result, err := service.CreateLease(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Lease)
	assert.Equal(t, unitID, result.Lease.UnitID)
	assert.Equal(t, tenantID, result.Lease.TenantID)
	assert.Equal(t, domain.LeaseStatusActive, result.Lease.Status)
	// Payments vazios pois paymentService é nil
	assert.Empty(t, result.Payments)
	mockLeaseRepo.AssertExpectations(t)
	mockUnitRepo.AssertExpectations(t)
	mockTenantRepo.AssertExpectations(t)
}

func TestCreateLease_UnitNotAvailable(t *testing.T) {
	// Arrange
	ctx := context.Background()
	unitID := uuid.New()
	tenantID := uuid.New()

	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil)

	// Unidade ocupada
	unit := createTestUnit(unitID, domain.UnitStatusOccupied)

	req := CreateLeaseRequest{
		UnitID:                  unitID,
		TenantID:                tenantID,
		ContractSignedDate:      time.Now(),
		StartDate:               time.Now().AddDate(0, 0, 1),
		PaymentDueDay:           5,
		MonthlyRentValue:        decimal.NewFromFloat(800),
		PaintingFeeTotal:        decimal.NewFromFloat(250),
		PaintingFeeInstallments: 3,
	}

	mockUnitRepo.On("GetByID", ctx, unitID).Return(unit, nil)

	// Act
	lease, err := service.CreateLease(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, lease)
	assert.Equal(t, ErrUnitNotAvailable, err)
	mockUnitRepo.AssertExpectations(t)
}

func TestCreateLease_UnitAlreadyHasActiveLease(t *testing.T) {
	// Arrange
	ctx := context.Background()
	unitID := uuid.New()
	tenantID := uuid.New()

	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil)

	unit := createTestUnit(unitID, domain.UnitStatusAvailable)
	existingLease, _ := domain.NewLease(
		unitID,
		uuid.New(),
		time.Now(),
		time.Now(),
		5,
		decimal.NewFromFloat(800),
		decimal.NewFromFloat(250),
		3,
	)

	req := CreateLeaseRequest{
		UnitID:                  unitID,
		TenantID:                tenantID,
		ContractSignedDate:      time.Now(),
		StartDate:               time.Now().AddDate(0, 0, 1),
		PaymentDueDay:           5,
		MonthlyRentValue:        decimal.NewFromFloat(800),
		PaintingFeeTotal:        decimal.NewFromFloat(250),
		PaintingFeeInstallments: 3,
	}

	mockUnitRepo.On("GetByID", ctx, unitID).Return(unit, nil)
	mockLeaseRepo.On("GetActiveByUnitID", ctx, unitID).Return(existingLease, nil)

	// Act
	lease, err := service.CreateLease(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, lease)
	assert.Equal(t, ErrUnitAlreadyHasActiveLease, err)
	mockUnitRepo.AssertExpectations(t)
	mockLeaseRepo.AssertExpectations(t)
}

func TestCancelLease_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	leaseID := uuid.New()
	unitID := uuid.New()
	tenantID := uuid.New()

	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil)

	lease, _ := domain.NewLease(
		unitID,
		tenantID,
		time.Now(),
		time.Now(),
		5,
		decimal.NewFromFloat(800),
		decimal.NewFromFloat(250),
		3,
	)
	lease.ID = leaseID

	mockLeaseRepo.On("GetByID", ctx, leaseID).Return(lease, nil)
	mockLeaseRepo.On("Update", ctx, mock.AnythingOfType("*domain.Lease")).Return(nil)
	mockUnitRepo.On("UpdateStatus", ctx, unitID, domain.UnitStatusAvailable).Return(nil)

	// Act
	err := service.CancelLease(ctx, leaseID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, domain.LeaseStatusCancelled, lease.Status)
	mockLeaseRepo.AssertExpectations(t)
	mockUnitRepo.AssertExpectations(t)
}

func TestCancelLease_AlreadyExpired(t *testing.T) {
	// Arrange
	ctx := context.Background()
	leaseID := uuid.New()

	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil)

	lease, _ := domain.NewLease(
		uuid.New(),
		uuid.New(),
		time.Now(),
		time.Now(),
		5,
		decimal.NewFromFloat(800),
		decimal.NewFromFloat(250),
		3,
	)
	lease.ID = leaseID
	lease.Status = domain.LeaseStatusExpired

	mockLeaseRepo.On("GetByID", ctx, leaseID).Return(lease, nil)

	// Act
	err := service.CancelLease(ctx, leaseID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrLeaseAlreadyExpired, err)
	mockLeaseRepo.AssertExpectations(t)
}

func TestUpdatePaintingFeePaid_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	leaseID := uuid.New()

	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil)

	lease, _ := domain.NewLease(
		uuid.New(),
		uuid.New(),
		time.Now(),
		time.Now(),
		5,
		decimal.NewFromFloat(800),
		decimal.NewFromFloat(250),
		3,
	)
	lease.ID = leaseID

	amountPaid := decimal.NewFromFloat(100)

	mockLeaseRepo.On("GetByID", ctx, leaseID).Return(lease, nil)
	// Use mock.MatchedBy para comparar decimals pelo valor, não pela representação interna
	mockLeaseRepo.On("UpdatePaintingFeePaid", ctx, leaseID, mock.MatchedBy(func(d decimal.Decimal) bool {
		return d.Equal(amountPaid)
	})).Return(nil)

	// Act
	err := service.UpdatePaintingFeePaid(ctx, leaseID, amountPaid)

	// Assert
	assert.NoError(t, err)
	assert.True(t, lease.PaintingFeePaid.Equal(amountPaid))
	mockLeaseRepo.AssertExpectations(t)
}

func TestGetLeaseStats_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil)

	mockLeaseRepo.On("Count", ctx).Return(int64(10), nil)
	mockLeaseRepo.On("CountByStatus", ctx, domain.LeaseStatusActive).Return(int64(7), nil)
	mockLeaseRepo.On("CountByStatus", ctx, domain.LeaseStatusExpiringSoon).Return(int64(2), nil)
	mockLeaseRepo.On("CountByStatus", ctx, domain.LeaseStatusExpired).Return(int64(1), nil)
	mockLeaseRepo.On("CountByStatus", ctx, domain.LeaseStatusCancelled).Return(int64(0), nil)

	// Act
	stats, err := service.GetLeaseStats(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(10), stats.Total)
	assert.Equal(t, int64(7), stats.Active)
	assert.Equal(t, int64(2), stats.ExpiringSoon)
	assert.Equal(t, int64(1), stats.Expired)
	assert.Equal(t, int64(0), stats.Cancelled)
	mockLeaseRepo.AssertExpectations(t)
}

func TestRenewLease_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	oldLeaseID := uuid.New()
	unitID := uuid.New()
	tenantID := uuid.New()

	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil)

	// Contrato antigo que está expirando em breve
	oldLease, _ := domain.NewLease(
		unitID,
		tenantID,
		time.Now().AddDate(0, -6, 0), // Assinado há 6 meses
		time.Now().AddDate(0, -6, 0), // Iniciou há 6 meses
		5,
		decimal.NewFromFloat(800),
		decimal.NewFromFloat(250),
		3,
	)
	oldLease.ID = oldLeaseID
	oldLease.Status = domain.LeaseStatusExpiringSoon // Marcado como expirando

	unit := createTestUnit(unitID, domain.UnitStatusOccupied)

	paintingFeeTotal := decimal.NewFromFloat(250)
	paintingFeeInstallments := 3

	mockLeaseRepo.On("GetByID", ctx, oldLeaseID).Return(oldLease, nil)
	mockUnitRepo.On("GetByID", ctx, unitID).Return(unit, nil)
	mockLeaseRepo.On("Update", ctx, mock.AnythingOfType("*domain.Lease")).Return(nil)
	mockLeaseRepo.On("Create", ctx, mock.AnythingOfType("*domain.Lease")).Return(nil)

	// Act
	newLease, err := service.RenewLease(ctx, oldLeaseID, paintingFeeTotal, paintingFeeInstallments)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, newLease)
	assert.Equal(t, unitID, newLease.UnitID)
	assert.Equal(t, tenantID, newLease.TenantID)
	assert.Equal(t, domain.LeaseStatusActive, newLease.Status)
	assert.Equal(t, domain.LeaseStatusExpired, oldLease.Status) // Antigo marcado como expirado
	assert.True(t, newLease.StartDate.After(oldLease.EndDate))  // Nova data de início após o fim do antigo
	mockLeaseRepo.AssertExpectations(t)
	mockUnitRepo.AssertExpectations(t)
}

func TestRenewLease_CannotRenewCancelled(t *testing.T) {
	// Arrange
	ctx := context.Background()
	oldLeaseID := uuid.New()

	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil)

	// Contrato cancelado
	oldLease, _ := domain.NewLease(
		uuid.New(),
		uuid.New(),
		time.Now(),
		time.Now(),
		5,
		decimal.NewFromFloat(800),
		decimal.NewFromFloat(250),
		3,
	)
	oldLease.ID = oldLeaseID
	oldLease.Status = domain.LeaseStatusCancelled // Cancelado - não pode renovar

	paintingFeeTotal := decimal.NewFromFloat(250)
	paintingFeeInstallments := 3

	mockLeaseRepo.On("GetByID", ctx, oldLeaseID).Return(oldLease, nil)

	// Act
	newLease, err := service.RenewLease(ctx, oldLeaseID, paintingFeeTotal, paintingFeeInstallments)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, newLease)
	assert.Equal(t, ErrCannotRenewLease, err)
	mockLeaseRepo.AssertExpectations(t)
}
