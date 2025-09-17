import { NextRequest, NextResponse } from 'next/server';

const BLOG_SERVICE_URL = process.env.BLOG_SERVICE_URL || 'http://localhost:3003/api/v1';

// GET /api/v1/blog/posts/popular - Get popular posts
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url);
    const limit = searchParams.get('limit') || '10';
    
    const response = await fetch(`${BLOG_SERVICE_URL}/posts/popular?limit=${limit}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    const data = await response.json();
    
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Error fetching popular posts:', error);
    return NextResponse.json(
      { 
        success: false, 
        message: 'Failed to fetch popular posts',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    );
  }
}
