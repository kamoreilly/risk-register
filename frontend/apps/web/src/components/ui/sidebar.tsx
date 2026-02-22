import * as React from "react";
import { MenuIcon, PanelLeftCloseIcon, PanelLeftIcon, XIcon } from "lucide-react";

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
      {/* Mobile overlay - glass effect */}
      <div
        aria-hidden="true"
        className={cn(
          "fixed inset-0 z-40 bg-black/20 backdrop-blur-[2px] transition-all duration-300 ease-out md:hidden",
          open ? "opacity-100" : "pointer-events-none opacity-0",
        )}
        onClick={() => setOpen(false)}
      />

      {/* Sidebar - fixed on all screens */}
      <aside
        className={cn(
          "bg-sidebar text-sidebar-foreground border-sidebar-border fixed left-0 top-0 z-50 flex h-dvh flex-col border-r transition-[width,transform] duration-300 ease-out",
          // Mobile: slide in/out
          "translate-x-0 md:translate-x-0",
          !open && "-translate-x-full md:translate-x-0",
          // Width states
          collapsed ? "w-[68px]" : "w-64",
          // Subtle shadow when open on mobile
          open && "shadow-xl md:shadow-none",
          className,
        )}
        {...props}
      >
        {children}
      </aside>
    </>
  );
}

function SidebarInset({ className, ...props }: React.ComponentProps<"div">) {
  const { collapsed } = useSidebar();

  return (
    <div
      className={cn(
        "flex min-h-dvh flex-1 flex-col transition-[margin-left] duration-300 ease-out",
        // Add margin on desktop to account for fixed sidebar
        collapsed ? "md:ml-[68px]" : "md:ml-64",
        className,
      )}
      {...props}
    />
  );
}

function SidebarHeader({ className, ...props }: React.ComponentProps<"div">) {
  const { collapsed } = useSidebar();
  return (
    <div
      className={cn(
        "flex h-16 items-center gap-2 border-b border-sidebar-border/50 px-4",
        collapsed && "justify-center px-3",
        className,
      )}
      {...props}
    />
  );
}

function SidebarContent({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn(
        "flex flex-1 flex-col gap-2 overflow-y-auto overflow-x-hidden px-3 py-4",
        "scrollbar-thin scrollbar-track-transparent scrollbar-thumb-sidebar-border/30 hover:scrollbar-thumb-sidebar-border/50",
        className,
      )}
      {...props}
    />
  );
}

function SidebarFooter({ className, ...props }: React.ComponentProps<"div">) {
  const { collapsed } = useSidebar();
  return (
    <div
      className={cn(
        "border-t border-sidebar-border/50 px-3 py-4",
        collapsed && "flex flex-col items-center",
        className,
      )}
      {...props}
    />
  );
}

function SidebarGroup({ className, ...props }: React.ComponentProps<"div">) {
  return <div className={cn("grid gap-1", className)} {...props} />;
}

function SidebarGroupLabel({ className, ...props }: React.ComponentProps<"div">) {
  const { collapsed } = useSidebar();
  if (collapsed) return null;
  return (
    <div
      className={cn(
        "px-2 pb-2 pt-1 text-[10px] font-bold uppercase tracking-widest text-sidebar-foreground/40",
        className,
      )}
      {...props}
    />
  );
}

function SidebarMenu({ className, ...props }: React.ComponentProps<"div">) {
  return <div className={cn("grid gap-0.5", className)} {...props} />;
}

function SidebarMenuItem({ className, ...props }: React.ComponentProps<"div">) {
  return <div className={cn("", className)} {...props} />;
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
  const [showTooltip, setShowTooltip] = React.useState(false);

  const baseClassName = cn(
    "group relative flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-all duration-200",
    isActive
      ? "bg-sidebar-accent text-sidebar-accent-foreground shadow-sm"
      : "text-sidebar-foreground/70 hover:bg-sidebar-accent/50 hover:text-sidebar-foreground",
    collapsed && "justify-center px-2.5",
    className,
  );

  // Tooltip for collapsed state
  const tooltipEl = collapsed && tooltip && showTooltip && (
    <div
      className={cn(
        "absolute left-full top-1/2 z-[100] ml-3 -translate-y-1/2",
        "animate-in fade-in-0 zoom-in-95 duration-150",
        "whitespace-nowrap rounded-md bg-popup-foreground px-3 py-1.5 text-xs font-medium text-popup shadow-lg",
      )}
    >
      {tooltip}
      <div className="absolute right-full top-1/2 -translate-y-1/2 border-4 border-transparent border-r-popup-foreground" />
    </div>
  );

  if (asChild && React.isValidElement<{ className?: string; children?: React.ReactNode }>(children)) {
    return (
      <div
        className="relative"
        onMouseEnter={() => setShowTooltip(true)}
        onMouseLeave={() => setShowTooltip(false)}
      >
        {React.cloneElement(children, {
          className: cn(baseClassName, children.props.className),
        })}
        {tooltipEl}
      </div>
    );
  }

  return (
    <div
      className="relative"
      onMouseEnter={() => setShowTooltip(true)}
      onMouseLeave={() => setShowTooltip(false)}
    >
      <button type="button" className={baseClassName} {...props}>
        {children}
      </button>
      {tooltipEl}
    </div>
  );
}

/**
 * Mobile menu toggle - hamburger/X icon
 * Only visible on mobile (< md breakpoint)
 */
function SidebarTrigger({ className, ...props }: React.ComponentProps<"button">) {
  const { toggle, open } = useSidebar();

  return (
    <Button
      aria-label={open ? "Close sidebar" : "Open sidebar"}
      variant="ghost"
      size="icon"
      className={cn(
        "md:hidden",
        className,
      )}
      onClick={toggle}
      {...props}
    >
      {open ? <XIcon className="size-5" /> : <MenuIcon className="size-5" />}
    </Button>
  );
}

/**
 * Desktop collapse toggle - panel icon
 * Only visible on desktop (>= md breakpoint)
 */
function SidebarCollapseTrigger({ className, ...props }: React.ComponentProps<"button">) {
  const { toggleCollapse, collapsed } = useSidebar();

  return (
    <Button
      aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
      variant="ghost"
      size="icon"
      className={cn(
        "hidden md:flex",
        className,
      )}
      onClick={toggleCollapse}
      {...props}
    >
      {collapsed ? (
        <PanelLeftIcon className="size-5" />
      ) : (
        <PanelLeftCloseIcon className="size-5" />
      )}
    </Button>
  );
}

/**
 * Close button for mobile sidebar footer
 */
function SidebarCloseButton({ className, ...props }: React.ComponentProps<"button">) {
  const { setOpen } = useSidebar();

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={() => setOpen(false)}
      className={cn("text-sidebar-foreground/70 hover:text-sidebar-foreground", className)}
      {...props}
    >
      <XIcon className="size-4 mr-2" />
      Close
    </Button>
  );
}

export {
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
};
