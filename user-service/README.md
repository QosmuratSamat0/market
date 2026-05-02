# User Service

A microservice responsible for user management (CRUD operations, profile management) within the marketplace platform. Built with Go, Chi router, and PostgreSQL following Clean Architecture principles.

## Features
- Full CRUD operations for user management.
- Authenticated user profile endpoints (`/me`).
- JWT-based authentication middleware (validates tokens issued by Auth Service).
- PostgreSQL storage with proper error handling.
- Role-Based Access Control (admin, user, manager).
- Clean Architecture (domain → usecase → repository → handler).

## API Endpoints

| Method | Endpoint       | Description                      | Auth Required |
|--------|----------------|----------------------------------|---------------|
| GET    | /users/me      | Get current user profile         | Yes           |
| PUT    | /users/me      | Update current user profile      | Yes           |
| GET    | /users/        | List all users                   | Yes           |
| POST   | /users/        | Create a new user                | Yes           |
| GET    | /users/{id}    | Get user by ID                   | Yes           |
| PUT    | /users/{id}    | Update user by ID                | Yes           |
| DELETE | /users/{id}    | Delete user by ID                | Yes           |

## How to Run

### Method 1: Using Docker Compose (Recommended)
Spins up the service alongside a PostgreSQL database. Migrations are applied automatically.

```bash
docker-compose up --build -d
```
The service will be available at `http://localhost:8081`.

### Method 2: Local Run (via Go)

1. Ensure you have a running PostgreSQL instance.
2. Create the `user_db` database and run `migrations/001_create_users.sql`.
3. Configure the `.env` file with the correct `DATABASE_URL` and `JWT_SECRET`.
4. Run the service:

   ```bash
   export $(grep -v '^#' .env | xargs) && go run cmd/main.go
   ```

## Project Structure

```
user-service/
├── cmd/
│   └── main.go                  # Entry point
├── internal/
│   ├── app/
│   │   └── app.go               # Application initialization & run
│   ├── config/
│   │   └── config.go            # Configuration loader (cleanenv)
│   ├── domain/
│   │   └── user/
│   │       └── model.go         # User domain model & roles
│   ├── handler/
│   │   └── http/
│   │       ├── middleware/
│   │       │   └── auth.go      # JWT auth middleware
│   │       ├── user.go          # User HTTP handlers
│   │       └── user_module.go   # Route registration
│   ├── lib/
│   │   ├── api/response/
│   │   │   └── response.go     # API response helpers
│   │   ├── errs/
│   │   │   └── errors.go       # Custom errors
│   │   ├── passwordUtils/
│   │   │   └── password.go     # Bcrypt password utilities
│   │   └── tokens/
│   │       └── jwt.go          # JWT parsing
│   ├── repository/
│   │   └── user/postgres/
│   │       └── repo.go         # PostgreSQL repository
│   └── usecase/
│       └── user/
│           ├── interfaces.go   # Repository interface
│           └── usecase.go      # Business logic
├── migrations/
│   └── 001_create_users.sql
├── Dockerfile
├── docker-compose.yml
├── .env
└── .gitignore
```
