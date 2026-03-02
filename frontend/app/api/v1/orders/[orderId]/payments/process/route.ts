import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

// POST /api/v1/orders/[orderId]/payments/process - Process payment for an order
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

    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/orders/${orderId}/payments/process`, {
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
    console.error('Order payment process proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}
