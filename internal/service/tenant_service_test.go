package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockTenantRepository é um mock do repository para testes
type MockTenantRepository struct {
	mock.Mock
}

func (m *MockTenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tenant), args.Error(1)
}

func (m *MockTenantRepository) GetByCPF(ctx context.Context, cpf string) (*domain.Tenant, error) {
	args := m.Called(ctx, cpf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tenant), args.Error(1)
}

func (m *MockTenantRepository) List(ctx context.Context) ([]*domain.Tenant, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Tenant), args.Error(1)
}

func (m *MockTenantRepository) SearchByName(ctx context.Context, name string) ([]*domain.Tenant, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Tenant), args.Error(1)
}

func (m *MockTenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTenantRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTenantRepository) ExistsByCPF(ctx context.Context, cpf string) (bool, error) {
	args := m.Called(ctx, cpf)
	return args.Get(0).(bool), args.Error(1)
}

// Tests

func TestTenantService_CreateTenant(t *testing.T) {
	ctx := context.Background()

	t.Run("should create tenant successfully", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		// Mock: CPF não existe
		mockRepo.On("ExistsByCPF", ctx, "123.456.789-00").Return(false, nil)
		// Mock: criação bem-sucedida
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Tenant")).Return(nil)

		tenant, err := service.CreateTenant(ctx, "João da Silva", "123.456.789-00", "11987654321", "joao@example.com", "RG", "123456")

		require.NoError(t, err)
		assert.NotNil(t, tenant)
		assert.Equal(t, "João da Silva", tenant.FullName)
		assert.Equal(t, "123.456.789-00", tenant.CPF)
		assert.Equal(t, "RG", tenant.IDDocumentType)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should create tenant without optional fields", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		mockRepo.On("ExistsByCPF", ctx, "987.654.321-00").Return(false, nil)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Tenant")).Return(nil)

		tenant, err := service.CreateTenant(ctx, "Maria Santos", "987.654.321-00", "11912345678", "", "", "")

		require.NoError(t, err)
		assert.NotNil(t, tenant)
		assert.Equal(t, "", tenant.Email)
		assert.Equal(t, "", tenant.IDDocumentType)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when CPF already exists", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		mockRepo.On("ExistsByCPF", ctx, "123.456.789-00").Return(true, nil)

		tenant, err := service.CreateTenant(ctx, "João da Silva", "123.456.789-00", "11987654321", "", "", "")

		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Equal(t, ErrCPFAlreadyExists, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail with invalid CPF format", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		mockRepo.On("ExistsByCPF", ctx, "12345678900").Return(false, nil)

		tenant, err := service.CreateTenant(ctx, "João da Silva", "12345678900", "11987654321", "", "", "")

		assert.Error(t, err)
		assert.Nil(t, tenant)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail with empty name", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		mockRepo.On("ExistsByCPF", ctx, "123.456.789-00").Return(false, nil)

		tenant, err := service.CreateTenant(ctx, "", "123.456.789-00", "11987654321", "", "", "")

		assert.Error(t, err)
		assert.Nil(t, tenant)
		mockRepo.AssertExpectations(t)
	})
}

func TestTenantService_GetTenantByID(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("should get tenant by id", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		expectedTenant, _ := domain.NewTenant("João da Silva", "123.456.789-00", "11987654321", "")
		expectedTenant.ID = tenantID

		mockRepo.On("GetByID", ctx, tenantID).Return(expectedTenant, nil)

		tenant, err := service.GetTenantByID(ctx, tenantID)

		require.NoError(t, err)
		assert.Equal(t, expectedTenant, tenant)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when tenant not found", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		mockRepo.On("GetByID", ctx, tenantID).Return(nil, nil)

		tenant, err := service.GetTenantByID(ctx, tenantID)

		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Equal(t, ErrTenantNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestTenantService_GetTenantByCPF(t *testing.T) {
	ctx := context.Background()

	t.Run("should get tenant by CPF", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		expectedTenant, _ := domain.NewTenant("João da Silva", "123.456.789-00", "11987654321", "")

		mockRepo.On("GetByCPF", ctx, "123.456.789-00").Return(expectedTenant, nil)

		tenant, err := service.GetTenantByCPF(ctx, "123.456.789-00")

		require.NoError(t, err)
		assert.Equal(t, expectedTenant, tenant)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when CPF not found", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		mockRepo.On("GetByCPF", ctx, "999.999.999-99").Return(nil, nil)

		tenant, err := service.GetTenantByCPF(ctx, "999.999.999-99")

		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Equal(t, ErrTenantNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestTenantService_ListTenants(t *testing.T) {
	ctx := context.Background()

	t.Run("should list all tenants", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		tenant1, _ := domain.NewTenant("João da Silva", "123.456.789-00", "11987654321", "")
		tenant2, _ := domain.NewTenant("Maria Santos", "987.654.321-00", "11912345678", "")
		expectedTenants := []*domain.Tenant{tenant1, tenant2}

		mockRepo.On("List", ctx).Return(expectedTenants, nil)

		tenants, err := service.ListTenants(ctx)

		require.NoError(t, err)
		assert.Len(t, tenants, 2)
		assert.Equal(t, expectedTenants, tenants)
		mockRepo.AssertExpectations(t)
	})
}

func TestTenantService_SearchTenantsByName(t *testing.T) {
	ctx := context.Background()

	t.Run("should search tenants by name", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		tenant1, _ := domain.NewTenant("João da Silva", "123.456.789-00", "11987654321", "")
		expectedTenants := []*domain.Tenant{tenant1}

		mockRepo.On("SearchByName", ctx, "João").Return(expectedTenants, nil)

		tenants, err := service.SearchTenantsByName(ctx, "João")

		require.NoError(t, err)
		assert.Len(t, tenants, 1)
		assert.Equal(t, "João da Silva", tenants[0].FullName)
		mockRepo.AssertExpectations(t)
	})
}

func TestTenantService_UpdateTenant(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("should update tenant successfully", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		existingTenant, _ := domain.NewTenant("João da Silva", "123.456.789-00", "11987654321", "joao@example.com")
		existingTenant.ID = tenantID

		mockRepo.On("GetByID", ctx, tenantID).Return(existingTenant, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Tenant")).Return(nil)

		tenant, err := service.UpdateTenant(ctx, tenantID, "João Silva Junior", "11999887766", "joao.junior@example.com", "CNH", "987654321")

		require.NoError(t, err)
		assert.Equal(t, "João Silva Junior", tenant.FullName)
		assert.Equal(t, "11999887766", tenant.Phone)
		assert.Equal(t, "CNH", tenant.IDDocumentType)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when tenant not found", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		mockRepo.On("GetByID", ctx, tenantID).Return(nil, nil)

		tenant, err := service.UpdateTenant(ctx, tenantID, "New Name", "11999999999", "", "", "")

		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Equal(t, ErrTenantNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestTenantService_DeleteTenant(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("should delete tenant successfully", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		existingTenant, _ := domain.NewTenant("João da Silva", "123.456.789-00", "11987654321", "")
		existingTenant.ID = tenantID

		mockRepo.On("GetByID", ctx, tenantID).Return(existingTenant, nil)
		mockRepo.On("Delete", ctx, tenantID).Return(nil)

		err := service.DeleteTenant(ctx, tenantID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when tenant not found", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		mockRepo.On("GetByID", ctx, tenantID).Return(nil, nil)

		err := service.DeleteTenant(ctx, tenantID)

		assert.Error(t, err)
		assert.Equal(t, ErrTenantNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestTenantService_GetTenantCount(t *testing.T) {
	ctx := context.Background()

	t.Run("should return tenant count", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantService(mockRepo)

		mockRepo.On("Count", ctx).Return(int64(15), nil)

		count, err := service.GetTenantCount(ctx)

		require.NoError(t, err)
		assert.Equal(t, int64(15), count)
		mockRepo.AssertExpectations(t)
	})
}
