import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { AuthShell, Field, PrimaryButton } from "@/components/auth-shell";
import { api } from "@/lib/api";
import { toast } from "sonner";

export const Route = createFileRoute("/sign-up")({
  head: () => ({ meta: [{ title: "Create account — Maison" }] }),
  component: SignUp,
});

function SignUp() {
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.register({ name, email, password });
      toast.success("Account created! Please sign in.");
      navigate({ to: "/sign-in" });
    } catch (err: any) {
      toast.error(err.message || "Failed to create account");
    }
  };

  return (
    <AuthShell
      title="Create your account."
      subtitle="Save favourites, track orders, hear about new editions first."
      footer={
        <>
          Already have one?{" "}
          <Link to="/sign-in" className="text-foreground underline underline-offset-4">
            Sign in
          </Link>
        </>
      }
    >
      <form className="space-y-6" onSubmit={handleSubmit}>
        <Field 
          label="Name" 
          autoComplete="name" 
          placeholder="Your name" 
          value={name}
          onChange={(e: any) => setName(e.target.value)}
        />
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
          autoComplete="new-password" 
          placeholder="At least 8 characters" 
          value={password}
          onChange={(e: any) => setPassword(e.target.value)}
        />
        <p className="text-xs leading-relaxed text-muted-foreground">
          By creating an account, you agree to our{" "}
          <a href="#" className="underline underline-offset-4">Terms</a> and{" "}
          <a href="#" className="underline underline-offset-4">Privacy Policy</a>.
        </p>
        <PrimaryButton>Create account</PrimaryButton>
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

