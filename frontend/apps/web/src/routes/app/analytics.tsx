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
      <div className="mb-6 flex items-center justify-end">
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
