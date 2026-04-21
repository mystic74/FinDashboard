package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"stock-screener/models"
	"stock-screener/utils"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"
)

// YahooFinanceService handles Yahoo Finance API interactions
type YahooFinanceService struct {
	client        *resty.Client
	cache         *CacheService
	baseURL       string
	cacheTTL      time.Duration
	maxConcurrent int
	// quoteDriver selects implementation: resty (raw Resty HTTP), ffeng (FFengIll yfinance-go), ampyfin (AmpyFin yfinance-go).
	quoteDriver string
}

// NewYahooFinanceService creates a Yahoo Finance service with default quote driver (resty).
func NewYahooFinanceService(cache *CacheService) *YahooFinanceService {
	return NewYahooFinanceServiceWithDriver(cache, "")
}

// NewYahooFinanceServiceWithDriver creates a Yahoo Finance service with an explicit quote driver.
func NewYahooFinanceServiceWithDriver(cache *CacheService, quoteDriver string) *YahooFinanceService {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	d := normalizeYahooQuoteDriver(quoteDriver)

	return &YahooFinanceService{
		client:        client,
		cache:         cache,
		baseURL:       "https://query1.finance.yahoo.com",
		cacheTTL:      5 * time.Minute,
		maxConcurrent: 10,
		quoteDriver:   d,
	}
}

// QuoteDriver returns the active Yahoo quote implementation (resty | ffeng | ampyfin).
func (y *YahooFinanceService) QuoteDriver() string {
	return y.quoteDriver
}

func normalizeYahooQuoteDriver(d string) string {
	d = strings.ToLower(strings.TrimSpace(d))
	if d == "" {
		d = YahooQuoteDriverResty
	}
	switch d {
	case YahooQuoteDriverResty, YahooQuoteDriverFFeng, YahooQuoteDriverAmpyFin:
		return d
	default:
		return YahooQuoteDriverResty
	}
}

// YahooQuoteResponse represents Yahoo Finance quote API response
type YahooQuoteResponse struct {
	QuoteResponse struct {
		Result []YahooQuote `json:"result"`
		Error  interface{}  `json:"error"`
	} `json:"quoteResponse"`
}

// YahooQuote represents a single quote from Yahoo Finance
type YahooQuote struct {
	Symbol                       string  `json:"symbol"`
	ShortName                    string  `json:"shortName"`
	LongName                     string  `json:"longName"`
	Exchange                     string  `json:"exchange"`
	Currency                     string  `json:"currency"`
	RegularMarketPrice           float64 `json:"regularMarketPrice"`
	RegularMarketChange          float64 `json:"regularMarketChange"`
	RegularMarketChangePercent   float64 `json:"regularMarketChangePercent"`
	RegularMarketVolume          int64   `json:"regularMarketVolume"`
	RegularMarketOpen            float64 `json:"regularMarketOpen"`
	RegularMarketDayHigh         float64 `json:"regularMarketDayHigh"`
	RegularMarketDayLow          float64 `json:"regularMarketDayLow"`
	RegularMarketPreviousClose   float64 `json:"regularMarketPreviousClose"`
	MarketCap                    int64   `json:"marketCap"`
	SharesOutstanding            int64   `json:"sharesOutstanding"`
	FiftyTwoWeekHigh             float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow              float64 `json:"fiftyTwoWeekLow"`
	FiftyDayAverage              float64 `json:"fiftyDayAverage"`
	TwoHundredDayAverage         float64 `json:"twoHundredDayAverage"`
	AverageDailyVolume3Month     int64   `json:"averageDailyVolume3Month"`
	AverageDailyVolume10Day      int64   `json:"averageDailyVolume10Day"`
	TrailingPE                   float64 `json:"trailingPE"`
	ForwardPE                    float64 `json:"forwardPE"`
	PriceToBook                  float64 `json:"priceToBook"`
	TrailingAnnualDividendYield  float64 `json:"trailingAnnualDividendYield"`
	TrailingAnnualDividendRate   float64 `json:"trailingAnnualDividendRate"`
	DividendYield                float64 `json:"dividendYield"`
	PayoutRatio                  float64 `json:"payoutRatio"`
	Beta                         float64 `json:"beta"`
	EpsTrailingTwelveMonths      float64 `json:"epsTrailingTwelveMonths"`
	EpsForward                   float64 `json:"epsForward"`
	BookValue                    float64 `json:"bookValue"`
	PriceToSalesTrailing12Months float64 `json:"priceToSalesTrailing12Months"`
	EnterpriseToRevenue          float64 `json:"enterpriseToRevenue"`
	EnterpriseToEbitda           float64 `json:"enterpriseToEbitda"`
	QuoteType                    string  `json:"quoteType"`
	Sector                       string  `json:"sector,omitempty"`
	Industry                     string  `json:"industry,omitempty"`
}

// GetQuotes fetches quotes for multiple symbols using the configured quote driver.
func (y *YahooFinanceService) GetQuotes(ctx context.Context, symbols []string) ([]models.Stock, error) {
	if len(symbols) == 0 {
		return []models.Stock{}, nil
	}

	// Check cache first
	cacheKey := fmt.Sprintf("quotes:%s:%s", y.quoteDriver, strings.Join(symbols, ","))
	if cached, found := y.cache.Get(cacheKey); found {
		if stocks, ok := cached.([]models.Stock); ok {
			return stocks, nil
		}
	}

	var allStocks []models.Stock
	var err error
	switch y.quoteDriver {
	case YahooQuoteDriverResty:
		allStocks, err = y.getQuotesResty(ctx, symbols)
	case YahooQuoteDriverAmpyFin:
		allStocks, err = y.getQuotesAmpyFin(ctx, symbols)
	default:
		allStocks, err = y.getQuotesFFeng(ctx, symbols)
	}
	if err != nil {
		return nil, err
	}
	if err := ValidateStockQuotes(allStocks); err != nil {
		return nil, err
	}

	y.cache.Set(cacheKey, allStocks, y.cacheTTL)
	return allStocks, nil
}

// fetchQuoteBatch fetches a batch of quotes
func (y *YahooFinanceService) fetchQuoteBatch(symbols []string) ([]models.Stock, error) {
	symbolsStr := strings.Join(symbols, ",")
	apiURL := fmt.Sprintf("%s/v7/finance/quote?symbols=%s", y.baseURL, url.QueryEscape(symbolsStr))

	resp, err := y.client.R().Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quotes: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("yahoo finance returned status %d", resp.StatusCode())
	}

	var quoteResp YahooQuoteResponse
	if err := json.Unmarshal(resp.Body(), &quoteResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	stocks := make([]models.Stock, 0, len(quoteResp.QuoteResponse.Result))
	for _, q := range quoteResp.QuoteResponse.Result {
		stock := y.convertQuoteToStock(q)
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// convertQuoteToStock converts a Yahoo quote to our Stock model
func (y *YahooFinanceService) convertQuoteToStock(q YahooQuote) models.Stock {
	name := q.LongName
	if name == "" {
		name = q.ShortName
	}

	dividendYield := q.TrailingAnnualDividendYield * 100
	if dividendYield == 0 && q.DividendYield > 0 {
		dividendYield = q.DividendYield
	}

	return models.Stock{
		Symbol:            q.Symbol,
		Name:              name,
		Exchange:          q.Exchange,
		Currency:          q.Currency,
		Price:             q.RegularMarketPrice,
		Open:              q.RegularMarketOpen,
		High:              q.RegularMarketDayHigh,
		Low:               q.RegularMarketDayLow,
		PreviousClose:     q.RegularMarketPreviousClose,
		Change:            q.RegularMarketChange,
		ChangePercent:     q.RegularMarketChangePercent,
		Volume:            q.RegularMarketVolume,
		AvgVolume:         q.AverageDailyVolume3Month,
		AvgVolume10Day:    q.AverageDailyVolume10Day,
		MarketCap:         q.MarketCap,
		SharesOutstanding: q.SharesOutstanding,
		Week52High:        q.FiftyTwoWeekHigh,
		Week52Low:         q.FiftyTwoWeekLow,
		MA50:              q.FiftyDayAverage,
		MA200:             q.TwoHundredDayAverage,
		PERatio:           q.TrailingPE,
		ForwardPE:         q.ForwardPE,
		PBRatio:           q.PriceToBook,
		PSRatio:           q.PriceToSalesTrailing12Months,
		EVToEBITDA:        q.EnterpriseToEbitda,
		EVToRevenue:       q.EnterpriseToRevenue,
		DividendYield:     dividendYield,
		DividendPerShare:  q.TrailingAnnualDividendRate,
		PayoutRatio:       q.PayoutRatio,
		Beta:              q.Beta,
		EPS:               q.EpsTrailingTwelveMonths,
		BookValuePerShare: q.BookValue,
		Sector:            q.Sector,
		Industry:          q.Industry,
		LastUpdated:       time.Now(),
	}
}

// GetStockFundamentals fetches detailed fundamentals for a single stock
func (y *YahooFinanceService) GetStockFundamentals(symbol string) (*models.Stock, error) {
	cacheKey := fmt.Sprintf("fundamentals:%s", symbol)
	if cached, found := y.cache.Get(cacheKey); found {
		if stock, ok := cached.(*models.Stock); ok {
			return stock, nil
		}
	}

	// Fetch quote summary with all modules
	modules := []string{
		"financialData",
		"defaultKeyStatistics",
		"summaryDetail",
		"price",
		"summaryProfile",
		"incomeStatementHistory",
		"balanceSheetHistory",
		"cashflowStatementHistory",
	}

	apiURL := fmt.Sprintf("%s/v10/finance/quoteSummary/%s?modules=%s",
		"https://query2.finance.yahoo.com",
		url.QueryEscape(symbol),
		strings.Join(modules, ","))

	resp, err := y.client.R().Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fundamentals: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("yahoo finance returned status %d", resp.StatusCode())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	stock, err := y.parseFundamentals(symbol, result)
	if err != nil {
		return nil, err
	}

	y.cache.Set(cacheKey, stock, y.cacheTTL)
	return stock, nil
}

// parseFundamentals extracts fundamental data from Yahoo response
func (y *YahooFinanceService) parseFundamentals(symbol string, data map[string]interface{}) (*models.Stock, error) {
	stock := &models.Stock{
		Symbol:      symbol,
		LastUpdated: time.Now(),
	}

	quoteSummary, ok := data["quoteSummary"].(map[string]interface{})
	if !ok {
		return stock, nil
	}

	results, ok := quoteSummary["result"].([]interface{})
	if !ok || len(results) == 0 {
		return stock, nil
	}

	result := results[0].(map[string]interface{})

	// Parse price data
	if priceData, ok := result["price"].(map[string]interface{}); ok {
		stock.Name = getStringValue(priceData, "longName")
		if stock.Name == "" {
			stock.Name = getStringValue(priceData, "shortName")
		}
		stock.Price = getRawValue(priceData, "regularMarketPrice")
		stock.Change = getRawValue(priceData, "regularMarketChange")
		stock.ChangePercent = getRawValue(priceData, "regularMarketChangePercent") * 100
		stock.MarketCap = int64(getRawValue(priceData, "marketCap"))
		stock.Exchange = getStringValue(priceData, "exchange")
		stock.Currency = getStringValue(priceData, "currency")
	}

	// Parse summary profile
	if profile, ok := result["summaryProfile"].(map[string]interface{}); ok {
		stock.Sector = getStringValue(profile, "sector")
		stock.Industry = getStringValue(profile, "industry")
		stock.Country = getStringValue(profile, "country")
	}

	// Parse financial data
	if finData, ok := result["financialData"].(map[string]interface{}); ok {
		stock.ROE = getRawValue(finData, "returnOnEquity") * 100
		stock.ROA = getRawValue(finData, "returnOnAssets") * 100
		stock.GrossMargin = getRawValue(finData, "grossMargins") * 100
		stock.OperatingMargin = getRawValue(finData, "operatingMargins") * 100
		stock.NetMargin = getRawValue(finData, "profitMargins") * 100
		stock.CurrentRatio = getRawValue(finData, "currentRatio")
		stock.QuickRatio = getRawValue(finData, "quickRatio")
		stock.DebtToEquity = getRawValue(finData, "debtToEquity")
		stock.RevenueGrowth = getRawValue(finData, "revenueGrowth") * 100
		stock.EPSGrowth = getRawValue(finData, "earningsGrowth") * 100
		stock.FreeCashFlow = int64(getRawValue(finData, "freeCashflow"))
		stock.OperatingCashFlow = int64(getRawValue(finData, "operatingCashflow"))
		stock.Revenue = int64(getRawValue(finData, "totalRevenue"))
		stock.TotalDebt = int64(getRawValue(finData, "totalDebt"))
		stock.TotalCash = int64(getRawValue(finData, "totalCash"))
	}

	// Parse key statistics
	if keyStats, ok := result["defaultKeyStatistics"].(map[string]interface{}); ok {
		stock.PERatio = getRawValue(keyStats, "trailingPE")
		stock.ForwardPE = getRawValue(keyStats, "forwardPE")
		stock.PEGRatio = getRawValue(keyStats, "pegRatio")
		stock.PBRatio = getRawValue(keyStats, "priceToBook")
		stock.Beta = getRawValue(keyStats, "beta")
		stock.SharesOutstanding = int64(getRawValue(keyStats, "sharesOutstanding"))
		stock.Float = int64(getRawValue(keyStats, "floatShares"))
		stock.BookValuePerShare = getRawValue(keyStats, "bookValue")
		stock.EPS = getRawValue(keyStats, "trailingEps")
		stock.NetIncome = int64(getRawValue(keyStats, "netIncomeToCommon"))
		stock.EBITDA = int64(getRawValue(keyStats, "ebitda"))
		stock.TotalAssets = int64(getRawValue(keyStats, "totalAssets"))
	}

	// Parse summary detail
	if summaryDetail, ok := result["summaryDetail"].(map[string]interface{}); ok {
		stock.DividendYield = getRawValue(summaryDetail, "dividendYield") * 100
		stock.DividendPerShare = getRawValue(summaryDetail, "dividendRate")
		stock.PayoutRatio = getRawValue(summaryDetail, "payoutRatio") * 100
		stock.Week52High = getRawValue(summaryDetail, "fiftyTwoWeekHigh")
		stock.Week52Low = getRawValue(summaryDetail, "fiftyTwoWeekLow")
		stock.MA50 = getRawValue(summaryDetail, "fiftyDayAverage")
		stock.MA200 = getRawValue(summaryDetail, "twoHundredDayAverage")
		stock.Volume = int64(getRawValue(summaryDetail, "volume"))
		stock.AvgVolume = int64(getRawValue(summaryDetail, "averageVolume"))
		stock.PreviousClose = getRawValue(summaryDetail, "previousClose")
	}

	// Calculate additional metrics
	if stock.TotalCash > 0 && stock.TotalDebt > 0 {
		stock.CashToDebt = float64(stock.TotalCash) / float64(stock.TotalDebt)
	}

	return stock, nil
}

// GetHistoricalPrices fetches historical price data
func (y *YahooFinanceService) GetHistoricalPrices(symbol string, period string) ([]models.HistoricalPrice, error) {
	cacheKey := fmt.Sprintf("history:%s:%s", symbol, period)
	if cached, found := y.cache.Get(cacheKey); found {
		if prices, ok := cached.([]models.HistoricalPrice); ok {
			return prices, nil
		}
	}

	// Period can be: 1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, ytd, max
	apiURL := fmt.Sprintf("%s/v8/finance/chart/%s?period1=0&period2=9999999999&interval=1d&range=%s",
		y.baseURL, url.QueryEscape(symbol), period)

	resp, err := y.client.R().Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical prices: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	prices := y.parseHistoricalPrices(result)
	y.cache.Set(cacheKey, prices, y.cacheTTL)

	return prices, nil
}

// parseHistoricalPrices extracts price history from Yahoo response
func (y *YahooFinanceService) parseHistoricalPrices(data map[string]interface{}) []models.HistoricalPrice {
	var prices []models.HistoricalPrice

	chart, ok := data["chart"].(map[string]interface{})
	if !ok {
		return prices
	}

	results, ok := chart["result"].([]interface{})
	if !ok || len(results) == 0 {
		return prices
	}

	result := results[0].(map[string]interface{})

	timestamps, ok := result["timestamp"].([]interface{})
	if !ok {
		return prices
	}

	indicators, ok := result["indicators"].(map[string]interface{})
	if !ok {
		return prices
	}

	quote, ok := indicators["quote"].([]interface{})
	if !ok || len(quote) == 0 {
		return prices
	}

	quoteData := quote[0].(map[string]interface{})
	opens := quoteData["open"].([]interface{})
	highs := quoteData["high"].([]interface{})
	lows := quoteData["low"].([]interface{})
	closes := quoteData["close"].([]interface{})
	volumes := quoteData["volume"].([]interface{})

	adjClose, hasAdj := indicators["adjclose"].([]interface{})
	var adjCloses []interface{}
	if hasAdj && len(adjClose) > 0 {
		adjCloses = adjClose[0].(map[string]interface{})["adjclose"].([]interface{})
	}

	for i, ts := range timestamps {
		if opens[i] == nil || closes[i] == nil {
			continue
		}

		price := models.HistoricalPrice{
			Date:   time.Unix(int64(ts.(float64)), 0),
			Open:   opens[i].(float64),
			High:   highs[i].(float64),
			Low:    lows[i].(float64),
			Close:  closes[i].(float64),
			Volume: int64(volumes[i].(float64)),
		}

		if hasAdj && i < len(adjCloses) && adjCloses[i] != nil {
			price.AdjClose = adjCloses[i].(float64)
		} else {
			price.AdjClose = price.Close
		}

		prices = append(prices, price)
	}

	return prices
}

// CalculateTechnicalIndicators calculates RSI, MACD, and returns for a stock
func (y *YahooFinanceService) CalculateTechnicalIndicators(stock *models.Stock, prices []models.HistoricalPrice) {
	if len(prices) < 26 {
		return
	}

	// Extract close prices
	closes := make([]float64, len(prices))
	for i, p := range prices {
		closes[i] = p.AdjClose
	}

	// Calculate RSI
	stock.RSI14 = utils.CalculateRSI(closes, 14)

	// Calculate MACD
	stock.MACD, stock.MACDSignal, stock.MACDHistogram = utils.CalculateMACD(closes)

	// Calculate returns
	if len(prices) >= 5 {
		stock.Return1W = utils.CalculatePercentageReturn(prices[len(prices)-6].AdjClose, prices[len(prices)-1].AdjClose)
	}
	if len(prices) >= 22 {
		stock.Return1M = utils.CalculatePercentageReturn(prices[len(prices)-23].AdjClose, prices[len(prices)-1].AdjClose)
	}
	if len(prices) >= 66 {
		stock.Return3M = utils.CalculatePercentageReturn(prices[len(prices)-67].AdjClose, prices[len(prices)-1].AdjClose)
	}
	if len(prices) >= 132 {
		stock.Return6M = utils.CalculatePercentageReturn(prices[len(prices)-133].AdjClose, prices[len(prices)-1].AdjClose)
	}
	if len(prices) >= 252 {
		stock.Return1Y = utils.CalculatePercentageReturn(prices[len(prices)-253].AdjClose, prices[len(prices)-1].AdjClose)
	}
}

// GetMultipleStocksWithFundamentals fetches multiple stocks with full fundamentals
func (y *YahooFinanceService) GetMultipleStocksWithFundamentals(ctx context.Context, symbols []string) ([]models.Stock, error) {
	// First get basic quotes
	stocks, err := y.GetQuotes(ctx, symbols)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup
	stockMap := make(map[string]*models.Stock)
	for i := range stocks {
		stockMap[stocks[i].Symbol] = &stocks[i]
	}

	// Fetch fundamentals concurrently
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(y.maxConcurrent)

	for _, symbol := range symbols {
		sym := symbol
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				fund, err := y.GetStockFundamentals(sym)
				if err != nil {
					// Log error but don't fail entire operation
					return nil
				}
				if stock, ok := stockMap[sym]; ok {
					// Merge fundamental data
					y.mergeStockData(stock, fund)
				}
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return stocks, nil
}

// mergeStockData merges fundamental data into stock
func (y *YahooFinanceService) mergeStockData(target *models.Stock, source *models.Stock) {
	if source.ROE != 0 {
		target.ROE = source.ROE
	}
	if source.ROA != 0 {
		target.ROA = source.ROA
	}
	if source.ROIC != 0 {
		target.ROIC = source.ROIC
	}
	if source.GrossMargin != 0 {
		target.GrossMargin = source.GrossMargin
	}
	if source.OperatingMargin != 0 {
		target.OperatingMargin = source.OperatingMargin
	}
	if source.NetMargin != 0 {
		target.NetMargin = source.NetMargin
	}
	if source.CurrentRatio != 0 {
		target.CurrentRatio = source.CurrentRatio
	}
	if source.QuickRatio != 0 {
		target.QuickRatio = source.QuickRatio
	}
	if source.DebtToEquity != 0 {
		target.DebtToEquity = source.DebtToEquity
	}
	if source.RevenueGrowth != 0 {
		target.RevenueGrowth = source.RevenueGrowth
	}
	if source.EPSGrowth != 0 {
		target.EPSGrowth = source.EPSGrowth
	}
	if source.FreeCashFlow != 0 {
		target.FreeCashFlow = source.FreeCashFlow
	}
	if source.OperatingCashFlow != 0 {
		target.OperatingCashFlow = source.OperatingCashFlow
	}
	if source.TotalDebt != 0 {
		target.TotalDebt = source.TotalDebt
	}
	if source.TotalCash != 0 {
		target.TotalCash = source.TotalCash
	}
	if source.CashToDebt != 0 {
		target.CashToDebt = source.CashToDebt
	}
	if source.Sector != "" {
		target.Sector = source.Sector
	}
	if source.Industry != "" {
		target.Industry = source.Industry
	}
	if source.Country != "" {
		target.Country = source.Country
	}
}

// Helper functions

func getRawValue(data map[string]interface{}, key string) float64 {
	if val, ok := data[key].(map[string]interface{}); ok {
		if raw, ok := val["raw"].(float64); ok {
			return raw
		}
	}
	if val, ok := data[key].(float64); ok {
		return val
	}
	return 0
}

func getStringValue(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

// GetDefaultStockSymbols returns a list of common stock symbols
func GetDefaultStockSymbols() []string {
	return buildDefaultStockUniverse()
}

// CalculatePiotroskiScore calculates F-Score for a stock based on available data
func CalculatePiotroskiScore(stock *models.Stock) int {
	score := 0

	// Profitability
	if stock.ROA > 0 {
		score++
	}
	if stock.OperatingCashFlow > 0 {
		score++
	}
	// Quality of earnings: CFO > Net Income
	if stock.OperatingCashFlow > stock.NetIncome {
		score++
	}

	// Leverage
	if stock.DebtToEquity < 0.5 && stock.DebtToEquity >= 0 {
		score++
	}
	if stock.CurrentRatio > 1 {
		score++
	}

	// Operating efficiency
	if stock.GrossMargin > 20 {
		score++
	}

	// For remaining criteria, use reasonable defaults
	// ROA increasing, no new shares, asset turnover - add if positive metrics
	if stock.EPSGrowth > 0 {
		score++
	}
	if stock.RevenueGrowth > 0 {
		score++
	}
	if stock.NetMargin > 0 {
		score++
	}

	// Cap at 9
	if score > 9 {
		score = 9
	}

	return score
}

// CalculateAltmanZ calculates Z-Score for a stock
func CalculateAltmanZ(stock *models.Stock) float64 {
	if stock.TotalAssets == 0 {
		return 0
	}

	// Simplified Z-Score calculation with available data
	workingCapital := float64(stock.TotalCash - stock.TotalDebt)
	x1 := workingCapital / float64(stock.TotalAssets)

	// Retained earnings approximation
	x2 := 0.0
	if stock.TotalEquity > 0 {
		x2 = float64(stock.TotalEquity) * 0.5 / float64(stock.TotalAssets)
	}

	// EBIT approximation
	x3 := 0.0
	if stock.EBITDA > 0 {
		x3 = float64(stock.EBITDA) * 0.8 / float64(stock.TotalAssets)
	}

	// Market cap / Total liabilities
	x4 := 0.0
	if stock.TotalLiabilities > 0 {
		x4 = float64(stock.MarketCap) / float64(stock.TotalLiabilities)
	}

	// Revenue / Total assets
	x5 := 0.0
	if stock.Revenue > 0 {
		x5 = float64(stock.Revenue) / float64(stock.TotalAssets)
	}

	zScore := 1.2*x1 + 1.4*x2 + 3.3*x3 + 0.6*x4 + 1.0*x5

	return math.Round(zScore*100) / 100
}
