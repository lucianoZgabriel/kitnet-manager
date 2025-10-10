package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/middleware"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(r chi.Router,
	unitService *service.UnitService,
	tenantService *service.TenantService,
	leaseService *service.LeaseService,
	paymentService *service.PaymentService,
	dashboardService *service.DashboardService,
	reportService *service.ReportService,
	authService *service.AuthService,
	authMiddleware *middleware.AuthMiddleware) {
	// Criar handlers
	unitHandler := NewUnitHandler(unitService)
	tenantHandler := NewTenantHandler(tenantService)
	leaseHandler := NewLeaseHandler(leaseService)
	paymentHandler := NewPaymentHandler(paymentService)
	dashboardHandler := NewDashboardHandler(dashboardService)
	reportHandler := NewReportHandler(reportService)
	authHandler := NewAuthHandler(authService)

	// Rotas públicas de autenticação
	r.Route("/api/v1/auth", func(r chi.Router) {
		// Rotas públicas (sem autenticação)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.RefreshToken)

		// Rotas protegidas (requerem autenticação)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			r.Get("/me", authHandler.GetCurrentUser)
			r.Post("/change-password", authHandler.ChangePassword)

			// Rotas admin apenas
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAdmin)

				r.Post("/users", authHandler.CreateUser)
				r.Get("/users", authHandler.ListUsers)
				r.Get("/users/{id}", authHandler.GetUser)
				r.Patch("/users/{id}/role", authHandler.ChangeUserRole)
				r.Post("/users/{id}/deactivate", authHandler.DeactivateUser)
				r.Post("/users/{id}/activate", authHandler.ActivateUser)
			})
		})
	})

	// Rotas protegidas da aplicação (requerem autenticação)
	r.Route("/api/v1", func(r chi.Router) {
		// Aplicar middleware de autenticação em todas as rotas
		r.Use(authMiddleware.Authenticate)

		// Rotas de unidades (Admin e Manager podem escrever, todos podem ler)
		r.Route("/units", func(r chi.Router) {
			// Rotas de leitura (todos autenticados)
			r.Get("/", unitHandler.ListUnits)
			r.Get("/stats/occupancy", unitHandler.GetOccupancyStats)
			r.Get("/{id}", unitHandler.GetUnit)

			// Rotas de escrita (Admin e Manager apenas)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAdminOrManager)
				r.Post("/", unitHandler.CreateUnit)
				r.Put("/{id}", unitHandler.UpdateUnit)
				r.Patch("/{id}/status", unitHandler.UpdateUnitStatus)
				r.Delete("/{id}", unitHandler.DeleteUnit)
			})
		})

		// Rotas de moradores (Admin e Manager podem escrever, todos podem ler)
		r.Route("/tenants", func(r chi.Router) {
			// Rotas de leitura
			r.Get("/", tenantHandler.ListTenants)
			r.Get("/cpf", tenantHandler.GetTenantByCPF)
			r.Get("/{id}", tenantHandler.GetTenant)

			// Rotas de escrita
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAdminOrManager)
				r.Post("/", tenantHandler.CreateTenant)
				r.Put("/{id}", tenantHandler.UpdateTenant)
				r.Delete("/{id}", tenantHandler.DeleteTenant)
			})
		})

		// Rotas de contratos (Admin e Manager podem escrever, todos podem ler)
		r.Route("/leases", func(r chi.Router) {
			// Rotas de leitura
			r.Get("/", leaseHandler.ListLeases)
			r.Get("/stats", leaseHandler.GetLeaseStats)
			r.Get("/expiring-soon", leaseHandler.GetExpiringSoonLeases)
			r.Get("/{id}", leaseHandler.GetLease)
			r.Get("/{lease_id}/payments", paymentHandler.GetPaymentsByLease)
			r.Get("/{lease_id}/payments/stats", paymentHandler.GetPaymentStatsByLease)

			// Rotas de escrita
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAdminOrManager)
				r.Post("/", leaseHandler.CreateLease)
				r.Post("/{id}/renew", leaseHandler.RenewLease)
				r.Post("/{id}/cancel", leaseHandler.CancelLease)
				r.Patch("/{id}/painting-fee", leaseHandler.UpdatePaintingFeePaid)
			})
		})

		// Rotas de pagamentos (Admin e Manager podem escrever, todos podem ler)
		r.Route("/payments", func(r chi.Router) {
			// Rotas de leitura
			r.Get("/overdue", paymentHandler.GetOverduePayments)
			r.Get("/upcoming", paymentHandler.GetUpcomingPayments)
			r.Get("/{id}", paymentHandler.GetPayment)

			// Rotas de escrita
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAdminOrManager)
				r.Put("/{id}/pay", paymentHandler.MarkPaymentAsPaid)
				r.Post("/{id}/cancel", paymentHandler.CancelPayment)
			})
		})

		// Rotas de dashboard (todos podem ler)
		r.Get("/dashboard", dashboardHandler.GetDashboard)

		// Rotas de relatórios (todos podem ler)
		r.Route("/reports", func(r chi.Router) {
			r.Get("/financial", reportHandler.GetFinancialReport)
			r.Get("/payments", reportHandler.GetPaymentHistoryReport)
		})
	})
}
