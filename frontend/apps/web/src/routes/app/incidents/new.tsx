import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/app/incidents/new")({
  component: NewIncident,
});

function NewIncident() {
  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-6">New Incident</h1>
      <p className="text-muted-foreground">New incident form will be implemented here.</p>
    </div>
  );
}
