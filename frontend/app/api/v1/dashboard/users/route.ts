import { NextRequest, NextResponse } from 'next/server'

const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'http://localhost:8081'

function getAuthHeaders(request: NextRequest) {
  return {
    'Content-Type': 'application/json',
    'Cookie': request.headers.get('cookie') || '',
  }
}

// GET /api/v1/dashboard/users - Get user analytics for dashboard
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const queryString = searchParams.toString()
    
    const response = await fetch(`${AUTH_SERVICE_URL}/api/v1/users/analytics${queryString ? `?${queryString}` : ''}`, {
      method: 'GET',
      headers: getAuthHeaders(request),
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        { 
          success: false, 
          message: 'Failed to fetch user analytics', 
          error: errorData 
        },
        { status: response.status }
      )
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('User analytics API error:', error)
    return NextResponse.json(
      { 
        success: false, 
        message: 'Failed to fetch user analytics',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    )
  }
}
