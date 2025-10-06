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
  Truck
} from "lucide-react"
import { formatCurrency } from "@/lib/currency"

export function AnalyticsStats() {
  const { 
    analytics, 
    isLoading, 
    error, 
    loadAnalytics 
  } = useAnalyticsStore()

  useEffect(() => {
    loadAnalytics('30d')
  }, [loadAnalytics])

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

  const { orders, products, users } = analytics

  const stats = [
    {
      title: "Total Revenue",
      value: formatCurrency(orders?.total_revenue || 0),
      icon: DollarSign,
      description: `${orders?.new_orders_this_month || 0} orders this month`,
      color: "text-green-600",
      trend: orders?.total_revenue ? "+12.5%" : "0%"
    },
    {
      title: "Total Orders",
      value: orders?.total_orders?.toLocaleString() || "0",
      icon: ShoppingCart,
      description: `${orders?.new_orders_today || 0} new today`,
      color: "text-blue-600",
      trend: orders?.new_orders_this_week ? "+8.2%" : "0%"
    },
    {
      title: "Active Products",
      value: products?.active_products?.toLocaleString() || "0",
      icon: Package,
      description: `${products?.low_stock_products || 0} low stock`,
      color: "text-purple-600",
      trend: products?.active_products ? "+5.1%" : "0%"
    },
    {
      title: "Active Users",
      value: users?.active_users?.toLocaleString() || "0",
      icon: Users,
      description: `${users?.new_users_today || 0} new today`,
      color: "text-orange-600",
      trend: users?.new_users_this_week ? "+15.3%" : "0%"
    }
  ]

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Analytics Overview</h2>
          <p className="text-muted-foreground">
            Key metrics for your business performance
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <Badge variant="secondary" className="flex items-center gap-1">
            <Activity className="h-3 w-3" />
            Live
          </Badge>
          <Badge variant="outline">
            {analytics.period}
          </Badge>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
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
