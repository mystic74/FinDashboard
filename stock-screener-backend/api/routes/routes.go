package routes

import (
	"log"
	"net/http"
	"stock-screener/api/handlers"
	"stock-screener/config"
	"stock-screener/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the Gin router
func SetupRouter(cfg *config.Config) *gin.Engine {
	if cfg.GinMode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize services
	cacheCleanup := cfg.CacheTTL * 2
	if cacheCleanup < time.Minute {
		cacheCleanup = 10 * time.Minute
	}
	cacheService := services.NewCacheService(cfg.CacheTTL, cacheCleanup)

	demoMode := cfg.DemoMode

	var screenerEngine *services.ScreenerEngine
	if demoMode {
		log.Println("Running in DEMO MODE with mock data")
		screenerEngine = services.NewScreenerEngineWithDemo(cacheService)
	} else {
		log.Println("Running in LIVE MODE with Yahoo Finance API")
		yahooService := services.NewYahooFinanceService(cacheService)
		screenerEngine = services.NewScreenerEngine(yahooService, cacheService)
	}

	// Yahoo service needed for stock handler (even in demo mode, just won't work)
	yahooService := services.NewYahooFinanceService(cacheService)

	// Initialize handlers
	profileHandler := handlers.NewProfileHandler()
	screenerHandler := handlers.NewScreenerHandler(screenerEngine, profileHandler)
	stockHandler := handlers.NewStockHandler(yahooService)
	filterHandler := handlers.NewFilterHandler()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"demoMode":  demoMode,
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Screener routes
		screeners := v1.Group("/screeners")
		{
			screeners.GET("", screenerHandler.GetAllScreeners)
			screeners.GET("/summary", screenerHandler.GetScreenersSummary)
			screeners.POST("/custom", screenerHandler.RunCustomScreener)
			screeners.GET("/:name", screenerHandler.RunScreener)
		}

		// Stock routes
		stocks := v1.Group("/stocks")
		{
			stocks.GET("", stockHandler.GetMultipleStocks)
			stocks.GET("/:symbol", stockHandler.GetStock)
			stocks.GET("/:symbol/fundamentals", stockHandler.GetStockFundamentals)
			stocks.GET("/:symbol/history", stockHandler.GetHistoricalPrices)
			stocks.GET("/:symbol/quote", stockHandler.GetStockQuote)
		}

		// Filter routes
		filters := v1.Group("/filters")
		{
			filters.GET("", filterHandler.GetAllFilters)
			filters.GET("/categories", filterHandler.GetFilterCategories)
			filters.GET("/marketcap-ranges", filterHandler.GetMarketCapRanges)
			filters.GET("/:category", filterHandler.GetFiltersByCategory)
		}

		// Sector routes
		v1.GET("/sectors", screenerHandler.GetSectorPerformance)
		v1.GET("/sectors/list", filterHandler.GetSectors)

		// Market Profile routes
		profiles := v1.Group("/profiles")
		{
			profiles.GET("", profileHandler.GetAllProfiles)
			profiles.POST("/reset", profileHandler.ResetAllProfiles)
			profiles.GET("/:country", profileHandler.GetProfile)
			profiles.PUT("/:country", profileHandler.UpdateProfile)
			profiles.POST("/:country/reset", profileHandler.ResetProfile)
		}
	}

	return router
}
