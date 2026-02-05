"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import {
    CreditCard,
    Eye,
    Search,
    Filter,
    Download,
    RefreshCw,
    TrendingUp,
    AlertCircle,
    CheckCircle2,
    Clock
} from "lucide-react"
import { formatCurrency } from "@/lib/currency"
import Link from "next/link"

interface Payment {
    id: number
    order_id: number
    transaction_id: string
    gateway_id: string
    amount: number
    currency: string
    status: string
    gateway_status: string
    created_at: string
}

interface Summary {
    total_payments: number
    successful_payments: number
    failed_payments: number
    pending_payments: number
    total_amount: number
    refunded_amount: number
    net_amount: number
}

export default function PaymentsPage() {
    const [payments, setPayments] = useState<Payment[]>([])
    const [summary, setSummary] = useState<Summary | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const [page, setPage] = useState(1)
    const [limit, setLimit] = useState(10)
    const [total, setTotal] = useState(0)
    const [totalPages, setTotalPages] = useState(0)

    // Filters
    const [search, setSearch] = useState("")
    const [statusFilter, setStatusFilter] = useState<string>("all")
    const [gatewayFilter, setGatewayFilter] = useState<string>("all")
    const [sortBy, setSortBy] = useState("created_at")
    const [sortOrder, setSortOrder] = useState("desc")

    useEffect(() => {
        fetchPayments()
        fetchSummary()
    }, [page, limit, statusFilter, gatewayFilter, sortBy, sortOrder])

    const fetchPayments = async () => {
        try {
            setLoading(true)
            setError(null)

            const params = new URLSearchParams({
                page: page.toString(),
                limit: limit.toString(),
                sort_by: sortBy,
                sort_order: sortOrder,
            })

            if (statusFilter !== "all") {
                params.append("status", statusFilter)
            }
            if (gatewayFilter !== "all") {
                params.append("gateway_id", gatewayFilter)
            }
            if (search) {
                params.append("search", search)
            }

            const response = await fetch(`/api/v1/payments?${params.toString()}`, {
                credentials: 'include',
            })

            if (!response.ok) {
                throw new Error(`Failed to fetch payments: ${response.statusText}`)
            }

            const result = await response.json()

            // The backend returns { payments: [], total: 0, page: 1, limit: 10, pages: 1 }
            setPayments(result.payments || [])
            setTotal(result.total || 0)
            setTotalPages(result.pages || 0)
        } catch (error) {
            console.error('Error loading payments:', error)
            setError(error instanceof Error ? error.message : 'Unknown error')
        } finally {
            setLoading(false)
        }
    }

    const fetchSummary = async () => {
        try {
            const response = await fetch('/api/v1/payments/summary', {
                credentials: 'include',
            })
            if (response.ok) {
                const result = await response.json()
                setSummary(result)
            }
        } catch (error) {
            console.error('Error loading summary:', error)
        }
    }

    const handleSearch = () => {
        setPage(1)
        fetchPayments()
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
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        })
    }

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Payments</h1>
                    <p className="text-muted-foreground">
                        Monitor and manage all payment transactions
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline" onClick={() => { fetchPayments(); fetchSummary(); }}>
                        <RefreshCw className="h-4 w-4 mr-2" />
                        Refresh
                    </Button>
                    <Button variant="outline">
                        <Download className="h-4 w-4 mr-2" />
                        Export
                    </Button>
                </div>
            </div>

            {/* Summary Cards */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Net Revenue</CardTitle>
                        <TrendingUp className="h-4 w-4 text-green-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{formatCurrency(summary?.net_amount || 0)}</div>
                        <p className="text-xs text-muted-foreground">
                            Total after refunds
                        </p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Successful</CardTitle>
                        <CheckCircle2 className="h-4 w-4 text-green-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{summary?.successful_payments || 0}</div>
                        <p className="text-xs text-muted-foreground">
                            {summary?.total_payments ? ((summary.successful_payments / summary.total_payments) * 100).toFixed(1) : 0}% success rate
                        </p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Pending</CardTitle>
                        <Clock className="h-4 w-4 text-blue-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{summary?.pending_payments || 0}</div>
                        <p className="text-xs text-muted-foreground">
                            Awaiting confirmation
                        </p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Failed</CardTitle>
                        <AlertCircle className="h-4 w-4 text-red-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{summary?.failed_payments || 0}</div>
                        <p className="text-xs text-muted-foreground">
                            Requiring attention
                        </p>
                    </CardContent>
                </Card>
            </div>

            {/* Filters */}
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Filter className="h-5 w-5" />
                        Filters
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                        <div className="space-y-2">
                            <label className="text-sm font-medium">Search</label>
                            <div className="flex gap-2">
                                <Input
                                    placeholder="Transaction detail..."
                                    value={search}
                                    onChange={(e) => setSearch(e.target.value)}
                                    onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                                />
                                <Button onClick={handleSearch}>
                                    <Search className="h-4 w-4" />
                                </Button>
                            </div>
                        </div>

                        <div className="space-y-2">
                            <label className="text-sm font-medium">Status</label>
                            <Select value={statusFilter} onValueChange={setStatusFilter}>
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="all">All Statuses</SelectItem>
                                    <SelectItem value="pending">Pending</SelectItem>
                                    <SelectItem value="processing">Processing</SelectItem>
                                    <SelectItem value="completed">Completed</SelectItem>
                                    <SelectItem value="failed">Failed</SelectItem>
                                    <SelectItem value="refunded">Refunded</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="space-y-2">
                            <label className="text-sm font-medium">Gateway</label>
                            <Select value={gatewayFilter} onValueChange={setGatewayFilter}>
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="all">All Gateways</SelectItem>
                                    <SelectItem value="stripe">Stripe</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="space-y-2">
                            <label className="text-sm font-medium">Sort By</label>
                            <Select value={`${sortBy}-${sortOrder}`} onValueChange={(value) => {
                                const [field, order] = value.split('-')
                                setSortBy(field)
                                setSortOrder(order)
                            }}>
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="created_at-desc">Newest First</SelectItem>
                                    <SelectItem value="created_at-asc">Oldest First</SelectItem>
                                    <SelectItem value="amount-desc">Highest Amount</SelectItem>
                                    <SelectItem value="amount-asc">Lowest Amount</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>
                    </div>
                </CardContent>
            </Card>

            {/* Payments Table */}
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <CreditCard className="h-5 w-5" />
                        All Transactions
                    </CardTitle>
                    <CardDescription>
                        {total} total transactions found
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    {error && (
                        <div className="mb-4 p-4 bg-destructive/10 text-destructive rounded-md">
                            {error}
                        </div>
                    )}

                    {payments.length === 0 && !loading ? (
                        <div className="text-center py-12 text-muted-foreground">
                            <CreditCard className="h-12 w-12 mx-auto mb-4 opacity-50" />
                            <p>No payments found</p>
                        </div>
                    ) : (
                        <>
                            <div className="rounded-md border">
                                <Table>
                                    <TableHeader>
                                        <TableRow>
                                            <TableHead>Transaction ID</TableHead>
                                            <TableHead>Order ID</TableHead>
                                            <TableHead>Amount</TableHead>
                                            <TableHead>Gateway</TableHead>
                                            <TableHead>Status</TableHead>
                                            <TableHead>Date</TableHead>
                                            <TableHead className="text-right">Actions</TableHead>
                                        </TableRow>
                                    </TableHeader>
                                    <TableBody>
                                        {payments.map((payment) => (
                                            <TableRow key={payment.id}>
                                                <TableCell className="font-mono text-xs max-w-[150px] truncate">
                                                    {payment.transaction_id || 'N/A'}
                                                </TableCell>
                                                <TableCell>
                                                    <Link href={`/dashboard/orders/${payment.order_id}`} className="text-blue-600 hover:underline">
                                                        #{payment.order_id}
                                                    </Link>
                                                </TableCell>
                                                <TableCell className="font-medium">
                                                    {formatCurrency(payment.amount)} {payment.currency}
                                                </TableCell>
                                                <TableCell className="capitalize text-sm">
                                                    {payment.gateway_id}
                                                </TableCell>
                                                <TableCell>
                                                    <Badge className={getStatusColor(payment.status)}>
                                                        {payment.status.charAt(0).toUpperCase() + payment.status.slice(1)}
                                                    </Badge>
                                                </TableCell>
                                                <TableCell className="text-sm text-muted-foreground">
                                                    {formatDate(payment.created_at)}
                                                </TableCell>
                                                <TableCell className="text-right">
                                                    <Button variant="ghost" size="sm" asChild>
                                                        <Link href={`/dashboard/payments/${payment.id}`}>
                                                            <Eye className="h-4 w-4 mr-2" />
                                                            View
                                                        </Link>
                                                    </Button>
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </div>

                            {/* Pagination */}
                            {totalPages > 1 && (
                                <div className="flex items-center justify-between mt-4">
                                    <div className="text-sm text-muted-foreground">
                                        Showing {(page - 1) * limit + 1} to {Math.min(page * limit, total)} of {total} transactions
                                    </div>
                                    <div className="flex gap-2">
                                        <Button
                                            variant="outline"
                                            size="sm"
                                            onClick={() => setPage(p => Math.max(1, p - 1))}
                                            disabled={page === 1}
                                        >
                                            Previous
                                        </Button>
                                        <div className="flex items-center gap-2">
                                            <span className="text-sm">
                                                Page {page} of {totalPages}
                                            </span>
                                        </div>
                                        <Button
                                            variant="outline"
                                            size="sm"
                                            onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                                            disabled={page === totalPages}
                                        >
                                            Next
                                        </Button>
                                    </div>
                                </div>
                            )}
                        </>
                    )}
                </CardContent>
            </Card>
        </div>
    )
}
