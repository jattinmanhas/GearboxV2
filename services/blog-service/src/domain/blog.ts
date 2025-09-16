import { BlogPost, CreateBlogPostRequest, UpdateBlogPostRequest, BlogPostFilters, BlogPostListResponse } from '../types/blog';
import { createSlug } from '../utils/slug';

export class BlogPostDomain {

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


  static prepareBlogPostForCreation(data: CreateBlogPostRequest): Omit<BlogPost, 'id' | 'createdAt' | 'updatedAt'> {
    const slug = createSlug(data.title);
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
      updateData.slug = createSlug(data.title);
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
