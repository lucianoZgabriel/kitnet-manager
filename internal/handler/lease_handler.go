package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/response"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
)

// LeaseHandler lida com requisições HTTP relacionadas a contratos
type LeaseHandler struct {
	leaseService *service.LeaseService
	validator    *validator.Validate
}

// NewLeaseHandler cria uma nova instância do handler
func NewLeaseHandler(leaseService *service.LeaseService) *LeaseHandler {
	return &LeaseHandler{
		leaseService: leaseService,
		validator:    validator.New(),
	}
}

// CreateLease godoc
// @Summary      Criar novo contrato
// @Description  Cria um novo contrato de locação
// @Tags         Leases
// @Accept       json
// @Produce      json
// @Param        lease body CreateLeaseRequestDTO true "Dados do contrato"
// @Success      201 {object} LeaseResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /leases [post]
func (h *LeaseHandler) CreateLease(w http.ResponseWriter, r *http.Request) {
	var req CreateLeaseRequestDTO

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

	// Converter DTO para service request
	serviceReq := service.CreateLeaseRequest{
		UnitID:                  req.UnitID,
		TenantID:                req.TenantID,
		ContractSignedDate:      req.ContractSignedDate,
		StartDate:               req.StartDate,
		PaymentDueDay:           req.PaymentDueDay,
		MonthlyRentValue:        req.MonthlyRentValue,
		PaintingFeeTotal:        req.PaintingFeeTotal,
		PaintingFeeInstallments: req.PaintingFeeInstallments,
	}

	// Chamar service
	result, err := h.leaseService.CreateLease(r.Context(), serviceReq)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// Retornar resposta com lease e pagamentos gerados
	responseData := map[string]interface{}{
		"lease":    ToLeaseResponse(result.Lease),
		"payments": ToPaymentResponseList(result.Payments),
	}
	response.Success(w, http.StatusCreated, "Lease created successfully with payments", responseData)
}

// GetLease godoc
// @Summary      Buscar contrato por ID
// @Description  Retorna os dados de um contrato específico
// @Tags         Leases
// @Produce      json
// @Param        id path string true "Lease ID (UUID)"
// @Success      200 {object} LeaseResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /leases/{id} [get]
func (h *LeaseHandler) GetLease(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid lease ID")
		return
	}

	// Buscar contrato
	lease, err := h.leaseService.GetLeaseByID(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Lease retrieved successfully", ToLeaseResponse(lease))
}

// ListLeases godoc
// @Summary      Listar contratos
// @Description  Retorna lista de contratos com filtros opcionais
// @Tags         Leases
// @Produce      json
// @Param        status query string false "Filter by status" Enums(active, expiring_soon, expired, cancelled)
// @Param        unit_id query string false "Filter by unit ID (UUID)"
// @Param        tenant_id query string false "Filter by tenant ID (UUID)"
// @Success      200 {array} LeaseResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /leases [get]
func (h *LeaseHandler) ListLeases(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Verificar filtros na query string
	statusFilter := r.URL.Query().Get("status")
	unitIDFilter := r.URL.Query().Get("unit_id")
	tenantIDFilter := r.URL.Query().Get("tenant_id")

	var leases []*domain.Lease
	var err error

	// Aplicar filtros
	if statusFilter != "" {
		status := domain.LeaseStatus(statusFilter)
		leases, err = h.leaseService.ListLeasesByStatus(ctx, status)
	} else if unitIDFilter != "" {
		unitID, parseErr := uuid.Parse(unitIDFilter)
		if parseErr != nil {
			response.Error(w, http.StatusBadRequest, "Invalid unit_id parameter")
			return
		}
		leases, err = h.leaseService.ListLeasesByUnitID(ctx, unitID)
	} else if tenantIDFilter != "" {
		tenantID, parseErr := uuid.Parse(tenantIDFilter)
		if parseErr != nil {
			response.Error(w, http.StatusBadRequest, "Invalid tenant_id parameter")
			return
		}
		leases, err = h.leaseService.ListLeasesByTenantID(ctx, tenantID)
	} else {
		leases, err = h.leaseService.ListLeases(ctx)
	}

	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Leases retrieved successfully", ToLeaseResponseList(leases))
}

// GetLeaseStats godoc
// @Summary      Estatísticas de contratos
// @Description  Retorna estatísticas agregadas dos contratos
// @Tags         Leases
// @Produce      json
// @Success      200 {object} service.LeaseStats
// @Failure      500 {object} response.ErrorResponse
// @Router       /leases/stats [get]
func (h *LeaseHandler) GetLeaseStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.leaseService.GetLeaseStats(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Lease stats retrieved successfully", stats)
}

// RenewLease godoc
// @Summary      Renovar contrato
// @Description  Renova um contrato existente criando um novo contrato
// @Tags         Leases
// @Accept       json
// @Produce      json
// @Param        id path string true "Lease ID (UUID)"
// @Param        renewal body RenewLeaseRequestDTO true "Dados da renovação"
// @Success      201 {object} LeaseResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /leases/{id}/renew [post]
func (h *LeaseHandler) RenewLease(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid lease ID")
		return
	}

	// Decodificar request
	var req RenewLeaseRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validar
	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Renovar contrato
	newLease, err := h.leaseService.RenewLease(
		r.Context(),
		id,
		req.PaintingFeeTotal,
		req.PaintingFeeInstallments,
	)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "Lease renewed successfully", ToLeaseResponse(newLease))
}

// CancelLease godoc
// @Summary      Cancelar contrato
// @Description  Cancela um contrato e libera a unidade
// @Tags         Leases
// @Produce      json
// @Param        id path string true "Lease ID (UUID)"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /leases/{id}/cancel [post]
func (h *LeaseHandler) CancelLease(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid lease ID")
		return
	}

	// Cancelar contrato
	if err := h.leaseService.CancelLease(r.Context(), id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Lease cancelled successfully", nil)
}

// UpdatePaintingFeePaid godoc
// @Summary      Atualizar taxa de pintura paga
// @Description  Registra pagamento da taxa de pintura
// @Tags         Leases
// @Accept       json
// @Produce      json
// @Param        id path string true "Lease ID (UUID)"
// @Param        payment body UpdatePaintingFeePaidRequestDTO true "Valor pago"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /leases/{id}/painting-fee [patch]
func (h *LeaseHandler) UpdatePaintingFeePaid(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid lease ID")
		return
	}

	// Decodificar request
	var req UpdatePaintingFeePaidRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validar
	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Atualizar valor pago
	if err := h.leaseService.UpdatePaintingFeePaid(r.Context(), id, req.AmountPaid); err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Painting fee updated successfully", nil)
}

// GetExpiringSoonLeases godoc
// @Summary      Listar contratos expirando em breve
// @Description  Retorna contratos que expiram nos próximos 45 dias
// @Tags         Leases
// @Produce      json
// @Success      200 {array} LeaseResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /leases/expiring-soon [get]
func (h *LeaseHandler) GetExpiringSoonLeases(w http.ResponseWriter, r *http.Request) {
	leases, err := h.leaseService.GetExpiringSoonLeases(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Expiring soon leases retrieved successfully", ToLeaseResponseList(leases))
}

// handleServiceError mapeia erros do service para respostas HTTP
func (h *LeaseHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrLeaseNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrUnitAlreadyHasActiveLease),
		errors.Is(err, service.ErrTenantAlreadyHasActiveLease),
		errors.Is(err, service.ErrUnitNotAvailable),
		errors.Is(err, service.ErrCannotCancelLease),
		errors.Is(err, service.ErrCannotRenewLease),
		errors.Is(err, service.ErrLeaseAlreadyExpired):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrUnitNotFound),
		errors.Is(err, service.ErrTenantNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidLeaseStatus),
		errors.Is(err, domain.ErrInvalidPaymentDueDay),
		errors.Is(err, domain.ErrInvalidPaintingFeeInstallments),
		errors.Is(err, domain.ErrInvalidMonthlyRentValue),
		errors.Is(err, domain.ErrInvalidDates),
		errors.Is(err, domain.ErrPaintingFeePaidExceedsTotal):
		response.Error(w, http.StatusBadRequest, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, "Internal server error")
	}
}

// PaymentResponse representa um pagamento na resposta HTTP
type PaymentResponse struct {
	ID             string  `json:"id"`
	LeaseID        string  `json:"lease_id"`
	PaymentType    string  `json:"payment_type"`
	ReferenceMonth string  `json:"reference_month"`
	Amount         float64 `json:"amount"`
	Status         string  `json:"status"`
	DueDate        string  `json:"due_date"`
	PaymentDate    *string `json:"payment_date,omitempty"`
	PaymentMethod  *string `json:"payment_method,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// ToPaymentResponse converte domain.Payment para PaymentResponse
func ToPaymentResponse(p *domain.Payment) *PaymentResponse {
	if p == nil {
		return nil
	}

	amount, _ := p.Amount.Float64()
	
	var paymentDate *string
	if p.PaymentDate != nil {
		dateStr := p.PaymentDate.Format("2006-01-02")
		paymentDate = &dateStr
	}

	var paymentMethod *string
	if p.PaymentMethod != nil {
		methodStr := string(*p.PaymentMethod)
		paymentMethod = &methodStr
	}

	return &PaymentResponse{
		ID:             p.ID.String(),
		LeaseID:        p.LeaseID.String(),
		PaymentType:    string(p.PaymentType),
		ReferenceMonth: p.ReferenceMonth.Format("2006-01-02"),
		Amount:         amount,
		Status:         string(p.Status),
		DueDate:        p.DueDate.Format("2006-01-02"),
		PaymentDate:    paymentDate,
		PaymentMethod:  paymentMethod,
		CreatedAt:      p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      p.UpdatedAt.Format(time.RFC3339),
	}
}

// ToPaymentResponseList converte slice de payments para slice de responses
func ToPaymentResponseList(payments []*domain.Payment) []*PaymentResponse {
	if payments == nil {
		return []*PaymentResponse{}
	}

	responses := make([]*PaymentResponse, len(payments))
	for i, payment := range payments {
		responses[i] = ToPaymentResponse(payment)
	}
	return responses
}
