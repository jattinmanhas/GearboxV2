import { NextRequest, NextResponse } from 'next/server';

const BLOG_SERVICE_URL = process.env.BLOG_SERVICE_URL || 'http://localhost:3003/api/v1';

// GET /api/v1/blog/posts/[id]/related - Get related posts
export async function GET(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const { id } = params;
    const { searchParams } = new URL(request.url);
    const limit = searchParams.get('limit') || '5';
    
    const response = await fetch(`${BLOG_SERVICE_URL}/posts/${id}/related?limit=${limit}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    const data = await response.json();
    
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Error fetching related posts:', error);
    return NextResponse.json(
      { 
        success: false, 
        message: 'Failed to fetch related posts',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    );
  }
}
