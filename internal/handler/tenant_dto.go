package handler

import (
	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
)

// CreateTenantRequest representa o payload para criar um morador
type CreateTenantRequest struct {
	FullName         string `json:"full_name" validate:"required,min=3,max=255"`
	CPF              string `json:"cpf" validate:"required,len=14"` // Formato: XXX.XXX.XXX-XX
	Phone            string `json:"phone" validate:"required,min=10,max=20"`
	Email            string `json:"email" validate:"omitempty,email,max=255"`
	IDDocumentType   string `json:"id_document_type" validate:"omitempty,max=10"`
	IDDocumentNumber string `json:"id_document_number" validate:"omitempty,max=50"`
}

// UpdateTenantRequest representa o payload para atualizar um morador
type UpdateTenantRequest struct {
	FullName         string `json:"full_name" validate:"required,min=3,max=255"`
	Phone            string `json:"phone" validate:"required,min=10,max=20"`
	Email            string `json:"email" validate:"omitempty,email,max=255"`
	IDDocumentType   string `json:"id_document_type" validate:"omitempty,max=10"`
	IDDocumentNumber string `json:"id_document_number" validate:"omitempty,max=50"`
}

// TenantResponse representa a resposta com dados de um morador
type TenantResponse struct {
	ID               uuid.UUID `json:"id"`
	FullName         string    `json:"full_name"`
	CPF              string    `json:"cpf"`
	Phone            string    `json:"phone"`
	Email            string    `json:"email,omitempty"`
	IDDocumentType   string    `json:"id_document_type,omitempty"`
	IDDocumentNumber string    `json:"id_document_number,omitempty"`
	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
}

// ToTenantResponse converte domain.Tenant para TenantResponse
func ToTenantResponse(tenant *domain.Tenant) *TenantResponse {
	return &TenantResponse{
		ID:               tenant.ID,
		FullName:         tenant.FullName,
		CPF:              tenant.CPF,
		Phone:            tenant.Phone,
		Email:            tenant.Email,
		IDDocumentType:   tenant.IDDocumentType,
		IDDocumentNumber: tenant.IDDocumentNumber,
		CreatedAt:        tenant.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        tenant.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToTenantResponseList converte slice de tenants para slice de responses
func ToTenantResponseList(tenants []*domain.Tenant) []*TenantResponse {
	responses := make([]*TenantResponse, len(tenants))
	for i, tenant := range tenants {
		responses[i] = ToTenantResponse(tenant)
	}
	return responses
}
