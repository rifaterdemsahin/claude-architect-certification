# Zero-Data Retention (ZDR) Verification Checklist

This security checklist helps enterprise architects audit and verify that deployments of Anthropic Claude models via API or cloud partners (AWS Bedrock, GCP Vertex AI) conform to zero-data retention policies.

## 🔐 Compliance Pillars

### 1. API Endpoint Audit
- [ ] Confirm all requests are sent to endpoint variants supporting ZDR.
- [ ] Check console dashboard to ensure "Data Training Opt-out" or similar is checked/contracted.
- [ ] Verify standard payload headers do not bypass security filters.

### 2. Network Boundary Controls
- [ ] Implement VPC endpoints (PrivateLink) to route traffic to LLM APIs without traversing the public internet.
- [ ] Block outgoing TCP connections from worker nodes to arbitrary external endpoints except through authorized proxies.
- [ ] Ensure transport layer security (TLS 1.3 preferred, TLS 1.2 minimum) with enterprise-approved cipher suites.

### 3. Log Redaction and Sanitization
- [ ] Strip Personally Identifiable Information (PII) and Protected Health Information (PHI) before tokenization.
- [ ] Sanitize system logs to prevent leak of prompt/response content in debug streams.
- [ ] Enable hash-based tracking for request correlation instead of logging full text.

### 4. Continuous Auditing
- [ ] Conduct automated daily compliance scans of network exit logs.
- [ ] Run mock threat scenarios simulating man-in-the-middle attacks on the API gateway.
- [ ] Review external vendor SOC2 Type II reports periodically.
