# Kitnet Manager 🏢

Sistema completo de gestão para administração de complexo de 31 kitnets, substituindo controles manuais em Excel por uma solução digital robusta e escalável.

[![Production](https://img.shields.io/badge/production-online-success)](https://kitnet-manager-production.up.railway.app)
[![API Docs](https://img.shields.io/badge/docs-swagger-blue)](https://kitnet-manager-production.up.railway.app/swagger/index.html)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

---

## 📋 Índice

- [Sobre o Projeto](#sobre-o-projeto)
- [Status](#status)
- [Funcionalidades](#funcionalidades)
- [Tecnologias](#tecnologias)
- [Arquitetura](#arquitetura)
- [Getting Started](#getting-started)
- [Documentação](#documentação)
- [Deploy](#deploy)
- [Comandos](#comandos)
- [Contribuindo](#contribuindo)
- [Licença](#licença)

---

## 🎯 Sobre o Projeto

O **Kitnet Manager** é um sistema de gestão imobiliária focado em administração de pequenos complexos residenciais. O sistema gerencia:

- **31 Unidades** (kitnets/apartamentos)
- **Contratos de locação** com duração fixa de 6 meses
- **Pagamentos mensais** com rastreamento de inadimplência
- **Renovação automática** de contratos
- **Taxa de pintura** parcelável em até 4x
- **Dashboard executivo** com métricas em tempo real
- **Relatórios financeiros** personalizados

### Problema Resolvido

Substituição de planilhas Excel dispersas por sistema centralizado com:
- ✅ Rastreamento automatizado de pagamentos
- ✅ Alertas de vencimento de contratos
- ✅ Cálculo automático de multas e juros
- ✅ Histórico completo de transações
- ✅ Métricas de ocupação e receita

---

## 🚀 Status

**Versão:** 1.0.0
**Status:** ✅ **Em Produção**
**URL:** https://kitnet-manager-production.up.railway.app

### Sprints Concluídas

- [x] **Sprint 0:** Setup inicial, database schema, migrations
- [x] **Sprint 1:** CRUD de Units e Tenants
- [x] **Sprint 2:** Gestão de contratos (Leases)
- [x] **Sprint 3:** Sistema de pagamentos
- [x] **Sprint 4:** Dashboard e relatórios financeiros
- [x] **Sprint 5:** Sistema de autenticação e autorização
- [x] **Deploy:** Produção na Railway

### Próximas Features

- [ ] **Sprint 6:** Notificações SMS (Twilio)
- [ ] **Sprint 7:** Exportação de relatórios (PDF/Excel)
- [ ] **Sprint 8:** Geração de contratos em PDF
- [ ] **Sprint 9:** Mobile responsiveness
- [ ] **Sprint 10:** Testes e2e, refinamentos finais

---

## ✨ Funcionalidades

### 🏠 Gestão de Unidades
- CRUD completo de unidades
- Status dinâmico (disponível, ocupada, manutenção, reforma)
- Controle de valores (base vs. reformado)
- Estatísticas de ocupação

### 👥 Gestão de Inquilinos
- Cadastro com validação de CPF
- Busca por nome ou CPF
- Documentos de identificação
- Histórico de locações

### 📝 Contratos de Locação
- Criação automática de cronograma de pagamentos
- Duração fixa de 6 meses
- Renovação de contratos
- Taxa de pintura parcelável (1-4x)
- Alertas de vencimento (45 dias)

### 💰 Sistema de Pagamentos
- Rastreamento de pagamentos mensais
- Múltiplos métodos (PIX, dinheiro, transferência, cartão)
- Detecção automática de atrasos
- Cálculo de multas (2% + 1% ao mês)
- Estatísticas por contrato

### 📊 Dashboard Executivo
- Métricas de ocupação em tempo real
- Receita mensal consolidada
- Pagamentos pendentes e atrasados
- Alertas inteligentes
- Contratos expirando

### 📈 Relatórios Financeiros
- Relatório por período customizável
- Filtros por tipo e status de pagamento
- Histórico completo de pagamentos
- Exportação de dados

### 🔐 Sistema de Autenticação
- Login com JWT
- Refresh token
- Controle de acesso baseado em roles
  - **Admin:** Acesso total + gestão de usuários
  - **Manager:** Leitura e escrita
  - **Viewer:** Apenas leitura

---

## 🛠️ Tecnologias

### Backend
- **[Go 1.21+](https://go.dev)** - Linguagem principal
- **[Chi Router](https://github.com/go-chi/chi)** - HTTP router leve e idiomático
- **[PostgreSQL 17.5](https://www.postgresql.org/)** - Banco de dados relacional
- **[SQLC](https://sqlc.dev/)** - Geração de código type-safe a partir de SQL
- **[golang-migrate](https://github.com/golang-migrate/migrate)** - Database migrations
- **[JWT](https://jwt.io/)** - Autenticação stateless

### Principais Dependências

```go
require (
    github.com/go-chi/chi/v5       // HTTP router
    github.com/go-chi/cors          // CORS middleware
    github.com/lib/pq               // PostgreSQL driver
    github.com/google/uuid          // UUID generation
    github.com/shopspring/decimal   // Monetary precision
    github.com/go-playground/validator/v10  // Validation
    github.com/golang-jwt/jwt/v5    // JWT tokens
    golang.org/x/crypto/bcrypt      // Password hashing
)
```

### Infraestrutura
- **[Railway](https://railway.app)** - Deploy e hosting
- **[Neon](https://neon.tech)** - PostgreSQL cloud
- **[Swagger](https://swagger.io/)** - Documentação interativa da API

### Frontend (Planejado)
- **Next.js 14+** - React framework
- **TypeScript** - Type safety
- **TailwindCSS** - Styling
- **React Query** - Data fetching

---

## 🏗️ Arquitetura

O projeto segue **Clean Architecture** com separação clara de responsabilidades:

```
kitnet-manager/
├── cmd/
│   └── api/                    # Application entry point
│       └── main.go             # Server initialization
│
├── internal/
│   ├── domain/                 # Business entities (pure Go)
│   │   ├── unit.go             # Unit entity + business rules
│   │   ├── tenant.go           # Tenant entity + validations
│   │   ├── lease.go            # Lease entity + calculations
│   │   ├── payment.go          # Payment entity + status
│   │   └── user.go             # User entity + auth
│   │
│   ├── repository/             # Data access layer
│   │   ├── interfaces.go       # Repository interfaces
│   │   └── postgres/           # PostgreSQL implementation
│   │       ├── unit_repository.go
│   │       ├── tenant_repository.go
│   │       ├── lease_repository.go
│   │       ├── payment_repository.go
│   │       └── user_repository.go
│   │
│   ├── service/                # Business logic layer
│   │   ├── unit_service.go
│   │   ├── tenant_service.go
│   │   ├── lease_service.go
│   │   ├── payment_service.go
│   │   ├── dashboard_service.go
│   │   ├── report_service.go
│   │   └── auth_service.go
│   │
│   ├── handler/                # HTTP handlers (controllers)
│   │   ├── routes.go           # Route registration
│   │   ├── unit_handler.go
│   │   ├── tenant_handler.go
│   │   ├── lease_handler.go
│   │   ├── payment_handler.go
│   │   ├── dashboard_handler.go
│   │   ├── report_handler.go
│   │   └── auth_handler.go
│   │
│   ├── pkg/                    # Internal packages
│   │   ├── database/           # DB connection & health
│   │   ├── middleware/         # Auth middleware
│   │   ├── response/           # Standardized responses
│   │   └── validator/          # Custom validators
│   │
│   └── config/                 # Configuration loading
│       └── config.go
│
├── migrations/                 # Database migrations
│   ├── 000001_create_units.up.sql
│   ├── 000002_create_tenants.up.sql
│   └── ...
│
├── docs/                       # Generated Swagger docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── frontend-docs/              # Frontend API documentation
│   ├── README.md               # Documentation index
│   ├── API.md                  # API overview
│   ├── types/                  # TypeScript definitions
│   ├── endpoints/              # Endpoint documentation
│   ├── validation-rules.md     # Business rules
│   └── examples.md             # Code examples
│
├── config/                     # Configuration files
│   └── .env.example
│
├── Makefile                    # Development commands
├── go.mod                      # Go dependencies
├── go.sum
├── railway.json                # Railway config
└── README.md                   # This file
```

### Princípios Arquiteturais

1. **Dependency Inversion:** Handlers → Services → Repositories (interfaces)
2. **Domain-Driven Design:** Business logic isolado na camada de domínio
3. **Repository Pattern:** Abstração de acesso a dados
4. **Clean Architecture:** Separação de concerns em camadas
5. **Type Safety:** SQLC para queries SQL type-safe

---

## 🚀 Getting Started

### Pré-requisitos

- [Go 1.21+](https://go.dev/dl/)
- [PostgreSQL 17+](https://www.postgresql.org/download/) ou conta no [Neon](https://neon.tech)
- [Make](https://www.gnu.org/software/make/) (opcional, mas recomendado)
- [migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) (para migrations)

### Instalação

1. **Clone o repositório:**
   ```bash
   git clone https://github.com/lucianogabriel/kitnet-manager.git
   cd kitnet-manager
   ```

2. **Instale as dependências:**
   ```bash
   go mod download
   ```

3. **Configure as variáveis de ambiente:**
   ```bash
   cp config/.env.example .env
   # Edite o .env com suas configurações
   ```

4. **Configure o banco de dados:**
   ```bash
   # Opção 1: Usar Make
   make db-setup

   # Opção 2: Manualmente
   make migrate-up
   ```

5. **Execute a aplicação:**
   ```bash
   # Modo desenvolvimento (com hot reload)
   make dev

   # Ou modo normal
   make run
   ```

6. **Acesse a aplicação:**
   - API: http://localhost:8080
   - Swagger: http://localhost:8080/swagger/index.html
   - Health: http://localhost:8080/health

### Configuração do Banco de Dados

#### Usando Neon (Cloud PostgreSQL)

1. Crie uma conta em [neon.tech](https://neon.tech)
2. Crie um novo projeto
3. Copie a connection string
4. Configure no `.env`:
   ```env
   DATABASE_URL=postgresql://user:password@host/database?sslmode=require
   ```

#### Usando PostgreSQL Local

```bash
# Criar database
createdb kitnet_manager

# Configurar .env
DATABASE_URL=postgresql://localhost/kitnet_manager?sslmode=disable
```

---

## 📚 Documentação

### Documentação da API

- **Swagger UI (Interativo):** https://kitnet-manager-production.up.railway.app/swagger/index.html
- **Frontend Docs:** [frontend-docs/](./frontend-docs/)
  - [README.md](./frontend-docs/README.md) - Visão geral
  - [API.md](./frontend-docs/API.md) - Guia da API
  - [types/](./frontend-docs/types/) - TypeScript types
  - [endpoints/](./frontend-docs/endpoints/) - Documentação detalhada
  - [examples.md](./frontend-docs/examples.md) - Exemplos de código

### Documentação do Projeto

- [Arquitetura Detalhada](./kitnet_architecture.md)
- [Roadmap e Sprints](./kitnet_roadmap.md)
- [Regras de Negócio](./BUSINESS_RULES.md)
- [Deploy e Produção](./DEPLOY.md)
- [Documentação de Produção](./PRODUCTION.md)

### Credenciais Padrão

```
Username: admin
Password: admin123
```

**⚠️ IMPORTANTE:** Altere essas credenciais em produção!

---

## 🚢 Deploy

### Produção (Railway)

O projeto está configurado para deploy automático na Railway:

1. **Push para main:**
   ```bash
   git push origin main
   ```

2. **Railway faz:**
   - Build automático
   - Deploy
   - Health checks
   - HTTPS automático

3. **Acompanhe:**
   - Dashboard: https://railway.app
   - Logs em tempo real
   - Métricas de uso

### Configuração de Produção

Variáveis de ambiente necessárias na Railway:

```env
DATABASE_URL=postgresql://...
PORT=8080
JWT_SECRET=your-secret-key-here
JWT_EXPIRY=24h
ENVIRONMENT=production
```

Veja [DEPLOY.md](./DEPLOY.md) para instruções completas.

---

## 🔧 Comandos

### Aplicação

```bash
make run           # Executar aplicação
make build         # Compilar binário
make dev           # Modo desenvolvimento (hot reload)
make test          # Executar testes
make clean         # Limpar arquivos gerados
make help          # Listar todos os comandos
```

### Banco de Dados

```bash
make migrate-create name=nome_da_migration  # Criar nova migration
make migrate-up                             # Aplicar migrations
make migrate-down                           # Reverter última migration
make migrate-status                         # Status das migrations
make db-setup                               # Setup completo
```

### Desenvolvimento

```bash
make swagger       # Gerar documentação Swagger
make lint          # Executar linter
make fmt           # Formatar código
```

### Workflow de Desenvolvimento

1. **Criar feature:**
   ```bash
   git checkout -b feature/nome-da-feature
   ```

2. **Desenvolver:**
   ```bash
   make dev  # Roda com hot reload
   ```

3. **Testar:**
   ```bash
   make test
   ```

4. **Criar migration se necessário:**
   ```bash
   make migrate-create name=add_new_field
   # Edite o arquivo SQL criado
   make migrate-up
   ```

5. **Commit e push:**
   ```bash
   git add .
   git commit -m "feat: descrição da feature"
   git push origin feature/nome-da-feature
   ```

---

## 🤝 Contribuindo

Este é um projeto educacional/pessoal, mas contribuições são bem-vindas!

### Como Contribuir

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'feat: Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

### Padrões de Commit

Seguimos [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: nova funcionalidade
fix: correção de bug
docs: mudanças na documentação
refactor: refatoração de código
test: adição de testes
chore: tarefas de manutenção
```

---

## 📊 Métricas do Projeto

- **Linhas de Código:** ~8,000+
- **Endpoints:** 41
- **Entidades de Domínio:** 6
- **Tabelas no Banco:** 6
- **Migrations:** 6
- **Testes:** Em desenvolvimento
- **Coverage:** Meta de 80%

---

## 🔐 Segurança

- ✅ Autenticação JWT
- ✅ Passwords hasheados com bcrypt
- ✅ Validação de entrada em todos os endpoints
- ✅ SQL injection protegido (queries parametrizadas)
- ✅ CORS configurado
- ✅ HTTPS em produção
- ✅ Rate limiting (planejado)

---

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

## 👨‍💻 Autor

**Luciano Gabriel**

- GitHub: [@lucianoZgabriel](https://github.com/lucianoZgabriel)
- LinkedIn: [Luciano Gabriel](https://linkedin.com/in/lucianogabriel)

---

## 🙏 Agradecimentos

- Comunidade Go pela excelente documentação
- Chi Router pela simplicidade
- Railway pela facilidade de deploy
- Neon pelo PostgreSQL cloud gratuito

---

## 📞 Suporte

- **Issues:** [GitHub Issues](https://github.com/lucianogabriel/kitnet-manager/issues)
- **Documentação:** [Frontend Docs](./frontend-docs/)
- **API Docs:** [Swagger UI](https://kitnet-manager-production.up.railway.app/swagger/index.html)

---

<div align="center">

**[⬆ Voltar ao topo](#kitnet-manager-)**

Made with ❤️ and Go

</div>
