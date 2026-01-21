import { Link } from "@tanstack/react-router";
import { ShieldCheckIcon } from "lucide-react";

import { ThemeToggle } from "@/components/theme-toggle";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";

function HeaderNav() {
  return (
    <header className="border-b bg-background">
      <div className="mx-auto flex h-14 w-full max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        <div className="flex items-center gap-4">
          <Link to="/" className="flex items-center gap-2">
            <span className="bg-primary text-primary-foreground inline-flex size-8 items-center justify-center border">
              <ShieldCheckIcon className="size-4" />
            </span>
            <span className="text-sm font-semibold">Risk Register</span>
          </Link>
        </div>

        <div className="flex items-center gap-2">
          <Link to="/login" className={cn(buttonVariants({ variant: "outline", size: "sm" }))}>
            Log in
          </Link>
          <ThemeToggle />
        </div>
      </div>
    </header>
  );
}

export { HeaderNav };
