import { createFileRoute } from "@tanstack/react-router";

import { useAuth } from "@/hooks/useAuth";

export const Route = createFileRoute("/app/")({
  component: Dashboard,
});

function Dashboard() {
  const { user } = useAuth();

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold">Dashboard</h1>
      <p className="text-muted-foreground">
        Welcome back, {user?.name}!
      </p>
      <p className="mt-4 text-sm text-muted-foreground">
        Dashboard widgets will be added in Slice 4.
      </p>
    </div>
  );
}
