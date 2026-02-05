import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

// GET /api/v1/orders/[orderId] - Get order by ID
export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ orderId: string }> }
) {
  try {
    const { orderId } = await params

    if (!orderId || isNaN(Number(orderId))) {
      return NextResponse.json(
        { message: 'Invalid order ID' },
        { status: 400 }
      )
    }

    // Forward the request to the product service with cookies
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/orders/${orderId}`, {
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
    console.error('Order proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}

// PUT /api/v1/orders/[orderId] - Update order
export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ orderId: string }> }
) {
  try {
    const { orderId } = await params
    const body = await request.json()

    if (!orderId || isNaN(Number(orderId))) {
      return NextResponse.json(
        { message: 'Invalid order ID' },
        { status: 400 }
      )
    }

    // Forward the request to the product service with cookies
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/orders/${orderId}`, {
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
    console.error('Order update proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}