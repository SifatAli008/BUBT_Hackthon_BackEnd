import { Module } from '@nestjs/common';
import { ConfigModule } from './config/config.module';
import { DatabaseModule } from './database/database.module';
import { UsersModule } from './users/users.module';
import { AuthModule } from './auth/auth.module';
import { FoodItemsModule } from './food-items/food-items.module';
import { InventoryModule } from './inventory/inventory.module';
import { AppController } from './app.controller';
import { AppService } from './app.service';

/**
 * Root application module
 * Imports all global modules and sets up application-level providers
 */
@Module({
  imports: [
    ConfigModule, // Configuration module (global)
    DatabaseModule, // Database module (global)
    UsersModule, // Users module (global)
    AuthModule, // Auth module
    FoodItemsModule, // Food items reference data module
    InventoryModule, // Inventory items module
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
