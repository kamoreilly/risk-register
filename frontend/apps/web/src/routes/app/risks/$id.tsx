import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/app/risks/$id")({
  component: RiskDetail,
});

function RiskDetail() {
  const { id } = Route.useParams();

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold">Risk Detail</h1>
      <p className="text-muted-foreground">Viewing risk: {id}</p>
      <p className="mt-4 text-sm text-muted-foreground">
        Risk detail page will be implemented in Task 2.9.
      </p>
    </div>
  );
}
