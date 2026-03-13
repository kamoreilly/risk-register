import { createFileRoute, Link } from "@tanstack/react-router";
import * as React from "react";

import {
  useIncident,
  useUpdateIncident,
  useDeleteIncident,
  useIncidentCategories,
  useIncidentRisks,
  useLinkIncidentRisk,
  useUnlinkIncidentRisk,
  useIncidentAuditLogs,
} from "@/hooks/useIncidents";
import { useRisks } from "@/hooks/useRisks";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
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
import { toast } from "sonner";
import type { IncidentStatus, IncidentPriority } from "@/types/incident";
import type { AuditAction } from "@/types/audit";

export const Route = createFileRoute("/app/incidents/$id")({
  component: IncidentDetail,
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
  p1: "P1 - Critical",
  p2: "P2 - High",
  p3: "P3 - Medium",
  p4: "P4 - Low",
};

function IncidentDetail() {
  const { id } = Route.useParams();
  const navigate = Route.useNavigate();

  const { data: incident, isLoading } = useIncident(id);
  const { data: categories } = useIncidentCategories();
  const updateIncident = useUpdateIncident(id);
  const deleteIncident = useDeleteIncident();

  // Risk linking hooks
  const { data: linkedRisks, isLoading: risksLoading } = useIncidentRisks(id);
  const linkRisk = useLinkIncidentRisk(id);
  const unlinkRisk = useUnlinkIncidentRisk(id);
  const { data: availableRisks } = useRisks({ limit: 100 });

  // Audit log hooks
  const { data: auditLogs, isLoading: auditLogsLoading } = useIncidentAuditLogs(id);

  const [isEditing, setIsEditing] = React.useState(false);
  const [title, setTitle] = React.useState("");
  const [description, setDescription] = React.useState("");
  const [status, setStatus] = React.useState<IncidentStatus>("new");
  const [priority, setPriority] = React.useState<IncidentPriority>("p3");
  const [categoryId, setCategoryId] = React.useState<string>("");
  const [serviceAffected, setServiceAffected] = React.useState("");
  const [rootCause, setRootCause] = React.useState("");
  const [resolutionNotes, setResolutionNotes] = React.useState("");
  const [occurredAt, setOccurredAt] = React.useState("");
  const [detectedAt, setDetectedAt] = React.useState("");

  // Risk linking state
  const [isAddingRisk, setIsAddingRisk] = React.useState(false);
  const [selectedRiskId, setSelectedRiskId] = React.useState("");

  React.useEffect(() => {
    if (incident) {
      setTitle(incident.title);
      setDescription(incident.description || "");
      setStatus(incident.status);
      setPriority(incident.priority);
      setCategoryId(incident.category_id || "");
      setServiceAffected(incident.service_affected || "");
      setRootCause(incident.root_cause || "");
      setResolutionNotes(incident.resolution_notes || "");
      setOccurredAt(incident.occurred_at ? incident.occurred_at.slice(0, 16) : "");
      setDetectedAt(incident.detected_at ? incident.detected_at.slice(0, 16) : "");
    }
  }, [incident]);

  const handleSave = async () => {
    try {
      await updateIncident.mutateAsync({
        title,
        description: description || undefined,
        status,
        priority,
        category_id: categoryId || undefined,
        service_affected: serviceAffected || undefined,
        root_cause: rootCause || undefined,
        resolution_notes: resolutionNotes || undefined,
        occurred_at: occurredAt || undefined,
        detected_at: detectedAt || undefined,
      });
      toast.success("Incident updated");
      setIsEditing(false);
    } catch (error) {
      toast.error("Failed to update incident");
    }
  };

  const handleDelete = async () => {
    if (!confirm("Are you sure you want to delete this incident?")) return;

    try {
      await deleteIncident.mutateAsync(id);
      toast.success("Incident deleted");
      navigate({ to: "/app/incidents" });
    } catch (error) {
      toast.error("Failed to delete incident");
    }
  };

  const handleLinkRisk = async () => {
    if (!selectedRiskId) {
      toast.error("Please select a risk to link");
      return;
    }

    try {
      await linkRisk.mutateAsync({ risk_id: selectedRiskId });
      toast.success("Risk linked");
      setSelectedRiskId("");
      setIsAddingRisk(false);
    } catch (error) {
      toast.error("Failed to link risk");
    }
  };

  const handleUnlinkRisk = async (riskId: string) => {
    if (!confirm("Are you sure you want to unlink this risk?")) return;

    try {
      await unlinkRisk.mutateAsync(riskId);
      toast.success("Risk unlinked");
    } catch (error) {
      toast.error("Failed to unlink risk");
    }
  };

  if (isLoading) {
    return (
      <div className="p-8">
        <div className="text-center text-muted-foreground">Loading...</div>
      </div>
    );
  }

  if (!incident) {
    return (
      <div className="p-8">
        <div className="text-center text-muted-foreground">Incident not found</div>
        <div className="text-center mt-4">
          <Link
            to="/app/incidents"
            className={cn(buttonVariants({ variant: "outline" }))}
          >
            Back to Incidents
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <div>
          <Link
            to="/app/incidents"
            className="text-sm text-muted-foreground hover:underline mb-2 block"
          >
            &larr; Back to Incidents
          </Link>
          <h1 className="text-2xl font-bold">
            {isEditing ? "Edit Incident" : incident.title}
          </h1>
        </div>
        <div className="flex gap-2">
          {isEditing ? (
            <>
              <Button variant="outline" onClick={() => setIsEditing(false)}>
                Cancel
              </Button>
              <Button onClick={handleSave} disabled={updateIncident.isPending}>
                Save
              </Button>
            </>
          ) : (
            <>
              <Button variant="outline" onClick={() => setIsEditing(true)}>
                Edit
              </Button>
              <Button
                variant="destructive"
                onClick={handleDelete}
                disabled={deleteIncident.isPending}
              >
                Delete
              </Button>
            </>
          )}
        </div>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center gap-2">
            <span
              className={cn(
                "px-2 py-1 rounded text-xs font-medium",
                STATUS_COLORS[incident.status],
              )}
            >
              {incident.status.replace("_", " ")}
            </span>
            <span
              className={cn(
                "px-2 py-1 rounded text-xs font-medium",
                PRIORITY_COLORS[incident.priority],
              )}
            >
              {PRIORITY_LABELS[incident.priority]}
            </span>
          </div>
        </CardHeader>
        <CardContent className="space-y-6">
          {isEditing ? (
            <>
              <div className="grid gap-2">
                <Label htmlFor="title">Title</Label>
                <Input
                  id="title"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="description">Description</Label>
                <textarea
                  id="description"
                  className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label>Status</Label>
                  <Select
                    value={status}
                    onValueChange={(v) => setStatus(v as IncidentStatus)}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
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
                    onValueChange={(v) => setPriority(v as IncidentPriority)}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="p1">P1 - Critical</SelectItem>
                      <SelectItem value="p2">P2 - High</SelectItem>
                      <SelectItem value="p3">P3 - Medium</SelectItem>
                      <SelectItem value="p4">P4 - Low</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label>Category</Label>
                  <Select
                    value={categoryId}
                    onValueChange={(value) => setCategoryId(value ?? "")}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select category" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="">None</SelectItem>
                      {categories?.map((cat) => (
                        <SelectItem key={cat.id} value={cat.id}>
                          {cat.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="serviceAffected">Service Affected</Label>
                  <Input
                    id="serviceAffected"
                    value={serviceAffected}
                    onChange={(e) => setServiceAffected(e.target.value)}
                    placeholder="e.g., API, Database, Frontend"
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="occurredAt">Occurred At</Label>
                  <Input
                    id="occurredAt"
                    type="datetime-local"
                    value={occurredAt}
                    onChange={(e) => setOccurredAt(e.target.value)}
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="detectedAt">Detected At</Label>
                  <Input
                    id="detectedAt"
                    type="datetime-local"
                    value={detectedAt}
                    onChange={(e) => setDetectedAt(e.target.value)}
                  />
                </div>
              </div>
              <div className="grid gap-2">
                <Label htmlFor="rootCause">Root Cause</Label>
                <textarea
                  id="rootCause"
                  className="flex min-h-[60px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value={rootCause}
                  onChange={(e) => setRootCause(e.target.value)}
                  placeholder="What caused this incident?"
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="resolutionNotes">Resolution Notes</Label>
                <textarea
                  id="resolutionNotes"
                  className="flex min-h-[60px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value={resolutionNotes}
                  onChange={(e) => setResolutionNotes(e.target.value)}
                  placeholder="How was this incident resolved?"
                />
              </div>
            </>
          ) : (
            <>
              <div>
                <h3 className="font-medium mb-2">Description</h3>
                <p className="text-muted-foreground">
                  {incident.description || "No description"}
                </p>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <h3 className="font-medium mb-2">Category</h3>
                  <p className="text-muted-foreground">
                    {incident.category?.name || "Uncategorized"}
                  </p>
                </div>
                <div>
                  <h3 className="font-medium mb-2">Service Affected</h3>
                  <p className="text-muted-foreground">
                    {incident.service_affected || "Not specified"}
                  </p>
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <h3 className="font-medium mb-2">Reporter</h3>
                  <p className="text-muted-foreground">
                    {incident.reporter?.name || "Unknown"}
                  </p>
                </div>
                <div>
                  <h3 className="font-medium mb-2">Assignee</h3>
                  <p className="text-muted-foreground">
                    {incident.assignee?.name || "Unassigned"}
                  </p>
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <h3 className="font-medium mb-2">Occurred At</h3>
                  <p className="text-muted-foreground">
                    {incident.occurred_at
                      ? new Date(incident.occurred_at).toLocaleString()
                      : "Not specified"}
                  </p>
                </div>
                <div>
                  <h3 className="font-medium mb-2">Detected At</h3>
                  <p className="text-muted-foreground">
                    {incident.detected_at
                      ? new Date(incident.detected_at).toLocaleString()
                      : "Not specified"}
                  </p>
                </div>
              </div>
              {incident.resolved_at && (
                <div>
                  <h3 className="font-medium mb-2">Resolved At</h3>
                  <p className="text-muted-foreground">
                    {new Date(incident.resolved_at).toLocaleString()}
                  </p>
                </div>
              )}
              <div>
                <h3 className="font-medium mb-2">Root Cause</h3>
                <p className="text-muted-foreground">
                  {incident.root_cause || "Not specified"}
                </p>
              </div>
              <div>
                <h3 className="font-medium mb-2">Resolution Notes</h3>
                <p className="text-muted-foreground">
                  {incident.resolution_notes || "Not specified"}
                </p>
              </div>
              <div className="text-sm text-muted-foreground pt-4 border-t">
                <p>Created: {new Date(incident.created_at).toLocaleString()}</p>
                <p>Updated: {new Date(incident.updated_at).toLocaleString()}</p>
              </div>
            </>
          )}
        </CardContent>
      </Card>

      {/* Linked Risks Section */}
      <Card className="mt-6">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Linked Risks</CardTitle>
              <CardDescription>
                Risks related to this incident
              </CardDescription>
            </div>
            {!isAddingRisk && (
              <Button onClick={() => setIsAddingRisk(true)}>
                Link Risk
              </Button>
            )}
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Add Risk Form */}
          {isAddingRisk && (
            <div className="border rounded-lg p-4 space-y-4 bg-muted/50">
              <h4 className="font-medium">Link New Risk</h4>
              <div className="grid gap-2">
                <Label htmlFor="risk-select">Select Risk</Label>
                <Select
                  value={selectedRiskId}
                  onValueChange={(value) => setSelectedRiskId(value ?? "")}
                >
                  <SelectTrigger id="risk-select">
                    <SelectValue placeholder="Select a risk to link" />
                  </SelectTrigger>
                  <SelectContent>
                    {(availableRisks?.data ?? [])
                      .filter(
                        (risk) =>
                          !(linkedRisks ?? []).some((lr) => lr.risk_id === risk.id)
                      )
                      .map((risk) => (
                        <SelectItem key={risk.id} value={risk.id}>
                          {risk.title}
                        </SelectItem>
                      ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="flex gap-2 justify-end">
                <Button variant="outline" onClick={() => {
                  setIsAddingRisk(false);
                  setSelectedRiskId("");
                }}>
                  Cancel
                </Button>
                <Button
                  onClick={handleLinkRisk}
                  disabled={linkRisk.isPending || !selectedRiskId}
                >
                  Link
                </Button>
              </div>
            </div>
          )}

          {/* Risks List */}
          {risksLoading ? (
            <div className="text-center text-muted-foreground py-4">
              Loading linked risks...
            </div>
          ) : linkedRisks && linkedRisks.length > 0 ? (
            <div className="flex flex-wrap gap-2">
              {linkedRisks.map((linkedRisk) => (
                <div
                  key={linkedRisk.id}
                  className="inline-flex items-center gap-2 px-3 py-1.5 bg-purple-100 text-purple-800 rounded-full text-sm font-medium group"
                >
                  <Link
                    to="/app/risks/$id"
                    params={{ id: linkedRisk.risk_id }}
                    className="hover:underline"
                  >
                    {linkedRisk.risk?.title ?? `Risk ${linkedRisk.risk_id}`}
                  </Link>
                  <button
                    onClick={() => handleUnlinkRisk(linkedRisk.risk_id)}
                    disabled={unlinkRisk.isPending}
                    className="ml-1 text-purple-600 hover:text-purple-900 opacity-0 group-hover:opacity-100 transition-opacity"
                    title="Unlink risk"
                  >
                    &times;
                  </button>
                </div>
              ))}
            </div>
          ) : (
            !isAddingRisk && (
              <div className="text-center text-muted-foreground py-8">
                No risks linked. Click "Link Risk" to add a related risk.
              </div>
            )
          )}
        </CardContent>
      </Card>

      {/* Audit History Section */}
      <Card className="mt-6">
        <CardHeader>
          <CardTitle>Audit History</CardTitle>
          <CardDescription>
            Timeline of changes made to this incident
          </CardDescription>
        </CardHeader>
        <CardContent>
          {auditLogsLoading ? (
            <div className="text-center text-muted-foreground py-4">
              Loading audit history...
            </div>
          ) : auditLogs && auditLogs.length > 0 ? (
            <div className="space-y-4">
              {auditLogs.map((log, index) => {
                const isLast = index === auditLogs.length - 1;
                const actionColorMap: Record<AuditAction, string> = {
                  created: "bg-green-500",
                  updated: "bg-blue-500",
                  deleted: "bg-red-500",
                };

                return (
                  <div key={log.id} className="flex gap-3">
                    <div className="flex flex-col items-center">
                      <div
                        className={cn(
                          "w-3 h-3 rounded-full",
                          actionColorMap[log.action],
                        )}
                      />
                      {!isLast && <div className="w-px flex-1 bg-border" />}
                    </div>
                    <div className={cn("pb-4", isLast && "pb-0")}>
                      <p className="text-sm">
                        <span className="font-medium">
                          {log.user_name || "Unknown"}
                        </span>{" "}
                        {log.action === "created" && "created this incident"}
                        {log.action === "updated" && "updated this incident"}
                        {log.action === "deleted" && "deleted this incident"}
                      </p>
                      {log.action === "updated" && log.changes && (
                        <ul className="mt-1 text-sm text-muted-foreground">
                          {Object.entries(log.changes).map(
                            ([field, change]) => {
                              const changeObj = change as
                                | { from?: unknown; to?: unknown }
                                | undefined;
                              return (
                                <li key={field}>
                                  {field}: &quot;{String(changeObj?.from ?? "")}
                                  &quot; &rarr; &quot;
                                  {String(changeObj?.to ?? "")}&quot;
                                </li>
                              );
                            },
                          )}
                        </ul>
                      )}
                      <p className="text-xs text-muted-foreground mt-1">
                        {new Date(log.created_at).toLocaleString()}
                      </p>
                    </div>
                  </div>
                );
              })}
            </div>
          ) : (
            <div className="text-center text-muted-foreground py-8">
              No audit history available.
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
