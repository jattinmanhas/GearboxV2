"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { 
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { 
  CreateInventoryRequest, 
  UpdateInventoryRequest, 
  Inventory,
  Product 
} from "@/lib/types"

interface InventoryFormProps {
  inventory?: Inventory | null
  products: Product[]
  onSave: (data: any) => void
  onCancel: () => void
}

export function InventoryForm({ inventory, products, onSave, onCancel }: InventoryFormProps) {
  const [formData, setFormData] = useState({
    product_id: 0,
    product_variant_id: 0,
    quantity: 0,
    min_stock_level: 0,
    max_stock_level: 0,
    reorder_point: 0,
  })

  const [errors, setErrors] = useState<Record<string, string>>({})
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [variants, setVariants] = useState<any[]>([])
  const [loadingVariants, setLoadingVariants] = useState(false)

  useEffect(() => {
    if (inventory) {
      setFormData({
        product_id: inventory.product_id,
        product_variant_id: inventory.product_variant_id || 0,
        quantity: inventory.quantity,
        min_stock_level: inventory.min_stock_level,
        max_stock_level: inventory.max_stock_level || 0,
        reorder_point: inventory.reorder_point,
      })
      
      // Find the selected product
      const product = products.find(p => p.id === inventory.product_id)
      if (product) {
        setSelectedProduct(product)
      }
    }
  }, [inventory, products])

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.product_id) {
      newErrors.product_id = "Product is required"
    }

    if (formData.quantity < 0) {
      newErrors.quantity = "Quantity cannot be negative"
    }

    if (formData.min_stock_level < 0) {
      newErrors.min_stock_level = "Minimum stock level cannot be negative"
    }

    if (formData.max_stock_level && formData.max_stock_level < formData.min_stock_level) {
      newErrors.max_stock_level = "Maximum stock level must be greater than minimum stock level"
    }

    if (formData.reorder_point < 0) {
      newErrors.reorder_point = "Reorder point cannot be negative"
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) return

    setIsSubmitting(true)
    try {
      const submitData = {
        ...formData,
        product_variant_id: formData.product_variant_id || undefined,
        max_stock_level: formData.max_stock_level || undefined,
      }
      await onSave(submitData)
    } catch (error) {
      console.error("Error saving inventory:", error)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleProductChange = async (productId: number) => {
    const product = products.find(p => p.id === productId)
    setSelectedProduct(product || null)
    setFormData(prev => ({ ...prev, product_id: productId, product_variant_id: 0 }))
    
    // Fetch variants for the selected product
    setLoadingVariants(true)
    try {
      const response = await fetch(`/api/v1/products/${productId}/variants`, {
        credentials: 'include'
      })
      const data = await response.json()
      
      console.log('Variants response:', data) // Debug log
      
      if (response.ok) {
        // Backend returns: { success: true, data: [variants_array] }
        const variants = data.data || []
        console.log('Parsed variants:', variants) // Debug log
        setVariants(variants)
      } else {
        console.log('Variants fetch failed:', response.status, data) // Debug log
        setVariants([])
      }
    } catch (error) {
      console.error('Error fetching variants:', error)
      setVariants([])
    } finally {
      setLoadingVariants(false)
    }
  }

  return (
    <Dialog open={true} onOpenChange={onCancel}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>
            {inventory ? "Edit Inventory" : "Create New Inventory"}
          </DialogTitle>
          <DialogDescription>
            {inventory 
              ? "Update the inventory information below." 
              : "Fill in the details to create a new inventory record."
            }
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Product Selection */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Product Selection</h3>
            <div className="space-y-2">
              <Label htmlFor="product_id">Product *</Label>
              <select
                id="product_id"
                value={formData.product_id}
                onChange={(e) => handleProductChange(Number(e.target.value))}
                className="w-full p-2 border border-input rounded-md bg-background"
              >
                <option value={0}>Select a product</option>
                {products.map((product) => (
                  <option key={product.id} value={product.id}>
                    {product.name} ({product.sku})
                  </option>
                ))}
              </select>
              {errors.product_id && (
                <p className="text-sm text-destructive">{errors.product_id}</p>
              )}
            </div>

            {selectedProduct && (
              <div className="space-y-2">
                <Label htmlFor="product_variant_id">Product Variant (Optional)</Label>
                <select
                  id="product_variant_id"
                  value={formData.product_variant_id}
                  onChange={(e) => setFormData(prev => ({ ...prev, product_variant_id: Number(e.target.value) }))}
                  className="w-full p-2 border border-input rounded-md bg-background"
                  disabled={loadingVariants}
                >
                  <option value={0}>
                    {loadingVariants ? "Loading variants..." : "No variant (base product)"}
                  </option>
                  {variants.map((variant) => (
                    <option key={variant.id} value={variant.id}>
                      {variant.name} ({variant.sku})
                    </option>
                  ))}
                </select>
                {loadingVariants && (
                  <p className="text-sm text-muted-foreground">Loading product variants...</p>
                )}
              </div>
            )}
          </div>

          {/* Stock Levels */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Stock Levels</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="quantity">Current Quantity *</Label>
                <Input
                  id="quantity"
                  type="number"
                  value={formData.quantity}
                  onChange={(e) => setFormData(prev => ({ ...prev, quantity: parseInt(e.target.value) || 0 }))}
                  className={errors.quantity ? "border-destructive" : ""}
                />
                {errors.quantity && (
                  <p className="text-sm text-destructive">{errors.quantity}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="min_stock_level">Minimum Stock Level *</Label>
                <Input
                  id="min_stock_level"
                  type="number"
                  value={formData.min_stock_level}
                  onChange={(e) => setFormData(prev => ({ ...prev, min_stock_level: parseInt(e.target.value) || 0 }))}
                  className={errors.min_stock_level ? "border-destructive" : ""}
                />
                {errors.min_stock_level && (
                  <p className="text-sm text-destructive">{errors.min_stock_level}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="max_stock_level">Maximum Stock Level (Optional)</Label>
                <Input
                  id="max_stock_level"
                  type="number"
                  value={formData.max_stock_level}
                  onChange={(e) => setFormData(prev => ({ ...prev, max_stock_level: parseInt(e.target.value) || 0 }))}
                  className={errors.max_stock_level ? "border-destructive" : ""}
                />
                {errors.max_stock_level && (
                  <p className="text-sm text-destructive">{errors.max_stock_level}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="reorder_point">Reorder Point *</Label>
                <Input
                  id="reorder_point"
                  type="number"
                  value={formData.reorder_point}
                  onChange={(e) => setFormData(prev => ({ ...prev, reorder_point: parseInt(e.target.value) || 0 }))}
                  className={errors.reorder_point ? "border-destructive" : ""}
                />
                {errors.reorder_point && (
                  <p className="text-sm text-destructive">{errors.reorder_point}</p>
                )}
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={onCancel}>
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Saving..." : inventory ? "Update Inventory" : "Create Inventory"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
