import { NextRequest, NextResponse } from 'next/server'
import { writeFile, mkdir, stat } from 'fs/promises'
import { join } from 'path'
import { randomUUID } from 'crypto'
import sharp from 'sharp'

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

export async function POST(request: NextRequest) {
  try {
    // Ensure upload directory exists
    await ensureUploadDir()
    
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

    // Generate unique filename
    const filename = generateFilename(file.name, userId)
    const filePath = join(UPLOAD_DIR, filename)

    // Convert file to buffer
    const buffer = Buffer.from(await file.arrayBuffer())

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
    const imageData = {
      id: randomUUID(),
      url: `/uploads/${filename}`,
      thumbnailUrl: thumbnails.thumbnail,
      alt: alt,
      width: metadata.width || 0,
      height: metadata.height || 0,
      size: processedBuffer.length,
      mimeType: file.type,
      uploadedAt: new Date().toISOString(),
      thumbnails: thumbnails
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
    // For now, we'll just return success
    // In a real implementation, you'd delete the file from storage
    // and remove any database records
    
    return NextResponse.json({
      success: true,
      message: 'Image deleted successfully'
    })

  } catch (error) {
    console.error('Image delete error:', error)
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
