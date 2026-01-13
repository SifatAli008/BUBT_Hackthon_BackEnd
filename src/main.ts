import { NestFactory } from '@nestjs/core';
import { ValidationPipe, Logger } from '@nestjs/common';
import { SwaggerModule, DocumentBuilder } from '@nestjs/swagger';
import { AppModule } from './app.module';
import { HttpExceptionFilter } from './common/filters/http-exception.filter';
import { AllExceptionsFilter } from './common/filters/all-exceptions.filter';
import { TransformInterceptor } from './common/interceptors/transform.interceptor';
import helmet from 'helmet';
import * as compression from 'compression';
import { ConfigService } from './config/config.service';

/**
 * Bootstrap the application
 * Sets up middleware, filters, pipes, and Swagger documentation
 */
async function bootstrap() {
  const logger = new Logger('Bootstrap');

  try {
    // Create NestJS application
    const app = await NestFactory.create(AppModule, {
      bufferLogs: true, // Buffer logs until logger is ready
    });

    // Get configuration service
    const configService = app.get(ConfigService);
    const config = configService.getConfig();

    // Enable CORS with proper configuration
    app.enableCors({
      origin:
        config.nodeEnv === 'production'
          ? process.env.CORS_ORIGIN?.split(',') || false
          : true, // Allow all origins in development
      credentials: true,
      methods: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS'],
      allowedHeaders: ['Content-Type', 'Authorization', 'Accept'],
      exposedHeaders: ['Authorization'],
    });

    // Security middleware
    app.use(
      helmet({
        contentSecurityPolicy: config.nodeEnv === 'production',
        crossOriginEmbedderPolicy: false,
      }),
    );

    // Compression middleware (gzip)
    app.use(compression());

    // Global prefix for all routes
    app.setGlobalPrefix('api/v1', {
      exclude: ['/health', '/'],
    });

    // Global validation pipe
    app.useGlobalPipes(
      new ValidationPipe({
        whitelist: true, // Strip properties that don't have decorators
        forbidNonWhitelisted: true, // Throw error if non-whitelisted properties are present
        transform: true, // Automatically transform payloads to DTO instances
        transformOptions: {
          enableImplicitConversion: true, // Enable implicit type conversion
        },
        disableErrorMessages: config.nodeEnv === 'production', // Hide error messages in production
      }),
    );

    // Global exception filters (order matters - more specific first)
    app.useGlobalFilters(new HttpExceptionFilter(), new AllExceptionsFilter());

    // Global interceptors
    app.useGlobalInterceptors(new TransformInterceptor());

    // Swagger API documentation (only in development/staging)
    if (config.nodeEnv !== 'production') {
      const swaggerConfig = new DocumentBuilder()
        .setTitle('FoodLink API')
        .setDescription(
          'FoodLink Backend API - A comprehensive food waste management platform',
        )
        .setVersion('1.0')
        .addBearerAuth(
          {
            type: 'http',
            scheme: 'bearer',
            bearerFormat: 'JWT',
            name: 'JWT',
            description: 'Enter JWT token',
            in: 'header',
          },
          'JWT-auth', // This name here is important for matching up with @ApiBearerAuth() in your controller!
        )
        .addTag('Health', 'Health check endpoints')
        .addTag('Auth', 'Authentication endpoints')
        .addTag('Users', 'User management endpoints')
        .addTag('Inventory', 'Inventory management endpoints')
        .addTag('Community', 'Community features endpoints')
        .addTag('Restaurant', 'Restaurant module endpoints')
        .addTag('NGO', 'NGO module endpoints')
        .addTag('Shop', 'Shop module endpoints')
        .build();

      const document = SwaggerModule.createDocument(app, swaggerConfig);
      SwaggerModule.setup('api/docs', app, document, {
        swaggerOptions: {
          persistAuthorization: true, // Persist authorization token in browser
          tagsSorter: 'alpha', // Sort tags alphabetically
          operationsSorter: 'alpha', // Sort operations alphabetically
        },
      });

      logger.log(`Swagger documentation available at /api/docs`);
    }

    // Start the application
    const port = config.port;
    await app.listen(port);

    logger.log(`Application is running on: http://localhost:${port}`);
    logger.log(`Environment: ${config.nodeEnv}`);
    logger.log(`API prefix: /api/v1`);
  } catch (error) {
    logger.error('Failed to start application', error);
    process.exit(1);
  }
}

bootstrap();
