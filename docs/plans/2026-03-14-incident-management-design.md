# Incident Management & Reporting - Design Document

**Date:** 2026-03-14
**Status:** Approved

## Overview

Add incident management capabilities to the risk register application. Incidents are adverse events that may or may not be linked to existing risks. This feature includes full lifecycle tracking, reporting, and a new responder role.

## Requirements Summary

- Generic events that can optionally link to risks (many-to-many)
- 6-state workflow: New → Acknowledged → In Progress → On Hold → Resolved → Closed
- IT priority classification (P1-P4)
- Comprehensive tracking: core fields, timestamps, impact/cause, resolution
- Reporting: summary stats, SLA metrics, trend analysis, export
- New 'responder' role for incident management
- Audit log tracking (existing pattern)
- Separate incident categories
- Single assignee model

## Architecture Approach

**Parallel Entity Model** - Incidents as a first-class entity alongside risks, sharing patterns but independent.

Rationale: Incidents and risks are fundamentally different concepts with different lifecycles, fields, and reporting needs. This approach provides clean separation while reusing established patterns.

---

## Data Model

### New Tables

#### `incident_categories`
```sql
- id (UUID, PK)
- name (VARCHAR, unique) - e.g., "Outage", "Breach", "Error", "External"
- description (TEXT)
- created_at (TIMESTAMP)
```

#### `incidents`
```sql
- id (UUID, PK)
- title (VARCHAR, not null)
- description (TEXT)
- category_id (UUID, FK → incident_categories.id, nullable)
- priority (ENUM: 'p1', 'p2', 'p3', 'p4')
- status (ENUM: 'new', 'acknowledged', 'in_progress', 'on_hold', 'resolved', 'closed')
- assignee_id (UUID, FK → users.id, nullable)
- reporter_id (UUID, FK → users.id)
- service_affected (VARCHAR) - which system/service
- root_cause (TEXT, nullable)
- resolution_notes (TEXT, nullable)
- occurred_at (TIMESTAMP) - when it happened
- detected_at (TIMESTAMP) - when we found out
- resolved_at (TIMESTAMP, nullable)
- created_at, updated_at (TIMESTAMP)
- created_by, updated_by (UUID, FK → users.id)
```

#### `incident_risks` (junction table)
```sql
- id (UUID, PK)
- incident_id (UUID, FK → incidents.id, cascade delete)
- risk_id (UUID, FK → risks.id, cascade delete)
- created_at (TIMESTAMP)
- created_by (UUID, FK → users.id)
- UNIQUE(incident_id, risk_id)
```

### Schema Changes

Add `responder` to the existing `user_role` enum.

---

## API Endpoints

### Incident Categories (Admin only)
```
GET    /api/v1/incident-categories        List categories
POST   /api/v1/incident-categories        Create (admin)
PUT    /api/v1/incident-categories/:id    Update (admin)
DELETE /api/v1/incident-categories/:id    Delete (admin)
```

### Incidents
```
GET    /api/v1/incidents                  List (filter: status, priority, category, assignee, search)
POST   /api/v1/incidents                  Create (responder+)
GET    /api/v1/incidents/:id              Get by ID
PUT    /api/v1/incidents/:id              Update (responder+, assignee can update own)
DELETE /api/v1/incidents/:id              Delete (admin only)
```

### Incident-Risk Links
```
GET    /api/v1/incidents/:incidentId/risks           List linked risks
POST   /api/v1/incidents/:incidentId/risks           Link risk
DELETE /api/v1/incidents/:incidentId/risks/:riskId   Unlink risk
```

### Reverse Lookup
```
GET    /api/v1/risks/:riskId/incidents     List incidents for a risk
```

### Audit (existing pattern)
```
GET    /api/v1/incidents/:incidentId/audit  Get audit logs
```

### Reporting
```
GET    /api/v1/incidents/report/summary    Summary stats (by status, priority, category)
GET    /api/v1/incidents/report/sla        SLA metrics (MTTR, time to acknowledge)
GET    /api/v1/incidents/report/trends     Trend data (incidents over time, recurring)
GET    /api/v1/incidents/export            Export to CSV/JSON
```

---

## Frontend

### Routes
```
/app/incidents                    Incidents list
/app/incidents/new                Create incident
/app/incidents/:id                View/edit incident
/app/incidents/report             Reporting dashboard
```

### Pages

#### Incidents List (`/app/incidents`)
- Data table: Priority, Title, Status, Category, Assignee, Occurred, Age
- Filters: status, priority, category, assignee, date range, search
- Sort: priority, created, occurred_at
- Export button

#### Incident Detail (`/app/incidents/:id`)
- Header: Title, priority badge, status badge
- Two-column layout:
  - Left: Description, service affected, root cause, resolution notes
  - Right: Metadata (reporter, assignee, category, timestamps)
- Timeline: Audit log entries
- Linked risks: Cards with add/remove capability

#### Incident Report (`/app/incidents/report`)
- Summary cards: Total, by status, by priority
- Charts: Incidents over time, by category, MTTR trend
- SLA table: Average acknowledge time, average resolve time (by priority)
- Export controls

### Components
- `IncidentPriorityBadge` - P1-P4 with colors
- `IncidentStatusBadge` - 6 states with colors
- `IncidentForm` - Create/edit form
- `IncidentFilters` - Filter bar
- `IncidentRiskLink` - Link/unlink risks component

---

## Permissions

### Role Hierarchy
```
admin     → Full access
responder → Create, view, update incidents
member    → View incidents only
```

### Permission Matrix

| Action | Admin | Responder | Member |
|--------|-------|-----------|--------|
| View incidents | ✓ | ✓ | ✓ |
| Create incidents | ✓ | ✓ | ✗ |
| Update incidents | ✓ | ✓ (own + assigned) | ✗ |
| Delete incidents | ✓ | ✗ | ✗ |
| Link/unlink risks | ✓ | ✓ | ✗ |
| View reports | ✓ | ✓ | ✓ |
| Export data | ✓ | ✓ | ✗ |
| Manage incident categories | ✓ | ✗ | ✗ |

---

## Implementation Phases

### Phase 1: Backend Foundation
1. Database migrations (tables, enum extension)
2. Go models and database layer
3. Incident categories CRUD handlers
4. Incidents CRUD handlers
5. Incident-risk linking handlers
6. Audit log integration

### Phase 2: Backend Reporting
1. Summary stats endpoint
2. SLA metrics endpoint
3. Trends endpoint
4. Export endpoint

### Phase 3: Frontend Core
1. Incident list page with filters
2. Incident detail page
3. Incident form (create/edit)
4. Risk linking UI

### Phase 4: Frontend Reporting
1. Report page layout
2. Summary cards
3. Charts (using existing chart library)
4. Export functionality

### Phase 5: Integration & Polish
1. Navigation updates
2. Dashboard integration (recent incidents)
3. Error handling
4. Testing
