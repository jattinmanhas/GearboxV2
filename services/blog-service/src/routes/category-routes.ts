import { FastifyInstance } from 'fastify';
import { CategoryHandler } from '../handlers/category-handler';

const categoryHandler = new CategoryHandler();

export async function categoryRoutes(fastify: FastifyInstance) {
  // Create a new category
  fastify.post('/categories', {
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
            success: { type: 'boolean' },
            data: { type: 'object' },
            message: { type: 'string' }
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

  // Update category
  fastify.put('/categories/:id', {
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
      }
    }
  }, categoryHandler.updateCategory.bind(categoryHandler));

  // Delete category
  fastify.delete('/categories/:id', {
    schema: {
      description: 'Delete a category',
      tags: ['Categories'],
      params: {
        type: 'object',
        properties: {
          id: { type: 'string', format: 'uuid' }
        },
        required: ['id']
      }
    }
  }, categoryHandler.deleteCategory.bind(categoryHandler));
}
