package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/response"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
)

// PaymentHandler lida com requisições HTTP relacionadas a pagamentos
type PaymentHandler struct {
	paymentService *service.PaymentService
	validator      *validator.Validate
}

// NewPaymentHandler cria uma nova instância do handler
func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		validator:      validator.New(),
	}
}

// GetPayment godoc
// @Summary      Buscar pagamento por ID
// @Description  Retorna os dados de um pagamento específico
// @Tags         Payments
// @Produce      json
// @Param        id path string true "Payment ID (UUID)"
// @Success      200 {object} PaymentResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /payments/{id} [get]
func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	// Buscar pagamento
	payment, err := h.paymentService.GetPaymentByID(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Payment retrieved successfully", ToPaymentResponse(payment))
}

// GetPaymentsByLease godoc
// @Summary      Listar pagamentos de um contrato
// @Description  Retorna todos os pagamentos de um contrato específico
// @Tags         Payments
// @Produce      json
// @Param        lease_id path string true "Lease ID (UUID)"
// @Success      200 {array} PaymentResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /leases/{lease_id}/payments [get]
func (h *PaymentHandler) GetPaymentsByLease(w http.ResponseWriter, r *http.Request) {
	// Extrair lease_id da URL
	leaseIDStr := chi.URLParam(r, "lease_id")
	leaseID, err := uuid.Parse(leaseIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid lease ID")
		return
	}

	// Buscar pagamentos do contrato
	payments, err := h.paymentService.GetPaymentsByLease(r.Context(), leaseID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Payments retrieved successfully", ToPaymentResponseList(payments))
}

// GetCancellablePayments godoc
// @Summary      Listar pagamentos canceláveis de um contrato
// @Description  Retorna pagamentos com status pending ou overdue que podem ser cancelados
// @Tags         Payments
// @Produce      json
// @Param        lease_id path string true "Lease ID (UUID)"
// @Success      200 {array} PaymentResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /leases/{lease_id}/cancellable-payments [get]
func (h *PaymentHandler) GetCancellablePayments(w http.ResponseWriter, r *http.Request) {
	// Extrair lease_id da URL
	leaseIDStr := chi.URLParam(r, "lease_id")
	leaseID, err := uuid.Parse(leaseIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid lease ID")
		return
	}

	// Buscar pagamentos canceláveis do contrato
	payments, err := h.paymentService.GetCancellablePayments(r.Context(), leaseID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Cancellable payments retrieved successfully", ToPaymentResponseList(payments))
}

// GetOverduePayments godoc
// @Summary      Listar pagamentos atrasados
// @Description  Retorna todos os pagamentos com status overdue
// @Tags         Payments
// @Produce      json
// @Success      200 {array} PaymentResponse
// @Failure      500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /payments/overdue [get]
func (h *PaymentHandler) GetOverduePayments(w http.ResponseWriter, r *http.Request) {
	payments, err := h.paymentService.GetOverduePayments(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Overdue payments retrieved successfully", ToPaymentResponseList(payments))
}

// GetUpcomingPayments godoc
// @Summary      Listar pagamentos próximos ao vencimento
// @Description  Retorna pagamentos que vencem nos próximos X dias (padrão: 7 dias)
// @Tags         Payments
// @Produce      json
// @Param        days query int false "Número de dias à frente (padrão: 7)"
// @Success      200 {array} PaymentResponse
// @Failure      500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /payments/upcoming [get]
func (h *PaymentHandler) GetUpcomingPayments(w http.ResponseWriter, r *http.Request) {
	// Obter parâmetro 'days' da query string (padrão: 7)
	daysStr := r.URL.Query().Get("days")
	days := 7 // valor padrão

	if daysStr != "" {
		parsedDays, err := strconv.Atoi(daysStr)
		if err != nil || parsedDays <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid days parameter")
			return
		}
		days = parsedDays
	}

	payments, err := h.paymentService.GetUpcomingPayments(r.Context(), days)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Upcoming payments retrieved successfully", ToPaymentResponseList(payments))
}

// MarkPaymentAsPaid godoc
// @Summary      Marcar pagamento como pago
// @Description  Registra um pagamento como pago com data e forma de pagamento
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id path string true "Payment ID (UUID)"
// @Param        payment body MarkPaymentAsPaidRequestDTO true "Dados do pagamento"
// @Success      200 {object} PaymentResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /payments/{id}/pay [put]
func (h *PaymentHandler) MarkPaymentAsPaid(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	// Decodificar request
	var req MarkPaymentAsPaidRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validar
	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Marcar como pago
	serviceReq := service.MarkPaymentAsPaidRequest{
		PaymentID:     id,
		PaymentDate:   req.PaymentDate,
		PaymentMethod: req.PaymentMethod,
	}

	updatedPayment, err := h.paymentService.MarkPaymentAsPaid(r.Context(), serviceReq)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Payment marked as paid successfully", ToPaymentResponse(updatedPayment))
}

// CancelPayment godoc
// @Summary      Cancelar pagamento
// @Description  Cancela um pagamento pendente
// @Tags         Payments
// @Produce      json
// @Param        id path string true "Payment ID (UUID)"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /payments/{id}/cancel [post]
func (h *PaymentHandler) CancelPayment(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	// Cancelar pagamento
	if err := h.paymentService.CancelPayment(r.Context(), id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Payment cancelled successfully", nil)
}

// GetPaymentStatsByLease godoc
// @Summary      Estatísticas de pagamentos de um contrato
// @Description  Retorna estatísticas agregadas dos pagamentos de um contrato
// @Tags         Payments
// @Produce      json
// @Param        lease_id path string true "Lease ID (UUID)"
// @Success      200 {object} PaymentStatsResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /leases/{lease_id}/payments/stats [get]
func (h *PaymentHandler) GetPaymentStatsByLease(w http.ResponseWriter, r *http.Request) {
	// Extrair lease_id da URL
	leaseIDStr := chi.URLParam(r, "lease_id")
	leaseID, err := uuid.Parse(leaseIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid lease ID")
		return
	}

	// Buscar estatísticas
	stats, err := h.paymentService.GetPaymentStatsByLease(r.Context(), leaseID)
	if err != nil {
		log.Printf("ERROR GetPaymentStatsByLease: %v", err)
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Payment stats retrieved successfully", ToPaymentStatsResponse(stats))
}

// handleServiceError mapeia erros do service para respostas HTTP
func (h *PaymentHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrPaymentNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrLeaseNotFoundForPayment):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrInvalidPaymentAmount),
		errors.Is(err, service.ErrInvalidInstallments),
		errors.Is(err, service.ErrPaymentCannotBePaid),
		errors.Is(err, service.ErrPaymentAlreadyPaid):
		response.Error(w, http.StatusBadRequest, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, "Internal server error")
	}
}
