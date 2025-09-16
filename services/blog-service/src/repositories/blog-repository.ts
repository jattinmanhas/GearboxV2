import { eq, and, or, desc, asc, ilike, sql, count, ne } from 'drizzle-orm';
import { db } from '../config/database';
import { blogPosts, categories, type BlogPost, type Category } from '../config/schema';
import { BlogPostFilters, BlogPostListResponse } from '../types/blog';

export class BlogRepository {
  async create(data: Omit<BlogPost, 'id' | 'createdAt' | 'updatedAt'>): Promise<BlogPost> {
    const [post] = await db.insert(blogPosts).values(data).returning();
    return post;
  }

  async findById(id: string): Promise<BlogPost | null> {
    const [post] = await db
      .select()
      .from(blogPosts)
      .leftJoin(categories, eq(blogPosts.categoryId, categories.id))
      .where(eq(blogPosts.id, id))
      .limit(1);
    
    if (!post) return null;
    
    return {
      ...post.blog_posts,
      categoryName: post.categories?.name,
    } as BlogPost;
  }

  async findBySlug(slug: string): Promise<BlogPost | null> {
    const [post] = await db
      .select()
      .from(blogPosts)
      .leftJoin(categories, eq(blogPosts.categoryId, categories.id))
      .where(eq(blogPosts.slug, slug))
      .limit(1);
    
    if (!post) return null;
    
    return {
      ...post.blog_posts,
      categoryName: post.categories?.name,
    } as BlogPost;
  }

  async findMany(filters: BlogPostFilters = {}): Promise<BlogPostListResponse> {
    const {
      status,
      authorId,
      categoryId,
      tags,
      search,
      page = 1,
      limit = 10,
      sortBy = 'createdAt',
      sortOrder = 'desc'
    } = filters;

    const offset = (page - 1) * limit;

    // Build where conditions
    const conditions = [];
    
    if (status) {
      conditions.push(eq(blogPosts.status, status));
    }
    
    if (authorId) {
      conditions.push(eq(blogPosts.authorId, authorId));
    }
    
    if (categoryId) {
      conditions.push(eq(blogPosts.categoryId, categoryId));
    }
    
    if (tags && tags.length > 0) {
      conditions.push(sql`${blogPosts.tags} && ${tags}`);
    }
    
    if (search) {
      conditions.push(
        or(
          ilike(blogPosts.title, `%${search}%`),
          ilike(blogPosts.content, `%${search}%`),
          ilike(blogPosts.excerpt, `%${search}%`)
        )
      );
    }

    const whereClause = conditions.length > 0 ? and(...conditions) : undefined;

    // Build order by clause
    let orderBy;
    switch (sortBy) {
      case 'title':
        orderBy = sortOrder === 'asc' ? asc(blogPosts.title) : desc(blogPosts.title);
        break;
      case 'updatedAt':
        orderBy = sortOrder === 'asc' ? asc(blogPosts.updatedAt) : desc(blogPosts.updatedAt);
        break;
      case 'publishedAt':
        orderBy = sortOrder === 'asc' ? asc(blogPosts.publishedAt) : desc(blogPosts.publishedAt);
        break;
      case 'viewCount':
        orderBy = sortOrder === 'asc' ? asc(blogPosts.viewCount) : desc(blogPosts.viewCount);
        break;
      default:
        orderBy = sortOrder === 'asc' ? asc(blogPosts.createdAt) : desc(blogPosts.createdAt);
    }

    // Get total count
    const [{ count: total }] = await db
      .select({ count: count() })
      .from(blogPosts)
      .where(whereClause);

    // Get posts
    const posts = await db
      .select({
        blog_posts: blogPosts,
        categories: {
          name: categories.name,
        }
      })
      .from(blogPosts)
      .leftJoin(categories, eq(blogPosts.categoryId, categories.id))
      .where(whereClause)
      .orderBy(orderBy)
      .limit(limit)
      .offset(offset);

    const blogPostsList = posts.map(post => ({
      ...post.blog_posts,
      categoryName: post.categories?.name,
    })) as BlogPost[];

    return {
      posts: blogPostsList,
      total: Number(total),
      page,
      limit,
      totalPages: Math.ceil(Number(total) / limit),
    };
  }

  async update(id: string, data: Partial<BlogPost>): Promise<BlogPost | null> {
    const [post] = await db
      .update(blogPosts)
      .set({ ...data, updatedAt: new Date() })
      .where(eq(blogPosts.id, id))
      .returning();
    
    return post || null;
  }

  async delete(id: string): Promise<boolean> {
    const result = await db
      .delete(blogPosts)
      .where(eq(blogPosts.id, id));
    
    return result.length > 0;
  }

  async incrementViewCount(id: string): Promise<void> {
    await db
      .update(blogPosts)
      .set({ 
        viewCount: sql`${blogPosts.viewCount} + 1`,
        updatedAt: new Date()
      })
      .where(eq(blogPosts.id, id));
  }

  async findByAuthor(authorId: string, limit: number = 10): Promise<BlogPost[]> {
    const posts = await db
      .select()
      .from(blogPosts)
      .leftJoin(categories, eq(blogPosts.categoryId, categories.id))
      .where(eq(blogPosts.authorId, authorId))
      .orderBy(desc(blogPosts.createdAt))
      .limit(limit);

    return posts.map(post => ({
      ...post.blog_posts,
      categoryName: post.categories?.name,
    })) as BlogPost[];
  }

  async findRelatedPosts(
    postId: string,
    categoryId: string | null,
    tags: string[],
    limit: number = 5
  ): Promise<BlogPost[]> {
    const conditions = [ne(blogPosts.id, postId)]; // exclude current post
  
    if (categoryId) {
      conditions.push(eq(blogPosts.categoryId, categoryId));
    }
  
    if (tags.length > 0) {
      const arrayLiteral = `{${tags.join(",")}}`; // e.g. "{fastify,typescript,nodejs}"
      conditions.push(
        sql`${blogPosts.tags} && ${sql.raw(`'${arrayLiteral}'::text[]`)}`
      );
    }
  
    const posts = await db
      .select()
      .from(blogPosts)
      .leftJoin(categories, eq(blogPosts.categoryId, categories.id))
      .where(and(...conditions))
      .orderBy(desc(blogPosts.createdAt))
      .limit(limit);
  
    return posts.map(post => ({
      ...post.blog_posts,
      categoryName: post.categories?.name,
    })) as BlogPost[];
  }
}