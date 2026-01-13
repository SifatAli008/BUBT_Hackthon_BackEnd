import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

/**
 * Food Item DTO
 * Food item reference data structure
 */
export class FoodItemDto {
  @ApiProperty({ description: 'Food item ID' })
  id: string;

  @ApiProperty({ description: 'Food item name' })
  name: string;

  @ApiProperty({ description: 'Food category' })
  category: string;

  @ApiProperty({ description: 'Typical expiry days' })
  typicalExpiryDays: number;

  @ApiPropertyOptional({ description: 'Storage tips' })
  storageTips?: string;

  @ApiProperty({ description: 'Creation timestamp' })
  createdAt: Date;

  @ApiProperty({ description: 'Last update timestamp' })
  updatedAt: Date;
}

/**
 * Create Food Item DTO
 * Validates food item creation data
 */
export class CreateFoodItemDto {
  @ApiProperty({ example: 'Banana', description: 'Food item name' })
  name: string;

  @ApiProperty({ example: 'Fruits', description: 'Food category' })
  category: string;

  @ApiProperty({ example: 7, description: 'Typical expiry days' })
  typicalExpiryDays: number;

  @ApiPropertyOptional({
    example: 'Store at room temperature until ripe, then refrigerate',
    description: 'Storage tips',
  })
  storageTips?: string;
}

/**
 * Update Food Item DTO
 * Validates food item update data
 */
export class UpdateFoodItemDto {
  @ApiPropertyOptional({ example: 'Banana', description: 'Food item name' })
  name?: string;

  @ApiPropertyOptional({ example: 'Fruits', description: 'Food category' })
  category?: string;

  @ApiPropertyOptional({ example: 7, description: 'Typical expiry days' })
  typicalExpiryDays?: number;

  @ApiPropertyOptional({
    example: 'Store at room temperature until ripe, then refrigerate',
    description: 'Storage tips',
  })
  storageTips?: string;
}
