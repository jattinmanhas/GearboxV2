import { NextRequest, NextResponse } from 'next/server'
import { writeFile, mkdir, stat, unlink } from 'fs/promises'
import { join } from 'path'
import { randomUUID } from 'crypto'
import sharp from 'sharp'
import { 
  uploadToCloudinary, 
  convertCloudinaryToUploadedImage
} from '@/lib/cloudinary'
import { isCloudinaryConfigured } from '@/lib/image-upload'

// Configuration
const UPLOAD_DIR = join(process.cwd(), 'public', 'uploads')
const MAX_FILE_SIZE = 5 * 1024 * 1024 // 5MB
const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/webp', 'image/gif']
const MAX_WIDTH = 2048
const MAX_HEIGHT = 2048
const QUALITY = 85

// Ensure upload directory exists
async function ensureUploadDir() {
  try {
    await stat(UPLOAD_DIR)
  } catch {
    await mkdir(UPLOAD_DIR, { recursive: true })
  }
}

// Generate unique filename
function generateFilename(originalName: string, userId?: string): string {
  const timestamp = Date.now()
  const randomString = Math.random().toString(36).substring(2, 8)
  const extension = originalName.split('.').pop()?.toLowerCase() || 'jpg'
  const prefix = userId ? `user_${userId}` : 'upload'
  
  return `${prefix}_${timestamp}_${randomString}.${extension}`
}

// Process and optimize image
async function processImage(
  buffer: Buffer,
  filename: string,
  maxWidth: number = MAX_WIDTH,
  maxHeight: number = MAX_HEIGHT,
  quality: number = QUALITY
) {
  const image = sharp(buffer)
  const metadata = await image.metadata()
  
  // Resize if needed
  if (metadata.width && metadata.height) {
    if (metadata.width > maxWidth || metadata.height > maxHeight) {
      image.resize(maxWidth, maxHeight, {
        fit: 'inside',
        withoutEnlargement: true
      })
    }
  }
  
  // Optimize based on format
  if (metadata.format === 'jpeg' || metadata.format === 'jpg') {
    image.jpeg({ quality, progressive: true })
  } else if (metadata.format === 'png') {
    image.png({ quality, progressive: true })
  } else if (metadata.format === 'webp') {
    image.webp({ quality })
  }
  
  return image.toBuffer()
}

// Generate thumbnails
async function generateThumbnails(
  buffer: Buffer,
  baseFilename: string,
  sizes: { width: number; height: number; name: string }[]
) {
  const thumbnails: { [key: string]: string } = {}
  
  for (const size of sizes) {
    const thumbnailFilename = `${baseFilename.split('.')[0]}_${size.name}.jpg`
    const thumbnailPath = join(UPLOAD_DIR, 'thumbnails', thumbnailFilename)
    
    // Ensure thumbnails directory exists
    await mkdir(join(UPLOAD_DIR, 'thumbnails'), { recursive: true })
    
    await sharp(buffer)
      .resize(size.width, size.height, {
        fit: 'cover',
        position: 'center'
      })
      .jpeg({ quality: 80 })
      .toFile(thumbnailPath)
    
    thumbnails[size.name] = `/uploads/thumbnails/${thumbnailFilename}`
  }
  
  return thumbnails
}

// Handle local file upload
async function handleLocalUpload(
  buffer: Buffer,
  file: File,
  alt: string,
  userId: string | null,
  maxWidth: number,
  maxHeight: number,
  quality: number,
  shouldGenerateThumbnails: boolean
) {
  // Ensure upload directory exists
  await ensureUploadDir()
  
  // Generate unique filename
  const filename = generateFilename(file.name, userId || undefined)
  const filePath = join(UPLOAD_DIR, filename)

  // Process and optimize image
  const processedBuffer = await processImage(buffer, filename, maxWidth, maxHeight, quality)

  // Save processed image
  await writeFile(filePath, processedBuffer)

  // Get final image metadata
  const finalImage = sharp(processedBuffer)
  const metadata = await finalImage.metadata()

  // Generate thumbnails if requested
  let thumbnails: { [key: string]: string } = {}
  if (shouldGenerateThumbnails) {
    const thumbnailSizes = [
      { width: 150, height: 150, name: 'thumbnail' },
      { width: 300, height: 300, name: 'small' },
      { width: 600, height: 600, name: 'medium' },
      { width: 1200, height: 1200, name: 'large' }
    ]
    thumbnails = await generateThumbnails(processedBuffer, filename, thumbnailSizes)
  }

  // Create response data
  return {
    id: randomUUID(),
    url: `${process.env.NEXTAUTH_URL || 'http://localhost:3000'}/uploads/${filename}`,
    thumbnailUrl: thumbnails.thumbnail ? `${process.env.NEXTAUTH_URL || 'http://localhost:3000'}${thumbnails.thumbnail}` : undefined,
    alt: alt,
    width: metadata.width || 0,
    height: metadata.height || 0,
    size: processedBuffer.length,
    mimeType: file.type,
    uploadedAt: new Date().toISOString(),
    thumbnails: thumbnails
  }
}

export async function POST(request: NextRequest) {
  try {
    // Parse form data
    const formData = await request.formData()
    const file = formData.get('file') as File
    const alt = formData.get('alt') as string || ''
    const mimeType = formData.get('mimeType') as string
    const size = parseInt(formData.get('size') as string)
    const maxWidth = parseInt(formData.get('maxWidth') as string) || MAX_WIDTH
    const maxHeight = parseInt(formData.get('maxHeight') as string) || MAX_HEIGHT
    const quality = parseInt(formData.get('quality') as string) || QUALITY
    const shouldGenerateThumbnails = formData.get('generateThumbnails') === 'true'
    const userId = formData.get('userId') as string
    const useCloudinary = formData.get('useCloudinary') === 'true'

    // Validate file
    if (!file) {
      return NextResponse.json(
        { success: false, message: 'No file provided' },
        { status: 400 }
      )
    }

    // Validate file type
    if (!ALLOWED_TYPES.includes(file.type)) {
      return NextResponse.json(
        { 
          success: false, 
          message: `Invalid file type. Allowed types: ${ALLOWED_TYPES.join(', ')}` 
        },
        { status: 400 }
      )
    }

    // Validate file size
    if (file.size > MAX_FILE_SIZE) {
      const maxSizeMB = Math.round(MAX_FILE_SIZE / (1024 * 1024))
      return NextResponse.json(
        { 
          success: false, 
          message: `File too large. Maximum size: ${maxSizeMB}MB` 
        },
        { status: 400 }
      )
    }

    // Convert file to buffer
    const buffer = Buffer.from(await file.arrayBuffer())

    // Check if we should use Cloudinary
    const shouldUseCloudinary = useCloudinary && isCloudinaryConfigured()

    let imageData: any

    if (shouldUseCloudinary) {
      // Upload to Cloudinary
      try {
        const cloudinaryResult = await uploadToCloudinary(buffer, {
          folder: 'gearbox-uploads',
          public_id: generateFilename(file.name, userId).split('.')[0], // Remove extension for Cloudinary
          quality: 'auto',
          transformation: [
            { width: maxWidth, height: maxHeight, crop: 'limit' }
          ],
          tags: userId ? [`user_${userId}`] : undefined,
          context: alt ? { alt } : undefined
        })

        // Convert Cloudinary result to our format
        imageData = convertCloudinaryToUploadedImage(cloudinaryResult, alt, file.name)
        
        // Add additional metadata
        imageData.id = cloudinaryResult.public_id
        imageData.mimeType = file.type
        imageData.uploadedAt = new Date().toISOString()
        
        // Ensure URL is complete (Cloudinary should already return complete URLs)
        if (imageData.url && !imageData.url.startsWith('http')) {
          imageData.url = `https:${imageData.url}`
        }
        if (imageData.secureUrl && !imageData.secureUrl.startsWith('http')) {
          imageData.secureUrl = `https:${imageData.secureUrl}`
        }
        
      } catch (cloudinaryError) {
        console.error('Cloudinary upload failed, falling back to local:', cloudinaryError)
        // Fall back to local upload if Cloudinary fails
        imageData = await handleLocalUpload(buffer, file, alt, userId, maxWidth, maxHeight, quality, shouldGenerateThumbnails)
      }
    } else {
      // Use local upload
      imageData = await handleLocalUpload(buffer, file, alt, userId, maxWidth, maxHeight, quality, shouldGenerateThumbnails)
    }

    return NextResponse.json({
      success: true,
      message: 'Image uploaded successfully',
      data: imageData
    })

  } catch (error) {
    console.error('Image upload error:', error)
    return NextResponse.json(
      { 
        success: false, 
        message: 'Upload failed',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    )
  }
}

export async function DELETE(request: NextRequest, { params }: { params: { id: string } }) {
  try {
    const { id } = params
    const body = await request.json()
    const isCloudinary = body?.isCloudinary || false

    if (isCloudinary && isCloudinaryConfigured()) {
      // Delete from Cloudinary
      try {
        const { deleteFromCloudinary } = await import('@/lib/cloudinary')
        await deleteFromCloudinary(id)
        
        return NextResponse.json({
          success: true,
          message: 'Image deleted from Cloudinary successfully'
        })
      } catch (error) {
        console.error('Cloudinary delete error:', error)
        return NextResponse.json(
          { 
            success: false, 
            message: 'Failed to delete from Cloudinary',
            error: error instanceof Error ? error.message : 'Unknown error'
          },
          { status: 500 }
        )
      }
    } else {
      // Delete local file
      const filename = id
      const filePath = join(UPLOAD_DIR, filename)
      
      try {
        await unlink(filePath)
        
        // Also delete thumbnails if they exist
        const thumbnailDir = join(UPLOAD_DIR, 'thumbnails')
        const thumbnailFiles = [
          `${filename.split('.')[0]}_thumbnail.jpg`,
          `${filename.split('.')[0]}_small.jpg`,
          `${filename.split('.')[0]}_medium.jpg`,
          `${filename.split('.')[0]}_large.jpg`
        ]
        
        for (const thumbnailFile of thumbnailFiles) {
          try {
            await unlink(join(thumbnailDir, thumbnailFile))
          } catch {
            // Thumbnail might not exist, ignore error
          }
        }
        
        return NextResponse.json({
          success: true,
          message: 'Image deleted successfully'
        })
      } catch (error) {
        console.error('Local delete error:', error)
        return NextResponse.json(
          { 
            success: false, 
            message: 'Failed to delete image',
            error: error instanceof Error ? error.message : 'Unknown error'
          },
          { status: 500 }
        )
      }
    }
  } catch (error) {
    console.error('Delete error:', error)
    return NextResponse.json(
      { 
        success: false, 
        message: 'Delete failed',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    )
  }
}
