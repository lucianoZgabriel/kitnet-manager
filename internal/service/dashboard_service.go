package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository"
	"github.com/shopspring/decimal"
)

// DashboardService contém a lógica de negócio
type DashboardService struct {
	dashboardRepo repository.DashboardRepository
	leaseRepo     repository.LeaseRepository
	paymentRepo   repository.PaymentRepository
	unitRepo      repository.UnitRepository
}

// NewDashboardService cria uma nova instância do serviço de dashboard
func NewDashboardService(
	dashboardRepo repository.DashboardRepository,
	leaseRepo repository.LeaseRepository,
	paymentRepo repository.PaymentRepository,
	unitRepo repository.UnitRepository,
) *DashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
		leaseRepo:     leaseRepo,
		paymentRepo:   paymentRepo,
		unitRepo:      unitRepo,
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

// ContractMetrics representa as métricas de contratos
type ContractMetrics struct {
	TotalActiveContracts  int64 `json:"total_active_contracts"`
	ContractsExpiringSoon int64 `json:"contracts_expiring_soon"`
	ExpiredContracts      int64 `json:"expired_contracts"`
	CancelledContracts    int64 `json:"cancelled_contracts"`
}

// Alert representa um alerta do sistema
type Alert struct {
	Type        string    `json:"type"`     // "overdue_payment", "expiring_lease", "vacant_unit"
	Severity    string    `json:"severity"` // "high", "medium", "low"
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EntityID    uuid.UUID `json:"entity_id"` // ID do payment, lease ou unit
	CreatedAt   time.Time `json:"created_at"`
}

// DashboardAlerts agrupa todos os alertas
type DashboardAlerts struct {
	OverduePayments []Alert `json:"overdue_payments"`
	ExpiringLeases  []Alert `json:"expiring_leases"`
	VacantUnits     []Alert `json:"vacant_units"`
	TotalAlerts     int     `json:"total_alerts"`
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

// GetContractMetrics retorna métricas sobre os contratos
func (s *DashboardService) GetContractMetrics(ctx context.Context) (*ContractMetrics, error) {
	// 1. Contar contratos ativos
	activeCount, err := s.leaseRepo.CountByStatus(ctx, domain.LeaseStatusActive)
	if err != nil {
		return nil, err
	}

	// 2. Contar contratos expirando em breve
	expiringSoonCount, err := s.leaseRepo.CountByStatus(ctx, domain.LeaseStatusExpiringSoon)
	if err != nil {
		return nil, err
	}

	// 3. Contar contratos expirados
	expiredCount, err := s.leaseRepo.CountByStatus(ctx, domain.LeaseStatusExpired)
	if err != nil {
		return nil, err
	}

	// 4. Contar contratos cancelados
	cancelledCount, err := s.leaseRepo.CountByStatus(ctx, domain.LeaseStatusCancelled)
	if err != nil {
		return nil, err
	}

	return &ContractMetrics{
		TotalActiveContracts:  activeCount,
		ContractsExpiringSoon: expiringSoonCount,
		ExpiredContracts:      expiredCount,
		CancelledContracts:    cancelledCount,
	}, nil
}

// GetAlerts retorna todos os alertas do dashboard
func (s *DashboardService) GetAlerts(ctx context.Context) (*DashboardAlerts, error) {
	alerts := &DashboardAlerts{
		OverduePayments: []Alert{},
		ExpiringLeases:  []Alert{},
		VacantUnits:     []Alert{},
	}

	// 1. Buscar pagamentos em atraso
	overduePayments, err := s.paymentRepo.GetOverdue(ctx)
	if err != nil {
		return nil, err
	}

	for _, payment := range overduePayments {
		daysOverdue := int(time.Since(payment.DueDate).Hours() / 24)

		severity := "medium"
		if daysOverdue > 30 {
			severity = "high"
		} else if daysOverdue > 15 {
			severity = "medium"
		} else {
			severity = "low"
		}

		alerts.OverduePayments = append(alerts.OverduePayments, Alert{
			Type:        "overdue_payment",
			Severity:    severity,
			Title:       fmt.Sprintf("Pagamento em atraso há %d dias", daysOverdue),
			Description: fmt.Sprintf("Valor: R$ %s - Vencimento: %s", payment.Amount.StringFixed(2), payment.DueDate.Format("02/01/2006")),
			EntityID:    payment.ID,
			CreatedAt:   time.Now(),
		})
	}

	// 2. Buscar contratos expirando em breve
	expiringLeases, err := s.leaseRepo.GetExpiringSoon(ctx)
	if err != nil {
		return nil, err
	}

	for _, lease := range expiringLeases {
		daysUntilExpiry := int(time.Until(lease.EndDate).Hours() / 24)

		severity := "medium"
		if daysUntilExpiry <= 15 {
			severity = "high"
		} else if daysUntilExpiry <= 30 {
			severity = "medium"
		} else {
			severity = "low"
		}

		alerts.ExpiringLeases = append(alerts.ExpiringLeases, Alert{
			Type:        "expiring_lease",
			Severity:    severity,
			Title:       fmt.Sprintf("Contrato expira em %d dias", daysUntilExpiry),
			Description: fmt.Sprintf("Data de término: %s", lease.EndDate.Format("02/01/2006")),
			EntityID:    lease.ID,
			CreatedAt:   time.Now(),
		})
	}

	// 3. Buscar unidades disponíveis (sem contrato)
	availableUnits, err := s.unitRepo.ListByStatus(ctx, domain.UnitStatusAvailable)
	if err != nil {
		return nil, err
	}

	// Considerar apenas unidades disponíveis há mais de 30 dias como alerta
	for _, unit := range availableUnits {
		// Calcular quantos dias a unidade está vaga
		daysVacant := int(time.Since(unit.UpdatedAt).Hours() / 24)

		// Só alerta se estiver vaga há mais de 30 dias
		if daysVacant > 30 {
			severity := "low"
			if daysVacant > 90 {
				severity = "high"
			} else if daysVacant > 60 {
				severity = "medium"
			}

			alerts.VacantUnits = append(alerts.VacantUnits, Alert{
				Type:        "vacant_unit",
				Severity:    severity,
				Title:       fmt.Sprintf("Unidade %s vaga há %d dias", unit.Number, daysVacant),
				Description: fmt.Sprintf("Andar %d - Valor: R$ %s", unit.Floor, unit.CurrentRentValue.StringFixed(2)),
				EntityID:    unit.ID,
				CreatedAt:   time.Now(),
			})
		}
	}

	// Calcular total de alertas
	alerts.TotalAlerts = len(alerts.OverduePayments) + len(alerts.ExpiringLeases) + len(alerts.VacantUnits)

	return alerts, nil
}
