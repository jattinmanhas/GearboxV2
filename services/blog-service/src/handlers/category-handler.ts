import { FastifyRequest, FastifyReply } from 'fastify';
import { CategoryService } from '../services/category-service';
import { 
  createCategorySchema, 
  updateCategorySchema, 
  categoryQuerySchema 
} from '../validation/blog-validation';
import { ResponseHelper } from '../utils/response';

const categoryService = new CategoryService();

export class CategoryHandler {
  async createCategory(request: FastifyRequest, reply: FastifyReply) {
    try {
      const data = createCategorySchema.parse(request.body);
      const category = await categoryService.createCategory(data);
      
      return ResponseHelper.created(reply, 'Category created successfully', category);
    } catch (error) {
      if (error instanceof Error) {
        return ResponseHelper.badRequest(reply, error.message);
      }
      
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async getCategoryById(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { id } = request.params as { id: string };
      const category = await categoryService.getCategoryById(id);
      
      if (!category) {
        return ResponseHelper.notFound(reply, 'Category not found');
      }
      
      return ResponseHelper.ok(reply, 'Category retrieved successfully', category);
    } catch (error) {
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async getCategoryBySlug(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { slug } = request.params as { slug: string };
      const category = await categoryService.getCategoryBySlug(slug);
      
      if (!category) {
        return ResponseHelper.notFound(reply, 'Category not found');
      }
      
      return ResponseHelper.ok(reply, 'Category retrieved successfully', category);
    } catch (error) {
      return ResponseHelper.internalServerError(reply, 'Internal server error');
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
      
      return ResponseHelper.ok(reply, 'Categories retrieved successfully', result);
    } catch (error) {
      if (error instanceof Error) {
        return ResponseHelper.badRequest(reply, error.message);
      }
      
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async getAllCategories(request: FastifyRequest, reply: FastifyReply) {
    try {
      const categories = await categoryService.getAllCategories();
      
      return ResponseHelper.ok(reply, 'All categories retrieved successfully', categories);
    } catch (error) {
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async updateCategory(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { id } = request.params as { id: string };
      const data = updateCategorySchema.parse(request.body);
      
      const category = await categoryService.updateCategory(id, data);
      
      if (!category) {
        return ResponseHelper.notFound(reply, 'Category not found');
      }
      
      return ResponseHelper.ok(reply, 'Category updated successfully', category);
    } catch (error) {
      if (error instanceof Error) {
        return ResponseHelper.badRequest(reply, error.message);
      }
      
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async deleteCategory(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { id } = request.params as { id: string };
      const deleted = await categoryService.deleteCategory(id);
      
      if (!deleted) {
        return ResponseHelper.notFound(reply, 'Category not found');
      }
      
      return ResponseHelper.ok(reply, 'Category deleted successfully');
    } catch (error) {
      if (error instanceof Error) {
        return ResponseHelper.badRequest(reply, error.message);
      }
      
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }
}
