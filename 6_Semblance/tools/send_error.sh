#!/bin/bash
# 🐛 send_error.sh - Auto-send errors to Axiom
# Usage: ./6_Semblance/send_error.sh "stage" "severity" "description of error"

# Find root of repository to locate .env
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Load environment variables
if [ -f "$REPO_ROOT/.env" ]; then
  # Read variables from .env file, ignoring comments and blank lines
  export $(grep -v '^#' "$REPO_ROOT/.env" | xargs)
fi

STAGE=$1
SEVERITY=$2
MESSAGE=$3

if [ -z "$STAGE" ] || [ -z "$SEVERITY" ] || [ -z "$MESSAGE" ]; then
  echo "❌ Usage: $0 <stage> <severity> <message>"
  echo "Example: $0 '5_Symbols' 'ERROR' 'Something went wrong'"
  exit 1
fi

if [ -z "$AXIOM_TOKEN" ] || [ -z "$AXIOM_ORG_ID" ]; then
  echo "❌ Error: AXIOM_TOKEN or AXIOM_ORG_ID not configured in .env"
  exit 1
fi

DATASET=${AXIOM_DATASET:-errors}
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Escape double quotes and backslashes in message for JSON safety
ESCAPED_MESSAGE=$(echo "$MESSAGE" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g')

JSON_PAYLOAD=$(cat <<EOF
[
  {
    "timestamp": "$TIMESTAMP",
    "stage": "$STAGE",
    "severity": "$SEVERITY",
    "message": "$ESCAPED_MESSAGE"
  }
]
EOF
)

API_URL=${AXIOM_API_URL:-https://api.axiom.co}

if [[ "$API_URL" == *".edge.axiom.co"* ]]; then
  INGEST_URL="$API_URL/v1/ingest/$DATASET"
else
  INGEST_URL="$API_URL/v1/datasets/$DATASET/ingest"
fi

echo "📡 Sending error to Axiom (dataset: $DATASET via $INGEST_URL)..."
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$INGEST_URL" \
  -H "Authorization: Bearer $AXIOM_TOKEN" \
  -H "Content-Type: application/json" \
  -H "X-Axiom-Org-Id: $AXIOM_ORG_ID" \
  -d "$JSON_PAYLOAD")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_CODE" -eq 200 ] || [ "$HTTP_CODE" -eq 201 ] || [ "$HTTP_CODE" -eq 204 ]; then
  echo "✅ Error successfully sent to Axiom."
else
  echo "❌ Failed to send error to Axiom. HTTP Code: $HTTP_CODE"
  echo "Response: $BODY"
  exit 1
fi
