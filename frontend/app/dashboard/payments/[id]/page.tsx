"use client"

import { useState, useEffect } from "react"
import { useParams, useRouter } from "next/navigation"
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { formatCurrency } from "@/lib/currency"
import {
    CreditCard,
    ArrowLeft,
    RefreshCcw,
    Calendar,
    ExternalLink,
    ShieldCheck,
    AlertCircle,
    FileText,
    DollarSign
} from "lucide-react"
import Link from "next/link"
import { toast } from "sonner"
import { paymentApi } from "@/lib/apiFunctions/payment.api"

interface Payment {
    id: number
    order_id: number
    payment_method_id: number
    transaction_id: string
    gateway_id: string
    amount: number
    currency: string
    status: string
    gateway_status: string
    gateway_response: string
    failure_reason: string
    processed_at: string
    created_at: string
    updated_at: string
    metadata?: any
}

export default function PaymentDetailsPage() {
    const params = useParams()
    const router = useRouter()
    const id = params.id

    const [payment, setPayment] = useState<Payment | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const [refundLoading, setRefundLoading] = useState(false)

    useEffect(() => {
        fetchPaymentDetails()
    }, [id])

    const fetchPaymentDetails = async () => {
        try {
            setLoading(true)
            const response = await paymentApi.getPayment(Number(id))
            if (response.error) throw new Error(response.message || "Failed to fetch payment details")
            setPayment(response)
        } catch (err) {
            console.error(err)
            const msg = err instanceof Error ? err.message : "Failed to load payment details"
            setError(msg)
            toast.error(msg)
        } finally {
            setLoading(false)
        }
    }

    const handleRefund = async () => {
        if (!payment) return

        // Simple verification
        if (!confirm(`Are you sure you want to refund ${formatCurrency(payment.amount)}?`)) return

        try {
            setRefundLoading(true)
            const response = await paymentApi.refundPayment({
                payment_id: payment.id,
                amount: payment.amount,
                reason: "Admin requested refund"
            })

            if (!response.error) {
                toast.success("Refund processed successfully")
                fetchPaymentDetails()
            } else {
                toast.error(response.message || "Failed to process refund")
            }
        } catch (err) {
            console.error(err)
            toast.error("An error occurred while processing refund")
        } finally {
            setRefundLoading(false)
        }
    }

    const getStatusColor = (status: string) => {
        switch (status.toLowerCase()) {
            case 'completed':
            case 'processed':
                return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
            case 'failed':
                return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
            case 'pending':
            case 'processing':
                return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'
            case 'refunded':
                return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
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

    if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Loading details...</div>

    if (error || !payment) {
        return (
            <div className="space-y-4">
                <Button variant="ghost" onClick={() => router.back()}>
                    <ArrowLeft className="h-4 w-4 mr-2" /> Back
                </Button>
                <div className="p-8 text-center bg-destructive/10 text-destructive rounded-lg border border-destructive/20">
                    <AlertCircle className="h-12 w-12 mx-auto mb-4 opacity-50" />
                    <h2 className="text-xl font-semibold">{error || "Payment not found"}</h2>
                </div>
            </div>
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
                        <h1 className="text-3xl font-bold tracking-tight">Payment Details</h1>
                        <div className="text-muted-foreground flex items-center gap-2">
                            <span className="font-mono text-xs">{payment.transaction_id}</span>
                            <Badge className={getStatusColor(payment.status)}>{payment.status}</Badge>
                        </div>
                    </div>
                </div>
                <div className="flex gap-2">
                    {payment.status === 'completed' && (
                        <Button variant="destructive" onClick={handleRefund} disabled={refundLoading}>
                            <RefreshCcw className={`h-4 w-4 mr-2 ${refundLoading ? 'animate-spin' : ''}`} />
                            Refund Payment
                        </Button>
                    )}
                    <Button variant="outline" asChild>
                        <Link href={`/dashboard/orders/${payment.order_id}`}>
                            <FileText className="h-4 w-4 mr-2" />
                            View Order
                        </Link>
                    </Button>
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {/* Main Info */}
                <Card className="md:col-span-2">
                    <CardHeader>
                        <CardTitle>Transaction Information</CardTitle>
                        <CardDescription>Full details of the payment event</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-6">
                        <div className="grid grid-cols-2 gap-4">
                            <div className="space-y-1">
                                <span className="text-sm text-muted-foreground">Amount</span>
                                <div className="text-lg font-bold flex items-center gap-1">
                                    <DollarSign className="h-4 w-4 text-green-600" />
                                    {formatCurrency(payment.amount)} {payment.currency}
                                </div>
                            </div>
                            <div className="space-y-1">
                                <span className="text-sm text-muted-foreground">Gateway</span>
                                <div className="flex items-center gap-2">
                                    <ShieldCheck className="h-4 w-4 text-blue-500" />
                                    <span className="capitalize">{payment.gateway_id}</span>
                                </div>
                            </div>
                            <div className="space-y-1">
                                <span className="text-sm text-muted-foreground">Created At</span>
                                <div className="flex items-center gap-2 text-sm">
                                    <Calendar className="h-4 w-4 text-muted-foreground" />
                                    {formatDate(payment.created_at)}
                                </div>
                            </div>
                            <div className="space-y-1">
                                <span className="text-sm text-muted-foreground">Processed At</span>
                                <div className="flex items-center gap-2 text-sm">
                                    <Calendar className="h-4 w-4 text-muted-foreground" />
                                    {formatDate(payment.processed_at)}
                                </div>
                            </div>
                        </div>

                        <div className="pt-4 border-t space-y-4">
                            <div className="space-y-1">
                                <span className="text-sm text-muted-foreground">Gateway Transaction Status</span>
                                <div className="text-sm font-mono p-2 bg-muted rounded">{payment.gateway_status || 'N/A'}</div>
                            </div>
                            {payment.failure_reason && (
                                <div className="space-y-1">
                                    <span className="text-sm text-muted-foreground text-red-500">Failure Reason</span>
                                    <div className="text-sm p-2 bg-red-50 dark:bg-red-950/20 text-red-600 rounded border border-red-100 dark:border-red-900/30">
                                        {payment.failure_reason}
                                    </div>
                                </div>
                            )}
                        </div>
                    </CardContent>
                </Card>

                {/* Sidebar Mini Cards */}
                <div className="space-y-6">
                    <Card>
                        <CardHeader>
                            <CardTitle className="text-sm font-medium">Customer & Order</CardTitle>
                        </CardHeader>
                        <CardContent className="text-sm space-y-2">
                            <div className="flex justify-between">
                                <span className="text-muted-foreground">Order ID</span>
                                <Link href={`/dashboard/orders/${payment.order_id}`} className="text-blue-600 hover:underline">
                                    #{payment.order_id}
                                </Link>
                            </div>
                            <div className="flex justify-between">
                                <span className="text-muted-foreground">Method ID</span>
                                <span>{payment.payment_method_id}</span>
                            </div>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader>
                            <CardTitle className="text-sm font-medium">Metadata</CardTitle>
                        </CardHeader>
                        <CardContent>
                            {payment.metadata ? (
                                <pre className="text-[10px] bg-muted p-2 rounded overflow-auto max-h-40">
                                    {JSON.stringify(payment.metadata, null, 2)}
                                </pre>
                            ) : (
                                <p className="text-xs text-muted-foreground italic">No metadata available</p>
                            )}
                        </CardContent>
                    </Card>
                </div>
            </div>

            {/* Gateway Response Section */}
            <Card>
                <CardHeader>
                    <CardTitle className="text-lg">Gateway Raw Response</CardTitle>
                </CardHeader>
                <CardContent>
                    <pre className="text-xs p-4 bg-black text-green-400 rounded-lg overflow-auto max-h-64 font-mono">
                        {payment.gateway_response ? (
                            (() => {
                                try {
                                    return JSON.stringify(JSON.parse(payment.gateway_response), null, 2)
                                } catch {
                                    return payment.gateway_response
                                }
                            })()
                        ) : (
                            "No raw response available"
                        )}
                    </pre>
                </CardContent>
            </Card>
        </div>
    )
}
