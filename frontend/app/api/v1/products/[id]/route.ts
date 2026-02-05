import { NextRequest, NextResponse } from 'next/server'
import { UpdateProductRequest } from '@/lib/types'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

// GET /api/v1/products/[id] - Get a specific product by ID
export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id: productId } = await params

    if (!productId || isNaN(Number(productId))) {
      return NextResponse.json(
        { message: 'Invalid product ID' },
        { status: 400 }
      )
    }

    // Forward the request to the product service
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/products/${productId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
        'Authorization': request.headers.get('authorization') || '',
      },
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

// PUT /api/v1/products/[id] - Update a specific product by ID
export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id: productId } = await params

    if (!productId || isNaN(Number(productId))) {
      return NextResponse.json(
        { message: 'Invalid product ID' },
        { status: 400 }
      )
    }

    const body: UpdateProductRequest = await request.json()

    // Validate required fields for update
    if (!body || Object.keys(body).length === 0) {
      return NextResponse.json(
        { message: 'Request body is required for update' },
        { status: 400 }
      )
    }

    // Forward the request to the product service
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/products/${productId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
        'Authorization': request.headers.get('authorization') || '',
      },
      body: JSON.stringify(body),
    })

    const data = await response.json()

    // Return the same status and data from the product service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Product update proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}

// DELETE /api/v1/products/[id] - Delete a specific product by ID
export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id: productId } = await params

    if (!productId || isNaN(Number(productId))) {
      return NextResponse.json(
        { message: 'Invalid product ID' },
        { status: 400 }
      )
    }

    // First, get the product data to extract image URLs
    const getResponse = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/products/${productId}`, {
      method: 'GET',
      headers: {
        'Cookie': request.headers.get('cookie') || '',
        'Authorization': request.headers.get('authorization') || '',
      },
    })

    let productData = null
    if (getResponse.ok) {
      productData = await getResponse.json()
    }

    // Forward the delete request to the product service
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/products/${productId}`, {
      method: 'DELETE',
      headers: {
        'Cookie': request.headers.get('cookie') || '',
        'Authorization': request.headers.get('authorization') || '',
      },
    })

    const data = await response.json()

    // If deletion was successful and we have product data, clean up images
    if (response.ok && productData) {
      try {
        const { cleanupEntityImages } = await import('@/lib/cloudinary-cleanup')
        const cleanupResult = await cleanupEntityImages(productData.data || productData)

        if (cleanupResult.deleted > 0) {
          console.log(`Cleaned up ${cleanupResult.deleted} images for deleted product ${productId}`)
        }

        if (cleanupResult.errors.length > 0) {
          console.warn('Some images could not be deleted:', cleanupResult.errors)
        }
      } catch (cleanupError) {
        console.error('Failed to cleanup images for deleted product:', cleanupError)
        // Don't fail the delete operation if image cleanup fails
      }
    }

    // Return the same status and data from the product service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Product delete proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}
