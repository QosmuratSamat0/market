import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { SiteHeader } from "@/components/site-header";
import { SiteFooter } from "@/components/site-footer";
import { ProductCard } from "@/components/product-card";
import { api, type Category, type Product } from "@/lib/api";

export const Route = createFileRoute("/shop")({
  head: () => ({
    meta: [
      { title: "Shop — Indigo Market" },
      { name: "description", content: "Browse all products from Indigo Market." },
    ],
  }),
  component: Shop,
});

function Shop() {
  const [active, setActive] = useState<string>("all");

  const { data: products = [], isLoading: loadingProducts } = useQuery<Product[]>({
    queryKey: ["products"],
    queryFn: api.getProducts,
  });

  const { data: categories = [] } = useQuery<Category[]>({
    queryKey: ["categories"],
    queryFn: api.getCategories,
  });

  const filtered =
    active === "all"
      ? products
      : products.filter((p) => p.categorySlug === active || p.categoryId === active);

  return (
    <div className="min-h-screen">
      <SiteHeader />
      <section className="container-x pt-16 pb-8">
        <span className="text-xs uppercase tracking-[0.25em] text-muted-foreground">Catalog</span>
        <h1 className="mt-4 font-display text-5xl md:text-6xl">All products</h1>
        <p className="mt-4 max-w-lg text-muted-foreground">
          {loadingProducts
            ? "Loading..."
            : `${filtered.length} item${filtered.length === 1 ? "" : "s"} available.`}
        </p>
      </section>
      <section className="container-x sticky top-16 z-30 border-y border-border/60 bg-background/85 backdrop-blur">
        <div className="flex gap-6 overflow-x-auto py-4 text-xs uppercase tracking-widest">
          <button
            onClick={() => setActive("all")}
            className={
              active === "all"
                ? "text-foreground underline underline-offset-8"
                : "text-muted-foreground hover:text-foreground"
            }
          >
            All
          </button>
          {categories.map((c) => (
            <button
              key={c.id || c.slug}
              onClick={() => setActive(c.id || c.slug)}
              className={
                active === (c.id || c.slug)
                  ? "text-foreground underline underline-offset-8"
                  : "text-muted-foreground hover:text-foreground"
              }
            >
              {c.name}
            </button>
          ))}
        </div>
      </section>
      <section className="container-x py-12">
        <div className="grid grid-cols-2 gap-x-6 gap-y-12 md:grid-cols-4">
          {filtered.map((p) => (
            <ProductCard key={p.id} {...p} />
          ))}
        </div>
      </section>
      <SiteFooter />
    </div>
  );
}
