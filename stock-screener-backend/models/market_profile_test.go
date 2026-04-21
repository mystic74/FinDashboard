package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMarketProfile(t *testing.T) {
	t.Run("Returns USA profile for USA", func(t *testing.T) {
		profile := GetMarketProfile("USA")
		assert.NotNil(t, profile)
		assert.Equal(t, "USA", profile.Country)
		assert.Equal(t, 1.0, profile.MarketCapMultiplier)
	})

	t.Run("Returns Israel profile", func(t *testing.T) {
		profile := GetMarketProfile("Israel")
		assert.NotNil(t, profile)
		assert.Equal(t, "Israel", profile.Country)
		assert.Equal(t, 0.1, profile.MarketCapMultiplier) // Much smaller market
	})

	t.Run("Returns USA profile for unknown country", func(t *testing.T) {
		profile := GetMarketProfile("Narnia")
		assert.NotNil(t, profile)
		assert.Equal(t, "USA", profile.Country) // Fallback to USA
	})

	t.Run("All defined markets have profiles", func(t *testing.T) {
		countries := []string{"USA", "Israel", "UK", "Germany", "Japan", "China", "India", "Brazil", "Canada", "Switzerland", "France", "Australia"}
		for _, country := range countries {
			profile := GetMarketProfile(country)
			assert.NotNil(t, profile, "Should have profile for %s", country)
			assert.Equal(t, country, profile.Country)
			assert.Greater(t, profile.MarketCapMultiplier, 0.0, "Multiplier should be positive for %s", country)
		}
	})
}

func TestAdjustFilterForMarket(t *testing.T) {
	israelProfile := GetMarketProfile("Israel")

	t.Run("Adjusts market cap filter", func(t *testing.T) {
		filter := Filter{
			Field:    "marketCap",
			Operator: OpGreaterThan,
			Value:    float64(1000000000), // $1B in USA terms
		}
		adjusted := AdjustFilterForMarket(filter, israelProfile)

		// Israel multiplier is 0.1, so $1B becomes $100M
		assert.Equal(t, float64(100000000), adjusted.Value)
	})

	t.Run("Adjusts market cap between filter", func(t *testing.T) {
		filter := Filter{
			Field:    "marketCap",
			Operator: OpBetween,
			Value:    float64(300000000),  // $300M
			Value2:   float64(2000000000), // $2B
		}
		adjusted := AdjustFilterForMarket(filter, israelProfile)

		// Both values should be scaled
		assert.Equal(t, float64(30000000), adjusted.Value)   // $30M
		assert.Equal(t, float64(200000000), adjusted.Value2) // $200M
	})

	t.Run("Adjusts volume filter", func(t *testing.T) {
		filter := Filter{
			Field:    "volume",
			Operator: OpGreaterThan,
			Value:    float64(1000000), // 1M shares in USA
		}
		adjusted := AdjustFilterForMarket(filter, israelProfile)

		// Israel volume multiplier is 0.1
		assert.Equal(t, float64(100000), adjusted.Value) // 100K shares
	})

	t.Run("Adjusts dividend yield filter", func(t *testing.T) {
		filter := Filter{
			Field:    "dividendYield",
			Operator: OpGreaterThan,
			Value:    float64(2.0), // 2% in USA
		}
		adjusted := AdjustFilterForMarket(filter, israelProfile)

		// Israel dividend multiplier is 1.2
		assert.Equal(t, float64(2.4), adjusted.Value) // 2.4%
	})

	t.Run("Does not adjust PE ratio filter", func(t *testing.T) {
		filter := Filter{
			Field:    "peRatio",
			Operator: OpLessThan,
			Value:    float64(20),
		}
		adjusted := AdjustFilterForMarket(filter, israelProfile)

		// PE ratio should not be adjusted
		assert.Equal(t, float64(20), adjusted.Value)
	})
}

func TestAdjustScreenerForMarket(t *testing.T) {
	t.Run("Does not adjust for USA", func(t *testing.T) {
		screener := SmallCapGrowthScreener()
		adjusted := AdjustScreenerForMarket(screener, "USA")

		// Should be identical
		assert.Equal(t, len(screener.Filters), len(adjusted.Filters))
		for i, filter := range screener.Filters {
			assert.Equal(t, filter.Value, adjusted.Filters[i].Value)
		}
	})

	t.Run("Does not adjust for empty country", func(t *testing.T) {
		screener := SmallCapGrowthScreener()
		adjusted := AdjustScreenerForMarket(screener, "")

		assert.Equal(t, len(screener.Filters), len(adjusted.Filters))
	})

	t.Run("Adjusts Small Cap Growth for Israel", func(t *testing.T) {
		screener := SmallCapGrowthScreener()
		adjusted := AdjustScreenerForMarket(screener, "Israel")

		// Find the market cap filter
		var originalMarketCapMin, originalMarketCapMax float64
		var adjustedMarketCapMin, adjustedMarketCapMax float64

		for _, f := range screener.Filters {
			if f.Field == "marketCap" && f.Operator == OpBetween {
				originalMarketCapMin, _ = toFloat64(f.Value)
				originalMarketCapMax, _ = toFloat64(f.Value2)
			}
		}

		for _, f := range adjusted.Filters {
			if f.Field == "marketCap" && f.Operator == OpBetween {
				adjustedMarketCapMin, _ = toFloat64(f.Value)
				adjustedMarketCapMax, _ = toFloat64(f.Value2)
			}
		}

		// Israel has 0.1 multiplier
		assert.Equal(t, originalMarketCapMin*0.1, adjustedMarketCapMin)
		assert.Equal(t, originalMarketCapMax*0.1, adjustedMarketCapMax)
	})

	t.Run("Adjusts High Beta Bulls for UK", func(t *testing.T) {
		screener := HighBetaBullsScreener()
		adjusted := AdjustScreenerForMarket(screener, "UK")

		// Find market cap filter
		var originalMarketCap, adjustedMarketCap float64

		for _, f := range screener.Filters {
			if f.Field == "marketCap" && f.Operator == OpGreaterThan {
				originalMarketCap, _ = toFloat64(f.Value)
			}
		}

		for _, f := range adjusted.Filters {
			if f.Field == "marketCap" && f.Operator == OpGreaterThan {
				adjustedMarketCap, _ = toFloat64(f.Value)
			}
		}

		// UK has 0.5 multiplier
		assert.Equal(t, originalMarketCap*0.5, adjustedMarketCap)
	})
}

func TestMarketProfilesHaveValidData(t *testing.T) {
	for country, profile := range MarketProfiles {
		t.Run(country+" has valid data", func(t *testing.T) {
			assert.NotEmpty(t, profile.Country)
			assert.NotEmpty(t, profile.DisplayName)
			assert.NotEmpty(t, profile.Currency)

			// Market cap thresholds should be ordered
			assert.LessOrEqual(t, profile.SmallCapMax, profile.MidCapMax)

			// Multipliers should be positive
			assert.Greater(t, profile.MarketCapMultiplier, 0.0)
			assert.Greater(t, profile.VolumeMultiplier, 0.0)
			assert.Greater(t, profile.DividendMultiplier, 0.0)
			assert.Greater(t, profile.GrowthMultiplier, 0.0)

			// PE ranges should make sense
			assert.Less(t, profile.TypicalPEMin, profile.TypicalPEMax)
		})
	}
}
