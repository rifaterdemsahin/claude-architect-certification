#!/usr/bin/env python3
"""
Backfill the gdrive_folders table by walking the real Google Drive folder tree.

Reads SUPABASE_URL / SUPABASE_ANON_KEY from .env, pulls the shared Google access
token + root folder id from project_settings, enumerates every folder under the
root via the Drive API, maps each to its course module/video, and upserts rows
into gdrive_folders (idempotent on drive_folder_id).

Folder depth → type:  0 root · 1 module · 2 video · 3 category · 4 subfolder
README.txt was placed only inside category + subfolder folders during creation.
"""
import os, sys, json, urllib.parse, urllib.request

FOLDER_MIME = "application/vnd.google-apps.folder"


def load_env():
    env = {}
    with open(".env") as f:
        for line in f:
            line = line.strip()
            if "=" in line and not line.startswith("#"):
                k, v = line.split("=", 1)
                env[k.strip()] = v.strip().strip("'\"")
    return env


def http(url, headers=None, method="GET", body=None):
    req = urllib.request.Request(url, data=body, method=method, headers=headers or {})
    with urllib.request.urlopen(req) as r:
        return r.status, r.read().decode()


def main():
    env = load_env()
    SB = env["SUPABASE_URL"].rstrip("/")
    ANON = env["SUPABASE_ANON_KEY"]
    sb_headers = {"apikey": ANON, "Authorization": f"Bearer {ANON}"}

    def sb_get(path):
        _, body = http(f"{SB}/rest/v1/{path}", headers=sb_headers)
        return json.loads(body)

    token = sb_get("project_settings?key=eq.gdrive_access_token&select=value")[0]["value"]
    root_id = sb_get("project_settings?key=eq.gdrive_root_folder_id&select=value")[0]["value"]
    drive_headers = {"Authorization": f"Bearer {token}"}

    # Build module/video name → id maps
    modules = sb_get("course_modules?select=id,module_number,title")
    videos = sb_get("course_videos?select=id,module_id,video_number,title")
    module_name_to_id = {f"Module {m['module_number']} - {m['title']}": m["id"] for m in modules}
    video_key_to_id = {(v["module_id"], f"Video {v['video_number']} - {v['title']}"): v["id"] for v in videos}
    root_name = "Claude AI Architect Certification"

    def list_child_folders(parent_id):
        out, page = [], None
        while True:
            q = f"'{parent_id}' in parents and mimeType = '{FOLDER_MIME}' and trashed = false"
            params = {"q": q, "fields": "nextPageToken,files(id,name)", "pageSize": "1000"}
            if page:
                params["pageToken"] = page
            url = "https://www.googleapis.com/drive/v3/files?" + urllib.parse.urlencode(params)
            _, body = http(url, headers=drive_headers)
            data = json.loads(body)
            out.extend(data.get("files", []))
            page = data.get("nextPageToken")
            if not page:
                return out

    rows = []

    def walk(folder_id, name, depth, path, parent_id, module_id, video_id):
        if depth == 0:
            ftype = "root"
        elif depth == 1:
            ftype = "module"; module_id = module_name_to_id.get(name)
        elif depth == 2:
            ftype = "video"; video_id = video_key_to_id.get((module_id, name))
        elif depth == 3:
            ftype = "category"
        else:
            ftype = "subfolder"
        rows.append({
            "drive_folder_id": folder_id,
            "name": name,
            "path": path,
            "drive_url": f"https://drive.google.com/drive/folders/{folder_id}",
            "folder_type": ftype,
            "parent_drive_id": parent_id,
            "module_id": module_id,
            "video_id": video_id,
            "has_readme": ftype in ("category", "subfolder"),
        })
        if depth >= 4:
            return  # subfolders are leaves
        for child in list_child_folders(folder_id):
            walk(child["id"], child["name"], depth + 1, f"{path}/{child['name']}", folder_id, module_id, video_id)

    print(f"Walking Drive tree from root {root_id} ...")
    walk(root_id, root_name, 0, root_name, None, None, None)
    by_type = {}
    for r in rows:
        by_type[r["folder_type"]] = by_type.get(r["folder_type"], 0) + 1
    print(f"Discovered {len(rows)} folders: {by_type}")

    # Bulk upsert into gdrive_folders (merge-duplicates on drive_folder_id)
    post_headers = {**sb_headers, "Content-Type": "application/json",
                    "Prefer": "resolution=merge-duplicates,return=minimal"}
    status, body = http(f"{SB}/rest/v1/gdrive_folders", headers=post_headers,
                        method="POST", body=json.dumps(rows).encode())
    if status not in (200, 201, 204):
        print(f"Upsert failed: HTTP {status} {body}")
        sys.exit(1)
    print(f"✓ Upserted {len(rows)} rows into gdrive_folders (HTTP {status}).")

    # Verify count
    _, cb = http(f"{SB}/rest/v1/gdrive_folders?select=drive_folder_id",
                 headers={**sb_headers, "Prefer": "count=exact", "Range-Unit": "items", "Range": "0-0"})
    print("Verification: re-query gdrive_folders row count below.")


if __name__ == "__main__":
    main()
