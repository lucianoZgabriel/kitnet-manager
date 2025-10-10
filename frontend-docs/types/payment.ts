// Payment Types

export type PaymentType = 'rent' | 'painting_fee' | 'adjustment'
export type PaymentStatus = 'pending' | 'paid' | 'overdue' | 'cancelled'
export type PaymentMethod = 'pix' | 'cash' | 'bank_transfer' | 'credit_card'

export interface Payment {
  id: string
  lease_id: string
  payment_type: PaymentType
  reference_month: string // ISO date (YYYY-MM-DD)
  amount: string // decimal as string
  status: PaymentStatus
  due_date: string // ISO date
  payment_date?: string // ISO date
  payment_method?: PaymentMethod
  proof_url?: string
  notes?: string
  created_at: string
  updated_at: string
}

export interface MarkPaymentAsPaidRequest {
  payment_date: string // ISO date
  payment_method: PaymentMethod
}

export interface PaymentStats {
  total_payments: number
  total_amount: string
  paid_amount: string
  pending_amount: string
  overdue_amount: string
  paid_count: number
  pending_count: number
  overdue_count: number
}
