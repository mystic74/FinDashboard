package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestYahooFinanceService_GetQuotes_Resty_UsesBaseURL(t *testing.T) {
	t.Parallel()
	respBody := YahooQuoteResponse{}
	respBody.QuoteResponse.Result = []YahooQuote{{
		Symbol:             "AAPL",
		ShortName:          "Apple",
		RegularMarketPrice: 100.5,
	}}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v7/finance/quote" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(respBody)
	}))
	t.Cleanup(srv.Close)

	cache := NewCacheService(time.Minute, 2*time.Minute)
	y := NewYahooFinanceServiceWithDriver(cache, YahooQuoteDriverResty)
	y.baseURL = srv.URL

	stocks, err := y.GetQuotes(context.Background(), []string{"AAPL"})
	if err != nil {
		t.Fatalf("GetQuotes: %v", err)
	}
	if len(stocks) != 1 {
		t.Fatalf("len=%d", len(stocks))
	}
	if stocks[0].Symbol != "AAPL" || stocks[0].Price != 100.5 {
		t.Fatalf("got %+v", stocks[0])
	}
	if err := ValidateStockQuote(&stocks[0]); err != nil {
		t.Fatalf("validate: %v", err)
	}
}
