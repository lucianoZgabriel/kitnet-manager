package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(r chi.Router,
	unitService *service.UnitService,
	tenantService *service.TenantService,
	leaseService *service.LeaseService,
	paymentService *service.PaymentService,
	dashboardService *service.DashboardService,
	reportService *service.ReportService) {
	// Criar handlers
	unitHandler := NewUnitHandler(unitService)
	tenantHandler := NewTenantHandler(tenantService)
	leaseHandler := NewLeaseHandler(leaseService)
	paymentHandler := NewPaymentHandler(paymentService)
	dashboardHandler := NewDashboardHandler(dashboardService)
	reportHandler := NewReportHandler(reportService)

	// Rotas de unidades sob /api/v1
	r.Route("/api/v1", func(r chi.Router) {
		// Rotas da unidades
		r.Route("/units", func(r chi.Router) {
			r.Post("/", unitHandler.CreateUnit)
			r.Get("/", unitHandler.ListUnits)
			r.Get("/stats/occupancy", unitHandler.GetOccupancyStats)
			r.Get("/{id}", unitHandler.GetUnit)
			r.Put("/{id}", unitHandler.UpdateUnit)
			r.Patch("/{id}/status", unitHandler.UpdateUnitStatus)
			r.Delete("/{id}", unitHandler.DeleteUnit)
		})

		// Rotas de moradores
		r.Route("/tenants", func(r chi.Router) {
			r.Post("/", tenantHandler.CreateTenant)
			r.Get("/", tenantHandler.ListTenants)
			r.Get("/cpf", tenantHandler.GetTenantByCPF)
			r.Get("/{id}", tenantHandler.GetTenant)
			r.Put("/{id}", tenantHandler.UpdateTenant)
			r.Delete("/{id}", tenantHandler.DeleteTenant)
		})

		// Rotas de contratos
		r.Route("/leases", func(r chi.Router) {
			r.Post("/", leaseHandler.CreateLease)
			r.Get("/", leaseHandler.ListLeases)
			r.Get("/stats", leaseHandler.GetLeaseStats)
			r.Get("/expiring-soon", leaseHandler.GetExpiringSoonLeases)
			r.Get("/{id}", leaseHandler.GetLease)
			r.Post("/{id}/renew", leaseHandler.RenewLease)
			r.Post("/{id}/cancel", leaseHandler.CancelLease)
			r.Patch("/{id}/painting-fee", leaseHandler.UpdatePaintingFeePaid)
			// Rotas de pagamentos por contrato
			r.Get("/{lease_id}/payments", paymentHandler.GetPaymentsByLease)
			r.Get("/{lease_id}/payments/stats", paymentHandler.GetPaymentStatsByLease)
		})

		// Rotas de pagamentos
		r.Route("/payments", func(r chi.Router) {
			r.Get("/overdue", paymentHandler.GetOverduePayments)
			r.Get("/upcoming", paymentHandler.GetUpcomingPayments)
			r.Get("/{id}", paymentHandler.GetPayment)
			r.Put("/{id}/pay", paymentHandler.MarkPaymentAsPaid)
			r.Post("/{id}/cancel", paymentHandler.CancelPayment)
		})

		// Rotas de dashboard
		r.Get("/dashboard", dashboardHandler.GetDashboard)

		// Rotas de relatórios
		r.Route("/reports", func(r chi.Router) {
			r.Get("/financial", reportHandler.GetFinancialReport)
			r.Get("/payments", reportHandler.GetPaymentHistoryReport)
		})
	})
}
