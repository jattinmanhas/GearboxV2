import { 
  BlogPost, 
  CreateBlogPostRequest, 
  UpdateBlogPostRequest, 
  BlogPostFilters, 
  BlogPostListResponse,
  Category,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  CategoryFilters,
  CategoryListResponse,
  APIResponse
} from './types/blog';

const BLOG_API_BASE = '/api/v1/blog';

class BlogAPIError extends Error {
  constructor(
    message: string,
    public status: number,
    public response?: APIResponse
  ) {
    super(message);
    this.name = 'BlogAPIError';
  }
}

async function handleResponse<T>(response: Response): Promise<T> {
  // Check if response has content before trying to parse JSON
  const contentType = response.headers.get('content-type');
  const hasJsonContent = contentType && contentType.includes('application/json');
  
  let data: APIResponse<T> | null = null;
  
  if (hasJsonContent) {
    try {
      const responseText = await response.text();
      if (responseText.trim()) {
        data = JSON.parse(responseText);
      }
    } catch (error) {
      console.error('Failed to parse JSON response:', error);
      console.error('Response text:', await response.text());
      data = null;
    }
  }
  
  if (!response.ok) {
    console.error('API Error Response:', {
      status: response.status,
      statusText: response.statusText,
      contentType,
      data
    });
    
    const errorMessage = data?.message || `HTTP ${response.status}: ${response.statusText}`;
    throw new BlogAPIError(
      errorMessage,
      response.status,
      data || undefined
    );
  }
  
  // For successful responses with no content (like DELETE operations)
  if (!data) {
    return undefined as T;
  }
  
  // Handle both wrapped responses and direct data
  if (data.data !== undefined) {
    return data.data as T;
  }
  
  // If no data field, return the entire response (for cases where data is returned directly)
  return data as T;
}

async function getAuthHeaders(): Promise<HeadersInit> {
  return {
    'Content-Type': 'application/json',
  };
}

// Blog Post API
export const blogPostAPI = {
  // Create a new blog post
  async createPost(data: CreateBlogPostRequest): Promise<BlogPost> {
    const response = await fetch(`${BLOG_API_BASE}/posts`, {
      method: 'POST',
      headers: await getAuthHeaders(),
      credentials: 'include', // Include cookies in the request
      body: JSON.stringify(data),
    });
    
    return handleResponse<BlogPost>(response);
  },

  // Get all blog posts with filters
  async getPosts(filters: BlogPostFilters = {}): Promise<BlogPostListResponse> {
    const params = new URLSearchParams();
    
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        if (Array.isArray(value)) {
          params.append(key, value.join(','));
        } else {
          params.append(key, String(value));
        }
      }
    });

    const response = await fetch(`${BLOG_API_BASE}/posts?${params.toString()}`);
    return handleResponse<BlogPostListResponse>(response);
  },

  // Get blog post by ID
  async getPostById(id: string): Promise<BlogPost> {
    const response = await fetch(`${BLOG_API_BASE}/posts/${id}`);
    return handleResponse<BlogPost>(response);
  },

  // Get blog post by slug
  async getPostBySlug(slug: string): Promise<BlogPost> {
    const response = await fetch(`${BLOG_API_BASE}/posts/slug/${slug}`);
    return handleResponse<BlogPost>(response);
  },

  // Update blog post
  async updatePost(id: string, data: UpdateBlogPostRequest): Promise<BlogPost> {
    const response = await fetch(`${BLOG_API_BASE}/posts/${id}`, {
      method: 'PUT',
      headers: await getAuthHeaders(),
      credentials: 'include', // Include cookies in the request
      body: JSON.stringify(data),
    });
    
    return handleResponse<BlogPost>(response);
  },

  // Delete blog post
  async deletePost(id: string): Promise<void> {
    const response = await fetch(`${BLOG_API_BASE}/posts/${id}`, {
      method: 'DELETE',
      headers: await getAuthHeaders(),
      credentials: 'include', // Include cookies in the request
    });
    
    await handleResponse<void>(response);
  },

  // Get posts by author
  async getPostsByAuthor(authorId: string, limit: number = 10): Promise<BlogPost[]> {
    const response = await fetch(`${BLOG_API_BASE}/posts/author/${authorId}?limit=${limit}`);
    return handleResponse<BlogPost[]>(response);
  },

  // Get related posts
  async getRelatedPosts(postId: string, limit: number = 5): Promise<BlogPost[]> {
    const response = await fetch(`${BLOG_API_BASE}/posts/${postId}/related?limit=${limit}`);
    return handleResponse<BlogPost[]>(response);
  },

  // Get popular posts
  async getPopularPosts(limit: number = 10): Promise<BlogPost[]> {
    const response = await fetch(`${BLOG_API_BASE}/posts/popular?limit=${limit}`);
    return handleResponse<BlogPost[]>(response);
  },

  // Get recent posts
  async getRecentPosts(limit: number = 10): Promise<BlogPost[]> {
    const response = await fetch(`${BLOG_API_BASE}/posts/recent?limit=${limit}`);
    return handleResponse<BlogPost[]>(response);
  },

  // Search posts
  async searchPosts(query: string, filters: Omit<BlogPostFilters, 'search'> = {}): Promise<BlogPostListResponse> {
    const params = new URLSearchParams({ q: query });
    
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        if (Array.isArray(value)) {
          params.append(key, value.join(','));
        } else {
          params.append(key, String(value));
        }
      }
    });

    const response = await fetch(`${BLOG_API_BASE}/posts/search?${params.toString()}`);
    return handleResponse<BlogPostListResponse>(response);
  },
};

// Category API
export const categoryAPI = {
  // Create a new category
  async createCategory(data: CreateCategoryRequest): Promise<Category> {
    const response = await fetch(`${BLOG_API_BASE}/categories`, {
      method: 'POST',
      headers: await getAuthHeaders(),
      credentials: 'include', // Include cookies in the request
      body: JSON.stringify(data),
    });
    
    return handleResponse<Category>(response);
  },

  // Get all categories with filters
  async getCategories(filters: CategoryFilters = {}): Promise<CategoryListResponse> {
    const params = new URLSearchParams();
    
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        params.append(key, String(value));
      }
    });

    const response = await fetch(`${BLOG_API_BASE}/categories?${params.toString()}`);
    return handleResponse<CategoryListResponse>(response);
  },

  // Get all categories (simple list)
  async getAllCategories(): Promise<Category[]> {
    const response = await fetch(`${BLOG_API_BASE}/categories/all`);
    return handleResponse<Category[]>(response);
  },

  // Get category by ID
  async getCategoryById(id: string): Promise<Category> {
    const response = await fetch(`${BLOG_API_BASE}/categories/${id}`);
    return handleResponse<Category>(response);
  },

  // Get category by slug
  async getCategoryBySlug(slug: string): Promise<Category> {
    const response = await fetch(`${BLOG_API_BASE}/categories/slug/${slug}`);
    return handleResponse<Category>(response);
  },

  // Update category
  async updateCategory(id: string, data: UpdateCategoryRequest): Promise<Category> {
    const response = await fetch(`${BLOG_API_BASE}/categories/${id}`, {
      method: 'PUT',
      headers: await getAuthHeaders(),
      credentials: 'include', // Include cookies in the request
      body: JSON.stringify(data),
    });
    
    return handleResponse<Category>(response);
  },

  // Delete category
  async deleteCategory(id: string): Promise<void> {
    const response = await fetch(`${BLOG_API_BASE}/categories/${id}`, {
      method: 'DELETE',
      headers: await getAuthHeaders(),
      credentials: 'include', // Include cookies in the request
    });
    
    await handleResponse<void>(response);
  },
};

export { BlogAPIError };
