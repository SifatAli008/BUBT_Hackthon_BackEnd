/**
 * Standard API response wrapper
 * Provides consistent response structure across all endpoints
 */
export class ApiResponse<T> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
  meta?: {
    page?: number;
    limit?: number;
    total?: number;
    totalPages?: number;
  };

  constructor(data?: T, message?: string, meta?: ApiResponse<T>['meta']) {
    this.success = true;
    this.data = data;
    this.message = message;
    this.meta = meta;
  }

  /**
   * Create success response
   */
  static success<T>(
    data: T,
    message?: string,
    meta?: ApiResponse<T>['meta'],
  ): ApiResponse<T> {
    return new ApiResponse(data, message, meta);
  }

  /**
   * Create error response
   */
  static error(message: string, error?: string): ApiResponse<null> {
    const response = new ApiResponse<null>();
    response.success = false;
    response.message = message;
    response.error = error;
    return response;
  }

  /**
   * Create paginated response
   */
  static paginated<T>(
    data: T[],
    total: number,
    page: number,
    limit: number,
    message?: string,
  ): ApiResponse<T[]> {
    const totalPages = Math.ceil(total / limit);
    return new ApiResponse(data, message, {
      page,
      limit,
      total,
      totalPages,
    });
  }
}
