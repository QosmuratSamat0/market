import { Link } from "@tanstack/react-router";
import { eur } from "@/lib/data";

type Product = { id?: string; slug?: string; name: string; price: number; image?: string; tag?: string };

export function ProductCard({ id, slug, name, price, image, tag }: Product) {
  const finalSlug = slug || id || "unknown";
  const finalImage = image || "https://placehold.co/400x500?text=No+Image";

  return (
    <Link to="/products/$slug" params={{ slug: finalSlug }} className="group block">
      <div className="relative aspect-[4/5] overflow-hidden rounded-md bg-secondary">
        <img
          src={finalImage}
          alt={name}
          loading="lazy"
          className="h-full w-full object-cover transition-transform duration-700 ease-out group-hover:scale-105"
        />
        {tag && (
          <span className="absolute left-3 top-3 rounded-full bg-background/90 px-2.5 py-1 text-[10px] font-medium uppercase tracking-widest text-accent">
            {tag}
          </span>
        )}
      </div>
      <div className="mt-4 flex items-baseline justify-between">
        <h3 className="font-display text-base">{name}</h3>
        <span className="text-sm font-medium text-foreground">{eur(price)}</span>
      </div>
    </Link>
  );
}
