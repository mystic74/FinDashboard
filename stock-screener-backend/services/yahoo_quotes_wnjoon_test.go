package services

import (
	"errors"
	"fmt"
	"testing"
	"time"

	yfclient "github.com/wnjoon/go-yfinance/pkg/client"
	yfq "github.com/wnjoon/go-yfinance/pkg/models"
)

func TestConvertWnjoonQuoteToStock(t *testing.T) {
	t.Parallel()
	q := &yfq.Quote{
		Symbol:                      "AAPL",
		ShortName:                   "Apple Inc.",
		LongName:                    "Apple Inc.",
		Exchange:                    "NMS",
		Currency:                    "USD",
		RegularMarketPrice:          100,
		RegularMarketOpen:           99,
		RegularMarketDayHigh:        101,
		RegularMarketDayLow:         98,
		RegularMarketPreviousClose:  97,
		RegularMarketChange:         3,
		RegularMarketChangePercent:  3.09,
		RegularMarketVolume:         1_000_000,
		AverageDailyVolume3Month:    900_000,
		AverageDailyVolume10Day:     950_000,
		MarketCap:                   2_000_000_000_000,
		SharesOutstanding:           15_000_000_000,
		FiftyTwoWeekHigh:            120,
		FiftyTwoWeekLow:             80,
		FiftyDayAverage:             95,
		TwoHundredDayAverage:        90,
		TrailingPE:                  28.5,
		ForwardPE:                   26,
		PriceToBook:                 40,
		TrailingAnnualDividendRate:  0.96,
		TrailingAnnualDividendYield: 0.005,
		EpsTrailingTwelveMonths:     3.5,
		BookValue:                   2.5,
	}
	s := convertWnjoonQuoteToStock(q)
	if s.Symbol != "AAPL" || s.Price != 100 || s.PERatio != 28.5 {
		t.Fatalf("unexpected mapped stock: %+v", s)
	}
	if s.DividendYield != 0.5 {
		t.Fatalf("dividend yield: got %v want 0.5 (yield * 100)", s.DividendYield)
	}
	if s.Name != "Apple Inc." {
		t.Fatalf("name: got %q", s.Name)
	}
	if time.Since(s.LastUpdated) > time.Minute {
		t.Fatalf("LastUpdated not recent")
	}
}

func TestWnjoonQuoteErrSkippable(t *testing.T) {
	t.Parallel()
	if !wnjoonQuoteErrSkippable(yfclient.WrapNotFoundError("MTRX.TA")) {
		t.Fatal("expected not found to be skippable")
	}
	if !wnjoonQuoteErrSkippable(fmt.Errorf("wrap: %w", yfclient.WrapNotFoundError("X"))) {
		t.Fatal("expected wrapped not found to be skippable")
	}
	if wnjoonQuoteErrSkippable(nil) || wnjoonQuoteErrSkippable(errors.New("other")) {
		t.Fatal("expected nil and generic errors not skippable")
	}
}

func TestConvertWnjoonQuoteToStock_Nil(t *testing.T) {
	t.Parallel()
	s := convertWnjoonQuoteToStock(nil)
	if s.Symbol != "" || s.Price != 0 {
		t.Fatalf("expected zero stock, got %+v", s)
	}
}
