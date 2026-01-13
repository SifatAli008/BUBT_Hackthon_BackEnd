import { Module } from '@nestjs/common';
import { APP_INTERCEPTOR } from '@nestjs/core';
import { ConfigModule } from './config/config.module';
import { DatabaseModule } from './database/database.module';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { TransformInterceptor } from './common/interceptors/transform.interceptor';

/**
 * Root application module
 * Imports all global modules and sets up application-level providers
 */
@Module({
  imports: [
    ConfigModule, // Configuration module (global)
    DatabaseModule, // Database module (global)
  ],
  controllers: [AppController],
  providers: [
    AppService,
    // Global interceptors can also be registered here instead of main.ts
    // {
    //   provide: APP_INTERCEPTOR,
    //   useClass: TransformInterceptor,
    // },
  ],
})
export class AppModule {}
