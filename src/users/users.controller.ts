import {
  Controller,
  Get,
  Put,
  Body,
  UseGuards,
  Request,
  HttpCode,
  HttpStatus,
  NotFoundException,
} from '@nestjs/common';
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiBearerAuth,
} from '@nestjs/swagger';
import { UsersService } from './users.service';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { UserDto } from './dto/user.dto';
import { UpdateUserDto } from './dto/update-user.dto';
import { ApiResponse as ApiResponseDto } from '../common/dto/response.dto';

/**
 * Users Controller
 * Handles user profile management endpoints
 */
@ApiTags('Users')
@Controller('users')
@UseGuards(JwtAuthGuard)
@ApiBearerAuth('JWT-auth')
export class UsersController {
  constructor(private usersService: UsersService) {}

  @Get('profile')
  @ApiOperation({ summary: 'Get current user profile' })
  @ApiResponse({
    status: 200,
    description: 'User profile retrieved',
    type: UserDto,
  })
  @ApiResponse({ status: 401, description: 'Unauthorized' })
  async getProfile(@Request() req): Promise<ApiResponseDto<UserDto>> {
    const userId = req.user.sub;
    const user = await this.usersService.getProfile(userId);

    if (!user) {
      throw new NotFoundException('User not found');
    }

    // Map to DTO
    const userDto: UserDto = {
      id: user.id,
      email: user.email,
      name: user.name,
      role: user.role,
      householdId: user.household_id || undefined,
      createdAt: user.created_at,
      updatedAt: user.updated_at,
    };

    return ApiResponseDto.success(userDto, 'User profile retrieved');
  }

  @Put('profile')
  @HttpCode(HttpStatus.OK)
  @ApiOperation({ summary: 'Update current user profile' })
  @ApiResponse({
    status: 200,
    description: 'User profile successfully updated',
    type: UserDto,
  })
  @ApiResponse({ status: 401, description: 'Unauthorized' })
  @ApiResponse({ status: 404, description: 'User not found' })
  async updateProfile(
    @Request() req,
    @Body() updateDto: UpdateUserDto,
  ): Promise<ApiResponseDto<UserDto>> {
    const userId = req.user.sub;
    const user = await this.usersService.updateProfile(
      userId,
      updateDto.name,
      updateDto.householdId,
    );

    // Map to DTO
    const userDto: UserDto = {
      id: user.id,
      email: user.email,
      name: user.name,
      role: user.role,
      householdId: user.household_id || undefined,
      createdAt: user.created_at,
      updatedAt: user.updated_at,
    };

    return ApiResponseDto.success(userDto, 'User profile successfully updated');
  }
}
