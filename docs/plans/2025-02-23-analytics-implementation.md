# Analytics Page Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a dedicated Analytics page with charts showing current risk state and trends over time.

**Architecture:** New backend `/api/v1/analytics` endpoint with pre-aggregated data. Frontend consumes endpoint via React Query and renders 4 chart components using shadcn/recharts.

**Tech Stack:** Go (Fiber), PostgreSQL, React, TanStack Router, TanStack Query, Recharts, shadcn/ui

---

## Task 1: Backend Analytics Models

**Files:**
- Create: `backend/internal/models/analytics.go`

**Step 1: Create the analytics models file**

```go
package models

// TimeDataPoint represents a single data point for time-series charts
type TimeDataPoint struct {
	Period string `json:"period"`
	Count  int    `json:"count"`
}

// StatusTimeDataPoint represents opened/closed counts for a period
type StatusTimeDataPoint struct {
	Period string `json:"period"`
	Open   int    `json:"open"`
	Closed int    `json:"closed"`
}

// AnalyticsResponse contains all analytics data for the frontend
type AnalyticsResponse struct {
	// Current State
	TotalRisks int               `json:"total_risks"`
	BySeverity map[string]int    `json:"by_severity"`
	ByStatus   map[string]int    `json:"by_status"`
	ByCategory []CategoryCount   `json:"by_category"`

	// Trends
	CreatedOverTime []TimeDataPoint       `json:"created_over_time"`
	StatusOverTime  []StatusTimeDataPoint `json:"status_over_time"`
}

// AnalyticsGranularity defines the time grouping for trend data
type AnalyticsGranularity string

const (
	GranularityMonthly AnalyticsGranularity = "monthly"
	GranularityWeekly  AnalyticsGranularity = "weekly"
)
```

**Step 2: Verify file compiles**

Run: `cd backend && go build ./...`
Expected: Success (no errors)

**Step 3: Commit**

```bash
git add backend/internal/models/analytics.go
git commit -m "feat(backend): add analytics models"
```

---

## Task 2: Backend Analytics Repository Interface

**Files:**
- Create: `backend/internal/database/analytics.go`

**Step 1: Create the analytics repository interface and implementation**

```go
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"backend/internal/models"
)

type AnalyticsRepository interface {
	GetAnalytics(ctx context.Context, granularity models.AnalyticsGranularity) (*models.AnalyticsResponse, error)
}

type analyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) GetAnalytics(ctx context.Context, granularity models.AnalyticsGranularity) (*models.AnalyticsResponse, error) {
	response := &models.AnalyticsResponse{
		BySeverity:      make(map[string]int),
		ByStatus:        make(map[string]int),
		ByCategory:      []models.CategoryCount{},
		CreatedOverTime: []models.TimeDataPoint{},
		StatusOverTime:  []models.StatusTimeDataPoint{},
	}

	// Get total count
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM risks").Scan(&response.TotalRisks); err != nil {
		return nil, fmt.Errorf("failed to get total risks: %w", err)
	}

	// Get counts by severity
	if err := r.populateByField(ctx, "severity", &response.BySeverity); err != nil {
		return nil, err
	}

	// Get counts by status
	if err := r.populateByField(ctx, "status", &response.ByStatus); err != nil {
		return nil, err
	}

	// Get counts by category
	if err := r.populateByCategory(ctx, &response.ByCategory); err != nil {
		return nil, err
	}

	// Get created over time
	if err := r.populateCreatedOverTime(ctx, granularity, &response.CreatedOverTime); err != nil {
		return nil, err
	}

	// Get status over time
	if err := r.populateStatusOverTime(ctx, granularity, &response.StatusOverTime); err != nil {
		return nil, err
	}

	return response, nil
}

func (r *analyticsRepository) populateByField(ctx context.Context, field string, target *map[string]int) error {
	query := fmt.Sprintf("SELECT %s, COUNT(*) FROM risks GROUP BY %s", field, field)
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to get counts by %s: %w", field, err)
	}
	defer rows.Close()

	for rows.Next() {
		var value string
		var count int
		if err := rows.Scan(&value, &count); err != nil {
			return err
		}
		(*target)[value] = count
	}
	return rows.Err()
}

func (r *analyticsRepository) populateByCategory(ctx context.Context, target *[]models.CategoryCount) error {
	query := `
		SELECT c.id, c.name, COUNT(r.id)
		FROM categories c
		LEFT JOIN risks r ON r.category_id = c.id
		GROUP BY c.id, c.name
		ORDER BY COUNT(r.id) DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to get counts by category: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cc models.CategoryCount
		if err := rows.Scan(&cc.CategoryID, &cc.CategoryName, &cc.Count); err != nil {
			return err
		}
		*target = append(*target, cc)
	}
	return rows.Err()
}

func (r *analyticsRepository) populateCreatedOverTime(ctx context.Context, granularity models.AnalyticsGranularity, target *[]models.TimeDataPoint) error {
	var dateFormat string
	if granularity == models.GranularityWeekly {
		dateFormat = "YYYY-\"W\"WW"
	} else {
		dateFormat = "YYYY-MM"
	}

	query := fmt.Sprintf(`
		SELECT TO_CHAR(created_at, '%s') as period, COUNT(*) as count
		FROM risks
		WHERE created_at >= NOW() - INTERVAL '12 months'
		GROUP BY period
		ORDER BY period ASC
	`, dateFormat)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to get created over time: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dp models.TimeDataPoint
		if err := rows.Scan(&dp.Period, &dp.Count); err != nil {
			return err
		}
		*target = append(*target, dp)
	}
	return rows.Err()
}

func (r *analyticsRepository) populateStatusOverTime(ctx context.Context, granularity models.AnalyticsGranularity, target *[]models.StatusTimeDataPoint) error {
	var dateFormat string
	if granularity == models.GranularityWeekly {
		dateFormat = "YYYY-\"W\"WW"
	} else {
		dateFormat = "YYYY-MM"
	}

	// Get opened (created) per period
	openedQuery := fmt.Sprintf(`
		SELECT TO_CHAR(created_at, '%s') as period, COUNT(*) as count
		FROM risks
		WHERE created_at >= NOW() - INTERVAL '12 months'
		GROUP BY period
	`, dateFormat)

	openedCounts := make(map[string]int)
	rows, err := r.db.QueryContext(ctx, openedQuery)
	if err != nil {
		return fmt.Errorf("failed to get opened counts: %w", err)
	}
	for rows.Next() {
		var period string
		var count int
		if err := rows.Scan(&period, &count); err != nil {
			rows.Close()
			return err
		}
		openedCounts[period] = count
	}
	rows.Close()

	// Get closed (status changed to resolved/accepted) per period
	closedQuery := fmt.Sprintf(`
		SELECT TO_CHAR(updated_at, '%s') as period, COUNT(*) as count
		FROM risks
		WHERE status IN ('resolved', 'accepted')
		  AND updated_at >= NOW() - INTERVAL '12 months'
		GROUP BY period
	`, dateFormat)

	closedCounts := make(map[string]int)
	rows, err = r.db.QueryContext(ctx, closedQuery)
	if err != nil {
		return fmt.Errorf("failed to get closed counts: %w", err)
	}
	for rows.Next() {
		var period string
		var count int
		if err := rows.Scan(&period, &count); err != nil {
			rows.Close()
			return err
		}
		closedCounts[period] = count
	}
	rows.Close()

	// Merge all periods
	allPeriods := make(map[string]bool)
	for p := range openedCounts {
		allPeriods[p] = true
	}
	for p := range closedCounts {
		allPeriods[p] = true
	}

	// Build sorted result
	for period := range allPeriods {
		*target = append(*target, models.StatusTimeDataPoint{
			Period: period,
			Open:   openedCounts[period],
			Closed: closedCounts[period],
		})
	}

	// Sort by period
	sortByPeriod(*target)

	return nil
}

func sortByPeriod(data []models.StatusTimeDataPoint) {
	for i := 0; i < len(data)-1; i++ {
		for j := i + 1; j < len(data); j++ {
			if data[i].Period > data[j].Period {
				data[i], data[j] = data[j], data[i]
			}
		}
	}
}

// Helper to get raw time for sorting
func parsePeriod(period string) time.Time {
	t, _ := time.Parse("2006-W01", period)
	if t.IsZero() {
		t, _ = time.Parse("2006-01", period)
	}
	return t
}

func init() {
	_ = parsePeriod
}
```

**Step 2: Verify file compiles**

Run: `cd backend && go build ./...`
Expected: Success (no errors)

**Step 3: Commit**

```bash
git add backend/internal/database/analytics.go
git commit -m "feat(backend): add analytics repository"
```

---

## Task 3: Backend Analytics Handler

**Files:**
- Create: `backend/internal/handlers/analytics.go`

**Step 1: Create the analytics handler**

```go
package handlers

import (
	"backend/internal/database"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	repo database.AnalyticsRepository
}

func NewAnalyticsHandler(repo database.AnalyticsRepository) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo}
}

// Get returns all analytics data
func (h *AnalyticsHandler) Get(c *fiber.Ctx) error {
	granularity := models.AnalyticsGranularity(c.Query("granularity", "monthly"))
	if granularity != models.GranularityMonthly && granularity != models.GranularityWeekly {
		granularity = models.GranularityMonthly
	}

	response, err := h.repo.GetAnalytics(c.Context(), granularity)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch analytics"})
	}

	return c.JSON(response)
}
```

**Step 2: Verify file compiles**

Run: `cd backend && go build ./...`
Expected: Success (no errors)

**Step 3: Commit**

```bash
git add backend/internal/handlers/analytics.go
git commit -m "feat(backend): add analytics handler"
```

---

## Task 4: Backend Handler Tests

**Files:**
- Create: `backend/internal/handlers/analytics_test.go`

**Step 1: Write the failing test**

```go
package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type mockAnalyticsRepo struct {
	response *models.AnalyticsResponse
	err      error
}

func (m *mockAnalyticsRepo) GetAnalytics(ctx context.Context, granularity models.AnalyticsGranularity) (*models.AnalyticsResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func TestAnalyticsHandler_Get(t *testing.T) {
	mockResp := &models.AnalyticsResponse{
		TotalRisks: 10,
		BySeverity: map[string]int{"critical": 2, "high": 3, "medium": 3, "low": 2},
		ByStatus:   map[string]int{"open": 6, "mitigating": 2, "resolved": 2},
		ByCategory: []models.CategoryCount{
			{CategoryID: "1", CategoryName: "Security", Count: 5},
			{CategoryID: "2", CategoryName: "Compliance", Count: 3},
		},
		CreatedOverTime: []models.TimeDataPoint{
			{Period: "2024-01", Count: 2},
			{Period: "2024-02", Count: 3},
		},
		StatusOverTime: []models.StatusTimeDataPoint{
			{Period: "2024-01", Open: 2, Closed: 1},
			{Period: "2024-02", Open: 3, Closed: 2},
		},
	}

	app := fiber.New()
	handler := NewAnalyticsHandler(&mockAnalyticsRepo{response: mockResp})
	app.Get("/analytics", handler.Get)

	tests := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{"default granularity", "/analytics", 200},
		{"monthly granularity", "/analytics?granularity=monthly", 200},
		{"weekly granularity", "/analytics?granularity=weekly", 200},
		{"invalid granularity defaults to monthly", "/analytics?granularity=invalid", 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			var response models.AnalyticsResponse
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.TotalRisks != 10 {
				t.Errorf("expected 10 total risks, got %d", response.TotalRisks)
			}
		})
	}
}
```

**Step 2: Run the test**

Run: `cd backend && go test ./internal/handlers -run TestAnalyticsHandler -v`
Expected: PASS

**Step 3: Commit**

```bash
git add backend/internal/handlers/analytics_test.go
git commit -m "test(backend): add analytics handler tests"
```

---

## Task 5: Register Analytics Route

**Files:**
- Modify: `backend/internal/server/server.go`
- Modify: `backend/internal/server/routes.go`

**Step 1: Add analytics handler and repository to FiberServer struct**

In `backend/internal/server/server.go`, add to the `FiberServer` struct:

```go
analyticsHandler   *handlers.AnalyticsHandler
```

And in the `New()` function, add after `dashboard :=` line:

```go
analytics := database.NewAnalyticsRepository(rawDB)
```

And in the server initialization:

```go
analyticsHandler:  handlers.NewAnalyticsHandler(analytics),
```

**Step 2: Add the analytics route in `routes.go`**

Add after the dashboard routes group:

```go
// Analytics routes
protected.Get("/analytics", s.analyticsHandler.Get)
```

**Step 3: Verify build**

Run: `cd backend && go build ./...`
Expected: Success (no errors)

**Step 4: Run all handler tests**

Run: `cd backend && go test ./internal/handlers -v`
Expected: All tests PASS

**Step 5: Commit**

```bash
git add backend/internal/server/server.go backend/internal/server/routes.go
git commit -m "feat(backend): register analytics route"
```

---

## Task 6: Frontend Analytics Types

**Files:**
- Create: `frontend/apps/web/src/types/analytics.ts`

**Step 1: Create the analytics types file**

```typescript
export interface TimeDataPoint {
  period: string;
  count: number;
}

export interface StatusTimeDataPoint {
  period: string;
  open: number;
  closed: number;
}

export interface AnalyticsResponse {
  total_risks: number;
  by_severity: Record<string, number>;
  by_status: Record<string, number>;
  by_category: CategoryCount[];
  created_over_time: TimeDataPoint[];
  status_over_time: StatusTimeDataPoint[];
}

export type Granularity = 'monthly' | 'weekly';
```

**Step 2: Verify TypeScript compiles**

Run: `cd frontend && bun run check-types`
Expected: Success (no errors)

**Step 3: Commit**

```bash
git add frontend/apps/web/src/types/analytics.ts
git commit -m "feat(frontend): add analytics types"
```

---

## Task 7: Frontend useAnalytics Hook

**Files:**
- Create: `frontend/apps/web/src/hooks/useAnalytics.ts`

**Step 1: Create the analytics hook**

```typescript
import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { AnalyticsResponse, Granularity } from '@/types/analytics';

export const ANALYTICS_KEY = ['analytics'];

export function useAnalytics(granularity: Granularity = 'monthly') {
  return useQuery({
    queryKey: [...ANALYTICS_KEY, granularity],
    queryFn: () =>
      api.get<AnalyticsResponse>(`/api/v1/analytics?granularity=${granularity}`),
  });
}
```

**Step 2: Verify TypeScript compiles**

Run: `cd frontend && bun run check-types`
Expected: Success (no errors)

**Step 3: Commit**

```bash
git add frontend/apps/web/src/hooks/useAnalytics.ts
git commit -m "feat(frontend): add useAnalytics hook"
```

---

## Task 8: Severity Radial Chart Component

**Files:**
- Create: `frontend/apps/web/src/components/analytics/severity-radial-chart.tsx`

**Step 1: Create the severity radial chart**

```tsx
import { Pie, PieChart } from "recharts";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";

interface SeverityRadialChartProps {
  data: Record<string, number>;
}

const chartConfig = {
  critical: {
    label: "Critical",
    color: "hsl(var(--destructive))",
  },
  high: {
    label: "High",
    color: "hsl(25, 95%, 53%)",
  },
  medium: {
    label: "Medium",
    color: "hsl(45, 93%, 47%)",
  },
  low: {
    label: "Low",
    color: "hsl(var(--muted-foreground))",
  },
} satisfies ChartConfig;

export function SeverityRadialChart({ data }: SeverityRadialChartProps) {
  const chartData = Object.entries(data).map(([severity, count]) => ({
    severity,
    count,
    fill: chartConfig[severity as keyof typeof chartConfig]?.color ?? "hsl(var(--muted))",
  }));

  const total = Object.values(data).reduce((sum, count) => sum + count, 0);

  return (
    <Card className="flex flex-col">
      <CardHeader className="items-center pb-0">
        <CardTitle>Risk Severity</CardTitle>
        <CardDescription>Distribution by severity level</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 pb-0">
        <ChartContainer
          config={chartConfig}
          className="mx-auto aspect-square max-h-[250px]"
        >
          <PieChart>
            <ChartTooltip
              content={<ChartTooltipContent nameKey="severity" hideLabel />}
            />
            <Pie
              data={chartData}
              dataKey="count"
              nameKey="severity"
              innerRadius={60}
              strokeWidth={5}
            />
          </PieChart>
        </ChartContainer>
        <div className="mt-4 flex flex-wrap justify-center gap-4 text-xs">
          {Object.entries(chartConfig).map(([key, config]) => (
            <div key={key} className="flex items-center gap-1.5">
              <div
                className="size-2.5 shrink-0 rounded-[2px]"
                style={{ backgroundColor: config.color }}
              />
              <span className="text-muted-foreground">{config.label}</span>
              <span className="font-medium">{data[key] ?? 0}</span>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
```

**Step 2: Verify TypeScript compiles**

Run: `cd frontend && bun run check-types`
Expected: Success (no errors)

**Step 3: Commit**

```bash
git add frontend/apps/web/src/components/analytics/severity-radial-chart.tsx
git commit -m "feat(frontend): add severity radial chart component"
```

---

## Task 9: Category Bar Chart Component

**Files:**
- Create: `frontend/apps/web/src/components/analytics/category-bar-chart.tsx`

**Step 1: Create the category bar chart**

```tsx
import { Bar, BarChart, XAxis, YAxis } from "recharts";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";
import type { CategoryCount } from "@/types/dashboard";

interface CategoryBarChartProps {
  data: CategoryCount[];
}

const chartConfig = {
  count: {
    label: "Risks",
    color: "hsl(var(--primary))",
  },
} satisfies ChartConfig;

export function CategoryBarChart({ data }: CategoryBarChartProps) {
  const chartData = data.map((cat) => ({
    name: cat.category_name,
    count: cat.count,
  }));

  return (
    <Card className="flex flex-col">
      <CardHeader>
        <CardTitle>Risks by Category</CardTitle>
        <CardDescription>Number of risks in each category</CardDescription>
      </CardHeader>
      <CardContent className="flex-1">
        {chartData.length === 0 ? (
          <div className="flex h-[200px] items-center justify-center text-muted-foreground">
            No categories available
          </div>
        ) : (
          <ChartContainer config={chartConfig} className="min-h-[200px] w-full">
            <BarChart
              layout="vertical"
              data={chartData}
              margin={{ left: 80 }}
            >
              <XAxis type="number" />
              <YAxis
                dataKey="name"
                type="category"
                tickLine={false}
                axisLine={false}
                width={70}
              />
              <ChartTooltip content={<ChartTooltipContent />} />
              <Bar dataKey="count" fill="hsl(var(--primary))" radius={4} />
            </BarChart>
          </ChartContainer>
        )}
      </CardContent>
    </Card>
  );
}
```

**Step 2: Verify TypeScript compiles**

Run: `cd frontend && bun run check-types`
Expected: Success (no errors)

**Step 3: Commit**

```bash
git add frontend/apps/web/src/components/analytics/category-bar-chart.tsx
git commit -m "feat(frontend): add category bar chart component"
```

---

## Task 10: Created Over Time Area Chart Component

**Files:**
- Create: `frontend/apps/web/src/components/analytics/created-over-time-chart.tsx`

**Step 1: Create the created over time chart**

```tsx
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from "recharts";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";
import type { TimeDataPoint } from "@/types/analytics";

interface CreatedOverTimeChartProps {
  data: TimeDataPoint[];
}

const chartConfig = {
  count: {
    label: "Risks Created",
    color: "hsl(var(--primary))",
  },
} satisfies ChartConfig;

export function CreatedOverTimeChart({ data }: CreatedOverTimeChartProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Risks Created Over Time</CardTitle>
        <CardDescription>Number of new risks created per period</CardDescription>
      </CardHeader>
      <CardContent>
        {data.length === 0 ? (
          <div className="flex h-[200px] items-center justify-center text-muted-foreground">
            No data available
          </div>
        ) : (
          <ChartContainer config={chartConfig} className="min-h-[200px] w-full">
            <AreaChart
              data={data}
              margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
            >
              <CartesianGrid strokeDasharray="3 3" vertical={false} />
              <XAxis
                dataKey="period"
                tickLine={false}
                axisLine={false}
                tickMargin={8}
              />
              <YAxis tickLine={false} axisLine={false} tickMargin={8} />
              <ChartTooltip content={<ChartTooltipContent />} />
              <Area
                type="monotone"
                dataKey="count"
                stroke="hsl(var(--primary))"
                fill="hsl(var(--primary) / 0.2)"
                strokeWidth={2}
              />
            </AreaChart>
          </ChartContainer>
        )}
      </CardContent>
    </Card>
  );
}
```

**Step 2: Verify TypeScript compiles**

Run: `cd frontend && bun run check-types`
Expected: Success (no errors)

**Step 3: Commit**

```bash
git add frontend/apps/web/src/components/analytics/created-over-time-chart.tsx
git commit -m "feat(frontend): add created over time chart component"
```

---

## Task 11: Status Over Time Bar Chart Component

**Files:**
- Create: `frontend/apps/web/src/components/analytics/status-over-time-chart.tsx`

**Step 1: Create the status over time chart**

```tsx
import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";
import type { StatusTimeDataPoint } from "@/types/analytics";

interface StatusOverTimeChartProps {
  data: StatusTimeDataPoint[];
}

const chartConfig = {
  open: {
    label: "Opened",
    color: "hsl(var(--primary))",
  },
  closed: {
    label: "Closed",
    color: "hsl(142, 76%, 36%)",
  },
} satisfies ChartConfig;

export function StatusOverTimeChart({ data }: StatusOverTimeChartProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Opened vs Closed Over Time</CardTitle>
        <CardDescription>Risks opened and closed per period</CardDescription>
      </CardHeader>
      <CardContent>
        {data.length === 0 ? (
          <div className="flex h-[200px] items-center justify-center text-muted-foreground">
            No data available
          </div>
        ) : (
          <ChartContainer config={chartConfig} className="min-h-[200px] w-full">
            <BarChart data={data}>
              <CartesianGrid strokeDasharray="3 3" vertical={false} />
              <XAxis
                dataKey="period"
                tickLine={false}
                axisLine={false}
                tickMargin={8}
              />
              <YAxis tickLine={false} axisLine={false} tickMargin={8} />
              <ChartTooltip content={<ChartTooltipContent />} />
              <Bar
                dataKey="open"
                fill="hsl(var(--primary))"
                radius={[4, 4, 0, 0]}
              />
              <Bar
                dataKey="closed"
                fill="hsl(142, 76%, 36%)"
                radius={[4, 4, 0, 0]}
              />
            </BarChart>
          </ChartContainer>
        )}
      </CardContent>
    </Card>
  );
}
```

**Step 2: Verify TypeScript compiles**

Run: `cd frontend && bun run check-types`
Expected: Success (no errors)

**Step 3: Commit**

```bash
git add frontend/apps/web/src/components/analytics/status-over-time-chart.tsx
git commit -m "feat(frontend): add status over time chart component"
```

---

## Task 12: Analytics Page Route

**Files:**
- Create: `frontend/apps/web/src/routes/app/analytics.tsx`

**Step 1: Create the analytics page**

```tsx
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";

import { useAnalytics } from "@/hooks/useAnalytics";
import { SeverityRadialChart } from "@/components/analytics/severity-radial-chart";
import { CategoryBarChart } from "@/components/analytics/category-bar-chart";
import { CreatedOverTimeChart } from "@/components/analytics/created-over-time-chart";
import { StatusOverTimeChart } from "@/components/analytics/status-over-time-chart";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import type { Granularity } from "@/types/analytics";

export const Route = createFileRoute("/app/analytics")({
  component: AnalyticsPage,
});

function AnalyticsPage() {
  const [granularity, setGranularity] = useState<Granularity>("monthly");
  const { data, isLoading, error } = useAnalytics(granularity);

  if (error) {
    return (
      <div className="p-8">
        <h1 className="text-2xl font-bold">Analytics</h1>
        <p className="mt-4 text-destructive">
          Error loading analytics: {error.message}
        </p>
      </div>
    );
  }

  return (
    <div className="p-8">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Analytics</h1>
          <p className="text-muted-foreground">
            Comprehensive risk metrics and trends
          </p>
        </div>
        <Select
          value={granularity}
          onValueChange={(value) => setGranularity(value as Granularity)}
        >
          <SelectTrigger className="w-[140px]">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="monthly">Monthly</SelectItem>
            <SelectItem value="weekly">Weekly</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Current State Section */}
      <div className="mb-8">
        <h2 className="mb-4 text-lg font-semibold">Current State</h2>
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
          {isLoading ? (
            <>
              <ChartSkeleton />
              <ChartSkeleton />
            </>
          ) : (
            <>
              <SeverityRadialChart data={data?.by_severity ?? {}} />
              <CategoryBarChart data={data?.by_category ?? []} />
            </>
          )}
        </div>
      </div>

      {/* Trends Section */}
      <div>
        <h2 className="mb-4 text-lg font-semibold">Trends Over Time</h2>
        <div className="grid grid-cols-1 gap-6">
          {isLoading ? (
            <>
              <ChartSkeleton />
              <ChartSkeleton />
            </>
          ) : (
            <>
              <CreatedOverTimeChart data={data?.created_over_time ?? []} />
              <StatusOverTimeChart data={data?.status_over_time ?? []} />
            </>
          )}
        </div>
      </div>
    </div>
  );
}

function ChartSkeleton() {
  return (
    <div className="rounded-lg border p-6">
      <Skeleton className="mb-4 h-6 w-32" />
      <Skeleton className="mb-2 h-4 w-48" />
      <Skeleton className="h-[200px] w-full" />
    </div>
  );
}
```

**Step 2: Verify TypeScript compiles**

Run: `cd frontend && bun run check-types`
Expected: Success (no errors)

**Step 3: Commit**

```bash
git add frontend/apps/web/src/routes/app/analytics.tsx
git commit -m "feat(frontend): add analytics page route"
```

---

## Task 13: Add Navigation Link to Analytics

**Files:**
- Modify: `frontend/apps/web/src/routes/app/route.tsx` (or navigation component)

**Step 1: Find and update navigation to include Analytics link**

Look for the navigation/sidebar component in your app and add a link to `/app/analytics`. This will vary based on your navigation structure.

If using a sidebar or navbar, add:
```tsx
<Link to="/app/analytics">Analytics</Link>
```

Or if there's a navigation items array, add:
```typescript
{ name: "Analytics", href: "/app/analytics" }
```

**Step 2: Verify TypeScript compiles**

Run: `cd frontend && bun run check-types`
Expected: Success (no errors)

**Step 3: Commit**

```bash
git add <modified-navigation-file>
git commit -m "feat(frontend): add analytics navigation link"
```

---

## Task 14: Integration Testing

**Step 1: Start the backend server**

Run: `cd backend && make watch`
Wait for server to start on port 8080

**Step 2: In a new terminal, start the frontend**

Run: `cd frontend && bun run dev`
Wait for server to start on port 3001

**Step 3: Test the analytics page manually**

1. Navigate to http://localhost:3001/app/analytics
2. Verify the page loads without errors
3. Verify charts render correctly
4. Test switching between Monthly and Weekly granularity
5. Check network tab for API call to `/api/v1/analytics`

**Step 4: Run backend tests**

Run: `cd backend && make test`
Expected: All tests PASS

**Step 5: Commit any fixes if needed**

If any fixes were required during testing:
```bash
git add .
git commit -m "fix: resolve analytics integration issues"
```

---

## Summary

### Backend Changes:
- `backend/internal/models/analytics.go` - Analytics data models
- `backend/internal/database/analytics.go` - Repository with SQL queries
- `backend/internal/handlers/analytics.go` - HTTP handler
- `backend/internal/handlers/analytics_test.go` - Unit tests
- `backend/internal/server/server.go` - Handler registration
- `backend/internal/server/routes.go` - Route registration

### Frontend Changes:
- `frontend/apps/web/src/types/analytics.ts` - TypeScript types
- `frontend/apps/web/src/hooks/useAnalytics.ts` - React Query hook
- `frontend/apps/web/src/components/analytics/*.tsx` - 4 chart components
- `frontend/apps/web/src/routes/app/analytics.tsx` - Main page
- Navigation file - Analytics link

### API Endpoint:
- `GET /api/v1/analytics?granularity=monthly|weekly`
