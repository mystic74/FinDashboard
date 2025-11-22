package services

import (
	"math/rand"
	"stock-screener/models"
	"time"
)

// MockDataService provides mock stock data for demo purposes
type MockDataService struct {
	stocks []models.Stock
}

// NewMockDataService creates a new mock data service with realistic stock data
func NewMockDataService() *MockDataService {
	m := &MockDataService{}
	m.generateMockStocks()
	return m
}

// GetQuotes returns mock quotes for the given symbols
func (m *MockDataService) GetQuotes(symbols []string) ([]models.Stock, error) {
	if len(symbols) == 0 {
		return m.stocks, nil
	}

	symbolSet := make(map[string]bool)
	for _, s := range symbols {
		symbolSet[s] = true
	}

	var result []models.Stock
	for _, stock := range m.stocks {
		if symbolSet[stock.Symbol] {
			result = append(result, stock)
		}
	}
	return result, nil
}

// GetAllStocks returns all mock stocks
func (m *MockDataService) GetAllStocks() []models.Stock {
	return m.stocks
}

func (m *MockDataService) generateMockStocks() {
	rand.Seed(time.Now().UnixNano())

	// Generate 100 realistic mock stocks
	mockData := []struct {
		symbol        string
		name          string
		sector        string
		industry      string
		priceBase     float64
		marketCapBase int64
	}{
		// Technology
		{"AAPL", "Apple Inc.", "Technology", "Consumer Electronics", 178.50, 2800000000000},
		{"MSFT", "Microsoft Corporation", "Technology", "Software - Infrastructure", 378.25, 2810000000000},
		{"GOOGL", "Alphabet Inc.", "Technology", "Internet Content & Information", 141.80, 1790000000000},
		{"AMZN", "Amazon.com Inc.", "Technology", "Internet Retail", 178.35, 1860000000000},
		{"META", "Meta Platforms Inc.", "Technology", "Internet Content & Information", 505.75, 1300000000000},
		{"NVDA", "NVIDIA Corporation", "Technology", "Semiconductors", 495.22, 1220000000000},
		{"TSLA", "Tesla Inc.", "Technology", "Auto Manufacturers", 238.45, 760000000000},
		{"AMD", "Advanced Micro Devices", "Technology", "Semiconductors", 121.33, 196000000000},
		{"INTC", "Intel Corporation", "Technology", "Semiconductors", 45.12, 191000000000},
		{"CRM", "Salesforce Inc.", "Technology", "Software - Application", 265.88, 258000000000},
		{"ORCL", "Oracle Corporation", "Technology", "Software - Infrastructure", 118.45, 326000000000},
		{"ADBE", "Adobe Inc.", "Technology", "Software - Application", 582.15, 259000000000},
		{"CSCO", "Cisco Systems Inc.", "Technology", "Communication Equipment", 48.92, 198000000000},
		{"IBM", "International Business Machines", "Technology", "IT Services", 166.78, 152000000000},
		{"QCOM", "Qualcomm Inc.", "Technology", "Semiconductors", 152.35, 170000000000},

		// Healthcare
		{"JNJ", "Johnson & Johnson", "Healthcare", "Drug Manufacturers", 156.42, 377000000000},
		{"UNH", "UnitedHealth Group", "Healthcare", "Healthcare Plans", 528.90, 489000000000},
		{"PFE", "Pfizer Inc.", "Healthcare", "Drug Manufacturers", 28.55, 161000000000},
		{"ABBV", "AbbVie Inc.", "Healthcare", "Drug Manufacturers", 154.32, 273000000000},
		{"MRK", "Merck & Co.", "Healthcare", "Drug Manufacturers", 108.65, 275000000000},
		{"LLY", "Eli Lilly and Company", "Healthcare", "Drug Manufacturers", 598.22, 568000000000},
		{"TMO", "Thermo Fisher Scientific", "Healthcare", "Diagnostics & Research", 535.45, 206000000000},
		{"ABT", "Abbott Laboratories", "Healthcare", "Medical Devices", 105.78, 183000000000},
		{"BMY", "Bristol-Myers Squibb", "Healthcare", "Drug Manufacturers", 51.22, 104000000000},
		{"AMGN", "Amgen Inc.", "Healthcare", "Drug Manufacturers", 278.95, 149000000000},

		// Financial Services
		{"JPM", "JPMorgan Chase & Co.", "Financial Services", "Banks - Diversified", 172.45, 498000000000},
		{"BAC", "Bank of America Corp", "Financial Services", "Banks - Diversified", 33.88, 265000000000},
		{"WFC", "Wells Fargo & Company", "Financial Services", "Banks - Diversified", 46.12, 168000000000},
		{"GS", "Goldman Sachs Group", "Financial Services", "Capital Markets", 385.62, 128000000000},
		{"MS", "Morgan Stanley", "Financial Services", "Capital Markets", 87.95, 146000000000},
		{"V", "Visa Inc.", "Financial Services", "Credit Services", 258.45, 531000000000},
		{"MA", "Mastercard Inc.", "Financial Services", "Credit Services", 422.80, 395000000000},
		{"AXP", "American Express", "Financial Services", "Credit Services", 178.92, 133000000000},
		{"BLK", "BlackRock Inc.", "Financial Services", "Asset Management", 782.35, 118000000000},
		{"SCHW", "Charles Schwab Corp", "Financial Services", "Capital Markets", 66.55, 121000000000},

		// Consumer
		{"WMT", "Walmart Inc.", "Consumer Defensive", "Discount Stores", 165.22, 446000000000},
		{"HD", "The Home Depot", "Consumer Cyclical", "Home Improvement Retail", 345.88, 343000000000},
		{"COST", "Costco Wholesale", "Consumer Defensive", "Discount Stores", 572.35, 254000000000},
		{"NKE", "Nike Inc.", "Consumer Cyclical", "Footwear & Accessories", 98.75, 150000000000},
		{"MCD", "McDonald's Corporation", "Consumer Cyclical", "Restaurants", 298.45, 214000000000},
		{"SBUX", "Starbucks Corporation", "Consumer Cyclical", "Restaurants", 95.62, 109000000000},
		{"TGT", "Target Corporation", "Consumer Defensive", "Discount Stores", 142.88, 66000000000},
		{"LOW", "Lowe's Companies", "Consumer Cyclical", "Home Improvement Retail", 218.45, 128000000000},
		{"DIS", "The Walt Disney Company", "Communication Services", "Entertainment", 91.35, 167000000000},
		{"NFLX", "Netflix Inc.", "Communication Services", "Entertainment", 478.92, 212000000000},
		{"PEP", "PepsiCo Inc.", "Consumer Defensive", "Beverages - Non-Alcoholic", 178.65, 245000000000},
		{"KO", "The Coca-Cola Company", "Consumer Defensive", "Beverages - Non-Alcoholic", 60.22, 261000000000},

		// Energy
		{"XOM", "Exxon Mobil Corporation", "Energy", "Oil & Gas Integrated", 104.55, 418000000000},
		{"CVX", "Chevron Corporation", "Energy", "Oil & Gas Integrated", 148.32, 276000000000},
		{"COP", "ConocoPhillips", "Energy", "Oil & Gas E&P", 115.78, 133000000000},
		{"SLB", "Schlumberger NV", "Energy", "Oil & Gas Equipment", 52.45, 75000000000},
		{"EOG", "EOG Resources", "Energy", "Oil & Gas E&P", 122.88, 72000000000},

		// Industrials
		{"CAT", "Caterpillar Inc.", "Industrials", "Farm & Heavy Construction", 278.45, 138000000000},
		{"DE", "Deere & Company", "Industrials", "Farm & Heavy Construction", 385.22, 112000000000},
		{"BA", "The Boeing Company", "Industrials", "Aerospace & Defense", 205.88, 125000000000},
		{"HON", "Honeywell International", "Industrials", "Conglomerates", 198.75, 130000000000},
		{"UPS", "United Parcel Service", "Industrials", "Integrated Freight", 155.32, 134000000000},
		{"GE", "General Electric", "Industrials", "Aerospace & Defense", 122.45, 133000000000},
		{"LMT", "Lockheed Martin", "Industrials", "Aerospace & Defense", 452.88, 108000000000},
		{"RTX", "RTX Corporation", "Industrials", "Aerospace & Defense", 88.92, 130000000000},

		// Real Estate
		{"AMT", "American Tower Corp", "Real Estate", "REIT - Specialty", 198.55, 92000000000},
		{"PLD", "Prologis Inc.", "Real Estate", "REIT - Industrial", 122.35, 113000000000},
		{"CCI", "Crown Castle Inc.", "Real Estate", "REIT - Specialty", 108.72, 47000000000},
		{"EQIX", "Equinix Inc.", "Real Estate", "REIT - Specialty", 782.45, 73000000000},
		{"SPG", "Simon Property Group", "Real Estate", "REIT - Retail", 142.88, 46000000000},
		{"O", "Realty Income Corp", "Real Estate", "REIT - Retail", 55.62, 44000000000},

		// Utilities
		{"NEE", "NextEra Energy", "Utilities", "Utilities - Regulated Electric", 68.45, 141000000000},
		{"DUK", "Duke Energy Corp", "Utilities", "Utilities - Regulated Electric", 98.22, 76000000000},
		{"SO", "Southern Company", "Utilities", "Utilities - Regulated Electric", 72.88, 79000000000},
		{"D", "Dominion Energy", "Utilities", "Utilities - Regulated Electric", 48.55, 40000000000},
		{"AEP", "American Electric Power", "Utilities", "Utilities - Regulated Electric", 85.32, 44000000000},

		// Small Cap Growth
		{"DDOG", "Datadog Inc.", "Technology", "Software - Application", 118.45, 38000000000},
		{"SNOW", "Snowflake Inc.", "Technology", "Software - Infrastructure", 162.88, 53000000000},
		{"CRWD", "CrowdStrike Holdings", "Technology", "Software - Infrastructure", 228.35, 55000000000},
		{"ZS", "Zscaler Inc.", "Technology", "Software - Infrastructure", 195.72, 29000000000},
		{"OKTA", "Okta Inc.", "Technology", "Software - Infrastructure", 88.45, 14000000000},
		{"NET", "Cloudflare Inc.", "Technology", "Software - Infrastructure", 78.92, 26000000000},
		{"MDB", "MongoDB Inc.", "Technology", "Software - Infrastructure", 398.55, 28000000000},

		// Value Stocks
		{"BRK.B", "Berkshire Hathaway B", "Financial Services", "Insurance - Diversified", 358.92, 789000000000},
		{"T", "AT&T Inc.", "Communication Services", "Telecom Services", 17.45, 125000000000},
		{"VZ", "Verizon Communications", "Communication Services", "Telecom Services", 38.88, 163000000000},
		{"TMUS", "T-Mobile US Inc.", "Communication Services", "Telecom Services", 162.55, 191000000000},

		// Dividend Aristocrats
		{"PG", "Procter & Gamble", "Consumer Defensive", "Household Products", 152.78, 361000000000},
		{"MMM", "3M Company", "Industrials", "Conglomerates", 98.45, 54000000000},
		{"EMR", "Emerson Electric", "Industrials", "Electrical Equipment", 98.22, 57000000000},
		{"CLX", "The Clorox Company", "Consumer Defensive", "Household Products", 142.55, 17600000000},
		{"KMB", "Kimberly-Clark", "Consumer Defensive", "Household Products", 128.35, 43000000000},
		{"SYY", "Sysco Corporation", "Consumer Defensive", "Food Distribution", 75.88, 38000000000},
		{"AFL", "Aflac Inc.", "Financial Services", "Insurance - Life", 82.45, 50000000000},
		{"CINF", "Cincinnati Financial", "Financial Services", "Insurance - Property", 112.72, 17500000000},
		{"ED", "Consolidated Edison", "Utilities", "Utilities - Regulated Electric", 92.35, 32000000000},
		{"XEL", "Xcel Energy", "Utilities", "Utilities - Regulated Electric", 62.88, 35000000000},
	}

	m.stocks = make([]models.Stock, 0, len(mockData))

	for _, data := range mockData {
		stock := m.generateStock(data.symbol, data.name, data.sector, data.industry, data.priceBase, data.marketCapBase)
		m.stocks = append(m.stocks, stock)
	}
}

func (m *MockDataService) generateStock(symbol, name, sector, industry string, priceBase float64, marketCapBase int64) models.Stock {
	// Add some randomness
	priceVariation := 0.95 + rand.Float64()*0.1
	price := priceBase * priceVariation

	// Generate realistic metrics
	peRatio := 10 + rand.Float64()*40
	pbRatio := 0.5 + rand.Float64()*10
	dividendYield := rand.Float64() * 6
	beta := 0.5 + rand.Float64()*1.5

	roe := 5 + rand.Float64()*30
	roa := 2 + rand.Float64()*15
	grossMargin := 20 + rand.Float64()*50
	netMargin := 5 + rand.Float64()*25
	operatingMargin := 10 + rand.Float64()*30
	currentRatio := 0.8 + rand.Float64()*2
	debtToEquity := rand.Float64() * 2

	revenueGrowth := -10 + rand.Float64()*40
	epsGrowth := -20 + rand.Float64()*50

	// Returns
	return1W := -5 + rand.Float64()*10
	return1M := -10 + rand.Float64()*20
	return3M := -15 + rand.Float64()*40
	return6M := -20 + rand.Float64()*50
	return1Y := -30 + rand.Float64()*60

	// RSI
	rsi := 30 + rand.Float64()*40

	// 52-week range
	week52Low := price * (0.6 + rand.Float64()*0.3)
	week52High := price * (1.05 + rand.Float64()*0.4)

	// Volume
	volume := int64(float64(marketCapBase) * (0.001 + rand.Float64()*0.005) / price)

	// Calculate Piotroski score based on metrics
	piotroskiScore := 0
	if roe > 0 {
		piotroskiScore++
	}
	if roa > 0 {
		piotroskiScore++
	}
	if grossMargin > 20 {
		piotroskiScore++
	}
	if currentRatio > 1 {
		piotroskiScore++
	}
	if debtToEquity < 0.5 {
		piotroskiScore++
	}
	if epsGrowth > 0 {
		piotroskiScore++
	}
	if revenueGrowth > 0 {
		piotroskiScore++
	}
	if netMargin > 10 {
		piotroskiScore++
	}
	if operatingMargin > 15 {
		piotroskiScore++
	}

	return models.Stock{
		Symbol:            symbol,
		Name:              name,
		Exchange:          "NASDAQ",
		Currency:          "USD",
		Price:             round2(price),
		Change:            round2(price * (return1W / 100 / 5)),
		ChangePercent:     round2(return1W / 5),
		Volume:            volume,
		AvgVolume:         int64(float64(volume) * (0.8 + rand.Float64()*0.4)),
		MarketCap:         marketCapBase,
		Week52High:        round2(week52High),
		Week52Low:         round2(week52Low),
		MA50:              round2(price * (0.95 + rand.Float64()*0.1)),
		MA200:             round2(price * (0.9 + rand.Float64()*0.2)),
		PERatio:           round2(peRatio),
		ForwardPE:         round2(peRatio * (0.8 + rand.Float64()*0.3)),
		PEGRatio:          round2(1 + rand.Float64()*2),
		PBRatio:           round2(pbRatio),
		PSRatio:           round2(1 + rand.Float64()*10),
		DividendYield:     round2(dividendYield),
		PayoutRatio:       round2(dividendYield * 10 * (0.5 + rand.Float64())),
		Beta:              round2(beta),
		ROE:               round2(roe),
		ROA:               round2(roa),
		GrossMargin:       round2(grossMargin),
		NetMargin:         round2(netMargin),
		OperatingMargin:   round2(operatingMargin),
		CurrentRatio:      round2(currentRatio),
		QuickRatio:        round2(currentRatio * (0.7 + rand.Float64()*0.3)),
		DebtToEquity:      round2(debtToEquity),
		RevenueGrowth:     round2(revenueGrowth),
		EPSGrowth:         round2(epsGrowth),
		RSI14:             round2(rsi),
		Return1W:          round2(return1W),
		Return1M:          round2(return1M),
		Return3M:          round2(return3M),
		Return6M:          round2(return6M),
		Return1Y:          round2(return1Y),
		Sector:            sector,
		Industry:          industry,
		PiotroskiFScore:   piotroskiScore,
		FreeCashFlow:      int64(float64(marketCapBase) * 0.03 * (0.5 + rand.Float64())),
		OperatingCashFlow: int64(float64(marketCapBase) * 0.05 * (0.5 + rand.Float64())),
		TotalCash:         int64(float64(marketCapBase) * 0.1 * (0.3 + rand.Float64())),
		TotalDebt:         int64(float64(marketCapBase) * 0.2 * rand.Float64()),
		LastUpdated:       time.Now(),
	}
}

func round2(val float64) float64 {
	return float64(int(val*100)) / 100
}
