import { Link } from "@tanstack/react-router";
import type { ReactNode } from "react";
import hero from "@/assets/hero.jpg";

export function AuthShell({
  title,
  subtitle,
  children,
  footer,
}: {
  title: string;
  subtitle: string;
  children: ReactNode;
  footer: ReactNode;
}) {
  return (
    <div className="grid min-h-screen md:grid-cols-2">
      <div className="relative hidden md:block">
        <img src={hero} alt="" className="h-full w-full object-cover" />
        <div className="absolute inset-0 bg-gradient-to-b from-transparent to-background/30" />
        <Link to="/" className="absolute left-8 top-8 font-display text-xl font-medium text-primary-foreground mix-blend-difference">
          indigo.
        </Link>
        <div className="absolute bottom-10 left-10 right-10 text-primary-foreground mix-blend-difference">
          <p className="font-display text-2xl leading-snug">
            "A simpler way to shop the things you love."
          </p>
        </div>
      </div>
      <div className="flex items-center justify-center px-6 py-16">
        <div className="w-full max-w-sm">
          <Link to="/" className="font-display text-xl font-medium md:hidden">indigo.</Link>
          <h1 className="mt-6 font-display text-3xl md:mt-0">{title}</h1>
          <p className="mt-2 text-sm text-muted-foreground">{subtitle}</p>
          <div className="mt-10 space-y-5">{children}</div>
          <div className="mt-8 text-sm text-muted-foreground">{footer}</div>
        </div>
      </div>
    </div>
  );
}

export function Field({
  label,
  type = "text",
  placeholder,
  autoComplete,
  ...rest
}: {
  label: string;
  type?: string;
  placeholder?: string;
  autoComplete?: string;
} & React.InputHTMLAttributes<HTMLInputElement>) {
  return (
    <label className="block">
      <span className="text-xs uppercase tracking-widest text-muted-foreground">{label}</span>
      <input
        type={type}
        placeholder={placeholder}
        autoComplete={autoComplete}
        {...rest}
        className="mt-2 block w-full rounded-md border border-border bg-background px-3 py-2.5 text-sm outline-none transition-colors placeholder:text-muted-foreground/60 focus:border-accent focus:ring-2 focus:ring-accent/20"
      />
    </label>
  );
}

export function PrimaryButton({ children }: { children: ReactNode }) {
  return (
    <button className="group flex w-full items-center justify-center gap-2 rounded-md bg-primary px-6 py-3.5 text-xs font-medium uppercase tracking-widest text-primary-foreground transition-all hover:bg-accent">
      {children}
    </button>
  );
}
