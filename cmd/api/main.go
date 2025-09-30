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
	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/database"
	"github.com/lucianoZgabriel/kitnet-manager/internal/pkg/response"
)

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
	defer dbConn.Close()

	log.Println("✅ Conectado ao banco de dados")

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
