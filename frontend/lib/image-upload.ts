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
  useCloudinary?: boolean; // Whether to use Cloudinary for uploads
}

export interface UploadedImage {
  id: string;
  url: string;
  secureUrl?: string; // Cloudinary secure URL
  publicId?: string; // Cloudinary public ID
  thumbnailUrl?: string;
  alt: string;
  width: number;
  height: number;
  size: number;
  mimeType: string;
  uploadedAt: string;
  thumbnails?: { [key: string]: string };
  folder?: string; // Cloudinary folder
  tags?: string[]; // Cloudinary tags
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
  ],
  useCloudinary: true // Default to Cloudinary if available
};

/**
 * Checks if Cloudinary is properly configured
 * Note: This function only works on the server side
 */
export function isCloudinaryConfigured(): boolean {
  // Only check on server side where process.env is available
  if (typeof window !== 'undefined') {
    throw new Error('isCloudinaryConfigured() can only be called on the server side')
  }

  return !!(
    process.env.CLOUDINARY_CLOUD_NAME &&
    process.env.CLOUDINARY_API_KEY &&
    process.env.CLOUDINARY_API_SECRET
  )
}

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

    // Create form data with Cloudinary preference
    const formData = createImageUploadFormData(file, alt, userId, config);

    // Add Cloudinary preference to form data
    formData.append('useCloudinary', (config.useCloudinary ?? true).toString());

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
export async function deleteImage(imageId: string, isCloudinary: boolean = false): Promise<{ success: boolean; error?: string }> {
  try {
    const response = await fetch(`/api/v1/upload/image/${imageId}`, {
      method: 'DELETE',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ isCloudinary })
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

/**
 * Extract public ID from Cloudinary URL
 * Client-side safe version
 */
export function getPublicIdFromUrl(url: string): string | null {
  if (!url) return null

  // Handle standard Cloudinary URLs
  // Example: https://res.cloudinary.com/cloud_name/image/upload/v1234567890/folder/filename.jpg
  const regex = /\/v\d+\/(.+)\.[a-zA-Z]+$/
  const match = url.match(regex)

  if (match && match[1]) {
    return match[1]
  }

  // Handle URLs without version
  // Example: https://res.cloudinary.com/cloud_name/image/upload/folder/filename.jpg
  const noVersionRegex = /\/upload\/(.+)\.[a-zA-Z]+$/
  const noVersionMatch = url.match(noVersionRegex)

  if (noVersionMatch && noVersionMatch[1]) {
    // If the match starts with 'v' followed by numbers and a slash, it might be a version that wasn't caught
    if (/^v\d+\//.test(noVersionMatch[1])) {
      return noVersionMatch[1].replace(/^v\d+\//, '')
    }
    return noVersionMatch[1]
  }

  return null
}
