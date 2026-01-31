import { Link, useRouterState } from "@tanstack/react-router";
import { HomeIcon, SearchXIcon, TriangleAlertIcon } from "lucide-react";

import { HeaderNav } from "@/components/header-nav";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";

function NotFoundPage() {
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });

  return (
    <div className="flex min-h-svh flex-col bg-gradient-to-b from-background via-background to-muted/20">
      <HeaderNav />
      <main className="grid flex-1 place-items-center px-4 py-10">
        <div className="w-full max-w-2xl text-center">
          <div className="mb-8 flex justify-center">
            <div className="relative">
              <div className="absolute inset-0 bg-gradient-to-r from-primary/20 via-primary/10 to-primary/20 blur-3xl" />
              <div className="relative flex items-center justify-center gap-4">
                <div className="flex size-24 items-center justify-center rounded-2xl bg-gradient-to-br from-primary to-primary/80 shadow-xl shadow-primary/25">
                  <SearchXIcon className="size-12 text-primary-foreground" />
                </div>
              </div>
            </div>
          </div>

          <h1 className="mb-4 text-6xl font-bold tracking-tight text-foreground sm:text-7xl md:text-8xl">
            404
          </h1>

          <h2 className="mb-3 text-2xl font-semibold text-foreground sm:text-3xl">
            Page not found
          </h2>

          <p className="mb-8 max-w-md mx-auto text-muted-foreground sm:text-lg">
            The page you’re looking for doesn’t exist, has been moved, or you don’t have access to it.
          </p>

          <div className="mb-8 inline-flex items-center gap-2 rounded-lg border bg-muted/50 px-4 py-3 text-sm">
            <TriangleAlertIcon className="size-4 text-muted-foreground" />
            <span className="text-muted-foreground">Requested path:</span>
            <code className="rounded bg-background px-2 py-0.5 font-mono text-xs text-foreground">
              {pathname}
            </code>
          </div>

          <div className="flex items-center justify-center">
            <Link to="/" className={cn(buttonVariants({ variant: "default", size: "lg" }), "gap-2")}>
              <HomeIcon className="size-4" />
              Home
            </Link>
          </div>

          <p className="mt-8 text-xs text-muted-foreground">
            Need help? Contact support if you believe this is an error.
          </p>
        </div>
      </main>
    </div>
  );
}

export { NotFoundPage };

