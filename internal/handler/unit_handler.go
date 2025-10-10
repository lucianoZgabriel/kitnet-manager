package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/response"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
)

// UnitHandler lida com requisições HTTP relacionadas a unidades
type UnitHandler struct {
	unitService *service.UnitService
	validator   *validator.Validate
}

// NewUnitHandler cria uma nova instância do handler
func NewUnitHandler(unitService *service.UnitService) *UnitHandler {
	return &UnitHandler{
		unitService: unitService,
		validator:   validator.New(),
	}
}

// CreateUnit godoc
// @Summary      Criar nova unidade
// @Description  Cria uma nova unidade no sistema
// @Tags         Units
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        unit body CreateUnitRequest true "Dados da unidade"
// @Success      201 {object} UnitResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /units [post]
func (h *UnitHandler) CreateUnit(w http.ResponseWriter, r *http.Request) {
	var req CreateUnitRequest

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
	unit, err := h.unitService.CreateUnit(
		r.Context(),
		req.Number,
		req.Floor,
		req.BaseRentValue,
		req.RenovatedRentValue,
	)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// Retornar resposta
	response.Success(w, http.StatusCreated, "Unit created successfully", ToUnitResponse(unit))
}

// GetUnit godoc
// @Summary      Buscar unidade por ID
// @Description  Retorna os dados de uma unidade específica
// @Tags         Units
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Unit ID (UUID)"
// @Success      200 {object} UnitResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /units/{id} [get]
func (h *UnitHandler) GetUnit(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid unit ID")
		return
	}

	// Buscar unidade
	unit, err := h.unitService.GetUnitByID(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Unit retrieved successfully", ToUnitResponse(unit))
}

// ListUnits godoc
// @Summary      Listar unidades
// @Description  Retorna lista de unidades com filtros opcionais
// @Tags         Units
// @Produce      json
// @Security     BearerAuth
// @Param        status query string false "Filter by status" Enums(available, occupied, maintenance, renovation)
// @Param        floor query int false "Filter by floor"
// @Success      200 {array} UnitResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /units [get]
func (h *UnitHandler) ListUnits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Verificar filtros na query string
	statusFilter := r.URL.Query().Get("status")
	floorFilter := r.URL.Query().Get("floor")

	var units []*domain.Unit
	var err error

	// Aplicar filtros
	if statusFilter != "" {
		status := domain.UnitStatus(statusFilter)
		units, err = h.unitService.ListUnitsByStatus(ctx, status)
	} else if floorFilter != "" {
		var floor int
		if _, parseErr := fmt.Sscanf(floorFilter, "%d", &floor); parseErr != nil {
			response.Error(w, http.StatusBadRequest, "Invalid floor parameter")
			return
		}
		units, err = h.unitService.ListUnitsByFloor(ctx, floor)
	} else {
		units, err = h.unitService.ListUnits(ctx)
	}

	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Units retrieved successfully", ToUnitResponseList(units))
}

// UpdateUnit godoc
// @Summary      Atualizar unidade
// @Description  Atualiza os dados de uma unidade existente
// @Tags         Units
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Unit ID (UUID)"
// @Param        unit body UpdateUnitRequest true "Dados atualizados"
// @Success      200 {object} UnitResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /units/{id} [put]
func (h *UnitHandler) UpdateUnit(w http.ResponseWriter, r *http.Request) {
	// Extrair ID
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid unit ID")
		return
	}

	// Decodificar request
	var req UpdateUnitRequest
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
	unit, err := h.unitService.UpdateUnit(
		r.Context(),
		id,
		req.Number,
		req.Floor,
		req.IsRenovated,
		req.BaseRentValue,
		req.RenovatedRentValue,
		req.Notes,
	)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Unit updated successfully", ToUnitResponse(unit))
}

// UpdateUnitStatus godoc
// @Summary      Atualizar status da unidade
// @Description  Atualiza apenas o status de uma unidade
// @Tags         Units
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Unit ID (UUID)"
// @Param        status body UpdateUnitStatusRequest true "Novo status"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Router       /units/{id}/status [patch]
func (h *UnitHandler) UpdateUnitStatus(w http.ResponseWriter, r *http.Request) {
	// Extrair ID
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid unit ID")
		return
	}

	// Decodificar request
	var req UpdateUnitStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validar
	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Atualizar status
	if err := h.unitService.UpdateUnitStatus(r.Context(), id, domain.UnitStatus(req.Status)); err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Unit status updated successfully", nil)
}

// DeleteUnit godoc
// @Summary      Deletar unidade
// @Description  Remove uma unidade do sistema
// @Tags         Units
// @Security     BearerAuth
// @Param        id path string true "Unit ID (UUID)"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /units/{id} [delete]
func (h *UnitHandler) DeleteUnit(w http.ResponseWriter, r *http.Request) {
	// Extrair ID
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid unit ID")
		return
	}

	// Deletar
	if err := h.unitService.DeleteUnit(r.Context(), id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Unit deleted successfully", nil)
}

// GetOccupancyStats godoc
// @Summary      Estatísticas de ocupação
// @Description  Retorna estatísticas de ocupação das unidades
// @Tags         Units
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} service.OccupancyStats
// @Failure      500 {object} response.ErrorResponse
// @Router       /units/stats/occupancy [get]
func (h *UnitHandler) GetOccupancyStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.unitService.GetOccupancyStats(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Occupancy stats retrieved successfully", stats)
}

// handleServiceError mapeia erros do service para respostas HTTP
func (h *UnitHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrUnitNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrUnitNumberAlreadyExists):
		response.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrCannotDeleteOccupiedUnit):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidUnitNumber),
		errors.Is(err, domain.ErrInvalidFloor),
		errors.Is(err, domain.ErrInvalidRentValue),
		errors.Is(err, domain.ErrInvalidStatus):
		response.Error(w, http.StatusBadRequest, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, "Internal server error")
	}
}
