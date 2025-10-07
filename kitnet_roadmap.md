# Roadmap - Kitnet Manager

> Plano de desenvolvimento detalhado com sprints e tarefas granulares

## 📅 Visão Geral

- **Duração estimada do MVP:** 8-10 semanas
- **Metodologia:** Desenvolvimento iterativo e incremental
- **Entregas:** Funcionalidades testáveis ao final de cada sprint

---

## Sprint 0: Setup & Infraestrutura
**Duração:** 2-3 dias  
**Objetivo:** Preparar ambiente de desenvolvimento e infraestrutura básica

### 0.1 Configuração do Controle de Versão
- [x] Criar conta no GitHub (se não tiver)
- [x] Criar novo repositório: `kitnet-manager`
- [x] Configurar visibilidade (privado ou público)
- [x] Adicionar `.gitignore` para Go
- [x] Criar branch `main` como padrão
- [ ] Configurar proteção da branch main (opcional)
- [x] Fazer commit inicial com README básico

### 0.2 Estrutura do Projeto Local
- [x] Criar diretório local do projeto
- [x] Inicializar Go module: `go mod init github.com/seu-usuario/kitnet-manager`
- [x] Criar estrutura de diretórios completa
- [x] Adicionar README.md detalhado
- [x] Adicionar ARCHITECTURE.md
- [x] Adicionar ROADMAP.md
- [ ] Adicionar LICENSE (se aplicável)

### 0.3 Configuração do Neon Database
- [x] Criar conta no Neon (https://neon.tech)
- [x] Criar novo projeto no Neon
- [x] Anotar connection string
- [x] Testar conexão localmente
- [x] Documentar credenciais no .env.example

### 0.4 Configuração de Dependências
- [x] Instalar Chi router: `go get github.com/go-chi/chi/v5`
- [x] Instalar driver PostgreSQL: `go get github.com/lib/pq`
- [x] Instalar godotenv: `go get github.com/joho/godotenv`
- [x] Instalar validator: `go get github.com/go-playground/validator/v10`
- [x] Instalar UUID: `go get github.com/google/uuid`
- [x] Instalar decimal: `go get github.com/shopspring/decimal`
- [x] Documentar todas as dependências no README

### 0.5 Setup de Migrations
- [x] Instalar golang-migrate CLI
- [x] Criar diretório `migrations/`
- [x] Adicionar comando no Makefile para criar migrations
- [x] Adicionar comando no Makefile para executar migrations up
- [x] Adicionar comando no Makefile para executar migrations down
- [x] Testar criação de migration de exemplo

### 0.6 Setup do SQLC
- [x] Instalar SQLC CLI
- [x] Criar diretório `internal/repository/queries/`
- [x] Criar arquivo `sqlc.yaml` com configurações
- [x] Adicionar comando no Makefile para gerar código SQLC
- [x] Documentar workflow do SQLC no README

### 0.7 Estrutura Base da Aplicação
- [x] Criar `cmd/api/main.go` com estrutura básica
- [x] Criar `internal/pkg/database/postgres.go` para conexão
- [x] Criar `internal/pkg/response/response.go` para padronização
- [x] Criar `.env.example` com variáveis necessárias
- [x] Criar `.env` local (gitignored)
- [x] Testar inicialização básica da aplicação

### 0.8 Makefile e Scripts
- [x] Criar Makefile com comandos úteis:
  - `make run` - executar aplicação
  - `make build` - compilar binário
  - `make test` - executar testes
  - `make migrate-up` - aplicar migrations
  - `make migrate-down` - reverter migrations
  - `make migrate-create` - criar nova migration
  - `make sqlc-generate` - gerar código SQLC
  - `make lint` - executar linter (futuro)
- [x] Documentar comandos no README
- [x] Testar todos os comandos

### 0.9 Commit e Push Inicial
- [x] Revisar todos os arquivos criados
- [x] Adicionar arquivos ao git
- [x] Fazer commit: "chore: initial project setup"
- [x] Push para repositório remoto
- [x] Verificar no GitHub

---

## Sprint 1: CRUD de Unidades e Moradores
**Duração:** 3-4 dias  
**Objetivo:** Implementar gestão completa de unidades e moradores

### 1.1 Migration e Schema - Units
- [x] Criar migration `000001_create_units_table.up.sql`
- [x] Definir tabela units com todos os campos
- [x] Adicionar constraints e checks
- [x] Criar índices necessários
- [x] Criar migration down correspondente
- [x] Executar migration e verificar no Neon
- [x] Adicionar schema no arquivo de referência SQLC

### 1.2 Domain Model - Unit
- [x] Criar arquivo `internal/domain/unit.go`
- [x] Definir struct Unit com todos os campos
- [x] Definir enum UnitStatus (available, occupied, maintenance, renovation)
- [x] Adicionar métodos de validação no domínio
- [x] Adicionar método CalculateCurrentRentValue()
- [x] Adicionar testes unitários do domínio

### 1.3 Repository - Unit (SQLC)
- [x] Criar arquivo `internal/repository/queries/units.sql`
- [x] Escrever query CreateUnit
- [x] Escrever query GetUnitByID
- [x] Escrever query ListUnits (com filtros opcionais)
- [x] Escrever query UpdateUnit
- [x] Escrever query UpdateUnitStatus
- [x] Escrever query DeleteUnit
- [x] Gerar código com SQLC
- [x] Criar `internal/repository/postgres/unit_repo.go`
- [x] Implementar interface UnitRepository
- [ ] Adicionar testes de integração (opcional neste momento)

### 1.4 Service - Unit
- [x] Criar arquivo `internal/service/unit_service.go`
- [x] Definir struct UnitService com dependências
- [x] Implementar CreateUnit com validações de negócio
- [x] Implementar GetUnitByID
- [x] Implementar ListUnits com filtros
- [x] Implementar UpdateUnit
- [x] Implementar UpdateUnitStatus (validar transições)
- [x] Implementar DeleteUnit (validar se não tem contrato ativo)
- [x] Adicionar testes unitários do service

### 1.5 Handler - Unit
- [x] Criar arquivo `internal/handler/unit_handler.go`
- [x] Definir struct UnitHandler
- [x] Criar DTOs (CreateUnitRequest, UpdateUnitRequest, UnitResponse)
- [x] Implementar CreateUnit handler (POST /api/units)
- [x] Implementar GetUnit handler (GET /api/units/:id)
- [x] Implementar ListUnits handler (GET /api/units)
- [x] Implementar UpdateUnit handler (PUT /api/units/:id)
- [x] Implementar UpdateUnitStatus handler (PATCH /api/units/:id/status)
- [x] Implementar DeleteUnit handler (DELETE /api/units/:id)
- [x] Adicionar validação de inputs

### 1.6 Router - Units
- [x] Criar ou atualizar `internal/handler/router.go`
- [x] Registrar rotas de units
- [x] Configurar middlewares básicos (logger, CORS)
- [x] Agrupar rotas sob /api/v1
- [x] Testar todas as rotas manualmente (Postman/cURL)

### 1.7 Migration e Schema - Tenants
- [x] Criar migration `000002_create_tenants_table.up.sql`
- [x] Definir tabela tenants com todos os campos
- [x] Adicionar constraint UNIQUE no CPF
- [x] Criar índice no CPF
- [x] Criar migration down correspondente
- [x] Executar migration e verificar no Neon
- [ ] Adicionar schema no arquivo de referência SQLC

### 1.8 Domain Model - Tenant
- [x] Criar arquivo `internal/domain/tenant.go`
- [x] Definir struct Tenant
- [x] Adicionar método ValidateCPF()
- [x] Adicionar método FormatPhone()
- [x] Adicionar testes unitários

### 1.9 Repository - Tenant (SQLC)
- [x] Criar arquivo `internal/repository/queries/tenants.sql`
- [x] Escrever queries: Create, GetByID, GetByCPF, List, Update, Delete
- [x] Gerar código com SQLC
- [x] Criar `internal/repository/postgres/tenant_repo.go`
- [x] Implementar interface TenantRepository

### 1.10 Service - Tenant
- [x] Criar arquivo `internal/service/tenant_service.go`
- [x] Implementar CreateTenant (validar CPF único)
- [x] Implementar GetTenantByID
- [x] Implementar GetTenantByCPF
- [x] Implementar ListTenants
- [x] Implementar UpdateTenant
- [x] Implementar DeleteTenant (validar se não tem contrato ativo)
- [x] Adicionar testes unitários

### 1.11 Handler - Tenant
- [x] Criar arquivo `internal/handler/tenant_handler.go`
- [x] Criar DTOs necessários
- [x] Implementar handlers para todas as operações CRUD
- [x] Adicionar validação de CPF no handler

### 1.12 Router - Tenants
- [x] Registrar rotas de tenants no router
- [x] Testar todas as rotas manualmente

### 1.13 Testes e Documentação
- [x] Testar fluxo completo de unidades
- [x] Testar fluxo completo de moradores
- [x] Documentar endpoints no README ou criar doc/api/
- [x] Commit: "feat: implement units and tenants CRUD"
- [x] Push para repositório

---

## Sprint 2: Gestão de Contratos (Leases)
**Duração:** 4-5 dias  
**Objetivo:** Implementar sistema completo de contratos com regras de negócio

### 2.1 Migration e Schema - Leases
- [x] Criar migration `000003_create_leases_table.up.sql`
- [x] Definir tabela leases com todas as colunas
- [x] Adicionar foreign keys para units e tenants
- [x] Adicionar checks (payment_due_day entre 1-31, installments 1-4)
- [x] Criar índices necessários (unit_id, tenant_id, status)
- [x] Criar migration down
- [x] Executar e verificar

### 2.2 Domain Model - Lease
- [x] Criar arquivo `internal/domain/lease.go`
- [x] Definir struct Lease completa
- [x] Definir enum LeaseStatus
- [x] Implementar método CalculateEndDate() (start + 6 meses)
- [x] Implementar método IsExpiringSoon() (< 45 dias)
- [x] Implementar método CanBeRenewed()
- [x] Implementar método RemainingPaintingFee()
- [x] Adicionar testes unitários de todos os métodos

### 2.3 Repository - Lease (SQLC)
- [x] Criar arquivo `internal/repository/queries/leases.sql`
- [x] Query CreateLease
- [x] Query GetLeaseByID (com JOIN de unit e tenant)
- [x] Query ListLeases (filtros: status, unit_id, tenant_id)
- [x] Query GetActiveLeaseByUnitID
- [x] Query GetActiveLeaseByTenantID
- [x] Query UpdateLease
- [x] Query UpdateLeaseStatus
- [x] Query GetExpiringSoonLeases (end_date < now + 45 days)
- [x] Gerar código SQLC
- [x] Implementar repository

### 2.4 Service - Lease (Parte 1: Criação)
- [x] Criar arquivo `internal/service/lease_service.go`
- [x] Definir dependências (leaseRepo, unitRepo, tenantRepo)
- [x] Implementar CreateLease:
  - Validar unidade existe e está disponível
  - Validar morador existe
  - Validar não há contrato ativo para essa unidade
  - Validar não há contrato ativo para esse morador
  - Validar datas (start_date < end_date)
  - Calcular end_date automaticamente
  - Criar lease
  - Atualizar unit.status = occupied
- [x] Adicionar testes do CreateLease

### 2.5 Service - Lease (Parte 2: Outras Operações)
- [x] Implementar GetLeaseByID
- [x] Implementar ListLeases com filtros
- [x] Implementar CancelLease:
  - Validar lease existe
  - Atualizar lease.status = cancelled
  - Atualizar unit.status = available
- [x] Implementar UpdatePaintingFeePaid
- [x] Implementar CheckExpiringSoonLeases (cronjob futuro)
- [x] Implementar MarkLeaseAsExpired
- [x] Implementar GetLeaseStats
- [x] Adicionar testes completos com mocks

### 2.6 Service - Lease (Parte 3: Renovação)
- [x] Implementar RenewLease:
  - Validar lease existe e está ativo ou expiring_soon
  - Criar novo lease com start_date = old.end_date + 1 dia
  - Calcular novo end_date (+ 6 meses)
  - Manter mesmo unit_id, tenant_id, payment_due_day
  - Usar valor atualizado da unidade (monthly_rent_value)
  - Nova taxa de pintura
  - Atualizar lease antigo para status = expired
  - Retornar novo lease
- [x] CheckExpiringSoonLeases já implementado na Task 2.5
- [x] Adicionar testes de renovação

### 2.7 Handler - Lease (Parte 1: CRUD Básico)
- [x] Criar arquivo `internal/handler/lease_handler.go`
- [x] Criar DTOs (CreateLeaseRequestDTO, LeaseResponse, etc)
- [x] Implementar CreateLease handler (POST /api/leases)
- [x] Implementar GetLease handler (GET /api/leases/:id)
- [x] Implementar ListLeases handler (GET /api/leases)
- [x] Adicionar query params para filtros (status, unit_id, tenant_id)
- [x] Implementar GetLeaseStats handler (GET /api/leases/stats)

### 2.8 Handler - Lease (Parte 2: Operações Especiais)
- [x] Implementar RenewLease handler (POST /api/leases/:id/renew)
- [x] Implementar CancelLease handler (POST /api/leases/:id/cancel)
- [x] Implementar UpdatePaintingFeePaid handler (PATCH /api/leases/:id/painting-fee)
- [x] Implementar GetExpiringSoonLeases handler (GET /api/leases/expiring-soon)
- [x] Validar inputs em todos os handlers
- [x] Mapear erros do service para HTTP status codes

### 2.9 Router e Testes
- [x] Registrar todas as rotas de leases no router.go
- [x] Atualizar main.go com LeaseRepository e LeaseService
- [x] Atualizar SetupRoutes com LeaseHandler
- [x] Adicionar tag @tag.name Leases no Swagger
- [x] Testar criação de contrato manualmente
- [x] Testar cancelamento
- [x] Testar renovação
- [x] Testar filtros de listagem
- [x] Verificar alteração de status das unidades
- [x] Testar atualização de taxa de pintura
- [x] Testar validações de negócio

### 2.10 Documentação e Commit
- [x] Gerar documentação Swagger
- [x] Corrigir mapeamento de erros (ErrPaintingFeePaidExceedsTotal)
- [x] Testar todos os endpoints via Swagger/cURL
- [x] Validar regras de negócio implementadas
- [x] Commit final: "feat: complete lease management system"
- [x] Push para repositório

---

## Sprint 3: Sistema de Pagamentos
**Duração:** 4-5 dias  
**Objetivo:** Implementar controle completo de pagamentos

### 3.1 Migration e Schema - Payments
- [x] Criar migration `000004_create_payments_table.up.sql`
- [x] Definir tabela payments
- [x] Adicionar foreign key para leases
- [x] Adicionar checks em payment_type
- [x] Criar índices (lease_id, status, due_date)
- [x] Criar índice composto (status, due_date) para queries de vencimento
- [x] Executar migration

### 3.2 Domain Model - Payment
- [x] Criar arquivo `internal/domain/payment.go`
- [x] Definir struct Payment
- [x] Definir enums: PaymentType, PaymentStatus, PaymentMethod
- [x] Implementar método IsOverdue()
- [x] Implementar método CanBePaid()
- [x] Implementar método MarkAsPaid()
- [x] Adicionar testes unitários

### 3.3 Repository - Payment (SQLC)
- [x] Criar arquivo `internal/repository/queries/payments.sql`
- [x] Query CreatePayment
- [x] Query GetPaymentByID
- [x] Query ListPaymentsByLeaseID
- [x] Query ListPaymentsByStatus
- [x] Query GetOverduePayments (status=pending AND due_date < now)
- [x] Query GetUpcomingPayments (due_date BETWEEN now AND now+X days)
- [x] Query UpdatePayment
- [x] Query UpdatePaymentStatus
- [x] Query MarkAsPaid
- [x] Gerar código e implementar repository

### 3.4 Service - Payment (Parte 1: Geração)
- [x] Criar arquivo `internal/service/payment_service.go`
- [x] Implementar GenerateMonthlyRentPayment:
  - Receber lease_id, reference_month
  - Calcular due_date baseado em payment_due_day
  - Criar Payment tipo "rent"
  - Amount = lease.monthly_rent_value
- [x] Implementar GeneratePaintingFeePayments:
  - Receber lease, installments
  - Se installments=1: 1 payment com amount=total
  - Se installments=2 ou 3: dividir amount igualmente
  - Calcular due_dates escalonadas
  - Criar múltiplos Payments
- [x] Implementar GenerateAdjustmentPayment (proporcional)
- [x] Adicionar testes

### 3.5 Service - Payment (Parte 2: Registro)
- [x] Implementar MarkPaymentAsPaid:
  - Validar payment existe e está pending/overdue
  - Atualizar payment_date, status=paid, payment_method
  - Se type=painting_fee: atualizar lease.painting_fee_paid
  - Retornar payment atualizado
- [x] Implementar GetPaymentsByLease
- [x] Implementar GetOverduePayments
- [x] Implementar GetPaymentsDueSoon (próximos X dias)
- [x] Adicionar testes

### 3.6 Service - Payment (Parte 3: Cronjob)
- [x] Implementar CheckOverduePayments:
  - Buscar payments com status=pending e due_date < hoje
  - Atualizar status para overdue
  - Retornar quantidade atualizada
- [x] Adicionar lógica para ser executado diariamente (scheduler futuro)

### 3.7 Integração Lease + Payment na Criação de Contrato
- [ ] Atualizar LeaseService.CreateLease:
  - Após criar lease, gerar primeiro pagamento de aluguel
  - Gerar pagamentos de taxa de pintura (1x ou 3x)
  - Retornar lease + lista de payments criados
- [ ] Atualizar LeaseHandler.CreateLease:
  - Retornar no response os payments gerados
- [ ] Testar criação de contrato com geração automática de pagamentos

### 3.8 Handler - Payment
- [ ] Criar arquivo `internal/handler/payment_handler.go`
- [ ] Criar DTOs (MarkPaymentAsPaidRequest, PaymentResponse)
- [ ] Implementar GetPaymentsByLease (GET /api/leases/:id/payments)
- [ ] Implementar GetPayment (GET /api/payments/:id)
- [ ] Implementar MarkAsPaid (PUT /api/payments/:id/pay)
- [ ] Implementar ListOverduePayments (GET /api/payments/overdue)
- [ ] Implementar ListUpcomingPayments (GET /api/payments/upcoming)

### 3.9 Router e Testes Manuais
- [ ] Registrar rotas de payments
- [ ] Testar criação de contrato e verificar payments gerados
- [ ] Testar marcação de pagamento como pago
- [ ] Testar listagem de atrasados
- [ ] Testar listagem de próximos vencimentos
- [ ] Verificar atualização de painting_fee_paid no lease

### 3.10 Documentação e Commit
- [ ] Documentar endpoints de payments
- [ ] Documentar lógica de geração de pagamentos
- [ ] Adicionar exemplos no README
- [ ] Commit: "feat: implement payment management"
- [ ] Push para repositório

---

## Sprint 4: Dashboard e Relatórios
**Duração:** 3-4 dias  
**Objetivo:** Criar visão executiva e relatórios financeiros

### 4.1 Service - Dashboard (Métricas Gerais)
- [ ] Criar arquivo `internal/service/dashboard_service.go`
- [ ] Implementar GetOccupancyMetrics:
  - Total de unidades
  - Unidades ocupadas
  - Unidades disponíveis
  - Unidades em manutenção/reforma
  - Taxa de ocupação (%)
- [ ] Implementar GetFinancialMetrics:
  - Receita mensal projetada (soma de todos alugueis ativos)
  - Receita mensal realizada (pagamentos pagos no mês)
  - Inadimplência (pagamentos overdue)
  - Taxa de inadimplência (%)
- [ ] Adicionar testes

### 4.2 Service - Dashboard (Contratos e Alertas)
- [ ] Implementar GetContractMetrics:
  - Total de contratos ativos
  - Contratos expirando em 45 dias
  - Contratos expirados
- [ ] Implementar GetAlerts:
  - Lista de pagamentos atrasados
  - Lista de contratos expirando
  - Unidades sem contrato há muito tempo
- [ ] Adicionar testes

### 4.3 Service - Reports (Relatório Financeiro)
- [ ] Criar arquivo `internal/service/report_service.go`
- [ ] Implementar GetFinancialReport:
  - Filtros: start_date, end_date, payment_type
  - Receita total por tipo (rent, painting_fee)
  - Receita por mês
  - Detalhamento por unidade
  - Retornar estrutura agregada
- [ ] Implementar GetPaymentHistoryReport:
  - Histórico completo de pagamentos
  - Filtros: lease_id, tenant_id, status, date_range
- [ ] Adicionar testes

### 4.4 Handler - Dashboard
- [ ] Criar arquivo `internal/handler/dashboard_handler.go`
- [ ] Criar DTO DashboardResponse com todas as métricas
- [ ] Implementar GetDashboard (GET /api/dashboard)
- [ ] Consolidar dados de múltiplos services
- [ ] Retornar JSON estruturado

### 4.5 Handler - Reports
- [ ] Criar arquivo `internal/handler/report_handler.go`
- [ ] Criar DTOs para requests e responses
- [ ] Implementar GetFinancialReport (GET /api/reports/financial)
- [ ] Adicionar query params para filtros
- [ ] Implementar GetPaymentHistory (GET /api/reports/payments)
- [ ] Validar filtros de data

### 4.6 Queries SQL Otimizadas
- [ ] Criar queries agregadas no SQLC para dashboard
- [ ] Query para receita mensal agrupada
- [ ] Query para contagem de unidades por status
- [ ] Query para pagamentos atrasados com detalhes
- [ ] Gerar código e testar performance

### 4.7 Router e Testes
- [ ] Registrar rotas de dashboard e reports
- [ ] Testar dashboard com dados reais
- [ ] Verificar cálculos de taxas e percentuais
- [ ] Testar filtros de relatórios
- [ ] Validar formato JSON das respostas

### 4.8 Melhorias de Performance (Opcional)
- [ ] Adicionar cache em memória para dashboard (5 minutos)
- [ ] Implementar paginação nos relatórios
- [ ] Adicionar campo de ordenação

### 4.9 Documentação e Commit
- [ ] Documentar endpoints de dashboard
- [ ] Documentar estrutura dos relatórios
- [ ] Adicionar exemplos de responses
- [ ] Commit: "feat: implement dashboard and reports"
- [ ] Push para repositório

---

## Sprint 5: Autenticação e Autorização
**Duração:** 2-3 dias
**Objetivo:** Implementar sistema de autenticação JWT e proteção de rotas

### 5.1 Migration e Schema - Users
- [ ] Criar migration `000005_create_users_table.up.sql`
- [ ] Definir tabela users (id, username, password_hash, role, created_at)
- [ ] Adicionar constraint UNIQUE no username
- [ ] Criar índice no username
- [ ] Seed inicial com usuário admin
- [ ] Executar migration

### 5.2 Domain Model - User
- [ ] Criar arquivo `internal/domain/user.go`
- [ ] Definir struct User
- [ ] Definir enum UserRole (admin, manager, viewer)
- [ ] Implementar método ValidatePassword()
- [ ] Implementar método HashPassword()
- [ ] Adicionar testes unitários

### 5.3 Repository - User (SQLC)
- [ ] Criar arquivo `internal/repository/queries/users.sql`
- [ ] Query GetUserByUsername
- [ ] Query CreateUser
- [ ] Query UpdateUser
- [ ] Gerar código SQLC
- [ ] Implementar UserRepository

### 5.4 Service - Auth
- [ ] Criar arquivo `internal/service/auth_service.go`
- [ ] Instalar dependências: `golang-jwt/jwt` e `golang.org/x/crypto/bcrypt`
- [ ] Implementar GenerateToken (JWT)
- [ ] Implementar ValidateToken
- [ ] Implementar Login(username, password)
- [ ] Implementar RefreshToken (opcional)
- [ ] Adicionar testes

### 5.5 Handler - Auth
- [ ] Criar arquivo `internal/handler/auth_handler.go`
- [ ] Criar DTOs (LoginRequest, LoginResponse, TokenResponse)
- [ ] Implementar Login handler (POST /api/auth/login)
- [ ] Implementar Refresh handler (POST /api/auth/refresh) - opcional
- [ ] Implementar GetCurrentUser (GET /api/auth/me)

### 5.6 Middleware - Authentication
- [ ] Criar arquivo `internal/pkg/middleware/auth.go`
- [ ] Implementar AuthMiddleware:
  - Extrair token do header Authorization
  - Validar token JWT
  - Adicionar user info no context
  - Retornar 401 se inválido
- [ ] Implementar RequireRole(roles ...string)
- [ ] Adicionar testes

### 5.7 Proteger Rotas Existentes
- [ ] Atualizar router.go para aplicar AuthMiddleware
- [ ] Proteger todas as rotas de /api/v1/*
- [ ] Deixar /health e /swagger públicos
- [ ] Deixar /api/auth/login público
- [ ] Testar autenticação em todas as rotas

### 5.8 Atualizar Swagger
- [ ] Adicionar securityDefinitions no main.go
- [ ] Adicionar @Security tags nos handlers
- [ ] Regenerar documentação Swagger
- [ ] Testar autenticação via Swagger UI

### 5.9 Testes e Documentação
- [ ] Testar fluxo de login
- [ ] Testar acesso sem token (401)
- [ ] Testar token expirado
- [ ] Testar token inválido
- [ ] Documentar processo de autenticação
- [ ] Commit: "feat: implement JWT authentication and authorization"
- [ ] Push para repositório

---

## Sprint 6: Sistema de Notificações
**Duração:** 2-3 dias
**Objetivo:** Implementar lembretes e alertas internos

### 6.1 Migration e Schema - Notifications
- [ ] Criar migration `000006_create_notifications_table.up.sql`
- [ ] Definir tabela notifications
- [ ] Adicionar foreign keys (lease_id, tenant_id)
- [ ] Criar índices (status, scheduled_date)
- [ ] Executar migration

### 6.2 Domain Model - Notification
- [ ] Criar arquivo `internal/domain/notification.go`
- [ ] Definir struct Notification
- [ ] Definir enums: NotificationType, NotificationStatus
- [ ] Implementar método IsReadyToSend()
- [ ] Implementar método MarkAsSent()
- [ ] Adicionar testes

### 6.3 Repository - Notification (SQLC)
- [ ] Criar arquivo `internal/repository/queries/notifications.sql`
- [ ] Query CreateNotification
- [ ] Query GetNotificationByID
- [ ] Query ListPendingNotifications (status=pending, scheduled_date <= now)
- [ ] Query ListNotificationsByLease
- [ ] Query UpdateNotificationStatus
- [ ] Query MarkAsSent
- [ ] Gerar código e implementar

### 6.4 Service - Notification (Criação)
- [ ] Criar arquivo `internal/service/notification_service.go`
- [ ] Implementar CreateRentReminderNotification:
  - Receber lease_id, due_date
  - Calcular scheduled_date = due_date - 3 dias
  - Gerar message_content com dados do morador
  - Criar notification tipo sms_rent_reminder
- [ ] Implementar CreateContractExpiringNotification:
  - Receber lease_id
  - Scheduled_date = end_date - 45 dias
  - Message_content personalizada
- [ ] Adicionar testes

### 6.5 Service - Notification (Processamento)
- [ ] Implementar ProcessDailyNotifications:
  - Buscar leases ativos
  - Para cada lease, buscar payments pendentes/overdue
  - Se payment.due_date = hoje + 3 dias e não existe notificação:
    - Criar notificação de lembrete
  - Buscar leases com end_date = hoje + 45 dias:
    - Criar notificação de contrato expirando
  - Retornar quantidade de notificações criadas
- [ ] Implementar GetPendingNotifications
- [ ] Implementar MarkNotificationAsSent
- [ ] Adicionar testes

### 6.6 Handler - Notification
- [ ] Criar arquivo `internal/handler/notification_handler.go`
- [ ] Criar DTOs necessários
- [ ] Implementar ListNotifications (GET /api/notifications)
- [ ] Implementar GetNotificationsByLease (GET /api/leases/:id/notifications)
- [ ] Implementar MarkAsSent (PUT /api/notifications/:id/mark-sent)
- [ ] Implementar TriggerDailyProcessing (POST /api/notifications/process) - endpoint administrativo

### 6.7 Scheduler/Cronjob Básico
- [ ] Criar arquivo `internal/pkg/scheduler/scheduler.go`
- [ ] Implementar função DailyNotificationJob:
  - Executar NotificationService.ProcessDailyNotifications()
  - Executar PaymentService.CheckOverduePayments()
  - Executar LeaseService.CheckExpiringSoonLeases()
  - Logar resultados
- [ ] Integrar scheduler no main.go (executar a cada X horas ou usar time.Ticker)
- [ ] Adicionar flag de enable/disable via config

### 6.8 Router e Testes
- [ ] Registrar rotas de notifications
- [ ] Testar criação manual de notificação
- [ ] Testar processamento diário (endpoint /process)
- [ ] Verificar notificações sendo criadas automaticamente
- [ ] Testar marcação como enviada

### 6.9 Logs e Monitoramento
- [ ] Adicionar logs detalhados no scheduler
- [ ] Logar quantidade de notificações processadas
- [ ] Logar erros de processamento
- [ ] Adicionar métricas básicas (opcional)

### 6.10 Documentação e Commit
- [ ] Documentar sistema de notificações
- [ ] Documentar cronjob e como funciona
- [ ] Adicionar instruções para teste manual
- [ ] Commit: "feat: implement notification system"
- [ ] Push para repositório

---

## Sprint 7: Refinamentos e MVP Final
**Duração:** 3-4 dias
**Objetivo:** Polir aplicação e preparar para uso

### 7.1 Tratamento de Erros Global
- [ ] Criar middleware de error handling
- [ ] Padronizar responses de erro (código, mensagem, detalhes)
- [ ] Implementar error types customizados
- [ ] Adicionar logging de erros

### 6.2 Validações Completas
- [ ] Revisar todas as validações de input nos handlers
- [ ] Adicionar validações de regras de negócio nos services
- [ ] Testar casos extremos e edge cases
- [ ] Documentar validações aplicadas

### 6.3 Middlewares
- [ ] Implementar middleware de logging (request/response)
- [ ] Implementar middleware de CORS
- [ ] Implementar middleware de recovery (panic)
- [ ] Implementar middleware de timeout (opcional)
- [ ] Aplicar middlewares no router

### 6.4 Configuração Avançada
- [ ] Externalizar todas as configurações hardcoded
- [ ] Criar arquivo config.yaml (opcional)
- [ ] Validar variáveis de ambiente obrigatórias
- [ ] Documentar todas as variáveis no .env.example

### 6.5 Testes Finais
- [ ] Executar suite completa de testes
- [ ] Testar todos os endpoints manualmente
- [ ] Criar collection do Postman/Insomnia (opcional)
- [ ] Testar fluxo completo end-to-end:
  - Criar unidade
  - Criar morador
  - Criar contrato
  - Gerar pagamentos
  - Registrar pagamento
  - Ver dashboard
  - Gerar relatórios

### 6.6 Documentação Final
- [ ] Atualizar README com instruções completas
- [ ] Criar documentação da API (endpoints, requests, responses)
- [ ] Adicionar exemplos de uso
- [ ] Criar guia de troubleshooting
- [ ] Documentar próximos passos e melhorias futuras

### 6.7 Preparação para Deploy (Futuro)
- [ ] Criar Dockerfile (básico)
- [ ] Adicionar health check endpoint (GET /health)
- [ ] Configurar variáveis para produção
- [ ] Documentar processo de deploy

### 6.8 Revisão de Código
- [ ] Revisar código de todos os módulos
- [ ] Refatorar duplicações
- [ ] Melhorar nomenclaturas se necessário
- [ ] Adicionar comentários onde necessário
- [ ] Executar linter (golangci-lint)

### 6.9 Commit Final do MVP
- [ ] Atualizar CHANGELOG (criar arquivo)
- [ ] Criar tag de versão v1.0.0
- [ ] Commit: "chore: MVP v1.0.0 release"
- [ ] Push com tags

### 6.10 Celebração! 🎉
- [ ] Fazer backup do banco de dados
- [ ] Documentar aprendizados
- [ ] Planejar próximas features (v2.0)

---

## 🚀 Próximas Versões (Pós-MVP)

### Versão 2.0 - Integrações e Automações
- Integração com gateway de SMS
- Geração automática de cobranças mensais
- Upload e armazenamento de comprovantes
- Exportação de relatórios (PDF/Excel)
- Integração com PIX para confirmação automática

### Versão 3.0 - Portal do Morador
- Autenticação e autorização
- Login para moradores
- Visualização de pagamentos e contratos
- Download de comprovantes
- Histórico completo

### Versão 4.0 - Analytics Avançado
- Dashboard avançado com gráficos
- Previsões de receita
- Análise de inadimplência
- KPIs e métricas de negócio
- Relatórios customizáveis

---

## 📊 Estimativas

| Sprint | Duração | Complexidade |
|--------|---------|--------------|
| Sprint 0 | 2-3 dias | Baixa |
| Sprint 1 | 3-4 dias | Média |
| Sprint 2 | 4-5 dias | Alta |
| Sprint 3 | 4-5 dias | Alta |
| Sprint 4 | 3-4 dias | Média |
| Sprint 5 | 2-3 dias | Média |
| Sprint 6 | 2-3 dias | Média |
| Sprint 7 | 3-4 dias | Baixa |
| **TOTAL** | **9-11 semanas** | - |

---

**Última atualização:** Setembro 2025
