# Units Endpoints

All endpoints require authentication via Bearer token.

Base URL: `https://kitnet-manager-production.up.railway.app/api/v1`

## Create Unit

```typescript
POST /units
```

**Request Body:**
```json
{
  "number": "101",
  "floor": 1,
  "base_rent_value": "800.00",
  "renovated_rent_value": "1000.00"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Unit created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "number": "101",
    "floor": 1,
    "status": "available",
    "is_renovated": false,
    "base_rent_value": "800.00",
    "renovated_rent_value": "1000.00",
    "current_rent_value": "800.00",
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

**Validation Rules:**
- `number`: required, cannot be empty
- `floor`: must be >= 1
- `base_rent_value`: must be > 0
- `renovated_rent_value`: must be >= base_rent_value

**Errors:**
- `400` - Invalid request body or validation error
- `500` - Internal server error

---

## Get Unit by ID

```typescript
GET /units/:id
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Unit retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "number": "101",
    "floor": 1,
    "status": "occupied",
    "is_renovated": true,
    "base_rent_value": "800.00",
    "renovated_rent_value": "1000.00",
    "current_rent_value": "1000.00",
    "notes": "Recently renovated",
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

**Errors:**
- `400` - Invalid unit ID
- `404` - Unit not found

---

## List Units

```typescript
GET /units
GET /units?status=available
GET /units?floor=1
```

**Query Parameters:**
- `status` (optional): Filter by status (`available`, `occupied`, `maintenance`, `renovation`)
- `floor` (optional): Filter by floor number

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Units retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "number": "101",
      "floor": 1,
      "status": "available",
      "is_renovated": false,
      "base_rent_value": "800.00",
      "renovated_rent_value": "1000.00",
      "current_rent_value": "800.00",
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

## Update Unit

```typescript
PUT /units/:id
```

**Request Body:**
```json
{
  "number": "101",
  "floor": 1,
  "is_renovated": true,
  "base_rent_value": "800.00",
  "renovated_rent_value": "1000.00",
  "notes": "Recently renovated with new paint"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Unit updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "number": "101",
    "floor": 1,
    "status": "available",
    "is_renovated": true,
    "base_rent_value": "800.00",
    "renovated_rent_value": "1000.00",
    "current_rent_value": "1000.00",
    "notes": "Recently renovated with new paint",
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T12:00:00Z"
  }
}
```

**Errors:**
- `400` - Invalid request body or validation error
- `404` - Unit not found

---

## Update Unit Status

```typescript
PATCH /units/:id/status
```

**Request Body:**
```json
{
  "status": "maintenance"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Unit status updated successfully",
  "data": null
}
```

**Valid Status Values:**
- `available`
- `occupied`
- `maintenance`
- `renovation`

**Errors:**
- `400` - Invalid status value
- `404` - Unit not found

---

## Delete Unit

```typescript
DELETE /units/:id
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Unit deleted successfully",
  "data": null
}
```

**Errors:**
- `400` - Cannot delete occupied unit
- `404` - Unit not found

---

## Get Occupancy Statistics

```typescript
GET /units/stats/occupancy
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Occupancy stats retrieved successfully",
  "data": {
    "total_units": 31,
    "available_units": 5,
    "occupied_units": 23,
    "maintenance_units": 2,
    "renovation_units": 1,
    "occupancy_rate": 74.19
  }
}
```

**Errors:**
- `500` - Internal server error
