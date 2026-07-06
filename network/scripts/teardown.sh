#!/bin/bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NETWORK_DIR="$(dirname "$SCRIPT_DIR")"
cd "$NETWORK_DIR"

docker-compose down --volumes --remove-orphans
docker rm -f $(docker ps -aq --filter "name=dev-peer") 2>/dev/null || true
docker rmi -f $(docker images -q "dev-peer*") 2>/dev/null || true
rm -rf crypto-config channel-artifacts

echo "网络已清理"
