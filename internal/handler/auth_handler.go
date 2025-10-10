package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/response"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
)

// AuthHandler lida com requisições HTTP relacionadas a autenticação
type AuthHandler struct {
	authService *service.AuthService
	validator   *validator.Validate
}

// NewAuthHandler cria uma nova instância do handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

// Login godoc
// @Summary      Login
// @Description  Autentica um usuário e retorna um token JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials body LoginRequest true "Credenciais de login"
// @Success      200 {object} LoginResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validar request
	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Chamar service
	token, user, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// Retornar resposta
	loginResponse := &LoginResponse{
		Token: token,
		User:  ToUserResponse(user),
	}

	response.Success(w, http.StatusOK, "Login successful", loginResponse)
}

// GetCurrentUser godoc
// @Summary      Obter usuário atual
// @Description  Retorna os dados do usuário autenticado pelo token
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} UserResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /auth/me [get]
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Extrair token do header
	token := h.extractToken(r)
	if token == "" {
		response.Error(w, http.StatusUnauthorized, "Missing authorization token")
		return
	}

	// Buscar usuário pelo token
	user, err := h.authService.GetUserFromToken(r.Context(), token)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "User retrieved successfully", ToUserResponse(user))
}

// RefreshToken godoc
// @Summary      Renovar token
// @Description  Gera um novo token JWT a partir de um token válido ou expirado
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        token body RefreshTokenRequest true "Token a ser renovado"
// @Success      200 {object} RefreshTokenResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest

	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validar request
	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Chamar service
	newToken, err := h.authService.RefreshToken(r.Context(), req.Token)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// Retornar resposta
	refreshResponse := &RefreshTokenResponse{
		Token: newToken,
	}

	response.Success(w, http.StatusOK, "Token refreshed successfully", refreshResponse)
}

// CreateUser godoc
// @Summary      Criar usuário
// @Description  Cria um novo usuário no sistema (apenas admin)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user body CreateUserRequest true "Dados do usuário"
// @Success      201 {object} UserResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      409 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /auth/users [post]
func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validar request
	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Converter role string para domain.UserRole
	role := domain.UserRole(req.Role)

	// Chamar service
	user, err := h.authService.CreateUser(r.Context(), req.Username, req.Password, role)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// Retornar resposta
	response.Success(w, http.StatusCreated, "User created successfully", ToUserResponse(user))
}

// ListUsers godoc
// @Summary      Listar usuários
// @Description  Retorna lista de todos os usuários (apenas admin)
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} UserResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /auth/users [get]
func (h *AuthHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.authService.ListUsers(r.Context())
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Users retrieved successfully", ToUserResponseList(users))
}

// GetUser godoc
// @Summary      Buscar usuário por ID
// @Description  Retorna os dados de um usuário específico (apenas admin)
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID (UUID)"
// @Success      200 {object} UserResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /auth/users/{id} [get]
func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Buscar usuário
	user, err := h.authService.GetUserByID(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "User retrieved successfully", ToUserResponse(user))
}

// ChangePassword godoc
// @Summary      Trocar senha
// @Description  Altera a senha do usuário autenticado
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        passwords body ChangePasswordRequest true "Senhas antiga e nova"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Router       /auth/change-password [post]
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req ChangePasswordRequest

	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validar request
	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Extrair token e buscar usuário
	token := h.extractToken(r)
	if token == "" {
		response.Error(w, http.StatusUnauthorized, "Missing authorization token")
		return
	}

	user, err := h.authService.GetUserFromToken(r.Context(), token)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// Chamar service
	if err := h.authService.ChangePassword(r.Context(), user.ID, req.OldPassword, req.NewPassword); err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Password changed successfully", nil)
}

// ChangeUserRole godoc
// @Summary      Alterar role de usuário
// @Description  Altera o papel/permissão de um usuário (apenas admin)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID (UUID)"
// @Param        role body ChangeRoleRequest true "Novo role"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /auth/users/{id}/role [patch]
func (h *AuthHandler) ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req ChangeRoleRequest

	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validar request
	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Converter role string para domain.UserRole
	role := domain.UserRole(req.Role)

	// Chamar service
	if err := h.authService.ChangeUserRole(r.Context(), id, role); err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "User role changed successfully", nil)
}

// DeactivateUser godoc
// @Summary      Desativar usuário
// @Description  Desativa um usuário (apenas admin)
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID (UUID)"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /auth/users/{id}/deactivate [post]
func (h *AuthHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Chamar service
	if err := h.authService.DeactivateUser(r.Context(), id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "User deactivated successfully", nil)
}

// ActivateUser godoc
// @Summary      Ativar usuário
// @Description  Ativa um usuário previamente desativado (apenas admin)
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID (UUID)"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /auth/users/{id}/activate [post]
func (h *AuthHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Chamar service
	if err := h.authService.ActivateUser(r.Context(), id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "User activated successfully", nil)
}

// extractToken extrai o token JWT do header Authorization
func (h *AuthHandler) extractToken(r *http.Request) string {
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

// handleServiceError mapeia erros do service para respostas HTTP
func (h *AuthHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrUsernameExists):
		response.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrInvalidCredentials):
		response.Error(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, service.ErrUserInactive):
		response.Error(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, service.ErrInvalidToken),
		errors.Is(err, service.ErrTokenExpired),
		errors.Is(err, service.ErrInvalidTokenClaims):
		response.Error(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrInvalidUsername),
		errors.Is(err, domain.ErrInvalidPassword),
		errors.Is(err, domain.ErrInvalidRole):
		response.Error(w, http.StatusBadRequest, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, "Internal server error")
	}
}
