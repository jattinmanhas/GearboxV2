"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Switch } from "@/components/ui/switch"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { 
  X, 
  Save, 
  Package,
  DollarSign,
  Weight,
  Hash,
  ArrowUpDown,
  Tag,
  TrendingUp,
  TrendingDown,
  Info,
  CheckCircle,
  AlertCircle,
  Sparkles,
  Layers,
  Settings
} from "lucide-react"
import { ProductVariant, CreateProductVariantRequest, UpdateProductVariantRequest } from "@/lib/types"
import { formatPrice } from "@/lib/currency"

interface ProductVariantFormProps {
  productId: number
  variant?: ProductVariant | null
  onSave: (variant: CreateProductVariantRequest | UpdateProductVariantRequest) => Promise<void>
  onCancel: () => void
  loading?: boolean
}

export function ProductVariantForm({ 
  productId, 
  variant, 
  onSave, 
  onCancel, 
  loading = false 
}: ProductVariantFormProps) {
  const [formData, setFormData] = useState({
    name: "",
    sku: "",
    price: 0,
    compare_price: 0,
    cost_price: 0,
    weight: 0,
    is_active: true,
    position: 0
  })
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    if (variant) {
      setFormData({
        name: variant.name,
        sku: variant.sku,
        price: variant.price,
        compare_price: variant.compare_price,
        cost_price: variant.cost_price,
        weight: variant.weight,
        is_active: variant.is_active,
        position: variant.position
      })
    }
  }, [variant])

  const handleInputChange = (field: string, value: string | number | boolean) => {
    setFormData(prev => ({
      ...prev,
      [field]: value
    }))
    
    // Clear error when user starts typing
    if (errors[field]) {
      setErrors(prev => ({
        ...prev,
        [field]: ""
      }))
    }
  }

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.name.trim()) {
      newErrors.name = "Variant name is required"
    }

    if (!formData.sku.trim()) {
      newErrors.sku = "SKU is required"
    }

    if (formData.price <= 0) {
      newErrors.price = "Price must be greater than 0"
    }

    if (formData.compare_price < 0) {
      newErrors.compare_price = "Compare price cannot be negative"
    }

    if (formData.cost_price < 0) {
      newErrors.cost_price = "Cost price cannot be negative"
    }

    if (formData.weight < 0) {
      newErrors.weight = "Weight cannot be negative"
    }

    if (formData.position < 0) {
      newErrors.position = "Position cannot be negative"
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) {
      return
    }

    setIsSubmitting(true)
    try {
      await onSave(formData)
    } catch (error) {
      console.error("Failed to save variant:", error)
    } finally {
      setIsSubmitting(false)
    }
  }

  const isEditing = !!variant

  // Calculate profit margin
  const profitMargin = formData.cost_price > 0 && formData.price > 0 
    ? ((formData.price - formData.cost_price) / formData.price * 100).toFixed(1)
    : null

  // Calculate discount percentage
  const discountPercentage = formData.compare_price > 0 && formData.price > 0 && formData.compare_price > formData.price
    ? (((formData.compare_price - formData.price) / formData.compare_price) * 100).toFixed(1)
    : null

  return (
    <div className="w-full max-w-4xl mx-auto">
      {/* Header */}
      <div className="mb-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-primary/10 rounded-lg">
              <Package className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h2 className="text-2xl font-bold flex items-center gap-2">
                {isEditing ? "Edit Variant" : "Create New Variant"}
                {isEditing && <Badge variant="secondary">Editing</Badge>}
              </h2>
              <p className="text-muted-foreground">
                {isEditing 
                  ? "Update the product variant details below" 
                  : "Add a new variant to expand your product options"
                }
              </p>
            </div>
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={onCancel}
            disabled={isSubmitting}
            className="h-8 w-8 p-0"
          >
            <X className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-8">
        {/* Basic Information Card */}
        <Card className="border-2">
          <CardHeader className="pb-4">
            <CardTitle className="flex items-center gap-2 text-lg">
              <Layers className="h-5 w-5 text-blue-600" />
              Basic Information
            </CardTitle>
            <CardDescription>
              Essential details for identifying this variant
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-3">
                <Label htmlFor="name" className="text-sm font-medium flex items-center gap-2">
                  <Tag className="h-4 w-4" />
                  Variant Name *
                </Label>
                <Input
                  id="name"
                  placeholder="e.g., Small, Red, 32GB, Premium"
                  value={formData.name}
                  onChange={(e) => handleInputChange("name", e.target.value)}
                  className={`h-11 ${errors.name ? "border-red-500 focus:border-red-500" : "focus:border-blue-500"}`}
                  disabled={isSubmitting}
                />
                {errors.name ? (
                  <div className="flex items-center gap-2 text-sm text-red-500">
                    <AlertCircle className="h-4 w-4" />
                    {errors.name}
                  </div>
                ) : (
                  <p className="text-xs text-muted-foreground">
                    A descriptive name for this variant
                  </p>
                )}
              </div>

              <div className="space-y-3">
                <Label htmlFor="sku" className="text-sm font-medium flex items-center gap-2">
                  <Hash className="h-4 w-4" />
                  SKU (Stock Keeping Unit) *
                </Label>
                <Input
                  id="sku"
                  placeholder="e.g., IPHONE-17-32GB-RED"
                  value={formData.sku}
                  onChange={(e) => handleInputChange("sku", e.target.value)}
                  className={`h-11 ${errors.sku ? "border-red-500 focus:border-red-500" : "focus:border-blue-500"}`}
                  disabled={isSubmitting}
                />
                {errors.sku ? (
                  <div className="flex items-center gap-2 text-sm text-red-500">
                    <AlertCircle className="h-4 w-4" />
                    {errors.sku}
                  </div>
                ) : (
                  <p className="text-xs text-muted-foreground">
                    Unique identifier for inventory tracking
                  </p>
                )}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Pricing Card */}
        <Card className="border-2">
          <CardHeader className="pb-4">
            <CardTitle className="flex items-center gap-2 text-lg">
              <DollarSign className="h-5 w-5 text-green-600" />
              Pricing & Profitability
            </CardTitle>
            <CardDescription>
              Set pricing and track profit margins
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="space-y-3">
                <Label htmlFor="price" className="text-sm font-medium flex items-center gap-2">
                  <TrendingUp className="h-4 w-4" />
                  Selling Price *
                </Label>
                <div className="relative">
                  <Input
                    id="price"
                    type="number"
                    step="0.01"
                    min="0"
                    placeholder="0.00"
                    value={formData.price}
                    onChange={(e) => handleInputChange("price", parseFloat(e.target.value) || 0)}
                    className={`h-11 pl-8 ${errors.price ? "border-red-500 focus:border-red-500" : "focus:border-green-500"}`}
                    disabled={isSubmitting}
                  />
                  <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                </div>
                {errors.price ? (
                  <div className="flex items-center gap-2 text-sm text-red-500">
                    <AlertCircle className="h-4 w-4" />
                    {errors.price}
                  </div>
                ) : (
                  <div className="flex items-center gap-2 text-sm">
                    <CheckCircle className="h-4 w-4 text-green-500" />
                    <span className="font-medium text-green-700">{formatPrice(formData.price)}</span>
                  </div>
                )}
              </div>

              <div className="space-y-3">
                <Label htmlFor="compare_price" className="text-sm font-medium flex items-center gap-2">
                  <TrendingDown className="h-4 w-4" />
                  Compare Price
                </Label>
                <div className="relative">
                  <Input
                    id="compare_price"
                    type="number"
                    step="0.01"
                    min="0"
                    placeholder="0.00"
                    value={formData.compare_price}
                    onChange={(e) => handleInputChange("compare_price", parseFloat(e.target.value) || 0)}
                    className={`h-11 pl-8 ${errors.compare_price ? "border-red-500 focus:border-red-500" : "focus:border-green-500"}`}
                    disabled={isSubmitting}
                  />
                  <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                </div>
                {errors.compare_price ? (
                  <div className="flex items-center gap-2 text-sm text-red-500">
                    <AlertCircle className="h-4 w-4" />
                    {errors.compare_price}
                  </div>
                ) : formData.compare_price > 0 ? (
                  <div className="flex items-center gap-2 text-sm">
                    <CheckCircle className="h-4 w-4 text-green-500" />
                    <span className="font-medium text-green-700">{formatPrice(formData.compare_price)}</span>
                    {discountPercentage && (
                      <Badge variant="secondary" className="text-xs">
                        {discountPercentage}% off
                      </Badge>
                    )}
                  </div>
                ) : (
                  <p className="text-xs text-muted-foreground">Optional - shows original price</p>
                )}
              </div>

              <div className="space-y-3">
                <Label htmlFor="cost_price" className="text-sm font-medium flex items-center gap-2">
                  <Settings className="h-4 w-4" />
                  Cost Price
                </Label>
                <div className="relative">
                  <Input
                    id="cost_price"
                    type="number"
                    step="0.01"
                    min="0"
                    placeholder="0.00"
                    value={formData.cost_price}
                    onChange={(e) => handleInputChange("cost_price", parseFloat(e.target.value) || 0)}
                    className={`h-11 pl-8 ${errors.cost_price ? "border-red-500 focus:border-red-500" : "focus:border-green-500"}`}
                    disabled={isSubmitting}
                  />
                  <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                </div>
                {errors.cost_price ? (
                  <div className="flex items-center gap-2 text-sm text-red-500">
                    <AlertCircle className="h-4 w-4" />
                    {errors.cost_price}
                  </div>
                ) : formData.cost_price > 0 ? (
                  <div className="flex items-center gap-2 text-sm">
                    <CheckCircle className="h-4 w-4 text-green-500" />
                    <span className="font-medium text-green-700">{formatPrice(formData.cost_price)}</span>
                    {profitMargin && (
                      <Badge variant="outline" className="text-xs">
                        {profitMargin}% margin
                      </Badge>
                    )}
                  </div>
                ) : (
                  <p className="text-xs text-muted-foreground">Optional - for profit tracking</p>
                )}
              </div>
            </div>

            {/* Profit Summary */}
            {(profitMargin || discountPercentage) && (
              <div className="bg-muted/50 rounded-lg p-4">
                <h4 className="text-sm font-medium mb-3 flex items-center gap-2">
                  <Sparkles className="h-4 w-4" />
                  Quick Summary
                </h4>
                <div className="grid grid-cols-2 gap-4 text-sm">
                  {profitMargin && (
                    <div className="flex items-center justify-between">
                      <span className="text-muted-foreground">Profit Margin:</span>
                      <span className="font-medium text-green-600">{profitMargin}%</span>
                    </div>
                  )}
                  {discountPercentage && (
                    <div className="flex items-center justify-between">
                      <span className="text-muted-foreground">Discount:</span>
                      <span className="font-medium text-blue-600">{discountPercentage}%</span>
                    </div>
                  )}
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Physical Properties Card */}
        <Card className="border-2">
          <CardHeader className="pb-4">
            <CardTitle className="flex items-center gap-2 text-lg">
              <Weight className="h-5 w-5 text-purple-600" />
              Physical Properties
            </CardTitle>
            <CardDescription>
              Weight and display order settings
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-3">
                <Label htmlFor="weight" className="text-sm font-medium flex items-center gap-2">
                  <Weight className="h-4 w-4" />
                  Weight (kg)
                </Label>
                <div className="relative">
                  <Input
                    id="weight"
                    type="number"
                    step="0.001"
                    min="0"
                    placeholder="0.000"
                    value={formData.weight}
                    onChange={(e) => handleInputChange("weight", parseFloat(e.target.value) || 0)}
                    className={`h-11 pl-8 ${errors.weight ? "border-red-500 focus:border-red-500" : "focus:border-purple-500"}`}
                    disabled={isSubmitting}
                  />
                  <Weight className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                </div>
                {errors.weight ? (
                  <div className="flex items-center gap-2 text-sm text-red-500">
                    <AlertCircle className="h-4 w-4" />
                    {errors.weight}
                  </div>
                ) : (
                  <p className="text-xs text-muted-foreground">
                    Used for shipping calculations
                  </p>
                )}
              </div>

              <div className="space-y-3">
                <Label htmlFor="position" className="text-sm font-medium flex items-center gap-2">
                  <ArrowUpDown className="h-4 w-4" />
                  Display Position
                </Label>
                <div className="relative">
                  <Input
                    id="position"
                    type="number"
                    min="0"
                    placeholder="0"
                    value={formData.position}
                    onChange={(e) => handleInputChange("position", parseInt(e.target.value) || 0)}
                    className={`h-11 pl-8 ${errors.position ? "border-red-500 focus:border-red-500" : "focus:border-purple-500"}`}
                    disabled={isSubmitting}
                  />
                  <Hash className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                </div>
                {errors.position ? (
                  <div className="flex items-center gap-2 text-sm text-red-500">
                    <AlertCircle className="h-4 w-4" />
                    {errors.position}
                  </div>
                ) : (
                  <p className="text-xs text-muted-foreground">
                    Lower numbers appear first in listings
                  </p>
                )}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Status Card */}
        <Card className="border-2">
          <CardHeader className="pb-4">
            <CardTitle className="flex items-center gap-2 text-lg">
              <Settings className="h-5 w-5 text-orange-600" />
              Status & Visibility
            </CardTitle>
            <CardDescription>
              Control variant availability and visibility
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between p-4 bg-muted/30 rounded-lg">
              <div className="flex items-center gap-3">
                <Switch
                  id="is_active"
                  checked={formData.is_active}
                  onCheckedChange={(checked) => handleInputChange("is_active", checked)}
                  disabled={isSubmitting}
                />
                <div>
                  <Label htmlFor="is_active" className="text-sm font-medium">
                    Active Variant
                  </Label>
                  <p className="text-xs text-muted-foreground">
                    {formData.is_active 
                      ? "This variant is visible to customers" 
                      : "This variant is hidden from customers"
                    }
                  </p>
                </div>
              </div>
              <Badge variant={formData.is_active ? "default" : "secondary"}>
                {formData.is_active ? "Active" : "Inactive"}
              </Badge>
            </div>
          </CardContent>
        </Card>

        {/* Error Alert */}
        {Object.keys(errors).length > 0 && (
          <Alert variant="destructive" className="border-red-200 bg-red-50">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription className="font-medium">
              Please fix the errors above before saving the variant
            </AlertDescription>
          </Alert>
        )}

        {/* Actions */}
        <div className="flex justify-end gap-3 pt-6">
          <Button
            type="button"
            variant="outline"
            onClick={onCancel}
            disabled={isSubmitting}
            className="h-11 px-6"
          >
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={isSubmitting || loading}
            className="h-11 px-6 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700"
          >
            {isSubmitting ? (
              <div className="flex items-center gap-2">
                <div className="h-4 w-4 animate-spin rounded-full border-2 border-background border-t-foreground" />
                {isEditing ? "Updating..." : "Creating..."}
              </div>
            ) : (
              <div className="flex items-center gap-2">
                <Save className="h-4 w-4" />
                {isEditing ? "Update Variant" : "Create Variant"}
              </div>
            )}
          </Button>
        </div>
      </form>
    </div>
  )
}
