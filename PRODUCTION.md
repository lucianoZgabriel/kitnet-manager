# 🚀 Kitnet Manager - Produção

Documentação da aplicação em produção no Railway.

---

## 🌐 URLs de Produção

### API Base URL
```
https://kitnet-manager-production.up.railway.app
```

### Endpoints Principais

#### Health Check
```
GET https://kitnet-manager-production.up.railway.app/health
```
**Status:** ✅ Online

#### Swagger Documentation
```
https://kitnet-manager-production.up.railway.app/swagger/index.html
```

#### API v1
```
https://kitnet-manager-production.up.railway.app/api/v1/
```

---

## 🔐 Credenciais Padrão

⚠️ **ALTERAR IMEDIATAMENTE APÓS PRIMEIRO ACESSO**

```
Username: admin
Password: admin123
```

### Como trocar a senha:

1. Acesse o Swagger: https://kitnet-manager-production.up.railway.app/swagger/index.html
2. Faça login (POST `/api/v1/auth/login`)
3. Autorize no Swagger (botão "Authorize" com o token)
4. Use POST `/api/v1/auth/change-password`:
   ```json
   {
     "old_password": "admin123",
     "new_password": "sua-nova-senha-forte"
   }
   ```

---

## 📊 Status do Deploy

### Última Deploy
- **Data:** 10/10/2025
- **Branch:** main
- **Commit:** Deployment configuration and authentication system
- **Status:** ✅ Online

### Build Info
- **Platform:** Railway
- **Region:** US East
- **Runtime:** Go 1.21+
- **Database:** Neon PostgreSQL (cloud)

### Health Check
```bash
curl https://kitnet-manager-production.up.railway.app/health
```
**Response:**
```json
{
  "success": true,
  "message": "Server is healthy"
}
```

---

## 🗄️ Banco de Dados

### Neon PostgreSQL
- **Status:** ✅ Conectado
- **Migrations:** 5/5 aplicadas
- **Host:** ep-flat-shape-adximx7a-pooler.c-2.us-east-1.aws.neon.tech
- **Database:** neondb
- **SSL:** Required

### Schema Atual
1. ✅ Units table (000001)
2. ✅ Tenants table (000002)
3. ✅ Leases table (000003)
4. ✅ Payments table (000004)
5. ✅ Users table (000005)

---

## 🔄 Workflow de Deploy

### Auto-Deploy Ativado
Qualquer push para `main` faz deploy automático:

```bash
# 1. Desenvolvimento local
git checkout -b feature/nova-funcionalidade
# ... código ...

# 2. Commit e push
git commit -m "feat: nova funcionalidade"
git push origin feature/nova-funcionalidade

# 3. Criar PR e merge para main
gh pr create
gh pr merge

# 4. Railway detecta e deploya automaticamente (1-2 min)
```

### Deploy Manual via Railway CLI
```bash
railway up
```

### Rollback
Railway Dashboard → Deployments → (deploy anterior) → Redeploy

---

## ⚙️ Variáveis de Ambiente

Configuradas no Railway:

```bash
DATABASE_URL=postgresql://...
PORT=8080
ENV=production
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5
DB_MAX_LIFETIME_MINUTES=5
JWT_SECRET=*** (configurado)
JWT_EXPIRY_HOURS=24
```

---

## 🧪 Testes Rápidos

### 1. Health Check
```bash
curl https://kitnet-manager-production.up.railway.app/health
```

### 2. Login
```bash
curl -X POST https://kitnet-manager-production.up.railway.app/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 3. Get Current User (com token)
```bash
TOKEN="seu-token-aqui"
curl https://kitnet-manager-production.up.railway.app/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Listar Unidades
```bash
curl https://kitnet-manager-production.up.railway.app/api/v1/units \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📈 Monitoramento

### Railway Dashboard
- **Metrics:** CPU, RAM, Network usage
- **Logs:** Real-time logs
- **Deployments:** History and status
- **Usage:** Billing and credits

### Health Endpoint
Railway faz health checks automáticos em:
```
GET /health
```
A cada 30 segundos.

---

## 🔒 Segurança

### ✅ Implementado
- [x] JWT authentication
- [x] HTTPS/SSL automático (Railway)
- [x] Bcrypt password hashing
- [x] Role-based access control
- [x] Environment variables para secrets
- [x] Database connection com SSL

### ⚠️ TODO Pós-Deploy
- [ ] Trocar senha do admin
- [ ] Configurar rate limiting (futuro)
- [ ] Setup de backup automático
- [ ] Configurar alertas de monitoramento
- [ ] Rotacionar JWT_SECRET periodicamente

---

## 💰 Custos

### Railway Free Tier
- $5 grátis/mês
- Suficiente para MVP com uso interno
- ~500 horas de execução

### Uso Atual
- Verificar em: Railway Dashboard → Usage

### Neon Database
- Free tier (sempre grátis)
- 3GB storage
- 1 branch

---

## 🆘 Troubleshooting

### App não responde?
1. Verificar logs: Railway Dashboard → Logs
2. Verificar variáveis de ambiente
3. Verificar health do database (Neon dashboard)

### Erro 502 Bad Gateway?
- App pode estar reiniciando (~30s)
- Verificar se build teve sucesso

### Erro de autenticação?
- Verificar JWT_SECRET está configurado
- Verificar DATABASE_URL está correto
- Verificar se migrations foram executadas

### Como ver logs?
```bash
# Via CLI
railway logs

# Ou no dashboard
Railway → Deployments → Ver logs
```

---

## 📞 Suporte

### Documentação
- [Deploy Guide](./DEPLOY.md)
- [API Documentation](https://kitnet-manager-production.up.railway.app/swagger/index.html)
- [Railway Docs](https://docs.railway.app)

### Links Úteis
- **Railway Dashboard:** https://railway.app/dashboard
- **Neon Dashboard:** https://console.neon.tech
- **GitHub Repo:** https://github.com/lucianoZgabriel/kitnet-manager

---

## 🎯 Próximos Passos

### Imediato
- [x] Deploy em produção
- [x] Gerar domínio público
- [x] Testar endpoints
- [ ] Trocar senha do admin
- [ ] Documentar para equipe

### Curto Prazo (1-2 semanas)
- [ ] Desenvolver frontend
- [ ] Criar usuários de produção
- [ ] Popular dados iniciais (unidades, moradores)
- [ ] Treinar usuários

### Médio Prazo (1 mês)
- [ ] Custom domain (opcional)
- [ ] Monitoramento avançado
- [ ] Backups automatizados
- [ ] Sistema de notificações (Sprint 6)

---

**Última atualização:** 10/10/2025
**Status:** ✅ Produção estável
**Versão:** 1.0.0 (MVP)
