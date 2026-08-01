#!/usr/bin/env python3
import os
import sys
import urllib.parse as urlparse
import importlib.util

def load_env():
    env_path = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.env"))
    if os.path.exists(env_path):
        with open(env_path, "r") as f:
            for line in f:
                if line.strip() and not line.startswith("#"):
                    parts = line.strip().split("=", 1)
                    if len(parts) == 2:
                        os.environ[parts[0]] = parts[1].strip(' "')

load_env()

# Try loading from Key Vault using dynamic loading
db_url = None
try:
    if not os.environ.get("AZURE_KEYVAULT_NAME"):
        os.environ["AZURE_KEYVAULT_NAME"] = "dp-kv-deliverypilot"

    kv_file = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../course_src/shared-utils/keyvault.py"))
    if os.path.exists(kv_file):
        spec = importlib.util.spec_from_file_location("keyvault", kv_file)
        keyvault = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(keyvault)
        db_url = keyvault.get_secret("SUPABASE_DB_URL")
except Exception as e:
    print(f"Keyvault fetch notice: {e}", file=sys.stderr)

if not db_url:
    db_url = os.environ.get("SUPABASE_DB_URL")

# If URL is missing or has placeholder, prompt user
if not db_url or "[password]" in db_url:
    print("Database URL not found in Key Vault or environment (or contains '[password]').")
    db_url = input("Please enter your Supabase Database URL (postgresql://...): ").strip()

if not db_url:
    print("Error: No database URL provided.")
    sys.exit(1)

try:
    import pg8000
except ImportError:
    print("Error: pg8000 python library is required. Install it using: pip install pg8000")
    sys.exit(1)

try:
    print("Connecting to Supabase Postgres database...")
    url = urlparse.urlparse(db_url)
    username = url.username
    password = url.password
    database = url.path[1:]
    hostname = url.hostname
    port = url.port or 5432

    conn = pg8000.connect(
        user=username,
        password=password,
        host=hostname,
        port=port,
        database=database
    )
    cursor = conn.cursor()

    migration_path = os.path.abspath(os.path.join(os.path.dirname(__file__), "../migrations/migration_research_assets_ocr_text.sql"))
    print(f"Reading migration from {migration_path}...")
    with open(migration_path, "r") as f:
        sql = f.read()

    print("Executing SQL migration...")
    cursor.execute(sql)
    conn.commit()
    print("Migration executed successfully!")
    cursor.close()
    conn.close()

except Exception as e:
    print(f"Failed to execute migration: {e}", file=sys.stderr)
    sys.exit(1)
