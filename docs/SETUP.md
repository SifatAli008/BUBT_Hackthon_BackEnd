# Setup Guide

## Prerequisites
- Node.js 18+
- npm or yarn
- PostgreSQL database (or Prisma Accelerate)

## Installation

1. Install dependencies:
```bash
npm install
```

2. Setup environment variables:
```bash
cp .env.example .env
```

3. Generate Prisma Client:
```bash
npx prisma generate
```

4. Run migrations:
```bash
npx prisma migrate dev
```

5. Start development server:
```bash
npm run start:dev
```
