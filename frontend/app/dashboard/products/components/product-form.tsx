"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Switch } from "@/components/ui/switch"
import { MultiSelect, MultiSelectOption } from "@/components/ui/multi-select"
import { 
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Product, Category, CreateProductRequest, UpdateProductRequest } from "@/lib/types"
import { FlexibleImageInput, ImageItem } from "@/components/ui/flexible-image-input"

interface ProductFormProps {
  product?: Product | null
  categories: Category[]
  onSave: (data: any) => void
  onCancel: () => void
}

export function ProductForm({ product, categories, onSave, onCancel }: ProductFormProps) {
  const [formData, setFormData] = useState({
    name: "",
    description: "",
    short_description: "",
    sku: "",
    price: 0,
    compare_price: 0,
    cost_price: 0,
    weight: 0,
    dimensions: "",
    is_active: true,
    is_digital: false,
    requires_shipping: true,
    taxable: true,
    track_quantity: true,
    min_quantity: 0,
    max_quantity: 0,
    meta_title: "",
    meta_description: "",
    tags: "",
    category_ids: [] as number[],
  })

  const [errors, setErrors] = useState<Record<string, string>>({})
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [productImages, setProductImages] = useState<ImageItem[]>([])

  useEffect(() => {
    if (product) {
      setFormData({
        name: product.name,
        description: product.description,
        short_description: product.short_description,
        sku: product.sku,
        price: product.price,
        compare_price: product.compare_price,
        cost_price: product.cost_price,
        weight: product.weight,
        dimensions: product.dimensions,
        is_active: product.is_active,
        is_digital: product.is_digital,
        requires_shipping: product.requires_shipping,
        taxable: product.taxable,
        track_quantity: product.track_quantity,
        min_quantity: product.min_quantity,
        max_quantity: product.max_quantity,
        meta_title: product.meta_title,
        meta_description: product.meta_description,
        tags: product.tags,
        category_ids: product.category_ids || [],
      })
      
      // Load existing product images
      if (product.images && product.images.length > 0) {
        const existingImages: ImageItem[] = product.images.map(img => ({
          id: img.id.toString(),
          url: img.url,
          alt: img.alt,
          source: 'upload' as const
        }))
        setProductImages(existingImages)
      } else {
        setProductImages([])
      }
    }
  }, [product])

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.name.trim()) {
      newErrors.name = "Name is required"
    }

    if (!formData.sku.trim()) {
      newErrors.sku = "SKU is required"
    }

    if (formData.price <= 0) {
      newErrors.price = "Price must be greater than 0"
    }

    if (formData.weight < 0) {
      newErrors.weight = "Weight cannot be negative"
    }

    if (formData.min_quantity < 0) {
      newErrors.min_quantity = "Minimum quantity cannot be negative"
    }

    if (formData.max_quantity < 0) {
      newErrors.max_quantity = "Maximum quantity cannot be negative"
    }

    if (formData.min_quantity > 0 && formData.max_quantity > 0 && formData.min_quantity > formData.max_quantity) {
      newErrors.max_quantity = "Maximum quantity must be greater than minimum quantity"
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) return

    setIsSubmitting(true)
    try {
      // Add image URLs to form data
      const formDataWithImages = {
        ...formData,
        images: productImages.map((img, index) => ({
          url: img.url,
          alt: img.alt || '',
          is_primary: index === 0, // First image is primary
          position: index + 1
        }))
      }
      
      console.log('[ProductForm] Submitting with images:', formDataWithImages.images)
      
      await onSave(formDataWithImages)
    } catch (error) {
      console.error("Error saving product:", error)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleCategoryChange = (selectedValues: (string | number)[]) => {
    setFormData(prev => ({
      ...prev,
      category_ids: selectedValues.map(id => Number(id))
    }))
  }

  const categoryOptions: MultiSelectOption[] = categories.map(category => ({
    value: category.id,
    label: category.name
  }))

  return (
    <Dialog open={true} onOpenChange={onCancel}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto overflow-x-hidden">
        <DialogHeader>
          <DialogTitle>
            {product ? "Edit Product" : "Create New Product"}
          </DialogTitle>
          <DialogDescription>
            {product 
              ? "Update the product information below." 
              : "Fill in the details to create a new product."
            }
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Basic Information */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Basic Information</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="name">Product Name *</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                  placeholder="Product name"
                  className={errors.name ? "border-destructive" : ""}
                />
                {errors.name && (
                  <p className="text-sm text-destructive">{errors.name}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="sku">SKU *</Label>
                <Input
                  id="sku"
                  value={formData.sku}
                  onChange={(e) => setFormData(prev => ({ ...prev, sku: e.target.value }))}
                  placeholder="Product SKU"
                  className={errors.sku ? "border-destructive" : ""}
                />
                {errors.sku && (
                  <p className="text-sm text-destructive">{errors.sku}</p>
                )}
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="short_description">Short Description</Label>
              <Input
                id="short_description"
                value={formData.short_description}
                onChange={(e) => setFormData(prev => ({ ...prev, short_description: e.target.value }))}
                placeholder="Brief product description"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                value={formData.description}
                onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
                placeholder="Detailed product description"
                rows={4}
              />
            </div>
          </div>

          {/* Product Images */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Product Images</h3>
            <FlexibleImageInput
              images={productImages}
              onImagesChange={setProductImages}
              multiple={true}
              maxImages={10}
              label="Product Images"
              description="Upload images or provide direct URLs. First image will be the primary image."
            />
          </div>

          {/* Pricing */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Pricing</h3>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label htmlFor="price">Price *</Label>
                <Input
                  id="price"
                  type="number"
                  step="0.01"
                  value={formData.price}
                  onChange={(e) => setFormData(prev => ({ 
                    ...prev, 
                    price: parseFloat(e.target.value) || 0 
                  }))}
                  className={errors.price ? "border-destructive" : ""}
                />
                {errors.price && (
                  <p className="text-sm text-destructive">{errors.price}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="compare_price">Compare Price</Label>
                <Input
                  id="compare_price"
                  type="number"
                  step="0.01"
                  value={formData.compare_price}
                  onChange={(e) => setFormData(prev => ({ 
                    ...prev, 
                    compare_price: parseFloat(e.target.value) || 0 
                  }))}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="cost_price">Cost Price</Label>
                <Input
                  id="cost_price"
                  type="number"
                  step="0.01"
                  value={formData.cost_price}
                  onChange={(e) => setFormData(prev => ({ 
                    ...prev, 
                    cost_price: parseFloat(e.target.value) || 0 
                  }))}
                />
              </div>
            </div>
          </div>

          {/* Physical Properties */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Physical Properties</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="weight">Weight (lbs)</Label>
                <Input
                  id="weight"
                  type="number"
                  step="0.01"
                  value={formData.weight}
                  onChange={(e) => setFormData(prev => ({ 
                    ...prev, 
                    weight: parseFloat(e.target.value) || 0 
                  }))}
                  className={errors.weight ? "border-destructive" : ""}
                />
                {errors.weight && (
                  <p className="text-sm text-destructive">{errors.weight}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="dimensions">Dimensions</Label>
                <Input
                  id="dimensions"
                  value={formData.dimensions}
                  onChange={(e) => setFormData(prev => ({ ...prev, dimensions: e.target.value }))}
                  placeholder="L x W x H"
                />
              </div>
            </div>
          </div>

          {/* Inventory */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Inventory</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="min_quantity">Minimum Quantity</Label>
                <Input
                  id="min_quantity"
                  type="number"
                  value={formData.min_quantity}
                  onChange={(e) => setFormData(prev => ({ 
                    ...prev, 
                    min_quantity: parseInt(e.target.value) || 0 
                  }))}
                  className={errors.min_quantity ? "border-destructive" : ""}
                />
                {errors.min_quantity && (
                  <p className="text-sm text-destructive">{errors.min_quantity}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="max_quantity">Maximum Quantity</Label>
                <Input
                  id="max_quantity"
                  type="number"
                  value={formData.max_quantity}
                  onChange={(e) => setFormData(prev => ({ 
                    ...prev, 
                    max_quantity: parseInt(e.target.value) || 0 
                  }))}
                  className={errors.max_quantity ? "border-destructive" : ""}
                />
                {errors.max_quantity && (
                  <p className="text-sm text-destructive">{errors.max_quantity}</p>
                )}
              </div>
            </div>
          </div>

          {/* Categories */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Categories</h3>
            <div className="space-y-2">
              <Label>Select Categories</Label>
              <MultiSelect
                options={categoryOptions}
                selectedValues={formData.category_ids}
                onSelectionChange={handleCategoryChange}
                placeholder="Search and select categories..."
                searchPlaceholder="Search categories..."
                maxDisplayed={2}
              />
            </div>
          </div>

          {/* SEO */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium">SEO</h3>
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="meta_title">Meta Title</Label>
                <Input
                  id="meta_title"
                  value={formData.meta_title}
                  onChange={(e) => setFormData(prev => ({ ...prev, meta_title: e.target.value }))}
                  placeholder="SEO meta title"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="meta_description">Meta Description</Label>
                <Textarea
                  id="meta_description"
                  value={formData.meta_description}
                  onChange={(e) => setFormData(prev => ({ ...prev, meta_description: e.target.value }))}
                  placeholder="SEO meta description"
                  rows={2}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="tags">Tags</Label>
                <Input
                  id="tags"
                  value={formData.tags}
                  onChange={(e) => setFormData(prev => ({ ...prev, tags: e.target.value }))}
                  placeholder="tag1, tag2, tag3"
                />
              </div>
            </div>
          </div>

          {/* Settings */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Settings</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-4">
                <div className="flex items-center space-x-2">
                  <Switch
                    id="is_active"
                    checked={formData.is_active}
                    onCheckedChange={(checked) => setFormData(prev => ({ ...prev, is_active: checked }))}
                  />
                  <Label htmlFor="is_active">Active</Label>
                </div>

                <div className="flex items-center space-x-2">
                  <Switch
                    id="is_digital"
                    checked={formData.is_digital}
                    onCheckedChange={(checked) => setFormData(prev => ({ ...prev, is_digital: checked }))}
                  />
                  <Label htmlFor="is_digital">Digital Product</Label>
                </div>

                <div className="flex items-center space-x-2">
                  <Switch
                    id="requires_shipping"
                    checked={formData.requires_shipping}
                    onCheckedChange={(checked) => setFormData(prev => ({ ...prev, requires_shipping: checked }))}
                  />
                  <Label htmlFor="requires_shipping">Requires Shipping</Label>
                </div>
              </div>

              <div className="space-y-4">
                <div className="flex items-center space-x-2">
                  <Switch
                    id="taxable"
                    checked={formData.taxable}
                    onCheckedChange={(checked) => setFormData(prev => ({ ...prev, taxable: checked }))}
                  />
                  <Label htmlFor="taxable">Taxable</Label>
                </div>

                <div className="flex items-center space-x-2">
                  <Switch
                    id="track_quantity"
                    checked={formData.track_quantity}
                    onCheckedChange={(checked) => setFormData(prev => ({ ...prev, track_quantity: checked }))}
                  />
                  <Label htmlFor="track_quantity">Track Quantity</Label>
                </div>
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={onCancel}>
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Saving..." : product ? "Update Product" : "Create Product"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
