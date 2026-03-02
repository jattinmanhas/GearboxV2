"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { useUserStore } from "@/lib/stores/user-store"
import { orderApi } from "@/lib/apiFunctions"
import { formatCurrency } from "@/lib/currency"
import { Loader2, ShoppingBag, Search, Eye } from "lucide-react"

interface OrderSummary {
  id: number
  order_number: string
  total_amount: number
  status: string
  payment_status: string
  created_at: string
  currency: string
}

export default function MyOrdersPage() {
  const { isAuthenticated } = useUserStore()
  const [isMounted, setIsMounted] = useState(false)
  const [orders, setOrders] = useState<OrderSummary[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [page, setPage] = useState(1)
  const [limit] = useState(10)
  const [total, setTotal] = useState(0)
  const [totalPages, setTotalPages] = useState(1)
  const [search, setSearch] = useState("")
  const [status, setStatus] = useState("all")
  const [paymentStatus, setPaymentStatus] = useState("all")

  useEffect(() => {
    setIsMounted(true)
  }, [])

  useEffect(() => {
    if (!isMounted || !isAuthenticated) {
      setLoading(false)
      return
    }
    fetchOrders()
  }, [isMounted, isAuthenticated, page, status, paymentStatus])

  const fetchOrders = async () => {
    try {
      setLoading(true)
      setError(null)

      const params: Record<string, string> = {
        page: String(page),
        limit: String(limit),
        sort: "created_at",
        order: "desc",
      }

      if (search.trim()) params.search = search.trim()
      if (status !== "all") params.status = status
      if (paymentStatus !== "all") params.payment_status = paymentStatus

      const response = await orderApi.getOrders(params)
      if (response.success && response.data) {
        setOrders(response.data.orders || [])
        setTotal(response.data.total || 0)
        setTotalPages(response.data.total_pages || 1)
      } else {
        setOrders([])
        setError(response.message || "Failed to load your orders")
      }
    } catch (err) {
      setOrders([])
      setError(err instanceof Error ? err.message : "Failed to load your orders")
    } finally {
      setLoading(false)
    }
  }

  const formatDate = (dateStr: string) =>
    new Date(dateStr).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    })

  const badgeClass = (value: string) => {
    const v = (value || "").toLowerCase()
    if (v === "delivered" || v === "paid") return "bg-green-100 text-green-800"
    if (v === "processing" || v === "confirmed" || v === "shipped") return "bg-blue-100 text-blue-800"
    if (v === "pending") return "bg-orange-100 text-orange-800"
    if (v === "failed" || v === "cancelled" || v === "refunded") return "bg-red-100 text-red-800"
    return "bg-gray-100 text-gray-800"
  }

  if (!isMounted) return null

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto px-4 py-12">
          <Card>
            <CardContent className="py-12 text-center">
              <ShoppingBag className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
              <h2 className="text-xl font-semibold mb-2">Login required</h2>
              <p className="text-muted-foreground mb-6">Please login to view your orders.</p>
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
      <div className="container mx-auto px-4 py-8">
        <div className="mb-6">
          <h1 className="text-3xl font-bold">My Orders</h1>
          <p className="text-muted-foreground">Track your recent purchases and payment status.</p>
        </div>

        <Card className="mb-6">
          <CardHeader>
            <CardTitle>Filters</CardTitle>
          </CardHeader>
          <CardContent className="grid grid-cols-1 md:grid-cols-4 gap-3">
            <div className="flex gap-2 md:col-span-2">
              <Input
                placeholder="Search order number..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    setPage(1)
                    fetchOrders()
                  }
                }}
              />
              <Button
                variant="outline"
                onClick={() => {
                  setPage(1)
                  fetchOrders()
                }}
              >
                <Search className="h-4 w-4" />
              </Button>
            </div>
            <Select
              value={status}
              onValueChange={(value) => {
                setStatus(value)
                setPage(1)
              }}
            >
              <SelectTrigger>
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Statuses</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="confirmed">Confirmed</SelectItem>
                <SelectItem value="processing">Processing</SelectItem>
                <SelectItem value="shipped">Shipped</SelectItem>
                <SelectItem value="delivered">Delivered</SelectItem>
                <SelectItem value="cancelled">Cancelled</SelectItem>
              </SelectContent>
            </Select>
            <Select
              value={paymentStatus}
              onValueChange={(value) => {
                setPaymentStatus(value)
                setPage(1)
              }}
            >
              <SelectTrigger>
                <SelectValue placeholder="Payment Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Payments</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="paid">Paid</SelectItem>
                <SelectItem value="failed">Failed</SelectItem>
                <SelectItem value="refunded">Refunded</SelectItem>
                <SelectItem value="partially_refunded">Partially Refunded</SelectItem>
              </SelectContent>
            </Select>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Your Orders</CardTitle>
            <CardDescription>{total} orders found</CardDescription>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="py-12 text-center">
                <Loader2 className="h-6 w-6 animate-spin mx-auto mb-2" />
                <p className="text-muted-foreground">Loading orders...</p>
              </div>
            ) : error ? (
              <div className="text-destructive">{error}</div>
            ) : orders.length === 0 ? (
              <div className="py-12 text-center text-muted-foreground">
                <ShoppingBag className="h-12 w-12 mx-auto mb-4 opacity-60" />
                <p>No orders found.</p>
              </div>
            ) : (
              <div className="space-y-3">
                {orders.map((order) => (
                  <div key={order.id} className="border rounded-lg p-4 flex flex-col md:flex-row md:items-center md:justify-between gap-3">
                    <div>
                      <p className="font-semibold">{order.order_number}</p>
                      <p className="text-sm text-muted-foreground">{formatDate(order.created_at)}</p>
                      <div className="flex gap-2 mt-2">
                        <Badge className={badgeClass(order.status)}>{order.status}</Badge>
                        <Badge className={badgeClass(order.payment_status)}>
                          {order.payment_status.replaceAll("_", " ")}
                        </Badge>
                      </div>
                    </div>
                    <div className="flex items-center gap-3">
                      <p className="font-semibold">{formatCurrency(order.total_amount)}</p>
                      <Button variant="outline" size="sm" asChild>
                        <Link href={`/orders/${order.id}`}>
                          <Eye className="h-4 w-4 mr-1" />
                          View
                        </Link>
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            )}

            {totalPages > 1 && !loading && !error && (
              <div className="flex items-center justify-between mt-6">
                <Button variant="outline" size="sm" disabled={page === 1} onClick={() => setPage((p) => p - 1)}>
                  Previous
                </Button>
                <span className="text-sm text-muted-foreground">Page {page} of {totalPages}</span>
                <Button variant="outline" size="sm" disabled={page >= totalPages} onClick={() => setPage((p) => p + 1)}>
                  Next
                </Button>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
