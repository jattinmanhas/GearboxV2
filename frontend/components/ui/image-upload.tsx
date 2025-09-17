'use client';

import { useState, useRef } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Upload, X, Image as ImageIcon, Link as LinkIcon } from 'lucide-react';

interface ImageUploadProps {
  onImageInsert: (imageUrl: string, altText?: string) => void;
  onClose: () => void;
}

export function ImageUpload({ onImageInsert, onClose }: ImageUploadProps) {
  const [imageUrl, setImageUrl] = useState('');
  const [altText, setAltText] = useState('');
  const [dragActive, setDragActive] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileUpload = (file: File) => {
    if (file && file.type.startsWith('image/')) {
      // For now, we'll use a simple approach with a data URL
      // In a real app, you'd upload to a cloud service like Cloudinary, AWS S3, etc.
      const reader = new FileReader();
      reader.onload = (e) => {
        const result = e.target?.result as string;
        if (result) {
          setImageUrl(result);
        }
      };
      reader.readAsDataURL(file);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setDragActive(false);
    
    const files = Array.from(e.dataTransfer.files);
    if (files.length > 0) {
      handleFileUpload(files[0]);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setDragActive(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setDragActive(false);
  };

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files && files.length > 0) {
      handleFileUpload(files[0]);
    }
  };

  const handleInsert = () => {
    if (imageUrl.trim()) {
      onImageInsert(imageUrl, altText.trim() || undefined);
      setImageUrl('');
      setAltText('');
      onClose();
    }
  };

  return (
    <Card className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <Card className="w-full max-w-md mx-4">
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle className="flex items-center gap-2">
            <ImageIcon className="h-5 w-5" />
            Insert Image
          </CardTitle>
          <Button variant="ghost" size="sm" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* File Upload Area */}
          <div
            className={`border-2 border-dashed rounded-lg p-6 text-center transition-colors ${
              dragActive 
                ? 'border-blue-500 bg-blue-50 dark:bg-blue-950' 
                : 'border-gray-300 dark:border-gray-600'
            }`}
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
          >
            <Upload className="h-8 w-8 mx-auto mb-2 text-gray-400" />
            <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
              Drag and drop an image here, or click to select
            </p>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => fileInputRef.current?.click()}
            >
              Choose File
            </Button>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              onChange={handleFileInputChange}
              className="hidden"
            />
          </div>

          {/* URL Input */}
          <div className="text-center text-sm text-gray-500 dark:text-gray-400">
            or
          </div>

          <div>
            <Label htmlFor="imageUrl">Image URL</Label>
            <Input
              id="imageUrl"
              value={imageUrl}
              onChange={(e) => setImageUrl(e.target.value)}
              placeholder="https://example.com/image.jpg"
            />
          </div>

          <div>
            <Label htmlFor="altText">Alt Text (optional)</Label>
            <Input
              id="altText"
              value={altText}
              onChange={(e) => setAltText(e.target.value)}
              placeholder="Describe the image for accessibility"
            />
          </div>

          {/* Preview */}
          {imageUrl && (
            <div className="space-y-2">
              <Label>Preview</Label>
              <div className="border rounded-lg p-2">
                <img
                  src={imageUrl}
                  alt={altText || 'Preview'}
                  className="max-w-full h-auto max-h-32 object-contain rounded"
                />
              </div>
            </div>
          )}

          {/* Actions */}
          <div className="flex gap-2 pt-4">
            <Button onClick={handleInsert} disabled={!imageUrl.trim()} className="flex-1">
              <LinkIcon className="h-4 w-4 mr-2" />
              Insert Image
            </Button>
            <Button variant="outline" onClick={onClose} className="flex-1">
              Cancel
            </Button>
          </div>

          {/* Help Text */}
          <div className="text-xs text-gray-500 dark:text-gray-400 space-y-1">
            <p>• Supported formats: JPG, PNG, GIF, WebP</p>
            <p>• For production, consider using a cloud storage service</p>
            <p>• Large images will be automatically resized</p>
          </div>
        </CardContent>
      </Card>
    </Card>
  );
}
