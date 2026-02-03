import { useEffect, useState } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
  TrendingUp,
  BarChart3,
  DollarSign,
  ShoppingCart
} from "lucide-react"
import { formatCurrency } from "@/lib/currency"
import { httpClient } from "@/lib/apiFunctions/http-client"

interface ChartDataPoint {
  date: string
  revenue: number
  orders: number
}

interface SalesChartProps {
  period?: string
}

export function SalesChart({ period }: SalesChartProps) {
  const [chartData, setChartData] = useState<ChartDataPoint[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedPeriod, setSelectedPeriod] = useState(period || '30d')
  const [hoveredIndex, setHoveredIndex] = useState<number | null>(null)

  useEffect(() => {
    if (period) {
      setSelectedPeriod(period)
    }
  }, [period])

  useEffect(() => {
    const fetchChartData = async () => {
      try {
        setLoading(true)

        const result = await httpClient.get<{
          success: boolean
          message: string
          data: ChartDataPoint[]
        }>(`/dashboard/sales-chart?period=${selectedPeriod}`)

        if (result.success) {
          const data = result.data || []
          setChartData(data)
        } else {
          throw new Error(result.message || 'Failed to load chart data')
        }
      } catch (error) {
        console.error('Error loading chart data:', error)
        const errorMessage = error instanceof Error ? error.message : 'Unknown error'
        setError(`Failed to load chart data: ${errorMessage}`)
      } finally {
        setLoading(false)
      }
    }

    fetchChartData()
  }, [selectedPeriod])

  const periods = [
    { label: '7D', value: '7d' },
    { label: '30D', value: '30d' },
    { label: '90D', value: '90d' },
    { label: '1Y', value: '1y' },
  ]

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    if (selectedPeriod === '7d') {
      return date.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' })
    } else if (selectedPeriod === '30d') {
      return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
    } else if (selectedPeriod === '90d') {
      return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
    } else {
      return date.toLocaleDateString('en-US', { month: 'short', year: '2-digit' })
    }
  }

  const maxRevenue = Math.max(...chartData.map(d => d.revenue), 1)
  const maxOrders = Math.max(...chartData.map(d => d.orders), 1)

  const totalRevenue = chartData.reduce((sum, d) => sum + d.revenue, 0)
  const totalOrders = chartData.reduce((sum, d) => sum + d.orders, 0)
  const avgRevenue = chartData.length > 0 ? totalRevenue / chartData.length : 0

  if (loading) {
    return (
      <Card className="col-span-4">
        <CardHeader>
          <CardTitle>Sales Overview</CardTitle>
          <CardDescription>Loading sales data...</CardDescription>
        </CardHeader>
        <CardContent className="pl-2">
          <div className="h-[350px] flex items-center justify-center">
            <div className="text-center">
              <BarChart3 className="h-12 w-12 mx-auto mb-4 opacity-50 animate-pulse" />
              <p className="text-muted-foreground">Loading chart data...</p>
            </div>
          </div>
        </CardContent>
      </Card>
    )
  }

  if (error) {
    return (
      <Card className="col-span-4">
        <CardHeader>
          <CardTitle>Sales Overview</CardTitle>
          <CardDescription>Failed to load sales data</CardDescription>
        </CardHeader>
        <CardContent className="pl-2">
          <div className="h-[350px] flex items-center justify-center">
            <div className="text-center text-destructive">
              <BarChart3 className="h-12 w-12 mx-auto mb-4 opacity-50" />
              <p>{error}</p>
            </div>
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="w-full">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <BarChart3 className="h-5 w-5" />
              Sales Overview
            </CardTitle>
            <CardDescription>
              Revenue and orders over time
            </CardDescription>
          </div>
          {!period && (
            <div className="flex items-center space-x-2">
              {periods.map((periodOption) => (
                <Badge
                  key={periodOption.value}
                  variant={selectedPeriod === periodOption.value ? "default" : "outline"}
                  className="cursor-pointer"
                  onClick={() => setSelectedPeriod(periodOption.value)}
                >
                  {periodOption.label}
                </Badge>
              ))}
            </div>
          )}
        </div>
      </CardHeader>
      <CardContent className="p-6">
        {chartData.length === 0 ? (
          <div className="h-[350px] flex items-center justify-center text-muted-foreground">
            <div className="text-center">
              <BarChart3 className="h-12 w-12 mx-auto mb-4 opacity-50" />
              <p>No sales data available for the selected period</p>
            </div>
          </div>
        ) : (
          <div className="space-y-6">
            {/* Summary Stats */}
            <div className="grid grid-cols-3 gap-4">
              <div className="flex items-center space-x-3 p-4 rounded-lg bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-950/20 dark:to-emerald-950/20 border border-green-200 dark:border-green-800">
                <div className="p-2 bg-green-500 rounded-lg">
                  <DollarSign className="h-5 w-5 text-white" />
                </div>
                <div>
                  <p className="text-2xl font-bold text-green-700 dark:text-green-400">{formatCurrency(totalRevenue)}</p>
                  <p className="text-xs text-green-600 dark:text-green-500">Total Revenue</p>
                </div>
              </div>
              <div className="flex items-center space-x-3 p-4 rounded-lg bg-gradient-to-br from-blue-50 to-cyan-50 dark:from-blue-950/20 dark:to-cyan-950/20 border border-blue-200 dark:border-blue-800">
                <div className="p-2 bg-blue-500 rounded-lg">
                  <ShoppingCart className="h-5 w-5 text-white" />
                </div>
                <div>
                  <p className="text-2xl font-bold text-blue-700 dark:text-blue-400">{totalOrders}</p>
                  <p className="text-xs text-blue-600 dark:text-blue-500">Total Orders</p>
                </div>
              </div>
              <div className="flex items-center space-x-3 p-4 rounded-lg bg-gradient-to-br from-purple-50 to-pink-50 dark:from-purple-950/20 dark:to-pink-950/20 border border-purple-200 dark:border-purple-800">
                <div className="p-2 bg-purple-500 rounded-lg">
                  <TrendingUp className="h-5 w-5 text-white" />
                </div>
                <div>
                  <p className="text-2xl font-bold text-purple-700 dark:text-purple-400">{formatCurrency(avgRevenue)}</p>
                  <p className="text-xs text-purple-600 dark:text-purple-500">Avg Daily Revenue</p>
                </div>
              </div>
            </div>

            {/* Chart */}
            <div className="relative">
              <div className="h-[280px] w-full">
                <svg className="w-full h-full" viewBox="0 0 800 280" preserveAspectRatio="none">
                  <defs>
                    <linearGradient id="revenueGradient" x1="0" x2="0" y1="0" y2="1">
                      <stop offset="0%" stopColor="rgb(34, 197, 94)" stopOpacity="0.3" />
                      <stop offset="100%" stopColor="rgb(34, 197, 94)" stopOpacity="0.05" />
                    </linearGradient>
                    <linearGradient id="ordersGradient" x1="0" x2="0" y1="0" y2="1">
                      <stop offset="0%" stopColor="rgb(59, 130, 246)" stopOpacity="0.3" />
                      <stop offset="100%" stopColor="rgb(59, 130, 246)" stopOpacity="0.05" />
                    </linearGradient>
                  </defs>

                  {/* Grid lines */}
                  {[0, 1, 2, 3, 4].map((i) => (
                    <line
                      key={i}
                      x1="0"
                      y1={i * 70}
                      x2="800"
                      y2={i * 70}
                      stroke="currentColor"
                      strokeOpacity="0.1"
                      strokeWidth="1"
                    />
                  ))}

                  {/* Revenue Area Chart */}
                  <path
                    d={`M 0 280 ${chartData.map((point, i) => {
                      const x = (i / (chartData.length - 1)) * 800
                      const y = 280 - (point.revenue / maxRevenue) * 260
                      return `L ${x} ${y}`
                    }).join(' ')} L 800 280 Z`}
                    fill="url(#revenueGradient)"
                  />
                  <path
                    d={`M ${chartData.map((point, i) => {
                      const x = (i / (chartData.length - 1)) * 800
                      const y = 280 - (point.revenue / maxRevenue) * 260
                      return `${i === 0 ? '' : 'L '}${x} ${y}`
                    }).join(' ')}`}
                    fill="none"
                    stroke="rgb(34, 197, 94)"
                    strokeWidth="3"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />

                  {/* Data points */}
                  {chartData.map((point, i) => {
                    const x = (i / (chartData.length - 1)) * 800
                    const y = 280 - (point.revenue / maxRevenue) * 260
                    return (
                      <circle
                        key={i}
                        cx={x}
                        cy={y}
                        r={hoveredIndex === i ? "6" : "4"}
                        fill="rgb(34, 197, 94)"
                        stroke="white"
                        strokeWidth="2"
                        className="cursor-pointer transition-all"
                        onMouseEnter={() => setHoveredIndex(i)}
                        onMouseLeave={() => setHoveredIndex(null)}
                      />
                    )
                  })}
                </svg>
              </div>

              {/* Tooltip */}
              {hoveredIndex !== null && (
                <div className="absolute top-0 left-1/2 transform -translate-x-1/2 bg-popover border rounded-lg shadow-lg p-3 z-10">
                  <p className="text-xs font-medium mb-2">{formatDate(chartData[hoveredIndex].date)}</p>
                  <div className="space-y-1">
                    <div className="flex items-center gap-2">
                      <div className="w-3 h-3 bg-green-500 rounded-full" />
                      <span className="text-xs">Revenue: {formatCurrency(chartData[hoveredIndex].revenue)}</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <div className="w-3 h-3 bg-blue-500 rounded-full" />
                      <span className="text-xs">Orders: {chartData[hoveredIndex].orders}</span>
                    </div>
                  </div>
                </div>
              )}
            </div>

            {/* X-axis labels */}
            <div className="flex justify-between text-xs text-muted-foreground px-2">
              {chartData.length > 0 && (
                <>
                  <span>{formatDate(chartData[0].date)}</span>
                  {chartData.length > 2 && (
                    <span>{formatDate(chartData[Math.floor(chartData.length / 2)].date)}</span>
                  )}
                  <span>{formatDate(chartData[chartData.length - 1].date)}</span>
                </>
              )}
            </div>

            {/* Legend */}
            <div className="flex items-center justify-center space-x-6 text-sm pt-2 border-t">
              <div className="flex items-center space-x-2">
                <div className="w-4 h-4 bg-green-500 rounded" />
                <span>Revenue</span>
              </div>
              <div className="flex items-center space-x-2">
                <div className="w-4 h-4 bg-blue-500 rounded" />
                <span>Orders</span>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
