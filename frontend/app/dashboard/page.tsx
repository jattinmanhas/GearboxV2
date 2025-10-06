"use client"

import { AnalyticsStats } from "@/components/dashboard/analytics-stats"
import { SalesChart } from "@/components/dashboard/sales-chart"
import { RecentOrders } from "@/components/dashboard/recent-orders"

export default function DashboardPage() {
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
      </div>

      {/* Dashboard Content */}
      <div className="space-y-6">
        {/* Analytics Stats */}
        <AnalyticsStats />

        {/* Main Content Grid */}
        <div className="space-y-6">
          {/* Sales Overview Chart - Full Width */}
          <SalesChart />
          
          {/* Recent Orders - Horizontal Layout */}
          <RecentOrders />
        </div>
      </div>
    </div>
  )
}
