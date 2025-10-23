package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/lucianoZgabriel/kitnet-manager/internal/service"
)

// Scheduler executa tarefas agendadas periodicamente
type Scheduler struct {
	paymentService *service.PaymentService
	leaseService   *service.LeaseService
	stopChan       chan struct{}
}

// New cria uma nova instância do Scheduler
func New(paymentService *service.PaymentService, leaseService *service.LeaseService) *Scheduler {
	return &Scheduler{
		paymentService: paymentService,
		leaseService:   leaseService,
		stopChan:       make(chan struct{}),
	}
}

// Start inicia o scheduler em background
func (s *Scheduler) Start(ctx context.Context) {
	log.Println("⏰ Scheduler iniciado")

	// Executar verificações imediatamente na inicialização
	s.runScheduledTasks(ctx)

	// Criar ticker para executar a cada 6 horas
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.runScheduledTasks(ctx)
		case <-s.stopChan:
			log.Println("⏹️ Scheduler parado")
			return
		case <-ctx.Done():
			log.Println("⏹️ Scheduler interrompido pelo contexto")
			return
		}
	}
}

// Stop para o scheduler
func (s *Scheduler) Stop() {
	close(s.stopChan)
}

// runScheduledTasks executa todas as tarefas agendadas
func (s *Scheduler) runScheduledTasks(ctx context.Context) {
	log.Println("🔄 Executando tarefas agendadas...")

	// Tarefa 1: Marcar pagamentos atrasados
	s.markOverduePayments(ctx)

	// Tarefa 2: Atualizar contratos expirando em breve
	s.checkExpiringSoonLeases(ctx)

	log.Println("✅ Tarefas agendadas concluídas")
}

// markOverduePayments marca pagamentos pendentes vencidos como atrasados
func (s *Scheduler) markOverduePayments(ctx context.Context) {
	log.Println("📅 Verificando pagamentos atrasados...")

	result, err := s.paymentService.CheckOverduePayments(ctx)
	if err != nil {
		log.Printf("❌ Erro ao marcar pagamentos atrasados: %v", err)
		return
	}

	if result.UpdatedCount > 0 {
		log.Printf("✅ %d pagamento(s) marcado(s) como atrasado(s)", result.UpdatedCount)
	} else {
		log.Println("✓ Nenhum pagamento atrasado encontrado")
	}
}

// checkExpiringSoonLeases verifica contratos próximos de expirar
func (s *Scheduler) checkExpiringSoonLeases(ctx context.Context) {
	log.Println("📅 Verificando contratos próximos de expirar...")

	updatedCount, err := s.leaseService.CheckExpiringSoonLeases(ctx)
	if err != nil {
		log.Printf("❌ Erro ao verificar contratos expirando: %v", err)
		return
	}

	if updatedCount > 0 {
		log.Printf("✅ %d contrato(s) marcado(s) como expirando em breve", updatedCount)
	} else {
		log.Println("✓ Nenhum contrato expirando em breve")
	}
}
