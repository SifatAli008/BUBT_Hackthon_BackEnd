import {
  ExceptionFilter,
  Catch,
  ArgumentsHost,
  HttpException,
  Logger,
} from '@nestjs/common';
import { Request, Response } from 'express';

/**
 * Type guard for exception response object
 */
interface ExceptionResponseObject {
  message?: string | string[];
  error?: string;
}

/**
 * Global HTTP exception filter
 * Catches and formats all HTTP exceptions for consistent error responses
 */
@Catch(HttpException)
export class HttpExceptionFilter implements ExceptionFilter {
  private readonly logger = new Logger(HttpExceptionFilter.name);

  catch(exception: HttpException, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();
    const request = ctx.getRequest<Request>();
    const status = exception.getStatus();
    const exceptionResponse = exception.getResponse();

    // Format error response
    const message =
      typeof exceptionResponse === 'string'
        ? exceptionResponse
        : (typeof exceptionResponse === 'object' &&
          exceptionResponse !== null &&
          'message' in exceptionResponse
            ? (exceptionResponse as ExceptionResponseObject).message
            : null) || exception.message;

    const error =
      typeof exceptionResponse === 'object' &&
      exceptionResponse !== null &&
      'error' in exceptionResponse
        ? (exceptionResponse as ExceptionResponseObject).error
        : exception.name;

    const errorResponse = {
      statusCode: status,
      timestamp: new Date().toISOString(),
      path: request.url,
      method: request.method,
      message: Array.isArray(message) ? message.join(', ') : message,
      error,
    };

    // Log error for debugging (only log non-4xx errors in production)
    if (status >= 500) {
      this.logger.error(
        `${request.method} ${request.url} - ${status} - ${exception.message}`,
        exception.stack,
      );
    } else if (status >= 400) {
      this.logger.warn(
        `${request.method} ${request.url} - ${status} - ${exception.message}`,
      );
    }

    response.status(status).json(errorResponse);
  }
}
