import { NextRequest, NextResponse } from 'next/server'
import { readdir, stat } from 'fs/promises'
import { join } from 'path'

// Get user's uploaded images
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const page = parseInt(searchParams.get('page') || '1')
    const limit = parseInt(searchParams.get('limit') || '20')
    const userId = searchParams.get('userId') // Optional: filter by user

    const UPLOAD_DIR = join(process.cwd(), 'public', 'uploads')
    
    // Read upload directory
    const files = await readdir(UPLOAD_DIR)
    
    // Filter image files and get metadata
    const imageFiles = []
    for (const file of files) {
      if (file.match(/\.(jpg|jpeg|png|webp|gif)$/i)) {
        // Skip thumbnails
        if (file.includes('_thumbnail') || file.includes('_small') || 
            file.includes('_medium') || file.includes('_large')) {
          continue
        }
        
        // Filter by user if specified
        if (userId && !file.startsWith(`user_${userId}_`)) {
          continue
        }
        
        const filePath = join(UPLOAD_DIR, file)
        const stats = await stat(filePath)
        
        imageFiles.push({
          id: file.split('.')[0], // Use filename as ID for now
          url: `/uploads/${file}`,
          filename: file,
          size: stats.size,
          uploadedAt: stats.birthtime.toISOString(),
          modifiedAt: stats.mtime.toISOString()
        })
      }
    }
    
    // Sort by upload date (newest first)
    imageFiles.sort((a, b) => new Date(b.uploadedAt).getTime() - new Date(a.uploadedAt).getTime())
    
    // Pagination
    const startIndex = (page - 1) * limit
    const endIndex = startIndex + limit
    const paginatedFiles = imageFiles.slice(startIndex, endIndex)
    
    // Calculate pagination info
    const total = imageFiles.length
    const totalPages = Math.ceil(total / limit)
    const hasNext = page < totalPages
    const hasPrev = page > 1
    
    return NextResponse.json({
      success: true,
      data: paginatedFiles,
      pagination: {
        page,
        limit,
        total,
        totalPages,
        hasNext,
        hasPrev
      }
    })

  } catch (error) {
    console.error('Get images error:', error)
    return NextResponse.json(
      { 
        success: false, 
        message: 'Failed to fetch images',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    )
  }
}
