package models

// FilterOperator represents comparison operators for filters
type FilterOperator string

const (
	OpEquals          FilterOperator = "eq"
	OpNotEquals       FilterOperator = "ne"
	OpGreaterThan     FilterOperator = "gt"
	OpGreaterOrEqual  FilterOperator = "gte"
	OpLessThan        FilterOperator = "lt"
	OpLessOrEqual     FilterOperator = "lte"
	OpBetween         FilterOperator = "between"
	OpIn              FilterOperator = "in"
	OpNotIn           FilterOperator = "notIn"
	OpContains        FilterOperator = "contains"
)

// FilterCategory represents a category of filters
type FilterCategory string

const (
	CategoryPriceVolume    FilterCategory = "price_volume"
	CategoryValuation      FilterCategory = "valuation"
	CategoryDividends      FilterCategory = "dividends"
	CategoryFinancialHealth FilterCategory = "financial_health"
	CategoryProfitability  FilterCategory = "profitability"
	CategoryGrowth         FilterCategory = "growth"
	CategoryTechnical      FilterCategory = "technical"
	CategoryProfile        FilterCategory = "profile"
)

// FilterType represents the data type of a filter
type FilterType string

const (
	TypeNumber   FilterType = "number"
	TypePercent  FilterType = "percent"
	TypeCurrency FilterType = "currency"
	TypeString   FilterType = "string"
	TypeBoolean  FilterType = "boolean"
)

// Filter represents a single filter condition
type Filter struct {
	Field    string         `json:"field"`
	Operator FilterOperator `json:"operator"`
	Value    interface{}    `json:"value"`
	Value2   interface{}    `json:"value2,omitempty"` // For between operator
}

// FilterDefinition describes a filterable field
type FilterDefinition struct {
	Field       string         `json:"field"`
	Label       string         `json:"label"`
	Description string         `json:"description"`
	Category    FilterCategory `json:"category"`
	Type        FilterType     `json:"type"`
	Unit        string         `json:"unit,omitempty"`
	Min         *float64       `json:"min,omitempty"`
	Max         *float64       `json:"max,omitempty"`
	Options     []string       `json:"options,omitempty"` // For categorical fields
	Operators   []FilterOperator `json:"operators"`
}

// FilterRequest represents a request with multiple filters
type FilterRequest struct {
	Filters     []Filter `json:"filters"`
	SortBy      string   `json:"sortBy,omitempty"`
	SortOrder   string   `json:"sortOrder,omitempty"` // "asc" or "desc"
	Limit       int      `json:"limit,omitempty"`
	Offset      int      `json:"offset,omitempty"`
}

// FilterResponse contains filtered results with metadata
type FilterResponse struct {
	Stocks      []Stock `json:"stocks"`
	Total       int     `json:"total"`
	Page        int     `json:"page"`
	PageSize    int     `json:"pageSize"`
	AppliedFilters []Filter `json:"appliedFilters"`
}

// GetAllFilterDefinitions returns all available filter definitions
func GetAllFilterDefinitions() []FilterDefinition {
	numericOperators := []FilterOperator{OpEquals, OpGreaterThan, OpGreaterOrEqual, OpLessThan, OpLessOrEqual, OpBetween}
	stringOperators := []FilterOperator{OpEquals, OpNotEquals, OpIn, OpNotIn, OpContains}

	minZero := 0.0
	maxHundred := 100.0

	return []FilterDefinition{
		// Price & Volume
		{Field: "price", Label: "Price", Description: "Current stock price", Category: CategoryPriceVolume, Type: TypeCurrency, Unit: "USD", Min: &minZero, Operators: numericOperators},
		{Field: "change", Label: "Price Change", Description: "Price change from previous close", Category: CategoryPriceVolume, Type: TypeCurrency, Unit: "USD", Operators: numericOperators},
		{Field: "changePercent", Label: "Price Change %", Description: "Percentage price change", Category: CategoryPriceVolume, Type: TypePercent, Operators: numericOperators},
		{Field: "volume", Label: "Volume", Description: "Current trading volume", Category: CategoryPriceVolume, Type: TypeNumber, Min: &minZero, Operators: numericOperators},
		{Field: "avgVolume", Label: "Avg Volume", Description: "Average trading volume", Category: CategoryPriceVolume, Type: TypeNumber, Min: &minZero, Operators: numericOperators},
		{Field: "marketCap", Label: "Market Cap", Description: "Market capitalization", Category: CategoryPriceVolume, Type: TypeCurrency, Unit: "USD", Min: &minZero, Operators: numericOperators},

		// Valuation
		{Field: "peRatio", Label: "P/E Ratio", Description: "Price to Earnings ratio", Category: CategoryValuation, Type: TypeNumber, Operators: numericOperators},
		{Field: "forwardPE", Label: "Forward P/E", Description: "Forward Price to Earnings ratio", Category: CategoryValuation, Type: TypeNumber, Operators: numericOperators},
		{Field: "pegRatio", Label: "PEG Ratio", Description: "Price/Earnings to Growth ratio", Category: CategoryValuation, Type: TypeNumber, Operators: numericOperators},
		{Field: "pbRatio", Label: "P/B Ratio", Description: "Price to Book ratio", Category: CategoryValuation, Type: TypeNumber, Operators: numericOperators},
		{Field: "psRatio", Label: "P/S Ratio", Description: "Price to Sales ratio", Category: CategoryValuation, Type: TypeNumber, Min: &minZero, Operators: numericOperators},
		{Field: "evToEbitda", Label: "EV/EBITDA", Description: "Enterprise Value to EBITDA", Category: CategoryValuation, Type: TypeNumber, Operators: numericOperators},
		{Field: "priceToFcf", Label: "Price/FCF", Description: "Price to Free Cash Flow", Category: CategoryValuation, Type: TypeNumber, Operators: numericOperators},

		// Dividends
		{Field: "dividendYield", Label: "Dividend Yield", Description: "Annual dividend yield", Category: CategoryDividends, Type: TypePercent, Min: &minZero, Max: &maxHundred, Operators: numericOperators},
		{Field: "payoutRatio", Label: "Payout Ratio", Description: "Dividend payout ratio", Category: CategoryDividends, Type: TypePercent, Min: &minZero, Operators: numericOperators},
		{Field: "consecutiveDivYears", Label: "Consecutive Div Years", Description: "Years of consecutive dividend payments", Category: CategoryDividends, Type: TypeNumber, Min: &minZero, Operators: numericOperators},
		{Field: "dividendGrowthYears", Label: "Div Growth Years", Description: "Years of dividend growth", Category: CategoryDividends, Type: TypeNumber, Min: &minZero, Operators: numericOperators},
		{Field: "dividendGrowthRate", Label: "Div Growth Rate", Description: "Dividend growth rate", Category: CategoryDividends, Type: TypePercent, Operators: numericOperators},

		// Financial Health
		{Field: "currentRatio", Label: "Current Ratio", Description: "Current assets / Current liabilities", Category: CategoryFinancialHealth, Type: TypeNumber, Min: &minZero, Operators: numericOperators},
		{Field: "quickRatio", Label: "Quick Ratio", Description: "Quick ratio (acid test)", Category: CategoryFinancialHealth, Type: TypeNumber, Min: &minZero, Operators: numericOperators},
		{Field: "debtToEquity", Label: "Debt/Equity", Description: "Debt to Equity ratio", Category: CategoryFinancialHealth, Type: TypeNumber, Min: &minZero, Operators: numericOperators},
		{Field: "interestCoverage", Label: "Interest Coverage", Description: "Interest coverage ratio", Category: CategoryFinancialHealth, Type: TypeNumber, Operators: numericOperators},
		{Field: "altmanZScore", Label: "Altman Z-Score", Description: "Bankruptcy risk indicator", Category: CategoryFinancialHealth, Type: TypeNumber, Operators: numericOperators},
		{Field: "cashToDebt", Label: "Cash/Debt", Description: "Cash to Total Debt ratio", Category: CategoryFinancialHealth, Type: TypeNumber, Min: &minZero, Operators: numericOperators},

		// Profitability
		{Field: "roe", Label: "ROE", Description: "Return on Equity", Category: CategoryProfitability, Type: TypePercent, Operators: numericOperators},
		{Field: "roa", Label: "ROA", Description: "Return on Assets", Category: CategoryProfitability, Type: TypePercent, Operators: numericOperators},
		{Field: "roic", Label: "ROIC", Description: "Return on Invested Capital", Category: CategoryProfitability, Type: TypePercent, Operators: numericOperators},
		{Field: "grossMargin", Label: "Gross Margin", Description: "Gross profit margin", Category: CategoryProfitability, Type: TypePercent, Operators: numericOperators},
		{Field: "operatingMargin", Label: "Operating Margin", Description: "Operating profit margin", Category: CategoryProfitability, Type: TypePercent, Operators: numericOperators},
		{Field: "netMargin", Label: "Net Margin", Description: "Net profit margin", Category: CategoryProfitability, Type: TypePercent, Operators: numericOperators},

		// Growth
		{Field: "revenueGrowth", Label: "Revenue Growth", Description: "Year-over-year revenue growth", Category: CategoryGrowth, Type: TypePercent, Operators: numericOperators},
		{Field: "epsGrowth", Label: "EPS Growth", Description: "Year-over-year EPS growth", Category: CategoryGrowth, Type: TypePercent, Operators: numericOperators},
		{Field: "bookValueGrowth", Label: "Book Value Growth", Description: "Book value growth", Category: CategoryGrowth, Type: TypePercent, Operators: numericOperators},
		{Field: "fcfGrowth", Label: "FCF Growth", Description: "Free cash flow growth", Category: CategoryGrowth, Type: TypePercent, Operators: numericOperators},

		// Technical
		{Field: "rsi14", Label: "RSI (14)", Description: "14-day Relative Strength Index", Category: CategoryTechnical, Type: TypeNumber, Min: &minZero, Max: &maxHundred, Operators: numericOperators},
		{Field: "ma50", Label: "50-Day MA", Description: "50-day moving average", Category: CategoryTechnical, Type: TypeCurrency, Min: &minZero, Operators: numericOperators},
		{Field: "ma200", Label: "200-Day MA", Description: "200-day moving average", Category: CategoryTechnical, Type: TypeCurrency, Min: &minZero, Operators: numericOperators},
		{Field: "week52High", Label: "52-Week High", Description: "52-week high price", Category: CategoryTechnical, Type: TypeCurrency, Min: &minZero, Operators: numericOperators},
		{Field: "week52Low", Label: "52-Week Low", Description: "52-week low price", Category: CategoryTechnical, Type: TypeCurrency, Min: &minZero, Operators: numericOperators},
		{Field: "beta", Label: "Beta", Description: "Stock beta vs market", Category: CategoryTechnical, Type: TypeNumber, Operators: numericOperators},

		// Performance
		{Field: "return1W", Label: "1-Week Return", Description: "One week return", Category: CategoryTechnical, Type: TypePercent, Operators: numericOperators},
		{Field: "return1M", Label: "1-Month Return", Description: "One month return", Category: CategoryTechnical, Type: TypePercent, Operators: numericOperators},
		{Field: "return3M", Label: "3-Month Return", Description: "Three month return", Category: CategoryTechnical, Type: TypePercent, Operators: numericOperators},
		{Field: "return6M", Label: "6-Month Return", Description: "Six month return", Category: CategoryTechnical, Type: TypePercent, Operators: numericOperators},
		{Field: "return1Y", Label: "1-Year Return", Description: "One year return", Category: CategoryTechnical, Type: TypePercent, Operators: numericOperators},

		// Profile
		{Field: "sector", Label: "Sector", Description: "Company sector", Category: CategoryProfile, Type: TypeString, Options: GetSectors(), Operators: stringOperators},
		{Field: "industry", Label: "Industry", Description: "Company industry", Category: CategoryProfile, Type: TypeString, Operators: stringOperators},
		{Field: "country", Label: "Country", Description: "Company country", Category: CategoryProfile, Type: TypeString, Options: []string{"USA", "China", "Japan", "UK", "Germany", "France", "Canada", "Switzerland", "Australia"}, Operators: stringOperators},
		{Field: "exchange", Label: "Exchange", Description: "Stock exchange", Category: CategoryProfile, Type: TypeString, Options: []string{"NYSE", "NASDAQ", "AMEX"}, Operators: stringOperators},

		// Special
		{Field: "piotroskiFScore", Label: "Piotroski F-Score", Description: "Piotroski financial score (0-9)", Category: CategoryFinancialHealth, Type: TypeNumber, Min: &minZero, Operators: numericOperators},
	}
}

// GetSectors returns all available sectors
func GetSectors() []string {
	return []string{
		"Technology",
		"Healthcare",
		"Financial Services",
		"Consumer Cyclical",
		"Consumer Defensive",
		"Industrials",
		"Energy",
		"Utilities",
		"Real Estate",
		"Basic Materials",
		"Communication Services",
	}
}

// GetMarketCapRanges returns predefined market cap ranges
func GetMarketCapRanges() map[string][2]int64 {
	return map[string][2]int64{
		"nano":   {0, 50_000_000},           // < $50M
		"micro":  {50_000_000, 300_000_000}, // $50M - $300M
		"small":  {300_000_000, 2_000_000_000}, // $300M - $2B
		"mid":    {2_000_000_000, 10_000_000_000}, // $2B - $10B
		"large":  {10_000_000_000, 200_000_000_000}, // $10B - $200B
		"mega":   {200_000_000_000, 0}, // > $200B (0 means no upper limit)
	}
}
