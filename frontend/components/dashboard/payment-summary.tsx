"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { CreditCard } from "lucide-react"
import { useAnalyticsStore } from "@/lib/stores/analytics-store"
import { formatCurrency } from "@/lib/currency"
import { DonutChart } from "@/components/dashboard/charts"

export function PaymentSummary() {
  const { analytics, isLoading } = useAnalyticsStore()
  const payments = analytics?.payments

  if (isLoading && !payments) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Payments</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-24 bg-muted animate-pulse rounded" />
        </CardContent>
      </Card>
    )
  }

  if (!payments) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Payments</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">No payment data available.</p>
        </CardContent>
      </Card>
    )
  }

  const segments = [
    { label: "Successful", value: payments.successful_payments, color: "hsl(142 76% 36%)" },
    { label: "Pending", value: payments.pending_payments, color: "hsl(43 96% 56%)" },
    { label: "Failed", value: payments.failed_payments, color: "hsl(0 84% 60%)" },
  ]

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0">
        <CardTitle className="flex items-center gap-2">
          <CreditCard className="h-5 w-5" />
          Payments
        </CardTitle>
        <Badge variant="secondary">Summary</Badge>
      </CardHeader>
      <CardContent className="grid gap-4 md:grid-cols-[160px_1fr] md:items-center">
        <div className="flex justify-center">
          <DonutChart segments={segments} />
        </div>
        <div className="space-y-3">
          <div>
            <p className="text-sm text-muted-foreground">Net Amount</p>
            <p className="text-2xl font-bold">{formatCurrency(payments.net_amount)}</p>
          </div>
          <div className="grid grid-cols-2 gap-3 text-sm">
            <div>
              <p className="text-muted-foreground">Total</p>
              <p className="font-semibold">{payments.total_payments.toLocaleString()}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Refunded</p>
              <p className="font-semibold">{formatCurrency(payments.refunded_amount)}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Successful</p>
              <p className="font-semibold text-green-600">{payments.successful_payments.toLocaleString()}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Failed</p>
              <p className="font-semibold text-red-600">{payments.failed_payments.toLocaleString()}</p>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
