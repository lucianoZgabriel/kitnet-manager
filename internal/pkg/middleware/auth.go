package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/response"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
)

// contextKey é um tipo customizado para keys do context
type contextKey string

const (
	// UserContextKey é a key usada para armazenar o usuário no context
	UserContextKey contextKey = "user"
)

// AuthMiddleware é o middleware de autenticação JWT
type AuthMiddleware struct {
	authService *service.AuthService
}

// NewAuthMiddleware cria uma nova instância do middleware de autenticação
func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Authenticate verifica se há um token JWT válido no header Authorization
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extrair token do header Authorization
		token := extractToken(r)
		if token == "" {
			response.Error(w, http.StatusUnauthorized, "Missing authorization token")
			return
		}

		// Validar token e obter usuário
		user, err := m.authService.GetUserFromToken(r.Context(), token)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Adicionar usuário ao context
		ctx := context.WithValue(r.Context(), UserContextKey, user)

		// Chamar próximo handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole retorna um middleware que verifica se o usuário tem um dos roles especificados
func (m *AuthMiddleware) RequireRole(roles ...domain.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obter usuário do context
			user, ok := r.Context().Value(UserContextKey).(*domain.User)
			if !ok || user == nil {
				response.Error(w, http.StatusUnauthorized, "User not authenticated")
				return
			}

			// Verificar se o usuário tem um dos roles permitidos
			hasPermission := false
			for _, role := range roles {
				if user.Role == role {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				response.Error(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			// Chamar próximo handler
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin é um atalho para RequireRole(UserRoleAdmin)
func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return m.RequireRole(domain.UserRoleAdmin)(next)
}

// RequireAdminOrManager é um atalho para RequireRole(UserRoleAdmin, UserRoleManager)
func (m *AuthMiddleware) RequireAdminOrManager(next http.Handler) http.Handler {
	return m.RequireRole(domain.UserRoleAdmin, domain.UserRoleManager)(next)
}

// extractToken extrai o token JWT do header Authorization
func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}

	// Formato esperado: "Bearer <token>"
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// GetUserFromContext retorna o usuário armazenado no context
// Esta função auxiliar pode ser usada em handlers para obter o usuário autenticado
func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*domain.User)
	return user, ok
}
