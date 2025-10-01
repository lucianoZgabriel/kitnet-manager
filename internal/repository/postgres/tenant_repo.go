package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository/sqlc"
)

// Compile-time check to ensure TenantRepository implements repository.TenantRepository
var _ repository.TenantRepository = (*TenantRepository)(nil)

// TenantRepository implementa repository.TenantRepository usando PostgreSQL via SQLC
type TenantRepository struct {
	queries *sqlc.Queries
	db      *sql.DB
}

// NewTenantRepository cria uma nova instância do repositório de moradores
func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{
		queries: sqlc.New(db),
		db:      db,
	}
}

// Create insere um novo morador no banco de dados
func (r *TenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	params := sqlc.CreateTenantParams{
		ID:               tenant.ID,
		FullName:         tenant.FullName,
		Cpf:              tenant.CPF,
		Phone:            tenant.Phone,
		Email:            toNullString(tenant.Email),
		IDDocumentType:   toNullString(tenant.IDDocumentType),
		IDDocumentNumber: toNullString(tenant.IDDocumentNumber),
		CreatedAt:        tenant.CreatedAt,
		UpdatedAt:        tenant.UpdatedAt,
	}

	created, err := r.queries.CreateTenant(ctx, params)
	if err != nil {
		return err
	}

	// Atualiza o objeto tenant com dados do banco
	tenant.CreatedAt = created.CreatedAt
	tenant.UpdatedAt = created.UpdatedAt

	return nil
}

// GetByID busca um morador pelo ID
func (r *TenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	dbTenant, err := r.queries.GetTenantByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Não encontrou, retorna nil sem erro
		}
		return nil, err
	}

	return r.toDomain(dbTenant), nil
}

// GetByCPF busca um morador pelo CPF
func (r *TenantRepository) GetByCPF(ctx context.Context, cpf string) (*domain.Tenant, error) {
	dbTenant, err := r.queries.GetTenantByCPF(ctx, cpf)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(dbTenant), nil
}

// List retorna todos os moradores
func (r *TenantRepository) List(ctx context.Context) ([]*domain.Tenant, error) {
	dbTenants, err := r.queries.ListTenants(ctx)
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(dbTenants), nil
}

// SearchByName busca moradores por nome (case-insensitive)
func (r *TenantRepository) SearchByName(ctx context.Context, name string) ([]*domain.Tenant, error) {
	// SQLC gerou o parâmetro como sql.NullString, então precisamos converter
	dbTenants, err := r.queries.SearchTenantsByName(ctx, toNullString(name))
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(dbTenants), nil
}

// Update atualiza um morador existente
func (r *TenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	tenant.UpdatedAt = time.Now()

	params := sqlc.UpdateTenantParams{
		ID:               tenant.ID,
		FullName:         tenant.FullName,
		Phone:            tenant.Phone,
		Email:            toNullString(tenant.Email),
		IDDocumentType:   toNullString(tenant.IDDocumentType),
		IDDocumentNumber: toNullString(tenant.IDDocumentNumber),
		UpdatedAt:        tenant.UpdatedAt,
	}

	updated, err := r.queries.UpdateTenant(ctx, params)
	if err != nil {
		return err
	}

	tenant.UpdatedAt = updated.UpdatedAt
	return nil
}

// Delete remove um morador do banco de dados
func (r *TenantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteTenant(ctx, id)
}

// Count retorna o total de moradores
func (r *TenantRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountTenants(ctx)
}

// ExistsByCPF verifica se já existe um morador com o CPF
func (r *TenantRepository) ExistsByCPF(ctx context.Context, cpf string) (bool, error) {
	return r.queries.TenantExistsByCPF(ctx, cpf)
}

// toDomain converte sqlc.Tenant para domain.Tenant
func (r *TenantRepository) toDomain(dbTenant sqlc.Tenant) *domain.Tenant {
	return &domain.Tenant{
		ID:               dbTenant.ID,
		FullName:         dbTenant.FullName,
		CPF:              dbTenant.Cpf,
		Phone:            dbTenant.Phone,
		Email:            dbTenant.Email.String,
		IDDocumentType:   dbTenant.IDDocumentType.String,
		IDDocumentNumber: dbTenant.IDDocumentNumber.String,
		CreatedAt:        dbTenant.CreatedAt,
		UpdatedAt:        dbTenant.UpdatedAt,
	}
}

// toDomainSlice converte []sqlc.Tenant para []*domain.Tenant
func (r *TenantRepository) toDomainSlice(dbTenants []sqlc.Tenant) []*domain.Tenant {
	tenants := make([]*domain.Tenant, len(dbTenants))
	for i, dbTenant := range dbTenants {
		tenants[i] = r.toDomain(dbTenant)
	}
	return tenants
}
