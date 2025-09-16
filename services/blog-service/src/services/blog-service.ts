import { BlogRepository } from '../repositories/blog-repository';
import { CategoryRepository } from '../repositories/category-repository';
import { BlogPostDomain } from '../domain/blog';
import { CreateBlogPostRequest, UpdateBlogPostRequest, BlogPostFilters, BlogPostListResponse, BlogPost } from '../types/blog';

export class BlogService {
  private blogRepository: BlogRepository;
  private categoryRepository: CategoryRepository;

  constructor() {
    this.blogRepository = new BlogRepository();
    this.categoryRepository = new CategoryRepository();
  }

  async createPost(data: CreateBlogPostRequest): Promise<BlogPost> {
    // Validate the blog post data
    const validation = BlogPostDomain.validateBlogPost(data);
    if (!validation.isValid) {
      throw new Error(`Validation failed: ${validation.errors.join(', ')}`);
    }

    // Check if category exists if provided
    if (data.categoryId) {
      const category = await this.categoryRepository.findById(data.categoryId);
      if (!category) {
        throw new Error('Category not found');
      }
    }

    // Check if slug already exists
    const slug = BlogPostDomain.createSlug(data.title);
    const existingPost = await this.blogRepository.findBySlug(slug);
    if (existingPost) {
      throw new Error('A blog post with this title already exists');
    }

    // Prepare the blog post data
    const blogPostData = BlogPostDomain.prepareBlogPostForCreation(data);

    // Create the blog post
    return await this.blogRepository.create(blogPostData);
  }

  async getPostById(id: string): Promise<BlogPost | null> {
    return await this.blogRepository.findById(id);
  }

  async getPostBySlug(slug: string): Promise<BlogPost | null> {
    return await this.blogRepository.findBySlug(slug);
  }

  async getPosts(filters: BlogPostFilters = {}): Promise<BlogPostListResponse> {
    return await this.blogRepository.findMany(filters);
  }

  async updatePost(id: string, data: UpdateBlogPostRequest): Promise<BlogPost | null> {
    // Check if post exists
    const existingPost = await this.blogRepository.findById(id);
    if (!existingPost) {
      throw new Error('Blog post not found');
    }

    // Validate category if provided
    if (data.categoryId) {
      const category = await this.categoryRepository.findById(data.categoryId);
      if (!category) {
        throw new Error('Category not found');
      }
    }

    // Check if new title creates a duplicate slug
    if (data.title && data.title !== existingPost.title) {
      const newSlug = BlogPostDomain.createSlug(data.title);
      const existingPostWithSlug = await this.blogRepository.findBySlug(newSlug);
      if (existingPostWithSlug && existingPostWithSlug.id !== id) {
        throw new Error('A blog post with this title already exists');
      }
    }

    // Prepare the update data
    const updateData = BlogPostDomain.prepareBlogPostForUpdate(data, existingPost);

    // Update the blog post
    return await this.blogRepository.update(id, updateData);
  }

  async deletePost(id: string): Promise<boolean> {
    const post = await this.blogRepository.findById(id);
    if (!post) {
      throw new Error('Blog post not found');
    }

    return await this.blogRepository.delete(id);
  }

  async incrementViewCount(id: string): Promise<void> {
    const post = await this.blogRepository.findById(id);
    if (!post) {
      throw new Error('Blog post not found');
    }

    await this.blogRepository.incrementViewCount(id);
  }

  async getPostsByAuthor(authorId: string, limit: number = 10): Promise<BlogPost[]> {
    return await this.blogRepository.findByAuthor(authorId, limit);
  }

  async getRelatedPosts(postId: string, limit: number = 5): Promise<BlogPost[]> {
    const post = await this.blogRepository.findById(postId);
    if (!post) {
      throw new Error('Blog post not found');
    }

    return await this.blogRepository.findRelatedPosts(
      postId,
      post.categoryId || null,
      post.tags || [],
      limit
    );
  }

  async getPopularPosts(limit: number = 10): Promise<BlogPost[]> {
    const result = await this.blogRepository.findMany({
      status: 'published',
      sortBy: 'viewCount',
      sortOrder: 'desc',
      limit,
    });

    return result.posts;
  }

  async getRecentPosts(limit: number = 10): Promise<BlogPost[]> {
    const result = await this.blogRepository.findMany({
      status: 'published',
      sortBy: 'publishedAt',
      sortOrder: 'desc',
      limit,
    });

    return result.posts;
  }

  async searchPosts(query: string, filters: Omit<BlogPostFilters, 'search'> = {}): Promise<BlogPostListResponse> {
    return await this.blogRepository.findMany({
      ...filters,
      search: query,
    });
  }
}
