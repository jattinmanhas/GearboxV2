import { NextResponse } from 'next/server'
import { NextRequest } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

export async function PUT(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const id = params.id;
    const body = await request.json()

    // Forward the request to the product service
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/categories/${id}`, {
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
    console.error('Product proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}

export async function DELETE(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const id = params.id;

    // Forward the request to the product service
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/categories/${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
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
