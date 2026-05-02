import { Link } from "@tanstack/react-router";
import { ShoppingBag, Search, User } from "lucide-react";
import { NotificationDropdown } from "@/components/notification-dropdown";

export function SiteHeader() {
  return (
    <header className="sticky top-0 z-40 border-b border-border/60 bg-background/85 backdrop-blur-md">
      <div className="container-x flex h-16 items-center justify-between">
        <Link to="/" className="font-display text-xl font-medium tracking-tight">
          indigo<span className="text-accent">.</span>
        </Link>
        <nav className="hidden items-center gap-8 text-sm md:flex">
          <Link to="/shop" className="text-muted-foreground transition-colors hover:text-foreground" activeProps={{ className: "text-foreground" }}>Shop</Link>
          <Link to="/categories" className="text-muted-foreground transition-colors hover:text-foreground" activeProps={{ className: "text-foreground" }}>Categories</Link>
          <Link to="/orders" className="text-muted-foreground transition-colors hover:text-foreground" activeProps={{ className: "text-foreground" }}>Orders</Link>
        </nav>
        <div className="flex items-center gap-1">
          <button aria-label="Search" className="rounded-full p-2 transition-colors hover:bg-secondary">
            <Search className="h-4 w-4" />
          </button>
          <NotificationDropdown />
          <Link to="/account" aria-label="Account" className="rounded-full p-2 transition-colors hover:bg-secondary">
            <User className="h-4 w-4" />
          </Link>
          <Link to="/cart" aria-label="Cart" className="relative rounded-full p-2 transition-colors hover:bg-secondary">
            <ShoppingBag className="h-4 w-4" />
            <span className="absolute right-1 top-1 h-1.5 w-1.5 rounded-full bg-accent" />
          </Link>
        </div>
      </div>
    </header>
  );
}
