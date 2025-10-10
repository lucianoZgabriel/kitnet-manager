// Unit Types

export type UnitStatus = 'available' | 'occupied' | 'maintenance' | 'renovation'

export interface Unit {
  id: string
  number: string
  floor: number
  status: UnitStatus
  is_renovated: boolean
  base_rent_value: string // decimal as string
  renovated_rent_value: string // decimal as string
  current_rent_value: string // decimal as string
  notes?: string
  created_at: string
  updated_at: string
}

export interface CreateUnitRequest {
  number: string
  floor: number
  base_rent_value: string
  renovated_rent_value: string
}

export interface UpdateUnitRequest {
  number: string
  floor: number
  is_renovated: boolean
  base_rent_value: string
  renovated_rent_value: string
  notes?: string
}

export interface UpdateUnitStatusRequest {
  status: UnitStatus
}

export interface OccupancyStats {
  total_units: number
  available_units: number
  occupied_units: number
  maintenance_units: number
  renovation_units: number
  occupancy_rate: number
}
