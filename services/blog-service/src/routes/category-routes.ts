import { FastifyInstance } from 'fastify';
import { CategoryHandler } from '../handlers/category-handler';
import { AuthMiddleware } from '../middleware/auth';
import { AuthService } from '../services/auth-service';

const categoryHandler = new CategoryHandler();
const authService = new AuthService();
const authMiddleware = new AuthMiddleware(authService);

export async function categoryRoutes(fastify: FastifyInstance) {
  // Create a new category (requires authentication)
  fastify.post('/categories', {
    preHandler: authMiddleware.requireAuth,
    schema: {
      description: 'Create a new category',
      tags: ['Categories'],
      body: {
        type: 'object',
        required: ['name'],
        properties: {
          name: { type: 'string', minLength: 1, maxLength: 255 },
          description: { type: 'string', maxLength: 1000 },
          color: { type: 'string', pattern: '^#[0-9A-Fa-f]{6}$' }
        }
      },
      response: {
        201: {
          type: 'object',
          properties: {
            timestamp: { type: 'string' },
            status: { type: 'number' },
            success: { type: 'boolean' },
            message: { type: 'string' },
            data: { type: 'object' }
          }
        }
      }
    }
  }, categoryHandler.createCategory.bind(categoryHandler));

  // Get all categories with filters
  fastify.get('/categories', {
    schema: {
      description: 'Get all categories with optional filters',
      tags: ['Categories'],
      querystring: {
        type: 'object',
        properties: {
          search: { type: 'string' },
          page: { type: 'string', pattern: '^[0-9]+$' },
          limit: { type: 'string', pattern: '^[0-9]+$' }
        }
      }
    }
  }, categoryHandler.getCategories.bind(categoryHandler));

  // Get all categories (simple list)
  fastify.get('/categories/all', {
    schema: {
      description: 'Get all categories as a simple list',
      tags: ['Categories']
    }
  }, categoryHandler.getAllCategories.bind(categoryHandler));

  // Get category by ID
  fastify.get('/categories/:id', {
    schema: {
      description: 'Get a category by ID',
      tags: ['Categories'],
      params: {
        type: 'object',
        properties: {
          id: { type: 'string', format: 'uuid' }
        },
        required: ['id']
      }
    }
  }, categoryHandler.getCategoryById.bind(categoryHandler));

  // Get category by slug
  fastify.get('/categories/slug/:slug', {
    schema: {
      description: 'Get a category by slug',
      tags: ['Categories'],
      params: {
        type: 'object',
        properties: {
          slug: { type: 'string' }
        },
        required: ['slug']
      }
    }
  }, categoryHandler.getCategoryBySlug.bind(categoryHandler));

  // Update category (requires authentication)
  fastify.put('/categories/:id', {
    preHandler: authMiddleware.requireAuth,
    schema: {
      description: 'Update a category',
      tags: ['Categories'],
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
          name: { type: 'string', minLength: 1, maxLength: 255 },
          description: { type: 'string', maxLength: 1000 },
          color: { type: 'string', pattern: '^#[0-9A-Fa-f]{6}$' }
        }
      },
      response: {
        200: {
          type: 'object',
          properties: {
            timestamp: { type: 'string' },
            status: { type: 'number' },
            success: { type: 'boolean' },
            message: { type: 'string' },
            data: { type: 'object' }
          }
        }
      }
    }
  }, categoryHandler.updateCategory.bind(categoryHandler));

  // Delete category (requires authentication)
  fastify.delete('/categories/:id', {
    preHandler: authMiddleware.requireAuth,
    schema: {
      description: 'Delete a category',
      tags: ['Categories'],
      params: {
        type: 'object',
        properties: {
          id: { type: 'string', format: 'uuid' }
        },
        required: ['id']
      },
      response: {
        200: {
          type: 'object',
          properties: {
            timestamp: { type: 'string' },
            status: { type: 'number' },
            success: { type: 'boolean' },
            message: { type: 'string' }
          }
        }
      }
    }
  }, categoryHandler.deleteCategory.bind(categoryHandler));
}
