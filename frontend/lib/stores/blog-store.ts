import { create } from 'zustand';
import { BlogPost, Category, BlogPostFilters, CategoryFilters, CreateBlogPostRequest, UpdateBlogPostRequest, CreateCategoryRequest, UpdateCategoryRequest } from '../types/blog';
import { blogPostAPI, categoryAPI, BlogAPIError } from '../blog-api';

interface BlogState {
  // Blog posts
  posts: BlogPost[];
  currentPost: BlogPost | null;
  popularPosts: BlogPost[];
  recentPosts: BlogPost[];
  relatedPosts: BlogPost[];
  
  // Categories
  categories: Category[];
  currentCategory: Category | null;
  
  // Pagination and filters
  currentPage: number;
  totalPages: number;
  totalPosts: number;
  filters: BlogPostFilters;
  
  // Loading states
  isLoading: boolean;
  isCreating: boolean;
  isUpdating: boolean;
  isDeleting: boolean;
  
  // Error handling
  error: string | null;
  
  // Actions
  setPosts: (posts: BlogPost[]) => void;
  setCurrentPost: (post: BlogPost | null) => void;
  setPopularPosts: (posts: BlogPost[]) => void;
  setRecentPosts: (posts: BlogPost[]) => void;
  setRelatedPosts: (posts: BlogPost[]) => void;
  setCategories: (categories: Category[]) => void;
  setCurrentCategory: (category: Category | null) => void;
  setFilters: (filters: Partial<BlogPostFilters>) => void;
  setCurrentPage: (page: number) => void;
  setTotalPages: (pages: number) => void;
  setTotalPosts: (total: number) => void;
  setLoading: (loading: boolean) => void;
  setCreating: (creating: boolean) => void;
  setUpdating: (updating: boolean) => void;
  setDeleting: (deleting: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
  
  // Blog post actions
  fetchPosts: (filters?: BlogPostFilters) => Promise<void>;
  fetchPostById: (id: string) => Promise<void>;
  fetchPostBySlug: (slug: string) => Promise<void>;
  createPost: (data: CreateBlogPostRequest) => Promise<void>;
  updatePost: (id: string, data: UpdateBlogPostRequest) => Promise<void>;
  deletePost: (id: string) => Promise<void>;
  fetchPopularPosts: (limit?: number) => Promise<void>;
  fetchRecentPosts: (limit?: number) => Promise<void>;
  fetchRelatedPosts: (postId: string, limit?: number) => Promise<void>;
  searchPosts: (query: string, filters?: Partial<BlogPostFilters>) => Promise<void>;
  
  // Category actions
  fetchCategories: (filters?: CategoryFilters) => Promise<void>;
  fetchAllCategories: () => Promise<void>;
  fetchCategoryById: (id: string) => Promise<void>;
  fetchCategoryBySlug: (slug: string) => Promise<void>;
  createCategory: (data: CreateCategoryRequest) => Promise<void>;
  updateCategory: (id: string, data: UpdateCategoryRequest) => Promise<void>;
  deleteCategory: (id: string) => Promise<void>;
}

export const useBlogStore = create<BlogState>((set, get) => ({
  // Initial state
  posts: [],
  currentPost: null,
  popularPosts: [],
  recentPosts: [],
  relatedPosts: [],
  categories: [],
  currentCategory: null,
  currentPage: 1,
  totalPages: 0,
  totalPosts: 0,
  filters: {
    page: 1,
    limit: 10,
    sortBy: 'createdAt',
    sortOrder: 'desc',
  },
  isLoading: false,
  isCreating: false,
  isUpdating: false,
  isDeleting: false,
  error: null,

  // Setters
  setPosts: (posts) => set({ posts }),
  setCurrentPost: (post) => set({ currentPost: post }),
  setPopularPosts: (posts) => set({ popularPosts: posts }),
  setRecentPosts: (posts) => set({ recentPosts: posts }),
  setRelatedPosts: (posts) => set({ relatedPosts: posts }),
  setCategories: (categories) => set({ categories }),
  setCurrentCategory: (category) => set({ currentCategory: category }),
  setFilters: (filters) => set((state) => ({ 
    filters: { ...state.filters, ...filters } 
  })),
  setCurrentPage: (page) => set({ currentPage: page }),
  setTotalPages: (pages) => set({ totalPages: pages }),
  setTotalPosts: (total) => set({ totalPosts: total }),
  setLoading: (loading) => set({ isLoading: loading }),
  setCreating: (creating) => set({ isCreating: creating }),
  setUpdating: (updating) => set({ isUpdating: updating }),
  setDeleting: (deleting) => set({ isDeleting: deleting }),
  setError: (error) => set({ error }),
  clearError: () => set({ error: null }),

  // Blog post actions
  fetchPosts: async (filters) => {
    try {
      set({ isLoading: true, error: null });
      const response = await blogPostAPI.getPosts(filters || get().filters);
      set({
        posts: response.posts,
        totalPages: response.totalPages,
        totalPosts: response.total,
        currentPage: response.page,
        isLoading: false,
      });
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to fetch posts';
      set({ error: message, isLoading: false });
    }
  },

  fetchPostById: async (id) => {
    try {
      set({ isLoading: true, error: null });
      const post = await blogPostAPI.getPostById(id);
      set({ currentPost: post, isLoading: false });
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to fetch post';
      set({ error: message, isLoading: false });
    }
  },

  fetchPostBySlug: async (slug) => {
    try {
      set({ isLoading: true, error: null });
      const post = await blogPostAPI.getPostBySlug(slug);
      set({ currentPost: post, isLoading: false });
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to fetch post';
      set({ error: message, isLoading: false });
    }
  },

  createPost: async (data) => {
    try {
      set({ isCreating: true, error: null });
      const post = await blogPostAPI.createPost(data);
      set((state) => ({
        posts: [post, ...state.posts],
        isCreating: false,
      }));
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to create post';
      set({ error: message, isCreating: false });
    }
  },

  updatePost: async (id, data) => {
    try {
      set({ isUpdating: true, error: null });
      const updatedPost = await blogPostAPI.updatePost(id, data);
      set((state) => ({
        posts: state.posts.map(post => post.id === id ? updatedPost : post),
        currentPost: state.currentPost?.id === id ? updatedPost : state.currentPost,
        isUpdating: false,
      }));
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to update post';
      set({ error: message, isUpdating: false });
    }
  },

  deletePost: async (id) => {
    try {
      set({ isDeleting: true, error: null });
      await blogPostAPI.deletePost(id);
      set((state) => ({
        posts: state.posts.filter(post => post.id !== id),
        currentPost: state.currentPost?.id === id ? null : state.currentPost,
        isDeleting: false,
      }));
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to delete post';
      set({ error: message, isDeleting: false });
    }
  },

  fetchPopularPosts: async (limit = 10) => {
    try {
      set({ error: null });
      const posts = await blogPostAPI.getPopularPosts(limit);
      set({ popularPosts: posts });
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to fetch popular posts';
      set({ error: message });
    }
  },

  fetchRecentPosts: async (limit = 10) => {
    try {
      set({ error: null });
      const posts = await blogPostAPI.getRecentPosts(limit);
      set({ recentPosts: posts });
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to fetch recent posts';
      set({ error: message });
    }
  },

  fetchRelatedPosts: async (postId, limit = 5) => {
    try {
      set({ error: null });
      const posts = await blogPostAPI.getRelatedPosts(postId, limit);
      set({ relatedPosts: posts });
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to fetch related posts';
      set({ error: message });
    }
  },

  searchPosts: async (query, filters = {}) => {
    try {
      set({ isLoading: true, error: null });
      const response = await blogPostAPI.searchPosts(query, filters);
      set({
        posts: response.posts,
        totalPages: response.totalPages,
        totalPosts: response.total,
        currentPage: response.page,
        isLoading: false,
      });
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to search posts';
      set({ error: message, isLoading: false });
    }
  },

  // Category actions
  fetchCategories: async (filters) => {
    try {
      set({ isLoading: true, error: null });
      const response = await categoryAPI.getCategories(filters);
      set({
        categories: response.categories,
        isLoading: false,
      });
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to fetch categories';
      set({ error: message, isLoading: false });
    }
  },

  fetchAllCategories: async () => {
    try {
      set({ error: null });
      const categories = await categoryAPI.getAllCategories();
      set({ categories });
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to fetch categories';
      set({ error: message });
    }
  },

  fetchCategoryById: async (id) => {
    try {
      set({ isLoading: true, error: null });
      const category = await categoryAPI.getCategoryById(id);
      set({ currentCategory: category, isLoading: false });
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to fetch category';
      set({ error: message, isLoading: false });
    }
  },

  fetchCategoryBySlug: async (slug) => {
    try {
      set({ isLoading: true, error: null });
      const category = await categoryAPI.getCategoryBySlug(slug);
      set({ currentCategory: category, isLoading: false });
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to fetch category';
      set({ error: message, isLoading: false });
    }
  },

  createCategory: async (data) => {
    try {
      set({ isCreating: true, error: null });
      const category = await categoryAPI.createCategory(data);
      set((state) => ({
        categories: [...state.categories, category],
        isCreating: false,
      }));
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to create category';
      set({ error: message, isCreating: false });
    }
  },

  updateCategory: async (id, data) => {
    try {
      set({ isUpdating: true, error: null });
      const updatedCategory = await categoryAPI.updateCategory(id, data);
      set((state) => ({
        categories: state.categories.map(cat => cat.id === id ? updatedCategory : cat),
        currentCategory: state.currentCategory?.id === id ? updatedCategory : state.currentCategory,
        isUpdating: false,
      }));
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to update category';
      set({ error: message, isUpdating: false });
    }
  },

  deleteCategory: async (id) => {
    try {
      set({ isDeleting: true, error: null });
      await categoryAPI.deleteCategory(id);
      set((state) => ({
        categories: state.categories.filter(cat => cat.id !== id),
        currentCategory: state.currentCategory?.id === id ? null : state.currentCategory,
        isDeleting: false,
      }));
    } catch (error) {
      const message = error instanceof BlogAPIError ? error.message : 'Failed to delete category';
      set({ error: message, isDeleting: false });
    }
  },
}));
