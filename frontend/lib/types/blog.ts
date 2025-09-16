// Blog post types
export interface BlogPost {
  id: string;
  title: string;
  slug: string;
  content: string;
  excerpt: string;
  authorId: string;
  authorName: string;
  authorEmail: string;
  status: 'draft' | 'published' | 'archived';
  featuredImage?: string;
  tags?: string[];
  categoryId?: string;
  categoryName?: string;
  publishedAt?: string;
  viewCount: number;
  readTime: number;
  createdAt: string;
  updatedAt: string;
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

// Category types
export interface Category {
  id: string;
  name: string;
  slug: string;
  description?: string;
  color?: string;
  createdAt: string;
  updatedAt: string;
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

export interface CategoryFilters {
  search?: string;
  page?: number;
  limit?: number;
}

export interface CategoryListResponse {
  categories: Category[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

// API Response types
export interface APIResponse<T = unknown> {
  timestamp: string;
  status: number;
  success: boolean;
  message?: string;
  data?: T;
  error?: unknown;
}
