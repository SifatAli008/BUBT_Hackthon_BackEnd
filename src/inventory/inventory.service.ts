import {
  Injectable,
  NotFoundException,
  ForbiddenException,
} from '@nestjs/common';
import { DatabaseService } from '../database/database.service';
import {
  CreateInventoryItemDto,
  UpdateInventoryItemDto,
  InventoryItemDto,
  InventoryFilterDto,
} from './dto/inventory-item.dto';

/**
 * Inventory Service
 * Handles inventory items database operations using raw SQL queries
 * Optimized queries with proper indexes for performance
 */
@Injectable()
export class InventoryService {
  constructor(private databaseService: DatabaseService) {}

  /**
   * Create a new inventory item
   * Uses parameterized query to prevent SQL injection
   */
  async create(
    userId: string,
    createDto: CreateInventoryItemDto,
  ): Promise<InventoryItemDto> {
    // Insert inventory item with optimized query
    const result = await this.databaseService.query<InventoryItemDto>(
      `
      INSERT INTO inventory_items (
        user_id, name, quantity, unit, expiry_date, 
        category, location, food_item_id
      )
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
      RETURNING 
        id,
        user_id as "userId",
        name,
        quantity,
        unit,
        expiry_date as "expiryDate",
        category,
        location,
        food_item_id as "foodItemId",
        created_at as "createdAt",
        updated_at as "updatedAt"
    `,
      [
        userId,
        createDto.name.trim(),
        createDto.quantity,
        createDto.unit?.trim() || null,
        createDto.expiryDate || null,
        createDto.category?.trim() || null,
        createDto.location?.trim() || null,
        createDto.foodItemId || null,
      ],
    );

    return result.rows[0];
  }

  /**
   * Find all inventory items for a user
   * Supports filtering, pagination, and sorting
   * Uses indexes on user_id, category, and expiry_date for performance
   */
  async findAll(
    userId: string,
    filter?: InventoryFilterDto,
    limit = 50,
    offset = 0,
  ): Promise<{
    items: InventoryItemDto[];
    total: number;
  }> {
    // Build WHERE clause dynamically
    const whereConditions: string[] = ['user_id = $1'];
    const params: any[] = [userId];
    let paramIndex = 2;

    if (filter?.category) {
      whereConditions.push(`category = $${paramIndex}`);
      params.push(filter.category.trim());
      paramIndex++;
    }

    if (filter?.location) {
      whereConditions.push(`location = $${paramIndex}`);
      params.push(filter.location.trim());
      paramIndex++;
    }

    if (filter?.expiringWithinDays !== undefined) {
      whereConditions.push(
        `expiry_date BETWEEN CURRENT_TIMESTAMP AND CURRENT_TIMESTAMP + ($${paramIndex}::text || ' days')::INTERVAL`,
      );
      params.push(filter.expiringWithinDays);
      paramIndex++;
    }

    if (filter?.expired === true) {
      whereConditions.push(`expiry_date < CURRENT_TIMESTAMP`);
    } else if (filter?.expired === false) {
      whereConditions.push(
        `(expiry_date IS NULL OR expiry_date >= CURRENT_TIMESTAMP)`,
      );
    }

    const whereClause =
      whereConditions.length > 0
        ? `WHERE ${whereConditions.join(' AND ')}`
        : '';

    // Get items with optimized query using indexes
    const query = `
      SELECT 
        id,
        user_id as "userId",
        name,
        quantity,
        unit,
        expiry_date as "expiryDate",
        category,
        location,
        food_item_id as "foodItemId",
        created_at as "createdAt",
        updated_at as "updatedAt"
      FROM inventory_items
      ${whereClause}
      ORDER BY expiry_date ASC NULLS LAST, created_at DESC
      LIMIT $${paramIndex} OFFSET $${paramIndex + 1}
    `;
    params.push(limit, offset);

    // Get total count
    const countQuery = `
      SELECT COUNT(*) as total
      FROM inventory_items
      ${whereClause}
    `;
    const countParams = params.slice(0, paramIndex - 2); // Exclude limit and offset

    const [itemsResult, countResult] = await Promise.all([
      this.databaseService.query<InventoryItemDto>(query, params),
      this.databaseService.query<{ total: string }>(countQuery, countParams),
    ]);

    return {
      items: itemsResult.rows,
      total: parseInt(countResult.rows[0].total, 10),
    };
  }

  /**
   * Find inventory item by ID
   * Verifies ownership before returning
   */
  async findOne(id: string, userId: string): Promise<InventoryItemDto> {
    const item = await this.databaseService.queryOne<InventoryItemDto>(
      `
      SELECT 
        id,
        user_id as "userId",
        name,
        quantity,
        unit,
        expiry_date as "expiryDate",
        category,
        location,
        food_item_id as "foodItemId",
        created_at as "createdAt",
        updated_at as "updatedAt"
      FROM inventory_items
      WHERE id = $1 AND user_id = $2
      LIMIT 1
    `,
      [id, userId],
    );

    if (!item) {
      throw new NotFoundException('Inventory item not found');
    }

    return item;
  }

  /**
   * Update inventory item
   * Verifies ownership before updating
   */
  async update(
    id: string,
    userId: string,
    updateDto: UpdateInventoryItemDto,
  ): Promise<InventoryItemDto> {
    // Build dynamic update query
    const updates: string[] = [];
    const params: any[] = [];
    let paramIndex = 1;

    if (updateDto.name !== undefined) {
      updates.push(`name = $${paramIndex}`);
      params.push(updateDto.name.trim());
      paramIndex++;
    }

    if (updateDto.quantity !== undefined) {
      updates.push(`quantity = $${paramIndex}`);
      params.push(updateDto.quantity);
      paramIndex++;
    }

    if (updateDto.unit !== undefined) {
      updates.push(`unit = $${paramIndex}`);
      params.push(updateDto.unit?.trim() || null);
      paramIndex++;
    }

    if (updateDto.expiryDate !== undefined) {
      updates.push(`expiry_date = $${paramIndex}`);
      params.push(updateDto.expiryDate || null);
      paramIndex++;
    }

    if (updateDto.category !== undefined) {
      updates.push(`category = $${paramIndex}`);
      params.push(updateDto.category?.trim() || null);
      paramIndex++;
    }

    if (updateDto.location !== undefined) {
      updates.push(`location = $${paramIndex}`);
      params.push(updateDto.location?.trim() || null);
      paramIndex++;
    }

    if (updateDto.foodItemId !== undefined) {
      updates.push(`food_item_id = $${paramIndex}`);
      params.push(updateDto.foodItemId || null);
      paramIndex++;
    }

    if (updates.length === 0) {
      return this.findOne(id, userId); // No updates, return current item
    }

    // Add updated_at timestamp
    updates.push(`updated_at = CURRENT_TIMESTAMP`);

    // Add id and user_id to params
    params.push(id, userId);

    const result = await this.databaseService.query<InventoryItemDto>(
      `
      UPDATE inventory_items
      SET ${updates.join(', ')}
      WHERE id = $${paramIndex} AND user_id = $${paramIndex + 1}
      RETURNING 
        id,
        user_id as "userId",
        name,
        quantity,
        unit,
        expiry_date as "expiryDate",
        category,
        location,
        food_item_id as "foodItemId",
        created_at as "createdAt",
        updated_at as "updatedAt"
    `,
      params,
    );

    if (result.rowCount === 0) {
      throw new NotFoundException('Inventory item not found');
    }

    return result.rows[0];
  }

  /**
   * Delete inventory item
   * Verifies ownership before deleting
   */
  async remove(id: string, userId: string): Promise<void> {
    const result = await this.databaseService.query(
      `
      DELETE FROM inventory_items
      WHERE id = $1 AND user_id = $2
    `,
      [id, userId],
    );

    if (result.rowCount === 0) {
      throw new NotFoundException('Inventory item not found');
    }
  }

  /**
   * Get expiring items
   * Returns items expiring within specified days
   * Uses index on expiry_date for fast query
   */
  async getExpiringItems(
    userId: string,
    days: number = 7,
  ): Promise<InventoryItemDto[]> {
    const items = await this.databaseService.queryMany<InventoryItemDto>(
      `
      SELECT 
        id,
        user_id as "userId",
        name,
        quantity,
        unit,
        expiry_date as "expiryDate",
        category,
        location,
        food_item_id as "foodItemId",
        created_at as "createdAt",
        updated_at as "updatedAt"
      FROM inventory_items
      WHERE user_id = $1
        AND expiry_date IS NOT NULL
        AND expiry_date BETWEEN CURRENT_TIMESTAMP 
          AND CURRENT_TIMESTAMP + ($2 || ' days')::INTERVAL
      ORDER BY expiry_date ASC
      LIMIT 50
    `,
      [userId, days.toString()],
    );

    return items;
  }

  /**
   * Get expired items
   * Returns items that have expired
   * Uses index on expiry_date for fast query
   */
  async getExpiredItems(userId: string): Promise<InventoryItemDto[]> {
    const items = await this.databaseService.queryMany<InventoryItemDto>(
      `
      SELECT 
        id,
        user_id as "userId",
        name,
        quantity,
        unit,
        expiry_date as "expiryDate",
        category,
        location,
        food_item_id as "foodItemId",
        created_at as "createdAt",
        updated_at as "updatedAt"
      FROM inventory_items
      WHERE user_id = $1
        AND expiry_date IS NOT NULL
        AND expiry_date < CURRENT_TIMESTAMP
      ORDER BY expiry_date DESC
      LIMIT 50
    `,
      [userId],
    );

    return items;
  }

  /**
   * Search inventory items by name
   * Uses ILIKE for case-insensitive search
   */
  async searchByName(
    userId: string,
    searchTerm: string,
    limit = 20,
  ): Promise<InventoryItemDto[]> {
    const items = await this.databaseService.queryMany<InventoryItemDto>(
      `
      SELECT 
        id,
        user_id as "userId",
        name,
        quantity,
        unit,
        expiry_date as "expiryDate",
        category,
        location,
        food_item_id as "foodItemId",
        created_at as "createdAt",
        updated_at as "updatedAt"
      FROM inventory_items
      WHERE user_id = $1
        AND name ILIKE $2
      ORDER BY expiry_date ASC NULLS LAST, name ASC
      LIMIT $3
    `,
      [userId, `%${searchTerm.trim()}%`, limit],
    );

    return items;
  }

  /**
   * Get inventory statistics
   * Returns summary statistics for user's inventory
   */
  async getStatistics(userId: string): Promise<{
    totalItems: number;
    totalCategories: number;
    expiringItems: number;
    expiredItems: number;
  }> {
    const result = await this.databaseService.queryOne<{
      total_items: string;
      total_categories: string;
      expiring_items: string;
      expired_items: string;
    }>(
      `
      SELECT 
        COUNT(*) as total_items,
        COUNT(DISTINCT category) as total_categories,
        COUNT(CASE 
          WHEN expiry_date IS NOT NULL 
            AND expiry_date BETWEEN CURRENT_TIMESTAMP 
              AND CURRENT_TIMESTAMP + INTERVAL '7 days' 
          THEN 1 END) as expiring_items,
        COUNT(CASE 
          WHEN expiry_date IS NOT NULL 
            AND expiry_date < CURRENT_TIMESTAMP 
          THEN 1 END) as expired_items
      FROM inventory_items
      WHERE user_id = $1
    `,
      [userId],
    );

    if (!result) {
      return {
        totalItems: 0,
        totalCategories: 0,
        expiringItems: 0,
        expiredItems: 0,
      };
    }

    return {
      totalItems: parseInt(result.total_items, 10),
      totalCategories: parseInt(result.total_categories, 10),
      expiringItems: parseInt(result.expiring_items, 10),
      expiredItems: parseInt(result.expired_items, 10),
    };
  }
}
