import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { 
  TrendingUp, 
  Activity, 
  AlertTriangle, 
  ShieldCheck, 
  LayoutDashboard,
  Calendar,
  Filter
} from "lucide-react";

import { useAnalytics } from "@/hooks/useAnalytics";
import { SeverityRadialChart } from "@/components/analytics/severity-radial-chart";
import { CategoryBarChart } from "@/components/analytics/category-bar-chart";
import { CreatedOverTimeChart } from "@/components/analytics/created-over-time-chart";
import { StatusOverTimeChart } from "@/components/analytics/status-over-time-chart";
import { ProgressReport } from "@/components/analytics/progress-report";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { 
  Card, 
  CardContent, 
  CardHeader, 
  CardTitle, 
  CardDescription 
} from "@/components/ui/card";
import { StatCard } from "@/components/ui/stat-card";
import type { Granularity } from "@/types/analytics";

export const Route = createFileRoute("/app/analytics")({
  component: AnalyticsPage,
});

function AnalyticsPage() {
  const [granularity, setGranularity] = useState<Granularity>("monthly");
  const { data, isLoading, error } = useAnalytics(granularity);

  if (isLoading) {
    return (
      <div className="p-8 space-y-8 animate-pulse">
        <div className="flex justify-end gap-3">
          <div className="h-10 w-40 bg-muted rounded" />
          <div className="h-10 w-[140px] bg-muted rounded" />
        </div>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {[1, 2, 3, 4].map(i => <Skeleton key={i} className="h-32 rounded-xl" />)}
        </div>
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <ChartSkeleton />
          <ChartSkeleton />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex h-[80vh] flex-col items-center justify-center p-8 text-center">
        <div className="bg-destructive/10 text-destructive mb-4 rounded-full p-4">
          <AlertTriangle size={48} />
        </div>
        <h1 className="text-2xl font-bold">Analytics Error</h1>
        <p className="mt-2 text-muted-foreground">
          {error.message || "Failed to load risk analytics data."}
        </p>
      </div>
    );
  }

  const totalRisks = data?.total_risks ?? 0;
  const criticalCount = data?.by_severity?.["critical"] ?? 0;
  const highCount = data?.by_severity?.["high"] ?? 0;
  const mitCount = data?.by_status?.["mitigating"] ?? 0;

  return (
    <div className="p-8 space-y-8">
      {/* Header and Controls */}
      <div className="flex flex-col md:flex-row md:items-center justify-end gap-4">
        <div className="flex items-center gap-3 self-end md:self-auto">
          <div className="flex items-center gap-2 text-sm text-muted-foreground bg-muted p-1 rounded-md px-3 h-10">
            <Calendar className="size-4" />
            <span>Range: {granularity === 'monthly' ? 'Last 12 Months' : 'Last 12 Weeks'}</span>
          </div>
          
          <Select
            value={granularity}
            onValueChange={(value) => setGranularity(value as Granularity)}
          >
            <SelectTrigger className="w-[140px] h-10 border-primary/20">
              <Filter className="size-4 mr-2 opacity-50" />
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="monthly">Monthly</SelectItem>
              <SelectItem value="weekly">Weekly</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Summary Stat Cards */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard 
          title="Total Active Risks" 
          value={totalRisks} 
          icon={<Activity className="size-4" />}
          description="Current inventory size"
        />
        <StatCard 
          title="Critical Exposure" 
          value={criticalCount} 
          icon={<AlertTriangle className="text-red-500 size-4" />}
          description="Requires top-level oversight"
          trend={criticalCount > 0 ? "High" : "Low"}
        />
        <StatCard 
          title="High Likelihood" 
          value={highCount} 
          icon={<TrendingUp className="text-orange-500 size-4" />}
          description="Common risk threshold"
        />
        <StatCard 
          title="Mitigation Rate" 
          value={`${totalRisks > 0 ? ((mitCount / totalRisks) * 100).toFixed(0) : 0}%`}
          icon={<ShieldCheck className="text-green-500 size-4" />}
          description="Active control coverage"
        />
      </div>

      {/* Main Content Sections */}
      <div className="grid grid-cols-1 gap-8">
        
        {/* Current State Section */}
        <div className="space-y-4">
          <div className="flex items-center gap-2 border-b pb-2">
            <LayoutDashboard className="size-5 text-muted-foreground" />
            <h2 className="text-xl font-semibold">Inventory Breakdown</h2>
          </div>
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            <Card className="flex flex-col">
              <SeverityRadialChart data={data?.by_severity ?? {}} />
              <CardContent className="pt-0 border-t mt-4 pt-6 space-y-4">
                <ProgressReport label="Critical Risks" value={criticalCount} total={totalRisks} color="bg-red-500" />
                <ProgressReport label="High Risks" value={highCount} total={totalRisks} color="bg-orange-500" />
              </CardContent>
            </Card>
            <CategoryBarChart data={data?.by_category ?? []} />
          </div>
        </div>

        {/* Trends Section */}
        <div className="space-y-4">
          <div className="flex items-center gap-2 border-b pb-2">
            <TrendingUp className="size-5 text-muted-foreground" />
            <h2 className="text-xl font-semibold">Temporal Trends</h2>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <CreatedOverTimeChart data={data?.created_over_time ?? []} />
            <StatusOverTimeChart data={data?.status_over_time ?? []} />
          </div>
        </div>
      </div>
    </div>
  );
}

function ChartSkeleton() {
  return (
    <div className="rounded-xl border p-6 bg-card/50">
      <div className="flex justify-between items-center mb-6">
        <div className="space-y-2">
          <Skeleton className="h-5 w-40" />
          <Skeleton className="h-4 w-64" />
        </div>
        <Skeleton className="h-10 w-10 rounded-full" />
      </div>
      <Skeleton className="h-[250px] w-full" />
    </div>
  );
}
