package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/response"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
)

// TenantHandler lida com requisições HTTP relacionadas a moradores
type TenantHandler struct {
	tenantService *service.TenantService
	validator     *validator.Validate
}

// NewTenantHandler cria uma nova instância do handler
func NewTenantHandler(tenantService *service.TenantService) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
		validator:     validator.New(),
	}
}

// CreateTenant godoc
// @Summary      Criar novo morador
// @Description  Cria um novo morador no sistema
// @Tags         Tenants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant body CreateTenantRequest true "Dados do morador"
// @Success      201 {object} TenantResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      409 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /tenants [post]
func (h *TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	var req CreateTenantRequest

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
	tenant, err := h.tenantService.CreateTenant(
		r.Context(),
		req.FullName,
		req.CPF,
		req.Phone,
		req.Email,
		req.IDDocumentType,
		req.IDDocumentNumber,
	)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// Retornar resposta
	response.Success(w, http.StatusCreated, "Tenant created successfully", ToTenantResponse(tenant))
}

// GetTenant godoc
// @Summary      Buscar morador por ID
// @Description  Retorna os dados de um morador específico
// @Tags         Tenants
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Tenant ID (UUID)"
// @Success      200 {object} TenantResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /tenants/{id} [get]
func (h *TenantHandler) GetTenant(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid tenant ID")
		return
	}

	// Buscar morador
	tenant, err := h.tenantService.GetTenantByID(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Tenant retrieved successfully", ToTenantResponse(tenant))
}

// GetTenantByCPF godoc
// @Summary      Buscar morador por CPF
// @Description  Retorna os dados de um morador pelo CPF
// @Tags         Tenants
// @Produce      json
// @Security     BearerAuth
// @Param        cpf query string true "CPF (formato: XXX.XXX.XXX-XX)"
// @Success      200 {object} TenantResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /tenants/cpf [get]
func (h *TenantHandler) GetTenantByCPF(w http.ResponseWriter, r *http.Request) {
	// Extrair CPF da query string
	cpf := r.URL.Query().Get("cpf")
	if cpf == "" {
		response.Error(w, http.StatusBadRequest, "CPF parameter is required")
		return
	}

	// Buscar morador
	tenant, err := h.tenantService.GetTenantByCPF(r.Context(), cpf)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Tenant retrieved successfully", ToTenantResponse(tenant))
}

// ListTenants godoc
// @Summary      Listar moradores
// @Description  Retorna lista de moradores com opção de busca por nome
// @Tags         Tenants
// @Produce      json
// @Security     BearerAuth
// @Param        name query string false "Search by name (case-insensitive)"
// @Success      200 {array} TenantResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /tenants [get]
func (h *TenantHandler) ListTenants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Verificar filtro de nome na query string
	nameFilter := r.URL.Query().Get("name")

	var tenants []*domain.Tenant
	var err error

	// Aplicar filtro se fornecido
	if nameFilter != "" {
		tenants, err = h.tenantService.SearchTenantsByName(ctx, nameFilter)
	} else {
		tenants, err = h.tenantService.ListTenants(ctx)
	}

	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Tenants retrieved successfully", ToTenantResponseList(tenants))
}

// UpdateTenant godoc
// @Summary      Atualizar morador
// @Description  Atualiza os dados de um morador existente
// @Tags         Tenants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Tenant ID (UUID)"
// @Param        tenant body UpdateTenantRequest true "Dados atualizados"
// @Success      200 {object} TenantResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /tenants/{id} [put]
func (h *TenantHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	// Extrair ID
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid tenant ID")
		return
	}

	// Decodificar request
	var req UpdateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validar
	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Atualizar
	tenant, err := h.tenantService.UpdateTenant(
		r.Context(),
		id,
		req.FullName,
		req.Phone,
		req.Email,
		req.IDDocumentType,
		req.IDDocumentNumber,
	)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Tenant updated successfully", ToTenantResponse(tenant))
}

// DeleteTenant godoc
// @Summary      Deletar morador
// @Description  Remove um morador do sistema
// @Tags         Tenants
// @Security     BearerAuth
// @Param        id path string true "Tenant ID (UUID)"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /tenants/{id} [delete]
func (h *TenantHandler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	// Extrair ID
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid tenant ID")
		return
	}

	// Deletar
	if err := h.tenantService.DeleteTenant(r.Context(), id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Tenant deleted successfully", nil)
}

// handleServiceError mapeia erros do service para respostas HTTP
func (h *TenantHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrTenantNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrCPFAlreadyExists):
		response.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrInvalidFullName),
		errors.Is(err, domain.ErrInvalidCPF),
		errors.Is(err, domain.ErrInvalidPhone),
		errors.Is(err, domain.ErrInvalidEmail):
		response.Error(w, http.StatusBadRequest, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, "Internal server error")
	}
}
