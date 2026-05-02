import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { SiteHeader } from "@/components/site-header";
import { SiteFooter } from "@/components/site-footer";
import { Field } from "@/components/auth-shell";
import { api } from "@/lib/api";
import { Package, MapPin, CreditCard, LogOut } from "lucide-react";
import { toast } from "sonner";

export const Route = createFileRoute("/account")({
  head: () => ({ meta: [{ title: "Account — Indigo Market" }] }),
  component: Account,
});

function Account() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");

  const {
    data: profile,
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ["user", "me"],
    queryFn: api.getMe,
    retry: false,
  });

  useEffect(() => {
    if (profile) {
      setName(profile.name || "");
      setEmail(profile.email || "");
    }
  }, [profile]);

  useEffect(() => {
    if (isError) {
      const message = (error as Error)?.message || "Please sign in first";
      if (message.toLowerCase().includes("sign in") || message.toLowerCase().includes("unauthorized")) {
        toast.error("Please sign in to view your account");
        navigate({ to: "/sign-in" });
      }
    }
  }, [isError, error, navigate]);

  const updateProfileMutation = useMutation({
    mutationFn: api.updateMe,
    onSuccess: (updatedProfile) => {
      queryClient.setQueryData(["user", "me"], updatedProfile);
      toast.success("Profile updated successfully");
    },
    onError: (err: unknown) => {
      toast.error(err instanceof Error ? err.message : "Failed to update profile");
    },
  });

  const handleSaveProfile = (e: React.FormEvent) => {
    e.preventDefault();
    updateProfileMutation.mutate({
      name: name.trim(),
      email: email.trim(),
    });
  };

  const handleSignOut = async () => {
    try {
      await api.logout();
    } catch {
      // Ignore backend logout failures, local token is already removed.
    }
    toast.success("Signed out successfully");
    queryClient.removeQueries({ queryKey: ["user", "me"] });
    navigate({ to: "/sign-in" });
  };

  return (
    <div className="min-h-screen">
      <SiteHeader />
      <section className="container-x pt-16 pb-10">
        <span className="text-xs uppercase tracking-[0.25em] text-muted-foreground">Profile</span>
        <h1 className="mt-4 font-display text-5xl md:text-6xl">Your account</h1>
        <p className="mt-4 max-w-lg text-muted-foreground">Manage your profile, addresses and payment methods.</p>
      </section>

      <section className="container-x grid gap-10 pb-24 md:grid-cols-[240px_1fr]">
        <aside className="space-y-1">
          <NavItem icon={<Package className="h-4 w-4" />} to="/orders" label="Orders" />
          <NavItem icon={<MapPin className="h-4 w-4" />} to="/account" label="Profile" active />
          <NavItem icon={<CreditCard className="h-4 w-4" />} to="/payments" label="Payment" />
          <button
            type="button"
            onClick={handleSignOut}
            className="mt-4 flex w-full items-center gap-3 rounded-md px-3 py-2 text-left text-sm text-muted-foreground hover:bg-secondary hover:text-foreground"
          >
            <LogOut className="h-4 w-4" /> Sign out
          </button>
        </aside>

        <div className="space-y-10">
          <form onSubmit={handleSaveProfile} className="rounded-lg border border-border/60 bg-card p-6">
            <h2 className="font-display text-2xl">Profile</h2>
            {isLoading && (
              <p className="mt-4 text-sm text-muted-foreground">Loading profile...</p>
            )}
            {!isLoading && !isError && profile && (
              <p className="mt-4 text-sm text-muted-foreground">
                Signed in as <span className="font-medium text-foreground">{profile.email}</span>
              </p>
            )}
            <div className="mt-6 grid gap-6 md:grid-cols-2">
              <Field
                label="Name"
                placeholder="Alex Morgan"
                value={name}
                onChange={(e) => setName(e.target.value)}
                disabled={isLoading || updateProfileMutation.isPending}
              />
              <Field
                label="Email"
                type="email"
                placeholder="alex@email.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                disabled={isLoading || updateProfileMutation.isPending}
              />
              <Field label="Phone" placeholder="+351 …" />
            </div>
            <button
              disabled={isLoading || updateProfileMutation.isPending || !name.trim() || !email.trim()}
              className="mt-8 rounded-md bg-primary px-6 py-3 text-xs font-medium uppercase tracking-widest text-primary-foreground hover:bg-accent disabled:cursor-not-allowed disabled:opacity-60"
            >
              {updateProfileMutation.isPending ? "Saving..." : "Save changes"}
            </button>
          </form>

          <form onSubmit={(e) => e.preventDefault()} className="rounded-lg border border-border/60 bg-card p-6">
            <h2 className="font-display text-2xl">Default address</h2>
            <div className="mt-6 space-y-6">
              <Field label="Street" placeholder="Rua Exemplo, 12" />
              <div className="grid gap-6 md:grid-cols-3">
                <Field label="City" placeholder="Lisbon" />
                <Field label="Postal code" placeholder="1000-001" />
                <Field label="Country" placeholder="Portugal" />
              </div>
            </div>
            <button className="mt-8 rounded-md border border-border px-6 py-3 text-xs font-medium uppercase tracking-widest hover:bg-secondary">
              Update address
            </button>
          </form>
        </div>
      </section>
      <SiteFooter />
    </div>
  );
}

function NavItem({ icon, to, label, active }: { icon: React.ReactNode; to: "/orders" | "/account" | "/payments"; label: string; active?: boolean }) {
  return (
    <Link to={to} className={`flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors ${active ? "bg-secondary text-foreground" : "text-muted-foreground hover:bg-secondary hover:text-foreground"}`}>
      {icon} {label}
    </Link>
  );
}
