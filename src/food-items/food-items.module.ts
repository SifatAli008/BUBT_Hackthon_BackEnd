import { Module } from '@nestjs/common';
import { FoodItemsService } from './food-items.service';
import { FoodItemsController } from './food-items.controller';

/**
 * Food Items Module
 * Provides food items reference data management
 */
@Module({
  controllers: [FoodItemsController],
  providers: [FoodItemsService],
  exports: [FoodItemsService],
})
export class FoodItemsModule {}
