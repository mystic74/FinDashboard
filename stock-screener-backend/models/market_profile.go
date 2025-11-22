package models

// MarketProfile defines market-specific adjustments for screener criteria
// Different markets have vastly different scales - a "large cap" in Israel
// would be a "small cap" in the USA
type MarketProfile struct {
	Country     string  `json:"country"`
	DisplayName string  `json:"displayName"`
	Currency    string  `json:"currency"`

	// Market cap thresholds (in local currency equivalent to USD scale)
	SmallCapMax int64 `json:"smallCapMax"` // Upper bound for small cap
	MidCapMin   int64 `json:"midCapMin"`   // Lower bound for mid cap
	MidCapMax   int64 `json:"midCapMax"`   // Upper bound for mid cap
	LargeCapMin int64 `json:"largeCapMin"` // Lower bound for large cap

	// Multiplier to apply to market cap filters (relative to USA = 1.0)
	MarketCapMultiplier float64 `json:"marketCapMultiplier"`

	// Typical valuation ranges for the market
	TypicalPEMin float64 `json:"typicalPEMin"`
	TypicalPEMax float64 `json:"typicalPEMax"`

	// Dividend expectations (some markets have higher/lower yields)
	TypicalDividendYield float64 `json:"typicalDividendYield"`
	DividendMultiplier   float64 `json:"dividendMultiplier"`

	// Volume expectations (smaller markets have less liquidity)
	VolumeMultiplier float64 `json:"volumeMultiplier"`

	// Growth expectations
	GrowthMultiplier float64 `json:"growthMultiplier"`
}

// GetMarketProfile returns the market profile for a given country
func GetMarketProfile(country string) *MarketProfile {
	profile, exists := MarketProfiles[country]
	if !exists {
		// Default to USA profile
		return MarketProfiles["USA"]
	}
	return profile
}

// MarketProfiles contains predefined profiles for supported markets
var MarketProfiles = map[string]*MarketProfile{
	"USA": {
		Country:              "USA",
		DisplayName:          "United States",
		Currency:             "USD",
		SmallCapMax:          2000000000,    // $2B
		MidCapMin:            2000000000,    // $2B
		MidCapMax:            10000000000,   // $10B
		LargeCapMin:          10000000000,   // $10B
		MarketCapMultiplier:  1.0,           // Baseline
		TypicalPEMin:         15,
		TypicalPEMax:         25,
		TypicalDividendYield: 1.5,
		DividendMultiplier:   1.0,
		VolumeMultiplier:     1.0,
		GrowthMultiplier:     1.0,
	},
	"Israel": {
		Country:              "Israel",
		DisplayName:          "Israel",
		Currency:             "ILS",
		SmallCapMax:          500000000,     // $500M (Israeli small cap)
		MidCapMin:            500000000,     // $500M
		MidCapMax:            2000000000,    // $2B
		LargeCapMin:          2000000000,    // $2B (large for Israel)
		MarketCapMultiplier:  0.1,           // 10x smaller thresholds
		TypicalPEMin:         10,
		TypicalPEMax:         20,
		TypicalDividendYield: 2.5,           // Israeli banks pay well
		DividendMultiplier:   1.2,
		VolumeMultiplier:     0.1,           // Much less liquidity
		GrowthMultiplier:     1.2,           // Tech-heavy market
	},
	"UK": {
		Country:              "UK",
		DisplayName:          "United Kingdom",
		Currency:             "GBP",
		SmallCapMax:          1000000000,    // GBP 1B
		MidCapMin:            1000000000,
		MidCapMax:            5000000000,
		LargeCapMin:          5000000000,
		MarketCapMultiplier:  0.5,
		TypicalPEMin:         12,
		TypicalPEMax:         20,
		TypicalDividendYield: 3.5,           // UK loves dividends
		DividendMultiplier:   1.5,
		VolumeMultiplier:     0.5,
		GrowthMultiplier:     0.9,
	},
	"Germany": {
		Country:              "Germany",
		DisplayName:          "Germany",
		Currency:             "EUR",
		SmallCapMax:          1000000000,
		MidCapMin:            1000000000,
		MidCapMax:            5000000000,
		LargeCapMin:          5000000000,
		MarketCapMultiplier:  0.5,
		TypicalPEMin:         12,
		TypicalPEMax:         22,
		TypicalDividendYield: 2.5,
		DividendMultiplier:   1.2,
		VolumeMultiplier:     0.4,
		GrowthMultiplier:     0.9,
	},
	"Japan": {
		Country:              "Japan",
		DisplayName:          "Japan",
		Currency:             "JPY",
		SmallCapMax:          100000000000,  // JPY scale
		MidCapMin:            100000000000,
		MidCapMax:            500000000000,
		LargeCapMin:          500000000000,
		MarketCapMultiplier:  0.6,
		TypicalPEMin:         12,
		TypicalPEMax:         25,
		TypicalDividendYield: 2.0,
		DividendMultiplier:   1.0,
		VolumeMultiplier:     0.7,
		GrowthMultiplier:     0.8,
	},
	"China": {
		Country:              "China",
		DisplayName:          "China",
		Currency:             "CNY",
		SmallCapMax:          2000000000,
		MidCapMin:            2000000000,
		MidCapMax:            10000000000,
		LargeCapMin:          10000000000,
		MarketCapMultiplier:  0.8,
		TypicalPEMin:         10,
		TypicalPEMax:         30,
		TypicalDividendYield: 1.5,
		DividendMultiplier:   0.8,
		VolumeMultiplier:     1.5,           // High retail participation
		GrowthMultiplier:     1.3,
	},
	"India": {
		Country:              "India",
		DisplayName:          "India",
		Currency:             "INR",
		SmallCapMax:          500000000,
		MidCapMin:            500000000,
		MidCapMax:            3000000000,
		LargeCapMin:          3000000000,
		MarketCapMultiplier:  0.2,
		TypicalPEMin:         15,
		TypicalPEMax:         35,            // India trades at premium
		TypicalDividendYield: 1.0,
		DividendMultiplier:   0.7,
		VolumeMultiplier:     0.3,
		GrowthMultiplier:     1.5,
	},
	"Brazil": {
		Country:              "Brazil",
		DisplayName:          "Brazil",
		Currency:             "BRL",
		SmallCapMax:          1000000000,
		MidCapMin:            1000000000,
		MidCapMax:            5000000000,
		LargeCapMin:          5000000000,
		MarketCapMultiplier:  0.3,
		TypicalPEMin:         8,
		TypicalPEMax:         15,            // Brazil trades cheaper
		TypicalDividendYield: 4.0,           // High yields
		DividendMultiplier:   1.5,
		VolumeMultiplier:     0.3,
		GrowthMultiplier:     1.0,
	},
	"Canada": {
		Country:              "Canada",
		DisplayName:          "Canada",
		Currency:             "CAD",
		SmallCapMax:          1000000000,
		MidCapMin:            1000000000,
		MidCapMax:            5000000000,
		LargeCapMin:          5000000000,
		MarketCapMultiplier:  0.4,
		TypicalPEMin:         12,
		TypicalPEMax:         22,
		TypicalDividendYield: 3.0,
		DividendMultiplier:   1.3,
		VolumeMultiplier:     0.4,
		GrowthMultiplier:     0.9,
	},
	"Switzerland": {
		Country:              "Switzerland",
		DisplayName:          "Switzerland",
		Currency:             "CHF",
		SmallCapMax:          1000000000,
		MidCapMin:            1000000000,
		MidCapMax:            10000000000,
		LargeCapMin:          10000000000,
		MarketCapMultiplier:  0.5,
		TypicalPEMin:         18,
		TypicalPEMax:         30,            // Quality premium
		TypicalDividendYield: 2.5,
		DividendMultiplier:   1.2,
		VolumeMultiplier:     0.3,
		GrowthMultiplier:     0.8,
	},
	"France": {
		Country:              "France",
		DisplayName:          "France",
		Currency:             "EUR",
		SmallCapMax:          1000000000,
		MidCapMin:            1000000000,
		MidCapMax:            5000000000,
		LargeCapMin:          5000000000,
		MarketCapMultiplier:  0.5,
		TypicalPEMin:         12,
		TypicalPEMax:         22,
		TypicalDividendYield: 2.5,
		DividendMultiplier:   1.2,
		VolumeMultiplier:     0.4,
		GrowthMultiplier:     0.9,
	},
	"Australia": {
		Country:              "Australia",
		DisplayName:          "Australia",
		Currency:             "AUD",
		SmallCapMax:          500000000,
		MidCapMin:            500000000,
		MidCapMax:            3000000000,
		LargeCapMin:          3000000000,
		MarketCapMultiplier:  0.3,
		TypicalPEMin:         14,
		TypicalPEMax:         22,
		TypicalDividendYield: 4.0,           // Franking credits culture
		DividendMultiplier:   1.5,
		VolumeMultiplier:     0.3,
		GrowthMultiplier:     0.9,
	},
}

// toFloat64 converts various numeric types to float64
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case int32:
		return float64(val), true
	default:
		return 0, false
	}
}

// AdjustFilterForMarket adjusts a filter's value based on the market profile
func AdjustFilterForMarket(filter Filter, profile *MarketProfile) Filter {
	adjusted := filter

	switch filter.Field {
	case "marketCap":
		// Adjust market cap thresholds
		if val, ok := toFloat64(filter.Value); ok {
			adjusted.Value = val * profile.MarketCapMultiplier
		}
		if filter.Value2 != nil {
			if val2, ok := toFloat64(filter.Value2); ok {
				adjusted.Value2 = val2 * profile.MarketCapMultiplier
			}
		}
	case "volume", "avgVolume":
		// Adjust volume expectations
		if val, ok := toFloat64(filter.Value); ok {
			adjusted.Value = val * profile.VolumeMultiplier
		}
	case "dividendYield":
		// Adjust dividend expectations
		if filter.Operator == OpGreaterThan || filter.Operator == OpGreaterOrEqual {
			if val, ok := toFloat64(filter.Value); ok {
				adjusted.Value = val * profile.DividendMultiplier
			}
		}
	case "revenueGrowth", "epsGrowth":
		// Adjust growth expectations
		if val, ok := toFloat64(filter.Value); ok {
			adjusted.Value = val * profile.GrowthMultiplier
		}
	}

	return adjusted
}

// AdjustScreenerForMarket returns a copy of the screener with filters adjusted for the market
func AdjustScreenerForMarket(screener Screener, country string) Screener {
	if country == "" || country == "USA" {
		return screener // USA is baseline, no adjustment needed
	}

	profile := GetMarketProfile(country)
	if profile == nil {
		return screener
	}

	adjusted := screener
	adjusted.Filters = make([]Filter, len(screener.Filters))

	for i, filter := range screener.Filters {
		adjusted.Filters[i] = AdjustFilterForMarket(filter, profile)
	}

	return adjusted
}
