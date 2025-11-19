"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Badge } from "@/components/ui/badge"
import { EnhancedImageUpload } from "@/components/ui/enhanced-image-upload"
import { UploadedImage } from "@/lib/image-upload"
import { X, Plus, Link as LinkIcon, Upload, Image as ImageIcon } from "lucide-react"

export interface ImageItem {
  id: string
  url: string
  alt: string
  source: 'upload' | 'url'
  publicId?: string
  isCloudinary?: boolean
  width?: number
  height?: number
  mimeType?: string
  size?: number
  file?: File
}

interface FlexibleImageInputProps {
  images: ImageItem[]
  onImagesChange: (images: ImageItem[]) => void
  multiple?: boolean
  maxImages?: number
  label?: string
  description?: string
  onUploadStatusChange?: (isUploading: boolean) => void
}

export function FlexibleImageInput({
  images,
  onImagesChange,
  multiple = true,
  maxImages = 10,
  label = "Images",
  description,
  onUploadStatusChange
}: FlexibleImageInputProps) {
  const [urlInput, setUrlInput] = useState("")
  const [urlAltInput, setUrlAltInput] = useState("")
  const [selectedImages, setSelectedImages] = useState<UploadedImage[]>([])

  const handleImageUpload = (uploadedImage: UploadedImage) => {
    const newImage: ImageItem = {
      id: uploadedImage.id,
      url: uploadedImage.secureUrl || uploadedImage.url || '',
      alt: uploadedImage.alt,
      source: 'upload',
      publicId: uploadedImage.publicId,
      isCloudinary: !!uploadedImage.publicId,
      width: uploadedImage.width,
      height: uploadedImage.height,
      mimeType: uploadedImage.mimeType,
      size: uploadedImage.size
    }

    if (multiple) {
      onImagesChange([...images, newImage])
    } else {
      onImagesChange([newImage])
    }

    setSelectedImages(prev => [...prev, uploadedImage])
  }

  const handleImageRemove = (imageId: string) => {
    setSelectedImages(prev => prev.filter(img => img.id !== imageId))
  }

  const handleAddUrl = () => {
    if (!urlInput.trim()) return

    const newImage: ImageItem = {
      id: `url-${Date.now()}`,
      url: urlInput.trim(),
      alt: urlAltInput.trim() || 'Product image',
      source: 'url'
    }

    if (multiple) {
      onImagesChange([...images, newImage])
    } else {
      onImagesChange([newImage])
    }

    setUrlInput("")
    setUrlAltInput("")
  }

  const handleRemoveImage = (imageId: string) => {
    onImagesChange(images.filter(img => img.id !== imageId))
    handleImageRemove(imageId)
  }

  const canAddMore = multiple ? images.length < maxImages : images.length === 0

  return (
    <div className="space-y-4">
      <div>
        <Label>{label}</Label>
        {description && (
          <p className="text-sm text-muted-foreground mt-1">{description}</p>
        )}
      </div>

      {/* Display existing images */}
      {images.length > 0 && (
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {images.map((image, index) => (
            <div key={image.id} className="relative group border rounded-lg overflow-hidden">
              <div className="aspect-square bg-muted flex items-center justify-center">
                <img
                  src={image.url}
                  alt={image.alt}
                  className="w-full h-full object-cover"
                  onError={(e) => {
                    const target = e.target as HTMLImageElement
                    target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgZmlsbD0iI2VlZSIvPjx0ZXh0IHg9IjUwJSIgeT0iNTAlIiBmb250LWZhbWlseT0ic2Fucy1zZXJpZiIgZm9udC1zaXplPSIxNCIgZmlsbD0iIzk5OSIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZG9taW5hbnQtYmFzZWxpbmU9Im1pZGRsZSI+SW1hZ2Ugbm90IGZvdW5kPC90ZXh0Pjwvc3ZnPg=='
                  }}
                />
              </div>

              {/* Image info */}
              <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/70 to-transparent p-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-1">
                    <Badge variant="secondary" className="text-xs">
                      {image.source === 'upload' ? (
                        <><Upload className="h-3 w-3 mr-1" /> Uploaded</>
                      ) : (
                        <><LinkIcon className="h-3 w-3 mr-1" /> URL</>
                      )}
                    </Badge>
                    {index === 0 && (
                      <Badge variant="default" className="text-xs">Primary</Badge>
                    )}
                  </div>
                </div>
              </div>

              {/* Remove button */}
              <Button
                variant="destructive"
                size="sm"
                className="absolute top-2 right-2 h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                onClick={() => handleRemoveImage(image.id)}
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
          ))}
        </div>
      )}

      {/* Add new images */}
      {canAddMore && (
        <Tabs defaultValue="upload" className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="upload">
              <Upload className="h-4 w-4 mr-2" />
              Upload Image
            </TabsTrigger>
            <TabsTrigger value="url">
              <LinkIcon className="h-4 w-4 mr-2" />
              Image URL
            </TabsTrigger>
          </TabsList>

          <TabsContent value="upload" className="space-y-4">
            <EnhancedImageUpload
              onImageSelect={handleImageUpload}
              onImageRemove={handleImageRemove}
              selectedImages={selectedImages}
              multiple={multiple}
              maxImages={maxImages}
              showPreview={false}
              showThumbnails={false}
              onUploadStart={() => onUploadStatusChange?.(true)}
              onUploadEnd={() => onUploadStatusChange?.(false)}
              config={{
                maxFileSize: 5 * 1024 * 1024,
                allowedTypes: ['image/jpeg', 'image/png', 'image/webp', 'image/gif'],
                maxWidth: 2048,
                maxHeight: 2048,
                quality: 85,
                generateThumbnails: true,
                thumbnailSizes: [
                  { width: 150, height: 150, name: 'thumbnail' },
                  { width: 300, height: 300, name: 'small' },
                  { width: 600, height: 600, name: 'medium' },
                  { width: 1200, height: 1200, name: 'large' }
                ]
              }}
            />
          </TabsContent>

          <TabsContent value="url" className="space-y-4">
            <div className="space-y-4 border rounded-lg p-4">
              <div className="space-y-2">
                <Label htmlFor="image-url">Image URL *</Label>
                <Input
                  id="image-url"
                  type="url"
                  placeholder="https://example.com/image.jpg"
                  value={urlInput}
                  onChange={(e) => setUrlInput(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                      e.preventDefault()
                      handleAddUrl()
                    }
                  }}
                />
                <p className="text-xs text-muted-foreground">
                  Enter a direct URL to an image (supports jpg, png, webp, gif)
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="image-alt">Alt Text (Optional)</Label>
                <Input
                  id="image-alt"
                  placeholder="Describe the image"
                  value={urlAltInput}
                  onChange={(e) => setUrlAltInput(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                      e.preventDefault()
                      handleAddUrl()
                    }
                  }}
                />
              </div>

              {urlInput && (
                <div className="border rounded-lg p-2 bg-muted">
                  <div className="text-xs font-medium mb-2">Preview:</div>
                  <div className="aspect-video bg-background rounded flex items-center justify-center overflow-hidden">
                    <img
                      src={urlInput}
                      alt="Preview"
                      className="max-w-full max-h-full object-contain"
                      onError={(e) => {
                        const target = e.target as HTMLImageElement
                        target.style.display = 'none'
                        target.parentElement!.innerHTML = '<div class="text-destructive text-sm">Invalid image URL</div>'
                      }}
                    />
                  </div>
                </div>
              )}

              <Button
                type="button"
                onClick={handleAddUrl}
                disabled={!urlInput.trim()}
                className="w-full"
              >
                <Plus className="h-4 w-4 mr-2" />
                Add Image from URL
              </Button>
            </div>
          </TabsContent>
        </Tabs>
      )}

      {!canAddMore && (
        <div className="text-sm text-muted-foreground text-center p-4 border border-dashed rounded-lg">
          {multiple
            ? `Maximum ${maxImages} images reached. Remove an image to add more.`
            : "Remove the current image to add a new one."}
        </div>
      )}
    </div>
  )
}

