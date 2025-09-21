"use client"

import { useState, useEffect } from "react"
import Link from "next/link"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { 
  Search, 
  Filter, 
  Grid, 
  List, 
  ShoppingCart,
  Star,
  Heart,
  Package
} from "lucide-react"
import { productApi } from "@/lib/api"
import { Product, Category } from "@/lib/types"
import { ProductCard } from "./components/product-card"
import { ProductFilters } from "./components/product-filters"
import { CartDrawer } from "./components/cart-drawer"
import { useCartStore } from "@/lib/stores/cart-store"
import { useWishlistStore } from "@/lib/stores/wishlist-store"

// Client-side only component to prevent hydration mismatch
function WishlistCount() {
  const { getWishlistItemCount } = useWishlistStore()
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  if (!mounted) {
    return <span>(0)</span>
  }

  return <span>({getWishlistItemCount()})</span>
}

export default function ShopPage() {
  const [products, setProducts] = useState<Product[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchTerm, setSearchTerm] = useState("")
  const [selectedCategory, setSelectedCategory] = useState<number | undefined>()
  const [priceRange, setPriceRange] = useState<[number, number]>([0, 1000])
  const [sortBy, setSortBy] = useState("name")
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid")
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [total, setTotal] = useState(0)
  
  // Feature filters
  const [inStock, setInStock] = useState<boolean | undefined>(undefined)
  const [onSale, setOnSale] = useState<boolean | undefined>(true)
  const [isDigital, setIsDigital] = useState<boolean | undefined>(undefined)

  // Cart store
  const { isOpen, setCartOpen, getItemCount, loadCart } = useCartStore()
  
  // Wishlist store
  const { loadWishlists, getWishlistItemCount } = useWishlistStore()

  const loadProducts = async (page: number = 1, search: string = "", categoryId?: number, priceRange?: [number, number], sortBy?: string, featureFilters?: {inStock?: boolean, onSale?: boolean, isDigital?: boolean}) => {
    try {
      setLoading(true)
      setError(null)
      
      const filters: any = {
        page,
        limit: 12,
        search: search || undefined,
        category_id: categoryId,
        sort_by: sortBy || "created_at",
        sort_order: "desc"
      }
      
      // Add price filters if provided
      if (priceRange) {
        filters.min_price = priceRange[0]
        filters.max_price = priceRange[1]
      }
      
      // Add feature filters if provided
      if (featureFilters) {
        if (featureFilters.inStock !== undefined) filters.in_stock = featureFilters.inStock
        if (featureFilters.onSale !== undefined) filters.on_sale = featureFilters.onSale
        if (featureFilters.isDigital !== undefined) filters.is_digital = featureFilters.isDigital
      }
      
      const response = await productApi.getProducts(filters)
      setProducts(response.data.products)
      setTotalPages(response.data.total_pages)
      setTotal(response.data.total)
      setCurrentPage(response.data.page)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load products")
    } finally {
      setLoading(false)
    }
  }

  const loadCategories = async () => {
    try {
      const response = await productApi.getCategories({
        page: 1,
        limit: 10
      })
      setCategories(response.data.categories)
    } catch (err) {
      console.error("Failed to load categories:", err)
    }
  }

  useEffect(() => {
    loadProducts()
    loadCategories()
    // Load cart to show correct item count
    loadCart().catch(error => {
      console.warn('Cart loading failed:', error)
      // Continue without cart functionality
    })
    // Try to load wishlists, but don't fail if it doesn't work
    loadWishlists().catch(error => {
      console.warn('Wishlist loading failed:', error)
      // Continue without wishlist functionality
    })
  }, [loadWishlists, loadCart])

  const handleSearch = (value: string) => {
    setSearchTerm(value)
    setCurrentPage(1)
    loadProducts(1, value, selectedCategory, priceRange, sortBy, {inStock, onSale, isDigital})
  }

  const handleCategoryFilter = (categoryId: number | undefined) => {
    setSelectedCategory(categoryId)
    setCurrentPage(1)
    loadProducts(1, searchTerm, categoryId, priceRange, sortBy, {inStock, onSale, isDigital})
  }

  const handlePriceFilter = (range: [number, number]) => {
    setPriceRange(range)
    setCurrentPage(1)
    loadProducts(1, searchTerm, selectedCategory, range, sortBy, {inStock, onSale, isDigital})
  }

  const handleInStockFilter = (value: boolean | undefined) => {
    setInStock(value)
    setCurrentPage(1)
    loadProducts(1, searchTerm, selectedCategory, priceRange, sortBy, {inStock: value, onSale, isDigital})
  }

  const handleOnSaleFilter = (value: boolean | undefined) => {
    setOnSale(value)
    setCurrentPage(1)
    loadProducts(1, searchTerm, selectedCategory, priceRange, sortBy, {inStock, onSale: value, isDigital})
  }

  const handleIsDigitalFilter = (value: boolean | undefined) => {
    setIsDigital(value)
    setCurrentPage(1)
    loadProducts(1, searchTerm, selectedCategory, priceRange, sortBy, {inStock, onSale, isDigital: value})
  }

  const handleSort = (sort: string) => {
    setSortBy(sort)
    setCurrentPage(1)
    
    // Map frontend sort values to backend sort values
    let backendSort = "created_at"
    switch (sort) {
      case "price-low":
        backendSort = "price"
        break
      case "price-high":
        backendSort = "price"
        break
      case "name":
        backendSort = "name"
        break
      case "newest":
        backendSort = "created_at"
        break
      default:
        backendSort = "created_at"
    }
    
    loadProducts(1, searchTerm, selectedCategory, priceRange, backendSort, {inStock, onSale, isDigital})
  }

  // No need for client-side filtering since we're doing it on the backend
  const filteredProducts = products

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold tracking-tight">Shop</h1>
              <p className="text-muted-foreground">
                Discover amazing products at great prices
              </p>
            </div>
            <div className="flex items-center gap-2">
              <Button 
                variant="outline"
                asChild
                className="flex items-center gap-2"
              >
                <Link href="/wishlist">
                  <Heart className="h-4 w-4" />
                  Wishlist <WishlistCount />
                </Link>
              </Button>
              <Button 
                onClick={() => setCartOpen(true)}
                className="flex items-center gap-2"
              >
                <ShoppingCart className="h-4 w-4" />
                Cart ({getItemCount()})
              </Button>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-6">
        <div className="flex flex-col lg:flex-row gap-6">
          {/* Sidebar Filters */}
          <div className="lg:w-64 space-y-4">
            <ProductFilters
              categories={categories}
              selectedCategory={selectedCategory}
              onCategoryChange={handleCategoryFilter}
              priceRange={priceRange}
              onPriceChange={handlePriceFilter}
              inStock={inStock}
              onInStockChange={handleInStockFilter}
              onSale={onSale}
              onOnSaleChange={handleOnSaleFilter}
              isDigital={isDigital}
              onIsDigitalChange={handleIsDigitalFilter}
            />
          </div>

          {/* Main Content */}
          <div className="flex-1 space-y-6">
            {/* Search and Controls */}
            <Card>
              <CardContent className="p-4">
                <div className="flex flex-col sm:flex-row gap-4">
                  <div className="relative flex-1">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                    <Input
                      placeholder="Search products..."
                      value={searchTerm}
                      onChange={(e) => handleSearch(e.target.value)}
                      className="pl-10"
                    />
                  </div>
                  <div className="flex items-center gap-2">
                    <select
                      value={sortBy}
                      onChange={(e) => handleSort(e.target.value)}
                      className="px-3 py-2 border border-input bg-background rounded-md text-sm"
                    >
                      <option value="name">Sort by Name</option>
                      <option value="price-low">Price: Low to High</option>
                      <option value="price-high">Price: High to Low</option>
                      <option value="newest">Newest First</option>
                    </select>
                    <div className="flex border rounded-md">
                      <Button
                        variant={viewMode === "grid" ? "default" : "ghost"}
                        size="sm"
                        onClick={() => setViewMode("grid")}
                        className="rounded-r-none"
                      >
                        <Grid className="h-4 w-4" />
                      </Button>
                      <Button
                        variant={viewMode === "list" ? "default" : "ghost"}
                        size="sm"
                        onClick={() => setViewMode("list")}
                        className="rounded-l-none"
                      >
                        <List className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Results Header */}
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Badge variant="secondary" className="flex items-center gap-1">
                  <Package className="h-3 w-3" />
                  {total} products
                </Badge>
                {selectedCategory && (
                  <Badge variant="outline">
                    {categories.find(c => c.id === selectedCategory)?.name}
                  </Badge>
                )}
              </div>
            </div>

            {/* Products Grid/List */}
            {loading ? (
              <div className="flex items-center justify-center h-64">
                <div className="text-muted-foreground">Loading products...</div>
              </div>
            ) : error ? (
              <div className="text-center py-8">
                <div className="text-destructive mb-4">{error}</div>
                <Button onClick={() => loadProducts()}>Try Again</Button>
              </div>
            ) : filteredProducts.length === 0 ? (
              <div className="text-center py-8">
                <Package className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
                <div className="text-muted-foreground mb-4">
                  No products found. Try adjusting your filters.
                </div>
              </div>
            ) : (
              <div className={
                viewMode === "grid" 
                  ? "grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6"
                  : "space-y-4"
              }>
                {filteredProducts.map((product) => (
                  <ProductCard
                    key={product.id}
                    product={product}
                    viewMode={viewMode}
                    categories={categories}
                  />
                ))}
              </div>
            )}

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex items-center justify-center space-x-2">
                <Button
                  variant="outline"
                  onClick={() => {
                    const newPage = Math.max(1, currentPage - 1)
                    setCurrentPage(newPage)
                    loadProducts(newPage, searchTerm, selectedCategory, priceRange, sortBy, {inStock, onSale, isDigital})
                  }}
                  disabled={currentPage <= 1 || loading}
                >
                  Previous
                </Button>
                <span className="text-sm text-muted-foreground">
                  Page {currentPage} of {totalPages}
                </span>
                <Button
                  variant="outline"
                  onClick={() => {
                    const newPage = Math.min(totalPages, currentPage + 1)
                    setCurrentPage(newPage)
                    loadProducts(newPage, searchTerm, selectedCategory, priceRange, sortBy, {inStock, onSale, isDigital})
                  }}
                  disabled={currentPage >= totalPages || loading}
                >
                  Next
                </Button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Cart Drawer */}
      <CartDrawer open={isOpen} onOpenChange={setCartOpen} />
    </div>
  )
}