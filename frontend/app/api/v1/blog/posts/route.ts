import { NextRequest, NextResponse } from 'next/server';

const BLOG_SERVICE_URL = process.env.BLOG_SERVICE_URL || 'http://localhost:3003/api/v1';

// Helper function to get auth headers
function getAuthHeaders(request: NextRequest, includeContentType: boolean = true) {
  const headers: Record<string, string> = {
    'Cookie': request.headers.get('cookie') || '',
  'Authorization': request.headers.get('authorization') || '',
  };

  if (includeContentType) {
    headers['Content-Type'] = 'application/json';
  }

  return headers;
}

// GET /api/v1/blog/posts - Get all blog posts with filters
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url);
    const queryString = searchParams.toString();
    
    const response = await fetch(`${BLOG_SERVICE_URL}/posts?${queryString}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Check if response is ok and has content
    if (!response.ok) {
      console.error('Blog service error:', response.status, response.statusText);
      return NextResponse.json(
        { 
          success: false, 
          message: `Blog service error: ${response.status} ${response.statusText}`,
          error: 'Service unavailable'
        },
        { status: response.status }
      );
    }

    // Check if response has JSON content
    const contentType = response.headers.get('content-type');
    if (!contentType || !contentType.includes('application/json')) {
      console.error('Non-JSON response from blog service:', contentType);
      return NextResponse.json(
        { 
          success: false, 
          message: 'Invalid response format from blog service',
          error: 'Service returned non-JSON response'
        },
        { status: 502 }
      );
    }

    const data = await response.json();
    
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Error fetching blog posts:', error);
    return NextResponse.json(
      { 
        success: false, 
        message: 'Failed to fetch blog posts',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    );
  }
}

// POST /api/v1/blog/posts - Create a new blog post
export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    
    const response = await fetch(`${BLOG_SERVICE_URL}/posts`, {
      method: 'POST',
      headers: getAuthHeaders(request),
      body: JSON.stringify(body),
    });

    const data = await response.json();
    
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Error creating blog post:', error);
    return NextResponse.json(
      { 
        success: false, 
        message: 'Failed to create blog post',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    );
  }
}
