// Analytics types for dashboard

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

export interface DashboardAnalytics {
  orders: OrderAnalytics | null
  products: ProductAnalytics | null
  users: UserAnalytics | null
  period: string
  lastUpdated: string
}

export interface AnalyticsPeriod {
  label: string
  value: string
  days: number
}

export const ANALYTICS_PERIODS: AnalyticsPeriod[] = [
  { label: 'Last 7 days', value: '7d', days: 7 },
  { label: 'Last 30 days', value: '30d', days: 30 },
  { label: 'Last 90 days', value: '90d', days: 90 },
  { label: 'Last year', value: '1y', days: 365 },
]

export interface ChartDataPoint {
  date: string
  value: number
  label?: string
}

export interface SalesChartData {
  revenue: ChartDataPoint[]
  orders: ChartDataPoint[]
  users: ChartDataPoint[]
}

export interface TopProduct {
  id: number
  name: string
  sku: string
  revenue: number
  quantity: number
  orders: number
}

export interface RecentOrder {
  id: number
  order_number: string
  customer_name: string
  customer_email: string
  total_amount: number
  status: string
  created_at: string
}
