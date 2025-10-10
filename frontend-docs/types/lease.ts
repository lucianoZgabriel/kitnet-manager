// Lease Types

export type LeaseStatus = 'active' | 'expiring_soon' | 'expired' | 'cancelled'

export interface Lease {
  id: string
  unit_id: string
  tenant_id: string
  contract_signed_date: string // ISO date
  start_date: string // ISO date
  end_date: string // ISO date
  payment_due_day: number // 1-31
  monthly_rent_value: string // decimal as string
  painting_fee_total: string // decimal as string
  painting_fee_installments: number // 1-4
  painting_fee_paid: string // decimal as string
  status: LeaseStatus
  created_at: string
  updated_at: string
}

export interface CreateLeaseRequest {
  unit_id: string
  tenant_id: string
  contract_signed_date: string // ISO date
  start_date: string // ISO date
  payment_due_day: number // 1-31
  monthly_rent_value: string
  painting_fee_total: string
  painting_fee_installments: number // 1-4
}

export interface RenewLeaseRequest {
  painting_fee_total: string
  painting_fee_installments: number // 1-4
}

export interface UpdatePaintingFeePaidRequest {
  amount_paid: string
}

export interface LeaseStats {
  total_leases: number
  active_leases: number
  expiring_soon_leases: number
  expired_leases: number
  cancelled_leases: number
}
