import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import * as React from "react";

import { useCreateIncident, useIncidentCategories } from "@/hooks/useIncidents";
import { useUsers } from "@/hooks/useUsers";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { toast } from "sonner";
import type { IncidentPriority, IncidentStatus } from "@/types/incident";

export const Route = createFileRoute("/app/incidents/new")({
  component: NewIncident,
});

function NewIncident() {
  const navigate = useNavigate();
  const createIncident = useCreateIncident();
  const { data: categories } = useIncidentCategories();
  const { data: users } = useUsers();

  const [title, setTitle] = React.useState("");
  const [description, setDescription] = React.useState("");
  const [priority, setPriority] = React.useState<IncidentPriority>("p3");
  const [status, setStatus] = React.useState<IncidentStatus>("new");
  const [categoryId, setCategoryId] = React.useState<string>("");
  const [assigneeId, setAssigneeId] = React.useState<string>("");
  const [serviceAffected, setServiceAffected] = React.useState("");
  const [occurredAt, setOccurredAt] = React.useState("");
  const [detectedAt, setDetectedAt] = React.useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!title.trim()) {
      toast.error("Title is required");
      return;
    }

    try {
      const result = await createIncident.mutateAsync({
        title,
        description: description || undefined,
        priority,
        status,
        category_id: categoryId || undefined,
        assignee_id: assigneeId || undefined,
        service_affected: serviceAffected || undefined,
        occurred_at: occurredAt || undefined,
        detected_at: detectedAt || undefined,
      });
      toast.success("Incident created");
      // Navigate to the incident detail page
      navigate({ to: "/app/incidents/$id", params: { id: result.id } });
    } catch (error) {
      toast.error("Failed to create incident");
    }
  };

  return (
    <div className="p-8 max-w-2xl mx-auto">
      <div className="mb-6">
        <Link
          to="/app/incidents"
          className="text-sm text-muted-foreground hover:underline mb-2 block"
        >
          &larr; Back to Incidents
        </Link>
        <h1 className="text-2xl font-bold">New Incident</h1>
        <p className="text-muted-foreground">Create a new incident entry</p>
      </div>

      <Card>
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-6 pt-6">
            <div className="grid gap-2">
              <Label htmlFor="title">Title *</Label>
              <Input
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Enter incident title"
                required
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="description">Description</Label>
              <textarea
                id="description"
                className="flex min-h-[100px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Describe the incident..."
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
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
                <Label>Assignee</Label>
                <Select
                  value={assigneeId}
                  onValueChange={(value) => setAssigneeId(value ?? "")}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select assignee" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="">Unassigned</SelectItem>
                    {users?.map((user) => (
                      <SelectItem key={user.id} value={user.id}>
                        {user.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="serviceAffected">Service Affected</Label>
                <Input
                  id="serviceAffected"
                  value={serviceAffected}
                  onChange={(e) => setServiceAffected(e.target.value)}
                  placeholder="e.g., API, Database, Frontend"
                />
              </div>

              <div className="grid gap-2">
                {/* Empty placeholder for grid alignment */}
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

            <div className="flex justify-end gap-2 pt-4">
              <Link to="/app/incidents">
                <Button type="button" variant="outline">
                  Cancel
                </Button>
              </Link>
              <Button type="submit" disabled={createIncident.isPending}>
                {createIncident.isPending ? "Creating..." : "Create Incident"}
              </Button>
            </div>
          </CardContent>
        </form>
      </Card>
    </div>
  );
}
