package handlers

import (
	"net/http"
	"stock-screener/models"
	"stock-screener/services"

	"github.com/gin-gonic/gin"
)

// ScreenerHandler handles screener-related requests
type ScreenerHandler struct {
	engine *services.ScreenerEngine
}

// NewScreenerHandler creates a new screener handler
func NewScreenerHandler(engine *services.ScreenerEngine) *ScreenerHandler {
	return &ScreenerHandler{engine: engine}
}

// GetAllScreeners returns all predefined screeners
// @Summary Get all predefined screeners
// @Description Returns a list of all available predefined screeners
// @Tags Screeners
// @Produce json
// @Success 200 {array} models.Screener
// @Router /api/v1/screeners [get]
func (h *ScreenerHandler) GetAllScreeners(c *gin.Context) {
	screeners := h.engine.GetAllPredefinedScreeners()
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"screeners": screeners,
		"count":     len(screeners),
	})
}

// GetScreenersSummary returns summaries of all predefined screeners
// @Summary Get screener summaries
// @Description Returns summaries of all predefined screeners with match counts
// @Tags Screeners
// @Produce json
// @Success 200 {array} models.ScreenerSummary
// @Router /api/v1/screeners/summary [get]
func (h *ScreenerHandler) GetScreenersSummary(c *gin.Context) {
	summaries, err := h.engine.GetScreenersSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"summaries": summaries,
	})
}

// RunScreener runs a specific predefined screener
// @Summary Run a predefined screener
// @Description Runs a predefined screener and returns matching stocks
// @Tags Screeners
// @Produce json
// @Param name path string true "Screener ID"
// @Success 200 {object} models.ScreenerResult
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/screeners/{name} [get]
func (h *ScreenerHandler) RunScreener(c *gin.Context) {
	screenerID := c.Param("name")

	screener, found := h.engine.GetScreenerByID(screenerID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Screener not found",
		})
		return
	}

	result, err := h.engine.RunScreener(c.Request.Context(), *screener)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}

// RunCustomScreener runs a custom screener with user-defined filters
// @Summary Run a custom screener
// @Description Runs a custom screener with user-defined filters
// @Tags Screeners
// @Accept json
// @Produce json
// @Param request body models.FilterRequest true "Filter request"
// @Success 200 {object} models.FilterResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/screeners/custom [post]
func (h *ScreenerHandler) RunCustomScreener(c *gin.Context) {
	var request models.FilterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body: " + err.Error(),
		})
		return
	}

	result, err := h.engine.RunCustomScreener(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}

// GetSectorPerformance returns sector performance data
// @Summary Get sector performance
// @Description Returns performance data grouped by sector
// @Tags Screeners
// @Produce json
// @Success 200 {array} models.SectorPerformance
// @Router /api/v1/sectors [get]
func (h *ScreenerHandler) GetSectorPerformance(c *gin.Context) {
	performance, err := h.engine.GetSectorPerformance(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"sectors":  performance,
		"count":    len(performance),
	})
}
