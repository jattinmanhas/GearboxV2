import { NextRequest, NextResponse } from 'next/server'

const BLOG_SERVICE_URL = process.env.BLOG_SERVICE_URL || 'http://localhost:3003'

function getAuthHeaders(request: NextRequest) {
  return {
    'Content-Type': 'application/json',
    'Cookie': request.headers.get('cookie') || '',
    'Authorization': request.headers.get('authorization') || '',
  }
}

async function fetchBlogList(
  request: NextRequest,
  params: Record<string, string | number | undefined>
) {
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

// GET /api/v1/dashboard/blog - Get blog analytics for dashboard
export async function GET(request: NextRequest) {
  try {
    const perPage = 100
    const maxPages = 20

    const totalResponse = await fetchBlogList(request, { limit: 1, page: 1 })
    const totalPosts = totalResponse?.data?.total || 0

    const publishedResponse = await fetchBlogList(request, { status: 'published', limit: 1, page: 1 })
    const publishedPosts = publishedResponse?.data?.total || 0

    const draftResponse = await fetchBlogList(request, { status: 'draft', limit: 1, page: 1 })
    const draftPosts = draftResponse?.data?.total || 0

    const archivedResponse = await fetchBlogList(request, { status: 'archived', limit: 1, page: 1 })
    const archivedPosts = archivedResponse?.data?.total || 0

    const topPostsResponse = await fetchBlogList(request, {
      status: 'published',
      sortBy: 'viewCount',
      sortOrder: 'desc',
      limit: 5,
      page: 1,
    })
    const topPosts = topPostsResponse?.data?.posts || []

    const recentPostsResponse = await fetchBlogList(request, {
      status: 'published',
      sortBy: 'publishedAt',
      sortOrder: 'desc',
      limit: 5,
      page: 1,
    })
    const recentPosts = recentPostsResponse?.data?.posts || []

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

    return NextResponse.json({
      success: true,
      message: 'Blog analytics retrieved successfully',
      data: {
        total_posts: totalPosts,
        published_posts: publishedPosts,
        draft_posts: draftPosts,
        archived_posts: archivedPosts,
        total_views: totalViews,
        average_read_time: Number(averageReadTime.toFixed(2)),
        top_posts: topPosts,
        recent_posts: recentPosts,
      },
    })
  } catch (error) {
    console.error('Blog analytics API error:', error)
    return NextResponse.json(
      {
        success: false,
        message: 'Failed to fetch blog analytics',
        error: error instanceof Error ? error.message : 'Unknown error',
      },
      { status: 500 }
    )
  }
}
