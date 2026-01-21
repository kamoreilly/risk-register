# AGENTS.md

This document provides guidelines for agentic coding agents working in this repository.

## Project Overview

Risk Register - A full-stack risk management application with a Go backend and React/TypeScript frontend.

**Structure:**
- `backend/` - Go API server (Fiber framework)
- `frontend/` - React web application (TanStack Start, Vite, Turbo monorepo)

## Build Commands

### Root Level
```bash
make all              # Install all dependencies (Go + Bun)
make dev              # Start backend + web + native dev servers
make dev-backend      # Start backend only with hot reload
make dev-frontend     # Start frontend only
make stop             # Stop all dev servers
make docker-run       # Start PostgreSQL database
make docker-down      # Stop PostgreSQL database
make clean            # Clean build artifacts
```

### Backend (Go)
```bash
cd backend
make build            # Build binary to tmp/main
make test             # Run all tests (unit + integration)
make itest            # Run integration tests only
make watch            # Start AIR hot reload
```

**Run single test:**
```bash
cd backend
go test ./internal/server -v                              # Single package
go test ./internal/server -run TestHandler -v             # Single test
go test ./internal/database -v -timeout 120s              # With longer timeout (DB tests)
```

### Frontend (TypeScript/Bun)
```bash
cd frontend
bun install           # Install dependencies
bun run build         # Build all apps/packages
bun run dev           # Start dev servers
bun run dev:web       # Start web app only
bun run dev:native    # Start native (Expo) app only
bun run check-types   # Type-check all packages
```

**Individual app:**
```bash
cd frontend/apps/web
bun run build         # Build web app
bun run dev           # Start web dev server
bun run serve         # Preview production build
```

## Code Style Guidelines

### TypeScript/React

**Imports:**
- Use `@/` alias for relative imports from app root (e.g., `@/components/ui/button`)
- Group imports: external → internal → CSS/styling
- Use named imports for libraries (`import { useState } from 'react'`)
- Use default imports for components (`import Button from '@/components/ui/button'`)

**Naming:**
- Components: PascalCase (e.g., `RiskCard`)
- Hooks: camelCase with `use` prefix (e.g., `useRisks`)
- Variables/functions: camelCase (e.g., `riskList`, `calculateSeverity`)
- Constants: SCREAMING_SNAKE_CASE (e.g., `MAX_SEVERITY`)
- Types/interfaces: PascalCase (e.g., `Risk`, `RiskStatus`)
- Files: kebab-case for utilities, PascalCase for components

**TypeScript:**
- Enable `strict: true` (see `tsconfig.base.json`)
- Use `noUncheckedIndexedAccess: true` - always define index signatures
- Use `noUnusedLocals: true` and `noUnusedParameters: true`
- Prefer interfaces over type aliases for object types
- Use `as const` for literal values that shouldn't widen
- Define explicit return types for complex functions
- Use `Record<K, V>` instead of `{ [key: K]: V }`

**React:**
- Use function components only (no class components)
- Use TypeScript types for props (no PropTypes)
- Memoize expensive computations with `useMemo`
- Memoize callbacks with `useCallback`
- Use `React.FC` sparingly; prefer explicit prop types
- Use `as const` assertion for literal string/number values in JSX

**Formatting:**
- 2 spaces for indentation (check `.editorconfig` if present)
- Use `cn()` utility for Tailwind class merging (see `lib/utils.ts`)
- Prefer Tailwind utility classes over custom CSS
- Use semantic HTML elements

**Error Handling:**
- Use error boundaries for component error catching
- Return typed error responses from API handlers
- Log errors with context; don't expose internals to users

### Go

**Imports:**
- Group: standard library → third-party → internal
- Use blank import (`_`) only for side effects (e.g., database drivers)
- Remove unused imports (Go will compile error)

**Naming:**
- Packages: lowercase, short, no underscores (e.g., `server`, `database`)
- Exported types/functions: PascalCase (e.g., `func New()`)
- Unexported: camelCase (e.g., `func newServer()`)
- Interfaces: -er suffix when possible (e.g., `Reader`, `Service`)
- Variables: camelCase; prefer brevity with clarity
- Constants: PascalCase for exported, camelCase for unexported
- Error variables: `Err` prefix (e.g., `ErrNotFound`)

**Error Handling:**
- Return errors as values; avoid panic for expected failures
- Wrap errors with context using `fmt.Errorf("doing X: %w", err)`
- Use custom error types for API responses
- Check errors immediately (`if err != nil { ... }`)

**Testing:**
- Use table-driven tests when possible
- Test file suffix: `_test.go`
- Test function prefix: `Test` (e.g., `TestHandler`)
- Use `t.Fatalf` for fatal failures, `t.Errorf` for assertions
- Integration tests in `internal/database` use testcontainers

**Code Structure:**
- One package per directory (or group related types)
- Keep main packages small; delegate to internal packages
- Use struct embedding sparingly; prefer explicit fields
- Pass context as first argument (`ctx context.Context`)

## Configuration Files

- `frontend/tsconfig.json` - Extends `@frontend/config/tsconfig.base.json`
- `frontend/turbo.json` - Turbo build configuration
- `backend/.air.toml` - AIR hot reload config (auto-generated)
- `backend/go.mod` - Go modules (Go 1.25.5)

## Key Libraries

**Frontend:**
- TanStack Router + TanStack Start for routing
- Recharts for data visualization
- Lucide React for icons
- Tailwind CSS + shadcn/ui for styling
- Zod for validation

**Backend:**
- Fiber v2 for HTTP server
- pgx for PostgreSQL
- testcontainers-go for integration tests

## Database

- PostgreSQL via Docker (see `backend/docker-compose.yml`)
- Start with `make docker-run` in backend directory
- Integration tests automatically start testcontainer

## Notes

- No existing Cursor rules or Copilot instructions found
- Frontend uses Bun package manager (v1.3.4)
- Frontend is a Turbo monorepo with workspaces
- Frontend apps use absolute imports with `@/` alias
- Backend uses AIR for hot reload during development
