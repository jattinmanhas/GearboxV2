"use client"

import { useState, useEffect } from "react"
import { useCouponStore } from "@/lib/stores/coupon-store"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from "@/components/ui/table"
import { 
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { 
  Plus, 
  Search, 
  MoreHorizontal, 
  Edit, 
  Trash2, 
  Eye,
  Filter,
  Download
} from "lucide-react"
import { format } from "date-fns"
import { CouponForm } from "./components/coupon-form"
import { CouponStats } from "./components/coupon-stats"

export default function CouponsPage() {
  const { 
    coupons, 
    isLoading, 
    error, 
    loadCoupons, 
    deleteCoupon 
  } = useCouponStore()
  
  const [searchTerm, setSearchTerm] = useState("")
  const [statusFilter, setStatusFilter] = useState<string>("all")
  const [typeFilter, setTypeFilter] = useState<string>("all")
  const [showForm, setShowForm] = useState(false)
  const [editingCoupon, setEditingCoupon] = useState<number | null>(null)
  const [currentPage, setCurrentPage] = useState(1)
  const [itemsPerPage] = useState(10)

  useEffect(() => {
    loadCoupons({
      page: currentPage,
      limit: itemsPerPage,
      search: searchTerm || undefined,
      status: statusFilter !== "all" ? statusFilter : undefined,
      type: typeFilter !== "all" ? typeFilter : undefined,
    })
  }, [loadCoupons, currentPage, searchTerm, statusFilter, typeFilter])

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    setCurrentPage(1)
    loadCoupons({
      page: 1,
      limit: itemsPerPage,
      search: searchTerm || undefined,
      status: statusFilter !== "all" ? statusFilter : undefined,
      type: typeFilter !== "all" ? typeFilter : undefined,
    })
  }

  const handleDelete = async (id: number) => {
    if (confirm("Are you sure you want to delete this coupon?")) {
      await deleteCoupon(id.toString())
    }
  }

  const handleEdit = (coupon: any) => {
    setEditingCoupon(coupon.id)
    setShowForm(true)
  }

  const getStatusBadge = (coupon: any) => {
    const now = new Date()
    const startsAt = new Date(coupon.starts_at)
    const expiresAt = coupon.expires_at ? new Date(coupon.expires_at) : null

    if (!coupon.is_active) {
      return <Badge variant="secondary">Inactive</Badge>
    }
    
    if (now < startsAt) {
      return <Badge variant="outline">Scheduled</Badge>
    }
    
    if (expiresAt && now > expiresAt) {
      return <Badge variant="destructive">Expired</Badge>
    }
    
    return <Badge variant="default">Active</Badge>
  }

  const getTypeBadge = (type: string) => {
    const typeMap = {
      percentage: { label: "Percentage", variant: "default" as const },
      fixed_amount: { label: "Fixed Amount", variant: "secondary" as const },
      free_shipping: { label: "Free Shipping", variant: "outline" as const },
    }
    
    const { label, variant } = typeMap[type as keyof typeof typeMap] || { label: type, variant: "default" as const }
    return <Badge variant={variant}>{label}</Badge>
  }

  const formatValue = (coupon: any) => {
    if (coupon.type === "percentage") {
      return `${coupon.value}%`
    } else if (coupon.type === "fixed_amount") {
      return `₹${coupon.value}`
    } else {
      return "Free"
    }
  }

  if (showForm) {
    return (
      <CouponForm 
        couponId={editingCoupon}
        onClose={() => {
          setShowForm(false)
          setEditingCoupon(null)
        }}
        onSuccess={() => {
          setShowForm(false)
          setEditingCoupon(null)
          loadCoupons()
        }}
      />
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Coupons</h1>
          <p className="text-muted-foreground">
            Manage discount coupons and promotional codes
          </p>
        </div>
        <Button onClick={() => setShowForm(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Create Coupon
        </Button>
      </div>

      <CouponStats />

      <Card>
        <CardHeader>
          <CardTitle>All Coupons</CardTitle>
          <CardDescription>
            View and manage all discount coupons
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <form onSubmit={handleSearch} className="flex gap-2">
              <div className="relative flex-1">
                <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search coupons..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-8"
                />
              </div>
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                className="px-3 py-2 border border-input bg-background rounded-md"
              >
                <option value="all">All Status</option>
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
                <option value="expired">Expired</option>
              </select>
              <select
                value={typeFilter}
                onChange={(e) => setTypeFilter(e.target.value)}
                className="px-3 py-2 border border-input bg-background rounded-md"
              >
                <option value="all">All Types</option>
                <option value="percentage">Percentage</option>
                <option value="fixed_amount">Fixed Amount</option>
                <option value="free_shipping">Free Shipping</option>
              </select>
              <Button type="submit" variant="outline">
                <Filter className="mr-2 h-4 w-4" />
                Filter
              </Button>
            </form>

            {isLoading ? (
              <div className="text-center py-8">Loading coupons...</div>
            ) : error ? (
              <div className="text-center py-8 text-destructive">{error}</div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Code</TableHead>
                    <TableHead>Name</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Value</TableHead>
                    <TableHead>Usage</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Expires</TableHead>
                    <TableHead className="w-[50px]"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {coupons.map((coupon) => (
                    <TableRow key={coupon.id}>
                      <TableCell className="font-mono font-medium">
                        {coupon.code}
                      </TableCell>
                      <TableCell>{coupon.name}</TableCell>
                      <TableCell>{getTypeBadge(coupon.type)}</TableCell>
                      <TableCell>{formatValue(coupon)}</TableCell>
                      <TableCell>
                        {coupon.used_count}
                        {coupon.usage_limit && ` / ${coupon.usage_limit}`}
                      </TableCell>
                      <TableCell>{getStatusBadge(coupon)}</TableCell>
                      <TableCell>
                        {coupon.expires_at 
                          ? format(new Date(coupon.expires_at), "MMM dd, yyyy")
                          : "No expiry"
                        }
                      </TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" className="h-8 w-8 p-0">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuLabel>Actions</DropdownMenuLabel>
                            <DropdownMenuItem onClick={() => handleEdit(coupon)}>
                              <Edit className="mr-2 h-4 w-4" />
                              Edit
                            </DropdownMenuItem>
                            <DropdownMenuItem>
                              <Eye className="mr-2 h-4 w-4" />
                              View Details
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem 
                              onClick={() => handleDelete(coupon.id)}
                              className="text-destructive"
                            >
                              <Trash2 className="mr-2 h-4 w-4" />
                              Delete
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
