import { createFileRoute, Link, notFound } from "@tanstack/react-router";
import { SiteHeader } from "@/components/site-header";
import { SiteFooter } from "@/components/site-footer";
import { ProductCard } from "@/components/product-card";
import { categories, products } from "@/lib/data";

export const Route = createFileRoute("/categories/$slug")({
  loader: ({ params }) => {
    const category = categories.find((c) => c.slug === params.slug);
    if (!category) throw notFound();
    return { category };
  },
  head: ({ loaderData }) => ({
    meta: [
      { title: `${loaderData?.category.name ?? "Category"} — Indigo Market` },
      { name: "description", content: loaderData?.category.description ?? "" },
    ],
  }),
  notFoundComponent: () => (
    <div className="min-h-screen">
      <SiteHeader />
      <div className="container-x py-32 text-center">
        <h1 className="font-display text-4xl">Category not found</h1>
        <Link to="/categories" className="mt-6 inline-block text-sm underline">Back to categories</Link>
      </div>
    </div>
  ),
  component: CategoryPage,
});

function CategoryPage() {
  const { category } = Route.useLoaderData();
  const items = products.filter((p) => p.categorySlug === category.slug);
  return (
    <div className="min-h-screen">
      <SiteHeader />
      <section className="container-x pt-16 pb-8">
        <Link to="/categories" className="text-xs uppercase tracking-widest text-muted-foreground hover:text-foreground">
          ← All categories
        </Link>
        <h1 className="mt-4 font-display text-5xl md:text-6xl">{category.name}</h1>
        <p className="mt-4 max-w-lg text-muted-foreground">{category.description}</p>
      </section>
      <section className="container-x py-12">
        <div className="grid grid-cols-2 gap-x-6 gap-y-12 md:grid-cols-4">
          {items.map((p) => <ProductCard key={p.id} {...p} />)}
        </div>
      </section>
      <SiteFooter />
    </div>
  );
}
