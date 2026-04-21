package services

import (
	"context"
	"fmt"
	"strings"

	"stock-screener/config"
)

// ProviderProbeName identifies a data source used for redundancy checks.
type ProviderProbeName string

const (
	ProbeFMP          ProviderProbeName = "fmp"
	ProbeAlphaVantage ProviderProbeName = "alpha_vantage"
	ProbeYahoo        ProviderProbeName = "yahoo"
	ProbeDemo         ProviderProbeName = "demo"
)

// ProviderProbeResult is one provider health outcome from a live probe.
type ProviderProbeResult struct {
	Name    ProviderProbeName `json:"name"`
	OK      bool              `json:"ok"`
	Message string            `json:"message,omitempty"`
}

// ProbeLiveDataProviders performs minimal live calls (typically AAPL) for each configured channel.
// Yahoo is always probed when not in demo mode. FMP and Alpha Vantage run only when API keys are set.
func ProbeLiveDataProviders(ctx context.Context, cfg *config.Config, cache *CacheService) []ProviderProbeResult {
	if cfg.DemoMode {
		return []ProviderProbeResult{
			{Name: ProbeDemo, OK: true, Message: "demo mode — live providers not probed"},
		}
	}

	var out []ProviderProbeResult
	testSym := []string{"AAPL"}

	if strings.TrimSpace(cfg.FMPAPIKey) != "" {
		fmp := NewFMPService(cfg.FMPAPIKey, cache)
		_, err := fmp.GetQuotes(testSym)
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		out = append(out, ProviderProbeResult{Name: ProbeFMP, OK: err == nil, Message: msg})
	}

	if strings.TrimSpace(cfg.AlphaVantageKey) != "" {
		av := NewAlphaVantageService(cfg.AlphaVantageKey, cache)
		_, err := av.GetQuote("AAPL")
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		out = append(out, ProviderProbeResult{Name: ProbeAlphaVantage, OK: err == nil, Message: msg})
	}

	yahoo := NewYahooFinanceServiceWithDriver(cache, cfg.YahooQuoteDriver)
	_, err := yahoo.GetQuotes(ctx, testSym)
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	out = append(out, ProviderProbeResult{Name: ProbeYahoo, OK: err == nil, Message: msg})

	return out
}

// ConfiguredProviderSlots returns how many independent provider channels are available (not whether they work).
func ConfiguredProviderSlots(cfg *config.Config) int {
	if cfg.DemoMode {
		return 0
	}
	n := 1 // Yahoo path always available in non-demo
	if strings.TrimSpace(cfg.FMPAPIKey) != "" {
		n++
	}
	if strings.TrimSpace(cfg.AlphaVantageKey) != "" {
		n++
	}
	return n
}

// CountWorkingProviders returns how many probes succeeded.
func CountWorkingProviders(results []ProviderProbeResult) int {
	n := 0
	for _, r := range results {
		if r.Name == ProbeDemo {
			continue
		}
		if r.OK {
			n++
		}
	}
	return n
}

// RedundancyCheckError describes why minimum redundancy was not met.
func RedundancyCheckError(cfg *config.Config, results []ProviderProbeResult, minRequired int) error {
	slots := ConfiguredProviderSlots(cfg)
	if slots < minRequired {
		return fmt.Errorf("only %d provider channel(s) configured (need at least %d: set FMP_API_KEY and/or ALPHA_VANTAGE_KEY for redundancy alongside Yahoo)", slots, minRequired)
	}
	w := CountWorkingProviders(results)
	if w < minRequired {
		return fmt.Errorf("only %d/%d provider(s) responded OK (need >= %d)", w, len(filterRealProbes(results)), minRequired)
	}
	return nil
}

func filterRealProbes(results []ProviderProbeResult) []ProviderProbeResult {
	var out []ProviderProbeResult
	for _, r := range results {
		if r.Name == ProbeDemo {
			continue
		}
		out = append(out, r)
	}
	return out
}
