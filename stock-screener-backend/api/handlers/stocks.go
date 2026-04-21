package handlers

import (
	"net/http"
	"stock-screener/services"
	"stock-screener/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// StockHandler handles stock-related requests
type StockHandler struct {
	yahooService *services.YahooFinanceService
}

// NewStockHandler creates a new stock handler
func NewStockHandler(yahooService *services.YahooFinanceService) *StockHandler {
	return &StockHandler{yahooService: yahooService}
}

// GetStock returns data for a single stock
// @Summary Get stock data
// @Description Returns current data for a single stock
// @Tags Stocks
// @Produce json
// @Param symbol path string true "Stock symbol"
// @Success 200 {object} models.Stock
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/stocks/{symbol} [get]
func (h *StockHandler) GetStock(c *gin.Context) {
	symbol := strings.ToUpper(c.Param("symbol"))

	if err := utils.ValidateSymbol(symbol); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	stocks, err := h.yahooService.GetQuotes([]string{symbol})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if len(stocks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Stock not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"stock":   stocks[0],
	})
}

// GetStockFundamentals returns detailed fundamental data for a stock
// @Summary Get stock fundamentals
// @Description Returns detailed fundamental data for a single stock
// @Tags Stocks
// @Produce json
// @Param symbol path string true "Stock symbol"
// @Success 200 {object} models.Stock
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/stocks/{symbol}/fundamentals [get]
func (h *StockHandler) GetStockFundamentals(c *gin.Context) {
	symbol := strings.ToUpper(c.Param("symbol"))

	if err := utils.ValidateSymbol(symbol); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	stock, err := h.yahooService.GetStockFundamentals(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Calculate additional metrics
	stock.PiotroskiFScore = services.CalculatePiotroskiScore(stock)
	stock.AltmanZScore = services.CalculateAltmanZ(stock)

	// Validate data quality
	validation := stock.Validate()

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"stock":      stock,
		"validation": validation,
	})
}

// GetMultipleStocks returns data for multiple stocks
// @Summary Get multiple stocks
// @Description Returns data for multiple stocks at once
// @Tags Stocks
// @Produce json
// @Param symbols query string true "Comma-separated stock symbols"
// @Success 200 {array} models.Stock
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/stocks [get]
func (h *StockHandler) GetMultipleStocks(c *gin.Context) {
	symbolsParam := c.Query("symbols")
	if symbolsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "symbols parameter is required",
		})
		return
	}

	symbols := strings.Split(symbolsParam, ",")
	for i := range symbols {
		symbols[i] = strings.TrimSpace(strings.ToUpper(symbols[i]))
	}

	// Validate symbols
	validSymbols, validationErrors := utils.ValidateSymbols(symbols)
	if len(validSymbols) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No valid symbols provided",
			"errors":  validationErrors,
		})
		return
	}

	stocks, err := h.yahooService.GetQuotes(validSymbols)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := gin.H{
		"success": true,
		"stocks":  stocks,
		"count":   len(stocks),
	}

	if len(validationErrors) > 0 {
		response["warnings"] = validationErrors
	}

	c.JSON(http.StatusOK, response)
}

// GetHistoricalPrices returns historical price data for a stock
// @Summary Get historical prices
// @Description Returns historical price data for a stock
// @Tags Stocks
// @Produce json
// @Param symbol path string true "Stock symbol"
// @Param period query string false "Time period (1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, ytd, max)" default(1y)
// @Success 200 {array} models.HistoricalPrice
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/stocks/{symbol}/history [get]
func (h *StockHandler) GetHistoricalPrices(c *gin.Context) {
	symbol := strings.ToUpper(c.Param("symbol"))
	period := c.DefaultQuery("period", "1y")

	if err := utils.ValidateSymbol(symbol); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	validPeriods := map[string]bool{
		"1d": true, "5d": true, "1mo": true, "3mo": true,
		"6mo": true, "1y": true, "2y": true, "5y": true,
		"10y": true, "ytd": true, "max": true,
	}
	if !validPeriods[period] {
		period = "1y"
	}

	prices, err := h.yahooService.GetHistoricalPrices(symbol, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"symbol":  symbol,
		"period":  period,
		"prices":  prices,
		"count":   len(prices),
	})
}

// GetStockQuote returns a quick quote for a stock
// @Summary Get stock quote
// @Description Returns a quick quote for a stock
// @Tags Stocks
// @Produce json
// @Param symbol path string true "Stock symbol"
// @Success 200 {object} models.StockQuote
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/stocks/{symbol}/quote [get]
func (h *StockHandler) GetStockQuote(c *gin.Context) {
	symbol := strings.ToUpper(c.Param("symbol"))

	if err := utils.ValidateSymbol(symbol); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	stocks, err := h.yahooService.GetQuotes([]string{symbol})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if len(stocks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Stock not found",
		})
		return
	}

	stock := stocks[0]
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"quote": gin.H{
			"symbol":        stock.Symbol,
			"name":          stock.Name,
			"price":         stock.Price,
			"change":        stock.Change,
			"changePercent": stock.ChangePercent,
			"volume":        stock.Volume,
			"marketCap":     stock.MarketCap,
		},
	})
}
