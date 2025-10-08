package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository"
	"github.com/shopspring/decimal"
)

// ReportService contém a lógica de negócio para geração de relatórios
type ReportService struct {
	paymentRepo repository.PaymentRepository
	leaseRepo   repository.LeaseRepository
	unitRepo    repository.UnitRepository
	tenantRepo  repository.TenantRepository
}

// NewReportService cria uma nova instância do serviço de relatórios
func NewReportService(
	paymentRepo repository.PaymentRepository,
	leaseRepo repository.LeaseRepository,
	unitRepo repository.UnitRepository,
	tenantRepo repository.TenantRepository,
) *ReportService {
	return &ReportService{
		paymentRepo: paymentRepo,
		leaseRepo:   leaseRepo,
		unitRepo:    unitRepo,
		tenantRepo:  tenantRepo,
	}
}

// FinancialReportRequest representa os filtros para o relatório financeiro
type FinancialReportRequest struct {
	StartDate   time.Time             `json:"start_date"`
	EndDate     time.Time             `json:"end_date"`
	PaymentType *domain.PaymentType   `json:"payment_type,omitempty"` // Opcional: filtrar por tipo
	Status      *domain.PaymentStatus `json:"status,omitempty"`       // Opcional: filtrar por status
}

// FinancialReportResponse representa o relatório financeiro consolidado
type FinancialReportResponse struct {
	Period        Period                 `json:"period"`
	Summary       FinancialSummary       `json:"summary"`
	ByType        map[string]TypeRevenue `json:"by_type"`
	ByMonth       []MonthlyRevenue       `json:"by_month"`
	ByUnit        []UnitRevenue          `json:"by_unit"`
	TotalPayments int                    `json:"total_payments"`
	GeneratedAt   time.Time              `json:"generated_at"`
}

// Period representa o período do relatório
type Period struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Days      int       `json:"days"`
}

// FinancialSummary representa o resumo financeiro
type FinancialSummary struct {
	TotalRevenue    decimal.Decimal `json:"total_revenue"`
	PaidAmount      decimal.Decimal `json:"paid_amount"`
	PendingAmount   decimal.Decimal `json:"pending_amount"`
	OverdueAmount   decimal.Decimal `json:"overdue_amount"`
	CancelledAmount decimal.Decimal `json:"cancelled_amount"`
}

// TypeRevenue representa receita por tipo de pagamento
type TypeRevenue struct {
	Type   string          `json:"type"`
	Amount decimal.Decimal `json:"amount"`
	Count  int             `json:"count"`
}

// MonthlyRevenue representa receita por mês
type MonthlyRevenue struct {
	Month  string          `json:"month"` // Formato: "2024-03"
	Year   int             `json:"year"`
	Amount decimal.Decimal `json:"amount"`
	Count  int             `json:"count"`
}

// UnitRevenue representa receita por unidade
type UnitRevenue struct {
	UnitID     uuid.UUID       `json:"unit_id"`
	UnitNumber string          `json:"unit_number"`
	Amount     decimal.Decimal `json:"amount"`
	Count      int             `json:"count"`
}

// Erro customizado
var ErrInvalidDateRange = errors.New("end date must be after start date")

// GetFinancialReport gera um relatório financeiro consolidado
func (s *ReportService) GetFinancialReport(ctx context.Context, req FinancialReportRequest) (*FinancialReportResponse, error) {
	// 1. Validar datas
	if req.EndDate.Before(req.StartDate) {
		return nil, ErrInvalidDateRange
	}

	// 2. Buscar todos os pagamentos no período
	allPayments, err := s.paymentRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Filtrar pagamentos pelo período e filtros opcionais
	filteredPayments := s.filterPayments(allPayments, req)

	// 4. Calcular resumo financeiro
	summary := s.calculateSummary(filteredPayments)

	// 5. Agrupar por tipo
	byType := s.groupByType(filteredPayments)

	// 6. Agrupar por mês
	byMonth := s.groupByMonth(filteredPayments, req.StartDate, req.EndDate)

	// 7. Agrupar por unidade
	byUnit, err := s.groupByUnit(ctx, filteredPayments)
	if err != nil {
		return nil, err
	}

	// 8. Calcular dias no período
	days := int(req.EndDate.Sub(req.StartDate).Hours() / 24)

	return &FinancialReportResponse{
		Period: Period{
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
			Days:      days,
		},
		Summary:       summary,
		ByType:        byType,
		ByMonth:       byMonth,
		ByUnit:        byUnit,
		TotalPayments: len(filteredPayments),
		GeneratedAt:   time.Now(),
	}, nil
}

// filterPayments filtra pagamentos baseado nos critérios da request
func (s *ReportService) filterPayments(payments []*domain.Payment, req FinancialReportRequest) []*domain.Payment {
	filtered := make([]*domain.Payment, 0)

	for _, p := range payments {
		// Filtrar por data (usar payment_date se pago, senão due_date)
		dateToCompare := p.DueDate
		if p.Status == domain.PaymentStatusPaid && p.PaymentDate != nil {
			dateToCompare = *p.PaymentDate
		}

		// Verificar se está no período
		if dateToCompare.Before(req.StartDate) || dateToCompare.After(req.EndDate) {
			continue
		}

		// Filtrar por tipo (se especificado)
		if req.PaymentType != nil && p.PaymentType != *req.PaymentType {
			continue
		}

		// Filtrar por status (se especificado)
		if req.Status != nil && p.Status != *req.Status {
			continue
		}

		filtered = append(filtered, p)
	}

	return filtered
}

// calculateSummary calcula o resumo financeiro
func (s *ReportService) calculateSummary(payments []*domain.Payment) FinancialSummary {
	summary := FinancialSummary{
		TotalRevenue:    decimal.Zero,
		PaidAmount:      decimal.Zero,
		PendingAmount:   decimal.Zero,
		OverdueAmount:   decimal.Zero,
		CancelledAmount: decimal.Zero,
	}

	for _, p := range payments {
		summary.TotalRevenue = summary.TotalRevenue.Add(p.Amount)

		switch p.Status {
		case domain.PaymentStatusPaid:
			summary.PaidAmount = summary.PaidAmount.Add(p.Amount)
		case domain.PaymentStatusPending:
			summary.PendingAmount = summary.PendingAmount.Add(p.Amount)
		case domain.PaymentStatusOverdue:
			summary.OverdueAmount = summary.OverdueAmount.Add(p.Amount)
		case domain.PaymentStatusCancelled:
			summary.CancelledAmount = summary.CancelledAmount.Add(p.Amount)
		}
	}

	return summary
}

// groupByType agrupa pagamentos por tipo
func (s *ReportService) groupByType(payments []*domain.Payment) map[string]TypeRevenue {
	typeMap := make(map[string]decimal.Decimal)
	typeCount := make(map[string]int)

	for _, p := range payments {
		typeStr := string(p.PaymentType)
		typeMap[typeStr] = typeMap[typeStr].Add(p.Amount)
		typeCount[typeStr]++
	}

	result := make(map[string]TypeRevenue)
	for typeStr, amount := range typeMap {
		result[typeStr] = TypeRevenue{
			Type:   typeStr,
			Amount: amount,
			Count:  typeCount[typeStr],
		}
	}

	return result
}

// groupByMonth agrupa pagamentos por mês
func (s *ReportService) groupByMonth(payments []*domain.Payment, startDate, endDate time.Time) []MonthlyRevenue {
	monthMap := make(map[string]decimal.Decimal)
	monthCount := make(map[string]int)

	for _, p := range payments {
		// Usar payment_date se pago, senão due_date
		dateToUse := p.DueDate
		if p.Status == domain.PaymentStatusPaid && p.PaymentDate != nil {
			dateToUse = *p.PaymentDate
		}

		monthKey := dateToUse.Format("2006-01")
		monthMap[monthKey] = monthMap[monthKey].Add(p.Amount)
		monthCount[monthKey]++
	}

	// Criar slice ordenado por mês
	result := make([]MonthlyRevenue, 0)
	current := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	for !current.After(end) {
		monthKey := current.Format("2006-01")
		amount := monthMap[monthKey]
		count := monthCount[monthKey]

		result = append(result, MonthlyRevenue{
			Month:  monthKey,
			Year:   current.Year(),
			Amount: amount,
			Count:  count,
		})

		current = current.AddDate(0, 1, 0)
	}

	return result
}

// groupByUnit agrupa pagamentos por unidade
func (s *ReportService) groupByUnit(ctx context.Context, payments []*domain.Payment) ([]UnitRevenue, error) {
	// Mapear lease_id -> unit_id e unit_number
	leaseUnitMap := make(map[uuid.UUID]struct {
		UnitID     uuid.UUID
		UnitNumber string
	})

	unitAmountMap := make(map[uuid.UUID]decimal.Decimal)
	unitCountMap := make(map[uuid.UUID]int)

	for _, p := range payments {
		// Se já temos o mapeamento do lease, usa
		if _, exists := leaseUnitMap[p.LeaseID]; !exists {
			// Buscar o lease para pegar o unit_id
			lease, err := s.leaseRepo.GetByID(ctx, p.LeaseID)
			if err != nil {
				continue // Skip se não encontrar o lease
			}

			// Buscar a unit para pegar o número
			unit, err := s.unitRepo.GetByID(ctx, lease.UnitID)
			if err != nil {
				continue // Skip se não encontrar a unit
			}

			leaseUnitMap[p.LeaseID] = struct {
				UnitID     uuid.UUID
				UnitNumber string
			}{
				UnitID:     unit.ID,
				UnitNumber: unit.Number,
			}
		}

		unitInfo := leaseUnitMap[p.LeaseID]
		unitAmountMap[unitInfo.UnitID] = unitAmountMap[unitInfo.UnitID].Add(p.Amount)
		unitCountMap[unitInfo.UnitID]++
	}

	// Converter map para slice
	result := make([]UnitRevenue, 0, len(unitAmountMap))
	for unitID, amount := range unitAmountMap {
		// Buscar o número da unidade do cache
		var unitNumber string
		for _, info := range leaseUnitMap {
			if info.UnitID == unitID {
				unitNumber = info.UnitNumber
				break
			}
		}

		result = append(result, UnitRevenue{
			UnitID:     unitID,
			UnitNumber: unitNumber,
			Amount:     amount,
			Count:      unitCountMap[unitID],
		})
	}

	return result, nil
}

// PaymentHistoryRequest representa os filtros para histórico de pagamentos
type PaymentHistoryRequest struct {
	LeaseID   *uuid.UUID            `json:"lease_id,omitempty"`
	TenantID  *uuid.UUID            `json:"tenant_id,omitempty"`
	Status    *domain.PaymentStatus `json:"status,omitempty"`
	StartDate *time.Time            `json:"start_date,omitempty"`
	EndDate   *time.Time            `json:"end_date,omitempty"`
}

// PaymentHistoryResponse representa o histórico de pagamentos
type PaymentHistoryResponse struct {
	Payments    []PaymentHistoryItem `json:"payments"`
	TotalCount  int                  `json:"total_count"`
	TotalAmount decimal.Decimal      `json:"total_amount"`
	GeneratedAt time.Time            `json:"generated_at"`
}

// PaymentHistoryItem representa um item do histórico
type PaymentHistoryItem struct {
	PaymentID     uuid.UUID             `json:"payment_id"`
	LeaseID       uuid.UUID             `json:"lease_id"`
	UnitNumber    string                `json:"unit_number"`
	TenantName    string                `json:"tenant_name"`
	PaymentType   domain.PaymentType    `json:"payment_type"`
	Amount        decimal.Decimal       `json:"amount"`
	Status        domain.PaymentStatus  `json:"status"`
	DueDate       time.Time             `json:"due_date"`
	PaymentDate   *time.Time            `json:"payment_date,omitempty"`
	PaymentMethod *domain.PaymentMethod `json:"payment_method,omitempty"`
	CreatedAt     time.Time             `json:"created_at"`
}

// GetPaymentHistoryReport gera um relatório de histórico de pagamentos
func (s *ReportService) GetPaymentHistoryReport(ctx context.Context, req PaymentHistoryRequest) (*PaymentHistoryResponse, error) {
	// 1. Buscar pagamentos baseado nos filtros
	var payments []*domain.Payment
	var err error

	if req.LeaseID != nil {
		payments, err = s.paymentRepo.ListByLeaseID(ctx, *req.LeaseID)
	} else if req.Status != nil {
		payments, err = s.paymentRepo.ListByStatus(ctx, *req.Status)
	} else {
		payments, err = s.paymentRepo.List(ctx)
	}

	if err != nil {
		return nil, err
	}

	// 2. Aplicar filtros adicionais
	filtered := s.filterPaymentHistory(payments, req)

	// 3. Enriquecer com dados de lease, unit e tenant
	items := make([]PaymentHistoryItem, 0, len(filtered))
	totalAmount := decimal.Zero

	for _, p := range filtered {
		// Buscar lease
		lease, err := s.leaseRepo.GetByID(ctx, p.LeaseID)
		if err != nil {
			continue // Skip se não encontrar
		}

		// Filtrar por tenant_id se especificado
		if req.TenantID != nil && lease.TenantID != *req.TenantID {
			continue
		}

		// Buscar unit
		unit, err := s.unitRepo.GetByID(ctx, lease.UnitID)
		if err != nil {
			continue
		}

		// Buscar tenant
		tenant, err := s.tenantRepo.GetByID(ctx, lease.TenantID)
		if err != nil {
			continue
		}

		items = append(items, PaymentHistoryItem{
			PaymentID:     p.ID,
			LeaseID:       p.LeaseID,
			UnitNumber:    unit.Number,
			TenantName:    tenant.FullName,
			PaymentType:   p.PaymentType,
			Amount:        p.Amount,
			Status:        p.Status,
			DueDate:       p.DueDate,
			PaymentDate:   p.PaymentDate,
			PaymentMethod: p.PaymentMethod,
			CreatedAt:     p.CreatedAt,
		})

		totalAmount = totalAmount.Add(p.Amount)
	}

	return &PaymentHistoryResponse{
		Payments:    items,
		TotalCount:  len(items),
		TotalAmount: totalAmount,
		GeneratedAt: time.Now(),
	}, nil
}

// filterPaymentHistory filtra pagamentos para o histórico
func (s *ReportService) filterPaymentHistory(payments []*domain.Payment, req PaymentHistoryRequest) []*domain.Payment {
	filtered := make([]*domain.Payment, 0)

	for _, p := range payments {
		// Filtrar por data range
		if req.StartDate != nil && p.DueDate.Before(*req.StartDate) {
			continue
		}
		if req.EndDate != nil && p.DueDate.After(*req.EndDate) {
			continue
		}

		filtered = append(filtered, p)
	}

	return filtered
}
