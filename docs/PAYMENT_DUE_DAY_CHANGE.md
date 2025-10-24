# Plano de Solução: Mudança de Data de Vencimento de Pagamentos

## 📋 Contexto do Problema

### Situação Atual
- Ao criar um contrato, define-se um `payment_due_day` (dia do mês para vencimento, ex: dia 15)
- São gerados automaticamente 6 pagamentos de aluguel + N parcelas de taxa de pintura
- Todos os pagamentos vencem no mesmo dia do mês definido em `payment_due_day`
- **LIMITAÇÃO:** O `payment_due_day` é fixo e não pode ser alterado durante o contrato

### Necessidade de Negócio
Pode acontecer de um inquilino precisar alterar a data de vencimento durante o contrato.

**Exemplo:**
1. Contrato iniciado em 15/10/2025 com vencimento no dia 15
2. Pagamentos gerados: 15/10, 15/11, 15/12, 15/01, 15/02, 15/03
3. Em novembro, inquilino quer mudar o vencimento para dia 5
4. No dia 05/12, ele faz o pagamento proporcional (período de 15/11 até 05/12 = 20 dias)
5. A partir daí, todos os vencimentos subsequentes passam a ser no dia 5

---

## 🎯 Objetivos da Solução

1. **Permitir mudança do dia de vencimento** durante a vigência do contrato
2. **Calcular pagamento proporcional** referente ao período entre a data antiga e a nova
3. **Recalcular datas de vencimento** de todos os pagamentos futuros pendentes
4. **Manter integridade** dos pagamentos já realizados (não alterar)
5. **Registrar histórico** da mudança para auditoria

---

## 🔍 Análise do Sistema Atual

### Estrutura de Dados

**Tabela: `leases`**
```sql
CREATE TABLE leases (
    id UUID PRIMARY KEY,
    unit_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    payment_due_day INTEGER NOT NULL,  -- ⚠️ CAMPO CRÍTICO
    monthly_rent_value DECIMAL(10,2),
    status VARCHAR(20),
    ...
);
```

**Tabela: `payments`**
```sql
CREATE TABLE payments (
    id UUID PRIMARY KEY,
    lease_id UUID NOT NULL REFERENCES leases(id),
    payment_type VARCHAR(20),  -- 'rent' ou 'painting_fee'
    reference_month DATE NOT NULL,
    amount DECIMAL(10,2),
    status VARCHAR(20),  -- 'pending', 'overdue', 'paid', 'cancelled'
    due_date DATE NOT NULL,  -- ⚠️ CAMPO CRÍTICO
    payment_date DATE,
    payment_method VARCHAR(20),
    ...
);
```

### Fluxo de Geração de Pagamentos (Atual)

```
CreateLease()
├── Valida dados do contrato
├── Cria Lease com payment_due_day = 15
├── Gera 6 pagamentos de aluguel:
│   ├── Mês 1: due_date = 15/10/2025
│   ├── Mês 2: due_date = 15/11/2025
│   ├── Mês 3: due_date = 15/12/2025
│   ├── Mês 4: due_date = 15/01/2026
│   ├── Mês 5: due_date = 15/02/2026
│   └── Mês 6: due_date = 15/03/2026
└── Gera N pagamentos de taxa de pintura (mesmo dia)
```

**Código Atual (payment_service.go:58-64):**
```go
dueDate := time.Date(
    req.ReferenceMonth.Year(),
    req.ReferenceMonth.Month(),
    lease.PaymentDueDay,  // ← Usa o payment_due_day do contrato
    0, 0, 0, 0,
    time.UTC,
)
```

---

## 💡 Solução Proposta

### 1. Nova Funcionalidade: Alterar Data de Vencimento

**Endpoint:** `POST /leases/{lease_id}/change-payment-due-day`

**Request:**
```json
{
  "new_payment_due_day": 5,
  "effective_date": "2025-12-05",
  "reason": "Solicitação do inquilino para ajuste de fluxo financeiro"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Payment due day changed successfully",
  "data": {
    "lease_id": "uuid-123",
    "old_payment_due_day": 15,
    "new_payment_due_day": 5,
    "effective_date": "2025-12-05",
    "proportional_payment": {
      "id": "payment-uuid-prop",
      "reference_period": "15/11/2025 - 05/12/2025",
      "days": 20,
      "amount": "800.00",
      "due_date": "2025-12-05",
      "status": "pending"
    },
    "updated_payments_count": 5,
    "updated_payments": [
      {
        "id": "payment-uuid-1",
        "reference_month": "2025-12-01",
        "old_due_date": "2025-12-15",
        "new_due_date": "2025-12-05"
      },
      {
        "id": "payment-uuid-2",
        "reference_month": "2026-01-01",
        "old_due_date": "2026-01-15",
        "new_due_date": "2026-01-05"
      },
      // ... demais pagamentos atualizados
    ]
  }
}
```

---

## 🏗️ Implementação Detalhada

### Passo 1: Criar Estruturas de Dados

**Arquivo:** `internal/service/lease_service.go`

```go
// ChangePaymentDueDayRequest representa a requisição para alterar dia de vencimento
type ChangePaymentDueDayRequest struct {
    LeaseID           uuid.UUID `json:"lease_id" validate:"required"`
    NewPaymentDueDay  int       `json:"new_payment_due_day" validate:"required,min=1,max=31"`
    EffectiveDate     time.Time `json:"effective_date" validate:"required"`
    Reason            string    `json:"reason"`
}

// ChangePaymentDueDayResponse representa a resposta da mudança
type ChangePaymentDueDayResponse struct {
    LeaseID                uuid.UUID                      `json:"lease_id"`
    OldPaymentDueDay       int                            `json:"old_payment_due_day"`
    NewPaymentDueDay       int                            `json:"new_payment_due_day"`
    EffectiveDate          time.Time                      `json:"effective_date"`
    ProportionalPayment    *ProportionalPaymentInfo       `json:"proportional_payment,omitempty"`
    UpdatedPaymentsCount   int                            `json:"updated_payments_count"`
    UpdatedPayments        []UpdatedPaymentInfo           `json:"updated_payments"`
}

// ProportionalPaymentInfo contém informações do pagamento proporcional gerado
type ProportionalPaymentInfo struct {
    ID              uuid.UUID       `json:"id"`
    ReferencePeriod string          `json:"reference_period"`
    Days            int             `json:"days"`
    Amount          decimal.Decimal `json:"amount"`
    DueDate         time.Time       `json:"due_date"`
    Status          string          `json:"status"`
}

// UpdatedPaymentInfo contém informações sobre pagamentos que tiveram data alterada
type UpdatedPaymentInfo struct {
    ID             uuid.UUID `json:"id"`
    ReferenceMonth time.Time `json:"reference_month"`
    OldDueDate     time.Time `json:"old_due_date"`
    NewDueDate     time.Time `json:"new_due_date"`
}
```

---

### Passo 2: Lógica de Negócio Principal

**Arquivo:** `internal/service/lease_service.go`

```go
// ChangePaymentDueDay altera o dia de vencimento de um contrato e recalcula pagamentos futuros
func (s *LeaseService) ChangePaymentDueDay(ctx context.Context, req ChangePaymentDueDayRequest) (*ChangePaymentDueDayResponse, error) {
    // ==================================================
    // ETAPA 1: VALIDAÇÕES
    // ==================================================

    // 1.1. Buscar o contrato
    lease, err := s.leaseRepo.GetByID(ctx, req.LeaseID)
    if err != nil {
        return nil, fmt.Errorf("error getting lease: %w", err)
    }
    if lease == nil {
        return nil, ErrLeaseNotFound
    }

    // 1.2. Validar que contrato está ativo
    if lease.Status != domain.LeaseStatusActive && lease.Status != domain.LeaseStatusExpiringSoon {
        return nil, errors.New("lease must be active to change payment due day")
    }

    // 1.3. Validar que o novo dia é diferente do atual
    if req.NewPaymentDueDay == lease.PaymentDueDay {
        return nil, errors.New("new payment due day must be different from current")
    }

    // 1.4. Validar que o novo dia está no range válido (1-31)
    if req.NewPaymentDueDay < 1 || req.NewPaymentDueDay > 31 {
        return nil, errors.New("payment due day must be between 1 and 31")
    }

    // 1.5. Validar que a data efetiva não está no passado
    if req.EffectiveDate.Before(time.Now()) {
        return nil, errors.New("effective date cannot be in the past")
    }

    // 1.6. Validar que a data efetiva está dentro da vigência do contrato
    if req.EffectiveDate.Before(lease.StartDate) || req.EffectiveDate.After(lease.EndDate) {
        return nil, errors.New("effective date must be within lease period")
    }

    // ==================================================
    // ETAPA 2: CALCULAR PAGAMENTO PROPORCIONAL
    // ==================================================

    // 2.1. Determinar o período proporcional
    // Exemplo: dia 15 → dia 5
    //   - Último vencimento no dia antigo: 15/11
    //   - Próximo vencimento no dia novo: 05/12
    //   - Período proporcional: 15/11 até 05/12 = 20 dias

    oldDueDay := lease.PaymentDueDay
    newDueDay := req.NewPaymentDueDay

    // Determinar a data do último vencimento no dia antigo
    // Se hoje é 20/11 e o antigo vencimento era dia 15, o último vencimento foi 15/11
    lastOldDueDate := time.Date(
        req.EffectiveDate.Year(),
        req.EffectiveDate.Month(),
        oldDueDay,
        0, 0, 0, 0,
        time.UTC,
    )

    // Se a data efetiva é antes do dia antigo no mês atual,
    // o último vencimento foi no mês anterior
    if req.EffectiveDate.Day() < oldDueDay {
        lastOldDueDate = lastOldDueDate.AddDate(0, -1, 0)
    }

    // A nova data de vencimento (data efetiva)
    firstNewDueDate := req.EffectiveDate

    // Calcular quantos dias entre o último vencimento antigo e o primeiro novo
    proportionalDays := int(firstNewDueDate.Sub(lastOldDueDate).Hours() / 24)

    // 2.2. Calcular valor proporcional
    // Valor proporcional = (valor_mensal / 30) * dias_proporcionais
    dailyRate := lease.MonthlyRentValue.Div(decimal.NewFromInt(30))
    proportionalAmount := dailyRate.Mul(decimal.NewFromInt(int64(proportionalDays)))

    // 2.3. Criar pagamento proporcional
    var proportionalPayment *domain.Payment
    if proportionalDays > 0 && proportionalAmount.GreaterThan(decimal.Zero) {
        // Usar o mês de referência da data efetiva
        referenceMonth := time.Date(
            firstNewDueDate.Year(),
            firstNewDueDate.Month(),
            1, 0, 0, 0, 0,
            time.UTC,
        )

        proportionalPayment, err = domain.NewPayment(
            lease.ID,
            domain.PaymentTypeAdjustment,  // ← Novo tipo: 'adjustment'
            referenceMonth,
            proportionalAmount,
            firstNewDueDate,
        )
        if err != nil {
            return nil, fmt.Errorf("error creating proportional payment: %w", err)
        }

        // Adicionar nota explicativa
        note := fmt.Sprintf(
            "Pagamento proporcional devido à mudança de vencimento do dia %d para dia %d. Período: %s a %s (%d dias)",
            oldDueDay,
            newDueDay,
            lastOldDueDate.Format("02/01/2006"),
            firstNewDueDate.Format("02/01/2006"),
            proportionalDays,
        )
        proportionalPayment.Notes = &note

        // Salvar no banco
        if err := s.paymentRepo.Create(ctx, proportionalPayment); err != nil {
            return nil, fmt.Errorf("error saving proportional payment: %w", err)
        }
    }

    // ==================================================
    // ETAPA 3: RECALCULAR PAGAMENTOS FUTUROS
    // ==================================================

    // 3.1. Buscar todos os pagamentos pendentes/atrasados do contrato
    allPayments, err := s.paymentRepo.GetByLeaseID(ctx, lease.ID)
    if err != nil {
        return nil, fmt.Errorf("error getting lease payments: %w", err)
    }

    // 3.2. Filtrar apenas pagamentos futuros que ainda não foram pagos
    var paymentsToUpdate []*domain.Payment
    for _, payment := range allPayments {
        // Só atualiza se:
        // - Status é pending ou overdue (não pago)
        // - A data de vencimento é após a data efetiva
        if (payment.Status == domain.PaymentStatusPending || payment.Status == domain.PaymentStatusOverdue) &&
           payment.DueDate.After(req.EffectiveDate) {
            paymentsToUpdate = append(paymentsToUpdate, payment)
        }
    }

    // 3.3. Atualizar a due_date de cada pagamento futuro
    updatedPaymentsInfo := make([]UpdatedPaymentInfo, 0, len(paymentsToUpdate))

    for _, payment := range paymentsToUpdate {
        oldDueDate := payment.DueDate

        // Calcular nova due_date mantendo o ano/mês, mas mudando o dia
        newDueDate := time.Date(
            payment.ReferenceMonth.Year(),
            payment.ReferenceMonth.Month(),
            req.NewPaymentDueDay,
            0, 0, 0, 0,
            time.UTC,
        )

        // Atualizar o pagamento
        payment.DueDate = newDueDate
        payment.UpdatedAt = time.Now()

        // Salvar no banco
        if err := s.paymentRepo.Update(ctx, payment); err != nil {
            return nil, fmt.Errorf("error updating payment %s: %w", payment.ID, err)
        }

        // Registrar a mudança
        updatedPaymentsInfo = append(updatedPaymentsInfo, UpdatedPaymentInfo{
            ID:             payment.ID,
            ReferenceMonth: payment.ReferenceMonth,
            OldDueDate:     oldDueDate,
            NewDueDate:     newDueDate,
        })
    }

    // ==================================================
    // ETAPA 4: ATUALIZAR O CONTRATO
    // ==================================================

    // 4.1. Atualizar payment_due_day no contrato
    oldPaymentDueDay := lease.PaymentDueDay
    lease.PaymentDueDay = req.NewPaymentDueDay
    lease.UpdatedAt = time.Now()

    if err := s.leaseRepo.Update(ctx, lease); err != nil {
        return nil, fmt.Errorf("error updating lease: %w", err)
    }

    // ==================================================
    // ETAPA 5: REGISTRAR HISTÓRICO (OPCIONAL)
    // ==================================================

    // TODO: Criar tabela de audit_log para registrar essa mudança
    // auditLog := &AuditLog{
    //     EntityType: "lease",
    //     EntityID:   lease.ID,
    //     Action:     "change_payment_due_day",
    //     OldValue:   oldPaymentDueDay,
    //     NewValue:   req.NewPaymentDueDay,
    //     Reason:     req.Reason,
    //     PerformedAt: time.Now(),
    // }
    // s.auditRepo.Create(ctx, auditLog)

    // ==================================================
    // ETAPA 6: MONTAR RESPOSTA
    // ==================================================

    response := &ChangePaymentDueDayResponse{
        LeaseID:              lease.ID,
        OldPaymentDueDay:     oldPaymentDueDay,
        NewPaymentDueDay:     req.NewPaymentDueDay,
        EffectiveDate:        req.EffectiveDate,
        UpdatedPaymentsCount: len(updatedPaymentsInfo),
        UpdatedPayments:      updatedPaymentsInfo,
    }

    // Incluir informações do pagamento proporcional se foi criado
    if proportionalPayment != nil {
        response.ProportionalPayment = &ProportionalPaymentInfo{
            ID:              proportionalPayment.ID,
            ReferencePeriod: fmt.Sprintf("%s - %s", lastOldDueDate.Format("02/01/2006"), firstNewDueDate.Format("02/01/2006")),
            Days:            proportionalDays,
            Amount:          proportionalAmount,
            DueDate:         firstNewDueDate,
            Status:          string(proportionalPayment.Status),
        }
    }

    return response, nil
}
```

---

### Passo 3: Novo Tipo de Pagamento (Ajuste)

**Arquivo:** `internal/domain/payment.go`

Adicionar novo tipo de pagamento para o valor proporcional:

```go
const (
    PaymentTypeRent        PaymentType = "rent"
    PaymentTypePaintingFee PaymentType = "painting_fee"
    PaymentTypeAdjustment  PaymentType = "adjustment"  // ← NOVO
)
```

---

### Passo 4: Handler HTTP

**Arquivo:** `internal/handler/lease_handler.go`

```go
// ChangePaymentDueDayRequest representa o request HTTP
type ChangePaymentDueDayHTTPRequest struct {
    NewPaymentDueDay int    `json:"new_payment_due_day" validate:"required,min=1,max=31"`
    EffectiveDate    string `json:"effective_date" validate:"required"`
    Reason           string `json:"reason"`
}

// ChangePaymentDueDay godoc
// @Summary      Alterar dia de vencimento de pagamentos
// @Description  Altera o dia de vencimento de um contrato e recalcula pagamentos futuros
// @Tags         leases
// @Accept       json
// @Produce      json
// @Param        id   path      string                           true  "Lease ID"
// @Param        request body   ChangePaymentDueDayHTTPRequest  true  "Dados da mudança"
// @Success      200  {object}  Response{data=service.ChangePaymentDueDayResponse}
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /leases/{id}/change-payment-due-day [post]
func (h *LeaseHandler) ChangePaymentDueDay(c *gin.Context) {
    // 1. Extrair lease_id da URL
    leaseIDStr := c.Param("id")
    leaseID, err := uuid.Parse(leaseIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Success: false,
            Error:   "invalid lease ID format",
        })
        return
    }

    // 2. Parse do body
    var httpReq ChangePaymentDueDayHTTPRequest
    if err := c.ShouldBindJSON(&httpReq); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Success: false,
            Error:   fmt.Sprintf("invalid request body: %v", err),
        })
        return
    }

    // 3. Validar campos
    if err := h.validator.Struct(httpReq); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Success: false,
            Error:   fmt.Sprintf("validation error: %v", err),
        })
        return
    }

    // 4. Parse da effective_date
    effectiveDate, err := time.Parse("2006-01-02", httpReq.EffectiveDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Success: false,
            Error:   "invalid effective_date format, expected YYYY-MM-DD",
        })
        return
    }

    // 5. Montar request do service
    serviceReq := service.ChangePaymentDueDayRequest{
        LeaseID:          leaseID,
        NewPaymentDueDay: httpReq.NewPaymentDueDay,
        EffectiveDate:    effectiveDate,
        Reason:           httpReq.Reason,
    }

    // 6. Chamar o service
    response, err := h.leaseService.ChangePaymentDueDay(c.Request.Context(), serviceReq)
    if err != nil {
        if errors.Is(err, service.ErrLeaseNotFound) {
            c.JSON(http.StatusNotFound, ErrorResponse{
                Success: false,
                Error:   "lease not found",
            })
            return
        }

        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Success: false,
            Error:   fmt.Sprintf("error changing payment due day: %v", err),
        })
        return
    }

    // 7. Retornar sucesso
    c.JSON(http.StatusOK, Response{
        Success: true,
        Message: "Payment due day changed successfully",
        Data:    response,
    })
}
```

---

### Passo 5: Adicionar Route

**Arquivo:** `internal/routes/routes.go`

```go
// Adicionar na seção de leases
leasesGroup.POST("/:id/change-payment-due-day", leaseHandler.ChangePaymentDueDay)
```

---

### Passo 6: Métodos de Repository

**Arquivo:** `internal/repository/postgres/payment_repo.go`

Adicionar método para atualizar pagamento (se não existir):

```go
// Update atualiza um pagamento existente
func (r *PaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
    query := `
        UPDATE payments
        SET
            due_date = $2,
            amount = $3,
            status = $4,
            payment_date = $5,
            payment_method = $6,
            notes = $7,
            updated_at = $8
        WHERE id = $1
    `

    _, err := r.db.ExecContext(
        ctx,
        query,
        payment.ID,
        payment.DueDate,
        payment.Amount,
        payment.Status,
        payment.PaymentDate,
        payment.PaymentMethod,
        payment.Notes,
        time.Now(),
    )

    if err != nil {
        return fmt.Errorf("error updating payment: %w", err)
    }

    return nil
}
```

---

## 📊 Exemplo de Uso Completo

### Cenário Real

**Contrato Inicial:**
```json
{
  "id": "lease-123",
  "start_date": "2025-10-15",
  "end_date": "2026-04-15",
  "payment_due_day": 15,
  "monthly_rent_value": "1200.00"
}
```

**Pagamentos Gerados:**
```json
[
  { "id": "pay-1", "reference_month": "2025-10-01", "due_date": "2025-10-15", "amount": "1200.00", "status": "paid" },
  { "id": "pay-2", "reference_month": "2025-11-01", "due_date": "2025-11-15", "amount": "1200.00", "status": "paid" },
  { "id": "pay-3", "reference_month": "2025-12-01", "due_date": "2025-12-15", "amount": "1200.00", "status": "pending" },
  { "id": "pay-4", "reference_month": "2026-01-01", "due_date": "2026-01-15", "amount": "1200.00", "status": "pending" },
  { "id": "pay-5", "reference_month": "2026-02-01", "due_date": "2026-02-15", "amount": "1200.00", "status": "pending" },
  { "id": "pay-6", "reference_month": "2026-03-01", "due_date": "2026-03-15", "amount": "1200.00", "status": "pending" }
]
```

**Solicitação de Mudança (em 20/11/2025):**
```bash
POST /leases/lease-123/change-payment-due-day
{
  "new_payment_due_day": 5,
  "effective_date": "2025-12-05",
  "reason": "Ajuste de fluxo de caixa do inquilino"
}
```

**Resultado:**

1. **Pagamento Proporcional Criado:**
```json
{
  "id": "pay-prop-1",
  "payment_type": "adjustment",
  "reference_month": "2025-12-01",
  "amount": "800.00",  // (1200 / 30) * 20 dias
  "due_date": "2025-12-05",
  "status": "pending",
  "notes": "Pagamento proporcional devido à mudança de vencimento do dia 15 para dia 5. Período: 15/11/2025 a 05/12/2025 (20 dias)"
}
```

2. **Pagamentos Futuros Atualizados:**
```json
[
  { "id": "pay-3", "due_date": "2025-12-05" },  // era 2025-12-15
  { "id": "pay-4", "due_date": "2026-01-05" },  // era 2026-01-15
  { "id": "pay-5", "due_date": "2026-02-05" },  // era 2026-02-15
  { "id": "pay-6", "due_date": "2026-03-05" }   // era 2026-03-15
]
```

3. **Contrato Atualizado:**
```json
{
  "id": "lease-123",
  "payment_due_day": 5  // era 15
}
```

---

## ✅ Validações e Regras de Negócio

### Validações Obrigatórias

1. **Contrato deve existir**
2. **Contrato deve estar ativo** (`active` ou `expiring_soon`)
3. **Novo dia deve ser diferente do atual**
4. **Novo dia deve estar entre 1 e 31**
5. **Data efetiva não pode estar no passado**
6. **Data efetiva deve estar dentro da vigência do contrato**
7. **Não deve haver outro processo de mudança em andamento** (se implementar controle de concorrência)

### Regras de Negócio

1. **Pagamentos já pagos não são alterados** (apenas pendentes/atrasados)
2. **Pagamento proporcional só é criado se houver dias a cobrar**
3. **Pagamento proporcional tem tipo "adjustment"**
4. **A mudança afeta todos os pagamentos futuros** (aluguel e taxa de pintura)
5. **O cálculo proporcional usa base de 30 dias** por mês
6. **A data efetiva define quando o novo vencimento passa a valer**

---

## 🧪 Testes Necessários

### Testes Unitários

1. **Cálculo de dias proporcionais:**
   - Mudança dentro do mesmo mês
   - Mudança para mês seguinte
   - Mudança com dia 31 (edge case)

2. **Cálculo de valor proporcional:**
   - Valores decimais corretos
   - Arredondamento adequado

3. **Recálculo de datas:**
   - Apenas pagamentos pendentes são alterados
   - Pagamentos pagos permanecem intactos
   - Mês de fevereiro (28/29 dias)

### Testes de Integração

1. **Fluxo completo de mudança**
2. **Concorrência:** múltiplas requisições simultâneas
3. **Rollback em caso de erro**

### Testes de API

```bash
# Teste 1: Mudança bem-sucedida
POST /leases/{id}/change-payment-due-day
Expect: 200 OK com pagamento proporcional e lista de atualizações

# Teste 2: Contrato não encontrado
POST /leases/invalid-uuid/change-payment-due-day
Expect: 404 Not Found

# Teste 3: Mesmo dia de vencimento
POST /leases/{id}/change-payment-due-day
{ "new_payment_due_day": 15 }  # sendo que já é 15
Expect: 400 Bad Request

# Teste 4: Data efetiva no passado
POST /leases/{id}/change-payment-due-day
{ "effective_date": "2020-01-01" }
Expect: 400 Bad Request

# Teste 5: Contrato cancelado
POST /leases/{cancelled-lease-id}/change-payment-due-day
Expect: 400 Bad Request
```

---

## 📝 Migrações de Banco de Dados

### Nova Migração (Opcional): Adicionar Tipo "adjustment"

**Arquivo:** `migrations/000XXX_add_payment_type_adjustment.up.sql`

```sql
-- Se a coluna payment_type usa ENUM, adicionar o novo valor
ALTER TYPE payment_type_enum ADD VALUE IF NOT EXISTS 'adjustment';

-- Ou se for VARCHAR, não precisa migration
-- Apenas documentar que o novo tipo é válido
```

### Nova Tabela (Opcional): Audit Log

**Arquivo:** `migrations/000XXX_create_audit_log.up.sql`

```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(50) NOT NULL,  -- 'lease', 'payment', etc
    entity_id UUID NOT NULL,
    action VARCHAR(100) NOT NULL,      -- 'change_payment_due_day', etc
    old_value TEXT,
    new_value TEXT,
    reason TEXT,
    performed_by UUID,                 -- user_id (futuro)
    performed_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_performed_at ON audit_logs(performed_at DESC);
```

---

## 🚀 Plano de Implementação (Passos)

### Fase 1: Backend Core (Prioritário)

1. ✅ **Criar estruturas de dados** (requests/responses)
2. ✅ **Implementar lógica no LeaseService** (ChangePaymentDueDay)
3. ✅ **Adicionar novo tipo de pagamento** (adjustment)
4. ✅ **Implementar método Update no PaymentRepository**
5. ✅ **Criar handler HTTP**
6. ✅ **Registrar rota**
7. ✅ **Escrever testes unitários**
8. ✅ **Testar manualmente com Postman/curl**

### Fase 2: Melhorias (Desejável)

9. ⚠️ **Adicionar auditoria** (tabela de logs)
10. ⚠️ **Adicionar validação de concorrência** (evitar mudanças simultâneas)
11. ⚠️ **Notificações** (email/SMS para inquilino sobre mudança)
12. ⚠️ **Webhook/evento** para sistemas externos

### Fase 3: Frontend (Subsequente)

13. 🎨 **Criar interface de mudança de vencimento**
14. 🎨 **Exibir histórico de mudanças**
15. 🎨 **Confirmação visual** antes de aplicar

---

## 🎯 Casos de Uso e Edge Cases

### Caso 1: Mudança Simples
- **Cenário:** Dia 15 → Dia 5, no meio do contrato
- **Comportamento:** Gera proporcional + atualiza futuros
- **Status:** ✅ Coberto pela solução

### Caso 2: Mudança no Último Mês
- **Cenário:** Contrato termina em março, mudança em fevereiro
- **Comportamento:** Atualiza apenas último pagamento pendente
- **Status:** ✅ Coberto pela solução

### Caso 3: Todos os Pagamentos Já Foram Pagos
- **Cenário:** Inquilino pagou tudo adiantado
- **Comportamento:** Apenas atualiza lease.payment_due_day, sem gerar proporcional
- **Status:** ✅ Coberto (não há pagamentos pendentes para atualizar)

### Caso 4: Dia 31 em Mês com 30 Dias
- **Cenário:** new_payment_due_day = 31, mas mês tem 30 dias
- **Comportamento:** Go ajusta automaticamente para dia 30
- **Status:** ⚠️ Documentar comportamento

### Caso 5: Múltiplas Mudanças no Mesmo Contrato
- **Cenário:** Inquilino muda de 15→5, depois de 5→10
- **Comportamento:** Cada mudança gera novo proporcional e recalcula futuros
- **Status:** ✅ Coberto (cada chamada é independente)

### Caso 6: Mudança em Contratos com Taxa de Pintura
- **Cenário:** Contrato tem pagamentos de aluguel + taxa de pintura pendentes
- **Comportamento:** Ambos os tipos são recalculados
- **Status:** ✅ Coberto (lógica não diferencia tipo)

---

## 📚 Documentação Adicional Necessária

### 1. Atualizar Swagger
- Adicionar endpoint `/leases/{id}/change-payment-due-day`
- Documentar request/response schemas

### 2. Atualizar README do Backend
- Explicar nova funcionalidade
- Adicionar exemplos de uso

### 3. Atualizar Documentação de Frontend
- Atualizar arquivo `API.md` com novo endpoint
- Adicionar exemplo de integração

### 4. Manual do Usuário
- Explicar quando usar a mudança de vencimento
- Alertas sobre impacto financeiro
- Passo a passo visual

---

## ⚠️ Considerações Importantes

### Segurança
- **Quem pode fazer essa mudança?** Apenas admin/manager?
- **Limite de mudanças:** Quantas vezes pode mudar no mesmo contrato?
- **Auditoria:** Registrar quem fez a mudança e quando

### Performance
- A operação pode atualizar vários pagamentos (até 6+)
- Considerar usar **transação** para garantir atomicidade
- Em caso de erro, fazer rollback completo

### Financeiro
- **Validar valor proporcional:** Garantir que não há cobranças duplicadas
- **Relatórios:** Incluir pagamentos proporcionais nos relatórios financeiros
- **Exportação:** Pagamentos de ajuste devem aparecer em extratos

### UX
- **Confirmação:** Mostrar preview antes de aplicar
- **Histórico:** Permitir visualizar mudanças anteriores
- **Notificação:** Avisar inquilino sobre a mudança

---

## 📊 Impacto em Outras Partes do Sistema

### Dashboard
- **Total a receber:** Incluir pagamentos proporcionais
- **Alertas:** Pagamento proporcional com vencimento próximo

### Relatórios
- **Receita mensal:** Considerar ajustes proporcionais
- **Inadimplência:** Pagamentos proporcionais também podem atrasar

### Notificações
- **Lembrete de vencimento:** Incluir pagamento proporcional
- **Email de mudança:** Notificar inquilino sobre alteração

---

## 🔧 Manutenção e Monitoramento

### Logs
```go
log.Info(
    "Payment due day changed",
    "lease_id", lease.ID,
    "old_day", oldDay,
    "new_day", newDay,
    "proportional_amount", proportionalAmount,
    "updated_payments", len(updatedPayments),
)
```

### Métricas
- Quantidade de mudanças por mês
- Valor médio de pagamentos proporcionais
- Tempo médio de processamento

### Alertas
- Erros ao processar mudança
- Valores proporcionais muito altos (possível bug)
- Tentativas de mudança em contratos inválidos

---

## ✅ Checklist de Implementação

- [ ] Criar estruturas de request/response
- [ ] Implementar lógica no LeaseService
- [ ] Adicionar tipo "adjustment" em PaymentType
- [ ] Implementar PaymentRepository.Update()
- [ ] Criar handler HTTP
- [ ] Registrar rota
- [ ] Adicionar validações
- [ ] Escrever testes unitários
- [ ] Testar com Postman
- [ ] Atualizar Swagger
- [ ] Documentar no README
- [ ] (Opcional) Criar tabela de auditoria
- [ ] (Opcional) Adicionar notificações
- [ ] Integrar com frontend
- [ ] Testes de integração
- [ ] Deploy em staging
- [ ] Testes com usuários
- [ ] Deploy em produção

---

## 📞 Contato e Suporte

Para dúvidas ou sugestões sobre essa implementação:
- **Documentação completa:** `/docs/PAYMENT_DUE_DAY_CHANGE.md`
- **API Reference:** `/docs/API.md`
- **Swagger:** `https://kitnet-manager-production.up.railway.app/swagger/index.html`

---

**Última Atualização:** 2025-10-23
**Versão:** 1.0
**Status:** 📝 Plano de Implementação
