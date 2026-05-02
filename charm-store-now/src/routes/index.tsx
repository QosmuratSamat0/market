import { createFileRoute, Link } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { SiteHeader } from "@/components/site-header";
import { SiteFooter } from "@/components/site-footer";
import { ProductCard } from "@/components/product-card";
import { api, type Product, type Category } from "@/lib/api";
import { LayoutGrid, ShoppingBag, Sparkles } from "lucide-react";

export const Route = createFileRoute("/")({
  head: () => ({
    meta: [
      { title: "Marketplace — Minimal & Fast" },
      { name: "description", content: "A minimal, blue e-commerce marketplace." },
    ],
  }),
  component: Index,
});

function Index() {
  const { data: products = [], isLoading: loadingProducts } = useQuery<Product[]>({
    queryKey: ["products"],
    queryFn: api.getProducts,
  });

  const { data: categories = [], isLoading: loadingCategories } = useQuery<Category[]>({
    queryKey: ["categories"],
    queryFn: api.getCategories,
  });

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteHeader />
      
      <main className="flex-1 container-x py-6 space-y-10">
        {/* Compact Hero Banner */}
        <section className="bg-primary/10 rounded-2xl p-6 md:p-8 border border-primary/20 flex flex-col md:flex-row items-center justify-between gap-6 shadow-sm">
          <div>
            <div className="flex items-center gap-2 text-primary text-sm font-semibold mb-2">
              <Sparkles className="w-4 h-4" /> <span>Next Generation</span>
            </div>
            <h1 className="text-3xl md:text-4xl font-display text-foreground font-bold leading-tight">
              A minimalist approach <br />to digital commerce.
            </h1>
            <p className="text-muted-foreground mt-3 max-w-lg text-sm md:text-base">
              Discover our curated collection of premium products, designed with simplicity and elegance in mind.
            </p>
          </div>
          <div className="flex gap-3 w-full md:w-auto">
            <Link to="/shop" className="flex items-center justify-center gap-2 bg-primary text-primary-foreground px-6 py-2.5 rounded-xl text-sm font-semibold hover:bg-primary/90 transition shadow-sm w-full md:w-auto">
              <ShoppingBag className="w-4 h-4" /> Shop Now
            </Link>
          </div>
        </section>

        {/* Compact Categories Slider */}
        <section>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold flex items-center gap-2">
              <LayoutGrid className="w-5 h-5 text-primary" /> Categories
            </h2>
            <Link to="/categories" className="text-sm text-primary hover:underline font-medium">
              Browse All
            </Link>
          </div>
          
          {loadingCategories ? (
            <div className="flex gap-3">
              {[1, 2, 3, 4].map(i => (
                <div key={i} className="h-9 w-24 bg-secondary rounded-full animate-pulse" />
              ))}
            </div>
          ) : (
            <div className="flex gap-3 overflow-x-auto pb-2 scrollbar-hide">
              <Link 
                to="/shop" 
                className="whitespace-nowrap px-4 py-1.5 rounded-full bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors shadow-sm"
              >
                All Products
              </Link>
              {categories.map(c => (
                <Link 
                  key={c.id || c.slug} 
                  to="/categories/$slug" 
                  params={{ slug: c.slug }}
                  className="whitespace-nowrap px-4 py-1.5 rounded-full border border-border bg-card text-sm font-medium hover:border-primary hover:text-primary transition-colors shadow-sm"
                >
                  {c.name}
                </Link>
              ))}
            </div>
          )}
        </section>

        {/* Compact Featured Products Grid */}
        <section>
          <h2 className="text-lg font-semibold mb-4">Featured Additions</h2>
          
          {loadingProducts ? (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 md:gap-6">
              {[1, 2, 3, 4].map(i => (
                <div key={i} className="aspect-[4/5] bg-secondary rounded-xl animate-pulse" />
              ))}
            </div>
          ) : (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 md:gap-6">
              {products.slice(0, 8).map(p => (
                <ProductCard key={p.id} {...p} />
              ))}
            </div>
          )}
        </section>
      </main>

      <SiteFooter />
    </div>
  );
}
