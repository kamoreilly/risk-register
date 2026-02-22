import * as React from "react";
import { MenuIcon, ChevronLeftIcon, ChevronRightIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

type SidebarContextValue = {
  open: boolean;
  setOpen: (open: boolean) => void;
  toggle: () => void;
  collapsed: boolean;
  setCollapsed: (collapsed: boolean) => void;
  toggleCollapse: () => void;
};

const SidebarContext = React.createContext<SidebarContextValue | null>(null);

function SidebarProvider({
  defaultOpen = false,
  defaultCollapsed = false,
  children,
}: {
  defaultOpen?: boolean;
  defaultCollapsed?: boolean;
  children: React.ReactNode;
}) {
  const [open, setOpen] = React.useState(defaultOpen);
  const [collapsed, setCollapsed] = React.useState(defaultCollapsed);

  const toggle = React.useCallback(() => {
    setOpen((prev) => !prev);
  }, []);

  const toggleCollapse = React.useCallback(() => {
    setCollapsed((prev) => !prev);
  }, []);

  return (
    <SidebarContext.Provider value={{ open, setOpen, toggle, collapsed, setCollapsed, toggleCollapse }}>
      {children}
    </SidebarContext.Provider>
  );
}

function useSidebar() {
  const context = React.useContext(SidebarContext);
  if (!context) {
    throw new Error("useSidebar must be used within SidebarProvider");
  }
  return context;
}

function Sidebar({ className, children, ...props }: React.ComponentProps<"aside">) {
  const { open, setOpen, collapsed } = useSidebar();

  return (
    <>
      <aside
        className={cn(
          "bg-sidebar text-sidebar-foreground border-sidebar-border fixed inset-y-0 left-0 z-40 flex flex-col border-r transition-[width,transform] duration-200 md:static md:translate-x-0",
          open ? "translate-x-0" : "-translate-x-full",
          collapsed ? "w-16" : "w-64",
          className,
        )}
        {...props}
      >
        {children}
      </aside>
      <div
        aria-hidden="true"
        className={cn(
          "fixed inset-0 z-30 bg-background/80 backdrop-blur-sm transition-opacity md:hidden",
          open ? "opacity-100" : "pointer-events-none opacity-0",
        )}
        onClick={() => setOpen(false)}
      />
    </>
  );
}

function SidebarInset({ className, ...props }: React.ComponentProps<"div">) {
  return <div className={cn("flex min-h-svh flex-1 flex-col", className)} {...props} />;
}

function SidebarHeader({ className, ...props }: React.ComponentProps<"div">) {
  const { collapsed } = useSidebar();
  return <div className={cn("flex h-14 items-center gap-2 border-b border-sidebar-border px-4", collapsed && "justify-center px-2", className)} {...props} />;
}

function SidebarContent({ className, ...props }: React.ComponentProps<"div">) {
  return <div className={cn("flex flex-1 flex-col gap-4 overflow-y-auto px-2 py-4", className)} {...props} />;
}

function SidebarFooter({ className, ...props }: React.ComponentProps<"div">) {
  const { collapsed } = useSidebar();
  return <div className={cn("border-t border-sidebar-border px-2 py-4", collapsed && "flex flex-col items-center", className)} {...props} />;
}

function SidebarGroup({ className, ...props }: React.ComponentProps<"div">) {
  return <div className={cn("grid gap-2", className)} {...props} />;
}

function SidebarGroupLabel({ className, ...props }: React.ComponentProps<"div">) {
  const { collapsed } = useSidebar();
  if (collapsed) return null;
  return (
    <div
      className={cn("px-3 text-[11px] font-semibold uppercase tracking-wide text-muted-foreground", className)}
      {...props}
    />
  );
}

function SidebarMenu({ className, ...props }: React.ComponentProps<"div">) {
  return <div className={cn("grid gap-1", className)} {...props} />;
}

function SidebarMenuItem({ className, ...props }: React.ComponentProps<"div">) {
  return <div className={cn("px-1", className)} {...props} />;
}

function SidebarMenuButton({
  asChild,
  isActive = false,
  className,
  children,
  tooltip,
  ...props
}: React.ComponentProps<"button"> & {
  asChild?: boolean;
  isActive?: boolean;
  tooltip?: string;
}) {
  const { collapsed } = useSidebar();
  const baseClassName = cn(
    "flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm font-medium transition-colors",
    isActive
      ? "bg-sidebar-accent text-sidebar-accent-foreground"
      : "text-sidebar-foreground/80 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground",
    collapsed && "justify-center px-2",
    className,
  );

  if (asChild && React.isValidElement<{ className?: string; children?: React.ReactNode }>(children)) {
    // When collapsed, we want to hide the text content (usually in a span)
    // This is a bit tricky with asChild. Usually the child is a Link which contains Icon and Span.
    // We can use CSS to hide the span when collapsed, or rely on the consumer to handle it.
    // But since we are modifying the library component, let's try to be smart.
    // Actually, simple CSS solution: "group-data-[collapsed=true]:hidden" if we add data attribute.
    // Or just use the passed className.
    
    return React.cloneElement(children, {
      className: cn(baseClassName, children.props.className),
      // We can't easily modify children of children here. 
      // Instead, we will wrap the content in a way that hides text.
    });
  }

  return (
    <button type="button" className={baseClassName} {...props}>
      {children}
    </button>
  );
}

function SidebarTrigger({ className, ...props }: React.ComponentProps<"button">) {
  const { toggle } = useSidebar();

  return (
    <Button
      aria-label="Toggle sidebar"
      variant="outline"
      size="icon-sm"
      className={cn("md:hidden", className)}
      onClick={toggle}
      {...props}
    >
      <MenuIcon className="size-4" />
    </Button>
  );
}

function SidebarCollapseTrigger({ className, ...props }: React.ComponentProps<"button">) {
  const { toggleCollapse, collapsed } = useSidebar();

  return (
    <Button
      aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
      variant="ghost"
      size="icon-sm"
      className={cn("hidden md:flex", className)}
      onClick={toggleCollapse}
      {...props}
    >
      {collapsed ? <ChevronRightIcon className="size-4" /> : <ChevronLeftIcon className="size-4" />}
    </Button>
  );
}

export {
  Sidebar,
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
};
