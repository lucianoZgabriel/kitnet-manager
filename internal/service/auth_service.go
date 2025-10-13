package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository"
)

// Auth service errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenExpired       = errors.New("token has expired")
	ErrInvalidTokenClaims = errors.New("invalid token claims")
)

// JWTClaims representa as claims customizadas do JWT
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthService contém a lógica de negócio para autenticação e autorização
type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
	jwtExpiry time.Duration
}

// NewAuthService cria uma nova instância do serviço de autenticação
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
		jwtExpiry: jwtExpiry,
	}
}

// Login autentica um usuário e retorna um token JWT
func (s *AuthService) Login(ctx context.Context, username, password string) (string, *domain.User, error) {
	// Buscar usuário pelo username
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return "", nil, ErrInvalidCredentials
	}

	// Verificar se usuário está ativo
	if !user.IsActive {
		return "", nil, ErrUserInactive
	}

	// Validar senha
	if err := user.ValidatePassword(password); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Atualizar último login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID, time.Now()); err != nil {
		// Log o erro mas não falhe o login por isso
		fmt.Printf("Warning: failed to update last login: %v\n", err)
	}

	// Gerar token JWT
	token, err := s.GenerateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("error generating token: %w", err)
	}

	return token, user, nil
}

// GenerateToken gera um token JWT para o usuário
func (s *AuthService) GenerateToken(user *domain.User) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.jwtExpiry)

	claims := JWTClaims{
		UserID:   user.ID.String(),
		Username: user.Username,
		Role:     string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "kitnet-manager",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken valida um token JWT e retorna as claims
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		// Verificar método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidTokenClaims
}

// GetUserFromToken busca o usuário completo a partir de um token
// DEPRECATED: Use GetUserFromTokenClaims para melhor performance
// Esta função faz query no banco a cada chamada, causando problemas de connection pool
func (s *AuthService) GetUserFromToken(ctx context.Context, tokenString string) (*domain.User, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, ErrInvalidTokenClaims
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Verificar se usuário ainda está ativo
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	return user, nil
}

// GetUserFromTokenClaims reconstrói o usuário apenas dos claims do JWT, sem query no banco
// RECOMENDADO: Use esta função para autenticação em requests (melhor performance)
// Nota: Não verifica se usuário foi desativado após emissão do token
// Para verificações críticas que precisam de dados atualizados, use GetUserFromToken
func (s *AuthService) GetUserFromTokenClaims(ctx context.Context, tokenString string) (*domain.User, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, ErrInvalidTokenClaims
	}

	// Reconstruir usuário a partir dos claims
	// Não faz query no banco, evitando overhead e connection pool exhaustion
	user := &domain.User{
		ID:       userID,
		Username: claims.Username,
		Role:     domain.UserRole(claims.Role),
		IsActive: true, // Assumimos ativo se token é válido
	}

	return user, nil
}

// CreateUser cria um novo usuário (apenas admin pode fazer isso)
func (s *AuthService) CreateUser(ctx context.Context, username, password string, role domain.UserRole) (*domain.User, error) {
	// Verificar se username já existe
	exists, err := s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("error checking username: %w", err)
	}
	if exists {
		return nil, ErrUsernameExists
	}

	// Criar usuário usando domain model
	user, err := domain.NewUser(username, password, role)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Persistir no banco
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("error saving user: %w", err)
	}

	return user, nil
}

// GetUserByID busca um usuário pelo ID
func (s *AuthService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// ListUsers retorna todos os usuários
func (s *AuthService) ListUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}
	return users, nil
}

// ChangePassword altera a senha de um usuário
func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	// Buscar usuário
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Validar senha antiga
	if err := user.ValidatePassword(oldPassword); err != nil {
		return ErrInvalidCredentials
	}

	// Definir nova senha
	if err := user.SetPassword(newPassword); err != nil {
		return fmt.Errorf("error setting new password: %w", err)
	}

	// Atualizar no banco
	if err := s.userRepo.UpdatePassword(ctx, userID, user.PasswordHash); err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}

	return nil
}

// ChangeUserRole altera o papel de um usuário (apenas admin)
func (s *AuthService) ChangeUserRole(ctx context.Context, userID uuid.UUID, newRole domain.UserRole) error {
	// Buscar usuário
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Alterar role usando método do domain
	if err := user.ChangeRole(newRole); err != nil {
		return fmt.Errorf("error changing role: %w", err)
	}

	// Persistir mudança
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

// DeactivateUser desativa um usuário (apenas admin)
func (s *AuthService) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.userRepo.Deactivate(ctx, userID); err != nil {
		return fmt.Errorf("error deactivating user: %w", err)
	}
	return nil
}

// ActivateUser ativa um usuário (apenas admin)
func (s *AuthService) ActivateUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.userRepo.Activate(ctx, userID); err != nil {
		return fmt.Errorf("error activating user: %w", err)
	}
	return nil
}

// RefreshToken gera um novo token a partir de um token válido (opcional)
func (s *AuthService) RefreshToken(ctx context.Context, oldToken string) (string, error) {
	// Validar token atual
	claims, err := s.ValidateToken(oldToken)
	if err != nil {
		// Se o token expirou mas é válido em outros aspectos, permitir refresh
		if !errors.Is(err, ErrTokenExpired) {
			return "", err
		}
	}

	// Buscar usuário
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return "", ErrInvalidTokenClaims
	}

	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	// Verificar se ainda está ativo
	if !user.IsActive {
		return "", ErrUserInactive
	}

	// Gerar novo token
	newToken, err := s.GenerateToken(user)
	if err != nil {
		return "", fmt.Errorf("error generating new token: %w", err)
	}

	return newToken, nil
}
