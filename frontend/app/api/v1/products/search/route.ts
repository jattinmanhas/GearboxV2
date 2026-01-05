import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)

    // Forward the request to the product service
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/products/search?q=${searchParams.get('q')}`, {
      method: 'GET',
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