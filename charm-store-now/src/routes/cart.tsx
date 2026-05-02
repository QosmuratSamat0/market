import { createFileRoute, Link } from "@tanstack/react-router";
import { SiteHeader } from "@/components/site-header";
import { SiteFooter } from "@/components/site-footer";
import { eur } from "@/lib/data";
import { api, type Product } from "@/lib/api";
import { readCart, writeCart, type CartLine } from "@/lib/cart";
import { Minus, Plus, X } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";

export const Route = createFileRoute("/cart")({
  head: () => ({ meta: [{ title: "Cart — Indigo Market" }] }),
  component: Cart,
});

function Cart() {
  const [lines, setLines] = useState<CartLine[]>([]);
  const { data: products = [], isLoading } = useQuery({
    queryKey: ["products"],
    queryFn: api.getProducts,
  });

  useEffect(() => {
    setLines(readCart());
  }, []);

  useEffect(() => {
    writeCart(lines);
  }, [lines]);

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
  const subtotal = items.reduce((s, i) => s + i.product.price * i.qty, 0);
  const shipping = subtotal > 100 || subtotal === 0 ? 0 : 8;
  const total = subtotal + shipping;

  const setQty = (id: string, qty: number) =>
    setLines((ls) => ls.map((l) => (l.id === id ? { ...l, qty: Math.max(1, qty) } : l)));
  const remove = (id: string) => setLines((ls) => ls.filter((l) => l.id !== id));

  return (
    <div className="min-h-screen">
      <SiteHeader />
      <section className="container-x pt-16 pb-10">
        <span className="text-xs uppercase tracking-[0.25em] text-muted-foreground">Checkout</span>
        <h1 className="mt-4 font-display text-5xl md:text-6xl">Your cart</h1>
      </section>

      {isLoading ? (
        <section className="container-x py-20 text-center text-muted-foreground">
          Loading cart...
        </section>
      ) : items.length === 0 ? (
        <section className="container-x py-20 text-center">
          <p className="text-muted-foreground">Your cart is empty.</p>
          <Link
            to="/shop"
            className="mt-6 inline-block rounded-md bg-primary px-6 py-3 text-xs font-medium uppercase tracking-widest text-primary-foreground hover:bg-accent"
          >
            Continue shopping
          </Link>
        </section>
      ) : (
        <section className="container-x grid gap-12 pb-24 md:grid-cols-3">
          <div className="md:col-span-2 divide-y divide-border/60 border-y border-border/60">
            {items.map((i) => (
              <div key={i.id} className="flex gap-4 py-6">
                <div className="h-28 w-24 shrink-0 overflow-hidden rounded-md bg-secondary">
                  <img
                    src={i.product.image}
                    alt={i.product.name}
                    className="h-full w-full object-cover"
                  />
                </div>
                <div className="flex flex-1 flex-col">
                  <div className="flex items-start justify-between gap-2">
                    <Link
                      to="/products/$slug"
                      params={{ slug: i.product.id }}
                      className="font-display text-lg hover:text-accent"
                    >
                      {i.product.name}
                    </Link>
                    <button
                      onClick={() => remove(i.id)}
                      aria-label="Remove"
                      className="rounded-full p-1 text-muted-foreground hover:bg-secondary hover:text-foreground"
                    >
                      <X className="h-4 w-4" />
                    </button>
                  </div>
                  <div className="text-sm text-muted-foreground">{eur(i.product.price)}</div>
                  <div className="mt-auto flex items-center justify-between">
                    <div className="inline-flex items-center rounded-md border border-border">
                      <button
                        onClick={() => setQty(i.id, i.qty - 1)}
                        className="p-2 hover:bg-secondary"
                      >
                        <Minus className="h-3 w-3" />
                      </button>
                      <span className="w-8 text-center text-sm">{i.qty}</span>
                      <button
                        onClick={() => setQty(i.id, i.qty + 1)}
                        className="p-2 hover:bg-secondary"
                      >
                        <Plus className="h-3 w-3" />
                      </button>
                    </div>
                    <div className="font-medium">{eur(i.product.price * i.qty)}</div>
                  </div>
                </div>
              </div>
            ))}
          </div>

          <aside className="h-fit rounded-lg border border-border/60 bg-secondary/40 p-6">
            <h2 className="font-display text-2xl">Summary</h2>
            <dl className="mt-6 space-y-3 text-sm">
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Subtotal</dt>
                <dd>{eur(subtotal)}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Shipping</dt>
                <dd>{shipping === 0 ? "Free" : eur(shipping)}</dd>
              </div>
              <div className="flex justify-between border-t border-border/60 pt-3 text-base font-medium">
                <dt>Total</dt>
                <dd>{eur(total)}</dd>
              </div>
            </dl>
            <Link
              to="/checkout"
              className="mt-6 flex w-full items-center justify-center rounded-md bg-primary px-6 py-3.5 text-xs font-medium uppercase tracking-widest text-primary-foreground transition-colors hover:bg-accent"
            >
              Checkout
            </Link>
            <Link
              to="/shop"
              className="mt-3 block text-center text-xs uppercase tracking-widest text-muted-foreground hover:text-foreground"
            >
              Continue shopping
            </Link>
          </aside>
        </section>
      )}

      <SiteFooter />
    </div>
  );
}
