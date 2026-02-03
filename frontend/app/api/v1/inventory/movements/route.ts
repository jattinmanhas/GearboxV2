import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

// GET /api/v1/inventory/movements - List stock movements
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)

    // Forward the request to the product service with cookies
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/inventory/movements?${searchParams}`, {
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
    console.error('Stock movements proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}

// POST /api/v1/inventory/movements - Record stock movement
export async function POST(request: NextRequest) {
  try {
    const body = await request.json()

    // Forward the request to the product service with cookies
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/inventory/movements`, {
      method: 'POST',
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
    console.error('Stock movement creation proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}
