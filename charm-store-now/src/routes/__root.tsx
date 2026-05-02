import { Outlet, Link, createRootRoute, HeadContent, Scripts } from "@tanstack/react-router";
import { QueryClientProvider } from "@tanstack/react-query";
import type { QueryClient } from "@tanstack/react-query";
import { Toaster } from "@/components/ui/sonner";
import { NotificationListener } from "@/components/notification-listener";
import { NotificationProvider } from "@/hooks/use-notifications";

import appCss from "../styles.css?url";

function NotFoundComponent() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4">
      <div className="max-w-md text-center">
        <h1 className="text-7xl font-bold text-foreground">404</h1>
        <h2 className="mt-4 text-xl font-semibold text-foreground">Page not found</h2>
        <p className="mt-2 text-sm text-muted-foreground">
          The page you're looking for doesn't exist or has been moved.
        </p>
        <div className="mt-6">
          <Link
            to="/"
            className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
          >
            Go home
          </Link>
        </div>
      </div>
    </div>
  );
}

export const Route = createRootRoute({
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { title: "market-space" },
      { name: "description", content: "Indigo Market is a frontend e-commerce application for browsing and purchasing products." },
      { name: "author", content: "Lovable" },
      { property: "og:title", content: "market-space" },
      { property: "og:description", content: "Indigo Market is a frontend e-commerce application for browsing and purchasing products." },
      { property: "og:type", content: "website" },
      { name: "twitter:card", content: "summary" },
      { name: "twitter:site", content: "@Lovable" },
      { name: "twitter:title", content: "market-space" },
      { name: "twitter:description", content: "Indigo Market is a frontend e-commerce application for browsing and purchasing products." },
      { property: "og:image", content: "https://pub-bb2e103a32db4e198524a2e9ed8f35b4.r2.dev/3f894aee-a3f1-4756-9309-18138c26dc6c/id-preview-a23841fd--6b7ca30f-ce8e-4ce6-9612-3b38db5878e2.lovable.app-1777121523391.png" },
      { name: "twitter:image", content: "https://pub-bb2e103a32db4e198524a2e9ed8f35b4.r2.dev/3f894aee-a3f1-4756-9309-18138c26dc6c/id-preview-a23841fd--6b7ca30f-ce8e-4ce6-9612-3b38db5878e2.lovable.app-1777121523391.png" },
    ],
    links: [
      {
        rel: "stylesheet",
        href: "https://fonts.googleapis.com/css2?family=Fraunces:opsz,wght@9..144,300;9..144,400;9..144,500;9..144,600&family=Inter:wght@300;400;500;600&display=swap",
      },
      {
        rel: "stylesheet",
        href: appCss,
      },
    ],
  }),
  shellComponent: RootShell,
  component: RootComponent,
  notFoundComponent: NotFoundComponent,
});

function RootShell({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body>
        {children}
        <Scripts />
      </body>
    </html>
  );
}

function RootComponent() {
  const { queryClient } = Route.useRouteContext() as { queryClient: QueryClient };

  return (
    <QueryClientProvider client={queryClient}>
      <NotificationProvider>
        <Outlet />
        <Toaster richColors closeButton />
        <NotificationListener />
      </NotificationProvider>
    </QueryClientProvider>
  );
}
