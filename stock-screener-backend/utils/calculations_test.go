package utils

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePiotroskiFScore(t *testing.T) {
	tests := []struct {
		name     string
		data     FinancialData
		expected int
	}{
		{
			name: "Perfect F-Score (9)",
			data: FinancialData{
				NetIncome:             1000000,
				TotalAssets:           5000000,
				OperatingCashFlow:     1500000, // CFO > Net Income
				Revenue:               8000000,
				GrossProfit:           3200000, // 40% margin
				CurrentAssets:         2000000,
				CurrentLiabilities:    800000, // Current ratio = 2.5
				LongTermDebt:          500000,
				SharesOutstanding:     1000000,
				PrevNetIncome:         800000,
				PrevTotalAssets:       4800000,
				PrevGrossProfit:       2880000, // 36% margin (increasing)
				PrevRevenue:           8000000,
				PrevCurrentRatio:      2.0,     // Increasing
				PrevLongTermDebt:      600000,  // Decreasing
				PrevSharesOutstanding: 1000000, // No new shares
			},
			expected: 9,
		},
		{
			name: "Zero F-Score",
			data: FinancialData{
				NetIncome:             -500000, // Negative ROA
				TotalAssets:           5000000,
				OperatingCashFlow:     -200000, // Negative CFO
				Revenue:               6000000,
				GrossProfit:           1800000,
				CurrentAssets:         1000000,
				CurrentLiabilities:    2000000, // Current ratio = 0.5 (decreasing)
				LongTermDebt:          1500000,
				SharesOutstanding:     1500000, // New shares issued
				PrevNetIncome:         -300000,
				PrevTotalAssets:       5000000,
				PrevGrossProfit:       2100000, // Decreasing margin
				PrevRevenue:           7000000,
				PrevCurrentRatio:      0.6,     // Still decreasing
				PrevLongTermDebt:      1000000, // Increasing debt
				PrevSharesOutstanding: 1000000,
			},
			expected: 0,
		},
		{
			name: "Middle F-Score (5)",
			data: FinancialData{
				NetIncome:             500000,
				TotalAssets:           5000000,
				OperatingCashFlow:     600000, // CFO > Net Income (2 points: positive ROA, positive CFO, quality)
				Revenue:               8000000,
				GrossProfit:           2400000,
				CurrentAssets:         1500000,
				CurrentLiabilities:    1000000,
				LongTermDebt:          800000,
				SharesOutstanding:     1000000,
				PrevNetIncome:         600000, // ROA decreasing
				PrevTotalAssets:       5000000,
				PrevGrossProfit:       2400000,
				PrevRevenue:           8000000,
				PrevCurrentRatio:      1.4, // Increasing
				PrevLongTermDebt:      900000, // Decreasing
				PrevSharesOutstanding: 1000000,
			},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePiotroskiFScore(tt.data)
			assert.Equal(t, tt.expected, result.Score, "F-Score should match expected")
		})
	}
}

func TestCalculatePiotroskiFScoreComponents(t *testing.T) {
	data := FinancialData{
		NetIncome:             1000000,
		TotalAssets:           5000000,
		OperatingCashFlow:     1200000,
		Revenue:               8000000,
		GrossProfit:           3200000,
		CurrentAssets:         2000000,
		CurrentLiabilities:    800000,
		LongTermDebt:          500000,
		SharesOutstanding:     1000000,
		PrevNetIncome:         800000,
		PrevTotalAssets:       5000000,
		PrevGrossProfit:       2400000,
		PrevRevenue:           8000000,
		PrevCurrentRatio:      2.0,
		PrevLongTermDebt:      600000,
		PrevSharesOutstanding: 1000000,
	}

	result := CalculatePiotroskiFScore(data)

	t.Run("PositiveROA", func(t *testing.T) {
		assert.True(t, result.PositiveROA, "ROA should be positive")
	})

	t.Run("PositiveCFO", func(t *testing.T) {
		assert.True(t, result.PositiveCFO, "CFO should be positive")
	})

	t.Run("IncreasingROA", func(t *testing.T) {
		// Current ROA: 1000000/5000000 = 0.2
		// Previous ROA: 800000/5000000 = 0.16
		assert.True(t, result.IncreasingROA, "ROA should be increasing")
	})

	t.Run("QualityOfEarnings", func(t *testing.T) {
		// CFO (1200000) > Net Income (1000000)
		assert.True(t, result.QualityOfEarnings, "Quality of earnings should be true")
	})

	t.Run("DecreasingDebt", func(t *testing.T) {
		// 500000 < 600000
		assert.True(t, result.DecreasingLongTermDebt, "Long-term debt should be decreasing")
	})

	t.Run("IncreasingCurrentRatio", func(t *testing.T) {
		// Current: 2000000/800000 = 2.5 > 2.0
		assert.True(t, result.IncreasingCurrentRatio, "Current ratio should be increasing")
	})

	t.Run("NoNewShares", func(t *testing.T) {
		assert.True(t, result.NoNewShares, "No new shares should be issued")
	})

	t.Run("IncreasingGrossMargin", func(t *testing.T) {
		// Current: 3200000/8000000 = 40% > 2400000/8000000 = 30%
		assert.True(t, result.IncreasingGrossMargin, "Gross margin should be increasing")
	})
}

func TestCalculateAltmanZScore(t *testing.T) {
	tests := []struct {
		name     string
		data     FinancialData
		expected float64
		zone     string // "safe", "grey", or "distress"
	}{
		{
			name: "Healthy Company (Safe Zone)",
			data: FinancialData{
				WorkingCapital:   2000000,
				TotalAssets:      10000000,
				RetainedEarnings: 3000000,
				EBIT:             1500000,
				MarketCap:        15000000,
				TotalLiabilities: 4000000,
				Revenue:          12000000,
			},
			expected: 3.59,
			zone:     "safe",
		},
		{
			name: "Grey Zone Company",
			data: FinancialData{
				WorkingCapital:   500000,
				TotalAssets:      10000000,
				RetainedEarnings: 1000000,
				EBIT:             500000,
				MarketCap:        8000000,
				TotalLiabilities: 6000000,
				Revenue:          10000000,
			},
			expected: 2.19,
			zone:     "grey",
		},
		{
			name: "Distressed Company",
			data: FinancialData{
				WorkingCapital:   -1000000,
				TotalAssets:      10000000,
				RetainedEarnings: -500000,
				EBIT:             -200000,
				MarketCap:        2000000,
				TotalLiabilities: 9000000,
				Revenue:          8000000,
			},
			expected: 0.65,
			zone:     "distress",
		},
		{
			name: "Zero Assets",
			data: FinancialData{
				TotalAssets: 0,
			},
			expected: 0,
			zone:     "distress",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateAltmanZScore(tt.data)
			assert.InDelta(t, tt.expected, result, 0.1, "Z-Score should be close to expected")

			// Verify zone classification
			switch tt.zone {
			case "safe":
				assert.GreaterOrEqual(t, result, 2.99, "Should be in safe zone")
			case "grey":
				assert.GreaterOrEqual(t, result, 1.81, "Should be at least grey zone")
				assert.Less(t, result, 2.99, "Should be below safe zone")
			case "distress":
				assert.Less(t, result, 1.81, "Should be in distress zone")
			}
		})
	}
}

func TestCalculateRSI(t *testing.T) {
	tests := []struct {
		name     string
		prices   []float64
		period   int
		expected float64
		valid    bool
	}{
		{
			name: "Overbought condition",
			prices: []float64{
				100, 102, 104, 106, 108, 110, 112, 114, 116, 118,
				120, 122, 124, 126, 128, 130,
			},
			period:   14,
			expected: 100,
			valid:    true,
		},
		{
			name: "Oversold condition",
			prices: []float64{
				100, 98, 96, 94, 92, 90, 88, 86, 84, 82,
				80, 78, 76, 74, 72, 70,
			},
			period:   14,
			expected: 0,
			valid:    true,
		},
		{
			name: "Neutral condition",
			prices: []float64{
				100, 102, 100, 102, 100, 102, 100, 102, 100, 102,
				100, 102, 100, 102, 100, 102,
			},
			period:   14,
			expected: 50,
			valid:    true,
		},
		{
			name:     "Insufficient data",
			prices:   []float64{100, 102, 104},
			period:   14,
			expected: 50, // Default neutral
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateRSI(tt.prices, tt.period)

			// RSI should always be between 0 and 100
			assert.GreaterOrEqual(t, result, 0.0, "RSI should be >= 0")
			assert.LessOrEqual(t, result, 100.0, "RSI should be <= 100")

			if tt.valid {
				assert.InDelta(t, tt.expected, result, 5.0, "RSI should be close to expected")
			}
		})
	}
}

func TestCalculateMovingAverages(t *testing.T) {
	prices := []float64{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	t.Run("SMA", func(t *testing.T) {
		sma5 := CalculateSMA(prices, 5)
		// Last 5 prices: 16, 17, 18, 19, 20 = 90/5 = 18
		assert.InDelta(t, 18.0, sma5, 0.01, "5-period SMA should be 18")

		sma3 := CalculateSMA(prices, 3)
		// Last 3 prices: 18, 19, 20 = 57/3 = 19
		assert.InDelta(t, 19.0, sma3, 0.01, "3-period SMA should be 19")
	})

	t.Run("EMA", func(t *testing.T) {
		ema5 := CalculateEMA(prices, 5)
		// EMA gives more weight to recent prices
		assert.Greater(t, ema5, 15.0, "5-period EMA should be > 15")
		assert.Less(t, ema5, 20.0, "5-period EMA should be < 20")
	})

	t.Run("Insufficient data for SMA", func(t *testing.T) {
		sma := CalculateSMA([]float64{1, 2}, 5)
		assert.Equal(t, 0.0, sma, "SMA with insufficient data should be 0")
	})
}

func TestCalculateRatios(t *testing.T) {
	t.Run("ROE", func(t *testing.T) {
		roe := CalculateROE(100000, 500000)
		assert.InDelta(t, 20.0, roe, 0.01, "ROE should be 20%")

		roe = CalculateROE(100000, 0)
		assert.Equal(t, 0.0, roe, "ROE with zero equity should be 0")

		roe = CalculateROE(100000, -100000)
		assert.Equal(t, 0.0, roe, "ROE with negative equity should be 0")
	})

	t.Run("ROA", func(t *testing.T) {
		roa := CalculateROA(100000, 1000000)
		assert.InDelta(t, 10.0, roa, 0.01, "ROA should be 10%")

		roa = CalculateROA(100000, 0)
		assert.Equal(t, 0.0, roa, "ROA with zero assets should be 0")
	})

	t.Run("DebtToEquity", func(t *testing.T) {
		de := CalculateDebtToEquity(500000, 1000000)
		assert.InDelta(t, 0.5, de, 0.01, "D/E should be 0.5")

		de = CalculateDebtToEquity(500000, 0)
		assert.Equal(t, 0.0, de, "D/E with zero equity should be 0")
	})

	t.Run("CurrentRatio", func(t *testing.T) {
		cr := CalculateCurrentRatio(2000000, 1000000)
		assert.InDelta(t, 2.0, cr, 0.01, "Current ratio should be 2.0")

		cr = CalculateCurrentRatio(2000000, 0)
		assert.Equal(t, 0.0, cr, "Current ratio with zero liabilities should be 0")
	})

	t.Run("QuickRatio", func(t *testing.T) {
		qr := CalculateQuickRatio(2000000, 500000, 1000000)
		assert.InDelta(t, 1.5, qr, 0.01, "Quick ratio should be 1.5")
	})

	t.Run("PERatio", func(t *testing.T) {
		pe := CalculatePERatio(100, 5)
		assert.InDelta(t, 20.0, pe, 0.01, "P/E should be 20")

		pe = CalculatePERatio(100, 0)
		assert.Equal(t, 0.0, pe, "P/E with zero EPS should be 0")

		pe = CalculatePERatio(100, -5)
		assert.Equal(t, 0.0, pe, "P/E with negative EPS should be 0")
	})

	t.Run("PEGRatio", func(t *testing.T) {
		peg := CalculatePEGRatio(20, 10)
		assert.InDelta(t, 2.0, peg, 0.01, "PEG should be 2.0")

		peg = CalculatePEGRatio(20, 0)
		assert.Equal(t, 0.0, peg, "PEG with zero growth should be 0")
	})

	t.Run("DividendYield", func(t *testing.T) {
		dy := CalculateDividendYield(4, 100)
		assert.InDelta(t, 4.0, dy, 0.01, "Dividend yield should be 4%")

		dy = CalculateDividendYield(4, 0)
		assert.Equal(t, 0.0, dy, "Dividend yield with zero price should be 0")
	})

	t.Run("PayoutRatio", func(t *testing.T) {
		pr := CalculatePayoutRatio(2, 4)
		assert.InDelta(t, 50.0, pr, 0.01, "Payout ratio should be 50%")

		pr = CalculatePayoutRatio(2, 0)
		assert.Equal(t, 0.0, pr, "Payout ratio with zero EPS should be 0")
	})

	t.Run("FreeCashFlow", func(t *testing.T) {
		fcf := CalculateFreeCashFlow(1000000, 300000)
		assert.Equal(t, 700000.0, fcf, "FCF should be 700000")
	})

	t.Run("Margins", func(t *testing.T) {
		gm := CalculateGrossMargin(400000, 1000000)
		assert.InDelta(t, 40.0, gm, 0.01, "Gross margin should be 40%")

		om := CalculateOperatingMargin(200000, 1000000)
		assert.InDelta(t, 20.0, om, 0.01, "Operating margin should be 20%")

		nm := CalculateNetMargin(100000, 1000000)
		assert.InDelta(t, 10.0, nm, 0.01, "Net margin should be 10%")
	})

	t.Run("GrowthRate", func(t *testing.T) {
		growth := CalculateGrowthRate(120, 100)
		assert.InDelta(t, 20.0, growth, 0.01, "Growth should be 20%")

		growth = CalculateGrowthRate(80, 100)
		assert.InDelta(t, -20.0, growth, 0.01, "Growth should be -20%")

		growth = CalculateGrowthRate(100, 0)
		assert.Equal(t, 0.0, growth, "Growth with zero previous should be 0")
	})

	t.Run("PercentageReturn", func(t *testing.T) {
		ret := CalculatePercentageReturn(100, 150)
		assert.InDelta(t, 50.0, ret, 0.01, "Return should be 50%")

		ret = CalculatePercentageReturn(100, 80)
		assert.InDelta(t, -20.0, ret, 0.01, "Return should be -20%")
	})

	t.Run("InterestCoverage", func(t *testing.T) {
		ic := CalculateInterestCoverage(500000, 100000)
		assert.InDelta(t, 5.0, ic, 0.01, "Interest coverage should be 5x")

		ic = CalculateInterestCoverage(500000, 0)
		assert.Equal(t, 100.0, ic, "Interest coverage with zero interest should be 100")
	})

	t.Run("MarketCap", func(t *testing.T) {
		mc := CalculateMarketCap(150.0, 1000000)
		assert.Equal(t, int64(150000000), mc, "Market cap should be 150M")
	})

	t.Run("EnterpriseValue", func(t *testing.T) {
		ev := CalculateEnterpriseValue(150000000, 50000000, 20000000)
		assert.Equal(t, int64(180000000), ev, "EV should be 180M")
	})
}

func TestRoundToDecimalPlaces(t *testing.T) {
	tests := []struct {
		value    float64
		places   int
		expected float64
	}{
		{3.14159, 2, 3.14},
		{3.145, 2, 3.15},
		{3.14, 0, 3},
		{100.999, 1, 101.0},
	}

	for _, tt := range tests {
		result := RoundToDecimalPlaces(tt.value, tt.places)
		assert.Equal(t, tt.expected, result)
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("Zero division protection", func(t *testing.T) {
		// All these should not panic and return sensible defaults
		assert.NotPanics(t, func() {
			CalculateROE(100, 0)
			CalculateROA(100, 0)
			CalculateDebtToEquity(100, 0)
			CalculateCurrentRatio(100, 0)
			CalculateQuickRatio(100, 0, 0)
			CalculatePERatio(100, 0)
			CalculatePEGRatio(100, 0)
			CalculateDividendYield(100, 0)
			CalculatePayoutRatio(100, 0)
			CalculateGrossMargin(100, 0)
			CalculateOperatingMargin(100, 0)
			CalculateNetMargin(100, 0)
			CalculateGrowthRate(100, 0)
			CalculatePercentageReturn(0, 100)
			CalculateAltmanZScore(FinancialData{})
		})
	})

	t.Run("NaN and Inf protection", func(t *testing.T) {
		result := CalculateROE(math.Inf(1), 100)
		assert.False(t, math.IsNaN(result) || math.IsInf(result, 0), "Should not return NaN or Inf")
	})
}
