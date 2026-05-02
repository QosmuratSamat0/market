import { createFileRoute, Link } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { SiteHeader } from "@/components/site-header";
import { SiteFooter } from "@/components/site-footer";
import { eur } from "@/lib/data";
import { api, type PaymentStatus } from "@/lib/api";
import { Package, MapPin, CreditCard, LogOut } from "lucide-react";

export const Route = createFileRoute("/payments")({
  head: () => ({ meta: [{ title: "Payments — Indigo Market" }] }),
  component: Payments,
});

const statusStyles: Record<PaymentStatus, string> = {
  pending: "bg-secondary text-foreground",
  processing: "bg-accent/15 text-accent",
  succeeded: "bg-primary/10 text-primary",
  failed: "bg-destructive/10 text-destructive",
  canceled: "bg-secondary text-muted-foreground",
  refunded: "bg-accent/10 text-accent",
};

function Payments() {
  const {
    data: payments = [],
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ["payments", "my"],
    queryFn: api.getMyPayments,
    retry: false,
  });

  return (
    <div className="min-h-screen">
      <SiteHeader />
      <section className="container-x pt-16 pb-10">
        <span className="text-xs uppercase tracking-[0.25em] text-muted-foreground">Account</span>
        <h1 className="mt-4 font-display text-5xl md:text-6xl">Your payments</h1>
      </section>

      <section className="container-x grid gap-10 pb-24 md:grid-cols-[240px_1fr]">
        <aside className="space-y-1">
          <NavItem icon={<Package className="h-4 w-4" />} to="/orders" label="Orders" />
          <NavItem icon={<MapPin className="h-4 w-4" />} to="/account" label="Profile" />
          <NavItem icon={<CreditCard className="h-4 w-4" />} to="/payments" label="Payment" active />
        </aside>

        <div className="space-y-6">
          {isLoading && (
            <div className="py-16 text-center text-muted-foreground">Loading payments...</div>
          )}
          {isError && (
            <div className="rounded-lg border border-border/60 bg-card p-6 text-sm text-muted-foreground">
              {(error as Error).message}
            </div>
          )}
          {!isLoading && !isError && payments.length === 0 && (
            <div className="rounded-lg border border-border/60 bg-card p-6 text-sm text-muted-foreground">
              You have no payments yet.
            </div>
          )}
          {payments.map((p: any) => (
            <article key={p.id} className="rounded-lg border border-border/60 bg-card p-6 flex items-center justify-between gap-6">
              <div>
                <div className="font-display text-lg">Transaction #{p.id.slice(0, 8)}</div>
                <div className="text-xs uppercase tracking-widest text-muted-foreground mt-1">
                  Order: {p.orderId.slice(0, 8)} · {new Date(p.createdAt).toLocaleString()}
                </div>
              </div>
              <div className="flex items-center gap-4">
                <span
                  className={`rounded-full px-3 py-1 text-[10px] font-medium uppercase tracking-widest ${statusStyles[p.status as PaymentStatus]}`}
                >
                  {p.status}
                </span>
                <div className="text-lg font-medium">{eur(p.amount)}</div>
              </div>
            </article>
          ))}
        </div>
      </section>
      <SiteFooter />
    </div>
  );
}

function NavItem({ icon, to, label, active }: { icon: React.ReactNode; to: string; label: string; active?: boolean }) {
  return (
    <Link to={to} className={`flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors ${active ? "bg-secondary text-foreground" : "text-muted-foreground hover:bg-secondary hover:text-foreground"}`}>
      {icon} {label}
    </Link>
  );
}
