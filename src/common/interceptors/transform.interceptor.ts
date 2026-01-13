import {
  Injectable,
  NestInterceptor,
  ExecutionContext,
  CallHandler,
} from '@nestjs/common';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { ApiResponse } from '../dto/response.dto';

/**
 * Transform interceptor
 * Wraps responses in standard ApiResponse format
 */
@Injectable()
export class TransformInterceptor<T> implements NestInterceptor<
  T,
  ApiResponse<T>
> {
  intercept(
    context: ExecutionContext,
    next: CallHandler,
  ): Observable<ApiResponse<T>> {
    return next.handle().pipe(
      map((data) => {
        // If data is already an ApiResponse, return it as is
        if (data && typeof data === 'object' && 'success' in data) {
          return data;
        }
        // Otherwise, wrap in ApiResponse
        return ApiResponse.success(data);
      }),
    );
  }
}
