package services

import (
	"testing"

	"stock-screener/config"
)

func TestConfiguredProviderSlots(t *testing.T) {
	t.Parallel()
	if n := ConfiguredProviderSlots(&config.Config{DemoMode: true}); n != 0 {
		t.Fatalf("demo: got %d want 0", n)
	}
	if n := ConfiguredProviderSlots(&config.Config{DemoMode: false}); n != 1 {
		t.Fatalf("yahoo only: got %d want 1", n)
	}
	if n := ConfiguredProviderSlots(&config.Config{DemoMode: false, FMPAPIKey: "k"}); n != 2 {
		t.Fatalf("yahoo+fmp: got %d want 2", n)
	}
	if n := ConfiguredProviderSlots(&config.Config{DemoMode: false, AlphaVantageKey: "k"}); n != 2 {
		t.Fatalf("yahoo+av: got %d want 2", n)
	}
	if n := ConfiguredProviderSlots(&config.Config{DemoMode: false, FMPAPIKey: "a", AlphaVantageKey: "b"}); n != 3 {
		t.Fatalf("all three: got %d want 3", n)
	}
}

func TestCountWorkingProviders_SkipsDemo(t *testing.T) {
	t.Parallel()
	n := CountWorkingProviders([]ProviderProbeResult{
		{Name: ProbeDemo, OK: true},
		{Name: ProbeYahoo, OK: true},
	})
	if n != 1 {
		t.Fatalf("got %d want 1", n)
	}
}

func TestRedundancyCheckError(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{DemoMode: false}
	results := []ProviderProbeResult{
		{Name: ProbeFMP, OK: true},
		{Name: ProbeYahoo, OK: true},
	}
	if err := RedundancyCheckError(cfg, results, 2); err == nil {
		t.Fatal("expected error: only one channel configured")
	}

	cfg2 := &config.Config{DemoMode: false, FMPAPIKey: "x", AlphaVantageKey: "y"}
	if err := RedundancyCheckError(cfg2, results, 2); err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	cfg3 := &config.Config{DemoMode: false, FMPAPIKey: "x"}
	bad := []ProviderProbeResult{
		{Name: ProbeFMP, OK: false},
		{Name: ProbeYahoo, OK: true},
	}
	if err := RedundancyCheckError(cfg3, bad, 2); err == nil {
		t.Fatal("expected error: only one working")
	}
}
