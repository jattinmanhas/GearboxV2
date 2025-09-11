"use client"

import { useState, useEffect } from "react"
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
import { Product, Category, CreateProductRequest, UpdateProductRequest } from "@/lib/types"
import { ProductForm } from "./components/product-form"
import { ProductTable } from "./components/product-table"

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchTerm, setSearchTerm] = useState("")
  const [selectedCategory, setSelectedCategory] = useState<number | undefined>()
  const [showForm, setShowForm] = useState(false)
  const [editingProduct, setEditingProduct] = useState<Product | null>(null)
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [total, setTotal] = useState(0)

  const loadProducts = async (page: number = 1, search: string = "", categoryId?: number) => {
    try {
      setLoading(true)
      setError(null)
      const response = await productApi.getProducts(page, 10, search, categoryId)
      setProducts(response.products)
      setTotalPages(response.total_pages)
      setTotal(response.total)
      setCurrentPage(response.page)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load products")
    } finally {
      setLoading(false)
    }
  }

  const loadCategories = async () => {
    try {
      const response = await productApi.getCategories(1, 100)
      setCategories(response.categories)
    } catch (err) {
      console.error("Failed to load categories:", err)
    }
  }

  useEffect(() => {
    loadProducts()
    loadCategories()
  }, [])

  const handleSearch = (value: string) => {
    setSearchTerm(value)
    setCurrentPage(1)
    loadProducts(1, value, selectedCategory)
  }

  const handleCategoryFilter = (categoryId: number | undefined) => {
    setSelectedCategory(categoryId)
    setCurrentPage(1)
    loadProducts(1, searchTerm, categoryId)
  }

  const handleCreateProduct = async (productData: CreateProductRequest) => {
    try {
      await productApi.createProduct(productData)
      setShowForm(false)
      loadProducts(currentPage, searchTerm, selectedCategory)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create product")
    }
  }

  const handleUpdateProduct = async (id: number, productData: UpdateProductRequest) => {
    try {
      await productApi.updateProduct(id, productData)
      setEditingProduct(null)
      loadProducts(currentPage, searchTerm, selectedCategory)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update product")
    }
  }

  const handleDeleteProduct = async (id: number) => {
    if (!confirm("Are you sure you want to delete this product?")) return
    
    try {
      await productApi.deleteProduct(id)
      loadProducts(currentPage, searchTerm, selectedCategory)
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

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Products</h1>
          <p className="text-muted-foreground">
            Manage your product catalog and inventory
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <Button variant="outline" className="flex items-center gap-2">
            <Download className="h-4 w-4" />
            Export
          </Button>
          <Button onClick={() => setShowForm(true)} className="flex items-center gap-2">
            <Plus className="h-4 w-4" />
            Add Product
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
          <CardTitle className="text-lg">Search & Filter</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center space-x-4">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
              <Input
                placeholder="Search products..."
                value={searchTerm}
                onChange={(e) => handleSearch(e.target.value)}
                className="pl-10"
              />
            </div>
            <div className="flex items-center space-x-2">
              <Filter className="h-4 w-4 text-muted-foreground" />
              <select
                value={selectedCategory || ""}
                onChange={(e) => handleCategoryFilter(e.target.value ? parseInt(e.target.value) : undefined)}
                className="px-3 py-2 border border-input bg-background rounded-md text-sm"
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
            <div className="flex items-center justify-center h-32">
              <div className="text-muted-foreground">Loading products...</div>
            </div>
          ) : (
            <ProductTable
              products={products}
              categories={categories}
              onEdit={handleEdit}
              onDelete={handleDeleteProduct}
              currentPage={currentPage}
              totalPages={totalPages}
              onPageChange={setCurrentPage}
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
    </div>
  )
}
