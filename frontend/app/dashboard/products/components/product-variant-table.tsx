"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { 
  Edit, 
  Trash2, 
  Eye, 
  EyeOff, 
  MoreHorizontal,
  Package,
  DollarSign,
  Weight,
  ArrowUpDown,
  Plus,
  AlertCircle,
  TrendingUp,
  TrendingDown,
  Hash,
  Settings,
  CheckCircle,
  XCircle,
  Sparkles
} from "lucide-react"
import { ProductVariant } from "@/lib/types"
import { formatPrice } from "@/lib/currency"
import { LoadingState } from "@/components/ui/loading"

interface ProductVariantTableProps {
  variants: ProductVariant[]
  onEdit: (variant: ProductVariant) => void
  onDelete: (id: number) => void
  onAdd: () => void
  loading?: boolean
  error?: string | null
}

export function ProductVariantTable({ 
  variants, 
  onEdit, 
  onDelete, 
  onAdd,
  loading = false,
  error = null
}: ProductVariantTableProps) {
  const [deletingId, setDeletingId] = useState<number | null>(null)

  const handleDelete = async (id: number) => {
    setDeletingId(id)
    try {
      await onDelete(id)
    } finally {
      setDeletingId(null)
    }
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
  }

  const getStatusBadge = (isActive: boolean) => {
    return (
      <Badge variant={isActive ? "default" : "secondary"}>
        {isActive ? (
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
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="border-2">
      <CardHeader className="pb-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-primary/10 rounded-lg">
              <Package className="h-5 w-5 text-primary" />
            </div>
            <div>
              <CardTitle className="text-lg">Product Variants</CardTitle>
              <CardDescription>
                Manage different variations of this product
              </CardDescription>
            </div>
          </div>
          <Button onClick={onAdd} className="flex items-center gap-2 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700">
            <Plus className="h-4 w-4" />
            Add Variant
          </Button>
        </div>
      </CardHeader>
      
      <CardContent>
        {!variants || variants.length === 0 ? (
          <div className="text-center py-12">
            <div className="p-4 bg-muted/30 rounded-full w-16 h-16 mx-auto mb-4 flex items-center justify-center">
              <Package className="h-8 w-8 text-muted-foreground" />
            </div>
            <h3 className="text-lg font-semibold mb-2">No variants found</h3>
            <p className="text-muted-foreground mb-6 max-w-sm mx-auto">
              Create your first product variant to expand your product options and reach more customers
            </p>
            <Button onClick={onAdd} className="flex items-center gap-2 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700">
              <Plus className="h-4 w-4" />
              Add First Variant
            </Button>
          </div>
        ) : (
          <div className="space-y-6">
            {/* Summary Stats */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 p-4 bg-muted/30 rounded-lg">
              <div className="flex items-center space-x-2">
                <Package className="h-4 w-4 text-blue-600" />
                <span className="text-sm text-muted-foreground">Total:</span>
                <span className="font-semibold text-blue-600">{variants.length}</span>
              </div>
              <div className="flex items-center space-x-2">
                <Eye className="h-4 w-4 text-green-600" />
                <span className="text-sm text-muted-foreground">Active:</span>
                <span className="font-semibold text-green-600">{variants.filter(v => v.is_active).length}</span>
              </div>
              <div className="flex items-center space-x-2">
                <DollarSign className="h-4 w-4 text-purple-600" />
                <span className="text-sm text-muted-foreground">Avg Price:</span>
                <span className="font-semibold text-purple-600">
                  {formatPrice(variants.reduce((sum, v) => sum + v.price, 0) / variants.length)}
                </span>
              </div>
              <div className="flex items-center space-x-2">
                <Weight className="h-4 w-4 text-orange-600" />
                <span className="text-sm text-muted-foreground">Avg Weight:</span>
                <span className="font-semibold text-orange-600">
                  {(variants.reduce((sum, v) => sum + v.weight, 0) / variants.length).toFixed(3)} kg
                </span>
              </div>
            </div>

            {/* Variants Grid */}
            <div className="grid gap-4">
              {variants
                .sort((a, b) => a.position - b.position)
                .map((variant) => {
                  const discountPercentage = variant.compare_price > 0 && variant.price > 0 && variant.compare_price > variant.price
                    ? (((variant.compare_price - variant.price) / variant.compare_price) * 100).toFixed(1)
                    : null

                  return (
                    <div
                      key={variant.id}
                      className="group relative p-6 border-2 rounded-xl hover:shadow-lg transition-all duration-200 bg-card hover:bg-muted/20"
                    >
                      <div className="flex items-start justify-between">
                        <div className="flex-1 space-y-4">
                          {/* Header */}
                          <div className="flex items-start justify-between">
                            <div className="space-y-1">
                              <div className="flex items-center gap-2">
                                <h3 className="text-lg font-semibold">{variant.name}</h3>
                                <Badge variant={variant.is_active ? "default" : "secondary"} className="text-xs">
                                  {variant.is_active ? (
                                    <div className="flex items-center gap-1">
                                      <CheckCircle className="h-3 w-3" />
                                      Active
                                    </div>
                                  ) : (
                                    <div className="flex items-center gap-1">
                                      <XCircle className="h-3 w-3" />
                                      Inactive
                                    </div>
                                  )}
                                </Badge>
                              </div>
                              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                <Hash className="h-3 w-3" />
                                <code className="bg-muted px-2 py-1 rounded text-xs font-mono">
                                  {variant.sku}
                                </code>
                              </div>
                            </div>
                            
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button variant="ghost" size="sm" className="h-8 w-8 p-0 opacity-0 group-hover:opacity-100 transition-opacity">
                                  <MoreHorizontal className="h-4 w-4" />
                                </Button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end">
                                <DropdownMenuItem onClick={() => onEdit(variant)}>
                                  <Edit className="mr-2 h-4 w-4" />
                                  Edit Variant
                                </DropdownMenuItem>
                                <DropdownMenuItem
                                  onClick={() => handleDelete(variant.id)}
                                  disabled={deletingId === variant.id}
                                  className="text-destructive"
                                >
                                  <Trash2 className="mr-2 h-4 w-4" />
                                  {deletingId === variant.id ? "Deleting..." : "Delete"}
                                </DropdownMenuItem>
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </div>

                          {/* Pricing */}
                          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                            <div className="space-y-1">
                              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                <TrendingUp className="h-3 w-3" />
                                Selling Price
                              </div>
                              <div className="text-xl font-bold text-green-600">
                                {formatPrice(variant.price)}
                              </div>
                            </div>

                            <div className="space-y-1">
                              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                <TrendingDown className="h-3 w-3" />
                                Compare Price
                              </div>
                              <div className="text-lg font-medium">
                                {variant.compare_price > 0 ? (
                                  <div className="flex items-center gap-2">
                                    <span className="line-through text-muted-foreground">
                                      {formatPrice(variant.compare_price)}
                                    </span>
                                    {discountPercentage && (
                                      <Badge variant="secondary" className="text-xs">
                                        {discountPercentage}% off
                                      </Badge>
                                    )}
                                  </div>
                                ) : (
                                  <span className="text-muted-foreground">No compare price</span>
                                )}
                              </div>
                            </div>

                            <div className="space-y-1">
                              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                <Settings className="h-3 w-3" />
                                Cost Price
                              </div>
                              <div className="text-lg font-medium">
                                {variant.cost_price > 0 ? (
                                  <span className="text-blue-600">{formatPrice(variant.cost_price)}</span>
                                ) : (
                                  <span className="text-muted-foreground">Not set</span>
                                )}
                              </div>
                            </div>
                          </div>

                          {/* Properties */}
                          <div className="flex items-center gap-6 text-sm">
                            <div className="flex items-center gap-2">
                              <Weight className="h-4 w-4 text-muted-foreground" />
                              <span className="text-muted-foreground">Weight:</span>
                              <span className="font-medium">{variant.weight} kg</span>
                            </div>
                            <div className="flex items-center gap-2">
                              <ArrowUpDown className="h-4 w-4 text-muted-foreground" />
                              <span className="text-muted-foreground">Position:</span>
                              <span className="font-medium">{variant.position}</span>
                            </div>
                            <div className="flex items-center gap-2">
                              <Sparkles className="h-4 w-4 text-muted-foreground" />
                              <span className="text-muted-foreground">Created:</span>
                              <span className="font-medium">{formatDate(variant.created_at)}</span>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  )
                })}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
