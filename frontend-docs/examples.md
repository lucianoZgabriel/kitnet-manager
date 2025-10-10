# Code Examples - Frontend Integration

Exemplos práticos de código para integração com a API do Kitnet Manager.

---

## Setup Inicial

### 1. Configurar Cliente API (Axios)

```typescript
// lib/api.ts
import axios from 'axios'

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'https://kitnet-manager-production.up.railway.app/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Interceptor: Adicionar token em todas as requisições
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Interceptor: Tratar respostas e erros
api.interceptors.response.use(
  (response) => response.data, // Retorna apenas data.data
  (error) => {
    const { response } = error

    if (response?.status === 401) {
      // Token expirado - redirecionar para login
      localStorage.removeItem('auth_token')
      window.location.href = '/login'
    }

    return Promise.reject(response?.data || error)
  }
)

export default api
```

---

## Authentication

### Login Hook

```typescript
// hooks/useAuth.ts
import { useState } from 'react'
import api from '@/lib/api'
import type { LoginRequest, LoginResponse, User } from '@/types/api/auth'

export const useAuth = () => {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(false)

  const login = async (credentials: LoginRequest) => {
    setLoading(true)
    try {
      const response = await api.post<any, { data: LoginResponse }>('/auth/login', credentials)

      // Salvar token
      localStorage.setItem('auth_token', response.data.token)

      // Salvar usuário
      setUser(response.data.user)

      return response.data
    } catch (error: any) {
      throw new Error(error.error || 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  const logout = () => {
    localStorage.removeItem('auth_token')
    setUser(null)
    window.location.href = '/login'
  }

  const getCurrentUser = async () => {
    try {
      const response = await api.get<any, { data: User }>('/auth/me')
      setUser(response.data)
      return response.data
    } catch (error) {
      logout()
    }
  }

  return {
    user,
    loading,
    login,
    logout,
    getCurrentUser,
  }
}
```

### Login Page

```typescript
// app/login/page.tsx
'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'

export default function LoginPage() {
  const router = useRouter()
  const { login, loading } = useAuth()
  const [credentials, setCredentials] = useState({
    username: '',
    password: '',
  })
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    try {
      await login(credentials)
      router.push('/dashboard')
    } catch (err: any) {
      setError(err.message)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center">
      <form onSubmit={handleSubmit} className="w-full max-w-md space-y-4">
        <h1 className="text-2xl font-bold">Kitnet Manager</h1>

        {error && (
          <div className="bg-red-100 text-red-700 p-3 rounded">
            {error}
          </div>
        )}

        <input
          type="text"
          placeholder="Username"
          value={credentials.username}
          onChange={(e) => setCredentials({ ...credentials, username: e.target.value })}
          className="w-full px-4 py-2 border rounded"
          required
        />

        <input
          type="password"
          placeholder="Password"
          value={credentials.password}
          onChange={(e) => setCredentials({ ...credentials, password: e.target.value })}
          className="w-full px-4 py-2 border rounded"
          required
        />

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 disabled:opacity-50"
        >
          {loading ? 'Logging in...' : 'Login'}
        </button>
      </form>
    </div>
  )
}
```

---

## Units

### Units Service

```typescript
// services/units.ts
import api from '@/lib/api'
import type { Unit, CreateUnitRequest, UpdateUnitRequest, OccupancyStats } from '@/types/api/unit'

export const unitsService = {
  getAll: async (filters?: { status?: string; floor?: number }) => {
    const params = new URLSearchParams()
    if (filters?.status) params.append('status', filters.status)
    if (filters?.floor) params.append('floor', filters.floor.toString())

    const response = await api.get<any, { data: Unit[] }>(`/units?${params}`)
    return response.data
  },

  getById: async (id: string) => {
    const response = await api.get<any, { data: Unit }>(`/units/${id}`)
    return response.data
  },

  create: async (data: CreateUnitRequest) => {
    const response = await api.post<any, { data: Unit }>('/units', data)
    return response.data
  },

  update: async (id: string, data: UpdateUnitRequest) => {
    const response = await api.put<any, { data: Unit }>(`/units/${id}`, data)
    return response.data
  },

  updateStatus: async (id: string, status: string) => {
    await api.patch(`/units/${id}/status`, { status })
  },

  delete: async (id: string) => {
    await api.delete(`/units/${id}`)
  },

  getOccupancyStats: async () => {
    const response = await api.get<any, { data: OccupancyStats }>('/units/stats/occupancy')
    return response.data
  },
}
```

### Units List with React Query

```typescript
// app/units/page.tsx
'use client'

import { useQuery } from '@tanstack/react-query'
import { unitsService } from '@/services/units'
import { useState } from 'react'

export default function UnitsPage() {
  const [statusFilter, setStatusFilter] = useState<string>('')

  const { data: units, isLoading, error } = useQuery({
    queryKey: ['units', statusFilter],
    queryFn: () => unitsService.getAll(statusFilter ? { status: statusFilter } : undefined),
    refetchInterval: 60000, // Refetch every minute
  })

  if (isLoading) return <div>Loading...</div>
  if (error) return <div>Error loading units</div>

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Units</h1>

        <select
          value={statusFilter}
          onChange={(e) => setStatusFilter(e.target.value)}
          className="px-4 py-2 border rounded"
        >
          <option value="">All Status</option>
          <option value="available">Available</option>
          <option value="occupied">Occupied</option>
          <option value="maintenance">Maintenance</option>
          <option value="renovation">Renovation</option>
        </select>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {units?.map((unit) => (
          <div key={unit.id} className="border rounded-lg p-4">
            <div className="flex justify-between items-start mb-2">
              <h3 className="text-lg font-semibold">Unit {unit.number}</h3>
              <span className={`px-2 py-1 text-xs rounded ${getStatusColor(unit.status)}`}>
                {unit.status}
              </span>
            </div>
            <p className="text-sm text-gray-600">Floor {unit.floor}</p>
            <p className="text-lg font-bold mt-2">
              R$ {parseFloat(unit.current_rent_value).toFixed(2)}
            </p>
            {unit.is_renovated && (
              <span className="inline-block mt-2 text-xs bg-green-100 text-green-800 px-2 py-1 rounded">
                Renovated
              </span>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}

function getStatusColor(status: string) {
  const colors = {
    available: 'bg-green-100 text-green-800',
    occupied: 'bg-blue-100 text-blue-800',
    maintenance: 'bg-yellow-100 text-yellow-800',
    renovation: 'bg-orange-100 text-orange-800',
  }
  return colors[status as keyof typeof colors] || 'bg-gray-100 text-gray-800'
}
```

---

## Tenants

### CPF Input Component

```typescript
// components/CPFInput.tsx
import { forwardRef } from 'react'

interface CPFInputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  error?: string
}

export const CPFInput = forwardRef<HTMLInputElement, CPFInputProps>(
  ({ error, onChange, value, ...props }, ref) => {
    const formatCPF = (value: string) => {
      const numbers = value.replace(/\D/g, '')
      return numbers
        .replace(/(\d{3})(\d)/, '$1.$2')
        .replace(/(\d{3})(\d)/, '$1.$2')
        .replace(/(\d{3})(\d{1,2})$/, '$1-$2')
        .slice(0, 14) // Limit to XXX.XXX.XXX-XX
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const formatted = formatCPF(e.target.value)
      e.target.value = formatted
      onChange?.(e)
    }

    return (
      <div>
        <input
          ref={ref}
          type="text"
          value={value}
          onChange={handleChange}
          placeholder="XXX.XXX.XXX-XX"
          {...props}
          className={`px-4 py-2 border rounded ${error ? 'border-red-500' : ''}`}
        />
        {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
      </div>
    )
  }
)
```

### Create Tenant Form

```typescript
// components/CreateTenantForm.tsx
'use client'

import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { tenantsService } from '@/services/tenants'
import { CPFInput } from '@/components/CPFInput'
import type { CreateTenantRequest } from '@/types/api/tenant'

export default function CreateTenantForm({ onSuccess }: { onSuccess?: () => void }) {
  const queryClient = useQueryClient()
  const [formData, setFormData] = useState<CreateTenantRequest>({
    full_name: '',
    cpf: '',
    phone: '',
    email: '',
  })

  const mutation = useMutation({
    mutationFn: tenantsService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tenants'] })
      onSuccess?.()
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    mutation.mutate(formData)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <input
        type="text"
        placeholder="Full Name"
        value={formData.full_name}
        onChange={(e) => setFormData({ ...formData, full_name: e.target.value })}
        className="w-full px-4 py-2 border rounded"
        required
      />

      <CPFInput
        value={formData.cpf}
        onChange={(e) => setFormData({ ...formData, cpf: e.target.value })}
        required
      />

      <input
        type="tel"
        placeholder="Phone"
        value={formData.phone}
        onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
        className="w-full px-4 py-2 border rounded"
        required
      />

      <input
        type="email"
        placeholder="Email (optional)"
        value={formData.email}
        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
        className="w-full px-4 py-2 border rounded"
      />

      {mutation.error && (
        <div className="bg-red-100 text-red-700 p-3 rounded">
          {(mutation.error as any).error}
        </div>
      )}

      <button
        type="submit"
        disabled={mutation.isPending}
        className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 disabled:opacity-50"
      >
        {mutation.isPending ? 'Creating...' : 'Create Tenant'}
      </button>
    </form>
  )
}
```

---

## Leases

### Create Lease Wizard

```typescript
// app/leases/create/page.tsx
'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useQuery, useMutation } from '@tanstack/react-query'
import { unitsService } from '@/services/units'
import { tenantsService } from '@/services/tenants'
import { leasesService } from '@/services/leases'
import type { CreateLeaseRequest } from '@/types/api/lease'

export default function CreateLeasePage() {
  const router = useRouter()
  const [step, setStep] = useState(1)
  const [formData, setFormData] = useState<Partial<CreateLeaseRequest>>({
    payment_due_day: 5,
    painting_fee_installments: 1,
  })

  // Step 1: Select Unit
  const { data: units } = useQuery({
    queryKey: ['units', 'available'],
    queryFn: () => unitsService.getAll({ status: 'available' }),
  })

  // Step 2: Select Tenant
  const { data: tenants } = useQuery({
    queryKey: ['tenants'],
    queryFn: () => tenantsService.getAll(),
  })

  // Create lease
  const mutation = useMutation({
    mutationFn: leasesService.create,
    onSuccess: () => {
      router.push('/leases')
    },
  })

  const handleSubmit = () => {
    mutation.mutate(formData as CreateLeaseRequest)
  }

  return (
    <div className="max-w-2xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-6">Create New Lease</h1>

      {/* Step Indicator */}
      <div className="flex mb-8">
        {[1, 2, 3].map((s) => (
          <div
            key={s}
            className={`flex-1 h-2 rounded ${s <= step ? 'bg-blue-600' : 'bg-gray-200'}`}
          />
        ))}
      </div>

      {/* Step 1: Select Unit */}
      {step === 1 && (
        <div>
          <h2 className="text-xl mb-4">Select Unit</h2>
          <div className="grid gap-4">
            {units?.map((unit) => (
              <button
                key={unit.id}
                onClick={() => {
                  setFormData({ ...formData, unit_id: unit.id, monthly_rent_value: unit.current_rent_value })
                  setStep(2)
                }}
                className="text-left p-4 border rounded hover:bg-gray-50"
              >
                <div className="font-semibold">Unit {unit.number}</div>
                <div className="text-sm text-gray-600">Floor {unit.floor}</div>
                <div className="font-bold mt-2">R$ {parseFloat(unit.current_rent_value).toFixed(2)}/month</div>
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Step 2: Select Tenant */}
      {step === 2 && (
        <div>
          <h2 className="text-xl mb-4">Select Tenant</h2>
          <div className="grid gap-4">
            {tenants?.map((tenant) => (
              <button
                key={tenant.id}
                onClick={() => {
                  setFormData({ ...formData, tenant_id: tenant.id })
                  setStep(3)
                }}
                className="text-left p-4 border rounded hover:bg-gray-50"
              >
                <div className="font-semibold">{tenant.full_name}</div>
                <div className="text-sm text-gray-600">{tenant.cpf}</div>
                <div className="text-sm text-gray-600">{tenant.phone}</div>
              </button>
            ))}
          </div>
          <button onClick={() => setStep(1)} className="mt-4 text-blue-600">
            ← Back
          </button>
        </div>
      )}

      {/* Step 3: Contract Details */}
      {step === 3 && (
        <div>
          <h2 className="text-xl mb-4">Contract Details</h2>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-1">Contract Signed Date</label>
              <input
                type="date"
                value={formData.contract_signed_date || ''}
                onChange={(e) => setFormData({ ...formData, contract_signed_date: e.target.value })}
                className="w-full px-4 py-2 border rounded"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-1">Start Date</label>
              <input
                type="date"
                value={formData.start_date || ''}
                onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                className="w-full px-4 py-2 border rounded"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-1">Payment Due Day (1-31)</label>
              <input
                type="number"
                min="1"
                max="31"
                value={formData.payment_due_day || 5}
                onChange={(e) => setFormData({ ...formData, payment_due_day: parseInt(e.target.value) })}
                className="w-full px-4 py-2 border rounded"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-1">Painting Fee Total</label>
              <input
                type="number"
                step="0.01"
                value={formData.painting_fee_total || ''}
                onChange={(e) => setFormData({ ...formData, painting_fee_total: e.target.value })}
                className="w-full px-4 py-2 border rounded"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-1">Painting Fee Installments</label>
              <select
                value={formData.painting_fee_installments || 1}
                onChange={(e) => setFormData({ ...formData, painting_fee_installments: parseInt(e.target.value) })}
                className="w-full px-4 py-2 border rounded"
              >
                <option value="1">1x</option>
                <option value="2">2x</option>
                <option value="3">3x</option>
                <option value="4">4x</option>
              </select>
            </div>

            {mutation.error && (
              <div className="bg-red-100 text-red-700 p-3 rounded">
                {(mutation.error as any).error}
              </div>
            )}

            <div className="flex gap-4">
              <button
                onClick={() => setStep(2)}
                className="flex-1 bg-gray-200 py-2 rounded hover:bg-gray-300"
              >
                ← Back
              </button>
              <button
                onClick={handleSubmit}
                disabled={mutation.isPending}
                className="flex-1 bg-blue-600 text-white py-2 rounded hover:bg-blue-700 disabled:opacity-50"
              >
                {mutation.isPending ? 'Creating...' : 'Create Lease'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
```

---

## Payments

### Mark Payment as Paid

```typescript
// components/PaymentCard.tsx
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { paymentsService } from '@/services/payments'
import type { Payment } from '@/types/api/payment'

export function PaymentCard({ payment }: { payment: Payment }) {
  const queryClient = useQueryClient()

  const markAsPaid = useMutation({
    mutationFn: () =>
      paymentsService.markAsPaid(payment.id, {
        payment_date: new Date().toISOString().split('T')[0],
        payment_method: 'pix',
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payments'] })
    },
  })

  const getStatusColor = (status: string) => {
    const colors = {
      pending: 'bg-yellow-100 text-yellow-800',
      paid: 'bg-green-100 text-green-800',
      overdue: 'bg-red-100 text-red-800',
      cancelled: 'bg-gray-100 text-gray-800',
    }
    return colors[status as keyof typeof colors]
  }

  return (
    <div className="border rounded-lg p-4">
      <div className="flex justify-between items-start mb-2">
        <div>
          <h3 className="font-semibold">{payment.payment_type}</h3>
          <p className="text-sm text-gray-600">
            {new Date(payment.reference_month).toLocaleDateString('pt-BR', {
              month: 'long',
              year: 'numeric',
            })}
          </p>
        </div>
        <span className={`px-2 py-1 text-xs rounded ${getStatusColor(payment.status)}`}>
          {payment.status}
        </span>
      </div>

      <p className="text-2xl font-bold mb-2">
        R$ {parseFloat(payment.amount).toFixed(2)}
      </p>

      <p className="text-sm text-gray-600">
        Due: {new Date(payment.due_date).toLocaleDateString('pt-BR')}
      </p>

      {payment.status === 'pending' || payment.status === 'overdue' ? (
        <button
          onClick={() => markAsPaid.mutate()}
          disabled={markAsPaid.isPending}
          className="w-full mt-4 bg-green-600 text-white py-2 rounded hover:bg-green-700 disabled:opacity-50"
        >
          {markAsPaid.isPending ? 'Processing...' : 'Mark as Paid'}
        </button>
      ) : payment.payment_date && (
        <p className="text-sm text-gray-600 mt-2">
          Paid on: {new Date(payment.payment_date).toLocaleDateString('pt-BR')}
        </p>
      )}
    </div>
  )
}
```

---

## Dashboard

### Dashboard Page

```typescript
// app/dashboard/page.tsx
'use client'

import { useQuery } from '@tanstack/react-query'
import { dashboardService } from '@/services/dashboard'

export default function DashboardPage() {
  const { data: dashboard, isLoading } = useQuery({
    queryKey: ['dashboard'],
    queryFn: dashboardService.getMetrics,
    refetchInterval: 60000, // Refresh every minute
  })

  if (isLoading) return <div>Loading dashboard...</div>

  return (
    <div className="p-6 space-y-6">
      <h1 className="text-3xl font-bold">Dashboard</h1>

      {/* Metrics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <MetricCard
          title="Occupancy Rate"
          value={`${dashboard?.occupancy.occupancy_rate.toFixed(1)}%`}
          subtitle={`${dashboard?.occupancy.occupied_units}/${dashboard?.occupancy.total_units} occupied`}
        />
        <MetricCard
          title="Monthly Revenue"
          value={`R$ ${parseFloat(dashboard?.financial.monthly_revenue || '0').toFixed(2)}`}
          subtitle="Expected"
        />
        <MetricCard
          title="Pending Payments"
          value={dashboard?.financial.pending_count.toString() || '0'}
          subtitle={`R$ ${parseFloat(dashboard?.financial.total_pending || '0').toFixed(2)}`}
        />
        <MetricCard
          title="Overdue Payments"
          value={dashboard?.financial.overdue_count.toString() || '0'}
          subtitle={`R$ ${parseFloat(dashboard?.financial.total_overdue || '0').toFixed(2)}`}
          alert={dashboard?.financial.overdue_count ? true : false}
        />
      </div>

      {/* Alerts */}
      {dashboard?.alerts && dashboard.alerts.length > 0 && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-bold mb-4">Alerts</h2>
          <div className="space-y-2">
            {dashboard.alerts.map((alert, i) => (
              <div
                key={i}
                className={`p-3 rounded ${
                  alert.severity === 'high'
                    ? 'bg-red-100 text-red-800'
                    : alert.severity === 'medium'
                    ? 'bg-yellow-100 text-yellow-800'
                    : 'bg-blue-100 text-blue-800'
                }`}
              >
                {alert.message}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

function MetricCard({
  title,
  value,
  subtitle,
  alert,
}: {
  title: string
  value: string
  subtitle?: string
  alert?: boolean
}) {
  return (
    <div className={`bg-white rounded-lg shadow p-6 ${alert ? 'border-2 border-red-500' : ''}`}>
      <h3 className="text-sm text-gray-600 mb-2">{title}</h3>
      <p className="text-3xl font-bold">{value}</p>
      {subtitle && <p className="text-sm text-gray-500 mt-1">{subtitle}</p>}
    </div>
  )
}
```

---

## React Query Provider Setup

```typescript
// app/providers.tsx
'use client'

import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { useState } from 'react'

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 60 * 1000, // 1 minute
            retry: 1,
          },
        },
      })
  )

  return (
    <QueryClientProvider client={queryClient}>
      {children}
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  )
}
```

```typescript
// app/layout.tsx
import { Providers } from './providers'

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="pt-BR">
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  )
}
```

---

Esses exemplos cobrem os principais casos de uso! Copie e adapte conforme necessário para seu projeto Next.js.
