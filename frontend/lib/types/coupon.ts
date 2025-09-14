export interface Coupon {
  id: number
  code: string
  name: string
  description?: string
  type: 'percentage' | 'fixed_amount' | 'free_shipping'
  value: number
  min_order_amount?: number
  max_discount_amount?: number
  usage_limit?: number
  used_count: number
  is_active: boolean
  starts_at: string
  expires_at?: string
  created_at: string
  updated_at: string
}

export interface CreateCouponRequest {
  code: string
  name: string
  description?: string
  type: 'percentage' | 'fixed_amount' | 'free_shipping'
  value: number
  min_order_amount?: number
  max_discount_amount?: number
  usage_limit?: number
  is_active: boolean
  starts_at: string
  expires_at?: string
}

export interface UpdateCouponRequest {
  name?: string
  description?: string
  type?: 'percentage' | 'fixed_amount' | 'free_shipping'
  value?: number
  min_order_amount?: number
  max_discount_amount?: number
  usage_limit?: number
  is_active?: boolean
  starts_at?: string
  expires_at?: string
}

export interface CouponUsage {
  id: number
  coupon_id: number
  order_id: number
  user_id: number
  discount_amount: number
  used_at: string
  coupon_code: string
  order_total: number
}

export interface CouponStats {
  total_coupons: number
  active_coupons: number
  expired_coupons: number
  total_usage: number
  total_discount_given: number
  most_used_coupon?: {
    id: number
    code: string
    usage_count: number
  }
}

export interface CouponListResponse {
  coupons: Coupon[]
  total: number
  page: number
  limit: number
  total_pages: number
}

export interface CouponUsageResponse {
  usage: CouponUsage[]
  total: number
  page: number
  limit: number
  total_pages: number
}

export interface ValidateCouponRequest {
  coupon_code: string
}

export interface ValidateCouponResponse {
  valid: boolean
  coupon?: Coupon
  discount_amount?: number
  message?: string
}
