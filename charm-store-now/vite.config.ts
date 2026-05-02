// @lovable.dev/vite-tanstack-config already includes the following — do NOT add them manually
// or the app will break with duplicate plugins:
//   - tanstackStart, viteReact, tailwindcss, tsConfigPaths, cloudflare (build-only),
//     componentTagger (dev-only), VITE_* env injection, @ path alias, React/TanStack dedupe,
//     error logger plugins, and sandbox detection (port/host/strictPort).
// You can pass additional config via defineConfig({ vite: { ... } }) if needed.
import { defineConfig } from "@lovable.dev/vite-tanstack-config";

export default defineConfig({
  vite: {
    server: {
      port: 80,
      proxy: {
        "/auth": "http://localhost:8080",
        "/users": "http://localhost:8081",
        "/products": "http://localhost:8082",
        "/categories": "http://localhost:8082",
        "/orders": "http://localhost:8083",
      },
    },
    preview: {
      allowedHosts: ["samat.work", "www.samat.work", ".samat.work"],
      port: 80,
    },
  },
});
