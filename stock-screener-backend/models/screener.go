package models

// Screener represents a predefined or custom screener
type Screener struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Filters     []Filter `json:"filters"`
	SortBy      string   `json:"sortBy,omitempty"`
	SortOrder   string   `json:"sortOrder,omitempty"`
	Icon        string   `json:"icon,omitempty"`
	IsCustom    bool     `json:"isCustom"`
	CreatedAt   string   `json:"createdAt,omitempty"`
}

// ScreenerResult contains the results of running a screener
type ScreenerResult struct {
	Screener     Screener `json:"screener"`
	Stocks       []Stock  `json:"stocks"`
	Total        int      `json:"total"`
	ExecutionMs  int64    `json:"executionMs"`
	LastUpdated  string   `json:"lastUpdated"`
}

// ScreenerSummary provides a quick overview of a screener
type ScreenerSummary struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Category     string `json:"category"`
	Icon         string `json:"icon"`
	MatchCount   int    `json:"matchCount"`
	TopStock     string `json:"topStock,omitempty"`
}

// GetPredefinedScreeners returns all predefined screeners
func GetPredefinedScreeners() []Screener {
	return []Screener{
		MomentumMastersScreener(),
		DividendAristocratsScreener(),
		ValueOpportunitiesScreener(),
		HighBetaBullsScreener(),
		CashIsKingScreener(),
		PiotroskiHighScoreScreener(),
		SmallCapGrowthScreener(),
		UndervaluedTechScreener(),
		GrowthAtReasonablePriceScreener(),
		QualityStocksScreener(),
		LowVolatilityScreener(),
		TurnaroundCandidatesScreener(),
	}
}

// MomentumMastersScreener returns the Momentum Masters screener
func MomentumMastersScreener() Screener {
	return Screener{
		ID:          "momentum-masters",
		Name:        "Momentum Masters",
		Description: "Stocks with strong price momentum and high volume",
		Category:    "Momentum",
		Icon:        "trending-up",
		Filters: []Filter{
			{Field: "return1W", Operator: OpGreaterThan, Value: 0},
			{Field: "return3M", Operator: OpGreaterThan, Value: 25},
			{Field: "return6M", Operator: OpGreaterThan, Value: 25},
			{Field: "volume", Operator: OpGreaterThan, Value: 100000},
		},
		SortBy:    "return3M",
		SortOrder: "desc",
	}
}

// DividendAristocratsScreener returns the Dividend Aristocrats screener
func DividendAristocratsScreener() Screener {
	return Screener{
		ID:          "dividend-aristocrats",
		Name:        "Dividend Aristocrats",
		Description: "Companies with long track records of dividend payments and growth",
		Category:    "Income",
		Icon:        "dollar-sign",
		Filters: []Filter{
			{Field: "dividendYield", Operator: OpGreaterThan, Value: 2},
			{Field: "payoutRatio", Operator: OpLessThan, Value: 80},
			{Field: "consecutiveDivYears", Operator: OpGreaterThan, Value: 10},
			{Field: "dividendGrowthYears", Operator: OpGreaterThan, Value: 5},
		},
		SortBy:    "dividendYield",
		SortOrder: "desc",
	}
}

// ValueOpportunitiesScreener returns the Value Opportunities screener
func ValueOpportunitiesScreener() Screener {
	return Screener{
		ID:          "value-opportunities",
		Name:        "Value Opportunities",
		Description: "Undervalued stocks with strong fundamentals",
		Category:    "Value",
		Icon:        "search",
		Filters: []Filter{
			{Field: "peRatio", Operator: OpLessThan, Value: 20},
			{Field: "peRatio", Operator: OpGreaterThan, Value: 0}, // Ensure positive earnings
			{Field: "pbRatio", Operator: OpLessThan, Value: 5},
			{Field: "debtToEquity", Operator: OpLessThan, Value: 1.5},
			{Field: "freeCashFlow", Operator: OpGreaterThan, Value: 0},
		},
		SortBy:    "peRatio",
		SortOrder: "asc",
	}
}

// HighBetaBullsScreener returns the High Beta Bulls screener
func HighBetaBullsScreener() Screener {
	return Screener{
		ID:          "high-beta-bulls",
		Name:        "High Beta Bulls",
		Description: "High-volatility stocks with strong recent performance",
		Category:    "Momentum",
		Icon:        "zap",
		Filters: []Filter{
			{Field: "beta", Operator: OpGreaterThan, Value: 1.5},
			{Field: "return1M", Operator: OpGreaterThan, Value: 10},
			{Field: "marketCap", Operator: OpGreaterThan, Value: 100000000}, // > $100M
		},
		SortBy:    "return1M",
		SortOrder: "desc",
	}
}

// CashIsKingScreener returns the Cash is King screener
func CashIsKingScreener() Screener {
	return Screener{
		ID:          "cash-is-king",
		Name:        "Cash is King",
		Description: "Companies with exceptional liquidity and strong cash positions",
		Category:    "Financial Health",
		Icon:        "shield",
		Filters: []Filter{
			{Field: "currentRatio", Operator: OpGreaterThan, Value: 2},
			{Field: "quickRatio", Operator: OpGreaterThan, Value: 1.5},
			{Field: "cashToDebt", Operator: OpGreaterThan, Value: 1},
			{Field: "operatingCashFlow", Operator: OpGreaterThan, Value: 0},
		},
		SortBy:    "currentRatio",
		SortOrder: "desc",
	}
}

// PiotroskiHighScoreScreener returns stocks with high Piotroski F-Score
func PiotroskiHighScoreScreener() Screener {
	return Screener{
		ID:          "piotroski-high-score",
		Name:        "Piotroski F-Score Leaders",
		Description: "Stocks with Piotroski F-Score of 8 or 9, indicating strong fundamentals",
		Category:    "Quality",
		Icon:        "award",
		Filters: []Filter{
			{Field: "piotroskiFScore", Operator: OpGreaterOrEqual, Value: 8},
		},
		SortBy:    "piotroskiFScore",
		SortOrder: "desc",
	}
}

// SmallCapGrowthScreener returns small cap growth stocks
func SmallCapGrowthScreener() Screener {
	return Screener{
		ID:          "small-cap-growth",
		Name:        "Small Cap Growth",
		Description: "High-growth small cap stocks with reasonable valuations",
		Category:    "Growth",
		Icon:        "sprout",
		Filters: []Filter{
			{Field: "marketCap", Operator: OpBetween, Value: 300000000, Value2: 2000000000},
			{Field: "revenueGrowth", Operator: OpGreaterThan, Value: 20},
			{Field: "epsGrowth", Operator: OpGreaterThan, Value: 25},
			{Field: "peRatio", Operator: OpLessThan, Value: 30},
			{Field: "peRatio", Operator: OpGreaterThan, Value: 0},
		},
		SortBy:    "epsGrowth",
		SortOrder: "desc",
	}
}

// UndervaluedTechScreener returns undervalued tech stocks
func UndervaluedTechScreener() Screener {
	return Screener{
		ID:          "undervalued-tech",
		Name:        "Undervalued Tech",
		Description: "Technology stocks trading below industry average valuations",
		Category:    "Value",
		Icon:        "cpu",
		Filters: []Filter{
			{Field: "sector", Operator: OpEquals, Value: "Technology"},
			{Field: "peRatio", Operator: OpLessThan, Value: 25},     // Below tech average
			{Field: "peRatio", Operator: OpGreaterThan, Value: 0},
			{Field: "pegRatio", Operator: OpLessThan, Value: 1},
			{Field: "revenueGrowth", Operator: OpGreaterThan, Value: 15},
		},
		SortBy:    "pegRatio",
		SortOrder: "asc",
	}
}

// GrowthAtReasonablePriceScreener (GARP) strategy screener
func GrowthAtReasonablePriceScreener() Screener {
	return Screener{
		ID:          "garp",
		Name:        "Growth at Reasonable Price",
		Description: "GARP strategy: high growth stocks with reasonable P/E and PEG ratios",
		Category:    "Growth",
		Icon:        "trending-up",
		Filters: []Filter{
			{Field: "epsGrowth", Operator: OpGreaterThan, Value: 15},
			{Field: "revenueGrowth", Operator: OpGreaterThan, Value: 10},
			{Field: "pegRatio", Operator: OpBetween, Value: 0.5, Value2: 1.5},
			{Field: "peRatio", Operator: OpLessThan, Value: 25},
			{Field: "peRatio", Operator: OpGreaterThan, Value: 0},
			{Field: "roe", Operator: OpGreaterThan, Value: 12},
		},
		SortBy:    "epsGrowth",
		SortOrder: "desc",
	}
}

// QualityStocksScreener returns high-quality stocks
func QualityStocksScreener() Screener {
	return Screener{
		ID:          "quality-stocks",
		Name:        "Quality Stocks",
		Description: "High-quality companies with strong profitability and financial health",
		Category:    "Quality",
		Icon:        "star",
		Filters: []Filter{
			{Field: "roe", Operator: OpGreaterThan, Value: 15},
			{Field: "roa", Operator: OpGreaterThan, Value: 8},
			{Field: "grossMargin", Operator: OpGreaterThan, Value: 40},
			{Field: "operatingMargin", Operator: OpGreaterThan, Value: 15},
			{Field: "debtToEquity", Operator: OpLessThan, Value: 1},
			{Field: "currentRatio", Operator: OpGreaterThan, Value: 1.5},
		},
		SortBy:    "roe",
		SortOrder: "desc",
	}
}

// LowVolatilityScreener returns low-volatility stocks
func LowVolatilityScreener() Screener {
	return Screener{
		ID:          "low-volatility",
		Name:        "Low Volatility",
		Description: "Stable stocks with low beta and consistent performance",
		Category:    "Defensive",
		Icon:        "shield",
		Filters: []Filter{
			{Field: "beta", Operator: OpLessThan, Value: 0.8},
			{Field: "beta", Operator: OpGreaterThan, Value: 0},
			{Field: "dividendYield", Operator: OpGreaterThan, Value: 1},
			{Field: "marketCap", Operator: OpGreaterThan, Value: 1000000000}, // > $1B
		},
		SortBy:    "beta",
		SortOrder: "asc",
	}
}

// TurnaroundCandidatesScreener identifies potential turnaround stocks
func TurnaroundCandidatesScreener() Screener {
	return Screener{
		ID:          "turnaround-candidates",
		Name:        "Turnaround Candidates",
		Description: "Beaten-down stocks with improving fundamentals",
		Category:    "Value",
		Icon:        "refresh-cw",
		Filters: []Filter{
			{Field: "return1Y", Operator: OpLessThan, Value: -20},     // Down significantly
			{Field: "return1M", Operator: OpGreaterThan, Value: 5},    // Recent uptick
			{Field: "currentRatio", Operator: OpGreaterThan, Value: 1}, // Not in distress
			{Field: "revenueGrowth", Operator: OpGreaterThan, Value: 0}, // Revenue stabilizing
		},
		SortBy:    "return1M",
		SortOrder: "desc",
	}
}
