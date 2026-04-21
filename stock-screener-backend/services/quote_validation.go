package services

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"stock-screener/models"
)

var ErrInvalidQuotePayload = errors.New("invalid quote payload")

// ValidateStockQuote checks a single stock payload is safe to return from the HTTP API.
func ValidateStockQuote(s *models.Stock) error {
	if s == nil {
		return fmt.Errorf("%w: stock is nil", ErrInvalidQuotePayload)
	}
	sym := strings.TrimSpace(s.Symbol)
	if sym == "" {
		return fmt.Errorf("%w: symbol is required", ErrInvalidQuotePayload)
	}
	if len(sym) > 32 {
		return fmt.Errorf("%w: symbol too long", ErrInvalidQuotePayload)
	}
	if math.IsNaN(s.Price) || math.IsInf(s.Price, 0) {
		return fmt.Errorf("%w: price is not a finite number", ErrInvalidQuotePayload)
	}
	if s.Price <= 0 {
		return fmt.Errorf("%w: price must be positive", ErrInvalidQuotePayload)
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
