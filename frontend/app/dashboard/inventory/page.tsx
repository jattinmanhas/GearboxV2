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
  AlertTriangle, 
  TrendingUp, 
  TrendingDown,
  MoreHorizontal,
  Plus,
  Search,
  Filter,
  Download,
  RefreshCw,
  ArrowUpDown,
  Bell,
  PieChart
} from "lucide-react"
import { 
  Inventory, 
  InventorySummary, 
  ListInventoryRequest,
  StockMovement,
  InventoryAlert,
  Product,
  CreateInventoryRequest,
  StockMovementRequest
} from "@/lib/types"
import { formatCurrency } from "@/lib/currency"
import { InventoryForm } from "./components/inventory-form"
import { StockMovementForm } from "./components/stock-movement-form"
import Link from "next/link"

export default function InventoryDashboard() {
  const [inventory, setInventory] = useState<Inventory[]>([])
  const [summary, setSummary] = useState<InventorySummary | null>(null)
  const [alerts, setAlerts] = useState<InventoryAlert[]>([])
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState("")
  const [filters, setFilters] = useState<ListInventoryRequest>({
    page: 1,
    limit: 20,
    low_stock: false,
    out_of_stock: false
  })
  
  // Modal states
  const [showInventoryForm, setShowInventoryForm] = useState(false)
  const [showStockMovementForm, setShowStockMovementForm] = useState(false)
  const [editingInventory, setEditingInventory] = useState<Inventory | null>(null)
  const [selectedInventoryForMovement, setSelectedInventoryForMovement] = useState<Inventory | null>(null)

  // Fetch inventory data
  const fetchInventory = async () => {
    try {
      setLoading(true)
      const params = new URLSearchParams()
      
      if (filters.product_id) params.append('product_id', filters.product_id.toString())
      if (filters.low_stock) params.append('low_stock', 'true')
      if (filters.out_of_stock) params.append('out_of_stock', 'true')
      if (filters.page) params.append('page', filters.page.toString())
      if (filters.limit) params.append('limit', filters.limit.toString())

      console.log('Fetching inventory with params:', params.toString())
      const response = await fetch(`/api/v1/inventory?${params}`, {
        credentials: 'include'
      })
      const data = await response.json()
      
      console.log('Inventory response:', { status: response.status, data })
      
      if (response.ok) {
        // Handle different possible response structures
        const inventoryData = data.data?.inventory || data.inventory || data.data || []
        console.log('Setting inventory to:', inventoryData)
        setInventory(inventoryData)
      } else {
        console.error('Inventory fetch failed:', data)
      }
    } catch (error) {
      console.error('Error fetching inventory:', error)
    } finally {
      setLoading(false)
    }
  }

  // Fetch inventory summary
  const fetchSummary = async () => {
    try {
      const response = await fetch('/api/v1/inventory/summary', {
        credentials: 'include'
      })
      const data = await response.json()
      
      if (response.ok) {
        setSummary(data.data)
      }
    } catch (error) {
      console.error('Error fetching inventory summary:', error)
    }
  }

  // Fetch inventory alerts
  const fetchAlerts = async () => {
    try {
      const response = await fetch('/api/v1/inventory/alerts', {
        credentials: 'include'
      })
      const data = await response.json()
      
      if (response.ok) {
        setAlerts(data.data || [])
      }
    } catch (error) {
      console.error('Error fetching inventory alerts:', error)
    }
  }

  // Fetch products for forms
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
    fetchInventory()
    fetchSummary()
    fetchAlerts()
    fetchProducts()
  }, [filters])

  // Handler functions
  const handleCreateInventory = () => {
    setEditingInventory(null)
    setShowInventoryForm(true)
  }

  const handleEditInventory = (inventoryItem: Inventory) => {
    setEditingInventory(inventoryItem)
    setShowInventoryForm(true)
  }

  const handleRecordMovement = (inventoryItem: Inventory) => {
    setSelectedInventoryForMovement(inventoryItem)
    setShowStockMovementForm(true)
  }

  const handleSaveInventory = async (data: CreateInventoryRequest) => {
    try {
      const response = await fetch('/api/v1/inventory', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(data)
      })

      if (response.ok) {
        setShowInventoryForm(false)
        fetchInventory()
        fetchSummary()
      }
    } catch (error) {
      console.error('Error saving inventory:', error)
    }
  }

  const handleSaveStockMovement = async (data: StockMovementRequest) => {
    try {
      const response = await fetch('/api/v1/inventory/movements', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(data)
      })

      if (response.ok) {
        setShowStockMovementForm(false)
        fetchInventory()
        fetchSummary()
      }
    } catch (error) {
      console.error('Error recording stock movement:', error)
    }
  }

  // Filter inventory based on search term
  const filteredInventory = inventory.filter(item => 
    item.product_name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    item.product_sku?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    item.variant_name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    item.variant_sku?.toLowerCase().includes(searchTerm.toLowerCase())
  )
  
  console.log('Total inventory items:', inventory.length)
  console.log('Filtered inventory items:', filteredInventory.length)

  const getStockStatus = (item: Inventory) => {
    if (item.available_quantity === 0) return { status: 'out', color: 'destructive', text: 'Out of Stock' }
    if (item.available_quantity <= item.min_stock_level) return { status: 'low', color: 'destructive', text: 'Low Stock' }
    if (item.max_stock_level && item.available_quantity >= item.max_stock_level) return { status: 'high', color: 'secondary', text: 'Overstock' }
    return { status: 'good', color: 'default', text: 'In Stock' }
  }

  const getAlertIcon = (alertType: string) => {
    switch (alertType) {
      case 'low_stock':
        return <TrendingDown className="h-4 w-4 text-yellow-500" />
      case 'out_of_stock':
        return <AlertTriangle className="h-4 w-4 text-red-500" />
      case 'overstock':
        return <TrendingUp className="h-4 w-4 text-blue-500" />
      default:
        return <AlertTriangle className="h-4 w-4 text-gray-500" />
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="space-y-4">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Inventory Management</h1>
            <p className="text-muted-foreground">
              Manage your product inventory, stock levels, and alerts
            </p>
          </div>
          
          {/* Primary Action */}
          <div className="flex items-center gap-2">
            <Button onClick={handleCreateInventory} className="flex-shrink-0">
              <Plus className="h-4 w-4 mr-2" />
              Add Inventory
            </Button>
          </div>
        </div>

        {/* Secondary Actions - Organized in rows */}
        <div className="flex flex-col lg:flex-row gap-4">
          {/* Navigation Links */}
          <div className="flex flex-wrap items-center gap-2">
            <Button variant="outline" size="sm" asChild>
              <Link href="/dashboard/inventory/summary">
                <PieChart className="h-4 w-4 mr-2" />
                Summary
              </Link>
            </Button>
            <Button variant="outline" size="sm" asChild>
              <Link href="/dashboard/inventory/movements">
                <ArrowUpDown className="h-4 w-4 mr-2" />
                Stock Movements
              </Link>
            </Button>
            <Button variant="outline" size="sm" asChild>
              <Link href="/dashboard/inventory/alerts">
                <Bell className="h-4 w-4 mr-2" />
                Alerts
              </Link>
            </Button>
          </div>

          {/* Utility Actions */}
          <div className="flex flex-wrap items-center gap-2">
            <Button variant="outline" size="sm">
              <Download className="h-4 w-4 mr-2" />
              Export
            </Button>
            <Button variant="outline" size="sm" onClick={() => { fetchInventory(); fetchSummary(); fetchAlerts(); }}>
              <RefreshCw className="h-4 w-4 mr-2" />
              Refresh
            </Button>
            <Button variant="outline" size="sm" onClick={async () => {
              try {
                console.log('Testing authentication...')
                const response = await fetch('/api/v1/auth/profile', {
                  credentials: 'include'
                })
                const data = await response.json()
                console.log('Auth test response:', { status: response.status, data })
              } catch (error) {
                console.error('Auth test error:', error)
              }
            }}>
              Test Auth
            </Button>
          </div>
        </div>
      </div>

      {/* Summary Cards */}
      {summary && (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Products</CardTitle>
              <Package className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{summary.total_products}</div>
              <p className="text-xs text-muted-foreground">
                {summary.total_variants} variants
              </p>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Stock</CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
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
              <CardTitle className="text-sm font-medium">Low Stock Items</CardTitle>
              <AlertTriangle className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-yellow-600">{summary.low_stock_items}</div>
              <p className="text-xs text-muted-foreground">
                {summary.out_of_stock_items} out of stock
              </p>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Value</CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatCurrency(summary.total_value)}</div>
              <p className="text-xs text-muted-foreground">
                Avg: {formatCurrency(summary.average_stock_level)}
              </p>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Alerts Section */}
      {alerts.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-yellow-500" />
              Inventory Alerts
            </CardTitle>
            <CardDescription>
              Items that require your attention
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {alerts.slice(0, 5).map((alert) => (
                <div key={alert.id} className="flex items-center justify-between p-3 border rounded-lg">
                  <div className="flex items-center gap-3">
                    {getAlertIcon(alert.alert_type)}
                    <div>
                      <p className="font-medium">
                        {alert.product_name} {alert.variant_name && `(${alert.variant_name})`}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {alert.alert_type.replace('_', ' ')} - {alert.current_quantity} units
                      </p>
                    </div>
                  </div>
                  <Badge variant={alert.is_resolved ? "secondary" : "destructive"}>
                    {alert.is_resolved ? "Resolved" : "Active"}
                  </Badge>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Filters and Search */}
      <Card>
        <CardHeader>
          <CardTitle>Inventory Items</CardTitle>
          <CardDescription>
            Manage your product inventory and stock levels
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col sm:flex-row gap-4 mb-6">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
              <Input
                placeholder="Search products, SKUs, or variants..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10"
              />
            </div>
            <div className="flex items-center gap-2">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="outline" className="flex-shrink-0">
                    <Filter className="h-4 w-4 mr-2" />
                    Filters
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56">
                  <DropdownMenuItem
                    onClick={() => setFilters(prev => ({ ...prev, low_stock: !prev.low_stock }))}
                  >
                    Low Stock Only
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onClick={() => setFilters(prev => ({ ...prev, out_of_stock: !prev.out_of_stock }))}
                  >
                    Out of Stock Only
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>

          {/* Inventory Table */}
          <div className="rounded-md border overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="min-w-[200px]">Product</TableHead>
                  <TableHead className="min-w-[120px]">SKU</TableHead>
                  <TableHead className="min-w-[100px]">Current Stock</TableHead>
                  <TableHead className="min-w-[100px]">Available</TableHead>
                  <TableHead className="min-w-[100px]">Reserved</TableHead>
                  <TableHead className="min-w-[100px]">Min Level</TableHead>
                  <TableHead className="min-w-[120px]">Status</TableHead>
                  <TableHead className="w-[50px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {loading ? (
                  <TableRow>
                    <TableCell colSpan={8} className="text-center py-8">
                      <RefreshCw className="h-6 w-6 animate-spin mx-auto mb-2" />
                      Loading inventory...
                    </TableCell>
                  </TableRow>
                ) : filteredInventory.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={8} className="text-center py-8 text-muted-foreground">
                      No inventory items found
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredInventory.map((item) => {
                    const stockStatus = getStockStatus(item)
                    return (
                      <TableRow key={item.id}>
                        <TableCell>
                          <div>
                            <p className="font-medium">{item.product_name}</p>
                            {item.variant_name && (
                              <p className="text-sm text-muted-foreground">{item.variant_name}</p>
                            )}
                          </div>
                        </TableCell>
                        <TableCell>
                          <div>
                            <p className="font-mono text-sm">{item.product_sku}</p>
                            {item.variant_sku && (
                              <p className="font-mono text-xs text-muted-foreground">{item.variant_sku}</p>
                            )}
                          </div>
                        </TableCell>
                        <TableCell className="font-mono">{item.quantity.toLocaleString()}</TableCell>
                        <TableCell className="font-mono">{item.available_quantity.toLocaleString()}</TableCell>
                        <TableCell className="font-mono">{item.reserved_quantity.toLocaleString()}</TableCell>
                        <TableCell className="font-mono">{item.min_stock_level.toLocaleString()}</TableCell>
                        <TableCell>
                          <Badge variant={stockStatus.color as any}>
                            {stockStatus.text}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="sm">
                                <MoreHorizontal className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem onClick={() => handleEditInventory(item)}>View Details</DropdownMenuItem>
                              <DropdownMenuItem onClick={() => handleEditInventory(item)}>Update Stock</DropdownMenuItem>
                              <DropdownMenuItem onClick={() => handleRecordMovement(item)}>Record Movement</DropdownMenuItem>
                              <DropdownMenuItem className="text-destructive">Delete</DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </TableCell>
                      </TableRow>
                    )
                  })
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      {/* Modals */}
      {showInventoryForm && (
        <InventoryForm
          inventory={editingInventory}
          products={products}
          onSave={handleSaveInventory}
          onCancel={() => setShowInventoryForm(false)}
        />
      )}

      {showStockMovementForm && (
        <StockMovementForm
          inventoryId={selectedInventoryForMovement?.id}
          products={products}
          onSave={handleSaveStockMovement}
          onCancel={() => setShowStockMovementForm(false)}
        />
      )}
    </div>
  )
}
