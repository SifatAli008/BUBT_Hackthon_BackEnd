-- FoodFlow/FoodLink Database Schema for Supabase PostgreSQL
-- Run this SQL in your Supabase SQL Editor
-- Based on all frontend data models

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable PostGIS for geolocation (if needed)
-- CREATE EXTENSION IF NOT EXISTS "postgis";

-- ============================================================================
-- CORE TABLES
-- ============================================================================

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    household_id UUID,
    role VARCHAR(50) DEFAULT 'family' CHECK (role IN ('family', 'restaurant', 'shop', 'ngo', 'admin')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Food items table (reference data)
CREATE TABLE IF NOT EXISTS food_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    typical_expiry_days INTEGER NOT NULL,
    storage_tips TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Inventory items table
CREATE TABLE IF NOT EXISTS inventory_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(50),
    expiry_date TIMESTAMP WITH TIME ZONE,
    category VARCHAR(100),
    location VARCHAR(100),
    food_item_id UUID REFERENCES food_items(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Consumption logs table
CREATE TABLE IF NOT EXISTS consumption_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    inventory_item_id UUID REFERENCES inventory_items(id) ON DELETE SET NULL,
    food_name VARCHAR(255) NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(50),
    category VARCHAR(100),
    consumed_at TIMESTAMP WITH TIME ZONE NOT NULL,
    was_wasted BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Shopping list items table
CREATE TABLE IF NOT EXISTS shopping_list_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    household_id UUID,
    name VARCHAR(255) NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(50),
    category VARCHAR(100),
    priority VARCHAR(20) DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high')),
    purchased BOOLEAN DEFAULT FALSE,
    purchased_at TIMESTAMP WITH TIME ZONE,
    estimated_price DECIMAL(10, 2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Meal plans table
CREATE TABLE IF NOT EXISTS meal_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    household_id UUID,
    date DATE NOT NULL,
    meal_type VARCHAR(20) NOT NULL CHECK (meal_type IN ('breakfast', 'lunch', 'dinner', 'snack')),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    ingredients TEXT[],
    servings INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Resources table
CREATE TABLE IF NOT EXISTS resources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    url TEXT,
    tags TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Uploads table
CREATE TABLE IF NOT EXISTS uploads (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(100) NOT NULL,
    file_size BIGINT NOT NULL,
    data TEXT, -- base64 encoded
    associated_type VARCHAR(50) CHECK (associated_type IN ('inventory', 'log', 'profile')),
    associated_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- GAMIFICATION & USER PROGRESS
-- ============================================================================

-- Badges table
CREATE TABLE IF NOT EXISTS badges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    badge_id VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(255),
    unlocked_at TIMESTAMP WITH TIME ZONE NOT NULL,
    xp_reward INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, badge_id)
);

-- User XP table
CREATE TABLE IF NOT EXISTS user_xp (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    total_xp INTEGER DEFAULT 0,
    level INTEGER DEFAULT 1,
    current_level_xp INTEGER DEFAULT 0,
    next_level_xp INTEGER DEFAULT 100,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- FAMILY & PREFERENCES
-- ============================================================================

-- Family preferences table (using JSONB for complex nested data)
CREATE TABLE IF NOT EXISTS family_preferences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    household_id UUID NOT NULL UNIQUE,
    household_size INTEGER NOT NULL,
    age_groups JSONB, -- {child: number, adult: number, senior: number}
    cooking_frequency VARCHAR(50) CHECK (cooking_frequency IN ('daily', 'few-times-week', 'weekly', 'occasional')),
    eating_schedule JSONB, -- {breakfast: string, lunch: string, dinner: string}
    dietary_type VARCHAR(50) CHECK (dietary_type IN ('vegan', 'vegetarian', 'halal', 'keto', 'low-sodium', 'general')),
    dietary_restrictions TEXT[],
    allergies TEXT[],
    health_conditions TEXT[],
    weekly_budget DECIMAL(10, 2),
    budget_range JSONB, -- {min: number, max: number}
    preferred_stores TEXT[],
    price_sensitivity VARCHAR(20) CHECK (price_sensitivity IN ('low', 'medium', 'high')),
    preferred_cuisines TEXT[],
    meal_prep_preference VARCHAR(50) CHECK (meal_prep_preference IN ('quick', 'diverse', 'budget', 'high-protein')),
    waste_sensitivity_level VARCHAR(20) CHECK (waste_sensitivity_level IN ('low', 'medium', 'high')),
    sustainability_preference VARCHAR(20) CHECK (sustainability_preference IN ('minimal', 'moderate', 'high')),
    leftover_comfort_level VARCHAR(20) CHECK (leftover_comfort_level IN ('low', 'medium', 'high')),
    daily_calories INTEGER,
    macro_goal JSONB, -- {protein: number, carbs: number, fats: number}
    vitamins_focus TEXT[],
    avoid_excess TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Nutrition data table
CREATE TABLE IF NOT EXISTS nutrition_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    calories DECIMAL(10, 2) DEFAULT 0,
    protein DECIMAL(10, 2) DEFAULT 0, -- grams
    carbs DECIMAL(10, 2) DEFAULT 0, -- grams
    fats DECIMAL(10, 2) DEFAULT 0, -- grams
    fiber DECIMAL(10, 2) DEFAULT 0, -- grams
    sugar DECIMAL(10, 2) DEFAULT 0, -- grams
    sodium DECIMAL(10, 2) DEFAULT 0, -- mg
    vitamin_a DECIMAL(10, 2) DEFAULT 0, -- IU
    vitamin_b DECIMAL(10, 2) DEFAULT 0, -- mg
    vitamin_c DECIMAL(10, 2) DEFAULT 0, -- mg
    vitamin_d DECIMAL(10, 2) DEFAULT 0, -- IU
    iron DECIMAL(10, 2) DEFAULT 0, -- mg
    calcium DECIMAL(10, 2) DEFAULT 0, -- mg
    nutrition_score INTEGER, -- 0-100
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, date)
);

-- Price comparisons table (using JSONB for stores array)
CREATE TABLE IF NOT EXISTS price_comparisons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_name VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    stores JSONB NOT NULL, -- [{storeName, price, unit, available}]
    best_price JSONB NOT NULL, -- {storeName, price}
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- COMMUNITY FEATURES
-- ============================================================================

-- Community surplus posts table
CREATE TABLE IF NOT EXISTS community_surplus_posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_name VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(100) NOT NULL,
    tags TEXT[],
    quantity DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    pickup_window JSONB NOT NULL, -- {start: string, end: string}
    pickup_location TEXT NOT NULL,
    distance_km DECIMAL(10, 2),
    image TEXT,
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'claimed', 'expired')),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Surplus requests table
CREATE TABLE IF NOT EXISTS surplus_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL REFERENCES community_surplus_posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_name VARCHAR(255) NOT NULL,
    message TEXT,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'declined')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Surplus comments table
CREATE TABLE IF NOT EXISTS surplus_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL REFERENCES community_surplus_posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_name VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Leftover items table
CREATE TABLE IF NOT EXISTS leftover_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_name VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    dish_name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    portions INTEGER NOT NULL,
    distance_km DECIMAL(10, 2) NOT NULL,
    dietary_tags TEXT[],
    allergens TEXT[],
    pickup_window TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'claimed')),
    image TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Leftover item claims table
CREATE TABLE IF NOT EXISTS leftover_item_claims (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    leftover_item_id UUID NOT NULL REFERENCES leftover_items(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_name VARCHAR(255) NOT NULL,
    message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Community kitchen events table
CREATE TABLE IF NOT EXISTS community_kitchen_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    date DATE NOT NULL,
    time TIME NOT NULL,
    location TEXT NOT NULL,
    tags TEXT[],
    volunteers_needed INTEGER DEFAULT 0,
    volunteers JSONB, -- [{id, userId, name, role, avatarUrl}]
    food_saved_kg DECIMAL(10, 2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'upcoming' CHECK (status IN ('upcoming', 'in-progress', 'completed')),
    image TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Community leaderboard table
CREATE TABLE IF NOT EXISTS community_leaderboard (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(50) NOT NULL CHECK (type IN ('top-sharers', 'zero-waste', 'volunteer-stars', 'building-impact', 'weekly-xp')),
    entries JSONB NOT NULL, -- [{id, name, household, value, unit, badge, trend}]
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Community impact table
CREATE TABLE IF NOT EXISTS community_impact (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    total_surplus_kg DECIMAL(10, 2) DEFAULT 0,
    donations INTEGER DEFAULT 0,
    co2_prevented_kg DECIMAL(10, 2) DEFAULT 0,
    water_saved_liters DECIMAL(10, 2) DEFAULT 0,
    meals_provided INTEGER DEFAULT 0,
    weekly_trend JSONB, -- [{label, value}]
    personal_contribution JSONB, -- [{label, value, unit}]
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Community notifications table
CREATE TABLE IF NOT EXISTS community_notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('claim', 'volunteer', 'announcement', 'surplus', 'reminder')),
    read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Community profiles table
CREATE TABLE IF NOT EXISTS community_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    username VARCHAR(255) NOT NULL UNIQUE,
    avatar_url TEXT,
    community_role VARCHAR(50) DEFAULT 'member' CHECK (community_role IN ('member', 'champion', 'organizer')),
    bio TEXT,
    preferred_items TEXT[],
    avoid_items TEXT[],
    dietary_restrictions TEXT[],
    allergens TEXT[],
    accepts_hot_meals BOOLEAN DEFAULT FALSE,
    distance_preference VARCHAR(10) CHECK (distance_preference IN ('1km', '3km', '5km', 'any')),
    visibility VARCHAR(20) DEFAULT 'public' CHECK (visibility IN ('public', 'community', 'private')),
    notifications_enabled BOOLEAN DEFAULT TRUE,
    notify_on_claim BOOLEAN DEFAULT TRUE,
    notify_on_messages BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- RESTAURANT MODULE
-- ============================================================================

-- Restaurant inventory items table
CREATE TABLE IF NOT EXISTS restaurant_inventory_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    category VARCHAR(100) NOT NULL,
    expiry_date TIMESTAMP WITH TIME ZONE NOT NULL,
    storage_type VARCHAR(20) NOT NULL CHECK (storage_type IN ('fresh', 'chilled', 'frozen', 'dry')),
    batch_code VARCHAR(100),
    alert_tags TEXT[],
    status VARCHAR(20) DEFAULT 'normal' CHECK (status IN ('normal', 'expiring', 'overstocked')),
    invoice_image TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Restaurant menu items table
CREATE TABLE IF NOT EXISTS restaurant_menu_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    ingredients JSONB NOT NULL, -- [{name, quantity}]
    predicted_waste_score VARCHAR(20) CHECK (predicted_waste_score IN ('low', 'medium', 'high')),
    price DECIMAL(10, 2) NOT NULL,
    margin DECIMAL(10, 2) NOT NULL,
    suggestions TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Restaurant surplus items table
CREATE TABLE IF NOT EXISTS restaurant_surplus_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    category VARCHAR(100) NOT NULL,
    storage_type VARCHAR(20) NOT NULL CHECK (storage_type IN ('fresh', 'chilled', 'frozen')),
    pickup_window JSONB NOT NULL, -- {start: string, end: string}
    tags TEXT[],
    image TEXT,
    assigned_to VARCHAR(50) CHECK (assigned_to IN ('ngo', 'kitchen')),
    recipient_name VARCHAR(255),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'picked-up', 'expired')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Restaurant donation logs table
CREATE TABLE IF NOT EXISTS restaurant_donation_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    recipient_type VARCHAR(50) NOT NULL CHECK (recipient_type IN ('ngo', 'community-kitchen')),
    recipient_name VARCHAR(255) NOT NULL,
    items TEXT NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    meals_provided INTEGER DEFAULT 0,
    co2_saved_kg DECIMAL(10, 2) DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Restaurant impact metrics table
CREATE TABLE IF NOT EXISTS restaurant_impact_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    waste_prevented_kg DECIMAL(10, 2) DEFAULT 0,
    surplus_donation_rate DECIMAL(5, 2) DEFAULT 0,
    water_saved_liters DECIMAL(10, 2) DEFAULT 0,
    co2_prevented_kg DECIMAL(10, 2) DEFAULT 0,
    sustainability_score INTEGER DEFAULT 0,
    weekly_trend JSONB, -- [{label, value}]
    monthly_trend JSONB, -- [{label, value}]
    category_breakdown JSONB, -- [{category, wasteKg}]
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Restaurant staff tasks table
CREATE TABLE IF NOT EXISTS restaurant_staff_tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    assignee VARCHAR(255) NOT NULL,
    shift VARCHAR(100) NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    priority VARCHAR(20) DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Restaurant shift schedule table
CREATE TABLE IF NOT EXISTS restaurant_shift_schedule (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(100) NOT NULL,
    staff VARCHAR(255) NOT NULL,
    time VARCHAR(100) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Restaurant preferences table
CREATE TABLE IF NOT EXISTS restaurant_preferences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    cuisine_type VARCHAR(100),
    operating_hours TEXT,
    donation_preferences TEXT[],
    storage_capabilities TEXT[],
    staff_roles TEXT[],
    notifications_enabled BOOLEAN DEFAULT TRUE,
    notify_on_pickup BOOLEAN DEFAULT TRUE,
    notify_on_expiry BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- NGO MODULE
-- ============================================================================

-- NGO capacity settings table
CREATE TABLE IF NOT EXISTS ngo_capacity_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    org_name VARCHAR(255) NOT NULL,
    location TEXT NOT NULL,
    geo_point JSONB, -- {lat: number, lng: number}
    manager_name VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(50) NOT NULL,
    contact_email VARCHAR(255),
    preferred_food_types TEXT[] CHECK (preferred_food_types <@ ARRAY['cooked', 'raw', 'produce', 'bakery', 'protein']),
    restricted_items TEXT[],
    storage_types TEXT[] CHECK (storage_types <@ ARRAY['refrigerated', 'frozen', 'dry']),
    safety_rules TEXT[],
    policy_notes TEXT,
    pickup_window JSONB NOT NULL, -- {start: string, end: string}
    daily_capacity_kg DECIMAL(10, 2) NOT NULL,
    refrigerated_capacity_kg DECIMAL(10, 2) DEFAULT 0,
    dry_capacity_kg DECIMAL(10, 2) DEFAULT 0,
    current_utilization_kg DECIMAL(10, 2) DEFAULT 0,
    xp_points INTEGER DEFAULT 0,
    level INTEGER DEFAULT 1,
    level_progress_pct DECIMAL(5, 2) DEFAULT 0,
    auto_acceptance JSONB, -- {allowPork: boolean, rejectExpired: boolean, temperatureChecks: boolean}
    preferred_pickup_radius_km DECIMAL(10, 2) DEFAULT 5.0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- NGO donation offers table
CREATE TABLE IF NOT EXISTS ngo_donation_offers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ngo_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    donor_name VARCHAR(255) NOT NULL,
    donor_type VARCHAR(50) NOT NULL CHECK (donor_type IN ('building', 'restaurant', 'household', 'kitchen')),
    partner_id UUID,
    distance_km DECIMAL(10, 2) NOT NULL,
    location_label TEXT NOT NULL,
    geo_point JSONB, -- {lat: number, lng: number}
    offer_title VARCHAR(255) NOT NULL,
    items JSONB NOT NULL, -- [{name, quantity, unit, type, temperature}]
    weight_kg DECIMAL(10, 2) NOT NULL,
    meals_estimated INTEGER DEFAULT 0,
    freshness_score INTEGER DEFAULT 0,
    pickup_window JSONB NOT NULL, -- {start: string, end: string}
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    urgency_level VARCHAR(20) DEFAULT 'medium' CHECK (urgency_level IN ('low', 'medium', 'high')),
    dietary_notes TEXT,
    safety_flags TEXT[],
    contact JSONB NOT NULL, -- {name, phone, email, channel}
    images TEXT[],
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined', 'scheduled', 'completed')),
    match_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- NGO pickup schedules table
CREATE TABLE IF NOT EXISTS ngo_pickup_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    offer_id UUID NOT NULL REFERENCES ngo_donation_offers(id) ON DELETE CASCADE,
    route_id UUID,
    scheduled_for TIMESTAMP WITH TIME ZONE NOT NULL,
    eta_minutes INTEGER,
    volunteer_name VARCHAR(255) NOT NULL,
    volunteer_contact VARCHAR(255) NOT NULL,
    vehicle_type VARCHAR(20) CHECK (vehicle_type IN ('van', 'bike', 'car', 'on-foot')),
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'en-route', 'picked-up', 'delivered', 'failed')),
    checkpoints JSONB, -- [{label, timestamp, status, note}]
    reminders JSONB, -- [{time, type, delivered}]
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- NGO donation history table
CREATE TABLE IF NOT EXISTS ngo_donation_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ngo_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    offer_id UUID REFERENCES ngo_donation_offers(id) ON DELETE SET NULL,
    donor_name VARCHAR(255) NOT NULL,
    donor_type VARCHAR(50) NOT NULL,
    items_summary TEXT NOT NULL,
    weight_kg DECIMAL(10, 2) NOT NULL,
    meals_provided INTEGER DEFAULT 0,
    co2_prevented_kg DECIMAL(10, 2) DEFAULT 0,
    beneficiaries INTEGER DEFAULT 0,
    pickup_time TIMESTAMP WITH TIME ZONE NOT NULL,
    delivered_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'delivered' CHECK (status IN ('delivered', 'partial', 'redirected', 'cancelled')),
    tags TEXT[],
    photo TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- NGO partner profiles table
CREATE TABLE IF NOT EXISTS ngo_partner_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ngo_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('community-kitchen', 'building', 'restaurant', 'ngo')),
    location TEXT NOT NULL,
    distance_km DECIMAL(10, 2),
    contact_name VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(50) NOT NULL,
    contact_email VARCHAR(255),
    operating_hours TEXT,
    acceptance_rate DECIMAL(5, 2) DEFAULT 0,
    last_donation_at TIMESTAMP WITH TIME ZONE,
    avg_donation_kg DECIMAL(10, 2) DEFAULT 0,
    storage_capabilities TEXT[],
    notes TEXT,
    avatar TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- NGO feedback entries table
CREATE TABLE IF NOT EXISTS ngo_feedback_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ngo_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_name VARCHAR(255) NOT NULL,
    partner_name VARCHAR(255) NOT NULL,
    delivery_date DATE NOT NULL,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    comment TEXT NOT NULL,
    tags TEXT[] CHECK (tags <@ ARRAY['quality', 'temperature', 'packaging', 'late', 'positive']),
    photo TEXT,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'acknowledged', 'resolved')),
    corrective_action TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- NGO impact stories table
CREATE TABLE IF NOT EXISTS ngo_impact_stories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ngo_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    excerpt TEXT NOT NULL,
    beneficiary_type VARCHAR(100) NOT NULL,
    image TEXT,
    metrics JSONB NOT NULL, -- {meals, families, smiles}
    published_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- NGO notifications table
CREATE TABLE IF NOT EXISTS ngo_notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ngo_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('urgent-offer', 'pickup', 'volunteer', 'feedback', 'message')),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    related_entity_id UUID,
    severity VARCHAR(20) DEFAULT 'info' CHECK (severity IN ('info', 'warning', 'critical')),
    read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- SHOP MODULE
-- ============================================================================

-- Shop inventory items table
CREATE TABLE IF NOT EXISTS shop_inventory_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    barcode VARCHAR(255) NOT NULL,
    stock_quantity DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    cost DECIMAL(10, 2) NOT NULL,
    expiry_date TIMESTAMP WITH TIME ZONE NOT NULL,
    storage_type VARCHAR(20) NOT NULL CHECK (storage_type IN ('frozen', 'chilled', 'ambient')),
    shelf_location VARCHAR(100),
    image_data TEXT,
    markdown_status VARCHAR(20) DEFAULT 'none' CHECK (markdown_status IN ('none', 'scheduled', 'active')),
    surplus_eligible BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Shop price map entries table
CREATE TABLE IF NOT EXISTS shop_price_map_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sku_id UUID NOT NULL REFERENCES shop_inventory_items(id) ON DELETE CASCADE,
    sku_name VARCHAR(255) NOT NULL,
    old_price DECIMAL(10, 2) NOT NULL,
    new_price DECIMAL(10, 2) NOT NULL,
    method VARCHAR(20) NOT NULL CHECK (method IN ('percentage', 'fixed')),
    change_value DECIMAL(10, 2) NOT NULL,
    effective_at TIMESTAMP WITH TIME ZONE NOT NULL,
    scheduled_by VARCHAR(255) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Shop discount suggestions table
CREATE TABLE IF NOT EXISTS shop_discount_suggestions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sku_id UUID NOT NULL REFERENCES shop_inventory_items(id) ON DELETE CASCADE,
    sku_name VARCHAR(255) NOT NULL,
    reason TEXT NOT NULL,
    suggested_discount_pct DECIMAL(5, 2) NOT NULL,
    predicted_sell_through DECIMAL(5, 2) DEFAULT 0,
    urgency VARCHAR(20) DEFAULT 'medium' CHECK (urgency IN ('low', 'medium', 'high')),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Shop surplus items table
CREATE TABLE IF NOT EXISTS shop_surplus_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sku_name VARCHAR(255) NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    expiry_window_start TIMESTAMP WITH TIME ZONE NOT NULL,
    expiry_window_end TIMESTAMP WITH TIME ZONE NOT NULL,
    condition VARCHAR(20) NOT NULL CHECK (condition IN ('fresh', 'near-expiry')),
    destination_type VARCHAR(50) CHECK (destination_type IN ('ngo', 'community-kitchen')),
    destination_name VARCHAR(255),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'picked', 'expired')),
    pickup_time TIMESTAMP WITH TIME ZONE,
    reminder_sent BOOLEAN DEFAULT FALSE,
    image_data TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Shop analytics records table
CREATE TABLE IF NOT EXISTS shop_analytics_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    waste_reduction_trend JSONB, -- [{label, value}]
    markdown_recovery_trend JSONB, -- [{label, value}]
    waste_by_category JSONB, -- [{category, value}]
    surplus_vs_sold JSONB, -- [{label, value}]
    expired_per_day JSONB, -- [{label, value}]
    total_co2_prevented DECIMAL(10, 2) DEFAULT 0,
    meals_donated INTEGER DEFAULT 0,
    waste_reduction_percent DECIMAL(5, 2) DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Shop staff members table
CREATE TABLE IF NOT EXISTS shop_staff_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(100) NOT NULL,
    shift VARCHAR(100) NOT NULL,
    contact VARCHAR(255) NOT NULL,
    avatar TEXT,
    responsibilities TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Shop staff tasks table
CREATE TABLE IF NOT EXISTS shop_staff_tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    assignee_id UUID REFERENCES shop_staff_members(id) ON DELETE SET NULL,
    due TIMESTAMP WITH TIME ZONE NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    category VARCHAR(50) CHECK (category IN ('expiry', 'pricing', 'surplus', 'cleaning')),
    priority VARCHAR(20) DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Shop shifts table
CREATE TABLE IF NOT EXISTS shop_shifts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    staff_id UUID NOT NULL REFERENCES shop_staff_members(id) ON DELETE CASCADE,
    day VARCHAR(20) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    station VARCHAR(100) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Shop profile table
CREATE TABLE IF NOT EXISTS shop_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    store_name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    contact_number VARCHAR(50) NOT NULL,
    manager_name VARCHAR(255) NOT NULL,
    operating_hours TEXT NOT NULL,
    notification_preferences JSONB, -- {expiryAlerts, surplusReminders, priceUpdates}
    donation_preferences TEXT[],
    category_priority TEXT[],
    barcode_prefix VARCHAR(50),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

-- Users indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_household_id ON users(household_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Inventory indexes
CREATE INDEX IF NOT EXISTS idx_inventory_user_id ON inventory_items(user_id);
CREATE INDEX IF NOT EXISTS idx_inventory_expiry_date ON inventory_items(expiry_date);
CREATE INDEX IF NOT EXISTS idx_inventory_category ON inventory_items(category);
CREATE INDEX IF NOT EXISTS idx_inventory_food_item_id ON inventory_items(food_item_id);

-- Consumption logs indexes
CREATE INDEX IF NOT EXISTS idx_logs_user_id ON consumption_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_logs_consumed_at ON consumption_logs(consumed_at);
CREATE INDEX IF NOT EXISTS idx_logs_inventory_item_id ON consumption_logs(inventory_item_id);

-- Shopping list indexes
CREATE INDEX IF NOT EXISTS idx_shopping_user_id ON shopping_list_items(user_id);
CREATE INDEX IF NOT EXISTS idx_shopping_household_id ON shopping_list_items(household_id);
CREATE INDEX IF NOT EXISTS idx_shopping_purchased ON shopping_list_items(purchased);

-- Meal plans indexes
CREATE INDEX IF NOT EXISTS idx_meal_plans_user_date ON meal_plans(user_id, date);
CREATE INDEX IF NOT EXISTS idx_meal_plans_household_id ON meal_plans(household_id);

-- Badges indexes
CREATE INDEX IF NOT EXISTS idx_badges_user_id ON badges(user_id);
CREATE INDEX IF NOT EXISTS idx_badges_badge_id ON badges(badge_id);

-- Nutrition data indexes
CREATE INDEX IF NOT EXISTS idx_nutrition_user_date ON nutrition_data(user_id, date);

-- Community indexes
CREATE INDEX IF NOT EXISTS idx_surplus_posts_user_id ON community_surplus_posts(user_id);
CREATE INDEX IF NOT EXISTS idx_surplus_posts_status ON community_surplus_posts(status);
CREATE INDEX IF NOT EXISTS idx_surplus_requests_post_id ON surplus_requests(post_id);
CREATE INDEX IF NOT EXISTS idx_leftover_items_user_id ON leftover_items(user_id);
CREATE INDEX IF NOT EXISTS idx_leftover_items_status ON leftover_items(status);
CREATE INDEX IF NOT EXISTS idx_kitchen_events_date ON community_kitchen_events(date);
CREATE INDEX IF NOT EXISTS idx_community_notifications_user_id ON community_notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_community_notifications_read ON community_notifications(read);

-- Restaurant indexes
CREATE INDEX IF NOT EXISTS idx_restaurant_inventory_user_id ON restaurant_inventory_items(user_id);
CREATE INDEX IF NOT EXISTS idx_restaurant_inventory_status ON restaurant_inventory_items(status);
CREATE INDEX IF NOT EXISTS idx_restaurant_surplus_user_id ON restaurant_surplus_items(user_id);
CREATE INDEX IF NOT EXISTS idx_restaurant_donations_user_id ON restaurant_donation_logs(user_id);

-- NGO indexes
CREATE INDEX IF NOT EXISTS idx_ngo_offers_ngo_user_id ON ngo_donation_offers(ngo_user_id);
CREATE INDEX IF NOT EXISTS idx_ngo_offers_status ON ngo_donation_offers(status);
CREATE INDEX IF NOT EXISTS idx_ngo_pickups_offer_id ON ngo_pickup_schedules(offer_id);
CREATE INDEX IF NOT EXISTS idx_ngo_history_ngo_user_id ON ngo_donation_history(ngo_user_id);
CREATE INDEX IF NOT EXISTS idx_ngo_notifications_ngo_user_id ON ngo_notifications(ngo_user_id);
CREATE INDEX IF NOT EXISTS idx_ngo_notifications_read ON ngo_notifications(read);

-- Shop indexes
CREATE INDEX IF NOT EXISTS idx_shop_inventory_user_id ON shop_inventory_items(user_id);
CREATE INDEX IF NOT EXISTS idx_shop_inventory_barcode ON shop_inventory_items(barcode);
CREATE INDEX IF NOT EXISTS idx_shop_inventory_expiry_date ON shop_inventory_items(expiry_date);
CREATE INDEX IF NOT EXISTS idx_shop_surplus_user_id ON shop_surplus_items(user_id);

-- ============================================================================
-- TRIGGERS FOR AUTO-UPDATING updated_at
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers to all tables with updated_at column
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_food_items_updated_at BEFORE UPDATE ON food_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_inventory_items_updated_at BEFORE UPDATE ON inventory_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_consumption_logs_updated_at BEFORE UPDATE ON consumption_logs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_shopping_list_items_updated_at BEFORE UPDATE ON shopping_list_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_meal_plans_updated_at BEFORE UPDATE ON meal_plans
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_resources_updated_at BEFORE UPDATE ON resources
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_uploads_updated_at BEFORE UPDATE ON uploads
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_xp_updated_at BEFORE UPDATE ON user_xp
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_family_preferences_updated_at BEFORE UPDATE ON family_preferences
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_nutrition_data_updated_at BEFORE UPDATE ON nutrition_data
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_price_comparisons_updated_at BEFORE UPDATE ON price_comparisons
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_community_surplus_posts_updated_at BEFORE UPDATE ON community_surplus_posts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_leftover_items_updated_at BEFORE UPDATE ON leftover_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_community_kitchen_events_updated_at BEFORE UPDATE ON community_kitchen_events
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_community_leaderboard_updated_at BEFORE UPDATE ON community_leaderboard
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_community_impact_updated_at BEFORE UPDATE ON community_impact
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_community_profiles_updated_at BEFORE UPDATE ON community_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_restaurant_inventory_items_updated_at BEFORE UPDATE ON restaurant_inventory_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_restaurant_menu_items_updated_at BEFORE UPDATE ON restaurant_menu_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_restaurant_surplus_items_updated_at BEFORE UPDATE ON restaurant_surplus_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_restaurant_impact_metrics_updated_at BEFORE UPDATE ON restaurant_impact_metrics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_restaurant_preferences_updated_at BEFORE UPDATE ON restaurant_preferences
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ngo_capacity_settings_updated_at BEFORE UPDATE ON ngo_capacity_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ngo_donation_offers_updated_at BEFORE UPDATE ON ngo_donation_offers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ngo_pickup_schedules_updated_at BEFORE UPDATE ON ngo_pickup_schedules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ngo_partner_profiles_updated_at BEFORE UPDATE ON ngo_partner_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ngo_feedback_entries_updated_at BEFORE UPDATE ON ngo_feedback_entries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_shop_inventory_items_updated_at BEFORE UPDATE ON shop_inventory_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_shop_surplus_items_updated_at BEFORE UPDATE ON shop_surplus_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_shop_analytics_records_updated_at BEFORE UPDATE ON shop_analytics_records
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_shop_profiles_updated_at BEFORE UPDATE ON shop_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
