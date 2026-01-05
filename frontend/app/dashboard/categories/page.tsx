"use client";

import { useState, useEffect, useCallback } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { toast } from "sonner";
import {
  Plus,
  Search,
  Package,
  Filter,
  SortAsc,
  SortDesc,
  Grid3X3,
  List,
  Eye,
  EyeOff,
  Folder,
  FolderOpen,
} from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { productApi } from "@/lib/apiFunctions";
import {
  Category,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  CategoryFilters,
} from "@/lib/types";
import { CategoryForm } from "./components/category-form";
import { CategoryTable } from "./components/category-table";
import { ParentCategorySelector } from "./components/parent-category-selector";
import { LoadingState } from "@/components/ui/loading";

export default function CategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [parentOptions, setParentOptions] = useState<Category[]>([]);

  const [loading, setLoading] = useState(true);

  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);
  const [showAdvancedFilters, setShowAdvancedFilters] = useState(false);
  const [viewMode, setViewMode] = useState<"grid" | "list">("list");

  // Search input state (separate from filters to prevent clearing)
  const [searchInput, setSearchInput] = useState("");

  // Filter states
  const [filters, setFilters] = useState<CategoryFilters>({
    page: 1,
    limit: 10,
    search: "",
    parent_id: undefined,
    is_active: undefined,
  });

  const loadCategories = async (newFilters?: Partial<CategoryFilters>) => {
    try {
      setLoading(true);
      setError(null);
      const updatedFilters = { ...filters, ...newFilters };
      const response = await productApi.getCategories(updatedFilters);
      setCategories(response.data?.categories || []);
      setTotalPages(response.data?.total_pages || 1);
      setTotal(response.data?.total || 0);
      setCurrentPage(response.data?.page || 1);
      setFilters(updatedFilters);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to load categories"
      );
      // Reset categories on error to prevent null reference
      setCategories([]);
      setTotalPages(1);
      setTotal(0);
      setCurrentPage(1);
    } finally {
      setLoading(false);
    }
  };

  const searchParentCategories = useCallback(async (query: string) => {
    try {
      const response = await productApi.getCategories({
        search: query,
        limit: 20,
        page: 1,
      });
      setParentOptions(response.data?.categories || []);
    } catch (err) {
      console.error("Failed to search parent categories:", err);
    }
  }, []);

  useEffect(() => {
    loadCategories();
  }, []);

  // Debounced search function
  const debouncedSearch = useCallback(
    (() => {
      let timeoutId: NodeJS.Timeout;
      return (value: string) => {
        clearTimeout(timeoutId);
        timeoutId = setTimeout(() => {
          loadCategories({ search: value, page: 1 });
        }, 300); // 300ms delay
      };
    })(),
    []
  );

  const handleSearchInput = (value: string) => {
    setSearchInput(value);
    debouncedSearch(value);
  };

  const handleFilterChange = (key: keyof CategoryFilters, value: any) => {
    loadCategories({ [key]: value, page: 1 });
  };

  const handlePageChange = (page: number) => {
    loadCategories({ page });
  };

  const clearFilters = () => {
    setSearchInput("");
    const clearedFilters: CategoryFilters = {
      page: 1,
      limit: 10,
      search: "",
      parent_id: undefined,
      is_active: undefined,
    };
    loadCategories(clearedFilters);
  };

  const handleCreateCategory = async (categoryData: CreateCategoryRequest) => {
    try {
      await productApi.createCategory(categoryData);
      setShowForm(false);
      loadCategories();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to create category"
      );
    }
  };

  const handleUpdateCategory = async (
    id: number,
    categoryData: UpdateCategoryRequest
  ) => {
    try {
      await productApi.updateCategory(id, categoryData);
      setEditingCategory(null);
      setShowForm(false);
      loadCategories();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to update category"
      );
    }
  };

  const handleDeleteCategory = async (id: number) => {
    if (!confirm("Are you sure you want to delete this category?")) return;

    try {
      // Image cleanup is handled automatically by the API route handler
      // via cleanupEntityImages() - no need for manual deletion here
      await productApi.deleteCategory(id);
      toast.success("Category deleted successfully");
      loadCategories();
    } catch (error) {
      console.error("Failed to delete category:", error);
      toast.error("Failed to delete category");
    }
  };

  const handleEdit = (category: Category) => {
    setEditingCategory(category);
    setShowForm(true);
  };

  const handleCancel = () => {
    setShowForm(false);
    setEditingCategory(null);
  };

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
        <Button
          onClick={() => setShowForm(true)}
          className="flex items-center gap-2"
        >
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

      {/* Enhanced Search and Filters */}
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
                {showAdvancedFilters ? "Hide" : "Show"} Filters
              </Button>
              <Button variant="outline" size="sm" onClick={clearFilters}>
                Clear All
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {/* Basic Search and View Controls */}
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-4 flex-1">
                <div className="relative flex-1 max-w-sm">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                  <Input
                    placeholder="Search categories..."
                    value={searchInput}
                    onChange={(e) => handleSearchInput(e.target.value)}
                    className="pl-10"
                  />
                </div>
                <Badge variant="secondary" className="flex items-center gap-1">
                  <Package className="h-3 w-3" />
                  {total} categories
                </Badge>
              </div>

              {/* View Mode Toggle */}
              <div className="flex items-center space-x-1">
                <Button
                  variant={viewMode === "list" ? "default" : "outline"}
                  size="sm"
                  onClick={() => setViewMode("list")}
                  className="h-8"
                >
                  <List className="h-4 w-4" />
                </Button>
                <Button
                  variant={viewMode === "grid" ? "default" : "outline"}
                  size="sm"
                  onClick={() => setViewMode("grid")}
                  className="h-8"
                >
                  <Grid3X3 className="h-4 w-4" />
                </Button>
              </div>
            </div>

            {/* Advanced Filters */}
            {showAdvancedFilters && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 pt-4 border-t">
                {/* Status Filter */}
                <div className="space-y-2">
                  <label className="text-sm font-medium">Status</label>
                  <Select
                    value={
                      filters.is_active === undefined
                        ? "all"
                        : filters.is_active.toString()
                    }
                    onValueChange={(value) =>
                      handleFilterChange(
                        "is_active",
                        value === "all" ? undefined : value === "true"
                      )
                    }
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select status" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All Status</SelectItem>
                      <SelectItem value="true">Active</SelectItem>
                      <SelectItem value="false">Inactive</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                {/* Parent Category Filter */}
                <div className="space-y-2">
                  <label className="text-sm font-medium">Parent Category</label>
                  <ParentCategorySelector
                    value={filters.parent_id}
                    onChange={(value) => handleFilterChange("parent_id", value)}
                    onSearch={searchParentCategories}
                    options={parentOptions}
                  />
                </div>

                {/* Results Summary */}
                <div className="space-y-2">
                  <label className="text-sm font-medium">Results</label>
                  <div className="text-sm text-muted-foreground">
                    {total} total categories
                    {filters.search && ` matching "${filters.search}"`}
                    {filters.is_active !== undefined &&
                      ` (${filters.is_active ? "Active" : "Inactive"})`}
                    {filters.parent_id !== undefined &&
                      ` (${filters.parent_id === 0
                        ? "Root"
                        : parentOptions.find(
                          (c) => c.id === filters.parent_id
                        )?.name || "Parent"
                      })`}
                  </div>
                </div>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Categories Table */}
      <Card>
        <CardHeader>
          <div>
            <CardTitle>All Categories</CardTitle>
            <CardDescription>
              Manage your product categories and their settings
            </CardDescription>
          </div>
        </CardHeader>
        <CardContent>
          {loading ? (
            <LoadingState
              type="skeleton"
              viewMode={viewMode}
              itemCount={viewMode === "list" ? 5 : 6}
              text="Loading Categories"
            />
          ) : (
            <CategoryTable
              categories={categories}
              onEdit={handleEdit}
              onDelete={handleDeleteCategory}
              currentPage={currentPage}
              totalPages={totalPages}
              onPageChange={handlePageChange}
              viewMode={viewMode}
              loading={loading}
            />
          )}
        </CardContent>
      </Card>

      {/* Category Form Modal */}
      {showForm && (
        <CategoryForm
          category={editingCategory}
          onSave={
            editingCategory
              ? (data: UpdateCategoryRequest) =>
                handleUpdateCategory(editingCategory.id, data)
              : (data: CreateCategoryRequest) => handleCreateCategory(data)
          }
          onCancel={handleCancel}
        />
      )}
    </div>
  );
}