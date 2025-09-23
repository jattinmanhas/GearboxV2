"use client"

import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { EnhancedImageUpload } from '@/components/ui/enhanced-image-upload'
import { DEFAULT_IMAGE_CONFIG, type UploadedImage, type ImageUploadConfig } from '@/lib/image-upload'
import { ImageIcon, Upload, Settings, Code } from 'lucide-react'

export default function ImageUploadDemoPage() {
  const [selectedImages, setSelectedImages] = useState<UploadedImage[]>([])
  const [customConfig, setCustomConfig] = useState<ImageUploadConfig>(DEFAULT_IMAGE_CONFIG)
  const [showConfig, setShowConfig] = useState(false)

  const handleImageSelect = (image: UploadedImage) => {
    setSelectedImages(prev => [...prev, image])
  }

  const handleImageRemove = (imageId: string) => {
    setSelectedImages(prev => prev.filter(img => img.id !== imageId))
  }

  const handleConfigChange = (key: keyof ImageUploadConfig, value: any) => {
    setCustomConfig(prev => ({
      ...prev,
      [key]: value
    }))
  }

  return (
    <div className="container mx-auto py-8 px-4">
      <div className="max-w-6xl mx-auto space-y-8">
        {/* Header */}
        <div className="text-center space-y-4">
          <h1 className="text-4xl font-bold">Image Upload System Demo</h1>
          <p className="text-xl text-muted-foreground">
            Test the enhanced image upload functionality with drag & drop, validation, and optimization
          </p>
        </div>

        <div className="grid gap-8 lg:grid-cols-2">
          {/* Upload Component */}
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Upload className="h-5 w-5" />
                  Image Upload
                </CardTitle>
                <CardDescription>
                  Drag and drop images or click to browse. Images are automatically optimized and resized.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <EnhancedImageUpload
                  onImageSelect={handleImageSelect}
                  onImageRemove={handleImageRemove}
                  selectedImages={selectedImages}
                  multiple={true}
                  maxImages={10}
                  config={customConfig}
                  showPreview={true}
                  showThumbnails={true}
                />
              </CardContent>
            </Card>

            {/* Configuration Panel */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Settings className="h-5 w-5" />
                  Upload Configuration
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setShowConfig(!showConfig)}
                  >
                    {showConfig ? 'Hide' : 'Show'} Config
                  </Button>
                </CardTitle>
              </CardHeader>
              {showConfig && (
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="text-sm font-medium">Max File Size (MB)</label>
                      <input
                        type="number"
                        value={Math.round(customConfig.maxFileSize / (1024 * 1024))}
                        onChange={(e) => handleConfigChange('maxFileSize', parseInt(e.target.value) * 1024 * 1024)}
                        className="w-full p-2 border rounded"
                      />
                    </div>
                    <div>
                      <label className="text-sm font-medium">Quality (0-100)</label>
                      <input
                        type="number"
                        min="1"
                        max="100"
                        value={customConfig.quality}
                        onChange={(e) => handleConfigChange('quality', parseInt(e.target.value))}
                        className="w-full p-2 border rounded"
                      />
                    </div>
                    <div>
                      <label className="text-sm font-medium">Max Width (px)</label>
                      <input
                        type="number"
                        value={customConfig.maxWidth}
                        onChange={(e) => handleConfigChange('maxWidth', parseInt(e.target.value))}
                        className="w-full p-2 border rounded"
                      />
                    </div>
                    <div>
                      <label className="text-sm font-medium">Max Height (px)</label>
                      <input
                        type="number"
                        value={customConfig.maxHeight}
                        onChange={(e) => handleConfigChange('maxHeight', parseInt(e.target.value))}
                        className="w-full p-2 border rounded"
                      />
                    </div>
                  </div>
                </CardContent>
              )}
            </Card>
          </div>

          {/* Results Panel */}
          <div className="space-y-6">
            {/* Selected Images */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <ImageIcon className="h-5 w-5" />
                  Selected Images ({selectedImages.length})
                </CardTitle>
                <CardDescription>
                  Images that have been uploaded and selected for use
                </CardDescription>
              </CardHeader>
              <CardContent>
                {selectedImages.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    No images selected yet
                  </div>
                ) : (
                  <div className="space-y-4">
                    {selectedImages.map((image) => (
                      <div key={image.id} className="flex items-center gap-4 p-4 border rounded-lg">
                        <img
                          src={image.url}
                          alt={image.alt}
                          className="w-16 h-16 object-cover rounded"
                        />
                        <div className="flex-1 min-w-0">
                          <p className="font-medium truncate">{image.alt || 'Untitled'}</p>
                          <div className="flex gap-2 text-sm text-muted-foreground">
                            <Badge variant="outline">{image.width}x{image.height}</Badge>
                            <Badge variant="outline">
                              {Math.round(image.size / 1024)}KB
                            </Badge>
                            <Badge variant="outline">{image.mimeType}</Badge>
                          </div>
                        </div>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleImageRemove(image.id)}
                        >
                          Remove
                        </Button>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>

            {/* JSON Output */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Code className="h-5 w-5" />
                  JSON Output
                </CardTitle>
                <CardDescription>
                  The data structure that would be sent to your backend
                </CardDescription>
              </CardHeader>
              <CardContent>
                <pre className="bg-muted p-4 rounded-lg overflow-auto text-sm">
                  {JSON.stringify(selectedImages, null, 2)}
                </pre>
              </CardContent>
            </Card>

            {/* Features List */}
            <Card>
              <CardHeader>
                <CardTitle>Features Included</CardTitle>
              </CardHeader>
              <CardContent>
                <ul className="space-y-2 text-sm">
                  <li className="flex items-center gap-2">
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                    Drag & drop file upload
                  </li>
                  <li className="flex items-center gap-2">
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                    File validation (type, size, dimensions)
                  </li>
                  <li className="flex items-center gap-2">
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                    Automatic image optimization with Sharp
                  </li>
                  <li className="flex items-center gap-2">
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                    Multiple thumbnail generation
                  </li>
                  <li className="flex items-center gap-2">
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                    Progress tracking during upload
                  </li>
                  <li className="flex items-center gap-2">
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                    Error handling and user feedback
                  </li>
                  <li className="flex items-center gap-2">
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                    Multiple file selection
                  </li>
                  <li className="flex items-center gap-2">
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                    Image preview and management
                  </li>
                </ul>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Usage Instructions */}
        <Card>
          <CardHeader>
            <CardTitle>How to Use in Your Components</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div>
                <h4 className="font-medium mb-2">1. Basic Usage</h4>
                <pre className="bg-muted p-3 rounded text-sm overflow-x-auto">
{`import { EnhancedImageUpload } from '@/components/ui/enhanced-image-upload'

<EnhancedImageUpload
  onImageSelect={(image) => console.log('Selected:', image)}
  multiple={true}
  maxImages={5}
/>`}
                </pre>
              </div>
              
              <div>
                <h4 className="font-medium mb-2">2. With Custom Configuration</h4>
                <pre className="bg-muted p-3 rounded text-sm overflow-x-auto">
{`const customConfig = {
  maxFileSize: 10 * 1024 * 1024, // 10MB
  maxWidth: 4096,
  maxHeight: 4096,
  quality: 90,
  generateThumbnails: true
}

<EnhancedImageUpload
  config={customConfig}
  onImageSelect={handleImageSelect}
/>`}
                </pre>
              </div>
              
              <div>
                <h4 className="font-medium mb-2">3. Integration with Forms</h4>
                <pre className="bg-muted p-3 rounded text-sm overflow-x-auto">
{`const [productImages, setProductImages] = useState([])

const handleSubmit = async (formData) => {
  const productData = {
    ...formData,
    images: productImages.map(img => ({
      url: img.url,
      alt: img.alt,
      is_primary: false
    }))
  }
  
  await fetch('/api/v1/products', {
    method: 'POST',
    body: JSON.stringify(productData)
  })
}`}
                </pre>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
