"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { 
  Package, 
  AlertTriangle, 
  TrendingUp, 
  TrendingDown,
  RefreshCw,
  BarChart3,
  DollarSign,
  ShoppingCart,
  Warehouse
} from "lucide-react"
import { InventorySummary } from "@/lib/types"
import { formatCurrency } from "@/lib/currency"

export default function InventorySummaryPage() {
  const [summary, setSummary] = useState<InventorySummary | null>(null)
  const [loading, setLoading] = useState(true)

  // Fetch inventory summary
  const fetchSummary = async () => {
    try {
      setLoading(true)
      const response = await fetch('/api/v1/inventory/summary', {
        credentials: 'include'
      })
      const data = await response.json()
      
      if (response.ok) {
        setSummary(data.data)
      }
    } catch (error) {
      console.error('Error fetching inventory summary:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchSummary()
  }, [])

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Inventory Summary</h1>
            <p className="text-muted-foreground">
              Overview of your inventory analytics and metrics
            </p>
          </div>
        </div>
        <div className="text-center py-8">
          <RefreshCw className="h-6 w-6 animate-spin mx-auto mb-2" />
          Loading summary...
        </div>
      </div>
    )
  }

  if (!summary) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Inventory Summary</h1>
            <p className="text-muted-foreground">
              Overview of your inventory analytics and metrics
            </p>
          </div>
        </div>
        <div className="text-center py-8 text-muted-foreground">
          Unable to load inventory summary
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Inventory Summary</h1>
          <p className="text-muted-foreground">
            Overview of your inventory analytics and metrics
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchSummary}>
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </Button>
        </div>
      </div>

      {/* Key Metrics */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Products</CardTitle>
            <Package className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{summary.total_products.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">
              {summary.total_variants.toLocaleString()} variants
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Stock</CardTitle>
            <Warehouse className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{summary.total_quantity.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">
              {summary.total_available.toLocaleString()} available
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Reserved Stock</CardTitle>
            <ShoppingCart className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{summary.total_reserved.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">
              In pending orders
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Value</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(summary.total_value)}</div>
            <p className="text-xs text-muted-foreground">
              Avg: {formatCurrency(summary.average_stock_level)}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Stock Status Overview */}
      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-yellow-500" />
              Low Stock Items
            </CardTitle>
            <CardDescription>
              Products that are running low on inventory
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold text-yellow-600">{summary.low_stock_items}</div>
            <p className="text-sm text-muted-foreground mt-2">
              Items below minimum stock level
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingDown className="h-5 w-5 text-red-500" />
              Out of Stock Items
            </CardTitle>
            <CardDescription>
              Products with zero available inventory
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold text-red-600">{summary.out_of_stock_items}</div>
            <p className="text-sm text-muted-foreground mt-2">
              Items requiring immediate restocking
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Inventory Health */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <BarChart3 className="h-5 w-5" />
            Inventory Health
          </CardTitle>
          <CardDescription>
            Overall health metrics for your inventory
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-3">
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Stock Utilization</span>
                <span className="text-sm text-muted-foreground">
                  {summary.total_quantity > 0 ? 
                    Math.round((summary.total_available / summary.total_quantity) * 100) : 0}%
                </span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div 
                  className="bg-blue-600 h-2 rounded-full" 
                  style={{ 
                    width: `${summary.total_quantity > 0 ? 
                      Math.round((summary.total_available / summary.total_quantity) * 100) : 0}%` 
                  }}
                ></div>
              </div>
            </div>

            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Reservation Rate</span>
                <span className="text-sm text-muted-foreground">
                  {summary.total_quantity > 0 ? 
                    Math.round((summary.total_reserved / summary.total_quantity) * 100) : 0}%
                </span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div 
                  className="bg-orange-600 h-2 rounded-full" 
                  style={{ 
                    width: `${summary.total_quantity > 0 ? 
                      Math.round((summary.total_reserved / summary.total_quantity) * 100) : 0}%` 
                  }}
                ></div>
              </div>
            </div>

            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Alert Rate</span>
                <span className="text-sm text-muted-foreground">
                  {summary.total_products > 0 ? 
                    Math.round(((summary.low_stock_items + summary.out_of_stock_items) / summary.total_products) * 100) : 0}%
                </span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div 
                  className="bg-red-600 h-2 rounded-full" 
                  style={{ 
                    width: `${summary.total_products > 0 ? 
                      Math.round(((summary.low_stock_items + summary.out_of_stock_items) / summary.total_products) * 100) : 0}%` 
                  }}
                ></div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle>Quick Actions</CardTitle>
          <CardDescription>
            Common inventory management tasks
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-3">
            <div className="p-4 border rounded-lg">
              <h3 className="font-medium mb-2">Low Stock Items</h3>
              <p className="text-sm text-muted-foreground mb-3">
                {summary.low_stock_items} items need restocking
              </p>
              <Button size="sm" variant="outline" className="w-full">
                View Low Stock
              </Button>
            </div>
            
            <div className="p-4 border rounded-lg">
              <h3 className="font-medium mb-2">Out of Stock</h3>
              <p className="text-sm text-muted-foreground mb-3">
                {summary.out_of_stock_items} items are out of stock
              </p>
              <Button size="sm" variant="outline" className="w-full">
                View Out of Stock
              </Button>
            </div>
            
            <div className="p-4 border rounded-lg">
              <h3 className="font-medium mb-2">Inventory Value</h3>
              <p className="text-sm text-muted-foreground mb-3">
                Total value: {formatCurrency(summary.total_value)}
              </p>
              <Button size="sm" variant="outline" className="w-full">
                View Analytics
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
