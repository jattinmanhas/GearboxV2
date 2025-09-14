import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

// GET /api/v1/inventory/[id] - Get inventory by ID
export async function GET(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const inventoryId = params.id

    if (!inventoryId || isNaN(Number(inventoryId))) {
      return NextResponse.json(
        { message: 'Invalid inventory ID' },
        { status: 400 }
      )
    }

    // Forward the request to the product service with cookies
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/inventory/${inventoryId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      },
    })

    const data = await response.json()

    // Return the same status and data from the product service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Inventory proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}

// PUT /api/v1/inventory/[id] - Update inventory by ID
export async function PUT(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const inventoryId = params.id

    if (!inventoryId || isNaN(Number(inventoryId))) {
      return NextResponse.json(
        { message: 'Invalid inventory ID' },
        { status: 400 }
      )
    }

    const body = await request.json()

    // Validate required fields for update
    if (!body || Object.keys(body).length === 0) {
      return NextResponse.json(
        { message: 'Request body is required for update' },
        { status: 400 }
      )
    }

    // Forward the request to the product service with cookies
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/inventory/${inventoryId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      },
      body: JSON.stringify(body),
    })

    const data = await response.json()

    // Return the same status and data from the product service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Inventory update proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}

// DELETE /api/v1/inventory/[id] - Delete inventory by ID
export async function DELETE(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const inventoryId = params.id

    if (!inventoryId || isNaN(Number(inventoryId))) {
      return NextResponse.json(
        { message: 'Invalid inventory ID' },
        { status: 400 }
      )
    }

    // Forward the request to the product service with cookies
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/inventory/${inventoryId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': request.headers.get('cookie') || '',
      },
    })

    const data = await response.json()

    // Return the same status and data from the product service
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Inventory delete proxy error:', error)
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    )
  }
}
