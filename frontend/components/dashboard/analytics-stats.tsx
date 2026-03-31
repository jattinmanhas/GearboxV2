"use client"

import { useEffect } from "react"
import { useAnalyticsStore } from "@/lib/stores/analytics-store"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { 
  Users, 
  ShoppingCart, 
  TrendingUp, 
  Package,
  DollarSign,
  Activity,
  AlertTriangle,
  CheckCircle,
  Clock,
  Truck,
  CreditCard,
  Eye
} from "lucide-react"
import { formatCurrency } from "@/lib/currency"

interface AnalyticsStatsProps {
  period?: string
}

export function AnalyticsStats({ period = '30d' }: AnalyticsStatsProps) {
  const { 
    analytics, 
    isLoading, 
    isRetrying,
    retryCount,
    error, 
    loadAnalytics 
  } = useAnalyticsStore()

  useEffect(() => {
    loadAnalytics(period)
  }, [loadAnalytics, period])

  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {[...Array(4)].map((_, i) => (
          <Card key={i}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Loading...</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="h-8 bg-muted animate-pulse rounded"></div>
            </CardContent>
          </Card>
        ))}
      </div>
    )
  }

  if (error) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card className="col-span-full">
          <CardContent className="pt-6">
            <div className="text-center text-destructive">
              <AlertTriangle className="h-8 w-8 mx-auto mb-2" />
              <p>Failed to load analytics: {error}</p>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  if (!analytics) {
    return null
  }

  const { orders, products, users, payments, blog, partial, errors: sectionErrors } = analytics
  const failedSections = Object.keys(sectionErrors || {})

  const stats = [
    {
      title: "Total Revenue",
      value: formatCurrency(orders?.total_revenue || 0),
      icon: DollarSign,
      description: `${orders?.new_orders_this_month || 0} orders this month`,
      color: "text-green-600",
      trend: orders?.total_revenue ? "Updated" : "No data"
    },
    {
      title: "Total Orders",
      value: orders?.total_orders?.toLocaleString() || "0",
      icon: ShoppingCart,
      description: `${orders?.new_orders_today || 0} new today`,
      color: "text-blue-600",
      trend: orders?.new_orders_this_week ? "Growing" : "No data"
    },
    {
      title: "Active Products",
      value: products?.active_products?.toLocaleString() || "0",
      icon: Package,
      description: `${products?.low_stock_products || 0} low stock`,
      color: "text-purple-600",
      trend: products?.active_products ? "Stable" : "No data"
    },
    {
      title: "Active Users",
      value: users?.active_users?.toLocaleString() || "0",
      icon: Users,
      description: `${users?.new_users_today || 0} new today`,
      color: "text-orange-600",
      trend: users?.new_users_this_week ? "Rising" : "No data"
    },
    {
      title: "Net Payments",
      value: formatCurrency(payments?.net_amount || 0),
      icon: CreditCard,
      description: `${payments?.successful_payments || 0} successful`,
      color: "text-emerald-600",
      trend: payments?.total_payments ? "Live" : "No data"
    },
    {
      title: "Blog Views",
      value: blog?.total_views?.toLocaleString() || "0",
      icon: Eye,
      description: `${blog?.published_posts || 0} published posts`,
      color: "text-sky-600",
      trend: blog?.total_views ? "Active" : "No data"
    }
  ]

  return (
    <div className="space-y-4">
      {isRetrying && (
        <Card className="border-blue-500/40">
          <CardContent className="pt-6">
            <div className="flex items-center gap-3 text-blue-700 dark:text-blue-400">
              <div className="h-4 w-4 animate-spin rounded-full border-2 border-blue-500 border-t-transparent" />
              <div>
                <p className="font-medium">Services are starting up...</p>
                <p className="text-sm text-muted-foreground">
                  Retrying automatically ({retryCount}/5). This usually takes a minute on first start.
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
      {!isRetrying && partial && failedSections.length > 0 && (
        <Card className="border-yellow-500/40">
          <CardContent className="pt-6">
            <div className="flex items-start gap-2 text-yellow-700 dark:text-yellow-400">
              <AlertTriangle className="h-5 w-5 mt-0.5" />
              <div>
                <p className="font-medium">Some analytics could not be loaded</p>
                <p className="text-sm">
                  Failed sections: {failedSections.join(", ")}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {stats.map((stat, index) => (
          <Card key={index}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{stat.title}</CardTitle>
              <stat.icon className={`h-4 w-4 ${stat.color}`} />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <p className="text-xs text-muted-foreground">
                {stat.description}
              </p>
              <div className="flex items-center mt-1">
                <TrendingUp className="h-3 w-3 text-green-500 mr-1" />
                <span className="text-xs text-green-500">{stat.trend}</span>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Order Status Breakdown */}
      {orders && (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Order Status</CardTitle>
              <CheckCircle className="h-4 w-4 text-green-600" />
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-sm">Delivered</span>
                  <span className="text-sm font-medium">{orders.delivered_orders}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm">Processing</span>
                  <span className="text-sm font-medium">{orders.processing_orders}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm">Pending</span>
                  <span className="text-sm font-medium">{orders.pending_orders}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm">Cancelled</span>
                  <span className="text-sm font-medium">{orders.cancelled_orders}</span>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Product Status</CardTitle>
              <Package className="h-4 w-4 text-blue-600" />
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-sm">Active</span>
                  <span className="text-sm font-medium">{products?.active_products || 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm">Inactive</span>
                  <span className="text-sm font-medium">{products?.inactive_products || 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm">Low Stock</span>
                  <span className="text-sm font-medium text-orange-600">{products?.low_stock_products || 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm">Out of Stock</span>
                  <span className="text-sm font-medium text-red-600">{products?.out_of_stock_products || 0}</span>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">User Growth</CardTitle>
              <Users className="h-4 w-4 text-purple-600" />
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-sm">Total Users</span>
                  <span className="text-sm font-medium">{users?.total_users || 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm">This Week</span>
                  <span className="text-sm font-medium">{users?.new_users_this_week || 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm">This Month</span>
                  <span className="text-sm font-medium">{users?.new_users_this_month || 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm">Today</span>
                  <span className="text-sm font-medium">{users?.new_users_today || 0}</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  )
}
