# 📡 Axiom Logging Integration Formula

This formula document describes the integration of Axiom error logging within the **Claude AI Certification for Architects** workspace. It outlines the architectural design, environment settings, and integration patterns used to capture errors from both shell scripts and the web user interface.

---

## 🧠 Regional Edge Architecture

Axiom partitions data storage and ingestion based on the primary region selected during organization creation:
*   **US Region (Default)**: Central endpoint is `https://api.axiom.co/v1/datasets/{dataset}/ingest`.
*   **EU AWS Region (`eu-central-1`)**: Event ingestion uses the Axiom Edge architecture, which routes logs to:
    `https://eu-central-1.aws.edge.axiom.co/v1/ingest/{dataset_name}`

> [!WARNING]
> Sending logs to `https://api.eu.axiom.co` for an EU-partitioned dataset returns a `403 Forbidden` credential error, and sending to `https://api.axiom.co` returns a `400 Bad Request` regional constraint error. The edge hostname `eu-central-1.aws.edge.axiom.co` and endpoint `/v1/ingest/{dataset}` must be used.

---

## ⚙️ Environment Configuration

Add the following parameters to your local `.env` file (ignored by version control) and configure them as secrets in **Azure Key Vault** (`AXIOM-TOKEN` and `AXIOM-ORG-ID`) for production deployment.

```ini
# --- AXIOM LOGGING CONFIGURATION ---
AXIOM_TOKEN=xaat-dcee4b59-5c6e-4ae4-9574-136f5986e84c
AXIOM_ORG_ID=rifaterdemsahin-stks
AXIOM_DATASET=videoproduction
AXIOM_API_URL=https://eu-central-1.aws.edge.axiom.co
```

---

## 💻 Integration Details

### 1. Ingestion via Shell Scripts (`6_Semblance/send_error.sh`)
The shell utility automatically determines the correct ingestion path based on the base domain:
```bash
if [[ "$API_URL" == *".edge.axiom.co"* ]]; then
  # Edge Ingestion Path
  INGEST_URL="$API_URL/v1/ingest/$DATASET"
else
  # Legacy Ingestion Path
  INGEST_URL="$API_URL/v1/datasets/$DATASET/ingest"
fi
```

### 2. Ingestion via Client-Side Web UI (`shared/debug-panel.js`)
The debug panel intercepts JavaScript errors, promise rejections, and fetch failures. It retrieves config variables from `localStorage` (falling back to default code parameters) and posts payloads using the native `_fetch` (original fetch) reference to prevent recursive reporting:
```javascript
const ingestUrl = apiUrl.includes('.edge.axiom.co') 
  ? `${apiUrl}/v1/ingest/${dataset}`
  : `${apiUrl}/v1/datasets/${dataset}/ingest`;

_fetch(ingestUrl, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
    'X-Axiom-Org-Id': orgId
  },
  body: JSON.stringify(payload)
});
```

---

## 🧪 Ingestion Verification

To verify that the ingestion pipeline works correctly:

### CLI Verification:
Run the script manually in the terminal:
```bash
./6_Semblance/send_error.sh "6_Semblance" "INFO" "Verify Axiom logging integration"
```

### UI Verification:
1. Open the Database Seed page at `5_Symbols/supabase/admin.html`.
2. Configure your Axiom parameters in the **Axiom Log Ingestion** card and click **Save Axiom Config**.
3. Click the **💥 Send Test Error** button.
4. Verify that the exception is caught, listed in the debug console overlay at the bottom, and transmitted to Axiom.
