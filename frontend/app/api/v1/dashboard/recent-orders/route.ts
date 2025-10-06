import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

function getAuthHeaders(request: NextRequest) {
  return {
    'Content-Type': 'application/json',
    'Cookie': request.headers.get('cookie') || '',
  }
}

// GET /api/v1/dashboard/recent-orders - Get recent orders for dashboard
export async function GET(request: NextRequest) {
  try {
    console.log('Recent orders API called')
    const { searchParams } = new URL(request.url)
    const limit = searchParams.get('limit') || '5'
    const sort = searchParams.get('sort') || 'created_at'
    const order = searchParams.get('order') || 'desc'

    // Try to fetch recent orders from product service
    let orders = []
    try {
      console.log('Attempting to fetch from product service:', `${PRODUCT_SERVICE_URL}/api/v1/orders?limit=${limit}&sort=${sort}&order=${order}`)
      const response = await fetch(`${PRODUCT_SERVICE_URL}/api/v1/orders?limit=${limit}&sort=${sort}&order=${order}`, {
        method: 'GET',
        headers: getAuthHeaders(request),
      })

      console.log('Product service response status:', response.status)
      if (response.ok) {
        const result = await response.json()
        console.log('Product service response:', result)
        if (result.success && result.data && result.data.orders) {
          orders = result.data.orders
          console.log('Using real orders data')
        }
      } else {
        console.warn('Product service returned error:', response.status, response.statusText)
      }
    } catch (error) {
      console.warn('Product service unavailable, using mock data:', error)
    }

    // If product service is unavailable or returns no data, use mock data
    if (orders.length === 0) {
      console.log('Using mock orders data')
      const mockOrders = [
        {
          id: 1,
          order_number: 'ORD-001',
          customer_name: 'John Doe',
          customer_email: 'john@example.com',
          total_amount: 299.99,
          status: 'delivered',
          created_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(), // 2 days ago
        },
        {
          id: 2,
          order_number: 'ORD-002',
          customer_name: 'Jane Smith',
          customer_email: 'jane@example.com',
          total_amount: 149.50,
          status: 'processing',
          created_at: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000).toISOString(), // 1 day ago
        },
        {
          id: 3,
          order_number: 'ORD-003',
          customer_name: 'Bob Johnson',
          customer_email: 'bob@example.com',
          total_amount: 89.99,
          status: 'shipped',
          created_at: new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString(), // 6 hours ago
        },
        {
          id: 4,
          order_number: 'ORD-004',
          customer_name: 'Alice Brown',
          customer_email: 'alice@example.com',
          total_amount: 199.99,
          status: 'pending',
          created_at: new Date(Date.now() - 3 * 60 * 60 * 1000).toISOString(), // 3 hours ago
        },
        {
          id: 5,
          order_number: 'ORD-005',
          customer_name: 'Charlie Wilson',
          customer_email: 'charlie@example.com',
          total_amount: 79.99,
          status: 'cancelled',
          created_at: new Date(Date.now() - 1 * 60 * 60 * 1000).toISOString(), // 1 hour ago
        },
      ]
      orders = mockOrders.slice(0, parseInt(limit))
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
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    )
  }
}
