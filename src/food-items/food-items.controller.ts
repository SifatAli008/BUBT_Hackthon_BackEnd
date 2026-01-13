import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
  Query,
  UseGuards,
  HttpCode,
  HttpStatus,
} from '@nestjs/common';
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiBearerAuth,
  ApiQuery,
} from '@nestjs/swagger';
import { FoodItemsService } from './food-items.service';
import {
  FoodItemDto,
  CreateFoodItemDto,
  UpdateFoodItemDto,
} from './dto/food-item.dto';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { Public } from '../common/decorators/public.decorator';
import { PaginationDto } from '../common/dto/pagination.dto';
import { ApiResponse as ApiResponseDto } from '../common/dto/response.dto';

/**
 * Food Items Controller
 * Handles food items reference data endpoints
 */
@ApiTags('Food Items')
@Controller('food-items')
@UseGuards(JwtAuthGuard)
@ApiBearerAuth('JWT-auth')
export class FoodItemsController {
  constructor(private readonly foodItemsService: FoodItemsService) {}

  @Post()
  @HttpCode(HttpStatus.CREATED)
  @ApiOperation({ summary: 'Create a new food item' })
  @ApiResponse({
    status: 201,
    description: 'Food item successfully created',
    type: FoodItemDto,
  })
  @ApiResponse({ status: 400, description: 'Validation error' })
  @ApiResponse({ status: 409, description: 'Food item already exists' })
  async create(
    @Body() createDto: CreateFoodItemDto,
  ): Promise<ApiResponseDto<FoodItemDto>> {
    const item = await this.foodItemsService.create(createDto);
    return ApiResponseDto.success(item, 'Food item successfully created');
  }

  @Public()
  @Get()
  @ApiOperation({ summary: 'Get all food items' })
  @ApiQuery({
    name: 'category',
    required: false,
    description: 'Filter by category',
  })
  @ApiResponse({
    status: 200,
    description: 'Food items retrieved successfully',
    type: [FoodItemDto],
  })
  async findAll(
    @Query('category') category?: string,
    @Query() pagination?: PaginationDto,
  ): Promise<ApiResponseDto<FoodItemDto[]>> {
    const limit = pagination?.getLimit() || 100;
    const offset = pagination?.getOffset() || 0;

    const result = await this.foodItemsService.findAll(category, limit, offset);
    return ApiResponseDto.paginated(
      result.items,
      result.total,
      Math.floor(offset / limit) + 1,
      limit,
      'Food items retrieved successfully',
    );
  }

  @Get('categories')
  @ApiOperation({ summary: 'Get all food item categories' })
  @ApiResponse({
    status: 200,
    description: 'Categories retrieved successfully',
  })
  async getCategories(): Promise<ApiResponseDto<string[]>> {
    const categories = await this.foodItemsService.getCategories();
    return ApiResponseDto.success(
      categories,
      'Categories retrieved successfully',
    );
  }

  @Get('search')
  @ApiOperation({ summary: 'Search food items by name' })
  @ApiQuery({ name: 'q', required: true, description: 'Search term' })
  @ApiQuery({ name: 'limit', required: false, description: 'Result limit' })
  @ApiResponse({
    status: 200,
    description: 'Search results retrieved successfully',
    type: [FoodItemDto],
  })
  async search(
    @Query('q') searchTerm: string,
    @Query('limit') limit?: number,
  ): Promise<ApiResponseDto<FoodItemDto[]>> {
    const items = await this.foodItemsService.searchByName(
      searchTerm,
      limit || 20,
    );
    return ApiResponseDto.success(
      items,
      'Search results retrieved successfully',
    );
  }

  @Get(':id')
  @ApiOperation({ summary: 'Get a food item by ID' })
  @ApiResponse({
    status: 200,
    description: 'Food item retrieved successfully',
    type: FoodItemDto,
  })
  @ApiResponse({ status: 404, description: 'Food item not found' })
  async findOne(@Param('id') id: string): Promise<ApiResponseDto<FoodItemDto>> {
    const item = await this.foodItemsService.findOne(id);
    return ApiResponseDto.success(item, 'Food item retrieved successfully');
  }

  @Patch(':id')
  @ApiOperation({ summary: 'Update a food item' })
  @ApiResponse({
    status: 200,
    description: 'Food item successfully updated',
    type: FoodItemDto,
  })
  @ApiResponse({ status: 404, description: 'Food item not found' })
  async update(
    @Param('id') id: string,
    @Body() updateDto: UpdateFoodItemDto,
  ): Promise<ApiResponseDto<FoodItemDto>> {
    const item = await this.foodItemsService.update(id, updateDto);
    return ApiResponseDto.success(item, 'Food item successfully updated');
  }

  @Delete(':id')
  @HttpCode(HttpStatus.NO_CONTENT)
  @ApiOperation({ summary: 'Delete a food item' })
  @ApiResponse({ status: 204, description: 'Food item successfully deleted' })
  @ApiResponse({ status: 404, description: 'Food item not found' })
  async remove(@Param('id') id: string): Promise<void> {
    await this.foodItemsService.remove(id);
  }
}
