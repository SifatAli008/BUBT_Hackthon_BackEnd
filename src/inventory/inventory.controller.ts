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
  Request,
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
import { InventoryService } from './inventory.service';
import {
  InventoryItemDto,
  CreateInventoryItemDto,
  UpdateInventoryItemDto,
  InventoryFilterDto,
} from './dto/inventory-item.dto';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { PaginationDto } from '../common/dto/pagination.dto';
import { ApiResponse as ApiResponseDto } from '../common/dto/response.dto';

/**
 * Inventory Controller
 * Handles inventory items endpoints
 */
@ApiTags('Inventory')
@Controller('inventory')
@UseGuards(JwtAuthGuard)
@ApiBearerAuth('JWT-auth')
export class InventoryController {
  constructor(private readonly inventoryService: InventoryService) {}

  @Post()
  @HttpCode(HttpStatus.CREATED)
  @ApiOperation({ summary: 'Create a new inventory item' })
  @ApiResponse({
    status: 201,
    description: 'Inventory item successfully created',
    type: InventoryItemDto,
  })
  @ApiResponse({ status: 400, description: 'Validation error' })
  async create(
    @Request() req,
    @Body() createDto: CreateInventoryItemDto,
  ): Promise<ApiResponseDto<InventoryItemDto>> {
    const userId = req.user.sub;
    const item = await this.inventoryService.create(userId, createDto);
    return ApiResponseDto.success(item, 'Inventory item successfully created');
  }

  @Get()
  @ApiOperation({ summary: 'Get all inventory items' })
  @ApiQuery({
    name: 'category',
    required: false,
    description: 'Filter by category',
  })
  @ApiQuery({
    name: 'location',
    required: false,
    description: 'Filter by location',
  })
  @ApiQuery({
    name: 'expiringWithinDays',
    required: false,
    description: 'Filter expiring items',
  })
  @ApiQuery({
    name: 'expired',
    required: false,
    description: 'Filter expired items',
  })
  @ApiResponse({
    status: 200,
    description: 'Inventory items retrieved successfully',
    type: [InventoryItemDto],
  })
  async findAll(
    @Request() req,
    @Query() filter?: InventoryFilterDto,
    @Query() pagination?: PaginationDto,
  ): Promise<ApiResponseDto<InventoryItemDto[]>> {
    const userId = req.user.sub;
    const limit = pagination?.getLimit() || 50;
    const offset = pagination?.getOffset() || 0;

    const result = await this.inventoryService.findAll(
      userId,
      filter,
      limit,
      offset,
    );
    return ApiResponseDto.paginated(
      result.items,
      result.total,
      Math.floor(offset / limit) + 1,
      limit,
      'Inventory items retrieved successfully',
    );
  }

  @Get('expiring')
  @ApiOperation({ summary: 'Get expiring inventory items' })
  @ApiQuery({
    name: 'days',
    required: false,
    description: 'Days to check ahead',
    type: Number,
  })
  @ApiResponse({
    status: 200,
    description: 'Expiring items retrieved successfully',
    type: [InventoryItemDto],
  })
  async getExpiringItems(
    @Request() req,
    @Query('days') days?: number,
  ): Promise<ApiResponseDto<InventoryItemDto[]>> {
    const userId = req.user.sub;
    const items = await this.inventoryService.getExpiringItems(
      userId,
      days || 7,
    );
    return ApiResponseDto.success(
      items,
      'Expiring items retrieved successfully',
    );
  }

  @Get('expired')
  @ApiOperation({ summary: 'Get expired inventory items' })
  @ApiResponse({
    status: 200,
    description: 'Expired items retrieved successfully',
    type: [InventoryItemDto],
  })
  async getExpiredItems(
    @Request() req,
  ): Promise<ApiResponseDto<InventoryItemDto[]>> {
    const userId = req.user.sub;
    const items = await this.inventoryService.getExpiredItems(userId);
    return ApiResponseDto.success(
      items,
      'Expired items retrieved successfully',
    );
  }

  @Get('search')
  @ApiOperation({ summary: 'Search inventory items by name' })
  @ApiQuery({ name: 'q', required: true, description: 'Search term' })
  @ApiQuery({ name: 'limit', required: false, description: 'Result limit' })
  @ApiResponse({
    status: 200,
    description: 'Search results retrieved successfully',
    type: [InventoryItemDto],
  })
  async search(
    @Request() req,
    @Query('q') searchTerm: string,
    @Query('limit') limit?: number,
  ): Promise<ApiResponseDto<InventoryItemDto[]>> {
    const userId = req.user.sub;
    const items = await this.inventoryService.searchByName(
      userId,
      searchTerm,
      limit || 20,
    );
    return ApiResponseDto.success(
      items,
      'Search results retrieved successfully',
    );
  }

  @Get('statistics')
  @ApiOperation({ summary: 'Get inventory statistics' })
  @ApiResponse({
    status: 200,
    description: 'Statistics retrieved successfully',
  })
  async getStatistics(@Request() req): Promise<ApiResponseDto<any>> {
    const userId = req.user.sub;
    const stats = await this.inventoryService.getStatistics(userId);
    return ApiResponseDto.success(stats, 'Statistics retrieved successfully');
  }

  @Get(':id')
  @ApiOperation({ summary: 'Get an inventory item by ID' })
  @ApiResponse({
    status: 200,
    description: 'Inventory item retrieved successfully',
    type: InventoryItemDto,
  })
  @ApiResponse({ status: 404, description: 'Inventory item not found' })
  async findOne(
    @Request() req,
    @Param('id') id: string,
  ): Promise<ApiResponseDto<InventoryItemDto>> {
    const userId = req.user.sub;
    const item = await this.inventoryService.findOne(id, userId);
    return ApiResponseDto.success(
      item,
      'Inventory item retrieved successfully',
    );
  }

  @Patch(':id')
  @ApiOperation({ summary: 'Update an inventory item' })
  @ApiResponse({
    status: 200,
    description: 'Inventory item successfully updated',
    type: InventoryItemDto,
  })
  @ApiResponse({ status: 404, description: 'Inventory item not found' })
  async update(
    @Request() req,
    @Param('id') id: string,
    @Body() updateDto: UpdateInventoryItemDto,
  ): Promise<ApiResponseDto<InventoryItemDto>> {
    const userId = req.user.sub;
    const item = await this.inventoryService.update(id, userId, updateDto);
    return ApiResponseDto.success(item, 'Inventory item successfully updated');
  }

  @Delete(':id')
  @HttpCode(HttpStatus.NO_CONTENT)
  @ApiOperation({ summary: 'Delete an inventory item' })
  @ApiResponse({
    status: 204,
    description: 'Inventory item successfully deleted',
  })
  @ApiResponse({ status: 404, description: 'Inventory item not found' })
  async remove(@Request() req, @Param('id') id: string): Promise<void> {
    const userId = req.user.sub;
    await this.inventoryService.remove(id, userId);
  }
}
