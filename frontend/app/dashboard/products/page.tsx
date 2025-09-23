"use client"

import { useState, useEffect, useCallback } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { 
  Plus, 
  Search, 
  Package,
  Filter,
  Download
} from "lucide-react"
import { productApi } from "@/lib/api"
import { Product, Category, CreateProductRequest, UpdateProductRequest, ProductFilters } from "@/lib/types"
import { ProductForm } from "./components/product-form"
import { ProductTable } from "./components/product-table"
import { ProductVariantManager } from "./components/product-variant-manager"
import { LoadingState } from "@/components/ui/loading"

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [editingProduct, setEditingProduct] = useState<Product | null>(null)
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [total, setTotal] = useState(0)
  const [showAdvancedFilters, setShowAdvancedFilters] = useState(false)
  const [showVariantManager, setShowVariantManager] = useState(false)
  const [selectedProductForVariants, setSelectedProductForVariants] = useState<Product | null>(null)
  
  // Search input state (separate from filters to prevent clearing)
  const [searchInput, setSearchInput] = useState("")
  
  // Filter states
  const [filters, setFilters] = useState<ProductFilters>({
    page: 1,
    limit: 10,
    search: "",
    category_id: undefined,
    is_active: undefined,
    is_digital: undefined,
    min_price: undefined,
    max_price: undefined,
    in_stock: undefined,
    tags: [],
    sort_by: "created_at",
    sort_order: "desc"
  })

  const loadProducts = async (newFilters?: Partial<ProductFilters>) => {
    try {
      setLoading(true)
      setError(null)
      const updatedFilters = { ...filters, ...newFilters }
      const response = await productApi.getProducts(updatedFilters)
      setProducts(response.data?.products || [])
      setTotalPages(response.data?.total_pages || 1)
      setTotal(response.data?.total || 0)
      setCurrentPage(response.data?.page || 1)
      setFilters(updatedFilters)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load products")
      // Reset products on error to prevent null reference
      setProducts([])
      setTotalPages(1)
      setTotal(0)
      setCurrentPage(1)
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
      setCategories(response.data?.categories || [])
    } catch (err) {
      console.error("Failed to load categories:", err)
    }
  }

  useEffect(() => {
    loadProducts()
    loadCategories()
  }, [])

  // Debounced search function
  const debouncedSearch = useCallback(
    (() => {
      let timeoutId: NodeJS.Timeout
      return (value: string) => {
        clearTimeout(timeoutId)
        timeoutId = setTimeout(() => {
          loadProducts({ search: value, page: 1 })
        }, 300) // 300ms delay
      }
    })(),
    []
  )

  const handleSearchInput = (value: string) => {
    setSearchInput(value)
    debouncedSearch(value)
  }

  const handleSearch = (value: string) => {
    setSearchInput(value)
    loadProducts({ search: value, page: 1 })
  }

  const handleFilterChange = (key: keyof ProductFilters, value: any) => {
    loadProducts({ [key]: value, page: 1 })
  }

  const handlePageChange = (page: number) => {
    loadProducts({ page })
  }

  const clearFilters = () => {
    setSearchInput("")
    const clearedFilters: ProductFilters = {
      page: 1,
      limit: 10,
      search: "",
      category_id: undefined,
      is_active: undefined,
      is_digital: undefined,
      min_price: undefined,
      max_price: undefined,
      in_stock: undefined,
      tags: [],
      sort_by: "created_at",
      sort_order: "desc"
    }
    loadProducts(clearedFilters)
  }

  const handleCreateProduct = async (productData: CreateProductRequest) => {
    try {
      await productApi.createProduct(productData)
      setShowForm(false)
      loadProducts()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create product")
    }
  }

  const handleUpdateProduct = async (id: number, productData: UpdateProductRequest) => {
    try {
      await productApi.updateProduct(id, productData)
      setEditingProduct(null)
      setShowForm(false)
      loadProducts()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update product")
    }
  }

  const handleDeleteProduct = async (id: number) => {
    if (!confirm("Are you sure you want to delete this product?")) return
    
    try {
      await productApi.deleteProduct(id)
      loadProducts()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete product")
    }
  }

  const handleEdit = (product: Product) => {
    setEditingProduct(product)
    setShowForm(true)
  }

  const handleCancel = () => {
    setShowForm(false)
    setEditingProduct(null)
  }

  const handleManageVariants = (product: Product) => {
    setSelectedProductForVariants(product)
    setShowVariantManager(true)
  }

  const handleVariantManagerCancel = () => {
    setShowVariantManager(false)
    setSelectedProductForVariants(null)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold tracking-tight">Products</h1>
          <p className="text-muted-foreground">
            Manage your product catalog and inventory
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <Button variant="outline" className="flex items-center gap-2">
            <Download className="h-4 w-4" />
            <span className="hidden sm:inline">Export</span>
          </Button>
          <Button onClick={() => setShowForm(true)} className="flex items-center gap-2">
            <Plus className="h-4 w-4" />
            <span className="hidden sm:inline">Add Product</span>
            <span className="sm:hidden">Add</span>
          </Button>
        </div>
      </div>

      {/* Error Alert */}
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {/* Search and Filters */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="text-lg">Search & Filter</CardTitle>
            <div className="flex items-center space-x-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowAdvancedFilters(!showAdvancedFilters)}
              >
                <Filter className="h-4 w-4 mr-2" />
                {showAdvancedFilters ? 'Hide' : 'Show'} Filters
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={clearFilters}
              >
                Clear All
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {/* Basic Search */}
            <div className="flex flex-col sm:flex-row sm:items-center gap-4">
              <div className="relative flex-1 max-w-sm">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input
                  placeholder="Search products..."
                  value={searchInput}
                  onChange={(e) => handleSearchInput(e.target.value)}
                  className="pl-10"
                />
              </div>
              <div className="flex items-center space-x-2">
                <select
                  value={filters.category_id || ""}
                  onChange={(e) => handleFilterChange('category_id', e.target.value ? parseInt(e.target.value) : undefined)}
                  className="px-3 py-2 border border-input bg-background rounded-md text-sm min-w-[150px]"
                >
                  <option value="">All Categories</option>
                  {categories.map((category) => (
                    <option key={category.id} value={category.id}>
                      {category.name}
                    </option>
                  ))}
                </select>
              </div>
              <Badge variant="secondary" className="flex items-center gap-1">
                <Package className="h-3 w-3" />
                {total} products
              </Badge>
            </div>

            {/* Advanced Filters */}
            {showAdvancedFilters && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 pt-4 border-t">
                {/* Status Filter */}
                <div className="space-y-2">
                  <label className="text-sm font-medium">Status</label>
                  <select
                    value={filters.is_active === undefined ? "" : filters.is_active.toString()}
                    onChange={(e) => handleFilterChange('is_active', e.target.value === "" ? undefined : e.target.value === "true")}
                    className="w-full px-3 py-2 border border-input bg-background rounded-md text-sm"
                  >
                    <option value="">All Status</option>
                    <option value="true">Active</option>
                    <option value="false">Inactive</option>
                  </select>
                </div>

                {/* Digital Filter */}
                <div className="space-y-2">
                  <label className="text-sm font-medium">Type</label>
                  <select
                    value={filters.is_digital === undefined ? "" : filters.is_digital.toString()}
                    onChange={(e) => handleFilterChange('is_digital', e.target.value === "" ? undefined : e.target.value === "true")}
                    className="w-full px-3 py-2 border border-input bg-background rounded-md text-sm"
                  >
                    <option value="">All Types</option>
                    <option value="false">Physical</option>
                    <option value="true">Digital</option>
                  </select>
                </div>

                {/* Stock Filter */}
                <div className="space-y-2">
                  <label className="text-sm font-medium">Stock</label>
                  <select
                    value={filters.in_stock === undefined ? "" : filters.in_stock.toString()}
                    onChange={(e) => handleFilterChange('in_stock', e.target.value === "" ? undefined : e.target.value === "true")}
                    className="w-full px-3 py-2 border border-input bg-background rounded-md text-sm"
                  >
                    <option value="">All Stock</option>
                    <option value="true">In Stock</option>
                    <option value="false">Out of Stock</option>
                  </select>
                </div>

                {/* Sort By */}
                <div className="space-y-2">
                  <label className="text-sm font-medium">Sort By</label>
                  <select
                    value={filters.sort_by || ""}
                    onChange={(e) => handleFilterChange('sort_by', e.target.value)}
                    className="w-full px-3 py-2 border border-input bg-background rounded-md text-sm"
                  >
                    <option value="created_at">Date Created</option>
                    <option value="name">Name</option>
                    <option value="price">Price</option>
                    <option value="sku">SKU</option>
                  </select>
                </div>

                {/* Price Range */}
                <div className="space-y-2">
                  <label className="text-sm font-medium">Min Price</label>
                  <Input
                    type="number"
                    placeholder="0.00"
                    value={filters.min_price || ""}
                    onChange={(e) => handleFilterChange('min_price', e.target.value ? parseFloat(e.target.value) : undefined)}
                    className="w-full"
                  />
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">Max Price</label>
                  <Input
                    type="number"
                    placeholder="999.99"
                    value={filters.max_price || ""}
                    onChange={(e) => handleFilterChange('max_price', e.target.value ? parseFloat(e.target.value) : undefined)}
                    className="w-full"
                  />
                </div>

                {/* Sort Order */}
                <div className="space-y-2">
                  <label className="text-sm font-medium">Order</label>
                  <select
                    value={filters.sort_order || "desc"}
                    onChange={(e) => handleFilterChange('sort_order', e.target.value as 'asc' | 'desc')}
                    className="w-full px-3 py-2 border border-input bg-background rounded-md text-sm"
                  >
                    <option value="desc">Descending</option>
                    <option value="asc">Ascending</option>
                  </select>
                </div>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Products Table */}
      <Card>
        <CardHeader>
          <CardTitle>All Products</CardTitle>
          <CardDescription>
            Manage your product catalog, pricing, and inventory
          </CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <LoadingState
              type="skeleton"
              viewMode="list"
              itemCount={5}
              text="Loading Products"
            />
          ) : (
            <ProductTable
              products={products}
              categories={categories}
              onEdit={handleEdit}
              onDelete={handleDeleteProduct}
              onManageVariants={handleManageVariants}
              currentPage={currentPage}
              totalPages={totalPages}
              onPageChange={handlePageChange}
            />
          )}
        </CardContent>
      </Card>

      {/* Product Form Modal */}
      {showForm && (
        <ProductForm
          product={editingProduct}
          categories={categories}
          onSave={editingProduct ? 
            (data: UpdateProductRequest) => handleUpdateProduct(editingProduct.id, data) : 
            (data: CreateProductRequest) => handleCreateProduct(data)
          }
          onCancel={handleCancel}
        />
      )}

      {/* Variant Management Modal */}
      {showVariantManager && selectedProductForVariants && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
          <div className="bg-background rounded-lg max-w-4xl w-full max-h-[90vh] overflow-y-auto">
            <div className="p-6">
              <div className="flex items-center justify-between mb-6">
                <div>
                  <h2 className="text-2xl font-bold">Manage Variants</h2>
                  <p className="text-muted-foreground">
                    Product: {selectedProductForVariants.name}
                  </p>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={handleVariantManagerCancel}
                >
                  ✕
                </Button>
              </div>
              <ProductVariantManager
                productId={selectedProductForVariants.id}
                onVariantChange={() => {
                  // Optionally refresh products or show notification
                }}
              />
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
