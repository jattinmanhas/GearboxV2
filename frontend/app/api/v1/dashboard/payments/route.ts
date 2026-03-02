import { NextRequest, NextResponse } from 'next/server'

const PAYMENT_SERVICE_URL = process.env.PAYMENT_SERVICE_URL || 'http://localhost:8083'

function getAuthHeaders(request: NextRequest) {
  return {
    'Content-Type': 'application/json',
    'Cookie': request.headers.get('cookie') || '',
    'Authorization': request.headers.get('authorization') || '',
  }
}

// GET /api/v1/dashboard/payments - Get payment summary for dashboard
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const queryString = searchParams.toString()

    const response = await fetch(
      `${PAYMENT_SERVICE_URL}/api/v1/admin/payments/summary${queryString ? `?${queryString}` : ''}`,
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
          message: 'Failed to fetch payment summary',
          error: errorData,
        },
        { status: response.status }
      )
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Payment summary API error:', error)
    return NextResponse.json(
      {
        success: false,
        message: 'Failed to fetch payment summary',
        error: error instanceof Error ? error.message : 'Unknown error',
      },
      { status: 500 }
    )
  }
}
