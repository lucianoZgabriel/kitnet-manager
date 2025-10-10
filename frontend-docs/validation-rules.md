# Validation Rules & Business Logic

This document consolidates all validation rules and business logic for the Kitnet Manager API.

---

## Units

### Field Validations

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `number` | string | ✅ | Cannot be empty, must be unique |
| `floor` | number | ✅ | Must be >= 1 |
| `base_rent_value` | decimal | ✅ | Must be > 0 |
| `renovated_rent_value` | decimal | ✅ | Must be >= base_rent_value |
| `status` | enum | ✅ | One of: `available`, `occupied`, `maintenance`, `renovation` |
| `is_renovated` | boolean | ✅ | - |
| `notes` | string | ❌ | Optional |

### Business Rules

1. **Current Rent Calculation**
   - If `is_renovated = true` → `current_rent_value = renovated_rent_value`
   - If `is_renovated = false` → `current_rent_value = base_rent_value`

2. **Status Changes**
   - Unit must be `available` to create a new lease
   - When lease is created → status automatically changes to `occupied`
   - When lease is cancelled → status automatically changes to `available`

3. **Deletion**
   - Cannot delete unit with `status = occupied`
   - Unit can only be deleted if no active leases exist

---

## Tenants

### Field Validations

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `full_name` | string | ✅ | Cannot be empty, trimmed |
| `cpf` | string | ✅ | Format: `XXX.XXX.XXX-XX`, must be unique, exactly 11 digits |
| `phone` | string | ✅ | Cannot be empty |
| `email` | string | ❌ | If provided, must be valid email format |
| `id_document_type` | string | ❌ | Optional |
| `id_document_number` | string | ❌ | Optional |

### Business Rules

1. **CPF Validation**
   - **Format:** `XXX.XXX.XXX-XX` (with dots and hyphen)
   - **Uniqueness:** Cannot register same CPF twice
   - **Immutability:** CPF cannot be changed after creation

2. **Phone Formatting**
   - Accepts various formats: `(XX) XXXXX-XXXX`, `(XX) XXXX-XXXX`
   - System stores as formatted string

3. **Email Validation**
   - Regex: `/^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/`
   - Optional but must be valid if provided

4. **Deletion**
   - Cannot delete tenant with active lease
   - System checks for active contracts before deletion

### CPF Format Examples

```typescript
// ✅ Valid
"123.456.789-00"
"987.654.321-99"

// ❌ Invalid
"12345678900"         // Missing formatting
"123.456.789"         // Incomplete
"123.456.789-0"       // Wrong length
"abc.def.ghi-jk"      // Non-numeric
```

---

## Leases

### Field Validations

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `unit_id` | UUID | ✅ | Must exist and unit must be `available` |
| `tenant_id` | UUID | ✅ | Must exist and tenant cannot have active lease |
| `contract_signed_date` | date | ✅ | ISO format |
| `start_date` | date | ✅ | Must be after signed date |
| `end_date` | date | Auto | Calculated as start_date + 6 months |
| `payment_due_day` | number | ✅ | Must be between 1-31 |
| `monthly_rent_value` | decimal | ✅ | Must be > 0 |
| `painting_fee_total` | decimal | ✅ | Must be >= 0 |
| `painting_fee_installments` | number | ✅ | Must be 1, 2, 3, or 4 |

### Business Rules

1. **Contract Duration**
   - Fixed at **6 months**
   - `end_date` is automatically calculated
   - Cannot be changed manually

2. **Unit Availability**
   - Unit must have status `available`
   - System automatically changes unit to `occupied` on lease creation
   - One unit can only have one active lease at a time

3. **Tenant Eligibility**
   - Tenant cannot have multiple active leases simultaneously
   - System checks for existing active leases before creation

4. **Payment Schedule Generation**
   - System automatically generates payment records for:
     - **Monthly rent**: 6 payments (one per month)
     - **Painting fee**: Divided into specified installments
   - Due date: `payment_due_day` of each month
   - First payment due on first occurrence of due day after start_date

5. **Painting Fee Installments**
   - Can be divided into 1, 2, 3, or 4 installments
   - Each installment = `painting_fee_total / painting_fee_installments`
   - Distributed across first months of the lease

6. **Lease Status Transitions**
   ```
   active → expiring_soon (45 days before end_date)
   active → expired (after end_date)
   active → cancelled (manual)

   expiring_soon → expired (after end_date)
   expiring_soon → cancelled (manual)
   expiring_soon → renewed (creates new lease)
   ```

7. **Renewal**
   - Can only renew leases with status `active` or `expiring_soon`
   - New lease starts on old lease's `end_date`
   - Same unit and tenant
   - Can have different painting fee
   - Old lease marked as `expired`

8. **Cancellation**
   - Lease status → `cancelled`
   - Unit status → `available`
   - Pending payments → `cancelled`

---

## Payments

### Field Validations

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `lease_id` | UUID | ✅ | Must exist |
| `payment_type` | enum | ✅ | One of: `rent`, `painting_fee`, `adjustment` |
| `reference_month` | date | ✅ | ISO format (YYYY-MM-DD) |
| `amount` | decimal | ✅ | Must be > 0 |
| `status` | enum | ✅ | One of: `pending`, `paid`, `overdue`, `cancelled` |
| `due_date` | date | ✅ | ISO format |
| `payment_date` | date | ❌ | Required when marking as paid |
| `payment_method` | enum | ❌ | Required when marking as paid. One of: `pix`, `cash`, `bank_transfer`, `credit_card` |

### Business Rules

1. **Automatic Payment Generation**
   - Payments are created automatically when lease is created
   - Cannot manually create payments (except adjustments)

2. **Payment Status Transitions**
   ```
   pending → paid (manual)
   pending → overdue (automatic after due_date)
   pending → cancelled (manual)

   overdue → paid (manual)
   overdue → cancelled (manual)
   ```

3. **Overdue Detection**
   - System automatically changes `pending` → `overdue` after `due_date`
   - Runs as background job or on-demand check

4. **Late Fees Calculation**
   ```typescript
   const calculateLateFees = (amount: number, daysOverdue: number) => {
     const penalty = amount * 0.02              // 2% flat penalty
     const monthlyInterest = amount * 0.01       // 1% monthly interest
     const dailyInterest = monthlyInterest / 30  // Pro-rata daily
     const interest = dailyInterest * daysOverdue

     return {
       penalty,
       interest,
       total: penalty + interest
     }
   }
   ```

5. **Painting Fee Payment**
   - When painting fee payment is marked as paid:
     - Updates `painting_fee_paid` in the lease
     - Validates that total paid doesn't exceed `painting_fee_total`

6. **Cancellation Rules**
   - Can only cancel `pending` or `overdue` payments
   - Cannot cancel `paid` payments

---

## Authentication

### Field Validations

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `username` | string | ✅ | Min 3 characters, lowercase, trimmed, unique |
| `password` | string | ✅ | Min 6 characters |
| `role` | enum | ✅ | One of: `admin`, `manager`, `viewer` |

### Business Rules

1. **Password Security**
   - Minimum 6 characters
   - Hashed using bcrypt (cost 10)
   - Never returned in API responses

2. **Roles & Permissions**
   ```typescript
   const permissions = {
     admin: {
       read: true,
       write: true,
       manage_users: true
     },
     manager: {
       read: true,
       write: true,
       manage_users: false
     },
     viewer: {
       read: true,
       write: false,
       manage_users: false
     }
   }
   ```

3. **User Status**
   - Active users: `is_active = true`
   - Deactivated users cannot login
   - Admin can activate/deactivate users

4. **Token Expiry**
   - JWT tokens have configurable expiry (default: 24 hours)
   - Use refresh endpoint to get new token

---

## General Validation Patterns

### UUIDs
All entity IDs use UUID v4 format:
```typescript
const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i
```

### Dates
ISO 8601 format: `YYYY-MM-DD` or `YYYY-MM-DDTHH:MM:SSZ`
```typescript
const DATE_REGEX = /^\d{4}-\d{2}-\d{2}$/
const DATETIME_REGEX = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$/
```

### Decimals
Monetary values as strings with 2 decimal places:
```typescript
const DECIMAL_REGEX = /^\d+\.\d{2}$/
// Examples: "1000.00", "123.45"
```

---

## Error Response Format

All validation errors follow this structure:

```json
{
  "success": false,
  "error": "Validation error message",
  "data": null
}
```

### HTTP Status Codes

| Code | Meaning | Usage |
|------|---------|-------|
| 200 | OK | Successful GET, PUT, PATCH, DELETE |
| 201 | Created | Successful POST |
| 400 | Bad Request | Validation error, business rule violation |
| 401 | Unauthorized | Missing/invalid token, invalid credentials |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Unique constraint violation (CPF, username) |
| 500 | Internal Server Error | Server error |

---

## Frontend Validation Examples

### CPF Validation

```typescript
const validateCPF = (cpf: string): boolean => {
  const cpfRegex = /^\d{3}\.\d{3}\.\d{3}-\d{2}$/
  if (!cpfRegex.test(cpf)) return false

  // Extract only digits
  const digits = cpf.replace(/[.\-]/g, '')
  return digits.length === 11
}

const formatCPF = (value: string): string => {
  const digits = value.replace(/\D/g, '')
  return digits
    .replace(/(\d{3})(\d)/, '$1.$2')
    .replace(/(\d{3})(\d)/, '$1.$2')
    .replace(/(\d{3})(\d{1,2})$/, '$1-$2')
}
```

### Email Validation

```typescript
const validateEmail = (email: string): boolean => {
  const emailRegex = /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/
  return emailRegex.test(email)
}
```

### Date Validation

```typescript
const validateDate = (dateStr: string): boolean => {
  const date = new Date(dateStr)
  return date instanceof Date && !isNaN(date.getTime())
}

const formatDateToISO = (date: Date): string => {
  return date.toISOString().split('T')[0] // Returns YYYY-MM-DD
}
```

### Decimal Validation

```typescript
const validateDecimal = (value: string): boolean => {
  const decimalRegex = /^\d+\.\d{2}$/
  return decimalRegex.test(value) && parseFloat(value) >= 0
}

const formatDecimal = (value: number): string => {
  return value.toFixed(2)
}
```

---

## Rate Limiting

Not currently implemented, but recommended limits for production:

```typescript
const RATE_LIMITS = {
  login: '5 requests per 15 minutes',
  api: '100 requests per minute',
  reports: '10 requests per minute'
}
```
