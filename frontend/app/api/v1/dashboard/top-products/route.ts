import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

function getAuthHeaders(request: NextRequest) {
  return {
    'Content-Type': 'application/json',
    'Cookie': request.headers.get('cookie') || '',
    'Authorization': request.headers.get('authorization') || '',
  }
}

// GET /api/v1/dashboard/top-products - Get top selling products for dashboard
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const limit = searchParams.get('limit') || '5'

    const response = await fetch(
      `${PRODUCT_SERVICE_URL}/api/v1/orders/analytics/top-products?limit=${limit}`,
      {
        method: 'GET',
        headers: getAuthHeaders(request),
      }
    )

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        {
          success: false,
          message: 'Failed to fetch top products',
          error: errorData,
        },
        { status: response.status }
      )
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Top products API error:', error)
    return NextResponse.json(
      {
        success: false,
        message: 'Failed to fetch top products',
        error: error instanceof Error ? error.message : 'Unknown error',
      },
      { status: 500 }
    )
  }
}
