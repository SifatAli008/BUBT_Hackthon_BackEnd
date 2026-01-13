import { Injectable } from '@nestjs/common';
import { ConfigService as NestConfigService } from '@nestjs/config';
import { IConfig } from './config.interface';

/**
 * Configuration service to manage environment variables
 * Provides type-safe access to configuration values
 */
@Injectable()
export class ConfigService {
  constructor(private nestConfigService: NestConfigService) {}

  /**
   * Get application configuration
   */
  getConfig(): IConfig {
    return {
      port: parseInt(this.nestConfigService.get<string>('PORT', '3000'), 10),
      nodeEnv: this.nestConfigService.get<string>('NODE_ENV', 'development'),
      database: {
        url: this.nestConfigService.get<string>('DATABASE_URL', ''),
        host: this.nestConfigService.get<string>('DB_HOST', 'localhost'),
        port: parseInt(
          this.nestConfigService.get<string>('DB_PORT', '5432'),
          10,
        ),
        user: this.nestConfigService.get<string>('DB_USER', ''),
        password: this.nestConfigService.get<string>('DB_PASSWORD', ''),
        database: this.nestConfigService.get<string>('DB_NAME', ''),
        maxPoolSize: parseInt(
          this.nestConfigService.get<string>('DB_MAX_POOL_SIZE', '20'),
          10,
        ),
        idleTimeoutMs: parseInt(
          this.nestConfigService.get<string>('DB_IDLE_TIMEOUT_MS', '30000'),
          10,
        ),
      },
      jwt: {
        secret: this.nestConfigService.get<string>('JWT_SECRET', ''),
        expiresIn: this.nestConfigService.get<string>('JWT_EXPIRES_IN', '3600'),
        refreshSecret: this.nestConfigService.get<string>(
          'JWT_REFRESH_SECRET',
          '',
        ),
        refreshExpiresIn: this.nestConfigService.get<string>(
          'JWT_REFRESH_EXPIRES_IN',
          '86400',
        ),
      },
      logging: {
        level: this.nestConfigService.get<string>('LOG_LEVEL', 'info'),
      },
    };
  }
}
