package config

import (
	"os"
	"time"
)

// Config holds the application configuration
type Config struct {
	ServerPort       string
	CacheTTL         time.Duration
	RateLimitPerMin  int
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
	return &Config{
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		CacheTTL:         5 * time.Minute,
		RateLimitPerMin:  60,
		YahooFinanceURL:  "https://query1.finance.yahoo.com",
		FMPBaseURL:       "https://financialmodelingprep.com/api/v3",
		FMPAPIKey:        getEnv("FMP_API_KEY", ""),
		AlphaVantageURL:  "https://www.alphavantage.co/query",
		AlphaVantageKey:  getEnv("ALPHA_VANTAGE_KEY", ""),
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		RequestTimeout:   30 * time.Second,
		MaxConcurrent:    10,
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
