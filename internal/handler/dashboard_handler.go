package handler

import (
	"net/http"

	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/response"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
)

// DashboardHandler lida com requisições HTTP relacionadas ao dashboard
type DashboardHandler struct {
	dashboardService *service.DashboardService
}

// NewDashboardHandler cria uma nova instância do handler
func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetDashboard godoc
// @Summary      Obter dados do dashboard
// @Description  Retorna métricas consolidadas de ocupação, financeiras, contratos e alertas
// @Tags         Dashboard
// @Produce      json
// @Success      200 {object} DashboardResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /dashboard [get]
func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Buscar métricas de ocupação
	occupancyMetrics, err := h.dashboardService.GetOccupancyMetrics(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve occupancy metrics")
		return
	}

	// 2. Buscar métricas financeiras
	financialMetrics, err := h.dashboardService.GetFinancialMetrics(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve financial metrics")
		return
	}

	// 3. Buscar métricas de contratos
	contractMetrics, err := h.dashboardService.GetContractMetrics(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve contract metrics")
		return
	}

	// 4. Buscar alertas
	alerts, err := h.dashboardService.GetAlerts(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve alerts")
		return
	}

	// 5. Montar resposta consolidada
	dashboardData := DashboardResponse{
		Occupancy: occupancyMetrics,
		Financial: financialMetrics,
		Contracts: contractMetrics,
		Alerts:    alerts,
	}

	response.Success(w, http.StatusOK, "Dashboard data retrieved successfully", dashboardData)
}
