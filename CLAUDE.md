# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## IMPORTANT: Educational Approach

**This project has an educational purpose. The working method is:**
- Claude provides guidance, explanations, and code snippets
- The user writes and executes all code
- Claude explains architectural decisions and best practices
- The user provides feedback after implementation
- Claude NEVER directly creates or modifies code files
- Focus on teaching industry best practices and patterns

## Project Status

This is a **Kitnet Manager** system for managing a 31-unit apartment complex. The project is currently in the **planning phase** with comprehensive documentation completed. Implementation follows Clean Architecture principles with Go backend.

## Tech Stack

- **Backend:** Go 1.21+ with Chi Router
- **Database:** PostgreSQL (Neon cloud) with SQLC for type-safe queries
- **Migrations:** golang-migrate
- **Validation:** go-playground/validator
- **Future Frontend:** Next.js 14+ with React and TailwindCSS

## Commands

### Development Commands (to be implemented)
```bash
# Initialize Go module
go mod init github.com/lucianogabriel/kitnet-manager

# Run application
make run

# Build binary
make build

# Run tests
make test

# Database migrations
make migrate-up
make migrate-down

# Generate SQLC code
make sqlc-generate

# Development with hot reload
air
```

### Database Setup
```bash
# Connect to Neon PostgreSQL
psql "postgresql://[user]:[password]@[host]/[database]?sslmode=require"
```

## Architecture

The project follows **Clean Architecture** with clear separation of concerns:

```
cmd/api/          → Application entry point
internal/
  domain/         → Business entities (Unit, Tenant, Lease, Payment)
  repository/     → Data access interfaces and implementations
  service/        → Business logic and use cases
  handler/        → HTTP handlers/controllers
  pkg/           → Shared utilities
migrations/       → Database migration files
config/          → Configuration files
```

### Key Design Principles

1. **Dependency Inversion:** Handlers depend on service interfaces, services depend on repository interfaces
2. **Domain-Centric:** Business logic isolated from infrastructure concerns
3. **SQLC over ORM:** Direct SQL control with type safety, no ORM overhead
4. **Repository Pattern:** Data access abstraction for testability

## Core Business Entities

- **Unit:** Apartment units with status (occupied/available/maintenance)
- **Tenant:** Residents with CPF validation and contact info
- **Lease:** 6-month contracts with auto-renewal capability
- **Payment:** Rent (R$800) and painting fee (R$250) tracking
- **PaymentSchedule:** Monthly payment expectations and status

## Implementation Roadmap

Current Sprint Plan (from kitnet_roadmap.md):

1. **Sprint 0:** Project setup, database schema, initial migrations
2. **Sprint 1:** Units and Tenants CRUD operations
3. **Sprint 2:** Lease management with contract lifecycle
4. **Sprint 3:** Payment system with tracking and validation
5. **Sprint 4:** Dashboard and financial reports
6. **Sprint 5:** SMS notifications via Twilio/similar
7. **Sprint 6:** Testing, refinements, and MVP deployment

## Database Schema Highlights

Key tables with relationships:
- `units` → Physical apartments with renovation status
- `tenants` → Resident information with unique CPF
- `leases` → Active/inactive contracts linking units and tenants
- `payments` → Financial transactions with type and status
- `payment_schedules` → Expected vs actual payment tracking

## API Structure

RESTful endpoints following pattern:
```
GET    /api/v1/units          → List all units
POST   /api/v1/units          → Create unit
GET    /api/v1/units/:id      → Get unit details
PUT    /api/v1/units/:id      → Update unit
DELETE /api/v1/units/:id      → Delete unit

Similar patterns for: /tenants, /leases, /payments
```

## Key Business Rules

1. **Lease Duration:** Fixed 6-month terms with renewal option
2. **Payment Types:** Monthly rent (R$800) + one-time painting fee (R$250)
3. **Late Fees:** 2% penalty + 1% monthly interest after due date
4. **CPF Validation:** Required for all tenants with format checking
5. **Unit Status:** Must be 'available' before new lease creation
6. **Payment Grace Period:** 5 days after due date before penalties

## Development Priorities

1. Focus on MVP features first (basic CRUD, payment tracking)
2. Implement comprehensive input validation at handler level
3. Use transactions for operations affecting multiple entities
4. Add proper error handling with meaningful messages
5. Implement logging for audit trail
6. Write unit tests for services and integration tests for handlers