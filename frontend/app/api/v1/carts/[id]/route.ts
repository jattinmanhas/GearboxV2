import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

// GET /api/v1/carts/[id] - Get specific cart
export async function GET(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/carts/${params.id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      },
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        { error: 'Failed to fetch cart', details: errorData },
        { status: response.status }
      )
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Cart API error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

// PUT /api/v1/carts/[id] - Update cart
export async function PUT(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const body = await request.json()
    
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/carts/${params.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      },
      body: JSON.stringify(body),
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        { error: 'Failed to update cart', details: errorData },
        { status: response.status }
      )
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Cart update error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

// DELETE /api/v1/carts/[id] - Delete cart
export async function DELETE(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/carts/${params.id}`, {
      method: 'DELETE',
      headers: {
        'Cookie': request.headers.get('cookie') || '',
      },
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        { error: 'Failed to delete cart', details: errorData },
        { status: response.status }
      )
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Cart deletion error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}
