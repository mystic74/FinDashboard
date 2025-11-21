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
			assert.Less(t, s.PERatio, 15.0, "P/E should be < 15")
			assert.Greater(t, s.PERatio, 0.0, "P/E should be > 0")
			assert.Less(t, s.PBRatio, 1.5, "P/B should be < 1.5")
			assert.Less(t, s.DebtToEquity, 0.5, "D/E should be < 0.5")
			assert.Greater(t, s.ROE, 10.0, "ROE should be > 10%")
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
