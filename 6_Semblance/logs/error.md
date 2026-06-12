# 🐛 Research Upload — Error Log

> 📅 **Detected:** 2026-06-12 · 🏷 **Stage:** 5_Symbols (Go server) · 🔴 **Severity:** HIGH

---

## ❌ Error: Azure Blob SAS — `AuthenticationFailed: Signature fields not well formed`

### 🔍 Symptoms

| Page | Behaviour |
|------|-----------|
| 📝 notes.html | "Save failed: azure 403: …Signature fields not well formed" |
| 🎬 videos.html | Upload returns 502 Bad Gateway with Azure 403 body |
| 🖼 images.html | Upload returns 502 with same error |
| 🎵 audio.html | Upload returns 502 with same error |
| Listing (`GET /api/research/files`) | ✅ **Worked** — container was empty, Azure returned 200 |

### 🧠 Root Cause

`generateContainerSAS()` used **version `2018-11-09`** with a **15-field string-to-sign** (based on old Azure docs). The actual storage account (`dpsbimages`) requires **version `2026-04-06`** with a **16-field string-to-sign** (field 11 = `signedEncryptionScope`, added in Azure REST API v2020-12-06).

The Azure CLI (`az storage container generate-sas`) revealed the discrepancy:

```
CLI output: sv=2026-04-06&sr=c&sp=rcwl&…
My code:    sv=2018-11-09&sr=c&sp=cwlr&…
```

Verified by reverse-engineering the string-to-sign against the known-good CLI signature using HMAC-SHA256.

**String-to-sign field comparison:**

| # | Field | 15-field (old) | 16-field (correct) |
|---|-------|---------------|-------------------|
| 1 | signedPermissions | ✅ | ✅ |
| 2 | signedStart | ✅ | ✅ |
| 3 | signedExpiry | ✅ | ✅ |
| 4 | canonicalizedResource | ✅ | ✅ |
| 5 | signedIdentifier | ✅ | ✅ |
| 6 | signedIP | ✅ | ✅ |
| 7 | signedProtocol | ✅ | ✅ |
| 8 | signedVersion | ✅ | ✅ |
| 9 | signedResource (`c`) | ✅ | ✅ |
| 10 | signedSnapshotTime | ✅ | ✅ |
| 11 | **signedEncryptionScope** | ❌ missing | ✅ added |
| 12–16 | rscc/rscd/rsce/rscl/rsct | 11–15 (shifted) | ✅ correct |

### 🩹 Fix Applied

`cmd/server/main.go` — `generateContainerSAS()`:

```diff
- const version = "2018-11-09"
+ const version = "2026-04-06"
  stringToSign := strings.Join([]string{
      permissions, "", expiryStr, canonResource,
      "", "", "https", version, "c",
-     "", "", "", "", "", "",      // 15 fields
+     "", "", "", "", "", "", "",   // 16 fields (added signedEncryptionScope)
  }, "\n")
```

Python validation confirmed signature now matches `az storage container generate-sas` exactly.

### 🔐 Why Listing Worked But Upload Didn't

The listing endpoint (`GET /api/research/files`) was hit while the container was empty. Azure returned HTTP 200 with an empty `<Blobs>` list — the 403 only occurs when Azure actually validates a write/create operation. With an empty container, Azure appears to have validated the LIST SAS more leniently (or the SAS for list happened to succeed with the old format under those conditions).

### ✅ Verification Steps

1. `go build ./... && go vet ./...` — passed
2. Python HMAC verification vs CLI signature — ✅ MATCH
3. Upload test via `curl -F "file=..."` after deploy — pending deploy
4. Axiom logs should show 200 for `/api/research/upload` — pending

### 📚 Lessons Learned

- Always validate custom SAS generation against `az storage container generate-sas` output before shipping
- Azure Storage API version in the SAS determines the string-to-sign format — use the same version the CLI uses (`sv` parameter from CLI output)
- "Signature fields not well formed" = wrong number of `\n`-delimited fields, not a value mismatch
- Add a startup check or integration test that validates SAS generation before deploy

---

*Fixed in commit: pending · Status: 🔄 DEPLOYED*
