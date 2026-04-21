package handlers

import (
	"net/http"
	"stock-screener/models"
	"sync"

	"github.com/gin-gonic/gin"
)

// ProfileHandler handles market profile-related requests
type ProfileHandler struct {
	customProfiles map[string]*models.MarketProfile
	mu             sync.RWMutex
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{
		customProfiles: make(map[string]*models.MarketProfile),
	}
}

// GetAllProfiles returns all market profiles
// @Summary Get all market profiles
// @Description Returns all available market profiles with their multipliers
// @Tags Profiles
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/profiles [get]
func (h *ProfileHandler) GetAllProfiles(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Build profiles list, applying any custom overrides
	profiles := make([]models.MarketProfile, 0)

	for country, profile := range models.MarketProfiles {
		p := *profile
		// Apply any custom overrides
		if custom, ok := h.customProfiles[country]; ok {
			p = *custom
		}
		profiles = append(profiles, p)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"profiles": profiles,
		"count":    len(profiles),
	})
}

// GetProfile returns a specific market profile
// @Summary Get a market profile
// @Description Returns the market profile for a specific country
// @Tags Profiles
// @Produce json
// @Param country path string true "Country code"
// @Success 200 {object} models.MarketProfile
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/profiles/{country} [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	country := c.Param("country")

	h.mu.RLock()
	defer h.mu.RUnlock()

	// Check for custom override first
	if custom, ok := h.customProfiles[country]; ok {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"profile": custom,
			"isCustom": true,
		})
		return
	}

	// Get default profile
	profile := models.GetMarketProfile(country)
	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Profile not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"profile":  profile,
		"isCustom": false,
	})
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	MarketCapMultiplier *float64 `json:"marketCapMultiplier"`
	VolumeMultiplier    *float64 `json:"volumeMultiplier"`
	DividendMultiplier  *float64 `json:"dividendMultiplier"`
	GrowthMultiplier    *float64 `json:"growthMultiplier"`
}

// UpdateProfile updates a market profile
// @Summary Update a market profile
// @Description Updates the multipliers for a specific country's market profile
// @Tags Profiles
// @Accept json
// @Produce json
// @Param country path string true "Country code"
// @Param request body UpdateProfileRequest true "Profile updates"
// @Success 200 {object} models.MarketProfile
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/profiles/{country} [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	country := c.Param("country")

	// USA cannot be modified (baseline)
	if country == "USA" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "USA is the baseline profile and cannot be modified",
		})
		return
	}

	// Get base profile
	baseProfile := models.GetMarketProfile(country)
	if baseProfile == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Profile not found",
		})
		return
	}

	var request UpdateProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body: " + err.Error(),
		})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Create or update custom profile
	customProfile, exists := h.customProfiles[country]
	if !exists {
		// Clone base profile
		cp := *baseProfile
		customProfile = &cp
	}

	// Apply updates
	if request.MarketCapMultiplier != nil {
		customProfile.MarketCapMultiplier = *request.MarketCapMultiplier
	}
	if request.VolumeMultiplier != nil {
		customProfile.VolumeMultiplier = *request.VolumeMultiplier
	}
	if request.DividendMultiplier != nil {
		customProfile.DividendMultiplier = *request.DividendMultiplier
	}
	if request.GrowthMultiplier != nil {
		customProfile.GrowthMultiplier = *request.GrowthMultiplier
	}

	h.customProfiles[country] = customProfile

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"profile": customProfile,
		"message": "Profile updated successfully",
	})
}

// ResetProfile resets a market profile to defaults
// @Summary Reset a market profile
// @Description Resets a country's market profile to default values
// @Tags Profiles
// @Produce json
// @Param country path string true "Country code"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/profiles/{country}/reset [post]
func (h *ProfileHandler) ResetProfile(c *gin.Context) {
	country := c.Param("country")

	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.customProfiles, country)

	profile := models.GetMarketProfile(country)
	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Profile not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"profile": profile,
		"message": "Profile reset to defaults",
	})
}

// ResetAllProfiles resets all market profiles to defaults
// @Summary Reset all market profiles
// @Description Resets all market profiles to their default values
// @Tags Profiles
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/profiles/reset [post]
func (h *ProfileHandler) ResetAllProfiles(c *gin.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.customProfiles = make(map[string]*models.MarketProfile)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All profiles reset to defaults",
	})
}

// GetCustomProfile returns the custom profile if it exists, otherwise the default
// This is used internally by the screener handler
func (h *ProfileHandler) GetCustomProfile(country string) *models.MarketProfile {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if custom, ok := h.customProfiles[country]; ok {
		return custom
	}
	return models.GetMarketProfile(country)
}

// GetProfileForCountry returns a custom/default profile only when the country exists.
// Unlike GetCustomProfile, this does not fallback unknown countries to USA.
func (h *ProfileHandler) GetProfileForCountry(country string) (*models.MarketProfile, bool) {
	if country == "" {
		return nil, false
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if custom, ok := h.customProfiles[country]; ok {
		return custom, true
	}

	defaultProfile, ok := models.MarketProfiles[country]
	return defaultProfile, ok
}
