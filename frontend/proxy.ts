import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export default function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // Handle Dashboard Authorization
  if (pathname.startsWith('/dashboard')) {
    const userRole = request.cookies.get('user-role')?.value;

    // Only allow admin and editor roles
    if (userRole !== 'admin' && userRole !== 'editor') {
      const url = request.nextUrl.clone();
      url.pathname = '/';
      return NextResponse.redirect(url);
    }
  }

  // Handle CORS for blog API routes
  if (pathname.startsWith('/api/v1/blog')) {
    const response = NextResponse.next();

    response.headers.set('Access-Control-Allow-Origin', '*');
    response.headers.set('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
    response.headers.set('Access-Control-Allow-Headers', 'Content-Type, Authorization');

    if (request.method === 'OPTIONS') {
      return new Response(null, { status: 200, headers: response.headers });
    }

    return response;
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/api/v1/blog/:path*', '/dashboard/:path*'],
};
