import { FastifyInstance } from 'fastify';
import { BlogHandler } from '../handlers/blog-handler';

const blogHandler = new BlogHandler();

export async function blogRoutes(fastify: FastifyInstance) {
  // Create a new blog post
  fastify.post('/posts', {
    schema: {
      description: 'Create a new blog post',
      tags: ['Blog Posts'],
      body: {
        type: 'object',
        required: ['title', 'content', 'authorId', 'authorName', 'authorEmail'],
        properties: {
          title: { type: 'string', minLength: 1, maxLength: 255 },
          content: { type: 'string', minLength: 1 },
          excerpt: { type: 'string', maxLength: 500 },
          authorId: { type: 'string', minLength: 1 },
          authorName: { type: 'string', minLength: 1, maxLength: 255 },
          authorEmail: { type: 'string', format: 'email' },
          status: { type: 'string', enum: ['draft', 'published'] },
          featuredImage: { type: 'string', format: 'uri' },
          tags: { type: 'array', items: { type: 'string' } },
          categoryId: { type: 'string', format: 'uuid' }
        }
      },
      response: {
        201: {
          type: 'object',
          properties: {
            success: { type: 'boolean' },
            data: { type: 'object' },
            message: { type: 'string' }
          }
        }
      }
    }
  }, blogHandler.createPost.bind(blogHandler));

  // Get all blog posts with filters
  fastify.get('/posts', {
    schema: {
      description: 'Get all blog posts with optional filters',
      tags: ['Blog Posts'],
      querystring: {
        type: 'object',
        properties: {
          status: { type: 'string', enum: ['draft', 'published', 'archived'] },
          authorId: { type: 'string' },
          categoryId: { type: 'string', format: 'uuid' },
          tags: { type: 'string' },
          search: { type: 'string' },
          page: { type: 'string', pattern: '^[0-9]+$' },
          limit: { type: 'string', pattern: '^[0-9]+$' },
          sortBy: { type: 'string', enum: ['createdAt', 'updatedAt', 'publishedAt', 'viewCount', 'title'] },
          sortOrder: { type: 'string', enum: ['asc', 'desc'] }
        }
      }
    }
  }, blogHandler.getPosts.bind(blogHandler));

  // Get blog post by ID
  fastify.get('/posts/:id', {
    schema: {
      description: 'Get a blog post by ID',
      tags: ['Blog Posts'],
      params: {
        type: 'object',
        properties: {
          id: { type: 'string', format: 'uuid' }
        },
        required: ['id']
      }
    }
  }, blogHandler.getPostById.bind(blogHandler));

  // Get blog post by slug
  fastify.get('/posts/slug/:slug', {
    schema: {
      description: 'Get a blog post by slug',
      tags: ['Blog Posts'],
      params: {
        type: 'object',
        properties: {
          slug: { type: 'string' }
        },
        required: ['slug']
      }
    }
  }, blogHandler.getPostBySlug.bind(blogHandler));

  // Update blog post
  fastify.put('/posts/:id', {
    schema: {
      description: 'Update a blog post',
      tags: ['Blog Posts'],
      params: {
        type: 'object',
        properties: {
          id: { type: 'string', format: 'uuid' }
        },
        required: ['id']
      },
      body: {
        type: 'object',
        properties: {
          title: { type: 'string', minLength: 1, maxLength: 255 },
          content: { type: 'string', minLength: 1 },
          excerpt: { type: 'string', maxLength: 500 },
          status: { type: 'string', enum: ['draft', 'published', 'archived'] },
          featuredImage: { type: 'string', format: 'uri' },
          tags: { type: 'array', items: { type: 'string' } },
          categoryId: { type: 'string', format: 'uuid' }
        }
      }
    }
  }, blogHandler.updatePost.bind(blogHandler));

  // Delete blog post
  fastify.delete('/posts/:id', {
    schema: {
      description: 'Delete a blog post',
      tags: ['Blog Posts'],
      params: {
        type: 'object',
        properties: {
          id: { type: 'string', format: 'uuid' }
        },
        required: ['id']
      }
    }
  }, blogHandler.deletePost.bind(blogHandler));

  // Get posts by author
  fastify.get('/posts/author/:authorId', {
    schema: {
      description: 'Get blog posts by author',
      tags: ['Blog Posts'],
      params: {
        type: 'object',
        properties: {
          authorId: { type: 'string' }
        },
        required: ['authorId']
      },
      querystring: {
        type: 'object',
        properties: {
          limit: { type: 'string', pattern: '^[0-9]+$' }
        }
      }
    }
  }, blogHandler.getPostsByAuthor.bind(blogHandler));

  // Get related posts
  fastify.get('/posts/:id/related', {
    schema: {
      description: 'Get related blog posts',
      tags: ['Blog Posts'],
      params: {
        type: 'object',
        properties: {
          id: { type: 'string', format: 'uuid' }
        },
        required: ['id']
      },
      querystring: {
        type: 'object',
        properties: {
          limit: { type: 'string', pattern: '^[0-9]+$' }
        }
      }
    }
  }, blogHandler.getRelatedPosts.bind(blogHandler));

  // Get popular posts
  fastify.get('/posts/popular', {
    schema: {
      description: 'Get popular blog posts',
      tags: ['Blog Posts'],
      querystring: {
        type: 'object',
        properties: {
          limit: { type: 'string', pattern: '^[0-9]+$' }
        }
      }
    }
  }, blogHandler.getPopularPosts.bind(blogHandler));

  // Get recent posts
  fastify.get('/posts/recent', {
    schema: {
      description: 'Get recent blog posts',
      tags: ['Blog Posts'],
      querystring: {
        type: 'object',
        properties: {
          limit: { type: 'string', pattern: '^[0-9]+$' }
        }
      }
    }
  }, blogHandler.getRecentPosts.bind(blogHandler));

  // Search posts
  fastify.get('/posts/search', {
    schema: {
      description: 'Search blog posts',
      tags: ['Blog Posts'],
      querystring: {
        type: 'object',
        required: ['q'],
        properties: {
          q: { type: 'string', minLength: 1 },
          status: { type: 'string', enum: ['draft', 'published', 'archived'] },
          authorId: { type: 'string' },
          categoryId: { type: 'string', format: 'uuid' },
          tags: { type: 'string' },
          page: { type: 'string', pattern: '^[0-9]+$' },
          limit: { type: 'string', pattern: '^[0-9]+$' },
          sortBy: { type: 'string', enum: ['createdAt', 'updatedAt', 'publishedAt', 'viewCount', 'title'] },
          sortOrder: { type: 'string', enum: ['asc', 'desc'] }
        }
      }
    }
  }, blogHandler.searchPosts.bind(blogHandler));
}
