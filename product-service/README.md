# Product Service

A microservice responsible for product and category management within the marketplace platform. Built with Go, Chi router, and PostgreSQL following Clean Architecture principles.

## Features
- Full CRUD for products (create, read, update, delete).
- Category management (create, list, delete).
- Seller ownership — only the product owner can update/delete their products.
- Admin override — admins can delete any product.
- Public endpoints for browsing products and categories (no auth required).
- Protected endpoints for creating/updating/deleting (JWT auth required).
- Shared JWT validation with auth-service (same `JWT_SECRET`).

## API Endpoints

### Public (no auth)
| Method | Endpoint                              | Description                    |
|--------|---------------------------------------|--------------------------------|
| GET    | /products/                            | List all products              |
| GET    | /products/{id}                        | Get product by ID              |
| GET    | /categories                           | List all categories            |
| GET    | /categories/{categoryID}/products     | Get products by category       |

### Protected (JWT required)
| Method | Endpoint          | Description                          |
|--------|-------------------|--------------------------------------|
| GET    | /products/my      | Get current seller's products        |
| POST   | /products/        | Create a new product                 |
| PUT    | /products/{id}    | Update product (owner only)          |
| DELETE | /products/{id}    | Delete product (owner or admin)      |
| POST   | /categories/      | Create a new category                |
| DELETE | /categories/{id}  | Delete a category                    |

## How to Run

### Method 1: Using Docker Compose (Recommended)
```bash
docker-compose up --build -d
```
The service will be available at `http://localhost:8082`.

### Method 2: Local Run (via Go)
1. Ensure PostgreSQL is running.
2. Create `product_db` and run `migrations/001_create_tables.sql`.
3. Configure `.env` with the correct `DATABASE_URL` and `JWT_SECRET`.
4. Run:
   ```bash
   export $(grep -v '^#' .env | xargs) && go run cmd/main.go
   ```

## Microservice Ports
| Service         | Port  |
|-----------------|-------|
| auth-service    | 8080  |
| user-service    | 8081  |
| product-service | 8082  |

All services share the same `JWT_SECRET` for token validation.
