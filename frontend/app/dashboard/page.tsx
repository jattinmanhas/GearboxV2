"use client"

import { useState } from "react"
import { AnalyticsStats } from "@/components/dashboard/analytics-stats"
import { SalesChart } from "@/components/dashboard/sales-chart"
import { RecentOrders } from "@/components/dashboard/recent-orders"
import { PaymentSummary } from "@/components/dashboard/payment-summary"
import { TopProductsChart } from "@/components/dashboard/top-products-chart"
import { UserGrowthChart } from "@/components/dashboard/user-growth-chart"
import { BlogInsights } from "@/components/dashboard/blog-insights"
import { Badge } from "@/components/ui/badge"

export default function DashboardPage() {
  const [period, setPeriod] = useState("7d")
  const periods = [
    { label: "7D", value: "7d" },
    { label: "30D", value: "30d" },
    { label: "90D", value: "90d" },
    { label: "1Y", value: "1y" },
  ]

  return (
    <div className="space-y-6">
      {/* Dashboard Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground">
            Welcome back! Here's what's happening with your business today.
          </p>
        </div>
        <div className="flex items-center gap-2">
          {periods.map((periodOption) => (
            <Badge
              key={periodOption.value}
              variant={period === periodOption.value ? "default" : "outline"}
              className="cursor-pointer"
              onClick={() => setPeriod(periodOption.value)}
            >
              {periodOption.label}
            </Badge>
          ))}
        </div>
      </div>

      {/* Dashboard Content */}
      <div className="space-y-6">
        {/* Analytics Stats */}
        <AnalyticsStats period={period} />

        {/* Main Content Grid */}
        <div className="space-y-6">
          {/* Sales Overview Chart - Full Width */}
          <SalesChart period={period} />
          
          {/* Recent Orders - Horizontal Layout */}
          <RecentOrders />
        </div>

        <div className="grid gap-6 lg:grid-cols-2">
          <PaymentSummary />
          <UserGrowthChart />
        </div>

        <div className="grid gap-6 lg:grid-cols-2">
          <TopProductsChart />
          <BlogInsights />
        </div>
      </div>
    </div>
  )
}
