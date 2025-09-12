"use client"

import { useState, useMemo } from "react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { 
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet"
import { 
  ShoppingCart, 
  Plus, 
  Minus, 
  Trash2,
  Package,
  ArrowRight,
  Search,
  Filter,
  SortAsc,
  SortDesc
} from "lucide-react"
import { useCartStore } from "@/lib/stores/cart-store"
import { formatPrice } from "@/lib/currency"

interface CartDrawerProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CartDrawer({ open, onOpenChange }: CartDrawerProps) {
  const [isUpdating, setIsUpdating] = useState(false)
  const [searchTerm, setSearchTerm] = useState("")
  const [sortBy, setSortBy] = useState<'name' | 'price' | 'quantity'>('name')
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc')
  const [showFilters, setShowFilters] = useState(false)
  const [minPrice, setMinPrice] = useState<number | undefined>()
  const [maxPrice, setMaxPrice] = useState<number | undefined>()
  const [priceFilter, setPriceFilter] = useState<'all' | 'low' | 'medium' | 'high'>('all')
  
  const { 
    items: cartItems, 
    updateQuantity, 
    removeItem, 
    clearCart, 
    getTotalPrice 
  } = useCartStore()


  // Filter and sort cart items
  const filteredAndSortedItems = useMemo(() => {
    let filtered = cartItems.filter(item => {
      // Search filter
      const matchesSearch = item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        item.sku.toLowerCase().includes(searchTerm.toLowerCase())
      
      if (!matchesSearch) return false

      // Price filters
      if (minPrice !== undefined && item.price < minPrice) return false
      if (maxPrice !== undefined && item.price > maxPrice) return false

      // Price range filter
      switch (priceFilter) {
        case 'low':
          if (item.price >= 50) return false
          break
        case 'medium':
          if (item.price < 50 || item.price >= 200) return false
          break
        case 'high':
          if (item.price < 200) return false
          break
        case 'all':
        default:
          break
      }

      return true
    })

    // Sort items
    filtered.sort((a, b) => {
      let aValue: string | number
      let bValue: string | number

      switch (sortBy) {
        case 'name':
          aValue = a.name.toLowerCase()
          bValue = b.name.toLowerCase()
          break
        case 'price':
          aValue = a.price
          bValue = b.price
          break
        case 'quantity':
          aValue = a.quantity
          bValue = b.quantity
          break
        default:
          aValue = a.name.toLowerCase()
          bValue = b.name.toLowerCase()
      }

      if (aValue < bValue) return sortOrder === 'asc' ? -1 : 1
      if (aValue > bValue) return sortOrder === 'asc' ? 1 : -1
      return 0
    })

    return filtered
  }, [cartItems, searchTerm, sortBy, sortOrder, minPrice, maxPrice, priceFilter])

  const handleSort = (field: 'name' | 'price' | 'quantity') => {
    if (sortBy === field) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')
    } else {
      setSortBy(field)
      setSortOrder('asc')
    }
  }

  const handleUpdateQuantity = async (id: number, newQuantity: number) => {
    setIsUpdating(true)
    try {
      updateQuantity(id, newQuantity)
    } catch (error) {
      console.error("Failed to update quantity:", error)
    } finally {
      setIsUpdating(false)
    }
  }

  const handleRemoveItem = async (id: number) => {
    setIsUpdating(true)
    try {
      removeItem(id)
    } catch (error) {
      console.error("Failed to remove item:", error)
    } finally {
      setIsUpdating(false)
    }
  }

  const handleClearCart = () => {
    clearCart()
  }

  const subtotal = getTotalPrice()
  const shipping = subtotal > 50 ? 0 : 9.99
  const tax = subtotal * 0.08 // 8% tax
  const total = subtotal + shipping + tax

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="w-full sm:max-w-lg">
        <SheetHeader>
          <SheetTitle className="flex items-center gap-2">
            <ShoppingCart className="h-5 w-5" />
            Shopping Cart
            {cartItems.length > 0 && (
              <Badge variant="secondary">{cartItems.length}</Badge>
            )}
          </SheetTitle>
          <SheetDescription>
            Review your items and proceed to checkout
          </SheetDescription>
        </SheetHeader>

        {/* Search and Filter Controls */}
        {cartItems.length > 0 && (
          <div className="space-y-3 py-4 border-b">
            {/* Search */}
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
              <Input
                placeholder="Search items..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10"
              />
            </div>

            {/* Sort Controls */}
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <span className="text-sm text-muted-foreground">Sort by:</span>
                <Button
                  variant={sortBy === 'name' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => handleSort('name')}
                  className="h-8"
                >
                  Name
                  {sortBy === 'name' && (
                    sortOrder === 'asc' ? <SortAsc className="h-3 w-3 ml-1" /> : <SortDesc className="h-3 w-3 ml-1" />
                  )}
                </Button>
                <Button
                  variant={sortBy === 'price' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => handleSort('price')}
                  className="h-8"
                >
                  Price
                  {sortBy === 'price' && (
                    sortOrder === 'asc' ? <SortAsc className="h-3 w-3 ml-1" /> : <SortDesc className="h-3 w-3 ml-1" />
                  )}
                </Button>
                <Button
                  variant={sortBy === 'quantity' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => handleSort('quantity')}
                  className="h-8"
                >
                  Qty
                  {sortBy === 'quantity' && (
                    sortOrder === 'asc' ? <SortAsc className="h-3 w-3 ml-1" /> : <SortDesc className="h-3 w-3 ml-1" />
                  )}
                </Button>
              </div>
              
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowFilters(!showFilters)}
                className="h-8"
              >
                <Filter className="h-3 w-3 mr-1" />
                Filters
              </Button>
            </div>

            {/* Results Count */}
            {(searchTerm || priceFilter !== 'all' || minPrice !== undefined || maxPrice !== undefined) && (
              <div className="text-sm text-muted-foreground">
                {filteredAndSortedItems.length} of {cartItems.length} items
              </div>
            )}
          </div>
        )}

        {/* Advanced Filters Panel */}
        {showFilters && cartItems.length > 0 && (
          <div className="border-b py-4 space-y-4">
            <h3 className="font-medium text-sm">Filters</h3>
            
            {/* Price Range Filter */}
            <div className="space-y-2">
              <label className="text-sm font-medium">Price Range</label>
              <div className="grid grid-cols-2 gap-2">
                <Button
                  variant={priceFilter === 'all' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setPriceFilter('all')}
                  className="h-8"
                >
                  All
                </Button>
                <Button
                  variant={priceFilter === 'low' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setPriceFilter('low')}
                  className="h-8"
                >
                  Under $50
                </Button>
                <Button
                  variant={priceFilter === 'medium' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setPriceFilter('medium')}
                  className="h-8"
                >
                  $50 - $200
                </Button>
                <Button
                  variant={priceFilter === 'high' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setPriceFilter('high')}
                  className="h-8"
                >
                  Over $200
                </Button>
              </div>
            </div>

            {/* Custom Price Range */}
            <div className="space-y-2">
              <label className="text-sm font-medium">Custom Price Range</label>
              <div className="grid grid-cols-2 gap-2">
                <Input
                  type="number"
                  placeholder="Min price"
                  value={minPrice || ''}
                  onChange={(e) => setMinPrice(e.target.value ? parseFloat(e.target.value) : undefined)}
                  className="h-8"
                />
                <Input
                  type="number"
                  placeholder="Max price"
                  value={maxPrice || ''}
                  onChange={(e) => setMaxPrice(e.target.value ? parseFloat(e.target.value) : undefined)}
                  className="h-8"
                />
              </div>
            </div>

            {/* Clear Filters */}
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                setSearchTerm('')
                setMinPrice(undefined)
                setMaxPrice(undefined)
                setPriceFilter('all')
              }}
              className="w-full h-8"
            >
              Clear All Filters
            </Button>
          </div>
        )}

        <div className="flex-1 overflow-hidden flex flex-col">
          {cartItems.length === 0 ? (
            <div className="flex-1 flex items-center justify-center">
              <div className="text-center">
                <Package className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
                <h3 className="text-lg font-semibold mb-2">Your cart is empty</h3>
                <p className="text-muted-foreground mb-4">
                  Add some products to get started
                </p>
                <Button onClick={() => onOpenChange(false)}>
                  Continue Shopping
                </Button>
              </div>
            </div>
          ) : (
            <>
              {/* Cart Items */}
              <div className="flex-1 overflow-y-auto space-y-4 py-4">
                {filteredAndSortedItems.length === 0 ? (
                  <div className="text-center py-8">
                    <Package className="h-8 w-8 mx-auto mb-2 text-muted-foreground" />
                    <p className="text-muted-foreground">
                      {searchTerm ? 'No items match your search' : 'No items in cart'}
                    </p>
                    {searchTerm && (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setSearchTerm('')}
                        className="mt-2"
                      >
                        Clear Search
                      </Button>
                    )}
                  </div>
                ) : (
                  filteredAndSortedItems.map((item) => (
                  <div key={item.id} className="flex gap-3 p-3 border rounded-lg">
                    {/* Product Image */}
                    <div className="w-16 h-16 bg-muted rounded-md flex items-center justify-center">
                      <Package className="h-6 w-6 text-muted-foreground" />
                    </div>
                    
                    {/* Product Info */}
                    <div className="flex-1 min-w-0">
                      <h4 className="font-medium text-sm line-clamp-2">{item.name}</h4>
                      <p className="text-xs text-muted-foreground">SKU: {item.sku}</p>
                      <p className="text-sm font-semibold">{formatPrice(item.price)}</p>
                      
                      {/* Quantity Controls */}
                      <div className="flex items-center gap-2 mt-2">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleUpdateQuantity(item.id, item.quantity - 1)}
                          disabled={isUpdating}
                          className="h-7 w-7 p-0"
                        >
                          <Minus className="h-3 w-3" />
                        </Button>
                        <span className="text-sm font-medium min-w-[2rem] text-center">
                          {item.quantity}
                        </span>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleUpdateQuantity(item.id, item.quantity + 1)}
                          disabled={isUpdating}
                          className="h-7 w-7 p-0"
                        >
                          <Plus className="h-3 w-3" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleRemoveItem(item.id)}
                          disabled={isUpdating}
                          className="h-7 w-7 p-0 text-destructive hover:text-destructive"
                        >
                          <Trash2 className="h-3 w-3" />
                        </Button>
                      </div>
                    </div>
                  </div>
                  ))
                )}
              </div>

              {/* Cart Summary */}
              <div className="border-t pt-4 space-y-4">
                {/* Clear Cart */}
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={handleClearCart}
                  className="w-full text-destructive hover:text-destructive"
                >
                  Clear Cart
                </Button>

                {/* Order Summary */}
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span>Subtotal</span>
                    <span>{formatPrice(subtotal)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Shipping</span>
                    <span>
                      {shipping === 0 ? "Free" : formatPrice(shipping)}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span>Tax</span>
                    <span>{formatPrice(tax)}</span>
                  </div>
                  <div className="flex justify-between font-semibold text-base border-t pt-2">
                    <span>Total</span>
                    <span>{formatPrice(total)}</span>
                  </div>
                </div>

                {/* Checkout Button */}
                <Button className="w-full" size="lg">
                  Proceed to Checkout
                  <ArrowRight className="h-4 w-4 ml-2" />
                </Button>

                {/* Continue Shopping */}
                <Button
                  variant="outline"
                  className="w-full"
                  onClick={() => onOpenChange(false)}
                >
                  Continue Shopping
                </Button>
              </div>
            </>
          )}
        </div>
      </SheetContent>
    </Sheet>
  )
}
