import { Link, createFileRoute, useNavigate } from "@tanstack/react-router";
import * as React from "react";
import {
  MoreHorizontalIcon,
  PencilIcon,
  PlusIcon,
  TrashIcon,
} from "lucide-react";
import { toast } from "sonner";

import Loader from "@/components/loader";
import { ApiError } from "@/lib/api";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  useControlLinkedRisks,
  useControls,
  useCreateControl,
  useDeleteControl,
  useUpdateControl,
} from "@/hooks/useControls";
import { useFrameworks } from "@/hooks/useFrameworks";
import { useAuth } from "@/hooks/useAuth";
import type { FrameworkControl } from "@/types/framework";

export const Route = createFileRoute("/app/controls")({
  component: ControlsPage,
});

function ControlsPage() {
  const navigate = useNavigate();
  const { user, isLoading: authLoading } = useAuth();
  const { data: frameworks, isLoading: frameworksLoading } = useFrameworks();
  const [selectedFrameworkId, setSelectedFrameworkId] = React.useState("all");
  const [search, setSearch] = React.useState("");
  const deferredSearch = React.useDeferredValue(search.trim());
  const {
    data: controls,
    isLoading: controlsLoading,
    error,
  } = useControls({
    frameworkId:
      selectedFrameworkId === "all" ? undefined : selectedFrameworkId,
    search: deferredSearch || undefined,
  });
  const createControl = useCreateControl();
  const updateControl = useUpdateControl();
  const deleteControl = useDeleteControl();

  const [isFormOpen, setIsFormOpen] = React.useState(false);
  const [expandedControlId, setExpandedControlId] = React.useState<string | null>(
    null,
  );
  const [editingControl, setEditingControl] =
    React.useState<FrameworkControl | null>(null);
  const [formData, setFormData] = React.useState({
    framework_id: "",
    control_ref: "",
    title: "",
    description: "",
  });
  const { data: linkedRisks, isLoading: linkedRisksLoading } =
    useControlLinkedRisks(expandedControlId ?? "", !!expandedControlId);

  React.useEffect(() => {
    if (!authLoading && user && user.role !== "admin") {
      navigate({ to: "/app" });
    }
  }, [authLoading, navigate, user]);

  if (authLoading || frameworksLoading) {
    return (
      <div className="flex items-center justify-center p-8">
        <Loader />
      </div>
    );
  }

  if (user?.role !== "admin") {
    return null;
  }

  const openCreateForm = () => {
    setEditingControl(null);
    setFormData({
      framework_id: selectedFrameworkId === "all" ? "" : selectedFrameworkId,
      control_ref: "",
      title: "",
      description: "",
    });
    setIsFormOpen(true);
  };

  const handleEdit = (control: FrameworkControl) => {
    setEditingControl(control);
    setFormData({
      framework_id: control.framework_id,
      control_ref: control.control_ref,
      title: control.title,
      description: control.description || "",
    });
    setIsFormOpen(true);
  };

  const handleCancel = () => {
    setEditingControl(null);
    setIsFormOpen(false);
    setFormData({
      framework_id: "",
      control_ref: "",
      title: "",
      description: "",
    });
  };

  const toggleLinkedRisks = (controlId: string) => {
    setExpandedControlId((current) =>
      current === controlId ? null : controlId,
    );
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!formData.framework_id) {
      toast.error("Framework is required");
      return;
    }
    if (!formData.control_ref.trim()) {
      toast.error("Control reference is required");
      return;
    }
    if (!formData.title.trim()) {
      toast.error("Control title is required");
      return;
    }

    try {
      if (editingControl) {
        await updateControl.mutateAsync({
          id: editingControl.id,
          control_ref: formData.control_ref,
          title: formData.title,
          description: formData.description || undefined,
        });
        toast.success("Control updated successfully");
      } else {
        await createControl.mutateAsync({
          framework_id: formData.framework_id,
          control_ref: formData.control_ref,
          title: formData.title,
          description: formData.description || undefined,
        });
        toast.success("Control created successfully");
      }

      handleCancel();
    } catch (error) {
      const message = error instanceof ApiError ? error.message : "Failed to save control";
      toast.error(message);
    }
  };

  const handleDelete = async (control: FrameworkControl) => {
    if (
      !window.confirm(
        `Delete control ${control.framework_name}: ${control.control_ref}?`,
      )
    ) {
      return;
    }

    try {
      await deleteControl.mutateAsync(control.id);
      toast.success("Control deleted successfully");
    } catch (error) {
      const message = error instanceof ApiError ? error.message : "Failed to delete control";
      toast.error(message);
    }
  };

  const isSubmitting = createControl.isPending || updateControl.isPending;

  return (
    <div className="space-y-6 p-8">
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="flex flex-1 flex-col gap-4 md:flex-row">
          <div className="grid gap-2 md:w-72">
            <Label htmlFor="controls-framework-filter">Framework</Label>
            <Select
              value={selectedFrameworkId}
              onValueChange={(value) => setSelectedFrameworkId(value ?? "all")}
            >
              <SelectTrigger id="controls-framework-filter">
                <SelectValue placeholder="All frameworks" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All frameworks</SelectItem>
                {frameworks?.map((framework) => (
                  <SelectItem key={framework.id} value={framework.id}>
                    {framework.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="grid gap-2 md:flex-1">
            <Label htmlFor="controls-search">Search</Label>
            <Input
              id="controls-search"
              value={search}
              onChange={(event) => setSearch(event.target.value)}
              placeholder="Search by reference, title, or description"
            />
          </div>
        </div>

        <Button onClick={openCreateForm}>
          <PlusIcon className="mr-2 size-4" />
          Add Control
        </Button>
      </div>

      {isFormOpen && (
        <Card>
          <CardHeader>
            <CardTitle>
              {editingControl ? "Edit Control" : "Create New Control"}
            </CardTitle>
            <CardDescription>
              {editingControl
                ? "Update the reusable control definition"
                : "Add a reusable control definition under a compliance framework"}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid gap-4 md:grid-cols-2">
                <div className="grid gap-2">
                  <Label htmlFor="control-framework">Framework</Label>
                  <Select
                    value={formData.framework_id}
                    onValueChange={(value) =>
                      setFormData((current) => ({
                        ...current,
                        framework_id: value ?? "",
                      }))
                    }
                    disabled={!!editingControl}
                  >
                    <SelectTrigger id="control-framework">
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
                    value={formData.control_ref}
                    onChange={(event) =>
                      setFormData((current) => ({
                        ...current,
                        control_ref: event.target.value,
                      }))
                    }
                    placeholder="e.g. CC6.1"
                  />
                </div>
              </div>

              <div className="grid gap-2">
                <Label htmlFor="control-title">Title</Label>
                <Input
                  id="control-title"
                  value={formData.title}
                  onChange={(event) =>
                    setFormData((current) => ({
                      ...current,
                      title: event.target.value,
                    }))
                  }
                  placeholder="Describe the control"
                />
              </div>

              <div className="grid gap-2">
                <Label htmlFor="control-description">Description</Label>
                <textarea
                  id="control-description"
                  className="flex min-h-[88px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value={formData.description}
                  onChange={(event) =>
                    setFormData((current) => ({
                      ...current,
                      description: event.target.value,
                    }))
                  }
                  placeholder="Optional implementation notes or detail"
                />
              </div>

              <div className="flex gap-2">
                <Button type="submit" disabled={isSubmitting}>
                  {isSubmitting
                    ? "Saving..."
                    : editingControl
                      ? "Update"
                      : "Create"}
                </Button>
                <Button type="button" variant="outline" onClick={handleCancel}>
                  Cancel
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {controlsLoading ? (
        <div className="flex items-center justify-center p-8">
          <Loader />
        </div>
      ) : error ? (
        <div className="text-center text-muted-foreground p-8">
          Failed to load controls.
        </div>
      ) : controls?.length ? (
        <div className="grid gap-4">
          {controls.map((control) => (
            <Card key={control.id}>
              <CardHeader className="flex flex-row items-start justify-between space-y-0">
                <div className="space-y-1">
                  <CardTitle className="text-lg">
                    {control.framework_name}: {control.control_ref}
                  </CardTitle>
                  <CardDescription>{control.title}</CardDescription>
                </div>

                <DropdownMenu>
                  <DropdownMenuTrigger className="-mr-2 flex size-8 items-center justify-center rounded-md hover:bg-muted hover:text-foreground dark:hover:bg-muted/50 aria-expanded:bg-muted aria-expanded:text-foreground">
                    <MoreHorizontalIcon className="h-4 w-4" />
                    <span className="sr-only">Actions</span>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem onClick={() => handleEdit(control)}>
                      <PencilIcon className="mr-2 h-4 w-4" />
                      Edit
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      onClick={() => handleDelete(control)}
                      className="text-destructive focus:text-destructive"
                    >
                      <TrashIcon className="mr-2 h-4 w-4" />
                      Delete
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </CardHeader>
              <CardContent className="space-y-3">
                {control.description && (
                  <p className="text-sm text-muted-foreground">
                    {control.description}
                  </p>
                )}
                <div className="text-xs text-muted-foreground">
                  Linked to {control.linked_risk_count}{" "}
                  {control.linked_risk_count === 1 ? "risk" : "risks"}
                </div>
                <div className="flex flex-wrap items-center gap-2">
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => toggleLinkedRisks(control.id)}
                    disabled={control.linked_risk_count === 0}
                  >
                    {expandedControlId === control.id
                      ? "Hide linked risks"
                      : "View linked risks"}
                  </Button>
                  {selectedFrameworkId !== control.framework_id && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => setSelectedFrameworkId(control.framework_id)}
                    >
                      Filter by {control.framework_name}
                    </Button>
                  )}
                </div>
                {expandedControlId === control.id && (
                  <div className="rounded-lg border bg-muted/40 p-3">
                    {linkedRisksLoading ? (
                      <p className="text-sm text-muted-foreground">
                        Loading linked risks...
                      </p>
                    ) : linkedRisks?.length ? (
                      <div className="space-y-2">
                        {linkedRisks.map((risk) => (
                          <Link
                            key={risk.id}
                            to="/app/risks/$id"
                            params={{ id: risk.id }}
                            className="block rounded-md border bg-background px-3 py-2 text-sm transition-colors hover:bg-accent"
                          >
                            <div className="font-medium">{risk.title}</div>
                            <div className="mt-1 flex flex-wrap gap-3 text-xs text-muted-foreground">
                              <span>Status: {risk.status}</span>
                              <span>Severity: {risk.severity}</span>
                              {risk.category_name && (
                                <span>Category: {risk.category_name}</span>
                              )}
                              {risk.owner_name && (
                                <span>Owner: {risk.owner_name}</span>
                              )}
                            </div>
                          </Link>
                        ))}
                      </div>
                    ) : (
                      <p className="text-sm text-muted-foreground">
                        No linked risks found.
                      </p>
                    )}
                  </div>
                )}
              </CardContent>
            </Card>
          ))}
        </div>
      ) : (
        <div className="text-center text-muted-foreground p-8">
          No controls found. Add one to start linking controls to risks.
        </div>
      )}
    </div>
  );
}
