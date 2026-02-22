# CLAUDE.md

Guidelines for Claude Code when working in this repository.

## Project Overview

Full-stack risk management application.

- **Backend:** Go API server (Fiber framework) in `backend/`
- **Frontend:** React web app (TanStack Start, Vite) in `frontend/`

## Quick Start

```bash
make all              # Install all dependencies
make docker-run       # Start PostgreSQL database
make dev              # Start backend + frontend dev servers
```

Ports: Frontend http://localhost:3001 | Backend API http://localhost:8080

## Commands

### Root
```bash
make dev              # Start all dev servers
make stop             # Stop dev servers
make docker-run       # Start PostgreSQL
make docker-down      # Stop PostgreSQL
```

### Backend (Go)
```bash
cd backend
make build            # Build binary to tmp/main
make test             # Run all tests
make itest            # Run integration tests (uses testcontainers)
make watch            # Start AIR hot reload
go test ./internal/server -run TestHandler -v  # Single test
```

### Frontend (TypeScript/Bun)
```bash
cd frontend
bun install           # Install dependencies
bun run build         # Build all apps/packages
bun run dev           # Start dev servers
bun run dev:web       # Start web app only
bun run check-types   # Type-check all packages
```

## Architecture

```
risk-register/
├── backend/
│   ├── cmd/api/main.go       # Entry point
│   ├── internal/
│   │   ├── database/         # DB layer (pgx)
│   │   └── server/           # HTTP handlers (Fiber)
│   └── docker-compose.yml    # PostgreSQL
└── frontend/                 # Turbo monorepo
    ├── apps/web/             # Main web app
    │   └── src/
    │       ├── routes/       # TanStack Router pages
    │       └── components/   # React components
    └── packages/
        ├── ui/               # Shared UI components (shadcn)
        ├── config/           # Shared TypeScript config
        └── env/              # Environment utilities
```

## Code Style

See [AGENTS.md](./AGENTS.md) for detailed style guidelines.

**Key points:**
- TypeScript: `strict: true`, `noUncheckedIndexedAccess: true`
- Imports: Use `@/` alias for app-relative imports
- React: Function components only, explicit prop types
- Go: Standard Go conventions, table-driven tests

## Gotchas

- Frontend uses **Bun** (not npm/yarn) - `packageManager: bun@1.3.4`
- Integration tests in backend use testcontainers (Docker required)
- Frontend imports use `@/` alias (configured in vite-tsconfig-paths)
- Run `make docker-run` before backend integration tests
- Backend hot reload uses AIR (auto-installed by Makefile)
