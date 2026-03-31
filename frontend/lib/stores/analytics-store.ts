import { create } from 'zustand'
import { devtools } from 'zustand/middleware'
import { httpClient } from '@/lib/apiFunctions/http-client'

// Types for analytics data
export interface OrderAnalytics {
  total_orders: number
  pending_orders: number
  confirmed_orders: number
  processing_orders: number
  shipped_orders: number
  delivered_orders: number
  cancelled_orders: number
  total_revenue: number
  average_order_value: number
  conversion_rate: number
  new_orders_today: number
  new_orders_this_week: number
  new_orders_this_month: number
}

export interface ProductAnalytics {
  total_products: number
  active_products: number
  inactive_products: number
  low_stock_products: number
  out_of_stock_products: number
  total_categories: number
  average_price: number
  total_inventory_value: number
  top_selling_products: Array<{
    product_id: number
    product_name: string
    sku: string
    total_quantity: number
    total_revenue: number
    order_count: number
  }>
}

export interface UserAnalytics {
  total_users: number
  active_users: number
  new_users_today: number
  new_users_this_week: number
  new_users_this_month: number
  users_by_role: Array<{
    role: string
    count: number
  }>
  user_registration_trend: Array<{
    date: string
    count: number
  }>
}

export interface PaymentSummary {
  total_payments: number
  successful_payments: number
  failed_payments: number
  pending_payments: number
  total_amount: number
  refunded_amount: number
  net_amount: number
}

export interface BlogAnalytics {
  total_posts: number
  published_posts: number
  draft_posts: number
  archived_posts: number
  total_views: number
  average_read_time: number
  top_posts: Array<{
    id: string
    title: string
    slug: string
    viewCount: number
    authorName: string
    publishedAt: string | null
  }>
  recent_posts: Array<{
    id: string
    title: string
    slug: string
    viewCount: number
    authorName: string
    publishedAt: string | null
  }>
}

export interface DashboardAnalytics {
  orders: OrderAnalytics | null
  products: ProductAnalytics | null
  users: UserAnalytics | null
  payments: PaymentSummary | null
  blog: BlogAnalytics | null
  top_products: Array<{
    product_id: number
    product_name: string
    sku: string
    total_quantity: number
    total_revenue: number
    order_count: number
    average_price?: number
  }> | null
  period: string
  lastUpdated: string
  partial?: boolean
  errors?: Record<string, string>
}

interface AnalyticsState {
  // State
  analytics: DashboardAnalytics | null
  isLoading: boolean
  isRetrying: boolean
  retryCount: number
  error: string | null
  
  // Actions
  loadAnalytics: (period?: string, attempt?: number) => Promise<void>
  loadOrderAnalytics: () => Promise<void>
  loadProductAnalytics: () => Promise<void>
  loadUserAnalytics: () => Promise<void>
  clearError: () => void
}

export const useAnalyticsStore = create<AnalyticsState>()(
  devtools(
    (set, get) => ({
      // Initial state
      analytics: null,
      isLoading: false,
      isRetrying: false,
      retryCount: 0,
      error: null,

      // Load comprehensive analytics with automatic retry when services are starting
      loadAnalytics: async (period = '30d', attempt = 0) => {
        const MAX_RETRIES = 5
        const RETRY_DELAYS = [2000, 4000, 8000, 16000, 30000]

        if (attempt === 0) {
          set({ isLoading: true, isRetrying: false, retryCount: 0, error: null })
        } else {
          set({ isLoading: false, isRetrying: true, retryCount: attempt })
        }
        
        try {
          const result = await httpClient.get<{
            success: boolean
            message?: string
            data?: DashboardAnalytics
            error?: boolean
            status?: number
          }>(`/dashboard/analytics?period=${period}`)
          
          if (result.success && result.data) {
            // Check if ALL sections failed (services still starting up)
            const data = result.data
            const hasAnyData = data.orders || data.products || data.users || data.payments
            const allSectionsFailed = data.partial && !hasAnyData

            if (allSectionsFailed && attempt < MAX_RETRIES) {
              // Services are still starting, schedule a retry
              console.log(`Analytics: all sections failed, retrying in ${RETRY_DELAYS[attempt]}ms (attempt ${attempt + 1}/${MAX_RETRIES})`)
              set({ analytics: data, isLoading: false, isRetrying: true, retryCount: attempt + 1 })
              setTimeout(() => {
                get().loadAnalytics(period, attempt + 1)
              }, RETRY_DELAYS[attempt])
            } else {
              set({ 
                analytics: data, 
                isLoading: false,
                isRetrying: false,
                retryCount: 0,
                error: null 
              })
            }
          } else {
            throw new Error(result.message || 'Failed to load analytics')
          }
        } catch (error) {
          console.error('Error loading analytics:', error)
          if (attempt < MAX_RETRIES) {
            // Network error - services likely still starting
            console.log(`Analytics: network error, retrying in ${RETRY_DELAYS[attempt]}ms (attempt ${attempt + 1}/${MAX_RETRIES})`)
            set({ isLoading: false, isRetrying: true, retryCount: attempt + 1 })
            setTimeout(() => {
              get().loadAnalytics(period, attempt + 1)
            }, RETRY_DELAYS[attempt])
          } else {
            set({ 
              error: error instanceof Error ? error.message : 'Unknown error',
              isLoading: false,
              isRetrying: false,
            })
          }
        }
      },

      // Load order analytics only
      loadOrderAnalytics: async () => {
        set({ isLoading: true, error: null })
        
        try {
          const result = await httpClient.get<{
            success: boolean
            message?: string
            data?: OrderAnalytics
            error?: boolean
            status?: number
          }>('/dashboard/orders')
          
          if (result.success && result.data) {
            const currentAnalytics = get().analytics
            set({ 
              analytics: {
                orders: result.data,
                products: currentAnalytics?.products || null,
                users: currentAnalytics?.users || null,
                payments: currentAnalytics?.payments || null,
                blog: currentAnalytics?.blog || null,
                top_products: currentAnalytics?.top_products || null,
                period: currentAnalytics?.period || '30d',
                lastUpdated: new Date().toISOString(),
              },
              isLoading: false,
              error: null 
            })
          } else {
            throw new Error(result.message || 'Failed to load order analytics')
          }
        } catch (error) {
          console.error('Error loading order analytics:', error)
          set({ 
            error: error instanceof Error ? error.message : 'Unknown error',
            isLoading: false 
          })
        }
      },

      // Load product analytics only
      loadProductAnalytics: async () => {
        set({ isLoading: true, error: null })
        
        try {
          const result = await httpClient.get<{
            success: boolean
            message?: string
            data?: ProductAnalytics
            error?: boolean
            status?: number
          }>('/dashboard/products')
          
          if (result.success && result.data) {
            const currentAnalytics = get().analytics
            set({ 
              analytics: {
                orders: currentAnalytics?.orders || null,
                products: result.data,
                users: currentAnalytics?.users || null,
                payments: currentAnalytics?.payments || null,
                blog: currentAnalytics?.blog || null,
                top_products: currentAnalytics?.top_products || null,
                period: currentAnalytics?.period || '30d',
                lastUpdated: new Date().toISOString(),
              },
              isLoading: false,
              error: null 
            })
          } else {
            throw new Error(result.message || 'Failed to load product analytics')
          }
        } catch (error) {
          console.error('Error loading product analytics:', error)
          set({ 
            error: error instanceof Error ? error.message : 'Unknown error',
            isLoading: false 
          })
        }
      },

      // Load user analytics only
      loadUserAnalytics: async () => {
        set({ isLoading: true, error: null })
        
        try {
          const result = await httpClient.get<{
            success: boolean
            message?: string
            data?: UserAnalytics
            error?: boolean
            status?: number
          }>('/dashboard/users')
          
          if (result.success && result.data) {
            const currentAnalytics = get().analytics
            set({ 
              analytics: {
                orders: currentAnalytics?.orders || null,
                products: currentAnalytics?.products || null,
                users: result.data,
                payments: currentAnalytics?.payments || null,
                blog: currentAnalytics?.blog || null,
                top_products: currentAnalytics?.top_products || null,
                period: currentAnalytics?.period || '30d',
                lastUpdated: new Date().toISOString(),
              },
              isLoading: false,
              error: null 
            })
          } else {
            throw new Error(result.message || 'Failed to load user analytics')
          }
        } catch (error) {
          console.error('Error loading user analytics:', error)
          set({ 
            error: error instanceof Error ? error.message : 'Unknown error',
            isLoading: false 
          })
        }
      },

      // Clear error
      clearError: () => set({ error: null }),
    }),
    {
      name: 'analytics-store',
    }
  )
)
