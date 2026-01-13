import { IsOptional, IsInt, Min, Max } from 'class-validator';
import { Type } from 'class-transformer';

/**
 * Pagination DTO for query parameters
 * Used across all list endpoints
 */
export class PaginationDto {
  @IsOptional()
  @Type(() => Number)
  @IsInt()
  @Min(1)
  page?: number = 1;

  @IsOptional()
  @Type(() => Number)
  @IsInt()
  @Min(1)
  @Max(100)
  limit?: number = 10;

  /**
   * Get offset for SQL queries
   */
  getOffset(): number {
    return ((this.page || 1) - 1) * (this.limit || 10);
  }

  /**
   * Get limit for SQL queries
   */
  getLimit(): number {
    return this.limit || 10;
  }
}
