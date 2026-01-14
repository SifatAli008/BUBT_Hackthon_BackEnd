# Architecture Documentation

## Project Structure

\\\
src/
├── common/          # Shared utilities, DTOs, filters
├── config/          # Configuration management
├── database/        # Database services (Prisma + raw SQL)
├── food-items/      # Food items module
├── inventory/       # Inventory module
├── users/           # Users module
└── main.ts          # Application entry point
\\\

## Database Layer

- PrismaService: ORM operations with Prisma Accelerate
- DatabaseService: Raw SQL queries with connection pooling

## Module Organization

Each feature module contains:
- Controller: HTTP endpoints
- Service: Business logic
- DTOs: Data transfer objects
- Module: NestJS module definition
