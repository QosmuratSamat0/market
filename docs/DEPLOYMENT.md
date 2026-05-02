# Deployment Guide

Follow these steps to deploy the Marketplace system to a production-like environment.

## Prerequisites
- Docker & Docker Compose
- Terraform >= 1.0
- OCI credentials configured for Terraform
- Make

## Step 1: Infrastructure Provisioning
1. Navigate to the `infra` directory.
2. Initialize Terraform:
   ```bash
   terraform init
   ```
3. Review the plan:
   ```bash
   terraform plan
   ```
4. Apply the configuration:
   ```bash
   terraform apply -auto-approve
   ```
5. Note the `public_ip` output.

## Step 2: Service Configuration
1. Update any environment-specific variables in `.env` files (if applicable).
2. Ensure the `DATABASE_URL` in `docker-compose.yml` points to the correct database containers.

## Step 3: Deployment
1. Build and start all containers:
   ```bash
   make build
   ```
2. Apply database migrations:
   ```bash
   make migrate
   ```

## Step 4: Verification
1. Check container status:
   ```bash
   sudo docker ps
   ```
2. Verify monitoring targets in Prometheus (`http://<IP>:9090/targets`).
3. Log in to Grafana (`http://<IP>:3000`) and verify the dashboard metrics.

## Step 5: Troubleshooting
- **Service DOWN:** Check logs using `make logs`.
- **DB Connection Error:** Verify `DATABASE_URL` and ensure Postgres containers are healthy.
- **Port Conflict:** Ensure ports 80, 3000, and 9090 are not already in use on the host.
