# 🔍 Axiom Query Guide — `videoproduction` Dataset

How to query, filter, and analyze logs from the `videoproduction` dataset in Axiom.

---

## 📊 Dataset Schema

Every ingested log entry has these fields:

| Field | Type | Example | Source |
|-------|------|---------|--------|
| `timestamp` | `timestamp` | `2026-06-08T14:30:00Z` | Auto-set at ingestion |
| `stage` | `string` | `5_Symbols` | Passed by `send_error.sh` (stage folder) |
| `severity` | `string` | `ERROR`, `WARN`, `INFO` | Passed by caller |
| `message` | `string` | `Failed to connect to MCP server` | Error description |

---

## 🖥️ Web UI Query

1. Open the [Axiom Query Editor](https://app.axiom.co/rifaterdemsahin-stks/query)
2. Select the **`videoproduction`** dataset from the dropdown
3. Write APL queries (see below)

### Saved Query (Direct Link)

[📌 Production Error Logs](https://app.axiom.co/rifaterdemsahin-stks/query?qid=VmhIX12pi7S-tgc7x0&relative=1)

This saved query shows recent errors grouped by stage and severity.

---

## 📝 APL Query Examples

### 1. All errors in the last hour

```apl
['videoproduction']
| where severity == 'ERROR'
| where @timestamp > now() - 1h
| order by @timestamp desc
```

### 2. Errors by stage (count)

```apl
['videoproduction']
| where severity == 'ERROR'
| summarize count() by stage
| order by count_ desc
```

### 3. Warnings in Semblance stage (last 24h)

```apl
['videoproduction']
| where stage == '6_Semblance'
| where severity in ('WARN', 'ERROR')
| where @timestamp > now() - 24h
| order by @timestamp desc
```

### 4. All logs grouped by severity, last 7 days

```apl
['videoproduction']
| where @timestamp > now() - 7d
| summarize count() by severity
| order by count_ desc
```

### 5. Search messages containing a keyword

```apl
['videoproduction']
| where message contains 'Supabase'
| where @timestamp > now() - 24h
| order by @timestamp desc
```

### 6. Timeline: errors per hour (last 24h)

```apl
['videoproduction']
| where severity == 'ERROR'
| where @timestamp > now() - 24h
| summarize count() by bin(@timestamp, 1h)
| order by bin_ asc
```

---

## 🧪 Axiom CLI (Command Line)

### Install

```bash
# macOS
brew install axiom/tap/axiom
```

### Authenticate

Create a personal token in Axiom → Settings → API Tokens, then:

```bash
export AXIOM_TOKEN="xaat-xxxxxxxxx"
export AXIOM_ORG_ID="rifaterdemsahin-stks"
```

### Query

```bash
# Run an APL query from the terminal
axiom dataset query videoproduction \
  '["videoproduction"] | where severity == "ERROR" | order by @timestamp desc | limit 20'
```

### Live Tail

```bash
# Watch incoming logs in real-time
axiom dataset tail videoproduction
```

---

## 🌐 REST API Query

```bash
curl -s -X POST "https://api.axiom.co/v1/datasets/videoproduction/query" \
  -H "Authorization: Bearer $AXIOM_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": {
      "startTime": "now-1h",
      "endTime": "now",
      "aggregations": [
        {"alias": "count", "op": "count"}
      ],
      "groupBy": ["severity"]
    }
  }'
```

### EU Edge Region

If using the EU edge region for queries:

```bash
curl -s -X POST "https://eu-central-1.aws.edge.axiom.co/v1/datasets/videoproduction/query" \
  -H "Authorization: Bearer $AXIOM_TOKEN" \
  -H "X-Axiom-Org-Id: rifaterdemsahin-stks" \
  -H "Content-Type: application/json" \
  -d '{"query": {"startTime": "now-24h", "endTime": "now"}}'
```

---

## 🏷 Time Ranges (Web UI & APL)

| Shortcut | Meaning |
|----------|---------|
| `now() - 5m` | Last 5 minutes |
| `now() - 1h` | Last hour |
| `now() - 24h` | Last 24 hours |
| `now() - 7d` | Last 7 days |
| `now() - 30d` | Last 30 days |

In the web UI, use the **Relative** time picker instead of typing these manually.

---

## 🔗 Quick Links

| Resource | URL |
|----------|-----|
| Axiom Query Editor | https://app.axiom.co/rifaterdemsahin-stks/query |
| Saved Error Query | https://app.axiom.co/rifaterdemsahin-stks/query?qid=VmhIX12pi7S-tgc7x0&relative=1 |
| Axiom Settings | https://app.axiom.co/rifaterdemsahin-stks/settings |
| Ingest Script | `./6_Semblance/send_error.sh` |
| Setup Guide | `./axiom_logging_setup.md` |

---

## ⚙️ Ingestion Quick Reference

```bash
# Send an error
./6_Semblance/send_error.sh "5_Symbols" "ERROR" "Something went wrong"

# Send a warning
./6_Semblance/send_error.sh "6_Semblance" "WARN" "Rate limit approaching"

# Send an info event
./6_Semblance/send_error.sh "7_Testing_Known" "INFO" "Validation run completed"
```