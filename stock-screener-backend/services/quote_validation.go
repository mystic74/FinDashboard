package services

import (
	"fmt"
	"math"
	"strings"

	"stock-screener/models"
)

// ValidateStockQuote checks a single stock payload is safe to return from the HTTP API.
func ValidateStockQuote(s *models.Stock) error {
	if s == nil {
		return fmt.Errorf("stock is nil")
	}
	sym := strings.TrimSpace(s.Symbol)
	if sym == "" {
		return fmt.Errorf("symbol is required")
	}
	if len(sym) > 32 {
		return fmt.Errorf("symbol too long")
	}
	if math.IsNaN(s.Price) || math.IsInf(s.Price, 0) {
		return fmt.Errorf("price is not a finite number")
	}
	if s.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	return nil
}

// ValidateStockQuotes validates each stock; used for batch endpoints.
func ValidateStockQuotes(stocks []models.Stock) error {
	for i := range stocks {
		if err := ValidateStockQuote(&stocks[i]); err != nil {
			return fmt.Errorf("stocks[%d] (%s): %w", i, stocks[i].Symbol, err)
		}
	}
	return nil
}
