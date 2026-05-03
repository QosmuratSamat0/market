#!/bin/bash
set -e

# Wait for postgres to be ready
# psql is available as postgres user

databases=("${AUTH_DB_NAME}" "${USER_DB_NAME}" "${PRODUCT_DB_NAME}" "${ORDER_DB_NAME}" "${PAYMENT_DB_NAME}")

for db in "${databases[@]}"; do
    if [ -n "$db" ]; then
        echo "Creating database: $db"
        psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_USER" <<-EOSQL
            CREATE DATABASE "$db";
EOSQL
    fi
done

echo "Running migrations..."

# Run migrations if the directory exists and contains sql files
run_migrations() {
    local service_db=$1
    local migration_dir=$2

    if [ -d "$migration_dir" ]; then
        for file in "$migration_dir"/*.sql; do
            if [ -e "$file" ]; then
                echo "Running migration $file on $service_db..."
                psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$service_db" -f "$file"
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
