package services

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"stock-screener/models"
	"strings"
	"time"
)

// ScreenerEngine handles stock screening operations
type ScreenerEngine struct {
	yahooService  *YahooFinanceService
	dataProvider  *DataProviderManager
	mockService   *MockDataService
	cache         *CacheService
	stockUniverse []string
	demoMode      bool
}

// NewScreenerEngine creates a new screener engine
func NewScreenerEngine(yahooService *YahooFinanceService, cache *CacheService) *ScreenerEngine {
	return &ScreenerEngine{
		yahooService:  yahooService,
		cache:         cache,
		stockUniverse: GetDefaultStockSymbols(),
		demoMode:      false,
	}
}

// NewScreenerEngineWithDemo creates a screener engine in demo mode
func NewScreenerEngineWithDemo(cache *CacheService) *ScreenerEngine {
	return &ScreenerEngine{
		mockService:   NewMockDataService(),
		cache:         cache,
		stockUniverse: GetDefaultStockSymbols(),
		demoMode:      true,
	}
}

// SetDemoMode enables or disables demo mode
func (e *ScreenerEngine) SetDemoMode(enabled bool) {
	e.demoMode = enabled
	if enabled && e.mockService == nil {
		e.mockService = NewMockDataService()
	}
}

// IsDemoMode returns whether demo mode is enabled
func (e *ScreenerEngine) IsDemoMode() bool {
	return e.demoMode
}

// getStocks fetches stocks from the appropriate source
func (e *ScreenerEngine) getStocks(ctx context.Context) ([]models.Stock, error) {
	if e.demoMode {
		return e.mockService.GetAllStocks(), nil
	}

	if e.dataProvider != nil {
		stocks, err := e.dataProvider.GetQuotes(ctx, e.stockUniverse)
		if err == nil {
			return stocks, nil
		}
		log.Printf("[ScreenerEngine] provider manager quote fetch failed: %v", err)
	}
	return e.yahooService.GetQuotes(ctx, e.stockUniverse)
}

// SetStockUniverse updates the stock universe for screening
func (e *ScreenerEngine) SetStockUniverse(symbols []string) {
	e.stockUniverse = symbols
}

// SetDataProviderManager configures resilient quote retrieval with provider fallback.
func (e *ScreenerEngine) SetDataProviderManager(provider *DataProviderManager) {
	e.dataProvider = provider
}

// GetStockUniverse returns the current stock universe
func (e *ScreenerEngine) GetStockUniverse() []string {
	return e.stockUniverse
}

// RunScreener runs a screener and returns matching stocks
func (e *ScreenerEngine) RunScreener(ctx context.Context, screener models.Screener) (*models.ScreenerResult, error) {
	startTime := time.Now()

	// Fetch all stocks
	stocks, err := e.getStocks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stocks: %w", err)
	}

	// Enrich with fundamentals for complex screeners (only in non-demo mode)
	if !e.demoMode {
		needsFundamentals := e.screenerNeedsFundamentals(screener)
		if needsFundamentals && e.yahooService != nil {
			stocks, err = e.yahooService.GetMultipleStocksWithFundamentals(ctx, e.stockUniverse)
			if err != nil {
				// Keep results available even when fundamentals endpoint is blocked.
				log.Printf("[ScreenerEngine] fundamentals fetch failed, continuing with quote-only data: %v", err)
			}
		}
	}

	// Calculate additional metrics
	for i := range stocks {
		stocks[i].PiotroskiFScore = CalculatePiotroskiScore(&stocks[i])
		stocks[i].AltmanZScore = CalculateAltmanZ(&stocks[i])
	}

	// Apply filters
	filteredStocks := e.ApplyFilters(stocks, screener.Filters)

	// Sort results
	if screener.SortBy != "" {
		e.SortStocks(filteredStocks, screener.SortBy, screener.SortOrder)
	}

	executionMs := time.Since(startTime).Milliseconds()

	return &models.ScreenerResult{
		Screener:    screener,
		Stocks:      filteredStocks,
		Total:       len(filteredStocks),
		ExecutionMs: executionMs,
		LastUpdated: time.Now().Format(time.RFC3339),
	}, nil
}

// RunCustomScreener runs a custom screener with arbitrary filters
func (e *ScreenerEngine) RunCustomScreener(ctx context.Context, request models.FilterRequest) (*models.FilterResponse, error) {
	// Fetch stocks
	stocks, err := e.getStocks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stocks: %w", err)
	}

	// Check if we need fundamentals (only in non-demo mode)
	if !e.demoMode {
		needsFundamentals := false
		for _, filter := range request.Filters {
			if e.filterNeedsFundamentals(filter.Field) {
				needsFundamentals = true
				break
			}
		}

		if needsFundamentals && e.yahooService != nil {
			stocks, err = e.yahooService.GetMultipleStocksWithFundamentals(ctx, e.stockUniverse)
			if err != nil {
				// Keep results available even when fundamentals endpoint is blocked.
				log.Printf("[ScreenerEngine] fundamentals fetch failed, continuing with quote-only data: %v", err)
			}
		}
	}

	// Calculate additional metrics
	for i := range stocks {
		stocks[i].PiotroskiFScore = CalculatePiotroskiScore(&stocks[i])
		stocks[i].AltmanZScore = CalculateAltmanZ(&stocks[i])
	}

	// Apply filters
	filteredStocks := e.ApplyFilters(stocks, request.Filters)

	// Sort results
	if request.SortBy != "" {
		e.SortStocks(filteredStocks, request.SortBy, request.SortOrder)
	}

	// Pagination
	total := len(filteredStocks)
	limit := request.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	offset := request.Offset
	if offset < 0 {
		offset = 0
	}

	// Apply pagination
	if offset >= len(filteredStocks) {
		filteredStocks = []models.Stock{}
	} else if offset+limit > len(filteredStocks) {
		filteredStocks = filteredStocks[offset:]
	} else {
		filteredStocks = filteredStocks[offset : offset+limit]
	}

	return &models.FilterResponse{
		Stocks:         filteredStocks,
		Total:          total,
		Page:           offset/limit + 1,
		PageSize:       limit,
		AppliedFilters: request.Filters,
	}, nil
}

// ApplyFilters applies a list of filters to stocks
func (e *ScreenerEngine) ApplyFilters(stocks []models.Stock, filters []models.Filter) []models.Stock {
	if len(filters) == 0 {
		return stocks
	}

	result := make([]models.Stock, 0, len(stocks))
	for _, stock := range stocks {
		if e.matchesAllFilters(&stock, filters) {
			result = append(result, stock)
		}
	}
	return result
}

// matchesAllFilters checks if a stock matches all filters
func (e *ScreenerEngine) matchesAllFilters(stock *models.Stock, filters []models.Filter) bool {
	for _, filter := range filters {
		if !e.matchesFilter(stock, filter) {
			return false
		}
	}
	return true
}

// matchesFilter checks if a stock matches a single filter
func (e *ScreenerEngine) matchesFilter(stock *models.Stock, filter models.Filter) bool {
	// Get the field value using reflection
	value := e.getFieldValue(stock, filter.Field)
	if value == nil {
		return false
	}

	switch filter.Operator {
	case models.OpEquals:
		return e.compareEqual(value, filter.Value)
	case models.OpNotEquals:
		return !e.compareEqual(value, filter.Value)
	case models.OpGreaterThan:
		return e.compareGreater(value, filter.Value, false)
	case models.OpGreaterOrEqual:
		return e.compareGreater(value, filter.Value, true)
	case models.OpLessThan:
		return e.compareLess(value, filter.Value, false)
	case models.OpLessOrEqual:
		return e.compareLess(value, filter.Value, true)
	case models.OpBetween:
		return e.compareBetween(value, filter.Value, filter.Value2)
	case models.OpIn:
		return e.compareIn(value, filter.Value)
	case models.OpNotIn:
		return !e.compareIn(value, filter.Value)
	case models.OpContains:
		return e.compareContains(value, filter.Value)
	default:
		return false
	}
}

// getFieldValue gets a field value from a stock using reflection
func (e *ScreenerEngine) getFieldValue(stock *models.Stock, field string) interface{} {
	v := reflect.ValueOf(stock).Elem()
	t := v.Type()

	// First, try to find by JSON tag
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		jsonTag := f.Tag.Get("json")
		jsonName := strings.Split(jsonTag, ",")[0]
		if jsonName == field {
			return v.Field(i).Interface()
		}
	}

	// Fall back to field name (case-insensitive)
	for i := 0; i < t.NumField(); i++ {
		if strings.EqualFold(t.Field(i).Name, field) {
			return v.Field(i).Interface()
		}
	}

	return nil
}

// Comparison functions

func (e *ScreenerEngine) compareEqual(value, target interface{}) bool {
	v, t := e.toFloat64(value), e.toFloat64(target)
	if v != nil && t != nil {
		return *v == *t
	}
	return fmt.Sprintf("%v", value) == fmt.Sprintf("%v", target)
}

func (e *ScreenerEngine) compareGreater(value, target interface{}, orEqual bool) bool {
	v, t := e.toFloat64(value), e.toFloat64(target)
	if v == nil || t == nil {
		return false
	}
	if orEqual {
		return *v >= *t
	}
	return *v > *t
}

func (e *ScreenerEngine) compareLess(value, target interface{}, orEqual bool) bool {
	v, t := e.toFloat64(value), e.toFloat64(target)
	if v == nil || t == nil {
		return false
	}
	if orEqual {
		return *v <= *t
	}
	return *v < *t
}

func (e *ScreenerEngine) compareBetween(value, min, max interface{}) bool {
	v := e.toFloat64(value)
	minVal := e.toFloat64(min)
	maxVal := e.toFloat64(max)
	if v == nil || minVal == nil || maxVal == nil {
		return false
	}
	return *v >= *minVal && *v <= *maxVal
}

func (e *ScreenerEngine) compareIn(value, targets interface{}) bool {
	// Handle string arrays
	if arr, ok := targets.([]string); ok {
		strVal := fmt.Sprintf("%v", value)
		for _, t := range arr {
			if strVal == t {
				return true
			}
		}
		return false
	}

	// Handle interface arrays
	if arr, ok := targets.([]interface{}); ok {
		strVal := fmt.Sprintf("%v", value)
		for _, t := range arr {
			if strVal == fmt.Sprintf("%v", t) {
				return true
			}
		}
		return false
	}

	return false
}

func (e *ScreenerEngine) compareContains(value, target interface{}) bool {
	strVal := strings.ToLower(fmt.Sprintf("%v", value))
	strTarget := strings.ToLower(fmt.Sprintf("%v", target))
	return strings.Contains(strVal, strTarget)
}

func (e *ScreenerEngine) toFloat64(v interface{}) *float64 {
	var result float64
	switch val := v.(type) {
	case float64:
		result = val
	case float32:
		result = float64(val)
	case int:
		result = float64(val)
	case int64:
		result = float64(val)
	case int32:
		result = float64(val)
	default:
		return nil
	}
	return &result
}

// SortStocks sorts stocks by a field
func (e *ScreenerEngine) SortStocks(stocks []models.Stock, sortBy, sortOrder string) {
	if sortBy == "" {
		return
	}

	ascending := sortOrder == "asc"

	sort.Slice(stocks, func(i, j int) bool {
		vi := e.getFieldValue(&stocks[i], sortBy)
		vj := e.getFieldValue(&stocks[j], sortBy)

		fi := e.toFloat64(vi)
		fj := e.toFloat64(vj)

		// Handle nil values
		if fi == nil && fj == nil {
			return false
		}
		if fi == nil {
			return ascending
		}
		if fj == nil {
			return !ascending
		}

		if ascending {
			return *fi < *fj
		}
		return *fi > *fj
	})
}

// screenerNeedsFundamentals checks if a screener needs fundamental data
func (e *ScreenerEngine) screenerNeedsFundamentals(screener models.Screener) bool {
	for _, filter := range screener.Filters {
		if e.filterNeedsFundamentals(filter.Field) {
			return true
		}
	}
	return false
}

// filterNeedsFundamentals checks if a filter field needs fundamental data
func (e *ScreenerEngine) filterNeedsFundamentals(field string) bool {
	fundamentalFields := map[string]bool{
		"roe": true, "roa": true, "roic": true,
		"grossMargin": true, "operatingMargin": true, "netMargin": true,
		"currentRatio": true, "quickRatio": true, "debtToEquity": true,
		"interestCoverage": true, "altmanZScore": true, "cashToDebt": true,
		"revenueGrowth": true, "epsGrowth": true, "fcfGrowth": true,
		"piotroskiFScore": true, "operatingCashFlow": true, "freeCashFlow": true,
		"consecutiveDivYears": true, "dividendGrowthYears": true,
	}
	return fundamentalFields[field]
}

// GetAllPredefinedScreeners returns all predefined screeners
func (e *ScreenerEngine) GetAllPredefinedScreeners() []models.Screener {
	return models.GetPredefinedScreeners()
}

// GetScreenerByID returns a screener by its ID
func (e *ScreenerEngine) GetScreenerByID(id string) (*models.Screener, bool) {
	for _, s := range models.GetPredefinedScreeners() {
		if s.ID == id {
			return &s, true
		}
	}
	return nil, false
}

// GetScreenersSummary returns summaries of all predefined screeners
func (e *ScreenerEngine) GetScreenersSummary(ctx context.Context) ([]models.ScreenerSummary, error) {
	screeners := models.GetPredefinedScreeners()
	summaries := make([]models.ScreenerSummary, len(screeners))

	for i, s := range screeners {
		summaries[i] = models.ScreenerSummary{
			ID:          s.ID,
			Name:        s.Name,
			Description: s.Description,
			Category:    s.Category,
			Icon:        s.Icon,
		}
	}

	return summaries, nil
}

// GetQuickScreenResults runs a screener without full fundamentals for speed
func (e *ScreenerEngine) GetQuickScreenResults(ctx context.Context, screenerID string) (*models.ScreenerResult, error) {
	screener, found := e.GetScreenerByID(screenerID)
	if !found {
		return nil, fmt.Errorf("screener not found: %s", screenerID)
	}

	return e.RunScreener(ctx, *screener)
}

// GetSectorPerformance returns performance by sector
func (e *ScreenerEngine) GetSectorPerformance(ctx context.Context) ([]models.SectorPerformance, error) {
	stocks, err := e.getStocks(ctx)
	if err != nil {
		return nil, err
	}

	// Group by sector
	sectorMap := make(map[string][]models.Stock)
	for _, stock := range stocks {
		if stock.Sector != "" {
			sectorMap[stock.Sector] = append(sectorMap[stock.Sector], stock)
		}
	}

	var performances []models.SectorPerformance
	for sector, sectorStocks := range sectorMap {
		perf := models.SectorPerformance{
			Sector:     sector,
			StockCount: len(sectorStocks),
		}

		// Calculate averages
		var totalChange, totalMarketCap float64
		var topPerf, worstPerf float64
		var topStock, worstStock string

		for _, s := range sectorStocks {
			totalChange += s.ChangePercent
			totalMarketCap += float64(s.MarketCap)

			if s.ChangePercent > topPerf || topStock == "" {
				topPerf = s.ChangePercent
				topStock = s.Symbol
			}
			if s.ChangePercent < worstPerf || worstStock == "" {
				worstPerf = s.ChangePercent
				worstStock = s.Symbol
			}
		}

		perf.Change1D = totalChange / float64(len(sectorStocks))
		perf.MarketCap = int64(totalMarketCap)
		perf.TopPerformer = topStock
		perf.WorstPerformer = worstStock

		performances = append(performances, perf)
	}

	// Sort by change
	sort.Slice(performances, func(i, j int) bool {
		return performances[i].Change1D > performances[j].Change1D
	})

	return performances, nil
}
