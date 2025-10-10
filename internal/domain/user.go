package domain

import (
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRole representa os papéis de usuários no sistema
type UserRole string

const (
	UserRoleAdmin   UserRole = "admin"
	UserRoleManager UserRole = "manager"
	UserRoleViewer  UserRole = "viewer"
)

// ValidRoles contém todos os roles válidos
var ValidRoles = []UserRole{
	UserRoleAdmin,
	UserRoleManager,
	UserRoleViewer,
}

// User representa um usuário do sistema
type User struct {
	ID           uuid.UUID  `json:"id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"-"` // nunca expor na API
	Role         UserRole   `json:"role"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Domain errors
var (
	ErrInvalidUsername       = errors.New("username must be at least 3 characters")
	ErrInvalidPassword       = errors.New("password must be at least 6 characters")
	ErrInvalidRole           = errors.New("invalid user role")
	ErrPasswordHashEmpty     = errors.New("password hash cannot be empty")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrUserInactive          = errors.New("user account is inactive")
)

// NewUser cria um novo usuário com password em texto plano (será hasheado)
func NewUser(username, password string, role UserRole) (*User, error) {
	user := &User{
		ID:        uuid.New(),
		Username:  strings.TrimSpace(strings.ToLower(username)),
		Role:      role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Hash da senha
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	// Valida o usuário
	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate verifica se o usuário possui dados válidos
func (u *User) Validate() error {
	// Valida username
	if len(strings.TrimSpace(u.Username)) < 3 {
		return ErrInvalidUsername
	}

	// Valida password hash
	if strings.TrimSpace(u.PasswordHash) == "" {
		return ErrPasswordHashEmpty
	}

	// Valida role
	if !u.IsValidRole() {
		return ErrInvalidRole
	}

	return nil
}

// IsValidRole verifica se o role do usuário é válido
func (u *User) IsValidRole() bool {
	return slices.Contains(ValidRoles, u.Role)
}

// SetPassword hasheia e define a senha do usuário
func (u *User) SetPassword(password string) error {
	// Valida tamanho mínimo
	if len(password) < 6 {
		return ErrInvalidPassword
	}

	// Gera hash bcrypt (cost 10)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)
	u.UpdatedAt = time.Now()
	return nil
}

// ValidatePassword compara a senha fornecida com o hash armazenado
func (u *User) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return ErrInvalidCredentials
	}
	return nil
}

// UpdateLastLogin atualiza o timestamp do último login
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}

// Deactivate desativa o usuário
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate ativa o usuário
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// ChangeRole altera o papel do usuário
func (u *User) ChangeRole(newRole UserRole) error {
	// Valida se o novo role é válido
	tempUser := &User{Role: newRole}
	if !tempUser.IsValidRole() {
		return ErrInvalidRole
	}

	u.Role = newRole
	u.UpdatedAt = time.Now()
	return nil
}

// IsAdmin verifica se o usuário é admin
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// IsManager verifica se o usuário é manager
func (u *User) IsManager() bool {
	return u.Role == UserRoleManager
}

// IsViewer verifica se o usuário é viewer
func (u *User) IsViewer() bool {
	return u.Role == UserRoleViewer
}

// CanManageUsers verifica se o usuário pode gerenciar outros usuários
func (u *User) CanManageUsers() bool {
	return u.IsAdmin()
}

// CanWrite verifica se o usuário tem permissão de escrita
func (u *User) CanWrite() bool {
	return u.IsAdmin() || u.IsManager()
}

// CanRead verifica se o usuário tem permissão de leitura
func (u *User) CanRead() bool {
	return u.IsActive // todos os usuários ativos podem ler
}

// String retorna uma representação em string do usuário
func (u *User) String() string {
	return u.Username + " (" + string(u.Role) + ")"
}
