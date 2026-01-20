# Foodlink Backend - Implementation Plan

## Overview
This document outlines the feature-based implementation plan for the Foodlink backend API. The project is organized into logical features/modules that can be developed independently while maintaining clear dependencies.

---

## Feature Breakdown

### Phase 0: Foundation & Infrastructure
**Priority: CRITICAL** | **Dependencies: None**

#### 0.1 Database Setup & Migrations
- **Tables**: All schema tables
- **Components**:
  - Database connection pool
  - Migration system
  - Schema initialization
  - Connection health checks
- **Files**:
  - `database/connection.go` - DB connection management
  - `database/migrations/` - Migration files
  - `database/schema.go` - Schema utilities

#### 0.2 Core Infrastructure
- **Components**:
  - Middleware (logging, CORS, error handling, request ID)
  - Response utilities
  - Validation helpers
  - Error types
- **Files**:
  - `middleware/` - All middleware
  - `utils/response.go` - Standardized responses
  - `utils/validator.go` - Validation helpers
  - `errors/` - Custom error types

---

### Phase 1: Authentication & User Management
**Priority: CRITICAL** | **Dependencies: Phase 0**

#### 1.1 Authentication Feature
- **Tables**: `users`
- **Endpoints**:
  - `POST /api/v1/auth/register` - User registration
  - `POST /api/v1/auth/login` - User login
  - `POST /api/v1/auth/logout` - User logout
  - `POST /api/v1/auth/refresh` - Refresh token
  - `GET /api/v1/auth/me` - Get current user
- **Components**:
  - JWT token generation & validation
  - Password hashing (bcrypt)
  - Session management
  - Authentication middleware
- **Files**:
  - `features/auth/models.go` - User model
  - `features/auth/repository.go` - Database operations
  - `features/auth/service.go` - Business logic
  - `features/auth/handlers.go` - HTTP handlers
  - `features/auth/routes.go` - Route definitions
  - `features/auth/middleware.go` - Auth middleware

---

### Phase 2: Core Family Features
**Priority: HIGH** | **Dependencies: Phase 1**

#### 2.1 Food Items (Reference Data)
- **Tables**: `food_items`
- **Endpoints**:
  - `GET /api/v1/food-items` - List all food items
  - `GET /api/v1/food-items/:id` - Get food item details
  - `POST /api/v1/food-items` - Create food item (admin)
  - `PUT /api/v1/food-items/:id` - Update food item (admin)
  - `DELETE /api/v1/food-items/:id` - Delete food item (admin)
- **Files**:
  - `features/food_items/models.go`
  - `features/food_items/repository.go`
  - `features/food_items/service.go`
  - `features/food_items/handlers.go`
  - `features/food_items/routes.go`

#### 2.2 Inventory Management
- **Tables**: `inventory_items`
- **Endpoints**:
  - `GET /api/v1/inventory` - List user's inventory
  - `GET /api/v1/inventory/:id` - Get inventory item
  - `POST /api/v1/inventory` - Add inventory item
  - `PUT /api/v1/inventory/:id` - Update inventory item
  - `DELETE /api/v1/inventory/:id` - Delete inventory item
  - `GET /api/v1/inventory/expiring` - Get expiring items
  - `GET /api/v1/inventory/expired` - Get expired items
- **Files**:
  - `features/inventory/models.go`
  - `features/inventory/repository.go`
  - `features/inventory/service.go`
  - `features/inventory/handlers.go`
  - `features/inventory/routes.go`

#### 2.3 Consumption Tracking
- **Tables**: `consumption_logs`
- **Endpoints**:
  - `GET /api/v1/consumption` - List consumption logs
  - `GET /api/v1/consumption/:id` - Get consumption log
  - `POST /api/v1/consumption` - Log consumption
  - `PUT /api/v1/consumption/:id` - Update consumption log
  - `DELETE /api/v1/consumption/:id` - Delete consumption log
  - `GET /api/v1/consumption/stats` - Get consumption statistics
- **Files**:
  - `features/consumption/models.go`
  - `features/consumption/repository.go`
  - `features/consumption/service.go`
  - `features/consumption/handlers.go`
  - `features/consumption/routes.go`

#### 2.4 Shopping Lists
- **Tables**: `shopping_list_items`
- **Endpoints**:
  - `GET /api/v1/shopping-list` - Get shopping list
  - `POST /api/v1/shopping-list` - Add item to shopping list
  - `PUT /api/v1/shopping-list/:id` - Update shopping list item
  - `DELETE /api/v1/shopping-list/:id` - Remove item
  - `PUT /api/v1/shopping-list/:id/purchase` - Mark as purchased
  - `POST /api/v1/shopping-list/from-inventory` - Generate from inventory
- **Files**:
  - `features/shopping_list/models.go`
  - `features/shopping_list/repository.go`
  - `features/shopping_list/service.go`
  - `features/shopping_list/handlers.go`
  - `features/shopping_list/routes.go`

#### 2.5 Meal Planning
- **Tables**: `meal_plans`
- **Endpoints**:
  - `GET /api/v1/meal-plans` - Get meal plans (with date filters)
  - `GET /api/v1/meal-plans/:id` - Get meal plan details
  - `POST /api/v1/meal-plans` - Create meal plan
  - `PUT /api/v1/meal-plans/:id` - Update meal plan
  - `DELETE /api/v1/meal-plans/:id` - Delete meal plan
  - `GET /api/v1/meal-plans/week` - Get weekly meal plan
- **Files**:
  - `features/meal_plans/models.go`
  - `features/meal_plans/repository.go`
  - `features/meal_plans/service.go`
  - `features/meal_plans/handlers.go`
  - `features/meal_plans/routes.go`

---

### Phase 3: User Preferences & Nutrition
**Priority: MEDIUM** | **Dependencies: Phase 1**

#### 3.1 Family Preferences
- **Tables**: `family_preferences`
- **Endpoints**:
  - `GET /api/v1/preferences` - Get family preferences
  - `POST /api/v1/preferences` - Create/Update preferences
  - `PUT /api/v1/preferences` - Update preferences
- **Files**:
  - `features/preferences/models.go`
  - `features/preferences/repository.go`
  - `features/preferences/service.go`
  - `features/preferences/handlers.go`
  - `features/preferences/routes.go`

#### 3.2 Nutrition Tracking
- **Tables**: `nutrition_data`
- **Endpoints**:
  - `GET /api/v1/nutrition` - Get nutrition data (with date range)
  - `GET /api/v1/nutrition/today` - Get today's nutrition
  - `POST /api/v1/nutrition` - Log nutrition data
  - `PUT /api/v1/nutrition/:id` - Update nutrition data
  - `GET /api/v1/nutrition/stats` - Get nutrition statistics
- **Files**:
  - `features/nutrition/models.go`
  - `features/nutrition/repository.go`
  - `features/nutrition/service.go`
  - `features/nutrition/handlers.go`
  - `features/nutrition/routes.go`

#### 3.3 Price Comparisons
- **Tables**: `price_comparisons`
- **Endpoints**:
  - `GET /api/v1/price-comparisons` - Get price comparisons
  - `GET /api/v1/price-comparisons/:id` - Get specific comparison
  - `POST /api/v1/price-comparisons` - Create price comparison
  - `PUT /api/v1/price-comparisons/:id` - Update comparison
- **Files**:
  - `features/price_comparisons/models.go`
  - `features/price_comparisons/repository.go`
  - `features/price_comparisons/service.go`
  - `features/price_comparisons/handlers.go`
  - `features/price_comparisons/routes.go`

---

### Phase 4: Gamification
**Priority: MEDIUM** | **Dependencies: Phase 1, Phase 2**

#### 4.1 Badges System
- **Tables**: `badges`
- **Endpoints**:
  - `GET /api/v1/badges` - Get user badges
  - `GET /api/v1/badges/available` - Get available badges
  - `POST /api/v1/badges/unlock` - Unlock badge (system)
- **Files**:
  - `features/badges/models.go`
  - `features/badges/repository.go`
  - `features/badges/service.go`
  - `features/badges/handlers.go`
  - `features/badges/routes.go`

#### 4.2 XP & Leveling System
- **Tables**: `user_xp`
- **Endpoints**:
  - `GET /api/v1/xp` - Get user XP and level
  - `POST /api/v1/xp/add` - Add XP (system)
  - `GET /api/v1/xp/leaderboard` - Get leaderboard
- **Files**:
  - `features/xp/models.go`
  - `features/xp/repository.go`
  - `features/xp/service.go`
  - `features/xp/handlers.go`
  - `features/xp/routes.go`

---

### Phase 5: Community Features
**Priority: HIGH** | **Dependencies: Phase 1**

#### 5.1 Community Surplus Posts
- **Tables**: `community_surplus_posts`, `surplus_requests`, `surplus_comments`
- **Endpoints**:
  - `GET /api/v1/community/surplus` - List surplus posts
  - `GET /api/v1/community/surplus/:id` - Get surplus post
  - `POST /api/v1/community/surplus` - Create surplus post
  - `PUT /api/v1/community/surplus/:id` - Update post
  - `DELETE /api/v1/community/surplus/:id` - Delete post
  - `POST /api/v1/community/surplus/:id/request` - Request surplus
  - `GET /api/v1/community/surplus/:id/requests` - Get requests
  - `PUT /api/v1/community/surplus/:id/requests/:requestId` - Approve/decline request
  - `POST /api/v1/community/surplus/:id/comments` - Add comment
  - `GET /api/v1/community/surplus/:id/comments` - Get comments
- **Files**:
  - `features/community/surplus/models.go`
  - `features/community/surplus/repository.go`
  - `features/community/surplus/service.go`
  - `features/community/surplus/handlers.go`
  - `features/community/surplus/routes.go`

#### 5.2 Leftover Items
- **Tables**: `leftover_items`, `leftover_item_claims`
- **Endpoints**:
  - `GET /api/v1/community/leftovers` - List leftover items
  - `GET /api/v1/community/leftovers/:id` - Get leftover item
  - `POST /api/v1/community/leftovers` - Create leftover item
  - `PUT /api/v1/community/leftovers/:id` - Update leftover item
  - `DELETE /api/v1/community/leftovers/:id` - Delete leftover item
  - `POST /api/v1/community/leftovers/:id/claim` - Claim leftover
  - `GET /api/v1/community/leftovers/:id/claims` - Get claims
- **Files**:
  - `features/community/leftovers/models.go`
  - `features/community/leftovers/repository.go`
  - `features/community/leftovers/service.go`
  - `features/community/leftovers/handlers.go`
  - `features/community/leftovers/routes.go`

#### 5.3 Community Kitchen Events
- **Tables**: `community_kitchen_events`
- **Endpoints**:
  - `GET /api/v1/community/kitchen-events` - List events
  - `GET /api/v1/community/kitchen-events/:id` - Get event details
  - `POST /api/v1/community/kitchen-events` - Create event
  - `PUT /api/v1/community/kitchen-events/:id` - Update event
  - `POST /api/v1/community/kitchen-events/:id/volunteer` - Volunteer for event
- **Files**:
  - `features/community/kitchen_events/models.go`
  - `features/community/kitchen_events/repository.go`
  - `features/community/kitchen_events/service.go`
  - `features/community/kitchen_events/handlers.go`
  - `features/community/kitchen_events/routes.go`

#### 5.4 Community Leaderboard & Impact
- **Tables**: `community_leaderboard`, `community_impact`
- **Endpoints**:
  - `GET /api/v1/community/leaderboard` - Get leaderboard
  - `GET /api/v1/community/impact` - Get community impact
  - `GET /api/v1/community/impact/personal` - Get personal contribution
- **Files**:
  - `features/community/leaderboard/models.go`
  - `features/community/leaderboard/repository.go`
  - `features/community/leaderboard/service.go`
  - `features/community/leaderboard/handlers.go`
  - `features/community/leaderboard/routes.go`

#### 5.5 Community Profiles
- **Tables**: `community_profiles`
- **Endpoints**:
  - `GET /api/v1/community/profile` - Get user's community profile
  - `GET /api/v1/community/profile/:username` - Get profile by username
  - `POST /api/v1/community/profile` - Create profile
  - `PUT /api/v1/community/profile` - Update profile
- **Files**:
  - `features/community/profiles/models.go`
  - `features/community/profiles/repository.go`
  - `features/community/profiles/service.go`
  - `features/community/profiles/handlers.go`
  - `features/community/profiles/routes.go`

---

### Phase 6: Restaurant Module
**Priority: MEDIUM** | **Dependencies: Phase 1**

#### 6.1 Restaurant Inventory
- **Tables**: `restaurant_inventory_items`
- **Endpoints**:
  - `GET /api/v1/restaurant/inventory` - List inventory
  - `GET /api/v1/restaurant/inventory/:id` - Get item
  - `POST /api/v1/restaurant/inventory` - Add item
  - `PUT /api/v1/restaurant/inventory/:id` - Update item
  - `DELETE /api/v1/restaurant/inventory/:id` - Delete item
  - `GET /api/v1/restaurant/inventory/expiring` - Get expiring items
- **Files**:
  - `features/restaurant/inventory/models.go`
  - `features/restaurant/inventory/repository.go`
  - `features/restaurant/inventory/service.go`
  - `features/restaurant/inventory/handlers.go`
  - `features/restaurant/inventory/routes.go`

#### 6.2 Restaurant Menu
- **Tables**: `restaurant_menu_items`
- **Endpoints**:
  - `GET /api/v1/restaurant/menu` - List menu items
  - `GET /api/v1/restaurant/menu/:id` - Get menu item
  - `POST /api/v1/restaurant/menu` - Add menu item
  - `PUT /api/v1/restaurant/menu/:id` - Update menu item
  - `DELETE /api/v1/restaurant/menu/:id` - Delete menu item
- **Files**:
  - `features/restaurant/menu/models.go`
  - `features/restaurant/menu/repository.go`
  - `features/restaurant/menu/service.go`
  - `features/restaurant/menu/handlers.go`
  - `features/restaurant/menu/routes.go`

#### 6.3 Restaurant Surplus
- **Tables**: `restaurant_surplus_items`
- **Endpoints**:
  - `GET /api/v1/restaurant/surplus` - List surplus items
  - `POST /api/v1/restaurant/surplus` - Create surplus item
  - `PUT /api/v1/restaurant/surplus/:id` - Update surplus item
  - `PUT /api/v1/restaurant/surplus/:id/assign` - Assign to NGO/kitchen
- **Files**:
  - `features/restaurant/surplus/models.go`
  - `features/restaurant/surplus/repository.go`
  - `features/restaurant/surplus/service.go`
  - `features/restaurant/surplus/handlers.go`
  - `features/restaurant/surplus/routes.go`

#### 6.4 Restaurant Donations & Impact
- **Tables**: `restaurant_donation_logs`, `restaurant_impact_metrics`
- **Endpoints**:
  - `GET /api/v1/restaurant/donations` - List donations
  - `POST /api/v1/restaurant/donations` - Log donation
  - `GET /api/v1/restaurant/impact` - Get impact metrics
- **Files**:
  - `features/restaurant/donations/models.go`
  - `features/restaurant/donations/repository.go`
  - `features/restaurant/donations/service.go`
  - `features/restaurant/donations/handlers.go`
  - `features/restaurant/donations/routes.go`

#### 6.5 Restaurant Staff Management
- **Tables**: `restaurant_staff_tasks`, `restaurant_shift_schedule`
- **Endpoints**:
  - `GET /api/v1/restaurant/tasks` - List tasks
  - `POST /api/v1/restaurant/tasks` - Create task
  - `PUT /api/v1/restaurant/tasks/:id` - Update task
  - `GET /api/v1/restaurant/shifts` - List shifts
  - `POST /api/v1/restaurant/shifts` - Create shift
- **Files**:
  - `features/restaurant/staff/models.go`
  - `features/restaurant/staff/repository.go`
  - `features/restaurant/staff/service.go`
  - `features/restaurant/staff/handlers.go`
  - `features/restaurant/staff/routes.go`

#### 6.6 Restaurant Preferences
- **Tables**: `restaurant_preferences`
- **Endpoints**:
  - `GET /api/v1/restaurant/preferences` - Get preferences
  - `POST /api/v1/restaurant/preferences` - Create/Update preferences
- **Files**:
  - `features/restaurant/preferences/models.go`
  - `features/restaurant/preferences/repository.go`
  - `features/restaurant/preferences/service.go`
  - `features/restaurant/preferences/handlers.go`
  - `features/restaurant/preferences/routes.go`

---

### Phase 7: NGO Module
**Priority: MEDIUM** | **Dependencies: Phase 1**

#### 7.1 NGO Capacity Settings
- **Tables**: `ngo_capacity_settings`
- **Endpoints**:
  - `GET /api/v1/ngo/capacity` - Get capacity settings
  - `POST /api/v1/ngo/capacity` - Create/Update capacity settings
- **Files**:
  - `features/ngo/capacity/models.go`
  - `features/ngo/capacity/repository.go`
  - `features/ngo/capacity/service.go`
  - `features/ngo/capacity/handlers.go`
  - `features/ngo/capacity/routes.go`

#### 7.2 NGO Donation Offers
- **Tables**: `ngo_donation_offers`
- **Endpoints**:
  - `GET /api/v1/ngo/offers` - List donation offers
  - `GET /api/v1/ngo/offers/:id` - Get offer details
  - `PUT /api/v1/ngo/offers/:id/accept` - Accept offer
  - `PUT /api/v1/ngo/offers/:id/decline` - Decline offer
- **Files**:
  - `features/ngo/offers/models.go`
  - `features/ngo/offers/repository.go`
  - `features/ngo/offers/service.go`
  - `features/ngo/offers/handlers.go`
  - `features/ngo/offers/routes.go`

#### 7.3 NGO Pickup Schedules
- **Tables**: `ngo_pickup_schedules`
- **Endpoints**:
  - `GET /api/v1/ngo/pickups` - List pickup schedules
  - `POST /api/v1/ngo/pickups` - Schedule pickup
  - `PUT /api/v1/ngo/pickups/:id` - Update pickup
  - `PUT /api/v1/ngo/pickups/:id/status` - Update pickup status
- **Files**:
  - `features/ngo/pickups/models.go`
  - `features/ngo/pickups/repository.go`
  - `features/ngo/pickups/service.go`
  - `features/ngo/pickups/handlers.go`
  - `features/ngo/pickups/routes.go`

#### 7.4 NGO Donation History
- **Tables**: `ngo_donation_history`
- **Endpoints**:
  - `GET /api/v1/ngo/history` - Get donation history
  - `GET /api/v1/ngo/history/:id` - Get donation details
- **Files**:
  - `features/ngo/history/models.go`
  - `features/ngo/history/repository.go`
  - `features/ngo/history/service.go`
  - `features/ngo/history/handlers.go`
  - `features/ngo/history/routes.go`

#### 7.5 NGO Partner Management
- **Tables**: `ngo_partner_profiles`
- **Endpoints**:
  - `GET /api/v1/ngo/partners` - List partners
  - `GET /api/v1/ngo/partners/:id` - Get partner details
  - `POST /api/v1/ngo/partners` - Add partner
  - `PUT /api/v1/ngo/partners/:id` - Update partner
- **Files**:
  - `features/ngo/partners/models.go`
  - `features/ngo/partners/repository.go`
  - `features/ngo/partners/service.go`
  - `features/ngo/partners/handlers.go`
  - `features/ngo/partners/routes.go`

#### 7.6 NGO Feedback & Impact
- **Tables**: `ngo_feedback_entries`, `ngo_impact_stories`
- **Endpoints**:
  - `GET /api/v1/ngo/feedback` - List feedback
  - `POST /api/v1/ngo/feedback` - Submit feedback
  - `GET /api/v1/ngo/stories` - Get impact stories
  - `POST /api/v1/ngo/stories` - Create impact story
- **Files**:
  - `features/ngo/feedback/models.go`
  - `features/ngo/feedback/repository.go`
  - `features/ngo/feedback/service.go`
  - `features/ngo/feedback/handlers.go`
  - `features/ngo/feedback/routes.go`

---

### Phase 8: Shop Module
**Priority: MEDIUM** | **Dependencies: Phase 1**

#### 8.1 Shop Inventory
- **Tables**: `shop_inventory_items`
- **Endpoints**:
  - `GET /api/v1/shop/inventory` - List inventory
  - `GET /api/v1/shop/inventory/:id` - Get item
  - `POST /api/v1/shop/inventory` - Add item
  - `PUT /api/v1/shop/inventory/:id` - Update item
  - `GET /api/v1/shop/inventory/barcode/:barcode` - Get by barcode
- **Files**:
  - `features/shop/inventory/models.go`
  - `features/shop/inventory/repository.go`
  - `features/shop/inventory/service.go`
  - `features/shop/inventory/handlers.go`
  - `features/shop/inventory/routes.go`

#### 8.2 Shop Price Management
- **Tables**: `shop_price_map_entries`, `shop_discount_suggestions`
- **Endpoints**:
  - `GET /api/v1/shop/prices` - List price mappings
  - `POST /api/v1/shop/prices` - Create price mapping
  - `GET /api/v1/shop/discounts` - Get discount suggestions
  - `POST /api/v1/shop/discounts/:id/apply` - Apply discount
- **Files**:
  - `features/shop/pricing/models.go`
  - `features/shop/pricing/repository.go`
  - `features/shop/pricing/service.go`
  - `features/shop/pricing/handlers.go`
  - `features/shop/pricing/routes.go`

#### 8.3 Shop Surplus Management
- **Tables**: `shop_surplus_items`
- **Endpoints**:
  - `GET /api/v1/shop/surplus` - List surplus items
  - `POST /api/v1/shop/surplus` - Create surplus item
  - `PUT /api/v1/shop/surplus/:id` - Update surplus item
- **Files**:
  - `features/shop/surplus/models.go`
  - `features/shop/surplus/repository.go`
  - `features/shop/surplus/service.go`
  - `features/shop/surplus/handlers.go`
  - `features/shop/surplus/routes.go`

#### 8.4 Shop Analytics
- **Tables**: `shop_analytics_records`
- **Endpoints**:
  - `GET /api/v1/shop/analytics` - Get analytics
  - `GET /api/v1/shop/analytics/waste-trend` - Get waste trend
  - `GET /api/v1/shop/analytics/markdown-recovery` - Get markdown recovery
- **Files**:
  - `features/shop/analytics/models.go`
  - `features/shop/analytics/repository.go`
  - `features/shop/analytics/service.go`
  - `features/shop/analytics/handlers.go`
  - `features/shop/analytics/routes.go`

#### 8.5 Shop Staff Management
- **Tables**: `shop_staff_members`, `shop_staff_tasks`, `shop_shifts`
- **Endpoints**:
  - `GET /api/v1/shop/staff` - List staff
  - `POST /api/v1/shop/staff` - Add staff
  - `GET /api/v1/shop/tasks` - List tasks
  - `POST /api/v1/shop/tasks` - Create task
  - `GET /api/v1/shop/shifts` - List shifts
  - `POST /api/v1/shop/shifts` - Create shift
- **Files**:
  - `features/shop/staff/models.go`
  - `features/shop/staff/repository.go`
  - `features/shop/staff/service.go`
  - `features/shop/staff/handlers.go`
  - `features/shop/staff/routes.go`

#### 8.6 Shop Profile
- **Tables**: `shop_profiles`
- **Endpoints**:
  - `GET /api/v1/shop/profile` - Get shop profile
  - `POST /api/v1/shop/profile` - Create/Update profile
- **Files**:
  - `features/shop/profile/models.go`
  - `features/shop/profile/repository.go`
  - `features/shop/profile/service.go`
  - `features/shop/profile/handlers.go`
  - `features/shop/profile/routes.go`

---

### Phase 9: Supporting Features
**Priority: LOW** | **Dependencies: Phase 1**

#### 9.1 File Uploads
- **Tables**: `uploads`
- **Endpoints**:
  - `POST /api/v1/uploads` - Upload file
  - `GET /api/v1/uploads/:id` - Get file
  - `DELETE /api/v1/uploads/:id` - Delete file
- **Files**:
  - `features/uploads/models.go`
  - `features/uploads/repository.go`
  - `features/uploads/service.go`
  - `features/uploads/handlers.go`
  - `features/uploads/routes.go`

#### 9.2 Resources
- **Tables**: `resources`
- **Endpoints**:
  - `GET /api/v1/resources` - List resources
  - `GET /api/v1/resources/:id` - Get resource
  - `POST /api/v1/resources` - Create resource (admin)
  - `PUT /api/v1/resources/:id` - Update resource (admin)
  - `DELETE /api/v1/resources/:id` - Delete resource (admin)
- **Files**:
  - `features/resources/models.go`
  - `features/resources/repository.go`
  - `features/resources/service.go`
  - `features/resources/handlers.go`
  - `features/resources/routes.go`

#### 9.3 Notifications
- **Tables**: `community_notifications`, `ngo_notifications`
- **Endpoints**:
  - `GET /api/v1/notifications` - List notifications
  - `GET /api/v1/notifications/unread` - Get unread notifications
  - `PUT /api/v1/notifications/:id/read` - Mark as read
  - `PUT /api/v1/notifications/read-all` - Mark all as read
- **Files**:
  - `features/notifications/models.go`
  - `features/notifications/repository.go`
  - `features/notifications/service.go`
  - `features/notifications/handlers.go`
  - `features/notifications/routes.go`

---

## Project Structure

```
foodlink_backend/
├── main.go                          # Application entry point
├── config/                          # Configuration
│   └── config.go
├── database/                        # Database layer
│   ├── connection.go               # DB connection pool
│   ├── migrations/                 # Migration files
│   └── schema.go                   # Schema utilities
├── middleware/                     # HTTP middleware
│   ├── auth.go                     # Authentication middleware
│   ├── cors.go                     # CORS middleware
│   ├── logging.go                  # Logging middleware
│   └── error_handler.go           # Error handling middleware
├── utils/                          # Utility functions
│   ├── response.go                 # Standardized HTTP responses
│   ├── validator.go                # Validation helpers
│   └── jwt.go                      # JWT utilities
├── errors/                         # Custom error types
│   └── errors.go
├── features/                       # Feature modules
│   ├── auth/                       # Authentication
│   │   ├── models.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handlers.go
│   │   └── routes.go
│   ├── inventory/                  # Inventory management
│   ├── consumption/                # Consumption tracking
│   ├── shopping_list/             # Shopping lists
│   ├── meal_plans/                 # Meal planning
│   ├── food_items/                 # Food items reference
│   ├── preferences/                # Family preferences
│   ├── nutrition/                  # Nutrition tracking
│   ├── price_comparisons/          # Price comparisons
│   ├── badges/                      # Badges system
│   ├── xp/                         # XP & leveling
│   ├── community/                  # Community features
│   │   ├── surplus/
│   │   ├── leftovers/
│   │   ├── kitchen_events/
│   │   ├── leaderboard/
│   │   └── profiles/
│   ├── restaurant/                 # Restaurant module
│   │   ├── inventory/
│   │   ├── menu/
│   │   ├── surplus/
│   │   ├── donations/
│   │   ├── staff/
│   │   └── preferences/
│   ├── ngo/                        # NGO module
│   │   ├── capacity/
│   │   ├── offers/
│   │   ├── pickups/
│   │   ├── history/
│   │   ├── partners/
│   │   └── feedback/
│   ├── shop/                       # Shop module
│   │   ├── inventory/
│   │   ├── pricing/
│   │   ├── surplus/
│   │   ├── analytics/
│   │   ├── staff/
│   │   └── profile/
│   ├── uploads/                    # File uploads
│   ├── resources/                  # Resources
│   └── notifications/              # Notifications
├── routes/                         # Main route setup
│   └── routes.go
├── docs/                           # Swagger docs (auto-generated)
├── schema.sql                      # Database schema
├── go.mod
├── go.sum
└── README.md
```

---

## Implementation Guidelines

### 1. Code Structure Pattern
Each feature follows this structure:
- **models.go**: Data models/structs
- **repository.go**: Database operations (queries)
- **service.go**: Business logic
- **handlers.go**: HTTP request handlers
- **routes.go**: Route definitions

### 2. Naming Conventions
- Packages: lowercase, singular (e.g., `inventory`, `auth`)
- Files: snake_case (e.g., `meal_plans.go`)
- Types: PascalCase (e.g., `InventoryItem`)
- Functions: PascalCase for exported, camelCase for private
- Variables: camelCase

### 3. Error Handling
- Use custom error types in `errors/` package
- Return appropriate HTTP status codes
- Provide meaningful error messages
- Log errors appropriately

### 4. Database
- Use prepared statements for queries
- Handle transactions properly
- Use connection pooling
- Implement proper migrations

### 5. Testing
- Unit tests for services
- Integration tests for handlers
- Repository tests with test database
- Test coverage > 70%

### 6. Documentation
- Swagger annotations for all endpoints
- Code comments for complex logic
- README for each feature module

---

## Implementation Order

### Sprint 1 (Week 1-2): Foundation
- Phase 0: Database setup & infrastructure
- Phase 1: Authentication & User Management

### Sprint 2 (Week 3-4): Core Features
- Phase 2: Core Family Features (Food Items, Inventory, Consumption, Shopping Lists, Meal Plans)

### Sprint 3 (Week 5-6): User Experience
- Phase 3: Preferences & Nutrition
- Phase 4: Gamification

### Sprint 4 (Week 7-8): Community
- Phase 5: Community Features

### Sprint 5 (Week 9-10): Business Modules
- Phase 6: Restaurant Module
- Phase 7: NGO Module

### Sprint 6 (Week 11-12): Shop & Polish
- Phase 8: Shop Module
- Phase 9: Supporting Features

---

## Dependencies & Prerequisites

### Required Packages
- Database: `github.com/lib/pq` (PostgreSQL driver)
- JWT: `github.com/golang-jwt/jwt/v5`
- Password: `golang.org/x/crypto/bcrypt`
- Environment: `github.com/joho/godotenv`
- Swagger: `github.com/swaggo/swag`, `github.com/swaggo/http-swagger`
- Validation: `github.com/go-playground/validator/v10`
- UUID: `github.com/google/uuid`

### Database Requirements
- PostgreSQL 12+
- UUID extension enabled
- SSL support (for Neon)

---

## Notes

1. **Role-Based Access Control**: Implement role checks in middleware
2. **Rate Limiting**: Add rate limiting middleware for production
3. **Caching**: Consider Redis for caching frequently accessed data
4. **Background Jobs**: Use goroutines/channels for async tasks
5. **Logging**: Structured logging with levels (info, warn, error)
6. **Monitoring**: Add health checks and metrics endpoints
7. **Security**: Input validation, SQL injection prevention, XSS protection

---

## Next Steps

1. Set up database connection and migrations
2. Implement authentication feature
3. Create base repository pattern
4. Implement first feature (Food Items or Inventory)
5. Set up CI/CD pipeline
6. Write tests for implemented features
