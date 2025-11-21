package handlers

import (
	"net/http"
	"stock-screener/models"

	"github.com/gin-gonic/gin"
)

// FilterHandler handles filter-related requests
type FilterHandler struct{}

// NewFilterHandler creates a new filter handler
func NewFilterHandler() *FilterHandler {
	return &FilterHandler{}
}

// GetAllFilters returns all available filter definitions
// @Summary Get all filters
// @Description Returns all available filter definitions with metadata
// @Tags Filters
// @Produce json
// @Success 200 {array} models.FilterDefinition
// @Router /api/v1/filters [get]
func (h *FilterHandler) GetAllFilters(c *gin.Context) {
	filters := models.GetAllFilterDefinitions()

	// Group by category
	grouped := make(map[string][]models.FilterDefinition)
	for _, f := range filters {
		category := string(f.Category)
		grouped[category] = append(grouped[category], f)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"filters":  filters,
		"grouped":  grouped,
		"count":    len(filters),
	})
}

// GetFiltersByCategory returns filters for a specific category
// @Summary Get filters by category
// @Description Returns filters for a specific category
// @Tags Filters
// @Produce json
// @Param category path string true "Filter category"
// @Success 200 {array} models.FilterDefinition
// @Router /api/v1/filters/{category} [get]
func (h *FilterHandler) GetFiltersByCategory(c *gin.Context) {
	category := c.Param("category")

	allFilters := models.GetAllFilterDefinitions()
	var categoryFilters []models.FilterDefinition

	for _, f := range allFilters {
		if string(f.Category) == category {
			categoryFilters = append(categoryFilters, f)
		}
	}

	if len(categoryFilters) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"category": category,
		"filters":  categoryFilters,
		"count":    len(categoryFilters),
	})
}

// GetFilterCategories returns all filter categories
// @Summary Get filter categories
// @Description Returns all available filter categories
// @Tags Filters
// @Produce json
// @Success 200 {array} string
// @Router /api/v1/filters/categories [get]
func (h *FilterHandler) GetFilterCategories(c *gin.Context) {
	categories := []gin.H{
		{"id": "price_volume", "name": "Price & Volume", "description": "Price and volume-related filters"},
		{"id": "valuation", "name": "Valuation", "description": "Valuation ratios and metrics"},
		{"id": "dividends", "name": "Dividends", "description": "Dividend-related filters"},
		{"id": "financial_health", "name": "Financial Health", "description": "Financial stability and health metrics"},
		{"id": "profitability", "name": "Profitability", "description": "Profitability ratios"},
		{"id": "growth", "name": "Growth", "description": "Growth metrics"},
		{"id": "technical", "name": "Technical", "description": "Technical indicators and performance"},
		{"id": "profile", "name": "Profile", "description": "Company profile information"},
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"categories": categories,
	})
}

// GetSectors returns all available sectors
// @Summary Get sectors
// @Description Returns all available stock sectors
// @Tags Filters
// @Produce json
// @Success 200 {array} string
// @Router /api/v1/sectors/list [get]
func (h *FilterHandler) GetSectors(c *gin.Context) {
	sectors := models.GetSectors()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"sectors": sectors,
		"count":   len(sectors),
	})
}

// GetMarketCapRanges returns predefined market cap ranges
// @Summary Get market cap ranges
// @Description Returns predefined market cap ranges
// @Tags Filters
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/filters/marketcap-ranges [get]
func (h *FilterHandler) GetMarketCapRanges(c *gin.Context) {
	ranges := models.GetMarketCapRanges()

	rangesList := []gin.H{
		{"id": "nano", "name": "Nano Cap", "min": ranges["nano"][0], "max": ranges["nano"][1]},
		{"id": "micro", "name": "Micro Cap", "min": ranges["micro"][0], "max": ranges["micro"][1]},
		{"id": "small", "name": "Small Cap", "min": ranges["small"][0], "max": ranges["small"][1]},
		{"id": "mid", "name": "Mid Cap", "min": ranges["mid"][0], "max": ranges["mid"][1]},
		{"id": "large", "name": "Large Cap", "min": ranges["large"][0], "max": ranges["large"][1]},
		{"id": "mega", "name": "Mega Cap", "min": ranges["mega"][0], "max": ranges["mega"][1]},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"ranges":  rangesList,
	})
}
