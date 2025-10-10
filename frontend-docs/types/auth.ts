// Authentication & User Types

export type UserRole = 'admin' | 'manager' | 'viewer'

export interface User {
  id: string
  username: string
  role: UserRole
  is_active: boolean
  last_login_at?: string
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface RefreshTokenRequest {
  token: string
}

export interface RefreshTokenResponse {
  token: string
}

export interface CreateUserRequest {
  username: string
  password: string
  role: UserRole
}

export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

export interface ChangeRoleRequest {
  role: UserRole
}
