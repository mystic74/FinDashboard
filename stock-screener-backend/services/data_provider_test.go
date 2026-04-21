package services

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestDataProviderManagerGetQuotes_ReturnsErrorWhenAllLiveProvidersFail(t *testing.T) {
	t.Parallel()

	cache := NewCacheService(time.Minute, 2*time.Minute)
	manager := NewDataProviderManager(DataProviderConfig{
		YahooQuoteDriver: YahooQuoteDriverResty,
	}, cache)
	manager.yahooService.client.SetTransport(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return nil, context.DeadlineExceeded
	}))
	manager.yahooService.baseURL = "https://unit.test"

	symbols := []string{"AAPL", "MSFT"}
	stocks, err := manager.GetQuotes(context.Background(), symbols)
	if err == nil {
		t.Fatal("expected live provider failure")
	}
	if manager.demoMode {
		t.Fatal("expected live mode manager, got demo mode")
	}
	if manager.mockService != nil {
		t.Fatal("did not expect mock service fallback in live mode")
	}
	if len(stocks) != 0 {
		t.Fatalf("expected no stocks on failure, got %#v", stocks)
	}
	if !strings.Contains(err.Error(), "all live data providers failed") {
		t.Fatalf("expected aggregated live-provider failure, got %v", err)
	}
}

func TestDataProviderManagerFetchAlphaVantageQuotes_ReturnsErrorOnRateLimit(t *testing.T) {
	t.Parallel()

	cache := NewCacheService(time.Minute, 2*time.Minute)
	av := NewAlphaVantageService("test-key", cache)
	av.client.SetTransport(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Query().Get("symbol") {
		case "AAPL":
			return jsonResponse(`{
				"Global Quote": {
					"01. symbol": "AAPL",
					"02. open": "120.00",
					"03. high": "124.00",
					"04. low": "119.00",
					"05. price": "123.45",
					"06. volume": "1000",
					"08. previous close": "121.00",
					"09. change": "2.45",
					"10. change percent": "2.02%"
				}
			}`), nil
		case "MSFT":
			return jsonResponse(`{"Note":"API call frequency is 5 calls per minute"}`), nil
		default:
			t.Fatalf("unexpected symbol %q", r.URL.Query().Get("symbol"))
			return nil, nil
		}
	}))

	manager := &DataProviderManager{alphaVantage: av}

	stocks, err := manager.fetchAlphaVantageQuotes([]string{"AAPL", "MSFT"})
	if err == nil {
		t.Fatal("expected rate limit error")
	}
	if len(stocks) != 0 {
		t.Fatalf("expected no partial results, got %#v", stocks)
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) || !errors.Is(apiErr.Underlying, ErrRateLimited) {
		t.Fatalf("expected rate-limited API error, got %v", err)
	}
}
