import { FastifyRequest, FastifyReply } from 'fastify';
import { CategoryService } from '../services/category-service';
import { 
  createCategorySchema, 
  updateCategorySchema, 
  categoryQuerySchema 
} from '../validation/blog-validation';

const categoryService = new CategoryService();

export class CategoryHandler {
  async createCategory(request: FastifyRequest, reply: FastifyReply) {
    try {
      const data = createCategorySchema.parse(request.body);
      const category = await categoryService.createCategory(data);
      
      return reply.status(201).send({
        success: true,
        data: category,
        message: 'Category created successfully'
      });
    } catch (error) {
      if (error instanceof Error) {
        return reply.status(400).send({
          success: false,
          message: error.message
        });
      }
      
      return reply.status(500).send({
        success: false,
        message: 'Internal server error'
      });
    }
  }

  async getCategoryById(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { id } = request.params as { id: string };
      const category = await categoryService.getCategoryById(id);
      
      if (!category) {
        return reply.status(404).send({
          success: false,
          message: 'Category not found'
        });
      }
      
      return reply.send({
        success: true,
        data: category
      });
    } catch (error) {
      return reply.status(500).send({
        success: false,
        message: 'Internal server error'
      });
    }
  }

  async getCategoryBySlug(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { slug } = request.params as { slug: string };
      const category = await categoryService.getCategoryBySlug(slug);
      
      if (!category) {
        return reply.status(404).send({
          success: false,
          message: 'Category not found'
        });
      }
      
      return reply.send({
        success: true,
        data: category
      });
    } catch (error) {
      return reply.status(500).send({
        success: false,
        message: 'Internal server error'
      });
    }
  }

  async getCategories(request: FastifyRequest, reply: FastifyReply) {
    try {
      const filters = categoryQuerySchema.parse(request.query);
      const result = await categoryService.getCategories(
        filters.search,
        filters.page,
        filters.limit
      );
      
      return reply.send({
        success: true,
        data: result
      });
    } catch (error) {
      if (error instanceof Error) {
        return reply.status(400).send({
          success: false,
          message: error.message
        });
      }
      
      return reply.status(500).send({
        success: false,
        message: 'Internal server error'
      });
    }
  }

  async getAllCategories(request: FastifyRequest, reply: FastifyReply) {
    try {
      const categories = await categoryService.getAllCategories();
      
      return reply.send({
        success: true,
        data: categories
      });
    } catch (error) {
      return reply.status(500).send({
        success: false,
        message: 'Internal server error'
      });
    }
  }

  async updateCategory(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { id } = request.params as { id: string };
      const data = updateCategorySchema.parse(request.body);
      
      const category = await categoryService.updateCategory(id, data);
      
      if (!category) {
        return reply.status(404).send({
          success: false,
          message: 'Category not found'
        });
      }
      
      return reply.send({
        success: true,
        data: category,
        message: 'Category updated successfully'
      });
    } catch (error) {
      if (error instanceof Error) {
        return reply.status(400).send({
          success: false,
          message: error.message
        });
      }
      
      return reply.status(500).send({
        success: false,
        message: 'Internal server error'
      });
    }
  }

  async deleteCategory(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { id } = request.params as { id: string };
      const deleted = await categoryService.deleteCategory(id);
      
      if (!deleted) {
        return reply.status(404).send({
          success: false,
          message: 'Category not found'
        });
      }
      
      return reply.send({
        success: true,
        message: 'Category deleted successfully'
      });
    } catch (error) {
      if (error instanceof Error) {
        return reply.status(400).send({
          success: false,
          message: error.message
        });
      }
      
      return reply.status(500).send({
        success: false,
        message: 'Internal server error'
      });
    }
  }
}
