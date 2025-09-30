# Makefile para kitnet-manager

# Carregar variáveis do .env
include .env
export

# Variáveis
MIGRATIONS_PATH = ./migrations
DB_URL = $(DATABASE_URL)

# Cores para output
GREEN = \033[0;32m
YELLOW = \033[0;33m
RED = \033[0;31m
NC = \033[0m # No Color

.PHONY: help
help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: migrate-create
migrate-create: ## Criar nova migration. Uso: make migrate-create name=create_users_table
	@if [ -z "$(name)" ]; then \
		echo "$(RED)Error: Nome da migration é obrigatório$(NC)"; \
		echo "$(YELLOW)Uso: make migrate-create name=nome_da_migration$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Criando migration: $(name)$(NC)"
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

.PHONY: migrate-up
migrate-up: ## Aplicar todas as migrations pendentes
	@echo "$(GREEN)Aplicando migrations...$(NC)"
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up
	@echo "$(GREEN)✓ Migrations aplicadas com sucesso$(NC)"

.PHONY: migrate-down
migrate-down: ## Reverter última migration
	@echo "$(YELLOW)Revertendo última migration...$(NC)"
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1
	@echo "$(GREEN)✓ Migration revertida$(NC)"

.PHONY: migrate-drop
migrate-drop: ## Reverter TODAS as migrations (CUIDADO!)
	@echo "$(RED)ATENÇÃO: Isso irá reverter TODAS as migrations!$(NC)"
	@echo "Pressione Ctrl+C para cancelar ou Enter para continuar..."
	@read confirm
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" drop -f
	@echo "$(GREEN)✓ Todas as migrations foram revertidas$(NC)"

.PHONY: migrate-status
migrate-status: ## Verificar status das migrations
	@echo "$(GREEN)Status das migrations:$(NC)"
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

.PHONY: migrate-force
migrate-force: ## Forçar versão da migration (resolver dirty state). Uso: make migrate-force version=1
	@if [ -z "$(version)" ]; then \
		echo "$(RED)Error: Versão é obrigatória$(NC)"; \
		echo "$(YELLOW)Uso: make migrate-force version=1$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Forçando migration para versão $(version)...$(NC)"
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(version)

.PHONY: sqlc-generate
sqlc-generate: ## Gerar código Go a partir das queries SQL
	@echo "$(GREEN)Gerando código SQLC...$(NC)"
	@sqlc generate
	@echo "$(GREEN)✓ Código SQLC gerado com sucesso$(NC)"

.PHONY: sqlc-verify
sqlc-verify: ## Verificar se as queries SQL estão corretas
	@echo "$(GREEN)Verificando queries SQL...$(NC)"
	@sqlc verify
	@echo "$(GREEN)✓ Queries SQL válidas$(NC)"

.PHONY: db-setup
db-setup: migrate-up sqlc-generate ## Setup completo do banco (migrations + sqlc)
	@echo "$(GREEN)✓ Banco de dados configurado$(NC)"

.PHONY: db-reset
db-reset: migrate-drop migrate-up sqlc-generate ## Reset completo do banco
	@echo "$(GREEN)✓ Banco de dados resetado$(NC)"

.PHONY: run
run: ## Executar a aplicação
	@echo "$(GREEN)Iniciando aplicação...$(NC)"
	@go run cmd/api/main.go

.PHONY: build
build: ## Compilar a aplicação
	@echo "$(GREEN)Compilando...$(NC)"
	@go build -o bin/api cmd/api/main.go
	@echo "$(GREEN)✓ Binário criado: bin/api$(NC)"

.PHONY: test
test: ## Executar testes
	@echo "$(GREEN)Executando testes...$(NC)"
	@go test -v ./...

# Configuração do linter
GOLANGCI_VERSION = 2.5.0
GOLANGCI_BIN = ./bin/golangci-lint

.PHONY: install-lint
install-lint: ## Instalar golangci-lint localmente no projeto
	@echo "$(GREEN)Instalando golangci-lint v$(GOLANGCI_VERSION)...$(NC)"
	@mkdir -p bin
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v$(GOLANGCI_VERSION)
	@echo "$(GREEN)✓ golangci-lint instalado em ./bin$(NC)"

.PHONY: lint
lint: ## Executar linter
	@echo "$(GREEN)Executando linter...$(NC)"
	@if [ -f $(GOLANGCI_BIN) ]; then \
		$(GOLANGCI_BIN) run; \
	else \
		echo "$(YELLOW)Por favor, execute 'make install-lint' primeiro$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✓ Lint concluído$(NC)"

.PHONY: lint-fix
lint-fix: ## Executar linter com correções automáticas
	@echo "$(GREEN)Executando linter com fix...$(NC)"
	@if [ -f $(GOLANGCI_BIN) ]; then \
		$(GOLANGCI_BIN) run --fix; \
	else \
		echo "$(YELLOW)Por favor, execute 'make install-lint' primeiro$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✓ Correções aplicadas$(NC)"