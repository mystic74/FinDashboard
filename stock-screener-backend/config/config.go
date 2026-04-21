package config

import (
	"os"
	"strings"
	"time"
)

// Config holds the application configuration
type Config struct {
	ServerPort       string
	CacheTTL         time.Duration
	GinMode          string
	RateLimitPerMin  int
	DemoMode         bool
	YahooFinanceURL  string
	FMPBaseURL       string
	FMPAPIKey        string
	AlphaVantageURL  string
	AlphaVantageKey  string
	AllowedOrigins   []string
	RequestTimeout   time.Duration
	MaxConcurrent    int
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	// Parse CORS origins from env (comma-separated)
	corsOrigins := getEnv("CORS_ORIGIN", "http://localhost:3000,http://localhost:5173")
	origins := strings.Split(corsOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	// Demo mode: default true unless explicitly set to "false"
	demoMode := getEnv("DEMO_MODE", "true") != "false"

	cacheTTL := parseDurationEnv("CACHE_TTL", 5*time.Minute)
	ginMode := strings.ToLower(strings.TrimSpace(getEnv("GIN_MODE", "release")))
	if ginMode == "" {
		ginMode = "release"
	}

	return &Config{
		ServerPort:       getEnv("PORT", "8080"),
		CacheTTL:         cacheTTL,
		GinMode:          ginMode,
		RateLimitPerMin:  60,
		DemoMode:         demoMode,
		YahooFinanceURL:  "https://query1.finance.yahoo.com",
		FMPBaseURL:       "https://financialmodelingprep.com/api/v3",
		FMPAPIKey:        getEnv("FMP_API_KEY", ""),
		AlphaVantageURL:  "https://www.alphavantage.co/query",
		AlphaVantageKey:  getEnv("ALPHA_VANTAGE_KEY", ""),
		AllowedOrigins:   origins,
		RequestTimeout:   30 * time.Second,
		MaxConcurrent:    10,
	}
}

// HasAPIKeys returns true if any API keys are configured
func (c *Config) HasAPIKeys() bool {
	return c.FMPAPIKey != "" || c.AlphaVantageKey != ""
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseDurationEnv(key string, defaultVal time.Duration) time.Duration {
	s := strings.TrimSpace(os.Getenv(key))
	if s == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return defaultVal
	}
	return d
}
