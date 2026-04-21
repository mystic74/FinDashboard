package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupProfileTestRouter() (*gin.Engine, *ProfileHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewProfileHandler()
	return router, handler
}

func TestGetAllProfiles(t *testing.T) {
	router, handler := setupProfileTestRouter()
	router.GET("/profiles", handler.GetAllProfiles)

	req := httptest.NewRequest(http.MethodGet, "/profiles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	profiles := response["profiles"].([]interface{})
	assert.GreaterOrEqual(t, len(profiles), 10, "Should have at least 10 market profiles")
}

func TestGetProfile(t *testing.T) {
	router, handler := setupProfileTestRouter()
	router.GET("/profiles/:country", handler.GetProfile)

	t.Run("Get USA profile", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/profiles/USA", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))
		assert.False(t, response["isCustom"].(bool))

		profile := response["profile"].(map[string]interface{})
		assert.Equal(t, "USA", profile["country"])
		assert.Equal(t, 1.0, profile["marketCapMultiplier"])
	})

	t.Run("Get Israel profile", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/profiles/Israel", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		profile := response["profile"].(map[string]interface{})
		assert.Equal(t, "Israel", profile["country"])
		assert.Equal(t, 0.1, profile["marketCapMultiplier"])
	})

	t.Run("Unknown country returns USA profile", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/profiles/Unknown", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		profile := response["profile"].(map[string]interface{})
		assert.Equal(t, "USA", profile["country"]) // Falls back to USA
	})
}

func TestUpdateProfile(t *testing.T) {
	router, handler := setupProfileTestRouter()
	router.PUT("/profiles/:country", handler.UpdateProfile)
	router.GET("/profiles/:country", handler.GetProfile)

	t.Run("Cannot update USA profile", func(t *testing.T) {
		body := []byte(`{"marketCapMultiplier": 0.5}`)
		req := httptest.NewRequest(http.MethodPut, "/profiles/USA", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"], "baseline")
	})

	t.Run("Can update Israel profile", func(t *testing.T) {
		body := []byte(`{"marketCapMultiplier": 0.15, "volumeMultiplier": 0.4}`)
		req := httptest.NewRequest(http.MethodPut, "/profiles/Israel", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		profile := response["profile"].(map[string]interface{})
		assert.Equal(t, 0.15, profile["marketCapMultiplier"])
		assert.Equal(t, 0.4, profile["volumeMultiplier"])
	})

	t.Run("Updated profile is persisted", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/profiles/Israel", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["isCustom"].(bool))

		profile := response["profile"].(map[string]interface{})
		assert.Equal(t, 0.15, profile["marketCapMultiplier"])
	})
}

func TestResetProfile(t *testing.T) {
	router, handler := setupProfileTestRouter()
	router.PUT("/profiles/:country", handler.UpdateProfile)
	router.POST("/profiles/:country/reset", handler.ResetProfile)
	router.GET("/profiles/:country", handler.GetProfile)

	// First update the profile
	body := []byte(`{"marketCapMultiplier": 0.5}`)
	req := httptest.NewRequest(http.MethodPut, "/profiles/UK", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Reset it
	req = httptest.NewRequest(http.MethodPost, "/profiles/UK/reset", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	profile := response["profile"].(map[string]interface{})
	assert.Equal(t, 0.5, profile["marketCapMultiplier"]) // Back to default

	// Verify it's no longer custom
	req = httptest.NewRequest(http.MethodGet, "/profiles/UK", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["isCustom"].(bool))
}

func TestResetAllProfiles(t *testing.T) {
	router, handler := setupProfileTestRouter()
	router.PUT("/profiles/:country", handler.UpdateProfile)
	router.POST("/profiles/reset", handler.ResetAllProfiles)
	router.GET("/profiles/:country", handler.GetProfile)

	// Update multiple profiles
	countries := []string{"UK", "Germany", "Japan"}
	for _, country := range countries {
		body := []byte(`{"marketCapMultiplier": 0.99}`)
		req := httptest.NewRequest(http.MethodPut, "/profiles/"+country, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// Reset all
	req := httptest.NewRequest(http.MethodPost, "/profiles/reset", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify all are back to defaults
	for _, country := range countries {
		req = httptest.NewRequest(http.MethodGet, "/profiles/"+country, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response["isCustom"].(bool), country+" should not be custom after reset")
	}
}
