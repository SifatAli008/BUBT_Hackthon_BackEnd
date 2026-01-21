package routes

import (
	"foodlink_backend/config"
	_ "foodlink_backend/docs" // Import docs for Swagger
	"foodlink_backend/features/auth"
	"foodlink_backend/features/badges"
	"foodlink_backend/features/community/kitchen_events"
	"foodlink_backend/features/community/leaderboard"
	"foodlink_backend/features/community/leftovers"
	"foodlink_backend/features/community/profiles"
	"foodlink_backend/features/community/surplus"
	"foodlink_backend/features/consumption"
	"foodlink_backend/features/food_items"
	"foodlink_backend/features/inventory"
	"foodlink_backend/features/meal_plans"
	ngo_capacity "foodlink_backend/features/ngo/capacity"
	ngo_feedback "foodlink_backend/features/ngo/feedback"
	ngo_history "foodlink_backend/features/ngo/history"
	ngo_offers "foodlink_backend/features/ngo/offers"
	ngo_partners "foodlink_backend/features/ngo/partners"
	ngo_pickups "foodlink_backend/features/ngo/pickups"
	"foodlink_backend/features/nutrition"
	"foodlink_backend/features/preferences"
	"foodlink_backend/features/price_comparisons"
	"foodlink_backend/features/shopping_list"
	restaurant_donations "foodlink_backend/features/restaurant/donations"
	restaurant_inventory "foodlink_backend/features/restaurant/inventory"
	restaurant_menu "foodlink_backend/features/restaurant/menu"
	restaurant_preferences "foodlink_backend/features/restaurant/preferences"
	restaurant_staff "foodlink_backend/features/restaurant/staff"
	restaurant_surplus "foodlink_backend/features/restaurant/surplus"
	"foodlink_backend/features/xp"
	"foodlink_backend/handlers"
	"foodlink_backend/middleware"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func mountWithOptionalSlash(mux *http.ServeMux, prefix string, handler http.Handler) {
	strip := http.StripPrefix(prefix, handler)

	// Match subtree requests: /prefix/...
	mux.Handle(prefix+"/", strip)

	// Match exact prefix: /prefix
	mux.Handle(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make StripPrefix produce "/" for subrouters that expect root.
		r.URL.Path = prefix + "/"
		strip.ServeHTTP(w, r)
	}))
}

func SetupRoutes(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint (no middleware needed)
	mux.HandleFunc("/health", handlers.HealthCheck)

	// API routes
	mux.HandleFunc("/api/v1/", handlers.APIV1)

	// Authentication routes
	authService := auth.NewService(cfg)
	authHandler := auth.NewHandler(authService)
	authRoutes := auth.SetupRoutes(authService, authHandler)
	mountWithOptionalSlash(mux, "/api/v1/auth", authRoutes)

	// Food Items routes (public, but admin-only for create/update/delete)
	foodItemsService := food_items.NewService()
	foodItemsHandler := food_items.NewHandler(foodItemsService)
	foodItemsRoutes := food_items.SetupRoutes(foodItemsHandler)
	mountWithOptionalSlash(mux, "/api/v1/food-items", foodItemsRoutes)

	// Inventory routes (protected)
	inventoryService := inventory.NewService()
	inventoryHandler := inventory.NewHandler(inventoryService)
	inventoryRoutes := inventory.SetupRoutes(inventoryService, inventoryHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/inventory", inventoryRoutes)

	// Shopping List routes (protected)
	shoppingListService := shopping_list.NewService()
	shoppingListHandler := shopping_list.NewHandler(shoppingListService)
	shoppingListRoutes := shopping_list.SetupRoutes(shoppingListService, shoppingListHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/shopping-list", shoppingListRoutes)

	// Consumption routes (protected)
	consumptionService := consumption.NewService()
	consumptionHandler := consumption.NewHandler(consumptionService)
	consumptionRoutes := consumption.SetupRoutes(consumptionService, consumptionHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/consumption", consumptionRoutes)

	// Meal Plans routes (protected)
	mealPlansService := meal_plans.NewService()
	mealPlansHandler := meal_plans.NewHandler(mealPlansService)
	mealPlansRoutes := meal_plans.SetupRoutes(mealPlansService, mealPlansHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/meal-plans", mealPlansRoutes)

	// Preferences routes (protected)
	preferencesService := preferences.NewService()
	preferencesHandler := preferences.NewHandler(preferencesService)
	preferencesRoutes := preferences.SetupRoutes(preferencesService, preferencesHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/preferences", preferencesRoutes)

	// Nutrition routes (protected)
	nutritionService := nutrition.NewService()
	nutritionHandler := nutrition.NewHandler(nutritionService)
	nutritionRoutes := nutrition.SetupRoutes(nutritionService, nutritionHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/nutrition", nutritionRoutes)

	// Price Comparisons routes (public)
	priceComparisonsService := price_comparisons.NewService()
	priceComparisonsHandler := price_comparisons.NewHandler(priceComparisonsService)
	priceComparisonsRoutes := price_comparisons.SetupRoutes(priceComparisonsHandler)
	mountWithOptionalSlash(mux, "/api/v1/price-comparisons", priceComparisonsRoutes)

	// Badges routes (protected)
	badgesService := badges.NewService()
	badgesHandler := badges.NewHandler(badgesService)
	badgesRoutes := badges.SetupRoutes(badgesService, badgesHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/badges", badgesRoutes)

	// XP routes (protected)
	xpService := xp.NewService()
	xpHandler := xp.NewHandler(xpService)
	xpRoutes := xp.SetupRoutes(xpService, xpHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/xp", xpRoutes)

	// Community Surplus routes (protected)
	surplusService := surplus.NewService()
	surplusHandler := surplus.NewHandler(surplusService)
	surplusRoutes := surplus.SetupRoutes(surplusService, surplusHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/community/surplus", surplusRoutes)

	// Community Leftovers routes (protected)
	leftoversService := leftovers.NewService()
	leftoversHandler := leftovers.NewHandler(leftoversService)
	leftoversRoutes := leftovers.SetupRoutes(leftoversService, leftoversHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/community/leftovers", leftoversRoutes)

	// Community Kitchen Events routes (protected)
	kitchenEventsService := kitchen_events.NewService()
	kitchenEventsHandler := kitchen_events.NewHandler(kitchenEventsService)
	kitchenEventsRoutes := kitchen_events.SetupRoutes(kitchenEventsService, kitchenEventsHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/community/kitchen-events", kitchenEventsRoutes)

	// Community Leaderboard & Impact routes (protected)
	leaderboardService := leaderboard.NewService()
	leaderboardHandler := leaderboard.NewHandler(leaderboardService)
	leaderboardRoutes := leaderboard.SetupRoutes(leaderboardService, leaderboardHandler, auth.AuthMiddleware(authService))
	mux.Handle("/api/v1/community/leaderboard", http.StripPrefix("/api/v1/community", leaderboardRoutes))
	mux.Handle("/api/v1/community/impact/", http.StripPrefix("/api/v1/community", leaderboardRoutes))

	// Community Profiles routes (protected)
	profilesService := profiles.NewService()
	profilesHandler := profiles.NewHandler(profilesService)
	profilesRoutes := profiles.SetupRoutes(profilesService, profilesHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/community/profile", profilesRoutes)

	// Restaurant Inventory routes (protected)
	restaurantInventoryService := restaurant_inventory.NewService()
	restaurantInventoryHandler := restaurant_inventory.NewHandler(restaurantInventoryService)
	restaurantInventoryRoutes := restaurant_inventory.SetupRoutes(restaurantInventoryService, restaurantInventoryHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/restaurant/inventory", restaurantInventoryRoutes)

	// Restaurant Menu routes (protected)
	restaurantMenuService := restaurant_menu.NewService()
	restaurantMenuHandler := restaurant_menu.NewHandler(restaurantMenuService)
	restaurantMenuRoutes := restaurant_menu.SetupRoutes(restaurantMenuService, restaurantMenuHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/restaurant/menu", restaurantMenuRoutes)

	// Restaurant Surplus routes (protected)
	restaurantSurplusService := restaurant_surplus.NewService()
	restaurantSurplusHandler := restaurant_surplus.NewHandler(restaurantSurplusService)
	restaurantSurplusRoutes := restaurant_surplus.SetupRoutes(restaurantSurplusService, restaurantSurplusHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/restaurant/surplus", restaurantSurplusRoutes)

	// Restaurant Donations & Impact routes (protected)
	restaurantDonationsService := restaurant_donations.NewService()
	restaurantDonationsHandler := restaurant_donations.NewHandler(restaurantDonationsService)
	restaurantDonationsRoutes := restaurant_donations.SetupRoutes(restaurantDonationsService, restaurantDonationsHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/restaurant/donations", restaurantDonationsRoutes)
	mux.Handle("/api/v1/restaurant/impact", http.StripPrefix("/api/v1/restaurant", restaurantDonationsRoutes))

	// Restaurant Staff Management routes (protected)
	restaurantStaffService := restaurant_staff.NewService()
	restaurantStaffHandler := restaurant_staff.NewHandler(restaurantStaffService)
	restaurantStaffRoutes := restaurant_staff.SetupRoutes(restaurantStaffService, restaurantStaffHandler, auth.AuthMiddleware(authService))
	mux.Handle("/api/v1/restaurant/tasks/", http.StripPrefix("/api/v1/restaurant", restaurantStaffRoutes))
	mux.Handle("/api/v1/restaurant/shifts", http.StripPrefix("/api/v1/restaurant", restaurantStaffRoutes))

	// Restaurant Preferences routes (protected)
	restaurantPreferencesService := restaurant_preferences.NewService()
	restaurantPreferencesHandler := restaurant_preferences.NewHandler(restaurantPreferencesService)
	restaurantPreferencesRoutes := restaurant_preferences.SetupRoutes(restaurantPreferencesService, restaurantPreferencesHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/restaurant/preferences", restaurantPreferencesRoutes)

	// NGO Capacity Settings routes (protected)
	ngoCapacityService := ngo_capacity.NewService()
	ngoCapacityHandler := ngo_capacity.NewHandler(ngoCapacityService)
	ngoCapacityRoutes := ngo_capacity.SetupRoutes(ngoCapacityService, ngoCapacityHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/ngo/capacity", ngoCapacityRoutes)

	// NGO Donation Offers routes (protected)
	ngoOffersService := ngo_offers.NewService()
	ngoOffersHandler := ngo_offers.NewHandler(ngoOffersService)
	ngoOffersRoutes := ngo_offers.SetupRoutes(ngoOffersService, ngoOffersHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/ngo/offers", ngoOffersRoutes)

	// NGO Pickup Schedules routes (protected)
	ngoPickupsService := ngo_pickups.NewService()
	ngoPickupsHandler := ngo_pickups.NewHandler(ngoPickupsService)
	ngoPickupsRoutes := ngo_pickups.SetupRoutes(ngoPickupsService, ngoPickupsHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/ngo/pickups", ngoPickupsRoutes)

	// NGO Donation History routes (protected)
	ngoHistoryService := ngo_history.NewService()
	ngoHistoryHandler := ngo_history.NewHandler(ngoHistoryService)
	ngoHistoryRoutes := ngo_history.SetupRoutes(ngoHistoryService, ngoHistoryHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/ngo/history", ngoHistoryRoutes)

	// NGO Partner Management routes (protected)
	ngoPartnersService := ngo_partners.NewService()
	ngoPartnersHandler := ngo_partners.NewHandler(ngoPartnersService)
	ngoPartnersRoutes := ngo_partners.SetupRoutes(ngoPartnersService, ngoPartnersHandler, auth.AuthMiddleware(authService))
	mountWithOptionalSlash(mux, "/api/v1/ngo/partners", ngoPartnersRoutes)

	// NGO Feedback & Impact routes (protected)
	ngoFeedbackService := ngo_feedback.NewService()
	ngoFeedbackHandler := ngo_feedback.NewHandler(ngoFeedbackService)
	ngoFeedbackRoutes := ngo_feedback.SetupRoutes(ngoFeedbackService, ngoFeedbackHandler, auth.AuthMiddleware(authService))
	mux.Handle("/api/v1/ngo/feedback", http.StripPrefix("/api/v1/ngo", ngoFeedbackRoutes))
	mux.Handle("/api/v1/ngo/stories", http.StripPrefix("/api/v1/ngo", ngoFeedbackRoutes))

	// Swagger documentation with CORS support
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // Use relative path
	)

	// Wrap Swagger handler with explicit CORS support
	corsSwaggerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers explicitly for Swagger
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "*")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		swaggerHandler.ServeHTTP(w, r)
	})

	mux.Handle("/swagger/", corsSwaggerHandler)

	// Redirect /swagger to /swagger/index.html
	mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	// Apply middleware chain
	handler := middleware.Chain(
		middleware.RecoverPanic,
		middleware.RequestID,
		middleware.Logging,
		middleware.CORS,
		middleware.ErrorHandler,
	)(mux)

	return handler
}
