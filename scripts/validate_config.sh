#!/bin/bash

# Configuration Validation Script
# Checks .env variables and connection strings before deployment

EXIT_CODE=0

echo "Starting configuration validation..."

# 1. Check if required variables are set
REQUIRED_VARS=(
    "JWT_SECRET"
    "DB_PASSWORD"
    "DB_USER"
    "AUTH_DB_NAME"
    "USER_DB_NAME"
    "PRODUCT_DB_NAME"
    "ORDER_DB_NAME"
    "PAYMENT_DB_NAME"
)

for VAR in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!VAR}" ]; then
        echo "Error: Required variable $VAR is missing or empty."
        EXIT_CODE=1
    fi
done

# 2. Validate Database Connection Format (if we were constructing it here)
# Since they are constructed in docker-compose.yml, we validate the components
if [[ ! "$DB_USER" =~ ^[a-zA-Z0-9_]+$ ]]; then
    echo "Error: DB_USER contains invalid characters."
    EXIT_CODE=1
fi

# 3. Check JWT_SECRET strength (min 16 chars)
if [ ${#JWT_SECRET} -lt 16 ]; then
    echo "Warning: JWT_SECRET is too short (less than 16 characters). This is not recommended for production."
    # We might not want to fail CI for this, but good to have a warning
fi

# 4. Validate Ports are numbers
PORT_VARS=("HTTP_PORT" "AUTH_PORT" "USER_PORT" "PRODUCT_PORT" "ORDER_PORT" "PAYMENT_PORT" "NOTIFICATION_PORT")
for PVAR in "${PORT_VARS[@]}"; do
    VAL="${!PVAR}"
    if [[ ! -z "$VAL" && ! "$VAL" =~ ^[0-9]+$ ]]; then
        echo "Error: $PVAR must be a number, got '$VAL'."
        EXIT_CODE=1
    fi
done

if [ $EXIT_CODE -eq 0 ]; then
    echo "Configuration validation successful."
else
    echo "Configuration validation failed."
fi

exit $EXIT_CODE
