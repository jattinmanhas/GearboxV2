import { z } from 'zod';

// Blog post validation schemas
export const createBlogPostSchema = z.object({
  title: z.string()
    .min(1, 'Title is required')
    .max(255, 'Title must be less than 255 characters')
    .trim(),
  content: z.string()
    .min(1, 'Content is required')
    .trim(),
  excerpt: z.string()
    .max(500, 'Excerpt must be less than 500 characters')
    .optional()
    .transform(val => val?.trim()),
  authorId: z.string()
    .min(1, 'Author ID is required')
    .trim(),
  authorName: z.string()
    .min(1, 'Author name is required')
    .max(255, 'Author name must be less than 255 characters')
    .trim(),
  authorEmail: z.string()
    .email('Author email must be valid')
    .max(255, 'Author email must be less than 255 characters')
    .trim(),
  status: z.enum(['draft', 'published'])
    .optional()
    .default('draft'),
  featuredImage: z.string()
    .url('Featured image must be a valid URL')
    .max(500, 'Featured image URL must be less than 500 characters')
    .optional(),
  tags: z.array(z.string().trim())
    .optional()
    .default([]),
  categoryId: z.string()
    .uuid('Category ID must be a valid UUID')
    .optional(),
});

export const updateBlogPostSchema = z.object({
  title: z.string()
    .min(1, 'Title is required')
    .max(255, 'Title must be less than 255 characters')
    .trim()
    .optional(),
  content: z.string()
    .min(1, 'Content is required')
    .trim()
    .optional(),
  excerpt: z.string()
    .max(500, 'Excerpt must be less than 500 characters')
    .optional()
    .transform(val => val?.trim()),
  status: z.enum(['draft', 'published', 'archived'])
    .optional(),
  featuredImage: z.string()
    .url('Featured image must be a valid URL')
    .max(500, 'Featured image URL must be less than 500 characters')
    .optional(),
  tags: z.array(z.string().trim())
    .optional(),
  categoryId: z.string()
    .uuid('Category ID must be a valid UUID')
    .optional(),
});

export const blogPostFiltersSchema = z.object({
  status: z.enum(['draft', 'published', 'archived'])
    .optional(),
  authorId: z.string()
    .trim()
    .optional(),
  categoryId: z.string()
    .uuid('Category ID must be a valid UUID')
    .optional(),
  tags: z.array(z.string().trim())
    .optional(),
  search: z.string()
    .trim()
    .optional(),
  page: z.number()
    .int('Page must be an integer')
    .min(1, 'Page must be at least 1')
    .optional()
    .default(1),
  limit: z.number()
    .int('Limit must be an integer')
    .min(1, 'Limit must be at least 1')
    .max(100, 'Limit must be at most 100')
    .optional()
    .default(10),
  sortBy: z.enum(['createdAt', 'updatedAt', 'publishedAt', 'viewCount', 'title'])
    .optional()
    .default('createdAt'),
  sortOrder: z.enum(['asc', 'desc'])
    .optional()
    .default('desc'),
});

// Category validation schemas
export const createCategorySchema = z.object({
  name: z.string()
    .min(1, 'Category name is required')
    .max(255, 'Category name must be less than 255 characters')
    .trim(),
  description: z.string()
    .max(1000, 'Description must be less than 1000 characters')
    .optional()
    .transform(val => val?.trim()),
  color: z.string()
    .regex(/^#[0-9A-Fa-f]{6}$/, 'Color must be a valid hex color code (e.g., #FF0000)')
    .optional(),
});

export const updateCategorySchema = z.object({
  name: z.string()
    .min(1, 'Category name is required')
    .max(255, 'Category name must be less than 255 characters')
    .trim()
    .optional(),
  description: z.string()
    .max(1000, 'Description must be less than 1000 characters')
    .optional()
    .transform(val => val?.trim()),
  color: z.string()
    .regex(/^#[0-9A-Fa-f]{6}$/, 'Color must be a valid hex color code (e.g., #FF0000)')
    .optional(),
});

export const categoryFiltersSchema = z.object({
  search: z.string()
    .trim()
    .optional(),
  page: z.number()
    .int('Page must be an integer')
    .min(1, 'Page must be at least 1')
    .optional()
    .default(1),
  limit: z.number()
    .int('Limit must be an integer')
    .min(1, 'Limit must be at least 1')
    .max(100, 'Limit must be at most 100')
    .optional()
    .default(10),
});

// Query parameter validation schemas
export const blogPostQuerySchema = z.object({
  status: z.string()
    .optional()
    .transform(val => val as 'draft' | 'published' | 'archived' | undefined),
  authorId: z.string()
    .optional(),
  categoryId: z.string()
    .optional(),
  tags: z.string()
    .optional()
    .transform(val => val ? val.split(',').map(tag => tag.trim()) : undefined),
  search: z.string()
    .optional(),
  page: z.string()
    .optional()
    .transform(val => val ? parseInt(val, 10) : undefined),
  limit: z.string()
    .optional()
    .transform(val => val ? parseInt(val, 10) : undefined),
  sortBy: z.string()
    .optional()
    .transform(val => val as 'createdAt' | 'updatedAt' | 'publishedAt' | 'viewCount' | 'title' | undefined),
  sortOrder: z.string()
    .optional()
    .transform(val => val as 'asc' | 'desc' | undefined),
});

export const categoryQuerySchema = z.object({
  search: z.string()
    .optional(),
  page: z.string()
    .optional()
    .transform(val => val ? parseInt(val, 10) : undefined),
  limit: z.string()
    .optional()
    .transform(val => val ? parseInt(val, 10) : undefined),
});
