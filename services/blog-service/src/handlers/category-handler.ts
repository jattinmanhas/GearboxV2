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
      
      if (!category) {
        return ResponseHelper.internalServerError(reply, 'Failed to create category - no data returned');
      }
      
      // Serialize dates to strings for JSON response
      const serializedCategory = {
        id: category.id,
        name: category.name,
        slug: category.slug,
        description: category.description,
        color: category.color,
        createdAt: category.createdAt.toISOString(),
        updatedAt: category.updatedAt.toISOString(),
      };
      
      return ResponseHelper.created(reply, 'Category created successfully', serializedCategory);
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
      
      // Serialize dates to strings for JSON response
      const serializedCategory = {
        ...category,
        createdAt: category.createdAt.toISOString(),
        updatedAt: category.updatedAt.toISOString(),
      };
      
      return ResponseHelper.ok(reply, 'Category retrieved successfully', serializedCategory);
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
      
      // Serialize dates to strings for JSON response
      const serializedCategory = {
        ...category,
        createdAt: category.createdAt.toISOString(),
        updatedAt: category.updatedAt.toISOString(),
      };
      
      return ResponseHelper.ok(reply, 'Category retrieved successfully', serializedCategory);
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
      
      // Serialize dates to strings for JSON response
      const serializedResult = {
        ...result,
        categories: result.categories.map(category => ({
          ...category,
          createdAt: category.createdAt.toISOString(),
          updatedAt: category.updatedAt.toISOString(),
        })),
      };
      
      return ResponseHelper.ok(reply, 'Categories retrieved successfully', serializedResult);
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
      
      // Serialize dates to strings for JSON response
      const serializedCategories = categories.map(category => ({
        ...category,
        createdAt: category.createdAt.toISOString(),
        updatedAt: category.updatedAt.toISOString(),
      }));
      
      return ResponseHelper.ok(reply, 'All categories retrieved successfully', serializedCategories);
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
      
      // Serialize dates to strings for JSON response
      const serializedCategory = {
        ...category,
        createdAt: category.createdAt.toISOString(),
        updatedAt: category.updatedAt.toISOString(),
      };
      
      return ResponseHelper.ok(reply, 'Category updated successfully', serializedCategory);
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
