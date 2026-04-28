.PHONY: up down build migrate logs

# Start all services
up:
	sudo docker compose up -d

# Stop all services
down:
	sudo docker compose down

# Stop all services and remove volumes (DELETES DATA)
down-clean:
	sudo docker compose down -v

# Rebuild and start all services
build:
	sudo docker compose up -d --build

# Run database migrations for all services manually
migrate:
	@echo "Running Auth Service migrations..."
	sudo docker exec -i auth-postgres psql -U postgres -d auth_db < ./auth-service/migrations/001_create_refresh_tokens.sql
	@echo "Running User Service migrations..."
	sudo docker exec -i user-postgres psql -U postgres -d user_db < ./user-service/migrations/001_create_users.sql
	sudo docker exec -i user-postgres psql -U postgres -d user_db < ./user-service/migrations/002_seed_users.sql
	@echo "Running Product Service migrations..."
	sudo docker exec -i product-postgres psql -U postgres -d product_db < ./product-service/migrations/001_create_tables.sql
	sudo docker exec -i product-postgres psql -U postgres -d product_db < ./product-service/migrations/002_seed_catalog.sql
	@echo "Running Order Service migrations..."
	sudo docker exec -i order-postgres psql -U postgres -d order_db < ./order-service/migrations/001_create_orders.sql
	sudo docker exec -i order-postgres psql -U postgres -d order_db < ./order-service/migrations/002_seed_orders.sql
	@echo "All migrations applied successfully!"

# Show logs for all services
logs:
	sudo docker compose logs -f
