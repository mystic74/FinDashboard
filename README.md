# Stock Screener Dashboard

A professional stock screener dashboard with predefined screeners and custom filtering capabilities. Built with Go (Gin) backend and React + TypeScript frontend.

## Features

### Predefined Screeners
- **Momentum Masters** - Stocks with strong price momentum and high volume
- **Dividend Aristocrats** - Companies with long track records of dividend payments
- **Value Opportunities** - Undervalued stocks with strong fundamentals
- **High Beta Bulls** - High-volatility stocks with strong recent performance
- **Cash is King** - Companies with exceptional liquidity
- **Piotroski F-Score Leaders** - Stocks with high F-Score (8-9)
- **Small Cap Growth** - High-growth small cap stocks
- **Undervalued Tech** - Technology stocks trading below industry average
- **GARP (Growth at Reasonable Price)** - Balanced growth and value
- **Quality Stocks** - High-quality companies with strong profitability
- **Low Volatility** - Stable stocks with low beta
- **Turnaround Candidates** - Beaten-down stocks with improving fundamentals

### Filter Categories
- **Price & Volume** - Price, change, volume, market cap
- **Valuation** - P/E, P/B, P/S, PEG, EV/EBITDA
- **Dividends** - Yield, payout ratio, consecutive years
- **Financial Health** - Current ratio, quick ratio, debt/equity, Altman Z-Score
- **Profitability** - ROE, ROA, ROIC, margins
- **Growth** - Revenue growth, EPS growth, FCF growth
- **Technical** - RSI, moving averages, 52-week range, beta
- **Profile** - Sector, industry, country, exchange

### Technical Indicators
- Piotroski F-Score (0-9)
- Altman Z-Score
- RSI (14-day)
- MACD
- 50-day and 200-day Moving Averages

## Tech Stack

### Backend (Go)
- **Framework**: Gin
- **HTTP Client**: Resty
- **Caching**: go-cache (in-memory)
- **Testing**: testify

### Frontend (React + TypeScript)
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **State Management**: React Query, Zustand
- **UI Components**: Lucide React icons
- **Charts**: Recharts (planned)

## Project Structure

```
FinDashboard/
├── stock-screener-backend/
│   ├── api/
│   │   ├── handlers/    # HTTP request handlers
│   │   └── routes/      # Route definitions
│   ├── models/          # Data models
│   ├── services/        # Business logic
│   ├── utils/           # Calculations, validators
│   └── config/          # Configuration
└── stock-screener-frontend/
    └── src/
        ├── components/  # React components
        ├── hooks/       # Custom hooks
        ├── services/    # API services
        ├── types/       # TypeScript types
        ├── utils/       # Utility functions
        └── context/     # React context
```

## Getting Started

### Backend

```bash
cd stock-screener-backend
go mod tidy
go run main.go
```

The API will be available at `http://localhost:8080`

### Frontend

```bash
cd stock-screener-frontend
npm install
npm run dev
```

The frontend will be available at `http://localhost:3000`

## API Endpoints

### Screeners
- `GET /api/v1/screeners` - Get all predefined screeners
- `GET /api/v1/screeners/:name` - Run a specific screener
- `POST /api/v1/screeners/custom` - Run custom screener with filters

### Stocks
- `GET /api/v1/stocks?symbols=AAPL,MSFT` - Get multiple stocks
- `GET /api/v1/stocks/:symbol` - Get stock data
- `GET /api/v1/stocks/:symbol/fundamentals` - Get fundamental data
- `GET /api/v1/stocks/:symbol/history` - Get historical prices

### Filters
- `GET /api/v1/filters` - Get all available filters
- `GET /api/v1/filters/categories` - Get filter categories

### Sectors
- `GET /api/v1/sectors` - Get sector performance
- `GET /api/v1/sectors/list` - Get sector list

## Running Tests

```bash
cd stock-screener-backend
go test ./... -v
```

Test coverage includes:
- Piotroski F-Score calculations
- Altman Z-Score calculations
- RSI and moving average calculations
- Financial ratio calculations
- Data validation and sanity checks
- Screener filter logic

## Data Sources

- Yahoo Finance API (primary)
- Financial Modeling Prep (alternative)
- Alpha Vantage (alternative)

## Environment Variables

```env
SERVER_PORT=8080
FMP_API_KEY=your_api_key
ALPHA_VANTAGE_KEY=your_api_key
```

## License

MIT
