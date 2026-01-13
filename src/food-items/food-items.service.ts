import {
  Injectable,
  NotFoundException,
  ConflictException,
} from '@nestjs/common';
import { DatabaseService } from '../database/database.service';
import { CreateFoodItemDto } from './dto/food-item.dto';
import { UpdateFoodItemDto } from './dto/food-item.dto';
import { FoodItemDto } from './dto/food-item.dto';

/**
 * Food Items Service
 * Handles food items reference data operations using raw SQL queries
 * Optimized queries for reference data management
 */
@Injectable()
export class FoodItemsService {
  constructor(private databaseService: DatabaseService) {}

  /**
   * Create a new food item
   * Uses parameterized query to prevent SQL injection
   */
  async create(createDto: CreateFoodItemDto): Promise<FoodItemDto> {
    // Check if food item with same name and category already exists
    const existing = await this.databaseService.queryOne<{ id: string }>(
      `
      SELECT id
      FROM food_items
      WHERE LOWER(name) = LOWER($1) AND LOWER(category) = LOWER($2)
      LIMIT 1
    `,
      [createDto.name.trim(), createDto.category.trim()],
    );

    if (existing) {
      throw new ConflictException(
        'Food item with this name and category already exists',
      );
    }

    // Create food item with optimized insert query
    const result = await this.databaseService.query<FoodItemDto>(
      `
      INSERT INTO food_items (name, category, typical_expiry_days, storage_tips)
      VALUES ($1, $2, $3, $4)
      RETURNING 
        id,
        name,
        category,
        typical_expiry_days as "typicalExpiryDays",
        storage_tips as "storageTips",
        created_at as "createdAt",
        updated_at as "updatedAt"
    `,
      [
        createDto.name.trim(),
        createDto.category.trim(),
        createDto.typicalExpiryDays,
        createDto.storageTips?.trim() || null,
      ],
    );

    return result.rows[0];
  }

  /**
   * Find all food items with optional filtering
   * Supports pagination and category filtering
   */
  async findAll(
    category?: string,
    limit = 100,
    offset = 0,
  ): Promise<{
    items: FoodItemDto[];
    total: number;
  }> {
    // Build query dynamically based on filters
    let query = `
      SELECT 
        id,
        name,
        category,
        typical_expiry_days as "typicalExpiryDays",
        storage_tips as "storageTips",
        created_at as "createdAt",
        updated_at as "updatedAt"
      FROM food_items
    `;
    const params: any[] = [];
    let paramIndex = 1;

    if (category) {
      query += ` WHERE LOWER(category) = LOWER($${paramIndex})`;
      params.push(category.trim());
      paramIndex++;
    }

    query += ` ORDER BY name ASC LIMIT $${paramIndex} OFFSET $${paramIndex + 1}`;
    params.push(limit, offset);

    // Get total count
    let countQuery = `SELECT COUNT(*) as total FROM food_items`;
    const countParams: any[] = [];

    if (category) {
      countQuery += ` WHERE LOWER(category) = LOWER($1)`;
      countParams.push(category.trim());
    }

    const [itemsResult, countResult] = await Promise.all([
      this.databaseService.query<FoodItemDto>(query, params),
      this.databaseService.query<{ total: string }>(countQuery, countParams),
    ]);

    return {
      items: itemsResult.rows,
      total: parseInt(countResult.rows[0].total, 10),
    };
  }

  /**
   * Find food item by ID
   * Uses primary key for fast lookup
   */
  async findOne(id: string): Promise<FoodItemDto> {
    const item = await this.databaseService.queryOne<FoodItemDto>(
      `
      SELECT 
        id,
        name,
        category,
        typical_expiry_days as "typicalExpiryDays",
        storage_tips as "storageTips",
        created_at as "createdAt",
        updated_at as "updatedAt"
      FROM food_items
      WHERE id = $1
      LIMIT 1
    `,
      [id],
    );

    if (!item) {
      throw new NotFoundException('Food item not found');
    }

    return item;
  }

  /**
   * Search food items by name
   * Uses ILIKE for case-insensitive search (optimized with index)
   */
  async searchByName(searchTerm: string, limit = 20): Promise<FoodItemDto[]> {
    const items = await this.databaseService.queryMany<FoodItemDto>(
      `
      SELECT 
        id,
        name,
        category,
        typical_expiry_days as "typicalExpiryDays",
        storage_tips as "storageTips",
        created_at as "createdAt",
        updated_at as "updatedAt"
      FROM food_items
      WHERE name ILIKE $1
      ORDER BY name ASC
      LIMIT $2
    `,
      [`%${searchTerm.trim()}%`, limit],
    );

    return items;
  }

  /**
   * Update food item
   * Uses parameterized query for security
   */
  async update(id: string, updateDto: UpdateFoodItemDto): Promise<FoodItemDto> {
    // Build dynamic update query
    const updates: string[] = [];
    const params: any[] = [];
    let paramIndex = 1;

    if (updateDto.name !== undefined) {
      updates.push(`name = $${paramIndex}`);
      params.push(updateDto.name.trim());
      paramIndex++;
    }

    if (updateDto.category !== undefined) {
      updates.push(`category = $${paramIndex}`);
      params.push(updateDto.category.trim());
      paramIndex++;
    }

    if (updateDto.typicalExpiryDays !== undefined) {
      updates.push(`typical_expiry_days = $${paramIndex}`);
      params.push(updateDto.typicalExpiryDays);
      paramIndex++;
    }

    if (updateDto.storageTips !== undefined) {
      updates.push(`storage_tips = $${paramIndex}`);
      params.push(updateDto.storageTips?.trim() || null);
      paramIndex++;
    }

    if (updates.length === 0) {
      return this.findOne(id); // No updates, return current item
    }

    // Add updated_at timestamp
    updates.push(`updated_at = CURRENT_TIMESTAMP`);

    // Add id to params
    params.push(id);

    const result = await this.databaseService.query<FoodItemDto>(
      `
      UPDATE food_items
      SET ${updates.join(', ')}
      WHERE id = $${paramIndex}
      RETURNING 
        id,
        name,
        category,
        typical_expiry_days as "typicalExpiryDays",
        storage_tips as "storageTips",
        created_at as "createdAt",
        updated_at as "updatedAt"
    `,
      params,
    );

    if (result.rowCount === 0) {
      throw new NotFoundException('Food item not found');
    }

    return result.rows[0];
  }

  /**
   * Delete food item
   * Uses parameterized query
   */
  async remove(id: string): Promise<void> {
    const result = await this.databaseService.query(
      `
      DELETE FROM food_items
      WHERE id = $1
    `,
      [id],
    );

    if (result.rowCount === 0) {
      throw new NotFoundException('Food item not found');
    }
  }

  /**
   * Get all categories
   * Returns distinct categories for filtering
   */
  async getCategories(): Promise<string[]> {
    const result = await this.databaseService.queryMany<{ category: string }>(
      `
      SELECT DISTINCT category
      FROM food_items
      ORDER BY category ASC
    `,
    );

    return result.map((row) => row.category);
  }
}
