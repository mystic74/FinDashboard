package services

import (
	"encoding/json"
	"fmt"
	"log"
	"stock-screener/models"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	fmpBaseURL  = "https://financialmodelingprep.com/api/v3"
	fmpProvider = "fmp"
)

// FMPService handles Financial Modeling Prep API interactions
type FMPService struct {
	client      *resty.Client
	apiKey      string
	rateLimiter *RateLimiter
	status      *ProviderStatus
	cache       *CacheService
}

// FMP API response types
type FMPQuote struct {
	Symbol                  string  `json:"symbol"`
	Name                    string  `json:"name"`
	Price                   float64 `json:"price"`
	Change                  float64 `json:"change"`
	ChangePercent           float64 `json:"changesPercentage"`
	DayLow                  float64 `json:"dayLow"`
	DayHigh                 float64 `json:"dayHigh"`
	YearLow                 float64 `json:"yearLow"`
	YearHigh                float64 `json:"yearHigh"`
	MarketCap               int64   `json:"marketCap"`
	PriceAvg50              float64 `json:"priceAvg50"`
	PriceAvg200             float64 `json:"priceAvg200"`
	Volume                  int64   `json:"volume"`
	AvgVolume               int64   `json:"avgVolume"`
	Exchange                string  `json:"exchange"`
	Open                    float64 `json:"open"`
	PreviousClose           float64 `json:"previousClose"`
	EPS                     float64 `json:"eps"`
	PE                      float64 `json:"pe"`
	EarningsAnnouncement    string  `json:"earningsAnnouncement"`
	SharesOutstanding       int64   `json:"sharesOutstanding"`
	Timestamp               int64   `json:"timestamp"`
}

type FMPProfile struct {
	Symbol            string  `json:"symbol"`
	CompanyName       string  `json:"companyName"`
	Currency          string  `json:"currency"`
	Exchange          string  `json:"exchange"`
	Industry          string  `json:"industry"`
	Sector            string  `json:"sector"`
	Country           string  `json:"country"`
	FullTimeEmployees int     `json:"fullTimeEmployees"`
	Price             float64 `json:"price"`
	Beta              float64 `json:"beta"`
	VolAvg            int64   `json:"volAvg"`
	MktCap            int64   `json:"mktCap"`
	LastDiv           float64 `json:"lastDiv"`
	Range             string  `json:"range"`
	DCF               float64 `json:"dcf"`
	IPODate           string  `json:"ipoDate"`
}

type FMPRatios struct {
	Symbol                         string  `json:"symbol"`
	Date                           string  `json:"date"`
	CurrentRatio                   float64 `json:"currentRatio"`
	QuickRatio                     float64 `json:"quickRatio"`
	CashRatio                      float64 `json:"cashRatio"`
	DaysOfSalesOutstanding         float64 `json:"daysOfSalesOutstanding"`
	DaysOfInventoryOutstanding     float64 `json:"daysOfInventoryOutstanding"`
	DaysOfPayablesOutstanding      float64 `json:"daysOfPayablesOutstanding"`
	GrossProfitMargin              float64 `json:"grossProfitMargin"`
	OperatingProfitMargin          float64 `json:"operatingProfitMargin"`
	NetProfitMargin                float64 `json:"netProfitMargin"`
	ReturnOnAssets                 float64 `json:"returnOnAssets"`
	ReturnOnEquity                 float64 `json:"returnOnEquity"`
	ReturnOnCapitalEmployed        float64 `json:"returnOnCapitalEmployed"`
	DebtEquityRatio                float64 `json:"debtEquityRatio"`
	DebtRatio                      float64 `json:"debtRatio"`
	PriceToBookRatio               float64 `json:"priceToBookRatio"`
	PriceToSalesRatio              float64 `json:"priceToSalesRatio"`
	PriceEarningsRatio             float64 `json:"priceEarningsRatio"`
	PriceToFreeCashFlowsRatio      float64 `json:"priceToFreeCashFlowsRatio"`
	PriceToOperatingCashFlowsRatio float64 `json:"priceToOperatingCashFlowsRatio"`
	PriceCashFlowRatio             float64 `json:"priceCashFlowRatio"`
	PriceEarningsToGrowthRatio     float64 `json:"priceEarningsToGrowthRatio"`
	DividendYield                  float64 `json:"dividendYield"`
	PayoutRatio                    float64 `json:"payoutRatio"`
	FreeCashFlowPerShare           float64 `json:"freeCashFlowPerShare"`
	BookValuePerShare              float64 `json:"bookValuePerShare"`
	OperatingCashFlowPerShare      float64 `json:"operatingCashFlowPerShare"`
}

type FMPErrorResponse struct {
	Message string `json:"Error Message"`
}

// NewFMPService creates a new FMP service
// Free tier: 250 calls/day
func NewFMPService(apiKey string, cache *CacheService) *FMPService {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(2)
	client.SetRetryWaitTime(2 * time.Second)

	return &FMPService{
		client:      client,
		apiKey:      apiKey,
		rateLimiter: NewRateLimiter(fmpProvider, 0, 250), // Free tier: 250/day, no minute limit
		status: &ProviderStatus{
			Name:      fmpProvider,
			Available: apiKey != "",
		},
		cache: cache,
	}
}

// IsAvailable returns whether the service is available
func (f *FMPService) IsAvailable() bool {
	return f.apiKey != "" && f.status.IsHealthy() && f.rateLimiter.CanMakeCall()
}

// GetStatus returns the current provider status
func (f *FMPService) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"provider":    fmpProvider,
		"available":   f.IsAvailable(),
		"hasAPIKey":   f.apiKey != "",
		"healthy":     f.status.IsHealthy(),
		"rateLimit":   f.rateLimiter.GetStatus(),
	}
}

// GetQuotes fetches quotes for multiple symbols
func (f *FMPService) GetQuotes(symbols []string) ([]models.Stock, error) {
	if len(symbols) == 0 {
		return []models.Stock{}, nil
	}

	// Check cache first
	cacheKey := fmt.Sprintf("fmp:quotes:%s", strings.Join(symbols, ","))
	if cached, found := f.cache.Get(cacheKey); found {
		if stocks, ok := cached.([]models.Stock); ok {
			return stocks, nil
		}
	}

	// Check rate limit
	if !f.rateLimiter.CanMakeCall() {
		waitTime := f.rateLimiter.GetWaitTime()
		return nil, NewRateLimitError(fmpProvider, waitTime)
	}

	// FMP allows batch quotes
	symbolsStr := strings.Join(symbols, ",")
	url := fmt.Sprintf("%s/quote/%s", fmpBaseURL, symbolsStr)

	resp, err := f.client.R().
		SetQueryParam("apikey", f.apiKey).
		Get(url)

	f.rateLimiter.RecordCall()

	if err != nil {
		f.status.RecordFailure(err)
		return nil, fmt.Errorf("FMP request failed: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode() == 401 {
		f.status.Available = false
		return nil, NewInvalidAPIKeyError(fmpProvider)
	}

	if resp.StatusCode() == 429 {
		f.rateLimiter.SetRetryAfter(time.Minute)
		return nil, NewRateLimitError(fmpProvider, time.Minute)
	}

	if resp.StatusCode() != 200 {
		f.status.RecordFailure(fmt.Errorf("HTTP %d", resp.StatusCode()))
		return nil, NewProviderUnavailableError(fmpProvider, fmt.Sprintf("HTTP %d", resp.StatusCode()))
	}

	// Check for error response
	if err := f.checkErrorResponse(resp.Body()); err != nil {
		return nil, err
	}

	// Parse response
	var quotes []FMPQuote
	if err := json.Unmarshal(resp.Body(), &quotes); err != nil {
		f.status.RecordFailure(err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to Stock models
	stocks := make([]models.Stock, 0, len(quotes))
	for _, q := range quotes {
		stocks = append(stocks, f.convertQuoteToStock(q))
	}

	f.status.RecordSuccess()

	// Cache the result
	f.cache.Set(cacheKey, stocks, 5*time.Minute)

	return stocks, nil
}

// GetProfile fetches company profile for a symbol
func (f *FMPService) GetProfile(symbol string) (*models.Stock, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("fmp:profile:%s", symbol)
	if cached, found := f.cache.Get(cacheKey); found {
		if stock, ok := cached.(*models.Stock); ok {
			return stock, nil
		}
	}

	// Check rate limit
	if !f.rateLimiter.CanMakeCall() {
		waitTime := f.rateLimiter.GetWaitTime()
		return nil, NewRateLimitError(fmpProvider, waitTime)
	}

	url := fmt.Sprintf("%s/profile/%s", fmpBaseURL, symbol)

	resp, err := f.client.R().
		SetQueryParam("apikey", f.apiKey).
		Get(url)

	f.rateLimiter.RecordCall()

	if err != nil {
		f.status.RecordFailure(err)
		return nil, fmt.Errorf("FMP request failed: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode() == 401 {
		f.status.Available = false
		return nil, NewInvalidAPIKeyError(fmpProvider)
	}

	if resp.StatusCode() == 429 {
		f.rateLimiter.SetRetryAfter(time.Minute)
		return nil, NewRateLimitError(fmpProvider, time.Minute)
	}

	// Check for error response
	if err := f.checkErrorResponse(resp.Body()); err != nil {
		return nil, err
	}

	// Parse response (returns array)
	var profiles []FMPProfile
	if err := json.Unmarshal(resp.Body(), &profiles); err != nil {
		f.status.RecordFailure(err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(profiles) == 0 {
		return nil, &APIError{
			Provider:   fmpProvider,
			Code:       "NOT_FOUND",
			Message:    fmt.Sprintf("Symbol '%s' not found", symbol),
			Underlying: ErrSymbolNotFound,
		}
	}

	stock := f.convertProfileToStock(profiles[0])
	f.status.RecordSuccess()

	// Cache the result
	f.cache.Set(cacheKey, stock, 15*time.Minute)

	return stock, nil
}

// GetRatios fetches financial ratios for a symbol
func (f *FMPService) GetRatios(symbol string) (*FMPRatios, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("fmp:ratios:%s", symbol)
	if cached, found := f.cache.Get(cacheKey); found {
		if ratios, ok := cached.(*FMPRatios); ok {
			return ratios, nil
		}
	}

	// Check rate limit
	if !f.rateLimiter.CanMakeCall() {
		waitTime := f.rateLimiter.GetWaitTime()
		return nil, NewRateLimitError(fmpProvider, waitTime)
	}

	url := fmt.Sprintf("%s/ratios-ttm/%s", fmpBaseURL, symbol)

	resp, err := f.client.R().
		SetQueryParam("apikey", f.apiKey).
		Get(url)

	f.rateLimiter.RecordCall()

	if err != nil {
		f.status.RecordFailure(err)
		return nil, fmt.Errorf("FMP request failed: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode() == 401 {
		f.status.Available = false
		return nil, NewInvalidAPIKeyError(fmpProvider)
	}

	if resp.StatusCode() == 429 {
		f.rateLimiter.SetRetryAfter(time.Minute)
		return nil, NewRateLimitError(fmpProvider, time.Minute)
	}

	// Check for error response
	if err := f.checkErrorResponse(resp.Body()); err != nil {
		return nil, err
	}

	// Parse response (returns array)
	var ratios []FMPRatios
	if err := json.Unmarshal(resp.Body(), &ratios); err != nil {
		f.status.RecordFailure(err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(ratios) == 0 {
		return nil, &APIError{
			Provider:   fmpProvider,
			Code:       "NOT_FOUND",
			Message:    fmt.Sprintf("Ratios for '%s' not found", symbol),
			Underlying: ErrSymbolNotFound,
		}
	}

	f.status.RecordSuccess()

	// Cache the result
	f.cache.Set(cacheKey, &ratios[0], 15*time.Minute)

	return &ratios[0], nil
}

// GetStockWithFundamentals fetches a stock with full fundamental data
func (f *FMPService) GetStockWithFundamentals(symbol string) (*models.Stock, error) {
	// Get quote first
	quotes, err := f.GetQuotes([]string{symbol})
	if err != nil {
		return nil, err
	}
	if len(quotes) == 0 {
		return nil, &APIError{
			Provider:   fmpProvider,
			Code:       "NOT_FOUND",
			Message:    fmt.Sprintf("Symbol '%s' not found", symbol),
			Underlying: ErrSymbolNotFound,
		}
	}

	stock := &quotes[0]

	// Get profile for sector/industry
	profile, err := f.GetProfile(symbol)
	if err == nil {
		stock.Sector = profile.Sector
		stock.Industry = profile.Industry
		stock.Country = profile.Country
		stock.Beta = profile.Beta
		if stock.Name == "" {
			stock.Name = profile.Name
		}
	} else {
		log.Printf("[FMP] Failed to get profile for %s: %v", symbol, err)
	}

	// Get ratios for fundamental data
	ratios, err := f.GetRatios(symbol)
	if err == nil {
		stock.CurrentRatio = ratios.CurrentRatio
		stock.QuickRatio = ratios.QuickRatio
		stock.DebtToEquity = ratios.DebtEquityRatio
		stock.GrossMargin = ratios.GrossProfitMargin * 100
		stock.OperatingMargin = ratios.OperatingProfitMargin * 100
		stock.NetMargin = ratios.NetProfitMargin * 100
		stock.ROA = ratios.ReturnOnAssets * 100
		stock.ROE = ratios.ReturnOnEquity * 100
		stock.ROIC = ratios.ReturnOnCapitalEmployed * 100
		stock.PBRatio = ratios.PriceToBookRatio
		stock.PSRatio = ratios.PriceToSalesRatio
		stock.PEGRatio = ratios.PriceEarningsToGrowthRatio
		stock.DividendYield = ratios.DividendYield * 100
		stock.PayoutRatio = ratios.PayoutRatio * 100
		stock.BookValuePerShare = ratios.BookValuePerShare
	} else {
		log.Printf("[FMP] Failed to get ratios for %s: %v", symbol, err)
	}

	return stock, nil
}

// checkErrorResponse checks for API error messages
func (f *FMPService) checkErrorResponse(body []byte) error {
	// Check if response is an error object
	var errResp FMPErrorResponse
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Message != "" {
		msg := strings.ToLower(errResp.Message)

		if strings.Contains(msg, "limit") || strings.Contains(msg, "exceeded") {
			if strings.Contains(msg, "daily") {
				return NewQuotaExceededError(fmpProvider, "daily")
			}
			return NewRateLimitError(fmpProvider, time.Minute)
		}

		if strings.Contains(msg, "invalid") || strings.Contains(msg, "api key") {
			f.status.Available = false
			return NewInvalidAPIKeyError(fmpProvider)
		}

		if strings.Contains(msg, "premium") || strings.Contains(msg, "upgrade") {
			return NewPremiumRequiredError(fmpProvider, "requested endpoint")
		}

		return &APIError{
			Provider: fmpProvider,
			Code:     "API_ERROR",
			Message:  errResp.Message,
		}
	}

	return nil
}

// convertQuoteToStock converts FMP quote to Stock model
func (f *FMPService) convertQuoteToStock(q FMPQuote) models.Stock {
	return models.Stock{
		Symbol:            q.Symbol,
		Name:              q.Name,
		Exchange:          q.Exchange,
		Price:             q.Price,
		Change:            q.Change,
		ChangePercent:     q.ChangePercent,
		Open:              q.Open,
		High:              q.DayHigh,
		Low:               q.DayLow,
		PreviousClose:     q.PreviousClose,
		Volume:            q.Volume,
		AvgVolume:         q.AvgVolume,
		MarketCap:         q.MarketCap,
		Week52High:        q.YearHigh,
		Week52Low:         q.YearLow,
		MA50:              q.PriceAvg50,
		MA200:             q.PriceAvg200,
		EPS:               q.EPS,
		PERatio:           q.PE,
		SharesOutstanding: q.SharesOutstanding,
		LastUpdated:       time.Now(),
	}
}

// convertProfileToStock converts FMP profile to Stock model
func (f *FMPService) convertProfileToStock(p FMPProfile) *models.Stock {
	return &models.Stock{
		Symbol:    p.Symbol,
		Name:      p.CompanyName,
		Exchange:  p.Exchange,
		Currency:  p.Currency,
		Sector:    p.Sector,
		Industry:  p.Industry,
		Country:   p.Country,
		Price:     p.Price,
		Beta:      p.Beta,
		AvgVolume: p.VolAvg,
		MarketCap: p.MktCap,
	}
}
