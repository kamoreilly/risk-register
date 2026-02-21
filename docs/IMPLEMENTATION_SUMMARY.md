# Risk Register Implementation Summary

**Branch:** `feature/risk-register-auth`
**Date:** 2026-02-20
**Total Commits:** 49

---

## Overview

Full-stack risk management application built using a layered slice approach. Each slice delivers a complete, testable feature end-to-end from database to UI.

### Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.22+, Fiber v2, PostgreSQL (pgx) |
| Auth | JWT (golang-jwt), bcrypt |
| Migrations | golang-migrate with embedded SQL |
| Frontend | React 18, TypeScript, Bun |
| Routing | TanStack Router |
| State | TanStack Query, Zustand |
| Forms | React Hook Form, Zod |
| UI | shadcn/ui, Tailwind CSS |
| Drag & Drop | @dnd-kit |

---

## Architecture

```
risk-register/
├── backend/
│   ├── cmd/api/main.go              # Entry point
│   └── internal/
│       ├── migrations/              # Embedded SQL migrations
│       │   └── migrations/
│       │       ├── 001_users.up.sql
│       │       ├── 002_categories.up.sql
│       │       ├── 003_risks.up.sql
│       │       ├── 004_mitigations.up.sql
│       │       ├── 005_frameworks.up.sql
│       │       ├── 006_risk_framework_controls.up.sql
│       │       └── 007_audit_logs.up.sql
│       ├── models/                  # Go structs
│       │   ├── user.go
│       │   ├── category.go
│       │   ├── risk.go
│       │   ├── mitigation.go
│       │   ├── framework.go
│       │   └── audit.go
│       ├── database/                # Repository layer
│       │   ├── users.go
│       │   ├── categories.go
│       │   ├── risks.go
│       │   ├── mitigations.go
│       │   ├── frameworks.go
│       │   └── audit.go
│       ├── handlers/                # HTTP handlers
│       │   ├── auth.go
│       │   ├── risks.go
│       │   ├── categories.go
│       │   ├── mitigations.go
│       │   ├── dashboard.go
│       │   ├── frameworks.go
│       │   ├── ai.go
│       │   └── audit.go
│       ├── middleware/              # Auth middleware
│       │   └── auth.go
│       ├── auth/                    # JWT, password hashing
│       │   └── auth.go
│       └── server/                  # Fiber server setup
│           ├── server.go
│           └── routes.go
│
└── frontend/
    └── apps/web/
        └── src/
            ├── routes/              # TanStack Router pages
            │   ├── login.tsx
            │   ├── register.tsx
            │   └── app/
            │       ├── __root.tsx   # Protected layout
            │       ├── index.tsx    # Dashboard
            │       ├── calendar.tsx # Review calendar
            │       ├── board.tsx    # Kanban board
            │       └── risks/
            │           ├── index.tsx  # Risk list
            │           ├── new.tsx    # Create risk
            │           └── $id.tsx    # Risk detail
            ├── components/
            │   └── ui/              # shadcn/ui components
            ├── hooks/
            │   ├── useAuth.ts
            │   ├── useRisks.ts
            │   ├── useMitigations.ts
            │   ├── useDashboard.ts
            │   ├── useFrameworks.ts
            │   ├── useAI.ts
            │   └── useAudit.ts
            ├── types/
            │   ├── auth.ts
            │   ├── risk.ts
            │   ├── mitigation.ts
            │   ├── dashboard.ts
            │   ├── framework.ts
            │   ├── ai.ts
            │   └── audit.ts
            └── lib/
                ├── api.ts           # API client
                └── utils.ts         # Utilities
```

---

## Data Model

### Entity Relationship Diagram

```
┌─────────────┐       ┌─────────────┐
│   users     │       │ categories  │
├─────────────┤       ├─────────────┤
│ id (PK)     │       │ id (PK)     │
│ email       │       │ name        │
│ password    │       │ description │
│ name        │       │ color       │
│ role        │       └──────┬──────┘
└──────┬──────┘              │
       │                     │
       │    ┌────────────────┘
       │    │
       ▼    ▼
┌─────────────────┐
│     risks       │
├─────────────────┤
│ id (PK)         │
│ title           │
│ description     │
│ owner_id (FK)   │──────┐
│ status          │      │
│ severity        │      │
│ category_id(FK) │      │
│ review_date     │      │
│ created_by (FK) │──┐   │
│ updated_by (FK) │──┼───┘
└────────┬────────┘  │
         │           │
         ▼           │
┌─────────────────┐  │
│  mitigations    │  │
├─────────────────┤  │
│ id (PK)         │  │
│ risk_id (FK)    │  │
│ description     │  │
│ owner           │  │
│ status          │  │
│ due_date        │  │
│ created_by (FK) │──┼──┐
│ updated_by (FK) │──┼──┤
└─────────────────┘  │  │
                     │  │
┌─────────────────┐  │  │
│  frameworks     │  │  │
├─────────────────┤  │  │
│ id (PK)         │  │  │
│ name            │  │  │
│ description     │  │  │
└────────┬────────┘  │  │
         │           │  │
         ▼           │  │
┌────────────────────┴──┴──┐
│ risk_framework_controls │
├──────────────────────────┤
│ id (PK)                  │
│ risk_id (FK)             │
│ framework_id (FK)        │
│ control_ref              │
│ notes                    │
│ created_by (FK)          │
└──────────────────────────┘

┌─────────────────┐
│   audit_logs    │
├─────────────────┤
│ id (PK)         │
│ entity_type     │
│ entity_id       │
│ action          │
│ changes (JSONB) │
│ user_id (FK)    │
│ created_at      │
└─────────────────┘
```

### Enum Types

```sql
-- User roles
CREATE TYPE user_role AS ENUM ('admin', 'member');

-- Risk status
CREATE TYPE risk_status AS ENUM ('open', 'mitigating', 'resolved', 'accepted');

-- Risk severity
CREATE TYPE risk_severity AS ENUM ('low', 'medium', 'high', 'critical');

-- Mitigation status
CREATE TYPE mitigation_status AS ENUM ('planned', 'in_progress', 'completed', 'cancelled');

-- Audit action
CREATE TYPE audit_action AS ENUM ('created', 'updated', 'deleted');
```

---

## API Endpoints

### Authentication (`/api/v1/auth`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/register` | Create new account |
| POST | `/login` | Authenticate, return JWT |
| GET | `/me` | Get current user profile |

### Risks (`/api/v1/risks`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | List risks (paginated, filterable) |
| POST | `/` | Create risk |
| GET | `/:id` | Get single risk |
| PUT | `/:id` | Update risk |
| DELETE | `/:id` | Delete risk |
| GET | `/:id/audit` | Get audit logs for risk |

### Mitigations (`/api/v1/risks/:riskId/mitigations`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | List mitigations for risk |
| POST | `/` | Create mitigation |
| PUT | `/:id` | Update mitigation |
| DELETE | `/:id` | Delete mitigation |

### Categories (`/api/v1/categories`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | List all categories |

### Frameworks (`/api/v1/frameworks`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | List all frameworks |
| POST | `/` | Create framework (admin) |

### Controls (`/api/v1/risks/:riskId/controls`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | List controls for risk |
| POST | `/` | Link control to risk |
| DELETE | `/:id` | Unlink control |

### Dashboard (`/api/v1/dashboard`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/summary` | Risk counts by status/severity/category |
| GET | `/reviews/upcoming` | Risks with upcoming review dates |
| GET | `/reviews/overdue` | Risks past review date |

### AI (`/api/v1/ai`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/summarize` | Generate risk summary (stub) |
| POST | `/draft-mitigation` | Draft mitigation text (stub) |

---

## Slice Breakdown

### Slice 1: Authentication (14 commits)

**Deliverable:** Users can register, log in, and access protected routes.

#### Backend
- Migrations infrastructure with embedded SQL files
- Users table with UUID, email, password_hash, name, role
- JWT generation and validation with 24-hour expiry
- bcrypt password hashing
- Auth middleware extracting user from Bearer token
- Auth handlers: Register, Login, Me
- User repository with Create, GetByEmail, GetByID

#### Frontend
- API client with token storage and Authorization header
- Auth types: User, AuthResponse, LoginInput, RegisterInput
- useAuth hook with TanStack Query
- Login page wired to backend
- Register page
- Protected route layout with auth check
- Dashboard placeholder

**Files Created:**
```
backend/internal/migrations/migrations.go
backend/internal/migrations/migrations/001_users.up.sql
backend/internal/migrations/migrations/001_users.down.sql
backend/internal/models/user.go
backend/internal/auth/auth.go
backend/internal/middleware/auth.go
backend/internal/middleware/auth_test.go
backend/internal/database/users.go
backend/internal/database/users_test.go
backend/internal/handlers/auth.go
backend/internal/server/server.go
backend/internal/server/routes.go
frontend/apps/web/src/lib/api.ts
frontend/apps/web/src/types/auth.ts
frontend/apps/web/src/hooks/useAuth.ts
frontend/apps/web/src/routes/login.tsx
frontend/apps/web/src/routes/register.tsx
frontend/apps/web/src/routes/app/__root.tsx
frontend/apps/web/src/routes/app/index.tsx
```

---

### Slice 2: Risks CRUD (8 commits)

**Deliverable:** Full risk management with categories.

#### Backend
- Categories table with 6 default categories
- Risks table with status/severity enums, FKs to users and categories
- Category and Risk models
- CategoryRepository and RiskRepository
- Risk and Category handlers with full CRUD
- Filtering: status, severity, category, owner, search
- Pagination: page, limit
- Sorting: sort, order

#### Frontend
- Risk and Category types
- useRisks, useRisk, useCreateRisk, useUpdateRisk, useDeleteRisk hooks
- Risk list page with filters and pagination
- Risk detail page with edit mode
- New risk form with validation

**Files Created:**
```
backend/internal/migrations/migrations/002_categories.up.sql
backend/internal/migrations/migrations/003_risks.up.sql
backend/internal/models/category.go
backend/internal/models/risk.go
backend/internal/database/categories.go
backend/internal/database/risks.go
backend/internal/handlers/categories.go
backend/internal/handlers/risks.go
frontend/apps/web/src/types/risk.ts
frontend/apps/web/src/hooks/useRisks.ts
frontend/apps/web/src/routes/app/risks/index.tsx
frontend/apps/web/src/routes/app/risks/$id.tsx
frontend/apps/web/src/routes/app/risks/new.tsx
frontend/apps/web/src/components/ui/select.tsx
```

---

### Slice 3: Mitigations (7 commits)

**Deliverable:** Track mitigation actions for each risk.

#### Backend
- Mitigations table with status enum, FK to risks
- Mitigation model with MitigationStatus enum
- MitigationRepository with CRUD operations
- MitigationHandler with nested routes under risks
- Cascade delete when parent risk is deleted

#### Frontend
- Mitigation types: Mitigation, MitigationStatus, CreateMitigationInput
- useMitigations, useCreateMitigation, useUpdateMitigation, useDeleteMitigation hooks
- Mitigation list section on risk detail page
- Add/edit/delete mitigation inline
- Status badges with colors (planned=yellow, in_progress=blue, completed=green, cancelled=gray)

**Files Created:**
```
backend/internal/migrations/migrations/004_mitigations.up.sql
backend/internal/models/mitigation.go
backend/internal/database/mitigations.go
backend/internal/database/mitigations_test.go
backend/internal/handlers/mitigations.go
frontend/apps/web/src/types/mitigation.ts
frontend/apps/web/src/hooks/useMitigations.ts
```

---

### Slice 4: Dashboard (3 commits)

**Deliverable:** At-a-glance overview with statistics.

#### Backend
- DashboardHandler with Summary endpoint
- Aggregation queries: total count, by status, by severity, by category
- Overdue reviews count

#### Frontend
- DashboardSummary and CategoryCount types
- useDashboardSummary hook
- Dashboard page with widgets:
  - Stat cards: Total Risks, Open, Mitigating, Overdue
  - Severity breakdown with colored indicators
  - Category breakdown
  - Quick links to risk list

**Files Created:**
```
backend/internal/handlers/dashboard.go
frontend/apps/web/src/types/dashboard.ts
frontend/apps/web/src/hooks/useDashboard.ts
```

---

### Slice 5: Review Calendar (2 commits)

**Deliverable:** Visual review schedule with calendar view.

#### Backend
- UpcomingReviews endpoint with configurable days parameter
- OverdueReviews endpoint
- Returns risk data with review_date, severity, status

#### Frontend
- ReviewRisk and ReviewListResponse types
- useUpcomingReviews and useOverdueReviews hooks
- Calendar page with:
  - Month navigation (prev/next)
  - Calendar grid with risk markers
  - Color-coded markers: red (overdue), yellow (<7 days), blue (>7 days)
  - Sidebar with overdue and upcoming lists
  - Click to navigate to risk detail

**Files Created:**
```
frontend/apps/web/src/routes/app/calendar.tsx
```

---

### Slice 6: Risk Board / Kanban (1 commit)

**Deliverable:** Visual workflow management with drag-and-drop.

#### Frontend Only
- @dnd-kit integration for drag-and-drop
- Board with 4 columns: Open, Mitigating, Resolved, Accepted
- Risk cards with severity badges and owner
- Drag-and-drop to update status
- useUpdateRiskStatus hook for API updates
- Column colors and drop zone highlighting

**Files Created:**
```
frontend/apps/web/src/routes/app/board.tsx
```

**Dependencies Added:**
```
@dnd-kit/core
@dnd-kit/sortable
@dnd-kit/utilities
```

---

### Slice 7: Framework Mapping (4 commits)

**Deliverable:** Compliance traceability with framework controls.

#### Backend
- Frameworks table with 5 default frameworks (ISO 27001, SOC 2, NIST CSF, GDPR, HIPAA)
- risk_framework_controls junction table
- Framework and RiskFrameworkControl models
- FrameworkRepository and RiskFrameworkControlRepository
- FrameworkHandler and ControlHandler
- Link/unlink controls to risks

#### Frontend
- Framework and RiskFrameworkControl types
- useFrameworks, useRiskControls, useLinkControl, useUnlinkControl hooks
- Compliance Controls section on risk detail
- Blue badges showing framework:control_ref
- Form to link new controls
- Delete to unlink

**Files Created:**
```
backend/internal/migrations/migrations/005_frameworks.up.sql
backend/internal/migrations/migrations/006_risk_framework_controls.up.sql
backend/internal/models/framework.go
backend/internal/database/frameworks.go
backend/internal/handlers/frameworks.go
frontend/apps/web/src/types/framework.ts
frontend/apps/web/src/hooks/useFrameworks.ts
```

---

### Slice 8: AI Stubs (2 commits)

**Deliverable:** UI ready for real AI integration.

#### Backend
- AIHandler with Summarize and DraftMitigation endpoints
- Accepts risk data, returns placeholder text
- Ready for replacement with real AI service

#### Frontend
- AI types: SummarizeRequest, SummarizeResponse, DraftMitigationRequest, DraftMitigationResponse
- useSummarize and useDraftMitigation hooks
- "Summarize" button on risk detail - displays summary in blue box
- "Draft with AI" button in mitigation form - pre-fills description

**Files Created:**
```
backend/internal/handlers/ai.go
frontend/apps/web/src/types/ai.ts
frontend/apps/web/src/hooks/useAI.ts
```

---

### Slice 9: Audit Trail (3 commits)

**Deliverable:** Full change history for risks.

#### Backend
- audit_logs table with action enum and JSONB changes
- AuditLog model with AuditAction enum
- AuditLogRepository with Create and ListByEntity
- Audit logging integrated into RiskHandler:
  - Created: logs initial field values
  - Updated: logs from/to values for changed fields
  - Deleted: logs deletion action
- AuditHandler with ListByRisk endpoint
- Joins with users table for user names

#### Frontend
- AuditLog and AuditAction types
- useAuditLogs hook
- Audit History section on risk detail:
  - Vertical timeline with colored dots
  - User name, action description, timestamp
  - For updates: list of changed fields with before/after values

**Files Created:**
```
backend/internal/migrations/migrations/007_audit_logs.up.sql
backend/internal/models/audit.go
backend/internal/database/audit.go
backend/internal/handlers/audit.go
frontend/apps/web/src/types/audit.ts
frontend/apps/web/src/hooks/useAudit.ts
```

---

## Environment Variables

### Backend (.env)
```
PORT=8080
RISK_REGISTER_DB_HOST=localhost
RISK_REGISTER_DB_PORT=5432
RISK_REGISTER_DB_USERNAME=postgres
RISK_REGISTER_DB_PASSWORD=postgres
RISK_REGISTER_DB_DATABASE=risk_register
RISK_REGISTER_DB_SCHEMA=public
JWT_SECRET=your-secret-key-here
```

### Frontend (.env)
```
VITE_API_URL=http://localhost:8080
```

---

## Running the Application

### Quick Start
```bash
# Install dependencies
make all

# Start PostgreSQL
make docker-run

# Start backend + frontend
make dev
```

### Individual Commands

**Backend:**
```bash
cd backend
make build      # Build binary
make test       # Run tests
make watch      # Hot reload with AIR
```

**Frontend:**
```bash
cd frontend
bun install     # Install dependencies
bun run dev     # Start dev server
bun run build   # Production build
```

### Ports
- Frontend: http://localhost:3001
- Backend API: http://localhost:8080

---

## Key Patterns Used

### Repository Pattern
All database access goes through repository interfaces, making the code testable and decoupled from the database implementation.

```go
type RiskRepository interface {
    Create(ctx context.Context, input *models.CreateRiskInput, ownerID string) (*models.Risk, error)
    GetByID(ctx context.Context, id string) (*models.Risk, error)
    List(ctx context.Context, params *RiskListParams) (*RiskListResponse, error)
    Update(ctx context.Context, id string, input *models.UpdateRiskInput, updatedBy string) (*models.Risk, error)
    Delete(ctx context.Context, id string) error
}
```

### Handler Pattern
Handlers are thin wrappers that parse requests, call repositories, and return responses.

```go
func (h *RiskHandler) Create(c *fiber.Ctx) error {
    user := middleware.GetUserFromContext(c)
    if user == nil {
        return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
    }

    var input models.CreateRiskInput
    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
    }

    risk, err := h.riskRepo.Create(c.Context(), &input, user.UserID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to create risk"})
    }

    return c.Status(201).JSON(risk)
}
```

### TanStack Query Pattern
All server state is managed through TanStack Query hooks.

```typescript
export function useRisks(params?: RiskListParams) {
  return useQuery({
    queryKey: ['risks', params],
    queryFn: async () => {
      const response = await api.get<RiskListResponse>(`/api/v1/risks${queryString}`);
      return response.data;
    },
  });
}
```

---

## Future Enhancements

1. **Real AI Integration** - Replace stub endpoints with actual LLM calls
2. **Multi-tenancy** - Add organization scoping
3. **SSO Integration** - OAuth2/SAML support
4. **Real-time Updates** - WebSocket for live collaboration
5. **Export/Import** - CSV, PDF exports
6. **Notifications** - Email alerts for review dates
7. **Advanced Reporting** - Custom reports and charts
8. **Bulk Operations** - Multi-select actions
9. **API Rate Limiting** - Protect against abuse
10. **Full Text Search** - PostgreSQL full-text search integration

---

## Git History

All 49 commits on `feature/risk-register-auth` branch:

```
8fedc6e feat(frontend): add audit timeline to risk detail
ce8c5c4 feat(backend): add audit logging
f9bb288 feat(backend): add audit_logs table migration
c4c03fb feat(frontend): add AI summarize and draft buttons
ef66aab feat(backend): add AI stub endpoints
96f84bf feat(frontend): add framework control badges to risk detail
66f9c9e feat(backend): add framework and control handlers
693d85d feat(backend): add Framework model and repository
6f5a907 feat(backend): add frameworks and controls tables
b4f36a8 feat(frontend): add Kanban board page
6514062 feat(frontend): add review calendar page
671411e feat(backend): add review calendar endpoints
23cefc3 feat(frontend): build dashboard page with widgets
2e9b7b1 feat(frontend): add dashboard types and hooks
0401e96 feat(backend): add dashboard summary endpoint
74e7acd feat(frontend): add mitigations to risk detail page
7f6e167 feat(frontend): add mitigation hooks
b0eebdb feat(frontend): add mitigation types
b84c9ac feat(backend): wire mitigation routes
bf0c05b feat(backend): add MitigationHandler
5ecd925 feat(backend): add MitigationRepository with tests
8790584 feat(backend): add Mitigation model
c23bd3a feat(backend): add mitigations table migration
37f0b00 feat(frontend): complete risk detail and new risk pages
d130d0a feat(frontend): add risk list page
8b8f7bc feat(frontend): add useRisks hook
2ba3f23 feat(frontend): add risk types
6c2b759 feat(backend): wire up risk routes
89b4c0a feat(backend): add risk and category handlers
3f3e295 feat(backend): add risk and category repositories
8e879a7 feat(backend): add Risk and Category models
fc18093 feat(backend): add categories and risks migrations
0042add feat(frontend): add dashboard placeholder
63af281 feat(frontend): add protected route layout
40aebb1 feat(frontend): add register page
e5a43ad feat(frontend): wire login page to backend auth
4eb4743 feat(frontend): add auth types and useAuth hook
111e34b feat(frontend): add API client
b3648f0 feat(backend): run migrations on startup
2dc223d feat(backend): wire up auth routes
acde624 feat(backend): add auth handlers with tests
5566367 feat(backend): add user repository with interface
d252920 feat(backend): add auth middleware with tests
bd364c9 feat(backend): add auth service with JWT and bcrypt
dd0fdd9 feat(backend): add User model
ec9561b feat(backend): add JWT and bcrypt dependencies
b3b2a70 fix(backend): add defer m.Close() and package docs
37a5f6d chore(backend): tidy go.mod after adding dependencies
88c2668 feat(backend): add migrations infrastructure and users table
```

---

*Generated on 2026-02-20*
