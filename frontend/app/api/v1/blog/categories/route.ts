import { NextRequest, NextResponse } from 'next/server';

const BLOG_SERVICE_URL = process.env.BLOG_SERVICE_URL || 'http://localhost:3003/api/v1';

// Helper function to get auth headers
function getAuthHeaders(request: NextRequest) {  
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Cookie': request.headers.get('cookie') || '',
  };


  return headers;
}

// GET /api/v1/blog/categories - Get all categories with filters
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url);
    const queryString = searchParams.toString();
    
    const response = await fetch(`${BLOG_SERVICE_URL}/categories?${queryString}`, {
      method: 'GET',
      headers: getAuthHeaders(request),
    });

    const data = await response.json();
    
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Error fetching categories:', error);
    return NextResponse.json(
      { 
        success: false, 
        message: 'Failed to fetch categories',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    );
  }
}

// POST /api/v1/blog/categories - Create a new category
export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    
    const response = await fetch(`${BLOG_SERVICE_URL}/categories`, {
      method: 'POST',
      headers: getAuthHeaders(request),
      body: JSON.stringify(body),
    });

    const data = await response.json();
    
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Error creating category:', error);
    return NextResponse.json(
      { 
        success: false, 
        message: 'Failed to create category',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    );
  }
}
