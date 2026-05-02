import { createFileRoute, Link } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { SiteHeader } from "@/components/site-header";
import { SiteFooter } from "@/components/site-footer";
import { eur } from "@/lib/data";
import { api } from "@/lib/api";
import { addToCart } from "@/lib/cart";
import { ShoppingBag, Truck, RefreshCw, ShieldCheck } from "lucide-react";
import { toast } from "sonner";

export const Route = createFileRoute("/products/$slug")({
  head: () => ({ meta: [{ title: "Product — Indigo Market" }] }),
  component: ProductPage,
});

function ProductPage() {
  const { slug } = Route.useParams();
  const {
    data: product,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["product", slug],
    queryFn: () => api.getProduct(slug),
  });

  const handleAddToCart = () => {
    if (!product) return;
    addToCart(product.id);
    toast.success("Added to cart");
  };

  if (isLoading) {
    return (
      <div className="min-h-screen">
        <SiteHeader />
        <div className="container-x py-32 text-center text-muted-foreground">
          Loading product...
        </div>
      </div>
    );
  }

  if (isError || !product) {
    return (
      <div className="min-h-screen">
        <SiteHeader />
        <div className="container-x py-32 text-center">
          <h1 className="font-display text-4xl">Product not found</h1>
          <Link to="/shop" className="mt-6 inline-block text-sm underline">
            Back to shop
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <SiteHeader />
      <section className="container-x grid gap-12 pt-10 pb-20 md:grid-cols-2">
        <div className="aspect-square overflow-hidden rounded-lg bg-secondary">
          <img src={product.image} alt={product.name} className="h-full w-full object-cover" />
        </div>
        <div className="md:pt-8">
          <div className="text-xs uppercase tracking-widest text-muted-foreground">
            {product.categoryId || "Product"}
          </div>
          <h1 className="mt-3 font-display text-4xl md:text-5xl">{product.name}</h1>
          <div className="mt-4 text-2xl font-medium">{eur(product.price)}</div>
          <p className="mt-6 max-w-md text-muted-foreground">{product.description}</p>

          <div className="mt-10 flex items-center gap-3">
            <button
              onClick={handleAddToCart}
              className="group inline-flex flex-1 items-center justify-center gap-2 rounded-md bg-primary px-6 py-3.5 text-xs font-medium uppercase tracking-widest text-primary-foreground transition-colors hover:bg-accent"
            >
              <ShoppingBag className="h-4 w-4" /> Add to cart
            </button>
            <button className="rounded-md border border-border px-6 py-3.5 text-xs font-medium uppercase tracking-widest hover:bg-secondary">
              Save
            </button>
          </div>

          <ul className="mt-10 grid gap-3 border-t border-border/60 pt-8 text-sm text-muted-foreground">
            <li className="flex items-center gap-3">
              <Truck className="h-4 w-4 text-accent" /> Free shipping over €100
            </li>
            <li className="flex items-center gap-3">
              <RefreshCw className="h-4 w-4 text-accent" /> 30-day returns
            </li>
            <li className="flex items-center gap-3">
              <ShieldCheck className="h-4 w-4 text-accent" /> 2-year warranty
            </li>
          </ul>
        </div>
      </section>
      <SiteFooter />
    </div>
  );
}
