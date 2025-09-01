import { NextRequest, NextResponse } from 'next/server'

const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'http://localhost:8080'

export async function POST(request: NextRequest) {
  try {
    const body = await request.json()

    // Forward the request to the auth service
    const response = await fetch(`${AUTH_SERVICE_URL}/api/v1/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    })

    const data = await response.json()

    // If login is successful, set the cookies from the auth service response
    if (response.ok) {
      const responseHeaders = new Headers()
      
      // Copy any Set-Cookie headers from the auth service
      const setCookieHeaders = response.headers.getSetCookie()
      setCookieHeaders.forEach(cookie => {
        responseHeaders.append('Set-Cookie', cookie)
      })

      return NextResponse.json(data, { 
        status: response.status,
        headers: responseHeaders
      })
    }

    // Return the same status and data from the auth service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Login proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}
