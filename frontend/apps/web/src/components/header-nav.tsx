import { Link } from "@tanstack/react-router";
import { ShieldCheckIcon } from "lucide-react";

import { ThemeToggle } from "@/components/theme-toggle";
import { buttonVariants } from "@/components/ui/button";
import { NavigationMenu, NavigationMenuItem, NavigationMenuList } from "@/components/ui/navigation-menu";
import { cn } from "@/lib/utils";

function HeaderNav() {
  return (
    <header className="border-b bg-background">
      <div className="mx-auto flex h-14 w-full max-w-5xl items-center justify-between px-4">
        <div className="flex items-center gap-4">
          <Link to="/" className="flex items-center gap-2">
            <span className="bg-primary text-primary-foreground inline-flex size-8 items-center justify-center border">
              <ShieldCheckIcon className="size-4" />
            </span>
            <span className="text-sm font-semibold">Risk Register</span>
          </Link>

          <NavigationMenu className="hidden sm:block">
            <NavigationMenuList>
              <NavigationMenuItem>
                <Link to="/" className="hover:bg-muted hover:text-foreground inline-flex h-8 items-center px-2.5 text-xs font-medium">
                  Home
                </Link>
              </NavigationMenuItem>
            </NavigationMenuList>
          </NavigationMenu>
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
