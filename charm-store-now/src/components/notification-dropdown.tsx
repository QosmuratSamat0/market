import { Bell, Check, Trash2 } from "lucide-react";
import { formatDistanceToNow } from "date-fns";
import { Link } from "@tanstack/react-router";
import { useNotifications } from "@/hooks/use-notifications";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function NotificationDropdown() {
  const { notifications, unreadCount, markAsRead, markAllAsRead, clearAll } = useNotifications();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button
          aria-label="Notifications"
          className="relative rounded-full p-2 transition-colors hover:bg-secondary focus:outline-none"
        >
          <Bell className="h-4 w-4" />
          {unreadCount > 0 && (
            <span className="absolute right-2 top-2 flex h-2 w-2">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-accent opacity-75"></span>
              <span className="relative inline-flex h-2 w-2 rounded-full bg-accent"></span>
            </span>
          )}
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-80">
        <div className="flex items-center justify-between px-4 py-2">
          <DropdownMenuLabel className="p-0 font-medium">Notifications</DropdownMenuLabel>
          {unreadCount > 0 && (
            <Button
              variant="ghost"
              size="sm"
              className="h-auto p-0 text-xs text-accent hover:bg-transparent"
              onClick={markAllAsRead}
            >
              Mark all as read
            </Button>
          )}
        </div>
        <DropdownMenuSeparator />
        <div className="max-h-96 overflow-y-auto">
          {notifications.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-8 text-center text-muted-foreground">
              <Bell className="mb-2 h-8 w-8 opacity-20" />
              <p className="text-sm">No notifications yet</p>
            </div>
          ) : (
            notifications.map((n) => (
              <DropdownMenuItem
                key={n.id}
                className={cn(
                  "flex flex-col items-start gap-1 p-4 focus:bg-secondary",
                  !n.read && "bg-secondary/30"
                )}
                onClick={() => markAsRead(n.id)}
              >
                <div className="flex w-full items-center justify-between">
                  <span className={cn("text-xs font-semibold", n.read ? "text-muted-foreground" : "text-foreground")}>
                    {n.title}
                  </span>
                  <span className="text-[10px] text-muted-foreground">
                    {formatDistanceToNow(n.time, { addSuffix: true })}
                  </span>
                </div>
                <p className={cn("text-xs line-clamp-2", n.read ? "text-muted-foreground" : "text-foreground")}>
                  {n.message}
                </p>
              </DropdownMenuItem>
            ))
          )}
        </div>
        <DropdownMenuSeparator />
        {notifications.length > 0 && (
          <div className="p-1">
            <Button
              variant="ghost"
              size="sm"
              className="w-full justify-start text-xs text-destructive hover:text-destructive hover:bg-destructive/10"
              onClick={clearAll}
            >
              <Trash2 className="mr-2 h-3 w-3" />
              Clear all
            </Button>
          </div>
        )}
        <DropdownMenuSeparator />
        <div className="p-1">
          <Link to="/notifications">
            <Button variant="ghost" size="sm" className="w-full text-xs">
              View all notifications
            </Button>
          </Link>
        </div>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
