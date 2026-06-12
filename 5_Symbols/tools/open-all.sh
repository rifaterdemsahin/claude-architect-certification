#!/usr/bin/env bash
# Open all project pages — builds the Go server, starts it, and opens the sitemap.
# Usage: bash 5_Symbols/tools/open-all.sh [port]
#
# The sitemap at http://localhost:<port>/5_Symbols/tools/sitemap.html
# lists every page in the project for one-click navigation.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
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

# Load .env
set -a
while IFS= read -r line || [ -n "$line" ]; do
  line="${line#export }"
  [[ "$line" =~ ^#.*$ || -z "$line" ]] && continue
  eval "$line" 2>/dev/null || true
done < "$ENV_FILE"
set +a

PORT="${1:-${PORT:-8080}}"
export PORT

SITEMAP_URL="http://localhost:$PORT/5_Symbols/tools/sitemap.html"

echo "🔨 Building server..."
go build -o /tmp/claude-architect-server ./cmd/server

echo "🚀 Starting server on http://localhost:$PORT"
echo "   Sitemap: $SITEMAP_URL"
echo "   Press Ctrl+C to stop."
echo ""

# Start server in background
/tmp/claude-architect-server &
SERVER_PID=$!

# Cleanup on exit
cleanup() {
  echo ""
  echo "⏹ Stopping server..."
  kill "$SERVER_PID" 2>/dev/null || true
  wait "$SERVER_PID" 2>/dev/null || true
  echo "✅ Server stopped."
}
trap cleanup EXIT INT TERM

# Wait for server to be ready
echo "⏳ Waiting for server..."
for i in $(seq 1 30); do
  if curl -s -o /dev/null "http://localhost:$PORT/" 2>/dev/null; then
    echo "✅ Server is ready."
    break
  fi
  if [ "$i" -eq 30 ]; then
    echo "❌ Server failed to start within 30 seconds."
    exit 1
  fi
  sleep 1
done

echo ""
echo "📂 Opening sitemap..."
open "$SITEMAP_URL"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  🗺  Sitemap loaded in your browser."
echo "     Browse all project pages from one place."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Keep server running in foreground
wait "$SERVER_PID"