"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Users } from "lucide-react"
import { useAnalyticsStore } from "@/lib/stores/analytics-store"
import { LineChart } from "@/components/dashboard/charts"

export function UserGrowthChart() {
  const { analytics, isLoading } = useAnalyticsStore()
  const trend = analytics?.users?.user_registration_trend || []

  const chartData = trend.map((point) => ({
    label: point.date,
    value: point.count,
  }))

  if (isLoading && trend.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>User Growth</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-24 bg-muted animate-pulse rounded" />
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0">
        <CardTitle className="flex items-center gap-2">
          <Users className="h-5 w-5" />
          User Growth
        </CardTitle>
        <Badge variant="secondary">Trend</Badge>
      </CardHeader>
      <CardContent>
        {chartData.length === 0 ? (
          <p className="text-sm text-muted-foreground">No user growth data available.</p>
        ) : (
          <LineChart data={chartData} height={120} />
        )}
      </CardContent>
    </Card>
  )
}
