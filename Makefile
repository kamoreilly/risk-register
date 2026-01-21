.PHONY: all dev dev-backend dev-frontend docker-run docker-down clean install-deps git-init help

all: install-deps

# Install all dependencies (Go modules and Bun packages)
install-deps:
	@echo "Installing backend dependencies..."
	@cd backend && go mod download
	@echo "Installing frontend dependencies..."
	@cd frontend && bun install

# Kill all running services first, then start both backend and frontend in development mode
dev:
	@echo "Killing all services..."
	@-killall -9 air 2>/dev/null || true
	@-killall -9 main 2>/dev/null || true
	@-killall -9 node 2>/dev/null || true
	@-lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	@-lsof -ti:3001 | xargs kill -9 2>/dev/null || true
	@sleep 1
	@echo "Starting backend and frontend..."
	@make dev-backend & make dev-frontend

# Start backend with hot reload (AIR)
dev-backend:
	@echo "Starting backend with AIR hot reload..."
	@cd backend && if ! command -v air > /dev/null; then go install github.com/air-verse/air@latest; fi && air

# Start frontend dev server
dev-frontend:
	@echo "Starting frontend..."
	@cd frontend && bun run dev

# Start PostgreSQL database
docker-run:
	@cd backend && make docker-run

# Stop PostgreSQL database
docker-down:
	@cd backend && make docker-down

# Clean build artifacts
clean:
	@cd backend && make clean

# Initialize git repository (without initial commit)
git-init:
	@if [ ! -d .git ]; then \
		git init; \
		echo "Git repository initialized"; \
	else \
		echo "Git repository already exists"; \
	fi

# Helper to check/install AIR
check-air:
	@if command -v air > /dev/null; then \
		echo ""; \
	else \
		go install github.com/air-verse/air@latest; \
	fi

help:
	@echo "Available targets:"
	@echo "  make all              - Install all dependencies"
	@echo "  make dev              - Start backend and frontend (kills existing processes)"
	@echo "  make dev-backend      - Start backend only with hot reload"
	@echo "  make dev-frontend     - Start frontend only"
	@echo "  make docker-run       - Start PostgreSQL database"
	@echo "  make docker-down      - Stop PostgreSQL database"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make git-init         - Initialize git repository"
	@echo "  make help             - Show this help message"
