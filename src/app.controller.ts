import { Controller, Get } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse } from '@nestjs/swagger';
import { AppService } from './app.service';
import { Public } from './common/decorators/public.decorator';
import { ApiResponse as ApiResponseDto } from './common/dto/response.dto';

/**
 * Root application controller
 * Handles basic health checks and API information
 */
@ApiTags('Health')
@Controller()
export class AppController {
  constructor(private readonly appService: AppService) {}

  @Public()
  @Get()
  @ApiOperation({ summary: 'Get API information' })
  @ApiResponse({ status: 200, description: 'API information' })
  getHello(): ApiResponseDto<{ message: string; version: string }> {
    return ApiResponseDto.success(
      {
        message: this.appService.getHello(),
        version: '1.0.0',
      },
      'Welcome to FoodLink API',
    );
  }

  @Public()
  @Get('health')
  @ApiOperation({ summary: 'Health check endpoint' })
  @ApiResponse({ status: 200, description: 'Service is healthy' })
  @ApiResponse({ status: 503, description: 'Service is unhealthy' })
  async getHealth(): Promise<
    ApiResponseDto<{ status: string; timestamp: string }>
  > {
    const isHealthy = await this.appService.healthCheck();
    const statusCode = isHealthy ? 200 : 503;

    return ApiResponseDto.success(
      {
        status: isHealthy ? 'healthy' : 'unhealthy',
        timestamp: new Date().toISOString(),
      },
      isHealthy ? 'Service is healthy' : 'Service is unhealthy',
    );
  }
}
