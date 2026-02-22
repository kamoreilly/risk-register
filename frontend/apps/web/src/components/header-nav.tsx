import { Link } from "@tanstack/react-router";
import { ShieldCheckIcon } from "lucide-react";

import { ThemeToggle } from "@/components/theme-toggle";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import type { User } from "@/types/auth";

interface HeaderNavProps {
  user?: User;
  logout?: () => void;
  fullWidth?: boolean;
}

const ADMIN_LINKS = [
  { to: "/app", label: "Dashboard", exact: true },
  { to: "/app/board", label: "Board" },
  { to: "/app/calendar", label: "Calendar" },
  { to: "/app/risks", label: "Risks" },
  { to: "/app/frameworks", label: "Frameworks" },
  { to: "/app/categories", label: "Categories" },
] as const;

const MEMBER_LINKS = [
  { to: "/app", label: "Dashboard", exact: true },
  { to: "/app/board", label: "Board" },
  { to: "/app/calendar", label: "Calendar" },
  { to: "/app/risks", label: "Risks" },
] as const;

function HeaderNav({ user, logout, fullWidth = false }: HeaderNavProps) {
  const links = user?.role === "admin" ? ADMIN_LINKS : MEMBER_LINKS;

  return (
    <header
      className={cn(
        "border-b bg-background",
        user &&
          "sticky top-0 z-50 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60",
      )}
    >
      <div
        className={cn(
          "flex h-14 items-center justify-between",
          fullWidth ? "px-4" : "mx-auto w-full max-w-7xl px-4 sm:px-6 lg:px-8",
        )}
      >
        <div className="flex items-center gap-8">
          <Link to={user ? "/app" : "/"} className="flex items-center gap-2">
            <span className="bg-primary text-primary-foreground inline-flex size-8 items-center justify-center rounded-md border">
              <ShieldCheckIcon className="size-4" />
            </span>
            <span className="text-sm font-semibold">Risk Register</span>
          </Link>

          {user && (
            <nav className="hidden md:flex items-center gap-6 text-sm font-medium">
              {links.map((link) => (
                <Link
                  key={link.to}
                  to={link.to}
                  activeProps={{ className: "text-foreground font-semibold" }}
                  inactiveProps={{
                    className: "text-muted-foreground hover:text-foreground/80",
                  }}
                  activeOptions={link.exact ? { exact: true } : undefined}
                >
                  {link.label}
                </Link>
              ))}
            </nav>
          )}
        </div>

        <div className="flex items-center gap-4">
          <ThemeToggle />

          {user ? (
            <div className="flex items-center gap-4">
              <span className="hidden text-sm text-muted-foreground sm:inline-block">
                {user.name}
              </span>
              <Button variant="outline" size="sm" onClick={logout}>
                Log out
              </Button>
            </div>
          ) : (
            <Link
              to="/login"
              className={cn(buttonVariants({ variant: "outline", size: "sm" }))}
            >
              Log in
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}

export { HeaderNav };
