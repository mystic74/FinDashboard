package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStockValidation(t *testing.T) {
	t.Run("Valid stock", func(t *testing.T) {
		stock := Stock{
			Symbol:        "AAPL",
			Name:          "Apple Inc.",
			Price:         150.0,
			Volume:        50000000,
			MarketCap:     2500000000000,
			PERatio:       25,
			PBRatio:       40,
			DividendYield: 0.5,
			GrossMargin:   43,
			NetMargin:     25,
			CurrentRatio:  1.0,
			Beta:          1.2,
			Week52High:    180,
			Week52Low:     120,
			RSI14:         55,
			LastUpdated:   time.Now(),
		}

		result := stock.Validate()
		assert.True(t, result.IsValid)
		assert.Empty(t, result.Errors)
		assert.Equal(t, 100.0, result.DataQuality)
	})

	t.Run("Invalid price", func(t *testing.T) {
		stock := Stock{Price: 0}
		result := stock.Validate()
		assert.False(t, result.IsValid)
		assert.Contains(t, result.Errors, "price must be positive")
	})

	t.Run("Negative price", func(t *testing.T) {
		stock := Stock{Price: -10}
		result := stock.Validate()
		assert.False(t, result.IsValid)
	})

	t.Run("Negative volume", func(t *testing.T) {
		stock := Stock{Price: 100, Volume: -1000}
		result := stock.Validate()
		assert.False(t, result.IsValid)
		assert.Contains(t, result.Errors, "volume cannot be negative")
	})

	t.Run("Negative market cap", func(t *testing.T) {
		stock := Stock{Price: 100, MarketCap: -1000000}
		result := stock.Validate()
		assert.False(t, result.IsValid)
		assert.Contains(t, result.Errors, "market cap cannot be negative")
	})

	t.Run("Invalid RSI", func(t *testing.T) {
		stock := Stock{Price: 100, RSI14: 150}
		result := stock.Validate()
		assert.False(t, result.IsValid)
		assert.Contains(t, result.Errors, "RSI must be between 0 and 100")
	})

	t.Run("Negative RSI", func(t *testing.T) {
		stock := Stock{Price: 100, RSI14: -10}
		result := stock.Validate()
		assert.False(t, result.IsValid)
	})

	t.Run("52-week inconsistency", func(t *testing.T) {
		stock := Stock{
			Price:      100,
			Week52High: 80,
			Week52Low:  120, // Low > High is invalid
		}
		result := stock.Validate()
		assert.False(t, result.IsValid)
		assert.Contains(t, result.Errors, "52-week low cannot be greater than 52-week high")
	})

	t.Run("Negative current ratio warning", func(t *testing.T) {
		stock := Stock{
			Price:        100,
			CurrentRatio: -1,
		}
		result := stock.Validate()
		assert.False(t, result.IsValid)
		assert.Contains(t, result.Errors, "current ratio cannot be negative")
	})

	t.Run("High dividend yield warning", func(t *testing.T) {
		stock := Stock{
			Price:         100,
			DividendYield: 50, // 50% yield is suspicious
		}
		result := stock.Validate()
		assert.True(t, result.IsValid)
		assert.NotEmpty(t, result.Warnings)
		assert.Less(t, result.DataQuality, 100.0)
	})

	t.Run("Extreme P/E ratio warning", func(t *testing.T) {
		stock := Stock{
			Price:   100,
			PERatio: 5000, // Extremely high P/E
		}
		result := stock.Validate()
		assert.True(t, result.IsValid)
		assert.NotEmpty(t, result.Warnings)
	})

	t.Run("Negative P/E ratio warning", func(t *testing.T) {
		stock := Stock{
			Price:   100,
			PERatio: -5000, // Extremely negative
		}
		result := stock.Validate()
		assert.True(t, result.IsValid) // Just a warning
		assert.NotEmpty(t, result.Warnings)
	})

	t.Run("Unusual price warning", func(t *testing.T) {
		stock := Stock{
			Price: 1500000, // Over $1M per share
		}
		result := stock.Validate()
		assert.True(t, result.IsValid)
		assert.NotEmpty(t, result.Warnings)
		assert.Less(t, result.DataQuality, 100.0)
	})

	t.Run("Negative P/B ratio warning", func(t *testing.T) {
		stock := Stock{
			Price:   100,
			PBRatio: -2,
		}
		result := stock.Validate()
		assert.True(t, result.IsValid)
		assert.NotEmpty(t, result.Warnings)
	})

	t.Run("Price above 52-week high warning", func(t *testing.T) {
		stock := Stock{
			Price:      150,
			Week52High: 100,
			Week52Low:  80,
		}
		result := stock.Validate()
		assert.True(t, result.IsValid)
		assert.NotEmpty(t, result.Warnings)
	})

	t.Run("Price below 52-week low warning", func(t *testing.T) {
		stock := Stock{
			Price:      50,
			Week52High: 150,
			Week52Low:  80,
		}
		result := stock.Validate()
		assert.True(t, result.IsValid)
		assert.NotEmpty(t, result.Warnings)
	})

	t.Run("Multiple errors", func(t *testing.T) {
		stock := Stock{
			Price:      0,
			Volume:     -1000,
			MarketCap:  -1000000,
			RSI14:      150,
			Week52High: 50,
			Week52Low:  100,
		}
		result := stock.Validate()
		assert.False(t, result.IsValid)
		assert.GreaterOrEqual(t, len(result.Errors), 4)
	})

	t.Run("Data quality score", func(t *testing.T) {
		// Stock with multiple warnings
		stock := Stock{
			Price:         100,
			PERatio:       5000,
			DividendYield: 50,
			PBRatio:       -2,
			GrossMargin:   150,
			NetMargin:     -500,
			Beta:          15,
		}
		result := stock.Validate()
		assert.True(t, result.IsValid)
		assert.Less(t, result.DataQuality, 100.0)
		assert.GreaterOrEqual(t, result.DataQuality, 0.0)
	})

	t.Run("Data quality never negative", func(t *testing.T) {
		stock := Stock{
			Price:         100,
			PERatio:       5000,
			DividendYield: 200,
			PBRatio:       -2,
			GrossMargin:   500,
			NetMargin:     -2000,
			Beta:          50,
			Week52High:    50,
			Week52Low:     40,
		}
		result := stock.Validate()
		assert.GreaterOrEqual(t, result.DataQuality, 0.0)
	})
}

func TestFilterDefinitions(t *testing.T) {
	filters := GetAllFilterDefinitions()

	t.Run("Has all expected categories", func(t *testing.T) {
		categories := make(map[FilterCategory]bool)
		for _, f := range filters {
			categories[f.Category] = true
		}

		assert.True(t, categories[CategoryPriceVolume])
		assert.True(t, categories[CategoryValuation])
		assert.True(t, categories[CategoryDividends])
		assert.True(t, categories[CategoryFinancialHealth])
		assert.True(t, categories[CategoryProfitability])
		assert.True(t, categories[CategoryGrowth])
		assert.True(t, categories[CategoryTechnical])
		assert.True(t, categories[CategoryProfile])
	})

	t.Run("All filters have required fields", func(t *testing.T) {
		for _, f := range filters {
			assert.NotEmpty(t, f.Field, "Filter should have field")
			assert.NotEmpty(t, f.Label, "Filter should have label")
			assert.NotEmpty(t, f.Description, "Filter should have description")
			assert.NotEmpty(t, f.Operators, "Filter should have operators")
			assert.NotEmpty(t, f.Type, "Filter should have type")
		}
	})

	t.Run("Numeric filters have valid operators", func(t *testing.T) {
		numericTypes := []FilterType{TypeNumber, TypePercent, TypeCurrency}
		for _, f := range filters {
			for _, nt := range numericTypes {
				if f.Type == nt {
					hasNumericOp := false
					for _, op := range f.Operators {
						if op == OpGreaterThan || op == OpLessThan || op == OpBetween {
							hasNumericOp = true
							break
						}
					}
					assert.True(t, hasNumericOp, "Numeric filter %s should have numeric operators", f.Field)
					break
				}
			}
		}
	})
}

func TestSectors(t *testing.T) {
	sectors := GetSectors()

	t.Run("Has minimum sectors", func(t *testing.T) {
		assert.GreaterOrEqual(t, len(sectors), 10, "Should have at least 10 sectors")
	})

	t.Run("Contains expected sectors", func(t *testing.T) {
		sectorMap := make(map[string]bool)
		for _, s := range sectors {
			sectorMap[s] = true
		}

		assert.True(t, sectorMap["Technology"])
		assert.True(t, sectorMap["Healthcare"])
		assert.True(t, sectorMap["Financial Services"])
		assert.True(t, sectorMap["Energy"])
		assert.True(t, sectorMap["Real Estate"])
	})
}

func TestMarketCapRanges(t *testing.T) {
	ranges := GetMarketCapRanges()

	t.Run("Has all ranges", func(t *testing.T) {
		assert.Contains(t, ranges, "nano")
		assert.Contains(t, ranges, "micro")
		assert.Contains(t, ranges, "small")
		assert.Contains(t, ranges, "mid")
		assert.Contains(t, ranges, "large")
		assert.Contains(t, ranges, "mega")
	})

	t.Run("Ranges are properly ordered", func(t *testing.T) {
		// Each range's max should be the next range's min
		assert.Equal(t, ranges["nano"][1], ranges["micro"][0])
		assert.Equal(t, ranges["micro"][1], ranges["small"][0])
		assert.Equal(t, ranges["small"][1], ranges["mid"][0])
		assert.Equal(t, ranges["mid"][1], ranges["large"][0])
	})

	t.Run("Nano starts at zero", func(t *testing.T) {
		assert.Equal(t, int64(0), ranges["nano"][0])
	})

	t.Run("Mega has no upper limit", func(t *testing.T) {
		assert.Equal(t, int64(0), ranges["mega"][1])
	})
}

func TestPredefinedScreeners(t *testing.T) {
	screeners := GetPredefinedScreeners()

	t.Run("Has minimum screeners", func(t *testing.T) {
		assert.GreaterOrEqual(t, len(screeners), 10)
	})

	t.Run("All screeners have required fields", func(t *testing.T) {
		for _, s := range screeners {
			assert.NotEmpty(t, s.ID, "Screener should have ID")
			assert.NotEmpty(t, s.Name, "Screener should have name")
			assert.NotEmpty(t, s.Description, "Screener should have description")
			assert.NotEmpty(t, s.Category, "Screener should have category")
			assert.NotEmpty(t, s.Filters, "Screener should have filters")
		}
	})

	t.Run("Unique IDs", func(t *testing.T) {
		ids := make(map[string]bool)
		for _, s := range screeners {
			assert.False(t, ids[s.ID], "Duplicate screener ID: %s", s.ID)
			ids[s.ID] = true
		}
	})

	t.Run("Momentum Masters has correct filters", func(t *testing.T) {
		screener := MomentumMastersScreener()
		assert.Equal(t, "momentum-masters", screener.ID)
		assert.Len(t, screener.Filters, 4)

		// Check for required filter fields
		fields := make(map[string]bool)
		for _, f := range screener.Filters {
			fields[f.Field] = true
		}
		assert.True(t, fields["return1W"])
		assert.True(t, fields["return3M"])
		assert.True(t, fields["return6M"])
		assert.True(t, fields["volume"])
	})

	t.Run("Value Opportunities has correct filters", func(t *testing.T) {
		screener := ValueOpportunitiesScreener()
		assert.Equal(t, "value-opportunities", screener.ID)

		fields := make(map[string]bool)
		for _, f := range screener.Filters {
			fields[f.Field] = true
		}
		assert.True(t, fields["peRatio"])
		assert.True(t, fields["pbRatio"])
		assert.True(t, fields["debtToEquity"])
		assert.True(t, fields["roe"])
	})
}
