.PHONY: all dev dev-all stop dev-backend dev-frontend docker-run docker-down clean install-deps git-init help

all: install-deps

# Install all dependencies (Go modules and Bun packages)
install-deps:
	@echo "Installing backend dependencies..."
	@cd backend && go mod download
	@echo "Installing frontend dependencies..."
	@cd frontend && bun install

# Stop dev servers (best-effort) by killing known dev ports.
# Note: this avoids broad process kills like `killall node`.
stop:
	@bash -eu -o pipefail -c 'echo "Stopping dev services (best-effort)..."; ports=(8080 3001 5173 8081); for port in "$${ports[@]}"; do pids="$$(lsof -ti :$$port 2>/dev/null || true)"; if [ -n "$$pids" ]; then echo "- stopping :$$port ($$pids)"; kill -TERM $$pids 2>/dev/null || true; sleep 1; kill -KILL $$pids 2>/dev/null || true; fi; done; echo "Stop complete."'

# Stop any running services first, then start backend + web.
# On Ctrl+C / termination, force shut down all servers.
dev: stop
	@bash -eu -o pipefail -c 'root_dir="$$(pwd)"; echo "Starting backend + web..."; (cd "$$root_dir/backend" && if ! command -v air >/dev/null; then go install github.com/air-verse/air@latest; fi && exec air) & backend_pid=$$!; (cd "$$root_dir/frontend/apps/web" && exec bun run dev) & web_pid=$$!; cleanup(){ echo ""; echo "Shutting down dev servers..."; pids="$$backend_pid $$web_pid"; kill -TERM $$pids 2>/dev/null || true; for i in 1 2 3 4 5; do alive=""; for pid in $$pids; do if kill -0 $$pid 2>/dev/null; then alive="$$alive $$pid"; fi; done; [ -z "$$alive" ] && break; sleep 1; done; for pid in $$pids; do if kill -0 $$pid 2>/dev/null; then kill -KILL $$pid 2>/dev/null || true; fi; done; $(MAKE) -C "$$root_dir" -s stop || true; }; trap cleanup INT TERM EXIT; echo "Backend PID: $$backend_pid"; echo "Web PID: $$web_pid"; wait'

# Alias for `dev` (kept for convenience).
dev-all: dev

# Start backend with hot reload (AIR)
dev-backend:
	@echo "Starting backend with AIR hot reload..."
	@cd backend && if ! command -v air > /dev/null; then go install github.com/air-verse/air@latest; fi && air

# Start frontend dev server
dev-frontend:
	@echo "Starting frontend..."
	@cd frontend/apps/web && bun run dev

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
	@echo "  make dev              - Start backend + web (stops existing processes; cleans up on exit)"
	@echo "  make dev-all          - Alias for make dev"
	@echo "  make stop             - Stop dev servers (kills known dev ports)"
	@echo "  make dev-backend      - Start backend only with hot reload"
	@echo "  make dev-frontend     - Start frontend only""
	@echo "  make docker-down      - Stop PostgreSQL database"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make git-init         - Initialize git repository"
	@echo "  make help             - Show this help message"
