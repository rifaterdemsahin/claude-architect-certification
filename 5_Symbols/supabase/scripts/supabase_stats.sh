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

# All 18 tables defined in schema.sql
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
  printf "  %-25s %10s\n" "TABLE" "ROWS"
  printf "  %-25s %10s\n" "-------------------------" "----------"

  TOTAL_ROWS=0

  for table in "${TABLES[@]}"; do
    if [ "$MODE" = "api" ]; then
      count=$(get_row_count_api "$table")
    else
      count=$(get_row_count_local "$table")
    fi
    case "$count" in
      ''|*[!0-9]*) : ;;
      *) TOTAL_ROWS=$((TOTAL_ROWS + count)) ;;
    esac
    printf "  %-25s %10s\n" "$table" "$count"
  done

  echo ""
  printf "  %-25s %10s\n" "TOTAL ROWS" "$TOTAL_ROWS"
  echo ""
  echo "========================================"
  if [ "$MODE" = "local" ]; then
    echo "  NOTE: Counts reflect seed file INSERT statements only."
    echo "  For live counts set SUPABASE_ANON_KEY or SUPABASE_SERVICE_KEY in .env"
  fi
  echo "  Report saved: $REPORT_FILE"
  echo "========================================"
} | tee "$REPORT_FILE"
