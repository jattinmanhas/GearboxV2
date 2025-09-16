import { FastifyRequest, FastifyReply } from 'fastify';
import { BlogService } from '../services/blog-service';
import { 
  createBlogPostSchema, 
  updateBlogPostSchema, 
  blogPostQuerySchema 
} from '../validation/blog-validation';
import { ResponseHelper } from '../utils/response';

const blogService = new BlogService();

export class BlogHandler {
  async createPost(request: FastifyRequest, reply: FastifyReply) {
    try {
      const data = createBlogPostSchema.parse(request.body);
      const post = await blogService.createPost(data);
      
      return ResponseHelper.created(reply, 'Blog post created successfully', post);
    } catch (error) {
      if (error instanceof Error) {
        return ResponseHelper.badRequest(reply, error.message);
      }
      
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async getPostById(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { id } = request.params as { id: string };
      const post = await blogService.getPostById(id);
      
      if (!post) {
        return ResponseHelper.notFound(reply, 'Blog post not found');
      }

      // Increment view count
      await blogService.incrementViewCount(id);
      
      return ResponseHelper.ok(reply, 'Blog post retrieved successfully', post);
    } catch (error) {
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async getPostBySlug(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { slug } = request.params as { slug: string };
      const post = await blogService.getPostBySlug(slug);
      
      if (!post) {
        return ResponseHelper.notFound(reply, 'Blog post not found');
      }

      // Increment view count
      await blogService.incrementViewCount(post.id);
      
      return ResponseHelper.ok(reply, 'Blog post retrieved successfully', post);
    } catch (error) {
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async getPosts(request: FastifyRequest, reply: FastifyReply) {
    try {
      const filters = blogPostQuerySchema.parse(request.query);
      const result = await blogService.getPosts(filters);
      
      return ResponseHelper.ok(reply, 'Blog posts retrieved successfully', result);
    } catch (error) {
      if (error instanceof Error) {
        return ResponseHelper.badRequest(reply, error.message);
      }
      
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async updatePost(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { id } = request.params as { id: string };
      const data = updateBlogPostSchema.parse(request.body);
      
      const post = await blogService.updatePost(id, data);
      
      if (!post) {
        return ResponseHelper.notFound(reply, 'Blog post not found');
      }
      
      return ResponseHelper.ok(reply, 'Blog post updated successfully', post);
    } catch (error) {
      if (error instanceof Error) {
        return ResponseHelper.badRequest(reply, error.message);
      }
      
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async deletePost(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { id } = request.params as { id: string };
      const deleted = await blogService.deletePost(id);
      
      if (!deleted) {
        return ResponseHelper.notFound(reply, 'Blog post not found');
      }
      
      return ResponseHelper.ok(reply, 'Blog post deleted successfully');
    } catch (error) {
      if (error instanceof Error) {
        return ResponseHelper.badRequest(reply, error.message);
      }
      
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async getPostsByAuthor(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { authorId } = request.params as { authorId: string };
      const { limit } = request.query as { limit?: string };
      
      const posts = await blogService.getPostsByAuthor(
        authorId, 
        limit ? parseInt(limit, 10) : 10
      );
      
      return ResponseHelper.ok(reply, 'Author posts retrieved successfully', posts);
    } catch (error) {
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async getRelatedPosts(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { id } = request.params as { id: string };
      const { limit } = request.query as { limit?: string };
      
      const posts = await blogService.getRelatedPosts(
        id, 
        limit ? parseInt(limit, 10) : 5
      );
      
      return ResponseHelper.ok(reply, 'Related posts retrieved successfully', posts);
    } catch (error) {
      if (error instanceof Error) {
        return ResponseHelper.badRequest(reply, error.message);
      }
      
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async getPopularPosts(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { limit } = request.query as { limit?: string };
      
      const posts = await blogService.getPopularPosts(
        limit ? parseInt(limit, 10) : 10
      );
      
      return ResponseHelper.ok(reply, 'Popular posts retrieved successfully', posts);
    } catch (error) {
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async getRecentPosts(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { limit } = request.query as { limit?: string };
      
      const posts = await blogService.getRecentPosts(
        limit ? parseInt(limit, 10) : 10
      );
      
      return ResponseHelper.ok(reply, 'Recent posts retrieved successfully', posts);
    } catch (error) {
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }

  async searchPosts(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { q } = request.query as { q: string };
      const filters = blogPostQuerySchema.parse(request.query);
      
      if (!q) {
        return ResponseHelper.badRequest(reply, 'Search query is required');
      }
      
      const result = await blogService.searchPosts(q, filters);
      
      return ResponseHelper.ok(reply, 'Search results retrieved successfully', result);
    } catch (error) {
      if (error instanceof Error) {
        return ResponseHelper.badRequest(reply, error.message);
      }
      
      return ResponseHelper.internalServerError(reply, 'Internal server error');
    }
  }
}
