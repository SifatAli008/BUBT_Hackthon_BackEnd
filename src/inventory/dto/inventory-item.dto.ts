import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { IsString, IsNumber, IsOptional, IsUUID, Min } from 'class-validator';
import { Type } from 'class-transformer';

/**
 * Inventory Item DTO
 * Inventory item structure
 */
export class InventoryItemDto {
  @ApiProperty({ description: 'Inventory item ID' })
  id: string;

  @ApiProperty({ description: 'User ID' })
  userId: string;

  @ApiProperty({ description: 'Item name' })
  name: string;

  @ApiProperty({ description: 'Quantity' })
  quantity: number;

  @ApiPropertyOptional({ description: 'Unit of measurement' })
  unit?: string;

  @ApiPropertyOptional({ description: 'Expiry date' })
  expiryDate?: Date;

  @ApiPropertyOptional({ description: 'Category' })
  category?: string;

  @ApiPropertyOptional({ description: 'Storage location' })
  location?: string;

  @ApiPropertyOptional({ description: 'Food item reference ID' })
  foodItemId?: string;

  @ApiProperty({ description: 'Creation timestamp' })
  createdAt: Date;

  @ApiProperty({ description: 'Last update timestamp' })
  updatedAt: Date;
}

/**
 * Create Inventory Item DTO
 * Validates inventory item creation data
 */
export class CreateInventoryItemDto {
  @ApiProperty({ example: 'Milk', description: 'Item name' })
  @IsString()
  name: string;

  @ApiProperty({ example: 2.5, description: 'Quantity' })
  @IsNumber()
  @Min(0)
  @Type(() => Number)
  quantity: number;

  @ApiPropertyOptional({
    example: 'liters',
    description: 'Unit of measurement',
  })
  @IsOptional()
  @IsString()
  unit?: string;

  @ApiPropertyOptional({
    example: '2024-12-31T00:00:00Z',
    description: 'Expiry date',
  })
  @IsOptional()
  expiryDate?: Date;

  @ApiPropertyOptional({ example: 'Dairy', description: 'Category' })
  @IsOptional()
  @IsString()
  category?: string;

  @ApiPropertyOptional({
    example: 'Refrigerator',
    description: 'Storage location',
  })
  @IsOptional()
  @IsString()
  location?: string;

  @ApiPropertyOptional({ description: 'Food item reference ID' })
  @IsOptional()
  @IsUUID()
  foodItemId?: string;
}

/**
 * Update Inventory Item DTO
 * Validates inventory item update data
 */
export class UpdateInventoryItemDto {
  @ApiPropertyOptional({ example: 'Milk', description: 'Item name' })
  @IsOptional()
  @IsString()
  name?: string;

  @ApiPropertyOptional({ example: 2.5, description: 'Quantity' })
  @IsOptional()
  @IsNumber()
  @Min(0)
  @Type(() => Number)
  quantity?: number;

  @ApiPropertyOptional({
    example: 'liters',
    description: 'Unit of measurement',
  })
  @IsOptional()
  @IsString()
  unit?: string;

  @ApiPropertyOptional({
    example: '2024-12-31T00:00:00Z',
    description: 'Expiry date',
  })
  @IsOptional()
  expiryDate?: Date;

  @ApiPropertyOptional({ example: 'Dairy', description: 'Category' })
  @IsOptional()
  @IsString()
  category?: string;

  @ApiPropertyOptional({
    example: 'Refrigerator',
    description: 'Storage location',
  })
  @IsOptional()
  @IsString()
  location?: string;

  @ApiPropertyOptional({ description: 'Food item reference ID' })
  @IsOptional()
  @IsUUID()
  foodItemId?: string;
}

/**
 * Inventory Filter DTO
 * Filter options for inventory queries
 */
export class InventoryFilterDto {
  @ApiPropertyOptional({ description: 'Filter by category' })
  category?: string;

  @ApiPropertyOptional({ description: 'Filter by location' })
  location?: string;

  @ApiPropertyOptional({ description: 'Filter expiring items (days)' })
  expiringWithinDays?: number;

  @ApiPropertyOptional({ description: 'Filter expired items' })
  expired?: boolean;
}
