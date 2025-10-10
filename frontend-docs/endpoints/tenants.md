# Tenants Endpoints

All endpoints require authentication via Bearer token.

Base URL: `https://kitnet-manager-production.up.railway.app/api/v1`

## Create Tenant

```typescript
POST /tenants
```

**Request Body:**
```json
{
  "full_name": "João Silva Santos",
  "cpf": "123.456.789-00",
  "phone": "(11) 98765-4321",
  "email": "joao.silva@email.com",
  "id_document_type": "RG",
  "id_document_number": "12.345.678-9"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Tenant created successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "full_name": "João Silva Santos",
    "cpf": "123.456.789-00",
    "phone": "(11) 98765-4321",
    "email": "joao.silva@email.com",
    "id_document_type": "RG",
    "id_document_number": "12.345.678-9",
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

**Validation Rules:**
- `full_name`: required, cannot be empty
- `cpf`: required, format `XXX.XXX.XXX-XX`, must be unique
- `phone`: required, cannot be empty
- `email`: optional, must be valid email format if provided
- `id_document_type`: optional
- `id_document_number`: optional

**Errors:**
- `400` - Invalid request body or validation error
- `409` - CPF already registered
- `500` - Internal server error

---

## Get Tenant by ID

```typescript
GET /tenants/:id
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Tenant retrieved successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "full_name": "João Silva Santos",
    "cpf": "123.456.789-00",
    "phone": "(11) 98765-4321",
    "email": "joao.silva@email.com",
    "id_document_type": "RG",
    "id_document_number": "12.345.678-9",
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

**Errors:**
- `400` - Invalid tenant ID
- `404` - Tenant not found

---

## Get Tenant by CPF

```typescript
GET /tenants/cpf?cpf=123.456.789-00
```

**Query Parameters:**
- `cpf` (required): CPF in format `XXX.XXX.XXX-XX`

**Response (200 OK):** Same as Get Tenant by ID

**Errors:**
- `400` - CPF parameter missing or invalid
- `404` - Tenant not found

---

## List Tenants

```typescript
GET /tenants
GET /tenants?name=joão
```

**Query Parameters:**
- `name` (optional): Search by name (case-insensitive, partial match)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Tenants retrieved successfully",
  "data": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "full_name": "João Silva Santos",
      "cpf": "123.456.789-00",
      "phone": "(11) 98765-4321",
      "email": "joao.silva@email.com",
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-15T10:00:00Z"
    }
  ]
}
```

**Errors:**
- `500` - Internal server error

---

## Update Tenant

```typescript
PUT /tenants/:id
```

**Request Body:**
```json
{
  "full_name": "João Silva Santos Junior",
  "phone": "(11) 99999-8888",
  "email": "joao.junior@email.com",
  "id_document_type": "CNH",
  "id_document_number": "12345678900"
}
```

**Note:** CPF cannot be updated (immutable)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Tenant updated successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "full_name": "João Silva Santos Junior",
    "cpf": "123.456.789-00",
    "phone": "(11) 99999-8888",
    "email": "joao.junior@email.com",
    "id_document_type": "CNH",
    "id_document_number": "12345678900",
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T12:00:00Z"
  }
}
```

**Errors:**
- `400` - Invalid request body or validation error
- `404` - Tenant not found

---

## Delete Tenant

```typescript
DELETE /tenants/:id
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Tenant deleted successfully",
  "data": null
}
```

**Errors:**
- `400` - Cannot delete tenant with active lease
- `404` - Tenant not found
- `500` - Internal server error

---

## CPF Format

The system expects CPF in the format: `XXX.XXX.XXX-XX`

**Examples:**
- ✅ Valid: `123.456.789-00`
- ❌ Invalid: `12345678900`
- ❌ Invalid: `123.456.789`
