# Arquitetura - Kitnet Manager

> Documentação técnica da arquitetura, decisões de design e padrões utilizados

## 📐 Visão Geral da Arquitetura

O Kitnet Manager segue os princípios da **Clean Architecture** adaptada para Go, com foco em:
- Separação clara de responsabilidades
- Independência de frameworks
- Testabilidade
- Manutenibilidade
- Escalabilidade futura

## 🏗 Padrões Arquiteturais

### Clean Architecture (Adaptada)

```
┌─────────────────────────────────────────┐
│         HTTP Handlers (API)             │  ← Camada de Entrada
├─────────────────────────────────────────┤
│           Services (Use Cases)          │  ← Lógica de Negócio
├─────────────────────────────────────────┤
│        Repositories (Interface)         │  ← Abstração de Dados
├─────────────────────────────────────────┤
│      Repository Implementation          │  ← Implementação Concreta
│         (PostgreSQL/SQLC)               │
├─────────────────────────────────────────┤
│             Domain Models               │  ← Entidades Core
└─────────────────────────────────────────┘
```

### Camadas

#### 1. Domain Layer (Domínio)
- **Responsabilidade:** Definir as entidades core e regras de negócio fundamentais
- **Localização:** `internal/domain/`
- **Características:**
  - Zero dependências externas
  - Modelos ricos com métodos de validação
  - Enums e constantes do domínio
  
**Exemplo:**
```go
type LeaseStatus string

const (
    LeaseStatusActive        LeaseStatus = "active"
    LeaseStatusExpiringsSoon LeaseStatus = "expiring_soon"
    LeaseStatusExpired       LeaseStatus = "expired"
    LeaseStatusCancelled     LeaseStatus = "cancelled"
)

type Lease struct {
    ID                      uuid.UUID
    UnitID                  uuid.UUID
    TenantID                uuid.UUID
    ContractSignedDate      time.Time
    StartDate               time.Time
    EndDate                 time.Time
    PaymentDueDay           int
    MonthlyRentValue        decimal.Decimal
    PaintingFeeTotal        decimal.Decimal
    PaintingFeeInstallments int
    PaintingFeePaid         decimal.Decimal
    Status                  LeaseStatus
    CreatedAt               time.Time
    UpdatedAt               time.Time
}

func (l *Lease) IsExpiringSoon() bool {
    daysUntilExpiry := int(time.Until(l.EndDate).Hours() / 24)
    return daysUntilExpiry <= 45 && daysUntilExpiry > 0
}
```

#### 2. Repository Layer (Repositórios)
- **Responsabilidade:** Abstração de acesso aos dados
- **Localização:** `internal/repository/`
- **Características:**
  - Interfaces definem contratos
  - Implementações específicas por tecnologia
  - SQLC para type-safe SQL queries

**Estrutura:**
```
repository/
├── interfaces.go           # Contratos dos repositórios
├── postgres/              # Implementação PostgreSQL
│   ├── unit_repo.go
│   ├── tenant_repo.go
│   ├── lease_repo.go
│   └── payment_repo.go
└── queries/               # SQLC queries
    ├── sqlc.yaml
    ├── schema.sql
    ├── units.sql
    ├── tenants.sql
    ├── leases.sql
    └── payments.sql
```

**Exemplo de Interface:**
```go
type UnitRepository interface {
    Create(ctx context.Context, unit *domain.Unit) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Unit, error)
    List(ctx context.Context, filters UnitFilters) ([]*domain.Unit, error)
    Update(ctx context.Context, unit *domain.Unit) error
    Delete(ctx context.Context, id uuid.UUID) error
    UpdateStatus(ctx context.Context, id uuid.UUID, status domain.UnitStatus) error
}
```

#### 3. Service Layer (Serviços/Casos de Uso)
- **Responsabilidade:** Orquestrar lógica de negócio complexa
- **Localização:** `internal/service/`
- **Características:**
  - Composição de repositórios
  - Validações de negócio
  - Transações multi-repository
  - Coordenação de operações

**Exemplo:**
```go
type LeaseService struct {
    leaseRepo   repository.LeaseRepository
    unitRepo    repository.UnitRepository
    paymentRepo repository.PaymentRepository
}

func (s *LeaseService) CreateLease(ctx context.Context, req CreateLeaseRequest) (*domain.Lease, error) {
    // 1. Validar unidade disponível
    // 2. Criar contrato
    // 3. Atualizar status da unidade para "occupied"
    // 4. Gerar pagamentos iniciais (aluguel + taxa pintura)
    // 5. Retornar contrato criado
}
```

#### 4. Handler Layer (Controladores HTTP)
- **Responsabilidade:** Lidar com requisições HTTP
- **Localização:** `internal/handler/`
- **Características:**
  - Parse de requests
  - Validação de inputs
  - Chamada aos services
  - Formatação de responses
  - Tratamento de erros HTTP

**Exemplo:**
```go
type UnitHandler struct {
    unitService *service.UnitService
    validator   *validator.Validate
}

func (h *UnitHandler) CreateUnit(w http.ResponseWriter, r *http.Request) {
    var req CreateUnitRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    if err := h.validator.Struct(req); err != nil {
        response.ValidationError(w, err)
        return
    }
    
    unit, err := h.unitService.CreateUnit(r.Context(), req)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    response.JSON(w, http.StatusCreated, unit)
}
```

## 🗂 Estrutura de Diretórios Detalhada

```
kitnet-manager/
│
├── cmd/
│   └── api/
│       └── main.go                    # Entry point, setup inicial
│
├── internal/                          # Código privado da aplicação
│   │
│   ├── domain/                        # Camada de Domínio
│   │   ├── unit.go                   # Entidade Unit + métodos
│   │   ├── tenant.go                 # Entidade Tenant + métodos
│   │   ├── lease.go                  # Entidade Lease + métodos
│   │   ├── payment.go                # Entidade Payment + métodos
│   │   ├── notification.go           # Entidade Notification + métodos
│   │   └── errors.go                 # Erros customizados do domínio
│   │
│   ├── repository/                    # Camada de Dados
│   │   ├── interfaces.go             # Contratos (interfaces)
│   │   ├── postgres/                 # Implementação PostgreSQL
│   │   │   ├── db.go                # Setup de conexão
│   │   │   ├── unit_repo.go
│   │   │   ├── tenant_repo.go
│   │   │   ├── lease_repo.go
│   │   │   ├── payment_repo.go
│   │   │   └── notification_repo.go
│   │   └── queries/                  # SQLC queries SQL
│   │       ├── sqlc.yaml            # Configuração SQLC
│   │       ├── schema.sql           # Schema completo (referência)
│   │       ├── units.sql
│   │       ├── tenants.sql
│   │       ├── leases.sql
│   │       ├── payments.sql
│   │       └── notifications.sql
│   │
│   ├── service/                       # Camada de Negócio
│   │   ├── unit_service.go
│   │   ├── tenant_service.go
│   │   ├── lease_service.go          # Lógica complexa de contratos
│   │   ├── payment_service.go        # Geração de pagamentos
│   │   ├── notification_service.go   # Lógica de notificações
│   │   └── dashboard_service.go      # Agregações para dashboard
│   │
│   ├── handler/                       # Camada HTTP
│   │   ├── unit_handler.go
│   │   ├── tenant_handler.go
│   │   ├── lease_handler.go
│   │   ├── payment_handler.go
│   │   ├── notification_handler.go
│   │   ├── dashboard_handler.go
│   │   └── router.go                 # Setup de rotas
│   │
│   └── pkg/                           # Utilitários internos
│       ├── database/
│       │   └── postgres.go           # Conexão com Neon
│       ├── validator/
│       │   └── validator.go          # Setup do go-playground/validator
│       ├── response/
│       │   └── response.go           # Padronização de respostas HTTP
│       └── middleware/
│           ├── logger.go             # Logging de requests
│           ├── cors.go               # CORS
│           └── recovery.go           # Panic recovery
│
├── migrations/                        # Database migrations
│   ├── 000001_create_units_table.up.sql
│   ├── 000001_create_units_table.down.sql
│   ├── 000002_create_tenants_table.up.sql
│   ├── 000002_create_tenants_table.down.sql
│   └── ...
│
├── config/
│   └── config.go                     # Gerenciamento de configurações
│
├── docs/
│   ├── api/                          # Documentação da API
│   └── database/                     # Diagramas ERD
│
├── .env.example                       # Template de variáveis de ambiente
├── .gitignore
├── go.mod
├── go.sum
├── Makefile                          # Comandos úteis
├── README.md
├── ARCHITECTURE.md
└── ROADMAP.md
```

## 🗄 Modelo de Dados (Database Schema)

### Diagrama ER (Simplificado)

```
┌──────────────┐         ┌──────────────┐
│    UNITS     │         │   TENANTS    │
├──────────────┤         ├──────────────┤
│ id (PK)      │         │ id (PK)      │
│ number       │         │ full_name    │
│ floor        │         │ cpf (UQ)     │
│ status       │         │ phone        │
│ is_renovated │         │ email        │
│ *_rent_value │         │ ...          │
└──────┬───────┘         └──────┬───────┘
       │                        │
       │    ┌──────────────┐    │
       └───▶│   LEASES     │◀───┘
            ├──────────────┤
            │ id (PK)      │
            │ unit_id (FK) │
            │ tenant_id(FK)│
            │ start_date   │
            │ end_date     │
            │ status       │
            │ ...          │
            └──────┬───────┘
                   │
                   │
            ┌──────▼───────┐
            │  PAYMENTS    │
            ├──────────────┤
            │ id (PK)      │
            │ lease_id (FK)│
            │ type         │
            │ due_date     │
            │ status       │
            │ amount       │
            └──────────────┘
```

### Tabelas Principais

#### units
```sql
CREATE TABLE units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number VARCHAR(10) NOT NULL UNIQUE,
    floor INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL,
    is_renovated BOOLEAN DEFAULT FALSE,
    base_rent_value DECIMAL(10,2) NOT NULL,
    renovated_rent_value DECIMAL(10,2) NOT NULL,
    current_rent_value DECIMAL(10,2) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

#### tenants
```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255),
    id_document_type VARCHAR(10),
    id_document_number VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);
```

#### leases
```sql
CREATE TABLE leases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    unit_id UUID NOT NULL REFERENCES units(id),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    contract_signed_date DATE NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    payment_due_day INTEGER NOT NULL CHECK (payment_due_day BETWEEN 1 AND 31),
    monthly_rent_value DECIMAL(10,2) NOT NULL,
    painting_fee_total DECIMAL(10,2) NOT NULL,
    painting_fee_installments INTEGER NOT NULL CHECK (painting_fee_installments IN (1,2,3)),
    painting_fee_paid DECIMAL(10,2) DEFAULT 0,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

#### payments
```sql
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lease_id UUID NOT NULL REFERENCES leases(id),
    payment_type VARCHAR(20) NOT NULL,
    reference_month VARCHAR(7),
    installment_number INTEGER,
    due_date DATE NOT NULL,
    payment_date DATE,
    amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    payment_method VARCHAR(20),
    pix_reference VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

#### notifications
```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lease_id UUID REFERENCES leases(id),
    tenant_id UUID REFERENCES tenants(id),
    type VARCHAR(50) NOT NULL,
    scheduled_date DATE NOT NULL,
    sent_date TIMESTAMP,
    status VARCHAR(20) NOT NULL,
    message_content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

## 🔄 Fluxos Principais

### 1. Criação de Novo Contrato

```
1. Handler recebe request com unit_id, tenant_id, start_date, etc
2. Service valida:
   - Unidade existe e está disponível
   - Morador existe
   - Datas são válidas
3. Service cria Lease no repository
4. Service atualiza Unit.status para "occupied"
5. Service gera Payment inicial (aluguel mês 1)
6. Service gera Payments da taxa de pintura (1x ou 3x)
7. Handler retorna Lease criado + lista de payments gerados
```

### 2. Registro de Pagamento

```
1. Handler recebe payment_id + payment_date + payment_method
2. Service busca Payment
3. Service valida status (deve estar "pending" ou "overdue")
4. Service atualiza Payment:
   - payment_date = data informada
   - status = "paid"
   - payment_method = método informado
5. Se for painting_fee, atualiza Lease.painting_fee_paid
6. Handler retorna Payment atualizado
```

### 3. Verificação Diária de Notificações

```
1. Cronjob executa NotificationService.ProcessDailyNotifications()
2. Service busca Payments pendentes com due_date = hoje + 3 dias
3. Para cada Payment, cria Notification tipo "sms_rent_reminder"
4. Service busca Leases com end_date = hoje + 45 dias
5. Para cada Lease, cria Notification tipo "sms_contract_expiring"
6. Service marca Notifications como "pending"
7. (Futuro) Service dispara SMS via gateway externo
```

## 🔧 Decisões Técnicas

### Por que SQLC ao invés de ORM?

- ✅ Type-safe em tempo de compilação
- ✅ SQL puro, sem abstrações "mágicas"
- ✅ Performance otimizada
- ✅ Facilita queries complexas
- ✅ Baixa curva de aprendizado
- ❌ Contras: Mais verboso que ORMs tradicionais

### Por que Chi Router?

- ✅ Leve e idiomático para Go
- ✅ Compatível com net/http
- ✅ Middlewares flexíveis
- ✅ Suporte a subrouters
- ✅ Grande adoção na comunidade

### Por que Neon PostgreSQL?

- ✅ Free tier generoso (0.5 GB storage)
- ✅ Serverless (pay-per-use no futuro)
- ✅ Setup rápido
- ✅ Interface web amigável
- ✅ PostgreSQL completo (sem limitações)

## 🧪 Estratégia de Testes

### Níveis de Teste

1. **Unit Tests:** Services e funções de domínio
2. **Integration Tests:** Repositories com banco real (test container)
3. **E2E Tests:** Handlers completos (futuro)

### Estrutura de Testes

```
internal/
├── service/
│   ├── lease_service.go
│   └── lease_service_test.go
└── repository/
    └── postgres/
        ├── lease_repo.go
        └── lease_repo_test.go
```

## 🚀 Deploy (Futuro)

### Opções Consideradas

1. **Railway:** Simples, free tier, suporta PostgreSQL
2. **Render:** Free tier, auto-deploy via GitHub
3. **Fly.io:** Edge computing, free tier limitado

### Requisitos Mínimos

- CPU: 0.5 vCPU
- RAM: 512 MB
- Storage: Neon cobre separadamente

## 📊 Performance Considerations

### Indexação

```sql
-- Indexes críticos
CREATE INDEX idx_leases_unit_id ON leases(unit_id);
CREATE INDEX idx_leases_tenant_id ON leases(tenant_id);
CREATE INDEX idx_leases_status ON leases(status);
CREATE INDEX idx_payments_lease_id ON payments(lease_id);
CREATE INDEX idx_payments_status_due_date ON payments(status, due_date);
```

### Caching (Futuro)

- Redis para dashboard aggregations
- Cache de configurações em memória

## 🔐 Segurança

### Autenticação (MVP: Não implementado)

- **v1.0:** Sistema interno, sem autenticação
- **v2.0:** JWT + Login básico
- **v3.0:** OAuth2 para portal do morador

### Validação

- Input validation em todos os handlers
- SQL injection protegido (SQLC)
- CORS configurado

## 📈 Escalabilidade

### Estratégias Futuras

1. **Horizontal scaling:** Múltiplas instâncias da API
2. **Read replicas:** Separar leitura de escrita (Neon suporta)
3. **Queue system:** Notificações async (RabbitMQ/Redis)
4. **CDN:** Frontend estático

---

**Última atualização:** Setembro 2025
