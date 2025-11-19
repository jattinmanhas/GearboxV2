/**
 * Cloudinary Configuration and Utilities
 * 
 * This module provides configuration and utility functions for Cloudinary
 * image uploads, transformations, and management.
 */

import { v2 as cloudinary } from 'cloudinary'

// Configure Cloudinary
cloudinary.config({
  cloud_name: process.env.CLOUDINARY_CLOUD_NAME,
  api_key: process.env.CLOUDINARY_API_KEY,
  api_secret: process.env.CLOUDINARY_API_SECRET,
  secure: true
})

export interface CloudinaryUploadResult {
  public_id: string
  version: number
  signature: string
  width: number
  height: number
  format: string
  resource_type: string
  created_at: string
  tags: string[]
  bytes: number
  type: string
  etag: string
  placeholder: boolean
  url: string
  secure_url: string
  folder?: string
  original_filename?: string
  eager?: Array<{
    transformation: string
    width: number
    height: number
    url: string
    secure_url: string
  }>
}

export interface CloudinaryUploadOptions {
  folder?: string
  public_id?: string
  overwrite?: boolean
  resource_type?: 'image' | 'video' | 'raw' | 'auto'
  transformation?: any[]
  eager?: any[]
  eager_async?: boolean
  tags?: string[]
  context?: Record<string, string>
  quality?: string | number
  format?: string
  width?: number
  height?: number
  crop?: string
  gravity?: string
}

export interface CloudinaryImageData {
  id: string
  url: string
  secureUrl: string
  publicId: string
  width: number
  height: number
  format: string
  bytes: number
  createdAt: string
  folder?: string
  tags?: string[]
  thumbnails?: { [key: string]: string }
}

/**
 * Upload an image to Cloudinary
 */
export async function uploadToCloudinary(
  file: Buffer | string,
  options: CloudinaryUploadOptions = {}
): Promise<CloudinaryUploadResult> {
  try {
    const uploadOptions: CloudinaryUploadOptions = {
      folder: 'gearbox-uploads',
      resource_type: 'image',
      overwrite: true,
      quality: 'auto',
      ...options
    }

    const result = await cloudinary.uploader.upload(
      typeof file === 'string' ? file : `data:image/jpeg;base64,${file.toString('base64')}`,
      uploadOptions
    )
    return result
  } catch (error) {
    console.error('Cloudinary upload error:', error)
    throw new Error(`Failed to upload to Cloudinary: ${error instanceof Error ? error.message : 'Unknown error'}`)
  }
}

/**
 * Upload multiple images to Cloudinary
 */
export async function uploadMultipleToCloudinary(
  files: Array<{ buffer: Buffer; filename: string; alt?: string }>,
  options: CloudinaryUploadOptions = {}
): Promise<CloudinaryUploadResult[]> {
  try {
    const uploadPromises = files.map((file, index) => {
      const fileOptions: CloudinaryUploadOptions = {
        ...options,
        public_id: options.public_id ? `${options.public_id}_${index}` : undefined,
        context: file.alt ? { alt: file.alt } : undefined
      }
      return uploadToCloudinary(file.buffer, fileOptions)
    })

    return await Promise.all(uploadPromises)
  } catch (error) {
    console.error('Cloudinary multiple upload error:', error)
    throw new Error(`Failed to upload multiple files to Cloudinary: ${error instanceof Error ? error.message : 'Unknown error'}`)
  }
}

/**
 * Delete an image from Cloudinary
 */
export async function deleteFromCloudinary(publicId: string): Promise<{ result: string }> {
  try {
    const result = await cloudinary.uploader.destroy(publicId)
    return result
  } catch (error) {
    console.error('Cloudinary delete error:', error)
    throw new Error(`Failed to delete from Cloudinary: ${error instanceof Error ? error.message : 'Unknown error'}`)
  }
}

/**
 * Generate Cloudinary URL with transformations
 */
export function getCloudinaryUrl(
  publicId: string,
  transformations: any = {}
): string {
  return cloudinary.url(publicId, {
    secure: true,
    ...transformations
  })
}

/**
 * Generate thumbnail URLs for different sizes
 */
export function generateThumbnailUrls(publicId: string): { [key: string]: string } {
  const sizes = {
    thumbnail: { width: 150, height: 150, crop: 'fill' },
    small: { width: 300, height: 300, crop: 'fill' },
    medium: { width: 600, height: 600, crop: 'fill' },
    large: { width: 1200, height: 1200, crop: 'fill' }
  }

  const thumbnails: { [key: string]: string } = {}

  for (const [size, config] of Object.entries(sizes)) {
    thumbnails[size] = getCloudinaryUrl(publicId, config)
  }

  return thumbnails
}

/**
 * Convert Cloudinary result to our UploadedImage format
 */
export function convertCloudinaryToUploadedImage(
  cloudinaryResult: CloudinaryUploadResult,
  alt: string = '',
  originalFilename?: string
): CloudinaryImageData {
  return {
    id: cloudinaryResult.public_id,
    url: cloudinaryResult.url,
    secureUrl: cloudinaryResult.secure_url,
    publicId: cloudinaryResult.public_id,
    width: cloudinaryResult.width,
    height: cloudinaryResult.height,
    format: cloudinaryResult.format,
    bytes: cloudinaryResult.bytes,
    createdAt: cloudinaryResult.created_at,
    folder: cloudinaryResult.folder,
    tags: cloudinaryResult.tags,
    thumbnails: generateThumbnailUrls(cloudinaryResult.public_id)
  }
}

/**
 * Get optimized image URL with automatic format and quality
 */
export function getOptimizedImageUrl(
  publicId: string,
  width?: number,
  height?: number,
  quality: string | number = 'auto'
): string {
  return getCloudinaryUrl(publicId, {
    width,
    height,
    quality,
    format: 'auto',
    crop: width && height ? 'fill' : 'scale'
  })
}

/**
 * Extract public ID from Cloudinary URL
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

export { cloudinary }
