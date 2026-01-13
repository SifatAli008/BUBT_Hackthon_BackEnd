import { Injectable, Logger } from '@nestjs/common';
import { DatabaseService } from './database/database.service';

/**
 * Root application service
 * Handles application-level business logic
 */
@Injectable()
export class AppService {
  private readonly logger = new Logger(AppService.name);

  constructor(private readonly databaseService: DatabaseService) {}

  /**
   * Get welcome message
   */
  getHello(): string {
    return 'FoodLink API is running!';
  }

  /**
   * Health check for the application
   * Checks database connectivity
   */
  async healthCheck(): Promise<boolean> {
    try {
      // Check database connection
      const dbHealthy = await this.databaseService.healthCheck();

      if (!dbHealthy) {
        this.logger.warn('Database health check failed');
        return false;
      }

      return true;
    } catch (error) {
      this.logger.error('Health check failed', error);
      return false;
    }
  }
}
