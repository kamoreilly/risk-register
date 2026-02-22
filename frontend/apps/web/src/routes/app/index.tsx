import { createFileRoute, Link } from "@tanstack/react-router";

import { useDashboardSummary } from "@/hooks/useDashboard";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import type { CategoryCount } from "@/types/dashboard";

export const Route = createFileRoute("/app/")({
  component: Dashboard,
});

function Dashboard() {
  const { data, isLoading, error } = useDashboardSummary();

  if (isLoading) {
    return (
      <div className="p-8">
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <p className="mt-4 text-muted-foreground">Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-8">
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <p className="mt-4 text-destructive">
          Error loading dashboard: {error.message}
        </p>
      </div>
    );
  }

  const openCount = data?.by_status?.["open"] ?? 0;
  const mitigatingCount = data?.by_status?.["mitigating"] ?? 0;

  return (
    <div className="p-8">
      {/* Stat Cards Row */}
      <div className="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard title="Total Risks" value={data?.total_risks ?? 0} />
        <StatCard title="Open" value={openCount} />
        <StatCard title="Mitigating" value={mitigatingCount} />
        <StatCard title="Overdue Reviews" value={data?.overdue_reviews ?? 0} />
      </div>

      {/* Breakdown Section */}
      <div className="mb-8 grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Severity Breakdown */}
        <Card>
          <CardHeader>
            <CardTitle>By Severity</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <SeverityRow
                label="Critical"
                count={data?.by_severity?.["critical"] ?? 0}
                color="bg-red-500"
              />
              <SeverityRow
                label="High"
                count={data?.by_severity?.["high"] ?? 0}
                color="bg-orange-500"
              />
              <SeverityRow
                label="Medium"
                count={data?.by_severity?.["medium"] ?? 0}
                color="bg-yellow-500"
              />
              <SeverityRow
                label="Low"
                count={data?.by_severity?.["low"] ?? 0}
                color="bg-gray-500"
              />
            </div>
          </CardContent>
        </Card>

        {/* Category Breakdown */}
        <Card>
          <CardHeader>
            <CardTitle>By Category</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {data?.by_category && data.by_category.length > 0 ? (
                data.by_category.map((category: CategoryCount) => (
                  <CategoryRow
                    key={category.category_id}
                    name={category.category_name}
                    count={category.count}
                  />
                ))
              ) : (
                <p className="text-muted-foreground text-sm">No categories</p>
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Quick Links */}
      <div className="flex gap-4">
        <Link
          to="/app/risks"
          className="bg-primary text-primary-foreground inline-flex h-8 items-center justify-center gap-1.5 rounded-md px-2.5 text-xs font-medium transition-colors hover:bg-primary/80"
        >
          View All Risks
        </Link>
        <Link
          to="/app/risks"
          className="border-border bg-background hover:bg-muted hover:text-foreground inline-flex h-8 items-center justify-center gap-1.5 rounded-md border px-2.5 text-xs font-medium transition-colors"
        >
          View Board
        </Link>
      </div>
    </div>
  );
}

function StatCard({ title, value }: { title: string; value: number }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="text-3xl font-bold">{value}</div>
      </CardContent>
    </Card>
  );
}

function SeverityRow({
  label,
  count,
  color,
}: {
  label: string;
  count: number;
  color: string;
}) {
  return (
    <div className="flex items-center justify-between">
      <div className="flex items-center gap-2">
        <div className={`size-3 rounded-full ${color}`} />
        <span className="text-sm">{label}</span>
      </div>
      <span className="font-medium">{count}</span>
    </div>
  );
}

function CategoryRow({ name, count }: { name: string; count: number }) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-sm">{name}</span>
      <span className="font-medium">{count}</span>
    </div>
  );
}
