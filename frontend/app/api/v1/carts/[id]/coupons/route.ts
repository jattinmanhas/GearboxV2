import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

// POST /api/v1/carts/[id]/coupons - Apply coupon to cart
export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id: cartId } = await params
    const body = await request.json()
    
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/carts/${cartId}/coupons`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      },
      body: JSON.stringify(body),
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        { error: 'Failed to apply coupon', details: errorData },
        { status: response.status }
      )
    }

    const data = await response.json()
    
    // Forward the Set-Cookie header from backend to frontend
    const setCookieHeader = response.headers.get('Set-Cookie')
    const nextResponse = NextResponse.json(data)
    
    if (setCookieHeader) {
      nextResponse.headers.set('Set-Cookie', setCookieHeader)
    }
    
    return nextResponse
  } catch (error) {
    console.error('Apply coupon API error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

// DELETE /api/v1/carts/[id]/coupons - Remove coupon from cart
export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id: cartId } = await params
    const body = await request.json()
    
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/carts/${cartId}/coupons`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      },
      body: JSON.stringify(body),
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        { error: 'Failed to remove coupon', details: errorData },
        { status: response.status }
      )
    }

    const data = await response.json()
    
    // Forward the Set-Cookie header from backend to frontend
    const setCookieHeader = response.headers.get('Set-Cookie')
    const nextResponse = NextResponse.json(data)
    
    if (setCookieHeader) {
      nextResponse.headers.set('Set-Cookie', setCookieHeader)
    }
    
    return nextResponse
  } catch (error) {
    console.error('Remove coupon API error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

// GET /api/v1/carts/[id]/coupons - Get cart coupons
export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id: cartId } = await params
    
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/carts/${cartId}/coupons`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      },
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        { error: 'Failed to get cart coupons', details: errorData },
        { status: response.status }
      )
    }

    const data = await response.json()
    
    // Forward the Set-Cookie header from backend to frontend
    const setCookieHeader = response.headers.get('Set-Cookie')
    const nextResponse = NextResponse.json(data)
    
    if (setCookieHeader) {
      nextResponse.headers.set('Set-Cookie', setCookieHeader)
    }
    
    return nextResponse
  } catch (error) {
    console.error('Get cart coupons API error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}
