# Marketplace Microservices & SRE Infrastructure

This repository contains a containerized microservices-based application with automated multi-orchestration deployment, infrastructure provisioning, and integrated monitoring.

## 🚀 Features
- **Microservices:** Auth, User, Product, Order, Payment, and Notification services (Go).
- **Frontend:** Modern UI built with TanStack Start.
- **IaC:** Terraform scripts for Yandex Cloud infrastructure.
- **Multi-Orchestration:** Support for **Docker Compose**, **Docker Swarm**, and **Kubernetes (K3s)**.
- **Unified Deployment:** Idempotent `infra/setup.sh` script for environment provisioning and app deployment.
- **Monitoring:** Prometheus and Grafana dashboards for real-time observability.
- **CI/CD:** GitHub Actions pipeline with automated build, push, and multi-mode deploy.

## 📂 Project Structure
- `api-gateway/`: Nginx reverse proxy (NodePort 30080 in K8s).
- `infra/`: 
  - `terraform/`: Yandex Cloud IaC.
  - `k8s/`: Kubernetes manifests (00-12).
  - `swarm/`: Docker Swarm stack configuration.
  - `setup.sh`: Universal deployment script.
- `scripts/`: Database initialization and migration scripts.
- `monitoring/`: Prometheus and Grafana configuration.

## 🛠 Quick Start

### 1. Local Development (Windows/macOS/Linux)
> [!IMPORTANT]
> For local development, use **Docker Compose**. The `setup.sh` script is designed for **Linux Servers only**.

1. Create a `.env` file in the root directory.
2. Run the services:
```bash
# Works in PowerShell, Command Prompt, or Git Bash
docker compose up --build
```

### 2. Infrastructure Provisioning (Yandex Cloud)
```bash
cd infra
terraform init
terraform apply
```

### 3. Server Deployment (Universal Script)
> [!NOTE]
> This script is intended to run on the **Remote VPS (Ubuntu/Debian)**. It is automatically executed by GitHub Actions.

The `infra/setup.sh` script handles everything from firewall rules to orchestration setup.

```bash
# Deploy using Docker Compose
sudo bash infra/setup.sh compose

# Deploy using Docker Swarm
sudo bash infra/setup.sh swarm

# Deploy using Kubernetes (K3s)
sudo bash infra/setup.sh k8s
```

## Manual Run

Before starting any mode, create `.env` in the repository root. You can copy `.env.example` and change values if needed:

```bash
cp .env.example .env
```

Required variables:

```env
DB_USER=postgres
DB_PASSWORD=postgres
AUTH_DB_NAME=auth_db
USER_DB_NAME=user_db
PRODUCT_DB_NAME=product_db
ORDER_DB_NAME=order_db
PAYMENT_DB_NAME=payment_db
JWT_SECRET=your-super-secret-key
REGISTRY=ghcr.io
IMAGE_OWNER=your_github_name
RABBIT_USER=guest
RABBIT_PASS=guest
```

### Docker Compose

Requirements:
- Docker Engine with Docker Compose plugin.
- Free port `80` on the host.

Start:

```bash
docker compose up --build -d
```

Check:

```bash
docker compose ps
curl http://localhost/health
```

Stop:

```bash
docker compose down
```

Access:
- API Gateway: `http://localhost`
- Prometheus: `http://localhost/prometheus/`
- Grafana: `http://localhost/grafana/`
- Grafana login: `admin/admin`

### Docker Swarm

Requirements:
- Docker Engine.
- Swarm mode initialized.
- Free port `80` on the host.
- On Linux servers, monitoring exporters use host mounts: `/proc`, `/sys`, `/`, `/var/run`, `/var/lib/docker`, `/dev/disk`.

Initialize Swarm if it is not initialized:

```bash
docker swarm init
```

Deploy:

```bash
docker stack deploy -c infra/swarm/docker-compose.swarm.yml market
```

Check:

```bash
docker stack services market
docker service ls
curl http://localhost/health
```

View logs:

```bash
docker service logs market_api-gateway
```

Remove stack:

```bash
docker stack rm market
```

Access:
- API Gateway: `http://localhost`

### Kubernetes

Requirements:
- Kubernetes cluster or local K3s/minikube/Docker Desktop Kubernetes.
- `kubectl` configured for the target cluster.
- Images must be available in `ghcr.io/your_github_name`.

Apply manifests:

```bash
kubectl apply -f infra/k8s
```

Check:

```bash
kubectl get pods -n market
kubectl get svc -n market
kubectl get endpoints -n market
```

Wait for rollout:

```bash
kubectl rollout status deployment/postgres -n market
kubectl rollout status deployment/rabbitmq -n market
kubectl rollout status deployment/api-gateway -n market
```

If you changed the nginx ConfigMap, restart the gateway pod:

```bash
kubectl rollout restart deployment/api-gateway -n market
```

Access through NodePort:

```bash
curl http://localhost:30080/health
```

URLs:
- API Gateway: `http://localhost:30080`

If you want to use port `80` locally instead of `30080`, run:

```bash
kubectl port-forward svc/api-gateway 80:80 -n market
```

Then open:
- API Gateway: `http://localhost`

Debug commands:

```bash
kubectl logs deployment/api-gateway -n market
```

## 🤖 CI/CD Orchestration
You can switch the deployment mode by changing the `DEPLOY_MODE` variable in `.github/workflows/main.yml`:
- `DEPLOY_MODE: compose`
- `DEPLOY_MODE: swarm`
- `DEPLOY_MODE: k8s`

## 📊 Monitoring & Access
- **Docker Compose / Swarm API Gateway:** `http://<SERVER_IP>`

- **Kubernetes API Gateway:** `http://<SERVER_IP>:30080`


## 📄 Documentation
- [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md): Step-by-step deployment guide.
- [docs/SHORT_DOCUMENTATION.md](docs/SHORT_DOCUMENTATION.md): System architecture and components.
