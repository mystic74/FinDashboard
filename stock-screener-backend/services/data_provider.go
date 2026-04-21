package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"stock-screener/models"
	"sync"
	"time"
)

// DataProvider defines the interface for stock data providers
type DataProvider interface {
	GetQuotes(symbols []string) ([]models.Stock, error)
	IsAvailable() bool
	GetStatus() map[string]interface{}
}

// DataProviderManager manages multiple data providers with fallback
type DataProviderManager struct {
	providers     map[string]DataProvider
	providerOrder []string
	cache         *CacheService
	mu            sync.RWMutex

	// Yahoo as primary (unofficial, might be blocked)
	yahooService *YahooFinanceService

	// Official APIs as fallback
	alphaVantage *AlphaVantageService
	fmpService   *FMPService

	// Mock data for demo mode
	mockService *MockDataService
	demoMode    bool
}

// DataProviderConfig holds API keys and configuration
type DataProviderConfig struct {
	AlphaVantageKey  string
	FMPKey           string
	DemoMode         bool
	YahooQuoteDriver string
}

// NewDataProviderManager creates a new data provider manager
func NewDataProviderManager(config DataProviderConfig, cache *CacheService) *DataProviderManager {
	m := &DataProviderManager{
		providers:     make(map[string]DataProvider),
		providerOrder: []string{},
		cache:         cache,
		demoMode:      config.DemoMode,
	}

	// Initialize providers based on available API keys
	if config.DemoMode {
		log.Println("[DataProvider] Running in DEMO MODE with mock data")
		m.mockService = NewMockDataService()
		return m
	}

	// Initialize Alpha Vantage if API key provided
	if config.AlphaVantageKey != "" {
		m.alphaVantage = NewAlphaVantageService(config.AlphaVantageKey, cache)
		m.providerOrder = append(m.providerOrder, "alpha_vantage")
		log.Println("[DataProvider] Alpha Vantage initialized (free tier: 25/day, 5/min)")
	}

	// Initialize FMP if API key provided
	if config.FMPKey != "" {
		m.fmpService = NewFMPService(config.FMPKey, cache)
		m.providerOrder = append(m.providerOrder, "fmp")
		log.Println("[DataProvider] Financial Modeling Prep initialized (free tier: 250/day)")
	}

	// Initialize Yahoo Finance as fallback (unofficial, may be blocked).
	m.yahooService = NewYahooFinanceServiceWithDriver(cache, config.YahooQuoteDriver)
	m.providerOrder = append(m.providerOrder, "yahoo")
	log.Printf("[DataProvider] Yahoo Finance fallback initialized (driver=%s)", m.yahooService.QuoteDriver())

	if len(m.providerOrder) == 0 {
		log.Println("[DataProvider] WARNING: No data providers configured, falling back to demo mode")
		m.demoMode = true
		m.mockService = NewMockDataService()
	}

	return m
}

// GetQuotes fetches quotes using available providers with fallback
func (m *DataProviderManager) GetQuotes(ctx context.Context, symbols []string) ([]models.Stock, error) {
	if m.demoMode {
		return m.mockService.GetQuotes(symbols)
	}

	var lastErr error

	// Try each provider in order
	for _, providerName := range m.providerOrder {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		stocks, err := m.tryProvider(ctx, providerName, symbols)
		if err == nil {
			return stocks, nil
		}

		lastErr = err
		log.Printf("[DataProvider] %s failed: %v, trying next provider...", providerName, err)

		// If rate limited, don't try other providers yet - they might have same limits
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			if errors.Is(apiErr.Underlying, ErrQuotaExceeded) {
				// Quota exceeded - this provider is done for the day
				log.Printf("[DataProvider] %s quota exceeded, skipping for remainder of session", providerName)
				continue
			}
			if errors.Is(apiErr.Underlying, ErrInvalidAPIKey) {
				// Invalid API key - don't retry this provider
				log.Printf("[DataProvider] %s has invalid API key, disabling", providerName)
				continue
			}
		}
	}

	// All providers failed - fall back to demo mode if available
	if m.mockService == nil {
		m.mockService = NewMockDataService()
	}

	log.Printf("[DataProvider] All providers failed, falling back to mock data. Last error: %v", lastErr)
	return m.mockService.GetQuotes(symbols)
}

// tryProvider attempts to get quotes from a specific provider
func (m *DataProviderManager) tryProvider(ctx context.Context, name string, symbols []string) ([]models.Stock, error) {
	switch name {
	case "yahoo":
		if m.yahooService != nil {
			return m.yahooService.GetQuotes(ctx, symbols)
		}
	case "alpha_vantage":
		if m.alphaVantage != nil && m.alphaVantage.IsAvailable() {
			// Alpha Vantage doesn't have batch quotes on free tier
			// Need to fetch one by one (expensive on rate limits)
			return m.fetchAlphaVantageQuotes(symbols)
		}
	case "fmp":
		if m.fmpService != nil && m.fmpService.IsAvailable() {
			return m.fmpService.GetQuotes(symbols)
		}
	}
	return nil, fmt.Errorf("provider %s not available", name)
}

// fetchAlphaVantageQuotes fetches quotes one by one (AV limitation)
func (m *DataProviderManager) fetchAlphaVantageQuotes(symbols []string) ([]models.Stock, error) {
	stocks := make([]models.Stock, 0, len(symbols))

	for _, symbol := range symbols {
		if !m.alphaVantage.IsAvailable() {
			// Rate limit hit during batch
			if len(stocks) > 0 {
				log.Printf("[AlphaVantage] Rate limited after %d symbols, returning partial results", len(stocks))
				return stocks, nil
			}
			return nil, NewRateLimitError("alpha_vantage", m.alphaVantage.rateLimiter.GetWaitTime())
		}

		stock, err := m.alphaVantage.GetQuote(symbol)
		if err != nil {
			var apiErr *APIError
			if errors.As(err, &apiErr) && errors.Is(apiErr.Underlying, ErrRateLimited) {
				// Return what we have so far
				if len(stocks) > 0 {
					return stocks, nil
				}
				return nil, err
			}
			// Skip this symbol but continue
			log.Printf("[AlphaVantage] Failed to get quote for %s: %v", symbol, err)
			continue
		}
		stocks = append(stocks, *stock)
	}

	return stocks, nil
}

// GetStockWithFundamentals fetches detailed stock data with fundamentals
func (m *DataProviderManager) GetStockWithFundamentals(ctx context.Context, symbol string) (*models.Stock, error) {
	if m.demoMode {
		stocks := m.mockService.GetAllStocks()
		for i := range stocks {
			if stocks[i].Symbol == symbol {
				return &stocks[i], nil
			}
		}
		return nil, ErrSymbolNotFound
	}

	// Try FMP first (has good fundamental data)
	if m.fmpService != nil && m.fmpService.IsAvailable() {
		stock, err := m.fmpService.GetStockWithFundamentals(symbol)
		if err == nil {
			return stock, nil
		}
		log.Printf("[DataProvider] FMP fundamentals failed for %s: %v", symbol, err)
	}

	// Try Alpha Vantage
	if m.alphaVantage != nil && m.alphaVantage.IsAvailable() {
		stock, err := m.alphaVantage.GetFundamentals(symbol)
		if err == nil {
			return stock, nil
		}
		log.Printf("[DataProvider] AlphaVantage fundamentals failed for %s: %v", symbol, err)
	}

	// Fall back to Yahoo
	if m.yahooService != nil {
		return m.yahooService.GetStockFundamentals(symbol)
	}

	return nil, fmt.Errorf("no providers available for fundamentals")
}

// GetProviderStatus returns status of all providers
func (m *DataProviderManager) GetProviderStatus() map[string]interface{} {
	status := map[string]interface{}{
		"demoMode": m.demoMode,
	}

	providers := []map[string]interface{}{}

	if m.yahooService != nil {
		providers = append(providers, map[string]interface{}{
			"name":      "yahoo",
			"type":      "unofficial",
			"available": true, // Always "available" but might fail
		})
	}

	if m.alphaVantage != nil {
		providers = append(providers, m.alphaVantage.GetStatus())
	}

	if m.fmpService != nil {
		providers = append(providers, m.fmpService.GetStatus())
	}

	status["providers"] = providers
	status["providerOrder"] = m.providerOrder

	return status
}

// SetDemoMode switches to/from demo mode
func (m *DataProviderManager) SetDemoMode(enabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.demoMode = enabled
	if enabled && m.mockService == nil {
		m.mockService = NewMockDataService()
	}
}

// IsDemoMode returns whether demo mode is active
func (m *DataProviderManager) IsDemoMode() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.demoMode
}

// HealthCheck performs a health check on all providers
func (m *DataProviderManager) HealthCheck(ctx context.Context) map[string]interface{} {
	results := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"demoMode":  m.demoMode,
	}

	if m.demoMode {
		results["status"] = "healthy"
		results["message"] = "Running in demo mode with mock data"
		return results
	}

	healthyProviders := 0
	providerResults := map[string]interface{}{}

	// Check Yahoo
	if m.yahooService != nil {
		// Quick test with a known symbol
		_, err := m.yahooService.GetQuotes(context.Background(), []string{"AAPL"})
		if err == nil {
			providerResults["yahoo"] = map[string]interface{}{
				"healthy": true,
			}
			healthyProviders++
		} else {
			providerResults["yahoo"] = map[string]interface{}{
				"healthy": false,
				"error":   err.Error(),
			}
		}
	}

	// Check Alpha Vantage
	if m.alphaVantage != nil {
		providerResults["alpha_vantage"] = map[string]interface{}{
			"healthy":   m.alphaVantage.IsAvailable(),
			"rateLimit": m.alphaVantage.rateLimiter.GetStatus(),
		}
		if m.alphaVantage.IsAvailable() {
			healthyProviders++
		}
	}

	// Check FMP
	if m.fmpService != nil {
		providerResults["fmp"] = map[string]interface{}{
			"healthy":   m.fmpService.IsAvailable(),
			"rateLimit": m.fmpService.rateLimiter.GetStatus(),
		}
		if m.fmpService.IsAvailable() {
			healthyProviders++
		}
	}

	results["providers"] = providerResults

	if healthyProviders > 0 {
		results["status"] = "healthy"
	} else {
		results["status"] = "degraded"
		results["message"] = "No healthy providers, will use mock data"
	}

	return results
}
