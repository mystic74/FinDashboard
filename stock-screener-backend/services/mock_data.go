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

	// Generate realistic mock stocks from multiple markets
	mockData := []struct {
		symbol        string
		name          string
		sector        string
		industry      string
		country       string
		exchange      string
		currency      string
		priceBase     float64
		marketCapBase int64
	}{
		// =====================================================================
		// USA - Technology
		// =====================================================================
		{"AAPL", "Apple Inc.", "Technology", "Consumer Electronics", "USA", "NASDAQ", "USD", 178.50, 2800000000000},
		{"MSFT", "Microsoft Corporation", "Technology", "Software - Infrastructure", "USA", "NASDAQ", "USD", 378.25, 2810000000000},
		{"GOOGL", "Alphabet Inc.", "Technology", "Internet Content & Information", "USA", "NASDAQ", "USD", 141.80, 1790000000000},
		{"AMZN", "Amazon.com Inc.", "Technology", "Internet Retail", "USA", "NASDAQ", "USD", 178.35, 1860000000000},
		{"META", "Meta Platforms Inc.", "Technology", "Internet Content & Information", "USA", "NASDAQ", "USD", 505.75, 1300000000000},
		{"NVDA", "NVIDIA Corporation", "Technology", "Semiconductors", "USA", "NASDAQ", "USD", 495.22, 1220000000000},
		{"TSLA", "Tesla Inc.", "Technology", "Auto Manufacturers", "USA", "NASDAQ", "USD", 238.45, 760000000000},

		// =====================================================================
		// USA - Financial Services / Banking
		// =====================================================================
		{"JPM", "JPMorgan Chase & Co.", "Financial Services", "Banks - Diversified", "USA", "NYSE", "USD", 172.45, 498000000000},
		{"BAC", "Bank of America Corp", "Financial Services", "Banks - Diversified", "USA", "NYSE", "USD", 33.88, 265000000000},
		{"WFC", "Wells Fargo & Company", "Financial Services", "Banks - Diversified", "USA", "NYSE", "USD", 46.12, 168000000000},
		{"GS", "Goldman Sachs Group", "Financial Services", "Capital Markets", "USA", "NYSE", "USD", 385.62, 128000000000},
		{"MS", "Morgan Stanley", "Financial Services", "Capital Markets", "USA", "NYSE", "USD", 87.95, 146000000000},

		// =====================================================================
		// USA - Energy
		// =====================================================================
		{"XOM", "Exxon Mobil Corporation", "Energy", "Oil & Gas Integrated", "USA", "NYSE", "USD", 104.55, 418000000000},
		{"CVX", "Chevron Corporation", "Energy", "Oil & Gas Integrated", "USA", "NYSE", "USD", 148.32, 276000000000},
		{"COP", "ConocoPhillips", "Energy", "Oil & Gas E&P", "USA", "NYSE", "USD", 115.78, 133000000000},

		// =====================================================================
		// USA - Healthcare
		// =====================================================================
		{"JNJ", "Johnson & Johnson", "Healthcare", "Drug Manufacturers", "USA", "NYSE", "USD", 156.42, 377000000000},
		{"UNH", "UnitedHealth Group", "Healthcare", "Healthcare Plans", "USA", "NYSE", "USD", 528.90, 489000000000},
		{"PFE", "Pfizer Inc.", "Healthcare", "Drug Manufacturers", "USA", "NYSE", "USD", 28.55, 161000000000},

		// =====================================================================
		// USA - Dividend Aristocrats (High Yield)
		// =====================================================================
		{"O", "Realty Income Corp", "Real Estate", "REIT - Retail", "USA", "NYSE", "USD", 55.62, 44000000000},
		{"T", "AT&T Inc.", "Communication Services", "Telecom Services", "USA", "NYSE", "USD", 17.45, 125000000000},
		{"VZ", "Verizon Communications", "Communication Services", "Telecom Services", "USA", "NYSE", "USD", 38.88, 163000000000},
		{"KO", "The Coca-Cola Company", "Consumer Defensive", "Beverages - Non-Alcoholic", "USA", "NYSE", "USD", 60.22, 261000000000},
		{"PG", "Procter & Gamble", "Consumer Defensive", "Household Products", "USA", "NYSE", "USD", 152.78, 361000000000},

		// =====================================================================
		// USA - Value Stocks
		// =====================================================================
		{"BRK.B", "Berkshire Hathaway B", "Financial Services", "Insurance - Diversified", "USA", "NYSE", "USD", 358.92, 789000000000},
		{"GM", "General Motors", "Consumer Cyclical", "Auto Manufacturers", "USA", "NYSE", "USD", 35.42, 48000000000},
		{"F", "Ford Motor Company", "Consumer Cyclical", "Auto Manufacturers", "USA", "NYSE", "USD", 12.15, 48000000000},
		{"INTC", "Intel Corporation", "Technology", "Semiconductors", "USA", "NASDAQ", "USD", 45.12, 191000000000},

		// =====================================================================
		// USA - High Risk / Growth
		// =====================================================================
		{"COIN", "Coinbase Global", "Financial Services", "Capital Markets", "USA", "NASDAQ", "USD", 142.55, 35000000000},
		{"RIVN", "Rivian Automotive", "Consumer Cyclical", "Auto Manufacturers", "USA", "NASDAQ", "USD", 18.22, 18000000000},
		{"PLTR", "Palantir Technologies", "Technology", "Software - Infrastructure", "USA", "NYSE", "USD", 22.35, 48000000000},
		{"SOFI", "SoFi Technologies", "Financial Services", "Credit Services", "USA", "NASDAQ", "USD", 8.45, 8500000000},
		{"HOOD", "Robinhood Markets", "Financial Services", "Capital Markets", "USA", "NASDAQ", "USD", 10.88, 9500000000},

		// =====================================================================
		// USA - Momentum (High Beta)
		// =====================================================================
		{"AMD", "Advanced Micro Devices", "Technology", "Semiconductors", "USA", "NASDAQ", "USD", 121.33, 196000000000},
		{"NFLX", "Netflix Inc.", "Communication Services", "Entertainment", "USA", "NASDAQ", "USD", 478.92, 212000000000},
		{"CRM", "Salesforce Inc.", "Technology", "Software - Application", "USA", "NYSE", "USD", 265.88, 258000000000},

		// =====================================================================
		// ISRAEL - Tel Aviv Stock Exchange
		// =====================================================================
		{"TEVA", "Teva Pharmaceutical", "Healthcare", "Drug Manufacturers - Generic", "Israel", "NYSE", "USD", 15.82, 17500000000},
		{"NICE", "NICE Ltd", "Technology", "Software - Application", "Israel", "NASDAQ", "USD", 185.45, 12200000000},
		{"CHKP", "Check Point Software", "Technology", "Software - Infrastructure", "Israel", "NASDAQ", "USD", 142.88, 16800000000},
		{"CYBR", "CyberArk Software", "Technology", "Software - Infrastructure", "Israel", "NASDAQ", "USD", 245.55, 10500000000},
		{"WIX", "Wix.com Ltd", "Technology", "Software - Application", "Israel", "NASDAQ", "USD", 142.35, 8200000000},
		{"MNDY", "monday.com Ltd", "Technology", "Software - Application", "Israel", "NASDAQ", "USD", 195.22, 9500000000},
		{"GLBE", "Global-e Online", "Technology", "Internet Retail", "Israel", "NASDAQ", "USD", 38.45, 6500000000},
		{"LPSN", "LivePerson Inc", "Technology", "Software - Application", "Israel", "NASDAQ", "USD", 3.22, 350000000},
		{"FVRR", "Fiverr International", "Technology", "Internet Content & Information", "Israel", "NYSE", "USD", 25.88, 920000000},
		{"LEUMI.TA", "Bank Leumi", "Financial Services", "Banks - Regional", "Israel", "TASE", "ILS", 34.55, 48000000000},
		{"HAPOALIM.TA", "Bank Hapoalim", "Financial Services", "Banks - Regional", "Israel", "TASE", "ILS", 38.22, 52000000000},
		{"ICL", "ICL Group", "Basic Materials", "Agricultural Inputs", "Israel", "NYSE", "USD", 5.45, 7000000000},

		// =====================================================================
		// UK - London Stock Exchange
		// =====================================================================
		{"HSBA.L", "HSBC Holdings", "Financial Services", "Banks - Diversified", "UK", "LSE", "GBP", 642.50, 128000000000},
		{"BP.L", "BP plc", "Energy", "Oil & Gas Integrated", "UK", "LSE", "GBP", 485.35, 88000000000},
		{"SHEL.L", "Shell plc", "Energy", "Oil & Gas Integrated", "UK", "LSE", "GBP", 2542.50, 185000000000},
		{"AZN.L", "AstraZeneca", "Healthcare", "Drug Manufacturers", "UK", "LSE", "GBP", 10285.00, 162000000000},
		{"GSK.L", "GSK plc", "Healthcare", "Drug Manufacturers", "UK", "LSE", "GBP", 1445.20, 60000000000},
		{"ULVR.L", "Unilever PLC", "Consumer Defensive", "Household Products", "UK", "LSE", "GBP", 4125.50, 105000000000},
		{"RIO.L", "Rio Tinto", "Basic Materials", "Other Industrial Metals", "UK", "LSE", "GBP", 5285.00, 85000000000},
		{"BARC.L", "Barclays PLC", "Financial Services", "Banks - Diversified", "UK", "LSE", "GBP", 175.82, 28000000000},
		{"LLOY.L", "Lloyds Banking Group", "Financial Services", "Banks - Regional", "UK", "LSE", "GBP", 52.45, 33000000000},
		{"VOD.L", "Vodafone Group", "Communication Services", "Telecom Services", "UK", "LSE", "GBP", 72.88, 20000000000},
		{"BT.A.L", "BT Group", "Communication Services", "Telecom Services", "UK", "LSE", "GBP", 128.55, 13000000000},
		{"TSCO.L", "Tesco PLC", "Consumer Defensive", "Grocery Stores", "UK", "LSE", "GBP", 285.40, 21000000000},
		{"DGE.L", "Diageo", "Consumer Defensive", "Beverages - Wineries & Distilleries", "UK", "LSE", "GBP", 2785.50, 62000000000},

		// =====================================================================
		// GERMANY - Frankfurt Stock Exchange
		// =====================================================================
		{"SAP.DE", "SAP SE", "Technology", "Software - Application", "Germany", "XETRA", "EUR", 145.88, 178000000000},
		{"SIE.DE", "Siemens AG", "Industrials", "Conglomerates", "Germany", "XETRA", "EUR", 168.42, 135000000000},
		{"ALV.DE", "Allianz SE", "Financial Services", "Insurance - Diversified", "Germany", "XETRA", "EUR", 245.55, 98000000000},
		{"BAS.DE", "BASF SE", "Basic Materials", "Chemicals", "Germany", "XETRA", "EUR", 45.22, 40000000000},
		{"BAYN.DE", "Bayer AG", "Healthcare", "Drug Manufacturers", "Germany", "XETRA", "EUR", 28.55, 28000000000},
		{"BMW.DE", "BMW AG", "Consumer Cyclical", "Auto Manufacturers", "Germany", "XETRA", "EUR", 98.45, 62000000000},
		{"MBG.DE", "Mercedes-Benz Group", "Consumer Cyclical", "Auto Manufacturers", "Germany", "XETRA", "EUR", 62.88, 67000000000},
		{"VOW3.DE", "Volkswagen AG", "Consumer Cyclical", "Auto Manufacturers", "Germany", "XETRA", "EUR", 108.55, 55000000000},
		{"DBK.DE", "Deutsche Bank", "Financial Services", "Banks - Diversified", "Germany", "XETRA", "EUR", 14.85, 29000000000},
		{"DTE.DE", "Deutsche Telekom", "Communication Services", "Telecom Services", "Germany", "XETRA", "EUR", 22.45, 112000000000},

		// =====================================================================
		// JAPAN - Tokyo Stock Exchange
		// =====================================================================
		{"7203.T", "Toyota Motor Corp", "Consumer Cyclical", "Auto Manufacturers", "Japan", "TSE", "JPY", 2845.00, 285000000000},
		{"6758.T", "Sony Group Corp", "Technology", "Consumer Electronics", "Japan", "TSE", "JPY", 12450.00, 155000000000},
		{"9984.T", "SoftBank Group", "Communication Services", "Telecom Services", "Japan", "TSE", "JPY", 6785.00, 98000000000},
		{"8306.T", "Mitsubishi UFJ Financial", "Financial Services", "Banks - Diversified", "Japan", "TSE", "JPY", 1285.50, 115000000000},
		{"9432.T", "Nippon Telegraph & Tel", "Communication Services", "Telecom Services", "Japan", "TSE", "JPY", 168.50, 152000000000},
		{"7267.T", "Honda Motor Co", "Consumer Cyclical", "Auto Manufacturers", "Japan", "TSE", "JPY", 1485.00, 78000000000},

		// =====================================================================
		// CHINA / HONG KONG
		// =====================================================================
		{"BABA", "Alibaba Group", "Consumer Cyclical", "Internet Retail", "China", "NYSE", "USD", 78.45, 195000000000},
		{"JD", "JD.com Inc", "Consumer Cyclical", "Internet Retail", "China", "NASDAQ", "USD", 28.55, 45000000000},
		{"PDD", "PDD Holdings", "Consumer Cyclical", "Internet Retail", "China", "NASDAQ", "USD", 132.88, 185000000000},
		{"BIDU", "Baidu Inc", "Communication Services", "Internet Content & Information", "China", "NASDAQ", "USD", 95.42, 33000000000},
		{"NIO", "NIO Inc", "Consumer Cyclical", "Auto Manufacturers", "China", "NYSE", "USD", 5.88, 11000000000},
		{"XPEV", "XPeng Inc", "Consumer Cyclical", "Auto Manufacturers", "China", "NYSE", "USD", 8.22, 7500000000},
		{"0700.HK", "Tencent Holdings", "Communication Services", "Internet Content & Information", "China", "HKEX", "HKD", 312.60, 360000000000},
		{"9988.HK", "Alibaba Group HK", "Consumer Cyclical", "Internet Retail", "China", "HKEX", "HKD", 78.85, 195000000000},

		// =====================================================================
		// INDIA
		// =====================================================================
		{"INFY", "Infosys Ltd", "Technology", "IT Services", "India", "NYSE", "USD", 18.22, 76000000000},
		{"WIT", "Wipro Ltd", "Technology", "IT Services", "India", "NYSE", "USD", 5.45, 28000000000},
		{"HDB", "HDFC Bank", "Financial Services", "Banks - Regional", "India", "NYSE", "USD", 58.88, 142000000000},
		{"IBN", "ICICI Bank", "Financial Services", "Banks - Regional", "India", "NYSE", "USD", 25.42, 88000000000},
		{"TTM", "Tata Motors", "Consumer Cyclical", "Auto Manufacturers", "India", "NYSE", "USD", 22.55, 85000000000},

		// =====================================================================
		// BRAZIL
		// =====================================================================
		{"VALE", "Vale S.A.", "Basic Materials", "Other Industrial Metals", "Brazil", "NYSE", "USD", 12.45, 52000000000},
		{"PBR", "Petrobras", "Energy", "Oil & Gas Integrated", "Brazil", "NYSE", "USD", 14.88, 98000000000},
		{"ITUB", "Itau Unibanco", "Financial Services", "Banks - Regional", "Brazil", "NYSE", "USD", 6.22, 60000000000},
		{"BBD", "Banco Bradesco", "Financial Services", "Banks - Regional", "Brazil", "NYSE", "USD", 2.88, 32000000000},
		{"NU", "Nu Holdings", "Financial Services", "Banks - Regional", "Brazil", "NYSE", "USD", 11.45, 55000000000},

		// =====================================================================
		// CANADA
		// =====================================================================
		{"TD", "Toronto-Dominion Bank", "Financial Services", "Banks - Diversified", "Canada", "NYSE", "USD", 58.22, 105000000000},
		{"RY", "Royal Bank of Canada", "Financial Services", "Banks - Diversified", "Canada", "NYSE", "USD", 98.45, 138000000000},
		{"ENB", "Enbridge Inc", "Energy", "Oil & Gas Midstream", "Canada", "NYSE", "USD", 35.88, 78000000000},
		{"CNQ", "Canadian Natural Resources", "Energy", "Oil & Gas E&P", "Canada", "NYSE", "USD", 32.42, 68000000000},
		{"SHOP", "Shopify Inc", "Technology", "Software - Application", "Canada", "NYSE", "USD", 68.55, 88000000000},

		// =====================================================================
		// AUSTRALIA
		// =====================================================================
		{"BHP", "BHP Group", "Basic Materials", "Other Industrial Metals", "Australia", "NYSE", "USD", 58.22, 145000000000},
		{"RIO", "Rio Tinto Group", "Basic Materials", "Other Industrial Metals", "Australia", "NYSE", "USD", 65.88, 105000000000},

		// =====================================================================
		// SWITZERLAND
		// =====================================================================
		{"NESN.SW", "Nestle S.A.", "Consumer Defensive", "Packaged Foods", "Switzerland", "SIX", "CHF", 98.42, 278000000000},
		{"NOVN.SW", "Novartis AG", "Healthcare", "Drug Manufacturers", "Switzerland", "SIX", "CHF", 92.55, 198000000000},
		{"ROG.SW", "Roche Holding", "Healthcare", "Drug Manufacturers", "Switzerland", "SIX", "CHF", 248.88, 212000000000},
		{"UBS", "UBS Group AG", "Financial Services", "Banks - Diversified", "Switzerland", "NYSE", "USD", 28.45, 92000000000},
		{"CS", "Credit Suisse Group", "Financial Services", "Capital Markets", "Switzerland", "NYSE", "USD", 0.82, 2800000000},

		// =====================================================================
		// FRANCE
		// =====================================================================
		{"MC.PA", "LVMH", "Consumer Cyclical", "Luxury Goods", "France", "EURONEXT", "EUR", 745.50, 375000000000},
		{"OR.PA", "L'Oreal S.A.", "Consumer Defensive", "Household Products", "France", "EURONEXT", "EUR", 425.88, 228000000000},
		{"TTE.PA", "TotalEnergies", "Energy", "Oil & Gas Integrated", "France", "EURONEXT", "EUR", 58.22, 142000000000},
		{"SAN.PA", "Sanofi S.A.", "Healthcare", "Drug Manufacturers", "France", "EURONEXT", "EUR", 92.45, 118000000000},
		{"BNP.PA", "BNP Paribas", "Financial Services", "Banks - Diversified", "France", "EURONEXT", "EUR", 58.88, 72000000000},

		// =====================================================================
		// NETHERLANDS
		// =====================================================================
		{"ASML", "ASML Holding", "Technology", "Semiconductor Equipment", "Netherlands", "NASDAQ", "USD", 685.42, 278000000000},

		// =====================================================================
		// SOUTH KOREA
		// =====================================================================
		{"005930.KS", "Samsung Electronics", "Technology", "Consumer Electronics", "South Korea", "KRX", "KRW", 71500.00, 425000000000},
		{"000660.KS", "SK Hynix", "Technology", "Semiconductors", "South Korea", "KRX", "KRW", 135000.00, 98000000000},

		// =====================================================================
		// TAIWAN
		// =====================================================================
		{"TSM", "Taiwan Semiconductor", "Technology", "Semiconductors", "Taiwan", "NYSE", "USD", 108.55, 562000000000},
	}

	m.stocks = make([]models.Stock, 0, len(mockData))

	for _, data := range mockData {
		stock := m.generateStock(data.symbol, data.name, data.sector, data.industry, data.country, data.exchange, data.currency, data.priceBase, data.marketCapBase)
		m.stocks = append(m.stocks, stock)
	}
}

func (m *MockDataService) generateStock(symbol, name, sector, industry, country, exchange, currency string, priceBase float64, marketCapBase int64) models.Stock {
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
		Exchange:          exchange,
		Currency:          currency,
		Country:           country,
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
