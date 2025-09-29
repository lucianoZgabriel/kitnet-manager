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

  ## Documentação

  - [Arquitetura](kitnet_architecture.md)
  - [Roadmap](kitnet_roadmap.md)

  ## Como executar

  Em breve...

  ## Licença

  Projeto privado

  4. Verificar estrutura criada