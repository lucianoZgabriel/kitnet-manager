package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserRepository é um mock do repository para testes
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) ListByRole(ctx context.Context, role domain.UserRole) ([]*domain.User, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	args := m.Called(ctx, id, passwordHash)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, lastLogin time.Time) error {
	args := m.Called(ctx, id, lastLogin)
	return args.Error(0)
}

func (m *MockUserRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Activate(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) CountActive(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// Tests

func TestAuthService_CreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should create user successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)

		// Mock: username não existe
		mockRepo.On("ExistsByUsername", ctx, "testuser").Return(false, nil)
		// Mock: criação bem-sucedida
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

		user, err := service.CreateUser(ctx, "testuser", "password123", domain.UserRoleAdmin)

		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, domain.UserRoleAdmin, user.Role)
		assert.True(t, user.IsActive)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when username exists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)

		mockRepo.On("ExistsByUsername", ctx, "existing").Return(true, nil)

		user, err := service.CreateUser(ctx, "existing", "password123", domain.UserRoleAdmin)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, ErrUsernameExists)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()

	t.Run("should login successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)

		// Criar usuário de teste
		user, _ := domain.NewUser("testuser", "password123", domain.UserRoleAdmin)

		mockRepo.On("GetByUsername", ctx, "testuser").Return(user, nil)
		mockRepo.On("UpdateLastLogin", ctx, user.ID, mock.Anything).Return(nil)

		token, returnedUser, err := service.Login(ctx, "testuser", "password123")

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotNil(t, returnedUser)
		assert.Equal(t, user.Username, returnedUser.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail with wrong password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)

		user, _ := domain.NewUser("testuser", "password123", domain.UserRoleAdmin)

		mockRepo.On("GetByUsername", ctx, "testuser").Return(user, nil)

		token, returnedUser, err := service.Login(ctx, "testuser", "wrongpassword")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, returnedUser)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when user is inactive", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)

		user, _ := domain.NewUser("testuser", "password123", domain.UserRoleAdmin)
		user.Deactivate()

		mockRepo.On("GetByUsername", ctx, "testuser").Return(user, nil)

		token, returnedUser, err := service.Login(ctx, "testuser", "password123")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, returnedUser)
		assert.ErrorIs(t, err, ErrUserInactive)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)

		mockRepo.On("GetByUsername", ctx, "nonexistent").Return(nil, nil)

		token, returnedUser, err := service.Login(ctx, "nonexistent", "password123")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, returnedUser)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_GenerateAndValidateToken(t *testing.T) {
	t.Run("should generate and validate token successfully", func(t *testing.T) {
		service := NewAuthService(nil, "test-secret", 24*time.Hour)

		user, _ := domain.NewUser("testuser", "password123", domain.UserRoleAdmin)

		// Gerar token
		token, err := service.GenerateToken(user)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Validar token
		claims, err := service.ValidateToken(token)
		require.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, user.ID.String(), claims.UserID)
		assert.Equal(t, user.Username, claims.Username)
		assert.Equal(t, string(user.Role), claims.Role)
	})

	t.Run("should fail to validate invalid token", func(t *testing.T) {
		service := NewAuthService(nil, "test-secret", 24*time.Hour)

		claims, err := service.ValidateToken("invalid-token")

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("should fail to validate expired token", func(t *testing.T) {
		service := NewAuthService(nil, "test-secret", -1*time.Hour) // Token já expira

		user, _ := domain.NewUser("testuser", "password123", domain.UserRoleAdmin)
		token, _ := service.GenerateToken(user)

		claims, err := service.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.ErrorIs(t, err, ErrTokenExpired)
	})
}

func TestAuthService_ChangePassword(t *testing.T) {
	ctx := context.Background()

	t.Run("should change password successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)

		user, _ := domain.NewUser("testuser", "oldpassword", domain.UserRoleAdmin)

		mockRepo.On("GetByID", ctx, user.ID).Return(user, nil)
		mockRepo.On("UpdatePassword", ctx, user.ID, mock.AnythingOfType("string")).Return(nil)

		err := service.ChangePassword(ctx, user.ID, "oldpassword", "newpassword123")

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail with wrong old password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)

		user, _ := domain.NewUser("testuser", "oldpassword", domain.UserRoleAdmin)

		mockRepo.On("GetByID", ctx, user.ID).Return(user, nil)

		err := service.ChangePassword(ctx, user.ID, "wrongpassword", "newpassword123")

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_ChangeUserRole(t *testing.T) {
	ctx := context.Background()

	t.Run("should change role successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)

		user, _ := domain.NewUser("testuser", "password123", domain.UserRoleViewer)

		mockRepo.On("GetByID", ctx, user.ID).Return(user, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

		err := service.ChangeUserRole(ctx, user.ID, domain.UserRoleManager)

		require.NoError(t, err)
		assert.Equal(t, domain.UserRoleManager, user.Role)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_ActivateDeactivate(t *testing.T) {
	ctx := context.Background()

	t.Run("should deactivate user successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)

		userID := uuid.New()
		mockRepo.On("Deactivate", ctx, userID).Return(nil)

		err := service.DeactivateUser(ctx, userID)

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should activate user successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)

		userID := uuid.New()
		mockRepo.On("Activate", ctx, userID).Return(nil)

		err := service.ActivateUser(ctx, userID)

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
