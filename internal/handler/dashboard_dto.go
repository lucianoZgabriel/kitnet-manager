package handler

import "github.com/lucianoZgabriel/kitnet-manager/internal/service"

// DashboardResponse representa a resposta consolidada do dashboard
type DashboardResponse struct {
	Occupancy *service.OccupancyMetrics `json:"occupancy"`
	Financial *service.FinancialMetrics `json:"financial"`
	Contracts *service.ContractMetrics  `json:"contracts"`
	Alerts    *service.DashboardAlerts  `json:"alerts"`
}
