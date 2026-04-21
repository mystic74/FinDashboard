package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	ampyfin "github.com/AmpyFin/yfinance-go"
	ffeng "github.com/FFengIll/yfinance-go"
	"stock-screener/models"

	"golang.org/x/sync/errgroup"
)

// Yahoo quote drivers (see config.YahooQuoteDriver).
const (
	YahooQuoteDriverResty   = "resty"
	YahooQuoteDriverFFeng   = "ffeng"
	YahooQuoteDriverAmpyFin = "ampyfin"
)

func (y *YahooFinanceService) getQuotesResty(_ context.Context, symbols []string) ([]models.Stock, error) {
	batchSize := 100
	var allStocks []models.Stock
	for i := 0; i < len(symbols); i += batchSize {
		end := i + batchSize
		if end > len(symbols) {
			end = len(symbols)
		}
		batch := symbols[i:end]
		stocks, err := y.fetchQuoteBatch(batch)
		if err != nil {
			return nil, err
		}
		allStocks = append(allStocks, stocks...)
	}
	return allStocks, nil
}

func (y *YahooFinanceService) getQuotesFFeng(ctx context.Context, symbols []string) ([]models.Stock, error) {
	if len(symbols) == 0 {
		return []models.Stock{}, nil
	}
	quotes, err := ffeng.GetQuotes(ctx, symbols)
	if err != nil {
		return nil, err
	}
	out := make([]models.Stock, 0, len(quotes))
	for _, q := range quotes {
		if q == nil {
			continue
		}
		out = append(out, convertFFengQuoteToStock(q))
	}
	return out, nil
}

func convertFFengQuoteToStock(q *ffeng.Quote) models.Stock {
	name := q.LongName
	if name == "" {
		name = q.ShortName
	}
	dividendYield := q.DividendYield * 100
	return models.Stock{
		Symbol:            q.Symbol,
		Name:              name,
		Exchange:          q.Exchange,
		Currency:          q.Currency,
		Price:             q.RegularMarketPrice,
		Open:              q.RegularMarketOpen,
		High:              q.RegularMarketDayHigh,
		Low:               q.RegularMarketDayLow,
		PreviousClose:     q.RegularMarketPreviousClose,
		Change:            q.RegularMarketChange,
		ChangePercent:     q.RegularMarketChangePercent,
		Volume:            q.RegularMarketVolume,
		MarketCap:         q.MarketCap,
		SharesOutstanding: q.SharesOutstanding,
		Week52High:        q.FiftyTwoWeekHigh,
		Week52Low:         q.FiftyTwoWeekLow,
		MA50:              q.FiftyDayAverage,
		MA200:             q.TwoHundredDayAverage,
		PERatio:           q.PE,
		ForwardPE:         q.ForwardPE,
		PBRatio:           q.PriceToBook,
		DividendYield:     dividendYield,
		EPS:               q.EPS,
		BookValuePerShare: q.BookValue,
		Beta:              q.Beta,
		LastUpdated:       time.Now(),
	}
}

// ampyQuoteJSON matches json from AmpyFin NormalizedQuote (internal/norm) via Marshal.
type ampyScaledMoney struct {
	Scaled int64 `json:"scaled"`
	Scale  int   `json:"scale"`
}

type ampyQuoteJSON struct {
	Security struct {
		Symbol string `json:"symbol"`
	} `json:"security"`
	RegularMarketPrice  *ampyScaledMoney `json:"regular_market_price"`
	RegularMarketHigh   *ampyScaledMoney `json:"regular_market_high"`
	RegularMarketLow    *ampyScaledMoney `json:"regular_market_low"`
	RegularMarketVolume *int64           `json:"regular_market_volume"`
	CurrencyCode        string           `json:"currency_code"`
}

func scaledToFloat(s *ampyScaledMoney) float64 {
	if s == nil {
		return 0
	}
	denom := math.Pow10(s.Scale)
	if denom == 0 {
		return float64(s.Scaled)
	}
	return float64(s.Scaled) / denom
}

func (y *YahooFinanceService) getQuotesAmpyFin(ctx context.Context, symbols []string) ([]models.Stock, error) {
	if len(symbols) == 0 {
		return []models.Stock{}, nil
	}
	client := ampyfin.NewClientWithSessionRotation()
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(y.maxConcurrent)
	var mu sync.Mutex
	out := make([]models.Stock, 0, len(symbols))

	for _, sym := range symbols {
		sym := strings.TrimSpace(sym)
		if sym == "" {
			continue
		}
		g.Go(func() error {
			nq, err := client.FetchQuote(ctx, sym, "stock-screener")
			if err != nil {
				return nil
			}
			raw, err := json.Marshal(nq)
			if err != nil {
				return nil
			}
			var wire ampyQuoteJSON
			if err := json.Unmarshal(raw, &wire); err != nil {
				return nil
			}
			st := models.Stock{
				Symbol:      wire.Security.Symbol,
				Currency:    wire.CurrencyCode,
				Price:       scaledToFloat(wire.RegularMarketPrice),
				High:        scaledToFloat(wire.RegularMarketHigh),
				Low:         scaledToFloat(wire.RegularMarketLow),
				Volume:      derefInt64(wire.RegularMarketVolume),
				LastUpdated: time.Now(),
			}
			mu.Lock()
			out = append(out, st)
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("ampyfin: no quotes returned")
	}
	return out, nil
}

func derefInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}
