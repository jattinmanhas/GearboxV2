import { NextRequest, NextResponse } from 'next/server'

const PAYMENT_SERVICE_URL = process.env.PAYMENT_SERVICE_URL || 'http://localhost:8083'

export async function POST(request: NextRequest) {
  try {
    const body = await request.json()
    
    const response = await fetch(`${PAYMENT_SERVICE_URL}/api/v1/protected/payments`, {
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
        { error: 'Failed to create payment', details: errorData },
        { status: response.status }
      )
    }

    const data = await response.json()
    // Extract data from wrapped response if present
    const result = data?.data || data
    return NextResponse.json(result)
  } catch (error) {
    console.error('Create payment API error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

