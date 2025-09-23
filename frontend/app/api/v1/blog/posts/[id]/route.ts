import { NextRequest, NextResponse } from 'next/server';

const BLOG_SERVICE_URL = process.env.BLOG_SERVICE_URL || 'http://localhost:3003/api/v1';

// Helper function to get auth headers
function getAuthHeaders(request: NextRequest, includeContentType: boolean = true) {
  const headers: Record<string, string> = {
    'Cookie': request.headers.get('cookie') || '',
  };

  if (includeContentType) {
    headers['Content-Type'] = 'application/json';
  }

  return headers;
}

// GET /api/v1/blog/posts/[id] - Get blog post by ID
export async function GET(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const { id } = params;
    
    const response = await fetch(`${BLOG_SERVICE_URL}/posts/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    const data = await response.json();
    
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Error fetching blog post:', error);
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

// PUT /api/v1/blog/posts/[id] - Update blog post
export async function PUT(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const { id } = params;
    const body = await request.json();
    
    const response = await fetch(`${BLOG_SERVICE_URL}/posts/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(request),
      body: JSON.stringify(body),
    });

    const data = await response.json();
    
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Error updating blog post:', error);
    return NextResponse.json(
      { 
        success: false, 
        message: 'Failed to update blog post',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    );
  }
}

// DELETE /api/v1/blog/posts/[id] - Delete blog post
export async function DELETE(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const { id } = params;
    
    // First, get the blog post data to extract image URLs
    const getResponse = await fetch(`${BLOG_SERVICE_URL}/posts/${id}`, {
      method: 'GET',
      headers: getAuthHeaders(request, false),
    });

    let postData = null
    if (getResponse.ok) {
      postData = await getResponse.json()
    }
    
    const response = await fetch(`${BLOG_SERVICE_URL}/posts/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders(request, false), // Don't include Content-Type for DELETE
    });

    const data = await response.json();

    // If deletion was successful and we have post data, clean up images
    if (response.ok && postData) {
      try {
        const { cleanupEntityImages } = await import('@/lib/cloudinary-cleanup')
        const cleanupResult = await cleanupEntityImages(postData.data || postData)
        
        if (cleanupResult.deleted > 0) {
          console.log(`Cleaned up ${cleanupResult.deleted} images for deleted blog post ${id}`)
        }
        
        if (cleanupResult.errors.length > 0) {
          console.warn('Some images could not be deleted:', cleanupResult.errors)
        }
      } catch (cleanupError) {
        console.error('Failed to cleanup images for deleted blog post:', cleanupError)
        // Don't fail the delete operation if image cleanup fails
      }
    }
    
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Error deleting blog post:', error);
    return NextResponse.json(
      { 
        success: false, 
        message: 'Failed to delete blog post',
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    );
  }
}
