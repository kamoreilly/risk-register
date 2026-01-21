# Risk Register Backend

Go backend API for the Risk Register application.

## Tech Stack
- Go 1.x
- PostgreSQL with pgx driver
- Docker Compose for database

## Environment Variables
Create a `.env` file with:
```
PORT=8080
APP_ENV=local
RISK_REGISTER_DB_HOST=localhost
RISK_REGISTER_DB_PORT=5432
RISK_REGISTER_DB_DATABASE=risk_register
RISK_REGISTER_DB_USERNAME=risk_register
RISK_REGISTER_DB_PASSWORD=risk_register
RISK_REGISTER_DB_SCHEMA=public
```

## Development

### Using the root Makefile (Recommended)
```bash
# From project root
make dev              # Start backend + frontend
make dev-backend      # Start backend only
make docker-run       # Start PostgreSQL
```

### Backend-specific commands
```bash
make build            # Build the application
make run              # Run the application
make watch            # Live reload with AIR
make test             # Run tests
make itest            # Integration tests
make clean            # Clean binaries
```

## Database
Start PostgreSQL: `make docker-run`
Stop PostgreSQL: `make docker-down`
