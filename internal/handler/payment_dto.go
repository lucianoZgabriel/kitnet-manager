package handler

import (
	"time"

	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
)

// MarkPaymentAsPaidRequestDTO representa os dados para marcar pagamento como pago
type MarkPaymentAsPaidRequestDTO struct {
	PaymentDate   time.Time            `json:"payment_date" validate:"required"`
	PaymentMethod domain.PaymentMethod `json:"payment_method" validate:"required"`
}

// PaymentStatsResponse representa as estatísticas de pagamentos
type PaymentStatsResponse struct {
	TotalPaid     float64 `json:"total_paid"`
	TotalPending  float64 `json:"total_pending"`
	TotalPayments int64   `json:"total_payments"`
	PaidCount     int64   `json:"paid_count"`
	PendingCount  int64   `json:"pending_count"`
	OverdueCount  int64   `json:"overdue_count"`
}

// ToPaymentStatsResponse converte service.PaymentStats para PaymentStatsResponse
func ToPaymentStatsResponse(stats *service.PaymentStats) *PaymentStatsResponse {
	totalPaid, _ := stats.TotalPaid.Float64()
	totalPending, _ := stats.TotalPending.Float64()

	return &PaymentStatsResponse{
		TotalPaid:     totalPaid,
		TotalPending:  totalPending,
		TotalPayments: stats.TotalPayments,
		PaidCount:     stats.PaidCount,
		PendingCount:  stats.PendingCount,
		OverdueCount:  stats.OverdueCount,
	}
}
