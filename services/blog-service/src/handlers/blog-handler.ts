import { FastifyRequest, FastifyReply } from 'fastify';
import { BlogService } from '../services/blog-service';
import { 
  createBlogPostSchema, 
  updateBlogPostSchema, 
  blogPostQuerySchema 
} from '../validation/blog-validation';

const blogService = new BlogService();

export class BlogHandler {
  async createPost(request: FastifyRequest, reply: FastifyReply) {
    try {
      const data = createBlogPostSchema.parse(request.body);
      const post = await blogService.createPost(data);
      
      return reply.status(201).send({
        success: true,
        data: post,
        message: 'Blog post created successfully'
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

  async getPostById(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { id } = request.params as { id: string };
      const post = await blogService.getPostById(id);
      
      if (!post) {
        return reply.status(404).send({
          success: false,
          message: 'Blog post not found'
        });
      }

      // Increment view count
      await blogService.incrementViewCount(id);
      
      return reply.send({
        success: true,
        data: post
      });
    } catch (error) {
      return reply.status(500).send({
        success: false,
        message: 'Internal server error'
      });
    }
  }

  async getPostBySlug(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { slug } = request.params as { slug: string };
      const post = await blogService.getPostBySlug(slug);
      
      if (!post) {
        return reply.status(404).send({
          success: false,
          message: 'Blog post not found'
        });
      }

      // Increment view count
      await blogService.incrementViewCount(post.id);
      
      return reply.send({
        success: true,
        data: post
      });
    } catch (error) {
      return reply.status(500).send({
        success: false,
        message: 'Internal server error'
      });
    }
  }

  async getPosts(request: FastifyRequest, reply: FastifyReply) {
    try {
      const filters = blogPostQuerySchema.parse(request.query);
      const result = await blogService.getPosts(filters);
      
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

  async updatePost(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { id } = request.params as { id: string };
      const data = updateBlogPostSchema.parse(request.body);
      
      const post = await blogService.updatePost(id, data);
      
      if (!post) {
        return reply.status(404).send({
          success: false,
          message: 'Blog post not found'
        });
      }
      
      return reply.send({
        success: true,
        data: post,
        message: 'Blog post updated successfully'
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

  async deletePost(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { id } = request.params as { id: string };
      const deleted = await blogService.deletePost(id);
      
      if (!deleted) {
        return reply.status(404).send({
          success: false,
          message: 'Blog post not found'
        });
      }
      
      return reply.send({
        success: true,
        message: 'Blog post deleted successfully'
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

  async getPostsByAuthor(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { authorId } = request.params as { authorId: string };
      const { limit } = request.query as { limit?: string };
      
      const posts = await blogService.getPostsByAuthor(
        authorId, 
        limit ? parseInt(limit, 10) : 10
      );
      
      return reply.send({
        success: true,
        data: posts
      });
    } catch (error) {
      return reply.status(500).send({
        success: false,
        message: 'Internal server error'
      });
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
      
      return reply.send({
        success: true,
        data: posts
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

  async getPopularPosts(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { limit } = request.query as { limit?: string };
      
      const posts = await blogService.getPopularPosts(
        limit ? parseInt(limit, 10) : 10
      );
      
      return reply.send({
        success: true,
        data: posts
      });
    } catch (error) {
      return reply.status(500).send({
        success: false,
        message: 'Internal server error'
      });
    }
  }

  async getRecentPosts(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { limit } = request.query as { limit?: string };
      
      const posts = await blogService.getRecentPosts(
        limit ? parseInt(limit, 10) : 10
      );
      
      return reply.send({
        success: true,
        data: posts
      });
    } catch (error) {
      return reply.status(500).send({
        success: false,
        message: 'Internal server error'
      });
    }
  }

  async searchPosts(request: FastifyRequest, reply: FastifyReply) {
    try {
      const { q } = request.query as { q: string };
      const filters = blogPostQuerySchema.parse(request.query);
      
      if (!q) {
        return reply.status(400).send({
          success: false,
          message: 'Search query is required'
        });
      }
      
      const result = await blogService.searchPosts(q, filters);
      
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
}
