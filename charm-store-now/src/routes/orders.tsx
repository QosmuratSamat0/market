import { createFileRoute, Link } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { SiteHeader } from "@/components/site-header";
import { SiteFooter } from "@/components/site-footer";
import { eur } from "@/lib/data";
import { api, type Order } from "@/lib/api";

export const Route = createFileRoute("/orders")({
  head: () => ({ meta: [{ title: "Orders — Indigo Market" }] }),
  component: Orders,
});

const statusStyles: Record<Order["status"], string> = {
  pending: "bg-secondary text-foreground",
  paid: "bg-primary/10 text-primary",
  shipped: "bg-accent/15 text-accent",
  delivered: "bg-primary/10 text-primary",
  cancelled: "bg-destructive/10 text-destructive",
};

function Orders() {
  const {
    data: orders = [],
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ["orders", "my"],
    queryFn: api.getMyOrders,
    retry: false,
  });

  return (
    <div className="min-h-screen">
      <SiteHeader />
      <section className="container-x pt-16 pb-10">
        <span className="text-xs uppercase tracking-[0.25em] text-muted-foreground">Account</span>
        <h1 className="mt-4 font-display text-5xl md:text-6xl">Your orders</h1>
      </section>
      <section className="container-x space-y-6 pb-24">
        {isLoading && (
          <div className="py-16 text-center text-muted-foreground">Loading orders...</div>
        )}
        {isError && (
          <div className="rounded-lg border border-border/60 bg-card p-6 text-sm text-muted-foreground">
            {(error as Error).message}
          </div>
        )}
        {!isLoading && !isError && orders.length === 0 && (
          <div className="rounded-lg border border-border/60 bg-card p-6 text-sm text-muted-foreground">
            You have no orders yet.
          </div>
        )}
        {orders.map((o) => (
          <article key={o.id} className="rounded-lg border border-border/60 bg-card p-6">
            <header className="flex flex-wrap items-center justify-between gap-4 border-b border-border/60 pb-4">
              <div>
                <div className="font-display text-lg">{o.id}</div>
                <div className="text-xs uppercase tracking-widest text-muted-foreground">
                  {new Date(o.createdAt).toLocaleDateString()}
                </div>
              </div>
              <div className="flex items-center gap-4">
                <span
                  className={`rounded-full px-3 py-1 text-[10px] font-medium uppercase tracking-widest ${statusStyles[o.status]}`}
                >
                  {o.status}
                </span>
                <div className="text-sm font-medium">{eur(o.total)}</div>
              </div>
            </header>
            <ul className="mt-4 space-y-3">
              {o.items.map((it, i) => (
                <li key={i} className="flex items-center gap-4">
                  <div className="flex-1">
                    <div className="text-sm">{it.productId}</div>
                    <div className="text-xs text-muted-foreground">
                      Qty {it.quantity} · {eur(it.price)}
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          </article>
        ))}
        <div className="pt-4 text-center">
          <Link
            to="/shop"
            className="text-xs uppercase tracking-widest underline-offset-4 hover:underline"
          >
            Continue shopping →
          </Link>
        </div>
      </section>
      <SiteFooter />
    </div>
  );
}
