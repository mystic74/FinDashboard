package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"stock-screener/api/routes"
	"stock-screener/config"
	"syscall"
)

func main() {
	// Load configuration
	cfg := config.DefaultConfig()

	// Setup router
	router := routes.SetupRouter(cfg)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting Stock Screener API server on %s", serverAddr)
	log.Printf("API available at http://localhost%s/api/v1", serverAddr)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
}
