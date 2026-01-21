import { NavigationMenu as NavigationMenuPrimitive } from "@base-ui/react/navigation-menu";
import { cva } from "class-variance-authority";
import * as React from "react";

import { cn } from "@/lib/utils";

type NavigationMenuProps = React.ComponentPropsWithoutRef<typeof NavigationMenuPrimitive.Root>;
type NavigationMenuListProps = React.ComponentPropsWithoutRef<typeof NavigationMenuPrimitive.List>;
type NavigationMenuItemProps = React.ComponentPropsWithoutRef<typeof NavigationMenuPrimitive.Item>;
type NavigationMenuLinkProps = React.ComponentPropsWithoutRef<typeof NavigationMenuPrimitive.Link>;

function NavigationMenu({ className, ...props }: NavigationMenuProps) {
  return (
    <NavigationMenuPrimitive.Root
      data-slot="navigation-menu"
      className={cn("relative z-10", className)}
      {...props}
    />
  );
}

function NavigationMenuList({ className, ...props }: NavigationMenuListProps) {
  return (
    <NavigationMenuPrimitive.List
      data-slot="navigation-menu-list"
      className={cn("flex items-center gap-1", className)}
      {...props}
    />
  );
}

function NavigationMenuItem({ ...props }: NavigationMenuItemProps) {
  return <NavigationMenuPrimitive.Item data-slot="navigation-menu-item" {...props} />;
}

const navigationMenuLinkStyle = cva(
  "focus-visible:border-ring focus-visible:ring-ring/50 inline-flex h-8 items-center justify-center whitespace-nowrap border border-transparent bg-clip-padding px-2.5 text-xs font-medium transition-all outline-none focus-visible:ring-1",
  {
    variants: {
      variant: {
        default: "hover:bg-muted hover:text-foreground",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  },
);

function NavigationMenuLink({ className, ...props }: NavigationMenuLinkProps) {
  return (
    <NavigationMenuPrimitive.Link
      data-slot="navigation-menu-link"
      className={cn(navigationMenuLinkStyle(), className)}
      {...props}
    />
  );
}

export { NavigationMenu, NavigationMenuItem, NavigationMenuLink, NavigationMenuList, navigationMenuLinkStyle };
