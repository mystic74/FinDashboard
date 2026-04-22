package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"time"

	ampyfin "github.com/AmpyFin/yfinance-go"
	ffeng "github.com/FFengIll/yfinance-go"
	yfclient "github.com/wnjoon/go-yfinance/pkg/client"
	yfq "github.com/wnjoon/go-yfinance/pkg/models"
	yfticker "github.com/wnjoon/go-yfinance/pkg/ticker"
	"stock-screener/models"

	"golang.org/x/sync/errgroup"
)

// Yahoo quote drivers (see config.YahooQuoteDriver).
const (
	YahooQuoteDriverResty   = "resty"
	YahooQuoteDriverFFeng   = "ffeng"
	YahooQuoteDriverAmpyFin = "ampyfin"
	// YahooQuoteDriverWnjoon uses github.com/wnjoon/go-yfinance (CycleTLS + crumb auth; avoids naive HTTP 401s).
	YahooQuoteDriverWnjoon = "wnjoon"
)

const wnjoonQuoteBatchSize = 50

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

func (y *YahooFinanceService) getQuotesWnjoon(ctx context.Context, symbols []string) ([]models.Stock, error) {
	if len(symbols) == 0 {
		return []models.Stock{}, nil
	}
	normalized := make([]string, 0, len(symbols))
	for _, sym := range symbols {
		s := strings.TrimSpace(sym)
		if s == "" {
			continue
		}
		normalized = append(normalized, strings.ToUpper(s))
	}
	if len(normalized) == 0 {
		return []models.Stock{}, nil
	}

	var all []models.Stock
	for i := 0; i < len(normalized); i += wnjoonQuoteBatchSize {
		end := i + wnjoonQuoteBatchSize
		if end > len(normalized) {
			end = len(normalized)
		}
		batch := normalized[i:end]
		stocks, err := y.getQuotesWnjoonBatch(ctx, batch)
		if err != nil {
			return nil, err
		}
		all = append(all, stocks...)
	}
	return all, nil
}

func wnjoonQuoteErrSkippable(err error) bool {
	if err == nil {
		return false
	}
	return yfclient.IsNotFoundError(err) || yfclient.IsNoDataError(err) || yfclient.IsInvalidSymbolError(err)
}

func (y *YahooFinanceService) getQuotesWnjoonBatch(ctx context.Context, batch []string) ([]models.Stock, error) {
	// One Ticker (and CycleTLS client) per symbol, bounded by maxConcurrent — safe in parallel and
	// much faster than sharing one client (sequential) across hundreds of S&P / index names.
	var out []models.Stock
	var mu sync.Mutex
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(y.maxConcurrent)
	for _, sym := range batch {
		sym := sym
		g.Go(func() error {
			if err := ctx.Err(); err != nil {
				return err
			}
			tkr, err := yfticker.New(sym)
			if err != nil {
				if wnjoonQuoteErrSkippable(err) {
					log.Printf("[Yahoo/wnjoon] skip %s: %v", sym, err)
					return nil
				}
				return fmt.Errorf("wnjoon ticker %s: %w", sym, err)
			}
			defer tkr.Close()
			q, err := tkr.Quote()
			if err != nil {
				if wnjoonQuoteErrSkippable(err) {
					log.Printf("[Yahoo/wnjoon] skip %s: %v", sym, err)
					return nil
				}
				return fmt.Errorf("wnjoon quote %s: %w", sym, err)
			}
			st := convertWnjoonQuoteToStock(q)
			mu.Lock()
			out = append(out, st)
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return out, nil
}

func convertWnjoonQuoteToStock(q *yfq.Quote) models.Stock {
	if q == nil {
		return models.Stock{LastUpdated: time.Now()}
	}
	name := q.LongName
	if name == "" {
		name = q.ShortName
	}
	dividendYield := q.TrailingAnnualDividendYield * 100
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
		AvgVolume:         q.AverageDailyVolume3Month,
		AvgVolume10Day:    q.AverageDailyVolume10Day,
		MarketCap:         q.MarketCap,
		SharesOutstanding: q.SharesOutstanding,
		Week52High:        q.FiftyTwoWeekHigh,
		Week52Low:         q.FiftyTwoWeekLow,
		MA50:              q.FiftyDayAverage,
		MA200:             q.TwoHundredDayAverage,
		PERatio:           q.TrailingPE,
		ForwardPE:         q.ForwardPE,
		PBRatio:           q.PriceToBook,
		DividendYield:     dividendYield,
		DividendPerShare:  q.TrailingAnnualDividendRate,
		EPS:               q.EpsTrailingTwelveMonths,
		BookValuePerShare: q.BookValue,
		LastUpdated:       time.Now(),
	}
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
				return fmt.Errorf("ampyfin fetch %s: %w", sym, err)
			}
			raw, err := json.Marshal(nq)
			if err != nil {
				return fmt.Errorf("ampyfin marshal %s: %w", sym, err)
			}
			var wire ampyQuoteJSON
			if err := json.Unmarshal(raw, &wire); err != nil {
				return fmt.Errorf("ampyfin unmarshal %s: %w", sym, err)
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
	expected := 0
	for _, sym := range symbols {
		if strings.TrimSpace(sym) != "" {
			expected++
		}
	}
	if len(out) != expected {
		return nil, fmt.Errorf("ampyfin: returned %d/%d quotes", len(out), expected)
	}
	return out, nil
}

func derefInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}
