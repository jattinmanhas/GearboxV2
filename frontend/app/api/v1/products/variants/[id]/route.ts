import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id: variantId } = await params

    // Forward the request to the product service
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/products/variants/${variantId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })

    const data = await response.json()

    // Return the same status and data from the product service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Product variant proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}

export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id: variantId } = await params
    const body = await request.json()

    // Forward the request to the product service
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/products/variants/${variantId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    })

    const data = await response.json()

    // Return the same status and data from the product service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Product variant proxy error:', error)
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
    const { id: variantId } = await params

    // Forward the request to the product service
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/products/variants/${variantId}`, {
      method: 'DELETE',
      headers: {},
    })

    const data = await response.json()

    // Return the same status and data from the product service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Product variant proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}
