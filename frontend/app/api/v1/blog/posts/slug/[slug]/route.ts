import { NextRequest, NextResponse } from 'next/server';

const BLOG_SERVICE_URL = process.env.BLOG_SERVICE_URL || 'http://localhost:3003/api/v1';

// GET /api/v1/blog/posts/slug/[slug] - Get blog post by slug
export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ slug: string }> }
) {
  try {
    const { slug } = await params;

    const response = await fetch(`${BLOG_SERVICE_URL}/posts/slug/${slug}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      'Authorization': request.headers.get('authorization') || '',
      },
    });

    const data = await response.json();

    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Error fetching blog post by slug:', error);
    return NextResponse.json(
      {
        success: false,
        message: 'Failed to fetch blog post',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    );
  }
}
