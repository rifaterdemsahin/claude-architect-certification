#!/usr/bin/env bash
# Supabase database backup script — dumps to local backup/ directory

set -euo pipefail

# Load env vars from .env if present
if [ -f "$(dirname "$0")/.env" ]; then
  # shellcheck disable=SC1090
  set -a; source "$(dirname "$0")/.env"; set +a
fi

SUPABASE_DB_URL="${SUPABASE_DB_URL:-}"
SUPABASE_PROJECT_ID="${SUPABASE_PROJECT_ID:-rmekfsdhglyiralxvkwc}"
BACKUP_DIR="$(dirname "$0")/backup"
TIMESTAMP="$(date +%Y%m%d_%H%M%S)"
BACKUP_FILE="${BACKUP_DIR}/supabase_backup_${TIMESTAMP}.sql"

mkdir -p "$BACKUP_DIR"

# ---- Strategy 1: direct pg_dump via DB URL ----
if [ -n "$SUPABASE_DB_URL" ]; then
  echo "Backing up via pg_dump..."
  pg_dump "$SUPABASE_DB_URL" \
    --no-owner \
    --no-acl \
    --schema=public \
    --file="$BACKUP_FILE"
  echo "Backup saved: $BACKUP_FILE"
  exit 0
fi

# ---- Strategy 2: supabase CLI (linked project) ----
if command -v supabase &>/dev/null; then
  echo "Backing up via Supabase CLI (linked)..."
  if supabase db dump --linked --file "$BACKUP_FILE" 2>/dev/null; then
    echo "Backup saved: $BACKUP_FILE"
    exit 0
  fi
  echo "Supabase CLI linked dump failed — falling through to snapshot..."
fi

# ---- Strategy 3: schema + seed files snapshot ----
echo "Neither SUPABASE_DB_URL nor supabase CLI found — snapshotting local SQL files..."
SNAPSHOT="${BACKUP_DIR}/supabase_snapshot_${TIMESTAMP}.sql"
cat \
  5_Symbols/src/supabase/schema.sql \
  5_Symbols/src/supabase/seed.sql 2>/dev/null \
  > "$SNAPSHOT" || true
echo "Snapshot saved: $SNAPSHOT"
echo ""
echo "To enable live backup, set SUPABASE_DB_URL in .env:"
echo "  SUPABASE_DB_URL=postgresql://postgres:[password]@db.${SUPABASE_PROJECT_ID}.supabase.co:5432/postgres"
