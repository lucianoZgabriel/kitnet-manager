package handler

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/response"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
)

// ReportHandler lida com requisições HTTP relacionadas a relatórios
type ReportHandler struct {
	reportService *service.ReportService
	validator     *validator.Validate
}

// NewReportHandler cria uma nova instância do handler
func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
		validator:     validator.New(),
	}
}

// GetFinancialReport godoc
// @Summary      Obter relatório financeiro
// @Description  Retorna relatório financeiro consolidado com filtros de período e tipo
// @Tags         Reports
// @Produce      json
// @Param        start_date query string true "Data inicial (YYYY-MM-DD)"
// @Param        end_date query string true "Data final (YYYY-MM-DD)"
// @Param        payment_type query string false "Tipo de pagamento" Enums(rent, painting_fee, adjustment)
// @Param        status query string false "Status do pagamento" Enums(pending, paid, overdue, cancelled)
// @Success      200 {object} service.FinancialReportResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /reports/financial [get]
func (h *ReportHandler) GetFinancialReport(w http.ResponseWriter, r *http.Request) {
	// 1. Extrair query params
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	paymentTypeStr := r.URL.Query().Get("payment_type")
	statusStr := r.URL.Query().Get("status")

	// 2. Validar parâmetros obrigatórios
	if startDateStr == "" || endDateStr == "" {
		response.Error(w, http.StatusBadRequest, "start_date and end_date are required")
		return
	}

	// 3. Parsear datas
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD")
		return
	}

	// 4. Construir request para o service
	req := service.FinancialReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// 5. Aplicar filtros opcionais
	if paymentTypeStr != "" {
		paymentType := domain.PaymentType(paymentTypeStr)
		req.PaymentType = &paymentType
	}

	if statusStr != "" {
		status := domain.PaymentStatus(statusStr)
		req.Status = &status
	}

	// 6. Gerar relatório
	report, err := h.reportService.GetFinancialReport(r.Context(), req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Financial report generated successfully", report)
}

// GetPaymentHistoryReport godoc
// @Summary      Obter histórico de pagamentos
// @Description  Retorna histórico detalhado de pagamentos com filtros
// @Tags         Reports
// @Produce      json
// @Param        lease_id query string false "Filtrar por ID do contrato (UUID)"
// @Param        tenant_id query string false "Filtrar por ID do morador (UUID)"
// @Param        status query string false "Filtrar por status" Enums(pending, paid, overdue, cancelled)
// @Param        start_date query string false "Data inicial (YYYY-MM-DD)"
// @Param        end_date query string false "Data final (YYYY-MM-DD)"
// @Success      200 {object} service.PaymentHistoryResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /reports/payments [get]
func (h *ReportHandler) GetPaymentHistoryReport(w http.ResponseWriter, r *http.Request) {
	// 1. Extrair query params
	leaseIDStr := r.URL.Query().Get("lease_id")
	tenantIDStr := r.URL.Query().Get("tenant_id")
	statusStr := r.URL.Query().Get("status")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// 2. Construir request
	req := service.PaymentHistoryRequest{}

	// 3. Parsear lease_id se fornecido
	if leaseIDStr != "" {
		leaseID, err := uuid.Parse(leaseIDStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid lease_id format")
			return
		}
		req.LeaseID = &leaseID
	}

	// 4. Parsear tenant_id se fornecido
	if tenantIDStr != "" {
		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid tenant_id format")
			return
		}
		req.TenantID = &tenantID
	}

	// 5. Aplicar filtro de status
	if statusStr != "" {
		status := domain.PaymentStatus(statusStr)
		req.Status = &status
	}

	// 6. Parsear datas se fornecidas
	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD")
			return
		}
		req.StartDate = &startDate
	}

	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD")
			return
		}
		req.EndDate = &endDate
	}

	// 7. Gerar relatório
	report, err := h.reportService.GetPaymentHistoryReport(r.Context(), req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Payment history report generated successfully", report)
}

// handleServiceError mapeia erros do service para respostas HTTP
func (h *ReportHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrInvalidDateRange:
		response.Error(w, http.StatusBadRequest, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, "Internal server error")
	}
}
