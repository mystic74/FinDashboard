package services

import (
	"encoding/json"
	"fmt"
	"log"
	"stock-screener/models"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	alphaVantageBaseURL = "https://www.alphavantage.co/query"
	alphaVantageProvider = "alpha_vantage"
)

// AlphaVantageService handles Alpha Vantage API interactions
type AlphaVantageService struct {
	client      *resty.Client
	apiKey      string
	rateLimiter *RateLimiter
	status      *ProviderStatus
	cache       *CacheService
}

// AlphaVantage API response types
type AVGlobalQuote struct {
	Symbol           string `json:"01. symbol"`
	Open             string `json:"02. open"`
	High             string `json:"03. high"`
	Low              string `json:"04. low"`
	Price            string `json:"05. price"`
	Volume           string `json:"06. volume"`
	LatestTradingDay string `json:"07. latest trading day"`
	PreviousClose    string `json:"08. previous close"`
	Change           string `json:"09. change"`
	ChangePercent    string `json:"10. change percent"`
}

type AVQuoteResponse struct {
	GlobalQuote AVGlobalQuote          `json:"Global Quote"`
	Note        string                 `json:"Note,omitempty"`        // Rate limit message
	Information string                 `json:"Information,omitempty"` // API key issues
	ErrorMessage string                `json:"Error Message,omitempty"`
}

type AVOverviewResponse struct {
	Symbol                     string `json:"Symbol"`
	Name                       string `json:"Name"`
	Exchange                   string `json:"Exchange"`
	Currency                   string `json:"Currency"`
	Country                    string `json:"Country"`
	Sector                     string `json:"Sector"`
	Industry                   string `json:"Industry"`
	MarketCapitalization       string `json:"MarketCapitalization"`
	PERatio                    string `json:"PERatio"`
	PEGRatio                   string `json:"PEGRatio"`
	BookValue                  string `json:"BookValue"`
	DividendPerShare           string `json:"DividendPerShare"`
	DividendYield              string `json:"DividendYield"`
	EPS                        string `json:"EPS"`
	RevenuePerShareTTM         string `json:"RevenuePerShareTTM"`
	ProfitMargin               string `json:"ProfitMargin"`
	OperatingMarginTTM         string `json:"OperatingMarginTTM"`
	ReturnOnAssetsTTM          string `json:"ReturnOnAssetsTTM"`
	ReturnOnEquityTTM          string `json:"ReturnOnEquityTTM"`
	RevenueTTM                 string `json:"RevenueTTM"`
	GrossProfitTTM             string `json:"GrossProfitTTM"`
	QuarterlyEarningsGrowthYOY string `json:"QuarterlyEarningsGrowthYOY"`
	QuarterlyRevenueGrowthYOY  string `json:"QuarterlyRevenueGrowthYOY"`
	Beta                       string `json:"Beta"`
	Week52High                 string `json:"52WeekHigh"`
	Week52Low                  string `json:"52WeekLow"`
	MA50                       string `json:"50DayMovingAverage"`
	MA200                      string `json:"200DayMovingAverage"`
	SharesOutstanding          string `json:"SharesOutstanding"`
	// Error fields
	Note        string `json:"Note,omitempty"`
	Information string `json:"Information,omitempty"`
}

// NewAlphaVantageService creates a new Alpha Vantage service
// Free tier: 25 calls/day, 5 calls/minute
func NewAlphaVantageService(apiKey string, cache *CacheService) *AlphaVantageService {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(2)
	client.SetRetryWaitTime(2 * time.Second)

	return &AlphaVantageService{
		client:      client,
		apiKey:      apiKey,
		rateLimiter: NewRateLimiter(alphaVantageProvider, 5, 25), // Free tier limits
		status: &ProviderStatus{
			Name:      alphaVantageProvider,
			Available: apiKey != "",
		},
		cache: cache,
	}
}

// IsAvailable returns whether the service is available
func (a *AlphaVantageService) IsAvailable() bool {
	return a.apiKey != "" && a.status.IsHealthy() && a.rateLimiter.CanMakeCall()
}

// GetStatus returns the current provider status
func (a *AlphaVantageService) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"provider":    alphaVantageProvider,
		"available":   a.IsAvailable(),
		"hasAPIKey":   a.apiKey != "",
		"healthy":     a.status.IsHealthy(),
		"rateLimit":   a.rateLimiter.GetStatus(),
	}
}

// GetQuote fetches a quote for a single symbol
func (a *AlphaVantageService) GetQuote(symbol string) (*models.Stock, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("av:quote:%s", symbol)
	if cached, found := a.cache.Get(cacheKey); found {
		if stock, ok := cached.(*models.Stock); ok {
			return stock, nil
		}
	}

	// Check rate limit
	if !a.rateLimiter.CanMakeCall() {
		waitTime := a.rateLimiter.GetWaitTime()
		return nil, NewRateLimitError(alphaVantageProvider, waitTime)
	}

	// Make API call
	resp, err := a.client.R().
		SetQueryParams(map[string]string{
			"function": "GLOBAL_QUOTE",
			"symbol":   symbol,
			"apikey":   a.apiKey,
		}).
		Get(alphaVantageBaseURL)

	a.rateLimiter.RecordCall()

	if err != nil {
		a.status.RecordFailure(err)
		return nil, fmt.Errorf("alpha vantage request failed: %w", err)
	}

	// Parse response
	var quoteResp AVQuoteResponse
	if err := json.Unmarshal(resp.Body(), &quoteResp); err != nil {
		a.status.RecordFailure(err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if err := a.checkAPIError(quoteResp.Note, quoteResp.Information, quoteResp.ErrorMessage); err != nil {
		return nil, err
	}

	// Convert to Stock model
	stock := a.convertQuoteToStock(quoteResp.GlobalQuote)
	a.status.RecordSuccess()

	// Cache the result
	a.cache.Set(cacheKey, stock, 5*time.Minute)

	return stock, nil
}

// GetFundamentals fetches fundamental data for a symbol
func (a *AlphaVantageService) GetFundamentals(symbol string) (*models.Stock, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("av:fundamentals:%s", symbol)
	if cached, found := a.cache.Get(cacheKey); found {
		if stock, ok := cached.(*models.Stock); ok {
			return stock, nil
		}
	}

	// Check rate limit
	if !a.rateLimiter.CanMakeCall() {
		waitTime := a.rateLimiter.GetWaitTime()
		return nil, NewRateLimitError(alphaVantageProvider, waitTime)
	}

	// Make API call
	resp, err := a.client.R().
		SetQueryParams(map[string]string{
			"function": "OVERVIEW",
			"symbol":   symbol,
			"apikey":   a.apiKey,
		}).
		Get(alphaVantageBaseURL)

	a.rateLimiter.RecordCall()

	if err != nil {
		a.status.RecordFailure(err)
		return nil, fmt.Errorf("alpha vantage request failed: %w", err)
	}

	// Parse response
	var overview AVOverviewResponse
	if err := json.Unmarshal(resp.Body(), &overview); err != nil {
		a.status.RecordFailure(err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if err := a.checkAPIError(overview.Note, overview.Information, ""); err != nil {
		return nil, err
	}

	// Check if symbol was found
	if overview.Symbol == "" {
		return nil, &APIError{
			Provider:   alphaVantageProvider,
			Code:       "NOT_FOUND",
			Message:    fmt.Sprintf("Symbol '%s' not found", symbol),
			Underlying: ErrSymbolNotFound,
		}
	}

	// Convert to Stock model
	stock := a.convertOverviewToStock(overview)
	a.status.RecordSuccess()

	// Cache the result
	a.cache.Set(cacheKey, stock, 15*time.Minute)

	return stock, nil
}

// checkAPIError checks response for API-level errors
func (a *AlphaVantageService) checkAPIError(note, info, errMsg string) error {
	// Check for rate limit (Note field)
	if note != "" {
		if strings.Contains(note, "API call frequency") || strings.Contains(note, "rate limit") {
			log.Printf("[AlphaVantage] Rate limit hit: %s", note)
			a.rateLimiter.SetRetryAfter(time.Minute)
			return NewRateLimitError(alphaVantageProvider, time.Minute)
		}
		// Daily limit exceeded
		if strings.Contains(note, "daily") || strings.Contains(note, "25 requests") {
			log.Printf("[AlphaVantage] Daily quota exceeded: %s", note)
			return NewQuotaExceededError(alphaVantageProvider, "daily")
		}
	}

	// Check for API key issues (Information field)
	if info != "" {
		if strings.Contains(info, "API key") || strings.Contains(info, "invalid") {
			log.Printf("[AlphaVantage] API key issue: %s", info)
			a.status.Available = false
			return NewInvalidAPIKeyError(alphaVantageProvider)
		}
		// Premium feature
		if strings.Contains(info, "premium") {
			log.Printf("[AlphaVantage] Premium required: %s", info)
			return NewPremiumRequiredError(alphaVantageProvider, "requested endpoint")
		}
	}

	// Check for error message
	if errMsg != "" {
		return &APIError{
			Provider: alphaVantageProvider,
			Code:     "API_ERROR",
			Message:  errMsg,
		}
	}

	return nil
}

// convertQuoteToStock converts Alpha Vantage quote to Stock model
func (a *AlphaVantageService) convertQuoteToStock(quote AVGlobalQuote) *models.Stock {
	return &models.Stock{
		Symbol:        quote.Symbol,
		Price:         parseFloat(quote.Price),
		Open:          parseFloat(quote.Open),
		High:          parseFloat(quote.High),
		Low:           parseFloat(quote.Low),
		PreviousClose: parseFloat(quote.PreviousClose),
		Change:        parseFloat(quote.Change),
		ChangePercent: parseFloat(strings.TrimSuffix(quote.ChangePercent, "%")),
		Volume:        parseInt64(quote.Volume),
		LastUpdated:   time.Now(),
	}
}

// convertOverviewToStock converts Alpha Vantage overview to Stock model
func (a *AlphaVantageService) convertOverviewToStock(ov AVOverviewResponse) *models.Stock {
	return &models.Stock{
		Symbol:            ov.Symbol,
		Name:              ov.Name,
		Exchange:          ov.Exchange,
		Currency:          ov.Currency,
		Country:           ov.Country,
		Sector:            ov.Sector,
		Industry:          ov.Industry,
		MarketCap:         parseInt64(ov.MarketCapitalization),
		PERatio:           parseFloat(ov.PERatio),
		PEGRatio:          parseFloat(ov.PEGRatio),
		BookValuePerShare: parseFloat(ov.BookValue),
		DividendPerShare:  parseFloat(ov.DividendPerShare),
		DividendYield:     parseFloat(ov.DividendYield) * 100,
		EPS:               parseFloat(ov.EPS),
		NetMargin:         parseFloat(ov.ProfitMargin) * 100,
		OperatingMargin:   parseFloat(ov.OperatingMarginTTM) * 100,
		ROA:               parseFloat(ov.ReturnOnAssetsTTM) * 100,
		ROE:               parseFloat(ov.ReturnOnEquityTTM) * 100,
		Revenue:           parseInt64(ov.RevenueTTM),
		GrossProfit:       parseInt64(ov.GrossProfitTTM),
		EPSGrowth:         parseFloat(ov.QuarterlyEarningsGrowthYOY) * 100,
		RevenueGrowth:     parseFloat(ov.QuarterlyRevenueGrowthYOY) * 100,
		Beta:              parseFloat(ov.Beta),
		Week52High:        parseFloat(ov.Week52High),
		Week52Low:         parseFloat(ov.Week52Low),
		MA50:              parseFloat(ov.MA50),
		MA200:             parseFloat(ov.MA200),
		SharesOutstanding: parseInt64(ov.SharesOutstanding),
		LastUpdated:       time.Now(),
	}
}

// Helper functions
func parseFloat(s string) float64 {
	if s == "" || s == "None" || s == "-" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func parseInt64(s string) int64 {
	if s == "" || s == "None" || s == "-" {
		return 0
	}
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
