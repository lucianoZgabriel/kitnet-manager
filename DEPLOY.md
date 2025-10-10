# 🚀 Guia de Deploy - Railway

Este guia contém instruções detalhadas para fazer deploy do Kitnet Manager no Railway.

## 📋 Pré-requisitos

- Conta no Railway (https://railway.app)
- Repositório GitHub conectado
- Banco Neon PostgreSQL já configurado

---

## 🛤️ Passo 1: Criar Projeto no Railway

### Via Dashboard (Recomendado)

1. Acesse https://railway.app/dashboard
2. Clique em **"New Project"**
3. Selecione **"Deploy from GitHub repo"**
4. Escolha o repositório: `lucianoZgabriel/kitnet-manager`
5. Branch: `main`
6. Railway vai detectar Go automaticamente ✅

### Via CLI (Alternativa)

```bash
# Instalar Railway CLI
npm install -g @railway/cli

# Login
railway login

# Criar projeto
railway init

# Linkar com GitHub
railway link
```

---

## ⚙️ Passo 2: Configurar Variáveis de Ambiente

No Dashboard do Railway, vá em:
**Settings → Variables → Raw Editor**

Cole as seguintes variáveis:

```bash
# Database (usar sua connection string do Neon)
DATABASE_URL=postgresql://neondb_owner:sua-senha@ep-xxx.aws.neon.tech/neondb?sslmode=require

# Server
PORT=8080
ENV=production

# Database Pool
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5
DB_MAX_LIFETIME_MINUTES=5

# JWT (GERAR NOVA SECRET!)
JWT_SECRET=gerar-um-secret-forte-aqui-min-32-caracteres
JWT_EXPIRY_HOURS=24
```

### 🔐 Como gerar JWT_SECRET seguro:

```bash
# Opção 1: OpenSSL
openssl rand -base64 32

# Opção 2: Go
go run -e 'package main; import ("crypto/rand"; "encoding/base64"; "fmt"); func main() { b := make([]byte, 32); rand.Read(b); fmt.Println(base64.StdEncoding.EncodeToString(b)) }'

# Opção 3: Online
# https://generate-secret.vercel.app/32
```

---

## 🔧 Passo 3: Configurar Build

Railway detecta Go automaticamente, mas vamos garantir:

**Settings → Build**
- Build Command: `go build -o bin/api cmd/api/main.go`
- Start Command: `./bin/api`

Ou edite o `railway.json` (já criado no projeto):
```json
{
  "build": {
    "builder": "NIXPACKS",
    "buildCommand": "go build -o bin/api cmd/api/main.go"
  },
  "deploy": {
    "startCommand": "./bin/api",
    "healthcheckPath": "/health",
    "healthcheckTimeout": 100
  }
}
```

---

## 🗄️ Passo 4: Executar Migrations

**IMPORTANTE:** Antes do primeiro deploy, execute as migrations no Neon:

### Opção 1: Local (Recomendado)
```bash
# No seu terminal local
export DATABASE_URL="sua-connection-string-do-neon"
make migrate-up
```

### Opção 2: Railway CLI
```bash
railway run make migrate-up
```

### ✅ Verificar migrations
```bash
make migrate-status
# Deve mostrar: version 5 (última migration de users)
```

---

## 🚀 Passo 5: Deploy

### Deploy Automático (Recomendado)
```bash
git add .
git commit -m "chore: configure Railway deployment"
git push origin main

# Railway detecta o push e deploya automaticamente! 🎉
```

### Deploy Manual via CLI
```bash
railway up
```

---

## 📊 Passo 6: Verificar Deploy

### 1. Acompanhar Build
- Railway Dashboard → Deployments
- Ver logs em tempo real
- Build deve levar ~1-2 minutos

### 2. Verificar Health Check
```bash
# Railway gera uma URL automática
curl https://kitnet-manager-production.up.railway.app/health

# Resposta esperada:
# {"success":true,"message":"Server is healthy","data":null}
```

### 3. Testar Login
```bash
curl -X POST https://seu-app.up.railway.app/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# Deve retornar token JWT
```

### 4. Acessar Swagger
```
https://seu-app.up.railway.app/swagger/index.html
```

---

## 🔄 Workflow de Deploy Contínuo

### Para novas features:

```bash
# 1. Desenvolver localmente
git checkout -b feature/nova-funcionalidade
# ... código ...
git commit -m "feat: nova funcionalidade"

# 2. Abrir PR
git push origin feature/nova-funcionalidade
# Abrir PR no GitHub

# 3. Revisar e aprovar PR

# 4. Merge para main
# GitHub → Merge Pull Request

# 5. Deploy automático! 🎉
# Railway detecta merge e deploya em ~1-2 min
```

---

## 🐛 Troubleshooting

### Build falhou?
```bash
# Ver logs detalhados
railway logs

# Verificar variáveis
railway variables

# Rebuild manualmente
railway up --detach
```

### App não inicia?
- Verificar se PORT está configurado como 8080
- Verificar DATABASE_URL
- Verificar se migrations foram executadas
- Ver logs: `railway logs`

### Erro 502 Bad Gateway?
- App pode estar demorando para iniciar (primeira vez ~30s)
- Verificar health check path: `/health`
- Aumentar timeout no railway.json

---

## 🔒 Segurança Pós-Deploy

### ⚠️ IMPORTANTE - Fazer após primeiro deploy:

1. **Trocar senha do admin**
   ```bash
   # Via Swagger ou curl
   POST /api/v1/auth/change-password
   {
     "old_password": "admin123",
     "new_password": "sua-senha-forte-aqui"
   }
   ```

2. **Rotacionar JWT_SECRET periodicamente**
   ```bash
   # Gerar novo secret
   openssl rand -base64 32

   # Atualizar no Railway
   # Settings → Variables → JWT_SECRET

   # Redeploy
   # Settings → Redeploy
   ```

3. **Configurar variáveis de ambiente separadas** (staging/prod)

---

## 💰 Custos Estimados

### Free Tier
- $5 grátis/mês
- Suficiente para MVP com uso interno
- ~500 horas de execução

### Pós Free Tier
- ~$5-10/mês para uso básico
- Baseado em uso real (CPU/RAM/Network)

### Monitorar uso:
Railway Dashboard → Usage → Billing

---

## 🔄 Rollback

Se algo der errado:

1. Railway Dashboard → Deployments
2. Encontrar último deploy estável
3. Click nos 3 pontos (...)
4. "Redeploy"
5. Deploy anterior volta em ~30s

---

## 📈 Próximos Passos

- [ ] Configurar domínio customizado
- [ ] Setup de monitoramento (Railway tem integrado)
- [ ] Configurar alertas
- [ ] Backup automático do banco (Neon tem built-in)
- [ ] Staging environment (branch deploy)

---

## 🆘 Suporte

- Railway Docs: https://docs.railway.app
- Railway Discord: https://discord.gg/railway
- GitHub Issues: https://github.com/lucianoZgabriel/kitnet-manager/issues
