import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { AuthShell, Field, PrimaryButton } from "@/components/auth-shell";
import { api } from "@/lib/api";
import { toast } from "sonner";

export const Route = createFileRoute("/sign-in")({
  head: () => ({ meta: [{ title: "Sign in — Maison" }] }),
  component: SignIn,
});

function SignIn() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.login({ email, password });
      toast.success("Signed in successfully");
      navigate({ to: "/shop" });
    } catch (err: any) {
      toast.error(err.message || "Failed to sign in");
    }
  };

  return (
    <AuthShell
      title="Welcome back."
      subtitle="Sign in to your Maison account."
      footer={
        <>
          New here?{" "}
          <Link to="/sign-up" className="text-foreground underline underline-offset-4">
            Create an account
          </Link>
        </>
      }
    >
      <form className="space-y-6" onSubmit={handleSubmit}>
        <Field 
          label="Email" 
          type="email" 
          autoComplete="email" 
          placeholder="your@email.com" 
          value={email}
          onChange={(e: any) => setEmail(e.target.value)}
        />
        <Field 
          label="Password" 
          type="password" 
          autoComplete="current-password" 
          placeholder="••••••••" 
          value={password}
          onChange={(e: any) => setPassword(e.target.value)}
        />
        <div className="flex items-center justify-between text-xs">
          <label className="flex items-center gap-2 text-muted-foreground">
            <input type="checkbox" className="h-3 w-3 accent-accent" /> Remember me
          </label>
          <a href="#" className="text-muted-foreground underline underline-offset-4 hover:text-foreground">
            Forgot password?
          </a>
        </div>
        <PrimaryButton>Sign in</PrimaryButton>
      </form>
      <div className="flex items-center gap-4 pt-2 text-xs uppercase tracking-widest text-muted-foreground">
        <span className="h-px flex-1 bg-border" /> or <span className="h-px flex-1 bg-border" />
      </div>
      <button className="w-full border border-border px-6 py-3.5 text-sm uppercase tracking-widest transition-colors hover:bg-secondary">
        Continue with Google
      </button>
    </AuthShell>
  );
}

