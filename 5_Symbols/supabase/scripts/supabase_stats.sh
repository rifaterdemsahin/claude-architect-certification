#!/usr/bin/env bash
# Supabase database stats — table count and row counts per table
# Run anytime: bash supabase_stats.sh
# Output saved to backup/stats_YYYYMMDD_HHMMSS.txt

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Load .env if present
if [ -f "$SCRIPT_DIR/.env" ]; then
  set -a; source "$SCRIPT_DIR/.env"; set +a
fi

SUPABASE_URL="${SUPABASE_URL:-https://rmekfsdhglyiralxvkwc.supabase.co}"
SUPABASE_KEY="${SUPABASE_SERVICE_KEY:-${SUPABASE_ANON_KEY:-}}"
SCHEMA_FILE="$SCRIPT_DIR/5_Symbols/supabase/schema.sql"
BACKUP_DIR="$SCRIPT_DIR/backup"
TIMESTAMP="$(date +%Y%m%d_%H%M%S)"
REPORT_FILE="$BACKUP_DIR/stats_${TIMESTAMP}.txt"

mkdir -p "$BACKUP_DIR"

# All tables defined in schema.sql
TABLES=(
  modules
  videos
  video_cues
  resource_links
  scenes
  scene_cues
  edl_entries
  checklist_items
  checklist_progress
  course_content
  scripts
  sentences
  outline
  course_modules
  course_videos
  milestones
  milestone_progress
  pricing
  courses
)

TABLE_COUNT="${#TABLES[@]}"

# ──────────────────────────────────────────────
# Helper: query row count via Supabase REST API
# ──────────────────────────────────────────────
get_row_count_api() {
  local table="$1"
  # %header{} requires curl 7.84+; macOS ships 8.x so this is safe
  local output
  output=$(curl -s -o /dev/null \
    -w "%{http_code} %header{content-range}" \
    -H "apikey: $SUPABASE_KEY" \
    -H "Authorization: Bearer $SUPABASE_KEY" \
    -H "Prefer: count=exact" \
    "${SUPABASE_URL}/rest/v1/${table}?select=*&limit=0" 2>/dev/null) || true

  local http_code="${output%% *}"
  if [ "$http_code" = "404" ]; then
    echo "not created"
    return
  fi
  # content-range format: */N or 0-N/N
  local count
  count=$(echo "$output" | grep -oE '[0-9]+$' | tail -1)
  echo "${count:-?}"
}

# ──────────────────────────────────────────────
# Helper: query word count via Supabase REST API
# ──────────────────────────────────────────────
get_word_count_api() {
  local table="$1"
  local column=""
  if [ "$table" = "scripts" ]; then column="script_text"; fi
  if [ "$table" = "sentences" ]; then column="sentence_text"; fi
  if [ "$table" = "scenes" ]; then column="script"; fi
  
  if [ -z "$column" ]; then echo "-"; return; fi

  # Use python for robust word counting to avoid regex issues in bash
  python3 -c "
import requests, os, re
url = \"$SUPABASE_URL\"
key = \"$SUPABASE_KEY\"
headers = {'apikey': key, 'Authorization': f'Bearer {key}'}
res = requests.get(f'{url}/rest/v1/$table?select=$column', headers=headers)
if res.status_code == 200:
    data = res.json()
    total = sum(len(re.findall(r\"\\w+\", item.get(\"$column\", \"\") or \"\")) for item in data)
    print(total)
else:
    print(0)
" 2>/dev/null || echo "0"
}

# ──────────────────────────────────────────────
# Helper: count rows from local SQL seed files
# ──────────────────────────────────────────────
get_row_count_local() {
  local table="$1"
  local seed_dir="$SCRIPT_DIR/5_Symbols/course_src/supabase"
  # Count INSERT INTO <table> lines across all seed files
  grep -rhi "INSERT INTO ${table}" "$seed_dir"/*.sql 2>/dev/null | wc -l | tr -d ' '
}

# ──────────────────────────────────────────────
# Detect mode
# ──────────────────────────────────────────────
MODE="local"
if [ -n "$SUPABASE_KEY" ]; then
  # Quick connectivity check
  HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
    -H "apikey: $SUPABASE_KEY" \
    "${SUPABASE_URL}/rest/v1/modules?select=*&limit=1" 2>/dev/null) || HTTP_CODE=0
  if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "206" ]; then
    MODE="api"
  fi
fi

# ──────────────────────────────────────────────
# Build report
# ──────────────────────────────────────────────
{
  echo "========================================"
  echo "  Supabase Database Stats"
  echo "  Project: rmekfsdhglyiralxvkwc"
  echo "  Generated: $(date '+%Y-%m-%d %H:%M:%S')"
  echo "  Mode: ${MODE} ($([ "$MODE" = 'api' ] && echo 'live REST API' || echo 'local SQL seed files'))"
  echo "========================================"
  echo ""
  echo "  Total tables in schema: $TABLE_COUNT"
  echo ""
  printf "  %-25s %10s %10s\n" "TABLE" "ROWS" "WORDS"
  printf "  %-25s %10s %10s\n" "-------------------------" "----------" "----------"

  TOTAL_ROWS=0
  TOTAL_WORDS=0

  for table in "${TABLES[@]}"; do
    if [ "$MODE" = "api" ]; then
      count=$(get_row_count_api "$table")
      words=$(get_word_count_api "$table")
    else
      count=$(get_row_count_local "$table")
      words="-"
    fi
    
    case "$count" in
      ''|*[!0-9]*) : ;;
      *) TOTAL_ROWS=$((TOTAL_ROWS + count)) ;;
    esac

    case "$words" in
      ''|*[!0-9]*) : ;;
      *) TOTAL_WORDS=$((TOTAL_WORDS + words)) ;;
    esac

    printf "  %-25s %10s %10s\n" "$table" "$count" "$words"
  done

  echo ""
  printf "  %-25s %10s %10s\n" "TOTAL" "$TOTAL_ROWS" "$TOTAL_WORDS"
  echo ""
  echo "========================================"
  if [ "$MODE" = "local" ]; then
    echo "  NOTE: Counts reflect seed file INSERT statements only."
    echo "  For live counts set SUPABASE_ANON_KEY or SUPABASE_SERVICE_KEY in .env"
  fi
  echo "  Report saved: $REPORT_FILE"
  echo "========================================"
} | tee "$REPORT_FILE"
