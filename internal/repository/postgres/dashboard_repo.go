package postgres

import (
	"context"
	"database/sql"

	"github.com/lucianoZgabriel/kitnet-manager/internal/repository"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository/sqlc"
	"github.com/shopspring/decimal"
)

// Compile-time check to ensure DashboardRepo implements repository.DashboardRepository
var _ repository.DashboardRepository = (*DashboardRepo)(nil)

// DashboardRepo implementa o repository de Dashboard usando SQLC
type DashboardRepo struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewDashboardRepo cria uma nova instância do repository de Dashboard
func NewDashboardRepo(db *sql.DB) repository.DashboardRepository {
	return &DashboardRepo{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *DashboardRepo) GetOccupancyMetrics(ctx context.Context) (*repository.OccupancyMetrics, error) {
	result, err := r.queries.GetOccupancyMetrics(ctx)
	if err != nil {
		return nil, err
	}

	return &repository.OccupancyMetrics{
		TotalUnits:       result.TotalUnits,
		OccupiedUnits:    result.OccupiedUnits,
		AvailableUnits:   result.AvailableUnits,
		MaintenanceUnits: result.MaintenanceUnits,
		RenovationUnits:  result.RenovationUnits,
	}, nil
}

func (r *DashboardRepo) GetMonthlyProjectedRevenue(ctx context.Context) (decimal.Decimal, error) {
	total, err := r.queries.GetMonthlyProjectedRevenue(ctx)
	if err != nil {
		return decimal.Zero, err
	}

	return decimal.NewFromString(total)
}

func (r *DashboardRepo) GetMonthlyRealizedRevenue(ctx context.Context) (decimal.Decimal, error) {
	total, err := r.queries.GetMonthlyRealizedRevenue(ctx)
	if err != nil {
		return decimal.Zero, err
	}

	return decimal.NewFromString(total)
}

func (r *DashboardRepo) GetOverdueAmount(ctx context.Context) (decimal.Decimal, error) {
	total, err := r.queries.GetOverdueAmount(ctx)
	if err != nil {
		return decimal.Zero, err
	}

	return decimal.NewFromString(total)
}

func (r *DashboardRepo) GetTotalPendingAmount(ctx context.Context) (decimal.Decimal, error) {
	total, err := r.queries.GetTotalPendingAmount(ctx)
	if err != nil {
		return decimal.Zero, err
	}

	return decimal.NewFromString(total)
}
