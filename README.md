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

### Quick Start with Make (Recommended)

```bash
# Install dependencies
make deps

# Run in development mode (no Docker)
make dev

# Or run with Docker
make prod
```

See all available commands:
```bash
make help
```

### Development Modes

#### Option 1: Local Development (Easier Debugging)

```bash
# Terminal 1 - Backend
cd stock-screener-backend
go run .

# Terminal 2 - Frontend
cd stock-screener-frontend
npm install
npm run dev
```

Or simply:
```bash
make dev
```

#### Option 2: Docker Development (Hot Reload)

```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

Or:
```bash
make dev-docker
```

#### Option 3: Docker Production

```bash
docker-compose up --build -d
```

Or:
```bash
make prod
```

### URLs

| Mode | Frontend | Backend API |
|------|----------|-------------|
| Local Dev | http://localhost:3000 | http://localhost:8080 |
| Docker Dev | http://localhost:3000 | http://localhost:8080 (via Vite `/api` proxy) |
| Docker Prod | http://localhost:3000 | http://localhost:8080 |

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
- `GET /api/v1/screeners/summary` - Summaries with match counts
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

### Market profiles
- `GET /api/v1/profiles` - List profiles (defaults + in-memory overrides)
- `GET /api/v1/profiles/:country` - Get one profile
- `PUT /api/v1/profiles/:country` - Update multipliers (non-USA; affects screener adjustment)
- `POST /api/v1/profiles/:country/reset` - Reset one country
- `POST /api/v1/profiles/reset` - Reset all overrides

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

With **`DEMO_MODE=false`**, the HTTP API uses **Yahoo Finance** for screener and stock data. Additional provider code (FMP, Alpha Vantage) exists under `services/` but is **not** wired into the live routes yet.

## Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

### Backend Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `GIN_MODE` | `release` | Gin mode (`debug` or `release`) |
| `DEMO_MODE` | `false` | Use mock data (`true`) or live Yahoo Finance (`false`) |
| `YAHOO_QUOTE_DRIVER` | `ffeng` | Yahoo quote path: `resty` (raw HTTP), `ffeng` (FFengIll yfinance-go), `ampyfin` (AmpyFin yfinance-go) |
| `CORS_ORIGIN` | `http://localhost:3000,http://localhost:5173` | Allowed CORS origins (comma-separated) |
| `CACHE_TTL` | `5m` | In-memory cache TTL (Go duration string) |

### Frontend Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `VITE_API_URL` | _(unset)_ | Optional full backend origin for direct browser calls (`.../api/v1` is appended). Default: relative `/api/v1` + Vite proxy (dev) or nginx (prod). |
| `VITE_DEV_PROXY_TARGET` | _(see compose)_ | Docker dev: Vite proxy target for `/api` (e.g. `http://backend:8080`). |

### Demo Mode

By default, the application runs in **Live Mode** with Yahoo Finance data.

Use **Demo Mode** (`DEMO_MODE=true`) with mock stock data when needed. This is useful for:
- Development without API rate limits
- Testing UI without network dependencies
- Running in environments where Yahoo Finance is blocked

To use live data from Yahoo Finance:
```bash
DEMO_MODE=false go run .
```

> **Note**: Yahoo Finance may block requests in some environments. Demo mode is recommended for development.

## License

MIT
