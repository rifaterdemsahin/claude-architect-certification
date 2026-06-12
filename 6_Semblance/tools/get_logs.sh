#!/bin/bash
# 📡 get_logs.sh - Pull logs from Axiom using APL
# Usage: ./6_Semblance/get_logs.sh [limit]

# Find root of repository to locate .env
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Load environment variables
if [ -f "$REPO_ROOT/.env" ]; then
  export $(grep -v '^#' "$REPO_ROOT/.env" | xargs)
fi

LIMIT=${1:-20}

if [ -z "$AXIOM_TOKEN" ] || [ -z "$AXIOM_ORG_ID" ]; then
  echo "❌ Error: AXIOM_TOKEN or AXIOM_ORG_ID not configured in .env"
  exit 1
fi

DATASET=${AXIOM_DATASET:-videoproduction}
API_URL=${AXIOM_QUERY_URL:-${AXIOM_API_URL:-https://api.axiom.co}}

# Construct APL query to select latest entries
APL_QUERY="['$DATASET'] | order by timestamp desc | limit $LIMIT"

JSON_PAYLOAD=$(cat <<EOF
{
  "apl": "$APL_QUERY",
  "startTime": "now-24h"
}
EOF
)

echo "🔍 Fetching last $LIMIT events from Axiom (dataset: $DATASET via $API_URL)..."

# Axiom APL query endpoint requires format parameter (legacy or tabular)
RESPONSE_BODY=$(curl -s -X POST "$API_URL/v1/query/_apl?format=legacy" \
  -H "Authorization: Bearer $AXIOM_TOKEN" \
  -H "Content-Type: application/json" \
  -H "X-Axiom-Org-Id: $AXIOM_ORG_ID" \
  -d "$JSON_PAYLOAD")

# Check if response indicates error
if echo "$RESPONSE_BODY" | grep -q '"message"'; then
  echo "❌ Failed to query logs from Axiom."
  echo "Response: $RESPONSE_BODY"
  echo "Tip: Make sure the API Token has the 'Query' permission role enabled in Axiom Settings."
  exit 1
fi

# Print formatted output using python json.tool
echo "$RESPONSE_BODY" | python3 -m json.tool
