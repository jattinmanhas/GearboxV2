import { NextRequest, NextResponse } from 'next/server';

const BLOG_SERVICE_URL = process.env.BLOG_SERVICE_URL || 'http://localhost:3003/api/v1';

// GET /api/v1/blog/posts/search - Search posts
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url);
    const queryString = searchParams.toString();
    
    const response = await fetch(`${BLOG_SERVICE_URL}/posts/search?${queryString}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      'Authorization': request.headers.get('authorization') || '',
      },
    });

    const data = await response.json();
    
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Error searching posts:', error);
    return NextResponse.json(
      { 
        success: false, 
        message: 'Failed to search posts',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    );
  }
}
