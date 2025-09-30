package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(r chi.Router, unitService *service.UnitService) {
	// Criar handlers
	unitHandler := NewUnitHandler(unitService)

	// Rotas de unidades sob /api/v1
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/units", func(r chi.Router) {
			r.Post("/", unitHandler.CreateUnit)
			r.Get("/", unitHandler.ListUnits)
			r.Get("/stats/occupancy", unitHandler.GetOccupancyStats) // DEVE vir ANTES do /{id}
			r.Get("/{id}", unitHandler.GetUnit)
			r.Put("/{id}", unitHandler.UpdateUnit)
			r.Patch("/{id}/status", unitHandler.UpdateUnitStatus)
			r.Delete("/{id}", unitHandler.DeleteUnit)
		})
	})
}
