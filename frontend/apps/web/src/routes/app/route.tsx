import { Link, Outlet, createFileRoute, useNavigate, useLocation } from "@tanstack/react-router";
import * as React from "react";
import {
  BarChart3Icon,
  CalendarIcon,
  FolderIcon,
  LayoutGridIcon,
  LineChartIcon,
  LogOut,
  ShieldCheckIcon,
  ShieldIcon,
} from "lucide-react";

import { useAuth } from "@/hooks/useAuth";
import { ThemeToggle } from "@/components/theme-toggle";
import { Button } from "@/components/ui/button";
import {
  Sidebar,
  SidebarCloseButton,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
  SidebarCollapseTrigger,
  useSidebar,
} from "@/components/ui/sidebar";
import Loader from "@/components/loader";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/app")({
  component: AppLayout,
});

type NavItem = {
  to: string;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  exact?: boolean;
};

const ADMIN_NAV_ITEMS: NavItem[] = [
  { to: "/app", label: "Dashboard", icon: BarChart3Icon, exact: true },
  { to: "/app/analytics", label: "Analytics", icon: LineChartIcon },
  { to: "/app/risks", label: "Risks", icon: ShieldIcon },
  { to: "/app/board", label: "Board", icon: LayoutGridIcon },
  { to: "/app/calendar", label: "Calendar", icon: CalendarIcon },
  { to: "/app/frameworks", label: "Frameworks", icon: ShieldCheckIcon },
  { to: "/app/categories", label: "Categories", icon: FolderIcon },
];

const MEMBER_NAV_ITEMS: NavItem[] = [
  { to: "/app", label: "Dashboard", icon: BarChart3Icon, exact: true },
  { to: "/app/analytics", label: "Analytics", icon: LineChartIcon },
  { to: "/app/risks", label: "Risks", icon: ShieldIcon },
  { to: "/app/board", label: "Board", icon: LayoutGridIcon },
  { to: "/app/calendar", label: "Calendar", icon: CalendarIcon },
];

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
      <div className="flex h-dvh items-center justify-center bg-background">
        <Loader />
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  return (
    <SidebarProvider defaultOpen>
      <AppLayoutContent user={user} logout={logout} />
    </SidebarProvider>
  );
}

function AppLayoutContent({ user, logout }: { user: any; logout: () => void }) {
  const navItems = user?.role === "admin" ? ADMIN_NAV_ITEMS : MEMBER_NAV_ITEMS;
  const { collapsed } = useSidebar();
  const location = useLocation();

  const getPageTitle = (pathname: string) => {
    if (pathname === "/app") return "Dashboard";
    if (pathname.startsWith("/app/analytics")) return "Analytics";
    if (pathname.startsWith("/app/board")) return "Board";
    if (pathname.startsWith("/app/calendar")) return "Calendar";
    if (pathname.startsWith("/app/risks")) return "Risks";
    if (pathname.startsWith("/app/frameworks")) return "Frameworks";
    if (pathname.startsWith("/app/categories")) return "Categories";
    return "Risk Register";
  };

  const pageTitle = getPageTitle(location.pathname);

  return (
    <>
      {/* Fixed Sidebar */}
      <Sidebar className="border-r">
        <SidebarHeader>
          <div className={cn("flex w-full items-center", collapsed ? "justify-center" : "justify-between")}>
            <Link to="/app" className="flex items-center gap-2.5">
              <span className="bg-primary text-primary-foreground inline-flex size-8 items-center justify-center rounded-lg shadow-sm">
                <ShieldCheckIcon className="size-4" />
              </span>
              {!collapsed && (
                <span className="text-sm font-semibold tracking-tight">Risk Register</span>
              )}
            </Link>
            {/* Only show collapse trigger in sidebar header when expanded */}
            {!collapsed && <SidebarCollapseTrigger />}
          </div>
        </SidebarHeader>

        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Navigation</SidebarGroupLabel>
            <SidebarMenu>
              {navItems.map((item) => (
                <SidebarMenuItem key={item.to}>
                  <SidebarMenuButton asChild isActive={false} tooltip={item.label}>
                    <Link
                      to={item.to}
                      activeProps={{ className: "bg-sidebar-accent text-sidebar-accent-foreground" }}
                      inactiveProps={{
                        className:
                          "text-sidebar-foreground/70 hover:bg-sidebar-accent/50 hover:text-sidebar-foreground",
                      }}
                      activeOptions={item.exact ? { exact: true } : undefined}
                    >
                      <item.icon className="size-[18px] shrink-0" />
                      {!collapsed && <span>{item.label}</span>}
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroup>
        </SidebarContent>

        <SidebarFooter>
          <div className={cn("flex w-full flex-col gap-3", collapsed && "items-center")}>
            {/* Theme Toggle */}
            <div className={cn("flex items-center", collapsed ? "justify-center" : "justify-between")}>
              {!collapsed && <span className="text-xs text-sidebar-foreground/50">Theme</span>}
              <ThemeToggle />
            </div>

            {/* User Info and Logout */}
            <div
              className={cn(
                "flex w-full items-center pt-3 border-t border-sidebar-border/50",
                collapsed ? "justify-center" : "justify-between",
              )}
            >
              {!collapsed && (
                <div className="flex flex-col">
                  <span className="text-sm font-medium">{user?.name}</span>
                  <span className="text-xs text-sidebar-foreground/50 capitalize">{user?.role}</span>
                </div>
              )}
              <Button
                variant="outline"
                size={collapsed ? "icon-sm" : "sm"}
                onClick={logout}
                title="Log out"
                className={cn(
                  "shrink-0",
                  collapsed && "size-8"
                )}
              >
                {collapsed ? <LogOut className="size-4" /> : "Log out"}
              </Button>
            </div>

            {/* Mobile Close */}
            <div className={cn("flex items-center", collapsed ? "justify-center" : "justify-end", "md:hidden")}>
              <SidebarCloseButton />
            </div>
          </div>
        </SidebarFooter>
      </Sidebar>

      {/* Main Content Area */}
      <SidebarInset>
        {/* Top Header Bar */}
        <header className="sticky top-0 z-30 flex h-16 items-center gap-2 border-b bg-background/80 backdrop-blur-sm px-4">
          {/* Mobile: Menu toggle | Desktop: Collapse trigger when collapsed */}
          <SidebarTrigger />
          {/* Only show collapse trigger in header when sidebar is collapsed */}
          {collapsed && <SidebarCollapseTrigger />}

          <div className="h-4 w-px bg-border mx-1 hidden md:block" />
          <span className="text-sm font-medium">{pageTitle}</span>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-auto">
          <div className="container py-6">
            <Outlet />
          </div>
        </main>
      </SidebarInset>
    </>
  );
}
