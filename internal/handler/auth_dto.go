package handler

import (
	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
)

// LoginRequest representa o payload para login
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse representa a resposta do login
type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

// RefreshTokenRequest representa o payload para refresh de token
type RefreshTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// RefreshTokenResponse representa a resposta do refresh
type RefreshTokenResponse struct {
	Token string `json:"token"`
}

// CreateUserRequest representa o payload para criar um usuário
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=admin manager viewer"`
}

// ChangePasswordRequest representa o payload para trocar senha
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=6"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// ChangeRoleRequest representa o payload para trocar role
type ChangeRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin manager viewer"`
}

// UserResponse representa a resposta com dados de um usuário
type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Role        string    `json:"role"`
	IsActive    bool      `json:"is_active"`
	LastLoginAt *string   `json:"last_login_at,omitempty"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

// ToUserResponse converte domain.User para UserResponse
func ToUserResponse(user *domain.User) *UserResponse {
	response := &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if user.LastLoginAt != nil {
		lastLogin := user.LastLoginAt.Format("2006-01-02T15:04:05Z07:00")
		response.LastLoginAt = &lastLogin
	}

	return response
}

// ToUserResponseList converte slice de users para slice de responses
func ToUserResponseList(users []*domain.User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}
