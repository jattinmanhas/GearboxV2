import { NextRequest, NextResponse } from 'next/server'

const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'http://localhost:8081'

// GET /api/v1/auth/oauth/[provider] - Initiate OAuth flow
export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ provider: string }> }
) {
  try {
    const { provider } = await params

    if (!provider) {
      return NextResponse.json(
        { message: 'Provider is required' },
        { status: 400 }
      )
    }

    // Forward the request to the auth service
    const response = await fetch(`${AUTH_SERVICE_URL}/api/v1/auth/oauth/${provider}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      },
      credentials: 'include',
    })

    const data = await response.json()

    // Copy any Set-Cookie headers from the auth service (for oauth_state cookie)
    const responseHeaders = new Headers()
    const setCookieHeaders = response.headers.getSetCookie()
    setCookieHeaders.forEach(cookie => {
      responseHeaders.append('Set-Cookie', cookie)
    })

    // Return the same status and data from the auth service
    return NextResponse.json(data, { 
      status: response.status,
      headers: responseHeaders
    })
  } catch (error) {
    console.error('OAuth initiate proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}

