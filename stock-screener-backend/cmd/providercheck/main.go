// Command providercheck verifies live data providers respond and (by default) at least two succeed.
//
// Usage:
//
//	cd stock-screener-backend && go run ./cmd/providercheck
//	MIN_WORKING_PROVIDERS=2 go run ./cmd/providercheck
//
// Environment: same as the API (FMP_API_KEY, ALPHA_VANTAGE_KEY, YAHOO_QUOTE_DRIVER, DEMO_MODE, etc.)
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"stock-screener/config"
	"stock-screener/services"
)

func main() {
	minFlag := flag.Int("min", 0, "minimum working providers required (0 = use MIN_WORKING_PROVIDERS env or default 2)")
	strictDemo := flag.Bool("strict-demo", false, "if set, exit non-zero when DEMO_MODE=true")
	flag.Parse()

	cfg := config.DefaultConfig()
	cacheCleanup := cfg.CacheTTL * 2
	if cacheCleanup < time.Minute {
		cacheCleanup = 10 * time.Minute
	}
	cache := services.NewCacheService(cfg.CacheTTL, cacheCleanup)

	min := *minFlag
	if min <= 0 {
		if s := os.Getenv("MIN_WORKING_PROVIDERS"); s != "" {
			if v, err := strconv.Atoi(s); err == nil && v > 0 {
				min = v
			}
		}
	}
	if min <= 0 {
		min = 2
	}

	if cfg.DemoMode {
		fmt.Println("DEMO_MODE=true: live provider probes skipped (mock data mode).")
		if *strictDemo {
			os.Exit(1)
		}
		os.Exit(0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	results := services.ProbeLiveDataProviders(ctx, cfg, cache)

	fmt.Println("Provider probe results:")
	for _, r := range results {
		status := "OK"
		if !r.OK {
			status = "FAIL"
		}
		msg := r.Message
		if len(msg) > 120 {
			msg = msg[:117] + "..."
		}
		fmt.Printf("  %-15s %s", r.Name, status)
		if msg != "" {
			fmt.Printf(" — %s", msg)
		}
		fmt.Println()
	}

	slots := services.ConfiguredProviderSlots(cfg)
	working := services.CountWorkingProviders(results)
	fmt.Printf("Configured channels: %d | Working now: %d | Required minimum: %d\n", slots, working, min)

	if err := services.RedundancyCheckError(cfg, results, min); err != nil {
		fmt.Fprintf(os.Stderr, "redundancy check failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("redundancy check passed.")
	os.Exit(0)
}
