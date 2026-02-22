# Analytics Page Design

## Overview

A dedicated Analytics page for visualizing risk management metrics using interactive charts. The page displays current state distributions and time-based trend analysis.

## Requirements

- Display current risk state: severity distribution, category breakdown
- Display trend analysis: risks created over time, opened vs closed over time
- Configurable time granularity (weekly/monthly)
- No filters for initial implementation

## Backend API

### Endpoint: `GET /api/analytics`

**Query Parameters:**
- `granularity`: `monthly` (default) | `weekly`

**Response:**

```go
type AnalyticsResponse struct {
    // Current State
    TotalRisks     int                     `json:"total_risks"`
    BySeverity     map[string]int          `json:"by_severity"`
    ByStatus       map[string]int          `json:"by_status"`
    ByCategory     []CategoryCount         `json:"by_category"`

    // Trends
    CreatedOverTime []TimeDataPoint        `json:"created_over_time"`
    StatusOverTime  []StatusTimeDataPoint  `json:"status_over_time"`
}

type TimeDataPoint struct {
    Period string `json:"period"`
    Count  int    `json:"count"`
}

type StatusTimeDataPoint struct {
    Period  string `json:"period"`
    Opened  int    `json:"opened"`
    Closed  int    `json:"closed"`
}
```

**Database Queries:**
- Current state: Reuse existing dashboard queries
- Created over time: `GROUP BY DATE_TRUNC('month'/'week', created_at)`
- Status over time: Created count per period, closed = risks where `status='closed'` and `updated_at` in period

## Frontend Structure

### Route

`/app/analytics` - `frontend/apps/web/src/routes/app/analytics.tsx`

### Page Layout

```
┌─────────────────────────────────────────────────────────────┐
│  Analytics                                    [Monthly ▼]   │
│  Comprehensive risk metrics and trends                      │
├─────────────────────────────────────────────────────────────┤
│  CURRENT STATE                                              │
│  ┌─────────────────────────┐  ┌─────────────────────────┐   │
│  │   Radial Chart          │  │   Bar Chart             │   │
│  │   (Severity dist.)      │  │   (By Category)         │   │
│  └─────────────────────────┘  └─────────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│  TRENDS OVER TIME                                           │
│  ┌─────────────────────────────────────────────────────┐    │
│  │   Area Chart (Risks Created Over Time)              │    │
│  └─────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │   Bar Chart (Opened vs Closed Over Time)            │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### Components

| Component | Path | Purpose |
|-----------|------|---------|
| `AnalyticsPage` | `routes/app/analytics.tsx` | Main page, data fetching, layout |
| `SeverityRadialChart` | `components/analytics/` | Radial chart for severity |
| `CategoryBarChart` | `components/analytics/` | Horizontal bar for categories |
| `CreatedOverTimeChart` | `components/analytics/` | Area chart for creation trends |
| `StatusOverTimeChart` | `components/analytics/` | Bar chart for opened/closed |
| `GranularitySelect` | `components/analytics/` | Monthly/Weekly toggle dropdown |

### Data Flow

```
Backend (Go) ──GET /api/analytics──► React Query (useAnalytics)
                                           │
                                           ▼
                                    Chart Components
```

### Libraries

- `@tanstack/react-query` - Data fetching
- `recharts` - Chart primitives
- `@/components/ui/chart.tsx` - shadcn chart wrapper
- `@/components/ui/select.tsx` - Granularity dropdown

### Chart Colors

| Severity | Color |
|----------|-------|
| Critical | Destructive (red) |
| High | Warning (orange) |
| Medium | Chart-2 (yellow) |
| Low | Chart-1 (muted) |

### States

- **Loading:** Skeleton placeholders for each chart card
- **Error:** Error message with retry button
- **Empty:** "No data available" message in chart areas

## Implementation Approach

Dedicated Analytics API (Approach 2) - New backend endpoint with pre-aggregated data for better performance and scalability.
