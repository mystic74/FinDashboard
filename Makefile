.PHONY: all help setup dev dev-docker prod backend frontend test clean

# Default target
all: help

help:
	@echo "Stock Screener Dashboard - Available Commands"
	@echo ""
	@echo "Setup:"
	@echo "  make setup        - First-time setup (copy .env and install deps)"
	@echo "  make deps         - Install all dependencies"
	@echo ""
	@echo "Development (No Docker - easier debugging):"
	@echo "  make dev          - Run both backend and frontend in dev mode"
	@echo "  make backend      - Run only backend (Go)"
	@echo "  make frontend     - Run only frontend (Vite)"
	@echo ""
	@echo "Development (Docker):"
	@echo "  make dev-docker   - Run with Docker (hot reload)"
	@echo ""
	@echo "Production (Docker):"
	@echo "  make prod         - Build and run production containers"
	@echo "  make prod-build   - Build production images only"
	@echo "  make prod-down    - Stop production containers"
	@echo "  make prod-logs    - View production logs"
	@echo ""
	@echo "Testing:"
	@echo "  make test         - Run all tests"
	@echo "  make test-verbose - Run tests with verbose output"
	@echo "  make test-coverage - Generate test coverage report"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean        - Clean build artifacts and containers"
	@echo ""
	@echo "Environment:"
	@echo "  Configure via .env file (see .env.example)"
	@echo "  DEMO_MODE=false (default) uses live Yahoo Finance API"
	@echo "  DEMO_MODE=true uses mock data"

# =============================================================================
# Setup
# =============================================================================

setup:
	@echo "Setting up Stock Screener Dashboard..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env from .env.example"; \
	else \
		echo ".env already exists, skipping..."; \
	fi
	@make deps
	@echo ""
	@echo "Setup complete! Run 'make dev' to start development servers."

# =============================================================================
# Development without Docker (recommended for debugging)
# =============================================================================

dev:
	@echo "Starting development servers..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@make -j2 backend frontend

backend:
	@echo "Starting Go backend..."
	cd stock-screener-backend && go run .

frontend:
	@echo "Starting Vite frontend..."
	cd stock-screener-frontend && npm run dev

# =============================================================================
# Development with Docker
# =============================================================================

dev-docker:
	@echo "Starting development with Docker (hot reload)..."
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

# =============================================================================
# Production with Docker
# =============================================================================

prod:
	@echo "Starting production containers..."
	docker-compose up --build -d
	@echo ""
	@echo "Application running at:"
	@echo "  Frontend: http://localhost:3000"
	@echo "  Backend:  http://localhost:8080"

prod-build:
	@echo "Building production images..."
	docker-compose build

prod-down:
	docker-compose down

prod-logs:
	docker-compose logs -f

# =============================================================================
# Testing
# =============================================================================

test: test-backend

test-backend:
	@echo "Running backend tests..."
	cd stock-screener-backend && go test ./...

test-verbose:
	@echo "Running backend tests (verbose)..."
	cd stock-screener-backend && go test ./... -v

test-coverage:
	@echo "Running backend tests with coverage..."
	cd stock-screener-backend && go test ./... -cover -coverprofile=coverage.out
	cd stock-screener-backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: stock-screener-backend/coverage.html"

# =============================================================================
# Dependencies
# =============================================================================

deps:
	@echo "Installing dependencies..."
	cd stock-screener-backend && go mod download
	cd stock-screener-frontend && npm install

deps-backend:
	cd stock-screener-backend && go mod download

deps-frontend:
	cd stock-screener-frontend && npm install

# =============================================================================
# Cleanup
# =============================================================================

clean:
	@echo "Cleaning up..."
	docker-compose down -v --remove-orphans 2>/dev/null || true
	rm -rf stock-screener-backend/tmp
	rm -rf stock-screener-frontend/dist
	rm -rf stock-screener-frontend/node_modules/.vite
	rm -f stock-screener-backend/coverage.out stock-screener-backend/coverage.html
	@echo "Clean complete."
