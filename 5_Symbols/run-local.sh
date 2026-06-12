#!/usr/bin/env bash
# Run the Go server locally from the project root.
# Usage: bash 5_Symbols/run-local.sh [port]
#
# Requires a .env file at the project root (copy from .env.example and fill in values).
# The server must be started from the project root so it can find all static assets.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

ENV_FILE="$PROJECT_ROOT/.env"
if [ ! -f "$ENV_FILE" ]; then
  if [ -f "$PROJECT_ROOT/.env.example" ]; then
    echo "⚠️  No .env found. Copying .env.example → .env"
    cp "$PROJECT_ROOT/.env.example" "$ENV_FILE"
    echo "✏️  Edit .env and fill in your SUPABASE_URL, SUPABASE_ANON_KEY, and AXIOM_DATASET, then re-run."
    exit 1
  else
    echo "❌ No .env or .env.example found at $PROJECT_ROOT"
    exit 1
  fi
fi

# Load .env — skip comments and blank lines, ignore export prefix
set -a
# shellcheck disable=SC1090
while IFS= read -r line || [ -n "$line" ]; do
  line="${line#export }"
  [[ "$line" =~ ^#.*$ || -z "$line" ]] && continue
  eval "$line" 2>/dev/null || true
done < "$ENV_FILE"
set +a

PORT="${1:-${PORT:-8080}}"
export PORT

echo "🔨 Building server..."
go build -o /tmp/claude-architect-server ./cmd/server

echo "🚀 Starting server on http://localhost:$PORT"
echo "   Press Ctrl+C to stop."
echo ""
exec /tmp/claude-architect-server
