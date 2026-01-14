import { Injectable, OnModuleInit, OnModuleDestroy, Logger } from '@nestjs/common';
import { PrismaClient } from '@prisma/client';
import { ConfigService } from '../config/config.service';

/**
 * Prisma Service
 * Provides Prisma Client with Prisma Accelerate support
 * Use this for Prisma ORM operations
 */
@Injectable()
export class PrismaService extends PrismaClient implements OnModuleInit, OnModuleDestroy {
  private readonly logger = new Logger(PrismaService.name);

  constructor(private configService: ConfigService) {
    const config = configService.getConfig();
    const databaseUrl = config.database.url;

    // Check if it's a Prisma Accelerate URL
    if (databaseUrl && databaseUrl.startsWith('prisma+postgres://')) {
      // For Prisma Accelerate, use the connection string directly
      // Prisma 7 supports prisma+postgres:// URLs
      super({
        datasources: {
          db: {
            url: databaseUrl,
          },
        },
      });
      this.logger.log('Initialized Prisma Client with Prisma Accelerate');
    } else if (databaseUrl) {
      // Regular PostgreSQL connection
      super({
        datasources: {
          db: {
            url: databaseUrl,
          },
        },
      });
      this.logger.log('Initialized Prisma Client with direct connection');
    } else {
      // No connection string provided
      super();
      this.logger.warn('No DATABASE_URL provided, Prisma Client may not work');
    }
  }

  async onModuleInit() {
    try {
      await this.$connect();
      this.logger.log('Prisma Client connected successfully');
    } catch (error) {
      this.logger.error('Failed to connect Prisma Client', error);
      throw error;
    }
  }

  async onModuleDestroy() {
    await this.$disconnect();
    this.logger.log('Prisma Client disconnected');
  }
}
