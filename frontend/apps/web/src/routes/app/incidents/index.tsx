import { createFileRoute, Link } from "@tanstack/react-router";
import * as React from "react";

import { useIncidents, useIncidentCategories } from "@/hooks/useIncidents";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { cn } from "@/lib/utils";
import type { IncidentStatus, IncidentPriority } from "@/types/incident";

export const Route = createFileRoute("/app/incidents/")({
  component: IncidentsList,
});

const STATUS_COLORS: Record<IncidentStatus, string> = {
  new: "bg-green-100 text-green-800",
  acknowledged: "bg-yellow-100 text-yellow-800",
  in_progress: "bg-blue-100 text-blue-800",
  on_hold: "bg-gray-100 text-gray-800",
  resolved: "bg-emerald-100 text-emerald-800",
  closed: "bg-slate-100 text-slate-800",
};

const PRIORITY_COLORS: Record<IncidentPriority, string> = {
  p1: "bg-red-100 text-red-800",
  p2: "bg-orange-100 text-orange-800",
  p3: "bg-blue-100 text-blue-800",
  p4: "bg-gray-100 text-gray-800",
};

const PRIORITY_LABELS: Record<IncidentPriority, string> = {
  p1: "P1",
  p2: "P2",
  p3: "P3",
  p4: "P4",
};

function formatAge(createdAt: string): string {
  const created = new Date(createdAt);
  const now = new Date();
  const diffMs = now.getTime() - created.getTime();
  const diffMins = Math.floor(diffMs / 60000);
  const diffHours = Math.floor(diffMs / 3600000);
  const diffDays = Math.floor(diffMs / 86400000);

  if (diffMins < 1) return "Just now";
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 30) return `${diffDays}d ago`;
  return created.toLocaleDateString();
}

function IncidentsList() {
  const [status, setStatus] = React.useState<IncidentStatus | "">("");
  const [priority, setPriority] = React.useState<IncidentPriority | "">("");
  const [categoryId, setCategoryId] = React.useState<string>("");
  const [searchInput, setSearchInput] = React.useState("");
  const [searchQuery, setSearchQuery] = React.useState("");
  const [page, setPage] = React.useState(1);

  // Debounce search input (300ms delay)
  React.useEffect(() => {
    const timeoutId = setTimeout(() => {
      setSearchQuery(searchInput);
      setPage(1);
    }, 300);

    return () => clearTimeout(timeoutId);
  }, [searchInput]);

  const { data: categoriesData } = useIncidentCategories();
  const { data, isLoading } = useIncidents({
    status: status || undefined,
    priority: priority || undefined,
    category_id: categoryId || undefined,
    search: searchQuery || undefined,
    page,
    limit: 10,
  });

  const incidents = data?.data ?? [];
  const meta = data?.meta;
  const categories = categoriesData ?? [];

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Incidents</h1>
        <Link to="/app/incidents/new">
          <Button>New Incident</Button>
        </Link>
      </div>

      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Filters</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-5">
            <div className="grid gap-2">
              <Label>Search</Label>
              <Input
                placeholder="Search incidents..."
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
              />
            </div>
            <div className="grid gap-2">
              <Label>Status</Label>
              <Select
                value={status}
                onValueChange={(v) => {
                  setStatus(v as IncidentStatus | "");
                  setPage(1);
                }}
              >
                <SelectTrigger>
                  <SelectValue placeholder="All statuses" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">All statuses</SelectItem>
                  <SelectItem value="new">New</SelectItem>
                  <SelectItem value="acknowledged">Acknowledged</SelectItem>
                  <SelectItem value="in_progress">In Progress</SelectItem>
                  <SelectItem value="on_hold">On Hold</SelectItem>
                  <SelectItem value="resolved">Resolved</SelectItem>
                  <SelectItem value="closed">Closed</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid gap-2">
              <Label>Priority</Label>
              <Select
                value={priority}
                onValueChange={(v) => {
                  setPriority(v as IncidentPriority | "");
                  setPage(1);
                }}
              >
                <SelectTrigger>
                  <SelectValue placeholder="All priorities" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">All priorities</SelectItem>
                  <SelectItem value="p1">P1 - Critical</SelectItem>
                  <SelectItem value="p2">P2 - High</SelectItem>
                  <SelectItem value="p3">P3 - Medium</SelectItem>
                  <SelectItem value="p4">P4 - Low</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid gap-2">
              <Label>Category</Label>
              <Select
                value={categoryId}
                onValueChange={(v) => {
                  setCategoryId(v ?? "");
                  setPage(1);
                }}
              >
                <SelectTrigger>
                  <SelectValue placeholder="All categories" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">All categories</SelectItem>
                  {categories.map((cat) => (
                    <SelectItem key={cat.id} value={cat.id}>
                      {cat.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="grid gap-2">
              <Label>Results</Label>
              <div className="flex items-center h-10 text-sm text-muted-foreground">
                {meta ? `${meta.total} incident${meta.total !== 1 ? "s" : ""}` : "Loading..."}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-0">
          {isLoading ? (
            <div className="p-8 text-center text-muted-foreground">Loading...</div>
          ) : incidents.length === 0 ? (
            <div className="p-8 text-center text-muted-foreground">No incidents found</div>
          ) : (
            <table className="w-full">
              <thead>
                <tr className="border-b">
                  <th className="text-left p-4 font-medium">Priority</th>
                  <th className="text-left p-4 font-medium">Title</th>
                  <th className="text-left p-4 font-medium">Status</th>
                  <th className="text-left p-4 font-medium">Category</th>
                  <th className="text-left p-4 font-medium">Assignee</th>
                  <th className="text-left p-4 font-medium">Occurred</th>
                  <th className="text-left p-4 font-medium">Age</th>
                  <th className="text-right p-4 font-medium">Actions</th>
                </tr>
              </thead>
              <tbody>
                {incidents.map((incident) => (
                  <tr key={incident.id} className="border-b last:border-0 hover:bg-muted/50">
                    <td className="p-4">
                      <span
                        className={cn(
                          "px-2 py-1 rounded text-xs font-bold",
                          PRIORITY_COLORS[incident.priority]
                        )}
                      >
                        {PRIORITY_LABELS[incident.priority]}
                      </span>
                    </td>
                    <td className="p-4">
                      <Link
                        to="/app/incidents/$id"
                        params={{ id: incident.id }}
                        className="font-medium hover:underline"
                      >
                        {incident.title}
                      </Link>
                    </td>
                    <td className="p-4">
                      <span
                        className={cn(
                          "px-2 py-1 rounded text-xs font-medium",
                          STATUS_COLORS[incident.status]
                        )}
                      >
                        {incident.status.replace("_", " ")}
                      </span>
                    </td>
                    <td className="p-4 text-sm text-muted-foreground">
                      {incident.category?.name ?? "-"}
                    </td>
                    <td className="p-4 text-sm text-muted-foreground">
                      {incident.assignee?.name ?? "-"}
                    </td>
                    <td className="p-4 text-sm text-muted-foreground">
                      {incident.occurred_at
                        ? new Date(incident.occurred_at).toLocaleDateString()
                        : "-"}
                    </td>
                    <td className="p-4 text-sm text-muted-foreground">
                      {formatAge(incident.created_at)}
                    </td>
                    <td className="p-4 text-right">
                      <Link
                        to="/app/incidents/$id"
                        params={{ id: incident.id }}
                        className={cn(buttonVariants({ variant: "outline", size: "sm" }))}
                      >
                        View
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </CardContent>
      </Card>

      {meta && meta.total > meta.limit && (
        <div className="flex items-center justify-between mt-4">
          <div className="text-sm text-muted-foreground">
            Page {meta.page} of {Math.ceil(meta.total / meta.limit)}
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={page === 1}
              onClick={() => setPage(page - 1)}
            >
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              disabled={page >= Math.ceil(meta.total / meta.limit)}
              onClick={() => setPage(page + 1)}
            >
              Next
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
