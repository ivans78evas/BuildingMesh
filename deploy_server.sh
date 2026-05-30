#!/bin/bash

# BUILDINGMESH ENTERPRISE | DEPLOYMENT SCRIPT
# Targeted for: Debian 13 (x86_64)
# Optimized for: Intel N100 / Low-power Edge Servers

set -e

echo "--- BuildingMesh Environment Setup ---"
sudo apt update && sudo apt upgrade -y
sudo apt install -y curl git lsof

# 1. Docker Installation
if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh && sudo sh get-docker.sh
fi

# 2. Infrastructure Provisioning
echo "--- Provisioning Hybrid Stack (Go Core + Twenty CRM) ---"
mkdir -p storage/redis_data storage/postgres_data storage/uploads

# Ensure .env exists or create mock for local run
if [ ! -f .env ]; then
    echo "Creating default environment config..."
    cat <<EOF > .env
DB_PATH=construction.db
REDIS_URL=redis:6379
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
TWENTY_API_URL=http://twenty-server:3000/graphql
ALLOWED_ORIGINS=*
EOF
fi

# 3. Launch Cluster
echo "--- Launching Containers ---"
docker compose up -d --build

# 4. Verification
echo "--- Verification ---"
sleep 15 # Wait for Twenty CRM to initialize metadata engine

if curl -s http://localhost:8080/api/v1/health/live > /dev/null; then
    echo "=========================================================="
    echo "🚀 BUILDINGMESH IS LIVE!"
    echo "----------------------------------------------------------"
    echo "Engineering Console: http://localhost:8080"
    echo "Twenty CRM (Metadata): http://localhost:3000"
    echo "API Documentation:   http://localhost:8080/api/docs"
    echo "----------------------------------------------------------"
    echo "Server IP for Mobile App: \$(hostname -I | awk '{print \$1}')"
    echo "=========================================================="
else
    echo "❌ Deployment failed. Checking logs..."
    docker compose logs backend
fi
