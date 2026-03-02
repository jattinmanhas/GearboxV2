"use client"

import { useEffect, useMemo, useState } from "react"
import Link from "next/link"
import { useParams } from "next/navigation"
import { useUserStore } from "@/lib/stores/user-store"
import { orderApi } from "@/lib/apiFunctions"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { formatCurrency } from "@/lib/currency"
import { ArrowLeft, Loader2, Package } from "lucide-react"

interface OrderItem {
  id: number
  product_id: number
  product_name: string
  quantity: number
  unit_price: number
  total_price: number
}

interface OrderDetails {
  id: number
  order_number: string
  status: string
  payment_status: string
  fulfillment_status: string
  subtotal: number
  discount_amount: number
  total_amount: number
  currency: string
  created_at: string
  items: OrderItem[]
}

export default function OrderDetailsPage() {
  const params = useParams<{ id: string }>()
  const { isAuthenticated } = useUserStore()
  const [isMounted, setIsMounted] = useState(false)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [order, setOrder] = useState<OrderDetails | null>(null)

  const orderId = useMemo(() => Number(params?.id || 0), [params])

  useEffect(() => {
    setIsMounted(true)
  }, [])

  useEffect(() => {
    const loadOrder = async () => {
      if (!isMounted || !isAuthenticated) {
        setLoading(false)
        return
      }
      if (!orderId || Number.isNaN(orderId)) {
        setError("Invalid order ID")
        setLoading(false)
        return
      }

      try {
        setLoading(true)
        setError(null)
        const response = await orderApi.getOrder(orderId)
        if (!response) {
          setError("Order not found")
          setOrder(null)
          return
        }
        setOrder(response as OrderDetails)
      } catch (err) {
        setOrder(null)
        setError(err instanceof Error ? err.message : "Failed to load order")
      } finally {
        setLoading(false)
      }
    }

    loadOrder()
  }, [isMounted, isAuthenticated, orderId])

  const formatDate = (dateStr: string) =>
    new Date(dateStr).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    })

  if (!isMounted) return null

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto px-4 py-12">
          <Card>
            <CardContent className="py-12 text-center">
              <h2 className="text-xl font-semibold mb-2">Login required</h2>
              <p className="text-muted-foreground mb-6">Please login to view order details.</p>
              <Button asChild>
                <Link href="/login">Login</Link>
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto px-4 py-8 max-w-4xl">
        <div className="mb-6">
          <Button variant="ghost" asChild>
            <Link href="/orders">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Orders
            </Link>
          </Button>
        </div>

        {loading ? (
          <Card>
            <CardContent className="py-12 text-center">
              <Loader2 className="h-6 w-6 animate-spin mx-auto mb-2" />
              <p className="text-muted-foreground">Loading order details...</p>
            </CardContent>
          </Card>
        ) : error || !order ? (
          <Card>
            <CardContent className="py-12 text-center text-destructive">
              {error || "Order not found"}
            </CardContent>
          </Card>
        ) : (
          <Card>
            <CardHeader>
              <CardTitle>Order {order.order_number}</CardTitle>
              <CardDescription>Placed on {formatDate(order.created_at)}</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex flex-wrap gap-2">
                <Badge>{order.status}</Badge>
                <Badge variant="outline">{order.payment_status}</Badge>
                <Badge variant="secondary">{order.fulfillment_status}</Badge>
              </div>

              <div className="space-y-3">
                <h3 className="font-semibold flex items-center gap-2">
                  <Package className="h-4 w-4" />
                  Items
                </h3>
                {order.items?.length ? (
                  <div className="space-y-2">
                    {order.items.map((item) => (
                      <div key={item.id} className="flex items-center justify-between border rounded-md p-3">
                        <div>
                          <p className="font-medium">{item.product_name || `Product ${item.product_id}`}</p>
                          <p className="text-sm text-muted-foreground">
                            {formatCurrency(item.unit_price)} × {item.quantity}
                          </p>
                        </div>
                        <p className="font-semibold">{formatCurrency(item.total_price)}</p>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-muted-foreground text-sm">No items found.</p>
                )}
              </div>

              <Separator />

              <div className="space-y-2 max-w-sm ml-auto">
                <div className="flex justify-between">
                  <span>Subtotal</span>
                  <span>{formatCurrency(order.subtotal)}</span>
                </div>
                {order.discount_amount > 0 && (
                  <div className="flex justify-between text-green-600">
                    <span>Discount</span>
                    <span>-{formatCurrency(order.discount_amount)}</span>
                  </div>
                )}
                <Separator />
                <div className="flex justify-between text-lg font-bold">
                  <span>Total</span>
                  <span>{formatCurrency(order.total_amount)}</span>
                </div>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  )
}
