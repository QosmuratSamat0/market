import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { SiteHeader } from "@/components/site-header";
import { SiteFooter } from "@/components/site-footer";
import { Field } from "@/components/auth-shell";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useMemo, useState } from "react";
import { api, type Product } from "@/lib/api";
import { clearCart, readCart, type CartLine } from "@/lib/cart";
import { eur } from "@/lib/data";
import { toast } from "sonner";

export const Route = createFileRoute("/checkout")({
  head: () => ({ meta: [{ title: "Checkout — Indigo Market" }] }),
  component: Checkout,
});

function Checkout() {
  const navigate = useNavigate();
  const [lines, setLines] = useState<CartLine[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isAuthorized, setIsAuthorized] = useState(false);
  const { data: products = [] } = useQuery({
    queryKey: ["products"],
    queryFn: api.getProducts,
  });
  const { data: profile } = useQuery({
    queryKey: ["me"],
    queryFn: api.getMe,
    enabled: isAuthorized,
    retry: false,
  });

  useEffect(() => {
    const token = typeof window !== "undefined" ? localStorage.getItem("token") : null;
    if (!token) {
      setIsAuthorized(false);
      toast.error("Please sign in or create an account before checkout");
      navigate({ to: "/sign-in" });
      return;
    }

    setIsAuthorized(true);
    setLines(readCart());
  }, [navigate]);

  const productByID = useMemo(
    () => new Map(products.map((product: Product) => [product.id, product])),
    [products],
  );
  const items = lines
    .map((line) => {
      const product = productByID.get(line.id);
      return product ? { ...line, product } : null;
    })
    .filter((item): item is CartLine & { product: Product } => item !== null);
  const subtotal = items.reduce((sum, item) => sum + item.product.price * item.qty, 0);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isAuthorized) {
      toast.error("Please sign in first");
      navigate({ to: "/sign-in" });
      return;
    }
    if (items.length === 0) {
      toast.error("Cart is empty");
      return;
    }

    try {
      setIsSubmitting(true);
      const order = await api.createOrder({
        items: items.map((item) => ({ product_id: item.product.id, quantity: item.qty })),
      });
      const currentUser = profile || (await api.getMe());

      const payment = await api.createPayment({
        order_id: order.id,
        user_id: order.userId,
        amount: Math.round(order.total || subtotal),
        currency: "EUR",
        provider: "mock",
        idempotency_key: order.id,
      });

      clearCart();
      toast.success("Order created, redirecting to payment...");

      if (payment.paymentUrl) {
        window.location.href = payment.paymentUrl;
      } else {
        navigate({ to: "/orders" });
      }
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Failed to place order");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen">
      <SiteHeader />
      <section className="container-x pt-16 pb-10">
        <Link
          to="/cart"
          className="text-xs uppercase tracking-widest text-muted-foreground hover:text-foreground"
        >
          ← Back to cart
        </Link>
        <h1 className="mt-4 font-display text-5xl md:text-6xl">Checkout</h1>
      </section>

      <section className="container-x grid gap-12 pb-24 md:grid-cols-3">
        <form onSubmit={handleSubmit} className="space-y-10 md:col-span-2">
          <fieldset className="space-y-6">
            <legend className="font-display text-2xl">Contact</legend>
            <Field
              label="Email"
              type="email"
              autoComplete="email"
              placeholder="you@email.com"
              defaultValue={profile?.email}
            />
          </fieldset>
          <fieldset className="space-y-6">
            <legend className="font-display text-2xl">Shipping</legend>
            <div className="grid gap-6 md:grid-cols-2">
              <Field label="First name" autoComplete="given-name" />
              <Field label="Last name" autoComplete="family-name" />
            </div>
            <Field label="Address" autoComplete="street-address" />
            <div className="grid gap-6 md:grid-cols-3">
              <Field label="City" autoComplete="address-level2" />
              <Field label="Postal code" autoComplete="postal-code" />
              <Field label="Country" autoComplete="country" />
            </div>
          </fieldset>
          <fieldset className="space-y-6">
            <legend className="font-display text-2xl">Payment</legend>
            <Field label="Card number" placeholder="1234 5678 9012 3456" />
            <div className="grid gap-6 md:grid-cols-2">
              <Field label="Expiry" placeholder="MM / YY" />
              <Field label="CVC" placeholder="123" />
            </div>
          </fieldset>
          <button
            disabled={isSubmitting || items.length === 0}
            className="w-full rounded-md bg-primary px-6 py-4 text-xs font-medium uppercase tracking-widest text-primary-foreground transition-colors hover:bg-accent disabled:cursor-not-allowed disabled:opacity-60"
          >
            {isSubmitting ? "Placing order..." : "Place order"}
          </button>
        </form>

        <aside className="h-fit rounded-lg border border-border/60 bg-secondary/40 p-6">
          <h2 className="font-display text-2xl">Order summary</h2>
          <dl className="mt-6 space-y-3 text-sm">
            <div className="flex justify-between">
              <dt className="text-muted-foreground">Subtotal</dt>
              <dd>{eur(subtotal)}</dd>
            </div>
            <div className="flex justify-between">
              <dt className="text-muted-foreground">Shipping</dt>
              <dd>Free</dd>
            </div>
            <div className="flex justify-between border-t border-border/60 pt-3 text-base font-medium">
              <dt>Total</dt>
              <dd>{eur(subtotal)}</dd>
            </div>
          </dl>
        </aside>
      </section>
      <SiteFooter />
    </div>
  );
}
