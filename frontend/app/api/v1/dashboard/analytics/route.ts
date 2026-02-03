import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'
const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'http://localhost:8081'
const PAYMENT_SERVICE_URL = process.env.PAYMENT_SERVICE_URL || 'http://localhost:8083'
const BLOG_SERVICE_URL = process.env.BLOG_SERVICE_URL || 'http://localhost:3003'

function getAuthHeaders(request: NextRequest) {
  return {
    'Content-Type': 'application/json',
    'Cookie': request.headers.get('cookie') || '',
    'Authorization': request.headers.get('authorization') || '',
  }
}

function getDateRange(period: string) {
  const endDate = new Date()
  const startDate = new Date()

  switch (period) {
    case '7d':
      startDate.setDate(endDate.getDate() - 7)
      break
    case '30d':
      startDate.setDate(endDate.getDate() - 30)
      break
    case '90d':
      startDate.setDate(endDate.getDate() - 90)
      break
    case '1y':
      startDate.setFullYear(endDate.getFullYear() - 1)
      break
    default:
      startDate.setDate(endDate.getDate() - 30)
  }

  const toDateString = (date: Date) => date.toISOString().split('T')[0]

  return {
    startDateStr: toDateString(startDate),
    endDateStr: toDateString(endDate),
  }
}

async function fetchBlogList(request: NextRequest, params: Record<string, string | number | undefined>) {
  const searchParams = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null) {
      searchParams.set(key, String(value))
    }
  })

  const response = await fetch(`${BLOG_SERVICE_URL}/api/v1/posts?${searchParams.toString()}`, {
    method: 'GET',
    headers: getAuthHeaders(request),
  })

  if (!response.ok) {
    const errorData = await response.text()
    throw new Error(errorData || 'Failed to fetch blog data')
  }

  return response.json()
}

// GET /api/v1/dashboard/analytics - Get comprehensive dashboard analytics
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const period = searchParams.get('period') || '30d' // 7d, 30d, 90d, 1y
    const { startDateStr, endDateStr } = getDateRange(period)

    let blogSummary = null
    try {
      const perPage = 100
      const maxPages = 10

      const totalResponse = await fetchBlogList(request, { limit: 1, page: 1 })
      const totalPosts = totalResponse?.data?.total || 0

      const publishedResponse = await fetchBlogList(request, { status: 'published', limit: 1, page: 1 })
      const publishedPosts = publishedResponse?.data?.total || 0

      const draftResponse = await fetchBlogList(request, { status: 'draft', limit: 1, page: 1 })
      const draftPosts = draftResponse?.data?.total || 0

      const archivedResponse = await fetchBlogList(request, { status: 'archived', limit: 1, page: 1 })
      const archivedPosts = archivedResponse?.data?.total || 0

      let page = 1
      let totalPages = Math.min(Math.ceil(totalPosts / perPage), maxPages)
      let totalViews = 0
      let totalReadTime = 0
      let countedPosts = 0

      while (page <= totalPages) {
        const pageResponse = await fetchBlogList(request, { limit: perPage, page })
        const posts = pageResponse?.data?.posts || []

        posts.forEach((post: any) => {
          totalViews += Number(post.viewCount) || 0
          totalReadTime += Number(post.readTime) || 0
          countedPosts += 1
        })

        page += 1
      }

      const averageReadTime = countedPosts > 0 ? totalReadTime / countedPosts : 0

      blogSummary = {
        total_posts: totalPosts,
        published_posts: publishedPosts,
        draft_posts: draftPosts,
        archived_posts: archivedPosts,
        total_views: totalViews,
        average_read_time: Number(averageReadTime.toFixed(2)),
      }
    } catch (error) {
      console.warn('Blog summary fetch failed:', error)
    }

    // Fetch analytics from multiple services in parallel
    const [
      orderAnalytics,
      productAnalytics,
      userAnalytics,
      paymentAnalytics,
      topProducts,
    ] = await Promise.allSettled([
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
      fetch(`${AUTH_SERVICE_URL}/api/v1/auth/users/analytics`, {
        method: 'GET',
        headers: getAuthHeaders(request),
      }),

      // Payment analytics
      fetch(`${PAYMENT_SERVICE_URL}/api/v1/protected/payments/summary?date_from=${startDateStr}&date_to=${endDateStr}`, {
        method: 'GET',
        headers: getAuthHeaders(request),
      }),

      // Top products
      fetch(`${PRODUCT_SERVICE_URL}/api/v1/orders/analytics/top-products?limit=5`, {
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

    // Process payment analytics
    let paymentData = null
    if (paymentAnalytics.status === 'fulfilled' && paymentAnalytics.value.ok) {
      const paymentResponse = await paymentAnalytics.value.json()
      paymentData = paymentResponse.data
    }

    const blogData = blogSummary

    // Process top products
    let topProductsData = null
    if (topProducts.status === 'fulfilled' && topProducts.value.ok) {
      const topProductsResponse = await topProducts.value.json()
      topProductsData = topProductsResponse.data
    }

    // Combine all analytics data
    const analyticsData = {
      orders: orderData,
      products: productData,
      users: userData,
      payments: paymentData,
      blog: blogData,
      top_products: topProductsData,
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
