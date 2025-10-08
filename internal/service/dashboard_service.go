package service

import (
	"context"

	"github.com/lucianoZgabriel/kitnet-manager/internal/repository"
	"github.com/shopspring/decimal"
)

// DashboardService contém a lógica de negócio
type DashboardService struct {
	dashboardRepo repository.DashboardRepository
	leaseRepo     repository.LeaseRepository
	paymentRepo   repository.PaymentRepository
}

// NewDashboardService cria uma nova instância do serviço de dashboard
func NewDashboardService(
	dashboardRepo repository.DashboardRepository,
	leaseRepo repository.LeaseRepository,
	paymentRepo repository.PaymentRepository,
) *DashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
		leaseRepo:     leaseRepo,
		paymentRepo:   paymentRepo,
	}
}

// OccupancyMetrics representa as métricas de ocupação
type OccupancyMetrics struct {
	TotalUnits       int64   `json:"total_units"`
	OccupiedUnits    int64   `json:"occupied_units"`
	AvailableUnits   int64   `json:"available_units"`
	MaintenanceUnits int64   `json:"maintenance_units"`
	RenovationUnits  int64   `json:"renovation_units"`
	OccupancyRate    float64 `json:"occupancy_rate"`    // Percentual
	AvailabilityRate float64 `json:"availability_rate"` // Percentual
}

// FinancialMetrics representa as métricas financeiras
type FinancialMetrics struct {
	MonthlyProjectedRevenue decimal.Decimal `json:"monthly_projected_revenue"`
	MonthlyRealizedRevenue  decimal.Decimal `json:"monthly_realized_revenue"`
	OverdueAmount           decimal.Decimal `json:"overdue_amount"`
	TotalPendingAmount      decimal.Decimal `json:"total_pending_amount"`
	DefaultRate             float64         `json:"default_rate"`    // Taxa de inadimplência
	CollectionRate          float64         `json:"collection_rate"` // Taxa de cobrança efetiva
}

// GetOccupancyMetrics retorna as métricas de ocupação das unidades
func (s *DashboardService) GetOccupancyMetrics(ctx context.Context) (*OccupancyMetrics, error) {
	metrics, err := s.dashboardRepo.GetOccupancyMetrics(ctx)
	if err != nil {
		return nil, err
	}

	occupancyRate := 0.0
	availabilityRate := 0.0

	if metrics.TotalUnits > 0 {
		occupancyRate = (float64(metrics.OccupiedUnits) / float64(metrics.TotalUnits)) * 100
		availabilityRate = (float64(metrics.AvailableUnits) / float64(metrics.TotalUnits)) * 100
	}

	return &OccupancyMetrics{
		TotalUnits:       metrics.TotalUnits,
		OccupiedUnits:    metrics.OccupiedUnits,
		AvailableUnits:   metrics.AvailableUnits,
		MaintenanceUnits: metrics.MaintenanceUnits,
		RenovationUnits:  metrics.RenovationUnits,
		OccupancyRate:    occupancyRate,
		AvailabilityRate: availabilityRate,
	}, nil
}

// GetFinancialMetrics retorna as métricas financeiras
func (s *DashboardService) GetFinancialMetrics(ctx context.Context) (*FinancialMetrics, error) {
	// 1. Buscar receita projetada (soma de todos os aluguéis ativos)
	projectedRevenue, err := s.dashboardRepo.GetMonthlyProjectedRevenue(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Buscar receita realizada (pagamentos pagos no mês atual)
	realizedRevenue, err := s.dashboardRepo.GetMonthlyRealizedRevenue(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Buscar valor em atraso (pagamentos overdue)
	overdueAmount, err := s.dashboardRepo.GetOverdueAmount(ctx)
	if err != nil {
		return nil, err
	}

	// 4. Buscar total pendente (pending + overdue)
	pendingAmount, err := s.dashboardRepo.GetTotalPendingAmount(ctx)
	if err != nil {
		return nil, err
	}

	// 5. Calcular taxa de inadimplência (overdue / projected * 100)
	defaultRate := 0.0
	if !projectedRevenue.IsZero() {
		defaultRate, _ = overdueAmount.Div(projectedRevenue).Mul(decimal.NewFromInt(100)).Float64()
	}

	// 6. Calcular taxa de cobrança (realized / projected * 100)
	collectionRate := 0.0
	if !projectedRevenue.IsZero() {
		collectionRate, _ = realizedRevenue.Div(projectedRevenue).Mul(decimal.NewFromInt(100)).Float64()
	}

	return &FinancialMetrics{
		MonthlyProjectedRevenue: projectedRevenue,
		MonthlyRealizedRevenue:  realizedRevenue,
		OverdueAmount:           overdueAmount,
		TotalPendingAmount:      pendingAmount,
		DefaultRate:             defaultRate,
		CollectionRate:          collectionRate,
	}, nil
}
