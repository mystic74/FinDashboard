package models

import "time"

// Stock represents comprehensive stock data with fundamentals
type Stock struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Exchange      string  `json:"exchange"`
	Currency      string  `json:"currency"`
	Country       string  `json:"country"`
	Sector        string  `json:"sector"`
	Industry      string  `json:"industry"`
	Description   string  `json:"description,omitempty"`

	// Price Data
	Price            float64 `json:"price"`
	Open             float64 `json:"open"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	PreviousClose    float64 `json:"previousClose"`
	Change           float64 `json:"change"`
	ChangePercent    float64 `json:"changePercent"`

	// Volume Data
	Volume           int64   `json:"volume"`
	AvgVolume        int64   `json:"avgVolume"`
	AvgVolume10Day   int64   `json:"avgVolume10Day"`

	// Market Data
	MarketCap        int64   `json:"marketCap"`
	SharesOutstanding int64  `json:"sharesOutstanding"`
	Float            int64   `json:"float"`

	// Valuation Ratios
	PERatio          float64 `json:"peRatio"`
	ForwardPE        float64 `json:"forwardPE"`
	PEGRatio         float64 `json:"pegRatio"`
	PBRatio          float64 `json:"pbRatio"`
	PSRatio          float64 `json:"psRatio"`
	PCFRatio         float64 `json:"pcfRatio"`          // Price to Cash Flow
	EVToEBITDA       float64 `json:"evToEbitda"`
	EVToRevenue      float64 `json:"evToRevenue"`
	PriceToFCF       float64 `json:"priceToFcf"`        // Price to Free Cash Flow

	// Dividend Data
	DividendYield        float64 `json:"dividendYield"`
	DividendPerShare     float64 `json:"dividendPerShare"`
	PayoutRatio          float64 `json:"payoutRatio"`
	DividendGrowthRate   float64 `json:"dividendGrowthRate"`
	ConsecutiveDivYears  int     `json:"consecutiveDivYears"`
	DividendGrowthYears  int     `json:"dividendGrowthYears"`
	ExDividendDate       string  `json:"exDividendDate,omitempty"`

	// Risk Metrics
	Beta             float64 `json:"beta"`
	Volatility       float64 `json:"volatility"`

	// Profitability
	ROE              float64 `json:"roe"`              // Return on Equity
	ROA              float64 `json:"roa"`              // Return on Assets
	ROIC             float64 `json:"roic"`             // Return on Invested Capital
	GrossMargin      float64 `json:"grossMargin"`
	OperatingMargin  float64 `json:"operatingMargin"`
	NetMargin        float64 `json:"netMargin"`
	FCFMargin        float64 `json:"fcfMargin"`        // Free Cash Flow Margin

	// Financial Health
	DebtToEquity     float64 `json:"debtToEquity"`
	CurrentRatio     float64 `json:"currentRatio"`
	QuickRatio       float64 `json:"quickRatio"`
	InterestCoverage float64 `json:"interestCoverage"`
	AltmanZScore     float64 `json:"altmanZScore"`
	CashToDebt       float64 `json:"cashToDebt"`

	// Growth Metrics
	RevenueGrowth       float64 `json:"revenueGrowth"`         // YoY
	RevenueGrowthQoQ    float64 `json:"revenueGrowthQoQ"`
	EPSGrowth           float64 `json:"epsGrowth"`             // YoY
	EPSGrowthQoQ        float64 `json:"epsGrowthQoQ"`
	BookValueGrowth     float64 `json:"bookValueGrowth"`
	FCFGrowth           float64 `json:"fcfGrowth"`
	EarningsGrowth5Y    float64 `json:"earningsGrowth5Y"`

	// Per Share Data
	EPS              float64 `json:"eps"`
	BookValuePerShare float64 `json:"bookValuePerShare"`
	CashPerShare     float64 `json:"cashPerShare"`
	RevenuePerShare  float64 `json:"revenuePerShare"`
	FCFPerShare      float64 `json:"fcfPerShare"`

	// Income Statement (TTM)
	Revenue          int64   `json:"revenue"`
	GrossProfit      int64   `json:"grossProfit"`
	OperatingIncome  int64   `json:"operatingIncome"`
	NetIncome        int64   `json:"netIncome"`
	EBITDA           int64   `json:"ebitda"`

	// Balance Sheet
	TotalAssets      int64   `json:"totalAssets"`
	TotalLiabilities int64   `json:"totalLiabilities"`
	TotalDebt        int64   `json:"totalDebt"`
	LongTermDebt     int64   `json:"longTermDebt"`
	TotalCash        int64   `json:"totalCash"`
	TotalEquity      int64   `json:"totalEquity"`
	WorkingCapital   int64   `json:"workingCapital"`

	// Cash Flow
	OperatingCashFlow int64  `json:"operatingCashFlow"`
	CapitalExpenditure int64 `json:"capitalExpenditure"`
	FreeCashFlow      int64  `json:"freeCashFlow"`

	// Technical Indicators
	MA50             float64 `json:"ma50"`
	MA200            float64 `json:"ma200"`
	Week52High       float64 `json:"week52High"`
	Week52Low        float64 `json:"week52Low"`
	RSI14            float64 `json:"rsi14"`
	MACD             float64 `json:"macd"`
	MACDSignal       float64 `json:"macdSignal"`
	MACDHistogram    float64 `json:"macdHistogram"`

	// Performance Returns
	Return1W         float64 `json:"return1W"`
	Return1M         float64 `json:"return1M"`
	Return3M         float64 `json:"return3M"`
	Return6M         float64 `json:"return6M"`
	ReturnYTD        float64 `json:"returnYtd"`
	Return1Y         float64 `json:"return1Y"`
	Return3Y         float64 `json:"return3Y"`
	Return5Y         float64 `json:"return5Y"`

	// Piotroski F-Score Components
	PiotroskiFScore  int     `json:"piotroskiFScore"`
	FScoreDetails    *FScoreDetails `json:"fScoreDetails,omitempty"`

	// Metadata
	LastUpdated      time.Time `json:"lastUpdated"`
}

// FScoreDetails contains the breakdown of Piotroski F-Score
type FScoreDetails struct {
	// Profitability
	PositiveROA          bool `json:"positiveRoa"`
	PositiveCFO          bool `json:"positiveCfo"`
	IncreasingROA        bool `json:"increasingRoa"`
	QualityOfEarnings    bool `json:"qualityOfEarnings"`    // CFO > Net Income

	// Leverage/Liquidity
	DecreasingDebt       bool `json:"decreasingDebt"`       // Long-term debt decreasing
	IncreasingCurrentRatio bool `json:"increasingCurrentRatio"`
	NoNewShares          bool `json:"noNewShares"`

	// Operating Efficiency
	IncreasingGrossMargin bool `json:"increasingGrossMargin"`
	IncreasingAssetTurnover bool `json:"increasingAssetTurnover"`
}

// StockQuote represents a simplified quote for listing
type StockQuote struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
	Volume        int64   `json:"volume"`
	MarketCap     int64   `json:"marketCap"`
}

// HistoricalPrice represents historical price data
type HistoricalPrice struct {
	Date     time.Time `json:"date"`
	Open     float64   `json:"open"`
	High     float64   `json:"high"`
	Low      float64   `json:"low"`
	Close    float64   `json:"close"`
	AdjClose float64   `json:"adjClose"`
	Volume   int64     `json:"volume"`
}

// SectorPerformance represents sector-level performance
type SectorPerformance struct {
	Sector        string  `json:"sector"`
	Change1D      float64 `json:"change1D"`
	Change1W      float64 `json:"change1W"`
	Change1M      float64 `json:"change1M"`
	Change3M      float64 `json:"change3M"`
	ChangeYTD     float64 `json:"changeYtd"`
	Change1Y      float64 `json:"change1Y"`
	StockCount    int     `json:"stockCount"`
	MarketCap     int64   `json:"marketCap"`
	TopPerformer  string  `json:"topPerformer"`
	WorstPerformer string `json:"worstPerformer"`
}

// StockValidationResult contains validation results for stock data
type StockValidationResult struct {
	IsValid       bool     `json:"isValid"`
	Warnings      []string `json:"warnings"`
	Errors        []string `json:"errors"`
	DataQuality   float64  `json:"dataQuality"` // 0-100 score
}

// ValidateStock performs sanity checks on stock data
func (s *Stock) Validate() *StockValidationResult {
	result := &StockValidationResult{
		IsValid:     true,
		Warnings:    []string{},
		Errors:      []string{},
		DataQuality: 100.0,
	}

	// Price validations
	if s.Price <= 0 {
		result.Errors = append(result.Errors, "price must be positive")
		result.IsValid = false
	}
	if s.Price > 1000000 {
		result.Warnings = append(result.Warnings, "unusually high price")
		result.DataQuality -= 10
	}

	// Volume validations
	if s.Volume < 0 {
		result.Errors = append(result.Errors, "volume cannot be negative")
		result.IsValid = false
	}

	// Market cap validations
	if s.MarketCap < 0 {
		result.Errors = append(result.Errors, "market cap cannot be negative")
		result.IsValid = false
	}

	// Ratio validations
	if s.PERatio < -1000 || s.PERatio > 10000 {
		result.Warnings = append(result.Warnings, "P/E ratio outside normal range")
		result.DataQuality -= 5
	}
	if s.PBRatio < 0 {
		result.Warnings = append(result.Warnings, "negative P/B ratio (book value issue)")
		result.DataQuality -= 5
	}
	if s.DividendYield < 0 || s.DividendYield > 100 {
		result.Warnings = append(result.Warnings, "dividend yield outside normal range")
		result.DataQuality -= 10
	}

	// Margin validations (should be between -100 and 100)
	if s.GrossMargin < -100 || s.GrossMargin > 100 {
		result.Warnings = append(result.Warnings, "gross margin outside valid range")
		result.DataQuality -= 5
	}
	if s.NetMargin < -1000 || s.NetMargin > 100 {
		result.Warnings = append(result.Warnings, "net margin outside normal range")
		result.DataQuality -= 5
	}

	// Current ratio validation
	if s.CurrentRatio < 0 {
		result.Errors = append(result.Errors, "current ratio cannot be negative")
		result.IsValid = false
	}

	// Beta validation
	if s.Beta < -10 || s.Beta > 10 {
		result.Warnings = append(result.Warnings, "beta outside normal range")
		result.DataQuality -= 5
	}

	// 52-week high/low validation
	if s.Week52High > 0 && s.Week52Low > 0 {
		if s.Week52Low > s.Week52High {
			result.Errors = append(result.Errors, "52-week low cannot be greater than 52-week high")
			result.IsValid = false
		}
		if s.Price > s.Week52High*1.1 {
			result.Warnings = append(result.Warnings, "price significantly above 52-week high")
			result.DataQuality -= 5
		}
		if s.Price < s.Week52Low*0.9 {
			result.Warnings = append(result.Warnings, "price significantly below 52-week low")
			result.DataQuality -= 5
		}
	}

	// RSI validation
	if s.RSI14 < 0 || s.RSI14 > 100 {
		result.Errors = append(result.Errors, "RSI must be between 0 and 100")
		result.IsValid = false
	}

	// Ensure data quality doesn't go below 0
	if result.DataQuality < 0 {
		result.DataQuality = 0
	}

	return result
}
