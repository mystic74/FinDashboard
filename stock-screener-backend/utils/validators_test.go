package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSymbol(t *testing.T) {
	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{"Valid simple symbol", "AAPL", false},
		{"Valid with numbers", "BRK", false},
		{"Valid two letter", "GE", false},
		{"Valid single letter", "T", false},
		{"Valid class A", "BRK.A", false},
		{"Valid class B", "BRK.B", false},
		{"Empty symbol", "", true},
		{"Too long symbol", "TOOLONGSYMBOL", true},
		{"Invalid characters", "AAP@L", true},
		{"Numbers only", "12345", true},
		{"Lowercase", "aapl", false}, // Should be normalized
		{"With spaces", " AAPL ", false}, // Should be trimmed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSymbol(tt.symbol)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSymbols(t *testing.T) {
	t.Run("All valid", func(t *testing.T) {
		symbols := []string{"AAPL", "MSFT", "GOOGL"}
		valid, errs := ValidateSymbols(symbols)
		assert.Len(t, valid, 3)
		assert.Len(t, errs, 0)
	})

	t.Run("Some invalid", func(t *testing.T) {
		symbols := []string{"AAPL", "INVALID123456", "MSFT", ""}
		valid, errs := ValidateSymbols(symbols)
		assert.Len(t, valid, 2)
		assert.Len(t, errs, 2)
	})

	t.Run("All invalid", func(t *testing.T) {
		symbols := []string{"", "TOOLONGSYMBOL", "@#$"}
		valid, errs := ValidateSymbols(symbols)
		assert.Len(t, valid, 0)
		assert.Len(t, errs, 3)
	})

	t.Run("Normalized to uppercase", func(t *testing.T) {
		symbols := []string{"aapl", "msft"}
		valid, _ := ValidateSymbols(symbols)
		assert.Equal(t, "AAPL", valid[0])
		assert.Equal(t, "MSFT", valid[1])
	})
}

func TestValidatePositive(t *testing.T) {
	tests := []struct {
		value   float64
		wantErr bool
	}{
		{100, false},
		{0.01, false},
		{0, true},
		{-1, true},
		{-0.01, true},
	}

	for _, tt := range tests {
		err := ValidatePositive(tt.value, "test")
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidateNonNegative(t *testing.T) {
	tests := []struct {
		value   float64
		wantErr bool
	}{
		{100, false},
		{0, false},
		{-1, true},
		{-0.01, true},
	}

	for _, tt := range tests {
		err := ValidateNonNegative(tt.value, "test")
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidateRange(t *testing.T) {
	tests := []struct {
		value   float64
		min     float64
		max     float64
		wantErr bool
	}{
		{50, 0, 100, false},
		{0, 0, 100, false},
		{100, 0, 100, false},
		{-1, 0, 100, true},
		{101, 0, 100, true},
	}

	for _, tt := range tests {
		err := ValidateRange(tt.value, tt.min, tt.max, "test")
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidatePercentage(t *testing.T) {
	tests := []struct {
		value   float64
		wantErr bool
	}{
		{50, false},
		{0, false},
		{100, false},
		{-1, true},
		{101, true},
	}

	for _, tt := range tests {
		err := ValidatePercentage(tt.value, "test")
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidateRSI(t *testing.T) {
	tests := []struct {
		value   float64
		wantErr bool
	}{
		{50, false},
		{0, false},
		{100, false},
		{-1, true},
		{101, true},
	}

	for _, tt := range tests {
		err := ValidateRSI(tt.value)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidatePERatio(t *testing.T) {
	tests := []struct {
		pe       float64
		valid    bool
		hasMsg   bool
	}{
		{15, true, false},
		{50, true, false},
		{-10, true, false},
		{-1001, false, true},
		{1001, false, true},
		{500, true, false},
	}

	for _, tt := range tests {
		valid, msg := ValidatePERatio(tt.pe)
		assert.Equal(t, tt.valid, valid)
		if tt.hasMsg {
			assert.NotEmpty(t, msg)
		}
	}
}

func TestValidateMarketCap(t *testing.T) {
	assert.NoError(t, ValidateMarketCap(1000000000))
	assert.NoError(t, ValidateMarketCap(0))
	assert.Error(t, ValidateMarketCap(-1))
}

func TestValidatePaginationParams(t *testing.T) {
	tests := []struct {
		limit      int
		offset     int
		wantLimit  int
		wantOffset int
	}{
		{50, 0, 50, 0},
		{0, 0, 50, 0},      // Default limit
		{-10, -5, 50, 0},   // Negative values
		{1000, 0, 500, 0},  // Max limit
		{100, 50, 100, 50}, // Valid values
	}

	for _, tt := range tests {
		limit, offset := ValidatePaginationParams(tt.limit, tt.offset)
		assert.Equal(t, tt.wantLimit, limit)
		assert.Equal(t, tt.wantOffset, offset)
	}
}

func TestValidateSortField(t *testing.T) {
	allowed := []string{"price", "volume", "marketCap"}

	assert.True(t, ValidateSortField("price", allowed))
	assert.True(t, ValidateSortField("volume", allowed))
	assert.True(t, ValidateSortField("", allowed)) // Empty is allowed
	assert.False(t, ValidateSortField("invalid", allowed))
}

func TestValidateSortOrder(t *testing.T) {
	assert.Equal(t, "asc", ValidateSortOrder("asc"))
	assert.Equal(t, "desc", ValidateSortOrder("desc"))
	assert.Equal(t, "asc", ValidateSortOrder("ASC"))
	assert.Equal(t, "desc", ValidateSortOrder("DESC"))
	assert.Equal(t, "desc", ValidateSortOrder("invalid")) // Default
	assert.Equal(t, "desc", ValidateSortOrder(""))        // Default
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "Hello World"},
		{"<script>alert('xss')</script>", "alert('xss')"},
		{"  trimmed  ", "trimmed"},
		{"<b>Bold</b>", "Bold"},
	}

	for _, tt := range tests {
		result := SanitizeString(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}

func TestValidateFilterOperator(t *testing.T) {
	validOps := []string{"eq", "ne", "gt", "gte", "lt", "lte", "between", "in", "notIn", "contains"}
	for _, op := range validOps {
		assert.True(t, ValidateFilterOperator(op), "Should be valid: %s", op)
	}

	assert.False(t, ValidateFilterOperator("invalid"))
	assert.False(t, ValidateFilterOperator(""))
	assert.False(t, ValidateFilterOperator("EQ")) // Case sensitive
}

func TestValidateFilterField(t *testing.T) {
	validFields := []string{
		"price", "change", "changePercent", "volume", "marketCap",
		"peRatio", "pbRatio", "dividendYield", "roe", "roa",
		"currentRatio", "debtToEquity", "sector", "industry",
	}

	for _, field := range validFields {
		assert.True(t, ValidateFilterField(field), "Should be valid: %s", field)
	}

	assert.False(t, ValidateFilterField("invalid"))
	assert.False(t, ValidateFilterField(""))
}

func TestValidateStockData(t *testing.T) {
	t.Run("Valid data", func(t *testing.T) {
		check := ValidateStockData(100, 1000000, 50000000000, 15, 2, 2.5, 120, 80, 50)
		assert.True(t, check.Passed)
		assert.Empty(t, check.Errors)
		assert.Equal(t, 100.0, check.Score)
	})

	t.Run("Invalid price", func(t *testing.T) {
		check := ValidateStockData(0, 1000000, 50000000000, 15, 2, 2.5, 120, 80, 50)
		assert.False(t, check.Passed)
		assert.Contains(t, check.Errors, "price must be positive")
	})

	t.Run("Negative volume", func(t *testing.T) {
		check := ValidateStockData(100, -1000, 50000000000, 15, 2, 2.5, 120, 80, 50)
		assert.False(t, check.Passed)
		assert.Contains(t, check.Errors, "volume cannot be negative")
	})

	t.Run("Negative market cap", func(t *testing.T) {
		check := ValidateStockData(100, 1000000, -50000000000, 15, 2, 2.5, 120, 80, 50)
		assert.False(t, check.Passed)
		assert.Contains(t, check.Errors, "market cap cannot be negative")
	})

	t.Run("Invalid RSI", func(t *testing.T) {
		check := ValidateStockData(100, 1000000, 50000000000, 15, 2, 2.5, 120, 80, 150)
		assert.False(t, check.Passed)
		assert.Contains(t, check.Errors, "RSI must be between 0 and 100")
	})

	t.Run("52-week inconsistency", func(t *testing.T) {
		check := ValidateStockData(100, 1000000, 50000000000, 15, 2, 2.5, 80, 120, 50)
		assert.False(t, check.Passed)
		assert.Contains(t, check.Errors, "52-week low cannot exceed 52-week high")
	})

	t.Run("High dividend yield warning", func(t *testing.T) {
		check := ValidateStockData(100, 1000000, 50000000000, 15, 2, 25, 120, 80, 50)
		assert.True(t, check.Passed)
		assert.NotEmpty(t, check.Warnings)
		assert.Less(t, check.Score, 100.0)
	})

	t.Run("Extreme P/E warning", func(t *testing.T) {
		check := ValidateStockData(100, 1000000, 50000000000, 2000, 2, 2.5, 120, 80, 50)
		assert.True(t, check.Passed)
		assert.NotEmpty(t, check.Warnings)
	})

	t.Run("Multiple issues", func(t *testing.T) {
		check := ValidateStockData(0, -100, -100, 15, 2, 2.5, 80, 120, 150)
		assert.False(t, check.Passed)
		assert.Len(t, check.Errors, 5)
		assert.Equal(t, 0.0, check.Score)
	})
}

func TestDataQualityScore(t *testing.T) {
	t.Run("Perfect score", func(t *testing.T) {
		check := ValidateStockData(100, 1000000, 50000000000, 15, 2, 2.5, 120, 80, 50)
		assert.Equal(t, 100.0, check.Score)
	})

	t.Run("Reduced score with warnings", func(t *testing.T) {
		check := ValidateStockData(100, 1000000, 50000000000, 2000, -1, 25, 120, 80, 50)
		assert.True(t, check.Passed)
		assert.Less(t, check.Score, 100.0)
	})

	t.Run("Score never negative", func(t *testing.T) {
		check := ValidateStockData(0, -1, -1, 2000, -1, 25, 80, 120, 150)
		assert.GreaterOrEqual(t, check.Score, 0.0)
	})
}
