package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lucianoZgabriel/kitnet-manager/internal/config"
	"github.com/lucianoZgabriel/kitnet-manager/internal/handler"
	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/database"
	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/response"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository/postgres"
	"github.com/lucianoZgabriel/kitnet-manager/internal/service"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/lucianoZgabriel/kitnet-manager/docs" // Swagger docs
)

// @title           Kitnet Manager API
// @version         1.0
// @description     API para gestão de complexo de kitnets com 31 unidades
// @description     Sistema completo de gerenciamento de unidades, moradores, contratos e pagamentos

// @contact.name   Luciano Gabriel
// @contact.email  contato@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes http https

// @tag.name Units
// @tag.description Operações relacionadas a unidades/kitnets

// @tag.name Health
// @tag.description Health check e status do sistema

func main() {
	// Carregar configurações
	cfg := config.Load()

	log.Printf("🚀 Iniciando Kitnet Manager [%s]", cfg.Environment)

	// Conectar ao banco de dados
	dbConn, err := database.NewConnection(database.Config{
		URL:            cfg.Database.URL,
		MaxConnections: cfg.Database.MaxConnections,
		MaxIdleConns:   cfg.Database.MaxIdleConns,
		MaxLifetime:    cfg.Database.MaxLifetime,
	})
	if err != nil {
		log.Fatal("Erro ao conectar com banco de dados:", err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("Erro ao fechar conexão com banco: %v", err)
		}
	}()

	log.Println("✅ Conectado ao banco de dados")

	// Inicializar camadas da aplicação
	// Repository
	unitRepo := postgres.NewUnitRepository(dbConn.DB)

	// Service
	unitService := service.NewUnitService(unitRepo)

	log.Println("✅ Serviços inicializados")

	// Configurar roteador
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := dbConn.Health(); err != nil {
			response.Error(w, http.StatusServiceUnavailable, "Database unhealthy")
			return
		}
		response.Success(w, http.StatusOK, "Server is healthy", nil)
	})

	// Rota de teste
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, "Kitnet Manager API", map[string]string{
			"version":     "1.0.0",
			"environment": cfg.Environment,
		})
	})

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+cfg.Port+"/swagger/doc.json"),
	))

	// Registrar rotas da aplicação
	handler.SetupRoutes(r, unitService)

	log.Println("✅ Rotas configuradas")
	log.Printf("📚 Documentação Swagger: http://localhost:%s/swagger/index.html", cfg.Port)

	// Configurar servidor
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("📡 Servidor rodando na porta %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Erro ao iniciar servidor:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Desligando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Erro ao desligar servidor:", err)
	}

	log.Println("✅ Servidor desligado com sucesso")
}
