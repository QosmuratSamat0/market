# AGENTS.md

Instructions for AI coding agents working in this monorepo.

## Scope
- Marketplace monorepo with Go microservices, an Nginx API gateway, and a TanStack Start frontend.
- Prefer small, service-local changes unless cross-service coordination is required.

## Fast Start Commands
- Start all services: `make up`
- Rebuild and start all services: `make build`
- Stop all services: `make down`
- Stop and remove volumes (destructive): `make down-clean`
- Apply SQL migrations manually: `make migrate`
- Tail logs: `make logs`

Notes:
- Root Makefile commands use `sudo docker`; this can fail without elevated privileges.
- Local service run pattern is `export $(grep -v '^#' .env | xargs) && go run cmd/main.go`.

## Project Map
- API gateway: [api-gateway](api-gateway) (`nginx.conf` path-based reverse proxy)
- Auth service: [auth-service](auth-service)
- User service: [user-service](user-service)
- Product service: [product-service](product-service)
- Order service: [order-service](order-service)
- Frontend app: [charm-store-now](charm-store-now)

## Backend Conventions (Go Services)
- Services follow layered structure under `internal/`: `domain` -> `usecase` -> `service` (optional) -> `repository` -> `handler/http`.
- App wiring happens in `internal/app/app.go` (dependency assembly, router setup, server start/shutdown).
- Route registration uses `*Module` types with `RegisterRoutes(chi.Router)` in `internal/handler/http/*_module.go`.
- Service-specific DB schema/data SQL files are in `migrations/` and mounted by docker-compose into Postgres init.

## Frontend Conventions
- Use scripts in [charm-store-now/package.json](charm-store-now/package.json): `dev`, `build`, `preview`, `lint`, `format`.
- Vite proxy routes and dev port are configured in [charm-store-now/vite.config.ts](charm-store-now/vite.config.ts).
- Do not manually edit generated router output: [charm-store-now/src/routeTree.gen.ts](charm-store-now/src/routeTree.gen.ts).

## Service Ports
- auth-service: `8080`
- user-service: `8081`
- product-service: `8082`
- order-service: `8083`
- frontend dev server: `80`

## Editing Guidance
- Keep changes scoped to one service when possible; update gateway/frontend proxies only when endpoint shapes or paths change.
- For backend API changes, verify call sites in both frontend `src/lib/api.ts` and any inter-service client under `internal/client`.
- Prefer adding migrations over editing existing applied migration files.

## Known Gaps / Pitfalls
- There are currently no Go test files (`*_test.go`) in services; add targeted tests when introducing business logic changes.
- Only [order-service/.env.example](order-service/.env.example) exists; other services rely on README and docker-compose defaults.

## Existing Docs (Link, Do Not Duplicate)
- Root commands: [Makefile](Makefile)
- Auth docs: [auth-service/README.md](auth-service/README.md)
- User docs: [user-service/README.md](user-service/README.md)
- Product docs: [product-service/README.md](product-service/README.md)
- Order docs: [order-service/README.md](order-service/README.md)
