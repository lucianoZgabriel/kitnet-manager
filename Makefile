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