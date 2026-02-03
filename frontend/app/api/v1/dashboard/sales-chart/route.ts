import { NextRequest, NextResponse } from 'next/server'

const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8082'

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
    startDate,
    endDate,
    startDateStr: toDateString(startDate),
    endDateStr: toDateString(endDate),
  }
}

function buildEmptySeries(startDate: Date, endDate: Date) {
  const points: Array<{ date: string; revenue: number; orders: number }> = []
  const current = new Date(startDate)

  while (current <= endDate) {
    points.push({
      date: current.toISOString().split('T')[0],
      revenue: 0,
      orders: 0,
    })
    current.setDate(current.getDate() + 1)
  }

  return points
}

// GET /api/v1/dashboard/sales-chart - Get sales chart data for a specific period
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const period = searchParams.get('period') || '30d'
    const { startDate, endDate, startDateStr, endDateStr } = getDateRange(period)

    const limit = 200
    const maxPages = 25
    let page = 1
    let totalPages = 1
    const orders: Array<{ created_at: string; total_amount: number }> = []

    while (page <= totalPages && page <= maxPages) {
      const response = await fetch(
        `${PRODUCT_SERVICE_URL}/api/v1/orders?date_from=${startDateStr}&date_to=${endDateStr}&page=${page}&limit=${limit}&sort=created_at&order=asc`,
        {
          method: 'GET',
          headers: getAuthHeaders(request),
        }
      )

      if (!response.ok) {
        const errorData = await response.text()
        return NextResponse.json(
          {
            success: false,
            message: 'Failed to fetch orders for sales chart',
            error: errorData,
          },
          { status: response.status }
        )
      }

      const result = await response.json()
      const data = result?.data
      const pageOrders = data?.orders || []

      pageOrders.forEach((order: any) => {
        orders.push({
          created_at: order.created_at,
          total_amount: Number(order.total_amount) || 0,
        })
      })

      totalPages = data?.pages || 1
      page += 1
    }

    const chartData = buildEmptySeries(startDate, endDate)
    const indexByDate = new Map<string, number>()
    chartData.forEach((point, index) => {
      indexByDate.set(point.date, index)
    })

    orders.forEach((order) => {
      const dateKey = new Date(order.created_at).toISOString().split('T')[0]
      const idx = indexByDate.get(dateKey)
      if (idx !== undefined) {
        chartData[idx].revenue += order.total_amount
        chartData[idx].orders += 1
      }
    })

    return NextResponse.json({
      success: true,
      message: 'Sales chart data retrieved successfully',
      data: chartData,
    })
  } catch (error) {
    console.error('Sales chart API error:', error)
    return NextResponse.json(
      {
        success: false,
        message: 'Failed to fetch sales chart data',
        error: error instanceof Error ? error.message : 'Unknown error',
      },
      { status: 500 }
    )
  }
}
