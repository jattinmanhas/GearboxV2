"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { 
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { 
  Package, 
  TrendingUp, 
  TrendingDown,
  MoreHorizontal,
  Plus,
  Search,
  Filter,
  RefreshCw,
  ArrowUpDown,
  ArrowUp,
  ArrowDown,
  RotateCcw
} from "lucide-react"
import { 
  StockMovement,
  ListStockMovementsRequest,
  Product,
  StockMovementRequest
} from "@/lib/types"
import { StockMovementForm } from "../components/stock-movement-form"

export default function StockMovementsPage() {
  const [movements, setMovements] = useState<StockMovement[]>([])
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState("")
  const [filters, setFilters] = useState<ListStockMovementsRequest>({
    page: 1,
    limit: 20
  })
  const [showMovementForm, setShowMovementForm] = useState(false)

  // Fetch stock movements
  const fetchMovements = async () => {
    try {
      setLoading(true)
      const params = new URLSearchParams()
      
      if (filters.product_id) params.append('product_id', filters.product_id.toString())
      if (filters.movement_type) params.append('movement_type', filters.movement_type)
      if (filters.start_date) params.append('start_date', filters.start_date)
      if (filters.end_date) params.append('end_date', filters.end_date)
      if (filters.page) params.append('page', filters.page.toString())
      if (filters.limit) params.append('limit', filters.limit.toString())

      const response = await fetch(`/api/v1/inventory/movements?${params}`, {
        credentials: 'include'
      })
      const data = await response.json()
      
      if (response.ok) {
        setMovements(data.data?.movements || [])
      }
    } catch (error) {
      console.error('Error fetching stock movements:', error)
    } finally {
      setLoading(false)
    }
  }

  // Fetch products
  const fetchProducts = async () => {
    try {
      const response = await fetch('/api/v1/products?limit=1000', {
        credentials: 'include'
      })
      const data = await response.json()
      
      if (response.ok) {
        setProducts(data.data?.products || [])
      }
    } catch (error) {
      console.error('Error fetching products:', error)
    }
  }

  useEffect(() => {
    fetchMovements()
    fetchProducts()
  }, [filters])

  const handleSaveStockMovement = async (data: StockMovementRequest) => {
    try {
      const response = await fetch('/api/v1/inventory/movements', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(data)
      })

      if (response.ok) {
        setShowMovementForm(false)
        fetchMovements()
      }
    } catch (error) {
      console.error('Error recording stock movement:', error)
    }
  }

  // Filter movements based on search term
  const filteredMovements = movements.filter(movement => 
    movement.product_name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    movement.product_sku?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    movement.variant_name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    movement.variant_sku?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    movement.reason?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    movement.reference?.toLowerCase().includes(searchTerm.toLowerCase())
  )

  const getMovementIcon = (type: string) => {
    switch (type) {
      case 'in':
        return <ArrowUp className="h-4 w-4 text-green-500" />
      case 'out':
        return <ArrowDown className="h-4 w-4 text-red-500" />
      case 'adjustment':
        return <RotateCcw className="h-4 w-4 text-blue-500" />
      case 'transfer':
        return <ArrowUpDown className="h-4 w-4 text-purple-500" />
      default:
        return <Package className="h-4 w-4 text-gray-500" />
    }
  }

  const getMovementColor = (type: string) => {
    switch (type) {
      case 'in':
        return 'bg-green-100 text-green-800'
      case 'out':
        return 'bg-red-100 text-red-800'
      case 'adjustment':
        return 'bg-blue-100 text-blue-800'
      case 'transfer':
        return 'bg-purple-100 text-purple-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getMovementTypeLabel = (type: string) => {
    switch (type) {
      case 'in': return 'Stock In'
      case 'out': return 'Stock Out'
      case 'adjustment': return 'Adjustment'
      case 'transfer': return 'Transfer'
      default: return type
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Stock Movements</h1>
          <p className="text-muted-foreground">
            Track and manage all inventory stock movements
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchMovements}>
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </Button>
          <Button onClick={() => setShowMovementForm(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Record Movement
          </Button>
        </div>
      </div>

      {/* Filters and Search */}
      <Card>
        <CardHeader>
          <CardTitle>Stock Movements</CardTitle>
          <CardDescription>
            View and manage all stock movements in your inventory
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-4 mb-4">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
              <Input
                placeholder="Search movements, products, SKUs, or references..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10"
              />
            </div>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline">
                  <Filter className="h-4 w-4 mr-2" />
                  Filters
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-56">
                <DropdownMenuItem
                  onClick={() => setFilters(prev => ({ ...prev, movement_type: prev.movement_type === 'in' ? undefined : 'in' }))}
                >
                  Stock In Only
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => setFilters(prev => ({ ...prev, movement_type: prev.movement_type === 'out' ? undefined : 'out' }))}
                >
                  Stock Out Only
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => setFilters(prev => ({ ...prev, movement_type: prev.movement_type === 'adjustment' ? undefined : 'adjustment' }))}
                >
                  Adjustments Only
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => setFilters(prev => ({ ...prev, movement_type: prev.movement_type === 'transfer' ? undefined : 'transfer' }))}
                >
                  Transfers Only
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          {/* Movements Table */}
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Type</TableHead>
                  <TableHead>Product</TableHead>
                  <TableHead>Quantity</TableHead>
                  <TableHead>Previous</TableHead>
                  <TableHead>New Total</TableHead>
                  <TableHead>Reason</TableHead>
                  <TableHead>Reference</TableHead>
                  <TableHead>Date</TableHead>
                  <TableHead className="w-[50px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {loading ? (
                  <TableRow>
                    <TableCell colSpan={9} className="text-center py-8">
                      <RefreshCw className="h-6 w-6 animate-spin mx-auto mb-2" />
                      Loading movements...
                    </TableCell>
                  </TableRow>
                ) : filteredMovements.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={9} className="text-center py-8 text-muted-foreground">
                      No stock movements found
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredMovements.map((movement) => (
                    <TableRow key={movement.id}>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          {getMovementIcon(movement.movement_type)}
                          <Badge className={getMovementColor(movement.movement_type)}>
                            {getMovementTypeLabel(movement.movement_type)}
                          </Badge>
                        </div>
                      </TableCell>
                      <TableCell>
                        <div>
                          <p className="font-medium">{movement.product_name}</p>
                          {movement.variant_name && (
                            <p className="text-sm text-muted-foreground">{movement.variant_name}</p>
                          )}
                          <p className="font-mono text-xs text-muted-foreground">
                            {movement.product_sku} {movement.variant_sku && `(${movement.variant_sku})`}
                          </p>
                        </div>
                      </TableCell>
                      <TableCell className="font-mono">
                        <span className={movement.movement_type === 'in' ? 'text-green-600' : 'text-red-600'}>
                          {movement.movement_type === 'in' ? '+' : '-'}{movement.quantity.toLocaleString()}
                        </span>
                      </TableCell>
                      <TableCell className="font-mono">{movement.previous_quantity.toLocaleString()}</TableCell>
                      <TableCell className="font-mono font-medium">{movement.new_quantity.toLocaleString()}</TableCell>
                      <TableCell>
                        <div className="max-w-[200px]">
                          <p className="truncate">{movement.reason}</p>
                          {movement.notes && (
                            <p className="text-xs text-muted-foreground truncate">{movement.notes}</p>
                          )}
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="max-w-[150px]">
                          {movement.reference && (
                            <p className="font-mono text-sm truncate">{movement.reference}</p>
                          )}
                          {movement.reference_type && (
                            <p className="text-xs text-muted-foreground truncate">{movement.reference_type}</p>
                          )}
                        </div>
                      </TableCell>
                      <TableCell className="text-sm text-muted-foreground">
                        {new Date(movement.created_at).toLocaleString()}
                      </TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="sm">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem>View Details</DropdownMenuItem>
                            <DropdownMenuItem>Edit Movement</DropdownMenuItem>
                            <DropdownMenuItem className="text-destructive">Delete</DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      {/* Stock Movement Form Modal */}
      {showMovementForm && (
        <StockMovementForm
          products={products}
          onSave={handleSaveStockMovement}
          onCancel={() => setShowMovementForm(false)}
        />
      )}
    </div>
  )
}
