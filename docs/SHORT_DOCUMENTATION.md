# Incident Response Simulation and Infrastructure as Code Implementation
**Project Documentation**

---

## 1. Project Title
**Design and Deployment of a Containerized Microservices System with Terraform-Based Infrastructure Provisioning and Incident Response Simulation**

## 2. Objectives
The primary objective of this project is to integrate Infrastructure as Code (IaC) principles with Site Reliability Engineering (SRE) practices. Key goals include:
1. Implementing a microservices architecture using Go and TanStack Start.
2. Provisioning infrastructure via Terraform for reproducibility.
3. Orchestrating services using Docker Compose.
4. Integrating monitoring via Prometheus and Grafana.
5. Simulating and responding to a production-level incident.
6. Conducting a structured postmortem analysis.

## 3. System Architecture
The application follows a microservices-based architecture, consisting of independent services communicating over HTTP.

### Core Components
- **Frontend Layer:** React-based interface (TanStack Start) served via Nginx.
- **API Gateway:** Nginx reverse proxy for request routing.
- **Microservices:**
  - **Authentication Service:** Handles user login and token issuance.
  - **User Service:** Manages user profiles and internal chat functionality.
  - **Product Service:** Manages the product catalog and categories.
  - **Order Service:** Handles transactional operations and order creation.
- **Database Layer:** Isolated PostgreSQL containers for each service.
- **Monitoring Layer:**
  - **Prometheus:** Metrics collection.
  - **Grafana:** Data visualization and alerting.

## 4. Technology Stack
| Layer | Technology |
|-------|------------|
| Backend | Go (Golang) |
| Frontend | React (TanStack Start), JavaScript |
| Database | PostgreSQL 15 |
| Containerization | Docker |
| Orchestration | Docker Compose |
| Monitoring | Prometheus, Grafana |
| IaC | Terraform (Docker Provider) |

## 5. Infrastructure as Code (Assignment 5)
Infrastructure was provisioned using Terraform to ensure a declarative and reproducible environment. For local simulation, the Docker provider was used to create a "VM container".

### Terraform Configuration
- **Provider:** kreuzwerker/docker
- **Resources:**
    - **Network:** `market_place_infra`
    - **Container:** `production-server-vm` (Ubuntu 22.04)
    - **Ports:** 80, 3000, 9090, 2222

### Provisioning Workflow
```bash
terraform init
terraform plan
terraform apply -auto-approve
```

## 6. Containerized Deployment
The entire system is deployed as a multi-container application using a unified `docker-compose.yml` file.

### Service Overview
| Service | Container Name | Port |
|---------|----------------|------|
| API Gateway | api-gateway | 80 |
| Frontend | frontend | Internal |
| Auth Service | auth-service | 8080 |
| User Service | user-service | 8081 |
| Product Service | product-service | 8082 |
| Order Service | order-service | 8083 |
| Prometheus | prometheus | 9091 |
| Grafana | grafana | 3001 |

## 7. Monitoring and Observability
The system integrates Prometheus to collect metrics from all microservices and Grafana for visualization.
- **Prometheus Targets:** Configured to monitor health endpoints.
- **Grafana Dashboards:** Visualize CPU, memory, and error rates.

## 8. Incident Response Simulation (Assignment 4)
### Incident Scenario
A failure was introduced in the **Order Service** caused by an incorrect database connection string (`DATABASE_URL`).

### 1. Incident Detection
- **Alert:** Grafana triggered a "Critical: Order Service Failure" alert.
- **Observation:** Prometheus reported the `order-service` target as `DOWN`.

### 2. Incident Analysis
- **Log Review:** `docker logs order-service` showed connection timeouts to the database.
- **Root Cause:** A misconfiguration in `docker-compose.yml` pointed the service to a non-existent database host.

### 3. Incident Mitigation
- **Fix:** Corrected the `DATABASE_URL` environment variable.
- **Restoration:** Redeployed the service using `make build`.

### 4. Service Restoration
- **Verification:** Health checks passed, and metrics returned to normal levels in Grafana.

## 9. Postmortem Analysis
### Lessons Learned
- **What went well:** Monitoring alerted the team immediately, preventing extended silent downtime.
- **What went wrong:** Manual configuration changes in the compose file introduced human error.
- **Action Items:**
  - Implement automated configuration validation.
  - Introduce centralized logging (Loki/ELK).
  - Automate rollbacks for unhealthy service deployments.

## 10. Conclusion
This project successfully demonstrates the integration of modern DevOps practices, combining containerization, infrastructure automation, and SRE principles. The system is robust, observable, and reproducible, reflecting real-world production standards.
