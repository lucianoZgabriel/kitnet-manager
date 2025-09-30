# Kitnet Manager

  Sistema de gestão para administração de kitnets.

  ## Descrição

  Sistema para gerenciamento de um complexo de 31 kitnets, substituindo controles manuais em
  Excel por uma solução digital robusta.

  ## Status do Projeto

  🚧 Em desenvolvimento - Sprint 0: Setup inicial

  ## Tecnologias e Dependências

  ### Backend
  - **Linguagem:** Go 1.21+
  - **Database:** PostgreSQL 17.5 (Neon)

  ### Principais Dependências
  - **Chi Router** (`go-chi/chi/v5`) - Roteador HTTP leve e idiomático
  - **pq** (`lib/pq`) - Driver PostgreSQL nativo para Go
  - **godotenv** (`joho/godotenv`) - Carregamento de variáveis de ambiente
  - **validator** (`go-playground/validator/v10`) - Validação de structs e campos
  - **uuid** (`google/uuid`) - Geração e manipulação de UUIDs
  - **decimal** (`shopspring/decimal`) - Precisão decimal para valores monetários

  ## Estrutura do Projeto

  kitnet-manager/
  ├── cmd/
  │   └── api/              # Ponto de entrada da aplicação
  ├── internal/
  │   ├── domain/           # Entidades de negócio
  │   ├── repository/       # Camada de acesso a dados
  │   │   ├── postgres/     # Implementação PostgreSQL
  │   │   └── queries/      # Queries SQL para SQLC
  │   ├── service/          # Lógica de negócio
  │   ├── handler/          # Handlers HTTP
  │   └── pkg/              # Pacotes internos reutilizáveis
  │       ├── database/     # Configuração de banco
  │       ├── validator/    # Validações customizadas
  │       └── response/     # Respostas HTTP padronizadas
  ├── migrations/           # Migrations do banco de dados
  ├── config/              # Arquivos de configuração
  └── docs/
      └── api/             # Documentação da API

## Workflow de Desenvolvimento

  ### Trabalhando com Banco de Dados

  1. **Criar nova tabela/alteração:**
     ```bash
     make migrate-create name=descrição_da_mudança

  2. Escrever queries SQL:
    - Adicione queries em internal/repository/queries/
    - Use comentários especiais do SQLC: -- name: NomeDaFuncao :tipo
  3. Gerar código:
  make sqlc-generate
  4. Aplicar no banco:
  make migrate-up

  Tipos de queries SQLC:

  - :one - Retorna um único registro
  - :many - Retorna múltiplos registros
  - :exec - Executa sem retorno (DELETE, UPDATE)
  - :execrows - Executa e retorna número de linhas afetadas
  - :copyfrom - Bulk insert eficiente

## Comandos Disponíveis

  ```bash
  make help          # Mostra todos os comandos disponíveis
  make run           # Executa a aplicação
  make build         # Compila o binário
  make test          # Executa os testes
  make clean         # Limpa arquivos gerados

  # Banco de dados
  make migrate-create name=nome_da_migration  # Cria nova migration
  make migrate-up    # Aplica todas as migrations
  make migrate-down  # Reverte última migration
  make migrate-status # Verifica status das migrations
  make sqlc-generate # Gera código a partir das queries SQL

  # Desenvolvimento
  make dev           # Roda em modo desenvolvimento
  make db-setup      # Setup completo do banco

  ### Testar comando clean

  ```bash
  # Testar o comando clean
  make build
  ls bin/
  make clean
  ls bin/  # Deve dar erro - pasta não existe mais

  ## Documentação

  - [Arquitetura](kitnet_architecture.md)
  - [Roadmap](kitnet_roadmap.md)

  ## Como executar

  Em breve...

  ## Licença

  Projeto privado

  1. Verificar estrutura criada