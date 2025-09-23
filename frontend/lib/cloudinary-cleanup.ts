/**
 * Cloudinary Cleanup Utilities
 * 
 * This module provides utilities for cleaning up Cloudinary images
 * when entities (products, categories, etc.) are deleted.
 */

import { deleteFromCloudinary } from './cloudinary'

/**
 * Extracts Cloudinary public ID from a URL
 */
export function extractCloudinaryPublicId(url: string): string | null {
  if (!url) return null
  
  // Match Cloudinary URLs
  const cloudinaryRegex = /res\.cloudinary\.com\/[^\/]+\/image\/upload\/(?:v\d+\/)?(.+?)(?:\.[^\/]+)?$/
  const match = url.match(cloudinaryRegex)
  
  if (match) {
    return match[1]
  }
  
  // Also handle URLs that might be in the format: gearbox-uploads/filename
  if (url.includes('gearbox-uploads/')) {
    return url.replace(/^.*gearbox-uploads\//, 'gearbox-uploads/')
  }
  
  return null
}

/**
 * Deletes a single image from Cloudinary
 */
export async function deleteCloudinaryImage(imageUrl: string): Promise<{ success: boolean; error?: string }> {
  try {
    const publicId = extractCloudinaryPublicId(imageUrl)
    
    if (!publicId) {
      return { success: true } // Not a Cloudinary image, nothing to delete
    }
    
    await deleteFromCloudinary(publicId)
    return { success: true }
  } catch (error) {
    console.error('Failed to delete Cloudinary image:', error)
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error'
    }
  }
}

/**
 * Deletes multiple images from Cloudinary
 */
export async function deleteCloudinaryImages(imageUrls: string[]): Promise<{ success: boolean; deleted: number; errors: string[] }> {
  const errors: string[] = []
  let deleted = 0
  
  for (const url of imageUrls) {
    const result = await deleteCloudinaryImage(url)
    if (result.success) {
      deleted++
    } else if (result.error) {
      errors.push(`${url}: ${result.error}`)
    }
  }
  
  return {
    success: errors.length === 0,
    deleted,
    errors
  }
}

/**
 * Extracts all image URLs from an entity (product, category, etc.)
 */
export function extractImageUrls(entity: any): string[] {
  const urls: string[] = []
  
  // Common image field names
  const imageFields = [
    'image_url',
    'imageUrl',
    'featured_image',
    'featuredImage',
    'thumbnail_url',
    'thumbnailUrl',
    'cover_image',
    'coverImage',
    'logo_url',
    'logoUrl'
  ]
  
  // Extract from direct fields
  for (const field of imageFields) {
    if (entity[field] && typeof entity[field] === 'string') {
      urls.push(entity[field])
    }
  }
  
  // Extract from images array
  if (entity.images && Array.isArray(entity.images)) {
    for (const image of entity.images) {
      if (typeof image === 'string') {
        urls.push(image)
      } else if (image && typeof image === 'object') {
        // Handle image objects with url field
        if (image.url) urls.push(image.url)
        if (image.image_url) urls.push(image.image_url)
        if (image.secureUrl) urls.push(image.secureUrl)
        if (image.secure_url) urls.push(image.secure_url)
      }
    }
  }
  
  // Extract from variants (for products)
  if (entity.variants && Array.isArray(entity.variants)) {
    for (const variant of entity.variants) {
      const variantUrls = extractImageUrls(variant)
      urls.push(...variantUrls)
    }
  }
  
  // Remove duplicates and filter out empty/null values
  return [...new Set(urls.filter(url => url && typeof url === 'string'))]
}

/**
 * Cleans up images for a deleted entity
 */
export async function cleanupEntityImages(entity: any): Promise<{ success: boolean; deleted: number; errors: string[] }> {
  const imageUrls = extractImageUrls(entity)
  
  if (imageUrls.length === 0) {
    return { success: true, deleted: 0, errors: [] }
  }
  
  console.log(`Cleaning up ${imageUrls.length} images for deleted entity:`, imageUrls)
  
  return await deleteCloudinaryImages(imageUrls)
}
