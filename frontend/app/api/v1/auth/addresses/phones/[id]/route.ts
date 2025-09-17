import { NextRequest, NextResponse } from 'next/server'

const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'http://localhost:8080'

export async function GET(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const { id } = params

    // Forward the request to the auth service
    const response = await fetch(`${AUTH_SERVICE_URL}/api/v1/auth/addresses/phones/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      },
    })

    const data = await response.json()

    // Return the same status and data from the auth service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Get phone number proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}

export async function PUT(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const { id } = params
    const body = await request.json()

    // Forward the request to the auth service
    const response = await fetch(`${AUTH_SERVICE_URL}/api/v1/auth/addresses/phones/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      },
      body: JSON.stringify(body),
    })

    const data = await response.json()

    // Return the same status and data from the auth service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Update phone number proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}

export async function DELETE(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const { id } = params

    // Forward the request to the auth service
    const response = await fetch(`${AUTH_SERVICE_URL}/api/v1/auth/addresses/phones/${id}`, {
      method: 'DELETE',
      headers: {
        'Cookie': request.headers.get('cookie') || '',
      },
    })

    const data = await response.json()

    // Return the same status and data from the auth service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Delete phone number proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}
