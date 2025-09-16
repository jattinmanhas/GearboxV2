import { CategoryRepository } from '../repositories/category-repository';
import { CreateCategoryRequest, UpdateCategoryRequest, Category } from '../types/blog';

export class CategoryService {
  private categoryRepository: CategoryRepository;

  constructor() {
    this.categoryRepository = new CategoryRepository();
  }

  async createCategory(data: CreateCategoryRequest): Promise<Category> {
    // Validate required fields
    if (!data.name || data.name.trim().length === 0) {
      throw new Error('Category name is required');
    }

    if (data.name.length > 255) {
      throw new Error('Category name must be less than 255 characters');
    }

    // Check if category with same name already exists
    const existingCategory = await this.categoryRepository.findBySlug(
      this.createSlug(data.name)
    );
    if (existingCategory) {
      throw new Error('A category with this name already exists');
    }

    return await this.categoryRepository.create(data);
  }

  async getCategoryById(id: string): Promise<Category | null> {
    return await this.categoryRepository.findById(id);
  }

  async getCategoryBySlug(slug: string): Promise<Category | null> {
    return await this.categoryRepository.findBySlug(slug);
  }

  async getCategories(search?: string, page: number = 1, limit: number = 10): Promise<{
    categories: Category[];
    total: number;
    page: number;
    limit: number;
    totalPages: number;
  }> {
    return await this.categoryRepository.findMany(search, page, limit);
  }

  async getAllCategories(): Promise<Category[]> {
    return await this.categoryRepository.getAll();
  }

  async updateCategory(id: string, data: UpdateCategoryRequest): Promise<Category | null> {
    // Check if category exists
    const existingCategory = await this.categoryRepository.findById(id);
    if (!existingCategory) {
      throw new Error('Category not found');
    }

    // Check if new name creates a duplicate slug
    if (data.name && data.name !== existingCategory.name) {
      const newSlug = this.createSlug(data.name);
      const existingCategoryWithSlug = await this.categoryRepository.findBySlug(newSlug);
      if (existingCategoryWithSlug && existingCategoryWithSlug.id !== id) {
        throw new Error('A category with this name already exists');
      }
    }

    return await this.categoryRepository.update(id, data);
  }

  async deleteCategory(id: string): Promise<boolean> {
    const category = await this.categoryRepository.findById(id);
    if (!category) {
      throw new Error('Category not found');
    }

    return await this.categoryRepository.delete(id);
  }

  private createSlug(name: string): string {
    return name
      .toLowerCase()
      .replace(/[^a-z0-9\s-]/g, '') // Remove special characters
      .replace(/\s+/g, '-') // Replace spaces with hyphens
      .replace(/-+/g, '-') // Replace multiple hyphens with single hyphen
      .trim();
  }
}
