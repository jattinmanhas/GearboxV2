# Blog Service

A modern, scalable blog service built with Node.js, TypeScript, Fastify, and PostgreSQL. This service provides a complete API for managing blog posts, categories, and related functionality.

## Features

- 📝 **Blog Post Management**: Create, read, update, and delete blog posts
- 🏷️ **Category System**: Organize posts with categories and tags
- 🔍 **Search & Filtering**: Advanced search and filtering capabilities
- 📊 **Analytics**: View count tracking and popular posts
- 🔗 **Related Posts**: Smart related post suggestions
- 📱 **RESTful API**: Clean, well-documented REST API
- 📚 **Swagger Documentation**: Interactive API documentation
- 🛡️ **Validation**: Comprehensive input validation with Zod
- 🚀 **Performance**: Optimized database queries and caching
- 🔒 **Security**: Rate limiting, CORS, and security headers

## Tech Stack

- **Runtime**: Node.js
- **Language**: TypeScript
- **Framework**: Fastify
- **Database**: PostgreSQL
- **ORM**: Drizzle ORM
- **Validation**: Zod
- **Documentation**: Swagger/OpenAPI

## Getting Started

### Prerequisites

- Node.js (v18 or higher)
- PostgreSQL (v12 or higher)
- npm or yarn

### Installation

1. Clone the repository and navigate to the blog service directory:
```bash
cd services/blog-service
```

2. Install dependencies:
```bash
npm install
```

3. Set up environment variables:
```bash
cp env.example .env
```

Edit the `.env` file with your database configuration:
```env
DATABASE_URL=postgresql://username:password@localhost:5432/blog_service
PORT=3003
HOST=0.0.0.0
```

4. Set up the database:
```bash
# Generate migrations
npm run db:generate

# Run migrations
npm run db:migrate

# Seed sample data
npm run db:seed
```

5. Start the development server:
```bash
npm run dev
```

The API will be available at `http://localhost:3003`

## API Documentation

Once the server is running, you can access the interactive API documentation at:
- Swagger UI: `http://localhost:3003/docs`

## API Endpoints

### Blog Posts

- `GET /api/v1/posts` - Get all blog posts with filters
- `POST /api/v1/posts` - Create a new blog post
- `GET /api/v1/posts/:id` - Get blog post by ID
- `GET /api/v1/posts/slug/:slug` - Get blog post by slug
- `PUT /api/v1/posts/:id` - Update blog post
- `DELETE /api/v1/posts/:id` - Delete blog post
- `GET /api/v1/posts/author/:authorId` - Get posts by author
- `GET /api/v1/posts/:id/related` - Get related posts
- `GET /api/v1/posts/popular` - Get popular posts
- `GET /api/v1/posts/recent` - Get recent posts
- `GET /api/v1/posts/search` - Search posts

### Categories

- `GET /api/v1/categories` - Get all categories with filters
- `POST /api/v1/categories` - Create a new category
- `GET /api/v1/categories/:id` - Get category by ID
- `GET /api/v1/categories/slug/:slug` - Get category by slug
- `PUT /api/v1/categories/:id` - Update category
- `DELETE /api/v1/categories/:id` - Delete category
- `GET /api/v1/categories/all` - Get all categories (simple list)

### Health Check

- `GET /health` - Health check endpoint

## Database Schema

### Blog Posts Table
- `id` (UUID, Primary Key)
- `title` (VARCHAR)
- `slug` (VARCHAR, Unique)
- `content` (TEXT)
- `excerpt` (TEXT)
- `author_id` (VARCHAR)
- `author_name` (VARCHAR)
- `author_email` (VARCHAR)
- `status` (ENUM: draft, published, archived)
- `featured_image` (VARCHAR)
- `tags` (TEXT[])
- `category_id` (UUID, Foreign Key)
- `published_at` (TIMESTAMP)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)
- `view_count` (INTEGER)
- `read_time` (INTEGER)

### Categories Table
- `id` (UUID, Primary Key)
- `name` (VARCHAR)
- `slug` (VARCHAR, Unique)
- `description` (TEXT)
- `color` (VARCHAR)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

## Query Parameters

### Blog Posts Filtering
- `status`: Filter by post status (draft, published, archived)
- `authorId`: Filter by author ID
- `categoryId`: Filter by category ID
- `tags`: Filter by tags (comma-separated)
- `search`: Search in title, content, and excerpt
- `page`: Page number for pagination
- `limit`: Number of items per page
- `sortBy`: Sort field (createdAt, updatedAt, publishedAt, viewCount, title)
- `sortOrder`: Sort direction (asc, desc)

### Categories Filtering
- `search`: Search in category name
- `page`: Page number for pagination
- `limit`: Number of items per page

## Example Usage

### Create a Blog Post
```bash
curl -X POST http://localhost:3003/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Blog Post",
    "content": "This is the content of my blog post...",
    "authorId": "user-123",
    "authorName": "John Doe",
    "authorEmail": "john@example.com",
    "status": "published",
    "tags": ["tutorial", "beginner"],
    "categoryId": "category-uuid"
  }'
```

### Get Blog Posts with Filters
```bash
curl "http://localhost:3003/api/v1/posts?status=published&page=1&limit=10&sortBy=createdAt&sortOrder=desc"
```

### Search Posts
```bash
curl "http://localhost:3003/api/v1/posts/search?q=typescript&status=published"
```

## Development

### Available Scripts

- `npm run dev` - Start development server with hot reload
- `npm run start` - Start production server
- `npm run build` - Build TypeScript to JavaScript
- `npm run db:generate` - Generate database migrations
- `npm run db:migrate` - Run database migrations
- `npm run db:push` - Push schema changes to database
- `npm run db:studio` - Open Drizzle Studio
- `npm run db:seed` - Seed database with sample data

### Project Structure

```
src/
├── config/          # Configuration files
├── domain/          # Domain logic and business rules
├── handlers/        # Request handlers/controllers
├── repositories/    # Data access layer
├── routes/          # API route definitions
├── services/        # Business logic layer
├── types/           # TypeScript type definitions
├── validation/      # Input validation schemas
├── scripts/         # Utility scripts
└── app.ts          # Main application file
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the ISC License.
