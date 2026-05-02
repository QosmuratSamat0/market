import { createFileRoute } from "@tanstack/react-router";
import { Bell, Check, Trash2, Clock } from "lucide-react";
import { formatDistanceToNow } from "date-fns";
import { useNotifications } from "@/hooks/use-notifications";
import { SiteHeader } from "@/components/site-header";
import { SiteFooter } from "@/components/site-footer";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/notifications")({
  head: () => ({ meta: [{ title: "Notifications — Indigo Market" }] }),
  component: NotificationsPage,
});

function NotificationsPage() {
  const { notifications, markAsRead, markAllAsRead, clearAll } = useNotifications();

  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      
      <main className="container-x pt-16 pb-24">
        <div className="mb-8 flex flex-wrap items-end justify-between gap-4">
          <div>
            <span className="text-xs uppercase tracking-[0.25em] text-muted-foreground">Activity</span>
            <h1 className="mt-4 font-display text-5xl md:text-6xl">Notifications</h1>
          </div>
          
          <div className="flex gap-2">
            {notifications.length > 0 && (
              <>
                <Button variant="outline" size="sm" onClick={markAllAsRead}>
                  <Check className="mr-2 h-4 w-4" />
                  Mark all as read
                </Button>
                <Button variant="outline" size="sm" className="text-destructive hover:bg-destructive/10" onClick={clearAll}>
                  <Trash2 className="mr-2 h-4 w-4" />
                  Clear all
                </Button>
              </>
            )}
          </div>
        </div>

        <div className="mx-auto max-w-3xl space-y-4">
          {notifications.length === 0 ? (
            <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed border-border/60 bg-card/30 py-24 text-center">
              <div className="mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-secondary">
                <Bell className="h-8 w-8 text-muted-foreground opacity-50" />
              </div>
              <h3 className="text-xl font-medium">All caught up!</h3>
              <p className="mt-2 text-muted-foreground">
                You don't have any notifications at the moment.
              </p>
            </div>
          ) : (
            notifications.map((n) => (
              <article
                key={n.id}
                className={cn(
                  "group relative overflow-hidden rounded-xl border border-border/60 p-6 transition-all hover:shadow-sm",
                  n.read ? "bg-card/50 opacity-80" : "bg-card shadow-sm ring-1 ring-accent/20"
                )}
                onClick={() => markAsRead(n.id)}
              >
                {!n.read && (
                  <div className="absolute left-0 top-0 h-full w-1 bg-accent" />
                )}
                
                <div className="flex items-start gap-4">
                  <div className={cn(
                    "mt-1 flex h-10 w-10 shrink-0 items-center justify-center rounded-full",
                    n.type === "success" ? "bg-primary/10 text-primary" : "bg-secondary text-foreground"
                  )}>
                    <Bell className="h-5 w-5" />
                  </div>
                  
                  <div className="flex-1 space-y-1">
                    <div className="flex items-center justify-between gap-2">
                      <h3 className={cn("font-medium", !n.read && "text-foreground")}>
                        {n.title}
                      </h3>
                      <time className="flex items-center gap-1 text-[10px] uppercase tracking-widest text-muted-foreground">
                        <Clock className="h-3 w-3" />
                        {formatDistanceToNow(n.time, { addSuffix: true })}
                      </time>
                    </div>
                    <p className="text-sm text-muted-foreground leading-relaxed">
                      {n.message}
                    </p>
                  </div>
                </div>
              </article>
            ))
          )}
        </div>
      </main>

      <SiteFooter />
    </div>
  );
}
