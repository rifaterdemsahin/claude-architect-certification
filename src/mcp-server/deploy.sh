#!/bin/sh
# Helper script to deploy the MCP server to Fly.io
set -e

echo "🚀 Starting Enterprise MCP Data Bridge Deployment to Fly.io..."

# Ensure we are in the correct directory
cd "$(dirname "$0")"

# Check if flyctl is installed
if ! command -v fly >/dev/null 2>&1 && ! command -v flyctl >/dev/null 2>&1; then
    echo "❌ Error: flyctl CLI is not installed."
    echo "Please install it from https://fly.io/docs/hands-on/install-flyctl/"
    exit 1
fi

# Determine the fly command name
FLY_CMD="fly"
if command -v flyctl >/dev/null 2>&1; then
    FLY_CMD="flyctl"
fi

# Check authentication
echo "🔑 Checking Fly.io authentication status..."
if ! $FLY_CMD auth whoami >/dev/null 2>&1; then
    echo "🔒 Not logged in. Prompting for Fly.io login..."
    $FLY_CMD auth login
fi

# Build TypeScript locally to ensure no compilation errors before deploying
echo "🛠️ Compiling TypeScript to verify build health..."
npm run build

# Deploy
echo "📦 Deploying application via Fly.io..."
$FLY_CMD deploy --remote-only

echo "✅ Deployment complete!"
