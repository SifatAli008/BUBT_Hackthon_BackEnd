# FoodLink Backend

A comprehensive NestJS backend application for food waste management, inventory tracking, and community sharing. Built with TypeScript, Prisma ORM, and PostgreSQL.

<p align="center">
  <a href="http://nestjs.com/" target="blank"><img src="https://nestjs.com/img/logo-small.svg" width="120" alt="Nest Logo" /></a>
</p>

## 🚀 Features

- **Food Inventory Management** - Track food items, expiry dates, and quantities
- **User Management** - Complete user profiles and authentication
- **Prisma ORM** - Type-safe database access with Prisma Accelerate support
- **Raw SQL Support** - Direct SQL queries for performance-critical operations
- **JWT Authentication** - Secure token-based authentication
- **RESTful API** - Well-structured REST endpoints
- **TypeScript** - Full type safety throughout the application
- **Error Handling** - Comprehensive error handling and validation
- **Logging** - Winston-based logging system
- **Swagger Documentation** - API documentation (ready for integration)

## 🛠️ Tech Stack

- **Framework**: NestJS 11.x
- **Language**: TypeScript 5.7
- **Database**: PostgreSQL
- **ORM**: Prisma 7.x with Accelerate support
- **Authentication**: JWT (Passport.js)
- **Validation**: class-validator, class-transformer
- **Logging**: Winston
- **Testing**: Jest

## 📋 Prerequisites

- Node.js 18+ 
- npm or yarn
- PostgreSQL database (or Prisma Accelerate account)
- Git

## 🔧 Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/SifatAli008/BUBT_Hackthon_BackEnd.git
   cd BUBT_Hackthon_BackEnd
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Setup environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` and configure:
   ```env
   DATABASE_URL=prisma+postgres://accelerate.prisma-data.net/?api_key=YOUR_API_KEY
   JWT_SECRET=your-secret-key-here
   PORT=3000
   NODE_ENV=development
   ```

4. **Generate Prisma Client**
   ```bash
   npx prisma generate
   ```

5. **Run database migrations**
   ```bash
   npx prisma migrate dev
   ```
   
   Or apply the SQL schema directly:
   ```bash
   # Connect to your database and run schema.sql
   psql -U postgres -d your_database -f schema.sql
   ```

## 🏃 Running the Application

### Development Mode
```bash
npm run start:dev
```

### Production Mode
```bash
npm run build
npm run start:prod
```

### Debug Mode
```bash
npm run start:debug
```

The server will start on `http://localhost:3000` (or the port specified in `.env`)

## 📁 Project Structure

```
src/
├── common/              # Shared utilities and common code
│   ├── decorators/      # Custom decorators
│   ├── dto/            # Data Transfer Objects
│   ├── filters/         # Exception filters
│   ├── interceptors/    # Response interceptors
│   ├── pipes/          # Validation pipes
│   └── utils/          # Utility functions
├── config/              # Configuration management
├── database/            # Database services
│   ├── database.service.ts    # Raw SQL service
│   └── prisma.service.ts      # Prisma ORM service
├── food-items/          # Food items module
├── inventory/           # Inventory module
├── users/              # Users module
├── app.module.ts       # Root application module
├── app.controller.ts    # Root controller
├── app.service.ts      # Root service
└── main.ts             # Application entry point
```

## 🗄️ Database

### Prisma Accelerate (Recommended)

This project supports Prisma Accelerate for:
- Connection pooling
- Query caching
- Global edge network
- Automatic failover

See [Prisma Accelerate Guide](./docs/PRISMA_ACCELERATE.md) for setup instructions.

### Direct PostgreSQL Connection

You can also use a direct PostgreSQL connection:
```env
DATABASE_URL=postgresql://user:password@host:5432/database?sslmode=require
```

### Database Schema

The complete database schema is available in:
- `prisma/schema.prisma` - Prisma schema definition
- `schema.sql` - Raw SQL schema file

## 📡 API Endpoints

### Food Items
- `GET /food-items` - Get all food items
- `POST /food-items` - Create food item
- `GET /food-items/:id` - Get food item by ID
- `PUT /food-items/:id` - Update food item
- `DELETE /food-items/:id` - Delete food item

### Inventory
- `GET /inventory` - Get all inventory items
- `POST /inventory` - Create inventory item
- `GET /inventory/:id` - Get inventory item by ID
- `PUT /inventory/:id` - Update inventory item
- `DELETE /inventory/:id` - Delete inventory item

### Users
- `GET /users/:id` - Get user by ID
- `PUT /users/:id/profile` - Update user profile

*Note: Full API documentation available in [docs/API.md](./docs/API.md)*

## 🧪 Testing

```bash
# Unit tests
npm run test

# E2E tests
npm run test:e2e

# Test coverage
npm run test:cov

# Watch mode
npm run test:watch
```

## 📚 Documentation

Comprehensive documentation is available in the `docs/` folder:

- [Setup Guide](./docs/SETUP.md) - Detailed setup instructions
- [API Documentation](./docs/API.md) - Complete API reference
- [Architecture](./docs/ARCHITECTURE.md) - System architecture overview
- [Deployment](./docs/DEPLOYMENT.md) - Deployment guide
- [Migrations](./docs/MIGRATIONS.md) - Database migration guide
- [Prisma Accelerate](./docs/PRISMA_ACCELERATE.md) - Prisma Accelerate setup
- [Troubleshooting](./docs/TROUBLESHOOTING.md) - Common issues and solutions
- [Contributing](./docs/CONTRIBUTING.md) - Contribution guidelines

## 🔐 Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `DATABASE_URL` | Database connection string | Yes | - |
| `JWT_SECRET` | Secret key for JWT tokens | Yes | - |
| `PORT` | Server port | No | 3000 |
| `NODE_ENV` | Environment (development/production) | No | development |
| `JWT_EXPIRES_IN` | JWT token expiration time | No | 3600 |
| `LOG_LEVEL` | Logging level | No | info |

## 🏗️ Development

### Code Style

- Follow TypeScript best practices
- Use ESLint and Prettier for code formatting
- Write meaningful commit messages
- Add tests for new features

### Commit Message Format

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes
- `refactor:` - Code refactoring
- `test:` - Adding tests
- `chore:` - Maintenance tasks

## 🚢 Deployment

See [Deployment Guide](./docs/DEPLOYMENT.md) for detailed deployment instructions.

### Quick Deploy

1. Build the application:
   ```bash
   npm run build
   ```

2. Set production environment variables

3. Start the server:
   ```bash
   npm run start:prod
   ```

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guidelines](./docs/CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## 📝 License

This project is licensed under the UNLICENSED License.

## 👥 Authors

- **Sifat Ali** - [GitHub](https://github.com/SifatAli008)

## 🙏 Acknowledgments

- [NestJS](https://nestjs.com/) - Progressive Node.js framework
- [Prisma](https://www.prisma.io/) - Next-generation ORM
- [TypeScript](https://www.typescriptlang.org/) - Typed JavaScript

## 📞 Support

For support, please open an issue in the [GitHub repository](https://github.com/SifatAli008/BUBT_Hackthon_BackEnd/issues).

## 🔗 Links

- [Repository](https://github.com/SifatAli008/BUBT_Hackthon_BackEnd)
- [NestJS Documentation](https://docs.nestjs.com)
- [Prisma Documentation](https://www.prisma.io/docs)

---

<p align="center">Made with ❤️ for BUBT Hackathon</p>
