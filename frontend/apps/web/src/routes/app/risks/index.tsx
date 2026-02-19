import { createFileRoute, Link } from "@tanstack/react-router";
import * as React from "react";

import { useRisks } from "@/hooks/useRisks";
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
import type { RiskStatus, RiskSeverity } from "@/types/risk";

export const Route = createFileRoute("/app/risks/")({
  component: RisksList,
});

const STATUS_COLORS: Record<RiskStatus, string> = {
  open: "bg-yellow-100 text-yellow-800",
  mitigating: "bg-blue-100 text-blue-800",
  resolved: "bg-green-100 text-green-800",
  accepted: "bg-gray-100 text-gray-800",
};

const SEVERITY_COLORS: Record<RiskSeverity, string> = {
  low: "bg-gray-100 text-gray-800",
  medium: "bg-yellow-100 text-yellow-800",
  high: "bg-orange-100 text-orange-800",
  critical: "bg-red-100 text-red-800",
};

function RisksList() {
  const [status, setStatus] = React.useState<RiskStatus | "">("");
  const [severity, setSeverity] = React.useState<RiskSeverity | "">("");
  const [search, setSearch] = React.useState("");
  const [page, setPage] = React.useState(1);

  const { data, isLoading } = useRisks({
    status: status || undefined,
    severity: severity || undefined,
    search: search || undefined,
    page,
    limit: 10,
  });

  const risks = data?.data ?? [];
  const meta = data?.meta;

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold">Risks</h1>
          <p className="text-muted-foreground">Manage your organization's risks</p>
        </div>
        <Link to="/app/risks/new">
          <Button>New Risk</Button>
        </Link>
      </div>

      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Filters</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-4">
            <div className="grid gap-2">
              <Label>Search</Label>
              <Input
                placeholder="Search risks..."
                value={search}
                onChange={(e) => {
                  setSearch(e.target.value);
                  setPage(1);
                }}
              />
            </div>
            <div className="grid gap-2">
              <Label>Status</Label>
              <Select value={status} onValueChange={(v) => { setStatus(v as RiskStatus | ""); setPage(1); }}>
                <SelectTrigger>
                  <SelectValue placeholder="All statuses" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">All statuses</SelectItem>
                  <SelectItem value="open">Open</SelectItem>
                  <SelectItem value="mitigating">Mitigating</SelectItem>
                  <SelectItem value="resolved">Resolved</SelectItem>
                  <SelectItem value="accepted">Accepted</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid gap-2">
              <Label>Severity</Label>
              <Select value={severity} onValueChange={(v) => { setSeverity(v as RiskSeverity | ""); setPage(1); }}>
                <SelectTrigger>
                  <SelectValue placeholder="All severities" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">All severities</SelectItem>
                  <SelectItem value="low">Low</SelectItem>
                  <SelectItem value="medium">Medium</SelectItem>
                  <SelectItem value="high">High</SelectItem>
                  <SelectItem value="critical">Critical</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid gap-2">
              <Label>Results</Label>
              <div className="flex items-center h-10 text-sm text-muted-foreground">
                {meta ? `${meta.total} risk${meta.total !== 1 ? 's' : ''}` : 'Loading...'}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-0">
          {isLoading ? (
            <div className="p-8 text-center text-muted-foreground">Loading...</div>
          ) : risks.length === 0 ? (
            <div className="p-8 text-center text-muted-foreground">No risks found</div>
          ) : (
            <table className="w-full">
              <thead>
                <tr className="border-b">
                  <th className="text-left p-4 font-medium">Title</th>
                  <th className="text-left p-4 font-medium">Status</th>
                  <th className="text-left p-4 font-medium">Severity</th>
                  <th className="text-left p-4 font-medium">Review Date</th>
                  <th className="text-right p-4 font-medium">Actions</th>
                </tr>
              </thead>
              <tbody>
                {risks.map((risk) => (
                  <tr key={risk.id} className="border-b last:border-0 hover:bg-muted/50">
                    <td className="p-4">
                      <Link to="/app/risks/$id" params={{ id: risk.id }} className="font-medium hover:underline">
                        {risk.title}
                      </Link>
                    </td>
                    <td className="p-4">
                      <span className={cn("px-2 py-1 rounded text-xs font-medium", STATUS_COLORS[risk.status])}>
                        {risk.status}
                      </span>
                    </td>
                    <td className="p-4">
                      <span className={cn("px-2 py-1 rounded text-xs font-medium", SEVERITY_COLORS[risk.severity])}>
                        {risk.severity}
                      </span>
                    </td>
                    <td className="p-4 text-sm text-muted-foreground">
                      {risk.review_date ? new Date(risk.review_date).toLocaleDateString() : '-'}
                    </td>
                    <td className="p-4 text-right">
                      <Link
                        to="/app/risks/$id"
                        params={{ id: risk.id }}
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
