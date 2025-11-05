import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

function getAuthHeaders(request: NextRequest) {
  return {
    'Content-Type': 'application/json',
    'Cookie': request.headers.get('cookie') || '',
    'Authorization': request.headers.get('authorization') || '',
  }
}

// GET /api/v1/dashboard/recent-orders - Get recent orders for dashboard
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const limit = searchParams.get('limit') || '5'
    const sort = searchParams.get('sort') || 'created_at'
    const order = searchParams.get('order') || 'desc'

    // Fetch recent orders from product service
    const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/orders?limit=${limit}&sort=${sort}&order=${order}`, {
      method: 'GET',
      headers: getAuthHeaders(request),
      credentials: 'include',
    })

    if (!response.ok) {
      const errorData = await response.text()
      return NextResponse.json(
        {
          success: false,
          message: 'Failed to fetch recent orders',
          error: errorData,
        },
        { status: response.status }
      )
    }

    const result = await response.json()
    
    // Transform the response to match what the frontend expects
    let orders = []
    if (result.success && result.data && result.data.orders) {
      orders = result.data.orders.map((order: any) => ({
        id: order.id,
        order_number: order.order_number,
        customer_name: order.customer_name || 'Unknown',
        customer_email: order.customer_email || '',
        total_amount: order.total_amount,
        status: order.status,
        created_at: order.created_at,
      }))
    }

    return NextResponse.json({
      success: true,
      message: 'Recent orders retrieved successfully',
      data: { orders },
    })
  } catch (error) {
    console.error('Recent orders API error:', error)
    return NextResponse.json(
      {
        success: false,
        message: 'Failed to fetch recent orders',
        error: error instanceof Error ? error.message : 'Unknown error',
      },
      { status: 500 }
    )
  }
}
