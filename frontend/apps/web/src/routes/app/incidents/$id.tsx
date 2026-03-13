import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/app/incidents/$id")({
  component: IncidentDetail,
});

function IncidentDetail() {
  const { id } = Route.useParams();

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-6">Incident Detail</h1>
      <p className="text-muted-foreground">Incident ID: {id}</p>
      <p className="text-muted-foreground mt-2">Detail view will be implemented here.</p>
    </div>
  );
}
