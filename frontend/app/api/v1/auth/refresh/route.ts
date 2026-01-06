import { NextRequest, NextResponse } from 'next/server'

const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'http://localhost:8080'

export async function POST(request: NextRequest) {
    try {
        // Get cookies from the request to forward refresh token
        const cookies = request.headers.get('cookie') || ''

        // Forward the request to the auth service
        const response = await fetch(`${AUTH_SERVICE_URL}/api/v1/auth/refresh`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Cookie': cookies, // Forward cookies (refresh token)
            },
        })

        const data = await response.json()

        // If refresh is successful, forward any Set-Cookie headers
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
        console.error('Refresh proxy error:', error)
        return NextResponse.json(
            { message: 'Internal server error' },
            { status: 500 }
        )
    }
}
