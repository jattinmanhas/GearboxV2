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
  Package
} from "lucide-react"
import { productApi } from "@/lib/api"
import { Category, CreateCategoryRequest, UpdateCategoryRequest } from "@/lib/types"
import { CategoryForm } from "./components/category-form"
import { CategoryTable } from "./components/category-table"

export default function CategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchTerm, setSearchTerm] = useState("")
  const [showForm, setShowForm] = useState(false)
  const [editingCategory, setEditingCategory] = useState<Category | null>(null)
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [total, setTotal] = useState(0)

  const loadCategories = async (page: number = 1, search: string = "") => {
    try {
      setLoading(true)
      setError(null)
      const response = await productApi.getCategories(page, 10, search)
      setCategories(response.categories)
      setTotalPages(response.total_pages)
      setTotal(response.total)
      setCurrentPage(response.page)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load categories")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadCategories()
  }, [])

  const handleSearch = (value: string) => {
    setSearchTerm(value)
    setCurrentPage(1)
    loadCategories(1, value)
  }

  const handleCreateCategory = async (categoryData: CreateCategoryRequest) => {
    try {
      await productApi.createCategory(categoryData)
      setShowForm(false)
      loadCategories(currentPage, searchTerm)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create category")
    }
  }

  const handleUpdateCategory = async (id: number, categoryData: UpdateCategoryRequest) => {
    try {
      await productApi.updateCategory(id, categoryData)
      setEditingCategory(null)
      loadCategories(currentPage, searchTerm)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update category")
    }
  }

  const handleDeleteCategory = async (id: number) => {
    if (!confirm("Are you sure you want to delete this category?")) return
    
    try {
      await productApi.deleteCategory(id)
      loadCategories(currentPage, searchTerm)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete category")
    }
  }

  const handleEdit = (category: Category) => {
    setEditingCategory(category)
    setShowForm(true)
  }

  const handleCancel = () => {
    setShowForm(false)
    setEditingCategory(null)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Categories</h1>
          <p className="text-muted-foreground">
            Manage product categories and their hierarchy
          </p>
        </div>
        <Button onClick={() => setShowForm(true)} className="flex items-center gap-2">
          <Plus className="h-4 w-4" />
          Add Category
        </Button>
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
                placeholder="Search categories..."
                value={searchTerm}
                onChange={(e) => handleSearch(e.target.value)}
                className="pl-10"
              />
            </div>
            <Badge variant="secondary" className="flex items-center gap-1">
              <Package className="h-3 w-3" />
              {total} categories
            </Badge>
          </div>
        </CardContent>
      </Card>

      {/* Categories Table */}
      <Card>
        <CardHeader>
          <CardTitle>All Categories</CardTitle>
          <CardDescription>
            Manage your product categories and their settings
          </CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="flex items-center justify-center h-32">
              <div className="text-muted-foreground">Loading categories...</div>
            </div>
          ) : (
            <CategoryTable
              categories={categories}
              onEdit={handleEdit}
              onDelete={handleDeleteCategory}
              currentPage={currentPage}
              totalPages={totalPages}
              onPageChange={setCurrentPage}
            />
          )}
        </CardContent>
      </Card>

      {/* Category Form Modal */}
      {showForm && (
        <CategoryForm
          category={editingCategory}
          onSave={editingCategory ? 
            (data: UpdateCategoryRequest) => handleUpdateCategory(editingCategory.id, data) : 
            (data: CreateCategoryRequest) => handleCreateCategory(data)
          }
          onCancel={handleCancel}
        />
      )}
    </div>
  )
}
