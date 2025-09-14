"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { 
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { 
  StockMovementRequest, 
  Product 
} from "@/lib/types"

interface StockMovementFormProps {
  inventoryId?: number
  products: Product[]
  onSave: (data: StockMovementRequest) => void
  onCancel: () => void
}

export function StockMovementForm({ inventoryId, products, onSave, onCancel }: StockMovementFormProps) {
  const [formData, setFormData] = useState({
    product_id: 0,
    product_variant_id: 0,
    movement_type: 'in' as 'in' | 'out' | 'adjustment' | 'transfer',
    quantity: 0,
    reference: '',
    reference_type: '',
    reason: '',
    notes: '',
  })

  const [errors, setErrors] = useState<Record<string, string>>({})
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [variants, setVariants] = useState<any[]>([])
  const [loadingVariants, setLoadingVariants] = useState(false)

  useEffect(() => {
    if (inventoryId) {
      // If we have an inventory ID, we can pre-populate the product
      // This would require fetching the inventory details
    }
  }, [inventoryId])

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.product_id) {
      newErrors.product_id = "Product is required"
    }

    if (formData.quantity <= 0) {
      newErrors.quantity = "Quantity must be greater than 0"
    }

    if (!formData.reason.trim()) {
      newErrors.reason = "Reason is required"
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) return

    setIsSubmitting(true)
    try {
      const submitData: StockMovementRequest = {
        product_id: formData.product_id,
        product_variant_id: formData.product_variant_id || undefined,
        movement_type: formData.movement_type,
        quantity: formData.quantity,
        reference: formData.reference || undefined,
        reference_type: formData.reference_type || undefined,
        reason: formData.reason,
        notes: formData.notes || undefined,
      }
      await onSave(submitData)
    } catch (error) {
      console.error("Error recording stock movement:", error)
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

  const getMovementTypeLabel = (type: string) => {
    switch (type) {
      case 'in': return 'Stock In'
      case 'out': return 'Stock Out'
      case 'adjustment': return 'Adjustment'
      case 'transfer': return 'Transfer'
      default: return type
    }
  }

  return (
    <Dialog open={true} onOpenChange={onCancel}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>Record Stock Movement</DialogTitle>
          <DialogDescription>
            Record a stock movement for inventory tracking
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

          {/* Movement Details */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Movement Details</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="movement_type">Movement Type *</Label>
                <select
                  id="movement_type"
                  value={formData.movement_type}
                  onChange={(e) => setFormData(prev => ({ ...prev, movement_type: e.target.value as any }))}
                  className="w-full p-2 border border-input rounded-md bg-background"
                >
                  <option value="in">Stock In</option>
                  <option value="out">Stock Out</option>
                  <option value="adjustment">Adjustment</option>
                  <option value="transfer">Transfer</option>
                </select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="quantity">Quantity *</Label>
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
                <Label htmlFor="reference">Reference (Optional)</Label>
                <Input
                  id="reference"
                  value={formData.reference}
                  onChange={(e) => setFormData(prev => ({ ...prev, reference: e.target.value }))}
                  placeholder="e.g., PO-12345, RMA-67890"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="reference_type">Reference Type (Optional)</Label>
                <Input
                  id="reference_type"
                  value={formData.reference_type}
                  onChange={(e) => setFormData(prev => ({ ...prev, reference_type: e.target.value }))}
                  placeholder="e.g., Purchase Order, Return"
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="reason">Reason *</Label>
              <Input
                id="reason"
                value={formData.reason}
                onChange={(e) => setFormData(prev => ({ ...prev, reason: e.target.value }))}
                placeholder="e.g., Restock, Sale, Damaged goods"
                className={errors.reason ? "border-destructive" : ""}
              />
              {errors.reason && (
                <p className="text-sm text-destructive">{errors.reason}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="notes">Notes (Optional)</Label>
              <Textarea
                id="notes"
                value={formData.notes}
                onChange={(e) => setFormData(prev => ({ ...prev, notes: e.target.value }))}
                placeholder="Additional notes about this movement"
                rows={3}
              />
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={onCancel}>
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Recording..." : "Record Movement"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
