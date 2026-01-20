# Foodlink Backend - Features Overview

Quick reference guide for all features in the Foodlink backend.

## Feature List

### 🔐 Authentication & User Management
- User registration and login
- JWT token management
- User profile management
- Role-based access control (family, restaurant, shop, ngo, admin)

### 📦 Core Family Features
- **Food Items**: Reference data for food items
- **Inventory Management**: Track household inventory
- **Consumption Tracking**: Log food consumption and waste
- **Shopping Lists**: Create and manage shopping lists
- **Meal Planning**: Plan meals and generate shopping lists

### 🥗 Nutrition & Preferences
- **Family Preferences**: Store household preferences and dietary restrictions
- **Nutrition Tracking**: Daily nutrition data and statistics
- **Price Comparisons**: Compare prices across stores

### 🎮 Gamification
- **Badges System**: Unlock badges for achievements
- **XP & Leveling**: Earn XP and level up
- **Leaderboards**: Community leaderboards

### 👥 Community Features
- **Surplus Posts**: Share surplus food with community
- **Leftover Items**: Share leftover meals
- **Kitchen Events**: Community kitchen events and volunteering
- **Leaderboard**: Community impact leaderboards
- **Community Profiles**: User community profiles

### 🍽️ Restaurant Module
- **Inventory Management**: Track restaurant inventory
- **Menu Management**: Manage menu items
- **Surplus Management**: Manage surplus food items
- **Donation Logging**: Track donations to NGOs/kitchens
- **Impact Metrics**: Track sustainability impact
- **Staff Management**: Tasks and shift scheduling

### 🏢 NGO Module
- **Capacity Settings**: Configure NGO capacity and preferences
- **Donation Offers**: View and manage donation offers
- **Pickup Scheduling**: Schedule and track pickups
- **Donation History**: Historical donation records
- **Partner Management**: Manage partner relationships
- **Feedback & Impact**: Feedback system and impact stories

### 🏪 Shop Module
- **Inventory Management**: Track shop inventory with barcodes
- **Price Management**: Price mapping and discount suggestions
- **Surplus Management**: Manage surplus items
- **Analytics**: Waste reduction and markdown recovery analytics
- **Staff Management**: Staff, tasks, and shifts
- **Shop Profile**: Shop information and preferences

### 📎 Supporting Features
- **File Uploads**: Handle file uploads (images, documents)
- **Resources**: Educational resources
- **Notifications**: User notifications system

---

## Feature Dependencies

```
Phase 0 (Foundation)
    ↓
Phase 1 (Authentication)
    ↓
Phase 2 (Core Family Features)
    ↓
Phase 3 (Preferences & Nutrition)
    ↓
Phase 4 (Gamification)
    ↓
Phase 5 (Community)
    ↓
Phase 6 (Restaurant)
    ↓
Phase 7 (NGO)
    ↓
Phase 8 (Shop)
    ↓
Phase 9 (Supporting Features)
```

---

## API Endpoint Structure

```
/api/v1/
├── /auth/                    # Authentication
├── /food-items/             # Food items reference
├── /inventory/              # Inventory management
├── /consumption/            # Consumption tracking
├── /shopping-list/          # Shopping lists
├── /meal-plans/             # Meal planning
├── /preferences/            # Family preferences
├── /nutrition/               # Nutrition tracking
├── /price-comparisons/      # Price comparisons
├── /badges/                 # Badges system
├── /xp/                     # XP & leveling
├── /community/               # Community features
│   ├── /surplus/
│   ├── /leftovers/
│   ├── /kitchen-events/
│   ├── /leaderboard/
│   └── /profile/
├── /restaurant/             # Restaurant module
│   ├── /inventory/
│   ├── /menu/
│   ├── /surplus/
│   ├── /donations/
│   ├── /staff/
│   └── /preferences/
├── /ngo/                    # NGO module
│   ├── /capacity/
│   ├── /offers/
│   ├── /pickups/
│   ├── /history/
│   ├── /partners/
│   └── /feedback/
├── /shop/                   # Shop module
│   ├── /inventory/
│   ├── /pricing/
│   ├── /surplus/
│   ├── /analytics/
│   ├── /staff/
│   └── /profile/
├── /uploads/                # File uploads
├── /resources/              # Resources
└── /notifications/          # Notifications
```

---

## Database Tables by Feature

### Authentication
- `users`

### Core Family
- `food_items`
- `inventory_items`
- `consumption_logs`
- `shopping_list_items`
- `meal_plans`

### Preferences & Nutrition
- `family_preferences`
- `nutrition_data`
- `price_comparisons`

### Gamification
- `badges`
- `user_xp`

### Community
- `community_surplus_posts`
- `surplus_requests`
- `surplus_comments`
- `leftover_items`
- `leftover_item_claims`
- `community_kitchen_events`
- `community_leaderboard`
- `community_impact`
- `community_notifications`
- `community_profiles`

### Restaurant
- `restaurant_inventory_items`
- `restaurant_menu_items`
- `restaurant_surplus_items`
- `restaurant_donation_logs`
- `restaurant_impact_metrics`
- `restaurant_staff_tasks`
- `restaurant_shift_schedule`
- `restaurant_preferences`

### NGO
- `ngo_capacity_settings`
- `ngo_donation_offers`
- `ngo_pickup_schedules`
- `ngo_donation_history`
- `ngo_partner_profiles`
- `ngo_feedback_entries`
- `ngo_impact_stories`
- `ngo_notifications`

### Shop
- `shop_inventory_items`
- `shop_price_map_entries`
- `shop_discount_suggestions`
- `shop_surplus_items`
- `shop_analytics_records`
- `shop_staff_members`
- `shop_staff_tasks`
- `shop_shifts`
- `shop_profiles`

### Supporting
- `uploads`
- `resources`

---

## User Roles

- **family**: Regular household users
- **restaurant**: Restaurant owners/managers
- **shop**: Shop owners/managers
- **ngo**: NGO administrators
- **admin**: System administrators

---

## Status Codes

- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict
- `500 Internal Server Error`: Server error
