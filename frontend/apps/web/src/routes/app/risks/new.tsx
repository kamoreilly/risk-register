import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/app/risks/new")({
  component: NewRisk,
});

function NewRisk() {
  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold">New Risk</h1>
      <p className="text-muted-foreground">Create a new risk</p>
      <p className="mt-4 text-sm text-muted-foreground">
        Risk form will be implemented in Task 2.10.
      </p>
    </div>
  );
}
