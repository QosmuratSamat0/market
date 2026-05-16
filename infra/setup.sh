#!/bin/bash
set -euo pipefail

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
ORCHESTRATOR="${1:-compose}"

echo "Project root: $PROJECT_ROOT"
echo "Orchestrator: $ORCHESTRATOR"

echo "Configuring firewall ports..."

# 30080 - K8s NodePort
# 2377 - Swarm cluster management
for port in 30080 2377; do
    sudo iptables -C INPUT -p tcp --dport "$port" -j ACCEPT 2>/dev/null \
      || sudo iptables -A INPUT -p tcp --dport "$port" -j ACCEPT
done

# 7946, 4789 - Swarm overlay network
for port in 7946 4789; do
    sudo iptables -C INPUT -p tcp --dport "$port" -j ACCEPT 2>/dev/null \
      || sudo iptables -A INPUT -p tcp --dport "$port" -j ACCEPT

    sudo iptables -C INPUT -p udp --dport "$port" -j ACCEPT 2>/dev/null \
      || sudo iptables -A INPUT -p udp --dport "$port" -j ACCEPT
done

cd "$PROJECT_ROOT"

case "$ORCHESTRATOR" in
  compose)
    echo "Deploying with Docker Compose..."

    docker compose pull
    docker compose up -d --remove-orphans
    docker compose ps
    ;;

  swarm)
    echo "Deploying with Docker Swarm..."

    if ! docker info | grep -q "Swarm: active"; then
      SERVER_IP=$(hostname -I | awk '{print $1}')
      docker swarm init --advertise-addr "$SERVER_IP"
    else
      echo "Docker Swarm already active."
    fi

    if [ -f .env ]; then
      set -a
      source .env
      set +a
    fi

    cd "$PROJECT_ROOT/infra/swarm"
    docker stack deploy -c docker-compose.swarm.yml market
    docker stack services market
    ;;

  k8s)
    echo "Deploying with Kubernetes / K3s..."

    if ! command -v kubectl >/dev/null 2>&1; then
      curl -sfL https://get.k3s.io | sudo sh -
    fi

    sudo KUBECONFIG=/etc/rancher/k3s/k3s.yaml kubectl apply -f "$PROJECT_ROOT/infra/k8s/"
    sudo KUBECONFIG=/etc/rancher/k3s/k3s.yaml kubectl get all -n market
    ;;

  *)
    echo "Unknown orchestrator: $ORCHESTRATOR"
    echo "Usage: bash infra/setup.sh {compose|swarm|k8s}"
    exit 1
    ;;
esac

echo "Deployment finished: $ORCHESTRATOR"
