# Auth Service

An authentication microservice responsible for generating and validating tokens (JWT & Refresh tokens) and managing user sessions within a distributed architecture. This service follows Clean Architecture principles.

## Features
- Issuance and validation of JWTs (Access Tokens).
- Management of Refresh Tokens securely stored in PostgreSQL.
- Integration with the User Service via RPC/HTTP (through a dedicated `UserClient`).
- Isolated configuration management (using `cleanenv` and `.env`).
- Role-Based Access Control (RBAC) middleware.

## How to Run

### Method 1: Using Docker Compose (Recommended)
This method spins up the application alongside a PostgreSQL database in isolated containers.

```bash
docker-compose up --build -d
```
The service will be available at `http://localhost:8080`.

### Method 2: Local Run (via Go)
If you want to run the project locally for development purposes:

1. Ensure you have a running PostgreSQL instance (you can run a standalone one via Docker).
2. Create the `auth_db` database and apply the necessary token tables.
3. Configure the `.env` file with the correct `DATABASE_URL` and other variables.
4. Run the service by first exporting the environment variables:
   
   ```bash
   # For Linux/macOS
   export $(grep -v '^#' .env | xargs) && go run cmd/main.go
   ```