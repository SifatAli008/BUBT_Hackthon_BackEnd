import {
  Injectable,
  OnModuleInit,
  OnModuleDestroy,
  Logger,
} from '@nestjs/common';
import { Pool, QueryResult, PoolClient } from 'pg';
import { ConfigService } from '../config/config.service';

/**
 * Database service using raw SQL queries with connection pooling
 * Optimized for production use with connection pooling and query optimization
 */
@Injectable()
export class DatabaseService implements OnModuleInit, OnModuleDestroy {
  private readonly logger = new Logger(DatabaseService.name);
  private pool: Pool;

  constructor(private configService: ConfigService) {}

  /**
   * Initialize database connection pool on module init
   */
  async onModuleInit() {
    const config = this.configService.getConfig();
    const dbConfig = config.database;

    // Use connection string if provided, otherwise build from individual config
    const connectionString =
      dbConfig.url ||
      `postgresql://${dbConfig.user}:${dbConfig.password}@${dbConfig.host}:${dbConfig.port}/${dbConfig.database}`;

    this.pool = new Pool({
      connectionString,
      max: dbConfig.maxPoolSize, // Maximum number of clients in the pool
      idleTimeoutMillis: dbConfig.idleTimeoutMs, // Close idle clients after 30 seconds
      connectionTimeoutMillis: 10000, // Return an error after 10 seconds if connection could not be established
      // SSL configuration for production (uncomment and configure as needed)
      // ssl: config.nodeEnv === 'production' ? { rejectUnauthorized: false } : false,
    });

    // Handle pool errors
    this.pool.on('error', (err) => {
      this.logger.error('Unexpected error on idle client', err);
    });

    // Test connection
    try {
      await this.pool.query('SELECT NOW()');
      this.logger.log('Database connection established successfully');
    } catch (error) {
      this.logger.error('Failed to connect to database', error);
      throw error;
    }
  }

  /**
   * Close database connection pool on module destroy
   */
  async onModuleDestroy() {
    if (this.pool) {
      await this.pool.end();
      this.logger.log('Database connection pool closed');
    }
  }

  /**
   * Execute a query with parameters (recommended - prevents SQL injection)
   * Uses parameterized queries for security and performance
   *
   * @param text SQL query string with $1, $2, etc. placeholders
   * @param params Array of parameter values
   * @returns Promise<QueryResult<T>>
   */
  async query<T extends Record<string, any> = any>(
    text: string,
    params?: any[],
  ): Promise<QueryResult<T>> {
    const start = Date.now();
    try {
      const result = await this.pool.query(text, params);
      const duration = Date.now() - start;

      // Log slow queries (>1000ms) for optimization
      if (duration > 1000) {
        this.logger.warn(
          `Slow query detected (${duration}ms): ${text.substring(0, 100)}`,
        );
      }

      return result;
    } catch (error) {
      this.logger.error(`Query error: ${text.substring(0, 100)}`, error);
      throw error;
    }
  }

  /**
   * Get a client from the pool for transactions
   * IMPORTANT: Always release the client after use with client.release()
   *
   * @returns Promise<PoolClient>
   */
  async getClient(): Promise<PoolClient> {
    return this.pool.connect();
  }

  /**
   * Execute a transaction
   * Automatically handles commit/rollback
   *
   * @param callback Function that receives a client and returns a promise
   * @returns Promise<T>
   */
  async transaction<T>(
    callback: (client: PoolClient) => Promise<T>,
  ): Promise<T> {
    const client = await this.getClient();
    try {
      await client.query('BEGIN');
      const result = await callback(client);
      await client.query('COMMIT');
      return result;
    } catch (error) {
      await client.query('ROLLBACK');
      throw error;
    } finally {
      client.release();
    }
  }

  /**
   * Execute a query and return a single row
   * Optimized for single-row queries
   *
   * @param text SQL query string
   * @param params Query parameters
   * @returns Promise<T | null>
   */
  async queryOne<T extends Record<string, any> = any>(
    text: string,
    params?: any[],
  ): Promise<T | null> {
    const result = await this.query<T>(text, params);
    return result.rows[0] || null;
  }

  /**
   * Execute a query and return all rows
   *
   * @param text SQL query string
   * @param params Query parameters
   * @returns Promise<T[]>
   */
  async queryMany<T extends Record<string, any> = any>(
    text: string,
    params?: any[],
  ): Promise<T[]> {
    const result = await this.query<T>(text, params);
    return result.rows;
  }

  /**
   * Execute a query that should return exactly one row
   * Throws error if zero or multiple rows returned
   *
   * @param text SQL query string
   * @param params Query parameters
   * @returns Promise<T>
   */
  async queryExactlyOne<T extends Record<string, any> = any>(
    text: string,
    params?: any[],
  ): Promise<T> {
    const result = await this.query<T>(text, params);
    if (result.rows.length === 0) {
      throw new Error('Expected exactly one row, but got zero');
    }
    if (result.rows.length > 1) {
      throw new Error(
        `Expected exactly one row, but got ${result.rows.length}`,
      );
    }
    return result.rows[0];
  }

  /**
   * Health check query
   * Used for health checks and connection testing
   *
   * @returns Promise<boolean>
   */
  async healthCheck(): Promise<boolean> {
    try {
      await this.pool.query('SELECT 1');
      return true;
    } catch {
      return false;
    }
  }
}
