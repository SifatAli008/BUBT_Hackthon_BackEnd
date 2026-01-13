import { Module, Global } from '@nestjs/common';
import { DatabaseService } from './database.service';

/**
 * Global database module
 * Provides database service throughout the application
 */
@Global()
@Module({
  providers: [DatabaseService],
  exports: [DatabaseService],
})
export class DatabaseModule {}
