// Dashboard & Reports Types

export interface OccupancyMetrics {
  total_units: number
  available_units: number
  occupied_units: number
  maintenance_units: number
  renovation_units: number
  occupancy_rate: number
}

export interface FinancialMetrics {
  monthly_revenue: string // decimal as string
  total_receivable: string
  total_received: string
  total_pending: string
  total_overdue: string
  pending_count: number
  overdue_count: number
}

export interface ContractMetrics {
  total_active_leases: number
  expiring_soon_count: number
  expired_count: number
}

export interface Alert {
  type: 'contract_expiring' | 'payment_overdue' | 'unit_maintenance'
  severity: 'high' | 'medium' | 'low'
  message: string
  entity_id: string
  entity_type: 'lease' | 'payment' | 'unit'
}

export interface DashboardResponse {
  occupancy: OccupancyMetrics
  financial: FinancialMetrics
  contracts: ContractMetrics
  alerts: Alert[]
}

export interface FinancialReportRequest {
  start_date: string // YYYY-MM-DD
  end_date: string // YYYY-MM-DD
  payment_type?: PaymentType
  status?: PaymentStatus
}

export interface FinancialReportResponse {
  period: {
    start_date: string
    end_date: string
  }
  summary: {
    total_amount: string
    paid_amount: string
    pending_amount: string
    overdue_amount: string
    payment_count: number
  }
  by_type: {
    [key in PaymentType]: {
      count: number
      total_amount: string
    }
  }
  by_status: {
    [key in PaymentStatus]: {
      count: number
      total_amount: string
    }
  }
  payments: Payment[]
}

export interface PaymentHistoryRequest {
  lease_id?: string
  tenant_id?: string
  status?: PaymentStatus
  start_date?: string // YYYY-MM-DD
  end_date?: string // YYYY-MM-DD
}

export interface PaymentHistoryResponse {
  total_count: number
  payments: Payment[]
}

// Import PaymentType, PaymentStatus, Payment from payment.ts
import type { PaymentType, PaymentStatus, Payment } from './payment'
