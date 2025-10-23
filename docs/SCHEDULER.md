# Scheduler de Tarefas Automáticas

## 📋 Visão Geral

O scheduler executa tarefas de manutenção periodicamente para manter o sistema atualizado:

- ✅ Marcar pagamentos vencidos como atrasados
- ✅ Atualizar status de contratos expirando em breve
- ✅ Logs detalhados de cada execução

## ⚙️ Configuração

### Variável de Ambiente

```env
SCHEDULER_INTERVAL_HOURS=24
```

**Valores Recomendados:**
- `24` (padrão) - Executa **1x ao dia** (ideal para produção)
- `12` - Executa **2x ao dia**
- `6` - Executa **4x ao dia** (para ambientes com alta demanda)
- `1` - Executa **a cada hora** (apenas para desenvolvimento/testes)

### Configuração no .env

```bash
# Copiar do exemplo
cp .env.example .env

# Editar
nano .env

# Adicionar/modificar:
SCHEDULER_INTERVAL_HOURS=24
```

## 🚀 Como Funciona

### Inicialização

Quando o servidor inicia:
```
⏰ Scheduler iniciado (intervalo: 24h)
🔄 Executando tarefas agendadas...
📅 Verificando pagamentos atrasados...
✅ 3 pagamento(s) marcado(s) como atrasado(s)
```

### Execução Periódica

O scheduler executa automaticamente:
- **Imediatamente** na inicialização (corrige status pendentes)
- **A cada X horas** configurado (mantém atualizado)

### Tarefas Executadas

#### 1. Marcar Pagamentos Atrasados
```sql
UPDATE payments
SET status = 'overdue', updated_at = NOW()
WHERE status = 'pending'
  AND due_date < CURRENT_DATE
```

#### 2. Atualizar Contratos Expirando
```sql
UPDATE leases
SET status = 'expiring_soon', updated_at = NOW()
WHERE status = 'active'
  AND end_date BETWEEN NOW() AND NOW() + INTERVAL '45 days'
```

## 📊 Logs

### Sucesso (com atualizações)
```log
2024-10-23 00:00:00 🔄 Executando tarefas agendadas...
2024-10-23 00:00:00 📅 Verificando pagamentos atrasados...
2024-10-23 00:00:01 ✅ 5 pagamento(s) marcado(s) como atrasado(s)
2024-10-23 00:00:01 📅 Verificando contratos próximos de expirar...
2024-10-23 00:00:01 ✅ 2 contrato(s) marcado(s) como expirando em breve
2024-10-23 00:00:01 ✅ Tarefas agendadas concluídas
```

### Sucesso (sem atualizações)
```log
2024-10-23 00:00:00 🔄 Executando tarefas agendadas...
2024-10-23 00:00:00 📅 Verificando pagamentos atrasados...
2024-10-23 00:00:00 ✓ Nenhum pagamento atrasado encontrado
2024-10-23 00:00:00 📅 Verificando contratos próximos de expirar...
2024-10-23 00:00:00 ✓ Nenhum contrato expirando em breve
2024-10-23 00:00:00 ✅ Tarefas agendadas concluídas
```

### Erro
```log
2024-10-23 00:00:00 📅 Verificando pagamentos atrasados...
2024-10-23 00:00:01 ❌ Erro ao marcar pagamentos atrasados: database connection lost
```

## 🛑 Graceful Shutdown

O scheduler para corretamente ao desligar o servidor:

```log
🛑 Desligando servidor...
⏹️ Parando scheduler...
⏹️ Scheduler interrompido pelo contexto
✅ Servidor desligado com sucesso
```

## 🧪 Testando

### Teste Manual

```bash
# 1. Configurar intervalo curto para teste
export SCHEDULER_INTERVAL_HOURS=1

# 2. Iniciar servidor
./kitnet-manager

# 3. Monitorar logs
tail -f logs/server.log | grep "Executando tarefas"

# 4. Criar pagamento com vencimento passado
curl -X POST http://localhost:8080/api/v1/leases/...

# 5. Aguardar 1 hora (ou reiniciar servidor)

# 6. Verificar que status mudou para 'overdue'
curl http://localhost:8080/api/v1/payments/...
```

### Forçar Execução Imediata

Para testar sem esperar, basta **reiniciar o servidor** - o scheduler executa imediatamente na inicialização.

## 📈 Recomendações por Ambiente

### Produção
```env
SCHEDULER_INTERVAL_HOURS=24  # 1x ao dia (00:00)
```
- Baixo overhead
- Suficiente para maioria dos casos
- Executar à noite (menos carga)

### Staging/Homologação
```env
SCHEDULER_INTERVAL_HOURS=12  # 2x ao dia
```
- Balance entre atualização e performance
- Testes realistas

### Desenvolvimento
```env
SCHEDULER_INTERVAL_HOURS=1   # Toda hora
```
- Testes rápidos
- Feedback imediato
- ⚠️ **NÃO usar em produção** (overhead desnecessário)

## ❓ FAQ

### Por que o intervalo mínimo é 1 hora?

Para evitar overhead excessivo no banco. Se configurar < 1, será automaticamente ajustado para 24h.

### Posso desabilitar o scheduler?

Não recomendado. Sem ele, pagamentos atrasados nunca serão marcados automaticamente.

### Como saber quando foi a última execução?

Monitore os logs do servidor. Cada execução registra timestamp completo.

### E se o servidor reiniciar no meio de uma execução?

Não há problema. Na próxima inicialização, o scheduler executa imediatamente e corrige qualquer status pendente.

## 🔧 Troubleshooting

### Scheduler não está executando

```bash
# Verificar configuração
echo $SCHEDULER_INTERVAL_HOURS

# Verificar logs de inicialização
grep "Scheduler iniciado" logs/server.log
```

### Pagamentos não estão sendo marcados

```bash
# Verificar logs de erro
grep "❌" logs/server.log

# Verificar conexão com banco
curl http://localhost:8080/health
```

### Intervalo não está sendo respeitado

```bash
# Verificar valor carregado
# Deve aparecer nos logs: "⏰ Scheduler iniciado (intervalo: XXh)"
```
