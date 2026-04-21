package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

	"stock-screener/services"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func setupStockTestRouter(t *testing.T, responseBody string) (*gin.Engine, *StockHandler) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	cache := services.NewCacheService(5*time.Minute, 10*time.Minute)
	yahoo := services.NewYahooFinanceServiceWithDriver(cache, services.YahooQuoteDriverResty)
	client := resty.New()
	client.SetTransport(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, "/v7/finance/quote", r.URL.Path)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(responseBody)),
		}, nil
	}))
	setUnexportedField(t, yahoo, "client", client)
	setUnexportedField(t, yahoo, "baseURL", "https://unit.test")

	return router, NewStockHandler(yahoo)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func setUnexportedField(t *testing.T, target interface{}, fieldName string, value interface{}) {
	t.Helper()

	v := reflect.ValueOf(target).Elem().FieldByName(fieldName)
	if !v.IsValid() {
		t.Fatalf("field %q not found", fieldName)
	}
	reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
}

func TestGetStock_InvalidQuotePayloadReturnsProviderError(t *testing.T) {
	router, handler := setupStockTestRouter(t, `{
		"quoteResponse": {
			"result": [
				{"symbol": "AAPL", "shortName": "Apple", "regularMarketPrice": -1}
			]
		}
	}`)
	router.GET("/stocks/:symbol", handler.GetStock)

	req := httptest.NewRequest(http.MethodGet, "/stocks/AAPL", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["error"], "invalid quote payload")
	assert.Contains(t, response["error"], "price must be positive")
}

func TestGetMultipleStocks_InvalidBatchQuotePayloadReturnsProviderError(t *testing.T) {
	router, handler := setupStockTestRouter(t, `{
		"quoteResponse": {
			"result": [
				{"symbol": "AAPL", "shortName": "Apple", "regularMarketPrice": 100.5},
				{"symbol": "MSFT", "shortName": "Microsoft", "regularMarketPrice": -10}
			]
		}
	}`)
	router.GET("/stocks", handler.GetMultipleStocks)

	req := httptest.NewRequest(http.MethodGet, "/stocks?symbols=AAPL,MSFT", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.True(t, strings.Contains(response["error"].(string), "invalid quote payload"))
	assert.True(t, strings.Contains(response["error"].(string), "stocks[1] (MSFT): invalid quote payload: price must be positive"))
}
