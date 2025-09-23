/**
 * Image Upload Configuration and Utilities
 * 
 * This module provides configuration and utility functions for image uploads
 * including validation, optimization, and cloud storage integration.
 */

export interface ImageUploadConfig {
  maxFileSize: number; // in bytes
  allowedTypes: string[];
  maxWidth: number;
  maxHeight: number;
  quality: number; // 0-100
  generateThumbnails: boolean;
  thumbnailSizes: { width: number; height: number; name: string }[];
}

export interface UploadedImage {
  id: string;
  url: string;
  thumbnailUrl?: string;
  alt: string;
  width: number;
  height: number;
  size: number;
  mimeType: string;
  uploadedAt: string;
  thumbnails?: { [key: string]: string };
}

export interface ImageUploadResponse {
  success: boolean;
  data?: UploadedImage;
  error?: string;
  message?: string;
}

// Default configuration
export const DEFAULT_IMAGE_CONFIG: ImageUploadConfig = {
  maxFileSize: 5 * 1024 * 1024, // 5MB
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
};

/**
 * Validates an image file before upload
 */
export function validateImageFile(file: File, config: ImageUploadConfig = DEFAULT_IMAGE_CONFIG): { valid: boolean; error?: string } {
  // Check file type
  if (!config.allowedTypes.includes(file.type)) {
    return {
      valid: false,
      error: `Invalid file type. Allowed types: ${config.allowedTypes.join(', ')}`
    };
  }

  // Check file size
  if (file.size > config.maxFileSize) {
    const maxSizeMB = Math.round(config.maxFileSize / (1024 * 1024));
    return {
      valid: false,
      error: `File too large. Maximum size: ${maxSizeMB}MB`
    };
  }

  return { valid: true };
}

/**
 * Generates a unique filename for uploaded images
 */
export function generateImageFilename(originalName: string, userId?: string): string {
  const timestamp = Date.now();
  const randomString = Math.random().toString(36).substring(2, 8);
  const extension = originalName.split('.').pop()?.toLowerCase() || 'jpg';
  const prefix = userId ? `user_${userId}` : 'upload';
  
  return `${prefix}_${timestamp}_${randomString}.${extension}`;
}

/**
 * Creates a FormData object for image upload
 */
export function createImageUploadFormData(
  file: File, 
  alt: string = '', 
  userId?: string,
  config: ImageUploadConfig = DEFAULT_IMAGE_CONFIG
): FormData {
  const formData = new FormData();
  
  // Add the file
  formData.append('file', file);
  
  // Add metadata
  formData.append('alt', alt);
  formData.append('mimeType', file.type);
  formData.append('size', file.size.toString());
  
  // Add configuration
  formData.append('maxWidth', config.maxWidth.toString());
  formData.append('maxHeight', config.maxHeight.toString());
  formData.append('quality', config.quality.toString());
  formData.append('generateThumbnails', config.generateThumbnails.toString());
  
  // Add user context if available
  if (userId) {
    formData.append('userId', userId);
  }
  
  return formData;
}

/**
 * Uploads an image to the server
 */
export async function uploadImage(
  file: File,
  alt: string = '',
  userId?: string,
  config: ImageUploadConfig = DEFAULT_IMAGE_CONFIG
): Promise<ImageUploadResponse> {
  try {
    // Validate file first
    const validation = validateImageFile(file, config);
    if (!validation.valid) {
      return {
        success: false,
        error: validation.error
      };
    }

    // Create form data
    const formData = createImageUploadFormData(file, alt, userId, config);

    // Upload to Next.js API route
    const response = await fetch('/api/v1/upload/image', {
      method: 'POST',
      body: formData,
      credentials: 'include' // Include cookies for authentication
    });

    const result = await response.json();

    if (!response.ok) {
      return {
        success: false,
        error: result.message || 'Upload failed',
        message: result.error
      };
    }

    return {
      success: true,
      data: result.data,
      message: result.message
    };
  } catch (error) {
    console.error('Image upload error:', error);
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Upload failed'
    };
  }
}

/**
 * Deletes an uploaded image
 */
export async function deleteImage(imageId: string): Promise<{ success: boolean; error?: string }> {
  try {
    const response = await fetch(`/api/v1/upload/image/${imageId}`, {
      method: 'DELETE',
      credentials: 'include'
    });

    const result = await response.json();

    if (!response.ok) {
      return {
        success: false,
        error: result.message || 'Delete failed'
      };
    }

    return { success: true };
  } catch (error) {
    console.error('Image delete error:', error);
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Delete failed'
    };
  }
}

/**
 * Gets a list of user's uploaded images
 */
export async function getUserImages(
  page: number = 1,
  limit: number = 20
): Promise<{ success: boolean; data?: UploadedImage[]; error?: string; pagination?: any }> {
  try {
    const response = await fetch(`/api/v1/upload/images?page=${page}&limit=${limit}`, {
      method: 'GET',
      credentials: 'include'
    });

    const result = await response.json();

    if (!response.ok) {
      return {
        success: false,
        error: result.message || 'Failed to fetch images'
      };
    }

    return {
      success: true,
      data: result.data,
      pagination: result.pagination
    };
  } catch (error) {
    console.error('Get images error:', error);
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Failed to fetch images'
    };
  }
}
