# Risk Register

A full-stack risk management application.

## Project Structure
- `backend/` - Go API server
- `frontend/` - React web application (with TanStack Start)

## Quick Start

### Prerequisites
- Go 1.x
- Bun
- Docker (for PostgreSQL)
- AIR (installed automatically)

### Setup
```bash
# Install all dependencies
make all

# Start PostgreSQL database
make docker-run

# Start backend and frontend
make dev
```

The application will be available at:
- Frontend: http://localhost:3001
- Backend API: http://localhost:8080

## Available Commands
```bash
make help              # Show all available commands
```

## Environment Setup
1. Copy `backend/.env` and update database credentials
2. Run `make docker-run` to start PostgreSQL
3. Run `make dev` to start the application
