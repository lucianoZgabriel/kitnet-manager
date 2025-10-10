# Kitnet Manager API Documentation

Complete API documentation for the Kitnet Manager frontend development.

---

## Base Information

**Production URL:** `https://kitnet-manager-production.up.railway.app`
**API Version:** `v1`
**Base Path:** `/api/v1`

---

## Quick Start

### 1. Authentication

All endpoints (except login) require a JWT Bearer token:

```typescript
const headers = {
  'Authorization': `Bearer ${token}`,
  'Content-Type': 'application/json'
}
```

### 2. Login Example

```typescript
const login = async (username: string, password: string) => {
  const response = await fetch(
    'https://kitnet-manager-production.up.railway.app/api/v1/auth/login',
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password })
    }
  )

  const data = await response.json()

  if (data.success) {
    // Store token
    localStorage.setItem('token', data.data.token)
    return data.data
  }

  throw new Error(data.error)
}
```

### 3. Authenticated Request Example

```typescript
const fetchUnits = async () => {
  const token = localStorage.getItem('token')

  const response = await fetch(
    'https://kitnet-manager-production.up.railway.app/api/v1/units',
    {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    }
  )

  const data = await response.json()
  return data.data // Array of units
}
```

---

## Response Format

All API responses follow this structure:

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { /* response data */ }
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message",
  "data": null
}
```

---

## API Endpoints Overview

### Authentication (`/auth`)
- `POST /auth/login` - Login and get JWT token
- `GET /auth/me` - Get current user info
- `POST /auth/refresh` - Refresh JWT token
- `POST /auth/users` - Create user (admin only)
- `GET /auth/users` - List users (admin only)
- `GET /auth/users/:id` - Get user by ID (admin only)
- `POST /auth/change-password` - Change password
- `PATCH /auth/users/:id/role` - Change user role (admin only)
- `POST /auth/users/:id/deactivate` - Deactivate user (admin only)
- `POST /auth/users/:id/activate` - Activate user (admin only)

### Units (`/units`)
- `POST /units` - Create unit
- `GET /units` - List units (with filters)
- `GET /units/:id` - Get unit by ID
- `PUT /units/:id` - Update unit
- `PATCH /units/:id/status` - Update unit status
- `DELETE /units/:id` - Delete unit
- `GET /units/stats/occupancy` - Get occupancy statistics

### Tenants (`/tenants`)
- `POST /tenants` - Create tenant
- `GET /tenants` - List tenants (with search)
- `GET /tenants/:id` - Get tenant by ID
- `GET /tenants/cpf?cpf=XXX.XXX.XXX-XX` - Get tenant by CPF
- `PUT /tenants/:id` - Update tenant
- `DELETE /tenants/:id` - Delete tenant

### Leases (`/leases`)
- `POST /leases` - Create lease (generates payments)
- `GET /leases` - List leases (with filters)
- `GET /leases/:id` - Get lease by ID
- `GET /leases/stats` - Get lease statistics
- `GET /leases/expiring-soon` - Get leases expiring in 45 days
- `POST /leases/:id/renew` - Renew lease
- `POST /leases/:id/cancel` - Cancel lease
- `PATCH /leases/:id/painting-fee` - Update painting fee paid

### Payments (`/payments`)
- `GET /payments/:id` - Get payment by ID
- `GET /leases/:lease_id/payments` - Get payments by lease
- `GET /payments/overdue` - Get overdue payments
- `GET /payments/upcoming?days=7` - Get upcoming payments
- `PUT /payments/:id/pay` - Mark payment as paid
- `POST /payments/:id/cancel` - Cancel payment
- `GET /leases/:lease_id/payments/stats` - Get payment statistics

### Dashboard (`/dashboard`)
- `GET /dashboard` - Get consolidated dashboard metrics

### Reports (`/reports`)
- `GET /reports/financial?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD` - Financial report
- `GET /reports/payments` - Payment history report

### Health Check
- `GET /health` - Check API health status

---

## Detailed Documentation

For detailed endpoint documentation, see:

- [📁 types/](./types/) - TypeScript type definitions
  - [auth.ts](./types/auth.ts) - Authentication & User types
  - [unit.ts](./types/unit.ts) - Unit types
  - [tenant.ts](./types/tenant.ts) - Tenant types
  - [lease.ts](./types/lease.ts) - Lease types
  - [payment.ts](./types/payment.ts) - Payment types
  - [dashboard.ts](./types/dashboard.ts) - Dashboard & Report types

- [📁 endpoints/](./endpoints/) - Endpoint documentation
  - [auth.md](./endpoints/auth.md) - Authentication endpoints
  - [units.md](./endpoints/units.md) - Units endpoints
  - [tenants.md](./endpoints/tenants.md) - Tenants endpoints
  - [leases.md](./endpoints/leases.md) - Leases endpoints
  - [payments.md](./endpoints/payments.md) - Payments endpoints
  - [dashboard.md](./endpoints/dashboard.md) - Dashboard & Reports endpoints

- [📄 validation-rules.md](./validation-rules.md) - Validation rules & business logic

---

## Authentication Flow

```
1. User submits credentials
   POST /auth/login
   { username, password }

2. API validates and returns token
   { token: "jwt...", user: {...} }

3. Frontend stores token
   localStorage.setItem('token', token)

4. All subsequent requests include token
   headers: { Authorization: 'Bearer jwt...' }

5. Token expires after 24h (configurable)
   - Use POST /auth/refresh to get new token
   - Or redirect to login
```

---

## User Roles & Permissions

| Role | Read | Write (Create/Update/Delete) | Manage Users |
|------|------|------------------------------|--------------|
| **admin** | ✅ | ✅ | ✅ |
| **manager** | ✅ | ✅ | ❌ |
| **viewer** | ✅ | ❌ | ❌ |

### Default Credentials
```
Username: admin
Password: admin123
```

**⚠️ Important:** Change the default admin password immediately in production!

---

## Common Workflows

### 1. Create New Lease

```typescript
// Step 1: Ensure unit is available
const units = await fetch('/api/v1/units?status=available')

// Step 2: Select or create tenant
const tenant = await fetch('/api/v1/tenants', {
  method: 'POST',
  body: JSON.stringify({
    full_name: "João Silva",
    cpf: "123.456.789-00",
    phone: "(11) 98765-4321"
  })
})

// Step 3: Create lease
const lease = await fetch('/api/v1/leases', {
  method: 'POST',
  body: JSON.stringify({
    unit_id: selectedUnit.id,
    tenant_id: tenant.data.id,
    contract_signed_date: "2025-01-15",
    start_date: "2025-02-01",
    payment_due_day: 5,
    monthly_rent_value: "1000.00",
    painting_fee_total: "250.00",
    painting_fee_installments: 2
  })
})

// Response includes lease + auto-generated payments
```

### 2. Process Payment

```typescript
// Step 1: Get payments for a lease
const payments = await fetch(`/api/v1/leases/${leaseId}/payments`)

// Step 2: Mark payment as paid
const updatedPayment = await fetch(`/api/v1/payments/${paymentId}/pay`, {
  method: 'PUT',
  body: JSON.stringify({
    payment_date: "2025-02-03",
    payment_method: "pix"
  })
})
```

### 3. Monitor Dashboard

```typescript
// Get all dashboard data
const dashboard = await fetch('/api/v1/dashboard')

// dashboard.data contains:
// - occupancy metrics
// - financial metrics
// - contract metrics
// - alerts

// Display alerts by priority
const sortedAlerts = dashboard.data.alerts.sort((a, b) => {
  const severityOrder = { high: 3, medium: 2, low: 1 }
  return severityOrder[b.severity] - severityOrder[a.severity]
})
```

---

## Error Handling

### HTTP Status Codes

| Code | Meaning | Action |
|------|---------|--------|
| 200 | OK | Success - process data |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Validation error - show user error message |
| 401 | Unauthorized | Invalid/expired token - redirect to login |
| 404 | Not Found | Resource doesn't exist - show not found message |
| 409 | Conflict | Duplicate entry (CPF, username) - show conflict message |
| 500 | Server Error | Server issue - show generic error, retry later |

### Error Handling Example

```typescript
const handleApiError = (error: any) => {
  switch (error.status) {
    case 400:
      // Validation error - show to user
      toast.error(error.message)
      break

    case 401:
      // Unauthorized - redirect to login
      localStorage.removeItem('token')
      router.push('/login')
      break

    case 404:
      // Not found
      toast.error('Resource not found')
      break

    case 409:
      // Conflict
      toast.error('This record already exists')
      break

    case 500:
      // Server error
      toast.error('Server error. Please try again later.')
      break

    default:
      toast.error('An unexpected error occurred')
  }
}
```

---

## Testing with Swagger

Interactive API documentation available at:
```
https://kitnet-manager-production.up.railway.app/swagger/index.html
```

You can:
1. Explore all endpoints
2. See request/response schemas
3. Test endpoints directly in the browser
4. Authenticate with JWT token

---

## Best Practices

### 1. Token Management

```typescript
// Store token securely
localStorage.setItem('token', token)

// Include in all requests
const apiClient = axios.create({
  baseURL: 'https://kitnet-manager-production.up.railway.app/api/v1',
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('token')}`
  }
})

// Refresh token before expiry
const checkTokenExpiry = () => {
  // Decode JWT and check exp claim
  // Refresh if expires in < 1 hour
}
```

### 2. Data Formatting

```typescript
// Always format decimals with 2 decimal places
const formatMoney = (value: number) => value.toFixed(2)

// Format dates to ISO
const formatDate = (date: Date) => date.toISOString().split('T')[0]

// Format CPF with mask
const formatCPF = (cpf: string) => {
  return cpf
    .replace(/\D/g, '')
    .replace(/(\d{3})(\d)/, '$1.$2')
    .replace(/(\d{3})(\d)/, '$1.$2')
    .replace(/(\d{3})(\d{1,2})$/, '$1-$2')
}
```

### 3. Optimistic Updates

```typescript
// Update UI immediately, rollback on error
const updateUnit = async (id: string, data: UpdateUnitRequest) => {
  // Update local state
  setUnits(prev => prev.map(u => u.id === id ? { ...u, ...data } : u))

  try {
    await api.put(`/units/${id}`, data)
  } catch (error) {
    // Rollback on error
    setUnits(prev => /* restore original state */)
    throw error
  }
}
```

### 4. Caching Strategy

```typescript
// Example with React Query
import { useQuery } from '@tanstack/react-query'

const useUnits = () => {
  return useQuery({
    queryKey: ['units'],
    queryFn: fetchUnits,
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: true
  })
}
```

---

## Support & Resources

- **Production API:** https://kitnet-manager-production.up.railway.app
- **Swagger Docs:** https://kitnet-manager-production.up.railway.app/swagger/index.html
- **GitHub Repository:** [Link to repository]
- **Issues:** Report at GitHub Issues

---

## Changelog

### Version 1.0.0 (Current)
- ✅ Complete CRUD for Units, Tenants, Leases, Payments
- ✅ Authentication with JWT
- ✅ Dashboard with consolidated metrics
- ✅ Financial and payment reports
- ✅ Automatic payment schedule generation
- ✅ Lease renewal system

### Future Enhancements
- [ ] SMS notifications (Twilio integration)
- [ ] Report export (PDF, CSV, Excel)
- [ ] Advanced filtering and sorting
- [ ] Bulk operations
- [ ] Payment receipts generation
- [ ] Contract PDF generation

---

## Quick Reference Card

```typescript
// Authentication
POST   /auth/login                    → Get JWT token
GET    /auth/me                       → Get current user

// Units
GET    /units                         → List all units
POST   /units                         → Create unit
GET    /units/:id                     → Get unit
PUT    /units/:id                     → Update unit
DELETE /units/:id                     → Delete unit

// Tenants
GET    /tenants                       → List all tenants
POST   /tenants                       → Create tenant
GET    /tenants/:id                   → Get tenant
PUT    /tenants/:id                   → Update tenant

// Leases
GET    /leases                        → List all leases
POST   /leases                        → Create lease + payments
POST   /leases/:id/renew              → Renew lease
POST   /leases/:id/cancel             → Cancel lease

// Payments
GET    /leases/:id/payments           → Get lease payments
PUT    /payments/:id/pay              → Mark as paid
GET    /payments/overdue              → Get overdue
GET    /payments/upcoming?days=7      → Get upcoming

// Dashboard
GET    /dashboard                     → All metrics + alerts

// Reports
GET    /reports/financial?start_date=...&end_date=...
GET    /reports/payments?lease_id=...
```
