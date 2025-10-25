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

func (m *MockLeaseRepo) UpdateAndCreateAtomic(ctx context.Context, oldLease, newLease *domain.Lease, adjustment *domain.LeaseRentAdjustment) error {
	args := m.Called(ctx, oldLease, newLease, adjustment)
	return args.Error(0)
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
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Unit), args.Error(1)
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

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil, nil)

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

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil, nil)

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

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil, nil)

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

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil, nil)

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

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil, nil)

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

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil, nil)

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

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil, nil)

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

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil, nil)

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

	req := RenewLeaseRequest{
		PaintingFeeTotal:        paintingFeeTotal,
		PaintingFeeInstallments: paintingFeeInstallments,
	}

	mockLeaseRepo.On("GetByID", ctx, oldLeaseID).Return(oldLease, nil)
	mockUnitRepo.On("GetByID", ctx, unitID).Return(unit, nil)
	mockLeaseRepo.On("UpdateAndCreateAtomic", ctx, mock.AnythingOfType("*domain.Lease"), mock.AnythingOfType("*domain.Lease"), mock.AnythingOfType("*domain.LeaseRentAdjustment")).Return(nil)

	// Act
	result, err := service.RenewLease(ctx, oldLeaseID, req, nil)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Lease)
	assert.Equal(t, unitID, result.Lease.UnitID)
	assert.Equal(t, tenantID, result.Lease.TenantID)
	assert.Equal(t, domain.LeaseStatusActive, result.Lease.Status)
	assert.Equal(t, domain.LeaseStatusExpired, oldLease.Status) // Antigo marcado como expirado
	assert.True(t, result.Lease.StartDate.After(oldLease.EndDate))  // Nova data de início após o fim do antigo
	// Payments vazios pois paymentService é nil
	assert.Empty(t, result.Payments)
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

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil, nil)

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

	req := RenewLeaseRequest{
		PaintingFeeTotal:        paintingFeeTotal,
		PaintingFeeInstallments: paintingFeeInstallments,
	}

	mockLeaseRepo.On("GetByID", ctx, oldLeaseID).Return(oldLease, nil)

	// Act
	result, err := service.RenewLease(ctx, oldLeaseID, req, nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrCannotRenewLease, err)
	mockLeaseRepo.AssertExpectations(t)
}

func TestRenewLease_WithRentAdjustment(t *testing.T) {
	// Arrange
	ctx := context.Background()
	oldLeaseID := uuid.New()
	unitID := uuid.New()
	tenantID := uuid.New()

	mockLeaseRepo := new(MockLeaseRepo)
	mockUnitRepo := new(MockUnitRepo)
	mockTenantRepo := new(MockTenantRepo)
	mockAdjustmentRepo := new(MockLeaseRentAdjustmentRepo)

	service := NewLeaseService(mockLeaseRepo, mockUnitRepo, mockTenantRepo, nil, mockAdjustmentRepo)

	// Contrato original (generation 1)
	oldLease, _ := domain.NewLease(
		unitID,
		tenantID,
		time.Now().AddDate(0, -6, 0),
		time.Now().AddDate(0, -6, 0),
		5,
		decimal.NewFromFloat(800),
		decimal.NewFromFloat(250),
		3,
	)
	oldLease.ID = oldLeaseID
	oldLease.Generation = 1
	oldLease.Status = domain.LeaseStatusExpiringSoon

	unit := createTestUnit(unitID, domain.UnitStatusOccupied)

	newRentValue := decimal.NewFromFloat(880)
	reason := "Reajuste anual IGPM"

	req := RenewLeaseRequest{
		PaintingFeeTotal:        decimal.NewFromFloat(250),
		PaintingFeeInstallments: 3,
		NewRentValue:            &newRentValue,
		AdjustmentReason:        &reason,
	}

	mockLeaseRepo.On("GetByID", ctx, oldLeaseID).Return(oldLease, nil)
	mockUnitRepo.On("GetByID", ctx, unitID).Return(unit, nil)
	mockLeaseRepo.On("UpdateAndCreateAtomic", ctx, mock.AnythingOfType("*domain.Lease"), mock.AnythingOfType("*domain.Lease"), mock.AnythingOfType("*domain.LeaseRentAdjustment")).Return(nil)
	// Nota: adjustment agora é criado dentro da transação UpdateAndCreateAtomic

	// Act
	result, err := service.RenewLease(ctx, oldLeaseID, req, nil)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Lease)
	assert.Equal(t, 2, result.Lease.Generation) // Segunda geração
	assert.Equal(t, &oldLeaseID, result.Lease.ParentLeaseID)
	assert.True(t, result.Lease.MonthlyRentValue.Equal(newRentValue))
	mockLeaseRepo.AssertExpectations(t)
	mockUnitRepo.AssertExpectations(t)
	// mockAdjustmentRepo não é mais usado diretamente, adjustment é criado via UpdateAndCreateAtomic
}

func TestLeaseDomain_ShouldApplyAnnualAdjustment(t *testing.T) {
	// Arrange
	unitID := uuid.New()
	tenantID := uuid.New()

	tests := []struct {
		name       string
		generation int
		expected   bool
	}{
		{
			name:       "Generation 1 (original) - no adjustment",
			generation: 1,
			expected:   false,
		},
		{
			name:       "Generation 2 (1st renewal) - apply adjustment",
			generation: 2,
			expected:   true,
		},
		{
			name:       "Generation 3 (2nd renewal) - no adjustment",
			generation: 3,
			expected:   false,
		},
		{
			name:       "Generation 4 (3rd renewal) - apply adjustment",
			generation: 4,
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar contrato
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
			lease.Generation = tt.generation

			// Act & Assert
			assert.Equal(t, tt.expected, lease.ShouldApplyAnnualAdjustment())
		})
	}
}

func TestLeaseDomain_GetTotalMonths(t *testing.T) {
	// Arrange
	unitID := uuid.New()
	tenantID := uuid.New()

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

	tests := []struct {
		generation    int
		expectedMonths int
	}{
		{1, 6},   // 1ª geração = 6 meses
		{2, 12},  // 2ª geração = 12 meses
		{3, 18},  // 3ª geração = 18 meses
		{4, 24},  // 4ª geração = 24 meses
	}

	for _, tt := range tests {
		t.Run("Generation "+string(rune(tt.generation+'0')), func(t *testing.T) {
			lease.Generation = tt.generation
			assert.Equal(t, tt.expectedMonths, lease.GetTotalMonths())
		})
	}
}

func TestLeaseRentAdjustment_PercentageCalculation(t *testing.T) {
	// Arrange
	leaseID := uuid.New()
	previousValue := decimal.NewFromFloat(800)
	newValue := decimal.NewFromFloat(880)
	reason := "Reajuste anual IGPM"

	// Act
	adjustment := domain.NewLeaseRentAdjustment(leaseID, previousValue, newValue, &reason, nil)

	// Assert
	assert.NotNil(t, adjustment)
	assert.Equal(t, leaseID, adjustment.LeaseID)
	assert.True(t, adjustment.PreviousRentValue.Equal(previousValue))
	assert.True(t, adjustment.NewRentValue.Equal(newValue))
	// Percentual esperado: (880-800)/800 * 100 = 10%
	expectedPercentage := decimal.NewFromFloat(10)
	assert.True(t, adjustment.AdjustmentPercentage.Equal(expectedPercentage))
	assert.Equal(t, &reason, adjustment.Reason)
}

// Mock para LeaseRentAdjustmentRepository
type MockLeaseRentAdjustmentRepo struct {
	mock.Mock
}

func (m *MockLeaseRentAdjustmentRepo) Create(ctx context.Context, adjustment *domain.LeaseRentAdjustment) error {
	args := m.Called(ctx, adjustment)
	return args.Error(0)
}

func (m *MockLeaseRentAdjustmentRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.LeaseRentAdjustment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LeaseRentAdjustment), args.Error(1)
}

func (m *MockLeaseRentAdjustmentRepo) ListByLeaseID(ctx context.Context, leaseID uuid.UUID) ([]*domain.LeaseRentAdjustment, error) {
	args := m.Called(ctx, leaseID)
	return args.Get(0).([]*domain.LeaseRentAdjustment), args.Error(1)
}

func (m *MockLeaseRentAdjustmentRepo) GetLatestByLeaseID(ctx context.Context, leaseID uuid.UUID) (*domain.LeaseRentAdjustment, error) {
	args := m.Called(ctx, leaseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LeaseRentAdjustment), args.Error(1)
}

func (m *MockLeaseRentAdjustmentRepo) CountByLeaseID(ctx context.Context, leaseID uuid.UUID) (int64, error) {
	args := m.Called(ctx, leaseID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLeaseRentAdjustmentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
