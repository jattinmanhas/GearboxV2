import { NextRequest, NextResponse } from 'next/server'

const PAYMENT_SERVICE_URL = process.env.PAYMENT_SERVICE_URL || 'http://localhost:8083'

export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const queryString = searchParams.toString()

    const authHeader = request.headers.get('authorization');
    console.log(`[PaymentsProxy] GET | Auth Header: ${authHeader ? 'Yes' : 'No'}`, authHeader ? `(Length: ${authHeader.length})` : '');

    const response = await fetch(`${PAYMENT_SERVICE_URL}/api/v1/admin/payments${queryString ? `?${queryString}` : ''}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
        'Authorization': authHeader || '',
      },
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        { error: 'Failed to fetch payments', details: errorData },
        { status: response.status }
      )
    }

    const data = await response.json()
    const result = data?.data || data
    return NextResponse.json(result)
  } catch (error) {
    console.error('List payments API error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

export async function POST(request: NextRequest) {
  try {
    const body = await request.json()

    const authHeader = request.headers.get('authorization');
    console.log(`[PaymentsProxy] POST | Auth Header: ${authHeader ? 'Yes' : 'No'}`, authHeader ? `(Length: ${authHeader.length})` : '');

    const response = await fetch(`${PAYMENT_SERVICE_URL}/api/v1/protected/payments`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
        'Authorization': authHeader || '',
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

