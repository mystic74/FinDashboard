package services

import (
	"math"
	"testing"

	"stock-screener/models"
)

func TestValidateStockQuote(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		stock   *models.Stock
		wantErr bool
	}{
		{name: "nil", stock: nil, wantErr: true},
		{name: "empty symbol", stock: &models.Stock{Symbol: "  ", Price: 1}, wantErr: true},
		{name: "symbol too long", stock: &models.Stock{Symbol: string(make([]byte, 33)), Price: 1}, wantErr: true},
		{name: "nan price", stock: &models.Stock{Symbol: "AAPL", Price: math.NaN()}, wantErr: true},
		{name: "inf price", stock: &models.Stock{Symbol: "AAPL", Price: math.Inf(1)}, wantErr: true},
		{name: "negative price", stock: &models.Stock{Symbol: "AAPL", Price: -0.01}, wantErr: true},
		{name: "ok zero price", stock: &models.Stock{Symbol: "AAPL", Price: 0}, wantErr: false},
		{name: "ok positive", stock: &models.Stock{Symbol: "AAPL", Price: 123.45}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateStockQuote(tt.stock)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateStockQuotes(t *testing.T) {
	t.Parallel()
	stocks := []models.Stock{
		{Symbol: "A", Price: 1},
		{Symbol: "B", Price: math.NaN()},
	}
	err := ValidateStockQuotes(stocks)
	if err == nil {
		t.Fatal("expected error on bad second row")
	}
}
