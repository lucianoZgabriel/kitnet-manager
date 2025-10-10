# Authentication Endpoints

Base URL: `https://kitnet-manager-production.up.railway.app/api/v1`

## Login

```typescript
POST /auth/login
```

**Request Body:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "admin",
      "role": "admin",
      "is_active": true,
      "created_at": "2025-01-01T10:00:00Z",
      "updated_at": "2025-01-01T10:00:00Z"
    }
  }
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Invalid credentials or inactive user
- `500` - Internal server error

---

## Get Current User

```typescript
GET /auth/me
```

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "admin",
    "role": "admin",
    "is_active": true,
    "last_login_at": "2025-01-15T10:00:00Z",
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

**Errors:**
- `401` - Missing or invalid token
- `500` - Internal server error

---

## Refresh Token

```typescript
POST /auth/refresh
```

**Request Body:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Invalid or expired token
- `500` - Internal server error

---

## Create User (Admin Only)

```typescript
POST /auth/users
```

**Headers:**
```
Authorization: Bearer {admin_token}
```

**Request Body:**
```json
{
  "username": "newuser",
  "password": "password123",
  "role": "manager"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "username": "newuser",
    "role": "manager",
    "is_active": true,
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

**Validation Rules:**
- `username`: min 3 characters
- `password`: min 6 characters
- `role`: must be "admin", "manager", or "viewer"

**Errors:**
- `400` - Invalid request body
- `409` - Username already exists
- `500` - Internal server error

---

## List Users (Admin Only)

```typescript
GET /auth/users
```

**Headers:**
```
Authorization: Bearer {admin_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "admin",
      "role": "admin",
      "is_active": true,
      "created_at": "2025-01-01T10:00:00Z",
      "updated_at": "2025-01-01T10:00:00Z"
    }
  ]
}
```

---

## Get User by ID (Admin Only)

```typescript
GET /auth/users/:id
```

**Headers:**
```
Authorization: Bearer {admin_token}
```

**Response (200 OK):** Same as List Users (single object)

**Errors:**
- `400` - Invalid user ID
- `404` - User not found

---

## Change Password

```typescript
POST /auth/change-password
```

**Headers:**
```
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "old_password": "oldpass123",
  "new_password": "newpass456"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password changed successfully",
  "data": null
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Invalid old password or missing token
- `500` - Internal server error

---

## Change User Role (Admin Only)

```typescript
PATCH /auth/users/:id/role
```

**Headers:**
```
Authorization: Bearer {admin_token}
```

**Request Body:**
```json
{
  "role": "manager"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User role changed successfully",
  "data": null
}
```

**Errors:**
- `400` - Invalid user ID or role
- `404` - User not found

---

## Deactivate User (Admin Only)

```typescript
POST /auth/users/:id/deactivate
```

**Headers:**
```
Authorization: Bearer {admin_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User deactivated successfully",
  "data": null
}
```

---

## Activate User (Admin Only)

```typescript
POST /auth/users/:id/activate
```

**Headers:**
```
Authorization: Bearer {admin_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User activated successfully",
  "data": null
}
```
