# Marketplace Microservices & Infrastructure

This repository contains a containerized microservices-based application with automated infrastructure provisioning and integrated monitoring.

## Features
- **Microservices:** Auth, User, Product, and Order services written in Go.
- **Frontend:** Modern UI built with TanStack Start.
- **IaC:** Terraform scripts for AWS infrastructure.
- **Orchestration:** Docker Compose for local and production deployment.
- **Monitoring:** Prometheus and Grafana dashboards for real-time observability.
- **Resilience:** Incident response and postmortem documentation included.

## Project Structure
- `api-gateway/`: Nginx reverse proxy configuration.
- `auth-service/`, `user-service/`, `product-service/`, `order-service/`: Microservices source code.
- `charm-store-now/`: Frontend application.
- `infra/`: Terraform IaC files.
- `monitoring/`: Prometheus and Grafana configuration.
- `assignments/`: Detailed reports for Assignments 4 & 5.

## Quick Start

### 1. Provision Infrastructure
```bash
cd infra
terraform init
terraform apply
```

### 2. Deploy Services
```bash
# In the root directory
make build
make migrate
```

### 3. Access the System
- **Frontend/API:** [http://localhost](http://localhost)
- **Prometheus:** [http://localhost:9090](http://localhost:9090)
- **Grafana:** [http://localhost:3000](http://localhost:3000) (Admin: `admin/admin`)

## Monitoring
Metrics are collected automatically from all services. Visit the Grafana dashboard to view system health and performance.

## Documentation
- [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md): Step-by-step deployment guide.
- [docs/AGENTS.md](docs/AGENTS.md): agents.
- [docs/SHORT_DOCUMENTATION](docs/SHORT_DOCUMENTATION.md): Short documentation.
