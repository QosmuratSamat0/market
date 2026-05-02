import { createFileRoute, Link } from "@tanstack/react-router";
import { SiteHeader } from "@/components/site-header";
import { SiteFooter } from "@/components/site-footer";
import { categories, products } from "@/lib/data";

export const Route = createFileRoute("/categories")({
  head: () => ({
    meta: [
      { title: "Categories — Indigo Market" },
      { name: "description", content: "Browse all product categories." },
    ],
  }),
  component: Categories,
});

function Categories() {
  return (
    <div className="min-h-screen">
      <SiteHeader />
      <section className="container-x pt-16 pb-10">
        <span className="text-xs uppercase tracking-[0.25em] text-muted-foreground">Browse</span>
        <h1 className="mt-4 font-display text-5xl md:text-6xl">Categories</h1>
      </section>
      <section className="container-x grid gap-6 pb-20 md:grid-cols-2">
        {categories.map((c) => {
          const count = products.filter((p) => p.categorySlug === c.slug).length;
          return (
            <Link
              key={c.slug}
              to="/categories/$slug"
              params={{ slug: c.slug }}
              className="group relative grid aspect-[16/9] overflow-hidden rounded-lg bg-secondary"
            >
              <img src={c.image} alt={c.name} className="absolute inset-0 h-full w-full object-cover opacity-80 transition-transform duration-700 group-hover:scale-105" />
              <div className="absolute inset-0 bg-gradient-to-t from-foreground/80 to-foreground/10" />
              <div className="relative z-10 mt-auto p-8 text-background">
                <div className="text-xs uppercase tracking-widest opacity-80">{count} products</div>
                <div className="mt-1 font-display text-3xl">{c.name}</div>
                <div className="mt-1 max-w-md text-sm opacity-80">{c.description}</div>
              </div>
            </Link>
          );
        })}
      </section>
      <SiteFooter />
    </div>
  );
}
