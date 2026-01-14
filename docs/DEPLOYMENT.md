# Deployment Guide

## Environment Variables

Required environment variables:
- DATABASE_URL: Prisma Accelerate connection string
- JWT_SECRET: Secret key for JWT tokens
- PORT: Server port (default: 3000)
- NODE_ENV: Environment (development/production)

## Production Build

1. Build the application:
```bash
npm run build
```

2. Start production server:
```bash
npm run start:prod
```

## Docker (Optional)

Build Docker image:
```bash
docker build -t foodlink-backend .
```
