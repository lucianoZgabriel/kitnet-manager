// Tenant Types

export interface Tenant {
  id: string
  full_name: string
  cpf: string // Format: XXX.XXX.XXX-XX
  phone: string
  email?: string
  id_document_type?: string
  id_document_number?: string
  created_at: string
  updated_at: string
}

export interface CreateTenantRequest {
  full_name: string
  cpf: string // Format: XXX.XXX.XXX-XX
  phone: string
  email?: string
  id_document_type?: string
  id_document_number?: string
}

export interface UpdateTenantRequest {
  full_name: string
  phone: string
  email?: string
  id_document_type?: string
  id_document_number?: string
}
