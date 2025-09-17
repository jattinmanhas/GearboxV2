import { eq, ilike, desc, asc, sql, count } from 'drizzle-orm';
import { db } from '../config/database';
import { categories, type Category } from '../config/schema';
import { CreateCategoryRequest, UpdateCategoryRequest } from '../types/blog';
import { createSlug } from '../utils/slug';

export class CategoryRepository {
  async create(data: CreateCategoryRequest): Promise<Category> {
    const slug = createSlug(data.name);
    const [category] = await db.insert(categories).values({
      ...data,
      slug,
    }).returning();
    return category;
  }

  async findById(id: string): Promise<Category | null> {
    const [category] = await db
      .select()
      .from(categories)
      .where(eq(categories.id, id))
      .limit(1);
    
    return category || null;
  }

  async findBySlug(slug: string): Promise<Category | null> {
    const [category] = await db
      .select()
      .from(categories)
      .where(eq(categories.slug, slug))
      .limit(1);
    
    return category || null;
  }

  async findMany(search?: string, page: number = 1, limit: number = 10): Promise<{ categories: Category[]; total: number; page: number; limit: number; totalPages: number }> {
    const offset = (page - 1) * limit;

    const whereClause = search ? ilike(categories.name, `%${search}%`) : undefined;

    // Get total count
    const [{ count: total }] = await db
      .select({ count: count() })
      .from(categories)
      .where(whereClause);

    // Get categories
    const categoriesList = await db
      .select()
      .from(categories)
      .where(whereClause)
      .orderBy(asc(categories.name))
      .limit(limit)
      .offset(offset);

    return {
      categories: categoriesList,
      total: Number(total),
      page,
      limit,
      totalPages: Math.ceil(Number(total) / limit),
    };
  }

  async update(id: string, data: UpdateCategoryRequest): Promise<Category | null> {
    const updateData: any = { ...data };
    
    if (data.name) {
      updateData.slug = createSlug(data.name);
    }

    const [category] = await db
      .update(categories)
      .set({ ...updateData, updatedAt: new Date() })
      .where(eq(categories.id, id))
      .returning();
    
    return category || null;
  }

  async delete(id: string): Promise<boolean> {
    // First check if the category exists
    const existingCategory = await this.findById(id);
    
    if (!existingCategory) {
      return false;
    }
    
    // Perform the delete operation
    await db
      .delete(categories)
      .where(eq(categories.id, id));
    
    return true;
  }

  async getAll(): Promise<Category[]> {
    return await db
      .select()
      .from(categories)
      .orderBy(asc(categories.name));
  }

}
