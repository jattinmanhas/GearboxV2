import { BlogPost, CreateBlogPostRequest, UpdateBlogPostRequest, BlogPostFilters, BlogPostListResponse } from '../types/blog';

export class BlogPostDomain {
  static createSlug(title: string): string {
    return title
      .toLowerCase()
      .replace(/[^a-z0-9\s-]/g, '') // Remove special characters
      .replace(/\s+/g, '-') // Replace spaces with hyphens
      .replace(/-+/g, '-') // Replace multiple hyphens with single hyphen
      .trim();
  }

  static calculateReadTime(content: string): number {
    const wordsPerMinute = 200;
    const wordCount = content.split(/\s+/).length;
    return Math.ceil(wordCount / wordsPerMinute);
  }

  static generateExcerpt(content: string, maxLength: number = 160): string {
    const plainText = content.replace(/<[^>]*>/g, ''); // Remove HTML tags
    if (plainText.length <= maxLength) {
      return plainText;
    }
    return plainText.substring(0, maxLength).replace(/\s+\S*$/, '') + '...';
  }

  static validateBlogPost(data: CreateBlogPostRequest): { isValid: boolean; errors: string[] } {
    const errors: string[] = [];

    if (!data.title || data.title.trim().length === 0) {
      errors.push('Title is required');
    } else if (data.title.length > 255) {
      errors.push('Title must be less than 255 characters');
    }

    if (!data.content || data.content.trim().length === 0) {
      errors.push('Content is required');
    }

    if (!data.authorId || data.authorId.trim().length === 0) {
      errors.push('Author ID is required');
    }

    if (!data.authorName || data.authorName.trim().length === 0) {
      errors.push('Author name is required');
    }

    if (!data.authorEmail || data.authorEmail.trim().length === 0) {
      errors.push('Author email is required');
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.authorEmail)) {
      errors.push('Author email must be valid');
    }

    if (data.status && !['draft', 'published'].includes(data.status)) {
      errors.push('Status must be either "draft" or "published"');
    }

    if (data.tags && !Array.isArray(data.tags)) {
      errors.push('Tags must be an array');
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  static prepareBlogPostForCreation(data: CreateBlogPostRequest): Omit<BlogPost, 'id' | 'createdAt' | 'updatedAt'> {
    const slug = this.createSlug(data.title);
    const excerpt = data.excerpt || this.generateExcerpt(data.content);
    const readTime = this.calculateReadTime(data.content);
    const publishedAt = data.status === 'published' ? new Date() : null;

    return {
      title: data.title.trim(),
      slug,
      content: data.content.trim(),
      excerpt,
      authorId: data.authorId,
      authorName: data.authorName.trim(),
      authorEmail: data.authorEmail.trim(),
      status: data.status || 'draft',
      featuredImage: data.featuredImage || null,
      tags: data.tags || null,
      categoryId: data.categoryId || null,
      publishedAt,
      viewCount: 0,
      readTime
    };
  }

  static prepareBlogPostForUpdate(data: UpdateBlogPostRequest, existingPost: BlogPost): Partial<BlogPost> {
    const updateData: Partial<BlogPost> = {};

    if (data.title !== undefined) {
      updateData.title = data.title.trim();
      updateData.slug = this.createSlug(data.title);
    }

    if (data.content !== undefined) {
      updateData.content = data.content.trim();
      updateData.readTime = this.calculateReadTime(data.content);
    }

    if (data.excerpt !== undefined) {
      updateData.excerpt = data.excerpt.trim();
    } else if (data.content !== undefined) {
      updateData.excerpt = this.generateExcerpt(data.content);
    }

    if (data.status !== undefined) {
      updateData.status = data.status;
      if (data.status === 'published' && existingPost.status !== 'published') {
        updateData.publishedAt = new Date();
      }
    }

    if (data.featuredImage !== undefined) {
      updateData.featuredImage = data.featuredImage || null;
    }

    if (data.tags !== undefined) {
      updateData.tags = data.tags || null;
    }

    if (data.categoryId !== undefined) {
      updateData.categoryId = data.categoryId || null;
    }

    return updateData;
  }
}
