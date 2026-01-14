import { Module, Global } from '@nestjs/common';
import { DatabaseService } from './database.service';
import { PrismaService } from './prisma.service';

/**
 * Global database module
 * Provides database services throughout the application
 * - DatabaseService: Raw SQL queries with pg
 * - PrismaService: Prisma ORM with Prisma Accelerate support
 */
@Global()
@Module({
  providers: [DatabaseService, PrismaService],
  exports: [DatabaseService, PrismaService],
})
export class DatabaseModule {}
