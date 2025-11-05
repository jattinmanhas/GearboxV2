import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

function getAuthHeaders(request: NextRequest) {
  return {
    'Content-Type': 'application/json',
    'Cookie': request.headers.get('cookie') || '',
    'Authorization': request.headers.get('authorization') || '',
  }
}

// GET /api/v1/orders - List orders with filters and pagination
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const queryString = searchParams.toString()

    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/orders?${queryString}`, {
      method: 'GET',
      headers: getAuthHeaders(request),
      credentials: 'include',
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        {
          success: false,
          message: 'Failed to fetch orders',
          error: errorData,
        },
        { status: response.status }
      )
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Orders API error:', error)
    return NextResponse.json(
      {
        success: false,
        message: 'Failed to fetch orders',
        error: error instanceof Error ? error.message : 'Unknown error',
      },
      { status: 500 }
    )
  }
}

// POST /api/v1/orders - Create a new order
export async function POST(request: NextRequest) {
  try {
    const body = await request.json()

    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/orders`, {
      method: 'POST',
      headers: getAuthHeaders(request),
      credentials: 'include',
      body: JSON.stringify(body),
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        {
          success: false,
          message: 'Failed to create order',
          error: errorData,
        },
        { status: response.status }
      )
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Create order API error:', error)
    return NextResponse.json(
      {
        success: false,
        message: 'Failed to create order',
        error: error instanceof Error ? error.message : 'Unknown error',
      },
      { status: 500 }
    )
  }
}

