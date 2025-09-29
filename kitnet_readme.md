# Kitnet Manager

> Sistema de gestão de condomínio de kitnets com controle de unidades, contratos, pagamentos e notificações

## 📋 Sobre o Projeto

**Kitnet Manager** é uma aplicação web desenvolvida para automatizar e otimizar a gestão de um condomínio residencial com 31 unidades tipo kitnet. O sistema substitui o controle manual via planilhas Excel, oferecendo uma solução robusta para gerenciamento de moradores, contratos, pagamentos e comunicação automatizada.

### Problema

A gestão manual através de planilhas apresenta diversos desafios:
- Dificuldade em rastrear pagamentos e inadimplência
- Falta de alertas automáticos para vencimentos
- Ausência de histórico estruturado
- Impossibilidade de gerar relatórios consolidados
- Risco de erros humanos no controle financeiro
- Dificuldade em controlar status de unidades e contratos

### Solução

Sistema web completo que oferece:
- ✅ Controle centralizado de unidades e seu status (disponível, ocupada, reforma)
- ✅ Gestão completa de moradores e contratos
- ✅ Registro e acompanhamento de pagamentos (aluguel + taxa de pintura)
- ✅ Dashboard com visão executiva do negócio
- ✅ Relatórios financeiros detalhados
- ✅ Notificações automáticas (lembretes de pagamento e renovação de contratos)
- ✅ Histórico completo de transações

## 🎯 Funcionalidades Principais

### MVP (Versão 1.0)

#### Gestão de Unidades
- Cadastro de unidades com número, andar e características
- Controle de status (disponível, ocupada, em manutenção, em reforma)
- Diferenciação entre unidades reformadas e não reformadas
- Valores de aluguel distintos por categoria

#### Gestão de Moradores
- Cadastro completo com dados pessoais e documentos
- Informações de contato (telefone, email)
- Validação de CPF único

#### Gestão de Contratos
- Contratos com duração padrão de 6 meses
- Renovação automática ao término (se não cancelado)
- Alerta 45 dias antes do vencimento para decisão de renovação
- Flexibilidade na data de assinatura vs data de início
- Definição de dia de vencimento personalizado por contrato
- Vinculação de unidade e morador

#### Gestão de Pagamentos
- Registro manual de pagamentos de aluguel
- Controle de taxa de pintura (à vista ou parcelada em até 3x)
- Tipos de pagamento separados (aluguel, taxa de pintura, ajustes)
- Status detalhado (pendente, pago, atrasado)
- Histórico completo por contrato
- Referências de PIX e comprovantes

#### Dashboard e Relatórios
- Visão geral de ocupação das unidades
- Indicadores de inadimplência
- Receita mensal (projetada vs realizada)
- Contratos próximos ao vencimento
- Relatórios financeiros por período

#### Notificações
- Alertas internos no sistema
- Lembretes SMS 3 dias antes do vencimento do aluguel
- Alertas 45 dias antes do vencimento de contratos

### Versão 2.0 (Futuro)

- Integração com gateway de SMS
- Geração automática de cobranças mensais
- Upload de comprovantes de pagamento
- Relatórios avançados e exportação
- Portal do morador (acesso limitado para consulta)
- Integração com PIX para confirmação automática

## 🛠 Tecnologias

### Backend
- **Linguagem:** Go 1.21+
- **Framework Web:** Chi Router
- **Database:** PostgreSQL (Neon - serverless)
- **Query Builder:** SQLC (type-safe SQL)
- **Migrations:** golang-migrate
- **Validação:** go-playground/validator
- **Configuração:** godotenv

### Frontend (Futuro)
- **Framework:** Next.js 14+
- **UI:** React + TailwindCSS
- **State Management:** Zustand ou React Query

### Infraestrutura
- **Database:** Neon PostgreSQL (free tier)
- **Versionamento:** Git + GitHub
- **Deploy (futuro):** Railway, Render ou Fly.io

## 📁 Estrutura do Projeto

```
kitnet-manager/
├── cmd/
│   └── api/              # Entry point da aplicação
├── internal/
│   ├── domain/           # Entidades e regras de negócio
│   ├── repository/       # Camada de dados (PostgreSQL)
│   ├── service/          # Casos de uso e lógica de negócio
│   ├── handler/          # HTTP handlers (controllers)
│   └── pkg/              # Utilitários internos
├── migrations/           # Database migrations
├── config/               # Arquivos de configuração
├── docs/                 # Documentação adicional
└── Makefile             # Comandos úteis
```

## 🚀 Como Executar

### Pré-requisitos

- Go 1.21 ou superior
- PostgreSQL (ou conta no Neon)
- Make (opcional, mas recomendado)

### Configuração

1. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/kitnet-manager.git
cd kitnet-manager
```

2. Configure as variáveis de ambiente:
```bash
cp .env.example .env
# Edite o .env com suas credenciais do Neon
```

3. Execute as migrations:
```bash
make migrate-up
```

4. Inicie o servidor:
```bash
make run
```

O servidor estará rodando em `http://localhost:8080`

## 📊 Modelo de Dados

### Entidades Principais

- **Unit (Unidade):** Representa cada kitnet do condomínio
- **Tenant (Morador):** Dados dos inquilinos
- **Lease (Contrato):** Vincula morador à unidade com período definido
- **Payment (Pagamento):** Registros de pagamentos de aluguel e taxas
- **Notification (Notificação):** Lembretes e alertas do sistema

Ver detalhes completos em [ARCHITECTURE.md](./ARCHITECTURE.md)

## 📝 Regras de Negócio

### Contratos
- Duração padrão de 6 meses com renovação automática
- Data de assinatura pode diferir da data de início
- Cada contrato define seu próprio dia de vencimento mensal
- Status da unidade muda automaticamente ao criar/cancelar contrato

### Pagamentos
- Aluguel cobrado no início do período ("paga para morar")
- Primeiro pagamento exigido antes da assinatura do contrato
- Taxa de pintura obrigatória (à vista ou 3x)
- Possibilidade de ajustes proporcionais para alterar data de vencimento

### Unidades
- Unidades reformadas têm valor de aluguel R$ 100,00 maior
- Status automático baseado em contratos ativos
- Controle de disponibilidade em tempo real

### Notificações
- SMS enviado 3 dias antes do vencimento do aluguel
- Alerta 45 dias antes do término do contrato
- Renovação automática se não houver ação manual

## 🗺 Roadmap

Ver o roadmap completo e detalhado em [ROADMAP.md](./ROADMAP.md)

### Fase Atual: MVP Development
- ✅ Planejamento e modelagem
- 🔄 Sprint 0: Setup e infraestrutura
- ⏳ Sprint 1: CRUD de Unidades e Moradores
- ⏳ Sprint 2: Gestão de Contratos
- ⏳ Sprint 3: Sistema de Pagamentos
- ⏳ Sprint 4: Dashboard e Relatórios
- ⏳ Sprint 5: Sistema de Notificações

## 🤝 Contribuindo

Este é um projeto pessoal de aprendizado, mas sugestões são bem-vindas!

## 📄 Licença

Este projeto é de uso pessoal e educacional.

## 👤 Autor

Desenvolvido como projeto de aprendizado em Go e gestão de software.

---

**Status do Projeto:** 🚧 Em desenvolvimento ativo (MVP)

**Última atualização:** Setembro 2025
