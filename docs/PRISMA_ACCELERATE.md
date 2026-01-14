# Prisma Accelerate Configuration

## Setup

1. Get your Prisma Accelerate connection string from Prisma Cloud
2. Add it to your \.env\ file:
\\\
DATABASE_URL=prisma+postgres://accelerate.prisma-data.net/?api_key=YOUR_API_KEY
\\\

## Benefits

- Connection pooling
- Query caching
- Global edge network
- Automatic failover

## Usage

The PrismaService automatically uses Prisma Accelerate when the connection string starts with \prisma+postgres://\.
