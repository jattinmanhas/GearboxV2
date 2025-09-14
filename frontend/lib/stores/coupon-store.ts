import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { couponApi } from '../coupon-api'
import { Coupon, CreateCouponRequest, UpdateCouponRequest, CouponUsage, CouponStats } from '../types/coupon'

interface CouponStore {
  // State
  coupons: Coupon[]
  currentCoupon: Coupon | null
  couponUsage: CouponUsage[]
  couponStats: CouponStats | null
  isLoading: boolean
  error: string | null

  // Actions
  setCoupons: (coupons: Coupon[]) => void
  setCurrentCoupon: (coupon: Coupon | null) => void
  setCouponUsage: (usage: CouponUsage[]) => void
  setCouponStats: (stats: CouponStats | null) => void
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void

  // API Actions
  loadCoupons: (params?: {
    page?: number
    limit?: number
    search?: string
    status?: string
    type?: string
  }) => Promise<void>
  loadCoupon: (id: string) => Promise<void>
  createCoupon: (couponData: CreateCouponRequest) => Promise<void>
  updateCoupon: (id: string, couponData: UpdateCouponRequest) => Promise<void>
  deleteCoupon: (id: string) => Promise<void>
  loadCouponUsage: (params?: {
    page?: number
    limit?: number
    coupon_id?: string
    start_date?: string
    end_date?: string
  }) => Promise<void>
  loadCouponStats: () => Promise<void>
  validateCoupon: (couponCode: string) => Promise<{ valid: boolean; coupon?: Coupon; discount_amount?: number; message?: string }>
}

export const useCouponStore = create<CouponStore>()(
  persist(
    (set, get) => ({
      // State
      coupons: [],
      currentCoupon: null,
      couponUsage: [],
      couponStats: null,
      isLoading: false,
      error: null,

      // Actions
      setCoupons: (coupons) => set({ coupons }),
      setCurrentCoupon: (coupon) => set({ currentCoupon: coupon }),
      setCouponUsage: (usage) => set({ couponUsage: usage }),
      setCouponStats: (stats) => set({ couponStats: stats }),
      setLoading: (loading) => set({ isLoading: loading }),
      setError: (error) => set({ error }),

      // API Actions
      loadCoupons: async (params) => {
        try {
          set({ isLoading: true, error: null })
          const response = await couponApi.getCoupons(params)
          set({ coupons: response.data.coupons, isLoading: false })
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to load coupons',
            isLoading: false 
          })
        }
      },

      loadCoupon: async (id) => {
        try {
          set({ isLoading: true, error: null })
          const response = await couponApi.getCoupon(id)
          set({ currentCoupon: response.data, isLoading: false })
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to load coupon',
            isLoading: false 
          })
        }
      },

      createCoupon: async (couponData) => {
        try {
          set({ isLoading: true, error: null })
          await couponApi.createCoupon(couponData)
          // Reload coupons after creation
          await get().loadCoupons()
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to create coupon',
            isLoading: false 
          })
        }
      },

      updateCoupon: async (id, couponData) => {
        try {
          set({ isLoading: true, error: null })
          await couponApi.updateCoupon(id, couponData)
          // Reload coupons after update
          await get().loadCoupons()
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to update coupon',
            isLoading: false 
          })
        }
      },

      deleteCoupon: async (id) => {
        try {
          set({ isLoading: true, error: null })
          await couponApi.deleteCoupon(id)
          // Reload coupons after deletion
          await get().loadCoupons()
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to delete coupon',
            isLoading: false 
          })
        }
      },

      loadCouponUsage: async (params) => {
        try {
          set({ isLoading: true, error: null })
          const response = await couponApi.getCouponUsage(params)
          set({ couponUsage: response.data.usage, isLoading: false })
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to load coupon usage',
            isLoading: false 
          })
        }
      },

      loadCouponStats: async () => {
        try {
          set({ isLoading: true, error: null })
          // This would need to be implemented in the backend
          // For now, we'll calculate basic stats from coupons
          const coupons = get().coupons
          const stats: CouponStats = {
            total_coupons: coupons.length,
            active_coupons: coupons.filter(c => c.is_active).length,
            expired_coupons: coupons.filter(c => c.expires_at && new Date(c.expires_at) < new Date()).length,
            total_usage: coupons.reduce((sum, c) => sum + c.used_count, 0),
            total_discount_given: 0, // This would need to be calculated from usage data
            most_used_coupon: coupons.length > 0 ? {
              id: coupons.reduce((max, c) => c.used_count > max.used_count ? c : max).id,
              code: coupons.reduce((max, c) => c.used_count > max.used_count ? c : max).code,
              usage_count: coupons.reduce((max, c) => c.used_count > max.used_count ? c : max).used_count
            } : undefined
          }
          set({ couponStats: stats, isLoading: false })
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to load coupon stats',
            isLoading: false 
          })
        }
      },

      validateCoupon: async (couponCode) => {
        try {
          set({ isLoading: true, error: null })
          const response = await couponApi.validateCoupon(couponCode)
          set({ isLoading: false })
          return response.data
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to validate coupon',
            isLoading: false 
          })
          return { valid: false, message: 'Failed to validate coupon' }
        }
      },
    }),
    {
      name: 'coupon-store',
      partialize: (state) => ({
        coupons: state.coupons,
        currentCoupon: state.currentCoupon,
        couponStats: state.couponStats,
      }),
    }
  )
)
