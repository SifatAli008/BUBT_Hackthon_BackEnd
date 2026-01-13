import {
  Injectable,
  NotFoundException,
  ConflictException,
} from '@nestjs/common';
import { DatabaseService } from '../database/database.service';
import * as bcrypt from 'bcrypt';

/**
 * User interface matching database structure
 */
export interface User {
  id: string;
  email: string;
  name: string;
  password_hash: string;
  household_id: string | null;
  role: string;
  created_at: Date;
  updated_at: Date;
}

/**
 * Users Service
 * Handles user-related database operations using raw SQL queries
 * Optimized queries with proper indexes
 */
@Injectable()
export class UsersService {
  constructor(private databaseService: DatabaseService) {}

  /**
   * Find user by email (optimized with index)
   * Uses email index for fast lookup
   */
  async findByEmail(email: string): Promise<User | null> {
    const user = await this.databaseService.queryOne<User>(
      `
      SELECT id, email, name, password_hash, household_id, role, created_at, updated_at
      FROM users
      WHERE email = $1
      LIMIT 1
    `,
      [email.toLowerCase().trim()],
    );

    return user || null;
  }

  /**
   * Find user by ID (optimized with primary key)
   */
  async findById(id: string): Promise<User | null> {
    const user = await this.databaseService.queryOne<User>(
      `
      SELECT id, email, name, password_hash, household_id, role, created_at, updated_at
      FROM users
      WHERE id = $1
      LIMIT 1
    `,
      [id],
    );

    return user || null;
  }

  /**
   * Create new user
   * Uses parameterized query to prevent SQL injection
   * Returns user without password hash
   */
  async create(
    email: string,
    name: string,
    passwordHash: string,
    role: string = 'family',
    householdId?: string,
  ): Promise<Omit<User, 'password_hash'>> {
    // Check if user already exists
    const existingUser = await this.findByEmail(email);
    if (existingUser) {
      throw new ConflictException('User with this email already exists');
    }

    // Create user with optimized insert query
    const result = await this.databaseService.query<
      Omit<User, 'password_hash'>
    >(
      `
      INSERT INTO users (email, name, password_hash, role, household_id)
      VALUES ($1, $2, $3, $4, $5)
      RETURNING id, email, name, household_id, role, created_at, updated_at
    `,
      [
        email.toLowerCase().trim(),
        name.trim(),
        passwordHash,
        role,
        householdId || null,
      ],
    );

    return result.rows[0];
  }

  /**
   * Update user password
   * Uses parameterized query for security
   */
  async updatePassword(userId: string, passwordHash: string): Promise<void> {
    const result = await this.databaseService.query(
      `
      UPDATE users
      SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
      WHERE id = $2
    `,
      [passwordHash, userId],
    );

    if (result.rowCount === 0) {
      throw new NotFoundException('User not found');
    }
  }

  /**
   * Get user profile (without password)
   * Optimized query with selected fields only
   */
  async getProfile(
    userId: string,
  ): Promise<Omit<User, 'password_hash'> | null> {
    const user = await this.databaseService.queryOne<
      Omit<User, 'password_hash'>
    >(
      `
      SELECT id, email, name, household_id, role, created_at, updated_at
      FROM users
      WHERE id = $1
      LIMIT 1
    `,
      [userId],
    );

    return user || null;
  }

  /**
   * Verify password against hash
   * Uses bcrypt for secure password comparison
   */
  async verifyPassword(
    plainPassword: string,
    hashedPassword: string,
  ): Promise<boolean> {
    return bcrypt.compare(plainPassword, hashedPassword);
  }

  /**
   * Hash password using bcrypt
   * Uses salt rounds for security
   */
  async hashPassword(password: string): Promise<string> {
    const saltRounds = 10;
    return bcrypt.hash(password, saltRounds);
  }

  /**
   * Update user profile
   * Updates name and/or household_id
   */
  async updateProfile(
    userId: string,
    name?: string,
    householdId?: string,
  ): Promise<Omit<User, 'password_hash'>> {
    // Build dynamic update query
    const updates: string[] = [];
    const params: any[] = [];
    let paramIndex = 1;

    if (name !== undefined) {
      updates.push(`name = $${paramIndex}`);
      params.push(name.trim());
      paramIndex++;
    }

    if (householdId !== undefined) {
      updates.push(`household_id = $${paramIndex}`);
      params.push(householdId || null);
      paramIndex++;
    }

    if (updates.length === 0) {
      const profile = await this.getProfile(userId);
      if (!profile) {
        throw new NotFoundException('User not found');
      }
      return profile; // No updates, return current profile
    }

    // Add updated_at timestamp
    updates.push(`updated_at = CURRENT_TIMESTAMP`);

    // Add userId to params
    params.push(userId);

    const result = await this.databaseService.query<
      Omit<User, 'password_hash'>
    >(
      `
      UPDATE users
      SET ${updates.join(', ')}
      WHERE id = $${paramIndex}
      RETURNING id, email, name, household_id, role, created_at, updated_at
    `,
      params,
    );

    if (result.rowCount === 0) {
      throw new NotFoundException('User not found');
    }

    return result.rows[0];
  }
}
