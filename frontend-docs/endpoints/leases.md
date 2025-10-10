# Leases Endpoints

All endpoints require authentication via Bearer token.

Base URL: `https://kitnet-manager-production.up.railway.app/api/v1`

## Create Lease

```typescript
POST /leases
```

**Request Body:**
```json
{
  "unit_id": "550e8400-e29b-41d4-a716-446655440000",
  "tenant_id": "660e8400-e29b-41d4-a716-446655440001",
  "contract_signed_date": "2025-01-15",
  "start_date": "2025-02-01",
  "payment_due_day": 5,
  "monthly_rent_value": "1000.00",
  "painting_fee_total": "250.00",
  "painting_fee_installments": 2
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Lease created successfully with payments",
  "data": {
    "lease": {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "unit_id": "550e8400-e29b-41d4-a716-446655440000",
      "tenant_id": "660e8400-e29b-41d4-a716-446655440001",
      "contract_signed_date": "2025-01-15",
      "start_date": "2025-02-01",
      "end_date": "2025-08-01",
      "payment_due_day": 5,
      "monthly_rent_value": "1000.00",
      "painting_fee_total": "250.00",
      "painting_fee_installments": 2,
      "painting_fee_paid": "0.00",
      "status": "active",
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-15T10:00:00Z"
    },
    "payments": [
      {
        "id": "880e8400-e29b-41d4-a716-446655440003",
        "lease_id": "770e8400-e29b-41d4-a716-446655440002",
        "payment_type": "rent",
        "reference_month": "2025-02-01",
        "amount": "1000.00",
        "status": "pending",
        "due_date": "2025-02-05",
        "created_at": "2025-01-15T10:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
      }
    ]
  }
}
```

**Validation Rules:**
- `unit_id`: must exist and be available
- `tenant_id`: must exist and not have active lease
- `contract_signed_date`: required
- `start_date`: required, must be after signed date
- `end_date`: automatically calculated (start_date + 6 months)
- `payment_due_day`: must be between 1-31
- `monthly_rent_value`: must be > 0
- `painting_fee_total`: must be >= 0
- `painting_fee_installments`: must be 1, 2, 3, or 4

**Business Rules:**
- Contract duration is fixed at 6 months
- Unit status automatically changes to "occupied"
- Payment schedules are automatically generated for the entire lease period
- Painting fee is divided into equal installments

**Errors:**
- `400` - Invalid request body, validation error, unit not available, or tenant already has active lease
- `404` - Unit or tenant not found
- `500` - Internal server error

---

## Get Lease by ID

```typescript
GET /leases/:id
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Lease retrieved successfully",
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "unit_id": "550e8400-e29b-41d4-a716-446655440000",
    "tenant_id": "660e8400-e29b-41d4-a716-446655440001",
    "contract_signed_date": "2025-01-15",
    "start_date": "2025-02-01",
    "end_date": "2025-08-01",
    "payment_due_day": 5,
    "monthly_rent_value": "1000.00",
    "painting_fee_total": "250.00",
    "painting_fee_installments": 2,
    "painting_fee_paid": "125.00",
    "status": "active",
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

**Errors:**
- `400` - Invalid lease ID
- `404` - Lease not found

---

## List Leases

```typescript
GET /leases
GET /leases?status=active
GET /leases?unit_id=550e8400-e29b-41d4-a716-446655440000
GET /leases?tenant_id=660e8400-e29b-41d4-a716-446655440001
```

**Query Parameters:**
- `status` (optional): Filter by status (`active`, `expiring_soon`, `expired`, `cancelled`)
- `unit_id` (optional): Filter by unit ID (UUID)
- `tenant_id` (optional): Filter by tenant ID (UUID)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Leases retrieved successfully",
  "data": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "unit_id": "550e8400-e29b-41d4-a716-446655440000",
      "tenant_id": "660e8400-e29b-41d4-a716-446655440001",
      "status": "active",
      "start_date": "2025-02-01",
      "end_date": "2025-08-01",
      "monthly_rent_value": "1000.00",
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-15T10:00:00Z"
    }
  ]
}
```

**Errors:**
- `400` - Invalid query parameters
- `500` - Internal server error

---

## Get Lease Statistics

```typescript
GET /leases/stats
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Lease stats retrieved successfully",
  "data": {
    "total_leases": 23,
    "active_leases": 20,
    "expiring_soon_leases": 3,
    "expired_leases": 5,
    "cancelled_leases": 2
  }
}
```

---

## Renew Lease

```typescript
POST /leases/:id/renew
```

**Request Body:**
```json
{
  "painting_fee_total": "250.00",
  "painting_fee_installments": 1
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Lease renewed successfully",
  "data": {
    "id": "990e8400-e29b-41d4-a716-446655440004",
    "unit_id": "550e8400-e29b-41d4-a716-446655440000",
    "tenant_id": "660e8400-e29b-41d4-a716-446655440001",
    "contract_signed_date": "2025-07-25",
    "start_date": "2025-08-01",
    "end_date": "2026-02-01",
    "payment_due_day": 5,
    "monthly_rent_value": "1000.00",
    "painting_fee_total": "250.00",
    "painting_fee_installments": 1,
    "painting_fee_paid": "0.00",
    "status": "active",
    "created_at": "2025-07-25T10:00:00Z",
    "updated_at": "2025-07-25T10:00:00Z"
  }
}
```

**Business Rules:**
- Can only renew leases with status `active` or `expiring_soon`
- Old lease is marked as `expired`
- New lease starts immediately after old lease ends
- New payment schedules are automatically generated

**Errors:**
- `400` - Lease cannot be renewed (wrong status or already expired)
- `404` - Lease not found

---

## Cancel Lease

```typescript
POST /leases/:id/cancel
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Lease cancelled successfully",
  "data": null
}
```

**Business Rules:**
- Lease status changes to `cancelled`
- Unit status changes to `available`
- Pending payments are cancelled

**Errors:**
- `400` - Lease cannot be cancelled
- `404` - Lease not found

---

## Update Painting Fee Paid

```typescript
PATCH /leases/:id/painting-fee
```

**Request Body:**
```json
{
  "amount_paid": "125.00"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Painting fee updated successfully",
  "data": null
}
```

**Validation:**
- `amount_paid`: must be > 0
- Total paid cannot exceed `painting_fee_total`

**Errors:**
- `400` - Invalid amount or exceeds total
- `404` - Lease not found

---

## Get Expiring Soon Leases

```typescript
GET /leases/expiring-soon
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Expiring soon leases retrieved successfully",
  "data": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "unit_id": "550e8400-e29b-41d4-a716-446655440000",
      "tenant_id": "660e8400-e29b-41d4-a716-446655440001",
      "status": "expiring_soon",
      "start_date": "2025-02-01",
      "end_date": "2025-08-01",
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-07-15T10:00:00Z"
    }
  ]
}
```

**Note:** Returns leases expiring in the next 45 days

---

## Lease Status Transitions

```
active → expiring_soon (45 days before end_date)
active → cancelled (manual cancellation)
active → expired (after end_date)
expiring_soon → active (renewed)
expiring_soon → cancelled (manual cancellation)
expiring_soon → expired (after end_date)
```
