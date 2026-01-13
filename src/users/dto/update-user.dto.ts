import { ApiPropertyOptional } from '@nestjs/swagger';
import { IsString, IsOptional, MinLength, IsUUID } from 'class-validator';

/**
 * Update User DTO
 * Validates user profile update data
 */
export class UpdateUserDto {
  @ApiPropertyOptional({ example: 'John Doe', description: 'User full name' })
  @IsOptional()
  @IsString()
  @MinLength(2, { message: 'Name must be at least 2 characters long' })
  name?: string;

  @ApiPropertyOptional({ description: 'Household ID' })
  @IsOptional()
  @IsUUID()
  householdId?: string;
}
