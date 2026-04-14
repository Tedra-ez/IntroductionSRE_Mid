#!/usr/bin/env sh
# Build images required by docker-stack.yml (Swarm does not run docker build).
set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
echo "Building clothes-store-backend:latest ..."
docker build -t clothes-store-backend:latest -f backend/Dockerfile .
echo "Building clothes-store-frontend:latest ..."
docker build -t clothes-store-frontend:latest -f frontend/Dockerfile .
echo "Done. Run: docker stack deploy -c docker-stack.yml clothes-store"
