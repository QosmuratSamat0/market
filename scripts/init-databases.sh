#!/bin/bash
set -e

# Wait for postgres to be ready (internal postgres image healthcheck is usually sufficient, 
# but we add this for extra robustness)

# Function to create database if it doesn't exist
create_db_if_not_exists() {
    local db=$1
    echo "Checking database: $db"
    
    # Check if database exists using system catalogs
    local exists=$(psql -U "$POSTGRES_USER" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname = '$db'")
    
    if [ "$exists" != "1" ]; then
        echo "Creating database: $db"
        psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" -c "CREATE DATABASE \"$db\";"
    else
        echo "Database $db already exists, skipping creation."
    fi
}

databases=("${AUTH_DB_NAME}" "${USER_DB_NAME}" "${PRODUCT_DB_NAME}" "${ORDER_DB_NAME}" "${PAYMENT_DB_NAME}" "grafana")

for db in "${databases[@]}"; do
    if [ -n "$db" ]; then
        create_db_if_not_exists "$db"
    fi
done

echo "Running migrations..."

# Modified function for idempotent migrations
run_migrations() {
    local service_db=$1
    local migration_dir=$2

    if [ -d "$migration_dir" ]; then
        # Ensure alphanumeric sorting (001, 002...)
        for file in $(ls "$migration_dir"/*.sql | sort); do
            if [ -e "$file" ]; then
                echo "Applying migration $file to $service_db..."
                # Use -1 to run everything in a single transaction
                # Use ON_ERROR_STOP to fail fast
                psql -v ON_ERROR_STOP=1 -1 --username "$POSTGRES_USER" --dbname "$service_db" -f "$file"
            fi
        done
    fi
}

run_migrations "${AUTH_DB_NAME}" "/migrations/auth"
run_migrations "${USER_DB_NAME}" "/migrations/user"
run_migrations "${PRODUCT_DB_NAME}" "/migrations/product"
run_migrations "${ORDER_DB_NAME}" "/migrations/order"
run_migrations "${PAYMENT_DB_NAME}" "/migrations/payment"

echo "Initialization complete!"
