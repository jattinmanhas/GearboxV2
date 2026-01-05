"use client";

import React, { useState, useRef, useCallback, useEffect } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
// import { Progress } from '@/components/ui/progress' // Progress component not available
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import {
  Upload,
  X,
  Image as ImageIcon,
  Check,
  AlertCircle,
  Loader2,
  Trash2,
  Eye,
  Download,
} from "lucide-react";
import { cn } from "@/lib/utils";
import {
  uploadImage,
  deleteImage,
  validateImageFile,
  DEFAULT_IMAGE_CONFIG,
  type UploadedImage,
  type ImageUploadConfig,
} from "@/lib/image-upload";
import { useUserStore } from "@/lib/stores/user-store";

interface EnhancedImageUploadProps {
  onImageSelect?: (image: UploadedImage) => void;
  onImagesSelect?: (images: UploadedImage[]) => void;
  onImageRemove?: (imageId: string) => void;
  selectedImages?: UploadedImage[];
  multiple?: boolean;
  maxImages?: number;
  config?: ImageUploadConfig;
  className?: string;
  showPreview?: boolean;
  showThumbnails?: boolean;
  disabled?: boolean;
  onUploadStart?: () => void;
  onUploadEnd?: () => void;
}

export function EnhancedImageUpload({
  onImageSelect,
  onImagesSelect,
  onImageRemove,
  selectedImages = [],
  multiple = false,
  maxImages = 5,
  config = DEFAULT_IMAGE_CONFIG,
  className,
  showPreview = true,
  showThumbnails = true,
  disabled = false,
  onUploadStart,
  onUploadEnd,
}: EnhancedImageUploadProps) {
  const [dragActive, setDragActive] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [previewImages, setPreviewImages] = useState<UploadedImage[]>([]);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { user } = useUserStore();

  // Update preview when selected images change
  useEffect(() => {
    setPreviewImages(selectedImages);
  }, [selectedImages]);

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  }, []);

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      setDragActive(false);

      if (disabled) return;

      const files = Array.from(e.dataTransfer.files);
      if (files.length > 0) {
        handleFiles(files);
      }
    },
    [disabled]
  );

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files && files.length > 0) {
      handleFiles(Array.from(files));
    }
  };

  const handleFiles = async (files: File[]) => {
    if (disabled) return;

    setError(null);

    // Check if adding these files would exceed the limit
    if (selectedImages.length + files.length > maxImages) {
      setError(`Maximum ${maxImages} images allowed`);
      return;
    }

    // Process files one by one
    const uploadedImages: UploadedImage[] = [];

    for (const file of files) {
      // If onImagesSelect is provided, we batch them (triggerSelect = false)
      // If NOT provided, we must trigger individual select (triggerSelect = true)
      const shouldTriggerIndividual = !onImagesSelect;
      const result = await handleSingleFile(file, shouldTriggerIndividual);
      if (result) {
        uploadedImages.push(result);
      }
    }

    if (uploadedImages.length > 0 && onImagesSelect) {
      onImagesSelect(uploadedImages);
    }
  };

  const handleSingleFile = async (
    file: File,
    triggerSelect = true
  ): Promise<UploadedImage | null> => {
    // Validate file
    const validation = validateImageFile(file, config);
    if (!validation.valid) {
      setError(validation.error || "Invalid file");
      return null;
    }

    setUploading(true);
    setUploadProgress(0);
    onUploadStart?.();

    try {
      // Simulate progress (in real app, you'd track actual upload progress)
      const progressInterval = setInterval(() => {
        setUploadProgress((prev) => {
          if (prev >= 90) {
            clearInterval(progressInterval);
            return prev;
          }
          return prev + 10;
        });
      }, 100);

      const result = await uploadImage(file, "", user?.id.toString(), config);

      clearInterval(progressInterval);
      setUploadProgress(100);

      if (result.success && result.data) {
        setPreviewImages((prev) => [...prev, result.data!]);
        if (triggerSelect) {
          onImageSelect?.(result.data);
        }
        return result.data;
      } else {
        setError(result.error || "Upload failed");
        return null;
      }
    } catch (error) {
      setError(error instanceof Error ? error.message : "Upload failed");
      return null;
    } finally {
      setUploading(false);
      setUploadProgress(0);
      onUploadEnd?.();
    }
  };

  const handleRemoveImage = async (image: UploadedImage) => {
    try {
      // Check if this is a Cloudinary image
      const isCloudinary = !!image.publicId;
      await deleteImage(image.id, isCloudinary);
      setPreviewImages((prev) => prev.filter((img) => img.id !== image.id));
      onImageRemove?.(image.id);
    } catch (error) {
      setError("Failed to remove image");
    }
  };

  const openFileDialog = () => {
    if (!disabled) {
      fileInputRef.current?.click();
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
  };

  return (
    <div className={cn("space-y-4", className)}>
      {/* Upload Area */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <ImageIcon className="h-5 w-5" />
            Image Upload
          </CardTitle>
          <CardDescription>
            Upload images for your content. Drag and drop or click to select.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {/* Drag and Drop Area */}
          <div
            className={cn(
              "relative border-2 border-dashed rounded-lg p-8 text-center transition-colors",
              dragActive
                ? "border-primary bg-primary/5"
                : "border-muted-foreground/25 hover:border-muted-foreground/50",
              disabled && "opacity-50 cursor-not-allowed"
            )}
            onDragEnter={handleDrag}
            onDragLeave={handleDrag}
            onDragOver={handleDrag}
            onDrop={handleDrop}
          >
            <input
              ref={fileInputRef}
              type="file"
              multiple={multiple}
              accept={config.allowedTypes.join(",")}
              onChange={handleFileInputChange}
              className="hidden"
              disabled={disabled}
            />

            <div className="space-y-4">
              <div className="mx-auto w-12 h-12 bg-muted rounded-full flex items-center justify-center">
                <Upload className="h-6 w-6 text-muted-foreground" />
              </div>

              <div>
                <p className="text-lg font-medium">
                  {dragActive ? "Drop images here" : "Drag & drop images here"}
                </p>
                <p className="text-sm text-muted-foreground">
                  or{" "}
                  <button
                    type="button"
                    onClick={openFileDialog}
                    className="text-primary hover:underline"
                    disabled={disabled}
                  >
                    browse files
                  </button>
                </p>
              </div>

              <div className="flex flex-wrap justify-center gap-2 text-xs text-muted-foreground">
                <Badge variant="outline">
                  Max {formatFileSize(config.maxFileSize)}
                </Badge>
                <Badge variant="outline">
                  {config.allowedTypes
                    .map((type) => type.split("/")[1])
                    .join(", ")}
                </Badge>
                <Badge variant="outline">
                  {config.maxWidth}x{config.maxHeight}px max
                </Badge>
              </div>
            </div>
          </div>

          {/* Upload Progress */}
          {uploading && (
            <div className="mt-4 space-y-2">
              <div className="flex items-center justify-between text-sm">
                <span>Uploading...</span>
                <span>{uploadProgress}%</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                  style={{ width: `${uploadProgress}%` }}
                ></div>
              </div>
            </div>
          )}

          {/* Error Display */}
          {error && (
            <Alert variant="destructive" className="mt-4">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
        </CardContent>
      </Card>

      {/* Image Preview */}
      {showPreview && previewImages.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">
              Uploaded Images ({previewImages.length})
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
              {previewImages.map((image) => (
                <div key={image.id} className="relative group">
                  <div className="aspect-square rounded-lg overflow-hidden bg-muted">
                    <img
                      src={image.secureUrl || image.url}
                      alt={image.alt}
                      className="w-full h-full object-cover"
                    />
                  </div>

                  {/* Image Info */}
                  <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                    <div className="text-white text-center p-2">
                      <p className="text-xs">
                        {image.width}x{image.height}
                      </p>
                      <p className="text-xs">{formatFileSize(image.size)}</p>
                    </div>
                  </div>

                  {/* Actions */}
                  <div className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity flex gap-1">
                    <Button
                      size="sm"
                      variant="destructive"
                      onClick={() => handleRemoveImage(image)}
                      className="h-6 w-6 p-0"
                    >
                      <X className="h-3 w-3" />
                    </Button>
                  </div>

                  {/* Thumbnails */}
                  {showThumbnails && image.thumbnails && (
                    <div className="mt-2 space-y-1">
                      {Object.entries(image.thumbnails).map(([size, url]) => (
                        <div
                          key={size}
                          className="flex items-center gap-2 text-xs text-muted-foreground"
                        >
                          <img
                            src={url}
                            alt={`${size} thumbnail`}
                            className="w-4 h-4 rounded object-cover"
                          />
                          <span className="capitalize">{size}</span>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}