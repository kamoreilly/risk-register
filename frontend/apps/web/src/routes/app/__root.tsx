import { Outlet, createFileRoute, useNavigate } from "@tanstack/react-router";
import * as React from "react";

import { useAuth } from "@/hooks/useAuth";
import { HeaderNav } from "@/components/header-nav";
import { Button } from "@/components/ui/button";
import Loader from "@/components/loader";

export const Route = createFileRoute("/app/__root")({
  component: AppLayout,
});

function AppLayout() {
  const navigate = useNavigate();
  const { user, isLoading, isAuthenticated, logout } = useAuth();

  React.useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate({ to: "/login" });
    }
  }, [isLoading, isAuthenticated, navigate]);

  if (isLoading) {
    return (
      <div className="flex min-h-svh items-center justify-center">
        <Loader />
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  return (
    <div className="flex min-h-svh flex-col bg-background">
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-14 items-center px-4">
          <HeaderNav />
          <div className="ml-auto flex items-center gap-4">
            <span className="text-sm text-muted-foreground">{user?.name}</span>
            <Button variant="outline" size="sm" onClick={logout}>
              Log out
            </Button>
          </div>
        </div>
      </header>
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  );
}
