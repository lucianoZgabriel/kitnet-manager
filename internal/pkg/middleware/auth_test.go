package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_Authenticate(t *testing.T) {
	t.Run("should extract token correctly", func(t *testing.T) {
		// Testar apenas a extração de token já que não podemos mockar facilmente o authService
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")

		token := extractToken(req)
		assert.Equal(t, "valid-token", token)
	})

	t.Run("should reject request without token", func(t *testing.T) {
		// Setup - criar um authService dummy
		authMiddleware := &AuthMiddleware{authService: &service.AuthService{}}

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called")
		})

		handler := authMiddleware.Authenticate(testHandler)

		// Criar request sem token
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should reject request with malformed token", func(t *testing.T) {
		// Setup
		authMiddleware := &AuthMiddleware{authService: &service.AuthService{}}

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called")
		})

		handler := authMiddleware.Authenticate(testHandler)

		// Criar request com token malformado
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthMiddleware_RequireRole(t *testing.T) {
	t.Run("should allow access with correct role", func(t *testing.T) {
		// Setup
		authMiddleware := &AuthMiddleware{authService: &service.AuthService{}}

		user, _ := domain.NewUser("admin", "password123", domain.UserRoleAdmin)

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler := authMiddleware.RequireRole(domain.UserRoleAdmin)(testHandler)

		// Criar request com usuário no context
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), UserContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should deny access without correct role", func(t *testing.T) {
		// Setup
		authMiddleware := &AuthMiddleware{authService: &service.AuthService{}}

		user, _ := domain.NewUser("viewer", "password123", domain.UserRoleViewer)

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called")
		})

		handler := authMiddleware.RequireRole(domain.UserRoleAdmin)(testHandler)

		// Criar request com usuário no context
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), UserContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("should deny access without user in context", func(t *testing.T) {
		// Setup
		authMiddleware := &AuthMiddleware{authService: &service.AuthService{}}

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called")
		})

		handler := authMiddleware.RequireRole(domain.UserRoleAdmin)(testHandler)

		// Criar request sem usuário no context
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
	}{
		{
			name:          "valid Bearer token",
			authHeader:    "Bearer valid-token-123",
			expectedToken: "valid-token-123",
		},
		{
			name:          "missing Bearer prefix",
			authHeader:    "valid-token-123",
			expectedToken: "",
		},
		{
			name:          "empty header",
			authHeader:    "",
			expectedToken: "",
		},
		{
			name:          "wrong format",
			authHeader:    "Basic user:pass",
			expectedToken: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			token := extractToken(req)
			assert.Equal(t, tt.expectedToken, token)
		})
	}
}

func TestGetUserFromContext(t *testing.T) {
	t.Run("should retrieve user from context", func(t *testing.T) {
		user, _ := domain.NewUser("testuser", "password123", domain.UserRoleAdmin)

		ctx := context.WithValue(context.Background(), UserContextKey, user)

		retrievedUser, ok := GetUserFromContext(ctx)

		assert.True(t, ok)
		assert.NotNil(t, retrievedUser)
		assert.Equal(t, user.Username, retrievedUser.Username)
	})

	t.Run("should return false when user not in context", func(t *testing.T) {
		ctx := context.Background()

		retrievedUser, ok := GetUserFromContext(ctx)

		assert.False(t, ok)
		assert.Nil(t, retrievedUser)
	})
}

func TestAuthMiddleware_RequireAdmin(t *testing.T) {
	authMiddleware := &AuthMiddleware{authService: &service.AuthService{}}

	adminUser, _ := domain.NewUser("admin", "password123", domain.UserRoleAdmin)
	managerUser, _ := domain.NewUser("manager", "password123", domain.UserRoleManager)

	t.Run("should allow admin", func(t *testing.T) {
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler := authMiddleware.RequireAdmin(testHandler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), UserContextKey, adminUser)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should deny non-admin", func(t *testing.T) {
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called")
		})

		handler := authMiddleware.RequireAdmin(testHandler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), UserContextKey, managerUser)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestAuthMiddleware_RequireAdminOrManager(t *testing.T) {
	authMiddleware := &AuthMiddleware{authService: &service.AuthService{}}

	adminUser, _ := domain.NewUser("admin", "password123", domain.UserRoleAdmin)
	managerUser, _ := domain.NewUser("manager", "password123", domain.UserRoleManager)
	viewerUser, _ := domain.NewUser("viewer", "password123", domain.UserRoleViewer)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("should allow admin", func(t *testing.T) {
		handler := authMiddleware.RequireAdminOrManager(testHandler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), UserContextKey, adminUser)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should allow manager", func(t *testing.T) {
		handler := authMiddleware.RequireAdminOrManager(testHandler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), UserContextKey, managerUser)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should deny viewer", func(t *testing.T) {
		handler := authMiddleware.RequireAdminOrManager(testHandler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), UserContextKey, viewerUser)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
