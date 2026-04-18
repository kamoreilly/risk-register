import { createFileRoute, Link } from "@tanstack/react-router";
import {
  ChevronRight,
  Activity,
  AlertTriangle,
  Calendar,
  LayoutDashboard,
  ListChecks
} from "lucide-react";
import {
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend
} from "recharts";

import { useDashboardSummary, useUpcomingReviews, useOverdueReviews } from "@/hooks/useDashboard";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { StatCard } from "@/components/ui/stat-card";
import { ProgressReport } from "@/components/analytics/progress-report";
import type { ReviewRisk } from "@/types/dashboard";

export const Route = createFileRoute("/app/")({
  component: Dashboard,
});

const SEVERITY_COLORS = {
  critical: "#ef4444",
  high: "#f97316",
  medium: "#eab308",
  low: "#6b7280",
} as const;

function Dashboard() {
  const { data, isLoading: isSummaryLoading, error: summaryError } = useDashboardSummary();
  const { data: reviewsData, isLoading: isReviewsLoading } = useUpcomingReviews(30);
  const { data: overdueData, isLoading: isOverdueLoading } = useOverdueReviews();

  const isLoading = isSummaryLoading || isReviewsLoading || isOverdueLoading;

  if (isLoading) {
    return (
      <div className="p-8 space-y-8 animate-pulse">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {[1, 2, 3, 4].map(i => <div key={i} className="h-32 bg-muted rounded-lg" />)}
        </div>
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
          <Card className="lg:col-span-2 h-80" />
          <Card className="h-80" />
        </div>
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <Card className="h-80" />
          <Card className="h-80" />
        </div>
      </div>
    );
  }

  if (summaryError) {
    return (
      <div className="flex h-[80vh] flex-col items-center justify-center p-8 text-center">
        <div className="bg-destructive/10 text-destructive mb-4 rounded-full p-4">
          <AlertTriangle size={48} />
        </div>
        <h1 className="text-2xl font-bold">Dashboard Error</h1>
        <p className="mt-2 text-muted-foreground max-w-md">
          {summaryError.message || "Something went wrong while loading dashboard data."}
        </p>
        <button
          onClick={() => window.location.reload()}
          className="bg-primary text-primary-foreground mt-6 rounded-md px-4 py-2 text-sm font-medium"
        >
          Retry Connection
        </button>
      </div>
    );
  }

  const totalRisks = data?.total_risks ?? 0;
  const criticalCount = data?.by_severity?.["critical"] ?? 0;
  const highCount = data?.by_severity?.["high"] ?? 0;
  const mediumCount = data?.by_severity?.["medium"] ?? 0;
  const lowCount = data?.by_severity?.["low"] ?? 0;
  const overdueCount = overdueData?.risks?.length ?? 0;

  const severityData = [
    { name: "Critical", value: criticalCount, fill: SEVERITY_COLORS.critical },
    { name: "High", value: highCount, fill: SEVERITY_COLORS.high },
    { name: "Medium", value: mediumCount, fill: SEVERITY_COLORS.medium },
    { name: "Low", value: lowCount, fill: SEVERITY_COLORS.low },
  ].filter(d => d.value > 0);

  const categoryData = (data?.by_category ?? []).map(cat => ({
    name: cat.category_name,
    count: cat.count
  })).sort((a, b) => b.count - a.count);

  return (
    <div className="p-8 space-y-8">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard 
          title="Total Active Risks" 
          value={totalRisks} 
          icon={<Activity className="size-4" />}
          description="Current inventory size"
        />
        <StatCard 
          title="Critical Risks" 
          value={criticalCount} 
          icon={<AlertTriangle className="text-red-500 size-4" />}
          description="Requires immediate attention"
          isUrgent={criticalCount > 0}
        />
        <StatCard 
          title="Overdue Reviews" 
          value={overdueCount} 
          icon={<Calendar className="text-orange-500 size-4" />}
          description="Past review due date"
          isUrgent={overdueCount > 0}
        />
        <StatCard 
          title="Upcoming Reviews" 
          value={reviewsData?.risks?.length ?? 0} 
          icon={<ListChecks className="text-blue-500 size-4" />}
          description="Due in next 30 days"
        />
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle className="text-base">Risks by Category</CardTitle>
            <CardDescription>Top functional areas affected</CardDescription>
          </CardHeader>
          <CardContent className="h-64">
            {categoryData.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={categoryData.slice(0, 8)} layout="vertical">
                  <CartesianGrid strokeDasharray="3 3" horizontal={true} vertical={false} opacity={0.3} />
                  <XAxis type="number" hide />
                  <YAxis
                    dataKey="name"
                    type="category"
                    width={120}
                    fontSize={12}
                    tickLine={false}
                    axisLine={false}
                  />
                  <Tooltip cursor={{ fill: 'transparent' }} />
                  <Bar
                    dataKey="count"
                    fill="hsl(var(--primary))"
                    radius={[0, 4, 4, 0]}
                    barSize={20}
                  />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex h-full flex-col items-center justify-center text-center">
                <p className="text-muted-foreground text-sm">No categories defined</p>
              </div>
            )}
          </CardContent>
        </Card>

        <Card className="lg:col-span-1">
          <CardHeader>
            <CardTitle className="text-base">Risk Severity Distribution</CardTitle>
            <CardDescription>Breakdown by impact level</CardDescription>
          </CardHeader>
          <CardContent className="h-64 flex flex-col justify-center">
            {severityData.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={severityData}
                    cx="50%"
                    cy="50%"
                    innerRadius={50}
                    outerRadius={70}
                    paddingAngle={5}
                    dataKey="value"
                  >
                    {severityData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.fill} />
                    ))}
                  </Pie>
                  <Tooltip
                    contentStyle={{ borderRadius: '8px', border: 'none', boxShadow: '0 4px 12px rgba(0,0,0,0.1)' }}
                  />
                  <Legend />
                </PieChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex h-full flex-col items-center justify-center text-center">
                <p className="text-muted-foreground text-sm">No data available</p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <div className="flex items-center gap-2">
              <div className="bg-primary/10 p-1.5 rounded-md">
                <Calendar className="size-4 text-primary" />
              </div>
              <div>
                <CardTitle className="text-base">Upcoming Reviews</CardTitle>
                <CardDescription>Risks due for review in next 30 days</CardDescription>
              </div>
            </div>
            <Link to="/app/risks" className="text-primary text-xs font-medium hover:underline flex items-center">
              View All <ChevronRight size={14} />
            </Link>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {reviewsData?.risks && reviewsData.risks.length > 0 ? (
                reviewsData.risks.slice(0, 5).map((risk: ReviewRisk) => (
                  <div key={risk.id} className="flex items-center justify-between group">
                    <div className="flex flex-col">
                      <span className="text-sm font-medium leading-none group-hover:text-primary transition-colors cursor-pointer capitalize">
                        {risk.title}
                      </span>
                      <span className="text-muted-foreground mt-1 text-xs">
                        Review: {new Date(risk.review_date).toLocaleDateString()}
                      </span>
                    </div>
                    <Badge variant={risk.severity === 'critical' || risk.severity === 'high' ? 'destructive' : 'secondary'} className="capitalize">
                      {risk.severity}
                    </Badge>
                  </div>
                ))
              ) : (
                <div className="flex h-32 flex-col items-center justify-center border-dashed border-2 rounded-lg border-muted">
                  <span className="text-muted-foreground text-sm italic">No upcoming reviews</span>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <div className="flex items-center gap-2">
              <div className="bg-primary/10 p-1.5 rounded-md">
                <LayoutDashboard className="size-4 text-primary" />
              </div>
              <div>
                <CardTitle className="text-base">Risk Distribution</CardTitle>
                <CardDescription>Severity breakdown overview</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-6 pt-4">
            <ProgressReport label="Critical" value={criticalCount} total={totalRisks} color="bg-red-500" />
            <ProgressReport label="High" value={highCount} total={totalRisks} color="bg-orange-500" />
            <ProgressReport label="Medium" value={mediumCount} total={totalRisks} color="bg-yellow-500" />
            <ProgressReport label="Low" value={lowCount} total={totalRisks} color="bg-gray-500" />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
