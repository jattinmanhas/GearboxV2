"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Switch } from "@/components/ui/switch"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Category } from "@/lib/types"
import { FlexibleImageInput, ImageItem } from "@/components/ui/flexible-image-input"
import { getPublicIdFromUrl } from "@/lib/image-upload"

interface CategoryFormProps {
  category?: Category | null
  onSave: (data: any) => void
  onCancel: () => void
}

export function CategoryForm({ category, onSave, onCancel }: CategoryFormProps) {
  const [formData, setFormData] = useState({
    name: "",
    description: "",
    slug: "",
    parent_id: undefined as number | undefined,
    is_active: true,
    sort_order: 0,
    image_url: "",
    image_public_id: "",
    meta_title: "",
    meta_description: "",
  })

  const [errors, setErrors] = useState<Record<string, string>>({})
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isUploading, setIsUploading] = useState(false)
  const [categoryImages, setCategoryImages] = useState<ImageItem[]>([])

  useEffect(() => {
    if (category) {
      setFormData({
        name: category.name,
        description: category.description,
        slug: category.slug || generateSlug(category.name),
        parent_id: category.parent_id,
        is_active: category.is_active,
        sort_order: category?.sort_order ?? 0,
        image_url: category?.image_url ?? "",
        image_public_id: category?.image_public_id ?? "",
        meta_title: category?.meta_title ?? "",
        meta_description: category.meta_description,
      })

      // Load existing category image
      if (category.image_url) {
        setCategoryImages([{
          id: 'category-image',
          url: category.image_url,
          alt: category.name,
          source: 'url' as const,
          publicId: category.image_public_id || undefined
        }])
      } else {
        setCategoryImages([])
      }
    } else {
      // For new categories, ensure slug is generated from name
      setFormData(prev => ({
        ...prev,
        slug: prev.name ? generateSlug(prev.name) : ""
      }))
      setCategoryImages([])
    }
  }, [category])

  const generateSlug = (name: string) => {
    return name
      .toLowerCase()
      .trim()
      .replace(/[^a-z0-9\s-]/g, '')
      .replace(/\s+/g, '-')
      .replace(/-+/g, '-')
      .replace(/^-+|-+$/g, '')
  }

  const handleNameChange = (value: string) => {
    setFormData(prev => {
      const newSlug = generateSlug(value)
      // Only auto-generate slug if it's empty or matches the previous auto-generated slug
      const shouldAutoGenerate = !prev.slug || prev.slug === generateSlug(prev.name)
      return {
        ...prev,
        name: value,
        slug: shouldAutoGenerate ? newSlug : prev.slug
      }
    })
  }

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.name.trim()) {
      newErrors.name = "Name is required"
    }

    // Auto-generate slug if it's empty
    const slug = formData.slug.trim() || generateSlug(formData.name)

    if (!slug) {
      newErrors.slug = "Slug is required"
    } else if (!/^[a-z0-9-]+$/.test(slug)) {
      newErrors.slug = "Slug can only contain lowercase letters, numbers, and hyphens"
    }

    if (formData.sort_order < 0) {
      newErrors.sort_order = "Sort order must be 0 or greater"
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleCategoryImagesChange = (newImages: ImageItem[]) => {
    setCategoryImages(newImages)
    // Update formData with the image URL and public_id
    if (newImages.length > 0) {
      const image = newImages[0]
      setFormData(prev => ({
        ...prev,
        image_url: image.url,
        image_public_id: image.publicId || ""
      }))
    } else {
      setFormData(prev => ({
        ...prev,
        image_url: "",
        image_public_id: ""
      }))
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) return

    setIsSubmitting(true)
    try {
      // Handle image cleanup if image has changed
      if (category && category.image_url && category.image_url !== formData.image_url) {
        // Prefer stored public_id, fallback to URL extraction for legacy data
        const publicId = category.image_public_id || getPublicIdFromUrl(category.image_url)

        if (publicId) {
          // Delete from Cloudinary using public_id
          console.log('Deleting old Cloudinary image with public_id:', publicId)
          try {
            const deleteResponse = await fetch(`/api/v1/upload/image/${publicId}`, {
              method: 'DELETE',
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({ isCloudinary: true })
            })

            const deleteResult = await deleteResponse.json()
            if (!deleteResponse.ok) {
              console.error('Failed to delete old image:', deleteResult)
            } else {
              console.log('Successfully deleted old image:', deleteResult)
            }
          } catch (err) {
            console.error("Failed to delete old Cloudinary image", err)
          }
        } else if (category.image_url.includes('/uploads/')) {
          // Fallback: Local image deletion
          const filename = category.image_url.split('/uploads/').pop()
          if (filename) {
            try {
              const deleteResponse = await fetch(`/api/v1/upload/image/${filename}`, {
                method: 'DELETE',
                headers: {
                  'Content-Type': 'application/json',
                },
                body: JSON.stringify({ isCloudinary: false })
              })

              const deleteResult = await deleteResponse.json()
              if (!deleteResponse.ok) {
                console.error('Failed to delete local image:', deleteResult)
              }
            } catch (err) {
              console.error("Failed to delete local image", err)
            }
          }
        }
      }

      await onSave(formData)
    } catch (error) {
      console.error("Error saving category:", error)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Dialog open={true} onOpenChange={onCancel}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {category ? "Edit Category" : "Create New Category"}
          </DialogTitle>
          <DialogDescription>
            {category
              ? "Update the category information below."
              : "Fill in the details to create a new product category."
            }
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* Name */}
            <div className="space-y-2">
              <Label htmlFor="name">Name *</Label>
              <Input
                id="name"
                value={formData.name}
                onChange={(e) => handleNameChange(e.target.value)}
                placeholder="Category name"
                className={errors.name ? "border-destructive" : ""}
              />
              {errors.name && (
                <p className="text-sm text-destructive">{errors.name}</p>
              )}
            </div>

            {/* Slug */}
            <div className="space-y-2">
              <Label htmlFor="slug">Slug *</Label>
              <Input
                id="slug"
                value={formData.slug}
                onChange={(e) => setFormData(prev => ({ ...prev, slug: e.target.value }))}
                placeholder="category-slug"
                className={errors.slug ? "border-destructive" : ""}
              />
              {errors.slug && (
                <p className="text-sm text-destructive">{errors.slug}</p>
              )}
            </div>
          </div>

          {/* Description */}
          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              value={formData.description}
              onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
              placeholder="Category description"
              rows={3}
            />
          </div>

          {/* Sort Order */}
          <div className="space-y-2">
            <Label htmlFor="sort_order">Sort Order</Label>
            <Input
              id="sort_order"
              type="number"
              value={formData.sort_order}
              onChange={(e) => setFormData(prev => ({
                ...prev,
                sort_order: parseInt(e.target.value) || 0
              }))}
              min="0"
              className={errors.sort_order ? "border-destructive" : ""}
            />
            {errors.sort_order && (
              <p className="text-sm text-destructive">{errors.sort_order}</p>
            )}
          </div>

          {/* Category Image - Full Width */}
          <div className="space-y-2">
            <FlexibleImageInput
              images={categoryImages}
              onImagesChange={handleCategoryImagesChange}
              multiple={false}
              maxImages={1}
              label="Category Image"
              description="Upload an image or provide a direct URL for the category banner."
              onUploadStatusChange={setIsUploading}
            />
          </div>

          {/* Meta Title */}
          <div className="space-y-2">
            <Label htmlFor="meta_title">Meta Title</Label>
            <Input
              id="meta_title"
              value={formData.meta_title}
              onChange={(e) => setFormData(prev => ({ ...prev, meta_title: e.target.value }))}
              placeholder="SEO meta title"
            />
          </div>

          {/* Meta Description */}
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

          {/* Active Status */}
          <div className="flex items-center space-x-2">
            <Switch
              id="is_active"
              checked={formData.is_active}
              onCheckedChange={(checked) => setFormData(prev => ({ ...prev, is_active: checked }))}
            />
            <Label htmlFor="is_active">Active</Label>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={onCancel}>
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting || isUploading}>
              {isSubmitting ? "Saving..." : isUploading ? "Uploading..." : category ? "Update Category" : "Create Category"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
