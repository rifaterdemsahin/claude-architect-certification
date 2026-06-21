# 🔐 Formula: Browser-to-CLI Google Drive Token Sharing

> **Stage:** `4_Formula` — Thinking & Planning  
> **Applies to:** [`5_Symbols/production/prod/google_drive_folder_creator.html`](../../5_Symbols/production/prod/google_drive_folder_creator.html) · [`7_Testing_Known/create_gdrive_folders.py`](../../7_Testing_Known/create_gdrive_folders.py)  
> **Related issues:** Port conflicts (8765 occupied) · GCP Web Client `redirect_uri_mismatch` on CLI loopback ports.

---

## 🧠 The Problem

Automated scripts running in the terminal need to perform write operations on Google Drive (creating folder structures for modules, videos, and subfolders like `raw`, `export`, `research`). However, standard terminal-based OAuth 2.0 flows face major roadblocks:

1. **Redirect URI Mismatch**: The OAuth Client ID configured in Google Cloud is a **Web Application Client ID**. Google Web Clients strictly validate redirect URIs. When terminal scripts start local loopback servers (e.g. `http://localhost:<random_port>/`), it results in a 400 Bad Request because the dynamic loopback port is not registered.
2. **Local Port Conflicts**: Attempting to use a static registered port like `8765` fails on development machines if the port is already bound to other local daemons or background python processes.
3. **Headless Limitations**: Terminal-based login triggers browser popups, which is not suitable for purely command-driven environments.

---

## 🏛️ Shared-Token Architecture

To solve these constraints, we implement a **hybrid token-sharing pattern** using the live database as a secure, transient handshake layer:

```
Rifat's Local Chrome (Authorized Session)
  │
  ├─ 1. Click "Connect Google Account" (Implicit Flow)
  ├─ 2. GIS Popup signs in via existing Google cookies (No redirects needed)
  ├─ 3. Upsert access_token to Supabase: `project_settings.gdrive_access_token`
  │
  ▼
Supabase Database (project_settings)
  ▲
  │
  ├─ 4. CLI Script polls Supabase every 3 seconds for `gdrive_access_token`
  ├─ 5. CLI Script extracts token and instantiates Google Drive API Client
  ├─ 6. CLI Script performs real API folder creation calls on Google Drive
  └─ 7. CLI Script writes real Google folder links to Supabase course tables
```

---

## 💻 Technical Implementation Details

### 1. Browser-Side Token Sharing (JavaScript)

In [`google_drive_folder_creator.html`](../../5_Symbols/production/prod/google_drive_folder_creator.html), when Google Identity Services returns the authorized token client response, the access token is written directly to the `project_settings` table:

```javascript
async function saveTokenToSupabase(token) {
  if (db) {
    try {
      await db.from('project_settings')
        .upsert({ key: 'gdrive_access_token', value: token });
      log('success', '✓ Shared access token with Supabase for terminal CLI tools.');
    } catch (err) {
      log('warn', 'Could not save access token to Supabase: ' + err.message);
    }
  }
}
```

### 2. CLI-Side Token Retrieval & Polling (Python)

In [`create_gdrive_folders.py`](../../7_Testing_Known/create_gdrive_folders.py), the script queries Supabase. If the token is missing or invalid, it prints steps instructing the user to sign in, then enters a non-blocking check loop:

```python
def fetch_token_from_supabase():
    endpoint = f"{SUPABASE_URL}/rest/v1/project_settings?key=eq.gdrive_access_token&select=value"
    headers = {
        'apikey': SUPABASE_ANON_KEY,
        'Authorization': f'Bearer {SUPABASE_ANON_KEY}'
    }
    try:
        res = requests.get(endpoint, headers=headers)
        if res.status_code == 200:
            data = res.json()
            if data and len(data) > 0:
                return data[0].get('value')
    except Exception as e:
        log_warn(f"Failed to fetch shared token from Supabase: {e}")
    return None

def authenticate_google_drive(poll_only=False):
    # Retrieve the shared token
    supabase_token = fetch_token_from_supabase()
    if supabase_token:
        try:
            return Credentials(token=supabase_token)
        except Exception as e:
            log_warn(f"Invalid token: {e}")
    
    # Fallback to local token.json
    ...
```

---

## 🧪 Quick Verification Checklist

- [x] Connect Google account on [Google Drive Folder Creator Tool](http://localhost:8080/5_Symbols/production/prod/google_drive_folder_creator.html).
- [x] Confirm `Shared access token with Supabase for terminal CLI tools` log is visible.
- [x] Run `python3 7_Testing_Known/create_gdrive_folders.py` in the terminal.
- [x] Verify script completes with `PRODUCTION MODE` and creates folders on Google Drive.
- [x] Verify Supabase tables contain real Drive links in the format `https://drive.google.com/drive/folders/<id>`.
