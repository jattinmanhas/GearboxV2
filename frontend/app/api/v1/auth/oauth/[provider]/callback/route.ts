import { NextRequest, NextResponse } from 'next/server'

const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'http://localhost:8081'
const FRONTEND_URL = process.env.NEXT_PUBLIC_FRONTEND_URL || 'http://localhost:3000'

// GET /api/v1/auth/oauth/[provider]/callback - Handle OAuth callback
export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ provider: string }> }
) {
  try {
    const { provider } = await params

    if (!provider) {
      return NextResponse.redirect(`${FRONTEND_URL}/auth/oauth/error?message=invalid_provider`)
    }

    // Get code and state from query params
    const { searchParams } = new URL(request.url)
    const code = searchParams.get('code')
    const state = searchParams.get('state')

    if (!code || !state) {
      return NextResponse.redirect(`${FRONTEND_URL}/auth/oauth/error?message=missing_parameters`)
    }

    // Forward the request to the auth service
    const response = await fetch(
      `${AUTH_SERVICE_URL}/api/v1/auth/oauth/${provider}/callback?code=${code}&state=${state}`,
      {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Cookie': request.headers.get('cookie') || '',
        'Authorization': request.headers.get('authorization') || '',
        },
        credentials: 'include',
        redirect: 'manual', // Don't follow redirects, we'll handle it
      }
    )

    // The backend will redirect, so we need to extract the redirect URL
    if (response.status >= 300 && response.status < 400) {
      const location = response.headers.get('Location')
      if (location) {
        return NextResponse.redirect(location)
      }
    }

    // If it's not a redirect, return the response
    const data = await response.json()

    // Copy any Set-Cookie headers from the auth service (for auth tokens)
    const responseHeaders = new Headers()
    const setCookieHeaders = response.headers.getSetCookie()
    setCookieHeaders.forEach(cookie => {
      responseHeaders.append('Set-Cookie', cookie)
    })

    return NextResponse.json(data, { 
      status: response.status,
      headers: responseHeaders
    })
  } catch (error) {
    console.error('OAuth callback proxy error:', error)
    return NextResponse.redirect(`${FRONTEND_URL}/auth/oauth/error?message=internal_server_error`)
  }
}

