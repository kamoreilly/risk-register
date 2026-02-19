# Risk Register Application Design

**Date:** 2026-02-19
**Status:** Approved
**Approach:** Layered Build

## Overview

Full-featured risk management application with:
- Risk tracking with categories, mitigations, and audit trail
- Email/password authentication
- Multiple views: dashboard, list, board, calendar
- Compliance framework mapping
- AI features (stubbed for future implementation)
- Single-tenant deployment

## Data Model

### users
| Column | Type | Description |
|--------|------|-------------|
| id | uuid | Primary key |
| email | string | Unique email |
| password_hash | string | Bcrypt hash |
| name | string | Display name |
| role | enum | admin, member |
| created_at | timestamp | |
| updated_at | timestamp | |

### categories
| Column | Type | Description |
|--------|------|-------------|
| id | uuid | Primary key |
| name | string | e.g., "Security", "Operational" |
| description | text | |

### risks
| Column | Type | Description |
|--------|------|-------------|
| id | uuid | Primary key |
| title | string | |
| description | text | |
| owner_id | uuid | FK to users |
| status | enum | open, mitigating, resolved, accepted |
| severity | enum | low, medium, high, critical |
| category_id | uuid | FK to categories |
| review_date | date | Next review date |
| created_at | timestamp | |
| updated_at | timestamp | |
| created_by | uuid | FK to users |
| updated_by | uuid | FK to users |

### mitigations
| Column | Type | Description |
|--------|------|-------------|
| id | uuid | Primary key |
| risk_id | uuid | FK to risks |
| description | text | |
| status | enum | planned, in_progress, completed |
| due_date | date | |
| owner_id | uuid | FK to users |
| created_at | timestamp | |
| updated_at | timestamp | |

### frameworks
| Column | Type | Description |
|--------|------|-------------|
| id | uuid | Primary key |
| name | string | e.g., "ISO 27001", "SOC 2" |
| description | text | |

### risk_framework_controls
| Column | Type | Description |
|--------|------|-------------|
| risk_id | uuid | FK to risks |
| framework_id | uuid | FK to frameworks |
| control_ref | string | e.g., "A.12.1.1", "CC6.1" |
| notes | text | |

### audit_logs
| Column | Type | Description |
|--------|------|-------------|
| id | uuid | Primary key |
| entity_type | string | "risk", "mitigation" |
| entity_id | uuid | |
| action | enum | created, updated, deleted |
| changes | jsonb | {field: {old, new}} |
| user_id | uuid | FK to users |
| created_at | timestamp | |

## API Endpoints

### Auth (`/api/v1/auth`)
- `POST /register` - Create account
- `POST /login` - Authenticate, return JWT
- `POST /logout` - Invalidate session
- `GET /me` - Current user profile

### Risks (`/api/v1/risks`)
- `GET /` - List risks (paginated, filterable)
- `POST /` - Create risk
- `GET /:id` - Get single risk with relations
- `PUT /:id` - Update risk
- `DELETE /:id` - Soft delete risk

**Query params for risk list:**
- `status`, `severity`, `category_id`, `owner_id` - filters
- `search` - full-text search on title/description
- `review_before`, `review_after` - date range
- `sort`, `order` - sorting
- `page`, `limit` - pagination

### Mitigations (`/api/v1/risks/:riskId/mitigations`)
- `GET /` - List mitigations for a risk
- `POST /` - Create mitigation
- `PUT /:id` - Update mitigation
- `DELETE /:id` - Delete mitigation

### Categories (`/api/v1/categories`)
- `GET /` - List all categories
- `POST /` - Create category (admin only)
- `PUT /:id` - Update category (admin only)
- `DELETE /:id` - Delete category (admin only)

### Frameworks (`/api/v1/frameworks`)
- `GET /` - List all frameworks
- `POST /` - Create framework (admin only)

### Dashboard (`/api/v1/dashboard`)
- `GET /summary` - Counts by status, severity, category
- `GET /reviews/upcoming` - Risks with review_date in next N days
- `GET /reviews/overdue` - Risks past review_date

### AI (`/api/v1/ai`) - Stubbed
- `POST /summarize` - Stub: returns placeholder summary
- `POST /draft-mitigation` - Stub: returns placeholder text

### Response Format
```json
{
  "data": { ... },
  "meta": { "page": 1, "limit": 20, "total": 100 }
}
```

## Frontend Architecture

### Pages/Routes
```
/                    - Landing page (existing)
/login               - Login page (existing, wire to backend)
/register            - New user registration

/app                 - Protected area wrapper
/app/dashboard       - Overview with widgets
/app/risks           - Risk list view
/app/risks/:id       - Risk detail page
/app/risks/new       - Create risk form
/app/board           - Kanban board view
/app/calendar        - Review calendar view
/app/settings        - User/org settings
```

### State Management
- TanStack Query for server state
- React Hook Form + Zod for form validation
- Zustand for client-only state (sidebar toggle, filters)

### Key Components
- `RiskTable` - Sortable, filterable data table
- `RiskCard` - Card for board view
- `RiskForm` - Create/edit form with validation
- `MitigationList` - Nested list under risk detail
- `DashboardWidget` - Reusable stat card
- `CalendarView` - Month view with risk markers
- `FrameworkBadge` - Shows linked compliance controls
- `AuditTimeline` - Shows change history

### AI Integration (Stubbed)
- `SummarizeButton` - Calls `/api/v1/ai/summarize`
- `DraftButton` - Calls `/api/v1/ai/draft-mitigation`

## Build Order (Layered Slices)

### Slice 1: Authentication
- DB: users table, password hashing
- API: register, login, logout, me endpoints
- API: JWT middleware
- Frontend: wire login, add register, auth context
- **Deliverable:** Users can create accounts and log in

### Slice 2: Risks CRUD
- DB: risks, categories tables
- API: full risks CRUD, categories list
- Frontend: risk list, create/edit forms, detail page
- **Deliverable:** Full risk management

### Slice 3: Mitigations
- DB: mitigations table
- API: nested mitigations endpoints
- Frontend: mitigation list, add/edit forms
- **Deliverable:** Track mitigation actions

### Slice 4: Dashboard
- API: dashboard summary endpoints
- Frontend: dashboard page with widgets
- **Deliverable:** At-a-glance overview

### Slice 5: Review Calendar
- API: upcoming/overdue review endpoints
- Frontend: calendar view
- **Deliverable:** Visual review schedule

### Slice 6: Risk Board
- Frontend: Kanban board (columns = status)
- Drag-and-drop to change status
- **Deliverable:** Visual workflow management

### Slice 7: Framework Mapping
- DB: frameworks, risk_framework_controls tables
- API: frameworks list, link/unlink controls
- Frontend: framework badges, mapping UI
- **Deliverable:** Compliance traceability

### Slice 8: AI Integration (Stub)
- API: stubbed AI endpoints
- Frontend: summarize/draft buttons
- **Deliverable:** UI ready for real AI

### Slice 9: Audit Trail
- DB: audit_logs table
- Backend: middleware to log changes
- Frontend: audit timeline component
- **Deliverable:** Full change history

## Technical Stack

### Backend
- Go 1.22+ with Fiber v2
- PostgreSQL with pgx driver
- JWT auth (golang-jwt)
- Password hashing: bcrypt
- Migrations: golang-migrate
- Validation: go-playground/validator

### Frontend
- React 18 with TypeScript
- TanStack Router (existing)
- TanStack Query for data fetching
- TanStack Table for risk list
- React Hook Form + Zod
- shadcn/ui components (existing)
- Tailwind CSS (existing)

### Auth Flow
1. User registers/logs in → backend returns JWT
2. Frontend stores JWT (httpOnly cookie or localStorage)
3. Protected API calls include `Authorization: Bearer <token>`
4. Backend middleware validates JWT, injects user context

## File Structure (New Additions)

```
backend/
├── migrations/
│   ├── 001_users.up.sql
│   ├── 002_categories.up.sql
│   ├── 003_risks.up.sql
│   └── ...
├── internal/
│   ├── auth/           # JWT, password hashing
│   ├── models/         # Go structs
│   ├── handlers/       # HTTP handlers
│   └── middleware/     # Auth middleware

frontend/apps/web/src/
├── routes/
│   ├── app/
│   │   ├── index.tsx   # Dashboard
│   │   ├── risks/
│   │   ├── board.tsx
│   │   └── calendar.tsx
│   └── register.tsx
├── hooks/              # useAuth, etc.
├── lib/                # API client
└── components/         # Risk components
```

## Deployment

Local/Docker only for now:
- Docker Compose for PostgreSQL (existing)
- Backend runs on port 8080
- Frontend runs on port 3001
