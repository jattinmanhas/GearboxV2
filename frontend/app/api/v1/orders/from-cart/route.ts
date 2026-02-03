import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

export async function POST(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const cartId = searchParams.get('cart_id')
    
    if (!cartId) {
      return NextResponse.json(
        { error: 'cart_id is required' },
        { status: 400 }
      )
    }

    const body = await request.json()
    
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/orders/from-cart?cart_id=${cartId}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      'Authorization': request.headers.get('authorization') || '',
      },
      body: JSON.stringify(body),
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        { error: 'Failed to create order', details: errorData },
        { status: response.status }
      )
    }

    const data = await response.json()
    // Extract data from wrapped response if present
    const result = data?.data || data
    return NextResponse.json(result)
  } catch (error) {
    console.error('Create order API error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

