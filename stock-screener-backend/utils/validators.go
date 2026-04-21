package utils

import (
	"errors"
	"regexp"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateSymbol validates a stock symbol
func ValidateSymbol(symbol string) error {
	if symbol == "" {
		return errors.New("symbol cannot be empty")
	}

	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	// Stock symbols are typically 1-5 characters, letters and sometimes numbers
	// Some have dots (BRK.A, BRK.B) or hyphens
	validPattern := regexp.MustCompile(`^[A-Z]{1,5}(\.[A-Z])?(-[A-Z]{1,2})?$`)
	if !validPattern.MatchString(symbol) {
		return errors.New("invalid symbol format")
	}

	return nil
}

// ValidateSymbols validates a list of symbols
func ValidateSymbols(symbols []string) ([]string, []ValidationError) {
	var valid []string
	var errs []ValidationError

	for _, s := range symbols {
		if err := ValidateSymbol(s); err != nil {
			errs = append(errs, ValidationError{
				Field:   s,
				Message: err.Error(),
			})
		} else {
			valid = append(valid, strings.ToUpper(strings.TrimSpace(s)))
		}
	}

	return valid, errs
}

// ValidatePositive checks if a value is positive
func ValidatePositive(value float64, fieldName string) error {
	if value <= 0 {
		return errors.New(fieldName + " must be positive")
	}
	return nil
}

// ValidateNonNegative checks if a value is non-negative
func ValidateNonNegative(value float64, fieldName string) error {
	if value < 0 {
		return errors.New(fieldName + " cannot be negative")
	}
	return nil
}

// ValidateRange checks if a value is within a range
func ValidateRange(value, min, max float64, fieldName string) error {
	if value < min || value > max {
		return errors.New(fieldName + " must be between specified range")
	}
	return nil
}

// ValidatePercentage validates a percentage value (0-100)
func ValidatePercentage(value float64, fieldName string) error {
	return ValidateRange(value, 0, 100, fieldName)
}

// ValidateRSI validates RSI value
func ValidateRSI(value float64) error {
	if value < 0 || value > 100 {
		return errors.New("RSI must be between 0 and 100")
	}
	return nil
}

// ValidatePERatio validates P/E ratio for sanity
func ValidatePERatio(pe float64) (bool, string) {
	if pe < -1000 {
		return false, "P/E ratio is extremely negative, data may be unreliable"
	}
	if pe > 1000 {
		return false, "P/E ratio is extremely high, data may be unreliable"
	}
	return true, ""
}

// ValidateMarketCap validates market cap
func ValidateMarketCap(marketCap int64) error {
	if marketCap < 0 {
		return errors.New("market cap cannot be negative")
	}
	return nil
}

// ValidatePaginationParams validates limit and offset parameters
func ValidatePaginationParams(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

// ValidateSortField validates sort field against allowed fields
func ValidateSortField(field string, allowedFields []string) bool {
	if field == "" {
		return true
	}
	for _, f := range allowedFields {
		if f == field {
			return true
		}
	}
	return false
}

// ValidateSortOrder validates sort order
func ValidateSortOrder(order string) string {
	order = strings.ToLower(strings.TrimSpace(order))
	if order != "asc" && order != "desc" {
		return "desc"
	}
	return order
}

// SanitizeString removes potentially harmful characters
func SanitizeString(s string) string {
	// Remove any HTML tags
	htmlPattern := regexp.MustCompile(`<[^>]*>`)
	s = htmlPattern.ReplaceAllString(s, "")

	// Remove special characters except common ones
	s = strings.TrimSpace(s)

	return s
}

// ValidateFilterOperator validates filter operator
func ValidateFilterOperator(op string) bool {
	validOperators := map[string]bool{
		"eq":       true,
		"ne":       true,
		"gt":       true,
		"gte":      true,
		"lt":       true,
		"lte":      true,
		"between":  true,
		"in":       true,
		"notIn":    true,
		"contains": true,
	}
	return validOperators[op]
}

// ValidateFilterField validates filter field
func ValidateFilterField(field string) bool {
	validFields := map[string]bool{
		// Price & Volume
		"price": true, "change": true, "changePercent": true,
		"volume": true, "avgVolume": true, "marketCap": true,
		// Valuation
		"peRatio": true, "forwardPE": true, "pegRatio": true,
		"pbRatio": true, "psRatio": true, "evToEbitda": true,
		"priceToFcf": true,
		// Dividends
		"dividendYield": true, "payoutRatio": true,
		"consecutiveDivYears": true, "dividendGrowthYears": true,
		"dividendGrowthRate": true,
		// Financial Health
		"currentRatio": true, "quickRatio": true, "debtToEquity": true,
		"interestCoverage": true, "altmanZScore": true, "cashToDebt": true,
		// Profitability
		"roe": true, "roa": true, "roic": true,
		"grossMargin": true, "operatingMargin": true, "netMargin": true,
		// Growth
		"revenueGrowth": true, "epsGrowth": true,
		"bookValueGrowth": true, "fcfGrowth": true,
		// Technical
		"rsi14": true, "ma50": true, "ma200": true,
		"week52High": true, "week52Low": true, "beta": true,
		// Returns
		"return1W": true, "return1M": true, "return3M": true,
		"return6M": true, "return1Y": true,
		// Profile
		"sector": true, "industry": true, "country": true, "exchange": true,
		// Special
		"piotroskiFScore": true, "operatingCashFlow": true, "freeCashFlow": true,
	}
	return validFields[field]
}

// DataQualityCheck performs comprehensive data quality validation
type DataQualityCheck struct {
	Passed   bool
	Score    float64 // 0-100
	Warnings []string
	Errors   []string
}

// ValidateStockData performs comprehensive validation on stock data
func ValidateStockData(
	price, volume float64,
	marketCap int64,
	peRatio, pbRatio, dividendYield float64,
	week52High, week52Low float64,
	rsi float64,
) DataQualityCheck {
	check := DataQualityCheck{
		Passed:   true,
		Score:    100,
		Warnings: []string{},
		Errors:   []string{},
	}

	// Critical validations (errors)
	if price <= 0 {
		check.Errors = append(check.Errors, "price must be positive")
		check.Score -= 20
	}
	if volume < 0 {
		check.Errors = append(check.Errors, "volume cannot be negative")
		check.Score -= 15
	}
	if marketCap < 0 {
		check.Errors = append(check.Errors, "market cap cannot be negative")
		check.Score -= 15
	}

	// RSI must be 0-100
	if rsi < 0 || rsi > 100 {
		check.Errors = append(check.Errors, "RSI must be between 0 and 100")
		check.Score -= 10
	}

	// 52-week consistency
	if week52High > 0 && week52Low > 0 && week52Low > week52High {
		check.Errors = append(check.Errors, "52-week low cannot exceed 52-week high")
		check.Score -= 15
	}

	// Warning validations
	if price > 10000 {
		check.Warnings = append(check.Warnings, "unusually high stock price")
		check.Score -= 5
	}

	if valid, msg := ValidatePERatio(peRatio); !valid {
		check.Warnings = append(check.Warnings, msg)
		check.Score -= 5
	}

	if dividendYield > 20 {
		check.Warnings = append(check.Warnings, "dividend yield over 20% may indicate data issue or distressed stock")
		check.Score -= 5
	}

	if pbRatio < 0 {
		check.Warnings = append(check.Warnings, "negative P/B ratio indicates negative book value")
		check.Score -= 5
	}

	// Ensure score doesn't go below 0
	if check.Score < 0 {
		check.Score = 0
	}

	// Mark as failed if there are critical errors
	if len(check.Errors) > 0 {
		check.Passed = false
	}

	return check
}
