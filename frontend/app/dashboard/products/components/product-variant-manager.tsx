"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Alert, AlertDescription } from "@/components/ui/alert"
import {
  Plus,
  Package,
  Edit,
  Trash2,
  MoreHorizontal,
  Eye,
  EyeOff,
  DollarSign,
  Weight,
  ArrowUpDown
} from "lucide-react"
import { ProductVariant, CreateProductVariantRequest, UpdateProductVariantRequest } from "@/lib/types"
import { productApi } from "@/lib/apiFunctions"
import { formatPrice } from "@/lib/currency"
import { ProductVariantForm } from "./product-variant-form"
import { LoadingState } from "@/components/ui/loading"

interface ProductVariantManagerProps {
  productId: number
  onVariantChange?: () => void
}

export function ProductVariantManager({ productId, onVariantChange }: ProductVariantManagerProps) {
  const [variants, setVariants] = useState<ProductVariant[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [editingVariant, setEditingVariant] = useState<ProductVariant | null>(null)
  const [deletingId, setDeletingId] = useState<number | null>(null)

  const loadVariants = async () => {
    try {
      setLoading(true)
      setError(null)
      const variantsData = await productApi.getProductVariants(productId)
      setVariants(Array.isArray(variantsData) ? variantsData : [])
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load variants")
      setVariants([]) // Ensure variants is always an array
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (productId) {
      loadVariants()
    }
  }, [productId])

  const handleCreateVariant = async (variantData: CreateProductVariantRequest | UpdateProductVariantRequest) => {
    try {
      await productApi.createProductVariant(productId, variantData as CreateProductVariantRequest)
      await loadVariants()
      setShowForm(false)
      onVariantChange?.()
    } catch (err) {
      throw err
    }
  }

  const handleUpdateVariant = async (variantData: CreateProductVariantRequest | UpdateProductVariantRequest) => {
    if (!editingVariant) return

    try {
      await productApi.updateProductVariant(editingVariant.id, variantData as UpdateProductVariantRequest)
      await loadVariants()
      setEditingVariant(null)
      onVariantChange?.()
    } catch (err) {
      throw err
    }
  }

  const handleDeleteVariant = async (id: number) => {
    setDeletingId(id)
    try {
      await productApi.deleteProductVariant(id)
      await loadVariants()
      onVariantChange?.()
    } catch (err) {
      console.error("Failed to delete variant:", err)
    } finally {
      setDeletingId(null)
    }
  }

  const handleEdit = (variant: ProductVariant) => {
    setEditingVariant(variant)
    setShowForm(true)
  }

  const handleCancel = () => {
    setShowForm(false)
    setEditingVariant(null)
  }

  if (showForm) {
    return (
      <ProductVariantForm
        productId={productId}
        variant={editingVariant}
        onSave={editingVariant ? handleUpdateVariant : handleCreateVariant}
        onCancel={handleCancel}
        loading={loading}
      />
    )
  }

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Package className="h-5 w-5" />
            Product Variants
          </CardTitle>
          <CardDescription>
            Manage different variations of this product
          </CardDescription>
        </CardHeader>
        <CardContent>
          <LoadingState type="skeleton" viewMode="list" itemCount={3} text="Loading variants..." />
        </CardContent>
      </Card>
    )
  }

  if (error) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Package className="h-5 w-5" />
            Product Variants
          </CardTitle>
          <CardDescription>
            Manage different variations of this product
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Package className="h-5 w-5" />
              Product Variants
            </CardTitle>
            <CardDescription>
              Manage different variations of this product
            </CardDescription>
          </div>
          <Button onClick={() => setShowForm(true)} className="flex items-center gap-2">
            <Plus className="h-4 w-4" />
            Add Variant
          </Button>
        </div>
      </CardHeader>

      <CardContent>
        {!variants || variants.length === 0 ? (
          <div className="text-center py-8">
            <Package className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
            <h3 className="text-lg font-medium mb-2">No variants found</h3>
            <p className="text-muted-foreground mb-4">
              Create your first product variant to get started
            </p>
            <Button onClick={() => setShowForm(true)} className="flex items-center gap-2">
              <Plus className="h-4 w-4" />
              Add First Variant
            </Button>
          </div>
        ) : (
          <div className="space-y-4">
            {/* Summary Stats */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div className="flex items-center space-x-2">
                <Package className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm text-muted-foreground">Total:</span>
                <span className="font-medium">{variants.length}</span>
              </div>
              <div className="flex items-center space-x-2">
                <Eye className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm text-muted-foreground">Active:</span>
                <span className="font-medium">{variants?.filter(v => v.is_active).length || 0}</span>
              </div>
              <div className="flex items-center space-x-2">
                <DollarSign className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm text-muted-foreground">Avg Price:</span>
                <span className="font-medium">
                  {formatPrice((variants?.reduce((sum, v) => sum + v.price, 0) || 0) / (variants?.length || 1))}
                </span>
              </div>
              <div className="flex items-center space-x-2">
                <Weight className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm text-muted-foreground">Avg Weight:</span>
                <span className="font-medium">
                  {((variants?.reduce((sum, v) => sum + v.weight, 0) || 0) / (variants?.length || 1)).toFixed(3)} kg
                </span>
              </div>
            </div>

            {/* Variants List */}
            <div className="space-y-2">
              {(variants || [])
                .sort((a, b) => a.position - b.position)
                .map((variant) => (
                  <div
                    key={variant.id}
                    className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50"
                  >
                    <div className="flex-1">
                      <div className="flex items-center gap-3">
                        <h4 className="font-medium">{variant.name}</h4>
                        <Badge variant={variant.is_active ? "default" : "secondary"}>
                          {variant.is_active ? (
                            <>
                              <Eye className="h-3 w-3 mr-1" />
                              Active
                            </>
                          ) : (
                            <>
                              <EyeOff className="h-3 w-3 mr-1" />
                              Inactive
                            </>
                          )}
                        </Badge>
                      </div>
                      <div className="flex items-center gap-4 mt-1 text-sm text-muted-foreground">
                        <span>SKU: {variant.sku}</span>
                        <span>•</span>
                        <span>Price: {formatPrice(variant.price)}</span>
                        {variant.compare_price > 0 && (
                          <>
                            <span>•</span>
                            <span>Compare: {formatPrice(variant.compare_price)}</span>
                          </>
                        )}
                        {variant.weight > 0 && (
                          <>
                            <span>•</span>
                            <span>Weight: {variant.weight} kg</span>
                          </>
                        )}
                        <span>•</span>
                        <span>Position: {variant.position}</span>
                      </div>
                    </div>

                    <div className="flex items-center gap-2">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleEdit(variant)}
                      >
                        <Edit className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDeleteVariant(variant.id)}
                        disabled={deletingId === variant.id}
                        className="text-red-600 hover:text-red-700"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
