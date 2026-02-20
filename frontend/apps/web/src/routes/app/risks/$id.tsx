import { createFileRoute, Link } from "@tanstack/react-router";
import * as React from "react";
import { useQueryClient } from "@tanstack/react-query";

import { useRisk, useUpdateRisk, useDeleteRisk, useCategories } from "@/hooks/useRisks";
import {
  useMitigations,
  useCreateMitigation,
  useDeleteMitigation,
} from "@/hooks/useMitigations";
import {
  useFrameworks,
  useRiskControls,
  useLinkControl,
  useUnlinkControl,
} from "@/hooks/useFrameworks";
import { useAuth } from "@/hooks/useAuth";
import { useSummarize, useDraftMitigation } from "@/hooks/useAI";
import { api } from "@/lib/api";
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
import type { RiskStatus, RiskSeverity } from "@/types/risk";
import type { MitigationStatus, Mitigation } from "@/types/mitigation";
import type { RiskFrameworkControl } from "@/types/framework";

export const Route = createFileRoute("/app/risks/$id")({
  component: RiskDetail,
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

const MITIGATION_STATUS_COLORS: Record<MitigationStatus, string> = {
  planned: "bg-yellow-100 text-yellow-800",
  in_progress: "bg-blue-100 text-blue-800",
  completed: "bg-green-100 text-green-800",
  cancelled: "bg-gray-100 text-gray-800",
};

function RiskDetail() {
  const { id } = Route.useParams();
  const { user } = useAuth();
  const navigate = Route.useNavigate();
  const queryClient = useQueryClient();

  const { data: risk, isLoading } = useRisk(id);
  const { data: categories } = useCategories();
  const updateRisk = useUpdateRisk(id);
  const deleteRisk = useDeleteRisk();

  // Mitigation hooks
  const { data: mitigations, isLoading: mitigationsLoading } = useMitigations(id);
  const createMitigation = useCreateMitigation(id);
  const deleteMitigation = useDeleteMitigation(id);

  // Framework control hooks
  const { data: frameworks } = useFrameworks();
  const { data: controls, isLoading: controlsLoading } = useRiskControls(id);
  const linkControl = useLinkControl(id);
  const unlinkControl = useUnlinkControl(id);

  const [isEditing, setIsEditing] = React.useState(false);
  const [title, setTitle] = React.useState("");
  const [description, setDescription] = React.useState("");
  const [status, setStatus] = React.useState<RiskStatus>("open");
  const [severity, setSeverity] = React.useState<RiskSeverity>("medium");
  const [categoryId, setCategoryId] = React.useState<string>("");
  const [reviewDate, setReviewDate] = React.useState("");

  // Mitigation state
  const [isAddingMitigation, setIsAddingMitigation] = React.useState(false);
  const [editingMitigationId, setEditingMitigationId] = React.useState<string | null>(null);
  const [mitigationDescription, setMitigationDescription] = React.useState("");
  const [mitigationOwner, setMitigationOwner] = React.useState("");
  const [mitigationStatus, setMitigationStatus] = React.useState<MitigationStatus>("planned");
  const [mitigationDueDate, setMitigationDueDate] = React.useState("");

  // Control state
  const [isAddingControl, setIsAddingControl] = React.useState(false);
  const [controlFrameworkId, setControlFrameworkId] = React.useState("");
  const [controlRef, setControlRef] = React.useState("");
  const [controlNotes, setControlNotes] = React.useState("");

  // AI state
  const [aiSummary, setAiSummary] = React.useState<string | null>(null);
  const summarize = useSummarize();
  const draftMitigation = useDraftMitigation();

  React.useEffect(() => {
    if (risk) {
      setTitle(risk.title);
      setDescription(risk.description || "");
      setStatus(risk.status);
      setSeverity(risk.severity);
      setCategoryId(risk.category_id || "");
      setReviewDate(risk.review_date ? risk.review_date.split("T")[0] : "");
    }
  }, [risk]);

  const handleSave = async () => {
    try {
      await updateRisk.mutateAsync({
        title,
        description: description || undefined,
        status,
        severity,
        category_id: categoryId || undefined,
        review_date: reviewDate || undefined,
      });
      toast.success("Risk updated");
      setIsEditing(false);
    } catch (error) {
      toast.error("Failed to update risk");
    }
  };

  const handleDelete = async () => {
    if (!confirm("Are you sure you want to delete this risk?")) return;

    try {
      await deleteRisk.mutateAsync(id);
      toast.success("Risk deleted");
      navigate({ to: "/app/risks" });
    } catch (error) {
      toast.error("Failed to delete risk");
    }
  };

  // Mitigation handlers
  const resetMitigationForm = () => {
    setMitigationDescription("");
    setMitigationOwner("");
    setMitigationStatus("planned");
    setMitigationDueDate("");
    setIsAddingMitigation(false);
    setEditingMitigationId(null);
  };

  const handleAddMitigation = async () => {
    if (!mitigationDescription.trim()) {
      toast.error("Description is required");
      return;
    }

    try {
      await createMitigation.mutateAsync({
        description: mitigationDescription,
        owner: mitigationOwner || undefined,
        status: mitigationStatus,
        due_date: mitigationDueDate || undefined,
      });
      toast.success("Mitigation added");
      resetMitigationForm();
    } catch (error) {
      toast.error("Failed to add mitigation");
    }
  };

  const startEditMitigation = (mitigation: Mitigation) => {
    setEditingMitigationId(mitigation.id);
    setMitigationDescription(mitigation.description);
    setMitigationOwner(mitigation.owner || "");
    setMitigationStatus(mitigation.status);
    setMitigationDueDate(mitigation.due_date ? mitigation.due_date.split("T")[0] : "");
    setIsAddingMitigation(false);
  };

  const handleUpdateMitigation = async () => {
    if (!editingMitigationId || !mitigationDescription.trim()) {
      toast.error("Description is required");
      return;
    }

    try {
      await api.put(`/api/v1/risks/${id}/mitigations/${editingMitigationId}`, {
        description: mitigationDescription,
        owner: mitigationOwner || undefined,
        status: mitigationStatus,
        due_date: mitigationDueDate || undefined,
      });
      toast.success("Mitigation updated");
      resetMitigationForm();
      queryClient.invalidateQueries({ queryKey: ["mitigations", id] });
    } catch (error) {
      toast.error("Failed to update mitigation");
    }
  };

  const handleDeleteMitigation = async (mitigationId: string) => {
    if (!confirm("Are you sure you want to delete this mitigation?")) return;

    try {
      await deleteMitigation.mutateAsync(mitigationId);
      toast.success("Mitigation deleted");
      if (editingMitigationId === mitigationId) {
        resetMitigationForm();
      }
    } catch (error) {
      toast.error("Failed to delete mitigation");
    }
  };

  // Control handlers
  const resetControlForm = () => {
    setControlFrameworkId("");
    setControlRef("");
    setControlNotes("");
    setIsAddingControl(false);
  };

  const handleAddControl = async () => {
    if (!controlFrameworkId) {
      toast.error("Framework is required");
      return;
    }
    if (!controlRef.trim()) {
      toast.error("Control reference is required");
      return;
    }

    try {
      await linkControl.mutateAsync({
        framework_id: controlFrameworkId,
        control_ref: controlRef,
        notes: controlNotes || undefined,
      });
      toast.success("Control linked");
      resetControlForm();
    } catch (error) {
      toast.error("Failed to link control");
    }
  };

  const handleDeleteControl = async (control: RiskFrameworkControl) => {
    if (!confirm(`Are you sure you want to unlink "${control.framework_name}: ${control.control_ref}"?`)) return;

    try {
      await unlinkControl.mutateAsync(control.id);
      toast.success("Control unlinked");
    } catch (error) {
      toast.error("Failed to unlink control");
    }
  };

  // AI handlers
  const handleSummarize = async () => {
    if (!risk) return;

    try {
      const result = await summarize.mutateAsync({
        title: risk.title,
        description: risk.description,
        severity: risk.severity,
        status: risk.status,
      });
      setAiSummary(result.summary);
      toast.success("Summary generated");
    } catch (error) {
      toast.error("Failed to generate summary");
    }
  };

  const handleDraftMitigation = async () => {
    if (!risk) return;

    try {
      const result = await draftMitigation.mutateAsync({
        risk_title: risk.title,
        risk_description: risk.description,
        severity: risk.severity,
      });
      setMitigationDescription(result.draft);
      toast.success("Draft generated");
    } catch (error) {
      toast.error("Failed to generate draft");
    }
  };

  if (isLoading) {
    return (
      <div className="p-8">
        <div className="text-center text-muted-foreground">Loading...</div>
      </div>
    );
  }

  if (!risk) {
    return (
      <div className="p-8">
        <div className="text-center text-muted-foreground">Risk not found</div>
        <div className="text-center mt-4">
          <Link to="/app/risks" className={cn(buttonVariants({ variant: "outline" }))}>
            Back to Risks
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <div>
          <Link to="/app/risks" className="text-sm text-muted-foreground hover:underline mb-2 block">
            &larr; Back to Risks
          </Link>
          <h1 className="text-2xl font-bold">{isEditing ? "Edit Risk" : risk.title}</h1>
        </div>
        <div className="flex gap-2">
          {isEditing ? (
            <>
              <Button variant="outline" onClick={() => setIsEditing(false)}>Cancel</Button>
              <Button onClick={handleSave} disabled={updateRisk.isPending}>Save</Button>
            </>
          ) : (
            <>
              <Button variant="outline" onClick={() => setIsEditing(true)}>Edit</Button>
              <Button variant="destructive" onClick={handleDelete} disabled={deleteRisk.isPending}>Delete</Button>
            </>
          )}
        </div>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center gap-2">
            <span className={cn("px-2 py-1 rounded text-xs font-medium", STATUS_COLORS[risk.status])}>
              {risk.status}
            </span>
            <span className={cn("px-2 py-1 rounded text-xs font-medium", SEVERITY_COLORS[risk.severity])}>
              {risk.severity}
            </span>
          </div>
        </CardHeader>
        <CardContent className="space-y-6">
          {isEditing ? (
            <>
              <div className="grid gap-2">
                <Label htmlFor="title">Title</Label>
                <Input id="title" value={title} onChange={(e) => setTitle(e.target.value)} />
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
                  <Select value={status} onValueChange={(v) => setStatus(v as RiskStatus)}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="open">Open</SelectItem>
                      <SelectItem value="mitigating">Mitigating</SelectItem>
                      <SelectItem value="resolved">Resolved</SelectItem>
                      <SelectItem value="accepted">Accepted</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="grid gap-2">
                  <Label>Severity</Label>
                  <Select value={severity} onValueChange={(v) => setSeverity(v as RiskSeverity)}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="low">Low</SelectItem>
                      <SelectItem value="medium">Medium</SelectItem>
                      <SelectItem value="high">High</SelectItem>
                      <SelectItem value="critical">Critical</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label>Category</Label>
                  <Select value={categoryId} onValueChange={setCategoryId}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select category" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="">None</SelectItem>
                      {categories?.map((cat) => (
                        <SelectItem key={cat.id} value={cat.id}>{cat.name}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="grid gap-2">
                  <Label>Review Date</Label>
                  <Input type="date" value={reviewDate} onChange={(e) => setReviewDate(e.target.value)} />
                </div>
              </div>
            </>
          ) : (
            <>
              <div>
                <div className="flex items-center justify-between mb-2">
                  <h3 className="font-medium">Description</h3>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleSummarize}
                    disabled={summarize.isPending}
                  >
                    {summarize.isPending ? "Generating..." : "Summarize"}
                  </Button>
                </div>
                <p className="text-muted-foreground">{risk.description || "No description"}</p>
              </div>
              {aiSummary && (
                <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                  <h4 className="font-medium text-blue-800 mb-2">AI Summary</h4>
                  <p className="text-sm text-blue-700">{aiSummary}</p>
                </div>
              )}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <h3 className="font-medium mb-2">Category</h3>
                  <p className="text-muted-foreground">{risk.category?.name || "Uncategorized"}</p>
                </div>
                <div>
                  <h3 className="font-medium mb-2">Review Date</h3>
                  <p className="text-muted-foreground">
                    {risk.review_date ? new Date(risk.review_date).toLocaleDateString() : "Not set"}
                  </p>
                </div>
              </div>
              <div className="text-sm text-muted-foreground pt-4 border-t">
                <p>Created: {new Date(risk.created_at).toLocaleString()}</p>
                <p>Updated: {new Date(risk.updated_at).toLocaleString()}</p>
              </div>
            </>
          )}
        </CardContent>
      </Card>

      {/* Mitigations Section */}
      <Card className="mt-6">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Mitigations</CardTitle>
              <CardDescription>
                Actions to reduce or eliminate this risk
              </CardDescription>
            </div>
            {!isAddingMitigation && !editingMitigationId && (
              <Button onClick={() => setIsAddingMitigation(true)}>
                Add Mitigation
              </Button>
            )}
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Add/Edit Mitigation Form */}
          {(isAddingMitigation || editingMitigationId) && (
            <div className="border rounded-lg p-4 space-y-4 bg-muted/50">
              <h4 className="font-medium">
                {isAddingMitigation ? "Add New Mitigation" : "Edit Mitigation"}
              </h4>
              <div className="grid gap-2">
                <div className="flex items-center justify-between">
                  <Label htmlFor="mitigation-description">Description</Label>
                  {isAddingMitigation && (
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={handleDraftMitigation}
                      disabled={draftMitigation.isPending}
                    >
                      {draftMitigation.isPending ? "Drafting..." : "Draft with AI"}
                    </Button>
                  )}
                </div>
                <textarea
                  id="mitigation-description"
                  className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value={mitigationDescription}
                  onChange={(e) => setMitigationDescription(e.target.value)}
                  placeholder="Describe the mitigation action..."
                />
              </div>
              <div className="grid grid-cols-3 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="mitigation-owner">Owner</Label>
                  <Input
                    id="mitigation-owner"
                    value={mitigationOwner}
                    onChange={(e) => setMitigationOwner(e.target.value)}
                    placeholder="Responsible person"
                  />
                </div>
                <div className="grid gap-2">
                  <Label>Status</Label>
                  <Select
                    value={mitigationStatus}
                    onValueChange={(v) => setMitigationStatus(v as MitigationStatus)}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="planned">Planned</SelectItem>
                      <SelectItem value="in_progress">In Progress</SelectItem>
                      <SelectItem value="completed">Completed</SelectItem>
                      <SelectItem value="cancelled">Cancelled</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="mitigation-due-date">Due Date</Label>
                  <Input
                    id="mitigation-due-date"
                    type="date"
                    value={mitigationDueDate}
                    onChange={(e) => setMitigationDueDate(e.target.value)}
                  />
                </div>
              </div>
              <div className="flex gap-2 justify-end">
                <Button variant="outline" onClick={resetMitigationForm}>
                  Cancel
                </Button>
                <Button
                  onClick={isAddingMitigation ? handleAddMitigation : handleUpdateMitigation}
                  disabled={createMitigation.isPending}
                >
                  {isAddingMitigation ? "Add" : "Save"}
                </Button>
              </div>
            </div>
          )}

          {/* Mitigations List */}
          {mitigationsLoading ? (
            <div className="text-center text-muted-foreground py-4">
              Loading mitigations...
            </div>
          ) : mitigations && mitigations.length > 0 ? (
            <div className="space-y-3">
              {mitigations.map((mitigation) => (
                <div
                  key={mitigation.id}
                  className="border rounded-lg p-4 flex items-start justify-between gap-4"
                >
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <span
                        className={cn(
                          "px-2 py-0.5 rounded text-xs font-medium",
                          MITIGATION_STATUS_COLORS[mitigation.status]
                        )}
                      >
                        {mitigation.status.replace("_", " ")}
                      </span>
                    </div>
                    <p className="text-sm mb-2">{mitigation.description}</p>
                    <div className="flex items-center gap-4 text-xs text-muted-foreground">
                      {mitigation.owner && (
                        <span>Owner: {mitigation.owner}</span>
                      )}
                      {mitigation.due_date && (
                        <span>Due: {new Date(mitigation.due_date).toLocaleDateString()}</span>
                      )}
                    </div>
                  </div>
                  <div className="flex gap-2 shrink-0">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => startEditMitigation(mitigation)}
                      disabled={editingMitigationId !== null && editingMitigationId !== mitigation.id}
                    >
                      Edit
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => handleDeleteMitigation(mitigation.id)}
                      disabled={deleteMitigation.isPending}
                    >
                      Delete
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            !isAddingMitigation && (
              <div className="text-center text-muted-foreground py-8">
                No mitigations yet. Click "Add Mitigation" to create one.
              </div>
            )
          )}
        </CardContent>
      </Card>

      {/* Compliance Controls Section */}
      <Card className="mt-6">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Compliance Controls</CardTitle>
              <CardDescription>
                Framework controls mapped to this risk
              </CardDescription>
            </div>
            {!isAddingControl && (
              <Button onClick={() => setIsAddingControl(true)}>
                Link Control
              </Button>
            )}
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Add Control Form */}
          {isAddingControl && (
            <div className="border rounded-lg p-4 space-y-4 bg-muted/50">
              <h4 className="font-medium">Link New Control</h4>
              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="control-framework">Framework</Label>
                  <Select
                    value={controlFrameworkId}
                    onValueChange={(v) => v && setControlFrameworkId(v)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select framework" />
                    </SelectTrigger>
                    <SelectContent>
                      {frameworks?.map((framework) => (
                        <SelectItem key={framework.id} value={framework.id}>
                          {framework.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="control-ref">Control Reference</Label>
                  <Input
                    id="control-ref"
                    value={controlRef}
                    onChange={(e) => setControlRef(e.target.value)}
                    placeholder="e.g., A.12.1.1"
                  />
                </div>
              </div>
              <div className="grid gap-2">
                <Label htmlFor="control-notes">Notes (optional)</Label>
                <textarea
                  id="control-notes"
                  className="flex min-h-[60px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value={controlNotes}
                  onChange={(e) => setControlNotes(e.target.value)}
                  placeholder="Additional notes about this control mapping..."
                />
              </div>
              <div className="flex gap-2 justify-end">
                <Button variant="outline" onClick={resetControlForm}>
                  Cancel
                </Button>
                <Button
                  onClick={handleAddControl}
                  disabled={linkControl.isPending}
                >
                  Link
                </Button>
              </div>
            </div>
          )}

          {/* Controls List */}
          {controlsLoading ? (
            <div className="text-center text-muted-foreground py-4">
              Loading controls...
            </div>
          ) : controls && controls.length > 0 ? (
            <div className="flex flex-wrap gap-2">
              {controls.map((control) => (
                <div
                  key={control.id}
                  className="inline-flex items-center gap-2 px-3 py-1.5 bg-blue-100 text-blue-800 rounded-full text-sm font-medium group"
                >
                  <span>{control.framework_name}: {control.control_ref}</span>
                  <button
                    onClick={() => handleDeleteControl(control)}
                    disabled={unlinkControl.isPending}
                    className="ml-1 text-blue-600 hover:text-blue-900 opacity-0 group-hover:opacity-100 transition-opacity"
                    title="Unlink control"
                  >
                    &times;
                  </button>
                  {control.notes && (
                    <span
                      className="text-xs text-blue-600 cursor-help"
                      title={control.notes}
                    >
                      ?
                    </span>
                  )}
                </div>
              ))}
            </div>
          ) : (
            !isAddingControl && (
              <div className="text-center text-muted-foreground py-8">
                No controls mapped. Click "Link Control" to add a compliance mapping.
              </div>
            )
          )}
        </CardContent>
      </Card>
    </div>
  );
}
