import { useEffect, useState } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  ShoppingCart,
  Eye,
  Clock,
  CheckCircle,
  XCircle,
  Truck,
  Package
} from "lucide-react"
import { formatCurrency } from "@/lib/currency"
import { httpClient } from "@/lib/apiFunctions/http-client"
import Link from "next/link"

interface RecentOrder {
  id: number
  order_number: string
  customer_name: string
  customer_email: string
  total_amount: number
  status: string
  created_at: string
}

export function RecentOrders() {
  const [orders, setOrders] = useState<RecentOrder[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchRecentOrders = async () => {
      try {
        setLoading(true)
        const result = await httpClient.get<{
          success: boolean
          message: string
          data: { orders: RecentOrder[] }
        }>('/dashboard/recent-orders?limit=5&sort=created_at&order=desc')

        if (result.success) {
          setOrders(result.data.orders || [])
        } else {
          throw new Error(result.message || 'Failed to load orders')
        }
      } catch (error) {
        console.error('Error loading recent orders:', error)
        const errorMessage = error instanceof Error ? error.message : 'Unknown error'
        console.error('Full error details:', {
          error,
          message: errorMessage,
          stack: error instanceof Error ? error.stack : undefined
        })
        setError(`Failed to load recent orders: ${errorMessage}`)
      } finally {
        setLoading(false)
      }
    }

    fetchRecentOrders()
  }, [])

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'delivered':
        return <CheckCircle className="h-4 w-4 text-green-600" />
      case 'processing':
        return <Package className="h-4 w-4 text-blue-600" />
      case 'shipped':
        return <Truck className="h-4 w-4 text-purple-600" />
      case 'cancelled':
        return <XCircle className="h-4 w-4 text-red-600" />
      default:
        return <Clock className="h-4 w-4 text-orange-600" />
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'delivered':
        return 'bg-green-100 text-green-800'
      case 'processing':
        return 'bg-blue-100 text-blue-800'
      case 'shipped':
        return 'bg-purple-100 text-purple-800'
      case 'cancelled':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-orange-100 text-orange-800'
    }
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  if (loading) {
    return (
      <Card className="w-full">
        <CardHeader>
          <CardTitle>Recent Orders</CardTitle>
          <CardDescription>Loading recent orders...</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="flex flex-col p-4 border rounded-lg">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center space-x-2">
                    <div className="h-4 w-4 bg-muted animate-pulse rounded"></div>
                    <div className="h-5 bg-muted animate-pulse rounded w-16"></div>
                  </div>
                  <div className="h-8 w-8 bg-muted animate-pulse rounded"></div>
                </div>
                <div className="space-y-2">
                  <div className="h-4 bg-muted animate-pulse rounded w-3/4"></div>
                  <div className="h-3 bg-muted animate-pulse rounded w-1/2"></div>
                  <div className="h-3 bg-muted animate-pulse rounded w-1/3"></div>
                  <div className="h-3 bg-muted animate-pulse rounded w-1/2"></div>
                  <div className="pt-2 border-t">
                    <div className="h-4 bg-muted animate-pulse rounded w-1/2"></div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    )
  }

  if (error) {
    return (
      <Card className="w-full">
        <CardHeader>
          <CardTitle>Recent Orders</CardTitle>
          <CardDescription>Failed to load recent orders</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-center text-destructive py-8">
            <p>{error}</p>
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
              <ShoppingCart className="h-5 w-5" />
              Recent Orders
            </CardTitle>
            <CardDescription>
              Latest {orders.length} orders from your customers
            </CardDescription>
          </div>
          <Button variant="outline" size="sm" asChild>
            <Link href="/dashboard/orders">View All</Link>
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        {orders.length === 0 ? (
          <div className="text-center py-12 text-muted-foreground">
            <ShoppingCart className="h-12 w-12 mx-auto mb-4 opacity-50" />
            <p>No recent orders found</p>
          </div>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5">
            {orders.map((order) => (
              <div key={order.id} className="flex flex-col p-4 border rounded-lg hover:bg-muted/50 transition-colors">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center space-x-2">
                    {getStatusIcon(order.status)}
                    <Badge className={getStatusColor(order.status)}>
                      {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                    </Badge>
                  </div>
                  <Button variant="ghost" size="sm" asChild>
                    <Link href={`/dashboard/orders/${order.id}`}>
                      <Eye className="h-4 w-4" />
                    </Link>
                  </Button>
                </div>

                <div className="space-y-2">
                  <div>
                    <p className="text-sm font-medium truncate">
                      {order.customer_name || 'Unknown Customer'}
                    </p>
                    <p className="text-xs text-muted-foreground truncate">
                      {order.customer_email || 'No email'}
                    </p>
                  </div>

                  <div className="text-xs text-muted-foreground">
                    #{order.order_number}
                  </div>

                  <div className="text-xs text-muted-foreground">
                    {formatDate(order.created_at)}
                  </div>

                  <div className="pt-2 border-t">
                    <p className="text-sm font-semibold">
                      {formatCurrency(order.total_amount)}
                    </p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {orders.length > 0 && (
          <div className="mt-6 pt-4 border-t">
            <Button variant="outline" className="w-full" asChild>
              <Link href="/dashboard/orders">View All Orders</Link>
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  )
}

