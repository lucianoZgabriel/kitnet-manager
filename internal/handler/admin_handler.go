package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/response"
	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/scheduler"
)

// AdminHandler lida com requisições administrativas
type AdminHandler struct {
	scheduler *scheduler.Scheduler
}

// NewAdminHandler cria uma nova instância do handler administrativo
func NewAdminHandler(scheduler *scheduler.Scheduler) *AdminHandler {
	return &AdminHandler{
		scheduler: scheduler,
	}
}

// ForceSchedulerRun godoc
// @Summary      Forçar execução do scheduler
// @Description  Executa todas as tarefas agendadas imediatamente (útil para testes)
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /admin/force-scheduler [post]
func (h *AdminHandler) ForceSchedulerRun(w http.ResponseWriter, r *http.Request) {
	// Criar contexto com timeout
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	// Executar tarefas do scheduler
	h.scheduler.RunScheduledTasks(ctx)

	response.Success(w, http.StatusOK, "Scheduler executado com sucesso", map[string]interface{}{
		"executed_at": time.Now().Format(time.RFC3339),
		"tasks": []string{
			"mark_overdue_payments",
			"check_expiring_soon_leases",
			"auto_renew_leases",
		},
	})
}
