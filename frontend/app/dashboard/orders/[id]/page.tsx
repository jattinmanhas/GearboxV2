"use client"

import { useState, useEffect } from "react"
import { useParams, useRouter } from "next/navigation"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import {
    ShoppingCart,
    ArrowLeft,
    User,
    Mail,
    MapPin,
    Package,
    Truck,
    CreditCard,
    Calendar,
    Clock,
    AlertCircle,
    CheckCircle2
} from "lucide-react"
import { formatCurrency } from "@/lib/currency"
import { orderApi } from "@/lib/apiFunctions/order.api"
import { toast } from "sonner"
import Link from "next/link"

interface OrderItem {
    id: number
    product_id: number
    product_name: string
    sku: string
    quantity: number
    unit_price: number
    total_price: number
    variant_id?: number
}

interface OrderAddress {
    first_name: string
    last_name: string
    email: string
    phone: string
    address_line1: string
    address_line2?: string
    city: string
    state: string
    postal_code: string
    country: string
}

interface Order {
    id: number
    order_number: string
    customer_id?: number
    customer_name: string
    customer_email: string
    total_amount: number
    subtotal: number
    tax_amount: number
    shipping_amount: number
    discount_amount: number
    status: string
    payment_status: string
    fulfillment_status: string
    currency: string
    notes?: string
    created_at: string
    updated_at: string
    items: OrderItem[]
    billing_address: OrderAddress
    shipping_address: OrderAddress
}

export default function OrderDetailsPage() {
    const params = useParams()
    const router = useRouter()
    const id = params.id as string

    const [order, setOrder] = useState<Order | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const [updatingStatus, setUpdatingStatus] = useState(false)

    useEffect(() => {
        if (id) {
            fetchOrderDetails()
        }
    }, [id])

    const fetchOrderDetails = async () => {
        try {
            setLoading(true)
            setError(null)
            const orderId = parseInt(id)
            if (isNaN(orderId)) throw new Error("Invalid order ID")

            const data = await orderApi.getOrder(orderId)
            setOrder(data)
        } catch (err) {
            console.error('Error fetching order:', err)
            setError(err instanceof Error ? err.message : 'Failed to load order details')
        } finally {
            setLoading(false)
        }
    }

    const handleUpdateStatus = async (newStatus: string) => {
        if (!order) return
        try {
            setUpdatingStatus(true)
            await orderApi.updateOrderStatus(order.id, newStatus)
            toast.success(`Order status updated to ${newStatus}`)
            fetchOrderDetails()
        } catch (err) {
            toast.error("Failed to update order status")
        } finally {
            setUpdatingStatus(false)
        }
    }

    const handleUpdatePaymentStatus = async (newStatus: string) => {
        if (!order) return
        try {
            setUpdatingStatus(true)
            await orderApi.updatePaymentStatus(order.id, newStatus)
            toast.success(`Payment status updated to ${newStatus}`)
            fetchOrderDetails()
        } catch (err) {
            toast.error("Failed to update payment status")
        } finally {
            setUpdatingStatus(false)
        }
    }

    const getStatusColor = (status: string) => {
        switch (status.toLowerCase()) {
            case 'delivered':
            case 'completed':
                return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
            case 'processing':
            case 'confirmed':
                return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'
            case 'shipped':
                return 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200'
            case 'cancelled':
            case 'refunded':
                return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
            case 'pending':
                return 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200'
            default:
                return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200'
        }
    }

    const getPaymentStatusColor = (status: string) => {
        switch (status.toLowerCase()) {
            case 'paid':
                return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
            case 'failed':
                return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
            case 'refunded':
                return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
            case 'pending':
                return 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200'
            default:
                return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200'
        }
    }

    const formatDate = (dateString: string) => {
        if (!dateString) return 'N/A'
        return new Date(dateString).toLocaleString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        })
    }

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-[400px]">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                <span className="ml-3 text-muted-foreground">Loading order details...</span>
            </div>
        )
    }

    if (error || !order) {
        return (
            <Card className="border-destructive">
                <CardHeader>
                    <CardTitle className="text-destructive flex items-center gap-2">
                        <AlertCircle className="h-5 w-5" />
                        Error Loading Order
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <p>{error || "Order not found"}</p>
                    <Button variant="outline" className="mt-4" onClick={() => router.push("/dashboard/orders")}>
                        Back to Orders
                    </Button>
                </CardContent>
            </Card>
        )
    }

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                    <Button variant="outline" size="icon" onClick={() => router.back()}>
                        <ArrowLeft className="h-4 w-4" />
                    </Button>
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight">Order {order.order_number}</h1>
                        <p className="text-muted-foreground flex items-center gap-2">
                            Placed on {formatDate(order.created_at)}
                        </p>
                    </div>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline" asChild>
                        <Link href="/dashboard/orders">
                            List Orders
                        </Link>
                    </Button>
                    <Button onClick={fetchOrderDetails} variant="ghost" size="icon">
                        <Clock className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
                    </Button>
                </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* Main Content - Left Side */}
                <div className="lg:col-span-2 space-y-6">
                    {/* Order Items */}
                    <Card>
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2">
                                <Package className="h-5 w-5" />
                                Order Items
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-4">
                                {order.items && order.items.map((item) => (
                                    <div key={item.id} className="flex items-center justify-between py-2 border-b last:border-0">
                                        <div className="flex gap-4">
                                            <div className="bg-muted w-12 h-12 rounded-md flex items-center justify-center">
                                                <ShoppingCart className="h-6 w-6 text-muted-foreground" />
                                            </div>
                                            <div>
                                                <h4 className="font-medium">{item.product_name}</h4>
                                                <p className="text-sm text-muted-foreground">SKU: {item.sku}</p>
                                                <p className="text-sm">Qty: {item.quantity}</p>
                                            </div>
                                        </div>
                                        <div className="text-right">
                                            <p className="font-medium">{formatCurrency(item.total_price)}</p>
                                            <p className="text-xs text-muted-foreground">{formatCurrency(item.unit_price)} each</p>
                                        </div>
                                    </div>
                                ))}
                                {!order.items?.length && <p className="text-muted-foreground italic">No items found in this order.</p>}
                            </div>

                            <div className="mt-6 space-y-2">
                                <div className="flex justify-between text-sm">
                                    <span>Subtotal</span>
                                    <span>{formatCurrency(order.subtotal)}</span>
                                </div>
                                <div className="flex justify-between text-sm">
                                    <span>Shipping</span>
                                    <span>{formatCurrency(order.shipping_amount)}</span>
                                </div>
                                <div className="flex justify-between text-sm">
                                    <span>Tax</span>
                                    <span>{formatCurrency(order.tax_amount)}</span>
                                </div>
                                {order.discount_amount > 0 && (
                                    <div className="flex justify-between text-sm text-green-600 font-medium">
                                        <span>Discount</span>
                                        <span>-{formatCurrency(order.discount_amount)}</span>
                                    </div>
                                )}
                                <Separator className="my-2" />
                                <div className="flex justify-between font-bold text-lg">
                                    <span>Total</span>
                                    <span>{formatCurrency(order.total_amount)}</span>
                                </div>
                            </div>
                        </CardContent>
                    </Card>

                    {/* Addresses */}
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <Card>
                            <CardHeader>
                                <CardTitle className="text-sm font-medium flex items-center gap-2">
                                    <Truck className="h-4 w-4" />
                                    Shipping Address
                                </CardTitle>
                            </CardHeader>
                            <CardContent className="text-sm space-y-1">
                                <p className="font-medium">{order.shipping_address.first_name} {order.shipping_address.last_name}</p>
                                <p>{order.shipping_address.address_line1}</p>
                                {order.shipping_address.address_line2 && <p>{order.shipping_address.address_line2}</p>}
                                <p>{order.shipping_address.city}, {order.shipping_address.state} {order.shipping_address.postal_code}</p>
                                <p>{order.shipping_address.country}</p>
                                <div className="pt-2 flex items-center gap-2 text-muted-foreground">
                                    <Mail className="h-3 w-3" /> {order.shipping_address.email}
                                </div>
                            </CardContent>
                        </Card>

                        <Card>
                            <CardHeader>
                                <CardTitle className="text-sm font-medium flex items-center gap-2">
                                    <CreditCard className="h-4 w-4" />
                                    Billing Address
                                </CardTitle>
                            </CardHeader>
                            <CardContent className="text-sm space-y-1">
                                <p className="font-medium">{order.billing_address.first_name} {order.billing_address.last_name}</p>
                                <p>{order.billing_address.address_line1}</p>
                                {order.billing_address.address_line2 && <p>{order.billing_address.address_line2}</p>}
                                <p>{order.billing_address.city}, {order.billing_address.state} {order.billing_address.postal_code}</p>
                                <p>{order.billing_address.country}</p>
                            </CardContent>
                        </Card>
                    </div>
                </div>

                {/* Sidebar - Right Side */}
                <div className="space-y-6">
                    {/* Status Management */}
                    <Card>
                        <CardHeader>
                            <CardTitle className="text-lg">Order Status</CardTitle>
                            <CardDescription>Manage the lifecycle of this order</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="space-y-2">
                                <label className="text-sm font-medium">General Status</label>
                                <Select
                                    value={order.status}
                                    onValueChange={handleUpdateStatus}
                                    disabled={updatingStatus}
                                >
                                    <SelectTrigger className={getStatusColor(order.status)}>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="pending">Pending</SelectItem>
                                        <SelectItem value="confirmed">Confirmed</SelectItem>
                                        <SelectItem value="processing">Processing</SelectItem>
                                        <SelectItem value="shipped">Shipped</SelectItem>
                                        <SelectItem value="delivered">Delivered</SelectItem>
                                        <SelectItem value="cancelled">Cancelled</SelectItem>
                                        <SelectItem value="refunded">Refunded</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>

                            <div className="space-y-2">
                                <label className="text-sm font-medium">Payment Status</label>
                                <Select
                                    value={order.payment_status}
                                    onValueChange={handleUpdatePaymentStatus}
                                    disabled={updatingStatus}
                                >
                                    <SelectTrigger className={getPaymentStatusColor(order.payment_status)}>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="pending">Pending</SelectItem>
                                        <SelectItem value="paid">Paid</SelectItem>
                                        <SelectItem value="failed">Failed</SelectItem>
                                        <SelectItem value="refunded">Refunded</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>
                        </CardContent>
                    </Card>

                    {/* Customer Info */}
                    <Card>
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2">
                                <User className="h-5 w-5" />
                                Customer
                            </CardTitle>
                        </CardHeader>
                        <CardContent className="text-sm space-y-4">
                            <div>
                                <p className="font-medium">{order.customer_name || 'Anonymous'}</p>
                                <p className="text-muted-foreground">{order.customer_email}</p>
                            </div>
                            {order.customer_id && (
                                <Button variant="outline" size="sm" className="w-full" asChild>
                                    <Link href={`/dashboard/users/${order.customer_id}`}>
                                        View Customer Profile
                                    </Link>
                                </Button>
                            )}
                        </CardContent>
                    </Card>

                    {/* Additional Notes */}
                    <Card>
                        <CardHeader>
                            <CardTitle className="text-sm font-medium">Notes</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <p className="text-sm text-muted-foreground italic">
                                {order.notes || "No notes provided for this order."}
                            </p>
                        </CardContent>
                    </Card>
                </div>
            </div>
        </div>
    )
}
