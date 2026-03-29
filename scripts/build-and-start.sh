#!/usr/bin/env bash
set -euo pipefail

# Build the agent-proxy binary and start the service.
# Usage:
#   ./scripts/build-and-start.sh [--config path/to/config.yaml] [--port 1234]

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

echo "[agent-proxy] Building..."
make build

echo "[agent-proxy] Starting with args: $*"
exec ./agent-proxy "$@"
