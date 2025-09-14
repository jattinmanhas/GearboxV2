import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

// GET /api/v1/carts/session - Get cart by session
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const endpoint = searchParams.get('endpoint') || 'session'
    
    let url = `${PRODUCT_SERVICE_URL}/api/v1/carts/${endpoint}`
    if (searchParams.toString()) {
      const params = new URLSearchParams(searchParams)
      params.delete('endpoint')
      if (params.toString()) {
        url += `?${params.toString()}`
      }
    }

    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      },
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        { error: 'Failed to fetch cart', details: errorData },
        { status: response.status }
      )
    }

    const data = await response.json()
    
    // Forward the Set-Cookie header from backend to frontend
    const setCookieHeader = response.headers.get('Set-Cookie')
    const nextResponse = NextResponse.json(data)
    
    if (setCookieHeader) {
      nextResponse.headers.set('Set-Cookie', setCookieHeader)
    }
    
    return nextResponse
  } catch (error) {
    console.error('Cart API error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

// POST /api/v1/carts - Create cart or get-or-create
export async function POST(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const endpoint = searchParams.get('endpoint') || 'get-or-create'
    
    let url = `${PRODUCT_SERVICE_URL}/api/v1/carts/${endpoint}`
    if (searchParams.toString()) {
      const params = new URLSearchParams(searchParams)
      params.delete('endpoint')
      if (params.toString()) {
        url += `?${params.toString()}`
      }
    }

    const body = await request.json()
    
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      },
      body: JSON.stringify(body),
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        { error: 'Failed to create/get cart', details: errorData },
        { status: response.status }
      )
    }

    const data = await response.json()
    
    // Forward the Set-Cookie header from backend to frontend
    const setCookieHeader = response.headers.get('Set-Cookie')
    const nextResponse = NextResponse.json(data)
    
    if (setCookieHeader) {
      nextResponse.headers.set('Set-Cookie', setCookieHeader)
    }
    
    return nextResponse
  } catch (error) {
    console.error('Cart creation error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}
