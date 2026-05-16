#!/bin/bash
set -e

create_db_if_not_exists() {
    local db=$1

    if [ -z "$db" ]; then
        return
    fi

    echo "Checking database: $db"

    local exists
    exists=$(psql -U "$POSTGRES_USER" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname = '$db'")

    if [ "$exists" != "1" ]; then
        echo "Creating database: $db"
        psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" -c "CREATE DATABASE \"$db\";"
    else
        echo "Database $db already exists, skipping creation."
    fi
}

databases=(
    "$AUTH_DB_NAME"
    "$USER_DB_NAME"
    "$PRODUCT_DB_NAME"
    "$ORDER_DB_NAME"
    "$PAYMENT_DB_NAME"
    "grafana"
)

for db in "${databases[@]}"; do
    create_db_if_not_exists "$db"
done

echo "Running migrations..."

run_migrations() {
    local service_db=$1
    local migration_dir=$2

    if [ -z "$service_db" ]; then
        echo "Database name is empty, skipping $migration_dir"
        return
    fi

    if [ ! -d "$migration_dir" ]; then
        echo "Migration directory not found: $migration_dir, skipping."
        return
    fi

    for file in "$migration_dir"/*.sql; do
        [ -e "$file" ] || continue

        echo "Applying migration $file to $service_db..."
        psql -v ON_ERROR_STOP=1 -1 --username "$POSTGRES_USER" --dbname "$service_db" -f "$file"
    done
}

run_migrations "$AUTH_DB_NAME" "/migrations/auth"
run_migrations "$USER_DB_NAME" "/migrations/user"
run_migrations "$PRODUCT_DB_NAME" "/migrations/product"
run_migrations "$ORDER_DB_NAME" "/migrations/order"
run_migrations "$PAYMENT_DB_NAME" "/migrations/payment"

echo "Initialization complete!"
