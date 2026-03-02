import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

// GET /api/v1/orders/[orderId]/payments - Get order payment details
export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ orderId: string }> }
) {
  try {
    const { orderId } = await params

    if (!orderId || Number.isNaN(Number(orderId))) {
      return NextResponse.json(
        { message: 'Invalid order ID' },
        { status: 400 }
      )
    }

    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/orders/${orderId}/payments`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
        'Authorization': request.headers.get('authorization') || '',
      },
    })

    const data = await response.json()
    const result = data?.data || data
    return NextResponse.json(result, { status: response.status })
  } catch (error) {
    console.error('Order payment proxy GET error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}

// POST /api/v1/orders/[orderId]/payments - Create order payment
export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ orderId: string }> }
) {
  try {
    const { orderId } = await params
    const body = await request.json()

    if (!orderId || Number.isNaN(Number(orderId))) {
      return NextResponse.json(
        { message: 'Invalid order ID' },
        { status: 400 }
      )
    }

    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/orders/${orderId}/payments`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
        'Authorization': request.headers.get('authorization') || '',
      },
      body: JSON.stringify(body),
    })

    const data = await response.json()
    const result = data?.data || data
    return NextResponse.json(result, { status: response.status })
  } catch (error) {
    console.error('Order payment proxy POST error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}
