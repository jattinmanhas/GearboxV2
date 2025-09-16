import { NextRequest, NextResponse } from 'next/server';

const BLOG_SERVICE_URL = process.env.BLOG_SERVICE_URL || 'http://localhost:3001/api/v1';

// GET /api/v1/blog/categories/slug/[slug] - Get category by slug
export async function GET(
  request: NextRequest,
  { params }: { params: { slug: string } }
) {
  try {
    const { slug } = params;
    
    const response = await fetch(`${BLOG_SERVICE_URL}/categories/slug/${slug}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    const data = await response.json();
    
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Error fetching category by slug:', error);
    return NextResponse.json(
      { 
        success: false, 
        message: 'Failed to fetch category',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    );
  }
}
