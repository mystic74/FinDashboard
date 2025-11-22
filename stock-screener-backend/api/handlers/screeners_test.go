package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"stock-screener/services"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *ScreenerHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	cache := services.NewCacheService(5*time.Minute, 10*time.Minute)
	engine := services.NewScreenerEngineWithDemo(cache)
	handler := NewScreenerHandler(engine)

	return router, handler
}

func TestRunScreenerWithCountryFilter(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/screeners/:name", handler.RunScreener)

	t.Run("Without country filter returns all stocks", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/screeners/piotroski-high-score", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		result := response["result"].(map[string]interface{})
		totalWithoutFilter := int(result["total"].(float64))
		t.Logf("Without filter: %d stocks", totalWithoutFilter)

		// Should have multiple countries in results
		stocks := result["stocks"].([]interface{})
		countries := make(map[string]bool)
		for _, s := range stocks {
			stock := s.(map[string]interface{})
			if country, ok := stock["country"].(string); ok {
				countries[country] = true
			}
		}
		assert.Greater(t, len(countries), 1, "Should have stocks from multiple countries")
	})

	t.Run("With country filter returns only that country", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/screeners/piotroski-high-score?country=USA", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		result := response["result"].(map[string]interface{})
		stocks := result["stocks"].([]interface{})

		// All stocks should be from USA
		for _, s := range stocks {
			stock := s.(map[string]interface{})
			assert.Equal(t, "USA", stock["country"], "All stocks should be from USA")
		}
		t.Logf("With USA filter: %d stocks", len(stocks))
	})

	t.Run("With Israel filter returns Israeli stocks", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/screeners/piotroski-high-score?country=Israel", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		result := response["result"].(map[string]interface{})
		stocks := result["stocks"].([]interface{})

		// All stocks should be from Israel
		for _, s := range stocks {
			stock := s.(map[string]interface{})
			assert.Equal(t, "Israel", stock["country"], "All stocks should be from Israel")
		}
		t.Logf("With Israel filter: %d stocks", len(stocks))
	})

	t.Run("With sector filter returns only that sector", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/screeners/piotroski-high-score?sector=Technology", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		result := response["result"].(map[string]interface{})
		stocks := result["stocks"].([]interface{})

		// All stocks should be Technology
		for _, s := range stocks {
			stock := s.(map[string]interface{})
			assert.Equal(t, "Technology", stock["sector"], "All stocks should be Technology")
		}
		t.Logf("With Technology filter: %d stocks", len(stocks))
	})

	t.Run("With country AND sector filter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/screeners/piotroski-high-score?country=USA&sector=Technology", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		result := response["result"].(map[string]interface{})
		stocks := result["stocks"].([]interface{})

		// All stocks should be USA + Technology
		for _, s := range stocks {
			stock := s.(map[string]interface{})
			assert.Equal(t, "USA", stock["country"], "All stocks should be from USA")
			assert.Equal(t, "Technology", stock["sector"], "All stocks should be Technology")
		}
		t.Logf("With USA+Technology filter: %d stocks", len(stocks))
	})
}

func TestRunScreenerNotFound(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/screeners/:name", handler.RunScreener)

	req := httptest.NewRequest(http.MethodGet, "/screeners/invalid-screener", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Screener not found", response["error"])
}

func TestGetAllScreeners(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/screeners", handler.GetAllScreeners)

	req := httptest.NewRequest(http.MethodGet, "/screeners", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, float64(12), response["count"], "Should have 12 screeners")
}

func TestScreenerMarketAdjustments(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/screeners/:name", handler.RunScreener)

	t.Run("USA screeners are not adjusted", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/screeners/small-cap-growth?country=USA", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		// USA should not be adjusted
		assert.False(t, response["adjusted"].(bool), "USA should not be adjusted")
	})

	t.Run("Israel screeners are adjusted with correct multiplier", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/screeners/small-cap-growth?country=Israel", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		// Israel should be adjusted
		assert.True(t, response["adjusted"].(bool), "Israel should be adjusted")

		// Check market profile
		profile := response["marketProfile"].(map[string]interface{})
		assert.Equal(t, "Israel", profile["country"])
		assert.Equal(t, 0.1, profile["marketCapMultiplier"], "Israel should have 0.1x market cap multiplier")
	})

	t.Run("UK screeners are adjusted with correct multiplier", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/screeners/high-beta-bulls?country=UK", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// UK should be adjusted
		assert.True(t, response["adjusted"].(bool), "UK should be adjusted")

		// Check market profile
		profile := response["marketProfile"].(map[string]interface{})
		assert.Equal(t, "UK", profile["country"])
		assert.Equal(t, 0.5, profile["marketCapMultiplier"], "UK should have 0.5x market cap multiplier")
	})

	t.Run("Can disable adjustments with adjust=false", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/screeners/small-cap-growth?country=Israel&adjust=false", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Should not be adjusted when explicitly disabled
		assert.False(t, response["adjusted"].(bool), "Should not be adjusted when adjust=false")
	})

	t.Run("All supported markets have correct profiles", func(t *testing.T) {
		countries := []string{"Israel", "UK", "Germany", "Japan", "China", "India", "Brazil", "Canada", "Switzerland", "France", "Australia"}

		for _, country := range countries {
			req := httptest.NewRequest(http.MethodGet, "/screeners/piotroski-high-score?country="+country, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Should succeed for "+country)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.True(t, response["adjusted"].(bool), country+" should be adjusted")

			profile := response["marketProfile"].(map[string]interface{})
			assert.Equal(t, country, profile["country"], "Profile country should match")
			assert.Greater(t, profile["marketCapMultiplier"].(float64), 0.0, "Multiplier should be positive for "+country)
			t.Logf("%s: marketCapMultiplier=%.2f", country, profile["marketCapMultiplier"].(float64))
		}
	})
}
