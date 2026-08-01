#!/usr/bin/env python3
import os
import sys
import json
import time
import urllib.request
import urllib.parse
import urllib.error

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

supabase_url = os.environ.get("SUPABASE_URL")
service_key = os.environ.get("SUPABASE_SERVICE_KEY")

if not supabase_url or not service_key:
    print("Error: SUPABASE_URL or SUPABASE_SERVICE_KEY not found in .env")
    sys.exit(1)

def http_req(url, method="GET", headers=None, data=None):
    if headers is None:
        headers = {}
    
    req_headers = {
        "apikey": service_key,
        "Authorization": f"Bearer {service_key}"
    }
    req_headers.update(headers)
    
    body = None
    if data is not None:
        body = json.dumps(data).encode("utf-8")
        req_headers["Content-Type"] = "application/json"
        
    req = urllib.request.Request(url, method=method, headers=req_headers, data=body)
    try:
        with urllib.request.urlopen(req) as resp:
            return resp.status, resp.read().decode("utf-8")
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode("utf-8")
    except Exception as e:
        return 0, str(e)

# 1. Fetch existing assets from Supabase
print("Fetching existing research assets from Supabase...")
query_url = f"{supabase_url}/rest/v1/research_assets?container=eq.research-images&select=item_name,ocr_text"
code, body = http_req(query_url)
if code != 200:
    print(f"Failed to fetch assets from Supabase: HTTP {code}\n{body}")
    sys.exit(1)

existing_assets = json.loads(body)
existing_ocr = {item["item_name"]: item["ocr_text"] for item in existing_assets if item.get("ocr_text")}
print(f"Found {len(existing_assets)} registered assets, {len(existing_ocr)} already have OCR text.")

# 2. Fetch list of files from local server (Azure container files)
print("Fetching file list from local server...")
server_url = "http://localhost:8080/api/research/files?container=research-images"
try:
    with urllib.request.urlopen(server_url) as resp:
        files = json.loads(resp.read().decode("utf-8"))
except Exception as e:
    print(f"Failed to connect to local server at {server_url}: {e}")
    print("Please make sure the server is running on port 8080.")
    sys.exit(1)

# Filter files (image content types and not thumbs)
images = []
for f in files:
    name = f.get("name", "")
    content_type = f.get("contentType", "")
    is_img = content_type.startswith("image/") or name.lower().endswith((".png", ".jpg", ".jpeg", ".gif", ".webp"))
    is_thumb = name.startswith("__thumb__")
    if is_img and not is_thumb:
        images.append(f)

print(f"Found {len(images)} images in Azure container.")

# 3. Filter to pending images (no OCR text in DB)
pending = [img for img in images if img["name"] not in existing_ocr]
print(f"{len(pending)} images are pending OCR scanning.")

if not pending:
    print("All images already have OCR text. Nothing to do!")
    sys.exit(0)

# Process pending images sequentially
success_count = 0
for idx, img in enumerate(pending):
    filename = img["name"]
    print(f"\n[{idx+1}/{len(pending)}] Running OCR on: {filename}...")
    
    ocr_api_url = "http://localhost:8080/api/research/ocr"
    req_body = json.dumps({"container": "research-images", "name": filename}).encode("utf-8")
    req = urllib.request.Request(
        ocr_api_url, 
        method="POST", 
        headers={"Content-Type": "application/json"}, 
        data=req_body
    )
    
    try:
        with urllib.request.urlopen(req) as resp:
            res_data = json.loads(resp.read().decode("utf-8"))
            text = res_data.get("text", "").strip()
            
            if not text:
                print("No text detected by AI OCR.")
                continue
                
            print(f"Extracted OCR Text length: {len(text)} characters.")
            
            # Upsert into Supabase
            upsert_url = f"{supabase_url}/rest/v1/research_assets?on_conflict=container,item_name"
            upsert_data = {
                "container": "research-images",
                "item_name": filename,
                "ocr_text": text,
                "updated_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
            }
            # PostgREST upsert requires Prefer header
            up_code, up_body = http_req(
                upsert_url, 
                method="POST", 
                headers={"Prefer": "resolution=merge-duplicates"}, 
                data=upsert_data
            )
            
            if up_code in (200, 201):
                print("Successfully updated in Supabase database!")
                success_count += 1
            else:
                print(f"Error saving to database: HTTP {up_code}\n{up_body}")
                
    except Exception as e:
        print(f"Error processing OCR: {e}")
    
    # Polite spacing to prevent rate limit limits
    time.sleep(2)

print(f"\nBulk OCR Scan finished. Successfully updated {success_count}/{len(pending)} images.")
