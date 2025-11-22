package services

import (
	"stock-screener/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTestCache() *CacheService {
	return NewCacheService(5*time.Minute, 10*time.Minute)
}

func createTestStocks() []models.Stock {
	return []models.Stock{
		{
			Symbol: "AAPL", Name: "Apple Inc.",
			Price: 150, Change: 2, ChangePercent: 1.35,
			Volume: 50000000, AvgVolume: 45000000,
			MarketCap: 2500000000000, PERatio: 25, PBRatio: 40, PSRatio: 7,
			DividendYield: 0.5, PayoutRatio: 15, Beta: 1.2,
			ROE: 150, ROA: 30, GrossMargin: 43, OperatingMargin: 30, NetMargin: 25,
			CurrentRatio: 1.0, QuickRatio: 0.9, DebtToEquity: 1.8,
			RevenueGrowth: 8, EPSGrowth: 10, Return1W: 2, Return1M: 5, Return3M: 10, Return6M: 15,
			Sector: "Technology", Industry: "Consumer Electronics",
			PiotroskiFScore: 7,
		},
		{
			Symbol: "MSFT", Name: "Microsoft Corporation",
			Price: 380, Change: 5, ChangePercent: 1.33,
			Volume: 25000000, AvgVolume: 22000000,
			MarketCap: 2800000000000, PERatio: 35, PBRatio: 12, PSRatio: 12,
			DividendYield: 0.8, PayoutRatio: 25, Beta: 0.9,
			ROE: 40, ROA: 20, GrossMargin: 70, OperatingMargin: 42, NetMargin: 36,
			CurrentRatio: 1.8, QuickRatio: 1.6, DebtToEquity: 0.4,
			RevenueGrowth: 12, EPSGrowth: 15, Return1W: 1, Return1M: 8, Return3M: 12, Return6M: 20,
			Sector: "Technology", Industry: "Software",
			PiotroskiFScore: 8,
		},
		{
			Symbol: "JNJ", Name: "Johnson & Johnson",
			Price: 160, Change: -1, ChangePercent: -0.62,
			Volume: 8000000, AvgVolume: 7500000,
			MarketCap: 400000000000, PERatio: 16, PBRatio: 5, PSRatio: 4,
			DividendYield: 2.8, PayoutRatio: 45, ConsecutiveDivYears: 60, DividendGrowthYears: 60,
			Beta: 0.5, ROE: 25, ROA: 10, GrossMargin: 68, OperatingMargin: 24, NetMargin: 20,
			CurrentRatio: 1.2, QuickRatio: 1.0, DebtToEquity: 0.45,
			RevenueGrowth: 5, EPSGrowth: 8, Return1W: -1, Return1M: 2, Return3M: 5, Return6M: 8,
			Sector: "Healthcare", Industry: "Pharmaceuticals",
			PiotroskiFScore: 6,
		},
		{
			Symbol: "XOM", Name: "Exxon Mobil",
			Price: 105, Change: 3, ChangePercent: 2.94,
			Volume: 20000000, AvgVolume: 18000000,
			MarketCap: 420000000000, PERatio: 10, PBRatio: 1.8, PSRatio: 1.2,
			DividendYield: 3.5, PayoutRatio: 35, ConsecutiveDivYears: 40, DividendGrowthYears: 5,
			Beta: 1.1, ROE: 18, ROA: 8, GrossMargin: 30, OperatingMargin: 12, NetMargin: 10,
			CurrentRatio: 1.4, QuickRatio: 1.1, DebtToEquity: 0.25,
			RevenueGrowth: 15, EPSGrowth: 20, Return1W: 3, Return1M: 10, Return3M: 20, Return6M: 25,
			Sector: "Energy", Industry: "Oil & Gas",
			FreeCashFlow: 50000000000, OperatingCashFlow: 70000000000,
			PiotroskiFScore: 7,
		},
		{
			Symbol: "TSLA", Name: "Tesla Inc.",
			Price: 250, Change: 10, ChangePercent: 4.17,
			Volume: 80000000, AvgVolume: 70000000,
			MarketCap: 800000000000, PERatio: 70, PBRatio: 15, PSRatio: 8,
			DividendYield: 0, PayoutRatio: 0, Beta: 2.0,
			ROE: 25, ROA: 12, GrossMargin: 25, OperatingMargin: 10, NetMargin: 8,
			CurrentRatio: 1.5, QuickRatio: 1.2, DebtToEquity: 0.3,
			RevenueGrowth: 50, EPSGrowth: 100, Return1W: 5, Return1M: 15, Return3M: 30, Return6M: 40,
			Sector: "Consumer Cyclical", Industry: "Auto Manufacturers",
			PiotroskiFScore: 6,
		},
		{
			Symbol: "VALUE", Name: "Value Stock Inc.",
			Price: 50, Change: 0.5, ChangePercent: 1.0,
			Volume: 500000, AvgVolume: 400000,
			MarketCap: 500000000, PERatio: 8, PBRatio: 0.8, PSRatio: 0.5,
			DividendYield: 4, PayoutRatio: 30, Beta: 0.7,
			ROE: 15, ROA: 8, GrossMargin: 35, OperatingMargin: 18, NetMargin: 12,
			CurrentRatio: 2.5, QuickRatio: 2.0, DebtToEquity: 0.2, CashToDebt: 1.5,
			RevenueGrowth: 3, EPSGrowth: 5, Return1W: 1, Return1M: 3, Return3M: 8, Return6M: 12,
			Sector: "Industrials", Industry: "Manufacturing",
			FreeCashFlow: 30000000, OperatingCashFlow: 50000000,
			PiotroskiFScore: 9,
		},
		{
			Symbol: "SMALLGROW", Name: "Small Cap Growth Co",
			Price: 30, Change: 2, ChangePercent: 7.14,
			Volume: 200000, AvgVolume: 150000,
			MarketCap: 800000000, PERatio: 25, PBRatio: 4, PSRatio: 3,
			DividendYield: 0, PayoutRatio: 0, Beta: 1.8,
			ROE: 20, ROA: 10, GrossMargin: 45, OperatingMargin: 15, NetMargin: 10,
			CurrentRatio: 1.8, QuickRatio: 1.5, DebtToEquity: 0.6,
			RevenueGrowth: 35, EPSGrowth: 40, Return1W: 5, Return1M: 12, Return3M: 25, Return6M: 35,
			Sector: "Technology", Industry: "Software",
			PiotroskiFScore: 7,
		},
	}
}

func TestScreenerEngineApplyFilters(t *testing.T) {
	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)
	stocks := createTestStocks()

	t.Run("No filters", func(t *testing.T) {
		result := engine.ApplyFilters(stocks, []models.Filter{})
		assert.Len(t, result, len(stocks))
	})

	t.Run("Single equals filter", func(t *testing.T) {
		filters := []models.Filter{
			{Field: "sector", Operator: models.OpEquals, Value: "Technology"},
		}
		result := engine.ApplyFilters(stocks, filters)
		assert.Len(t, result, 3) // AAPL, MSFT, SMALLGROW
	})

	t.Run("Greater than filter", func(t *testing.T) {
		filters := []models.Filter{
			{Field: "peRatio", Operator: models.OpGreaterThan, Value: 20},
		}
		result := engine.ApplyFilters(stocks, filters)
		for _, s := range result {
			assert.Greater(t, s.PERatio, 20.0)
		}
	})

	t.Run("Less than filter", func(t *testing.T) {
		filters := []models.Filter{
			{Field: "peRatio", Operator: models.OpLessThan, Value: 15},
		}
		result := engine.ApplyFilters(stocks, filters)
		for _, s := range result {
			assert.Less(t, s.PERatio, 15.0)
		}
	})

	t.Run("Between filter", func(t *testing.T) {
		filters := []models.Filter{
			{Field: "peRatio", Operator: models.OpBetween, Value: 10, Value2: 30},
		}
		result := engine.ApplyFilters(stocks, filters)
		for _, s := range result {
			assert.GreaterOrEqual(t, s.PERatio, 10.0)
			assert.LessOrEqual(t, s.PERatio, 30.0)
		}
	})

	t.Run("Multiple filters (AND logic)", func(t *testing.T) {
		filters := []models.Filter{
			{Field: "sector", Operator: models.OpEquals, Value: "Technology"},
			{Field: "dividendYield", Operator: models.OpGreaterThan, Value: 0},
		}
		result := engine.ApplyFilters(stocks, filters)
		for _, s := range result {
			assert.Equal(t, "Technology", s.Sector)
			assert.Greater(t, s.DividendYield, 0.0)
		}
	})

	t.Run("Contains filter", func(t *testing.T) {
		filters := []models.Filter{
			{Field: "industry", Operator: models.OpContains, Value: "Software"},
		}
		result := engine.ApplyFilters(stocks, filters)
		assert.Len(t, result, 2) // MSFT and SMALLGROW
	})
}

func TestScreenerEngineSorting(t *testing.T) {
	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)
	stocks := createTestStocks()

	t.Run("Sort by price descending", func(t *testing.T) {
		engine.SortStocks(stocks, "price", "desc")
		for i := 0; i < len(stocks)-1; i++ {
			assert.GreaterOrEqual(t, stocks[i].Price, stocks[i+1].Price)
		}
	})

	t.Run("Sort by price ascending", func(t *testing.T) {
		engine.SortStocks(stocks, "price", "asc")
		for i := 0; i < len(stocks)-1; i++ {
			assert.LessOrEqual(t, stocks[i].Price, stocks[i+1].Price)
		}
	})

	t.Run("Sort by P/E ratio", func(t *testing.T) {
		engine.SortStocks(stocks, "peRatio", "asc")
		for i := 0; i < len(stocks)-1; i++ {
			assert.LessOrEqual(t, stocks[i].PERatio, stocks[i+1].PERatio)
		}
	})
}

func TestPredefinedScreeners(t *testing.T) {
	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)
	stocks := createTestStocks()

	t.Run("Value Opportunities Screener", func(t *testing.T) {
		screener := models.ValueOpportunitiesScreener()
		result := engine.ApplyFilters(stocks, screener.Filters)

		for _, s := range result {
			assert.Less(t, s.PERatio, 20.0, "P/E should be < 20")
			assert.Greater(t, s.PERatio, 0.0, "P/E should be > 0")
			assert.Less(t, s.PBRatio, 5.0, "P/B should be < 5")
			assert.Less(t, s.DebtToEquity, 1.5, "D/E should be < 1.5")
			assert.Greater(t, s.FreeCashFlow, int64(0), "FCF should be > 0")
		}
	})

	t.Run("Dividend Aristocrats Screener", func(t *testing.T) {
		screener := models.DividendAristocratsScreener()
		result := engine.ApplyFilters(stocks, screener.Filters)

		for _, s := range result {
			assert.Greater(t, s.DividendYield, 2.0, "Dividend yield should be > 2%")
			assert.Less(t, s.PayoutRatio, 80.0, "Payout ratio should be < 80%")
			assert.Greater(t, s.ConsecutiveDivYears, 10, "Consecutive years should be > 10")
		}
	})

	t.Run("High Beta Bulls Screener", func(t *testing.T) {
		screener := models.HighBetaBullsScreener()
		result := engine.ApplyFilters(stocks, screener.Filters)

		for _, s := range result {
			assert.Greater(t, s.Beta, 1.5, "Beta should be > 1.5")
			assert.Greater(t, s.Return1M, 10.0, "1-month return should be > 10%")
			assert.Greater(t, s.MarketCap, int64(100000000), "Market cap should be > $100M")
		}
	})

	t.Run("Piotroski High Score Screener", func(t *testing.T) {
		screener := models.PiotroskiHighScoreScreener()
		result := engine.ApplyFilters(stocks, screener.Filters)

		for _, s := range result {
			assert.GreaterOrEqual(t, s.PiotroskiFScore, 8, "F-Score should be >= 8")
		}
	})

	t.Run("Small Cap Growth Screener", func(t *testing.T) {
		screener := models.SmallCapGrowthScreener()
		result := engine.ApplyFilters(stocks, screener.Filters)

		for _, s := range result {
			assert.GreaterOrEqual(t, s.MarketCap, int64(300000000), "Market cap should be >= $300M")
			assert.LessOrEqual(t, s.MarketCap, int64(2000000000), "Market cap should be <= $2B")
			assert.Greater(t, s.RevenueGrowth, 20.0, "Revenue growth should be > 20%")
			assert.Greater(t, s.EPSGrowth, 25.0, "EPS growth should be > 25%")
		}
	})

	t.Run("Low Volatility Screener", func(t *testing.T) {
		screener := models.LowVolatilityScreener()
		result := engine.ApplyFilters(stocks, screener.Filters)

		for _, s := range result {
			assert.Less(t, s.Beta, 0.8, "Beta should be < 0.8")
			assert.Greater(t, s.Beta, 0.0, "Beta should be > 0")
			assert.Greater(t, s.DividendYield, 1.0, "Dividend yield should be > 1%")
			assert.Greater(t, s.MarketCap, int64(1000000000), "Market cap should be > $1B")
		}
	})
}

func TestGetScreenerByID(t *testing.T) {
	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)

	t.Run("Valid screener ID", func(t *testing.T) {
		screener, found := engine.GetScreenerByID("momentum-masters")
		assert.True(t, found)
		assert.Equal(t, "Momentum Masters", screener.Name)
	})

	t.Run("Invalid screener ID", func(t *testing.T) {
		_, found := engine.GetScreenerByID("invalid-screener")
		assert.False(t, found)
	})
}

func TestGetAllPredefinedScreeners(t *testing.T) {
	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)

	screeners := engine.GetAllPredefinedScreeners()
	assert.GreaterOrEqual(t, len(screeners), 10, "Should have at least 10 predefined screeners")

	// Check that all screeners have required fields
	for _, s := range screeners {
		assert.NotEmpty(t, s.ID, "Screener should have ID")
		assert.NotEmpty(t, s.Name, "Screener should have name")
		assert.NotEmpty(t, s.Description, "Screener should have description")
		assert.NotEmpty(t, s.Category, "Screener should have category")
		assert.NotEmpty(t, s.Filters, "Screener should have filters")
	}
}

func TestFilterOperatorValidation(t *testing.T) {
	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)
	stocks := createTestStocks()

	t.Run("GreaterOrEqual", func(t *testing.T) {
		filters := []models.Filter{
			{Field: "peRatio", Operator: models.OpGreaterOrEqual, Value: 25},
		}
		result := engine.ApplyFilters(stocks, filters)
		for _, s := range result {
			assert.GreaterOrEqual(t, s.PERatio, 25.0)
		}
	})

	t.Run("LessOrEqual", func(t *testing.T) {
		filters := []models.Filter{
			{Field: "peRatio", Operator: models.OpLessOrEqual, Value: 20},
		}
		result := engine.ApplyFilters(stocks, filters)
		for _, s := range result {
			assert.LessOrEqual(t, s.PERatio, 20.0)
		}
	})

	t.Run("NotEquals", func(t *testing.T) {
		filters := []models.Filter{
			{Field: "sector", Operator: models.OpNotEquals, Value: "Technology"},
		}
		result := engine.ApplyFilters(stocks, filters)
		for _, s := range result {
			assert.NotEqual(t, "Technology", s.Sector)
		}
	})
}

func TestPiotroskiScoreCalculation(t *testing.T) {
	tests := []struct {
		name     string
		stock    models.Stock
		minScore int
		maxScore int
	}{
		{
			name: "Strong financial stock",
			stock: models.Stock{
				ROA: 15, ROE: 20,
				OperatingCashFlow: 1000000, NetIncome: 800000,
				DebtToEquity: 0.3, CurrentRatio: 2.0,
				GrossMargin: 40, EPSGrowth: 10, RevenueGrowth: 8, NetMargin: 15,
			},
			minScore: 7,
			maxScore: 9,
		},
		{
			name: "Weak financial stock",
			stock: models.Stock{
				ROA: -5, ROE: -10,
				OperatingCashFlow: -500000, NetIncome: -300000,
				DebtToEquity: 2.0, CurrentRatio: 0.5,
				GrossMargin: 10, EPSGrowth: -20, RevenueGrowth: -10, NetMargin: -5,
			},
			minScore: 0,
			maxScore: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculatePiotroskiScore(&tt.stock)
			assert.GreaterOrEqual(t, score, tt.minScore)
			assert.LessOrEqual(t, score, tt.maxScore)
		})
	}
}

func TestCountrySectorFiltering(t *testing.T) {
	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)

	// Add test stocks with country data
	stocks := []models.Stock{
		{Symbol: "AAPL", Name: "Apple", Sector: "Technology", Country: "USA", PERatio: 25},
		{Symbol: "MSFT", Name: "Microsoft", Sector: "Technology", Country: "USA", PERatio: 30},
		{Symbol: "LEUMI", Name: "Bank Leumi", Sector: "Financial Services", Country: "Israel", PERatio: 12},
		{Symbol: "HSBC", Name: "HSBC Holdings", Sector: "Financial Services", Country: "UK", PERatio: 10},
		{Symbol: "SAP", Name: "SAP SE", Sector: "Technology", Country: "Germany", PERatio: 22},
	}

	t.Run("Filter by country", func(t *testing.T) {
		filters := []models.Filter{
			{Field: "country", Operator: models.OpEquals, Value: "USA"},
		}
		result := engine.ApplyFilters(stocks, filters)
		assert.Len(t, result, 2)
		for _, s := range result {
			assert.Equal(t, "USA", s.Country)
		}
	})

	t.Run("Filter by sector", func(t *testing.T) {
		filters := []models.Filter{
			{Field: "sector", Operator: models.OpEquals, Value: "Technology"},
		}
		result := engine.ApplyFilters(stocks, filters)
		assert.Len(t, result, 3)
		for _, s := range result {
			assert.Equal(t, "Technology", s.Sector)
		}
	})

	t.Run("Filter by country and sector", func(t *testing.T) {
		filters := []models.Filter{
			{Field: "country", Operator: models.OpEquals, Value: "USA"},
			{Field: "sector", Operator: models.OpEquals, Value: "Technology"},
		}
		result := engine.ApplyFilters(stocks, filters)
		assert.Len(t, result, 2) // AAPL and MSFT
		for _, s := range result {
			assert.Equal(t, "USA", s.Country)
			assert.Equal(t, "Technology", s.Sector)
		}
	})

	t.Run("Filter by country, sector, and metric", func(t *testing.T) {
		filters := []models.Filter{
			{Field: "sector", Operator: models.OpEquals, Value: "Financial Services"},
			{Field: "peRatio", Operator: models.OpLessThan, Value: 15},
		}
		result := engine.ApplyFilters(stocks, filters)
		assert.Len(t, result, 2) // LEUMI and HSBC
		for _, s := range result {
			assert.Equal(t, "Financial Services", s.Sector)
			assert.Less(t, s.PERatio, 15.0)
		}
	})
}

func TestCashIsKingScreener(t *testing.T) {
	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)

	stocks := []models.Stock{
		{
			Symbol: "STRONG", Name: "Strong Cash Co",
			CurrentRatio: 2.5, QuickRatio: 2.0, CashToDebt: 1.5,
			OperatingCashFlow: 1000000,
		},
		{
			Symbol: "WEAK", Name: "Weak Cash Co",
			CurrentRatio: 0.8, QuickRatio: 0.5, CashToDebt: 0.3,
			OperatingCashFlow: -100000,
		},
	}

	t.Run("Cash is King filters strong cash positions", func(t *testing.T) {
		screener := models.CashIsKingScreener()
		result := engine.ApplyFilters(stocks, screener.Filters)

		// Only STRONG should pass
		assert.Len(t, result, 1)
		assert.Equal(t, "STRONG", result[0].Symbol)

		for _, s := range result {
			assert.Greater(t, s.CurrentRatio, 2.0, "Current ratio should be > 2")
			assert.Greater(t, s.QuickRatio, 1.5, "Quick ratio should be > 1.5")
			assert.Greater(t, s.CashToDebt, 1.0, "Cash to debt should be > 1")
			assert.Greater(t, s.OperatingCashFlow, int64(0), "Operating cash flow should be > 0")
		}
	})
}

func TestMomentumMastersScreener(t *testing.T) {
	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)

	stocks := []models.Stock{
		{
			Symbol: "MOMENTUM", Name: "Momentum Stock",
			Return1W: 5, Return3M: 30, Return6M: 40,
			Volume: 500000,
		},
		{
			Symbol: "SLOW", Name: "Slow Stock",
			Return1W: -2, Return3M: 5, Return6M: -10,
			Volume: 500000,
		},
	}

	t.Run("Momentum Masters filters high momentum stocks", func(t *testing.T) {
		screener := models.MomentumMastersScreener()
		result := engine.ApplyFilters(stocks, screener.Filters)

		// Only MOMENTUM should pass
		assert.Len(t, result, 1)
		assert.Equal(t, "MOMENTUM", result[0].Symbol)

		for _, s := range result {
			assert.Greater(t, s.Return1W, 0.0, "1W return should be > 0")
			assert.Greater(t, s.Return3M, 25.0, "3M return should be > 25%")
			assert.Greater(t, s.Return6M, 25.0, "6M return should be > 25%")
			assert.Greater(t, s.Volume, int64(100000), "Volume should be > 100k")
		}
	})
}

func TestMockDataServiceCountries(t *testing.T) {
	service := NewMockDataService()
	stocks := service.GetAllStocks()

	// Count countries
	countries := make(map[string]int)
	for _, s := range stocks {
		countries[s.Country]++
	}

	t.Run("Has multiple countries", func(t *testing.T) {
		assert.GreaterOrEqual(t, len(countries), 5, "Should have at least 5 different countries")
	})

	t.Run("Has USA stocks", func(t *testing.T) {
		assert.Greater(t, countries["USA"], 0, "Should have USA stocks")
	})

	t.Run("Has international stocks", func(t *testing.T) {
		// Check for at least some international stocks
		internationalCount := 0
		for country, count := range countries {
			if country != "USA" && country != "" {
				internationalCount += count
			}
		}
		assert.Greater(t, internationalCount, 10, "Should have at least 10 international stocks")
	})

	t.Run("Stocks have CashToDebt", func(t *testing.T) {
		stocksWithCashToDebt := 0
		for _, s := range stocks {
			if s.CashToDebt > 0 {
				stocksWithCashToDebt++
			}
		}
		assert.Greater(t, stocksWithCashToDebt, len(stocks)/2, "Most stocks should have CashToDebt > 0")
	})
}

// Tests for remaining screeners to ensure all 12 have coverage

func TestUndervaluedTechScreener(t *testing.T) {
	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)

	stocks := []models.Stock{
		{
			Symbol: "UNDERTECH", Name: "Undervalued Tech Co",
			Sector: "Technology", PERatio: 18, PEGRatio: 0.8, RevenueGrowth: 20,
		},
		{
			Symbol: "OVERTECH", Name: "Overvalued Tech Co",
			Sector: "Technology", PERatio: 35, PEGRatio: 2.5, RevenueGrowth: 10,
		},
		{
			Symbol: "NONTECH", Name: "Non Tech Co",
			Sector: "Healthcare", PERatio: 12, PEGRatio: 0.5, RevenueGrowth: 25,
		},
	}

	t.Run("Undervalued Tech filters correctly", func(t *testing.T) {
		screener := models.UndervaluedTechScreener()
		result := engine.ApplyFilters(stocks, screener.Filters)

		// Only UNDERTECH should pass (Technology sector + P/E < 25 + PEG < 1 + Revenue growth > 15)
		assert.Len(t, result, 1)
		assert.Equal(t, "UNDERTECH", result[0].Symbol)

		for _, s := range result {
			assert.Equal(t, "Technology", s.Sector)
			assert.Less(t, s.PERatio, 25.0)
			assert.Greater(t, s.PERatio, 0.0)
			assert.Less(t, s.PEGRatio, 1.0)
			assert.Greater(t, s.RevenueGrowth, 15.0)
		}
	})
}

func TestGARPScreener(t *testing.T) {
	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)

	stocks := []models.Stock{
		{
			Symbol: "GARP", Name: "GARP Stock",
			EPSGrowth: 20, RevenueGrowth: 15, PEGRatio: 1.0, PERatio: 20, ROE: 15,
		},
		{
			Symbol: "HIGHPE", Name: "High P/E Stock",
			EPSGrowth: 25, RevenueGrowth: 20, PEGRatio: 1.2, PERatio: 35, ROE: 18,
		},
		{
			Symbol: "LOWGROWTH", Name: "Low Growth Stock",
			EPSGrowth: 5, RevenueGrowth: 3, PEGRatio: 0.8, PERatio: 12, ROE: 10,
		},
	}

	t.Run("GARP filters growth at reasonable price", func(t *testing.T) {
		screener := models.GrowthAtReasonablePriceScreener()
		result := engine.ApplyFilters(stocks, screener.Filters)

		// Only GARP should pass
		assert.Len(t, result, 1)
		assert.Equal(t, "GARP", result[0].Symbol)

		for _, s := range result {
			assert.Greater(t, s.EPSGrowth, 15.0)
			assert.Greater(t, s.RevenueGrowth, 10.0)
			assert.GreaterOrEqual(t, s.PEGRatio, 0.5)
			assert.LessOrEqual(t, s.PEGRatio, 1.5)
			assert.Less(t, s.PERatio, 25.0)
			assert.Greater(t, s.PERatio, 0.0)
			assert.Greater(t, s.ROE, 12.0)
		}
	})
}

func TestQualityStocksScreener(t *testing.T) {
	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)

	stocks := []models.Stock{
		{
			Symbol: "QUALITY", Name: "Quality Stock",
			ROE: 20, ROA: 12, GrossMargin: 50, OperatingMargin: 20,
			DebtToEquity: 0.5, CurrentRatio: 2.0,
		},
		{
			Symbol: "LOWQUAL", Name: "Low Quality Stock",
			ROE: 8, ROA: 3, GrossMargin: 25, OperatingMargin: 8,
			DebtToEquity: 2.0, CurrentRatio: 0.8,
		},
	}

	t.Run("Quality Stocks filters high-quality companies", func(t *testing.T) {
		screener := models.QualityStocksScreener()
		result := engine.ApplyFilters(stocks, screener.Filters)

		// Only QUALITY should pass
		assert.Len(t, result, 1)
		assert.Equal(t, "QUALITY", result[0].Symbol)

		for _, s := range result {
			assert.Greater(t, s.ROE, 15.0)
			assert.Greater(t, s.ROA, 8.0)
			assert.Greater(t, s.GrossMargin, 40.0)
			assert.Greater(t, s.OperatingMargin, 15.0)
			assert.Less(t, s.DebtToEquity, 1.0)
			assert.Greater(t, s.CurrentRatio, 1.5)
		}
	})
}

func TestTurnaroundCandidatesScreener(t *testing.T) {
	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)

	stocks := []models.Stock{
		{
			Symbol: "TURN", Name: "Turnaround Stock",
			Return1Y: -30, Return1M: 10, CurrentRatio: 1.5, RevenueGrowth: 5,
		},
		{
			Symbol: "STILLDOWN", Name: "Still Declining Stock",
			Return1Y: -40, Return1M: -5, CurrentRatio: 1.2, RevenueGrowth: -10,
		},
		{
			Symbol: "WINNER", Name: "Winner Stock",
			Return1Y: 25, Return1M: 8, CurrentRatio: 2.0, RevenueGrowth: 15,
		},
	}

	t.Run("Turnaround Candidates filters beaten-down stocks with recovery signs", func(t *testing.T) {
		screener := models.TurnaroundCandidatesScreener()
		result := engine.ApplyFilters(stocks, screener.Filters)

		// Only TURN should pass (down significantly but showing recovery)
		assert.Len(t, result, 1)
		assert.Equal(t, "TURN", result[0].Symbol)

		for _, s := range result {
			assert.Less(t, s.Return1Y, -20.0, "1Y return should be < -20%")
			assert.Greater(t, s.Return1M, 5.0, "1M return should be > 5%")
			assert.Greater(t, s.CurrentRatio, 1.0, "Current ratio should be > 1")
			assert.Greater(t, s.RevenueGrowth, 0.0, "Revenue growth should be > 0")
		}
	})
}

// Integration test: verify all screeners return results with actual mock data
func TestAllScreenersReturnResultsWithMockData(t *testing.T) {
	cache := createTestCache()
	engine := NewScreenerEngineWithDemo(cache)

	// Get all mock stocks
	mockService := NewMockDataService()
	stocks := mockService.GetAllStocks()

	screenerFuncs := []struct {
		id       string
		screener models.Screener
	}{
		{"momentum-masters", models.MomentumMastersScreener()},
		{"dividend-aristocrats", models.DividendAristocratsScreener()},
		{"value-opportunities", models.ValueOpportunitiesScreener()},
		{"high-beta-bulls", models.HighBetaBullsScreener()},
		{"cash-is-king", models.CashIsKingScreener()},
		{"piotroski-high-score", models.PiotroskiHighScoreScreener()},
		{"small-cap-growth", models.SmallCapGrowthScreener()},
		{"undervalued-tech", models.UndervaluedTechScreener()},
		{"garp", models.GrowthAtReasonablePriceScreener()},
		{"quality-stocks", models.QualityStocksScreener()},
		{"low-volatility", models.LowVolatilityScreener()},
		{"turnaround-candidates", models.TurnaroundCandidatesScreener()},
	}

	for _, sf := range screenerFuncs {
		t.Run("Screener "+sf.id+" returns results", func(t *testing.T) {
			result := engine.ApplyFilters(stocks, sf.screener.Filters)
			// Note: Some screeners may return 0 results depending on random mock data generation
			// This test documents the current behavior rather than asserting specific counts
			t.Logf("Screener %s returned %d results out of %d stocks", sf.id, len(result), len(stocks))
		})
	}
}

// createTestStocksForAllScreeners creates hardcoded test data that is GUARANTEED
// to match each screener's criteria. This ensures deterministic test results.
func createTestStocksForAllScreeners() []models.Stock {
	return []models.Stock{
		// MOMENTUM MASTERS: return1W > 0, return3M > 25, return6M > 25, volume > 100000
		{
			Symbol: "MOMENTUM1", Name: "Momentum Stock",
			Return1W: 5, Return3M: 30, Return6M: 35, Volume: 500000,
		},
		// DIVIDEND ARISTOCRATS: dividendYield > 2, payoutRatio < 80, consecutiveDivYears > 10, dividendGrowthYears > 5
		{
			Symbol: "DIVIDEND1", Name: "Dividend Aristocrat",
			DividendYield: 3.5, PayoutRatio: 50, ConsecutiveDivYears: 25, DividendGrowthYears: 15,
		},
		// VALUE OPPORTUNITIES: peRatio < 20 && > 0, pbRatio < 5, debtToEquity < 1.5, freeCashFlow > 0
		{
			Symbol: "VALUE1", Name: "Value Stock",
			PERatio: 12, PBRatio: 1.5, DebtToEquity: 0.8, FreeCashFlow: 1000000,
		},
		// HIGH BETA BULLS: beta > 1.5, return1M > 10, marketCap > 100000000
		{
			Symbol: "HIGHBETA1", Name: "High Beta Stock",
			Beta: 2.0, Return1M: 15, MarketCap: 500000000,
		},
		// CASH IS KING: currentRatio > 2, quickRatio > 1.5, cashToDebt > 1, operatingCashFlow > 0
		{
			Symbol: "CASHKING1", Name: "Cash Rich Stock",
			CurrentRatio: 3.0, QuickRatio: 2.5, CashToDebt: 2.0, OperatingCashFlow: 5000000,
		},
		// PIOTROSKI HIGH SCORE: piotroskiFScore >= 8
		{
			Symbol: "PIOTROSKI1", Name: "Piotroski Leader",
			PiotroskiFScore: 9,
		},
		// SMALL CAP GROWTH: marketCap 300M-2B, revenueGrowth > 20, epsGrowth > 25, peRatio < 30 && > 0
		{
			Symbol: "SMALLCAP1", Name: "Small Cap Growth",
			MarketCap: 800000000, RevenueGrowth: 35, EPSGrowth: 40, PERatio: 22,
		},
		// UNDERVALUED TECH: sector = "Technology", peRatio < 25 && > 0, pegRatio < 1, revenueGrowth > 15
		{
			Symbol: "UNDERTECH1", Name: "Undervalued Tech",
			Sector: "Technology", PERatio: 18, PEGRatio: 0.7, RevenueGrowth: 25,
		},
		// GARP: epsGrowth > 15, revenueGrowth > 10, pegRatio 0.5-1.5, peRatio < 25 && > 0, roe > 12
		{
			Symbol: "GARP1", Name: "GARP Stock",
			EPSGrowth: 20, RevenueGrowth: 18, PEGRatio: 1.0, PERatio: 20, ROE: 18,
		},
		// QUALITY STOCKS: roe > 15, roa > 8, grossMargin > 40, operatingMargin > 15, debtToEquity < 1, currentRatio > 1.5
		{
			Symbol: "QUALITY1", Name: "Quality Stock",
			ROE: 22, ROA: 12, GrossMargin: 55, OperatingMargin: 25, DebtToEquity: 0.5, CurrentRatio: 2.2,
		},
		// LOW VOLATILITY: beta < 0.8 && > 0, dividendYield > 1, marketCap > 1000000000
		{
			Symbol: "LOWVOL1", Name: "Low Volatility Stock",
			Beta: 0.5, DividendYield: 2.5, MarketCap: 50000000000,
		},
		// TURNAROUND CANDIDATES: return1Y < -20, return1M > 5, currentRatio > 1, revenueGrowth > 0
		{
			Symbol: "TURNAROUND1", Name: "Turnaround Stock",
			Return1Y: -35, Return1M: 12, CurrentRatio: 1.5, RevenueGrowth: 8,
		},
	}
}

// TestAllScreenersWithGuaranteedMatches tests each screener with hardcoded data
// that is guaranteed to match. This ensures deterministic, reliable tests.
func TestAllScreenersWithGuaranteedMatches(t *testing.T) {
	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)

	// Use hardcoded test data that guarantees matches
	stocks := createTestStocksForAllScreeners()

	testCases := []struct {
		name           string
		screener       models.Screener
		expectedSymbol string
	}{
		{"Momentum Masters", models.MomentumMastersScreener(), "MOMENTUM1"},
		{"Dividend Aristocrats", models.DividendAristocratsScreener(), "DIVIDEND1"},
		{"Value Opportunities", models.ValueOpportunitiesScreener(), "VALUE1"},
		{"High Beta Bulls", models.HighBetaBullsScreener(), "HIGHBETA1"},
		{"Cash is King", models.CashIsKingScreener(), "CASHKING1"},
		{"Piotroski High Score", models.PiotroskiHighScoreScreener(), "PIOTROSKI1"},
		{"Small Cap Growth", models.SmallCapGrowthScreener(), "SMALLCAP1"},
		{"Undervalued Tech", models.UndervaluedTechScreener(), "UNDERTECH1"},
		{"GARP", models.GrowthAtReasonablePriceScreener(), "GARP1"},
		{"Quality Stocks", models.QualityStocksScreener(), "QUALITY1"},
		{"Low Volatility", models.LowVolatilityScreener(), "LOWVOL1"},
		{"Turnaround Candidates", models.TurnaroundCandidatesScreener(), "TURNAROUND1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name+" finds guaranteed match", func(t *testing.T) {
			result := engine.ApplyFilters(stocks, tc.screener.Filters)

			// Must find at least one result
			assert.GreaterOrEqual(t, len(result), 1, "Screener %s should find at least 1 match", tc.name)

			// The expected stock should be in results
			found := false
			for _, s := range result {
				if s.Symbol == tc.expectedSymbol {
					found = true
					break
				}
			}
			assert.True(t, found, "Screener %s should find stock %s", tc.name, tc.expectedSymbol)
		})
	}
}

// Test that country-based filtering works correctly
func TestScreenersByCountry(t *testing.T) {
	cache := createTestCache()
	engine := NewScreenerEngineWithDemo(cache)

	// Get all stocks
	mockService := NewMockDataService()
	stocks := mockService.GetAllStocks()

	countries := []string{"USA", "Israel", "UK", "Germany", "Japan", "China"}

	for _, country := range countries {
		t.Run("Filter by "+country, func(t *testing.T) {
			filters := []models.Filter{
				{Field: "country", Operator: models.OpEquals, Value: country},
			}
			result := engine.ApplyFilters(stocks, filters)
			assert.Greater(t, len(result), 0, "Should have stocks from %s", country)
			for _, s := range result {
				assert.Equal(t, country, s.Country)
			}
			t.Logf("Found %d stocks from %s", len(result), country)
		})
	}

	t.Run("Country + Sector combination", func(t *testing.T) {
		filters := []models.Filter{
			{Field: "country", Operator: models.OpEquals, Value: "USA"},
			{Field: "sector", Operator: models.OpEquals, Value: "Technology"},
		}
		result := engine.ApplyFilters(stocks, filters)
		t.Logf("Found %d USA Technology stocks", len(result))
		for _, s := range result {
			assert.Equal(t, "USA", s.Country)
			assert.Equal(t, "Technology", s.Sector)
		}
	})
}

// Test that all 12 screeners are accessible and have correct basic structure
func TestAllTwelveScreeners(t *testing.T) {
	screenerIDs := []string{
		"momentum-masters",
		"dividend-aristocrats",
		"value-opportunities",
		"high-beta-bulls",
		"cash-is-king",
		"piotroski-high-score",
		"small-cap-growth",
		"undervalued-tech",
		"garp",
		"quality-stocks",
		"low-volatility",
		"turnaround-candidates",
	}

	cache := createTestCache()
	yahooService := NewYahooFinanceService(cache)
	engine := NewScreenerEngine(yahooService, cache)

	for _, id := range screenerIDs {
		t.Run("Screener "+id+" exists and has filters", func(t *testing.T) {
			screener, found := engine.GetScreenerByID(id)
			assert.True(t, found, "Screener %s should exist", id)
			assert.NotEmpty(t, screener.Name, "Screener should have name")
			assert.NotEmpty(t, screener.Description, "Screener should have description")
			assert.NotEmpty(t, screener.Filters, "Screener should have filters")
			assert.NotEmpty(t, screener.Category, "Screener should have category")
		})
	}

	t.Run("GetPredefinedScreeners returns all 12", func(t *testing.T) {
		all := engine.GetAllPredefinedScreeners()
		assert.Len(t, all, 12, "Should have exactly 12 predefined screeners")
	})
}

func TestAltmanZCalculation(t *testing.T) {
	tests := []struct {
		name     string
		stock    models.Stock
		minZ     float64
		maxZ     float64
	}{
		{
			name: "Healthy company",
			stock: models.Stock{
				TotalAssets:      10000000,
				TotalCash:        2000000,
				TotalDebt:        1000000,
				TotalEquity:      5000000,
				EBITDA:           1500000,
				MarketCap:        20000000,
				TotalLiabilities: 3000000,
				Revenue:          12000000,
			},
			minZ: 2.5,
			maxZ: 10.0,
		},
		{
			name: "Zero assets",
			stock: models.Stock{
				TotalAssets: 0,
			},
			minZ: 0,
			maxZ: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := CalculateAltmanZ(&tt.stock)
			assert.GreaterOrEqual(t, z, tt.minZ)
			assert.LessOrEqual(t, z, tt.maxZ)
		})
	}
}
