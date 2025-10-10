# Payments Endpoints

All endpoints require authentication via Bearer token.

Base URL: `https://kitnet-manager-production.up.railway.app/api/v1`

## Get Payment by ID

```typescript
GET /payments/:id
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Payment retrieved successfully",
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440003",
    "lease_id": "770e8400-e29b-41d4-a716-446655440002",
    "payment_type": "rent",
    "reference_month": "2025-02-01",
    "amount": "1000.00",
    "status": "paid",
    "due_date": "2025-02-05",
    "payment_date": "2025-02-03",
    "payment_method": "pix",
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-02-03T14:30:00Z"
  }
}
```

**Errors:**
- `400` - Invalid payment ID
- `404` - Payment not found

---

## Get Payments by Lease

```typescript
GET /leases/:lease_id/payments
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Payments retrieved successfully",
  "data": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "lease_id": "770e8400-e29b-41d4-a716-446655440002",
      "payment_type": "rent",
      "reference_month": "2025-02-01",
      "amount": "1000.00",
      "status": "paid",
      "due_date": "2025-02-05",
      "payment_date": "2025-02-03",
      "payment_method": "pix",
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-02-03T14:30:00Z"
    },
    {
      "id": "990e8400-e29b-41d4-a716-446655440004",
      "lease_id": "770e8400-e29b-41d4-a716-446655440002",
      "payment_type": "painting_fee",
      "reference_month": "2025-02-01",
      "amount": "125.00",
      "status": "pending",
      "due_date": "2025-02-05",
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-15T10:00:00Z"
    }
  ]
}
```

**Errors:**
- `400` - Invalid lease ID
- `404` - Lease not found

---

## Get Overdue Payments

```typescript
GET /payments/overdue
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Overdue payments retrieved successfully",
  "data": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "lease_id": "770e8400-e29b-41d4-a716-446655440002",
      "payment_type": "rent",
      "reference_month": "2025-01-01",
      "amount": "1000.00",
      "status": "overdue",
      "due_date": "2025-01-05",
      "created_at": "2024-12-15T10:00:00Z",
      "updated_at": "2025-01-06T00:00:00Z"
    }
  ]
}
```

**Note:** Payments are automatically marked as `overdue` after the due_date passes

---

## Get Upcoming Payments

```typescript
GET /payments/upcoming
GET /payments/upcoming?days=14
```

**Query Parameters:**
- `days` (optional): Number of days ahead to look (default: 7)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Upcoming payments retrieved successfully",
  "data": [
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
```

**Errors:**
- `400` - Invalid days parameter
- `500` - Internal server error

---

## Mark Payment as Paid

```typescript
PUT /payments/:id/pay
```

**Request Body:**
```json
{
  "payment_date": "2025-02-03",
  "payment_method": "pix"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Payment marked as paid successfully",
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440003",
    "lease_id": "770e8400-e29b-41d4-a716-446655440002",
    "payment_type": "rent",
    "reference_month": "2025-02-01",
    "amount": "1000.00",
    "status": "paid",
    "due_date": "2025-02-05",
    "payment_date": "2025-02-03",
    "payment_method": "pix",
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-02-03T14:30:00Z"
  }
}
```

**Validation:**
- `payment_date`: required, ISO date format (YYYY-MM-DD)
- `payment_method`: required, must be one of: `pix`, `cash`, `bank_transfer`, `credit_card`
- Payment status must be `pending` or `overdue`

**Business Rules:**
- If payment type is `painting_fee`, automatically updates the `painting_fee_paid` in the lease

**Errors:**
- `400` - Invalid request, payment already paid, or cannot be paid
- `404` - Payment not found

---

## Cancel Payment

```typescript
POST /payments/:id/cancel
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Payment cancelled successfully",
  "data": null
}
```

**Business Rules:**
- Can only cancel payments with status `pending` or `overdue`
- Cannot cancel already paid payments

**Errors:**
- `400` - Payment already paid or cannot be cancelled
- `404` - Payment not found

---

## Get Payment Statistics by Lease

```typescript
GET /leases/:lease_id/payments/stats
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Payment stats retrieved successfully",
  "data": {
    "total_payments": 8,
    "total_amount": "6250.00",
    "paid_amount": "3125.00",
    "pending_amount": "2125.00",
    "overdue_amount": "1000.00",
    "paid_count": 3,
    "pending_count": 4,
    "overdue_count": 1
  }
}
```

**Errors:**
- `400` - Invalid lease ID
- `404` - Lease not found

---

## Payment Types

```typescript
type PaymentType = 'rent' | 'painting_fee' | 'adjustment'
```

- **rent**: Monthly rent payment (automatically generated)
- **painting_fee**: One-time painting fee (divided into installments)
- **adjustment**: Manual adjustment (credits/debits)

---

## Payment Status

```typescript
type PaymentStatus = 'pending' | 'paid' | 'overdue' | 'cancelled'
```

**Status Transitions:**
```
pending → paid (manual mark as paid)
pending → overdue (automatic after due_date)
pending → cancelled (manual cancellation)
overdue → paid (manual mark as paid)
overdue → cancelled (manual cancellation)
```

---

## Payment Methods

```typescript
type PaymentMethod = 'pix' | 'cash' | 'bank_transfer' | 'credit_card'
```

---

## Late Fees (Business Rule)

Late fees are calculated but not automatically added to the payment amount. The calculation is:

- **2% penalty** on the payment amount
- **1% monthly interest** (pro-rata daily)

Example:
```typescript
// Payment amount: R$1000.00
// Days overdue: 10 days

const penalty = amount * 0.02 // R$20.00
const dailyInterest = (amount * 0.01) / 30 // R$0.33 per day
const interest = dailyInterest * daysOverdue // R$3.30
const totalWithFees = amount + penalty + interest // R$1023.30
```

**Note:** The frontend should display late fees separately for transparency
