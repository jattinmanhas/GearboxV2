import { NextRequest, NextResponse } from 'next/server'

const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'http://localhost:8080'

export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)

    // Forward the request to the auth service
    const response = await fetch(`${AUTH_SERVICE_URL}/api/v1/auth/roles/user?${searchParams}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      'Authorization': request.headers.get('authorization') || '',
      },
    })

    const data = await response.json()

    // Return the same status and data from the auth service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Get user role proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}
