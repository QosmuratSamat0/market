#!/bin/bash
# Manual deployment script - run this on the server to deploy immediately
# Usage: ./scripts/manual_deploy.sh
set -e

echo "=== Manual Deploy ==="

# 1. Pull latest code
echo "📥 Pulling latest code..."
git pull origin main

# 2. Login to GHCR (if token is available)
if [ -n "$GHCR_TOKEN" ] && [ -n "$GHCR_USER" ]; then
  echo "🔑 Logging into GHCR..."
  echo "$GHCR_TOKEN" | sudo docker login ghcr.io -u "$GHCR_USER" --password-stdin
else
  echo "⚠️ GHCR_TOKEN/GHCR_USER not set. Trying anonymous pull (packages must be public)."
  echo "   Set with: export GHCR_TOKEN=ghp_xxx GHCR_USER=QosmuratSamat0"
fi

# 3. Pull new images
echo "📦 Pulling new images..."
sudo docker compose pull

# 4. Recreate containers
echo "🔄 Recreating containers..."
sudo docker compose up -d --force-recreate --remove-orphans

# 5. Wait and verify
echo "⏳ Waiting for containers to start..."
sleep 15

echo ""
echo "=== Container Status ==="
sudo docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Image}}'

echo ""
echo "=== Health Checks ==="
for svc in auth-service user-service product-service order-service; do
  port=$(sudo docker inspect --format='{{range $p, $conf := .Config.ExposedPorts}}{{$p}}{{end}}' "$svc" 2>/dev/null | grep -oP '\d+' | head -1)
  if [ -n "$port" ]; then
    health=$(sudo docker exec "$svc" curl -sf "http://localhost:$port/health" 2>/dev/null && echo "✅ OK" || echo "❌ FAIL")
    echo "  $svc (port $port): $health"
  fi
done

echo ""
echo "🧹 Cleaning up old images..."
sudo docker image prune -f

echo ""
echo "✅ Deploy complete!"
