# Blog API Routes

This document describes the Next.js API routes created for the blog functionality. These routes act as a proxy between the frontend and the blog service backend.

## 🏗️ Architecture

The API routes are located in `app/api/v1/blog/` and provide a clean interface for the frontend to interact with the blog service backend.

## 📡 API Routes

### Blog Posts

#### `GET /api/v1/blog/posts`
- **Description**: Get all blog posts with optional filters
- **Query Parameters**: 
  - `status` - Filter by status (draft, published, archived)
  - `authorId` - Filter by author ID
  - `categoryId` - Filter by category ID
  - `tags` - Filter by tags (comma-separated)
  - `search` - Search in title, content, excerpt
  - `page` - Page number (default: 1)
  - `limit` - Items per page (default: 10)
  - `sortBy` - Sort field (createdAt, updatedAt, publishedAt, viewCount, title)
  - `sortOrder` - Sort direction (asc, desc)
- **Authentication**: Not required
- **Response**: `BlogPostListResponse`

#### `POST /api/v1/blog/posts`
- **Description**: Create a new blog post
- **Authentication**: Required (Bearer token or cookie)
- **Request Body**: `CreateBlogPostRequest`
- **Response**: `BlogPost`

#### `GET /api/v1/blog/posts/[id]`
- **Description**: Get blog post by ID
- **Authentication**: Not required
- **Response**: `BlogPost`

#### `PUT /api/v1/blog/posts/[id]`
- **Description**: Update blog post
- **Authentication**: Required (Bearer token or cookie)
- **Request Body**: `UpdateBlogPostRequest`
- **Response**: `BlogPost`

#### `DELETE /api/v1/blog/posts/[id]`
- **Description**: Delete blog post
- **Authentication**: Required (Bearer token or cookie)
- **Response**: Success message

#### `GET /api/v1/blog/posts/slug/[slug]`
- **Description**: Get blog post by slug
- **Authentication**: Not required
- **Response**: `BlogPost`

#### `GET /api/v1/blog/posts/author/[authorId]`
- **Description**: Get posts by author
- **Query Parameters**: `limit` (default: 10)
- **Authentication**: Not required
- **Response**: `BlogPost[]`

#### `GET /api/v1/blog/posts/[id]/related`
- **Description**: Get related posts
- **Query Parameters**: `limit` (default: 5)
- **Authentication**: Not required
- **Response**: `BlogPost[]`

#### `GET /api/v1/blog/posts/popular`
- **Description**: Get popular posts
- **Query Parameters**: `limit` (default: 10)
- **Authentication**: Not required
- **Response**: `BlogPost[]`

#### `GET /api/v1/blog/posts/recent`
- **Description**: Get recent posts
- **Query Parameters**: `limit` (default: 10)
- **Authentication**: Not required
- **Response**: `BlogPost[]`

#### `GET /api/v1/blog/posts/search`
- **Description**: Search posts
- **Query Parameters**: 
  - `q` - Search query (required)
  - All other filter parameters from GET /posts
- **Authentication**: Not required
- **Response**: `BlogPostListResponse`

### Categories

#### `GET /api/v1/blog/categories`
- **Description**: Get all categories with optional filters
- **Query Parameters**:
  - `search` - Search in category name
  - `page` - Page number (default: 1)
  - `limit` - Items per page (default: 10)
- **Authentication**: Not required
- **Response**: `CategoryListResponse`

#### `POST /api/v1/blog/categories`
- **Description**: Create a new category
- **Authentication**: Required (Bearer token or cookie)
- **Request Body**: `CreateCategoryRequest`
- **Response**: `Category`

#### `GET /api/v1/blog/categories/all`
- **Description**: Get all categories as a simple list
- **Authentication**: Not required
- **Response**: `Category[]`

#### `GET /api/v1/blog/categories/[id]`
- **Description**: Get category by ID
- **Authentication**: Not required
- **Response**: `Category`

#### `PUT /api/v1/blog/categories/[id]`
- **Description**: Update category
- **Authentication**: Required (Bearer token or cookie)
- **Request Body**: `UpdateCategoryRequest`
- **Response**: `Category`

#### `DELETE /api/v1/blog/categories/[id]`
- **Description**: Delete category
- **Authentication**: Required (Bearer token or cookie)
- **Response**: Success message

#### `GET /api/v1/blog/categories/slug/[slug]`
- **Description**: Get category by slug
- **Authentication**: Not required
- **Response**: `Category`

## 🔐 Authentication

### How Authentication Works

The API routes automatically handle authentication by:

1. **Checking Authorization Header**: Looks for `Bearer <token>` in the `Authorization` header
2. **Checking Cookies**: Falls back to `access_token` cookie if no header is present
3. **Forwarding to Backend**: Passes the authentication token to the blog service backend
4. **Returning Response**: Returns the backend response with appropriate status codes

### Authentication Required Routes

- `POST /api/v1/blog/posts` - Create post
- `PUT /api/v1/blog/posts/[id]` - Update post
- `DELETE /api/v1/blog/posts/[id]` - Delete post
- `POST /api/v1/blog/categories` - Create category
- `PUT /api/v1/blog/categories/[id]` - Update category
- `DELETE /api/v1/blog/categories/[id]` - Delete category

### Public Routes

All other routes are public and don't require authentication.

## ⚙️ Configuration

### Environment Variables

Add to your `.env.local`:

```env
# Blog service backend URL
BLOG_SERVICE_URL=http://localhost:3001/api/v1
```

### Backend Requirements

- Blog service must be running on the configured URL
- Blog service must have CORS enabled for the frontend domain
- Blog service must support the same authentication methods

## 🚀 Usage Examples

### Frontend API Calls

```typescript
// Get all posts
const response = await fetch('/api/v1/blog/posts');
const data = await response.json();

// Get posts with filters
const response = await fetch('/api/v1/blog/posts?status=published&page=1&limit=10');
const data = await response.json();

// Create a new post (requires authentication)
const response = await fetch('/api/v1/blog/posts', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}` // or cookie will be used automatically
  },
  body: JSON.stringify({
    title: 'My Post',
    content: '<p>Content here</p>',
    authorId: 'user-id',
    authorName: 'Author Name',
    authorEmail: 'author@example.com'
  })
});

// Search posts
const response = await fetch('/api/v1/blog/posts/search?q=javascript&limit=5');
const data = await response.json();
```

### Error Handling

All API routes include proper error handling:

```typescript
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error information"
}
```

## 🔄 Request Flow

1. **Frontend** makes request to `/api/v1/blog/*`
2. **Next.js API Route** receives request
3. **API Route** extracts authentication (header or cookie)
4. **API Route** forwards request to blog service backend
5. **Blog Service** processes request and returns response
6. **API Route** returns response to frontend

## 📊 Benefits

- **Centralized Authentication**: Handles auth logic in one place
- **CORS Handling**: No CORS issues between frontend and backend
- **Error Standardization**: Consistent error responses
- **Request Logging**: Easy to log and monitor API calls
- **Caching**: Can add caching layer in the future
- **Rate Limiting**: Can add rate limiting if needed
- **Request Transformation**: Can modify requests/responses as needed

## 🐛 Troubleshooting

### Common Issues

1. **502 Bad Gateway**
   - Check if blog service is running
   - Verify `BLOG_SERVICE_URL` environment variable
   - Check blog service logs

2. **401 Unauthorized**
   - Check if user is logged in
   - Verify JWT token is valid
   - Check if token is being sent correctly

3. **404 Not Found**
   - Verify the API route path is correct
   - Check if the blog service endpoint exists

4. **500 Internal Server Error**
   - Check Next.js API route logs
   - Check blog service logs
   - Verify request body format

## 📝 Notes

- All API routes are located in `app/api/v1/blog/`
- Authentication is handled automatically
- Error responses follow the same format as the backend
- Query parameters are passed through to the backend
- Request bodies are forwarded as-is to the backend
