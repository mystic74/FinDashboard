package utils

import (
	"math"
)

// FinancialData holds data needed for financial calculations
type FinancialData struct {
	// Current Year Data
	NetIncome           float64
	TotalAssets         float64
	OperatingCashFlow   float64
	Revenue             float64
	GrossProfit         float64
	CurrentAssets       float64
	CurrentLiabilities  float64
	LongTermDebt        float64
	SharesOutstanding   float64

	// Previous Year Data
	PrevNetIncome       float64
	PrevTotalAssets     float64
	PrevGrossProfit     float64
	PrevRevenue         float64
	PrevCurrentRatio    float64
	PrevLongTermDebt    float64
	PrevSharesOutstanding float64

	// Balance Sheet
	TotalLiabilities    float64
	TotalEquity         float64
	WorkingCapital      float64
	RetainedEarnings    float64
	EBIT                float64
	MarketCap           float64
	TotalDebt           float64
	TotalCash           float64
}

// PiotroskiFScoreResult contains the F-Score and its components
type PiotroskiFScoreResult struct {
	Score                   int
	PositiveROA             bool
	PositiveCFO             bool
	IncreasingROA           bool
	QualityOfEarnings       bool
	DecreasingLongTermDebt  bool
	IncreasingCurrentRatio  bool
	NoNewShares             bool
	IncreasingGrossMargin   bool
	IncreasingAssetTurnover bool
}

// CalculatePiotroskiFScore calculates the Piotroski F-Score (0-9)
func CalculatePiotroskiFScore(data FinancialData) PiotroskiFScoreResult {
	result := PiotroskiFScoreResult{}

	// Profitability Signals (4 points)

	// 1. ROA > 0 (Net Income / Total Assets)
	roa := 0.0
	if data.TotalAssets > 0 {
		roa = data.NetIncome / data.TotalAssets
	}
	result.PositiveROA = roa > 0
	if result.PositiveROA {
		result.Score++
	}

	// 2. Operating Cash Flow > 0
	result.PositiveCFO = data.OperatingCashFlow > 0
	if result.PositiveCFO {
		result.Score++
	}

	// 3. ROA is increasing (current ROA > previous ROA)
	prevROA := 0.0
	if data.PrevTotalAssets > 0 {
		prevROA = data.PrevNetIncome / data.PrevTotalAssets
	}
	result.IncreasingROA = roa > prevROA
	if result.IncreasingROA {
		result.Score++
	}

	// 4. Quality of Earnings: CFO > Net Income (accrual)
	result.QualityOfEarnings = data.OperatingCashFlow > data.NetIncome
	if result.QualityOfEarnings {
		result.Score++
	}

	// Leverage, Liquidity and Source of Funds Signals (3 points)

	// 5. Long-term debt is decreasing
	result.DecreasingLongTermDebt = data.LongTermDebt < data.PrevLongTermDebt
	if result.DecreasingLongTermDebt {
		result.Score++
	}

	// 6. Current ratio is increasing
	currentRatio := 0.0
	if data.CurrentLiabilities > 0 {
		currentRatio = data.CurrentAssets / data.CurrentLiabilities
	}
	result.IncreasingCurrentRatio = currentRatio > data.PrevCurrentRatio
	if result.IncreasingCurrentRatio {
		result.Score++
	}

	// 7. No new shares issued (shares outstanding not increased)
	result.NoNewShares = data.SharesOutstanding <= data.PrevSharesOutstanding
	if result.NoNewShares {
		result.Score++
	}

	// Operating Efficiency Signals (2 points)

	// 8. Gross margin is increasing
	grossMargin := 0.0
	if data.Revenue > 0 {
		grossMargin = data.GrossProfit / data.Revenue
	}
	prevGrossMargin := 0.0
	if data.PrevRevenue > 0 {
		prevGrossMargin = data.PrevGrossProfit / data.PrevRevenue
	}
	result.IncreasingGrossMargin = grossMargin > prevGrossMargin
	if result.IncreasingGrossMargin {
		result.Score++
	}

	// 9. Asset turnover is increasing (Revenue / Total Assets)
	assetTurnover := 0.0
	if data.TotalAssets > 0 {
		assetTurnover = data.Revenue / data.TotalAssets
	}
	prevAssetTurnover := 0.0
	if data.PrevTotalAssets > 0 {
		prevAssetTurnover = data.PrevRevenue / data.PrevTotalAssets
	}
	result.IncreasingAssetTurnover = assetTurnover > prevAssetTurnover
	if result.IncreasingAssetTurnover {
		result.Score++
	}

	return result
}

// CalculateAltmanZScore calculates the Altman Z-Score for bankruptcy prediction
// Z-Score > 2.99: Safe zone
// Z-Score 1.81-2.99: Grey zone
// Z-Score < 1.81: Distress zone
func CalculateAltmanZScore(data FinancialData) float64 {
	if data.TotalAssets == 0 {
		return 0
	}

	// X1 = Working Capital / Total Assets
	x1 := data.WorkingCapital / data.TotalAssets

	// X2 = Retained Earnings / Total Assets
	x2 := data.RetainedEarnings / data.TotalAssets

	// X3 = EBIT / Total Assets
	x3 := data.EBIT / data.TotalAssets

	// X4 = Market Value of Equity / Total Liabilities
	x4 := 0.0
	if data.TotalLiabilities > 0 {
		x4 = data.MarketCap / data.TotalLiabilities
	}

	// X5 = Revenue / Total Assets
	x5 := data.Revenue / data.TotalAssets

	// Z = 1.2*X1 + 1.4*X2 + 3.3*X3 + 0.6*X4 + 1.0*X5
	zScore := 1.2*x1 + 1.4*x2 + 3.3*x3 + 0.6*x4 + 1.0*x5

	return math.Round(zScore*100) / 100
}

// CalculateRSI calculates the Relative Strength Index
func CalculateRSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		return 50 // Default to neutral if not enough data
	}

	var gains, losses float64

	// Calculate initial average gain and loss
	for i := 1; i <= period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	// Calculate subsequent values using smoothing
	for i := period + 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			avgGain = (avgGain*float64(period-1) + change) / float64(period)
			avgLoss = (avgLoss * float64(period-1)) / float64(period)
		} else {
			avgGain = (avgGain * float64(period-1)) / float64(period)
			avgLoss = (avgLoss*float64(period-1) - change) / float64(period)
		}
	}

	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return math.Round(rsi*100) / 100
}

// CalculateMACD calculates MACD, Signal, and Histogram
func CalculateMACD(prices []float64) (macd, signal, histogram float64) {
	if len(prices) < 26 {
		return 0, 0, 0
	}

	ema12 := CalculateEMA(prices, 12)
	ema26 := CalculateEMA(prices, 26)

	macd = ema12 - ema26

	// For signal line, we need MACD values history
	// Simplified: use last MACD value with EMA smoothing factor
	signal = macd * 0.15 // Approximation for 9-period EMA of MACD

	histogram = macd - signal

	return math.Round(macd*100) / 100, math.Round(signal*100) / 100, math.Round(histogram*100) / 100
}

// CalculateEMA calculates Exponential Moving Average
func CalculateEMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}

	multiplier := 2.0 / float64(period+1)

	// Start with SMA
	var sum float64
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema := sum / float64(period)

	// Calculate EMA
	for i := period; i < len(prices); i++ {
		ema = (prices[i]-ema)*multiplier + ema
	}

	return ema
}

// CalculateSMA calculates Simple Moving Average
func CalculateSMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}

	var sum float64
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i]
	}

	return sum / float64(period)
}

// CalculateROE calculates Return on Equity
func CalculateROE(netIncome, totalEquity float64) float64 {
	if totalEquity <= 0 {
		return 0
	}
	return (netIncome / totalEquity) * 100
}

// CalculateROA calculates Return on Assets
func CalculateROA(netIncome, totalAssets float64) float64 {
	if totalAssets <= 0 {
		return 0
	}
	return (netIncome / totalAssets) * 100
}

// CalculateROIC calculates Return on Invested Capital
func CalculateROIC(nopat, investedCapital float64) float64 {
	if investedCapital <= 0 {
		return 0
	}
	return (nopat / investedCapital) * 100
}

// CalculateDebtToEquity calculates Debt to Equity ratio
func CalculateDebtToEquity(totalDebt, totalEquity float64) float64 {
	if totalEquity <= 0 {
		return 0
	}
	return totalDebt / totalEquity
}

// CalculateCurrentRatio calculates Current Ratio
func CalculateCurrentRatio(currentAssets, currentLiabilities float64) float64 {
	if currentLiabilities <= 0 {
		return 0
	}
	return currentAssets / currentLiabilities
}

// CalculateQuickRatio calculates Quick Ratio (Acid Test)
func CalculateQuickRatio(currentAssets, inventory, currentLiabilities float64) float64 {
	if currentLiabilities <= 0 {
		return 0
	}
	return (currentAssets - inventory) / currentLiabilities
}

// CalculatePERatio calculates Price to Earnings ratio
func CalculatePERatio(price, eps float64) float64 {
	if eps <= 0 {
		return 0
	}
	return price / eps
}

// CalculatePEGRatio calculates PEG ratio
func CalculatePEGRatio(peRatio, epsGrowth float64) float64 {
	if epsGrowth <= 0 {
		return 0
	}
	return peRatio / epsGrowth
}

// CalculatePBRatio calculates Price to Book ratio
func CalculatePBRatio(price, bookValuePerShare float64) float64 {
	if bookValuePerShare <= 0 {
		return 0
	}
	return price / bookValuePerShare
}

// CalculateDividendYield calculates dividend yield as percentage
func CalculateDividendYield(annualDividend, price float64) float64 {
	if price <= 0 {
		return 0
	}
	return (annualDividend / price) * 100
}

// CalculatePayoutRatio calculates dividend payout ratio
func CalculatePayoutRatio(dividendPerShare, eps float64) float64 {
	if eps <= 0 {
		return 0
	}
	return (dividendPerShare / eps) * 100
}

// CalculateFreeCashFlow calculates Free Cash Flow
func CalculateFreeCashFlow(operatingCashFlow, capitalExpenditure float64) float64 {
	return operatingCashFlow - capitalExpenditure
}

// CalculateGrossMargin calculates Gross Margin percentage
func CalculateGrossMargin(grossProfit, revenue float64) float64 {
	if revenue <= 0 {
		return 0
	}
	return (grossProfit / revenue) * 100
}

// CalculateOperatingMargin calculates Operating Margin percentage
func CalculateOperatingMargin(operatingIncome, revenue float64) float64 {
	if revenue <= 0 {
		return 0
	}
	return (operatingIncome / revenue) * 100
}

// CalculateNetMargin calculates Net Margin percentage
func CalculateNetMargin(netIncome, revenue float64) float64 {
	if revenue <= 0 {
		return 0
	}
	return (netIncome / revenue) * 100
}

// CalculatePercentageReturn calculates percentage return between two prices
func CalculatePercentageReturn(startPrice, endPrice float64) float64 {
	if startPrice <= 0 {
		return 0
	}
	return ((endPrice - startPrice) / startPrice) * 100
}

// CalculateGrowthRate calculates YoY growth rate
func CalculateGrowthRate(current, previous float64) float64 {
	if previous == 0 {
		return 0
	}
	return ((current - previous) / math.Abs(previous)) * 100
}

// CalculateInterestCoverage calculates Interest Coverage Ratio
func CalculateInterestCoverage(ebit, interestExpense float64) float64 {
	if interestExpense <= 0 {
		return 100 // No interest expense means excellent coverage
	}
	return ebit / interestExpense
}

// CalculateMarketCap calculates Market Capitalization
func CalculateMarketCap(price float64, sharesOutstanding int64) int64 {
	return int64(price * float64(sharesOutstanding))
}

// CalculateEnterpriseValue calculates Enterprise Value
func CalculateEnterpriseValue(marketCap, totalDebt, cash int64) int64 {
	return marketCap + totalDebt - cash
}

// RoundToDecimalPlaces rounds a number to specified decimal places
func RoundToDecimalPlaces(value float64, places int) float64 {
	multiplier := math.Pow(10, float64(places))
	return math.Round(value*multiplier) / multiplier
}
