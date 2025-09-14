"use client"

import { useEffect } from "react"
import { useCouponStore } from "@/lib/stores/coupon-store"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { 
  Ticket, 
  TrendingUp, 
  Clock, 
  DollarSign,
  Users,
  Award
} from "lucide-react"

export function CouponStats() {
  const { 
    couponStats, 
    isLoading, 
    loadCouponStats 
  } = useCouponStore()

  useEffect(() => {
    loadCouponStats()
  }, [loadCouponStats])

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

  if (!couponStats) {
    return null
  }

  const stats = [
    {
      title: "Total Coupons",
      value: couponStats.total_coupons,
      icon: Ticket,
      description: "All created coupons",
      color: "text-blue-600"
    },
    {
      title: "Active Coupons",
      value: couponStats.active_coupons,
      icon: TrendingUp,
      description: "Currently active",
      color: "text-green-600"
    },
    {
      title: "Total Usage",
      value: couponStats.total_usage,
      icon: Users,
      description: "Times used",
      color: "text-purple-600"
    },
    {
      title: "Discount Given",
      value: `₹${couponStats.total_discount_given.toLocaleString()}`,
      icon: DollarSign,
      description: "Total savings",
      color: "text-orange-600"
    }
  ]

  return (
    <div className="space-y-4">
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                {stat.title}
              </CardTitle>
              <stat.icon className={`h-4 w-4 ${stat.color}`} />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <p className="text-xs text-muted-foreground">
                {stat.description}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>

      {couponStats.most_used_coupon && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Award className="h-5 w-5 text-yellow-600" />
              Most Used Coupon
            </CardTitle>
            <CardDescription>
              The most popular coupon by usage count
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-between">
              <div>
                <p className="font-mono font-medium text-lg">
                  {couponStats.most_used_coupon.code}
                </p>
                <p className="text-sm text-muted-foreground">
                  Used {couponStats.most_used_coupon.usage_count} times
                </p>
              </div>
              <Badge variant="secondary">
                {couponStats.most_used_coupon.usage_count} uses
              </Badge>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
