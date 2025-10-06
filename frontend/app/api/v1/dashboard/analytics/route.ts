import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'
const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'http://localhost:8081'

function getAuthHeaders(request: NextRequest) {
  return {
    'Content-Type': 'application/json',
    'Cookie': request.headers.get('cookie') || '',
  }
}

// GET /api/v1/dashboard/analytics - Get comprehensive dashboard analytics
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const period = searchParams.get('period') || '30d' // 7d, 30d, 90d, 1y

    // Fetch analytics from multiple services in parallel
    const [orderAnalytics, productAnalytics, userAnalytics] = await Promise.allSettled([
      // Order analytics
      fetch(`${PRODUCT_SERVICE_URL}/api/v1/orders/analytics`, {
        method: 'GET',
        headers: getAuthHeaders(request),
      }),
      
      // Product analytics (we'll create this endpoint)
      fetch(`${PRODUCT_SERVICE_URL}/api/v1/products/analytics`, {
        method: 'GET',
        headers: getAuthHeaders(request),
      }),
      
      // User analytics (we'll create this endpoint)
      fetch(`${AUTH_SERVICE_URL}/api/v1/users/analytics`, {
        method: 'GET',
        headers: getAuthHeaders(request),
      }),
    ])

    // Process order analytics
    let orderData = null
    if (orderAnalytics.status === 'fulfilled' && orderAnalytics.value.ok) {
      const orderResponse = await orderAnalytics.value.json()
      orderData = orderResponse.data
    }

    // Process product analytics
    let productData = null
    if (productAnalytics.status === 'fulfilled' && productAnalytics.value.ok) {
      const productResponse = await productAnalytics.value.json()
      productData = productResponse.data
    }

    // Process user analytics
    let userData = null
    if (userAnalytics.status === 'fulfilled' && userAnalytics.value.ok) {
      const userResponse = await userAnalytics.value.json()
      userData = userResponse.data
    }

    // Combine all analytics data
    const analyticsData = {
      orders: orderData,
      products: productData,
      users: userData,
      period,
      lastUpdated: new Date().toISOString(),
    }

    return NextResponse.json({
      success: true,
      message: 'Dashboard analytics retrieved successfully',
      data: analyticsData,
    })
  } catch (error) {
    console.error('Dashboard analytics API error:', error)
    return NextResponse.json(
      { 
        success: false, 
        message: 'Failed to fetch dashboard analytics',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    )
  }
}
