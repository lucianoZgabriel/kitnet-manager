# Dashboard & Reports Endpoints

All endpoints require authentication via Bearer token.

Base URL: `https://kitnet-manager-production.up.railway.app/api/v1`

## Get Dashboard Metrics

```typescript
GET /dashboard
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Dashboard data retrieved successfully",
  "data": {
    "occupancy": {
      "total_units": 31,
      "available_units": 5,
      "occupied_units": 23,
      "maintenance_units": 2,
      "renovation_units": 1,
      "occupancy_rate": 74.19
    },
    "financial": {
      "monthly_revenue": "23000.00",
      "total_receivable": "138000.00",
      "total_received": "115000.00",
      "total_pending": "18000.00",
      "total_overdue": "5000.00",
      "pending_count": 12,
      "overdue_count": 3
    },
    "contracts": {
      "total_active_leases": 23,
      "expiring_soon_count": 3,
      "expired_count": 0
    },
    "alerts": [
      {
        "type": "contract_expiring",
        "severity": "high",
        "message": "Lease for Unit 101 expires in 15 days",
        "entity_id": "770e8400-e29b-41d4-a716-446655440002",
        "entity_type": "lease"
      },
      {
        "type": "payment_overdue",
        "severity": "high",
        "message": "Payment overdue for Unit 205 (10 days)",
        "entity_id": "880e8400-e29b-41d4-a716-446655440003",
        "entity_type": "payment"
      },
      {
        "type": "unit_maintenance",
        "severity": "medium",
        "message": "Unit 310 is under maintenance",
        "entity_id": "550e8400-e29b-41d4-a716-446655440000",
        "entity_type": "unit"
      }
    ]
  }
}
```

**Metrics Breakdown:**

### Occupancy Metrics
- `total_units`: Total number of units in the system
- `available_units`: Units available for rent
- `occupied_units`: Units currently rented
- `maintenance_units`: Units under maintenance
- `renovation_units`: Units being renovated
- `occupancy_rate`: Percentage of occupied units (occupied / total * 100)

### Financial Metrics
- `monthly_revenue`: Expected monthly revenue from all active leases
- `total_receivable`: Total amount to be received (all pending + overdue)
- `total_received`: Total amount already paid
- `total_pending`: Sum of pending payments
- `total_overdue`: Sum of overdue payments
- `pending_count`: Number of pending payments
- `overdue_count`: Number of overdue payments

### Contract Metrics
- `total_active_leases`: Number of active leases
- `expiring_soon_count`: Leases expiring in next 45 days
- `expired_count`: Recently expired leases

### Alerts
- `type`: Type of alert (`contract_expiring`, `payment_overdue`, `unit_maintenance`)
- `severity`: Alert severity (`high`, `medium`, `low`)
- `message`: Human-readable alert message
- `entity_id`: ID of the related entity
- `entity_type`: Type of entity (`lease`, `payment`, `unit`)

**Errors:**
- `500` - Internal server error

---

## Get Financial Report

```typescript
GET /reports/financial?start_date=2025-01-01&end_date=2025-01-31
GET /reports/financial?start_date=2025-01-01&end_date=2025-01-31&payment_type=rent
GET /reports/financial?start_date=2025-01-01&end_date=2025-01-31&status=paid
```

**Query Parameters:**
- `start_date` (required): Start date in format YYYY-MM-DD
- `end_date` (required): End date in format YYYY-MM-DD
- `payment_type` (optional): Filter by type (`rent`, `painting_fee`, `adjustment`)
- `status` (optional): Filter by status (`pending`, `paid`, `overdue`, `cancelled`)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Financial report generated successfully",
  "data": {
    "period": {
      "start_date": "2025-01-01",
      "end_date": "2025-01-31"
    },
    "summary": {
      "total_amount": "25000.00",
      "paid_amount": "20000.00",
      "pending_amount": "3000.00",
      "overdue_amount": "2000.00",
      "payment_count": 25
    },
    "by_type": {
      "rent": {
        "count": 20,
        "total_amount": "20000.00"
      },
      "painting_fee": {
        "count": 5,
        "total_amount": "5000.00"
      },
      "adjustment": {
        "count": 0,
        "total_amount": "0.00"
      }
    },
    "by_status": {
      "paid": {
        "count": 20,
        "total_amount": "20000.00"
      },
      "pending": {
        "count": 3,
        "total_amount": "3000.00"
      },
      "overdue": {
        "count": 2,
        "total_amount": "2000.00"
      },
      "cancelled": {
        "count": 0,
        "total_amount": "0.00"
      }
    },
    "payments": [
      {
        "id": "880e8400-e29b-41d4-a716-446655440003",
        "lease_id": "770e8400-e29b-41d4-a716-446655440002",
        "payment_type": "rent",
        "reference_month": "2025-01-01",
        "amount": "1000.00",
        "status": "paid",
        "due_date": "2025-01-05",
        "payment_date": "2025-01-03",
        "payment_method": "pix",
        "created_at": "2024-12-15T10:00:00Z",
        "updated_at": "2025-01-03T14:30:00Z"
      }
    ]
  }
}
```

**Use Cases:**
- Monthly financial closing
- Revenue analysis by payment type
- Identify payment trends
- Export data for accounting

**Errors:**
- `400` - Invalid date format or missing required parameters
- `500` - Internal server error

---

## Get Payment History Report

```typescript
GET /reports/payments
GET /reports/payments?lease_id=770e8400-e29b-41d4-a716-446655440002
GET /reports/payments?tenant_id=660e8400-e29b-41d4-a716-446655440001
GET /reports/payments?status=paid&start_date=2025-01-01&end_date=2025-01-31
```

**Query Parameters:**
- `lease_id` (optional): Filter by lease ID (UUID)
- `tenant_id` (optional): Filter by tenant ID (UUID)
- `status` (optional): Filter by status (`pending`, `paid`, `overdue`, `cancelled`)
- `start_date` (optional): Filter by date range (YYYY-MM-DD)
- `end_date` (optional): Filter by date range (YYYY-MM-DD)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Payment history report generated successfully",
  "data": {
    "total_count": 25,
    "payments": [
      {
        "id": "880e8400-e29b-41d4-a716-446655440003",
        "lease_id": "770e8400-e29b-41d4-a716-446655440002",
        "payment_type": "rent",
        "reference_month": "2025-01-01",
        "amount": "1000.00",
        "status": "paid",
        "due_date": "2025-01-05",
        "payment_date": "2025-01-03",
        "payment_method": "pix",
        "created_at": "2024-12-15T10:00:00Z",
        "updated_at": "2025-01-03T14:30:00Z"
      }
    ]
  }
}
```

**Use Cases:**
- Tenant payment history
- Lease payment tracking
- Payment compliance verification
- Audit trail

**Errors:**
- `400` - Invalid query parameters
- `500` - Internal server error

---

## Dashboard Best Practices

### Refresh Intervals

Recommended refresh intervals for real-time dashboard:

```typescript
const REFRESH_INTERVALS = {
  dashboard: 60000,      // 1 minute
  alerts: 30000,         // 30 seconds
  financial: 300000,     // 5 minutes
}
```

### Caching Strategy

For performance, implement client-side caching:

```typescript
// Example with React Query
const { data } = useQuery({
  queryKey: ['dashboard'],
  queryFn: fetchDashboard,
  refetchInterval: 60000,
  staleTime: 30000,
})
```

### Error Handling

Always handle partial failures gracefully:

```typescript
try {
  const dashboard = await fetchDashboard()
  // Use dashboard data
} catch (error) {
  // Fallback to cached data or show partial UI
  console.error('Dashboard fetch failed:', error)
}
```

### Alert Prioritization

Display alerts in this order:
1. `high` severity first
2. Sort by type: `payment_overdue` > `contract_expiring` > `unit_maintenance`
3. Most recent first

---

## Export Options (Future Enhancement)

The reports endpoints can be extended to support export formats:

```typescript
GET /reports/financial?start_date=2025-01-01&end_date=2025-01-31&format=csv
GET /reports/financial?start_date=2025-01-01&end_date=2025-01-31&format=pdf
GET /reports/financial?start_date=2025-01-01&end_date=2025-01-31&format=xlsx
```

**Note:** Export functionality is not yet implemented in the current version
