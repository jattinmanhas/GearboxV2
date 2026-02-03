import { NextRequest, NextResponse } from 'next/server'

const PAYMENT_SERVICE_URL = process.env.PAYMENT_SERVICE_URL || 'http://localhost:8083'

export async function GET(request: NextRequest) {
  try {
    const response = await fetch(`${PAYMENT_SERVICE_URL}/api/v1/public/payment-gateways`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      'Authorization': request.headers.get('authorization') || '',
      },
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        { error: 'Failed to fetch payment gateways', details: errorData },
        { status: response.status }
      )
    }

    const data = await response.json()
    // Extract data from wrapped response if present
    const result = data?.data || data
    return NextResponse.json(result)
  } catch (error) {
    console.error('Payment gateways API error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

