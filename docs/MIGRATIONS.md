# Database Migration Guide

## Using Prisma Migrate

1. Create a new migration:
\\\ash
npx prisma migrate dev --name migration_name
\\\

2. Apply migrations to production:
\\\ash
npx prisma migrate deploy
\\\

3. Reset database (development only):
\\\ash
npx prisma migrate reset
\\\

## Using Raw SQL

The \schema.sql\ file contains the complete database schema. You can apply it directly to your database.
