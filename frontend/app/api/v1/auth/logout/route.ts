import { NextRequest, NextResponse } from 'next/server'

const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'http://localhost:8080'

export async function POST(request: NextRequest) {
  try {
    // Forward the request to the backend auth service
    // The cookies will be automatically forwarded by the browser
    const response = await fetch(`${AUTH_SERVICE_URL}/api/v1/auth/logout`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        // Forward cookies from the original request
        'Cookie': request.headers.get('cookie') || '',
      },
    })

    const data = await response.json()

    // Create response with the same status
    const nextResponse = NextResponse.json(data, { status: response.status })

    // Forward any Set-Cookie headers from the backend to clear cookies
    const setCookieHeaders = response.headers.get('set-cookie')
    if (setCookieHeaders) {
      nextResponse.headers.set('Set-Cookie', setCookieHeaders)
    }

    return nextResponse
  } catch (error) {
    console.error('Logout proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' }, 
      { status: 500 }
    )
  }
}
