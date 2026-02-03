import { NextResponse } from 'next/server'
import { NextRequest } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id: id } = await params;
    const body = await request.json()

    // Forward the request to the product service
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/categories/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      'Authorization': request.headers.get('authorization') || '',
      },
      body: JSON.stringify(body),
    })

    const data = await response.json()

    // Return the same status and data from the product service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Product proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}

export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id: id } = await params;

    // First, get the category data to extract image URLs
    const getResponse = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/categories/${id}`, {
      method: 'GET',
      headers: {},
    })

    let categoryData = null
    if (getResponse.ok) {
      categoryData = await getResponse.json()
    }

    // Forward the delete request to the product service
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/categories/${id}`, {
      method: 'DELETE',
      headers: {},
    })

    const data = await response.json()

    // If deletion was successful and we have category data, clean up images
    if (response.ok && categoryData) {
      try {
        const { cleanupEntityImages } = await import('@/lib/cloudinary-cleanup')
        const cleanupResult = await cleanupEntityImages(categoryData.data || categoryData)
        
        if (cleanupResult.deleted > 0) {
          console.log(`✅ Cleaned up ${cleanupResult.deleted} images for deleted category ${id}`)
        }
        
        if (cleanupResult.errors.length > 0) {
          console.warn('⚠️ Some images could not be deleted:', cleanupResult.errors)
        }
      } catch (cleanupError) {
        console.error('❌ Failed to cleanup images for deleted category:', cleanupError)
        // Don't fail the delete operation if image cleanup fails
      }
    }

    // Return the same status and data from the product service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Product proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}
