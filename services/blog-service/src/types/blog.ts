export interface BlogPost {
  id: string;
  title: string;
  slug: string;
  content: string;
  excerpt: string | null;
  authorId: string;
  authorName: string;
  authorEmail: string;
  status: 'draft' | 'published' | 'archived';
  featuredImage: string | null;
  tags: string[] | null;
  categoryId: string | null;
  categoryName?: string;
  publishedAt: Date | null;
  createdAt: Date;
  updatedAt: Date;
  viewCount: number;
  readTime: number; // in minutes
}

export interface CreateBlogPostRequest {
  title: string;
  content: string;
  excerpt?: string;
  authorId: string;
  authorName: string;
  authorEmail: string;
  status?: 'draft' | 'published';
  featuredImage?: string;
  tags?: string[];
  categoryId?: string;
}

export interface UpdateBlogPostRequest {
  title?: string;
  content?: string;
  excerpt?: string;
  status?: 'draft' | 'published' | 'archived';
  featuredImage?: string;
  tags?: string[];
  categoryId?: string;
}

export interface BlogPostFilters {
  status?: 'draft' | 'published' | 'archived';
  authorId?: string;
  categoryId?: string;
  tags?: string[];
  search?: string;
  page?: number;
  limit?: number;
  sortBy?: 'createdAt' | 'updatedAt' | 'publishedAt' | 'viewCount' | 'title';
  sortOrder?: 'asc' | 'desc';
}

export interface BlogPostListResponse {
  posts: BlogPost[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  description: string | null;
  color: string | null;
  createdAt: Date;
  updatedAt: Date;
}

export interface CreateCategoryRequest {
  name: string;
  description?: string;
  color?: string;
}

export interface UpdateCategoryRequest {
  name?: string;
  description?: string;
  color?: string;
}
