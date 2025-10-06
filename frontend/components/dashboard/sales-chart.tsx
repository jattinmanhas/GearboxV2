"use client"

import { useEffect, useState } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { 
  TrendingUp, 
  BarChart3,
  Calendar,
  DollarSign
} from "lucide-react"
import { formatCurrency } from "@/lib/currency"

interface ChartDataPoint {
  date: string
  revenue: number
  orders: number
}

export function SalesChart() {
  const [chartData, setChartData] = useState<ChartDataPoint[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedPeriod, setSelectedPeriod] = useState('30d')
  const [displayedData, setDisplayedData] = useState<ChartDataPoint[]>([])
  const [currentPage, setCurrentPage] = useState(0)

  useEffect(() => {
    const fetchChartData = async () => {
      try {
        setLoading(true)
        
        const response = await fetch(`/api/v1/dashboard/sales-chart?period=${selectedPeriod}`, {
          credentials: 'include',
        })
        
        if (!response.ok) {
          throw new Error(`Failed to fetch chart data: ${response.statusText}`)
        }
        
        const result = await response.json()
        
        if (result.success) {
          const data = result.data || []
          setChartData(data)
          
          // For very long periods, show only a subset of data
          if (data.length > 60) {
            setDisplayedData(data.slice(0, 60))
          } else {
            setDisplayedData(data)
          }
          setCurrentPage(0)
        } else {
          throw new Error(result.message || 'Failed to load chart data')
        }
      } catch (error) {
        console.error('Error loading chart data:', error)
        const errorMessage = error instanceof Error ? error.message : 'Unknown error'
        console.error('Full error details:', {
          error,
          message: errorMessage,
          stack: error instanceof Error ? error.stack : undefined
        })
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
      return date.toLocaleDateString('en-US', { weekday: 'short' })
    } else if (selectedPeriod === '30d') {
      return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
    } else if (selectedPeriod === '90d') {
      return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
    } else {
      return date.toLocaleDateString('en-US', { month: 'short' })
    }
  }

  const formatDateCompact = (dateString: string) => {
    const date = new Date(dateString)
    if (selectedPeriod === '7d') {
      return date.toLocaleDateString('en-US', { weekday: 'short' })
    } else if (selectedPeriod === '30d') {
      return `${date.getMonth() + 1}/${date.getDate()}`
    } else if (selectedPeriod === '90d') {
      return `${date.getMonth() + 1}/${date.getDate()}`
    } else {
      return date.toLocaleDateString('en-US', { month: 'short' })
    }
  }

  const maxRevenue = Math.max(...displayedData.map(d => d.revenue), 1)
  const maxOrders = Math.max(...displayedData.map(d => d.orders), 1)

  const totalRevenue = chartData.reduce((sum, d) => sum + d.revenue, 0)
  const totalOrders = chartData.reduce((sum, d) => sum + d.orders, 0)

  const itemsPerPage = 60
  const totalPages = Math.ceil(chartData.length / itemsPerPage)
  const hasPagination = chartData.length > itemsPerPage

  const handlePageChange = (page: number) => {
    const start = page * itemsPerPage
    const end = start + itemsPerPage
    setDisplayedData(chartData.slice(start, end))
    setCurrentPage(page)
  }

  if (loading) {
    return (
      <Card className="col-span-4">
        <CardHeader>
          <CardTitle>Sales Overview</CardTitle>
          <CardDescription>Loading sales data...</CardDescription>
        </CardHeader>
        <CardContent className="pl-2">
          <div className="h-[300px] flex items-center justify-center">
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
          <div className="h-[300px] flex items-center justify-center">
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
          <div className="flex items-center space-x-2">
            {periods.map((period) => (
              <Badge
                key={period.value}
                variant={selectedPeriod === period.value ? "default" : "outline"}
                className="cursor-pointer"
                onClick={() => setSelectedPeriod(period.value)}
              >
                {period.label}
              </Badge>
            ))}
          </div>
        </div>
      </CardHeader>
      <CardContent className="p-6">
        {displayedData.length === 0 ? (
          <div className="h-[300px] flex items-center justify-center text-muted-foreground">
            <div className="text-center">
              <BarChart3 className="h-12 w-12 mx-auto mb-4 opacity-50" />
              <p>No sales data available for the selected period</p>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            {/* Summary Stats */}
            <div className="grid grid-cols-2 gap-4">
              <div className="flex items-center space-x-2">
                <DollarSign className="h-4 w-4 text-green-600" />
                <div>
                  <p className="text-sm font-medium">{formatCurrency(totalRevenue)}</p>
                  <p className="text-xs text-muted-foreground">Total Revenue</p>
                </div>
              </div>
              <div className="flex items-center space-x-2">
                <TrendingUp className="h-4 w-4 text-blue-600" />
                <div>
                  <p className="text-sm font-medium">{totalOrders}</p>
                  <p className="text-xs text-muted-foreground">Total Orders</p>
                </div>
              </div>
            </div>

            {/* Chart Container with proper overflow handling */}
            <div className="w-full overflow-x-auto overflow-y-hidden scrollbar-thin scrollbar-thumb-gray-300 scrollbar-track-gray-100">
              <div className={`h-[200px] flex items-end px-2 ${
                displayedData.length > 60 ? 'space-x-0' : 
                displayedData.length > 30 ? 'space-x-0.5' : 
                displayedData.length > 15 ? 'space-x-1' : 'space-x-2'
              } ${displayedData.length <= 30 ? 'w-full' : ''}`}>
                {displayedData.map((data, index) => (
                  <div key={index} className={`flex flex-col items-center space-y-1 ${
                    displayedData.length > 60 ? 'w-2' :
                    displayedData.length > 30 ? 'w-3' :
                    displayedData.length > 15 ? 'w-4' : 'flex-1'
                  } min-w-0`}>
                    <div className="w-full flex flex-col space-y-1">
                      {/* Revenue Bar */}
                      <div
                        className="bg-green-500 rounded-t-sm transition-all duration-300 hover:bg-green-600"
                        style={{
                          height: `${(data.revenue / maxRevenue) * 120}px`,
                          minHeight: '2px'
                        }}
                        title={`${formatDate(data.date)}: ${formatCurrency(data.revenue)}`}
                      />
                      {/* Orders Bar */}
                      <div
                        className="bg-blue-500 rounded-b-sm transition-all duration-300 hover:bg-blue-600"
                        style={{
                          height: `${(data.orders / maxOrders) * 80}px`,
                          minHeight: '2px'
                        }}
                        title={`${formatDate(data.date)}: ${data.orders} orders`}
                      />
                    </div>
                    {/* Date labels - more compact and properly contained */}
                    <div className="text-xs text-muted-foreground text-center leading-tight max-w-full overflow-hidden h-8 flex items-end justify-center">
                      <div className="transform -rotate-45 origin-center whitespace-nowrap max-w-full">
                        {formatDateCompact(data.date)}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
            
            {/* Pagination Controls for long periods */}
            {hasPagination && (
              <div className="flex items-center justify-between">
                <div className="text-sm text-muted-foreground">
                  Showing {currentPage * itemsPerPage + 1}-{Math.min((currentPage + 1) * itemsPerPage, chartData.length)} of {chartData.length} data points
                </div>
                <div className="flex items-center space-x-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handlePageChange(currentPage - 1)}
                    disabled={currentPage === 0}
                  >
                    Previous
                  </Button>
                  <span className="text-sm text-muted-foreground">
                    Page {currentPage + 1} of {totalPages}
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handlePageChange(currentPage + 1)}
                    disabled={currentPage === totalPages - 1}
                  >
                    Next
                  </Button>
                </div>
              </div>
            )}

            {/* Legend */}
            <div className="flex items-center justify-center space-x-4 text-xs">
              <div className="flex items-center space-x-1">
                <div className="w-3 h-3 bg-green-500 rounded"></div>
                <span>Revenue</span>
              </div>
              <div className="flex items-center space-x-1">
                <div className="w-3 h-3 bg-blue-500 rounded"></div>
                <span>Orders</span>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
