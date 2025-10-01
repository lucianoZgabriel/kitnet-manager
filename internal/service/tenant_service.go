package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository"
)

// Service layer errors
var (
	ErrTenantNotFound                       = errors.New("tenant not found")
	ErrCPFAlreadyExists                     = errors.New("CPF already registered")
	ErrCannotDeleteTenantWithActiveContract = errors.New("cannot delete tenant with active contract")
)

// TenantService contém a lógica de negócio para gestão de moradores
type TenantService struct {
	tenantRepo repository.TenantRepository
	// leaseRepo será adicionado futuramente para validar contratos ativos
}

// NewTenantService cria uma nova instância do serviço de moradors
func NewTenantService(tenantRepo repository.TenantRepository) *TenantService {
	return &TenantService{
		tenantRepo: tenantRepo,
	}
}

// CreateTenant cria um novo morador com validações de negócio
func (s *TenantService) CreateTenant(ctx context.Context, fullName, cpf, phone, email, idDocType, idDocNumber string) (*domain.Tenant, error) {
	// Verifica se CPF já existe
	exists, err := s.tenantRepo.ExistsByCPF(ctx, cpf)
	if err != nil {
		return nil, fmt.Errorf("error checking CPF: %w", err)
	}
	if exists {
		return nil, ErrCPFAlreadyExists
	}

	// Cria novo morador usando o domain model
	tenant, err := domain.NewTenant(fullName, cpf, phone, email)
	if err != nil {
		return nil, fmt.Errorf("error creating tenant: %w", err)
	}

	// Define documentos se fornecidos
	if idDocType != "" {
		tenant.IDDocumentType = idDocType
	}
	if idDocNumber != "" {
		tenant.IDDocumentNumber = idDocNumber
	}

	// Persistir no banco
	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, fmt.Errorf("error saving tenant: %w", err)
	}

	return tenant, nil
}

// GetTenantByID busca um morador pelo ID
func (s *TenantService) GetTenantByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting tenant: %w", err)
	}
	if tenant == nil {
		return nil, ErrTenantNotFound
	}

	return tenant, nil
}

// GetTenantByCPF busca um morador pelo CPF
func (s *TenantService) GetTenantByCPF(ctx context.Context, cpf string) (*domain.Tenant, error) {
	tenant, err := s.tenantRepo.GetByCPF(ctx, cpf)
	if err != nil {
		return nil, fmt.Errorf("error getting tenant: %w", err)
	}
	if tenant == nil {
		return nil, ErrTenantNotFound
	}

	return tenant, nil
}

// ListTenants retorna todos os moradores
func (s *TenantService) ListTenants(ctx context.Context) ([]*domain.Tenant, error) {
	tenants, err := s.tenantRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing tenants: %w", err)
	}
	return tenants, nil
}

// SearchTenantsByName busca moradores por nome (case-insensitive)
func (s *TenantService) SearchTenantsByName(ctx context.Context, name string) ([]*domain.Tenant, error) {
	tenants, err := s.tenantRepo.SearchByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("error searching tenants: %w", err)
	}
	return tenants, nil
}

// UpdateTenant atualiza um morador existente
func (s *TenantService) UpdateTenant(ctx context.Context, id uuid.UUID, fullName, phone, email, idDocType, idDocNumber string) (*domain.Tenant, error) {
	// Busca morador
	tenant, err := s.GetTenantByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Atualizar usando método do domain
	if err := tenant.UpdateInfo(fullName, phone, email, idDocType, idDocNumber); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Persistir mudanças
	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return nil, fmt.Errorf("error updating tenant: %w", err)
	}

	return tenant, nil
}

// DeleteTenant remove um morador
func (s *TenantService) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	// Buscar morador
	tenant, err := s.GetTenantByID(ctx, id)
	if err != nil {
		return err
	}

	// TODO: Regra de negócio: não pode deletar morador com contrato ativo
	// Isso será implementado quando criarmos o LeaseRepository
	// Por enquanto, deixamos um placeholder
	_ = tenant

	// Deletar
	if err := s.tenantRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting tenant: %w", err)
	}

	return nil
}

// GetTenantCount retorna o total de moradores cadastrados
func (s *TenantService) GetTenantCount(ctx context.Context) (int64, error) {
	count, err := s.tenantRepo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("error counting tenants: %w", err)
	}
	return count, nil
}
