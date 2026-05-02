export function SiteFooter() {
  return (
    <footer className="mt-24 border-t border-border/60">
      <div className="container-x grid gap-12 py-16 md:grid-cols-4">
        <div>
          <div className="font-display text-xl font-medium">indigo<span className="text-accent">.</span></div>
          <p className="mt-3 max-w-xs text-sm text-muted-foreground">
            A modern marketplace for everyday essentials. Shipped fast, returned easily.
          </p>
        </div>
        <FooterCol title="Shop" items={["All products", "Audio", "Wearables", "Smart Home"]} />
        <FooterCol title="Account" items={["Sign in", "Create account", "Orders", "Profile"]} />
        <FooterCol title="Help" items={["Shipping", "Returns", "Contact", "FAQ"]} />
      </div>
      <div className="container-x flex flex-col gap-2 border-t border-border/60 py-6 text-xs text-muted-foreground md:flex-row md:justify-between">
        <span>© {new Date().getFullYear()} Indigo Market</span>
        <span>Built with care</span>
      </div>
    </footer>
  );
}

function FooterCol({ title, items }: { title: string; items: string[] }) {
  return (
    <div>
      <div className="text-xs font-medium uppercase tracking-widest text-foreground">{title}</div>
      <ul className="mt-4 space-y-2 text-sm text-muted-foreground">
        {items.map((i) => (
          <li key={i}><a href="#" className="transition-colors hover:text-foreground">{i}</a></li>
        ))}
      </ul>
    </div>
  );
}
